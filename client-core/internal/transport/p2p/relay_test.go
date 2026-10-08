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
	alice := newRelayHarness(t, "alice", r, turnSrv)
	bob := newRelayHarness(t, "bob", r, turnSrv)

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

func newRelayHarness(t *testing.T, id string, r *relay, turnSrv *turnrelay.Server) *harness {
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
		ForceRelay: true,
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
