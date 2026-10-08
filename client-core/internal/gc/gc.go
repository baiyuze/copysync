// Package gc 回收过期的记录与缓存文件。
//
// 记录和文件分别存放在数据库与磁盘上，两者可能因进程崩溃或事务回滚而失配，
// 因此回收是双向的：既按记录删文件，也扫描磁盘清掉无人引用的孤儿目录。
package gc

import (
	"context"
	"log/slog"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/cache"
	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/store"
)

// 启动后先等一会儿再首次回收：此时 daemon 正在建连、恢复状态，
// 没必要和它抢 I/O。
const (
	startupDelay = 20 * time.Second
	interval     = time.Hour
)

type Collector struct {
	log        *slog.Logger
	store      *store.Store
	cache      *cache.Cache
	loadConfig func() config.Config
	onRemoved  func(ids []string)
	// 便于测试注入假时钟
	now func() time.Time
}

type Options struct {
	Logger     *slog.Logger
	Store      *store.Store
	Cache      *cache.Cache
	LoadConfig func() config.Config
	OnRemoved  func(ids []string)
}

func New(opts Options) *Collector {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	return &Collector{
		log:        log,
		store:      opts.Store,
		cache:      opts.Cache,
		loadConfig: opts.LoadConfig,
		onRemoved:  opts.OnRemoved,
		now:        time.Now,
	}
}

// Run 周期性回收，直到 ctx 取消。
func (c *Collector) Run(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(startupDelay):
	}
	c.Collect()

	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			c.Collect()
		}
	}
}

// Result 汇报一次回收的成果。
type Result struct {
	RecordsRemoved int
	CachesRemoved  int
	OrphansRemoved int
}

// Collect 执行一次完整回收。
func (c *Collector) Collect() Result {
	var res Result
	now := c.now()

	// ── 第一步：删除过期记录及其缓存 ──
	expired, err := c.store.ExpiredClips(now)
	if err != nil {
		c.log.Warn("查询过期记录失败", "err", err)
		return res
	}

	ids := make([]string, 0, len(expired))
	for _, clip := range expired {
		if clip.CachePath != "" {
			if err := c.cache.Remove(clip.CachePath); err != nil {
				c.log.Warn("删除缓存失败", "clip", clip.ID, "err", err)
			} else {
				res.CachesRemoved++
			}
		}
		ids = append(ids, clip.ID)
	}
	if len(ids) > 0 {
		// 先删文件再删记录：反过来一旦中途失败，文件就永远没人认领了
		if err := c.store.DeleteClips(ids); err != nil {
			c.log.Warn("删除过期记录失败", "err", err)
		} else {
			res.RecordsRemoved = len(ids)
			if c.onRemoved != nil {
				c.onRemoved(ids)
			}
		}
	}

	// ── 第二步：清理仅缓存过期、但记录仍在有效期内的内容 ──
	// cache_ttl 可以短于 history_ttl：记录留着让用户看得见历史，
	// 但占地方的文件先腾出来，需要时可重新拉取。
	cfg := c.loadConfig()
	if cfg.CacheTTL < cfg.HistoryTTL {
		res.CachesRemoved += c.expireStaleCaches(now, cfg.CacheTTL.Std())
	}

	// ── 第三步：与磁盘对账，清掉无人引用的孤儿目录 ──
	referenced, err := c.store.CachePaths()
	if err != nil {
		c.log.Warn("读取缓存引用失败", "err", err)
	} else {
		orphans, err := c.cache.Orphans(referenced)
		if err != nil {
			c.log.Warn("扫描孤儿缓存失败", "err", err)
		} else {
			for _, o := range orphans {
				if err := c.cache.Remove(o); err == nil {
					res.OrphansRemoved++
				}
			}
		}
	}

	if res.RecordsRemoved+res.CachesRemoved+res.OrphansRemoved > 0 {
		c.log.Info("已回收过期内容",
			"记录", res.RecordsRemoved,
			"缓存", res.CachesRemoved,
			"孤儿目录", res.OrphansRemoved)
	}
	return res
}

// expireStaleCaches 释放超过 cacheTTL 的缓存文件，但保留记录本身。
func (c *Collector) expireStaleCaches(now time.Time, cacheTTL time.Duration) int {
	clips, err := c.store.ListClips(500, time.Time{}, store.KindUnspecified)
	if err != nil {
		c.log.Warn("读取记录失败", "err", err)
		return 0
	}

	var removed int
	cutoff := now.Add(-cacheTTL)
	for _, clip := range clips {
		if clip.CachePath == "" || clip.CreatedAt.After(cutoff) {
			continue
		}
		// 本机复制的记录指向用户自己的原始文件，缓存目录里通常没东西；
		// 即便有（比如截图），删掉也只影响对端再次拉取。
		if err := c.cache.Remove(clip.CachePath); err != nil {
			continue
		}
		removed++
		// 标记为已过期：记录仍在列表里可见，但点击时会提示内容已被清理
		if err := c.store.SetClipStatus(clip.ID, store.StatusExpired, "缓存已过期清理"); err != nil {
			c.log.Debug("更新过期状态失败", "clip", clip.ID, "err", err)
		}
	}
	return removed
}
