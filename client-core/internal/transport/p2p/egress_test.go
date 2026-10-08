package p2p

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"path/filepath"
	"testing"
	"time"

	"github.com/pion/ice/v4"
	"github.com/pion/stun/v4"
)

// fakeProber 按服务器地址返回预设的公网映射，模拟「按目标 IP 选出口」的网络。
type fakeProber map[string]string // 服务器 ip:port -> 本机公网 ip:port

func (f fakeProber) GetXORMappedAddr(addr net.Addr, _ time.Duration) (*stun.XORMappedAddress, error) {
	mapped, ok := f[addr.String()]
	if !ok {
		return nil, errors.New("超时")
	}
	ap := netip.MustParseAddrPort(mapped)
	return &stun.XORMappedAddress{IP: ap.Addr().AsSlice(), Port: int(ap.Port())}, nil
}

func TestParseProbeServer(t *testing.T) {
	cases := []struct {
		in   string
		host string
		port uint16
	}{
		{"stun:stun.l.google.com:19302", "stun.l.google.com", 19302},
		{"stun:106.12.83.160:3478", "106.12.83.160", 3478},
		{"stun.cloudflare.com:3478", "stun.cloudflare.com", 3478},
		{"stun:example.com", "example.com", 3478},
		{"stun:example.com:3478?transport=udp", "example.com", 3478},
	}
	for _, c := range cases {
		host, port, err := parseProbeServer(c.in)
		if err != nil || host != c.host || port != c.port {
			t.Errorf("parseProbeServer(%q) = %q, %d, %v；期望 %q, %d", c.in, host, port, err, c.host, c.port)
		}
	}
}

func TestUsableProbeIP(t *testing.T) {
	cases := []struct {
		ip           string
		allowPrivate bool
		want         bool
	}{
		{"106.12.83.160", false, true},
		{"162.159.207.0", false, true},
		{"192.0.2.42", false, false}, // 实测 stun.syncthing.net 被 DNS 污染到这里
		{"198.18.3.7", false, false}, // 代理软件 fake-ip 的默认网段
		{"100.64.1.1", false, false}, // 运营商 NAT 内部地址
		{"0.0.0.0", false, false},
		{"10.1.2.3", false, false},
		{"10.1.2.3", true, true}, // 实验环境允许内网
		{"127.0.0.1", false, false},
		{"127.0.0.1", true, true},
	}
	for _, c := range cases {
		if got := usableProbeIP(netip.MustParseAddr(c.ip), c.allowPrivate); got != c.want {
			t.Errorf("usableProbeIP(%s, allowPrivate=%v) = %v，期望 %v", c.ip, c.allowPrivate, got, c.want)
		}
	}
}

// 公司网络的情形：不同的服务器落在不同出口上，同一出口被多台服务器看到。
func TestProbeEgressMergesByMappedAddress(t *testing.T) {
	prober := fakeProber{
		"1.1.1.1:3478": "9.9.9.1:5000", // 出口 1，不改端口
		"2.2.2.2:3478": "9.9.9.2:6000", // 出口 2，改端口
		"3.3.3.3:3478": "9.9.9.2:6000", // 也走出口 2：同一线路上映射固定
		// 4.4.4.4 不响应
	}
	servers := []string{
		"stun:1.1.1.1:3478", "2.2.2.2:3478", "3.3.3.3:3478", "4.4.4.4:3478",
		"1.1.1.1:19302",   // 与第一个同 IP：出口只按 IP 选，不该重复探测
		"192.0.2.42:3478", // 被污染的地址：直接丢弃
	}
	r := probeEgress(context.Background(), prober, 5000, servers, probeOptions{})

	if r.Probed != 4 || r.Answered != 3 {
		t.Fatalf("探测 %d 个、响应 %d 个，期望 4 与 3", r.Probed, r.Answered)
	}
	if len(r.Egresses) != 2 {
		t.Fatalf("出口数 = %d，期望 2：%+v", len(r.Egresses), r.Egresses)
	}
	// 被更多服务器看到的出口排在前面
	if e := r.Egresses[0]; e.Addr.String() != "9.9.9.2:6000" || len(e.Via) != 2 || e.PortPreserved {
		t.Errorf("第一个出口 = %+v，期望 9.9.9.2:6000、2 个来源、改端口", e)
	}
	if e := r.Egresses[1]; e.Addr.String() != "9.9.9.1:5000" || !e.PortPreserved {
		t.Errorf("第二个出口 = %+v，期望 9.9.9.1:5000、不改端口", e)
	}
}

