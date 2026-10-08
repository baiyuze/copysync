package clipboard

// Windows 剪贴板格式与 CopySync 内部表示之间的转换。
//
// 都是纯函数，不依赖 Win32 API，放在没有构建标签的文件里：任何平台都能编译和测试，
// 在 Mac 上开发时就能验证。实际调用在 windows.go。

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"
	"unicode/utf16"
)

// ─────────────────────────── 文本 ───────────────────────────

// 网络上传输的文本一律用 LF；Windows 程序（记事本等）要的是 CRLF。

func toCRLF(s string) string {
	return strings.ReplaceAll(fromCRLF(s), "\n", "\r\n")
}

func fromCRLF(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// utf16Bytes 把字符串编码成以 NUL 结尾的 UTF-16LE，即 CF_UNICODETEXT 的内存布局。
func utf16Bytes(s string) []byte {
	u := utf16.Encode([]rune(s))
	b := make([]byte, 2*len(u)+2)
	for i, v := range u {
		binary.LittleEndian.PutUint16(b[2*i:], v)
	}
	return b
}

// utf16String 解析 UTF-16LE，遇到 NUL 截止。剪贴板里的内存块常比内容大，结尾是填充。
func utf16String(b []byte) string {
	u := make([]uint16, 0, len(b)/2)
	for i := 0; i+1 < len(b); i += 2 {
		v := binary.LittleEndian.Uint16(b[i:])
		if v == 0 {
			break
		}
		u = append(u, v)
	}
	return string(utf16.Decode(u))
}

// ─────────────────────────── HTML ───────────────────────────

// CF_HTML（注册名 "HTML Format"）是 UTF-8 的 HTML，前面加一段文本头部，
// 用字节偏移量标出整个 HTML 与其中「片段」的位置：
//
//	Version:0.9
//	StartHTML:0000000105
//	EndHTML:0000000199
//	StartFragment:0000000141
//	EndFragment:0000000163
//	<html><body><!--StartFragment-->…<!--EndFragment--></body></html>
//
// 见 https://learn.microsoft.com/windows/win32/dataxchg/html-clipboard-format

const (
	fragmentStart = "<!--StartFragment-->"
	fragmentEnd   = "<!--EndFragment-->"
)

// encodeCFHTML 生成 CF_HTML。html 可以是完整文档，也可以只是片段。
//
// 片段的范围：html 里已有 StartFragment/EndFragment 注释（多半是从 Windows 复制来的）
// 就沿用；没有的话，完整文档取 <body> 的内容，否则整段都是片段。
func encodeCFHTML(html string) []byte {
	doc := html
	if !strings.Contains(doc, fragmentStart) || !strings.Contains(doc, fragmentEnd) {
		lower := strings.ToLower(doc)
		bodyOpen := strings.Index(lower, "<body")
		bodyClose := strings.LastIndex(lower, "</body>")
		if bodyOpen >= 0 && bodyClose > bodyOpen {
			if end := strings.IndexByte(doc[bodyOpen:], '>'); end >= 0 {
				inner := bodyOpen + end + 1
				if inner <= bodyClose {
					doc = doc[:inner] + fragmentStart + doc[inner:bodyClose] + fragmentEnd + doc[bodyClose:]
				}
			}
		} else {
			doc = "<html><body>" + fragmentStart + doc + fragmentEnd + "</body></html>"
		}
	}

	// 头部长度固定：偏移量都写成 10 位数，先算好再填
	const header = "Version:0.9\r\nStartHTML:%010d\r\nEndHTML:%010d\r\n" +
		"StartFragment:%010d\r\nEndFragment:%010d\r\n"
	headerLen := len(fmt.Sprintf(header, 0, 0, 0, 0))
	startHTML := headerLen
	endHTML := startHTML + len(doc)
	startFrag := startHTML + strings.Index(doc, fragmentStart) + len(fragmentStart)
	endFrag := startHTML + strings.Index(doc, fragmentEnd)

	var b bytes.Buffer
	fmt.Fprintf(&b, header, startHTML, endHTML, startFrag, endFrag)
	b.WriteString(doc)
	b.WriteByte(0)
	return b.Bytes()
}

// decodeCFHTML 取出 CF_HTML 里的 HTML。
//
// 返回 StartHTML 到 EndHTML 之间的完整文档，而不只是片段：片段常常依赖外层的
// <style>、表格结构（从 Excel 复制时片段只是几个 <tr>），单独拿出来格式会丢。
func decodeCFHTML(b []byte) (string, error) {
	if i := bytes.IndexByte(b, 0); i >= 0 {
		b = b[:i]
	}
	fields := map[string]int{}
	rest := b
	for {
		line, after, ok := bytes.Cut(rest, []byte("\n"))
		key, val, isField := strings.Cut(strings.TrimRight(string(line), "\r"), ":")
		if !ok || !isField || strings.ContainsAny(key, "<> ") {
			break
		}
		if n, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			fields[key] = n
		}
		rest = after
	}

	start, okStart := fields["StartHTML"]
	end, okEnd := fields["EndHTML"]
	// 有的程序把 StartHTML 写成 -1（只有片段），或者偏移量越界：退回片段
	if !okStart || !okEnd || start < 0 || end > len(b) || start >= end {
		fs, okFS := fields["StartFragment"]
		fe, okFE := fields["EndFragment"]
		if !okFS || !okFE || fs < 0 || fe > len(b) || fs >= fe {
			return "", errors.New("CF_HTML 头部缺少有效的偏移量")
		}
		return string(b[fs:fe]), nil
	}
	return string(b[start:end]), nil
}

