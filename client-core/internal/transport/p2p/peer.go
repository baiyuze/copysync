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

	"github.com/pion/ice/v4"
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
	// relayAcceptanceWait 是中转路径最早可以被选中的时间，给直连留出打通的机会。
	// 要长于对端探测出口、把出口地址发过来所需的时间（见 Manager.advertisedCandidates）
	relayAcceptanceWait = 5 * time.Second
	// restartTimeout 是 ICE 重启后等待恢复的最长时间，超时就重置连接
	restartTimeout = 30 * time.Second
	// remoteEgressWait 是「先等对端的出口地址、再发自己的」最多等多久，见 sendExtraCandidates。
	// 要给对端探测出口留出时间（Manager.advertisedCandidates 最多等 4 秒，通常不到 1 秒），
	// 又要短于中转的等待时限
	remoteEgressWait = 3 * time.Second
	// punchHeadStart：收到对端的出口地址后，等本端发往那里的包先出去，再公布自己的
	punchHeadStart = 200 * time.Millisecond
)

var (
	// upgradeDelays 是走中转时，每次尝试换成直连之前等待的时间；用完后按最后一项重复
	upgradeDelays = []time.Duration{10 * time.Second, 30 * time.Second, 2 * time.Minute, 5 * time.Minute, 15 * time.Minute}
	// upgradeIdle：这么久之内有过同步就先不换，ICE 重启会让传输停顿几秒
	upgradeIdle = 10 * time.Second
)

var errNegotiating = errors.New("正在协商中")

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
	// extraCands 返回出口探测得到的公网地址，见 sendExtraCandidates
	extraCands func() []string
	// refreshEgress 重新探测一轮出口，见 tryUpgrade
	refreshEgress func()

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

	forceRelay bool
	// remoteEgress 在收到本轮协商中对端的出口地址时关闭，见 sendExtraCandidates
	remoteEgress     chan struct{}
	remoteEgressSeen bool
	// 生成并发出 offer/answer 期间，新收集到的 candidate 先攒着，见 holdCandidates
	holdCands bool
	heldCands []*pb.IceCandidate
	// 走中转时定期尝试换成直连，见 tryUpgrade
	upgradeTimer   *time.Timer
	upgradeAttempt int
	restarts       int // ICE 重启的次数，watchRestart 用它识别过期的计时
	// lastActivity 与 streams 用来判断连接是否空闲，见 busy
	lastActivity time.Time
	streams      []*webrtc.DataChannel
}

type PeerOptions struct {
	DeviceID   string
	SelfID     string
	Logger     *slog.Logger
	ICEServers []webrtc.ICEServer
	// ForceRelay 强制只用 TURN 中转，不尝试直连。
	// 用于验证中转链路是否可用——正常运行时应保持关闭。
	ForceRelay bool
	// UDPMux 是所有连接共用的打洞端口；为空则由 pion 为每个连接各开端口
	UDPMux ice.UDPMux
	// ExtraCandidates 返回出口探测得到的公网地址（srflx candidate），握手时一并发给对端
	ExtraCandidates func() []string
	// RefreshEgress 立即重新探测出口并等它结束，尝试换成直连之前调用
	RefreshEgress func()
	SendSignal    func(ctx context.Context, msg *pb.SignalMessage) error
	OnState       func(ConnState)
	OnMessage     func(*pb.PeerMessage)
	OnStream      func(label string, stream *Stream)
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
	if opts.UDPMux != nil {
		settings.SetICEUDPMux(opts.UDPMux)
	}

	config := webrtc.Configuration{ICEServers: opts.ICEServers}
	if opts.ForceRelay {
		config.ICETransportPolicy = webrtc.ICETransportPolicyRelay
	} else {
		// 中转路径握手快，若不加等待，往往在直连打通之前就被选中，之后再也不会换成直连
		settings.SetRelayAcceptanceMinWait(relayAcceptanceWait)
	}

	return &Peer{
		deviceID:      opts.DeviceID,
		log:           log.With("peer", opts.DeviceID),
		sendSignal:    opts.SendSignal,
		onState:       opts.OnState,
		onMessage:     opts.OnMessage,
		onStream:      opts.OnStream,
		extraCands:    opts.ExtraCandidates,
		refreshEgress: opts.RefreshEgress,
		api:           webrtc.NewAPI(webrtc.WithSettingEngine(settings)),
		config:        config,
		forceRelay:    opts.ForceRelay,
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

	p.holdCandidates()
	defer p.releaseCandidates()
	remoteEgress := p.expectRemoteEgress()

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		return fmt.Errorf("创建 offer: %w", err)
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		return fmt.Errorf("设置本地描述: %w", err)
	}

	p.setState(StateConnecting)
	p.armConnectTimeout()

	if err := p.sendSignal(ctx, offerMessage(offer, false)); err != nil {
		return err
	}
	// 初次握手由发起方等：应答方先公布出口地址，发起方先往那边发包
	go p.sendExtraCandidates(remoteEgress)
	return nil
}

