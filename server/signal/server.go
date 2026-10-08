package signal

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/copysync/proto/deviceid"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
	"github.com/baiyuze/copysync/server/pairing"
)

const (
	// 信令消息都很小（最大的是 SDP，通常几 KB）
	maxMessageBytes  = 256 << 10
	handshakeTimeout = 15 * time.Second
	writeTimeout     = 10 * time.Second
	// 心跳：及时发现半开连接，否则对端会一直显示"在线"
	pingInterval = 30 * time.Second
	pongTimeout  = 60 * time.Second
)

type Server struct {
	hub     *Hub
	broker  *pairing.Broker
	iceFunc func(deviceID string) []*pb.IceServer
	log     *slog.Logger
}

type Options struct {
	// ICEServers 为每个设备签发 STUN/TURN 配置（含短期 TURN 凭证）
	ICEServers func(deviceID string) []*pb.IceServer
	Logger     *slog.Logger
}

func NewServer(opts Options) *Server {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	ice := opts.ICEServers
	if ice == nil {
		ice = func(string) []*pb.IceServer { return nil }
	}
	return &Server{
		hub:     NewHub(),
		broker:  pairing.NewBroker(),
		iceFunc: ice,
		log:     log,
	}
}

func (s *Server) Hub() *Hub { return s.hub }

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// 客户端是桌面应用而非浏览器，不存在跨站请求伪造的场景
		InsecureSkipVerify: true,
	})
	if err != nil {
		s.log.Warn("websocket 握手失败", "err", err, "remote", r.RemoteAddr)
		return
	}
	ws.SetReadLimit(maxMessageBytes)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	conn, err := s.handshake(ctx, ws)
	if err != nil {
		s.log.Warn("设备认证失败", "err", err, "remote", r.RemoteAddr)
		ws.Close(websocket.StatusPolicyViolation, "认证失败")
		return
	}

	s.log.Info("设备已连接",
		"device", conn.DeviceID, "name", conn.DeviceName,
		"platform", conn.Platform, "在线数", s.hub.Count())

	defer func() {
		if s.hub.Unregister(conn) {
			s.broadcastPresence(ctx, conn.DeviceID, false)
			s.log.Info("设备已断开", "device", conn.DeviceID, "在线数", s.hub.Count())
		}
	}()

	// 出站泵：所有写操作集中在这一个 goroutine，websocket 不支持并发写
	go s.writePump(ctx, ws, conn)

	s.readLoop(ctx, ws, conn)
}

// handshake 完成「证明你持有这把私钥」的挑战应答。
//
// 服务器不维护用户表，所以认证是自证式的：设备提交公钥和对 challenge 的签名，
// 服务器验签并校验 device_id 确由该公钥派生。这挡住了冒用他人 device_id 的连接。
func (s *Server) handshake(ctx context.Context, ws *websocket.Conn) (*Conn, error) {
	ctx, cancel := context.WithTimeout(ctx, handshakeTimeout)
	defer cancel()

	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, err
	}

	hello := &pb.ServerEnvelope{
		Payload: &pb.ServerEnvelope_Hello{Hello: &pb.ServerHello{
			Challenge:      challenge,
			ServerTimeUnix: time.Now().Unix(),
		}},
	}
	if err := writeProto(ctx, ws, hello); err != nil {
		return nil, fmt.Errorf("发送 ServerHello: %w", err)
	}

	var env pb.ClientEnvelope
	if err := readProto(ctx, ws, &env); err != nil {
		return nil, fmt.Errorf("读取 ClientHello: %w", err)
	}
	ch := env.GetHello()
	if ch == nil {
		return nil, errors.New("首条消息必须是 ClientHello")
	}

	pub := ed25519.PublicKey(ch.GetPublicKey())
	if len(pub) != ed25519.PublicKeySize {
		return nil, errors.New("公钥长度不合法")
	}
	if !ed25519.Verify(pub, challenge, ch.GetSignature()) {
		return nil, errors.New("challenge 签名验证失败")
	}
	if !deviceid.Valid(ch.GetDeviceId(), pub) {
		return nil, errors.New("device_id 与公钥不匹配")
	}

	conn := &Conn{
		DeviceID:   ch.GetDeviceId(),
		DeviceName: ch.GetDeviceName(),
		Platform:   ch.GetPlatform(),
		PublicKey:  pub,
		send:       make(chan []byte, 32),
		watching:   make(map[string]struct{}, len(ch.GetPairedDeviceIds())),
	}
	for _, p := range ch.GetPairedDeviceIds() {
		conn.watching[p] = struct{}{}
	}

	if replaced := s.hub.Register(conn); replaced != nil {
		s.log.Info("顶替同设备的旧连接", "device", conn.DeviceID)
		close(replaced.send)
	}

	// ICE 配置含按设备签发的短期 TURN 凭证，必须等认证通过后才能下发
	if servers := s.iceFunc(conn.DeviceID); len(servers) > 0 {
		conn.Send(mustMarshal(&pb.ServerEnvelope{
			Payload: &pb.ServerEnvelope_IceConfig{
				IceConfig: &pb.IceConfig{IceServers: servers},
			},
		}))
	}

	// 先告诉新连接：它关心的对端里谁在线
	s.sendCurrentPresence(conn)
	// 再告诉关心它的对端：我上线了
	s.broadcastPresence(ctx, conn.DeviceID, true)

	return conn, nil
}

