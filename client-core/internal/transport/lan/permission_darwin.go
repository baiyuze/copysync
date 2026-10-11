//go:build darwin

package lan

import (
	"errors"
	"os/exec"
	"syscall"
)

// RequestPermission 打开「隐私与安全性 → 本地网络」，让用户放行后台服务。
//
// macOS 15 起，后台服务主动往局域网地址发包要经用户允许，没允许时 sendto
// 直接返回 EHOSTUNREACH（对方先发来的包仍能收到、能回复）。
func RequestPermission() error {
	return exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_LocalNetwork").Run()
}

// deniedByOS 判断一次发送失败是不是被系统的本地网络权限拦下的。
// 只对本机直连网段里的地址下这个结论：真的没有路由时也会是同一个错误。
func deniedByOS(err error) bool {
	return errors.Is(err, syscall.EHOSTUNREACH)
}
