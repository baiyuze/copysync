// Package lan 让同一局域网里的设备不经服务器，直接互相发现、配对、交换信令。
//
// 做法：
//   - 每台设备在固定的 UDP 端口（默认 8788）上收发 LanPacket，并用 mDNS 广播自己。
//     广播里只有一个每天都变的令牌，已配对的设备能认出来，别人认不出。
//   - 对每台已配对但还没连通的设备，定期往它可能在的地址发签名的 hello；收到的一方回一条。
//     两边都在发，有状态的防火墙（Windows）会把对方的包当作回包放行，所以不需要入站规则。
//   - 收到对方的有效消息，就记下来源地址：之后的信令发往那里，地址也存进地址簿，
//     跨网段或重启之后靠它重新找到对方。
//
// 信令仍是 transport/signed 签名的 Signed，信任模型与经服务器转发时完全一样。
package lan

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/copysync/client-core/internal/identity"
	"github.com/baiyuze/copysync/client-core/internal/transport/signed"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

const (
	// DefaultPort 是局域网信令的 UDP 端口。被占用时（同一台电脑上两个用户都登录着）改用随机端口，
	// 端口号随 mDNS 广播和 hello 里的地址告诉对方。
	DefaultPort = 8788

	magic     = "CSL1"
	maxPacket = 16 << 10

	helloInterval     = 10 * time.Second // 对方还没连通时
	keepaliveInterval = 30 * time.Second // 已连通时只需低频保活
	tickInterval      = 2 * time.Second
	maxCandidates     = 12
	resendDelay       = 150 * time.Millisecond // 信令补发一次，抵消 UDP 偶尔丢包
	perIPLimit        = 200                    // 每个来源每秒最多处理的包数
	blockedFor        = 5 * time.Minute
	// 新版本刚启动时，系统要先核对已有的授权，实测有十几秒最初的包也会被拦下。
	// 持续被拦这么久才算真的没授权，免得每次升级都闪一下提示
	blockGrace    = 30 * time.Second
	instanceTTL   = 5 * time.Minute
	saveAddrEvery = 10 * time.Minute
	// mDNS 启动失败后的重试间隔。网络刚变化时网卡可能还没拿到地址，这时启动会失败
	discoveryRetryEvery = 5 * time.Second
	// 本地网络权限的探测间隔（只有 macOS 探测）：被拦着时勤一些，用户一允许就能很快撤掉提示。
	// 刚被拦的宽限期里每一轮都探测：新程序启动时的拦截一解除就要重启 mDNS，见 rediscover
	probeEvery        = 10 * time.Minute
	probeEveryBlocked = 20 * time.Second
	discardPort       = 9
)

var (
	ErrDisabled    = errors.New("局域网发现已关闭")
	ErrUnreachable = errors.New("对方不在局域网里")
)

// PeerInfo 是配对时从对方那里拿到的信息。
type PeerInfo struct {
	ID        string
	Name      string
	Platform  string
	PublicKey []byte
	// Addrs 是对方自报的地址加上实际收到回复的地址，确认配对后写进地址簿
	Addrs []string
}

// AddrBook 持久化已配对设备的地址，由 store.Store 实现。
type AddrBook interface {
	DeviceAddrs(deviceID string) ([]string, error)
	PutDeviceAddr(deviceID, addr string, seen time.Time) error
}

type Options struct {
	Identity *identity.Identity
	Logger   *slog.Logger
	// Port 为 0 用 DefaultPort（被占用时退回随机端口），小于 0 直接用随机端口（测试用）
	Port      int
	Discovery Discovery // 为空时用 mDNS
	Verifier  *signed.Verifier
	Book      AddrBook

	// Peers 返回要在局域网里找的已配对设备
	Peers func() []string
	// Connected 报告与某台设备的 P2P 连接是否已建立：已连通的设备只低频保活
	Connected  func(deviceID string) bool
	DeviceInfo func() (name, platform string)

	// OnSignal 收到已验签的信令（hello 除外）
	OnSignal func(from string, msg *pb.SignalMessage)
	// OnReachable 某台设备在局域网里变得可达或不可达
	OnReachable func(deviceID string, reachable bool)
	// OnRestarted 发现某台已配对设备重启过（hello 里的启动标识变了）时调用：旧连接已经失效
	OnRestarted func(deviceID string)
	// AcceptPair 在有人用正确的配对码发来请求时调用。返回 false 表示这个码已经在别处
	// （服务器那边）用掉了，请求会被拒绝。
	AcceptPair func(peer PeerInfo) bool
	// OnChange 在开关、系统拦截状态变化时调用，用来刷新界面上的状态
	OnChange func()
}