// sendCurrentPresence 把「该连接所关注的对端中当前在线的那些」推给它。
func (s *Server) sendCurrentPresence(conn *Conn) {
	for _, peer := range s.hub.PresenceFor(conn.DeviceID) {
		conn.Send(mustMarshal(&pb.ServerEnvelope{
			Payload: &pb.ServerEnvelope_Presence{Presence: &pb.PeerPresence{
				DeviceId: peer, Online: true,
			}},
		}))
	}
}

func (s *Server) readLoop(ctx context.Context, ws *websocket.Conn, conn *Conn) {
	for {
		var env pb.ClientEnvelope
		if err := readProto(ctx, ws, &env); err != nil {
			if ctx.Err() == nil && websocket.CloseStatus(err) == -1 {
				s.log.Debug("读取中断", "device", conn.DeviceID, "err", err)
			}
			return
		}

		switch p := env.GetPayload().(type) {
		case *pb.ClientEnvelope_Forward:
			s.forward(conn, p.Forward)
		case *pb.ClientEnvelope_PairingCreate:
			s.handlePairingCreate(ctx, conn, p.PairingCreate)
		case *pb.ClientEnvelope_PairingRedeem:
			s.handlePairingRedeem(conn, p.PairingRedeem)
		case *pb.ClientEnvelope_Hello:
			// 重复的 Hello 用来刷新关注列表（刚完成配对时会发）
			s.hub.SetWatching(conn.DeviceID, p.Hello.GetPairedDeviceIds())
			// 必须补发一次当前在线状态：presence 只在设备上下线的瞬间广播，
			// 新加入关注列表的对端此刻可能早已在线，不补发的话本端会一直显示它离线。
			s.sendCurrentPresence(conn)
			// 同理，对端也可能正关注着本设备却错过了上线广播
			s.broadcastPresence(ctx, conn.DeviceID, true)
		default:
			s.sendError(conn, "unknown_message", "未知的消息类型")
		}
	}
}

// forward 把带签名的信令原样转发。
//
// 服务器不验签也不解析内容——验签是接收方的事（它有对端公钥，服务器没有）。
// 服务器只确保 from_device_id 不能伪造：强制改写为连接自身的身份。
func (s *Server) forward(from *Conn, signed *pb.Signed) {
	if signed == nil {
		return
	}
	// 防止 A 冒充 B 向 C 发消息
	signed.FromDeviceId = from.DeviceID

	var inner pb.SignalMessage
	if err := proto.Unmarshal(signed.GetPayload(), &inner); err != nil {
		s.sendError(from, "bad_payload", "signal payload 无法解析")
		return
	}
	to := inner.GetToDeviceId()
	if to == "" {
		s.sendError(from, "missing_target", "缺少 to_device_id")
		return
	}

	target, ok := s.hub.Get(to)
	if !ok {
		s.sendError(from, "peer_offline", fmt.Sprintf("对端 %s 不在线", to))
		return
	}
	if !target.Send(mustMarshal(&pb.ServerEnvelope{
		Payload: &pb.ServerEnvelope_Forwarded{Forwarded: signed},
	})) {
		s.log.Warn("转发丢弃：对端出站队列已满", "from", from.DeviceID, "to", to)
	}
}