// restartICE 在已连通的连接上做一次 ICE 重启：两端重新收集地址、重新做连通性检查，
// DTLS 与数据通道保持不变。重启期间发出的数据会暂时丢失，由 SCTP 重传补上，
// 所以只停顿、不断开。
//
// offererWaits 决定这一轮谁先发包：为真时本端等对端的出口地址、先往那边发包再公布
// 自己的；为假时请应答方这样等（见 wire.proto 的 SdpOffer.answerer_waits）。
func (p *Peer) restartICE(ctx context.Context, offererWaits bool) error {
	p.mu.Lock()
	pc := p.pc
	if pc == nil || p.makingOffer || pc.SignalingState() != webrtc.SignalingStateStable {
		p.mu.Unlock()
		return errNegotiating
	}
	p.makingOffer = true
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		p.makingOffer = false
		p.mu.Unlock()
	}()

	p.holdCandidates()
	defer p.releaseCandidates()
	remoteEgress := p.expectRemoteEgress()

	// 先盯住：从 CreateOffer 起本端的 ICE 就已重启，之后任何一步失败都要能恢复
	p.watchRestart(pc)
	offer, err := pc.CreateOffer(&webrtc.OfferOptions{ICERestart: true})
	if err != nil {
		return fmt.Errorf("创建 offer: %w", err)
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		return fmt.Errorf("设置本地描述: %w", err)
	}
	if err := p.sendSignal(ctx, offerMessage(offer, !offererWaits)); err != nil {
		return err
	}
	if !offererWaits {
		remoteEgress = nil
	}
	go p.sendExtraCandidates(remoteEgress)
	return nil
}

// watchRestart 在 ICE 重启后迟迟恢复不了时重置连接，交给上层重连。
// 重启用的 offer 或 answer 若在信令里丢了，ICE 会一直停在检查中，自己不会失败。
func (p *Peer) watchRestart(pc *webrtc.PeerConnection) {
	p.mu.Lock()
	p.restarts++
	gen := p.restarts
	p.mu.Unlock()
	time.AfterFunc(restartTimeout, func() {
		p.mu.Lock()
		stale := p.pc != pc || p.restarts != gen
		p.mu.Unlock()
		if !stale && pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			p.log.Warn("ICE 重启后迟迟没有恢复，重置连接")
			p.Close()
		}
	})
}

func offerMessage(offer webrtc.SessionDescription, answererWaits bool) *pb.SignalMessage {
	return &pb.SignalMessage{Payload: &pb.SignalMessage_Offer{Offer: &pb.SdpOffer{
		Sdp:             offer.SDP,
		DtlsFingerprint: fingerprintFromSDP(offer.SDP),
		AnswererWaits:   answererWaits,
	}}}
}

// expectRemoteEgress 开始一轮协商时调用，返回的通道在收到对端这一轮的出口地址时关闭。
func (p *Peer) expectRemoteEgress() <-chan struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.remoteEgress = make(chan struct{})
	p.remoteEgressSeen = false
	return p.remoteEgress
}

// noteRemoteCandidate 在加入一条对端 candidate 后调用：是出口地址就通知 sendExtraCandidates。
func (p *Peer) noteRemoteCandidate(candidate string) {
	if !isEgressCandidate(candidate) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.remoteEgress != nil && !p.remoteEgressSeen {
		p.remoteEgressSeen = true
		close(p.remoteEgress)
	}
}

// isEgressCandidate 判断是不是出口探测得到的地址（见 srflxCandidates）：类型是 srflx，
// 且不带 ufrag 扩展。pion 自己收集的 candidate 都带 ufrag。
func isEgressCandidate(candidate string) bool {
	return strings.Contains(candidate, " typ srflx") && candidateUfrag(candidate) == ""
}

// holdCandidates 让新收集到的 candidate 先攒着，releaseCandidates 时再发。
//
// 生成 offer/answer 时 pion 就开始收集地址，candidate 可能抢在描述之前发出。
// 初次握手时无妨，对端会先缓存；ICE 重启时对端手里还是旧的远端描述，新一代
// candidate 的 ufrag 对不上，会被直接丢弃。
func (p *Peer) holdCandidates() {
	p.mu.Lock()
	p.holdCands = true
	p.mu.Unlock()
}

