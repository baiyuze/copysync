// Package signaling 是信令服务器的客户端：维持长连接、交换 SDP/candidate、
// 走配对流程。
//
// 信任模型：服务器只是个转发通道。所有经它转发的消息都由发送方用设备私钥签名，
// 接收方用配对时存下的对端公钥验签——服务器改不了内容，也冒充不了任何一方。
// 验签失败的消息直接丢弃。签名与校验见 transport/signed。
//
// 服务器是可选的：地址为空时客户端什么也不做，设备只在局域网里互相发现。
package signaling

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/copysync/client-core/internal/identity"
	"github.com/baiyuze/copysync/client-core/internal/transport/signed"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

const (
	maxMessageBytes = 256 << 10
	dialTimeout     = 15 * time.Second
	writeTimeout    = 10 * time.Second
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

	verifier *signed.Verifier
	// DeviceInfo 提供每次握手时上报的设备名与已配对列表
	deviceInfo func() (name, platform string, pairedIDs []string)

	mu        sync.Mutex
	conn      *websocket.Conn
	url       string
	connected bool
	// urlChanged 在地址变化时唤醒正在等待的 Run，立即按新地址连接
	urlChanged chan struct{}
}

type Options struct {
	Identity *identity.Identity
	URL      string
	Logger   *slog.Logger
	Handlers Handlers
	// Verifier 校验转发来的信令。与局域网通道共用一个，同一条信令两条路都到时只处理一次。
	// 为空时按 PublicKeyOf 自建一个。
	Verifier *signed.Verifier
	// PublicKeyOf 返回已配对设备的公钥，用于验签。返回 nil 表示未配对，
	// 此时来自该设备的消息一律丢弃。
	PublicKeyOf func(deviceID string) ed25519.PublicKey
	DeviceInfo  func() (name, platform string, pairedIDs []string)
}

func New(opts Options) *Client {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	verifier := opts.Verifier
	if verifier == nil {
		verifier = signed.NewVerifier(opts.Identity.DeviceID, opts.PublicKeyOf)
	}
	return &Client{
		identity:   opts.Identity,
		url:        opts.URL,
		log:        log,
		handlers:   opts.Handlers,
		verifier:   verifier,
		deviceInfo: opts.DeviceInfo,
		urlChanged: make(chan struct{}, 1),
	}
}

// SetURL 更换信令服务器地址，立即生效。
// 当前连接会被断开，以便立即切换；地址为空表示不用服务器。
func (c *Client) SetURL(url string) {
	c.mu.Lock()
	changed := c.url != url
	c.url = url
	conn := c.conn
	c.mu.Unlock()

	if !changed {
		return
	}
	if conn != nil {
		conn.Close(websocket.StatusNormalClosure, "切换服务器")
	}
	select {
	case c.urlChanged <- struct{}{}:
	default:
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
// 没配地址时空等，直到 SetURL 给出地址。
func (c *Client) Run(ctx context.Context) {
	const initialBackoff = time.Second
	const maxBackoff = 30 * time.Second
	backoff := initialBackoff

	for ctx.Err() == nil {
		if c.URL() == "" {
			select {
			case <-ctx.Done():
				return
			case <-c.urlChanged:
				backoff = initialBackoff
				continue
			}
		}

		established, err := c.session(ctx)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			c.log.Debug("信令连接中断", "err", err, "重试间隔", backoff)
		}
		// 连上过就从头计时：服务器重启后应当很快重连，而不是沿用上次的长间隔
		if established {
			backoff = initialBackoff
		}

		// 加抖动，避免服务器重启后所有设备同时涌入
		jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
		select {
		case <-ctx.Done():
			return
		case <-c.urlChanged:
			backoff = initialBackoff
			continue
		case <-time.After(backoff + jitter):
		}
		if backoff < maxBackoff {
			backoff *= 2
		}
	}
}

// session 完成一次完整的连接生命周期。established 表示握手成功过。
func (c *Client) session(ctx context.Context) (established bool, err error) {
	url := c.URL()
	if url == "" {
		return false, errors.New("未配置信令服务器地址")
	}

	dialCtx, cancel := context.WithTimeout(ctx, dialTimeout)
	conn, _, err := websocket.Dial(dialCtx, url, nil)
	cancel()
	if err != nil {
		return false, fmt.Errorf("连接 %s: %w", url, err)
	}
	conn.SetReadLimit(maxMessageBytes)
	defer conn.Close(websocket.StatusNormalClosure, "")

	if err := c.handshake(ctx, conn); err != nil {
		return false, fmt.Errorf("握手: %w", err)
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

	err = c.readLoop(ctx, conn)
	return true, err
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
func (c *Client) handleForwarded(s *pb.Signed) {
	from, msg, err := c.verifier.Verify(s)
	switch {
	case errors.Is(err, signed.ErrReplay):
		// 局域网那条路往往先到，这份是重复的，正常现象
		c.log.Debug("丢弃重复的消息", "from", from)
		return
	case err != nil:
		c.log.Warn("丢弃转发来的消息", "from", from, "原因", err)
		return
	}

	if c.handlers.OnSignal != nil {
		c.handlers.OnSignal(from, msg)
	}
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

// Signal 签名后经服务器发给对端。
func (c *Client) Signal(ctx context.Context, to string, msg *pb.SignalMessage) error {
	s, err := signed.Sign(c.identity, to, msg)
	if err != nil {
		return err
	}
	return c.SendSigned(ctx, s)
}

// SendSigned 请服务器转发一条已签名的信令。
func (c *Client) SendSigned(ctx context.Context, s *pb.Signed) error {
	return c.send(ctx, &pb.ClientEnvelope{
		Payload: &pb.ClientEnvelope_Forward{Forward: s},
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
