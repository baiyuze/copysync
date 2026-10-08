package signaling_test

import (
	"context"
	"crypto/ed25519"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/identity"
	"github.com/baiyuze/copysync/client-core/internal/transport/signaling"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
	"github.com/baiyuze/copysync/server/signal"
)

// peer 把一个设备身份、它的信令客户端、以及它观察到的事件打包在一起，
// 方便在测试里像操作真实设备那样操作它。
type peer struct {
	name     string
	identity *identity.Identity
	client   *signaling.Client

	mu sync.Mutex
	// 已配对对端的公钥。真实实现中来自 SQLite，这里用内存模拟。
	known map[string]ed25519.PublicKey

	presence chan presenceEvent
	signals  chan signalEvent
	codes    chan string
	paired   chan pairedEvent
}

type presenceEvent struct {
	deviceID string
	online   bool
}

type signalEvent struct {
	from string
	msg  *pb.SignalMessage
}

type pairedEvent struct {
	deviceID  string
	name      string
	publicKey ed25519.PublicKey
}

func newPeer(t *testing.T, name, url string) *peer {
	t.Helper()

	id, err := identity.LoadOrCreate(filepath.Join(t.TempDir(), "identity.json"))
	if err != nil {
		t.Fatalf("创建身份: %v", err)
	}

	p := &peer{
		name:     name,
		identity: id,
		known:    make(map[string]ed25519.PublicKey),
		presence: make(chan presenceEvent, 16),
		signals:  make(chan signalEvent, 16),
		codes:    make(chan string, 4),
		paired:   make(chan pairedEvent, 4),
	}

	p.client = signaling.New(signaling.Options{
		Identity: id,
		URL:      url,
		PublicKeyOf: func(deviceID string) ed25519.PublicKey {
			p.mu.Lock()
			defer p.mu.Unlock()
			return p.known[deviceID]
		},
		DeviceInfo: func() (string, string, []string) {
			p.mu.Lock()
			defer p.mu.Unlock()
			ids := make([]string, 0, len(p.known))
			for k := range p.known {
				ids = append(ids, k)
			}
			return name, "test", ids
		},
		Handlers: signaling.Handlers{
			OnPresence: func(deviceID string, online bool) {
				p.presence <- presenceEvent{deviceID, online}
			},
			OnSignal: func(from string, msg *pb.SignalMessage) {
				p.signals <- signalEvent{from, msg}
			},
			OnPairingCreated: func(code string, _ time.Time) {
				p.codes <- code
			},
			// 发起方：有人兑换了我的配对码
			OnPairingRedeemed: func(ev *pb.PairingRedeemed) {
				p.remember(ev.GetPeerDeviceId(), ev.GetPeerPublicKey())
				p.paired <- pairedEvent{
					ev.GetPeerDeviceId(), ev.GetPeerDeviceName(), ev.GetPeerPublicKey(),
				}
			},
			// 兑换方：我兑换成功，拿到发起方信息
			OnPairingResult: func(ev *pb.PairingRedeemResponse) {
				p.remember(ev.GetPeerDeviceId(), ev.GetPeerPublicKey())
				p.paired <- pairedEvent{
					ev.GetPeerDeviceId(), ev.GetPeerDeviceName(), ev.GetPeerPublicKey(),
				}
			},
		},
	})
	return p
}

func (p *peer) remember(deviceID string, pub []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.known[deviceID] = ed25519.PublicKey(pub)
}

func (p *peer) forget(deviceID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.known, deviceID)
}

// startServer 起一个真实的信令服务器（httptest + websocket）。
func startServer(t *testing.T) string {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle("/signal", signal.NewServer(signal.Options{}))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return "ws" + strings.TrimPrefix(srv.URL, "http") + "/signal"
}

