//go:build windows

package clipboard

// Windows 的剪贴板实现。纯 Go 调用 Win32，不用 cgo，在 Mac 上也能交叉编译。
//
// 线程模型与 macOS 相同：主线程（已 LockOSThread）创建一个隐藏窗口作为剪贴板的
// 拥有者，并在 MainLoop 的空闲间隙泵消息。所有剪贴板操作都经 MainLoop 投递到这个
// 线程执行。
//
// 隐藏窗口用普通的顶层窗口，而不是「仅消息窗口」（HWND_MESSAGE）：后者收不到
// 注销、关机、睡眠唤醒这些只发给顶层窗口的消息。

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"log/slog"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procOpenClipboard               = user32.NewProc("OpenClipboard")
	procCloseClipboard              = user32.NewProc("CloseClipboard")
	procEmptyClipboard              = user32.NewProc("EmptyClipboard")
	procGetClipboardData            = user32.NewProc("GetClipboardData")
	procEnumClipboardFormats        = user32.NewProc("EnumClipboardFormats")
	procSetClipboardData            = user32.NewProc("SetClipboardData")
	procIsClipboardFormatAvailable  = user32.NewProc("IsClipboardFormatAvailable")
	procRegisterClipboardFormatW    = user32.NewProc("RegisterClipboardFormatW")
	procGetClipboardSequenceNumber  = user32.NewProc("GetClipboardSequenceNumber")
	procAddClipboardFormatListener  = user32.NewProc("AddClipboardFormatListener")
	procRegisterClassExW            = user32.NewProc("RegisterClassExW")
	procCreateWindowExW             = user32.NewProc("CreateWindowExW")
	procDefWindowProcW              = user32.NewProc("DefWindowProcW")
	procPeekMessageW                = user32.NewProc("PeekMessageW")
	procTranslateMessage            = user32.NewProc("TranslateMessage")
	procDispatchMessageW            = user32.NewProc("DispatchMessageW")
	procMsgWaitForMultipleObjectsEx = user32.NewProc("MsgWaitForMultipleObjectsEx")

	procGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	procGlobalAlloc      = kernel32.NewProc("GlobalAlloc")
	procGlobalFree       = kernel32.NewProc("GlobalFree")
	procGlobalLock       = kernel32.NewProc("GlobalLock")
	procGlobalUnlock     = kernel32.NewProc("GlobalUnlock")
	procGlobalSize       = kernel32.NewProc("GlobalSize")
)

// 标准剪贴板格式
const (
	cfUnicodeText = 13
	cfHDROP       = 15
	cfDIB         = 8
	cfDIBV5       = 17
)

// 窗口消息
const (
	wmQueryEndSession  = 0x0011
	wmEndSession       = 0x0016
	wmPowerBroadcast   = 0x0218
	wmClipboardUpdate  = 0x031D
	pbtAPMResumeAuto   = 0x0012
	pmRemove           = 0x0001
	qsAllInput         = 0x04FF
	mwmoInputAvailable = 0x0004
	gmemMoveable       = 0x0002
	dropEffectCopy     = 1
)

// 注册格式：按名字向系统换一个编号，各程序用同一个名字就能互通。
var formats struct {
	once sync.Once

	png, html, dropEffect uint32
	// 密码管理器用来声明「别记录、别同步」的标记，见 sensitive
	exclude, viewerIgnore, canIncludeHistory, canUploadCloud uint32
}

func registerFormats() {
	formats.once.Do(func() {
		reg := func(name string) uint32 {
			p, _ := windows.UTF16PtrFromString(name)
			r, _, _ := procRegisterClipboardFormatW.Call(uintptr(unsafe.Pointer(p)))
			return uint32(r)
		}
		formats.png = reg("PNG")
		formats.html = reg("HTML Format")
		formats.dropEffect = reg("Preferred DropEffect")
		formats.exclude = reg("ExcludeClipboardContentFromMonitorProcessing")
		formats.viewerIgnore = reg("Clipboard Viewer Ignore")
		formats.canIncludeHistory = reg("CanIncludeInClipboardHistory")
		formats.canUploadCloud = reg("CanUploadToCloudClipboard")
	})
}

// owner 是隐藏窗口的句柄。写剪贴板必须有拥有者窗口：用空句柄打开剪贴板时，
// EmptyClipboard 会把拥有者设为空，随后的 SetClipboardData 一律失败。
var owner uintptr

type windowsClipboard struct{}

func New() Clipboard { return windowsClipboard{} }

