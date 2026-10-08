// Package peers 管理已配对设备：配对流程、在线状态、信令收发。
package peers

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/identity"
	"github.com/baiyuze/copysync/client-core/internal/store"
	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	"github.com/baiyuze/copysync/client-core/internal/transport/signaling"
	"github.com/baiyuze/copysync/proto/deviceid"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// pendingPair 是一次等待用户核对指纹的配对。
//
// 配对必须由双方各自确认指纹后才落库：信令服务器能替换转发中的公钥，
// 但无法让两端显示出相同的指纹——用户的肉眼核对是防中间人的最后一道闸。
type pendingPair struct {
	session   string
	peer      store.Device
	expiresAt time.Time
	// initiated 表示本机是配对码的发起方（对方兑换了我的码）
	initiated bool
}

const pendingTTL = 5 * time.Minute

type Manager struct {
	identity *identity.Identity
	store    *store.Store
	log      *slog.Logger

	loadConfig func() config.Config
	// 配对关系或在线状态变化时通知上层（推给 UI）
	onDeviceChanged func(*pb.Device)
	// 信令连上/断开、对端上下线时通知上层，UI 的连接状态靠它刷新
	onStatusChanged func()
	// 收到对端信令（SDP/candidate），M2 的 WebRTC 层在此接管
	onSignal func(from string, msg *pb.SignalMessage)
	// onPeerReady 见 Options.OnPeerReady
	onPeerReady func(deviceID string)

	client *signaling.Client
	p2p    *p2p.Manager

	mu      sync.RWMutex
	known   map[string]store.Device // 已配对设备
	online  map[string]bool
	conn    map[string]p2p.ConnState // P2P 连接状态（直连/中转）
	pending map[string]*pendingPair  // session -> 待确认配对
	// stashed 暂存配对确认前到达的信令，按对端 device_id 归集
	stashed map[string][]*pb.SignalMessage
	// 等待配对码的请求方（CreatePairingCode 是异步的，靠这个回传结果）
	codeWaiters []chan pairingCode
}

type pairingCode struct {
	code      string
	expiresAt time.Time
}

type Options struct {
	Identity        *identity.Identity
	Store           *store.Store
	Logger          *slog.Logger
	LoadConfig      func() config.Config
	OnDeviceChanged func(*pb.Device)
	OnStatusChanged func()
	OnSignal        func(from string, msg *pb.SignalMessage)
	// OnPeerMessage / OnPeerStream 是 P2P 数据面的入口，M4 的同步逻辑挂在这里
	OnPeerMessage func(deviceID string, msg *pb.PeerMessage)
	OnPeerStream  func(deviceID, label string, stream *p2p.Stream)
	// OnPeerReady 在与某台设备的 P2P 连接建立或恢复时调用（同步引擎据此补发）
	OnPeerReady func(deviceID string)
	// StatePath 保存网络出口历史的文件，见 p2p/egress_history.go
	StatePath string
}

func NewManager(opts Options) (*Manager, error) {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}

	m := &Manager{
		identity:        opts.Identity,
		store:           opts.Store,
		log:             log,
		loadConfig:      opts.LoadConfig,
		onDeviceChanged: opts.OnDeviceChanged,
		onStatusChanged: opts.OnStatusChanged,
		onSignal:        opts.OnSignal,
		onPeerReady:     opts.OnPeerReady,
		known:           make(map[string]store.Device),
		online:          make(map[string]bool),
		conn:            make(map[string]p2p.ConnState),
		pending:         make(map[string]*pendingPair),
		stashed:         make(map[string][]*pb.SignalMessage),
	}

	devices, err := opts.Store.ListDevices()
	if err != nil {
		return nil, fmt.Errorf("读取已配对设备: %w", err)
	}
	for _, d := range devices {
		m.known[d.ID] = d
	}

	m.p2p = p2p.NewManager(p2p.ManagerOptions{
		SelfID:     opts.Identity.DeviceID,
		Logger:     log,
		SendSignal: m.sendSignal,
		OnState:    m.handleP2PState,
		OnMessage:  opts.OnPeerMessage,
		OnStream:   opts.OnPeerStream,
		StatePath:  opts.StatePath,
		// 公共 STUN 服务器 IP 分散，能探测到更多出口；用户可以在设置里关掉
		ProbeServers: func() []string {
			if opts.LoadConfig().OnlyOwnSTUN {
				return nil
			}
			return p2p.DefaultProbeServers
		},
	})

	m.client = signaling.New(signaling.Options{
		Identity:    opts.Identity,
		URL:         opts.LoadConfig().SignalingURL,
		Logger:      log,
		PublicKeyOf: m.publicKeyOf,
		DeviceInfo:  m.deviceInfo,
		Handlers: signaling.Handlers{
			OnConnected:       m.handleConnected,
			OnDisconnected:    m.handleDisconnected,
			OnPresence:        m.handlePresence,
			OnSignal:          m.handleSignal,
			OnPairingCreated:  m.handlePairingCreated,
			OnPairingRedeemed: m.handlePairingRedeemed,
			OnPairingResult:   m.handlePairingResult,
			OnIceConfig: func(servers []*pb.IceServer) {
				m.p2p.SetICEServers(servers)
			},
			OnServerError: func(code, msg string) {
				log.Warn("信令服务器错误", "code", code, "msg", msg)
			},
		},
	})
	return m, nil
}

