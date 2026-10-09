//go:build darwin

package imagepreview

/*
#cgo LDFLAGS: -framework ImageIO -framework CoreGraphics -framework CoreFoundation
#include <ImageIO/ImageIO.h>
#include <stdlib.h>
#include <string.h>

// 读入 ImageIO 认识的图片，按 EXIF 方向转正、长边缩到 maxPixels 以内（不放大），写成 PNG。
static int cs_convert_png(const char *src, const char *dst, int maxPixels) {
    int ok = 0;
    CFURLRef inURL = CFURLCreateFromFileSystemRepresentation(NULL, (const UInt8 *)src, (CFIndex)strlen(src), false);
    CFURLRef outURL = CFURLCreateFromFileSystemRepresentation(NULL, (const UInt8 *)dst, (CFIndex)strlen(dst), false);
    CGImageSourceRef source = inURL ? CGImageSourceCreateWithURL(inURL, NULL) : NULL;
    if (source != NULL && outURL != NULL) {
        CFNumberRef size = CFNumberCreate(NULL, kCFNumberIntType, &maxPixels);
        const void *keys[] = {
            kCGImageSourceCreateThumbnailFromImageAlways,
            kCGImageSourceCreateThumbnailWithTransform,
            kCGImageSourceThumbnailMaxPixelSize,
        };
        const void *values[] = {kCFBooleanTrue, kCFBooleanTrue, size};
        CFDictionaryRef options = CFDictionaryCreate(NULL, keys, values, 3,
            &kCFTypeDictionaryKeyCallBacks, &kCFTypeDictionaryValueCallBacks);
        CGImageRef image = CGImageSourceCreateThumbnailAtIndex(source, 0, options);
        if (image != NULL) {
            CGImageDestinationRef dest = CGImageDestinationCreateWithURL(outURL, CFSTR("public.png"), 1, NULL);
            if (dest != NULL) {
                CGImageDestinationAddImage(dest, image, NULL);
                ok = CGImageDestinationFinalize(dest) ? 1 : 0;
                CFRelease(dest);
            }
            CGImageRelease(image);
        }
        CFRelease(options);
        CFRelease(size);
    }
    if (source != NULL) CFRelease(source);
    if (inURL != NULL) CFRelease(inURL);
    if (outURL != NULL) CFRelease(outURL);
    return ok;
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Mac 上由 ImageIO 转换的格式
var converted = map[string]bool{".heic": true, ".heif": true, ".tif": true, ".tiff": true}

func convert(src, dst string, maxPixels int) error {
	cSrc, cDst := C.CString(src), C.CString(dst)
	defer C.free(unsafe.Pointer(cSrc))
	defer C.free(unsafe.Pointer(cDst))
	if C.cs_convert_png(cSrc, cDst, C.int(maxPixels)) == 0 {
		return errors.New("无法读取这张图片，文件可能已损坏")
	}
	return nil
}
