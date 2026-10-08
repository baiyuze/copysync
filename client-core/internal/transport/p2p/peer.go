// Package p2p 在设备之间建立 WebRTC 连接：优先 ICE 直连打洞，
// 打不通时自动降级到 TURN 中转。
//
// 安全性依赖 DTLS：DataChannel 的内容经 DTLS+SCTP 加密，TURN 只转发 UDP
// 数据包，看不到明文。而 DTLS 证书指纹经信令传输时由设备私钥签名，
// 信令服务器无法替换——这封死了中间人攻击的入口。
package p2p

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
	"google.golang.org/protobuf/proto"

	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// ConnState 是对端连接的状态，会上报给界面。
type ConnState int

const (
	StateOffline ConnState = iota
	StateConnecting
	StateDirect // ICE 打洞成功，点对点直连
	StateRelay  // 打洞失败，经 TURN 中转
)

func (s ConnState) String() string {
	switch s {
	case StateConnecting:
		return "connecting"
	case StateDirect:
		return "direct"
	case StateRelay:
		return "relay"
	default:
		return "offline"
	}
}

func (s ConnState) Proto() pb.ConnectionKind {
	switch s {
	case StateDirect:
		return pb.ConnectionKind_CONNECTION_KIND_DIRECT
	case StateRelay:
		return pb.ConnectionKind_CONNECTION_KIND_RELAY
	case StateOffline:
		return pb.ConnectionKind_CONNECTION_KIND_OFFLINE
	default:
		return pb.ConnectionKind_CONNECTION_KIND_UNSPECIFIED
	}
}

const (
	controlChannelLabel = "control"
	connectTimeout      = 30 * time.Second
)

// Peer 管理与单个对端的 WebRTC 连接。
type Peer struct {
	deviceID string
	log      *slog.Logger

	// sendSignal 把信令发给对端（经信令服务器转发，带签名）
	sendSignal func(ctx context.Context, msg *pb.SignalMessage) error
	onState    func(ConnState)
	onMessage  func(*pb.PeerMessage)
	// onStream 处理对端新开的数据流（大文件传输用）
	onStream func(label string, stream *Stream)

	api    *webrtc.API
	config webrtc.Configuration

	mu      sync.Mutex
	pc      *webrtc.PeerConnection
	control *webrtc.DataChannel
	state   ConnState
	// 对端 DTLS 指纹，来自签名过的信令。握手完成后必须与实际使用的证书一致。
	expectFingerprint string
	// 已缓存但尚未能加入的 candidate（remote description 还没设置时会发生）
	pendingCandidates []webrtc.ICECandidateInit
	// polite 决定 glare（双方同时发起）时谁让步
	polite bool
	// makingOffer 用于识别 glare：本端正在发 offer 时又收到对方的 offer
	makingOffer bool

	// 连接就绪需要两个条件同时满足：ICE/DTLS 握手完成，且控制通道已打开。
	// 只看前者会导致状态报成"已连接"但 Send 立刻失败。
	pcConnected bool
	pendingKind ConnState
	// remoteRelayed 是对端经 LinkInfo 告知的「它那边是否走中转」，见 wire.proto 的 LinkInfo
	remoteRelayed bool
	// connectTimer 在握手迟迟不完成时触发重置，避免永久卡在 connecting
	connectTimer *time.Timer
	// dialing 表示 Connect 正在进行中，防止并发重入建出两个 PeerConnection
	dialing bool
}

type PeerOptions struct {
	DeviceID   string
	SelfID     string
	Logger     *slog.Logger
	ICEServers []webrtc.ICEServer
	// ForceRelay 强制只用 TURN 中转，不尝试直连。
	// 用于验证中转链路是否可用——正常运行时应保持关闭。
	ForceRelay bool
	SendSignal func(ctx context.Context, msg *pb.SignalMessage) error
	OnState    func(ConnState)
	OnMessage  func(*pb.PeerMessage)
	OnStream   func(label string, stream *Stream)
}

