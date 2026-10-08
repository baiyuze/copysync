//go:build windows

package clipboard

// 在真实的 Windows 剪贴板上做读写往返。会覆盖本机剪贴板的内容，
// 需要交互式桌面（没有时自动跳过）。

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	testLoop     *MainLoop
	testLoopOnce sync.Once
)

// mainLoop 与 copysyncd 一样：一个锁定的线程创建窗口、泵消息，剪贴板操作都投递过去。
func mainLoop(t *testing.T) *MainLoop {
	t.Helper()
	testLoopOnce.Do(func() {
		l := NewMainLoop()
		ready := make(chan struct{})
		go func() {
			runtime.LockOSThread()
			InitPlatform()
			close(ready)
			l.Run(context.Background())
		}()
		<-ready
		testLoop = l
	})
	if owner == 0 {
		t.Skip("没有交互式桌面，创建不了剪贴板窗口")
	}
	return testLoop
}

func onMain[T any](t *testing.T, fn func() (T, error)) T {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	v, err := CallResult(ctx, mainLoop(t), fn)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func write(t *testing.T, c Content) Snapshot {
	t.Helper()
	cb := New()
	onMain(t, func() (struct{}, error) { return struct{}{}, cb.Write(c) })
	return onMain(t, cb.Peek)
}

func TestWindowsText(t *testing.T) {
	before := onMain(t, New().Peek).ChangeCount
	snap := write(t, Content{Kind: KindText, Text: "第一行\n第二行 😀"})
	if snap.Kind != KindText || snap.ChangeCount == before {
		t.Fatalf("Peek = %+v（写入前序号 %d）", snap, before)
	}
	got := onMain(t, New().Read)
	if got.Text != "第一行\n第二行 😀" {
		t.Errorf("读回 %q", got.Text)
	}
	// 剪贴板里存的是 CRLF，记事本才认
	raw := onMain(t, func() (string, error) {
		if err := openClipboard(); err != nil {
			return "", err
		}
		defer procCloseClipboard.Call()
		b, err := getData(cfUnicodeText)
		return utf16String(b), err
	})
	if !strings.Contains(raw, "\r\n") {
		t.Errorf("剪贴板里的原始文本 %q 没有 CRLF", raw)
	}
}

func TestWindowsHTML(t *testing.T) {
	snap := write(t, Content{Kind: KindHTML, HTML: "<p><b>粗体</b>与普通</p>", Text: "粗体与普通"})
	if snap.Kind != KindHTML {
		t.Fatalf("Kind = %v，格式 %v", snap.Kind, snap.Types)
	}
	got := onMain(t, New().Read)
	if !strings.Contains(got.HTML, "<p><b>粗体</b>与普通</p>") || got.Text != "粗体与普通" {
		t.Errorf("读回 HTML=%q Text=%q", got.HTML, got.Text)
	}
}

func TestWindowsImage(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	src.SetNRGBA(0, 0, color.NRGBA{255, 0, 0, 255})
	src.SetNRGBA(1, 1, color.NRGBA{0, 0, 255, 100})
	var buf bytes.Buffer
	if err := png.Encode(&buf, src); err != nil {
		t.Fatal(err)
	}
	snap := write(t, Content{Kind: KindImage, Image: buf.Bytes()})
	if snap.Kind != KindImage || !strings.Contains(strings.Join(snap.Types, ","), "CF_DIB") {
		t.Fatalf("Peek = %+v（系统应从 CF_DIBV5 合成 CF_DIB）", snap)
	}
	got := onMain(t, New().Read)
	img, err := png.Decode(bytes.NewReader(got.Image))
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []image.Point{{0, 0}, {1, 1}} {
		want := src.NRGBAAt(p.X, p.Y)
		if c := color.NRGBAModel.Convert(img.At(p.X, p.Y)); c != want {
			t.Errorf("%v = %v，应为 %v", p, c, want)
		}
	}
}

func TestWindowsFiles(t *testing.T) {
	dir := t.TempDir()
	var paths []string
	for _, name := range []string{"报表 2026.xlsx", "a.txt"} {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(name), 0o600); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, p)
	}
	snap := write(t, Content{Kind: KindFile, Files: paths})
	if snap.Kind != KindFile {
		t.Fatalf("Kind = %v，格式 %v", snap.Kind, snap.Types)
	}
	got := onMain(t, New().Read)
	if !reflect.DeepEqual(got.Files, paths) {
		t.Errorf("读回 %v", got.Files)
	}
}

// 剪贴板里同时有图片和文本（Excel、Word 复制的内容）时按文本处理。
func TestWindowsRichContentIsNotImage(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	snap := onMain(t, func() (Snapshot, error) {
		if err := openClipboard(); err != nil {
			return Snapshot{}, err
		}
		procEmptyClipboard.Call()
		_ = setData(formats.html, encodeCFHTML("<table><tr><td>1</td></tr></table>"))
		_ = setData(cfUnicodeText, utf16Bytes("1"))
		_ = setData(cfDIBV5, imageToDIBV5(src))
		procCloseClipboard.Call()
		return New().Peek()
	})
	if snap.Kind != KindHTML {
		t.Errorf("Kind = %v，应按带格式文本处理", snap.Kind)
	}
}

func TestWindowsSensitive(t *testing.T) {
	cases := []struct {
		name   string
		marker func() (uint32, []byte)
		want   bool
	}{
		{"无标记", nil, false},
		{"ExcludeClipboardContentFromMonitorProcessing", func() (uint32, []byte) { return formats.exclude, []byte{0} }, true},
		{"Clipboard Viewer Ignore", func() (uint32, []byte) { return formats.viewerIgnore, []byte{0} }, true},
		{"CanIncludeInClipboardHistory=0", func() (uint32, []byte) { return formats.canIncludeHistory, dwordBytes(0) }, true},
		{"CanIncludeInClipboardHistory=1", func() (uint32, []byte) { return formats.canIncludeHistory, dwordBytes(1) }, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			snap := onMain(t, func() (Snapshot, error) {
				if err := openClipboard(); err != nil {
					return Snapshot{}, err
				}
				procEmptyClipboard.Call()
				_ = setData(cfUnicodeText, utf16Bytes("密码"))
				if c.marker != nil {
					f, data := c.marker()
					if err := setData(f, data); err != nil {
						procCloseClipboard.Call()
						return Snapshot{}, err
					}
				}
				procCloseClipboard.Call()
				return New().Peek()
			})
			if snap.Sensitive != c.want {
				t.Errorf("Sensitive = %v，应为 %v", snap.Sensitive, c.want)
			}
		})
	}
}

// 剪贴板变化时系统会通知，Watcher 不必等下一轮轮询。
func TestWindowsChangeNotify(t *testing.T) {
	mainLoop(t)
	select { // 清掉之前的测试留下的信号
	case <-changeNotify:
	default:
	}
	write(t, Content{Kind: KindText, Text: "通知"})
	select {
	case <-changeNotify:
	case <-time.After(2 * time.Second):
		t.Error("写入后没有收到剪贴板变化通知")
	}
}
