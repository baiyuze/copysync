//go:build !darwin && !windows && !linux

package clipboard

import "time"

// macOS、Windows、Linux 以外的平台的占位实现。
type stubClipboard struct{}

func New() Clipboard { return stubClipboard{} }

func InitPlatform() {}

func (stubClipboard) Peek() (Snapshot, error)       { return Snapshot{}, ErrUnsupported }
func (stubClipboard) Read() (Content, error)        { return Content{}, ErrUnsupported }
func (stubClipboard) Write(Content) error           { return ErrUnsupported }
func (stubClipboard) Permission() Permission        { return PermissionNotApplicable }
func (stubClipboard) OpenPermissionSettings() error { return ErrUnsupported }

// 无原生事件循环可驱动时，单纯让出 CPU。
func pumpRunLoop(seconds float64) {
	time.Sleep(time.Duration(seconds * float64(time.Second)))
}
