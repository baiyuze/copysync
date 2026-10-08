package p2p

import (
	"context"
	"log/slog"
	"net"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/pion/ice/v4"
	"github.com/pion/webrtc/v4"

	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// Manager 管理与所有已配对设备的 P2P 连接。
type Manager struct {
	selfID string
	log    *slog.Logger

	sendSignal func(ctx context.Context, to string, msg *pb.SignalMessage) error
	onState    func(deviceID string, state ConnState)
	onMessage  func(deviceID string, msg *pb.PeerMessage)
	onStream   func(deviceID, label string, stream *Stream)

	mu         sync.RWMutex
	peers      map[string]*Peer
	iceServers []webrtc.ICEServer
	forceRelay bool

	// 打洞共用的那一个本地 UDP 端口，所有连接都从它收发。
	// 只有共用一个端口，各台探测服务器问到的公网地址才对同一个端口成立（见 egress.go）。
	mux          *ice.UniversalUDPMuxDefault
	localPort    int
	probeOpts    probeOptions
	probeServers func() []string
	history      *egressHistory
	stop         chan struct{}
	closeOnce    sync.Once

	egressMu   sync.Mutex
	egress     EgressReport
	serverSTUN []string      // 信令服务器下发的 STUN 地址，只用于出口探测
	probing    chan struct{} // 非空表示正在探测，探测结束时关闭
	probeAgain bool          // 探测期间又有新的触发，结束后再探一轮
}

type ManagerOptions struct {
	SelfID     string
	Logger     *slog.Logger
	SendSignal func(ctx context.Context, to string, msg *pb.SignalMessage) error
	OnState    func(deviceID string, state ConnState)
	OnMessage  func(deviceID string, msg *pb.PeerMessage)
	OnStream   func(deviceID, label string, stream *Stream)
	// ForceRelay 强制走 TURN 中转，用于验证中转链路
	ForceRelay bool

	// StatePath 保存出口历史的文件；为空则只记在内存里
	StatePath string
	// ProbeServers 返回额外用于出口探测的公共 STUN 服务器；用户关闭时返回空
	ProbeServers func() []string
	// ProbeAllowPrivate 允许内网与回环地址的探测服务器，仅供测试与 NAT 实验室使用
	ProbeAllowPrivate bool
	// DisableEgressProbe 退回 pion 默认的候选收集方式：每个 STUN 服务器各开一个端口，
	// 不做多出口探测。NAT 实验室用它复现旧版本的行为作对照
	DisableEgressProbe bool
}

func NewManager(opts ManagerOptions) *Manager {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	m := &Manager{
		selfID:       opts.SelfID,
		log:          log,
		sendSignal:   opts.SendSignal,
		onState:      opts.OnState,
		onMessage:    opts.OnMessage,
		onStream:     opts.OnStream,
		peers:        make(map[string]*Peer),
		forceRelay:   opts.ForceRelay,
		probeOpts:    probeOptions{AllowPrivate: opts.ProbeAllowPrivate},
		probeServers: opts.ProbeServers,
		history:      loadEgressHistory(opts.StatePath),
		stop:         make(chan struct{}),
	}

	if opts.DisableEgressProbe {
		return m
	}
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		// 开不出共享端口就退回 pion 的默认做法：功能不受影响，只是多出口网络里难以直连
		log.Warn("无法打开打洞端口，退回默认的候选收集方式", "err", err)
		return m
	}
	m.mux = ice.NewUniversalUDPMuxDefault(ice.UniversalUDPMuxParams{UDPConn: conn})
	m.localPort = conn.LocalAddr().(*net.UDPAddr).Port

	go m.egressLoop()
	m.triggerProbe()
	return m
}

// SetICEServers 更新 STUN/TURN 配置。信令服务器在握手时下发，
// TURN 凭证是短期的，需要定期刷新。
//
// STUN 地址不交给 pion：pion 每问一个 STUN 服务器都单独开一个端口，问到的地址
// 只对那个端口成立。改由本端在共享端口上探测（见 egress.go），TURN 照常交给 pion。
func (m *Manager) SetICEServers(servers []*pb.IceServer) {
	converted := make([]webrtc.ICEServer, 0, len(servers))
	var stunURLs []string
	for _, s := range servers {
		urls := s.GetUrls()
		if m.mux != nil {
			var turnURLs []string
			for _, u := range urls {
				if strings.HasPrefix(u, "stun:") || strings.HasPrefix(u, "stuns:") {
					stunURLs = append(stunURLs, u)
				} else {
					turnURLs = append(turnURLs, u)
				}
			}
			if len(turnURLs) == 0 {
				continue
			}
			urls = turnURLs
		}
		converted = append(converted, webrtc.ICEServer{
			URLs:       urls,
			Username:   s.GetUsername(),
			Credential: s.GetCredential(),
		})
	}

	m.mu.Lock()
	m.iceServers = converted
	m.mu.Unlock()

	m.egressMu.Lock()
	changed := !slices.Equal(m.serverSTUN, stunURLs)
	m.serverSTUN = stunURLs
	m.egressMu.Unlock()
	if changed {
		m.triggerProbe()
	}

	m.log.Debug("ICE 配置已更新", "servers", len(converted), "stun", len(stunURLs))
}