type peerState struct {
	addr      netip.AddrPort // 最近一次收到对方有效消息的来源
	lastSeen  time.Time
	ttl       time.Duration
	reachable bool

	reported   []netip.AddrPort             // 对方 hello 里自报的地址
	boot       []byte                       // 对方 hello 里的启动标识
	discovered map[netip.AddrPort]time.Time // mDNS 发现的地址
	book       []netip.AddrPort             // 地址簿里存的地址
	bookLoaded bool

	nextHello time.Time
	lastReply time.Time
	savedAddr netip.AddrPort
	savedAt   time.Time
}

type seenInstance struct {
	inst Instance
	seen time.Time
}

type Transport struct {
	opts      Options
	log       *slog.Logger
	id        *identity.Identity
	discovery Discovery
	instance  string
	boot      []byte // 本次启动的标识，随 hello 发给对方

	mu        sync.Mutex
	running   bool
	conn      *net.UDPConn
	port      int
	cancel    context.CancelFunc
	peers     map[string]*peerState
	instances map[string]seenInstance
	tokens    map[string]string // 令牌 → 已配对设备
	tokensDay string
	selfAddrs []string
	netKey    string
	// deniedSince 是这一轮连续被系统拦下的开始时间（中间有一次发成功就清零），
	// lastDenied 是最近一次被拦的时间；reportedBlocked 是上次告诉上层的状态
	deniedSince     time.Time
	lastDenied      time.Time
	reportedBlocked bool
	rate            map[netip.Addr]int
	rateAt          time.Time

	pairing   *pairingState
	responses map[string]cachedResponse
	waiters   map[string]*pairWaiter
	// seeking 是正在进行的 Pair 调用数：大于零时广播里带 s=1，
	// 显示配对码的一方据此先往本机发一个空包，见 found
	seeking int
	pokedAt map[string]time.Time
	// nextProbe 是下一次探测本地网络权限的时间，probeAddr 是探测的地址，见 probeTarget
	nextProbe time.Time
	probeAddr netip.Addr
	// rediscover：被系统拦过之后又发成功了，下一轮 tick 要重启 mDNS。
	// pion/mdns 启动时在哪块网卡上加入组播失败，就再也不用那块网卡；
	// 被拦期间启动的 mDNS 会漏掉局域网网卡，拦截解除后也不会自己恢复
	rediscover bool
	// discoveryRetry 不为零表示 mDNS 没启动起来，到这个时间再试一次
	discoveryRetry time.Time

	// reach 是各设备可达状态的副本，由叶子锁保护：peers 层会在持有自己的锁时查询它，
	// 不能让它去等 t.mu（t.mu 持有期间可能回调 peers 层）
	reachMu sync.RWMutex
	reach   map[string]bool
}

func New(opts Options) *Transport {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	d := opts.Discovery
	if d == nil {
		d = NewMDNS()
	}
	suffix := make([]byte, 6)
	_, _ = cryptorand.Read(suffix)
	boot := make([]byte, 8)
	_, _ = cryptorand.Read(boot)
	return &Transport{
		opts:      opts,
		log:       log.With("模块", "局域网"),
		id:        opts.Identity,
		discovery: d,
		// 实例名每次启动随机，不带设备名：同一网络里的其他人看不出这是谁的电脑
		instance:  "copysync-" + hex.EncodeToString(suffix),
		boot:      boot,
		peers:     make(map[string]*peerState),
		instances: make(map[string]seenInstance),
		tokens:    make(map[string]string),
		rate:      make(map[netip.Addr]int),
		responses: make(map[string]cachedResponse),
		waiters:   make(map[string]*pairWaiter),
		pokedAt:   make(map[string]time.Time),
		reach:     make(map[string]bool),
	}
}