// ─────────────────────────── 文件列表 ───────────────────────────

// CF_HDROP 的内存布局：DROPFILES 结构（20 字节），紧跟以 NUL 分隔、
// 以两个 NUL 结尾的路径表。
//
//	typedef struct { DWORD pFiles; POINT pt; BOOL fNC; BOOL fWide; } DROPFILES;
const dropFilesHeader = 20

func encodeDropFiles(paths []string) []byte {
	var b bytes.Buffer
	var hdr [dropFilesHeader]byte
	binary.LittleEndian.PutUint32(hdr[0:], dropFilesHeader) // pFiles：路径表的偏移
	binary.LittleEndian.PutUint32(hdr[16:], 1)              // fWide：UTF-16
	b.Write(hdr[:])
	for _, p := range paths {
		b.Write(utf16Bytes(p))
	}
	b.Write([]byte{0, 0}) // 表尾再一个 NUL
	return b.Bytes()
}

func decodeDropFiles(b []byte) ([]string, error) {
	if len(b) < dropFilesHeader {
		return nil, errors.New("CF_HDROP 数据过短")
	}
	off := int(binary.LittleEndian.Uint32(b[0:]))
	wide := binary.LittleEndian.Uint32(b[16:]) != 0
	if off < dropFilesHeader || off > len(b) {
		return nil, errors.New("CF_HDROP 偏移量无效")
	}
	if !wide {
		return nil, errors.New("不支持 ANSI 编码的文件列表")
	}
	var paths []string
	rest := b[off:]
	for len(rest) >= 2 {
		s := utf16String(rest)
		if s == "" {
			break
		}
		paths = append(paths, s)
		rest = rest[2*len(utf16.Encode([]rune(s)))+2:]
	}
	return paths, nil
}

