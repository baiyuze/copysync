package clipboard

import (
	"context"
	"log/slog"
	"time"
)

// Watcher 监听本机剪贴板变化。
//
// macOS 不提供剪贴板变更通知，只能轮询 changeCount。好在 M-1 实测表明
// changeCount、类型列表、detectMetadata 这三项都免授权且合计不到 0.03 秒，
// 因此 300ms 的轮询既跟手又不会造成负担。
type Watcher struct {
	cb   Clipboard
	loop *MainLoop
	log  *slog.Logger

	interval time.Duration
	onChange func(Snapshot)
	// onPermission 在授权状态变化时回调，界面据此引导用户去系统设置
	onPermission func(Permission)

	lastChange     int64
	lastPermission Permission
	// selfWrites 记录本机主动写入剪贴板时的 changeCount，
	// 避免把"自己刚写进去的内容"当成用户的新复制再广播一遍。
	selfWrites map[int64]struct{}
}

type WatcherOptions struct {
	Clipboard    Clipboard
	MainLoop     *MainLoop
	Logger       *slog.Logger
	Interval     time.Duration
	OnChange     func(Snapshot)
	OnPermission func(Permission)
}

func NewWatcher(opts WatcherOptions) *Watcher {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	interval := opts.Interval
	if interval <= 0 {
		interval = 300 * time.Millisecond
	}
	return &Watcher{
		cb:           opts.Clipboard,
		loop:         opts.MainLoop,
		log:          log,
		interval:     interval,
		onChange:     opts.OnChange,
		onPermission: opts.OnPermission,
		lastChange:   -1,
		selfWrites:   make(map[int64]struct{}),
	}
}

// SetOnChange 设置剪贴板变化的回调。
//
// 单独提供这个方法，是因为同步引擎需要在 Watcher 之后构造
// （引擎依赖 Watcher 读写剪贴板），无法在 NewWatcher 时就把回调传进去。
func (w *Watcher) SetOnChange(fn func(Snapshot)) { w.onChange = fn }

// Run 持续轮询直到 ctx 取消。可在任意 goroutine 中启动——
// 真正接触剪贴板的操作都会被投递到主线程执行。
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.poll(ctx)
		}
	}
}

func (w *Watcher) poll(ctx context.Context) {
	// 给足超时：万一主线程正卡在授权弹窗上，这里要能干净退出而不是堆积
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	snap, err := CallResult(ctx, w.loop, w.cb.Peek)
	if err != nil {
		w.log.Debug("探测剪贴板失败", "err", err)
		return
	}

	if perm, err := CallResult(ctx, w.loop, func() (Permission, error) {
		return w.cb.Permission(), nil
	}); err == nil && perm != w.lastPermission {
		w.lastPermission = perm
		w.log.Info("剪贴板授权状态变化", "state", perm)
		if w.onPermission != nil {
			w.onPermission(perm)
		}
	}

	if snap.ChangeCount == w.lastChange {
		return
	}
	// 首轮只记录基线，不把启动前就存在的内容当作新复制广播出去
	first := w.lastChange < 0
	w.lastChange = snap.ChangeCount
	if first {
		w.log.Debug("剪贴板监听已就绪", "changeCount", snap.ChangeCount)
		return
	}

	// 本机刚写进去的内容不该被当成用户的新复制再同步一轮
	if _, mine := w.selfWrites[snap.ChangeCount]; mine {
		delete(w.selfWrites, snap.ChangeCount)
		w.log.Debug("跳过本机写入触发的变化", "changeCount", snap.ChangeCount)
		return
	}

	w.log.Debug("检测到剪贴板变化",
		"changeCount", snap.ChangeCount, "kind", snap.Kind,
		"type", snap.ContentType, "types", snap.Types)
	if w.onChange != nil {
		w.onChange(snap)
	}
}

// Read 在主线程读取剪贴板内容，并先行检查授权状态。
func (w *Watcher) Read(ctx context.Context) (Content, error) {
	perm, err := CallResult(ctx, w.loop, func() (Permission, error) {
		return w.cb.Permission(), nil
	})
	if err != nil {
		return Content{}, err
	}
	// 未授权时不要贸然读取：会阻塞在一个对后台进程可能不可见的弹窗上
	if !perm.CanReadSilently() {
		return Content{}, ErrPermissionRequired
	}
	return CallResult(ctx, w.loop, w.cb.Read)
}

// Write 在主线程写入剪贴板，并记下由此产生的 changeCount，
// 使下一轮轮询不会把它误判为用户的新复制。
func (w *Watcher) Write(ctx context.Context, content Content) error {
	return w.loop.Call(ctx, func() {
		if err := w.cb.Write(content); err != nil {
			w.log.Warn("写入剪贴板失败", "err", err)
			return
		}
		if snap, err := w.cb.Peek(); err == nil {
			w.selfWrites[snap.ChangeCount] = struct{}{}
			w.lastChange = snap.ChangeCount
		}
	})
}

func (w *Watcher) Permission(ctx context.Context) Permission {
	perm, err := CallResult(ctx, w.loop, func() (Permission, error) {
		return w.cb.Permission(), nil
	})
	if err != nil {
		return PermissionNotApplicable
	}
	return perm
}

func (w *Watcher) OpenPermissionSettings(ctx context.Context) error {
	return w.loop.Call(ctx, func() {
		if err := w.cb.OpenPermissionSettings(); err != nil {
			w.log.Warn("打开授权设置失败", "err", err)
		}
	})
}