// Start 开始在局域网里收发。已经在运行时什么也不做。
func (t *Transport) Start() error {
	peers := t.peerIDs()
	t.mu.Lock()
	if t.running {
		t.mu.Unlock()
		return nil
	}
	conn, err := t.listen()
	if err != nil {
		t.mu.Unlock()
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.conn, t.cancel, t.running = conn, cancel, true
	t.port = conn.LocalAddr().(*net.UDPAddr).Port
	t.netKey = networkKey()
	t.refreshSelfLocked()
	t.rebuildTokensLocked(time.Now(), peers)
	t.mu.Unlock()

	// 启动 mDNS 之前先探测一次：新程序刚启动的几秒里系统会拦着发往局域网的包，
	// 这时启动的 mDNS 用不了局域网网卡。先记下被拦，拦截一解除 tick 就重启 mDNS
	t.probe()
	t.startDiscovery(ctx)
	go t.readLoop(conn)
	go t.tickLoop(ctx)
	t.log.Info("局域网直连已开启", "端口", t.Port())
	t.changed()
	return nil
}

func (t *Transport) listen() (*net.UDPConn, error) {
	port := t.opts.Port
	if port == 0 {
		port = DefaultPort
	}
	if port > 0 {
		if conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: port}); err == nil {
			return conn, nil
		} else {
			t.log.Warn("局域网端口被占用，改用随机端口", "端口", port, "err", err)
		}
	}
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		return nil, fmt.Errorf("监听局域网端口: %w", err)
	}
	return conn, nil
}

// Stop 停止收发。所有对端都变为不可达。
func (t *Transport) Stop() {
	t.mu.Lock()
	if !t.running {
		t.mu.Unlock()
		return
	}
	t.running = false
	t.cancel()
	_ = t.conn.Close()
	var lost []string
	for id, st := range t.peers {
		if st.reachable {
			t.setReachableLocked(id, st, false)
			lost = append(lost, id)
		}
	}
	t.instances = make(map[string]seenInstance)
	t.pairing = nil
	t.mu.Unlock()

	t.discovery.Stop()
	for _, id := range lost {
		t.reachable(id, false)
	}
	t.log.Info("局域网直连已关闭")
	t.changed()
}

func (t *Transport) Running() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running
}

// Port 返回实际监听的端口，没在运行时为 0。
func (t *Transport) Port() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.running {
		return 0
	}
	return t.port
}

// Blocked 报告系统是不是拦下了本机发往局域网的包（macOS 的本地网络权限）。
func (t *Transport) Blocked() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running && t.blockedLocked(time.Now())
}

func (t *Transport) blockedLocked(now time.Time) bool {
	return !t.deniedSince.IsZero() && now.Sub(t.deniedSince) >= blockGrace && now.Sub(t.lastDenied) < blockedFor
}

// Reachable 报告某台设备当前能否在局域网里收到本机的信令。
// 可以在持有其他锁时调用。
func (t *Transport) Reachable(deviceID string) bool {
	t.reachMu.RLock()
	defer t.reachMu.RUnlock()
	return t.reach[deviceID]
}

// Send 把一条已签名的信令发到对方最近所在的地址。
func (t *Transport) Send(to string, s *pb.Signed) error {
	t.mu.Lock()
	st := t.peers[to]
	if !t.running || st == nil || !st.reachable {
		t.mu.Unlock()
		return ErrUnreachable
	}
	dst := st.addr
	t.mu.Unlock()

	pkt, err := encode(&pb.LanPacket{Payload: &pb.LanPacket_Signal{Signal: s}})
	if err != nil {
		return err
	}
	err = t.write(dst, pkt)
	time.AfterFunc(resendDelay, func() { _ = t.write(dst, pkt) })
	return err
}

// DisplayAddrs 返回给用户看的本机地址：不在同一网段时，另一台要输入其中一个。
func (t *Transport) DisplayAddrs() []string {
	port := t.Port()
	if port == 0 {
		return nil
	}
	var out []string
	for _, ip := range LocalAddrs() {
		out = append(out, DisplayAddr(netip.AddrPortFrom(ip, uint16(port))))
		if len(out) == 3 {
			break
		}
	}
	return out
}

// PeersChanged 在配对关系变化后调用：重算令牌，马上向新设备打招呼。
func (t *Transport) PeersChanged() {
	peers := t.peerIDs()
	t.mu.Lock()
	t.rebuildTokensLocked(time.Now(), peers)
	for _, st := range t.peers {
		st.bookLoaded = false
		st.nextHello = time.Time{}
	}
	t.mu.Unlock()
	t.tick()
}

// Remember 把地址写进某台设备的地址簿（配对确认后调用）。
func (t *Transport) Remember(deviceID string, addrs []string) {
	now := time.Now()
	for _, a := range addrs {
		ap, err := netip.ParseAddrPort(a)
		if err != nil || !ap.Addr().Is4() {
			continue
		}
		if t.opts.Book != nil {
			if err := t.opts.Book.PutDeviceAddr(deviceID, ap.String(), now); err != nil {
				t.log.Debug("保存地址失败", "device", deviceID, "addr", ap, "err", err)
			}
		}
	}
	t.mu.Lock()
	st := t.peerLocked(deviceID)
	st.bookLoaded = false
	st.nextHello = time.Time{}
	t.mu.Unlock()
}

