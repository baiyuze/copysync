package main

import (
	"log/slog"
	"sync"
	"time"
)

// logLevel 是后台服务的日志级别：平时记 Info；用 -v 启动、或者界面里临时打开详细日志时记 Debug。
var logLevel slog.LevelVar

// verboseLog 让界面临时打开详细日志，到时间自动关回去：详细日志量大，一直开着会占满磁盘。
// 有些连接问题只在详细日志里看得出原因，反馈问题时先打开它、重现一次。
type verboseLog struct {
	mu      sync.Mutex
	base    slog.Level // 关回去时的级别：用 -v 启动的一直是 Debug
	until   time.Time
	timer   *time.Timer
	changed func() // 打开、关闭时调用，让界面刷新状态
}

// Set 打开详细日志 d 这么久，d 不大于 0 时立即关闭。返回打开到什么时候，关闭时为零值。
func (v *verboseLog) Set(d time.Duration) time.Time {
	v.mu.Lock()
	if v.timer != nil {
		v.timer.Stop()
		v.timer = nil
	}
	if d <= 0 {
		v.until = time.Time{}
		logLevel.Set(v.base)
		slog.Info("详细日志已关闭")
	} else {
		until := time.Now().Add(d)
		v.until = until
		logLevel.Set(slog.LevelDebug)
		slog.Info("详细日志已打开", "到", until.Format("15:04:05"))
		v.timer = time.AfterFunc(d, func() { v.expire(until) })
	}
	until := v.until
	v.mu.Unlock()
	if v.changed != nil {
		v.changed()
	}
	return until
}

// expire 是计时器到点后的关闭。期间用户又打开过一次的话，until 已经变了，这次到点不算数
func (v *verboseLog) expire(until time.Time) {
	v.mu.Lock()
	current := v.until.Equal(until)
	v.mu.Unlock()
	if current {
		v.Set(0)
	}
}

func (v *verboseLog) Until() time.Time {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.until
}
