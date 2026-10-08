//go:build !windows

package main

import (
	"io"
	"os"
)

// Mac 上这些都由 launchd 负责：日志重定向到文件、只运行一个实例、崩溃自动拉起。

func logOutput(string) io.Writer { return os.Stderr }

func acquireInstance(string) (release func(), ok bool) { return func() {}, true }

func runSupervisor(string) int {
	os.Stderr.WriteString("-supervise 只在 Windows 上使用；这个平台由系统服务管理器负责拉起\n")
	return 2
}
