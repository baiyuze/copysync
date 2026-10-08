package p2p_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
	turnrelay "github.com/baiyuze/copysync/server/relay"
)

// TestManagerAdvertisesProbedEgress 验证：管理器在共享端口上做出口探测，
// 把探测到的公网地址作为 srflx candidate 发给对端，连接在共享端口上照常建立。
//
// 本机回环没有 NAT，探测到的「公网地址」就是 127.0.0.1 加上共享端口，
// 这正好能核对两件事：探测确实发生在共享端口上，结果确实发给了对端。
func TestManagerAdvertisesProbedEgress(t *testing.T) {
	stunSrv, err := turnrelay.New(turnrelay.Options{
		PublicIP: "127.0.0.1", Port: 34782, Secret: "test-secret", CredentialTTL: time.Hour,
	})
	if err != nil {
		t.Skipf("无法启动 STUN 服务（端口可能被占用）: %v", err)
	}
	defer stunSrv.Close()

	r := newRelay()
	var mu sync.Mutex
	var sentByAlice []string
	spy := func(from string) func(context.Context, string, *pb.SignalMessage) error {
		forward := r.send(from)
		return func(ctx context.Context, to string, msg *pb.SignalMessage) error {
			if c := msg.GetCandidate(); c != nil && from == "alice" {
				mu.Lock()
				sentByAlice = append(sentByAlice, c.GetCandidate())
				mu.Unlock()
			}
			return forward(ctx, to, msg)
		}
	}

	newProbingHarness := func(id string) *harness {
		h := &harness{id: id, states: make(chan p2p.ConnState, 16), messages: make(chan *pb.PeerMessage, 16)}
		h.manager = p2p.NewManager(p2p.ManagerOptions{
			SelfID:            id,
			SendSignal:        spy(id),
			ProbeAllowPrivate: true, // 探测服务器在回环地址上
			OnState: func(_ string, s p2p.ConnState) {
				select {
				case h.states <- s:
				default:
				}
			},
		})
		h.manager.SetICEServers([]*pb.IceServer{{Urls: []string{"stun:127.0.0.1:34782"}}})
		r.register(id, h.manager)
		t.Cleanup(h.manager.Close)
		return h
	}
	alice := newProbingHarness("alice")
	bob := newProbingHarness("bob")

	// 等首轮探测结束，核对结果
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) && len(alice.manager.Egress().Egresses) == 0 {
		time.Sleep(50 * time.Millisecond)
	}
	egress := alice.manager.Egress()
	if len(egress.Egresses) != 1 {
		t.Fatalf("探测结果 = %+v，期望恰好一个出口", egress)
	}
	e := egress.Egresses[0]
	if e.Addr.String() != fmt.Sprintf("127.0.0.1:%d", egress.LocalPort) || !e.PortPreserved {
		t.Fatalf("出口 = %v，期望 127.0.0.1:%d 且不改端口（回环没有 NAT）", e.Addr, egress.LocalPort)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := alice.manager.Connect(ctx, "bob"); err != nil {
		t.Fatalf("发起连接: %v", err)
	}
	waitConnected(t, alice, "bob")
	waitConnected(t, bob, "alice")

	// 出口地址必须作为 srflx candidate 发给了对端
	want := fmt.Sprintf("127.0.0.1 %d typ srflx", egress.LocalPort)
	mu.Lock()
	defer mu.Unlock()
	for _, c := range sentByAlice {
		if strings.Contains(c, want) {
			return
		}
	}
	t.Errorf("alice 没有把出口地址发给对端；发出的 candidate：\n%s", strings.Join(sentByAlice, "\n"))
}
