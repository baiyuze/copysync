// Package clipboard 是跨平台剪贴板访问层。
//
// ⚠️ 线程约束（M-1 实测，见 spikes/results.md）
//
// macOS 的剪贴板隐私机制开启后，剪贴板 API 会经 ViewBridge 与系统 UI 服务
// 通信以便显示授权弹窗，其内部 NSCFRunLoopSemaphore 依赖 run loop。
// 在裸 dispatch queue / 非主线程上调用会直接 SIGSEGV：
//
//	reportAccessBehavior → AppKit → ViewBridge -[NSCFRunLoopSemaphore wait]
//	                     → NSException raise → EXC_BAD_ACCESS
//
// 因此本包的所有方法**必须**在主线程调用。用 MainLoop 来保证这一点：
// main() 里 runtime.LockOSThread() 锁住主线程跑 MainLoop，
// 业务 goroutine 通过 Post/Call 把操作投递过去。
package clipboard

import (
	"context"
	"errors"
	"fmt"
)

// Kind 是剪贴板内容的粗分类。
type Kind int

const (
	KindUnknown Kind = iota
	KindText
	KindHTML
	KindImage
	KindFile // 一个或多个文件/文件夹
)

func (k Kind) String() string {
	switch k {
	case KindText:
		return "text"
	case KindHTML:
		return "html"
	case KindImage:
		return "image"
	case KindFile:
		return "file"
	default:
		return "unknown"
	}
}

// Permission 对应 macOS 的 NSPasteboardAccessBehavior。
// 其他平台恒为 PermissionNotApplicable。
type Permission int

const (
	PermissionNotApplicable Permission = iota // 系统无此机制
	PermissionDefault                         // 尚未触发过授权弹窗
	PermissionAsk                             // 每次程序化读取都询问
	PermissionAlwaysAllow                     // 永久放行
	PermissionAlwaysDeny                      // 永久拒绝
)

// CanReadSilently 表示此刻读取内容是否不会弹窗阻塞。
//
// 只有 AlwaysAllow 才安全。其余状态下 daemon 不应贸然读取——
// 后台进程的授权弹窗可能根本不可见，一旦读取就会无限阻塞。
// 正确做法是上报状态给 UI，由 UI 引导用户去系统设置授权。
func (p Permission) CanReadSilently() bool {
	return p == PermissionNotApplicable || p == PermissionAlwaysAllow
}

func (p Permission) String() string {
	switch p {
	case PermissionNotApplicable:
		return "not_applicable"
	case PermissionDefault:
		return "default"
	case PermissionAsk:
		return "ask"
	case PermissionAlwaysAllow:
		return "always_allow"
	case PermissionAlwaysDeny:
		return "always_deny"
	default:
		return "unknown"
	}
}

// Snapshot 是对剪贴板的**免提示**探测结果。
//
// 只含类型信息，不含任何内容。M-1 实测这些操作全部免授权且耗时 < 0.03s：
// changeCount / types / detectMetadataForTypes。
type Snapshot struct {
	ChangeCount int64
	Kind        Kind
	ContentType string   // UTType 标识符，如 public.plain-text
	Types       []string // 原始类型列表，排障用
	// Sensitive 表示复制它的程序（多为密码管理器）标明了「不要记录、不要同步」。
	// 目前只有 Windows 实现会设置，见 windows.go 的 sensitive
	Sensitive bool
}

// Content 是剪贴板的实际内容。读取它可能触发授权弹窗。
type Content struct {
	Kind  Kind
	Text  string   // KindText / KindHTML 的纯文本形式
	HTML  string   // KindHTML 的原始标记
	Image []byte   // KindImage 的 PNG 字节
	Files []string // KindFile 的绝对路径
}

// Clipboard 是平台实现要满足的接口。所有方法必须在主线程调用。
type Clipboard interface {
	// Peek 免提示地探测剪贴板类型。可安全地高频调用。
	Peek() (Snapshot, error)

	// Read 读取实际内容。调用前务必先检查 Permission().CanReadSilently()，
	// 否则可能阻塞在不可见的授权弹窗上。
	Read() (Content, error)

	// Write 覆盖本机剪贴板。文件类型要求路径已真实存在于本地磁盘。
	Write(Content) error

	// Permission 返回当前授权状态，本身免提示。
	Permission() Permission

	// OpenPermissionSettings 打开系统设置中对应的隐私面板，引导用户授权。
	OpenPermissionSettings() error
}

var (
	ErrUnsupported = errors.New("当前平台不支持该剪贴板操作")
	// ErrPermissionRequired 表示用户尚未授予剪贴板读取权限。
	// daemon 遇到它应当上报给界面去引导授权，而不是反复重试。
	ErrPermissionRequired = errors.New("需要授予剪贴板访问权限")
)

// ─────────────────────── 主线程调度 ───────────────────────

// MainLoop 在主线程上运行平台的事件循环，并执行投递进来的任务。
//
// 用法（main 函数）：
//
//	func main() {
//	    runtime.LockOSThread()
//	    loop := clipboard.NewMainLoop()
//	    go daemon.Run(loop)   // 业务逻辑跑在别的 goroutine
//	    loop.Run(ctx)          // 主线程在此阻塞
//	}
type MainLoop struct {
	tasks chan func()
}

func NewMainLoop() *MainLoop {
	// 带缓冲，避免投递方在主线程忙时阻塞
	return &MainLoop{tasks: make(chan func(), 64)}
}

// Post 异步投递任务到主线程，不等待结果。
func (l *MainLoop) Post(fn func()) {
	select {
	case l.tasks <- fn:
	default:
		// 队列满说明主线程卡住了（很可能阻塞在授权弹窗上）。
		// 丢弃而非阻塞调用方——剪贴板轮询丢一次无所谓，卡死业务则不可接受。
	}
}

// Call 投递任务并等待其在主线程执行完毕。
//
// 注意：若主线程正阻塞在授权弹窗上，这里会一直等到 ctx 取消。
// 所以调用方必须带一个有超时的 ctx。
func (l *MainLoop) Call(ctx context.Context, fn func()) error {
	done := make(chan struct{})
	wrapped := func() {
		defer close(done)
		fn()
	}
	select {
	case l.tasks <- wrapped:
	case <-ctx.Done():
		return ctx.Err()
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Run 在当前线程（必须是已 LockOSThread 的主线程）处理投递进来的任务，
// 直到 ctx 取消。
//
// 等待任务的间隙会驱动平台事件循环：macOS 上剪贴板 API 的跨进程回调与
// 授权弹窗都依赖 run loop 转动，只 select 等 channel 会让它们永远得不到处理。
func (l *MainLoop) Run(ctx context.Context) {
	const pumpInterval = 0.05 // 秒。足够跟手，又不会空转烧 CPU

	for {
		select {
		case <-ctx.Done():
			return
		case fn := <-l.tasks:
			fn()
		default:
			// 没有待办任务时转动一小段 run loop，然后回来再看队列
			pumpRunLoop(pumpInterval)
		}
	}
}

// CallResult 是 Call 的泛型版本，用于需要返回值的场景。
func CallResult[T any](ctx context.Context, l *MainLoop, fn func() (T, error)) (T, error) {
	var (
		val T
		err error
	)
	if cerr := l.Call(ctx, func() { val, err = fn() }); cerr != nil {
		var zero T
		return zero, fmt.Errorf("投递到主线程失败: %w", cerr)
	}
	return val, err
}
