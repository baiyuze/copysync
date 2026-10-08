// Package signaling 是信令服务器的客户端：维持长连接、交换 SDP/candidate、
// 走配对流程。
//
// 信任模型：服务器只是个转发通道。所有经它转发的消息都由发送方用设备私钥签名，
// 接收方用配对时存下的对端公钥验签——服务器改不了内容，也冒充不了任何一方。
// 验签失败的消息直接丢弃。
package signaling

import (
	"context"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/copysync/client-core/internal/identity"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

const (
	maxMessageBytes = 256 << 10
	dialTimeout     = 15 * time.Second
	writeTimeout    = 10 * time.Second
	// 防重放：拒绝时间戳偏离过大的消息
	maxClockSkew = 5 * time.Minute
)

// Handlers 是调用方关心的事件。均在客户端的读循环中同步调用，
// 实现里不要做耗时操作。
type Handlers struct {
	OnConnected       func()
	OnDisconnected    func(err error)
	OnPresence        func(deviceID string, online bool)
	OnSignal          func(from string, msg *pb.SignalMessage)
	OnPairingCreated  func(code string, expiresAt time.Time)
	OnPairingRedeemed func(peer *pb.PairingRedeemed)
	OnPairingResult   func(peer *pb.PairingRedeemResponse)
	OnServerError     func(code, message string)
	// OnIceConfig 在认证后收到，携带短期 TURN 凭证
	OnIceConfig func(servers []*pb.IceServer)
}

type Client struct {
	identity *identity.Identity
	log      *slog.Logger
	handlers Handlers

	// PublicKeyOf 返回已配对设备的公钥，用于验签。返回 nil 表示未配对，
	// 此时来自该设备的消息一律丢弃。
	publicKeyOf func(deviceID string) ed25519.PublicKey
	// DeviceInfo 提供每次握手时上报的设备名与已配对列表
	deviceInfo func() (name, platform string, pairedIDs []string)

	mu        sync.Mutex
	conn      *websocket.Conn
	url       string
	connected bool

	// 防重放：记录近期见过的 nonce
	seenNonces map[string]time.Time
}

type Options struct {
	Identity    *identity.Identity
	URL         string
	Logger      *slog.Logger
	Handlers    Handlers
	PublicKeyOf func(deviceID string) ed25519.PublicKey
	DeviceInfo  func() (name, platform string, pairedIDs []string)
}

func New(opts Options) *Client {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	return &Client{
		identity:    opts.Identity,
		url:         opts.URL,
		log:         log,
		handlers:    opts.Handlers,
		publicKeyOf: opts.PublicKeyOf,
		deviceInfo:  opts.DeviceInfo,
		seenNonces:  make(map[string]time.Time),
	}
}

// SetURL 更换信令服务器地址，下次重连生效。
// 当前连接会被断开，以便立即切换。
func (c *Client) SetURL(url string) {
	c.mu.Lock()
	changed := c.url != url
	c.url = url
	conn := c.conn
	c.mu.Unlock()

	if changed && conn != nil {
		conn.Close(websocket.StatusNormalClosure, "切换服务器")
	}
}

func (c *Client) URL() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.url
}

func (c *Client) Connected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

// Run 维持与信令服务器的连接，断线后自动重连，直到 ctx 取消。
func (c *Client) Run(ctx context.Context) {
	backoff := time.Second
	const maxBackoff = 30 * time.Second

	for ctx.Err() == nil {
		err := c.session(ctx)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			c.log.Debug("信令连接中断", "err", err, "重试间隔", backoff)
		}

		// 加抖动，避免服务器重启后所有设备同时涌入
		jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff + jitter):
		}
		if backoff < maxBackoff {
			backoff *= 2
		}
	}
}