func NewPeer(opts PeerOptions) (*Peer, error) {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}

	// 刻意不启用 SettingEngine.DetachDataChannels()：
	// 它是全局开关，一旦开启，所有通道的 OnMessage/Send 都失效，必须改用
	// dc.Detach() 拿裸流。控制通道用消息回调更自然，而文件传输所需的
	// 分片与背压（BufferedAmount / OnBufferedAmountLow）本来就要自己处理，
	// 统一走消息接口反而让两类通道的代码保持一致。
	settings := webrtc.SettingEngine{}

	config := webrtc.Configuration{ICEServers: opts.ICEServers}
	if opts.ForceRelay {
		config.ICETransportPolicy = webrtc.ICETransportPolicyRelay
	}

	return &Peer{
		deviceID:   opts.DeviceID,
		log:        log.With("peer", opts.DeviceID),
		sendSignal: opts.SendSignal,
		onState:    opts.OnState,
		onMessage:  opts.OnMessage,
		onStream:   opts.OnStream,
		api:        webrtc.NewAPI(webrtc.WithSettingEngine(settings)),
		config:     config,
		// 用 device_id 字典序决定礼让方：双方规则一致且无需协商，
		// 天然避免两端同时发 offer 时的死锁。
		polite: opts.SelfID > opts.DeviceID,
	}, nil
}

func (p *Peer) State() ConnState {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.state
}

func (p *Peer) setState(s ConnState) {
	p.mu.Lock()
	changed := p.state != s
	p.state = s
	p.mu.Unlock()

	if changed {
		p.log.Debug("连接状态变化", "state", s)
		if p.onState != nil {
			p.onState(s)
		}
	}
}

// Connect 主动发起连接。已在连接中或已连通时直接返回。
func (p *Peer) Connect(ctx context.Context) error {
	// 必须在同一个临界区里完成"检查 + 占位"。
	//
	// 否则两个 goroutine（配对确认与上线通知常常同时触发）都会看到
	// pc == nil，各建一个 PeerConnection，后者覆盖前者；对端回来的
	// answer 属于前一个，落到后一个身上就成了"时序不符"而被丢弃，
	// 握手随之永久卡住。
	p.mu.Lock()
	if p.pc != nil || p.dialing {
		p.mu.Unlock()
		return nil
	}
	p.dialing = true
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		p.dialing = false
		p.mu.Unlock()
	}()

	pc, err := p.newPeerConnection()
	if err != nil {
		return err
	}

	// 发起方创建控制通道；应答方靠 OnDataChannel 拿到它
	ordered := true
	dc, err := pc.CreateDataChannel(controlChannelLabel, &webrtc.DataChannelInit{
		Ordered: &ordered,
	})
	if err != nil {
		return fmt.Errorf("创建控制通道: %w", err)
	}
	p.bindControl(dc)

	p.mu.Lock()
	p.makingOffer = true
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		p.makingOffer = false
		p.mu.Unlock()
	}()

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		return fmt.Errorf("创建 offer: %w", err)
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		return fmt.Errorf("设置本地描述: %w", err)
	}

	p.setState(StateConnecting)
	p.armConnectTimeout()

	return p.sendSignal(ctx, &pb.SignalMessage{
		Payload: &pb.SignalMessage_Offer{Offer: &pb.SdpOffer{
			Sdp:             offer.SDP,
			DtlsFingerprint: fingerprintFromSDP(offer.SDP),
		}},
	})
}

// armConnectTimeout 在握手迟迟不完成时重置连接，让上层得以重试。
//
// 现实中有多种情况会让 offer 石沉大海：对端刚完成配对、尚未把本端加入
// 已配对列表，于是把 offer 当作未配对设备的消息丢弃；或者信令在网络切换
// 时丢包。没有这道兜底，Peer 会永远停在 connecting。
func (p *Peer) armConnectTimeout() {
	p.mu.Lock()
	if p.connectTimer != nil {
		p.connectTimer.Stop()
	}
	p.connectTimer = time.AfterFunc(connectTimeout, func() {
		if p.State() == StateConnecting {
			p.log.Warn("连接握手超时，重置以便重试")
			p.Close()
		}
	})
	p.mu.Unlock()
}
func (p *Peer) HandleSignal(ctx context.Context, msg *pb.SignalMessage) error {
	switch m := msg.GetPayload().(type) {
	case *pb.SignalMessage_Offer:
		return p.handleOffer(ctx, m.Offer)
	case *pb.SignalMessage_Answer:
		return p.handleAnswer(m.Answer)
	case *pb.SignalMessage_Candidate:
		return p.handleCandidate(m.Candidate)
	default:
		return nil
	}
}

