//go:build windows

package main

// Windows 上 launchd 管的那些事要自己做：日志写文件、只运行一个实例、崩溃自动重启。
//
// 开机自启由界面在 HKCU\...\Run 下登记「copysyncd.exe -supervise」，
// 守护进程再拉起真正干活的子进程（见 runSupervisor）。

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

const (
	logMaxSize = 10 << 20 // 10 MB，轮转后最多占 20 MB
	// restartThrottle 与 Mac 上 LaunchAgent 的 ThrottleInterval 一致：崩溃后至少隔这么久再拉起
	restartThrottle = 10 * time.Second
)

// logOutput 打开日志文件。后台服务以 GUI 子系统运行、没有控制台，标准错误输出不可用。
// 未处理的 panic 另写到 crash.log：日志文件会轮转，崩溃信息不能跟着被挤掉。
func logOutput(dataDir string) io.Writer {
	paths, err := resolvePaths(dataDir)
	if err != nil {
		return io.Discard
	}
	dir := filepath.Join(paths.Root, "logs")
	w, err := openRotating(filepath.Join(dir, "daemon.log"), logMaxSize)
	if err != nil {
		return io.Discard
	}
	if f, err := os.OpenFile(filepath.Join(dir, "crash.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600); err == nil {
		_ = debug.SetCrashOutput(f, debug.CrashOptions{})
		f.Close() // SetCrashOutput 会自己复制一份句柄
	}
	return w
}

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

// runSupervisor 是守护模式：拉起后台服务，它异常退出就隔一会儿再拉起。
//
// 子进程正常退出（退出码 0：用户注销、或界面要求停止）或发现已有实例在运行时，
// 守护进程也随之退出。界面要重启服务时，把两者一起结束再重新启动守护进程。
func runSupervisor(dataDir string) int {
	paths, err := resolvePaths(dataDir)
	if err != nil {
		return 1
	}
	release, ok := acquireMutex("CopySync.Supervisor." + dirKey(paths.Root))
	if !ok {
		return 0 // 已经有守护进程在看着了
	}
	defer release()

	logFile, err := openRotating(filepath.Join(paths.Root, "logs", "supervisor.log"), 1<<20)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err == nil {
		defer logFile.Close()
		log = slog.New(slog.NewTextHandler(logFile, nil))
	}

	exe, err := os.Executable()
	if err != nil {
		log.Error("找不到自身的路径", "err", err)
		return 1
	}
	var args []string
	if dataDir != "" {
		args = append(args, "-data-dir", dataDir)
	}

	for {
		started := time.Now()
		cmd := exec.Command(exe, args...)
		err := cmd.Run()
		code := 0
		var exitErr *exec.ExitError
		switch {
		case errors.As(err, &exitErr):
			code = exitErr.ExitCode()
		case err != nil:
			log.Error("启动后台服务失败", "err", err)
			code = 1
		}
		if code == 0 || code == exitAlreadyRunning {
			log.Info("后台服务已退出，守护结束", "code", code)
			return 0
		}
		wait := restartThrottle - time.Since(started)
		log.Warn("后台服务异常退出，稍后重启", "code", code, "等待", wait.Round(time.Second))
		if wait > 0 {
			time.Sleep(wait)
		}
	}
}