func (m *Manager) Run(ctx context.Context) {
	go m.gcPending(ctx)
	go m.reconnectLoop(ctx)
	m.client.Run(ctx)
}

// reconnectLoop 定期把"在线但尚未建立 P2P 通道"的对端重新连上。
//
// 单靠上线通知触发连接是不够的：配对刚完成时两端确认有先后，
// 先确认的一方发出的 offer 会被后确认的一方当作未配对消息丢弃；
// 网络切换、信令重连也会让某次握手失败。有了这道巡检，
// 连接状态最终总会收敛，而不需要用户手动重启。
func (m *Manager) reconnectLoop(ctx context.Context) {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.connectPending(ctx)
		}
	}
}

// connectPending 对所有在线但未连通的对端发起连接。
func (m *Manager) connectPending(ctx context.Context) {
	m.mu.RLock()
	var candidates []string
	for id := range m.known {
		if m.online[id] && m.p2p.State(id) == p2p.StateOffline {
			candidates = append(candidates, id)
		}
	}
	m.mu.RUnlock()

	for _, id := range candidates {
		connectCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		if err := m.p2p.Connect(connectCtx, id); err != nil {
			m.log.Debug("重连失败，稍后重试", "device", id, "err", err)
		}
		cancel()
	}
}

func (m *Manager) Connected() bool { return m.client.Connected() }

// SetSignalingURL 在用户于设置里改了服务器地址后调用，立即重连。
func (m *Manager) SetSignalingURL(url string) { m.client.SetURL(url) }

// ─────────────────────── 信令回调 ───────────────────────

// publicKeyOf 供信令层验签。
//
// 除已配对设备外，也认待确认配对中的对端公钥——那把公钥已随配对流程
// 取得并校验过与 device_id 相符。能验签不等于信任：是否建立互信仍由
// 用户在 ConfirmPairing 里决定，这里只是让早到的信令不至于被直接丢弃。
func (m *Manager) publicKeyOf(deviceID string) ed25519.PublicKey {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if d, ok := m.known[deviceID]; ok {
		return ed25519.PublicKey(d.PublicKey)
	}
	for _, p := range m.pending {
		if p.peer.ID == deviceID {
			return ed25519.PublicKey(p.peer.PublicKey)
		}
	}
	return nil
}

func (m *Manager) deviceInfo() (name, platform string, paired []string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make([]string, 0, len(m.known))
	for id := range m.known {
		ids = append(ids, id)
	}
	return m.loadConfig().DeviceName, runtime.GOOS, ids
}

func (m *Manager) handleConnected() {
	m.log.Info("信令已连接", "url", m.client.URL())
	m.notifyAll()
	m.notifyStatus()
}

func (m *Manager) handleDisconnected(error) {
	m.mu.Lock()
	// 信令断开后无法得知对端状态，一律视为离线
	m.online = make(map[string]bool)
	m.mu.Unlock()
	m.notifyAll()
	m.notifyStatus()
}

func (m *Manager) handlePresence(deviceID string, online bool) {
	m.mu.Lock()
	m.online[deviceID] = online
	m.mu.Unlock()

	m.log.Debug("对端状态变化", "device", deviceID, "online", online)

	if online {
		// 对端上线即建立 P2P 通道，这样真正要传数据时已经就绪。
		// 谁发起由 p2p.Manager 内部按 device_id 字典序决定，这里无脑调用即可。
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// 只在完全没有连接时发起。正在握手中就不要插手——
			// 重置会换掉 ICE ufrag，让对端已经发出的 answer 全部失配。
			if m.p2p.State(deviceID) != p2p.StateOffline {
				return
			}
			if err := m.p2p.Connect(ctx, deviceID); err != nil {
				m.log.Warn("建立 P2P 连接失败", "device", deviceID, "err", err)
			}
		}()
	} else {
		m.p2p.Disconnect(deviceID)
		m.mu.Lock()
		delete(m.conn, deviceID)
		m.mu.Unlock()
	}

	if d, ok := m.device(deviceID); ok {
		m.notify(d)
	}
	m.notifyStatus() // 在线对端数随之变化
}