func waitConnected(t *testing.T, peers ...*peer) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for _, p := range peers {
		for !p.client.Connected() {
			if time.Now().After(deadline) {
				t.Fatalf("%s 未能连上信令服务器", p.name)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func recvPaired(t *testing.T, p *peer) pairedEvent {
	t.Helper()
	select {
	case ev := <-p.paired:
		return ev
	case <-time.After(5 * time.Second):
		t.Fatalf("%s 未收到配对结果", p.name)
		return pairedEvent{}
	}
}

// TestPairingExchangesKeysBothWays 验证配对的核心语义：
// 双方都要拿到对方公钥，否则只有单向互信，签名验证会失败。
func TestPairingExchangesKeysBothWays(t *testing.T) {
	url := startServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	alice := newPeer(t, "Alice", url)
	bob := newPeer(t, "Bob", url)
	go alice.client.Run(ctx)
	go bob.client.Run(ctx)
	waitConnected(t, alice, bob)

	if err := alice.client.CreatePairingCode(ctx); err != nil {
		t.Fatalf("生成配对码: %v", err)
	}
	var code string
	select {
	case code = <-alice.codes:
	case <-time.After(5 * time.Second):
		t.Fatal("未收到配对码")
	}
	if len(code) != 6 {
		t.Errorf("配对码长度应为 6，实际 %d（%q）", len(code), code)
	}

	if err := bob.client.RedeemPairingCode(ctx, code); err != nil {
		t.Fatalf("兑换配对码: %v", err)
	}

	// 兑换方拿到发起方信息
	got := recvPaired(t, bob)
	if got.deviceID != alice.identity.DeviceID {
		t.Errorf("Bob 拿到的对端 ID = %q，期望 %q", got.deviceID, alice.identity.DeviceID)
	}
	if got.name != "Alice" {
		t.Errorf("Bob 拿到的对端名 = %q，期望 Alice", got.name)
	}

	// 发起方也要拿到兑换方信息——否则 Alice 无法验证 Bob 的签名
	got = recvPaired(t, alice)
	if got.deviceID != bob.identity.DeviceID {
		t.Errorf("Alice 拿到的对端 ID = %q，期望 %q", got.deviceID, bob.identity.DeviceID)
	}

	// 双方存下的公钥必须与对方真实公钥一致
	alice.mu.Lock()
	bobKeyAtAlice := alice.known[bob.identity.DeviceID]
	alice.mu.Unlock()
	if !bobKeyAtAlice.Equal(bob.identity.PublicKey) {
		t.Error("Alice 存下的 Bob 公钥与真实公钥不一致")
	}
}

// TestSignalForwardingVerifiesSignature 验证信任模型的落点：
// 已配对设备间的消息能送达，未配对设备的消息必须被丢弃。
func TestSignalForwardingVerifiesSignature(t *testing.T) {
	url := startServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	alice := newPeer(t, "Alice", url)
	bob := newPeer(t, "Bob", url)

	// 直接互相植入公钥，跳过配对流程
	alice.remember(bob.identity.DeviceID, bob.identity.PublicKey)
	bob.remember(alice.identity.DeviceID, alice.identity.PublicKey)

	go alice.client.Run(ctx)
	go bob.client.Run(ctx)
	waitConnected(t, alice, bob)

	send := &pb.SignalMessage{
		Payload: &pb.SignalMessage_Offer{Offer: &pb.SdpOffer{
			Sdp:             "v=0 test-offer",
			DtlsFingerprint: "sha-256 AA:BB:CC",
		}},
	}
	if err := alice.client.Signal(ctx, bob.identity.DeviceID, send); err != nil {
		t.Fatalf("发送信令: %v", err)
	}

	select {
	case ev := <-bob.signals:
		if ev.from != alice.identity.DeviceID {
			t.Errorf("发送方 = %q，期望 %q", ev.from, alice.identity.DeviceID)
		}
		if got := ev.msg.GetOffer().GetSdp(); got != "v=0 test-offer" {
			t.Errorf("SDP = %q，期望 %q", got, "v=0 test-offer")
		}
		// DTLS 指纹经签名传递，是防中间人的锚点，必须完整到达
		if got := ev.msg.GetOffer().GetDtlsFingerprint(); got != "sha-256 AA:BB:CC" {
			t.Errorf("DTLS 指纹 = %q，期望 %q", got, "sha-256 AA:BB:CC")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Bob 未收到 Alice 的信令")
	}

	// Bob 忘掉 Alice 后，来自 Alice 的消息应被丢弃（模拟未配对设备的消息）
	bob.forget(alice.identity.DeviceID)
	if err := alice.client.Signal(ctx, bob.identity.DeviceID, &pb.SignalMessage{
		Payload: &pb.SignalMessage_Offer{Offer: &pb.SdpOffer{Sdp: "should-be-dropped"}},
	}); err != nil {
		t.Fatalf("发送信令: %v", err)
	}
	select {
	case ev := <-bob.signals:
		t.Errorf("未配对设备的消息本应被丢弃，却收到了：%v", ev.msg)
	case <-time.After(800 * time.Millisecond):
		// 符合预期
	}
}

// TestPresenceNotifiesWatchers 验证在线状态推送：
// 设备只会收到它关心（已配对）的对端的状态变化。
func TestPresenceNotifiesWatchers(t *testing.T) {
	url := startServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	alice := newPeer(t, "Alice", url)
	bob := newPeer(t, "Bob", url)
	alice.remember(bob.identity.DeviceID, bob.identity.PublicKey)
	bob.remember(alice.identity.DeviceID, alice.identity.PublicKey)

	go alice.client.Run(ctx)
	waitConnected(t, alice)

	// Bob 后上线，Alice 应当收到上线通知
	bobCtx, bobCancel := context.WithCancel(ctx)
	go bob.client.Run(bobCtx)
	waitConnected(t, bob)

	select {
	case ev := <-alice.presence:
		if ev.deviceID != bob.identity.DeviceID || !ev.online {
			t.Errorf("期望 Bob 上线通知，实际 %+v", ev)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Alice 未收到 Bob 的上线通知")
	}

	// Bob 断开，Alice 应当收到下线通知
	bobCancel()
	select {
	case ev := <-alice.presence:
		if ev.deviceID != bob.identity.DeviceID || ev.online {
			t.Errorf("期望 Bob 下线通知，实际 %+v", ev)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Alice 未收到 Bob 的下线通知")
	}
}

// TestRedeemInvalidCode 验证错误的配对码不会建立任何关系。
func TestRedeemInvalidCode(t *testing.T) {
	url := startServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bob := newPeer(t, "Bob", url)
	go bob.client.Run(ctx)
	waitConnected(t, bob)

	if err := bob.client.RedeemPairingCode(ctx, "ZZZZZZ"); err != nil {
		t.Fatalf("发送兑换请求: %v", err)
	}
	select {
	case ev := <-bob.paired:
		t.Errorf("无效配对码本不该成功，却返回了 %+v", ev)
	case <-time.After(1 * time.Second):
		// 符合预期：服务器回错误，不产生配对
	}
}
