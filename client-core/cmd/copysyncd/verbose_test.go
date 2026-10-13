package main

import (
	"log/slog"
	"testing"
	"time"
)

func TestVerboseLog(t *testing.T) {
	logLevel.Set(slog.LevelInfo)
	t.Cleanup(func() { logLevel.Set(slog.LevelInfo) })
	changes := 0
	v := &verboseLog{base: slog.LevelInfo, changed: func() { changes++ }}

	until := v.Set(50 * time.Millisecond)
	if until.IsZero() || logLevel.Level() != slog.LevelDebug || !v.Until().Equal(until) {
		t.Fatalf("打开后：until=%v level=%v", until, logLevel.Level())
	}
	// 到点自动关回平时的级别
	deadline := time.Now().Add(2 * time.Second)
	for logLevel.Level() != slog.LevelInfo && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if logLevel.Level() != slog.LevelInfo || !v.Until().IsZero() {
		t.Fatalf("到点后没有关闭：level=%v until=%v", logLevel.Level(), v.Until())
	}

	// 上一轮的计时器晚到一步，不能把刚刚重新打开的详细日志关掉
	old := v.Set(time.Hour)
	v.Set(time.Hour)
	v.expire(old)
	if logLevel.Level() != slog.LevelDebug || v.Until().IsZero() {
		t.Error("过期的计时器把新打开的详细日志关掉了")
	}
	v.Set(0)
	if logLevel.Level() != slog.LevelInfo || changes < 4 {
		t.Errorf("关闭后：level=%v，通知 %d 次", logLevel.Level(), changes)
	}
}

// 用 -v 启动的，关掉临时的详细日志后仍是 Debug
func TestVerboseLogKeepsFlag(t *testing.T) {
	logLevel.Set(slog.LevelDebug)
	t.Cleanup(func() { logLevel.Set(slog.LevelInfo) })
	v := &verboseLog{base: slog.LevelDebug}
	v.Set(time.Hour)
	v.Set(0)
	if logLevel.Level() != slog.LevelDebug {
		t.Errorf("level=%v，应保持 Debug", logLevel.Level())
	}
}