// SetPeerAddr 记下用户手动输入的地址，并立即往那里打招呼。
func (t *Transport) SetPeerAddr(ctx context.Context, deviceID, addr string) error {
	aps, err := ResolveAddress(ctx, addr)
	if err != nil {
		return err
	}
	var list []string
	for _, ap := range aps {
		list = append(list, ap.String())
	}
	t.Remember(deviceID, list)
	t.sendHello(deviceID, aps, true, helloInterval)
	return nil
}

// Forget 在解除配对后丢掉这台设备的状态（地址簿随设备记录一起由 store 删除）。
func (t *Transport) Forget(deviceID string) {
	peers := t.peerIDs()
	t.mu.Lock()
	if st := t.peers[deviceID]; st != nil {
		t.setReachableLocked(deviceID, st, false)
	}
	delete(t.peers, deviceID)
	t.rebuildTokensLocked(time.Now(), peers)
	t.mu.Unlock()
}

// PeerStatus 是某台已配对设备在局域网里的状态，供诊断展示。
type PeerStatus struct {
	ID        string
	Reachable bool
	Addr      string // 最近一次收到对方消息的地址
	Known     []string
	LastSeen  time.Time
}

// Info 是局域网直连的诊断信息，见 copysync-cli nat。
type Info struct {
	Running    bool
	Port       int
	LocalAddrs []string
	Discovered int // mDNS 发现的实例数，含未配对的
	Blocked    bool
	Peers      []PeerStatus
}

func (t *Transport) Info() Info {
	peers := t.peerIDs()
	blocked := t.Blocked()
	t.mu.Lock()
	info := Info{
		Running:    t.running,
		Port:       t.port,
		LocalAddrs: append([]string(nil), t.selfAddrs...),
		Discovered: len(t.instances),
		Blocked:    blocked,
	}
	if !t.running {
		info.Port = 0
	}
	for _, id := range peers {
		ps := PeerStatus{ID: id}
		if st := t.peers[id]; st != nil {
			ps.Reachable = st.reachable
			ps.LastSeen = st.lastSeen
			if st.addr.IsValid() {
				ps.Addr = st.addr.String()
			}
		}
		info.Peers = append(info.Peers, ps)
	}
	t.mu.Unlock()

	if t.opts.Book != nil {
		for i := range info.Peers {
			info.Peers[i].Known, _ = t.opts.Book.DeviceAddrs(info.Peers[i].ID)
		}
	}
	return info
}

// ─────────────────────────── 接收 ───────────────────────────

func (t *Transport) readLoop(conn *net.UDPConn) {
	buf := make([]byte, 64<<10)
	for {
		n, src, err := conn.ReadFromUDPAddrPort(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			// Windows 上对端端口不可达时，下一次读会收到 ICMP 引起的错误，忽略即可
			continue
		}
		src = netip.AddrPortFrom(src.Addr().Unmap(), src.Port())
		if n < len(magic) || n > maxPacket || string(buf[:len(magic)]) != magic || !t.allow(src.Addr()) {
			continue
		}
		var pkt pb.LanPacket
		if err := proto.Unmarshal(buf[len(magic):n], &pkt); err != nil {
			continue
		}
		t.handle(&pkt, src)
	}
}

// allow 按来源限速，挡住同一网络里的刷包。
func (t *Transport) allow(ip netip.Addr) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	if now.Sub(t.rateAt) >= time.Second {
		t.rateAt = now
		clear(t.rate)
	}
	t.rate[ip]++
	return t.rate[ip] <= perIPLimit
}

func (t *Transport) handle(pkt *pb.LanPacket, src netip.AddrPort) {
	switch p := pkt.GetPayload().(type) {
	case *pb.LanPacket_Signal:
		from, msg, err := t.opts.Verifier.Verify(p.Signal)
		switch {
		case errors.Is(err, signed.ErrReplay), errors.Is(err, signed.ErrUnknownSender):
			// 重复的包（同时发往对方的几个地址）、或是解除配对后对方还在打招呼：都正常
			return
		case err != nil:
			t.log.Debug("丢弃局域网信令", "from", from, "src", src, "原因", err)
			return
		}
		hello := msg.GetLanHello()
		t.sawPeer(from, src, hello)
		if hello != nil {
			t.answerHello(from, src, hello)
			return
		}
		if t.opts.OnSignal != nil {
			t.opts.OnSignal(from, msg)
		}

	case *pb.LanPacket_PairRequest:
		t.handlePairRequest(p.PairRequest, src)
	case *pb.LanPacket_PairResponse:
		t.handlePairResponse(p.PairResponse, src)
	case *pb.LanPacket_PairReject:
		t.handlePairReject(p.PairReject)
	}
}

