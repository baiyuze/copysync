package p2p_test

import (
	"context"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
	// 测试内已有一个名为 relay 的信令转发桩，这里给 TURN 包起别名以示区分
	turnrelay "github.com/baiyuze/copysync/server/relay"
)

// TestRelayFallback 验证「打洞失败就走服务器中转」这条兜底路径确实可用。
//
// 本地回环环境下 host candidate 总能直连，测不出真实的打洞失败，
// 因此用 ICETransportPolicy=relay 强制只走 TURN——若这条路能通，
// 真实网络中打洞失败时的降级就有保障。
func TestRelayFallback(t *testing.T) {
	turnSrv, err := turnrelay.New(turnrelay.Options{
		PublicIP:      "127.0.0.1",
		Port:          34780, // 避开默认 3478，防止与本机已运行的服务冲突
		Secret:        "test-secret",
		CredentialTTL: time.Hour,
	})
	if err != nil {
		t.Skipf("无法启动 TURN 服务（端口可能被占用）: %v", err)
	}
	defer turnSrv.Close()

	r := newRelay()
	alice := newRelayHarness(t, "alice", r, turnSrv, true)
	bob := newRelayHarness(t, "bob", r, turnSrv, true)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if err := alice.manager.Connect(ctx, "bob"); err != nil {
		t.Fatalf("发起连接: %v", err)
	}

	state := waitConnected(t, alice, "bob")
	waitConnected(t, bob, "alice")

	// 强制 relay 策略下，连接方式必须被识别为中转——
	// 若这里报 direct，说明 detectConnectionKind 的判定是错的，
	// 用户界面上就会把中转误显示成直连。
	if state != p2p.StateRelay {
		t.Errorf("连接方式 = %v，强制 relay 时应为 relay", state)
	}

	// 中转链路同样要能正常收发
	if err := alice.manager.Send("bob", &pb.PeerMessage{
		Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: "via-relay"}},
	}); err != nil {
		t.Fatalf("经中转发送失败: %v", err)
	}
	select {
	case got := <-bob.messages:
		if got.GetFetch().GetClipId() != "via-relay" {
			t.Errorf("收到的 clip_id = %q", got.GetFetch().GetClipId())
		}
	case <-time.After(10 * time.Second):
		t.Fatal("经中转的消息未送达")
	}
}

// TestRelayOneSided 验证只有一端走中转时，两端都识别为中转。
//
// 真实网络里常见的情形：一端在对称 NAT 后面打不通，用自己的 TURN 中转地址发送；
// 另一端仍从普通地址收发。只看本端 candidate 的话，后者会把中转误报成直连，
// 两台设备一台显示「直连」一台显示「中转」。
func TestRelayOneSided(t *testing.T) {
	turnSrv, err := turnrelay.New(turnrelay.Options{
		PublicIP:      "127.0.0.1",
		Port:          34781,
		Secret:        "test-secret",
		CredentialTTL: time.Hour,
	})
	if err != nil {
		t.Skipf("无法启动 TURN 服务（端口可能被占用）: %v", err)
	}
	defer turnSrv.Close()

	r := newRelay()
	alice := newRelayHarness(t, "alice", r, turnSrv, true) // 只能走中转
	bob := newRelayHarness(t, "bob", r, turnSrv, false)    // 普通设备

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if err := alice.manager.Connect(ctx, "bob"); err != nil {
		t.Fatalf("发起连接: %v", err)
	}

	waitConnected(t, alice, "bob")
	waitConnected(t, bob, "alice")

	// bob 自己看到的对端地址可能只是个普通地址，要等 alice 发来 LinkInfo 才改报中转，
	// 所以这里等最终状态，而不是第一次报告的状态
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) &&
		(alice.manager.State("bob") != p2p.StateRelay || bob.manager.State("alice") != p2p.StateRelay) {
		time.Sleep(50 * time.Millisecond)
	}
	if s := alice.manager.State("bob"); s != p2p.StateRelay {
		t.Errorf("alice 的连接方式 = %v，应为 relay", s)
	}
	if s := bob.manager.State("alice"); s != p2p.StateRelay {
		t.Errorf("bob 的连接方式 = %v，对端经中转发送时也应识别为 relay", s)
	}
}

func newRelayHarness(t *testing.T, id string, r *relay, turnSrv *turnrelay.Server, forceRelay bool) *harness {
	t.Helper()

	h := &harness{
		id:       id,
		states:   make(chan p2p.ConnState, 16),
		messages: make(chan *pb.PeerMessage, 16),
		streams:  nil,
	}
	h.manager = p2p.NewManager(p2p.ManagerOptions{
		SelfID:     id,
		SendSignal: r.send(id),
		ForceRelay: forceRelay,
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
	})

	// 每台设备各自签发凭证，与真实流程一致
	cred := turnSrv.Issue(id)
	h.manager.SetICEServers([]*pb.IceServer{{
		Urls:       cred.URLs,
		Username:   cred.Username,
		Credential: cred.Password,
	}})

	r.register(id, h.manager)
	t.Cleanup(h.manager.CloseAll)
	return h
}