// InitPlatform 在主线程上创建隐藏窗口并订阅剪贴板变化。
// 必须与 MainLoop.Run 在同一个线程上调用：窗口的消息只投递给创建它的线程。
func InitPlatform() {
	registerFormats()
	changeNotify = make(chan struct{}, 1)
	if err := createOwnerWindow(); err != nil {
		slog.Error("创建剪贴板窗口失败，无法写入剪贴板", "err", err)
		return
	}
	if r, _, err := procAddClipboardFormatListener.Call(owner); r == 0 {
		// 不致命：Watcher 的轮询兜底
		slog.Warn("订阅剪贴板变化失败，退回轮询", "err", err)
	}
}

type wndClassEx struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

type winMsg struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
	Private uint32
}

// wndProcCallback 必须是包级变量：NewCallback 生成的回调数量有上限，且不能被回收。
var wndProcCallback = windows.NewCallback(wndProc)

func createOwnerWindow() error {
	instance, _, _ := procGetModuleHandleW.Call(0)
	className, _ := windows.UTF16PtrFromString("CopySyncClipboardOwner")
	wc := wndClassEx{
		Size:      uint32(unsafe.Sizeof(wndClassEx{})),
		WndProc:   wndProcCallback,
		Instance:  instance,
		ClassName: className,
	}
	if r, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return fmt.Errorf("注册窗口类: %w", err)
	}
	title, _ := windows.UTF16PtrFromString("CopySync")
	// 样式为 0、尺寸为 0，且从不 ShowWindow：用户看不到，任务栏上也没有
	hwnd, _, err := procCreateWindowExW.Call(0,
		uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(title)),
		0, 0, 0, 0, 0, 0, 0, instance, 0)
	if hwnd == 0 {
		return fmt.Errorf("创建窗口: %w", err)
	}
	owner = hwnd
	return nil
}

func wndProc(hwnd, msg, wparam, lparam uintptr) uintptr {
	switch msg {
	case wmClipboardUpdate:
		select {
		case changeNotify <- struct{}{}:
		default: // 已有一个未处理的信号，合并
		}
		return 0
	case wmQueryEndSession:
		return 1 // 同意结束会话
	case wmEndSession:
		if wparam != 0 && onSessionEnd != nil {
			onSessionEnd()
		}
		return 0
	case wmPowerBroadcast:
		if wparam == pbtAPMResumeAuto && onResume != nil {
			go onResume()
		}
		return 1
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wparam, lparam)
	return r
}