// sawPeer 记下某台设备的一条有效消息：它在 src，此刻可达。
func (t *Transport) sawPeer(id string, src netip.AddrPort, hello *pb.LanHello) {
	now := time.Now()
	t.mu.Lock()
	if !t.running {
		t.mu.Unlock()
		return
	}
	st := t.peerLocked(id)
	st.addr, st.lastSeen = src, now
	restarted := false
	if hello != nil {
		st.ttl = ttlFor(hello.GetIntervalSeconds())
		st.reported = parseAddrs(hello.GetAddresses())
		if b := hello.GetBoot(); len(b) > 0 {
			restarted = len(st.boot) > 0 && !bytes.Equal(st.boot, b)
			st.boot = append(st.boot[:0], b...)
		}
	} else if st.ttl == 0 {
		st.ttl = ttlFor(0)
	}
	became := !st.reachable
	t.setReachableLocked(id, st, true)
	save := src != st.savedAddr || now.Sub(st.savedAt) > saveAddrEvery
	if save {
		st.savedAddr, st.savedAt = src, now
	}
	t.mu.Unlock()

	if save && t.opts.Book != nil {
		// 正在配对、尚未确认的设备还没落库，这里会因外键失败，确认时再存
		_ = t.opts.Book.PutDeviceAddr(id, src.String(), now)
	}
	if became {
		t.log.Info("局域网里找到了设备", "device", id, "地址", src)
		t.reachable(id, true)
	}
	if restarted {
		t.log.Info("局域网里的设备重启过", "device", id)
		if t.opts.OnRestarted != nil {
			t.opts.OnRestarted(id)
		}
	}
}

func (t *Transport) answerHello(from string, src netip.AddrPort, hello *pb.LanHello) {
	if !hello.GetReplyWanted() {
		return
	}
	t.mu.Lock()
	st := t.peerLocked(from)
	// 对方往几个地址各发了一份时只回一次
	if time.Since(st.lastReply) < time.Second {
		t.mu.Unlock()
		return
	}
	st.lastReply = time.Now()
	t.mu.Unlock()
	t.sendHello(from, []netip.AddrPort{src}, false, t.intervalFor(from))
}

func ttlFor(intervalSeconds uint32) time.Duration {
	if intervalSeconds == 0 {
		intervalSeconds = uint32(helloInterval / time.Second)
	}
	return 3*time.Duration(intervalSeconds)*time.Second + 5*time.Second
}

// ─────────────────────────── 定时任务 ───────────────────────────

func (t *Transport) tickLoop(ctx context.Context) {
	tk := time.NewTicker(tickInterval)
	defer tk.Stop()
	t.tick()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			t.tick()
		}
	}
}

