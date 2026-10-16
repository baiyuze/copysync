package diag

import (
	"cmp"
	"fmt"
	"maps"
	"net/netip"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// redactor 把报告里能认出人的信息换成占位符：设备名、服务器地址、公网 IP、用户目录、文件名。
//
// 局域网与 Tailscale 的地址（192.168.x.x、10.x、100.64/10 等）保留：只在用户自己的网络里有意义，
// 「不在同一网段」「访客网络隔离」这类问题全靠它们判断。
type redactor struct {
	words *strings.Replacer
}

var (
	ipv4Pattern = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	// 先粗找由十六进制数字和冒号组成的一段，能不能算 IPv6 交给 netip 判断：
	// 时间 20:50:01、MAC 地址都解析不成 IPv6，不会被误换
	ipv6Pattern = regexp.MustCompile(`[0-9A-Fa-f]*:[0-9A-Fa-f:]*:[0-9A-Fa-f]*`)
	cgnat       = netip.MustParsePrefix("100.64.0.0/10") // 运营商级 NAT，Tailscale 也用这一段
)

// 太短的名字（比如叫「PC」的设备、叫「a」的文件）替换掉会误伤日志里的普通字词，宁可放过
const minNameRunes = 3

func newRedactor(in Input) *redactor {
	type pair struct{ old, new string }
	var pairs []pair
	add := func(old, new string) {
		old = strings.TrimSpace(old)
		if utf8.RuneCountInString(old) >= minNameRunes {
			pairs = append(pairs, pair{old, new})
		}
	}

	// 设备名：本机一个，配对的设备按 ID 排序编号，同一份报告里同一台设备始终是同一个代号
	add(in.Config.DeviceName, "<this-device>")
	if in.Self != nil {
		add(in.Self.GetName(), "<this-device>")
	}
	names := map[string]string{} // 设备 ID → 名字
	for _, d := range in.Peers {
		names[d.GetId()] = d.GetName()
	}
	for _, p := range in.Network.GetLan().GetPeers() {
		if _, ok := names[p.GetDeviceId()]; !ok {
			names[p.GetDeviceId()] = p.GetDeviceName()
		}
	}
	for i, id := range slices.Sorted(maps.Keys(names)) {
		add(names[id], fmt.Sprintf("<device-%d>", i+1))
	}

	// 服务器：主机名或公网 IP 换掉，协议、端口、路径留着，排查时用得上
	for _, address := range []string{in.Config.EffectiveSignalingURL(), in.Config.SignalingURL, in.Config.CustomSignalingURL} {
		if u, err := url.Parse(address); err == nil {
			if host := u.Hostname(); host != "" && !keepIP(host) {
				add(host, "<server>")
			}
		}
	}

	// 用户目录里有用户名；Windows 上的路径两种斜杠都可能出现
	if in.Home != "" {
		add(in.Home, "~")
		add(strings.ReplaceAll(in.Home, `\`, "/"), "~")
	}

	for _, name := range in.FileNames {
		add(name, "<file>")
	}

	// 长的先换：设备名「张三的 MacBook Pro」不能先被它的一部分换掉
	slices.SortStableFunc(pairs, func(a, b pair) int { return cmp.Compare(len(b.old), len(a.old)) })
	args := make([]string, 0, 2*len(pairs))
	for _, p := range pairs {
		args = append(args, p.old, p.new)
	}
	return &redactor{words: strings.NewReplacer(args...)}
}

func (r *redactor) apply(s string) string {
	s = r.words.Replace(s)
	s = ipv4Pattern.ReplaceAllStringFunc(s, func(m string) string {
		if keepIP(m) {
			return m
		}
		return "<public-ip>"
	})
	return ipv6Pattern.ReplaceAllStringFunc(s, func(m string) string {
		if _, err := netip.ParseAddr(m); err != nil || keepIP(m) {
			return m
		}
		return "<public-ip>"
	})
}

// keepIP 判断一个地址能不能原样留在报告里：局域网、本机、链路本地、运营商级 NAT（含 Tailscale）可以。
// 解析不成 IP 的（比如域名）返回 false。
func keepIP(s string) bool {
	a, err := netip.ParseAddr(s)
	if err != nil {
		return false
	}
	a = a.Unmap()
	return a.IsPrivate() || a.IsLoopback() || a.IsLinkLocalUnicast() || a.IsUnspecified() ||
		a.IsMulticast() || cgnat.Contains(a)
}
