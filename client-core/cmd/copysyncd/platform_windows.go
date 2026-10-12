//go:build windows

package main

// Windows 上的单实例保护用命名互斥量。日志与守护进程见 supervisor.go。
//
// 开机自启由界面在 HKCU\...\Run 下登记「copysyncd.exe -supervise」，
// 守护进程再拉起真正干活的子进程（见 runSupervisor）。

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// acquireInstance 保证同一个数据目录只有一个后台服务在运行。
// 用命名互斥量：进程不管怎么退出，系统都会释放它，不会留下僵死的锁。
func acquireInstance(root string) (release func(), ok bool) {
	return acquireMutex("CopySync.Daemon." + dirKey(root))
}

func acquireMutex(name string) (release func(), ok bool) {
	p, err := windows.UTF16PtrFromString(`Local\` + name)
	if err != nil {
		return func() {}, true
	}
	h, err := windows.CreateMutex(nil, false, p)
	if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		windows.CloseHandle(h)
		return nil, false
	}
	if err != nil {
		// 建不了互斥量就不做单实例保护，总比起不来好
		return func() {}, true
	}
	return func() { windows.CloseHandle(h) }, true
}

// dirKey 把数据目录变成互斥量名字的一部分：开发时用 -data-dir 起的第二个实例不受影响。
func dirKey(root string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(filepath.Clean(root))))
	return hex.EncodeToString(sum[:6])
}

// acquireSupervisorLock 保证同一个数据目录只有一个守护进程。
func acquireSupervisorLock(root string) (release func(), ok bool) {
	return acquireMutex("CopySync.Supervisor." + dirKey(root))
}

// childAttr 是守护进程拉起子进程时的进程属性。Windows 上不需要特别设置。
func childAttr() *syscall.SysProcAttr { return nil }
