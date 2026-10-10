package main

import (
	"os"
	"path/filepath"
	"sync"
)

// 日志文件超过这个大小就换新，旧的改名为 .1：作为系统服务常年运行，不能把磁盘写满。
const maxLogBytes = 10 << 20

// logFile 是按大小轮换的日志文件，只保留上一份。
type logFile struct {
	mu   sync.Mutex
	path string
	f    *os.File
	size int64
}

func openLog(path string) (*logFile, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	l := &logFile{path: path}
	if err := l.open(); err != nil {
		return nil, err
	}
	if l.size >= maxLogBytes {
		if err := l.rotate(); err != nil {
			l.f.Close()
			return nil, err
		}
	}
	return l, nil
}

func (l *logFile) open() error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	l.f, l.size = f, fi.Size()
	return nil
}

// rotate 把当前文件改名为 .1 并重新打开。Windows 上不能改名打开着的文件，所以先关。
func (l *logFile) rotate() error {
	l.f.Close()
	if err := os.Rename(l.path, l.path+".1"); err != nil && !os.IsNotExist(err) {
		return err
	}
	return l.open()
}

func (l *logFile) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.size+int64(len(p)) > maxLogBytes {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := l.f.Write(p)
	l.size += int64(n)
	return n, err
}

func (l *logFile) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.f.Close()
}