func (p *Peer) handleOffer(ctx context.Context, offer *pb.SdpOffer) error {
	p.mu.Lock()
	pc := p.pc
	polite := p.polite
	making := p.makingOffer
	p.mu.Unlock()

	// Perfect Negotiation：双方同时发 offer 时（glare），
	// 不礼让方忽略对方的 offer，礼让方回滚自己的 local description。
	//
	// 关键是回滚而非关闭重建——重建会换掉 ICE ufrag，
	// 让已经在途的 candidate 全部失配，连接就再也建不起来了。
	if pc != nil {
		collision := making || pc.SignalingState() != webrtc.SignalingStateStable
		if collision {
			if !polite {
				p.log.Debug("忽略对端 offer（本端不礼让）")
				return nil
			}
			p.log.Debug("礼让对端 offer，回滚本端的 local description")
			if err := pc.SetLocalDescription(webrtc.SessionDescription{
				Type: webrtc.SDPTypeRollback,
			}); err != nil {
				return fmt.Errorf("回滚本地描述: %w", err)
			}
		}
	} else {
		var err error
		if pc, err = p.newPeerConnection(); err != nil {
			return err
		}
	}

	p.mu.Lock()
	p.expectFingerprint = offer.GetDtlsFingerprint()
	p.mu.Unlock()

	if err := pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer.GetSdp(),
	}); err != nil {
		return fmt.Errorf("设置远端描述: %w", err)
	}
	p.flushPendingCandidates()

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return fmt.Errorf("创建 answer: %w", err)
	}
	if err := pc.SetLocalDescription(answer); err != nil {
		return fmt.Errorf("设置本地描述: %w", err)
	}

	p.setState(StateConnecting)
	p.armConnectTimeout()

	return p.sendSignal(ctx, &pb.SignalMessage{
		Payload: &pb.SignalMessage_Answer{Answer: &pb.SdpAnswer{
			Sdp:             answer.SDP,
			DtlsFingerprint: fingerprintFromSDP(answer.SDP),
		}},
	})
}

func (p *Peer) handleAnswer(answer *pb.SdpAnswer) error {
	p.mu.Lock()
	pc := p.pc
	p.expectFingerprint = answer.GetDtlsFingerprint()
	p.mu.Unlock()

	if pc == nil {
		return errors.New("收到 answer 但本端没有进行中的连接")
	}
	if pc.SignalingState() != webrtc.SignalingStateHaveLocalOffer {
		p.log.Debug("忽略时序不符的 answer", "state", pc.SignalingState())
		return nil
	}
	if err := pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answer.GetSdp(),
	}); err != nil {
		return fmt.Errorf("设置远端描述: %w", err)
	}
	p.flushPendingCandidates()
	return nil
}

func (p *Peer) handleCandidate(c *pb.IceCandidate) error {
	init := webrtc.ICECandidateInit{Candidate: c.GetCandidate()}
	if mid := c.GetSdpMid(); mid != "" {
		init.SDPMid = &mid
	}
	if idx := uint16(c.GetSdpMlineIndex()); true {
		init.SDPMLineIndex = &idx
	}

	p.mu.Lock()
	pc := p.pc
	// remote description 尚未设置时 AddICECandidate 会失败，先缓存
	if pc == nil || pc.RemoteDescription() == nil {
		p.pendingCandidates = append(p.pendingCandidates, init)
		p.mu.Unlock()
		return nil
	}
	p.mu.Unlock()

	return pc.AddICECandidate(init)
}

func (p *Peer) flushPendingCandidates() {
	p.mu.Lock()
	pc := p.pc
	pending := p.pendingCandidates
	p.pendingCandidates = nil
	p.mu.Unlock()

	if pc == nil {
		return
	}
	for _, c := range pending {
		if err := pc.AddICECandidate(c); err != nil {
			p.log.Debug("补加 candidate 失败", "err", err)
		}
	}
}