func (m *Manager) handleSignal(from string, msg *pb.SignalMessage) {
	m.mu.RLock()
	_, paired := m.known[from]
	pendingPair := !paired && m.hasPendingLocked(from)
	m.mu.RUnlock()

	if !paired {
		if pendingPair {
			// 配对尚未确认，但对方已经在发起连接了。
			//
			// 两端确认配对必然有先后，先确认的一方发出的 offer 一定早于
			// 后确认方建立互信。直接丢弃会让握手石沉大海，只能等超时重试；
			// 暂存下来在确认后重放，用户确认完就能立刻连上。
			m.stashSignal(from, msg)
			return
		}
		m.log.Warn("忽略未配对设备的信令", "from", from)
		return
	}

	m.dispatchSignal(from, msg)
}

func (m *Manager) dispatchSignal(from string, msg *pb.SignalMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := m.p2p.HandleSignal(ctx, from, msg); err != nil {
		m.log.Warn("处理信令失败", "from", from, "err", err)
	}

	if m.onSignal != nil {
		m.onSignal(from, msg)
	}
}

func (m *Manager) hasPendingLocked(deviceID string) bool {
	for _, p := range m.pending {
		if p.peer.ID == deviceID {
			return true
		}
	}
	return false
}

// stashSignal 暂存配对确认前到达的信令。
func (m *Manager) stashSignal(from string, msg *pb.SignalMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 只留少量：握手信令很快，堆积说明对端在异常重试，没必要全存
	const maxStashed = 32
	if len(m.stashed[from]) >= maxStashed {
		return
	}
	if m.stashed == nil {
		m.stashed = make(map[string][]*pb.SignalMessage)
	}
	m.stashed[from] = append(m.stashed[from], msg)
	m.log.Debug("暂存配对确认前到达的信令", "from", from, "已存", len(m.stashed[from]))
}

// replaySignals 在配对确认后重放暂存的信令。
func (m *Manager) replaySignals(from string) {
	m.mu.Lock()
	msgs := m.stashed[from]
	delete(m.stashed, from)
	m.mu.Unlock()

	if len(msgs) == 0 {
		return
	}
	m.log.Info("重放配对期间暂存的信令", "from", from, "条数", len(msgs))
	for _, msg := range msgs {
		m.dispatchSignal(from, msg)
	}
}

// handleP2PState 在直连/中转状态变化时刷新界面。
func (m *Manager) handleP2PState(deviceID string, state p2p.ConnState) {
	m.mu.Lock()
	prev := m.conn[deviceID]
	m.conn[deviceID] = state
	m.mu.Unlock()

	m.log.Info("P2P 连接状态", "device", deviceID, "state", state)
	if d, ok := m.device(deviceID); ok {
		m.notify(d)
	}
	// 直连与中转之间切换不算重新连上
	if connected(state) && !connected(prev) && m.onPeerReady != nil {
		m.onPeerReady(deviceID)
	}
}

func connected(s p2p.ConnState) bool {
	return s == p2p.StateDirect || s == p2p.StateRelay
}

func (m *Manager) sendSignal(ctx context.Context, to string, msg *pb.SignalMessage) error {
	return m.client.Signal(ctx, to, msg)
}

// P2P 暴露数据面给上层（M4 的同步逻辑用）。
func (m *Manager) P2P() *p2p.Manager { return m.p2p }

func (m *Manager) handlePairingCreated(code string, expiresAt time.Time) {
	m.mu.Lock()
	waiters := m.codeWaiters
	m.codeWaiters = nil
	m.mu.Unlock()

	for _, ch := range waiters {
		ch <- pairingCode{code, expiresAt}
		close(ch)
	}
}

