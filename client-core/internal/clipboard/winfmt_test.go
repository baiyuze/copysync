package clipboard

import (
	"encoding/binary"
	"image"
	"image/color"
	"reflect"
	"strings"
	"testing"
)

func TestCRLF(t *testing.T) {
	if got := toCRLF("a\nb\r\nc"); got != "a\r\nb\r\nc" {
		t.Errorf("toCRLF = %q", got)
	}
	if got := fromCRLF("a\r\nb\nc"); got != "a\nb\nc" {
		t.Errorf("fromCRLF = %q", got)
	}
}

func TestUTF16(t *testing.T) {
	for _, s := range []string{"", "hello", "中文：复制", "emoji 😀 与组合字符 が"} {
		b := utf16Bytes(s)
		if b[len(b)-1] != 0 || b[len(b)-2] != 0 {
			t.Errorf("%q 没有以 NUL 结尾", s)
		}
		// 剪贴板的内存块常比内容大，结尾有填充
		if got := utf16String(append(b, 'x', 0, 'y', 0)); got != s {
			t.Errorf("往返 %q 得到 %q", s, got)
		}
	}
}

func TestCFHTMLRoundTrip(t *testing.T) {
	cases := []struct {
		name, in, fragment string
	}{
		{"片段", "<b>加粗</b> 文字", "<b>加粗</b> 文字"},
		{"完整文档", "<html><head><meta charset=utf-8></head><body class=x><p>段落</p></body></html>", "<p>段落</p>"},
		{"自带标记", "<html><body><!--StartFragment--><i>x</i><!--EndFragment--></body></html>", "<i>x</i>"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := encodeCFHTML(c.in)
			if b[len(b)-1] != 0 {
				t.Fatal("没有以 NUL 结尾")
			}
			s := string(b[:len(b)-1])
			off := func(key string) int {
				i := strings.Index(s, key+":")
				var n int
				for _, r := range s[i+len(key)+1 : i+len(key)+11] {
					n = n*10 + int(r-'0')
				}
				return n
			}
			// 偏移量是字节偏移，中文是多字节，正好检验
			if got := s[off("StartFragment"):off("EndFragment")]; got != c.fragment {
				t.Errorf("片段 = %q，应为 %q", got, c.fragment)
			}
			if !strings.HasPrefix(s[off("StartHTML"):], "<html") || off("EndHTML") != len(s) {
				t.Errorf("StartHTML/EndHTML 不对：%q", s)
			}
			html, err := decodeCFHTML(b)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(html, c.fragment) || !strings.HasPrefix(html, "<html") {
				t.Errorf("解析得到 %q", html)
			}
		})
	}
}

func TestDecodeCFHTMLFromWindows(t *testing.T) {
	// Chrome 写入的样子：带 SourceURL，值里有冒号
	doc := "<html>\r\n<body>\r\n<!--StartFragment--><p>你好</p><!--EndFragment-->\r\n</body>\r\n</html>"
	head := "Version:0.9\r\nStartHTML:0000000000\r\nEndHTML:0000000000\r\nStartFragment:0000000000\r\nEndFragment:0000000000\r\nSourceURL:https://example.com/a\r\n"
	start := len(head)
	head = strings.Replace(head, "StartHTML:0000000000", "StartHTML:"+pad10(start), 1)
	head = strings.Replace(head, "EndHTML:0000000000", "EndHTML:"+pad10(start+len(doc)), 1)
	fs := start + strings.Index(doc, "<p>")
	fe := start + strings.Index(doc, "<!--End")
	head = strings.Replace(head, "StartFragment:0000000000", "StartFragment:"+pad10(fs), 1)
	head = strings.Replace(head, "EndFragment:0000000000", "EndFragment:"+pad10(fe), 1)

	got, err := decodeCFHTML([]byte(head + doc + "\x00\x00"))
	if err != nil || got != doc {
		t.Fatalf("得到 %q, %v", got, err)
	}

	// StartHTML 为 -1 的程序：只给片段
	onlyFrag := strings.Replace(head, "StartHTML:"+pad10(start), "StartHTML:-000000001", 1)
	got, err = decodeCFHTML([]byte(onlyFrag + doc))
	if err != nil || got != "<p>你好</p>" {
		t.Fatalf("只有片段时得到 %q, %v", got, err)
	}

	if _, err := decodeCFHTML([]byte("<p>没有头部</p>")); err == nil {
		t.Error("没有头部时应报错")
	}
}

func pad10(n int) string {
	s := strings.Repeat("0", 10)
	d := []byte(s)
	for i := 9; i >= 0 && n > 0; i-- {
		d[i] = byte('0' + n%10)
		n /= 10
	}
	return string(d)
}

func TestDropFiles(t *testing.T) {
	paths := []string{`C:\Users\光光\Desktop\报表.xlsx`, `D:\a b\c`}
	b := encodeDropFiles(paths)
	if binary.LittleEndian.Uint32(b) != 20 || binary.LittleEndian.Uint32(b[16:]) != 1 {
		t.Fatal("DROPFILES 头部不对")
	}
	if !strings.HasSuffix(string(b), "\x00\x00\x00\x00") {
		t.Error("路径表应以两个 NUL 结尾")
	}
	got, err := decodeDropFiles(b)
	if err != nil || !reflect.DeepEqual(got, paths) {
		t.Fatalf("往返得到 %v, %v", got, err)
	}
	if _, err := decodeDropFiles(b[:10]); err == nil {
		t.Error("过短的数据应报错")
	}
}

