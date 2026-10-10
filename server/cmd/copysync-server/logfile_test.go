package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLogFileRotates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "server.log")
	l, err := openLog(path)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	line := bytes.Repeat([]byte("x"), 1<<20)
	for range 10 {
		if _, err := l.Write(line); err != nil {
			t.Fatal(err)
		}
	}
	// 第 11 MB 写不下，先轮换
	if _, err := l.Write([]byte("new\n")); err != nil {
		t.Fatal(err)
	}

	old, err := os.Stat(path + ".1")
	if err != nil {
		t.Fatal("没有轮换出 .1：", err)
	}
	if old.Size() != 10<<20 {
		t.Errorf(".1 的大小 = %d，应为 10 MB", old.Size())
	}
	cur, _ := os.ReadFile(path)
	if string(cur) != "new\n" {
		t.Errorf("新文件内容 = %q", cur)
	}
}

func TestLogFileRotatesOversizedOnOpen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.log")
	if err := os.WriteFile(path, bytes.Repeat([]byte("x"), maxLogBytes), 0o644); err != nil {
		t.Fatal(err)
	}
	l, err := openLog(path)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if l.size != 0 {
		t.Errorf("打开已满的日志后 size = %d，应先轮换", l.size)
	}
	if _, err := os.Stat(path + ".1"); err != nil {
		t.Error("没有轮换出 .1：", err)
	}
}
