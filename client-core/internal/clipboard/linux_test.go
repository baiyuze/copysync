//go:build linux

package clipboard

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// 需要 X11 的测试用 xclip 扮演「另一个程序」，在 Xvfb 里运行：
//
//	xvfb-run -a go test ./internal/clipboard/
//
// 没有 DISPLAY 或没装 xclip 时跳过；纯函数的测试照常运行。

func TestMain(m *testing.M) {
	if os.Getenv("DISPLAY") != "" {
		InitPlatform()
	}
	os.Exit(m.Run())
}

func requireX(t *testing.T) {
	t.Helper()
	if os.Getenv("DISPLAY") == "" {
		t.Skip("需要 X11 显示：用 xvfb-run 运行")
	}
	if _, err := exec.LookPath("xclip"); err != nil {
		t.Skip("需要 xclip")
	}
	if _, _, err := x.session(); err != nil {
		t.Fatalf("连接 X11: %v", err)
	}
}

// xclipIn 让 xclip 成为剪贴板的持有者，以 target 格式提供 data，并等到本机收到变化通知。
func xclipIn(t *testing.T, target string, data []byte) {
	t.Helper()
	before := x.count.Load()
	// xclip 读完输入后分出子进程常驻提供内容，所以不能接它的输出管道，否则会一直等子进程退出
	cmd := exec.Command("xclip", "-selection", "clipboard", "-t", target, "-i")
	cmd.Stdin = bytes.NewReader(data)
	if err := cmd.Run(); err != nil {
		t.Fatalf("xclip -i: %v", err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for x.count.Load() == before {
		if time.Now().After(deadline) {
			t.Fatal("没收到剪贴板变化的通知")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func xclipOut(t *testing.T, target string) []byte {
	t.Helper()
	out, err := exec.Command("xclip", "-selection", "clipboard", "-t", target, "-o").Output()
	if err != nil {
		t.Fatalf("xclip -o -t %s: %v", target, err)
	}
	return out
}

func TestLinuxReadText(t *testing.T) {
	requireX(t)
	xclipIn(t, "UTF8_STRING", []byte("你好，X11\r\n第二行"))
	snap, err := New().Peek()
	if err != nil || snap.Kind != KindText {
		t.Fatalf("Peek: %+v %v", snap, err)
	}
	c, err := New().Read()
	if err != nil || c.Kind != KindText || c.Text != "你好，X11\n第二行" {
		t.Fatalf("Read: %+v %v", c, err)
	}
}

func TestLinuxWriteText(t *testing.T) {
	requireX(t)
	before := x.count.Load()
	if err := New().Write(Content{Kind: KindText, Text: "从 CopySync 写入"}); err != nil {
		t.Fatal(err)
	}
	snap, err := New().Peek()
	if err != nil || snap.ChangeCount <= before || snap.Kind != KindText {
		t.Fatalf("写入后 Peek 应当看到新的计数: %+v %v（之前 %d）", snap, err, before)
	}
	if got := string(xclipOut(t, "UTF8_STRING")); got != "从 CopySync 写入" {
		t.Fatalf("别的程序读到 %q", got)
	}
	targets := string(xclipOut(t, "TARGETS"))
	for _, want := range []string{"UTF8_STRING", "text/plain;charset=utf-8", "TIMESTAMP"} {
		if !strings.Contains(targets, want) {
			t.Errorf("TARGETS 里缺少 %s：%q", want, targets)
		}
	}
	// 别人复制了新内容：本机不再持有，读到的是别人的
	xclipIn(t, "UTF8_STRING", []byte("别人复制的"))
	c, _ := New().Read()
	if c.Text != "别人复制的" {
		t.Fatalf("应当读到别人的内容，得到 %+v", c)
	}
}

// noisyPNG 生成一张压缩不下去的大图，超过 256 KB，读写都要走 INCR 分段。
func noisyPNG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	r := rand.New(rand.NewSource(1))
	for i := range img.Pix {
		img.Pix[i] = byte(r.Intn(256))
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestLinuxLargeImageBothWays(t *testing.T) {
	requireX(t)
	big := noisyPNG(t, 600, 600)
	if len(big) <= maxChunk {
		t.Fatalf("测试图片太小（%d 字节），测不到 INCR", len(big))
	}

	if err := New().Write(Content{Kind: KindImage, Image: big}); err != nil {
		t.Fatal(err)
	}
	if got := xclipOut(t, "image/png"); !bytes.Equal(got, big) {
		t.Fatalf("别的程序读到的图片不一致：%d 字节，应为 %d", len(got), len(big))
	}

	other := noisyPNG(t, 500, 700)
	xclipIn(t, "image/png", other)
	snap, _ := New().Peek()
	c, err := New().Read()
	if err != nil || snap.Kind != KindImage || !bytes.Equal(c.Image, other) {
		t.Fatalf("读到的图片不一致：%v %v，%d 字节，应为 %d", snap.Kind, err, len(c.Image), len(other))
	}
}

func TestLinuxFiles(t *testing.T) {
	requireX(t)
	dir := t.TempDir()
	a, b := filepath.Join(dir, "报告 1.txt"), filepath.Join(dir, "b")
	os.WriteFile(a, []byte("a"), 0o644)
	os.Mkdir(b, 0o755)

	if err := New().Write(Content{Kind: KindFile, Files: []string{a, b}}); err != nil {
		t.Fatal(err)
	}
	gnome := string(xclipOut(t, "x-special/gnome-copied-files"))
	if !strings.HasPrefix(gnome, "copy\n") || !strings.Contains(gnome, "%E6%8A%A5%E5%91%8A%201.txt") {
		t.Fatalf("文件管理器要的格式不对：%q", gnome)
	}
	if got := parseURIs(string(xclipOut(t, "text/uri-list"))); len(got) != 2 || got[0] != a || got[1] != b {
		t.Fatalf("text/uri-list 解出 %v", got)
	}

	xclipIn(t, "text/uri-list", []byte("file://"+strings.ReplaceAll(a, " ", "%20")+"\r\n"))
	snap, _ := New().Peek()
	c, err := New().Read()
	if err != nil || snap.Kind != KindFile || c.Kind != KindFile || len(c.Files) != 1 || c.Files[0] != a {
		t.Fatalf("读文件：%+v %+v %v", snap, c, err)
	}

	// 网址不是本地文件：当作文本
	xclipIn(t, "text/uri-list", []byte("https://example.com/a\r\n"))
	c, _ = New().Read()
	if c.Kind != KindText || c.Text != "https://example.com/a" {
		t.Fatalf("网址应当按文本处理：%+v", c)
	}
}

func TestLinuxHTML(t *testing.T) {
	requireX(t)
	xclipIn(t, "text/html", []byte("<b>粗体</b>"))
	snap, _ := New().Peek()
	c, err := New().Read()
	if err != nil || snap.Kind != KindHTML || c.HTML != "<b>粗体</b>" {
		t.Fatalf("HTML：%+v %+v %v", snap, c, err)
	}
}

func TestLinuxPasswordHintIsSensitive(t *testing.T) {
	requireX(t)
	xclipIn(t, "x-kde-passwordManagerHint", []byte("secret"))
	if snap, _ := New().Peek(); !snap.Sensitive {
		t.Fatalf("密码管理器复制的内容应当标为敏感：%+v", snap)
	}
}

// ─────────────────────────── 不需要 X11 ───────────────────────────

func TestClassify(t *testing.T) {
	cases := []struct {
		targets []string
		want    Kind
	}{
		{[]string{"x-special/gnome-copied-files", "text/uri-list", "UTF8_STRING"}, KindFile},
		{[]string{"text/uri-list", "text/plain"}, KindFile},
		{[]string{"text/uri-list", "text/html", "image/png"}, KindImage}, // 浏览器里「复制图片」
		{[]string{"image/png"}, KindImage},
		{[]string{"image/png", "text/html", "UTF8_STRING"}, KindHTML}, // 办公软件附带的渲染图不算
		{[]string{"text/html", "UTF8_STRING"}, KindHTML},
		{[]string{"UTF8_STRING", "STRING"}, KindText},
		{[]string{"TARGETS"}, KindUnknown},
	}
	for _, c := range cases {
		if got := classify(c.targets).Kind; got != c.want {
			t.Errorf("%v：得到 %v，应为 %v", c.targets, got, c.want)
		}
	}
}

func TestDecodeHTMLUTF16(t *testing.T) {
	b := []byte{0xFF, 0xFE}
	for _, r := range "<i>斜体</i>" {
		b = append(b, byte(r), byte(r>>8))
	}
	if got := decodeHTML(append(b, 0, 0)); got != "<i>斜体</i>" {
		t.Fatalf("得到 %q", got)
	}
	if got := decodeText(tString, []byte{'c', 'a', 'f', 0xE9}); got != "café" {
		t.Fatalf("Latin-1 得到 %q", got)
	}
}

func TestToPNGFromBMP(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 2))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	dib := imageToDIBV5(img)
	bmp := append([]byte{'B', 'M', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, dib...)
	out, err := toPNG("image/bmp", bmp)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := png.Decode(bytes.NewReader(out))
	if err != nil || decoded.Bounds().Dx() != 3 || decoded.Bounds().Dy() != 2 {
		t.Fatalf("BMP 转 PNG：%v %v", decoded, err)
	}
}

func TestParseURIs(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "a b.txt")
	os.WriteFile(p, nil, 0o644)
	in := "copy\n# 注释\nfile://" + strings.ReplaceAll(p, " ", "%20") + "\nfile:///不存在\nhttps://example.com\n"
	if got := parseURIs(in); len(got) != 1 || got[0] != p {
		t.Fatalf("得到 %v", got)
	}
}

// Wayland 桌面上，后台服务经 XWayland 的 X11 选区读写剪贴板，靠混成器在两边同步。
// 这里用 wl-copy、wl-paste 扮演原生 Wayland 程序，验证两个方向都通。只在 Wayland 会话里运行：
//
//	WAYLAND_DISPLAY=wayland-0 DISPLAY=:0 XAUTHORITY=… ./clipboard.test -test.run Wayland
func TestLinuxWaylandBridge(t *testing.T) {
	requireX(t)
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		t.Skip("需要 Wayland 会话")
	}
	if _, err := exec.LookPath("wl-copy"); err != nil {
		t.Skip("需要 wl-clipboard")
	}
	waitCount := func(before int64) {
		t.Helper()
		deadline := time.Now().Add(5 * time.Second)
		for x.count.Load() == before {
			if time.Now().After(deadline) {
				t.Fatal("Wayland 程序复制后，X11 这边没收到变化通知")
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	wlCopy := func(mime string, data []byte) {
		t.Helper()
		before := x.count.Load()
		cmd := exec.Command("wl-copy", "--type", mime)
		cmd.Stdin = bytes.NewReader(data)
		if err := cmd.Run(); err != nil {
			t.Fatalf("wl-copy: %v", err)
		}
		waitCount(before)
	}
	wlPaste := func(mime string) []byte {
		t.Helper()
		out, err := exec.Command("wl-paste", "--no-newline", "--type", mime).Output()
		if err != nil {
			t.Fatalf("wl-paste --type %s: %v", mime, err)
		}
		return out
	}

	// Wayland → X11
	wlCopy("text/plain;charset=utf-8", []byte("从 Wayland 复制"))
	if c, err := New().Read(); err != nil || c.Text != "从 Wayland 复制" {
		t.Fatalf("读 Wayland 程序复制的文本：%+v %v", c, err)
	}
	img := noisyPNG(t, 400, 400)
	wlCopy("image/png", img)
	if c, err := New().Read(); err != nil || c.Kind != KindImage || !bytes.Equal(c.Image, img) {
		t.Fatalf("读 Wayland 程序复制的图片：%v %v，%d 字节，应为 %d", c.Kind, err, len(c.Image), len(img))
	}

	// X11 → Wayland
	if err := New().Write(Content{Kind: KindText, Text: "从 CopySync 写入，Wayland 程序来粘贴"}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond) // 混成器接管选区需要一点时间
	if got := string(wlPaste("text/plain;charset=utf-8")); got != "从 CopySync 写入，Wayland 程序来粘贴" {
		t.Fatalf("Wayland 程序粘贴到 %q", got)
	}
	if err := New().Write(Content{Kind: KindImage, Image: img}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	if got := wlPaste("image/png"); !bytes.Equal(got, img) {
		t.Fatalf("Wayland 程序粘贴到的图片不一致：%d 字节，应为 %d", len(got), len(img))
	}
}
