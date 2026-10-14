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
	session int    // 每打开一次加一，计时器据此认出自己是不是这一轮的
	changed func() // 打开、关闭时调用，让界面刷新状态
}

// Set 打开详细日志 d 这么久，d 不大于 0 时立即关闭。返回打开到什么时候，关闭时为零值。
func (v *verboseLog) Set(d time.Duration) time.Time {
	v.mu.Lock()
	if v.timer != nil {
		v.timer.Stop()
		v.timer = nil
	}
	v.session++
	if d <= 0 {
		v.until = time.Time{}
		logLevel.Set(v.base)
		slog.Info("详细日志已关闭")
	} else {
		until := time.Now().Add(d)
		v.until = until
		logLevel.Set(slog.LevelDebug)
		slog.Info("详细日志已打开", "到", until.Format("15:04:05"))
		session := v.session
		v.timer = time.AfterFunc(d, func() { v.expire(session) })
	}
	until := v.until
	v.mu.Unlock()
	if v.changed != nil {
		v.changed()
	}
	return until
}

// expire 是计时器到点后的关闭。期间用户又打开或关闭过的话，这次到点不算数。
// 不按截止时间比较：Windows 的时钟精度粗，连着打开两次，两次的截止时间可能一模一样
func (v *verboseLog) expire(session int) {
	v.mu.Lock()
	current := v.session == session
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
