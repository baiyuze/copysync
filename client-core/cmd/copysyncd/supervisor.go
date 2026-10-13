//go:build windows || linux

package main

// Windows 与 Linux 上，launchd 在 Mac 上管的那些事要自己做：日志写文件、崩溃自动重启。
// 开机自启由界面登记「copysyncd -supervise」（Windows 写注册表的 Run 项，Linux 写 XDG 自启动项），
// 守护进程再拉起真正干活的子进程。单实例保护各平台不同，见 platform_windows.go、platform_linux.go。

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"time"
)

const (
	logMaxSize = 10 << 20 // 10 MB，轮转后最多占 20 MB
	// restartThrottle 与 Mac 上 LaunchAgent 的 ThrottleInterval 一致：崩溃后至少隔这么久再拉起
	restartThrottle = 10 * time.Second
)

// logOutput 打开日志文件。Windows 上后台服务以 GUI 子系统运行、没有控制台；
// Linux 上由桌面的自启动拉起，标准错误输出也没人看。
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

// logFiles 是诊断报告要读的日志，与 logOutput 写的是同一处：轮转出去的旧文件在前。
func logFiles(dataDir string) (logs []string, crash string) {
	paths, err := resolvePaths(dataDir)
	if err != nil {
		return nil, ""
	}
	dir := filepath.Join(paths.Root, "logs")
	return []string{filepath.Join(dir, "daemon.log.1"), filepath.Join(dir, "daemon.log")}, filepath.Join(dir, "crash.log")
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
	release, ok := acquireSupervisorLock(paths.Root)
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
		cmd.SysProcAttr = childAttr()
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