// handlePairingRedeemed：本机是发起方，有人兑换了我的配对码。
// 同样要经用户确认指纹后才落库。
func (m *Manager) handlePairingRedeemed(ev *pb.PairingRedeemed) {
	m.addPending(store.Device{
		ID:        ev.GetPeerDeviceId(),
		Name:      ev.GetPeerDeviceName(),
		Platform:  ev.GetPeerPlatform(),
		PublicKey: ev.GetPeerPublicKey(),
	}, true)
}

// handlePairingResult：本机是兑换方，拿到了发起方信息。
func (m *Manager) handlePairingResult(ev *pb.PairingRedeemResponse) {
	m.addPending(store.Device{
		ID:        ev.GetPeerDeviceId(),
		Name:      ev.GetPeerDeviceName(),
		Platform:  ev.GetPeerPlatform(),
		PublicKey: ev.GetPeerPublicKey(),
	}, false)
}

func (m *Manager) addPending(peer store.Device, initiated bool) {
	// 校验对端 ID 确由其公钥派生，挡住服务器伪造的 device_id
	if !deviceid.Valid(peer.ID, peer.PublicKey) {
		m.log.Warn("拒绝配对：device_id 与公钥不匹配", "device", peer.ID)
		return
	}

	session := randomToken()
	m.mu.Lock()
	m.pending[session] = &pendingPair{
		session:   session,
		peer:      peer,
		expiresAt: time.Now().Add(pendingTTL),
		initiated: initiated,
	}
	m.mu.Unlock()

	m.log.Info("等待用户核对配对指纹",
		"peer", peer.Name, "device", peer.ID,
		"fingerprint", deviceid.Fingerprint(peer.PublicKey))

	// 发起方侧是被动收到的，主动推给 UI 让它弹出确认
	if initiated {
		m.notify(m.toProto(peer, false, session))
	}
}

// ─────────────────────── 对外 API ───────────────────────

var (
	ErrNotConnected   = errors.New("未连接到信令服务器")
	ErrSessionUnknown = errors.New("配对会话不存在或已过期")
)

// CreatePairingCode 生成配对码。服务器异步返回，这里等待结果。
func (m *Manager) CreatePairingCode(ctx context.Context) (string, time.Time, error) {
	if !m.client.Connected() {
		return "", time.Time{}, ErrNotConnected
	}

	ch := make(chan pairingCode, 1)
	m.mu.Lock()
	m.codeWaiters = append(m.codeWaiters, ch)
	m.mu.Unlock()

	if err := m.client.CreatePairingCode(ctx); err != nil {
		return "", time.Time{}, err
	}

	select {
	case got := <-ch:
		return got.code, got.expiresAt, nil
	case <-ctx.Done():
		return "", time.Time{}, ctx.Err()
	case <-time.After(10 * time.Second):
		return "", time.Time{}, errors.New("等待配对码超时")
	}
}

// RedeemPairingCode 兑换配对码，返回待用户核对的对端信息与会话 ID。
// 此时尚未建立互信——必须再调 ConfirmPairing。
func (m *Manager) RedeemPairingCode(ctx context.Context, code string) (*pb.Device, string, error) {
	if !m.client.Connected() {
		return nil, "", ErrNotConnected
	}
	if err := m.client.RedeemPairingCode(ctx, code); err != nil {
		return nil, "", err
	}

	// 等待服务器回传发起方信息（经 handlePairingResult 落入 pending）
	deadline := time.After(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil, "", ctx.Err()
		case <-deadline:
			return nil, "", errors.New("兑换配对码超时（配对码可能已失效）")
		case <-time.After(50 * time.Millisecond):
			m.mu.RLock()
			var found *pendingPair
			for _, p := range m.pending {
				if !p.initiated {
					found = p
					break
				}
			}
			m.mu.RUnlock()
			if found != nil {
				return m.toProto(found.peer, false, found.session), found.session, nil
			}
		}
	}
}

