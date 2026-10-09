//go:build windows

package main

// 托盘由 daemon 持有，界面关闭或尚未打开时仍可访问。
// 独立线程泵消息：用户打开右键菜单时不能阻塞剪贴板线程。

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	trayUser32           = windows.NewLazySystemDLL("user32.dll")
	trayShell32          = windows.NewLazySystemDLL("shell32.dll")
	trayKernel32         = windows.NewLazySystemDLL("kernel32.dll")
	trayNotify           = trayShell32.NewProc("Shell_NotifyIconW")
	trayExtractIcon      = trayShell32.NewProc("ExtractIconExW")
	trayRegisterClass    = trayUser32.NewProc("RegisterClassExW")
	trayUnregisterClass  = trayUser32.NewProc("UnregisterClassW")
	trayCreateWindow     = trayUser32.NewProc("CreateWindowExW")
	trayDestroyWindow    = trayUser32.NewProc("DestroyWindow")
	trayDefWindowProc    = trayUser32.NewProc("DefWindowProcW")
	trayRegisterMessage  = trayUser32.NewProc("RegisterWindowMessageW")
	trayGetMessage       = trayUser32.NewProc("GetMessageW")
	trayTranslateMessage = trayUser32.NewProc("TranslateMessage")
	trayDispatchMessage  = trayUser32.NewProc("DispatchMessageW")
	trayPostMessage      = trayUser32.NewProc("PostMessageW")
	trayPostQuitMessage  = trayUser32.NewProc("PostQuitMessage")
	traySetTimer         = trayUser32.NewProc("SetTimer")
	trayKillTimer        = trayUser32.NewProc("KillTimer")
	trayCreateMenu       = trayUser32.NewProc("CreatePopupMenu")
	trayAppendMenu       = trayUser32.NewProc("AppendMenuW")
	trayDestroyMenu      = trayUser32.NewProc("DestroyMenu")
	trayTrackMenu        = trayUser32.NewProc("TrackPopupMenu")
	trayGetCursorPos     = trayUser32.NewProc("GetCursorPos")
	traySetForeground    = trayUser32.NewProc("SetForegroundWindow")
	trayFindWindow       = trayUser32.NewProc("FindWindowW")
	trayLoadIcon         = trayUser32.NewProc("LoadIconW")
	trayDestroyIcon      = trayUser32.NewProc("DestroyIcon")
	trayGetModule        = trayKernel32.NewProc("GetModuleHandleW")
)

const (
	trayCallbackMessage = 0x8001 // WM_APP + 1
	trayWMClose         = 0x0010
	trayWMDestroy       = 0x0002
	trayWMTimer         = 0x0113
	trayOpenCommand     = 1
	trayQuitCommand     = 2
)

type trayNotifyData struct {
	Size                uint32
	Window              uintptr
	ID, Flags, Callback uint32
	Icon                uintptr
	Tip                 [128]uint16
	State, StateMask    uint32
	Info                [256]uint16
	Version             uint32
	InfoTitle           [64]uint16
	InfoFlags           uint32
	GUID                windows.GUID
	BalloonIcon         uintptr
}

type trayWindowClass struct {
	Size, Style                        uint32
	WndProc                            uintptr
	ClsExtra, WndExtra                 int32
	Instance, Icon, Cursor, Background uintptr
	MenuName, ClassName                *uint16
	IconSmall                          uintptr
}

type trayMessage struct {
	Window         uintptr
	Message        uint32
	WParam, LParam uintptr
	Time           uint32
	X, Y           int32
	Private        uint32
}

type windowsTray struct {
	data           trayNotifyData
	taskbarCreated uintptr
	uiPath         string
	quit           context.CancelFunc
}

// 回调和状态只在托盘线程使用。回调 thunk 不能反复创建（Windows 数量有限）。
var trayCallback = windows.NewCallback(trayWindowProc)
var activeTray *windowsTray

