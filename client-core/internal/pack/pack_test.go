package pack

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestPackUnpackSingleFile(t *testing.T) {
	src := t.TempDir()
	file := filepath.Join(src, "hello.txt")
	writeFile(t, file, "hello copysync")

	var buf bytes.Buffer
	if err := Pack(&buf, []string{file}); err != nil {
		t.Fatalf("打包: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("打包结果为空")
	}

	dst := t.TempDir()
	tops, err := Unpack(&buf, dst)
	if err != nil {
		t.Fatalf("解包: %v", err)
	}

	// 顶层条目是写入剪贴板时要用的路径，缺了它对端就无从粘贴
	if len(tops) != 1 {
		t.Fatalf("顶层条目数 = %d，期望 1；实际 %v", len(tops), tops)
	}
	if filepath.Base(tops[0]) != "hello.txt" {
		t.Errorf("顶层条目 = %q，期望 hello.txt", tops[0])
	}

	got, err := os.ReadFile(filepath.Join(dst, "hello.txt"))
	if err != nil {
		t.Fatalf("读取解包结果: %v", err)
	}
	if string(got) != "hello copysync" {
		t.Errorf("内容 = %q，期望 %q", got, "hello copysync")
	}
}

func TestPackUnpackDirectory(t *testing.T) {
	src := t.TempDir()
	root := filepath.Join(src, "docs")
	writeFile(t, filepath.Join(root, "a.txt"), "AAA")
	writeFile(t, filepath.Join(root, "sub", "b.txt"), "BBB")

	var buf bytes.Buffer
	if err := Pack(&buf, []string{root}); err != nil {
		t.Fatalf("打包: %v", err)
	}

	dst := t.TempDir()
	tops, err := Unpack(&buf, dst)
	if err != nil {
		t.Fatalf("解包: %v", err)
	}
	if len(tops) != 1 || filepath.Base(tops[0]) != "docs" {
		t.Fatalf("顶层条目 = %v，期望只有 docs", tops)
	}

	// 目录结构要原样还原，用户复制的是什么，粘贴出来就该是什么
	for path, want := range map[string]string{
		filepath.Join(dst, "docs", "a.txt"):        "AAA",
		filepath.Join(dst, "docs", "sub", "b.txt"): "BBB",
	} {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("缺少 %s: %v", path, err)
			continue
		}
		if string(got) != want {
			t.Errorf("%s 内容 = %q，期望 %q", path, got, want)
		}
	}
}

func TestPackMultipleTopLevelEntries(t *testing.T) {
	src := t.TempDir()
	f1 := filepath.Join(src, "one.txt")
	f2 := filepath.Join(src, "two.txt")
	dir := filepath.Join(src, "folder")
	writeFile(t, f1, "1")
	writeFile(t, f2, "2")
	writeFile(t, filepath.Join(dir, "inner.txt"), "3")

	var buf bytes.Buffer
	if err := Pack(&buf, []string{f1, f2, dir}); err != nil {
		t.Fatalf("打包: %v", err)
	}

	dst := t.TempDir()
	tops, err := Unpack(&buf, dst)
	if err != nil {
		t.Fatalf("解包: %v", err)
	}

	names := make([]string, 0, len(tops))
	for _, p := range tops {
		names = append(names, filepath.Base(p))
	}
	sort.Strings(names)
	want := []string{"folder", "one.txt", "two.txt"}
	if len(names) != 3 || names[0] != want[0] || names[1] != want[1] || names[2] != want[2] {
		t.Errorf("顶层条目 = %v，期望 %v", names, want)
	}
}

func TestPackPreservesEmptyDirectory(t *testing.T) {
	src := t.TempDir()
	root := filepath.Join(src, "empty-root")
	if err := os.MkdirAll(filepath.Join(root, "nothing"), 0o700); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := Pack(&buf, []string{root}); err != nil {
		t.Fatalf("打包: %v", err)
	}
	dst := t.TempDir()
	if _, err := Unpack(&buf, dst); err != nil {
		t.Fatalf("解包: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "empty-root", "nothing")); err != nil {
		t.Errorf("空目录未被还原: %v", err)
	}
}

// 归档来自对端设备，不能因为已配对就无条件信任其中的路径。
func TestUnpackRejectsPathEscape(t *testing.T) {
	for _, name := range []string{"../evil.txt", "/etc/passwd", "a/../../evil"} {
		if _, err := safeJoin("/tmp/base", name); err == nil {
			t.Errorf("路径 %q 应被拒绝", name)
		}
	}
	// 正常路径仍要放行
	if _, err := safeJoin("/tmp/base", "docs/a.txt"); err != nil {
		t.Errorf("正常路径被误拒: %v", err)
	}
}

func TestScanCountsBytes(t *testing.T) {
	src := t.TempDir()
	writeFile(t, filepath.Join(src, "a.txt"), "12345")
	writeFile(t, filepath.Join(src, "sub", "b.txt"), "678")

	s, err := Scan([]string{src}, 0)
	if err != nil {
		t.Fatalf("扫描: %v", err)
	}
	if s.Files != 2 {
		t.Errorf("文件数 = %d，期望 2", s.Files)
	}
	if s.Bytes != 8 {
		t.Errorf("字节数 = %d，期望 8", s.Bytes)
	}
}

// 超大目录全量扫描会让复制动作明显卡顿，因此要能在达到上限时及时收手。
func TestScanRespectsLimit(t *testing.T) {
	src := t.TempDir()
	for i := range 50 {
		writeFile(t, filepath.Join(src, string(rune('a'+i%26))+string(rune('0'+i/26))+".txt"), "x")
	}
	s, err := Scan([]string{src}, 10)
	if err != nil {
		t.Fatalf("扫描: %v", err)
	}
	if s.Files+s.Dirs > 12 { // 允许少量越界：上限在遍历回调里判定
		t.Errorf("扫描条目 = %d，未在上限附近停止", s.Files+s.Dirs)
	}
}

func TestTopLevel(t *testing.T) {
	cases := map[string]string{
		"hello.txt":       "hello.txt",
		"docs/":           "docs",
		"docs/a.txt":      "docs",
		"docs/sub/b.txt":  "docs",
		"./hello.txt":     "hello.txt",
		".":               "",
	}
	for in, want := range cases {
		if got := topLevel(in); got != want {
			t.Errorf("topLevel(%q) = %q，期望 %q", in, got, want)
		}
	}
}
