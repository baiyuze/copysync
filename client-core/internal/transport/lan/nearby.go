package lan

import (
	"context"
	"encoding/base64"
	"errors"
	"net/netip"
	"sort"
	"sync"
	"time"
	"unicode/utf8"

	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// 「搜索附近设备」：窗口打开期间在局域网里公开设备名，看得到附近同样打开了窗口的设备，
// 点一下就能配对，不用输配对码。
//
//   - 只在窗口打开期间公开：广播里平时只有匿名令牌，打开窗口才加上设备名（n）和系统（o）。
//     设备名用 base64url 编码，经 Windows 的系统 DNS-SD 服务转一道也不会走样。
//   - mDNS 不会通知实例下线（pion 只报新出现和有变化的实例），所以两台都打开窗口时
//     每 2 秒互发一个 LanNearbyPing，几秒收不到就当对方走了。它同时是敲门：
//     没放行防火墙的 Windows 先往对方发过包，对方随后的配对请求才进得来。
//     ping 自带名字、系统和令牌：对方开着窗口很久没有变化时，mDNS 不会再报它，
//     本机可能早把它的实例过期丢掉了，收到 ping 就能重新显示。
//   - 配对请求不带配对码，只在窗口打开期间接受，由对方的用户在弹窗里决定接不接受；
//     之后双方照常核对指纹，防中间人的办法与用配对码时相同。

const (
	// 这么久没收到对方的 ping、mDNS 也没有它的新消息，就从列表里拿掉。对方每 2 秒（tick）发一次
	nearbyStale    = 7 * time.Second
	maxNearbyName  = 60 // 设备名在广播里最多占的字节数（编码前）
	maxNearbyAddrs = 4
)

// NearbyInstance 是附近一台打开了「搜索附近设备」窗口的设备。
type NearbyInstance struct {
	ID       string // mDNS 实例名，本次启动内有效，用来发起配对
	Name     string
	Platform string
	// PeerID 非空表示这是已配对的设备（按广播里的令牌认出来的）
	PeerID string
}

type nearbyState struct {
	addrs map[netip.AddrPort]time.Time // 对方的 ping 实际来自的地址
}

// ShowNearby 让本机出现在附近设备的「搜索附近设备」窗口里，直到调用返回的 hide。
// 可以叠加调用（同时开着两个界面），全部 hide 之后才撤下。
func (t *Transport) ShowNearby() (hide func()) {
	name, platform := "", ""
	if t.opts.DeviceInfo != nil {
		name, platform = t.opts.DeviceInfo()
	}
	t.mu.Lock()
	t.visible++
	t.visibleName, t.visiblePlatform = name, platform
	t.mu.Unlock()
	t.updateAdvert()
	t.pingNearby()

	var once sync.Once
	return func() {
		once.Do(func() {
			t.mu.Lock()
			t.visible--
			if t.visible == 0 {
				t.nearby = make(map[string]*nearbyState)
			}
			t.mu.Unlock()
			t.updateAdvert()
		})
	}
}

// Nearby 返回附近打开了「搜索附近设备」窗口的设备，按名字排序。本机没打开窗口时为空。
func (t *Transport) Nearby() []NearbyInstance {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.running || t.visible == 0 {
		return nil
	}
	var out []NearbyInstance
	for id, si := range t.instances {
		name, ok := nearbyName(si.inst.Text)
		if !ok || now.Sub(si.seen) > nearbyStale {
			continue
		}
		out = append(out, NearbyInstance{
			ID:       id,
			Name:     name,
			Platform: si.inst.Text["o"],
			PeerID:   t.tokens[si.inst.Text["t"]],
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// PairNearby 向附近的一台设备发起配对，直到 ctx 结束。对方关掉了窗口或一直没有回音时返回 ErrNotFound。
func (t *Transport) PairNearby(ctx context.Context, id string) (PeerInfo, error) {
	if !t.Running() {
		return PeerInfo{}, ErrDisabled
	}
	targets := func() []netip.AddrPort {
		t.mu.Lock()
		defer t.mu.Unlock()
		return t.nearbyAddrsLocked(id)
	}
	if len(targets()) == 0 {
		return PeerInfo{}, ErrNotFound
	}
	peer, err := t.exchange(ctx, t.selfPairInfo(), nearbyRequest, targets, false)
	if errors.Is(err, ErrCodeClosed) {
		// 对方明确拒收：它已经关掉了窗口
		return PeerInfo{}, ErrNotFound
	}
	return peer, err
}

func (t *Transport) handleNearbyRequest(ps *pb.LanPairSigned, src netip.AddrPort) {
	info, ok := t.openPairInfo(ps)
	if !ok || len(info.GetNonce()) < 8 || t.resendResponse(info.GetNonce(), src) {
		return
	}
	t.mu.Lock()
	visible := t.visible > 0
	t.mu.Unlock()
	if !visible {
		t.reject(src, info.GetNonce(), "closed")
		return
	}
	peer := peerFromInfo(info, src)
	if t.opts.AcceptNearby == nil || !t.opts.AcceptNearby(peer) {
		t.reject(src, info.GetNonce(), "closed")
		return
	}
	t.log.Info("附近的设备请求配对，等待用户核对指纹", "peer", peer.Name, "src", src)
	t.respondPair(info, peer, src)
}

// handleNearbyPing 记下对方还在、它实际在哪个地址，以及 ping 里自带的名字、系统和令牌。
func (t *Transport) handleNearbyPing(p *pb.LanNearbyPing, src netip.AddrPort) {
	id := p.GetInstance()
	if id == "" || id == t.instance || p.GetName() == "" {
		return
	}
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.visible == 0 {
		return
	}
	si, ok := t.instances[id]
	if !ok {
		// mDNS 只在实例有变化时才报：本机可能早把它过期丢掉了，或者还没收到它的广播
		si = seenInstance{inst: Instance{Name: id, Addr: src}}
	}
	text := copyText(si.inst.Text)
	text["v"] = "1"
	text["n"] = encodeNearbyName(p.GetName())
	text["o"] = p.GetPlatform()
	if p.GetToken() != "" {
		text["t"] = p.GetToken()
	}
	si.inst.Text = text
	// 对方一直开着窗口时，靠 ping 让它保持新鲜
	si.seen = now
	t.instances[id] = si
	ns := t.nearby[id]
	if ns == nil {
		ns = &nearbyState{addrs: make(map[netip.AddrPort]time.Time)}
		t.nearby[id] = ns
	}
	ns.addrs[src] = now
}

// pingNearby 告诉附近打开了窗口的每台设备：本机还在。本机没打开窗口时什么也不做。
func (t *Transport) pingNearby() {
	t.mu.Lock()
	if !t.running || t.visible == 0 {
		t.mu.Unlock()
		return
	}
	var targets []netip.AddrPort
	for id, si := range t.instances {
		if _, ok := nearbyName(si.inst.Text); ok {
			targets = append(targets, t.nearbyAddrsLocked(id)...)
		}
	}
	t.mu.Unlock()
	t.sendNearbyPing(targets)
}

func (t *Transport) sendNearbyPing(targets []netip.AddrPort) {
	if len(targets) == 0 {
		return
	}
	t.mu.Lock()
	name, platform := t.nearbySelfLocked()
	t.mu.Unlock()
	pkt, err := encode(&pb.LanPacket{Payload: &pb.LanPacket_NearbyPing{NearbyPing: &pb.LanNearbyPing{
		Instance: t.instance,
		Name:     name,
		Platform: platform,
		Token:    presenceToken(t.id.DeviceID, time.Now()),
	}}})
	if err != nil {
		return
	}
	for _, dst := range targets {
		_ = t.write(dst, pkt)
	}
}

// nearbyAddrsLocked 是某个实例可能在的地址：ping 实际来自的地址在前，mDNS 报的地址在后。
func (t *Transport) nearbyAddrsLocked(id string) []netip.AddrPort {
	si, ok := t.instances[id]
	if !ok {
		return nil
	}
	var out []netip.AddrPort
	if ns := t.nearby[id]; ns != nil {
		for ap := range ns.addrs {
			out = append(out, ap)
		}
		sort.Slice(out, func(i, j int) bool { return ns.addrs[out[i]].After(ns.addrs[out[j]]) })
	}
	if si.inst.Addr.IsValid() && !containsAddr(out, si.inst.Addr) {
		out = append(out, si.inst.Addr)
	}
	if len(out) > maxNearbyAddrs {
		out = out[:maxNearbyAddrs]
	}
	return out
}

// nearbyAdvertLocked 把本机的设备名和系统加进广播（窗口打开期间）。
func (t *Transport) nearbyAdvertLocked(text map[string]string) {
	if t.visible == 0 {
		return
	}
	name, platform := t.nearbySelfLocked()
	text["n"] = encodeNearbyName(name)
	if platform != "" {
		text["o"] = platform
	}
}

// nearbySelfLocked 是本机在「搜索附近设备」里显示的名字和系统。
func (t *Transport) nearbySelfLocked() (name, platform string) {
	name = t.visibleName
	if name == "" {
		name = "CopySync"
	}
	return truncateUTF8(name, maxNearbyName), t.visiblePlatform
}

func encodeNearbyName(name string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(truncateUTF8(name, maxNearbyName)))
}

// nearbyName 取出广播里的设备名。没有设备名（对方没打开窗口，或是旧版本）时返回 false。
func nearbyName(text map[string]string) (string, bool) {
	enc := text["n"]
	if enc == "" {
		return "", false
	}
	b, err := base64.RawURLEncoding.DecodeString(enc)
	if err != nil || len(b) == 0 || !utf8.Valid(b) {
		return "", false
	}
	return string(b), true
}

func truncateUTF8(s string, max int) string {
	if len(s) <= max {
		return s
	}
	for max > 0 && !utf8.RuneStart(s[max]) {
		max--
	}
	return s[:max]
}

func containsAddr(list []netip.AddrPort, ap netip.AddrPort) bool {
	for _, x := range list {
		if x == ap {
			return true
		}
	}
	return false
}