// session 完成一次完整的连接生命周期。
func (c *Client) session(ctx context.Context) error {
	url := c.URL()
	if url == "" {
		return errors.New("未配置信令服务器地址")
	}

	dialCtx, cancel := context.WithTimeout(ctx, dialTimeout)
	conn, _, err := websocket.Dial(dialCtx, url, nil)
	cancel()
	if err != nil {
		return fmt.Errorf("连接 %s: %w", url, err)
	}
	conn.SetReadLimit(maxMessageBytes)
	defer conn.Close(websocket.StatusNormalClosure, "")

	if err := c.handshake(ctx, conn); err != nil {
		return fmt.Errorf("握手: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.connected = true
	c.mu.Unlock()

	if c.handlers.OnConnected != nil {
		c.handlers.OnConnected()
	}

	defer func() {
		c.mu.Lock()
		c.conn = nil
		c.connected = false
		c.mu.Unlock()
		if c.handlers.OnDisconnected != nil {
			c.handlers.OnDisconnected(err)
		}
	}()

	return c.readLoop(ctx, conn)
}

// handshake 应答服务器的 challenge，证明本机持有该公钥对应的私钥。
func (c *Client) handshake(ctx context.Context, conn *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(ctx, dialTimeout)
	defer cancel()

	var env pb.ServerEnvelope
	if err := readProto(ctx, conn, &env); err != nil {
		return fmt.Errorf("读取 ServerHello: %w", err)
	}
	hello := env.GetHello()
	if hello == nil {
		return errors.New("服务器首条消息不是 ServerHello")
	}

	name, platform, paired := c.deviceInfo()
	return writeProto(ctx, conn, &pb.ClientEnvelope{
		Payload: &pb.ClientEnvelope_Hello{Hello: &pb.ClientHello{
			DeviceId:        c.identity.DeviceID,
			DeviceName:      name,
			Platform:        platform,
			PublicKey:       c.identity.PublicKey,
			Signature:       c.identity.Sign(hello.GetChallenge()),
			PairedDeviceIds: paired,
		}},
	})
}

func (c *Client) readLoop(ctx context.Context, conn *websocket.Conn) error {
	for {
		var env pb.ServerEnvelope
		if err := readProto(ctx, conn, &env); err != nil {
			return err
		}

		switch p := env.GetPayload().(type) {
		case *pb.ServerEnvelope_Presence:
			if c.handlers.OnPresence != nil {
				c.handlers.OnPresence(p.Presence.GetDeviceId(), p.Presence.GetOnline())
			}

		case *pb.ServerEnvelope_Forwarded:
			c.handleForwarded(p.Forwarded)

		case *pb.ServerEnvelope_PairingCreated:
			if c.handlers.OnPairingCreated != nil {
				c.handlers.OnPairingCreated(
					p.PairingCreated.GetCode(),
					time.Unix(p.PairingCreated.GetExpiresAtUnix(), 0))
			}

		case *pb.ServerEnvelope_PairingRedeemed:
			if c.handlers.OnPairingRedeemed != nil {
				c.handlers.OnPairingRedeemed(p.PairingRedeemed)
			}

		case *pb.ServerEnvelope_PairingRedeemedResult:
			if c.handlers.OnPairingResult != nil {
				c.handlers.OnPairingResult(p.PairingRedeemedResult)
			}

		case *pb.ServerEnvelope_IceConfig:
			if c.handlers.OnIceConfig != nil {
				c.handlers.OnIceConfig(p.IceConfig.GetIceServers())
			}

		case *pb.ServerEnvelope_Error:
			c.log.Warn("服务器返回错误",
				"code", p.Error.GetCode(), "msg", p.Error.GetMessage())
			if c.handlers.OnServerError != nil {
				c.handlers.OnServerError(p.Error.GetCode(), p.Error.GetMessage())
			}
		}
	}
}

// handleForwarded 验签后再交给上层。
//
// 这是整个信任模型的落点：服务器能转发任意字节，但伪造不出对端的签名。
// 任何验签不过的消息都当作攻击丢弃，不向上层暴露。
func (c *Client) handleForwarded(signed *pb.Signed) {
	from := signed.GetFromDeviceId()
	pub := c.publicKeyOf(from)
	if pub == nil {
		c.log.Warn("丢弃未配对设备的消息", "from", from)
		return
	}
	if !ed25519.Verify(pub, signed.GetPayload(), signed.GetSignature()) {
		c.log.Warn("丢弃签名无效的消息", "from", from)
		return
	}

	var msg pb.SignalMessage
	if err := proto.Unmarshal(signed.GetPayload(), &msg); err != nil {
		c.log.Warn("丢弃无法解析的消息", "from", from, "err", err)
		return
	}

	// 防重放：时间戳超出容忍范围，或 nonce 重复出现过
	ts := time.Unix(msg.GetTimestampUnix(), 0)
	if skew := time.Since(ts); skew > maxClockSkew || skew < -maxClockSkew {
		c.log.Warn("丢弃时间戳异常的消息", "from", from, "偏差", skew)
		return
	}
	if c.seenNonce(msg.GetNonce()) {
		c.log.Warn("丢弃重放的消息", "from", from)
		return
	}

	if c.handlers.OnSignal != nil {
		c.handlers.OnSignal(from, &msg)
	}
}

// seenNonce 记录并判断 nonce 是否重复，同时清理过期条目。
func (c *Client) seenNonce(nonce []byte) bool {
	if len(nonce) == 0 {
		return false
	}
	key := string(nonce)

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if _, dup := c.seenNonces[key]; dup {
		return true
	}
	// 超出时钟容忍窗口的 nonce 不可能再被接受，可以安全丢弃
	for k, t := range c.seenNonces {
		if now.Sub(t) > maxClockSkew*2 {
			delete(c.seenNonces, k)
		}
	}
	c.seenNonces[key] = now
	return false
}

// ─────────────────────────── 发送 ───────────────────────────

var ErrNotConnected = errors.New("信令未连接")

func (c *Client) send(ctx context.Context, env *pb.ClientEnvelope) error {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn == nil {
		return ErrNotConnected
	}
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	return writeProto(ctx, conn, env)
}

// Signal 向对端发送一条已签名的信令。
func (c *Client) Signal(ctx context.Context, to string, msg *pb.SignalMessage) error {
	msg.ToDeviceId = to
	msg.TimestampUnix = time.Now().Unix()
	if msg.Nonce == nil {
		nonce := make([]byte, 16)
		if _, err := cryptorand.Read(nonce); err != nil {
			return err
		}
		msg.Nonce = nonce
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return c.send(ctx, &pb.ClientEnvelope{
		Payload: &pb.ClientEnvelope_Forward{Forward: &pb.Signed{
			Payload:      payload,
			Signature:    c.identity.Sign(payload),
			FromDeviceId: c.identity.DeviceID,
		}},
	})
}

func (c *Client) CreatePairingCode(ctx context.Context) error {
	name, _, _ := c.deviceInfo()
	return c.send(ctx, &pb.ClientEnvelope{
		Payload: &pb.ClientEnvelope_PairingCreate{PairingCreate: &pb.PairingCreateRequest{
			DeviceName: name,
			PublicKey:  c.identity.PublicKey,
		}},
	})
}

func (c *Client) RedeemPairingCode(ctx context.Context, code string) error {
	name, _, _ := c.deviceInfo()
	return c.send(ctx, &pb.ClientEnvelope{
		Payload: &pb.ClientEnvelope_PairingRedeem{PairingRedeem: &pb.PairingRedeemRequest{
			Code:       code,
			DeviceId:   c.identity.DeviceID,
			DeviceName: name,
			PublicKey:  c.identity.PublicKey,
		}},
	})
}

// RefreshPaired 在配对关系变化后重新上报关注列表，
// 这样服务器才会推送新配对设备的在线状态。
func (c *Client) RefreshPaired(ctx context.Context) error {
	name, platform, paired := c.deviceInfo()
	return c.send(ctx, &pb.ClientEnvelope{
		Payload: &pb.ClientEnvelope_Hello{Hello: &pb.ClientHello{
			DeviceId:        c.identity.DeviceID,
			DeviceName:      name,
			Platform:        platform,
			PublicKey:       c.identity.PublicKey,
			PairedDeviceIds: paired,
		}},
	})
}

// ─────────────────────── proto over websocket ───────────────────────

func readProto(ctx context.Context, conn *websocket.Conn, m proto.Message) error {
	typ, data, err := conn.Read(ctx)
	if err != nil {
		return err
	}
	if typ != websocket.MessageBinary {
		return fmt.Errorf("期望二进制消息，收到 %v", typ)
	}
	return proto.Unmarshal(data, m)
}

func writeProto(ctx context.Context, conn *websocket.Conn, m proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return conn.Write(ctx, websocket.MessageBinary, data)
}
