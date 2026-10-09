// Package imagepreview 决定复制记录里哪些图片能在界面上预览，并把图片准备成界面能直接解码的文件。
//
// 界面（Flutter）自己能解码 PNG、JPEG、GIF、WebP、BMP。iPhone 照片常见的 HEIC 和 TIFF 解不了：
// Mac 上由系统的 ImageIO 转成 PNG（见 convert_darwin.go），其他平台不支持，也就不算可预览。
package imagepreview

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/baiyuze/copysync/client-core/internal/store"
)

// MaxPixels 是转换出来的预览图长边的上限，与界面解码时的限制一致。
const MaxPixels = 2048

// 界面能直接解码的格式
var native = map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true, ".bmp": true}

// Previewable 按文件名判断这是不是一张能预览的图片。
func Previewable(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return native[ext] || converted[ext]
}

// ImageCount 返回一条记录里能预览的图片数：剪贴板图片算一张，文件记录数其中的图片文件（不进文件夹）。
func ImageCount(c store.Clip) int32 {
	switch c.Kind {
	case store.KindImage:
		return 1
	case store.KindFile:
		var n int32
		for _, it := range c.Items {
			if !it.IsDir && Previewable(it.Name) {
				n++
			}
		}
		return n
	}
	return 0
}

// Prepare 返回界面可以直接读取的图片路径。原生格式原样返回；其他格式转换成 PNG，
// 放在系统临时目录里，按原文件的路径、大小、修改时间命名，同一张图第二次打开不再转换。
func Prepare(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if native[ext] {
		return path, nil
	}
	if !converted[ext] {
		return "", errors.New("不支持预览这种图片格式")
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d", path, info.Size(), info.ModTime().UnixNano())))
	dir := filepath.Join(os.TempDir(), "CopySync-preview")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	out := filepath.Join(dir, hex.EncodeToString(sum[:10])+".png")
	if _, err := os.Stat(out); err == nil {
		return out, nil
	}
	// 先写到临时名再改名：转换到一半时界面读到的不会是残缺的文件
	tmp := out + ".partial"
	if err := convert(path, tmp, MaxPixels); err != nil {
		os.Remove(tmp)
		return "", err
	}
	if err := os.Rename(tmp, out); err != nil {
		os.Remove(tmp)
		return "", err
	}
	return out, nil
}
