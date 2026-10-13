//go:build !windows && !linux

package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

func startTray(context.Context, context.CancelFunc, bool) func() { return func() {} }

// Mac 上这些都由 launchd 负责：日志重定向到文件、只运行一个实例、崩溃自动拉起。

func logOutput(string) io.Writer { return os.Stderr }

// logFiles 是诊断报告要读的日志：launchd 把标准错误输出写到这里（见界面里写 LaunchAgent 的地方），
// 不轮转；panic 也直接写进去，没有单独的 crash.log。
func logFiles(string) (logs []string, crash string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, ""
	}
	return []string{filepath.Join(home, "Library", "Logs", "CopySync", "daemon.log")}, ""
}

func acquireInstance(string) (release func(), ok bool) { return func() {}, true }

func runSupervisor(string) int {
	os.Stderr.WriteString("-supervise 只在 Windows 和 Linux 上使用；这个平台由系统服务管理器负责拉起\n")
	return 2
}