func startTray(ctx context.Context, quit context.CancelFunc, enabled bool) func() {
	if !enabled {
		return func() {}
	}
	exe, err := os.Executable()
	if err != nil {
		return func() {}
	}
	uiPath := filepath.Join(filepath.Dir(exe), "CopySync.exe")
	if _, err := os.Stat(uiPath); err != nil {
		return func() {} // 单独运行源码构建的 daemon，没有可打开的界面
	}
	trayCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		if err := runTray(trayCtx, uiPath, quit); err != nil {
			slog.Warn("启动托盘图标失败", "err", err)
		}
	}()
	return func() {
		cancel()
		<-done // 删除图标后再退出进程，避免残留无响应的图标
	}
}

func runTray(ctx context.Context, uiPath string, quit context.CancelFunc) error {
	class, _ := windows.UTF16PtrFromString("CopySyncTrayWindow")
	instance, _, _ := trayGetModule.Call(0)
	wc := trayWindowClass{Size: uint32(unsafe.Sizeof(trayWindowClass{})), WndProc: trayCallback,
		Instance: instance, ClassName: class}
	if r, _, err := trayRegisterClass.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return fmt.Errorf("注册托盘窗口: %w", err)
	}
	defer trayUnregisterClass.Call(uintptr(unsafe.Pointer(class)), instance)
	t := &windowsTray{uiPath: uiPath, quit: quit}
	activeTray = t
	defer func() { activeTray = nil }()
	title, _ := windows.UTF16PtrFromString("CopySync Tray")
	hwnd, _, err := trayCreateWindow.Call(0, uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(title)),
		0, 0, 0, 0, 0, 0, 0, instance, 0)
	if hwnd == 0 {
		return fmt.Errorf("创建托盘窗口: %w", err)
	}
	defer trayDestroyWindow.Call(hwnd)
	t.data = trayNotifyData{Size: uint32(unsafe.Sizeof(trayNotifyData{})), Window: hwnd,
		ID: 1, Flags: 1 | 2 | 4 | 0x20 | 0x80, Callback: trayCallbackMessage,
		// NIF_MESSAGE | NIF_ICON | NIF_TIP | NIF_GUID | NIF_SHOWTIP。
		// 稳定 GUID 保留用户的托盘显示偏好，也能清掉崩溃前遗留的图标。
		GUID: windows.GUID{Data1: 0x33e590c2, Data2: 0x9bfd, Data3: 0x47a3,
			Data4: [8]byte{0xae, 0x76, 0x84, 0x2b, 0xef, 0xa2, 0xf6, 0xca}}}
	tip, _ := windows.UTF16FromString("CopySync · 后台同步正在运行")
	copy(t.data.Tip[:], tip)
	path, _ := windows.UTF16PtrFromString(uiPath)
	trayExtractIcon.Call(uintptr(unsafe.Pointer(path)), 0, 0, uintptr(unsafe.Pointer(&t.data.Icon)), 1)
	if t.data.Icon != 0 {
		defer trayDestroyIcon.Call(t.data.Icon)
	} else {
		t.data.Icon, _, _ = trayLoadIcon.Call(0, 32512) // IDI_APPLICATION（共享句柄）
	}
	message, _ := windows.UTF16PtrFromString("TaskbarCreated")
	t.taskbarCreated, _, _ = trayRegisterMessage.Call(uintptr(unsafe.Pointer(message)))
	// 异常退出时 Explorer 可能留有旧图标，用同一 GUID 删除后再注册。
	trayNotify.Call(2, uintptr(unsafe.Pointer(&t.data))) // NIM_DELETE
	t.addIcon()
	defer trayNotify.Call(2, uintptr(unsafe.Pointer(&t.data)))

	loopDone := make(chan struct{})
	defer close(loopDone)
	go func() {
		select {
		case <-ctx.Done():
			trayPostMessage.Call(hwnd, trayWMClose, 0, 0)
		case <-loopDone:
		}
	}()
	var m trayMessage
	for {
		r, _, err := trayGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) == -1 {
			return fmt.Errorf("读取托盘消息: %w", err)
		}
		if r == 0 {
			return nil
		}
		trayTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		trayDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func (t *windowsTray) addIcon() {
	if r, _, _ := trayNotify.Call(0, uintptr(unsafe.Pointer(&t.data))); r == 0 {
		// 重复通知时可能已经存在，先尝试更新；登录时 Explorer 尚未就绪则稍后重试。
		if r, _, _ := trayNotify.Call(1, uintptr(unsafe.Pointer(&t.data))); r == 0 {
			traySetTimer.Call(t.data.Window, 1, 2000, 0)
			return
		}
	}
	t.data.Version = 4
	trayNotify.Call(4, uintptr(unsafe.Pointer(&t.data))) // NIM_SETVERSION
	trayKillTimer.Call(t.data.Window, 1)
}