func TestDIBRoundTrip(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 3, 2))
	src.SetNRGBA(0, 0, color.NRGBA{255, 0, 0, 255})
	src.SetNRGBA(1, 0, color.NRGBA{0, 255, 0, 128}) // 半透明
	src.SetNRGBA(2, 1, color.NRGBA{0, 0, 255, 0})   // 全透明

	b := imageToDIBV5(src)
	if len(b) != 124+3*2*4 {
		t.Fatalf("长度 = %d", len(b))
	}
	got, err := dibToImage(b)
	if err != nil {
		t.Fatal(err)
	}
	for y := range 2 {
		for x := range 3 {
			if a, g := src.NRGBAAt(x, y), got.(*image.NRGBA).NRGBAAt(x, y); a != g {
				t.Errorf("(%d,%d) = %v，应为 %v", x, y, g, a)
			}
		}
	}
}

// dib 拼一个最简的 BITMAPINFOHEADER 位图。
func dib(width, height, bitCount int, compression uint32, extra, pixels []byte) []byte {
	b := make([]byte, 40)
	le := binary.LittleEndian
	le.PutUint32(b[0:], 40)
	le.PutUint32(b[4:], uint32(int32(width)))
	le.PutUint32(b[8:], uint32(int32(height)))
	le.PutUint16(b[12:], 1)
	le.PutUint16(b[14:], uint16(bitCount))
	le.PutUint32(b[16:], compression)
	return append(append(b, extra...), pixels...)
}

func TestDIBVariants(t *testing.T) {
	t.Run("24 位，自下而上，行对齐", func(t *testing.T) {
		// 宽 1：每行 3 字节，补到 4 字节。第一行存的是图片的最后一行
		px := []byte{0, 0, 255, 0 /* 红，在底部 */, 255, 0, 0, 0 /* 蓝，在顶部 */}
		img, err := dibToImage(dib(1, 2, 24, biRGB, nil, px))
		if err != nil {
			t.Fatal(err)
		}
		n := img.(*image.NRGBA)
		if n.NRGBAAt(0, 0) != (color.NRGBA{0, 0, 255, 255}) || n.NRGBAAt(0, 1) != (color.NRGBA{255, 0, 0, 255}) {
			t.Errorf("像素 = %v %v", n.NRGBAAt(0, 0), n.NRGBAAt(0, 1))
		}
	})
	t.Run("32 位透明度全为 0 当作不透明", func(t *testing.T) {
		img, err := dibToImage(dib(1, -1, 32, biRGB, nil, []byte{10, 20, 30, 0}))
		if err != nil {
			t.Fatal(err)
		}
		if c := img.(*image.NRGBA).NRGBAAt(0, 0); c != (color.NRGBA{30, 20, 10, 255}) {
			t.Errorf("像素 = %v", c)
		}
	})
	t.Run("BI_BITFIELDS 掩码跟在 40 字节头部之后", func(t *testing.T) {
		masks := make([]byte, 12)
		binary.LittleEndian.PutUint32(masks[0:], 0x00FF0000)
		binary.LittleEndian.PutUint32(masks[4:], 0x0000FF00)
		binary.LittleEndian.PutUint32(masks[8:], 0x000000FF)
		img, err := dibToImage(dib(1, 1, 32, biBitfields, masks, []byte{1, 2, 3, 0}))
		if err != nil {
			t.Fatal(err)
		}
		if c := img.(*image.NRGBA).NRGBAAt(0, 0); c != (color.NRGBA{3, 2, 1, 255}) {
			t.Errorf("像素 = %v", c)
		}
	})
	t.Run("8 位调色板", func(t *testing.T) {
		pal := []byte{0, 0, 0, 0, 255, 255, 255, 0} // 黑、白
		le := binary.LittleEndian
		b := dib(2, 1, 8, biRGB, pal, []byte{1, 0, 0, 0}) // 白、黑，补到 4 字节
		le.PutUint32(b[32:], 2)                           // biClrUsed
		img, err := dibToImage(b)
		if err != nil {
			t.Fatal(err)
		}
		n := img.(*image.NRGBA)
		if n.NRGBAAt(0, 0) != (color.NRGBA{255, 255, 255, 255}) || n.NRGBAAt(1, 0) != (color.NRGBA{0, 0, 0, 255}) {
			t.Errorf("像素 = %v %v", n.NRGBAAt(0, 0), n.NRGBAAt(1, 0))
		}
	})
	t.Run("数据不完整", func(t *testing.T) {
		if _, err := dibToImage(dib(10, 10, 32, biRGB, nil, []byte{1, 2, 3})); err == nil {
			t.Error("应报错")
		}
		if _, err := dibToImage([]byte{1, 2, 3}); err == nil {
			t.Error("应报错")
		}
	})
}