// tick 处理到期的事：网络变化、令牌换天、可达状态过期、该打招呼的设备。
func (t *Transport) tick() {
	now := time.Now()
	peers := t.peerIDs()
	// 回调 P2P 层要在锁外：它发信令时会反过来调用 Send
	intervals := make(map[string]time.Duration, len(peers))
	for _, id := range peers {
		intervals[id] = t.intervalFor(id)
	}

	t.mu.Lock()
	if !t.running {
		t.mu.Unlock()
		return
	}
	netChanged := false
	if k := networkKey(); k != t.netKey {
		t.netKey = k
		netChanged = true
		t.refreshSelfLocked()
	}
	dayChanged := now.UTC().Format("2006-01-02") != t.tokensDay
	if dayChanged {
		t.rebuildTokensLocked(now, peers)
	}
	probe := !now.Before(t.nextProbe)
	rediscover := t.rediscover
	t.rediscover = false
	retryDiscovery := !t.discoveryRetry.IsZero() && !now.Before(t.discoveryRetry)
	blocked := t.blockedLocked(now)
	blockedChanged := blocked != t.reportedBlocked
	t.reportedBlocked = blocked
	pairingExpired := t.pairing != nil && now.After(t.pairing.expires)
	if pairingExpired {
		t.pairing = nil
	}
	for k, r := range t.responses {
		if now.Sub(r.at) > 2*time.Minute {
			delete(t.responses, k)
		}
	}
	for k, at := range t.pokedAt {
		if now.Sub(at) > instanceTTL {
			delete(t.pokedAt, k)
		}
	}
	for k, si := range t.instances {
		if now.Sub(si.seen) > instanceTTL {
			delete(t.instances, k)
		}
	}

	var lost []string
	for id, st := range t.peers {
		if st.reachable && now.Sub(st.lastSeen) > st.ttl {
			t.setReachableLocked(id, st, false)
			lost = append(lost, id)
		}
	}

	type due struct {
		id       string
		targets  []netip.AddrPort
		interval time.Duration
	}
	var hellos []due
	for _, id := range peers {
		st := t.peerLocked(id)
		if netChanged {
			st.nextHello = time.Time{}
		}
		if now.Before(st.nextHello) {
			continue
		}
		interval := intervals[id]
		st.nextHello = now.Add(interval)
		if targets := t.candidatesLocked(id, st, now); len(targets) > 0 {
			hellos = append(hellos, due{id, targets, interval})
		}
	}
	ctx := context.Background()
	t.mu.Unlock()

	for _, id := range lost {
		t.log.Info("局域网里的设备不再可达", "device", id)
		t.reachable(id, false)
	}
	if blockedChanged {
		if blocked {
			t.log.Warn("系统拦下了发往局域网的包，需要在系统设置里允许本地网络访问")
		} else {
			t.log.Info("系统不再拦着发往局域网的包")
		}
		t.changed()
	}
	switch {
	case netChanged:
		// 网卡变了：mDNS 要在新网卡上重新加入组播
		t.log.Info("网络已变化，重新开始局域网发现")
		t.discovery.Stop()
		t.startDiscovery(ctx)
	case rediscover:
		t.log.Info("系统放行了发往局域网的包，重新开始局域网发现")
		t.discovery.Stop()
		t.startDiscovery(ctx)
	case retryDiscovery:
		t.discovery.Stop()
		t.startDiscovery(ctx)
	case dayChanged || pairingExpired:
		t.updateAdvert()
	}
	for _, h := range hellos {
		t.sendHello(h.id, h.targets, true, h.interval)
	}
	if probe {
		t.probe()
	}
}

// probe 往本网段里的一个地址发几个字节，看系统是否拦着（只有 macOS 有目标地址），
// 并排好下一次探测：被拦时 write 会把它提前。
func (t *Transport) probe() {
	ip, ok := probeTarget()
	t.mu.Lock()
	t.nextProbe = time.Now().Add(probeEvery)
	t.probeAddr = ip // 没有可探测的地址时是零值
	t.mu.Unlock()
	if ok {
		_ = t.send(netip.AddrPortFrom(ip, discardPort), []byte(magic), true)
	}
}

// intervalFor 返回向某台设备打招呼的间隔。不能在持有 t.mu 时调用。
func (t *Transport) intervalFor(id string) time.Duration {
	if t.opts.Connected != nil && t.opts.Connected(id) {
		return keepaliveInterval
	}
	return helloInterval
}

// candidatesLocked 汇总某台设备可能在的地址：最近的来源、对方自报的、mDNS 发现的、地址簿里的。
func (t *Transport) candidatesLocked(id string, st *peerState, now time.Time) []netip.AddrPort {
	if !st.bookLoaded && t.opts.Book != nil {
		st.bookLoaded = true
		st.book = nil
		if list, err := t.opts.Book.DeviceAddrs(id); err == nil {
			st.book = parseAddrs(list)
		}
	}
	var out []netip.AddrPort
	seen := map[netip.AddrPort]bool{}
	add := func(ap netip.AddrPort) {
		if ap.IsValid() && ap.Addr().Is4() && !seen[ap] && len(out) < maxCandidates {
			seen[ap] = true
			out = append(out, ap)
		}
	}
	add(st.addr)
	for ap, at := range st.discovered {
		if now.Sub(at) < instanceTTL {
			add(ap)
		}
	}
	for _, ap := range st.reported {
		add(ap)
	}
	for _, ap := range st.book {
		add(ap)
	}
	return out
}

func (t *Transport) sendHello(to string, targets []netip.AddrPort, replyWanted bool, interval time.Duration) {
	t.mu.Lock()
	addrs := append([]string(nil), t.selfAddrs...)
	t.mu.Unlock()

	s, err := signed.Sign(t.id, to, &pb.SignalMessage{Payload: &pb.SignalMessage_LanHello{LanHello: &pb.LanHello{
		Addresses:       addrs,
		ReplyWanted:     replyWanted,
		IntervalSeconds: uint32(interval / time.Second),
		Boot:            t.boot,
	}}})
	if err != nil {
		return
	}
	pkt, err := encode(&pb.LanPacket{Payload: &pb.LanPacket_Signal{Signal: s}})
	if err != nil {
		return
	}
	for _, dst := range targets {
		_ = t.write(dst, pkt)
	}
}

