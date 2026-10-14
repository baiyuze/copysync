package main

import (
	"os"
	"path/filepath"
	"sync"
)

// rotatingFile 是按大小轮转的日志文件：超过上限就把当前文件改名为 .1（覆盖更早的），
// 再从头写一个新的。最多占用两倍上限的空间。
//
// Mac 上日志由 launchd 重定向到文件，用不到它；Windows 上后台服务没有控制台，
// 日志只能自己写。
type rotatingFile struct {
	mu   sync.Mutex
	path string
	max  int64
	f    *os.File
	size int64
}

func openRotating(path string, max int64) (*rotatingFile, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	r := &rotatingFile{path: path, max: max}
	if err := r.open(); err != nil {
		return nil, err
	}
	if r.size > max {
		r.rotate()
	}
	return r, nil
}

func (r *rotatingFile) open() error {
	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	r.f = f
	r.size = 0
	if info, err := f.Stat(); err == nil {
		r.size = info.Size()
	}
	return nil
}

func (r *rotatingFile) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.size+int64(len(p)) > r.max {
		r.rotate()
	}
	n, err := r.f.Write(p)
	r.size += int64(n)
	return n, err
}

// rotate 换一个新文件。改名失败（比如正被别的程序以独占方式打开）就接着往原文件写，
// 下次再试：日志写不进去比日志过大更糟。
func (r *rotatingFile) rotate() {
	r.f.Close()
	_ = os.Remove(r.path + ".1")
	_ = os.Rename(r.path, r.path+".1")
	if err := r.open(); err != nil {
		// 连新文件都开不了，只能丢弃后续日志
		r.f, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	}
}

func (r *rotatingFile) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.f.Close()
}
