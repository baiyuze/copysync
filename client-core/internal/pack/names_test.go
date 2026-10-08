package pack

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/klauspost/compress/zstd"
)

func TestWindowsSafeName(t *testing.T) {
	cases := map[string]string{
		"会议:纪要.txt":       "会议：纪要.txt",
		`a*b?c"d<e>f|g\h`: "a＊b？c＂d＜e＞f｜g＼h",
		"notes.":          "notes._",
		"trailing ":       "trailing _",
		"CON":             "_CON",
		"nul.txt":         "_nul.txt",
		"com1.tar.gz":     "_com1.tar.gz",
		"console.txt":     "console.txt", // 只是以保留名开头，不算
		"が.txt":          "が.txt",       // か + 浊点 → が
		"中文文件名.docx":      "中文文件名.docx",
		"tab\tname":       "tab_name",
	}
	for in, want := range cases {
		if got := windowsSafeName(in); got != want {
			t.Errorf("windowsSafeName(%q) = %q，应为 %q", in, got, want)
		}
	}
}

func TestDestNames(t *testing.T) {
	d := newDestNames()
	steps := []struct {
		name  string
		isDir bool
		want  string
	}{
		{"项目/", true, "项目"},
		{"项目/Readme.md", false, "项目/Readme.md"},
		{"项目/README.md", false, "项目/README (2).md"}, // 只差大小写的文件改名
		{"项目/readme.md", false, "项目/readme (3).md"},
		{"项目/a:b/c.txt", false, "项目/a：b/c.txt"}, // 归档里没有单独的目录条目
		{"项目/A：B/", true, "项目/a：b"},             // 与上面的目录只差大小写：合并，沿用先占用的名字
		{"Docs/x", false, "Docs/x"},
		{"docs/y", false, "Docs/y"}, // 父目录只差大小写，落到同一个目录里
	}
	for _, s := range steps {
		got, err := d.rel(s.name, s.isDir)
		if err != nil || got != s.want {
			t.Errorf("rel(%q) = %q, %v，应为 %q", s.name, got, err, s.want)
		}
	}
	for _, bad := range []string{"../x", "/etc/passwd", "..", "a/../../b"} {
		if _, err := d.rel(bad, false); err == nil {
			t.Errorf("rel(%q) 应拒绝", bad)
		}
	}
}

// writeTar 手工拼一个归档：Mac 的文件系统多半不区分大小写，造不出只差大小写的两个文件。
func writeTar(t *testing.T, entries []tar.Header, bodies map[string]string) io.Reader {
	t.Helper()
	var buf bytes.Buffer
	enc, _ := zstd.NewWriter(&buf)
	tw := tar.NewWriter(enc)
	for _, h := range entries {
		body := bodies[h.Name]
		h.Size = int64(len(body))
		if h.Mode == 0 {
			h.Mode = 0o644
		}
		if err := tw.WriteHeader(&h); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	enc.Close()
	return &buf
}

func TestUnpackForWindows(t *testing.T) {
	old := windowsNames
	windowsNames = true
	t.Cleanup(func() { windowsNames = old })

	r := writeTar(t, []tar.Header{
		{Name: "会议:纪要/", Typeflag: tar.TypeDir, Mode: 0o755},
		{Name: "会议:纪要/Readme.md", Typeflag: tar.TypeReg},
		{Name: "会议:纪要/README.md", Typeflag: tar.TypeReg},
		{Name: "会议:纪要/.DS_Store", Typeflag: tar.TypeReg},
		{Name: "会议:纪要/._Readme.md", Typeflag: tar.TypeReg},
		// 上级目录不存在，符号链接一定建不成：模拟 Windows 上没有权限
		{Name: "会议:纪要/sub/link", Typeflag: tar.TypeSymlink, Linkname: "x"},
		{Name: "CON", Typeflag: tar.TypeReg},
	}, map[string]string{
		"会议:纪要/Readme.md": "1", "会议:纪要/README.md": "2", "会议:纪要/.DS_Store": "x",
		"会议:纪要/._Readme.md": "x", "CON": "3",
	})

	dest := t.TempDir()
	tops, err := Unpack(r, dest)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(tops)
	want := []string{filepath.Join(dest, "_CON"), filepath.Join(dest, "会议：纪要")}
	sort.Strings(want)
	if len(tops) != 2 || tops[0] != want[0] || tops[1] != want[1] {
		t.Errorf("顶层条目 = %v，应为 %v", tops, want)
	}

	entries, _ := os.ReadDir(filepath.Join(dest, "会议：纪要"))
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	wantNames := []string{"README (2).md", "Readme.md"}
	if len(names) != 2 || names[0] != wantNames[0] || names[1] != wantNames[1] {
		t.Errorf("目录内容 = %v，应为 %v（.DS_Store、._ 文件不落盘，建不了的符号链接跳过）", names, wantNames)
	}
	for name, body := range map[string]string{"Readme.md": "1", "README (2).md": "2"} {
		b, _ := os.ReadFile(filepath.Join(dest, "会议：纪要", name))
		if string(b) != body {
			t.Errorf("%s 内容 = %q，应为 %q（不能互相覆盖）", name, b, body)
		}
	}
}

// 打不开的文件要整个跳过，不能写了头却没有内容，把整个归档弄坏。
func TestPackSkipsUnreadableFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows 上权限位挡不住读取，造不出打不开的文件")
	}
	if os.Geteuid() == 0 {
		t.Skip("root 能读任何文件")
	}
	dir := t.TempDir()
	root := filepath.Join(dir, "d")
	os.MkdirAll(root, 0o755)
	os.WriteFile(filepath.Join(root, "ok.txt"), []byte("ok"), 0o644)
	locked := filepath.Join(root, "locked.txt")
	os.WriteFile(locked, []byte("secret"), 0o000)

	var buf bytes.Buffer
	if err := Pack(&buf, []string{root}); err != nil {
		t.Fatal(err)
	}
	dest := t.TempDir()
	if _, err := Unpack(&buf, dest); err != nil {
		t.Fatalf("归档损坏：%v", err)
	}
	if b, err := os.ReadFile(filepath.Join(dest, "d", "ok.txt")); err != nil || string(b) != "ok" {
		t.Errorf("能读的文件没传过来：%q, %v", b, err)
	}
	if _, err := os.Stat(filepath.Join(dest, "d", "locked.txt")); err == nil {
		t.Error("打不开的文件不该出现")
	}
}
