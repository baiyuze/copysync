package p2p

import (
	"context"
	"fmt"
	"hash/fnv"
	"net"
	"net/netip"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pion/ice/v4"
	"github.com/pion/stun/v4"
)

// 出口探测：在打洞用的那一个本地端口上，向多台 STUN 服务器询问本机的公网地址。
//
// 有的网络不止一个出口——公司双线、运营商的 NAT 地址池——并且按目标地址选出口。
// 只问一台服务器，得到的只是「去往那台服务器」的出口地址；发往对端的包可能从另一个
// 出口出去，对端收到的源地址与告诉它的对不上，就会被它的路由器丢弃，打洞随之失败。
//
// 实测（见仓库提交记录）：这类网络里出口只按目标 IP 选择，与目标端口无关；同一条
// 线路上同一个本地端口的公网映射是固定的。所以向 IP 足够分散的一组服务器探测，
// 就能拿到本机在每条线路上的地址，把它们都告诉对端，对端挨个打过来总有一个对得上。

// DefaultProbeServers 是内置的公共 STUN 服务器，按「IP 分散、国内可达」挑选。
// 信令服务器下发的 STUN 地址总会被一并探测；这一组只在用户允许时使用。
var DefaultProbeServers = []string{
	"stun.chat.bilibili.com:3478",
	"stun.cloudflare.com:3478",
	"global.stun.twilio.com:3478",
	"stun.nextcloud.com:3478",
	"stun.l.google.com:19302",
	"stun.freeswitch.org:3478",
}

const (
	probeTimeout      = 1500 * time.Millisecond
	maxSrflxAdvertise = 8 // 发给对端的出口地址上限，太多只会拖慢连通性检查
	historyTTL        = 7 * 24 * time.Hour
)

// Egress 是本机在某个出口上的公网地址。
type Egress struct {
	Addr netip.AddrPort
	// PortPreserved 表示这个出口不改端口（公网端口等于本地端口）。
	// 这样的出口即使某一轮没探测到，也能用「出口 IP + 本地端口」推算出来。
	PortPreserved bool
	// Via 是经由哪些探测服务器发现的；为空表示来自历史记录的推算。
	Via []string
}

// Guessed 表示这个地址是按历史记录推算的，而非本轮实测。
func (e Egress) Guessed() bool { return len(e.Via) == 0 }

// EgressReport 是一轮探测的结果。
type EgressReport struct {
	LocalPort int
	Egresses  []Egress
	Probed    int // 发出探测的服务器数（按 IP 去重后）
	Answered  int // 有响应的服务器数
	At        time.Time
}

// stunProber 由 ice.UniversalUDPMux 实现：在共享的那个端口上做一次 STUN 探测。
type stunProber interface {
	GetXORMappedAddr(stunAddr net.Addr, deadline time.Duration) (*stun.XORMappedAddress, error)
}

var _ stunProber = (*ice.UniversalUDPMuxDefault)(nil)

type probeOptions struct {
	// AllowPrivate 允许内网与回环地址的探测服务器，仅供测试与实验环境使用
	AllowPrivate bool
	Resolver     *net.Resolver
}

// probeEgress 并发探测所有服务器，按公网地址归并结果。
func probeEgress(ctx context.Context, prober stunProber, localPort int, servers []string,
	opts probeOptions) EgressReport {
	targets := resolveProbeTargets(ctx, servers, opts)
	report := EgressReport{LocalPort: localPort, Probed: len(targets), At: time.Now()}

	type answer struct {
		name string
		addr netip.AddrPort
	}
	answers := make(chan answer, len(targets))
	var wg sync.WaitGroup
	for _, t := range targets {
		wg.Add(1)
		go func(t probeTarget) {
			defer wg.Done()
			mapped, err := prober.GetXORMappedAddr(net.UDPAddrFromAddrPort(t.addr), probeTimeout)
			if err != nil || mapped == nil {
				return
			}
			ip, ok := netip.AddrFromSlice(mapped.IP)
			if !ok {
				return
			}
			answers <- answer{name: t.name, addr: netip.AddrPortFrom(ip.Unmap(), uint16(mapped.Port))}
		}(t)
	}
	wg.Wait()
	close(answers)

	byAddr := map[netip.AddrPort]*Egress{}
	var order []netip.AddrPort
	for a := range answers {
		report.Answered++
		e, ok := byAddr[a.addr]
		if !ok {
			e = &Egress{Addr: a.addr, PortPreserved: int(a.addr.Port()) == localPort}
			byAddr[a.addr] = e
			order = append(order, a.addr)
		}
		e.Via = append(e.Via, a.name)
	}
	// 被越多服务器看到的出口越可能是常用出口，排在前面
	sort.SliceStable(order, func(i, j int) bool {
		return len(byAddr[order[i]].Via) > len(byAddr[order[j]].Via)
	})
	for _, a := range order {
		sort.Strings(byAddr[a].Via)
		report.Egresses = append(report.Egresses, *byAddr[a])
	}
	return report
}

type probeTarget struct {
	name string
	addr netip.AddrPort
}