// ─────────────────────────── 发现 ───────────────────────────

func (t *Transport) startDiscovery(ctx context.Context) {
	t.mu.Lock()
	port := t.port
	text := t.advertTextLocked()
	t.mu.Unlock()
	err := t.discovery.Start(ctx, t.instance, port, text, t.found)

	t.mu.Lock()
	failedBefore := !t.discoveryRetry.IsZero()
	if err != nil {
		t.discoveryRetry = time.Now().Add(discoveryRetryEvery)
	} else {
		t.discoveryRetry = time.Time{}
	}
	t.mu.Unlock()
	switch {
	case err != nil && !failedBefore:
		// 没有 mDNS 也能用：输入地址、地址簿里的地址照样能连。过一会儿再试
		t.log.Warn("局域网发现启动失败，稍后重试", "err", err)
	case err != nil:
		t.log.Debug("局域网发现仍然启动失败", "err", err)
	case failedBefore:
		t.log.Info("局域网发现已启动")
	}
}

func (t *Transport) updateAdvert() {
	t.mu.Lock()
	if !t.running {
		t.mu.Unlock()
		return
	}
	text := t.advertTextLocked()
	t.mu.Unlock()
	if err := t.discovery.Update(text); err != nil {
		t.log.Debug("更新局域网广播失败", "err", err)
	}
}

func (t *Transport) advertTextLocked() map[string]string {
	text := map[string]string{
		"v": "1",
		"t": presenceToken(t.id.DeviceID, time.Now()),
	}
	if t.pairing != nil {
		text["p"] = "1"
	}
	if t.seeking > 0 {
		text["s"] = "1"
	}
	return text
}

// found 处理 mDNS 发现的实例：认出已配对的设备就马上打招呼。
func (t *Transport) found(inst Instance) {
	if inst.Text["v"] != "1" || !inst.Addr.IsValid() {
		return
	}
	now := time.Now()
	t.mu.Lock()
	if !t.running {
		t.mu.Unlock()
		return
	}
	prev, had := t.instances[inst.Name]
	t.instances[inst.Name] = seenInstance{inst: inst, seen: now}
	// pion 每秒都会报一次同样的结果，只在有新信息时行动
	fresh := !had || prev.inst.Addr != inst.Addr || prev.inst.Text["t"] != inst.Text["t"] ||
		now.Sub(prev.seen) > 30*time.Second

	// 本机在显示配对码，对方在找人配对：先往它发一个空包。本机若是没放行防火墙的 Windows，
	// 对方主动发来的配对请求会被挡住；本机先发过包，有状态的防火墙就会把它当回包放行
	// 同一个实例可能有好几个地址（多块网卡），每个地址都要敲到
	var poke netip.AddrPort
	pokeKey := inst.Name + "|" + inst.Addr.String()
	if t.pairing != nil && inst.Text["s"] == "1" && now.Sub(t.pokedAt[pokeKey]) > 2*time.Second {
		t.pokedAt[pokeKey] = now
		poke = inst.Addr
	}
	if fresh && (inst.Text["p"] == "1" || inst.Text["s"] == "1") {
		t.log.Debug("发现正在配对的实例", "实例", inst.Name, "地址", inst.Addr, "开放配对码", inst.Text["p"] == "1", "在找人配对", inst.Text["s"] == "1")
	}

	var hello string
	if id := t.tokens[inst.Text["t"]]; id != "" {
		st := t.peerLocked(id)
		if st.discovered == nil {
			st.discovered = make(map[netip.AddrPort]time.Time)
		}
		st.discovered[inst.Addr] = now
		// 实例名或地址变了，说明对方重启过或换了网络：即使还当它可达，也马上打招呼。
		// 对方若是没放行防火墙的 Windows，它收不到本机主动发去的包，得靠本机先回应它
		restarted := !had || prev.inst.Addr != inst.Addr
		if fresh && (!st.reachable || restarted) {
			hello = id
		}
	}
	t.mu.Unlock()

	if poke.IsValid() {
		t.log.Debug("向在找人配对的实例敲门", "地址", poke)
		_ = t.write(poke, []byte(magic))
	}
	if hello != "" {
		t.sendHello(hello, []netip.AddrPort{inst.Addr}, true, t.intervalFor(hello))
	}
}