// ConfirmPairing 在用户核对指纹后落库，正式建立互信。
func (m *Manager) ConfirmPairing(ctx context.Context, session string, accept bool) error {
	m.mu.Lock()
	p, ok := m.pending[session]
	if ok {
		delete(m.pending, session)
	}
	m.mu.Unlock()

	if !ok {
		return ErrSessionUnknown
	}
	if !accept {
		m.log.Info("用户拒绝配对", "peer", p.peer.Name)
		m.mu.Lock()
		delete(m.stashed, p.peer.ID)
		m.mu.Unlock()
		return nil
	}

	p.peer.PairedAt = time.Now()
	if err := m.store.PutDevice(p.peer); err != nil {
		return fmt.Errorf("保存设备: %w", err)
	}

	m.mu.Lock()
	m.known[p.peer.ID] = p.peer
	m.mu.Unlock()

	m.log.Info("配对成功", "peer", p.peer.Name, "device", p.peer.ID)

	// 让服务器知道我们现在关心这个设备的在线状态
	if err := m.client.RefreshPaired(ctx); err != nil {
		m.log.Warn("刷新关注列表失败", "err", err)
	}
	if d, ok := m.device(p.peer.ID); ok {
		m.notify(d)
	}

	// 先重放对方在本端确认前发来的信令，这样握手能一次成功，
	// 不必等连接超时后再重来一轮。
	m.replaySignals(p.peer.ID)

	// 对方没抢先发起时，由本端补发一次。
	go func() {
		connectCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if m.p2p.State(p.peer.ID) != p2p.StateOffline {
			return
		}
		if err := m.p2p.Connect(connectCtx, p.peer.ID); err != nil {
			m.log.Debug("配对后建立连接失败，将由巡检重试", "err", err)
		}
	}()
	return nil
}

func (m *Manager) Unpair(ctx context.Context, deviceID string) error {
	if err := m.store.DeleteDevice(deviceID); err != nil {
		return err
	}
	m.mu.Lock()
	delete(m.known, deviceID)
	delete(m.online, deviceID)
	m.mu.Unlock()

	if err := m.client.RefreshPaired(ctx); err != nil {
		m.log.Warn("刷新关注列表失败", "err", err)
	}
	m.log.Info("已解除配对", "device", deviceID)
	return nil
}

// Devices 返回本机与所有已配对设备。
func (m *Manager) Devices() (*pb.Device, []*pb.Device) {
	cfg := m.loadConfig()
	self := &pb.Device{
		Id:                   m.identity.DeviceID,
		Name:                 cfg.DeviceName,
		Platform:             runtime.GOOS,
		Online:               m.client.Connected(),
		PublicKeyFingerprint: m.identity.Fingerprint(),
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	peers := make([]*pb.Device, 0, len(m.known))
	for _, d := range m.known {
		peers = append(peers, m.toProtoLocked(d, m.online[d.ID], ""))
	}
	return self, peers
}

// Signal 向对端发送信令。M2 的 WebRTC 层用它交换 SDP 与 candidate。
func (m *Manager) Signal(ctx context.Context, to string, msg *pb.SignalMessage) error {
	return m.client.Signal(ctx, to, msg)
}

func (m *Manager) IsOnline(deviceID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.online[deviceID]
}

// ─────────────────────── 内部辅助 ───────────────────────

func (m *Manager) device(deviceID string) (*pb.Device, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.known[deviceID]
	if !ok {
		return nil, false
	}
	return m.toProtoLocked(d, m.online[deviceID], ""), true
}

func (m *Manager) toProto(d store.Device, online bool, session string) *pb.Device {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.toProtoLocked(d, online, session)
}

func (m *Manager) toProtoLocked(d store.Device, online bool, session string) *pb.Device {
	conn := pb.ConnectionKind_CONNECTION_KIND_OFFLINE
	if online {
		// 已有 P2P 通道时显示具体是直连还是中转，便于用户排查网络问题
		conn = m.conn[d.ID].Proto()
	}
	return &pb.Device{
		Id:                   d.ID,
		Name:                 d.Name,
		Platform:             d.Platform,
		Online:               online,
		Connection:           conn,
		PublicKeyFingerprint: deviceid.Fingerprint(d.PublicKey),
		PairedAtUnix:         d.PairedAt.Unix(),
		PairingSession:       session,
	}
}

func (m *Manager) notify(d *pb.Device) {
	if m.onDeviceChanged != nil && d != nil {
		m.onDeviceChanged(d)
	}
}

func (m *Manager) notifyAll() {
	_, peers := m.Devices()
	for _, d := range peers {
		m.notify(d)
	}
}

func (m *Manager) notifyStatus() {
	if m.onStatusChanged != nil {
		m.onStatusChanged()
	}
}

// gcPending 清理用户迟迟未确认的配对会话。
func (m *Manager) gcPending(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			m.mu.Lock()
			for k, p := range m.pending {
				if now.After(p.expiresAt) {
					delete(m.pending, k)
				}
			}
			m.mu.Unlock()
		}
	}
}

func randomToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// 密码学随机源不可用时没有安全的降级方案
		panic(fmt.Sprintf("生成会话 ID 失败: %v", err))
	}
	return hex.EncodeToString(b)
}