func (p *Peer) releaseCandidates() {
	p.mu.Lock()
	held := p.heldCands
	p.heldCands = nil
	p.holdCands = false
	p.mu.Unlock()
	for _, c := range held {
		p.sendCandidate(c)
	}
}

func (p *Peer) sendCandidate(c *pb.IceCandidate) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := p.sendSignal(ctx, &pb.SignalMessage{
		Payload: &pb.SignalMessage_Candidate{Candidate: c},
	}); err != nil {
		p.log.Debug("发送 candidate 失败", "err", err)
	}
}

// sendExtraCandidates 把出口探测得到的各个公网地址作为 srflx candidate 发给对端。
//
// pion 自己只会探测出一个公网地址；在多出口网络里，发往对端的包很可能从另一个出口
// 出去，那个地址对不上，打洞就失败了。把每个出口的地址都告诉对端，对端逐个检查连通性，
// 总有一个与实际出口吻合。必须在 offer/answer 之后发：对端设好远端描述前收到的
// candidate 会先缓存起来。
//
// remoteEgress 不为空时，先等对端的出口地址到达、本端往那里发的包出去之后再发。
// 打洞时谁的包先到，谁的路由器就先收到「陌生来源」的包。有的路由器（实测公司那台）
// 会为它留下一条连接记录，占住这个端口；本端随后往外发时源端口只好换掉，事先告诉
// 对端的地址就失效了。本端先发包，路由器里就已有本端发起的记录，对端的包回来正好对上。
func (p *Peer) sendExtraCandidates(remoteEgress <-chan struct{}) {
	if p.extraCands == nil {
		return
	}
	if remoteEgress != nil {
		select {
		case <-remoteEgress:
			time.Sleep(punchHeadStart)
		case <-time.After(remoteEgressWait):
			p.log.Debug("没等到对端的出口地址，照常发出本端的")
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, c := range p.extraCands() {
		if err := p.sendSignal(ctx, &pb.SignalMessage{
			Payload: &pb.SignalMessage_Candidate{Candidate: &pb.IceCandidate{Candidate: c, SdpMid: "0"}},
		}); err != nil {
			p.log.Debug("发送出口地址失败", "err", err)
			return
		}
	}
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
	// 已连通时收到同一证书的 offer，是对端在做 ICE 重启（见 tryUpgrade）：
	// 连接还在，不能当成新握手把状态打回「连接中」
	restart := pc != nil && p.pcConnected &&
		strings.EqualFold(offer.GetDtlsFingerprint(), p.expectFingerprint)
	// 证书变了，说明对端已经重建了 PeerConnection（重启、或它那边判定连接失效）。
	// 本端的旧连接与它再也接不上，在旧连接上协商只会卡到超时，直接换新的
	if pc != nil && p.expectFingerprint != "" &&
		!strings.EqualFold(offer.GetDtlsFingerprint(), p.expectFingerprint) {
		p.log.Info("对端已重建连接，丢弃本端的旧连接")
		pending := p.pendingCandidates // 对端新连接的 candidate 可能先到，留着
		p.closeLocked()
		p.pendingCandidates = pending
		pc = nil
		making = false
	}
	p.mu.Unlock()

	p.holdCandidates()
	defer p.releaseCandidates()

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
	remoteEgress := p.expectRemoteEgress()
	if !offer.GetAnswererWaits() {
		remoteEgress = nil
	}

	if restart {
		// 设置远端描述时本端的 ICE 随之重启，从这里起就要盯住
		p.watchRestart(pc)
	}
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

	if !restart {
		p.setState(StateConnecting)
		p.armConnectTimeout()
	}

	if err := p.sendSignal(ctx, &pb.SignalMessage{
		Payload: &pb.SignalMessage_Answer{Answer: &pb.SdpAnswer{
			Sdp:             answer.SDP,
			DtlsFingerprint: fingerprintFromSDP(answer.SDP),
		}},
	}); err != nil {
		return err
	}
	go p.sendExtraCandidates(remoteEgress)
	return nil
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
	// 先缓存、等设好远端描述再加的两种情况：
	//   - 远端描述尚未设置，AddICECandidate 会失败；
	//   - ICE 重启时对端新一代的 candidate 抢在它的 offer/answer 之前到达，
	//     ufrag 与手里的远端描述对不上，pion 会直接丢弃。
	// 缓存里若混进旧一代的 candidate，补加时由 pion 按 ufrag 丢弃。
	if pc == nil || pc.RemoteDescription() == nil ||
		!sdpHasUfrag(pc.RemoteDescription().SDP, candidateUfrag(init.Candidate)) {
		p.pendingCandidates = append(p.pendingCandidates, init)
		p.mu.Unlock()
		return nil
	}
	p.mu.Unlock()

	if err := pc.AddICECandidate(init); err != nil {
		return err
	}
	p.noteRemoteCandidate(init.Candidate)
	return nil
}

// candidateUfrag 取出 candidate 里的 ufrag 扩展。pion 生成的 candidate 都带，
// 出口探测得到的那几条不带（返回空）。
func candidateUfrag(candidate string) string {
	f := strings.Fields(candidate)
	for i := 6; i+1 < len(f); i++ { // 前 6 项是 foundation、component、协议、优先级、地址、端口
		if f[i] == "ufrag" {
			return f[i+1]
		}
	}
	return ""
}

// sdpHasUfrag 判断 SDP 里是否有这个 ufrag；ufrag 为空视为匹配。
func sdpHasUfrag(sdp, ufrag string) bool {
	if ufrag == "" {
		return true
	}
	for _, line := range strings.Split(sdp, "\n") {
		if strings.TrimSpace(line) == "a=ice-ufrag:"+ufrag {
			return true
		}
	}
	return false
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
			continue
		}
		p.noteRemoteCandidate(c.Candidate)
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
		cand := &pb.IceCandidate{Candidate: init.Candidate}
		if init.SDPMid != nil {
			cand.SdpMid = *init.SDPMid
		}
		if init.SDPMLineIndex != nil {
			cand.SdpMlineIndex = uint32(*init.SDPMLineIndex)
		}
		p.mu.Lock()
		if p.holdCands {
			p.heldCands = append(p.heldCands, cand)
			p.mu.Unlock()
			return
		}
		p.mu.Unlock()
		p.sendCandidate(cand)
	})

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		if dc.Label() == controlChannelLabel {
			p.bindControl(dc)
			return
		}
		p.trackStream(dc)
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
		p.mu.Lock()
		current := p.pc == pc
		p.mu.Unlock()
		if !current {
			// 已被替换或关闭的连接：它迟到的「已关闭」不能把新连接的状态改成离线
			return
		}
		switch s {
		case webrtc.PeerConnectionStateConnected:
			p.onConnected(pc)
		case webrtc.PeerConnectionStateFailed, webrtc.PeerConnectionStateClosed:
			// 失败与关闭都是终态（本端主动关闭的在上面已被过滤），必须清理掉。
			// 留着的话 Connect 以为连接还在、不会重建；控制通道看上去仍是打开的，
			// 发出的消息石沉大海
			p.log.Warn("连接已失效，关闭以便重连", "state", s)
			p.Close()
		case webrtc.PeerConnectionStateDisconnected:
			// 断开可能是暂时的（网络抖动），pion 会继续检测，恢复不了时转为失败
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
	kind := p.combinedKindLocked()
	p.mu.Unlock()

	if ready {
		p.mu.Lock()
		if p.connectTimer != nil {
			p.connectTimer.Stop()
			p.connectTimer = nil
		}
		p.mu.Unlock()
		p.settle(kind)
		// 只有本端能确定自己发出的数据走没走 TURN，告诉对端，让两边显示一致
		_ = p.Send(&pb.PeerMessage{Payload: &pb.PeerMessage_Link{
			Link: &pb.LinkInfo{Relayed: ownView == StateRelay},
		}})
	}
}