// pumpRunLoop 等待最多 seconds 秒，期间有消息就分发掉。
func pumpRunLoop(seconds float64) {
	procMsgWaitForMultipleObjectsEx.Call(0, 0, uintptr(seconds*1000), qsAllInput, mwmoInputAvailable)
	var m winMsg
	for {
		if r, _, _ := procPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, pmRemove); r == 0 {
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

// ─────────────────────────── 探测 ───────────────────────────

func available(format uint32) bool {
	if format == 0 {
		return false
	}
	r, _, _ := procIsClipboardFormatAvailable.Call(uintptr(format))
	return r != 0
}

// Peek 只看有哪些格式，不读内容，也不打开剪贴板（判定敏感标记时例外）。
func (windowsClipboard) Peek() (Snapshot, error) {
	registerFormats()
	seq, _, _ := procGetClipboardSequenceNumber.Call()
	snap := Snapshot{ChangeCount: int64(seq)}

	known := []struct {
		id   uint32
		name string
	}{
		{cfHDROP, "CF_HDROP"}, {formats.png, "PNG"}, {cfDIBV5, "CF_DIBV5"}, {cfDIB, "CF_DIB"},
		{formats.html, "HTML Format"}, {cfUnicodeText, "CF_UNICODETEXT"},
	}
	for _, f := range known {
		if available(f.id) {
			snap.Types = append(snap.Types, f.name)
		}
	}

	hasText := available(cfUnicodeText)
	hasImage := available(formats.png) || available(cfDIBV5) || available(cfDIB)
	switch {
	case available(cfHDROP):
		// 资源管理器复制文件时同时带着文件名文本，文件优先
		snap.Kind, snap.ContentType = KindFile, "public.file-url"
	case hasImage && !hasText:
		// 截图、画图、浏览器里「复制图片」只有图片格式。Excel、Word 复制的内容
		// 同时带着一份渲染出来的图片，但有文本，应当按文本处理
		snap.Kind, snap.ContentType = KindImage, "public.png"
	case available(formats.html):
		snap.Kind, snap.ContentType = KindHTML, "public.html"
	case hasText:
		snap.Kind, snap.ContentType = KindText, "public.utf8-plain-text"
	case hasImage:
		snap.Kind, snap.ContentType = KindImage, "public.png"
	}
	snap.Sensitive = sensitive()
	return snap, nil
}

// sensitive 判断复制内容的程序是否要求「别记录、别同步」。
//
// 约定见微软文档「Cloud Clipboard and Clipboard History Formats」：
//   - ExcludeClipboardContentFromMonitorProcessing、Clipboard Viewer Ignore：存在即是
//   - CanIncludeInClipboardHistory、CanUploadToCloudClipboard：值为 0 才是。
//     Bitwarden 只写前者，只看存在与否会漏掉、也会误伤
func sensitive() bool {
	if available(formats.exclude) || available(formats.viewerIgnore) {
		return true
	}
	for _, f := range []uint32{formats.canIncludeHistory, formats.canUploadCloud} {
		if !available(f) {
			continue
		}
		v, err := readDWORD(f)
		if err != nil || v == 0 {
			return true // 有标记却读不出来，宁可当作敏感
		}
	}
	return false
}

func readDWORD(format uint32) (uint32, error) {
	if err := openClipboard(); err != nil {
		return 0, err
	}
	defer procCloseClipboard.Call()
	b, err := getData(format)
	if err != nil {
		return 0, err
	}
	if len(b) < 4 {
		return 0, errors.New("数据过短")
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, nil
}

// ─────────────────────────── 读取 ───────────────────────────

func (c windowsClipboard) Read() (Content, error) {
	snap, err := c.Peek()
	if err != nil {
		return Content{}, err
	}
	if err := openClipboard(); err != nil {
		return Content{}, err
	}
	defer procCloseClipboard.Call()

	out := Content{Kind: snap.Kind}
	switch snap.Kind {
	case KindFile:
		b, err := getData(cfHDROP)
		if err != nil {
			return Content{}, err
		}
		if out.Files, err = decodeDropFiles(b); err != nil {
			return Content{}, err
		}
		if len(out.Files) == 0 {
			return Content{}, errors.New("剪贴板中没有可读取的文件路径")
		}

	case KindText:
		out.Text = readText()

	case KindHTML:
		b, err := getData(formats.html)
		if err != nil {
			return Content{}, err
		}
		if out.HTML, err = decodeCFHTML(b); err != nil {
			return Content{}, err
		}
		out.Text = readText()

	case KindImage:
		if out.Image, err = readImage(); err != nil {
			return Content{}, err
		}

	default:
		return Content{}, fmt.Errorf("不支持的剪贴板类型：%v", snap.Types)
	}
	return out, nil
}

func readText() string {
	b, err := getData(cfUnicodeText)
	if err != nil {
		return ""
	}
	return fromCRLF(utf16String(b))
}

// readImage 优先取 PNG（保留透明度、无损），没有再从位图转。
func readImage() ([]byte, error) {
	if b, err := getData(formats.png); err == nil {
		if _, err := png.DecodeConfig(bytes.NewReader(b)); err == nil {
			return b, nil
		}
	}
	// Windows 枚举时先返回应用提供的格式，再返回由它合成的格式。
	// 保留这个顺序：系统把普通 CF_DIB 转成 CF_DIBV5 时可能保留多余的
	// 颜色掩码，导致像素偏移；而原生 CF_DIBV5 应优先读取，以保留透明度。
	bitmapFormats := []uint32{cfDIBV5, cfDIB}
	for f := uintptr(0); ; {
		f, _, _ = procEnumClipboardFormats.Call(f)
		if f == 0 || f == cfDIBV5 {
			break
		}
		if f == cfDIB {
			bitmapFormats = []uint32{cfDIB, cfDIBV5}
			break
		}
	}
	var lastErr error
	for _, f := range bitmapFormats {
		b, err := getData(f)
		if err != nil {
			lastErr = err
			continue
		}
		img, err := dibToImage(b)
		if err != nil {
			lastErr = err
			continue
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("读取剪贴板图片失败: %w", lastErr)
}

// ─────────────────────────── 写入 ───────────────────────────

// Write 覆盖本机剪贴板。各格式的数据先准备好，再打开剪贴板一次写完：
// 剪贴板打开期间别的程序都用不了它，要尽量短。
func (windowsClipboard) Write(content Content) error {
	registerFormats()
	type item struct {
		format uint32
		data   []byte
	}
	var items []item

	switch content.Kind {
	case KindFile:
		if len(content.Files) == 0 {
			return errors.New("没有要写入的文件路径")
		}
		items = append(items,
			item{cfHDROP, encodeDropFiles(content.Files)},
			// 资源管理器粘贴时据此执行复制而不是移动：缓存里的文件不能被挪走
			item{formats.dropEffect, dwordBytes(dropEffectCopy)})

	case KindText:
		items = append(items, item{cfUnicodeText, utf16Bytes(toCRLF(content.Text))})

	case KindHTML:
		items = append(items, item{formats.html, encodeCFHTML(content.HTML)})
		if content.Text != "" {
			items = append(items, item{cfUnicodeText, utf16Bytes(toCRLF(content.Text))})
		}

	case KindImage:
		if len(content.Image) == 0 {
			return errors.New("没有要写入的图片数据")
		}
		img, err := png.Decode(bytes.NewReader(content.Image))
		if err != nil {
			return fmt.Errorf("解析图片: %w", err)
		}
		// PNG 给 Office、浏览器等新程序；位图给只认位图的程序（画图等），
		// CF_DIB 与 CF_BITMAP 由系统从 CF_DIBV5 自动合成
		items = append(items, item{formats.png, content.Image}, item{cfDIBV5, imageToDIBV5(img)})

	default:
		return fmt.Errorf("不支持写入的类型：%v", content.Kind)
	}

	// 别让 Windows 自带的云剪贴板再把它同步到用户的其他电脑：内容本来就是从别的设备来的。
	// 带着这个标记，Peek 会把它判为敏感；Watcher 先按 selfWrites 跳过本端写入，不受影响
	items = append(items, item{formats.canUploadCloud, dwordBytes(0)})

	if owner == 0 {
		return errors.New("剪贴板窗口未创建，无法写入")
	}
	if err := openClipboard(); err != nil {
		return err
	}
	defer procCloseClipboard.Call()
	if r, _, err := procEmptyClipboard.Call(); r == 0 {
		return fmt.Errorf("清空剪贴板: %w", err)
	}
	for _, it := range items {
		if err := setData(it.format, it.data); err != nil {
			return err
		}
	}
	return nil
}

func (windowsClipboard) Permission() Permission { return PermissionNotApplicable }

func (windowsClipboard) OpenPermissionSettings() error { return nil }

// ─────────────────────────── 底层 ───────────────────────────

// openClipboard 打开剪贴板。别的程序正打开着时会失败，稍等重试。
func openClipboard() error {
	var lastErr error
	for range 10 {
		if r, _, err := procOpenClipboard.Call(owner); r != 0 {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(20 * time.Millisecond)
	}
	return fmt.Errorf("剪贴板被其他程序占用: %w", lastErr)
}

// getData 复制出某个格式的数据。剪贴板必须已打开。
func getData(format uint32) ([]byte, error) {
	h, _, err := procGetClipboardData.Call(uintptr(format))
	if h == 0 {
		return nil, fmt.Errorf("读取剪贴板格式 %d: %w", format, err)
	}
	size, _, _ := procGlobalSize.Call(h)
	p, _, err := procGlobalLock.Call(h)
	if p == 0 {
		return nil, fmt.Errorf("锁定剪贴板数据: %w", err)
	}
	defer procGlobalUnlock.Call(h)
	return bytes.Clone(unsafe.Slice((*byte)(sysPtr(p)), size)), nil
}

// setData 把数据放进剪贴板。成功后内存归系统所有，失败才由本端释放。
func setData(format uint32, data []byte) error {
	if format == 0 || len(data) == 0 {
		return fmt.Errorf("无效的剪贴板数据（格式 %d）", format)
	}
	h, _, err := procGlobalAlloc.Call(gmemMoveable, uintptr(len(data)))
	if h == 0 {
		return fmt.Errorf("分配剪贴板内存: %w", err)
	}
	p, _, err := procGlobalLock.Call(h)
	if p == 0 {
		procGlobalFree.Call(h)
		return fmt.Errorf("锁定剪贴板内存: %w", err)
	}
	copy(unsafe.Slice((*byte)(sysPtr(p)), len(data)), data)
	procGlobalUnlock.Call(h)
	if r, _, err := procSetClipboardData.Call(uintptr(format), h); r == 0 {
		procGlobalFree.Call(h)
		return fmt.Errorf("写入剪贴板格式 %d: %w", format, err)
	}
	return nil
}

// sysPtr 把系统调用返回的地址转成指针。
//
// 这块内存由系统分配（GlobalLock 锁定期间地址固定），不归 Go 的垃圾回收管。
// 写成这样而不是 unsafe.Pointer(p)，是为了不让 go vet 把它误报成指针误用。
func sysPtr(p uintptr) unsafe.Pointer { return *(*unsafe.Pointer)(unsafe.Pointer(&p)) }
