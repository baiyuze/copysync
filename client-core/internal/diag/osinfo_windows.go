package diag

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// osVersion 按内部版本号区分 Windows 10 与 11：两者的主版本号都是 10，11 从 22000 开始
func osVersion() string {
	v := windows.RtlGetVersion()
	name := "Windows 10"
	if v.MajorVersion == 10 && v.BuildNumber >= 22000 {
		name = "Windows 11"
	}
	return fmt.Sprintf("%s (build %d)", name, v.BuildNumber)
}