// handleLinkInfo 记下对端的视角，与本端的合起来重新决定显示直连还是中转。
// ICE 重启换成直连后，也是靠对端发来的这条消息从「中转」改回「直连」。
func (p *Peer) handleLinkInfo(l *pb.LinkInfo) {
	p.mu.Lock()
	p.remoteRelayed = l.GetRelayed()
	settled := p.state == StateDirect || p.state == StateRelay
	kind := p.combinedKindLocked()
	p.mu.Unlock()
	if settled {
		p.settle(kind)
	}
}

// combinedKindLocked：两端任一端经中转，数据就经过服务器，算中转。
func (p *Peer) combinedKindLocked() ConnState {
	if p.remoteRelayed {
		return StateRelay
	}
	return p.pendingKind
}

// settle 报告连接就绪后的状态。走中转时安排稍后尝试换成直连，走直连时取消。
func (p *Peer) settle(kind ConnState) {
	p.setState(kind)
	p.mu.Lock()
	defer p.mu.Unlock()
	if kind == StateDirect {
		p.upgradeAttempt = 0
		p.stopUpgradeLocked()
		return
	}
	p.scheduleUpgradeLocked(upgradeDelays[min(p.upgradeAttempt, len(upgradeDelays)-1)])
}

// scheduleUpgradeLocked 安排一次换直连的尝试。只由不礼让的一方（也就是发起连接的
// 一方）来做，两端不会同时发起 ICE 重启。
func (p *Peer) scheduleUpgradeLocked(delay time.Duration) {
	if p.polite || p.forceRelay || p.upgradeTimer != nil {
		return
	}
	p.upgradeTimer = time.AfterFunc(delay, p.tryUpgrade)
}

