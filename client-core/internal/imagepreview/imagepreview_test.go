package imagepreview

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/baiyuze/copysync/client-core/internal/store"
)

func TestPreviewableAndCount(t *testing.T) {
	for name, want := range map[string]bool{
		"a.png": true, "B.JPG": true, "c.jpeg": true, "d.gif": true, "e.webp": true, "f.bmp": true,
		"notes.txt": false, "deck.key": false, "png": false, "": false,
		"IMG_2041.HEIC": runtime.GOOS == "darwin", "scan.tiff": runtime.GOOS == "darwin",
	} {
		if got := Previewable(name); got != want {
			t.Errorf("Previewable(%q) = %v，应为 %v", name, got, want)
		}
	}
	clip := store.Clip{Kind: store.KindFile, Items: []store.Item{
		{Name: "a.png"}, {Name: "b.jpg"}, {Name: "notes.txt"}, {Name: "photos.png", IsDir: true},
	}}
	if got := ImageCount(clip); got != 2 {
		t.Errorf("文件记录的图片数 = %d，应为 2（文件夹不算）", got)
	}
	if got := ImageCount(store.Clip{Kind: store.KindImage}); got != 1 {
		t.Errorf("剪贴板图片 = %d，应为 1", got)
	}
	if got := ImageCount(store.Clip{Kind: store.KindText}); got != 0 {
		t.Errorf("文本 = %d，应为 0", got)
	}
}

func writePNG(t *testing.T, path string, w, h int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		img.Set(x, 0, color.RGBA{R: 255, A: 255})
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func TestPrepareKeepsNativeFormats(t *testing.T) {
	path := filepath.Join(t.TempDir(), "原图.png")
	writePNG(t, path, 4, 3)
	got, err := Prepare(path)
	if err != nil || got != path {
		t.Fatalf("PNG 应原样返回：%q, %v", got, err)
	}
	if _, err := Prepare(filepath.Join(t.TempDir(), "a.txt")); err == nil {
		t.Fatal("不是图片也返回了路径")
	}
}

// HEIC 与 TIFF 在 Mac 上转成 PNG：大图缩到 MaxPixels，小图不放大。用系统自带的 sips 造测试图片。
func TestPrepareConvertsOnMac(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("只有 Mac 转换 HEIC、TIFF")
	}
	dir := t.TempDir()
	for _, tc := range []struct {
		format, ext string
		w, h        int
		wantW       int
	}{
		{"heic", ".heic", 3000, 2000, MaxPixels},
		{"tiff", ".tiff", 40, 30, 40},
	} {
		src := filepath.Join(dir, "src-"+tc.format+".png")
		writePNG(t, src, tc.w, tc.h)
		in := filepath.Join(dir, "照片"+tc.ext)
		if out, err := exec.Command("sips", "-s", "format", tc.format, src, "--out", in).CombinedOutput(); err != nil {
			t.Skipf("sips 无法生成 %s：%v %s", tc.format, err, out)
		}
		got, err := Prepare(in)
		if err != nil {
			t.Fatalf("%s 转换失败：%v", tc.format, err)
		}
		f, err := os.Open(got)
		if err != nil {
			t.Fatal(err)
		}
		cfg, err := png.DecodeConfig(f)
		f.Close()
		if err != nil {
			t.Fatalf("%s 转出来的不是 PNG：%v", tc.format, err)
		}
		if cfg.Width != tc.wantW {
			t.Errorf("%s 转换后宽 %d，应为 %d", tc.format, cfg.Width, tc.wantW)
		}
		// 第二次直接用上次转好的文件
		again, err := Prepare(in)
		if err != nil || again != got {
			t.Errorf("第二次应复用 %q，得到 %q, %v", got, again, err)
		}
	}
}
