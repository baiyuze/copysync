// Package store 是 daemon 的本地持久化层：剪贴板历史、已配对设备、杂项元数据。
//
// 用 modernc.org/sqlite（纯 Go 实现）而非 mattn/go-sqlite3——剪贴板层已经
// 不得不引入 cgo，数据库这层能省则省，也让交叉编译简单些。
package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Kind int32

const (
	KindUnspecified Kind = 0
	KindText        Kind = 1
	KindHTML        Kind = 2
	KindImage       Kind = 3
	KindFile        Kind = 4
)

type Status int32

const (
	StatusUnspecified Status = 0
	StatusReady       Status = 1 // 内容已在本地，可直接粘贴
	StatusRemoteOnly  Status = 2 // 超过阈值，只有元数据
	StatusFetching    Status = 3
	StatusFailed      Status = 4
	StatusExpired     Status = 5
)

type Item struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	IsDir       bool   `json:"is_dir"`
	ContentType string `json:"content_type"`
}

type Clip struct {
	ID               string
	Kind             Kind
	Status           Status
	OriginDeviceID   string
	OriginDeviceName string
	Outgoing         bool // true=本机复制  false=从对端收到
	Items            []Item
	TotalSize        int64
	TextPreview      string
	TextContent      string   // 文本/HTML 的完整内容，直接存库
	TextPlain        string   // HTML 的纯文本形式，供只认纯文本的输入框粘贴
	CachePath        string   // 相对 cache 目录的落盘位置
	SourcePaths      []string // 本机复制时的源路径，供对端 fetch 时回源
	Error            string
	CreatedAt        time.Time
	ExpiresAt        time.Time
}

type Device struct {
	ID        string
	Name      string
	Platform  string
	PublicKey []byte
	PairedAt  time.Time
}

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	// busy_timeout 避免 GC 与写入并发时立刻报 SQLITE_BUSY；
	// WAL 让读写互不阻塞；foreign_keys 让 clip_items 的级联删除生效。
	dsn := path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库: %w", err)
	}
	// modernc 的驱动在并发写时容易踩锁，单连接最省心，且本地负载极低
	db.SetMaxOpenConns(1)

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	const schema = `
CREATE TABLE IF NOT EXISTS clips (
    id                 TEXT PRIMARY KEY,
    kind               INTEGER NOT NULL,
    status             INTEGER NOT NULL,
    origin_device_id   TEXT NOT NULL DEFAULT '',
    origin_device_name TEXT NOT NULL DEFAULT '',
    outgoing           INTEGER NOT NULL DEFAULT 0,
    total_size         INTEGER NOT NULL DEFAULT 0,
    text_preview       TEXT NOT NULL DEFAULT '',
    text_content       TEXT NOT NULL DEFAULT '',
    text_plain         TEXT NOT NULL DEFAULT '',
    cache_path         TEXT NOT NULL DEFAULT '',
    source_paths       TEXT NOT NULL DEFAULT '[]',
    error              TEXT NOT NULL DEFAULT '',
    created_at         INTEGER NOT NULL,
    expires_at         INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_clips_created  ON clips(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_clips_expires  ON clips(expires_at);

CREATE TABLE IF NOT EXISTS clip_items (
    clip_id      TEXT NOT NULL REFERENCES clips(id) ON DELETE CASCADE,
    idx          INTEGER NOT NULL,
    name         TEXT NOT NULL,
    size         INTEGER NOT NULL DEFAULT 0,
    is_dir       INTEGER NOT NULL DEFAULT 0,
    content_type TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (clip_id, idx)
);

CREATE TABLE IF NOT EXISTS devices (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL DEFAULT '',
    platform   TEXT NOT NULL DEFAULT '',
    public_key BLOB NOT NULL,
    paired_at  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS meta (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);`
	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("建表: %w", err)
	}
	// 建表语句对已有的库不生效，后加的列要单独补上
	return s.addColumnIfMissing("clips", "text_plain", "TEXT NOT NULL DEFAULT ''")
}

