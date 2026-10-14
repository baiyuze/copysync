package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRotatingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "logs", "daemon.log")
	r, err := openRotating(path, 100)
	if err != nil {
		t.Fatal(err)
	}
	line := strings.Repeat("x", 39) + "\n" // 40 字节
	for range 5 {
		if _, err := r.Write([]byte(line)); err != nil {
			t.Fatal(err)
		}
	}
	r.Close()

	cur, _ := os.ReadFile(path)
	old, _ := os.ReadFile(path + ".1")
	// 写到第 3 行时超出 100 字节，前两行转到 .1；第 5 行时再转一次
	if len(cur) != 40 || len(old) != 80 {
		t.Errorf("当前 %d 字节、上一份 %d 字节，应为 40、80", len(cur), len(old))
	}

	// 重新打开时接着已有的大小算
	r, _ = openRotating(path, 100)
	r.Write([]byte(line))
	r.Write([]byte(line))
	r.Close()
	if cur, _ := os.ReadFile(path); len(cur) != 40 {
		t.Errorf("重新打开后当前文件 %d 字节，应已轮转为 40", len(cur))
	}
}