// Close 断开所有连接并释放打洞端口。
func (m *Manager) Close() {
	m.CloseAll()
	m.closeOnce.Do(func() {
		close(m.stop)
		if m.mux != nil {
			_ = m.mux.Close()
		}
	})
}

// Reprobe 立即重新探测出口。从睡眠中唤醒时调用：网络多半已经变了，不必等下一轮。
func (m *Manager) Reprobe() { m.triggerProbe() }

// Egress 返回最近一轮出口探测的结果（含按历史推算的出口），供诊断展示。
func (m *Manager) Egress() EgressReport {
	m.egressMu.Lock()
	defer m.egressMu.Unlock()
	return m.egress
}

// egressLoop 定期重新探测；本机所在网络变化时立即探测。
func (m *Manager) egressLoop() {
	periodic := time.NewTicker(5 * time.Minute)
	netCheck := time.NewTicker(15 * time.Second)
	defer periodic.Stop()
	defer netCheck.Stop()

	lastNet := localNetworkKey()
	for {
		select {
		case <-m.stop:
			return
		case <-periodic.C:
			m.triggerProbe()
		case <-netCheck.C:
			if k := localNetworkKey(); k != lastNet {
				lastNet = k
				m.log.Info("网络已变化，重新探测出口")
				m.triggerProbe()
			}
		}
	}
}

// triggerProbe 发起一轮出口探测。正在探测时只做标记，结束后再补一轮，
// 不会并发探测。
func (m *Manager) triggerProbe() {
	if m.mux == nil {
		return
	}
	m.egressMu.Lock()
	if m.probing != nil {
		m.probeAgain = true
		m.egressMu.Unlock()
		return
	}
	done := make(chan struct{})
	m.probing = done
	m.egressMu.Unlock()

	go m.runProbe(done)
}

func (m *Manager) runProbe(done chan struct{}) {
	for {
		m.egressMu.Lock()
		servers := slices.Clone(m.serverSTUN)
		m.egressMu.Unlock()
		if m.probeServers != nil {
			servers = append(servers, m.probeServers()...)
		}

		report := EgressReport{LocalPort: m.localPort, At: time.Now()}
		if len(servers) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			report = probeEgress(ctx, m.mux, m.localPort, servers, m.probeOpts)
			cancel()
		}
		network := localNetworkKey()
		m.history.record(network, report)
		report = m.history.augment(network, report)

		var summary []string
		for _, e := range report.Egresses {
			s := e.Addr.String()
			switch {
			case e.Guessed():
				s += "（按历史推算）"
			case e.PortPreserved:
				s += "（不改端口）"
			}
			summary = append(summary, s)
		}
		if len(servers) > 0 {
			m.log.Info("出口探测完成", "本地端口", m.localPort, "出口", summary,
				"响应", report.Answered, "探测", report.Probed)
		}

		m.egressMu.Lock()
		m.egress = report
		if m.probeAgain {
			m.probeAgain = false
			m.egressMu.Unlock()
			continue
		}
		m.probing = nil
		m.egressMu.Unlock()
		close(done)
		return
	}
}

// refreshEgress 立即重新探测一轮出口并等它结束（最多 5 秒）。
// 尝试从中转换成直连之前调用：公司网络里每条线路分到哪些目标会随时间变化，
// 上一轮没探测到的线路这一轮可能就有了。
func (m *Manager) refreshEgress() {
	m.triggerProbe()
	m.egressMu.Lock()
	wait := m.probing
	m.egressMu.Unlock()
	if wait != nil {
		select {
		case <-wait:
		case <-time.After(5 * time.Second):
		}
	}
}

// advertisedCandidates 返回要额外发给对端的出口地址。
//
// 正在探测时最多等 4 秒，免得刚连上信令就用一份空的结果去握手。这是在发出 offer/answer
// 之后的后台里等，不耽误握手；但要短于中转的等待时限（relayAcceptanceWait），
// 出口地址到了对端还来得及打洞。
func (m *Manager) advertisedCandidates() []string {
	if m.mux == nil {
		return nil
	}
	m.egressMu.Lock()
	wait := m.probing
	m.egressMu.Unlock()
	if wait != nil {
		select {
		case <-wait:
		case <-time.After(4 * time.Second):
		}
	}
	return srflxCandidates(m.Egress())
}

