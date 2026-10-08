// Package cache 管理落盘内容的存放与回收。
package cache

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Cache 是缓存目录的封装。
//
// 每条记录占一个子目录，目录名即 clip id 派生的短标识。
// 用目录而非单文件，是因为一次复制可能包含多个文件或整个文件夹。
type Cache struct {
	root string
}

func New(root string) (*Cache, error) {
	if err := os.MkdirAll(root, 0o700); err != nil {
		return nil, fmt.Errorf("创建缓存目录: %w", err)
	}
	return &Cache{root: root}, nil
}

func (c *Cache) Root() string { return c.root }

// Dir 返回某条记录的缓存目录（相对路径），必要时创建。
func (c *Cache) Dir(clipID string) (rel, abs string, err error) {
	rel = safeName(clipID)
	abs = filepath.Join(c.root, rel)
	if err := os.MkdirAll(abs, 0o700); err != nil {
		return "", "", fmt.Errorf("创建缓存子目录: %w", err)
	}
	return rel, abs, nil
}

func (c *Cache) Abs(rel string) string {
	if rel == "" {
		return ""
	}
	return filepath.Join(c.root, rel)
}

// Exists 判断某条记录的缓存是否还在。
// 用户可能手动清理过目录，所以每次使用前都该确认。
func (c *Cache) Exists(rel string) bool {
	if rel == "" {
		return false
	}
	_, err := os.Stat(c.Abs(rel))
	return err == nil
}

// Entries 列出某条记录缓存目录下的顶层条目的绝对路径。
// 写入剪贴板时需要这些路径。
func (c *Cache) Entries(rel string) ([]string, error) {
	abs := c.Abs(rel)
	items, err := os.ReadDir(abs)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, filepath.Join(abs, it.Name()))
	}
	return out, nil
}

// Remove 删除某条记录的缓存。
func (c *Cache) Remove(rel string) error {
	if rel == "" {
		return nil
	}
	abs := c.Abs(rel)
	// 防御路径逃逸：rel 理论上都是自己生成的，但删除操作值得多一道保险
	if !strings.HasPrefix(filepath.Clean(abs), filepath.Clean(c.root)) {
		return fmt.Errorf("拒绝删除缓存目录之外的路径: %s", abs)
	}
	return os.RemoveAll(abs)
}

// Orphans 找出磁盘上存在、但已无记录引用的目录。
//
// 记录与文件分别存储，进程崩溃或数据库回滚都可能让两者失配，
// 因此 GC 需要双向对账，而不只是按记录删文件。
func (c *Cache) Orphans(referenced map[string]struct{}) ([]string, error) {
	items, err := os.ReadDir(c.root)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, it := range items {
		if !it.IsDir() {
			continue
		}
		if _, ok := referenced[it.Name()]; !ok {
			out = append(out, it.Name())
		}
	}
	return out, nil
}

// Size 统计缓存总占用。
func (c *Cache) Size() int64 {
	var total int64
	_ = filepath.Walk(c.root, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 单个文件读不到不影响整体统计
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total
}

// TempFile 在缓存目录下创建临时文件，用于边下边写。
// 与最终位置同盘，完成后可以原子改名。
func (c *Cache) TempFile(prefix string) (*os.File, error) {
	dir := filepath.Join(c.root, ".tmp")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return os.CreateTemp(dir, prefix+"-*")
}

// CleanTemp 清掉上次异常退出残留的临时文件。
func (c *Cache) CleanTemp() error {
	return os.RemoveAll(filepath.Join(c.root, ".tmp"))
}

// safeName 把 clip id 转成安全的目录名。
// clip id 是本地生成的，但仍然过滤一遍，避免路径分隔符带来的意外。
func safeName(id string) string {
	var b strings.Builder
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9',
			r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return randomName()
	}
	return b.String()
}

func randomName() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "unnamed"
	}
	return hex.EncodeToString(buf)
}

// WriteAtomic 先写临时文件再改名，避免中断留下半截文件。
func WriteAtomic(path string, r io.Reader, perm os.FileMode) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return 0, err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".partial-*")
	if err != nil {
		return 0, err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // 成功改名后这里是空操作

	n, err := io.Copy(tmp, r)
	if err != nil {
		tmp.Close()
		return n, err
	}
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		return n, err
	}
	if err := tmp.Close(); err != nil {
		return n, err
	}
	return n, os.Rename(tmpName, path)
}