func (p *Peer) newPeerConnection() (*webrtc.PeerConnection, error) {
	pc, err := p.api.NewPeerConnection(p.config)
	if err != nil {
		return nil, fmt.Errorf("创建 PeerConnection: %w", err)
	}

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return // candidate 收集完毕
		}
		init := c.ToJSON()
		var mid string
		if init.SDPMid != nil {
			mid = *init.SDPMid
		}
		var idx uint32
		if init.SDPMLineIndex != nil {
			idx = uint32(*init.SDPMLineIndex)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := p.sendSignal(ctx, &pb.SignalMessage{
			Payload: &pb.SignalMessage_Candidate{Candidate: &pb.IceCandidate{
				Candidate:     init.Candidate,
				SdpMid:        mid,
				SdpMlineIndex: idx,
			}},
		}); err != nil {
			p.log.Debug("发送 candidate 失败", "err", err)
		}
	})

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		if dc.Label() == controlChannelLabel {
			p.bindControl(dc)
			return
		}
		// 数据流必须在回调里**立即**挂上接收器。
		//
		// pion 的 OnMessage 只是存一个 handler；在它被设置之前到达的消息
		// 会被直接丢弃，不做缓冲。上层拿到通道后还要查表、建管道，
		// 这中间的几毫秒足以让小文件的全部内容消失得无影无踪。
		// NewStream 先把消息接住，上层再决定怎么消费。
		stream := NewStream(dc)
		if p.onStream != nil {
			p.onStream(dc.Label(), stream)
		}
	})

	pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		p.log.Debug("PeerConnection 状态", "state", s)
		switch s {
		case webrtc.PeerConnectionStateConnected:
			p.onConnected(pc)
		case webrtc.PeerConnectionStateFailed,
			webrtc.PeerConnectionStateClosed,
			webrtc.PeerConnectionStateDisconnected:
			p.setState(StateOffline)
		}
	})

	p.mu.Lock()
	p.pc = pc
	p.mu.Unlock()
	return pc, nil
}

// onConnected 在握手完成后校验 DTLS 指纹，并判定走的是直连还是中转。
func (p *Peer) onConnected(pc *webrtc.PeerConnection) {
	// ★ 安全校验：实际协商用的证书指纹必须与签名信令里声明的一致。
	// 不一致说明有人在中间替换了 SDP——立即断开。
	p.mu.Lock()
	expect := p.expectFingerprint
	p.mu.Unlock()

	if remote := pc.RemoteDescription(); remote != nil && expect != "" {
		actual := fingerprintFromSDP(remote.SDP)
		if !strings.EqualFold(actual, expect) {
			p.log.Error("DTLS 指纹不匹配，疑似中间人攻击，断开连接",
				"expect", expect, "actual", actual)
			p.Close()
			return
		}
	}

	kind := p.detectConnectionKind(pc)
	p.mu.Lock()
	p.pcConnected = true
	p.pendingKind = kind
	p.mu.Unlock()

	p.maybeReady()
}

// maybeReady 只有在传输握手完成**且**控制通道已打开时才对外报告"已连接"。
//
// 两者缺一不可：PeerConnection 进入 Connected 时 SCTP 可能仍在建流，
// 此时报告已连接会让调用方立刻发送并收到 ErrNotConnected。
func (p *Peer) maybeReady() {
	p.mu.Lock()
	ready := p.pcConnected && p.control != nil &&
		p.control.ReadyState() == webrtc.DataChannelStateOpen
	ownView := p.pendingKind
	kind := ownView
	if p.remoteRelayed {
		kind = StateRelay
	}
	p.mu.Unlock()

	if ready {
		p.mu.Lock()
		if p.connectTimer != nil {
			p.connectTimer.Stop()
			p.connectTimer = nil
		}
		p.mu.Unlock()
		p.setState(kind)
		// 只有本端能确定自己发出的数据走没走 TURN，告诉对端，让两边显示一致
		_ = p.Send(&pb.PeerMessage{Payload: &pb.PeerMessage_Link{
			Link: &pb.LinkInfo{Relayed: ownView == StateRelay},
		}})
	}
}

// handleLinkInfo 记下对端的视角；对端经中转时，本端也改报中转。
func (p *Peer) handleLinkInfo(l *pb.LinkInfo) {
	p.mu.Lock()
	p.remoteRelayed = l.GetRelayed()
	upgrade := l.GetRelayed() && p.state == StateDirect
	p.mu.Unlock()
	if upgrade {
		p.setState(StateRelay)
	}
}