func (p *Peer) stopUpgradeLocked() {
	if p.upgradeTimer != nil {
		p.upgradeTimer.Stop()
		p.upgradeTimer = nil
	}
}

// tryUpgrade 在走中转时用 ICE 重启再打一次洞。
//
// pion 一旦选定 candidate 对就不会再换。走中转常常只是输了一场赛跑：对端刚启动，
// 出口还没探测完，中转的等待时限就到了；之后直连其实打得通，却再没有机会。
// ICE 重启让两端重新收集地址、同时再打一次洞。仍打不通就照旧落回中转，
// 下一次隔得更久再试（upgradeDelays）。
func (p *Peer) tryUpgrade() {
	p.mu.Lock()
	p.upgradeTimer = nil
	relay := p.pc != nil && p.state == StateRelay
	p.mu.Unlock()
	if !relay {
		return
	}
	if p.busy() {
		// 有数据在传，重启会让它停顿几秒，等空闲了再试
		p.mu.Lock()
		p.scheduleUpgradeLocked(upgradeIdle)
		p.mu.Unlock()
		return
	}

	p.mu.Lock()
	p.upgradeAttempt++
	attempt := p.upgradeAttempt
	p.mu.Unlock()
	// 谁先发包轮流来：初次握手是本端（发起方）先发，没打通的话第一次重试换对端先发，
	// 再下一次换回来。哪一边的路由器会被对方先到的包占住端口，总有一种顺序能避开
	offererWaits := attempt%2 == 0
	p.log.Info("当前经中转，尝试换成直连", "attempt", attempt, "本端先发包", offererWaits)

	// 先重新探测出口，新发现的线路地址随重启一起发给对端
	if p.refreshEgress != nil {
		p.refreshEgress()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := p.restartICE(ctx, offererWaits); err != nil {
		p.log.Debug("ICE 重启没有发起", "err", err)
		p.mu.Lock()
		p.scheduleUpgradeLocked(upgradeDelays[min(p.upgradeAttempt, len(upgradeDelays)-1)])
		p.mu.Unlock()
	}
	// 发起成功时，重启完成后 onConnected → maybeReady → settle 会报告新状态；
	// 仍是中转的话，settle 会按下一档间隔再安排
}

// busy 判断连接上是否有数据在传：有文件流未关闭，或者刚有过控制消息。
func (p *Peer) busy() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if time.Since(p.lastActivity) < upgradeIdle {
		return true
	}
	live := p.streams[:0]
	for _, dc := range p.streams {
		if dc.ReadyState() != webrtc.DataChannelStateClosed {
			live = append(live, dc)
		}
	}
	clear(p.streams[len(live):])
	p.streams = live
	return len(live) > 0
}

func (p *Peer) trackStream(dc *webrtc.DataChannel) {
	p.mu.Lock()
	p.streams = append(p.streams, dc)
	p.lastActivity = time.Now()
	p.mu.Unlock()
}

func (p *Peer) touch() {
	p.mu.Lock()
	p.lastActivity = time.Now()
	p.mu.Unlock()
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
		p.touch()
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
	// 连接断开（或已失效）时，控制通道可能还显示「打开」，往里写只会石沉大海。
	// 报告未连接，让上层记下来等连接恢复后补发
	up := p.state == StateDirect || p.state == StateRelay
	p.mu.Unlock()

	if !up || dc == nil || dc.ReadyState() != webrtc.DataChannelStateOpen {
		return ErrNotConnected
	}
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	if m.GetLink() == nil {
		p.touch()
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
	dc, err := pc.CreateDataChannel(label, &webrtc.DataChannelInit{Ordered: &ordered})
	if err == nil {
		p.trackStream(dc)
	}
	return dc, err
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
	p.holdCands = false
	p.heldCands = nil
	p.streams = nil
	// upgradeAttempt 保留：连接反复重建时，换直连的尝试不必每次都从最密的间隔开始
	p.stopUpgradeLocked()
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
