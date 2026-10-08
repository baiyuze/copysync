package p2p

import (
	"context"
	"log/slog"
	"sync"

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
}

func NewManager(opts ManagerOptions) *Manager {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	return &Manager{
		selfID:     opts.SelfID,
		log:        log,
		sendSignal: opts.SendSignal,
		onState:    opts.OnState,
		onMessage:  opts.OnMessage,
		onStream:   opts.OnStream,
		peers:      make(map[string]*Peer),
		forceRelay: opts.ForceRelay,
	}
}

// SetICEServers 更新 STUN/TURN 配置。信令服务器在握手时下发，
// TURN 凭证是短期的，需要定期刷新。
func (m *Manager) SetICEServers(servers []*pb.IceServer) {
	converted := make([]webrtc.ICEServer, 0, len(servers))
	for _, s := range servers {
		converted = append(converted, webrtc.ICEServer{
			URLs:       s.GetUrls(),
			Username:   s.GetUsername(),
			Credential: s.GetCredential(),
		})
	}

	m.mu.Lock()
	m.iceServers = converted
	m.mu.Unlock()

	m.log.Debug("ICE 配置已更新", "servers", len(converted))
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
	ice := m.iceServers
	forceRelay := m.forceRelay
	m.mu.Unlock()

	p, err := NewPeer(PeerOptions{
		DeviceID:   deviceID,
		SelfID:     m.selfID,
		Logger:     m.log,
		ICEServers: ice,
		ForceRelay: forceRelay,
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