// detectConnectionKind 判断当前选中的 candidate 对是直连还是经 TURN 中转。
//
// 两端的 candidate 都要看：只要有一端是 relay，数据就经过 TURN。
// 常见的情形是一端在对称 NAT 后面，用自己的中转地址发送，另一端仍从普通地址收发；
// 只看本端的话，后者会把中转误报成直连，两台设备一台显示直连、一台显示中转。
//
// 选中的 candidate 对直接从 ICE 传输层取：受控一方（answer 方）的统计数据里
// 往往找不到标记为 nominated 的 candidate 对，只靠统计数据会落到默认值上。
func (p *Peer) detectConnectionKind(pc *webrtc.PeerConnection) ConnState {
	if sctp := pc.SCTP(); sctp != nil && sctp.Transport() != nil {
		if ice := sctp.Transport().ICETransport(); ice != nil {
			if pair, err := ice.GetSelectedCandidatePair(); err == nil && pair != nil &&
				pair.Local != nil && pair.Remote != nil {
				if pair.Local.Typ == webrtc.ICECandidateTypeRelay ||
					pair.Remote.Typ == webrtc.ICECandidateTypeRelay {
					return StateRelay
				}
				return StateDirect
			}
		}
	}

	// 取不到选中的 candidate 对时，退回统计数据
	stats := pc.GetStats()
	for _, s := range stats {
		pair, ok := s.(webrtc.ICECandidatePairStats)
		if !ok || pair.State != webrtc.StatsICECandidatePairStateSucceeded || !pair.Nominated {
			continue
		}
		for _, cs := range stats {
			c, ok := cs.(webrtc.ICECandidateStats)
			if ok && (c.ID == pair.LocalCandidateID || c.ID == pair.RemoteCandidateID) &&
				c.CandidateType == webrtc.ICECandidateTypeRelay {
				return StateRelay
			}
		}
		return StateDirect
	}
	// 拿不到明细时报直连——功能不受影响，只是显示可能不准
	return StateDirect
}

func (p *Peer) bindControl(dc *webrtc.DataChannel) {
	p.mu.Lock()
	p.control = dc
	p.mu.Unlock()

	dc.OnOpen(func() {
		p.log.Debug("控制通道已打开")
		p.maybeReady()
	})
	// 应答方绑定控制通道时，传输层可能早已就绪，此时不会再触发 OnOpen
	if dc.ReadyState() == webrtc.DataChannelStateOpen {
		p.maybeReady()
	}
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		var m pb.PeerMessage
		if err := proto.Unmarshal(msg.Data, &m); err != nil {
			p.log.Warn("控制消息解析失败", "err", err)
			return
		}
		if l := m.GetLink(); l != nil {
			p.handleLinkInfo(l) // 连接层自己的消息，不交给同步引擎
			return
		}
		if p.onMessage != nil {
			p.onMessage(&m)
		}
	})
}

var ErrNotConnected = errors.New("与对端的连接尚未建立")

// Send 经控制通道发送一条消息。
func (p *Peer) Send(m *pb.PeerMessage) error {
	p.mu.Lock()
	dc := p.control
	p.mu.Unlock()

	if dc == nil || dc.ReadyState() != webrtc.DataChannelStateOpen {
		return ErrNotConnected
	}
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return dc.Send(data)
}

// OpenStream 为一次文件传输新开一条数据通道。
//
// 与控制通道分离，避免大文件把控制消息堵在队列后面。
func (p *Peer) OpenStream(label string) (*webrtc.DataChannel, error) {
	p.mu.Lock()
	pc := p.pc
	p.mu.Unlock()

	if pc == nil {
		return nil, ErrNotConnected
	}
	ordered := true
	return pc.CreateDataChannel(label, &webrtc.DataChannelInit{Ordered: &ordered})
}

func (p *Peer) Close() {
	p.mu.Lock()
	p.closeLocked()
	p.mu.Unlock()
	p.setState(StateOffline)
}

func (p *Peer) closeLocked() {
	if p.control != nil {
		p.control.Close()
		p.control = nil
	}
	if p.pc != nil {
		p.pc.Close()
		p.pc = nil
	}
	p.pendingCandidates = nil
	p.pcConnected = false
	p.makingOffer = false
	p.pendingKind = StateOffline
	p.remoteRelayed = false
	p.dialing = false
	if p.connectTimer != nil {
		p.connectTimer.Stop()
		p.connectTimer = nil
	}
}

// fingerprintFromSDP 从 SDP 中提取 DTLS 证书指纹（a=fingerprint 行）。
func fingerprintFromSDP(sdp string) string {
	for _, line := range strings.Split(sdp, "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "a=fingerprint:"); ok {
			return strings.TrimSpace(after)
		}
	}
	return ""
}
