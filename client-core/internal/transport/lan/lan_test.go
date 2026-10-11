package lan

import (
	"context"
	"crypto/ed25519"
	"errors"
	netpkg "net"
	"net/netip"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/identity"
	"github.com/baiyuze/copysync/client-core/internal/transport/signed"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// node 是测试里的一台设备：一个 Transport 加上它认识的公钥和收到的东西。
type node struct {
	id *identity.Identity
	t  *Transport

	mu        sync.Mutex
	known     map[string]ed25519.PublicKey
	book      map[string][]string
	accepted  []PeerInfo
	signals   chan *pb.SignalMessage
	reach     chan string
	restarted chan string
}

func newNode(t *testing.T, net *MemoryNetwork) *node {
	t.Helper()
	id, err := identity.LoadOrCreate(filepath.Join(t.TempDir(), "identity.json"))
	if err != nil {
		t.Fatal(err)
	}
	n := &node{
		id:        id,
		known:     map[string]ed25519.PublicKey{},
		book:      map[string][]string{},
		signals:   make(chan *pb.SignalMessage, 16),
		reach:     make(chan string, 16),
		restarted: make(chan string, 16),
	}
	var d Discovery
	if net != nil {
		d = net.Discovery()
	} else {
		d = NewMemoryNetwork().Discovery() // 自成一个网络：谁也发现不了它
	}
	n.t = New(Options{
		Identity:   id,
		Port:       -1,
		Discovery:  d,
		Verifier:   signed.NewVerifier(id.DeviceID, n.keyOf),
		Book:       n,
		Peers:      n.peers,
		DeviceInfo: func() (string, string) { return "node-" + id.DeviceID[:4], "test" },
		OnSignal:   func(_ string, msg *pb.SignalMessage) { n.signals <- msg },
		OnReachable: func(id string, ok bool) {
			if ok {
				n.reach <- id
			}
		},
		OnRestarted: func(id string) { n.restarted <- id },
		AcceptPair: func(p PeerInfo) bool {
			n.mu.Lock()
			defer n.mu.Unlock()
			n.accepted = append(n.accepted, p)
			return true
		},
	})
	if err := n.t.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(n.t.Stop)
	return n
}

func (n *node) keyOf(id string) ed25519.PublicKey {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.known[id]
}

func (n *node) peers() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	var out []string
	for id := range n.known {
		out = append(out, id)
	}
	return out
}

func (n *node) DeviceAddrs(id string) ([]string, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]string(nil), n.book[id]...), nil
}

func (n *node) PutDeviceAddr(id, addr string, _ time.Time) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, a := range n.book[id] {
		if a == addr {
			return nil
		}
	}
	n.book[id] = append(n.book[id], addr)
	return nil
}

// pair 让两台设备互相认识（相当于已经配对过）。
func pair(a, b *node) {
	a.mu.Lock()
	a.known[b.id.DeviceID] = b.id.PublicKey
	a.mu.Unlock()
	b.mu.Lock()
	b.known[a.id.DeviceID] = a.id.PublicKey
	b.mu.Unlock()
	a.t.PeersChanged()
	b.t.PeersChanged()
}

func (n *node) addr() netip.AddrPort {
	return netip.AddrPortFrom(netip.MustParseAddr("127.0.0.1"), uint16(n.t.Port()))
}

func waitReachable(t *testing.T, n *node, peer string) {
	t.Helper()
	deadline := time.After(5 * time.Second)
	for !n.t.Reachable(peer) {
		select {
		case <-n.reach:
		case <-time.After(50 * time.Millisecond):
		case <-deadline:
			t.Fatalf("%s 等不到 %s 可达", n.id.DeviceID[:6], peer[:6])
		}
	}
}

