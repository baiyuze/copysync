//go:build darwin

package clipboard

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework UniformTypeIdentifiers

#include <stdlib.h>

long  cs_change_count(void);
int   cs_access_behavior(void);
char *cs_types_json(void);
char *cs_detect_content_type(void);
char *cs_read_text(void);
char *cs_read_html(void);
char *cs_read_file_urls_json(void);
int   cs_read_image_png(void **out, int *outLen);
int   cs_write_text(const char *utf8);
int   cs_write_html(const char *html, const char *plain);
int   cs_write_image_png(const void *data, int len);
int   cs_write_files_json(const char *pathsJSON);
void  cs_open_permission_settings(void);
void  cs_run_loop_once(double seconds);
void  cs_init_app(void);
*/
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// darwinClipboard 是 macOS 的实现。所有方法必须在主线程调用。
type darwinClipboard struct{}

// New 返回当前平台的剪贴板实现。
func New() Clipboard { return darwinClipboard{} }

// InitPlatform 初始化 NSApplication。必须在主线程、进入 MainLoop 之前调用。
func InitPlatform() { C.cs_init_app() }

// goString 接管 C 侧 strdup 出来的字符串并释放它。
func goString(c *C.char) string {
	if c == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(c))
	return C.GoString(c)
}

func goStringSlice(c *C.char) []string {
	raw := goString(c)
	if raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func (darwinClipboard) Permission() Permission {
	switch C.cs_access_behavior() {
	case -1:
		return PermissionNotApplicable
	case 0:
		return PermissionDefault
	case 1:
		return PermissionAsk
	case 2:
		return PermissionAlwaysAllow
	case 3:
		return PermissionAlwaysDeny
	default:
		return PermissionNotApplicable
	}
}

// Peek 只做免授权的探测：变更计数、类型列表、内容类型。
//
// M-1 实测这三项在隐私机制开启时合计耗时不到 0.03 秒且不触发弹窗，
// 因此可以按 300ms 的频率安全地轮询。
func (c darwinClipboard) Peek() (Snapshot, error) {
	s := Snapshot{
		ChangeCount: int64(C.cs_change_count()),
		Types:       goStringSlice(C.cs_types_json()),
		ContentType: goString(C.cs_detect_content_type()),
	}
	s.Kind = classify(s.Types, s.ContentType)
	return s, nil
}

// classify 由类型列表和 UTType 推断内容分类。
//
// 顺序有讲究：文件优先于文本——Finder 复制文件时剪贴板里往往同时带着
// 文件名的纯文本表示，先判文本会把文件误判成文本。
func classify(types []string, contentType string) Kind {
	has := func(t string) bool {
		for _, v := range types {
			if v == t {
				return true
			}
		}
		return false
	}

	switch {
	case has("public.file-url"), has("NSFilenamesPboardType"):
		return KindFile
	case has("public.png"), has("public.tiff"):
		return KindImage
	case has("public.html"):
		return KindHTML
	case has("public.utf8-plain-text"), has("NSStringPboardType"):
		return KindText
	}

	// 类型列表没命中时退回 detectMetadata 拿到的 UTType
	switch {
	case strings.HasPrefix(contentType, "public.image"),
		strings.HasSuffix(contentType, ".png"),
		strings.HasSuffix(contentType, ".jpeg"):
		return KindImage
	case contentType == "public.html":
		return KindHTML
	case contentType != "":
		return KindFile // 有具体 UTType 通常意味着这是个文件引用
	}
	return KindUnknown
}

// Read 读取剪贴板实际内容。
//
// 调用前务必确认 Permission().CanReadSilently()，否则在隐私机制开启且
// 未授权时会阻塞在授权弹窗上——而 daemon 是后台进程，弹窗未必可见。
func (c darwinClipboard) Read() (Content, error) {
	snap, err := c.Peek()
	if err != nil {
		return Content{}, err
	}

	out := Content{Kind: snap.Kind}
	switch snap.Kind {
	case KindFile:
		out.Files = goStringSlice(C.cs_read_file_urls_json())
		if len(out.Files) == 0 {
			return Content{}, errors.New("剪贴板中没有可读取的文件路径")
		}

	case KindText:
		out.Text = goString(C.cs_read_text())

	case KindHTML:
		out.HTML = goString(C.cs_read_html())
		out.Text = goString(C.cs_read_text())

	case KindImage:
		var buf unsafe.Pointer
		var n C.int
		if C.cs_read_image_png(&buf, &n) == 0 || buf == nil {
			return Content{}, errors.New("读取剪贴板图片失败")
		}
		defer C.free(buf)
		out.Image = C.GoBytes(buf, n)

	default:
		return Content{}, fmt.Errorf("不支持的剪贴板类型：%v", snap.Types)
	}
	return out, nil
}

// Write 覆盖本机剪贴板。
//
// 文件类型要求路径已真实存在——M-1 验证过延迟渲染在 macOS 的粘贴路径上
// 不成立，所以内容必须先落盘再写剪贴板。
func (c darwinClipboard) Write(content Content) error {
	switch content.Kind {
	case KindFile:
		if len(content.Files) == 0 {
			return errors.New("没有要写入的文件路径")
		}
		data, err := json.Marshal(content.Files)
		if err != nil {
			return err
		}
		cs := C.CString(string(data))
		defer C.free(unsafe.Pointer(cs))
		if C.cs_write_files_json(cs) == 0 {
			return errors.New("写入文件到剪贴板失败（路径可能不存在）")
		}

	case KindText:
		cs := C.CString(content.Text)
		defer C.free(unsafe.Pointer(cs))
		if C.cs_write_text(cs) == 0 {
			return errors.New("写入文本到剪贴板失败")
		}

	case KindHTML:
		chtml := C.CString(content.HTML)
		defer C.free(unsafe.Pointer(chtml))
		cplain := C.CString(content.Text)
		defer C.free(unsafe.Pointer(cplain))
		if C.cs_write_html(chtml, cplain) == 0 {
			return errors.New("写入 HTML 到剪贴板失败")
		}

	case KindImage:
		if len(content.Image) == 0 {
			return errors.New("没有要写入的图片数据")
		}
		if C.cs_write_image_png(
			unsafe.Pointer(&content.Image[0]), C.int(len(content.Image))) == 0 {
			return errors.New("写入图片到剪贴板失败")
		}

	default:
		return fmt.Errorf("不支持写入的类型：%v", content.Kind)
	}
	return nil
}

func (darwinClipboard) OpenPermissionSettings() error {
	C.cs_open_permission_settings()
	return nil
}

// pumpRunLoop 驱动 macOS 主 run loop。
//
// 剪贴板 API 的跨进程回调与授权弹窗都依赖 run loop 转动，
// 所以 MainLoop 在等待任务的间隙必须调用它。
func pumpRunLoop(seconds float64) {
	C.cs_run_loop_once(C.double(seconds))
}