func TestSrflxCandidatesAreValidAndOrdered(t *testing.T) {
	r := EgressReport{LocalPort: 5000}
	for i := 0; i < 10; i++ {
		r.Egresses = append(r.Egresses, Egress{
			Addr: netip.AddrPortFrom(netip.AddrFrom4([4]byte{9, 9, 9, byte(i + 1)}), 5000),
			Via:  []string{"stun"},
		})
	}
	// 推算的地址排在实测之后；与实测重复的不再发
	r.Egresses = append([]Egress{{Addr: netip.MustParseAddrPort("8.8.8.8:5000"), PortPreserved: true}}, r.Egresses...)

	cands := srflxCandidates(r)
	if len(cands) != maxSrflxAdvertise {
		t.Fatalf("候选数 = %d，应截断到 %d", len(cands), maxSrflxAdvertise)
	}
	var last uint32 = 1<<32 - 1
	for i, s := range cands {
		c, err := ice.UnmarshalCandidate(s)
		if err != nil {
			t.Fatalf("pion 解析不了第 %d 个候选 %q: %v", i, s, err)
		}
		if c.Type() != ice.CandidateTypeServerReflexive {
			t.Errorf("第 %d 个候选类型 = %v，应为 srflx", i, c.Type())
		}
		if c.Priority() > last {
			t.Errorf("候选 %d 的优先级高于前一个：实测地址应排在前面", i)
		}
		last = c.Priority()
		if c.Address() == "8.8.8.8" {
			t.Errorf("推算的地址挤掉了实测地址：%q", s)
		}
	}
}

func TestHistoryFillsInMissedPreservedEgress(t *testing.T) {
	path := filepath.Join(t.TempDir(), "egress.json")
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)

	h := loadEgressHistory(path)
	h.now = func() time.Time { return now }
	h.record("office", EgressReport{LocalPort: 5000, Egresses: []Egress{
		{Addr: netip.MustParseAddrPort("9.9.9.1:5000"), PortPreserved: true, Via: []string{"a"}},
		{Addr: netip.MustParseAddrPort("9.9.9.2:6000"), Via: []string{"b"}},
	}})

	// 重新加载：历史要能持久化
	h = loadEgressHistory(path)
	h.now = func() time.Time { return now.Add(time.Hour) }

	// 新一轮换了本地端口，且只探测到出口 2
	got := h.augment("office", EgressReport{LocalPort: 7000, Egresses: []Egress{
		{Addr: netip.MustParseAddrPort("9.9.9.2:6100"), Via: []string{"b"}},
	}})
	if len(got.Egresses) != 2 {
		t.Fatalf("出口数 = %d，期望补上不改端口的出口 1：%+v", len(got.Egresses), got.Egresses)
	}
	if g := got.Egresses[1]; g.Addr.String() != "9.9.9.1:7000" || !g.Guessed() {
		t.Errorf("补上的出口 = %+v，期望 9.9.9.1:7000（出口 IP + 新的本地端口）", g)
	}

	// 别的网络不受影响；过期的历史不再使用
	if other := h.augment("home", EgressReport{LocalPort: 7000}); len(other.Egresses) != 0 {
		t.Errorf("其他网络不该用这份历史：%+v", other.Egresses)
	}
	h.now = func() time.Time { return now.Add(historyTTL + time.Hour) }
	if stale := h.augment("office", EgressReport{LocalPort: 7000}); len(stale.Egresses) != 0 {
		t.Errorf("过期的历史不该再用：%+v", stale.Egresses)
	}
}