func (t *windowsTray) openUI() {
	cmd := exec.Command(t.uiPath)
	if err := cmd.Start(); err != nil {
		slog.Warn("从托盘打开界面失败", "err", err)
		return
	}
	go cmd.Wait()
}

func (t *windowsTray) showMenu() {
	menu, _, _ := trayCreateMenu.Call()
	if menu == 0 {
		return
	}
	defer trayDestroyMenu.Call(menu)
	openLabel, quitLabel := trayLabels(trayLanguage())
	open, _ := windows.UTF16PtrFromString(openLabel)
	exit, _ := windows.UTF16PtrFromString(quitLabel)
	trayAppendMenu.Call(menu, 0, trayOpenCommand, uintptr(unsafe.Pointer(open)))
	trayAppendMenu.Call(menu, 0x800, 0, 0) // MF_SEPARATOR
	trayAppendMenu.Call(menu, 0, trayQuitCommand, uintptr(unsafe.Pointer(exit)))
	var point struct{ X, Y int32 }
	trayGetCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	traySetForeground.Call(t.data.Window)
	// TPM_RETURNCMD | TPM_NONOTIFY | TPM_RIGHTBUTTON。
	command, _, _ := trayTrackMenu.Call(menu, 0x100|0x80|2, uintptr(point.X), uintptr(point.Y), 0, t.data.Window, 0)
	trayPostMessage.Call(t.data.Window, 0, 0, 0) // WM_NULL：点击菜单外部后正确收起
	switch command {
	case trayOpenCommand:
		t.openUI()
	case trayQuitCommand:
		class, _ := windows.UTF16PtrFromString("FLUTTER_RUNNER_WIN32_WINDOW")
		title, _ := windows.UTF16PtrFromString("CopySync")
		if hwnd, _, _ := trayFindWindow.Call(uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(title))); hwnd != 0 {
			trayPostMessage.Call(hwnd, trayWMClose, 0, 0)
		}
		t.quit() // 正常退出码 0：守护进程也退出，不会误当成崩溃再拉起
	}
}

func trayWindowProc(hwnd, message, wparam, lparam uintptr) uintptr {
	t := activeTray
	if t != nil {
		if t.taskbarCreated != 0 && message == t.taskbarCreated {
			t.addIcon()
			return 0
		}
		switch message {
		case trayCallbackMessage:
			switch lparam & 0xffff { // NOTIFYICON_VERSION_4：事件在低 16 位
			case 0x400, 0x401, 0x203: // NIN_SELECT、NIN_KEYSELECT、WM_LBUTTONDBLCLK
				t.openUI()
			case 0x7b: // WM_CONTEXTMENU
				t.showMenu()
			}
			return 0
		case trayWMTimer:
			t.addIcon()
			return 0
		case trayWMClose:
			trayDestroyWindow.Call(hwnd)
			return 0
		case trayWMDestroy:
			trayPostQuitMessage.Call(0)
			return 0
		}
	}
	r, _, _ := trayDefWindowProc.Call(hwnd, message, wparam, lparam)
	return r
}
