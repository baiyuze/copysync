//go:build linux

package main

import (
	"context"
	"os"
	"path/filepath"
	"syscall"
)

// Linux 上没有托盘：主流桌面（GNOME）默认不显示托盘图标，界面从应用菜单打开。
func startTray(context.Context, context.CancelFunc, bool) func() { return func() {} }

// acquireInstance 保证同一个数据目录只有一个后台服务在运行。用文件锁：进程不管怎么退出，
// 内核都会释放它，不会留下僵死的锁。
func acquireInstance(root string) (release func(), ok bool) {
	return lockFile(filepath.Join(root, "daemon.lock"))
}

// acquireSupervisorLock 保证同一个数据目录只有一个守护进程。
func acquireSupervisorLock(root string) (release func(), ok bool) {
	return lockFile(filepath.Join(root, "supervisor.lock"))
}

func lockFile(path string) (release func(), ok bool) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return func() {}, true // 建不了锁就不做单实例保护，总比起不来好
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return func() {}, true
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		if err == syscall.EWOULDBLOCK {
			return nil, false
		}
		return func() {}, true
	}
	return func() { f.Close() }, true
}

// childAttr 让子进程在守护进程退出时随之结束：界面重启服务时只要结束守护进程，
// 不会留下一个没人看着的后台服务。
func childAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Pdeathsig: syscall.SIGTERM}
}
