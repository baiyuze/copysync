//go:build darwin

package lan

import (
	"errors"
	"net/netip"
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
// 直连网段里的对方不在线时，ARP 解析失败也会返回一次同样的错误，见 writeOnLink。
func deniedByOS(err error) bool {
	return errors.Is(err, syscall.EHOSTUNREACH)
}

// probeTarget 返回探测本地网络权限用的地址：默认路由所在网段里、靠近末尾的一个地址
// （不是本机，也不像是网关）。
//
// 还没有配对任何设备时，后台服务不会往局域网地址发单播，也就发现不了自己被拦着，
// 用户只会看到「找不到设备」。往这个地址的丢弃端口发一个几字节的包，被拦与否立刻就知道。
// 不用网关：发往网关的包可能不受这项限制，探测出来的「没被拦」不可信。
func probeTarget() (netip.Addr, bool) {
	primary, ok := defaultRouteIP()
	if !ok {
		return netip.Addr{}, false
	}
	for _, p := range localNetworks() {
		if p.Addr() != primary || p.Bits() > 30 {
			continue
		}
		b := p.Masked().Addr().As4()
		host := uint32(1)<<(32-p.Bits()) - 1 // 网段内的最大偏移，即广播地址
		base := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
		for _, off := range []uint32{host - 2, host - 3} {
			v := base + off
			ip := netip.AddrFrom4([4]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
			if ip != primary {
				return ip, true
			}
		}
	}
	return netip.Addr{}, false
}
