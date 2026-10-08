package store

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"
)

func openTest(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("打开数据库: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func sampleClip(id string, created time.Time) Clip {
	return Clip{
		ID:               id,
		Kind:             KindFile,
		Status:           StatusReady,
		OriginDeviceID:   "dev-1",
		OriginDeviceName: "笔记本",
		TotalSize:        2048,
		CachePath:        id + ".bin",
		SourcePaths:      []string{"/tmp/a.txt", "/tmp/b.txt"},
		Items: []Item{
			{Name: "a.txt", Size: 1024, ContentType: "public.plain-text"},
			{Name: "docs", Size: 1024, IsDir: true},
		},
		CreatedAt: created,
		ExpiresAt: created.Add(72 * time.Hour),
	}
}

func TestPutAndGetRoundTrip(t *testing.T) {
	s := openTest(t)
	now := time.Now().Truncate(time.Second)
	want := sampleClip("clip-1", now)

	if err := s.PutClip(want); err != nil {
		t.Fatalf("写入: %v", err)
	}
	got, err := s.GetClip("clip-1")
	if err != nil {
		t.Fatalf("读取: %v", err)
	}

	if got.TotalSize != want.TotalSize || got.OriginDeviceName != want.OriginDeviceName {
		t.Errorf("基本字段不一致: %+v", got)
	}
	if len(got.Items) != 2 || got.Items[0].Name != "a.txt" || !got.Items[1].IsDir {
		t.Errorf("items 未正确往返: %+v", got.Items)
	}
	// 源路径供对端 fetch 时回源，丢了就再也拉不到原文件
	if len(got.SourcePaths) != 2 || got.SourcePaths[1] != "/tmp/b.txt" {
		t.Errorf("source_paths 未正确往返: %v", got.SourcePaths)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("创建时间 = %v，期望 %v", got.CreatedAt, now)
	}
}

func TestGetMissingClip(t *testing.T) {
	s := openTest(t)
	if _, err := s.GetClip("nope"); err != ErrNotFound {
		t.Errorf("err = %v，期望 ErrNotFound", err)
	}
}

// 同一条记录重复写入应当是更新而非报错——收到进度更新时会反复写。
func TestPutClipIsUpsert(t *testing.T) {
	s := openTest(t)
	now := time.Now()
	c := sampleClip("clip-1", now)
	if err := s.PutClip(c); err != nil {
		t.Fatal(err)
	}

	c.Status = StatusFetching
	c.TotalSize = 4096
	c.Items = []Item{{Name: "only.txt", Size: 4096}}
	if err := s.PutClip(c); err != nil {
		t.Fatalf("二次写入: %v", err)
	}

	got, _ := s.GetClip("clip-1")
	if got.Status != StatusFetching || got.TotalSize != 4096 {
		t.Errorf("更新未生效: status=%v size=%d", got.Status, got.TotalSize)
	}
	// items 是整体替换，不能残留旧条目
	if len(got.Items) != 1 || got.Items[0].Name != "only.txt" {
		t.Errorf("items 未被整体替换: %+v", got.Items)
	}
}

func TestListClipsOrderAndPaging(t *testing.T) {
	s := openTest(t)
	base := time.Now().Truncate(time.Second)
	for i := range 5 {
		c := sampleClip("clip-"+string(rune('a'+i)), base.Add(time.Duration(i)*time.Minute))
		if err := s.PutClip(c); err != nil {
			t.Fatal(err)
		}
	}

	got, err := s.ListClips(10, time.Time{}, KindUnspecified)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 {
		t.Fatalf("返回 %d 条，期望 5", len(got))
	}
	// 倒序：最新的在最前，界面据此直接渲染
	for i := 1; i < len(got); i++ {
		if got[i].CreatedAt.After(got[i-1].CreatedAt) {
			t.Errorf("第 %d 条比前一条更新，排序不是倒序", i)
		}
	}

	// 游标分页：只返回早于 before 的记录
	page, err := s.ListClips(10, got[2].CreatedAt, KindUnspecified)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range page {
		if !c.CreatedAt.Before(got[2].CreatedAt) {
			t.Errorf("分页返回了不早于游标的记录: %v", c.ID)
		}
	}
}

func TestListClipsKindFilter(t *testing.T) {
	s := openTest(t)
	now := time.Now()
	fileClip := sampleClip("f", now)
	textClip := sampleClip("t", now.Add(time.Minute))
	textClip.Kind = KindText
	for _, c := range []Clip{fileClip, textClip} {
		if err := s.PutClip(c); err != nil {
			t.Fatal(err)
		}
	}

	got, err := s.ListClips(10, time.Time{}, KindText)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "t" {
		t.Errorf("类型过滤失效: %+v", got)
	}
}

// GC 依赖这个查询找出该删的记录；边界判定错了会导致缓存永不释放。
func TestExpiredClips(t *testing.T) {
	s := openTest(t)
	now := time.Now().Truncate(time.Second)

	fresh := sampleClip("fresh", now)
	fresh.ExpiresAt = now.Add(time.Hour)
	stale := sampleClip("stale", now.Add(-48*time.Hour))
	stale.ExpiresAt = now.Add(-time.Hour)
	for _, c := range []Clip{fresh, stale} {
		if err := s.PutClip(c); err != nil {
			t.Fatal(err)
		}
	}

	got, err := s.ExpiredClips(now)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "stale" {
		t.Fatalf("过期查询结果 = %+v，期望只有 stale", got)
	}
	// 过期记录要带着 cache_path 返回，GC 才知道该删哪个文件
	if got[0].CachePath == "" {
		t.Error("过期记录缺少 cache_path，GC 无法清理对应文件")
	}
}

// GC 要和磁盘对账，清掉没有记录引用的孤儿文件。
func TestCachePaths(t *testing.T) {
	s := openTest(t)
	now := time.Now()
	withCache := sampleClip("a", now)
	noCache := sampleClip("b", now)
	noCache.CachePath = ""
	for _, c := range []Clip{withCache, noCache} {
		if err := s.PutClip(c); err != nil {
			t.Fatal(err)
		}
	}

	paths, err := s.CachePaths()
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("返回 %d 条，期望 1（空 cache_path 应被排除）", len(paths))
	}
	if _, ok := paths["a.bin"]; !ok {
		t.Errorf("缺少 a.bin: %v", paths)
	}
}