// 同一网络里的两台已配对设备靠 mDNS 互相找到，之后信令直接在局域网里收发。
func TestPairedDevicesFindEachOtherAndSignal(t *testing.T) {
	net := NewMemoryNetwork()
	a, b := newNode(t, net), newNode(t, net)
	pair(a, b)

	waitReachable(t, a, b.id.DeviceID)
	waitReachable(t, b, a.id.DeviceID)

	s, err := signed.Sign(a.id, b.id.DeviceID, &pb.SignalMessage{
		Payload: &pb.SignalMessage_Ping{Ping: &pb.Ping{SentAtUnix: 42}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := a.t.Send(b.id.DeviceID, s); err != nil {
		t.Fatal(err)
	}
	select {
	case msg := <-b.signals:
		if msg.GetPing().GetSentAtUnix() != 42 {
			t.Fatalf("收到 %v", msg)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("没收到信令")
	}
	// 补发的那一份会被当作重复消息丢掉
	select {
	case msg := <-b.signals:
		t.Fatalf("同一条信令交给上层两次: %v", msg)
	case <-time.After(400 * time.Millisecond):
	}

	if got, _ := a.DeviceAddrs(b.id.DeviceID); len(got) == 0 {
		t.Error("对方的地址应当记进地址簿")
	}
}

// 跨网段时 mDNS 过不去：只要一边知道另一边的地址，打过招呼后双方都可达。
func TestReachableByRememberedAddress(t *testing.T) {
	a, b := newNode(t, nil), newNode(t, nil)
	pair(a, b)
	a.t.Remember(b.id.DeviceID, []string{b.addr().String()})

	waitReachable(t, a, b.id.DeviceID)
	waitReachable(t, b, a.id.DeviceID)
}

// 未配对的设备发来的 hello 不算数。
func TestStrangerIsIgnored(t *testing.T) {
	a, stranger := newNode(t, nil), newNode(t, nil)
	stranger.mu.Lock()
	stranger.known[a.id.DeviceID] = a.id.PublicKey // 陌生人知道 a，但 a 不认识它
	stranger.mu.Unlock()
	stranger.t.Remember(a.id.DeviceID, []string{a.addr().String()})
	stranger.t.PeersChanged()

	time.Sleep(500 * time.Millisecond)
	if a.t.Reachable(stranger.id.DeviceID) {
		t.Fatal("未配对设备不该被当成可达")
	}
}

func TestPairingOverLAN(t *testing.T) {
	net := NewMemoryNetwork()
	shower, joiner := newNode(t, net), newNode(t, net)

	shower.t.OpenPairing("abc-234", time.Now().Add(time.Minute))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	peer, err := joiner.t.Pair(ctx, "ABC234", nil)
	if err != nil {
		t.Fatalf("配对失败: %v", err)
	}
	if peer.ID != shower.id.DeviceID || string(peer.PublicKey) != string(shower.id.PublicKey) {
		t.Fatalf("拿到的不是对方的身份: %+v", peer)
	}
	if len(peer.Addrs) == 0 || peer.Addrs[0] != shower.addr().String() {
		t.Errorf("回复的来源地址应排第一，得到 %v", peer.Addrs)
	}

	shower.mu.Lock()
	accepted := shower.accepted
	shower.mu.Unlock()
	if len(accepted) != 1 || accepted[0].ID != joiner.id.DeviceID {
		t.Fatalf("显示码的一方应当收到对方的身份，得到 %+v", accepted)
	}

	// 码只能用一次
	late := newNode(t, net)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel2()
	if _, err := late.t.Pair(ctx2, "ABC234", []netip.AddrPort{shower.addr()}); !errors.Is(err, ErrCodeClosed) {
		t.Fatalf("用过的码应当失效，得到 %v", err)
	}
}

// 不在同一网段：mDNS 找不到，靠用户输入的地址配对。
func TestPairingByAddress(t *testing.T) {
	shower, joiner := newNode(t, nil), newNode(t, nil)
	shower.t.OpenPairing("K7MP2X", time.Now().Add(time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	if _, err := joiner.t.Pair(ctx, "K7MP2X", nil); !errors.Is(err, ErrNotFound) {
		t.Fatalf("不输地址应当找不到，得到 %v", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel2()
	peer, err := joiner.t.Pair(ctx2, "K7MP2X", []netip.AddrPort{shower.addr()})
	if err != nil || peer.ID != shower.id.DeviceID {
		t.Fatalf("按地址配对失败: %+v, %v", peer, err)
	}
}

func TestWrongCodeLocksAfterTooManyAttempts(t *testing.T) {
	shower, joiner := newNode(t, nil), newNode(t, nil)
	shower.t.OpenPairing("K7MP2X", time.Now().Add(time.Minute))
	try := func(code string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := joiner.t.Pair(ctx, code, []netip.AddrPort{shower.addr()})
		return err
	}

	if err := try("WRONG2"); !errors.Is(err, ErrWrongCode) {
		t.Fatalf("错码: %v", err)
	}
	for i := 0; i < maxPairAttempts; i++ {
		_ = try("WRONG2")
	}
	if err := try("K7MP2X"); !errors.Is(err, ErrCodeClosed) {
		t.Fatalf("错太多次之后，正确的码也应当失效，得到 %v", err)
	}
}

// 可达状态超过对方声明的间隔三倍没刷新就过期。
func TestReachabilityExpires(t *testing.T) {
	net := NewMemoryNetwork()
	a, b := newNode(t, net), newNode(t, net)
	pair(a, b)
	waitReachable(t, a, b.id.DeviceID)

	b.t.Stop()
	a.t.mu.Lock()
	a.t.peers[b.id.DeviceID].lastSeen = time.Now().Add(-time.Hour)
	a.t.mu.Unlock()
	a.t.tick()
	if a.t.Reachable(b.id.DeviceID) {
		t.Fatal("对方长时间没消息，应当变为不可达")
	}
}

func TestResolveAddress(t *testing.T) {
	cases := map[string]string{
		"192.168.1.20":       "192.168.1.20:8788",
		" 192.168.1.20:9000": "192.168.1.20:9000",
	}
	for in, want := range cases {
		got, err := ResolveAddress(context.Background(), in)
		if err != nil || len(got) != 1 || got[0].String() != want {
			t.Errorf("%q → %v, %v；期望 %s", in, got, err, want)
		}
	}
	for _, bad := range []string{"", "192.168.1.20:0", "[::1]:8788"} {
		if _, err := ResolveAddress(context.Background(), bad); !errors.Is(err, ErrBadAddress) {
			t.Errorf("%q 应当报地址格式不对，得到 %v", bad, err)
		}
	}
	if DisplayAddr(netip.MustParseAddrPort("10.0.0.2:8788")) != "10.0.0.2" ||
		DisplayAddr(netip.MustParseAddrPort("10.0.0.2:9000")) != "10.0.0.2:9000" {
		t.Error("默认端口应当省略，其他端口保留")
	}
}

// 广播里的令牌每天都变，不含设备 ID 本身。
func TestPresenceTokenRotatesDaily(t *testing.T) {
	day := time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC)
	a := presenceToken("dev", day)
	if a == presenceToken("dev", day.Add(24*time.Hour)) {
		t.Error("令牌应当每天变化")
	}
	if a != presenceToken("dev", day.Add(5*time.Hour)) {
		t.Error("同一天的令牌应当相同")
	}
	if a == presenceToken("other", day) {
		t.Error("不同设备的令牌应当不同")
	}
}

// 显示配对码的一方看到有人在找人配对（s=1），先往它发一个空包，好让本机的防火墙放行对方的请求。
func TestShowerPokesSeekers(t *testing.T) {
	net := NewMemoryNetwork()
	shower := newNode(t, net)
	shower.t.OpenPairing("K7MP2X", time.Now().Add(time.Minute))

	seeker, err := netpkg.ListenUDP("udp4", &netpkg.UDPAddr{IP: netpkg.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer seeker.Close()
	port := seeker.LocalAddr().(*netpkg.UDPAddr).Port
	d := net.Discovery()
	if err := d.Start(context.Background(), "seeker", port, map[string]string{"v": "1", "s": "1"}, func(Instance) {}); err != nil {
		t.Fatal(err)
	}
	defer d.Stop()

	_ = seeker.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 64)
	n, _, err := seeker.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("没收到敲门的空包: %v", err)
	}
	if string(buf[:n]) != magic {
		t.Fatalf("收到的不是空包: %q", buf[:n])
	}
}

// 正在配对的一方在广播里带 s=1，配对结束后撤下。
func TestJoinerAdvertisesSeeking(t *testing.T) {
	net := NewMemoryNetwork()
	joiner := newNode(t, net)
	watcher := net.Discovery()
	seen := make(chan map[string]string, 16)
	if err := watcher.Start(context.Background(), "watcher", 1, map[string]string{"v": "1"}, func(i Instance) { seen <- i.Text }); err != nil {
		t.Fatal(err)
	}
	defer watcher.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()
	done := make(chan struct{})
	go func() { _, _ = joiner.t.Pair(ctx, "K7MP2X", nil); close(done) }()

	sawSeeking := false
	for !sawSeeking {
		select {
		case txt := <-seen:
			sawSeeking = txt["s"] == "1"
		case <-time.After(2 * time.Second):
			t.Fatal("配对时广播里没有 s=1")
		}
	}
	<-done
	deadline := time.After(2 * time.Second)
	for {
		select {
		case txt := <-seen:
			if txt["s"] == "" {
				return
			}
		case <-deadline:
			t.Fatal("配对结束后 s=1 没有撤下")
		}
	}
}

// 对方重启后带着新的启动标识打招呼：本端要知道旧连接已经失效。
func TestPeerRestartIsNoticed(t *testing.T) {
	net := NewMemoryNetwork()
	a, b := newNode(t, net), newNode(t, net)
	pair(a, b)
	waitReachable(t, a, b.id.DeviceID)

	// b 用同一个身份重新启动：新的 Transport 有新的启动标识
	b.t.Stop()
	b.t = New(Options{
		Identity:   b.id,
		Port:       -1,
		Discovery:  net.Discovery(),
		Verifier:   signed.NewVerifier(b.id.DeviceID, b.keyOf),
		Book:       b,
		Peers:      b.peers,
		DeviceInfo: func() (string, string) { return "b", "test" },
	})
	if err := b.t.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(b.t.Stop)

	select {
	case id := <-a.restarted:
		if id != b.id.DeviceID {
			t.Fatalf("报告重启的设备不对: %s", id)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("对方重启后没有察觉")
	}
}
