//go:build windows

package lan

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

// RequestPermission 申请一次管理员权限（弹出 UAC），在 Windows 防火墙里放行后台服务。
//
// 第一次运行时 Windows 会问是否允许后台服务通过防火墙；点了取消会留下阻止规则，
// 而且那条规则只对当时的网络类型生效。不放行时本机仍能主动连别人，
// 但别人连不上本机——跨网段配对时要输入地址的那一方正是被连的一方。
// 这里删掉本程序原有的入站规则（包括阻止规则，阻止优先于允许），
// 再加一条对所有网络类型生效的 UDP 允许规则。
func RequestPermission() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if strings.ContainsRune(exe, '"') {
		return fmt.Errorf("程序路径含引号，无法放行: %s", exe)
	}
	script := fmt.Sprintf(
		`/c netsh advfirewall firewall delete rule name=all dir=in program="%[1]s" & `+
			`netsh advfirewall firewall add rule name="CopySync" dir=in action=allow program="%[1]s" protocol=UDP profile=any`,
		exe)
	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString("cmd.exe")
	args, err := windows.UTF16PtrFromString(script)
	if err != nil {
		return err
	}
	// 用户在 UAC 里点了「否」时返回 ERROR_CANCELLED
	return windows.ShellExecute(0, verb, file, args, nil, windows.SW_HIDE)
}

// Windows 不会在发送时拦截，防火墙只挡入站。
func deniedByOS(error) bool { return false }
