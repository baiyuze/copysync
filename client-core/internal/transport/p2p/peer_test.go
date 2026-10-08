package p2p_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// relay 把两个 Manager 直接对接起来，替代真实的信令服务器。
//
// 这样测试聚焦于 WebRTC 握手本身，不受网络与信令实现干扰；
// 信令链路已由 signaling 包的集成测试覆盖。
type relay struct {
	mu      sync.Mutex
	targets map[string]*p2p.Manager
}

func newRelay() *relay {
	return &relay{targets: make(map[string]*p2p.Manager)}
}

func (r *relay) register(id string, m *p2p.Manager) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.targets[id] = m
}

// send 模拟信令服务器的转发：异步投递，贴近真实的网络时序。
func (r *relay) send(from string) func(context.Context, string, *pb.SignalMessage) error {
	return func(ctx context.Context, to string, msg *pb.SignalMessage) error {
		r.mu.Lock()
		target := r.targets[to]
		r.mu.Unlock()
		if target == nil {
			return nil
		}
		go func() {
			bg, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			_ = target.HandleSignal(bg, from, msg)
		}()
		return nil
	}
}

type harness struct {
	id       string
	manager  *p2p.Manager
	states   chan p2p.ConnState
	messages chan *pb.PeerMessage
	streams  chan *p2p.Stream
}

func newHarness(t *testing.T, id string, r *relay) *harness {
	t.Helper()

	h := &harness{
		id:       id,
		states:   make(chan p2p.ConnState, 16),
		messages: make(chan *pb.PeerMessage, 16),
		streams:  make(chan *p2p.Stream, 8),
	}
	h.manager = p2p.NewManager(p2p.ManagerOptions{
		SelfID:     id,
		SendSignal: r.send(id),
		OnState: func(_ string, s p2p.ConnState) {
			select {
			case h.states <- s:
			default:
			}
		},
		OnMessage: func(_ string, m *pb.PeerMessage) {
			select {
			case h.messages <- m:
			default:
			}
		},
		OnStream: func(_, _ string, st *p2p.Stream) {
			select {
			case h.streams <- st:
			default:
			}
		},
	})
	// 本地回环测试不需要 STUN/TURN：host candidate 足以直连
	h.manager.SetICEServers(nil)
	r.register(id, h.manager)
	t.Cleanup(h.manager.CloseAll)
	return h
}

// waitConnected 等到握手完成（直连或中转皆可）。
func waitConnected(t *testing.T, h *harness, peerID string) p2p.ConnState {
	t.Helper()
	deadline := time.After(30 * time.Second)
	for {
		select {
		case s := <-h.states:
			if s == p2p.StateDirect || s == p2p.StateRelay {
				return s
			}
		case <-deadline:
			t.Fatalf("%s 未能在 30 秒内与 %s 建立连接（当前 %v）",
				h.id, peerID, h.manager.State(peerID))
			return p2p.StateOffline
		}
	}
}

// TestPeersConnectAndExchangeMessages 是 M2 的核心验收：
// 两端能真正建立 WebRTC 连接，并在控制通道上互发消息。
func TestPeersConnectAndExchangeMessages(t *testing.T) {
	r := newRelay()
	// device_id 的字典序决定谁发起、谁礼让，这里刻意让 alice < bob
	alice := newHarness(t, "alice", r)
	bob := newHarness(t, "bob", r)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if err := alice.manager.Connect(ctx, "bob"); err != nil {
		t.Fatalf("发起连接: %v", err)
	}

	aliceState := waitConnected(t, alice, "bob")
	waitConnected(t, bob, "alice")
	t.Logf("连接方式：%v", aliceState)

	// 控制通道：alice → bob
	offer := &pb.PeerMessage{
		Payload: &pb.PeerMessage_Offer{Offer: &pb.ClipOffer{
			ClipId:    "clip-1",
			Kind:      pb.ClipKind_CLIP_KIND_TEXT,
			TotalSize: 11,
			// 文本内容直接内联，不需要二次传输
			TextContent: "hello peer",
			WillPush:    true,
		}},
	}
	if err := alice.manager.Send("bob", offer); err != nil {
		t.Fatalf("发送控制消息: %v", err)
	}

	select {
	case got := <-bob.messages:
		if got.GetOffer().GetClipId() != "clip-1" {
			t.Errorf("clip_id = %q，期望 clip-1", got.GetOffer().GetClipId())
		}
		if got.GetOffer().GetTextContent() != "hello peer" {
			t.Errorf("内容 = %q，期望 hello peer", got.GetOffer().GetTextContent())
		}
	case <-time.After(10 * time.Second):
		t.Fatal("bob 未收到控制消息")
	}

	// 反方向也要通：拉取请求是由接收方发起的
	if err := bob.manager.Send("alice", &pb.PeerMessage{
		Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: "clip-1"}},
	}); err != nil {
		t.Fatalf("回发控制消息: %v", err)
	}
	select {
	case got := <-alice.messages:
		if got.GetFetch().GetClipId() != "clip-1" {
			t.Errorf("fetch clip_id = %q", got.GetFetch().GetClipId())
		}
	case <-time.After(10 * time.Second):
		t.Fatal("alice 未收到回发的控制消息")
	}
}

