//go:build !windows

package main

import "log/slog"

// 只有 Windows 有服务管理器；其他系统用 systemd、launchd 直接运行可执行文件。

func isService() bool { return false }

func defaultServiceLog() string { return "" }

func runService(serverConfig, *slog.Logger) error { return nil }