// Connect 与对端建立连接。双方都可以调用，不会因此产生 glare。
//
// 由 device_id 字典序较小的一方发起 offer，另一方只准备好接收。
// 这个约定两端各自可算出、无需协商，从根上避免了双方同时发 offer 的竞争——
// 比事后用 rollback 去化解冲突简单得多，也更不容易出错。
func (m *Manager) Connect(ctx context.Context, deviceID string) error {
	p, err := m.peer(deviceID)
	if err != nil {
		return err
	}
	if m.selfID > deviceID {
		// 本端负责应答：Peer 对象已就绪，等对方的 offer 到来即可
		m.log.Debug("等待对端发起连接", "peer", deviceID)
		return nil
	}
	return p.Connect(ctx)
}

// HandleSignal 把验签后的信令分派给对应的 Peer。
func (m *Manager) HandleSignal(ctx context.Context, from string, msg *pb.SignalMessage) error {
	p, err := m.peer(from)
	if err != nil {
		return err
	}
	return p.HandleSignal(ctx, msg)
}

// Send 向对端发送控制消息。
func (m *Manager) Send(deviceID string, msg *pb.PeerMessage) error {
	m.mu.RLock()
	p := m.peers[deviceID]
	m.mu.RUnlock()

	if p == nil {
		return ErrNotConnected
	}
	return p.Send(msg)
}

// Broadcast 把消息发给所有已连接的对端，返回成功送达的数量。
func (m *Manager) Broadcast(msg *pb.PeerMessage) int {
	m.mu.RLock()
	peers := make([]*Peer, 0, len(m.peers))
	for _, p := range m.peers {
		peers = append(peers, p)
	}
	m.mu.RUnlock()

	var sent int
	for _, p := range peers {
		if err := p.Send(msg); err == nil {
			sent++
		}
	}
	return sent
}

func (m *Manager) OpenStream(deviceID, label string) (*webrtc.DataChannel, error) {
	m.mu.RLock()
	p := m.peers[deviceID]
	m.mu.RUnlock()

	if p == nil {
		return nil, ErrNotConnected
	}
	return p.OpenStream(label)
}

func (m *Manager) State(deviceID string) ConnState {
	m.mu.RLock()
	p := m.peers[deviceID]
	m.mu.RUnlock()

	if p == nil {
		return StateOffline
	}
	return p.State()
}

// Disconnect 断开并丢弃与某个对端的连接（解除配对或对端离线时）。
func (m *Manager) Disconnect(deviceID string) {
	m.mu.Lock()
	p := m.peers[deviceID]
	delete(m.peers, deviceID)
	m.mu.Unlock()

	if p != nil {
		p.Close()
	}
}

func (m *Manager) CloseAll() {
	m.mu.Lock()
	peers := m.peers
	m.peers = make(map[string]*Peer)
	m.mu.Unlock()

	for _, p := range peers {
		p.Close()
	}
}

// peer 取出或按需创建与某设备的连接对象。
func (m *Manager) peer(deviceID string) (*Peer, error) {
	m.mu.Lock()
	if p, ok := m.peers[deviceID]; ok {
		m.mu.Unlock()
		return p, nil
	}
	iceServers := m.iceServers
	forceRelay := m.forceRelay
	m.mu.Unlock()

	var mux ice.UDPMux
	if m.mux != nil {
		mux = m.mux
	}
	p, err := NewPeer(PeerOptions{
		DeviceID:        deviceID,
		SelfID:          m.selfID,
		Logger:          m.log,
		ICEServers:      iceServers,
		ForceRelay:      forceRelay,
		UDPMux:          mux,
		ExtraCandidates: m.advertisedCandidates,
		RefreshEgress:   m.refreshEgress,
		SendSignal: func(ctx context.Context, msg *pb.SignalMessage) error {
			return m.sendSignal(ctx, deviceID, msg)
		},
		OnState: func(s ConnState) {
			if m.onState != nil {
				m.onState(deviceID, s)
			}
		},
		OnMessage: func(msg *pb.PeerMessage) {
			if m.onMessage != nil {
				m.onMessage(deviceID, msg)
			}
		},
		OnStream: func(label string, stream *Stream) {
			if m.onStream != nil {
				m.onStream(deviceID, label, stream)
			}
		},
	})
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	// 双检：并发调用时可能已有别的 goroutine 建好了
	if existing, ok := m.peers[deviceID]; ok {
		m.mu.Unlock()
		p.Close()
		return existing, nil
	}
	m.peers[deviceID] = p
	m.mu.Unlock()

	return p, nil
}