// TestSeparateStreamChannel 验证数据流与控制通道相互独立——
// 大文件传输不能把控制消息堵在队列后面。
func TestSeparateStreamChannel(t *testing.T) {
	r := newRelay()
	alice := newHarness(t, "alice", r)
	bob := newHarness(t, "bob", r)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if err := alice.manager.Connect(ctx, "bob"); err != nil {
		t.Fatalf("发起连接: %v", err)
	}
	waitConnected(t, alice, "bob")
	waitConnected(t, bob, "alice")

	dc, err := alice.manager.OpenStream("bob", "clip-1")
	if err != nil {
		t.Fatalf("开数据流: %v", err)
	}
	if dc.Label() != "clip-1" {
		t.Errorf("数据流 label = %q，期望 clip-1", dc.Label())
	}

	select {
	case got := <-bob.streams:
		if got.Label() != "clip-1" {
			t.Errorf("对端收到的 label = %q，期望 clip-1", got.Label())
		}
		// 数据流不该被当成控制消息处理
		select {
		case m := <-bob.messages:
			t.Errorf("数据流被误当作控制消息: %v", m)
		case <-time.After(300 * time.Millisecond):
		}
	case <-time.After(10 * time.Second):
		t.Fatal("bob 未收到新开的数据流")
	}
}

// TestSimultaneousConnectIsSafe 验证双方同时调用 Connect 不会互相干扰。
//
// 真实场景就是这样：两端几乎同时收到对方的上线通知，各自调用 Connect。
// Manager 内部按 device_id 字典序决定由谁发起，因此不会出现两个 offer 相撞。
func TestSimultaneousConnectIsSafe(t *testing.T) {
	r := newRelay()
	alice := newHarness(t, "alice", r)
	bob := newHarness(t, "bob", r)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); _ = alice.manager.Connect(ctx, "bob") }()
	go func() { defer wg.Done(); _ = bob.manager.Connect(ctx, "alice") }()
	wg.Wait()

	waitConnected(t, alice, "bob")
	waitConnected(t, bob, "alice")

	// 连接可用即证明没有因竞争而留下半开状态
	if err := alice.manager.Send("bob", &pb.PeerMessage{
		Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: "after-race"}},
	}); err != nil {
		t.Fatalf("并发 Connect 后仍无法发送: %v", err)
	}
	select {
	case got := <-bob.messages:
		if got.GetFetch().GetClipId() != "after-race" {
			t.Errorf("收到的 clip_id = %q", got.GetFetch().GetClipId())
		}
	case <-time.After(10 * time.Second):
		t.Fatal("并发 Connect 后消息未送达")
	}
}

// TestConnectFromAnswererSideAlone 验证只有应答方调用 Connect 时不会误发 offer。
// 它应当安静地等待，而不是自己也发起一轮协商。
func TestConnectFromAnswererSideAlone(t *testing.T) {
	r := newRelay()
	// "zoe" > "alice"，所以 zoe 是应答方
	zoe := newHarness(t, "zoe", r)
	newHarness(t, "alice", r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := zoe.manager.Connect(ctx, "alice"); err != nil {
		t.Fatalf("应答方调用 Connect 不该报错: %v", err)
	}
	// 没有对方发起，就应当一直停在 offline，而不是进入 connecting
	select {
	case s := <-zoe.states:
		t.Errorf("应答方不该主动改变状态，却变成了 %v", s)
	case <-time.After(1500 * time.Millisecond):
	}
}