// resolveProbeTargets 解析服务器地址，过滤掉不可能是公网 STUN 的地址，并按 IP 去重。
//
// 按 IP 去重是因为出口只按目标 IP 选择：同一 IP 的不同端口、同一 IP 的多个域名
// （Google 的 stun1–4 就是同一个 IP）只会落在同一个出口上，多探无益。
//
// 各域名并发解析：逐个解析时一台 DNS 慢就拖住整轮探测，而刚开机、刚换网络时
// 正是 DNS 最慢、最需要尽快探完的时候（探完才能把出口地址发给对端）。
func resolveProbeTargets(ctx context.Context, servers []string, opts probeOptions) []probeTarget {
	resolver := opts.Resolver
	if resolver == nil {
		resolver = net.DefaultResolver
	}
	type resolved struct {
		host string
		port uint16
		ips  []netip.Addr
	}
	results := make([]resolved, len(servers))
	var wg sync.WaitGroup
	for i, s := range servers {
		host, port, err := parseProbeServer(s)
		if err != nil {
			continue
		}
		results[i] = resolved{host: host, port: port}
		if ip, err := netip.ParseAddr(host); err == nil {
			results[i].ips = []netip.Addr{ip}
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			lctx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			if ips, err := resolver.LookupNetIP(lctx, "ip4", host); err == nil {
				results[i].ips = ips
			}
		}()
	}
	wg.Wait()

	// 按服务器原本的顺序去重，结果与解析快慢无关
	seen := map[netip.Addr]bool{}
	var out []probeTarget
	for _, r := range results {
		host, port := r.host, r.port
		for _, ip := range r.ips {
			ip = ip.Unmap()
			if !ip.Is4() || seen[ip] || !usableProbeIP(ip, opts.AllowPrivate) {
				continue
			}
			seen[ip] = true
			out = append(out, probeTarget{name: host, addr: netip.AddrPortFrom(ip, port)})
			break // 一个域名取一个可用 IP 就够了
		}
	}
	return out
}

// parseProbeServer 接受 "stun:host:port"、"stun:host"、"host:port" 等写法。
func parseProbeServer(s string) (string, uint16, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "stuns:"), "stun:")
	if i := strings.IndexByte(s, '?'); i >= 0 {
		s = s[:i]
	}
	if s == "" {
		return "", 0, fmt.Errorf("空地址")
	}
	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
		return s, 3478, nil // 没写端口，用 STUN 默认端口
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil || port == 0 {
		return "", 0, fmt.Errorf("端口无效: %q", portStr)
	}
	return host, uint16(port), nil
}

// 不可能是公网 STUN 服务器的地址段。DNS 污染与代理软件的 fake-ip 会把域名解析到这些地方，
// 向它们发探测只会超时，或者更糟——被当成一个出口。
var bogusPrefixes = mustPrefixes(
	"0.0.0.0/8",       // 本网络
	"100.64.0.0/10",   // 运营商 NAT 内部地址
	"192.0.0.0/24",    // IETF 协议分配
	"192.0.2.0/24",    // 文档示例（实测某些域名被污染到这里）
	"198.18.0.0/15",   // 基准测试段，Clash 等代理的 fake-ip 默认用它
	"198.51.100.0/24", // 文档示例
	"203.0.113.0/24",  // 文档示例
	"240.0.0.0/4",     // 保留
)

func usableProbeIP(ip netip.Addr, allowPrivate bool) bool {
	if !ip.IsValid() || ip.IsUnspecified() || ip.IsMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip.IsLoopback() || ip.IsPrivate() {
		return allowPrivate
	}
	for _, p := range bogusPrefixes {
		if p.Contains(ip) {
			return false
		}
	}
	return true
}

func mustPrefixes(ss ...string) []netip.Prefix {
	out := make([]netip.Prefix, len(ss))
	for i, s := range ss {
		out[i] = netip.MustParsePrefix(s)
	}
	return out
}

// srflxCandidates 把探测结果转成发给对端的 srflx candidate 字符串。
//
// 实测到的地址排在前面、优先级更高；推算的地址垫底。数量有上限：每多一个地址，
// 对端就要多做一轮连通性检查。
func srflxCandidates(r EgressReport) []string {
	var out []string
	seen := map[netip.AddrPort]bool{}
	add := func(e Egress, localPref int) {
		if len(out) >= maxSrflxAdvertise || seen[e.Addr] {
			return
		}
		seen[e.Addr] = true
		// RFC 8445 5.1.2.1：priority = 2^24*type + 2^8*local + (256 - component)
		priority := (1<<24)*100 + (1<<8)*localPref + 255
		out = append(out, fmt.Sprintf("candidate:%d 1 udp %d %s %d typ srflx raddr 0.0.0.0 rport %d",
			candidateFoundation(e.Addr), priority, e.Addr.Addr(), e.Addr.Port(), r.LocalPort))
	}
	var measured, guessed []Egress
	for _, e := range r.Egresses {
		if e.Guessed() {
			guessed = append(guessed, e)
		} else {
			measured = append(measured, e)
		}
	}
	for i, e := range measured {
		add(e, 65535-i)
	}
	for i, e := range guessed {
		add(e, 32767-i)
	}
	return out
}

func candidateFoundation(a netip.AddrPort) uint32 {
	h := fnv.New32a()
	h.Write([]byte(a.String()))
	return h.Sum32()
}