func (s *Store) addColumnIfMissing(table, column, decl string) error {
	rows, err := s.db.Query(`SELECT name FROM pragma_table_info(?)`, table)
	if err != nil {
		return fmt.Errorf("读取 %s 表结构: %w", table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		if name == column {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if _, err := s.db.Exec(`ALTER TABLE ` + table + ` ADD COLUMN ` + column + ` ` + decl); err != nil {
		return fmt.Errorf("为 %s 添加列 %s: %w", table, column, err)
	}
	return nil
}

// ─────────────────────────── clips ───────────────────────────

func (s *Store) PutClip(c Clip) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	srcJSON, err := json.Marshal(c.SourcePaths)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO clips (id, kind, status, origin_device_id, origin_device_name, outgoing,
                   total_size, text_preview, text_content, text_plain, cache_path,
                   source_paths, error, created_at, expires_at)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
ON CONFLICT(id) DO UPDATE SET
    status=excluded.status, total_size=excluded.total_size,
    text_preview=excluded.text_preview, text_content=excluded.text_content,
    text_plain=excluded.text_plain,
    cache_path=excluded.cache_path, source_paths=excluded.source_paths,
    error=excluded.error, expires_at=excluded.expires_at`,
		c.ID, c.Kind, c.Status, c.OriginDeviceID, c.OriginDeviceName, c.Outgoing,
		c.TotalSize, c.TextPreview, c.TextContent, c.TextPlain, c.CachePath,
		string(srcJSON), c.Error, c.CreatedAt.Unix(), c.ExpiresAt.Unix())
	if err != nil {
		return fmt.Errorf("写入 clip: %w", err)
	}

	// items 整体替换，避免增量更新的边界情况
	if _, err := tx.Exec(`DELETE FROM clip_items WHERE clip_id=?`, c.ID); err != nil {
		return err
	}
	for i, it := range c.Items {
		if _, err := tx.Exec(`
INSERT INTO clip_items (clip_id, idx, name, size, is_dir, content_type)
VALUES (?,?,?,?,?,?)`, c.ID, i, it.Name, it.Size, it.IsDir, it.ContentType); err != nil {
			return fmt.Errorf("写入 clip_item: %w", err)
		}
	}
	return tx.Commit()
}

var ErrNotFound = errors.New("记录不存在")

func (s *Store) GetClip(id string) (Clip, error) {
	rows, err := s.queryClips(`SELECT `+clipCols+` FROM clips WHERE id=?`, id)
	if err != nil {
		return Clip{}, err
	}
	if len(rows) == 0 {
		return Clip{}, ErrNotFound
	}
	return rows[0], nil
}

// ListClips 按时间倒序分页。before 为零值表示从最新开始。
func (s *Store) ListClips(limit int, before time.Time, kind Kind) ([]Clip, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	q := `SELECT ` + clipCols + ` FROM clips WHERE 1=1`
	var args []any
	if !before.IsZero() {
		q += ` AND created_at < ?`
		args = append(args, before.Unix())
	}
	if kind != KindUnspecified {
		q += ` AND kind = ?`
		args = append(args, kind)
	}
	q += ` ORDER BY created_at DESC LIMIT ?`
	args = append(args, limit)
	return s.queryClips(q, args...)
}

const clipCols = `id, kind, status, origin_device_id, origin_device_name, outgoing,
                  total_size, text_preview, text_content, text_plain, cache_path,
                  source_paths, error, created_at, expires_at`

func (s *Store) queryClips(query string, args ...any) ([]Clip, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询 clips: %w", err)
	}
	defer rows.Close()

	var out []Clip
	for rows.Next() {
		var c Clip
		var srcJSON string
		var created, expires int64
		if err := rows.Scan(&c.ID, &c.Kind, &c.Status, &c.OriginDeviceID, &c.OriginDeviceName,
			&c.Outgoing, &c.TotalSize, &c.TextPreview, &c.TextContent, &c.TextPlain,
			&c.CachePath, &srcJSON, &c.Error, &created, &expires); err != nil {
			return nil, err
		}
		c.CreatedAt = time.Unix(created, 0)
		c.ExpiresAt = time.Unix(expires, 0)
		_ = json.Unmarshal([]byte(srcJSON), &c.SourcePaths)
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// items 单独查，避免 JOIN 后的行展开处理
	for i := range out {
		if out[i].Items, err = s.itemsOf(out[i].ID); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (s *Store) itemsOf(clipID string) ([]Item, error) {
	rows, err := s.db.Query(`
SELECT name, size, is_dir, content_type FROM clip_items
WHERE clip_id=? ORDER BY idx`, clipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.Name, &it.Size, &it.IsDir, &it.ContentType); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (s *Store) SetClipStatus(id string, st Status, errMsg string) error {
	_, err := s.db.Exec(`UPDATE clips SET status=?, error=? WHERE id=?`, st, errMsg, id)
	return err
}

func (s *Store) SetClipCache(id, cachePath string, totalSize int64) error {
	_, err := s.db.Exec(`
UPDATE clips SET cache_path=?, total_size=?, status=?, error='' WHERE id=?`,
		cachePath, totalSize, StatusReady, id)
	return err
}

func (s *Store) DeleteClips(ids []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, id := range ids {
		if _, err := tx.Exec(`DELETE FROM clips WHERE id=?`, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) DeleteAllClips() error {
	_, err := s.db.Exec(`DELETE FROM clips`)
	return err
}

// ExpiredClips 返回已过期的记录，交由 GC 先删缓存文件再删行。
func (s *Store) ExpiredClips(now time.Time) ([]Clip, error) {
	return s.queryClips(`SELECT `+clipCols+` FROM clips WHERE expires_at <= ?`, now.Unix())
}

// CachePaths 返回所有仍被引用的缓存路径，供 GC 与磁盘对账、清理孤儿文件。
func (s *Store) CachePaths() (map[string]struct{}, error) {
	rows, err := s.db.Query(`SELECT cache_path FROM clips WHERE cache_path <> ''`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]struct{}{}
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		out[p] = struct{}{}
	}
	return out, rows.Err()
}

// ─────────────────────────── devices ───────────────────────────

func (s *Store) PutDevice(d Device) error {
	_, err := s.db.Exec(`
INSERT INTO devices (id, name, platform, public_key, paired_at) VALUES (?,?,?,?,?)
ON CONFLICT(id) DO UPDATE SET
    name=excluded.name, platform=excluded.platform, public_key=excluded.public_key`,
		d.ID, d.Name, d.Platform, d.PublicKey, d.PairedAt.Unix())
	return err
}

func (s *Store) ListDevices() ([]Device, error) {
	rows, err := s.db.Query(`
SELECT id, name, platform, public_key, paired_at FROM devices ORDER BY paired_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Device
	for rows.Next() {
		var d Device
		var paired int64
		if err := rows.Scan(&d.ID, &d.Name, &d.Platform, &d.PublicKey, &paired); err != nil {
			return nil, err
		}
		d.PairedAt = time.Unix(paired, 0)
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *Store) DeleteDevice(id string) error {
	_, err := s.db.Exec(`DELETE FROM devices WHERE id=?`, id)
	return err
}

// ─────────────────────────── meta ───────────────────────────

func (s *Store) GetMeta(key string) (string, error) {
	var v string
	err := s.db.QueryRow(`SELECT value FROM meta WHERE key=?`, key).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return v, err
}

func (s *Store) SetMeta(key, value string) error {
	_, err := s.db.Exec(`
INSERT INTO meta (key, value) VALUES (?,?)
ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, value)
	return err
}