// 删除记录时 clip_items 必须级联删除，否则会留下孤儿行。
func TestDeleteCascadesItems(t *testing.T) {
	s := openTest(t)
	if err := s.PutClip(sampleClip("clip-1", time.Now())); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteClips([]string{"clip-1"}); err != nil {
		t.Fatal(err)
	}

	items, err := s.itemsOf("clip-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Errorf("删除记录后仍残留 %d 条 item", len(items))
	}
}

func TestSetClipCacheMarksReady(t *testing.T) {
	s := openTest(t)
	c := sampleClip("clip-1", time.Now())
	c.Status = StatusRemoteOnly
	c.CachePath = ""
	c.Error = "上次失败了"
	if err := s.PutClip(c); err != nil {
		t.Fatal(err)
	}

	if err := s.SetClipCache("clip-1", "downloaded.bin", 9999); err != nil {
		t.Fatal(err)
	}
	got, _ := s.GetClip("clip-1")
	if got.Status != StatusReady || got.CachePath != "downloaded.bin" || got.TotalSize != 9999 {
		t.Errorf("落盘后状态不对: %+v", got)
	}
	// 上一次的错误信息必须清掉，否则界面会一直显示旧的失败原因
	if got.Error != "" {
		t.Errorf("落盘成功后仍残留错误信息: %q", got.Error)
	}
}

