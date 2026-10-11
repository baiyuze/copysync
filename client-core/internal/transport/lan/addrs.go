package lan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// cgnat 是运营商级 NAT 的地址段，Tailscale 的 100.x 地址也在这里面。
var cgnat = netip.MustParsePrefix("100.64.0.0/10")

// usableAddr 判断一个本机地址能不能给对端用：只要局域网私有地址和 Tailscale 地址。
func usableAddr(ip netip.Addr) bool {
	return ip.Is4() && (ip.IsPrivate() || cgnat.Contains(ip))
}

// skippedAdapter 判断一块网卡是否不该拿来给对端连：WSL、Hyper-V、Docker 的虚拟网卡
// 只通本机里的虚拟机，地址还每次开机都变。
func skippedAdapter(name string) bool {
	if runtime.GOOS == "windows" {
		return strings.HasPrefix(name, "vEthernet")
	}
	return strings.HasPrefix(name, "docker") || strings.HasPrefix(name, "br-") || strings.HasPrefix(name, "veth")
}

// LocalAddrs 返回本机可供对端连接的 IPv4 地址：默认路由所在的网卡排第一，
// 然后是 Tailscale 地址，最后是其他网卡。
func LocalAddrs() []netip.Addr {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var out []netip.Addr
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 || skippedAdapter(ifc.Name) {
			continue
		}
		list, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, a := range list {
			ipn, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			ip, ok := netip.AddrFromSlice(ipn.IP)
			if ok && usableAddr(ip.Unmap()) {
				out = append(out, ip.Unmap())
			}
		}
	}
	primary, _ := defaultRouteIP()
	rank := func(ip netip.Addr) int {
		switch {
		case ip == primary:
			return 0
		case cgnat.Contains(ip):
			return 1
		}
		return 2
	}
	sort.SliceStable(out, func(i, j int) bool { return rank(out[i]) < rank(out[j]) })
	return dedupeAddrs(out)
}

// defaultRouteIP 返回默认路由所用的本机地址。UDP 的「连接」只查路由表，不会发包。
func defaultRouteIP() (netip.Addr, bool) {
	c, err := net.Dial("udp4", "192.0.2.1:9") // 文档保留地址，不会真的有人收
	if err != nil {
		return netip.Addr{}, false
	}
	defer c.Close()
	ap, err := netip.ParseAddrPort(c.LocalAddr().String())
	if err != nil || !usableAddr(ap.Addr()) {
		return netip.Addr{}, false
	}
	return ap.Addr(), true
}

// localNetworks 返回本机各网卡直连的网段，Addr() 是本机在那个网段里的地址。
func localNetworks() []netip.Prefix {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var out []netip.Prefix
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		list, _ := ifc.Addrs()
		for _, a := range list {
			if ipn, ok := a.(*net.IPNet); ok && ipn.IP.To4() != nil {
				// 不做 Masked：保留本机地址，onLink 要用它排除自己
				ones, _ := ipn.Mask.Size()
				ip, _ := netip.AddrFromSlice(ipn.IP.To4())
				out = append(out, netip.PrefixFrom(ip, ones))
			}
		}
	}
	return out
}

// onLink 判断地址是否在本机直连的网段里（不经路由器就能到），本机自己的地址除外：
// 发往本机地址的包不受 macOS 本地网络权限限制，发成功了也证明不了什么。
// 虚拟机经 NAT 发来的包，源地址就可能是本机的地址。
func onLink(ip netip.Addr) bool {
	nets := localNetworks()
	for _, p := range nets {
		if p.Addr() == ip {
			return false
		}
	}
	for _, p := range nets {
		if p.Contains(ip) {
			return true
		}
	}
	return false
}

// networkKey 标识本机当前所在的网络：网卡地址变了，就要重启发现、重新打招呼。
func networkKey() string {
	var parts []string
	for _, ip := range LocalAddrs() {
		parts = append(parts, ip.String())
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

// DisplayAddr 把地址写成给人看的形式：端口是默认值时省略端口。
func DisplayAddr(ap netip.AddrPort) string {
	if ap.Port() == DefaultPort {
		return ap.Addr().String()
	}
	return ap.String()
}

// ErrBadAddress 表示用户输入的地址无法识别。
var ErrBadAddress = errors.New("地址格式不对")

// ResolveAddress 把用户输入的地址解析成 ip:port。可以是 ip、ip:port 或主机名
// （如 Tailscale 的 MagicDNS 名字），没写端口时用默认端口。
func ResolveAddress(ctx context.Context, s string) ([]netip.AddrPort, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, ErrBadAddress
	}
	host, port := s, DefaultPort
	if h, p, err := net.SplitHostPort(s); err == nil {
		n, err := strconv.Atoi(p)
		if err != nil || n <= 0 || n > 65535 {
			return nil, ErrBadAddress
		}
		host, port = h, n
	}
	if ip, err := netip.ParseAddr(host); err == nil {
		if !ip.Unmap().Is4() {
			return nil, fmt.Errorf("%w：只支持 IPv4", ErrBadAddress)
		}
		return []netip.AddrPort{netip.AddrPortFrom(ip.Unmap(), uint16(port))}, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip4", host)
	if err != nil || len(ips) == 0 {
		return nil, fmt.Errorf("%w：解析 %s 失败", ErrBadAddress, host)
	}
	out := make([]netip.AddrPort, 0, len(ips))
	for _, ip := range ips {
		out = append(out, netip.AddrPortFrom(ip.Unmap(), uint16(port)))
	}
	return out, nil
}

// presenceToken 是设备在某一天（UTC）广播的令牌。
//
// 广播里不放设备名和 ID：同一网络里的其他人看到的只是一串每天都变的字符。
// 已配对的设备知道对方的 ID，自己算一遍就能认出来。
func presenceToken(deviceID string, day time.Time) string {
	sum := sha256.Sum256([]byte("copysync-lan|" + deviceID + "|" + day.UTC().Format("2006-01-02")))
	return hex.EncodeToString(sum[:8])
}

func dedupeAddrs(in []netip.Addr) []netip.Addr {
	seen := make(map[netip.Addr]bool, len(in))
	out := in[:0]
	for _, a := range in {
		if !seen[a] {
			seen[a] = true
			out = append(out, a)
		}
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