// peerIDs 向上层要已配对设备的列表。不能在持有 t.mu 时调用。
func (t *Transport) peerIDs() []string {
	if t.opts.Peers == nil {
		return nil
	}
	return t.opts.Peers()
}

func (t *Transport) rebuildTokensLocked(now time.Time, peers []string) {
	t.tokens = make(map[string]string)
	t.tokensDay = now.UTC().Format("2006-01-02")
	for _, id := range peers {
		// 两边的时钟可能不在同一天，前后各算一天
		for _, d := range []time.Duration{-24 * time.Hour, 0, 24 * time.Hour} {
			t.tokens[presenceToken(id, now.Add(d))] = id
		}
	}
	// 配对之前就发现了的实例，现在可能认得出来了
	for _, si := range t.instances {
		if id := t.tokens[si.inst.Text["t"]]; id != "" {
			st := t.peerLocked(id)
			if st.discovered == nil {
				st.discovered = make(map[netip.AddrPort]time.Time)
			}
			st.discovered[si.inst.Addr] = si.seen
		}
	}
}

// ─────────────────────────── 内部辅助 ───────────────────────────

func (t *Transport) setReachableLocked(id string, st *peerState, ok bool) {
	st.reachable = ok
	t.reachMu.Lock()
	if ok {
		t.reach[id] = true
	} else {
		delete(t.reach, id)
	}
	t.reachMu.Unlock()
}

func (t *Transport) peerLocked(id string) *peerState {
	st := t.peers[id]
	if st == nil {
		st = &peerState{}
		t.peers[id] = st
	}
	return st
}

func (t *Transport) refreshSelfLocked() {
	t.selfAddrs = t.selfAddrs[:0]
	for _, ip := range LocalAddrs() {
		t.selfAddrs = append(t.selfAddrs, netip.AddrPortFrom(ip, uint16(t.port)).String())
	}
}

// write 发一个包。macOS 拒绝本地网络访问时记下来，界面据此提示用户去授权。
func (t *Transport) write(dst netip.AddrPort, pkt []byte) error { return t.send(dst, pkt, false) }

// send 发一个包，probe 表示这是权限探测包。
//
// 判断拦截解除只信探测包：被拦着时，发给先联系过本机的设备照样能成功（算作回包），
// 拿它当证据会让提示时有时无、mDNS 反复重启。找不到探测地址时才退回认任何一次成功。
func (t *Transport) send(dst netip.AddrPort, pkt []byte, probe bool) error {
	t.mu.Lock()
	conn := t.conn
	running := t.running
	t.mu.Unlock()
	if !running {
		return ErrDisabled
	}
	_, err := conn.WriteToUDPAddrPort(pkt, dst)
	if onLink(dst.Addr()) {
		t.noteSend(err, probe, time.Now())
	}
	return err
}

// noteSend 记下一次发往直连网段的结果。只记下来，是否算「被拦着」由 tick 按宽限期判断并通知上层。
func (t *Transport) noteSend(err error, probe bool, now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	switch {
	case err != nil && deniedByOS(err):
		if t.deniedSince.IsZero() {
			t.deniedSince = now
		}
		t.lastDenied = now
		next := now.Add(probeEveryBlocked)
		if now.Sub(t.deniedSince) < blockGrace {
			next = now // 下一轮 tick 就再探测
		}
		if t.nextProbe.After(next) {
			t.nextProbe = next
		}
	case err == nil && (probe || !t.probeAddr.IsValid()):
		if !t.deniedSince.IsZero() {
			t.rediscover = true
		}
		t.deniedSince = time.Time{}
	}
}

func (t *Transport) reachable(id string, ok bool) {
	if t.opts.OnReachable != nil {
		t.opts.OnReachable(id, ok)
	}
}

func (t *Transport) changed() {
	if t.opts.OnChange != nil {
		t.opts.OnChange()
	}
}

func encode(pkt *pb.LanPacket) ([]byte, error) {
	body, err := proto.Marshal(pkt)
	if err != nil {
		return nil, err
	}
	if len(magic)+len(body) > maxPacket {
		return nil, fmt.Errorf("局域网消息过大: %d 字节", len(body))
	}
	return append([]byte(magic), body...), nil
}

func parseAddrs(list []string) []netip.AddrPort {
	var out []netip.AddrPort
	for _, s := range list {
		if ap, err := netip.ParseAddrPort(s); err == nil && ap.Addr().Is4() {
			out = append(out, ap)
		}
	}
	return out
}
