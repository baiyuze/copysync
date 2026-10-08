//go:build !darwin

package clipboard

import "time"

// 非 macOS 平台的占位实现。Windows 的实现是二期工作（M6/M7）：
// 文本与图片走 SetClipboardData + WM_RENDERFORMAT，
// 文件粘贴需要在 Go 里实现 COM IDataObject。
type stubClipboard struct{}

func New() Clipboard { return stubClipboard{} }

func InitPlatform() {}

func (stubClipboard) Peek() (Snapshot, error)  { return Snapshot{}, ErrUnsupported }
func (stubClipboard) Read() (Content, error)   { return Content{}, ErrUnsupported }
func (stubClipboard) Write(Content) error      { return ErrUnsupported }
func (stubClipboard) Permission() Permission   { return PermissionNotApplicable }
func (stubClipboard) OpenPermissionSettings() error { return ErrUnsupported }

// 无原生事件循环可驱动时，单纯让出 CPU。
func pumpRunLoop(seconds float64) {
	time.Sleep(time.Duration(seconds * float64(time.Second)))
}