// dwordBytes 是 "Preferred DropEffect"、"CanUploadToCloudClipboard" 这类格式的内容：一个 DWORD。
func dwordBytes(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

// ─────────────────────────── 图片 ───────────────────────────

// Windows 的位图格式 CF_DIB / CF_DIBV5：BITMAPINFOHEADER（40 字节）或
// BITMAPV5HEADER（124 字节），之后是调色板或颜色掩码，最后是像素。
// 行按 4 字节对齐；高度为正时自下而上存放。

const (
	biRGB       = 0
	biBitfields = 3
	lcsSRGB     = 0x73524742 // 'sRGB'
)

// dibToImage 解析 CF_DIB / CF_DIBV5。支持 32、24 位，以及 1、4、8 位调色板。
func dibToImage(b []byte) (image.Image, error) {
	if len(b) < 40 {
		return nil, errors.New("位图数据过短")
	}
	le := binary.LittleEndian
	hdrSize := int(le.Uint32(b[0:]))
	width := int(int32(le.Uint32(b[4:])))
	height := int(int32(le.Uint32(b[8:])))
	bitCount := int(le.Uint16(b[14:]))
	compression := le.Uint32(b[16:])
	clrUsed := int(le.Uint32(b[32:]))

	if hdrSize < 40 || hdrSize > len(b) || width <= 0 || height == 0 || width > 1<<15 {
		return nil, fmt.Errorf("位图头部无效（头部 %d 字节，%d×%d）", hdrSize, width, height)
	}
	topDown := height < 0
	if topDown {
		height = -height
	}
	if height > 1<<15 {
		return nil, fmt.Errorf("位图尺寸无效：%d×%d", width, height)
	}

	pos := hdrSize
	// 颜色掩码：V4/V5 头部里自带；只有 40 字节头部时紧跟在头部之后
	masks := [4]uint32{0x00FF0000, 0x0000FF00, 0x000000FF, 0}
	if compression == biBitfields {
		if hdrSize >= 56 {
			for i := range 4 {
				masks[i] = le.Uint32(b[40+4*i:])
			}
		} else {
			if len(b) < pos+12 {
				return nil, errors.New("位图缺少颜色掩码")
			}
			for i := range 3 {
				masks[i] = le.Uint32(b[pos+4*i:])
			}
			pos += 12
		}
	} else if compression != biRGB {
		return nil, fmt.Errorf("不支持压缩的位图（compression=%d）", compression)
	}

	var palette []color.NRGBA
	if bitCount <= 8 {
		n := clrUsed
		if n == 0 {
			n = 1 << bitCount
		}
		if len(b) < pos+4*n {
			return nil, errors.New("位图调色板不完整")
		}
		for i := range n {
			p := b[pos+4*i:]
			palette = append(palette, color.NRGBA{R: p[2], G: p[1], B: p[0], A: 0xFF})
		}
		pos += 4 * n
	}

	stride := ((width*bitCount + 31) / 32) * 4
	if len(b) < pos+stride*height {
		return nil, errors.New("位图像素数据不完整")
	}
	pixels := b[pos:]

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	alphaSeen := false
	for y := range height {
		row := pixels[y*stride:]
		dy := height - 1 - y
		if topDown {
			dy = y
		}
		for x := range width {
			var c color.NRGBA
			switch bitCount {
			case 32:
				v := le.Uint32(row[4*x:])
				c = color.NRGBA{
					R: maskChannel(v, masks[0]), G: maskChannel(v, masks[1]),
					B: maskChannel(v, masks[2]), A: 0xFF,
				}
				if compression == biRGB {
					c.A = row[4*x+3] // BI_RGB 的第四个字节约定俗成是透明度，但常常全是 0
				} else if masks[3] != 0 {
					c.A = maskChannel(v, masks[3])
				}
				if c.A != 0 {
					alphaSeen = true
				}
			case 24:
				c = color.NRGBA{R: row[3*x+2], G: row[3*x+1], B: row[3*x], A: 0xFF}
			case 8, 4, 1:
				bit := x * bitCount
				idx := int(row[bit/8]>>(8-bitCount-bit%8)) & (1<<bitCount - 1)
				if idx < len(palette) {
					c = palette[idx]
				}
			default:
				return nil, fmt.Errorf("不支持 %d 位的位图", bitCount)
			}
			img.SetNRGBA(x, dy, c)
		}
	}
	// 透明度全为 0 说明这个通道没被使用（很多程序这样写 32 位位图），当作不透明
	if bitCount == 32 && !alphaSeen {
		for i := 3; i < len(img.Pix); i += 4 {
			img.Pix[i] = 0xFF
		}
	}
	return img, nil
}

// maskChannel 按掩码取出一个通道，并缩放到 8 位。
func maskChannel(v, mask uint32) uint8 {
	if mask == 0 {
		return 0
	}
	shift := 0
	for mask&1 == 0 {
		mask >>= 1
		shift++
	}
	val := (v >> shift) & mask
	if mask == 0xFF {
		return uint8(val)
	}
	return uint8(val * 0xFF / mask)
}

// imageToDIBV5 生成 CF_DIBV5：32 位、带透明通道、自下而上。
// 写入时同时放 PNG，这份给只认位图的程序（画图、旧版 Office）用；
// CF_DIB 与 CF_BITMAP 由系统从它自动合成。
func imageToDIBV5(img image.Image) []byte {
	const hdrSize = 124
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	size := w * h * 4

	b := make([]byte, hdrSize+size)
	le := binary.LittleEndian
	le.PutUint32(b[0:], hdrSize)
	le.PutUint32(b[4:], uint32(int32(w)))
	le.PutUint32(b[8:], uint32(int32(h))) // 正数：自下而上
	le.PutUint16(b[12:], 1)               // planes
	le.PutUint16(b[14:], 32)
	le.PutUint32(b[16:], biBitfields)
	le.PutUint32(b[20:], uint32(size))
	le.PutUint32(b[40:], 0x00FF0000) // 红
	le.PutUint32(b[44:], 0x0000FF00) // 绿
	le.PutUint32(b[48:], 0x000000FF) // 蓝
	le.PutUint32(b[52:], 0xFF000000) // 透明度
	le.PutUint32(b[56:], lcsSRGB)
	le.PutUint32(b[108:], 4) // bV5Intent = LCS_GM_IMAGES

	px := b[hdrSize:]
	for y := range h {
		row := px[(h-1-y)*w*4:]
		for x := range w {
			c := color.NRGBAModel.Convert(img.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.NRGBA)
			row[4*x+0] = c.B
			row[4*x+1] = c.G
			row[4*x+2] = c.R
			row[4*x+3] = c.A
		}
	}
	return b
}
