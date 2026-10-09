//go:build !darwin

package imagepreview

import "errors"

// 只有 Mac 能转换 HEIC、TIFF；这里没有需要转换的格式，Previewable 也就不会认它们。
var converted = map[string]bool{}

func convert(string, string, int) error {
	return errors.New("不支持预览这种图片格式")
}