func TestDeviceRoundTrip(t *testing.T) {
	s := openTest(t)
	d := Device{
		ID:        "dev-1",
		Name:      "台式机",
		Platform:  "windows",
		PublicKey: []byte{1, 2, 3, 4},
		PairedAt:  time.Now().Truncate(time.Second),
	}
	if err := s.PutDevice(d); err != nil {
		t.Fatal(err)
	}

	list, err := s.ListDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "台式机" || string(list[0].PublicKey) != "\x01\x02\x03\x04" {
		t.Fatalf("设备未正确往返: %+v", list)
	}

	// 重连时设备名可能变了，应当更新而非插入重复行
	d.Name = "改名后的台式机"
	if err := s.PutDevice(d); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListDevices()
	if len(list) != 1 || list[0].Name != "改名后的台式机" {
		t.Errorf("重复写入未更新: %+v", list)
	}

	if err := s.DeleteDevice("dev-1"); err != nil {
		t.Fatal(err)
	}
	if list, _ = s.ListDevices(); len(list) != 0 {
		t.Errorf("删除后仍有 %d 台设备", len(list))
	}
}

func TestMetaRoundTrip(t *testing.T) {
	s := openTest(t)
	if _, err := s.GetMeta("missing"); err != ErrNotFound {
		t.Errorf("err = %v，期望 ErrNotFound", err)
	}
	if err := s.SetMeta("k", "v1"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetMeta("k", "v2"); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetMeta("k")
	if err != nil || got != "v2" {
		t.Errorf("got = %q, err = %v，期望 v2", got, err)
	}
}

func TestHTMLPlainTextRoundTrip(t *testing.T) {
	s := openTest(t)
	now := time.Now().Truncate(time.Second)
	c := Clip{
		ID: "html-1", Kind: KindHTML, Status: StatusReady,
		TextContent: "<meta charset='utf-8'><b>你好</b>", TextPlain: "你好",
		CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	}
	if err := s.PutClip(c); err != nil {
		t.Fatalf("写入: %v", err)
	}
	got, err := s.GetClip("html-1")
	if err != nil {
		t.Fatalf("读取: %v", err)
	}
	if got.TextContent != c.TextContent || got.TextPlain != "你好" {
		t.Errorf("HTML 与纯文本未正确往返: content=%q plain=%q", got.TextContent, got.TextPlain)
	}
}

// 升级前的库没有 text_plain 列，打开时要自动补上，且旧记录照常可读。
func TestOpenMigratesOldSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "old.db")
	old, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := old.Exec(`
CREATE TABLE clips (
    id TEXT PRIMARY KEY, kind INTEGER NOT NULL, status INTEGER NOT NULL,
    origin_device_id TEXT NOT NULL DEFAULT '', origin_device_name TEXT NOT NULL DEFAULT '',
    outgoing INTEGER NOT NULL DEFAULT 0, total_size INTEGER NOT NULL DEFAULT 0,
    text_preview TEXT NOT NULL DEFAULT '', text_content TEXT NOT NULL DEFAULT '',
    cache_path TEXT NOT NULL DEFAULT '', source_paths TEXT NOT NULL DEFAULT '[]',
    error TEXT NOT NULL DEFAULT '', created_at INTEGER NOT NULL, expires_at INTEGER NOT NULL
);
INSERT INTO clips (id, kind, status, text_content, created_at, expires_at)
VALUES ('legacy', 2, 1, '<p>旧记录</p>', 1, 9999999999);`); err != nil {
		t.Fatalf("构造旧库: %v", err)
	}
	old.Close()

	// 打开两次：第二次列已存在，迁移必须是幂等的
	for i := range 2 {
		s, err := Open(path)
		if err != nil {
			t.Fatalf("第 %d 次打开: %v", i+1, err)
		}
		got, err := s.GetClip("legacy")
		if err != nil {
			t.Fatalf("读取旧记录: %v", err)
		}
		if got.TextContent != "<p>旧记录</p>" || got.TextPlain != "" {
			t.Errorf("旧记录内容异常: %+v", got)
		}
		s.Close()
	}
}