func (s *Server) handlePairingCreate(ctx context.Context, conn *Conn, req *pb.PairingCreateRequest) {
	code, expires, redeemed, err := s.broker.Create(pairing.Peer{
		DeviceID:   conn.DeviceID,
		DeviceName: req.GetDeviceName(),
		Platform:   conn.Platform,
		PublicKey:  req.GetPublicKey(),
	})
	if err != nil {
		s.sendError(conn, "pairing_create_failed", err.Error())
		return
	}

	conn.Send(mustMarshal(&pb.ServerEnvelope{
		Payload: &pb.ServerEnvelope_PairingCreated{PairingCreated: &pb.PairingCreateResponse{
			Code:          code,
			ExpiresAtUnix: expires.Unix(),
		}},
	}))

	// 等待有人兑换。兑换方的信息要回传给发起方，
	// 否则发起方拿不到对端公钥，配对只能单向完成。
	go func() {
		select {
		case peer, ok := <-redeemed:
			if !ok {
				return // 已过期或被取消
			}
			conn.Send(mustMarshal(&pb.ServerEnvelope{
				Payload: &pb.ServerEnvelope_PairingRedeemed{PairingRedeemed: &pb.PairingRedeemed{
					PeerDeviceId:   peer.DeviceID,
					PeerDeviceName: peer.DeviceName,
					PeerPlatform:   peer.Platform,
					PeerPublicKey:  peer.PublicKey,
				}},
			}))
			s.log.Info("配对完成", "initiator", conn.DeviceID, "peer", peer.DeviceID)
		case <-ctx.Done():
			s.broker.Cancel(code)
		}
	}()
}

func (s *Server) handlePairingRedeem(conn *Conn, req *pb.PairingRedeemRequest) {
	initiator, err := s.broker.Redeem(req.GetCode(), pairing.Peer{
		DeviceID:   conn.DeviceID,
		DeviceName: req.GetDeviceName(),
		Platform:   conn.Platform,
		PublicKey:  req.GetPublicKey(),
	})
	if err != nil {
		s.sendError(conn, "pairing_redeem_failed", err.Error())
		return
	}

	conn.Send(mustMarshal(&pb.ServerEnvelope{
		Payload: &pb.ServerEnvelope_PairingRedeemedResult{
			PairingRedeemedResult: &pb.PairingRedeemResponse{
				PeerDeviceId:   initiator.DeviceID,
				PeerDeviceName: initiator.DeviceName,
				PeerPlatform:   initiator.Platform,
				PeerPublicKey:  initiator.PublicKey,
			},
		},
	}))
}

func (s *Server) broadcastPresence(_ context.Context, deviceID string, online bool) {
	data := mustMarshal(&pb.ServerEnvelope{
		Payload: &pb.ServerEnvelope_Presence{Presence: &pb.PeerPresence{
			DeviceId: deviceID, Online: online,
		}},
	})
	for _, w := range s.hub.WatchersOf(deviceID) {
		w.Send(data)
	}
}

func (s *Server) sendError(conn *Conn, code, msg string) {
	conn.Send(mustMarshal(&pb.ServerEnvelope{
		Payload: &pb.ServerEnvelope_Error{Error: &pb.ServerError{Code: code, Message: msg}},
	}))
}

// writePump 独占写端。websocket 连接不允许并发写，所有出站消息都经此串行化。
func (s *Server) writePump(ctx context.Context, ws *websocket.Conn, conn *Conn) {
	ping := time.NewTicker(pingInterval)
	defer ping.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-conn.send:
			if !ok {
				ws.Close(websocket.StatusNormalClosure, "连接被顶替")
				return
			}
			wctx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := ws.Write(wctx, websocket.MessageBinary, data)
			cancel()
			if err != nil {
				return
			}
		case <-ping.C:
			pctx, cancel := context.WithTimeout(ctx, pongTimeout)
			err := ws.Ping(pctx)
			cancel()
			if err != nil {
				s.log.Debug("心跳失败，断开", "device", conn.DeviceID, "err", err)
				ws.Close(websocket.StatusGoingAway, "心跳超时")
				return
			}
		}
	}
}

// ─────────────────────── proto over websocket ───────────────────────

func readProto(ctx context.Context, ws *websocket.Conn, m proto.Message) error {
	typ, data, err := ws.Read(ctx)
	if err != nil {
		return err
	}
	if typ != websocket.MessageBinary {
		return fmt.Errorf("期望二进制消息，收到 %v", typ)
	}
	return proto.Unmarshal(data, m)
}

func writeProto(ctx context.Context, ws *websocket.Conn, m proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return ws.Write(ctx, websocket.MessageBinary, data)
}

// mustMarshal 用于出站信封：这些消息都是本地构造的，
// 序列化失败意味着程序有 bug，不该被静默吞掉。
func mustMarshal(m proto.Message) []byte {
	data, err := proto.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("序列化出站消息失败: %v", err))
	}
	return data
}
