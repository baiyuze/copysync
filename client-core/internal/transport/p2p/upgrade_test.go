package p2p

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pion/webrtc/v4"

	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// signalBus 把几个 Manager 直接对接起来。异步、不保序地投递，
// 比真实的信令服务器更苛刻：ICE 重启时新一代 candidate 常常抢在 offer/answer 之前到达。
type signalBus struct {
	mu      sync.Mutex
	targets map[string]*Manager
}

func (b *signalBus) send(from string) func(context.Context, string, *pb.SignalMessage) error {
	return func(_ context.Context, to string, msg *pb.SignalMessage) error {
		b.mu.Lock()
		target := b.targets[to]
		b.mu.Unlock()
		if target != nil {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				_ = target.HandleSignal(ctx, from, msg)
			}()
		}
		return nil
	}
}

type testNode struct {
	m        *Manager
	mu       sync.Mutex
	states   []ConnState
	messages chan *pb.PeerMessage
}

func (n *testNode) history() []ConnState {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]ConnState(nil), n.states...)
}

func (n *testNode) peer(id string) *Peer {
	n.m.mu.RLock()
	defer n.m.mu.RUnlock()
	return n.m.peers[id]
}

func newTestNode(t *testing.T, bus *signalBus, id string) *testNode {
	t.Helper()
	n := &testNode{messages: make(chan *pb.PeerMessage, 16)}
	n.m = NewManager(ManagerOptions{
		SelfID:     id,
		SendSignal: bus.send(id),
		OnState: func(_ string, s ConnState) {
			n.mu.Lock()
			n.states = append(n.states, s)
			n.mu.Unlock()
		},
		OnMessage: func(_ string, m *pb.PeerMessage) {
			select {
			case n.messages <- m:
			default:
			}
		},
	})
	n.m.SetICEServers(nil) // 回环环境下 host candidate 足以直连
	bus.mu.Lock()
	bus.targets[id] = n.m
	bus.mu.Unlock()
	t.Cleanup(n.m.Close)
	return n
}

// connectPair 让 alice（不礼让、负责发起）与 bob 直连，返回两端。
func connectPair(t *testing.T) (alice, bob *testNode) {
	t.Helper()
	bus := &signalBus{targets: map[string]*Manager{}}
	alice = newTestNode(t, bus, "alice")
	bob = newTestNode(t, bus, "bob")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := alice.m.Connect(ctx, "bob"); err != nil {
		t.Fatalf("发起连接: %v", err)
	}
	waitState(t, alice, "bob", StateDirect, 20*time.Second)
	waitState(t, bob, "alice", StateDirect, 20*time.Second)
	return alice, bob
}

func waitState(t *testing.T, n *testNode, peerID string, want ConnState, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if n.m.State(peerID) == want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("%s 看 %s：等待 %v 超时，当前 %v，经历 %v", n.m.selfID, peerID, want, n.m.State(peerID), n.history())
}

func localUfrag(pc *webrtc.PeerConnection) string {
	for _, line := range strings.Split(pc.LocalDescription().SDP, "\n") {
		if u, ok := strings.CutPrefix(strings.TrimSpace(line), "a=ice-ufrag:"); ok {
			return u
		}
	}
	return ""
}

func expectMessage(t *testing.T, n *testNode, clipID string) {
	t.Helper()
	timeout := time.After(15 * time.Second)
	for {
		select {
		case m := <-n.messages:
			if m.GetFetch().GetClipId() == clipID {
				return
			}
		case <-timeout:
			t.Fatalf("%s 没有收到 %s", n.m.selfID, clipID)
		}
	}
}

// TestICERestartKeepsConnection 验证 ICE 重启只是让传输停顿，不会断开：
// 两端状态不经过「连接中」「离线」，重启期间发出的消息也能送达。
func TestICERestartKeepsConnection(t *testing.T) {
	alice, bob := connectPair(t)
	ap, bp := alice.peer("bob"), bob.peer("alice")
	before := localUfrag(ap.pc)
	aliceSeen, bobSeen := len(alice.history()), len(bob.history())

	// 记下两端 ICE 的状态变化：先回到「检查中」再「已连通」，才说明重启真的走完了
	restarted := func(pc *webrtc.PeerConnection) <-chan struct{} {
		done := make(chan struct{})
		var once sync.Once
		checking := false
		pc.OnICEConnectionStateChange(func(s webrtc.ICEConnectionState) {
			switch s {
			case webrtc.ICEConnectionStateChecking:
				checking = true
			case webrtc.ICEConnectionStateConnected:
				if checking {
					once.Do(func() { close(done) })
				}
			}
		})
		return done
	}
	aliceDone, bobDone := restarted(ap.pc), restarted(bp.pc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := ap.restartICE(ctx); err != nil {
		t.Fatalf("ICE 重启: %v", err)
	}
	// 重启进行中就发：数据暂时丢失，由 SCTP 重传补上
	fetch := func(id string) *pb.PeerMessage {
		return &pb.PeerMessage{Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: id}}}
	}
	if err := alice.m.Send("bob", fetch("during-restart")); err != nil {
		t.Fatalf("重启期间发送: %v", err)
	}

	for name, done := range map[string]<-chan struct{}{"alice": aliceDone, "bob": bobDone} {
		select {
		case <-done:
		case <-time.After(15 * time.Second):
			t.Fatalf("%s 的 ICE 没有完成重启（当前 %v）", name, ap.pc.ICEConnectionState())
		}
	}
	if localUfrag(ap.pc) == before {
		t.Fatal("本地 ufrag 没有变化，ICE 没有重启")
	}
	expectMessage(t, bob, "during-restart")

	if err := bob.m.Send("alice", fetch("after-restart")); err != nil {
		t.Fatalf("重启后发送: %v", err)
	}
	expectMessage(t, alice, "after-restart")

	for _, s := range alice.history()[aliceSeen:] {
		if s != StateDirect {
			t.Errorf("alice 在重启期间报告了 %v，应一直是直连", s)
		}
	}
	for _, s := range bob.history()[bobSeen:] {
		if s != StateDirect {
			t.Errorf("bob 在重启期间报告了 %v，应一直是直连", s)
		}
	}
}

// TestUpgradeFromRelay 验证走中转时，发起方会自己做 ICE 重启，打通后两端都改回直连。
//
// 回环环境里造不出真正打不通的网络，这里让 bob 谎报「我这边走中转」，
// 两端随之显示中转；重启后 bob 如实报告，两端回到直连。
// 真实的「先中转、后直连」由 NAT 实验室的「出口地址晚到」场景验证。
func TestUpgradeFromRelay(t *testing.T) {
	saveDelays, saveIdle := upgradeDelays, upgradeIdle
	upgradeDelays, upgradeIdle = []time.Duration{300 * time.Millisecond}, 100*time.Millisecond
	t.Cleanup(func() { upgradeDelays, upgradeIdle = saveDelays, saveIdle })

	alice, bob := connectPair(t)
	ap := alice.peer("bob")

	if err := bob.peer("alice").Send(&pb.PeerMessage{Payload: &pb.PeerMessage_Link{
		Link: &pb.LinkInfo{Relayed: true},
	}}); err != nil {
		t.Fatalf("发送 LinkInfo: %v", err)
	}

	waitState(t, alice, "bob", StateRelay, 5*time.Second)
	waitState(t, alice, "bob", StateDirect, 20*time.Second)

	ap.mu.Lock()
	restarts := ap.restarts
	ap.mu.Unlock()
	if restarts == 0 {
		t.Fatal("alice 没有做 ICE 重启就回到了直连")
	}
	if s := bob.m.State("alice"); s != StateDirect {
		t.Errorf("bob 的状态 = %v，应为直连", s)
	}
}

// TestPoliteSideWaitsForPeer 验证礼让的一方不主动重启，免得两端同时发起。
func TestPoliteSideWaitsForPeer(t *testing.T) {
	saveDelays := upgradeDelays
	upgradeDelays = []time.Duration{100 * time.Millisecond}
	t.Cleanup(func() { upgradeDelays = saveDelays })

	alice, bob := connectPair(t)
	if err := alice.peer("bob").Send(&pb.PeerMessage{Payload: &pb.PeerMessage_Link{
		Link: &pb.LinkInfo{Relayed: true},
	}}); err != nil {
		t.Fatalf("发送 LinkInfo: %v", err)
	}
	waitState(t, bob, "alice", StateRelay, 5*time.Second)
	time.Sleep(time.Second)

	bp := bob.peer("alice")
	bp.mu.Lock()
	restarts := bp.restarts
	bp.mu.Unlock()
	if restarts != 0 {
		t.Errorf("礼让方做了 %d 次 ICE 重启，应该等对方发起", restarts)
	}
}

func TestBusy(t *testing.T) {
	p := &Peer{}
	if p.busy() {
		t.Fatal("新连接不应算忙")
	}
	p.touch()
	if !p.busy() {
		t.Fatal("刚有过消息往来，应算忙")
	}
	p.lastActivity = time.Now().Add(-2 * upgradeIdle)
	if p.busy() {
		t.Fatal("空闲够久，不应算忙")
	}
}

func TestCandidateUfrag(t *testing.T) {
	cases := map[string]string{
		"candidate:1 1 udp 2130706431 192.168.1.2 50000 typ host ufrag AbCd network-cost 999": "AbCd",
		"candidate:2 1 udp 1694498815 1.2.3.4 50000 typ srflx raddr 0.0.0.0 rport 50000":      "",
		// 地址、端口等位置上即使碰巧是 "ufrag" 也不算
		"candidate:ufrag 1 udp 1 ufrag 50000 typ host": "",
	}
	for c, want := range cases {
		if got := candidateUfrag(c); got != want {
			t.Errorf("candidateUfrag(%q) = %q，应为 %q", c, got, want)
		}
	}

	sdp := "v=0\r\na=ice-ufrag:AbCd\r\na=ice-pwd:x\r\n"
	if !sdpHasUfrag(sdp, "AbCd") || !sdpHasUfrag(sdp, "") {
		t.Error("应匹配")
	}
	if sdpHasUfrag(sdp, "Ab") || sdpHasUfrag(sdp, "Other") {
		t.Error("不应匹配")
	}
}

// TestRebuildWhenPeerRestarts 验证对端重建了连接（例如重启了后台服务）时，
// 本端丢弃旧连接、跟着建新的，而不是在旧连接上协商到超时。
func TestRebuildWhenPeerRestarts(t *testing.T) {
	alice, bob := connectPair(t)
	bp := bob.peer("alice")
	bp.mu.Lock()
	oldBob := bp.pc
	bp.mu.Unlock()

	// alice 悄悄丢下旧连接、建一个新的：模拟它重启，bob 对此一无所知
	ap := alice.peer("bob")
	ap.mu.Lock()
	oldAlice := ap.pc
	ap.pc, ap.control, ap.pcConnected = nil, nil, false
	ap.mu.Unlock()
	t.Cleanup(func() { _ = oldAlice.Close() })
	ap.setState(StateOffline)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := alice.m.Connect(ctx, "bob"); err != nil {
		t.Fatalf("重新发起: %v", err)
	}
	waitState(t, alice, "bob", StateDirect, 20*time.Second)
	waitState(t, bob, "alice", StateDirect, 20*time.Second)

	bp.mu.Lock()
	rebuilt := bp.pc != oldBob
	bp.mu.Unlock()
	if !rebuilt {
		t.Error("bob 还在用旧连接")
	}
	if err := alice.m.Send("bob", &pb.PeerMessage{Payload: &pb.PeerMessage_Fetch{
		Fetch: &pb.PeerFetch{ClipId: "after-rebuild"}}}); err != nil {
		t.Fatalf("发送: %v", err)
	}
	expectMessage(t, bob, "after-rebuild")
}

// TestSendRefusedWhileDown 验证连接断开期间发送会报错，而不是写进一个收不到的通道。
func TestSendRefusedWhileDown(t *testing.T) {
	alice, _ := connectPair(t)
	ap := alice.peer("bob")
	ap.setState(StateOffline) // 模拟 ICE 断开：控制通道此时仍显示「打开」
	err := alice.m.Send("bob", &pb.PeerMessage{Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{}}})
	if err != ErrNotConnected {
		t.Errorf("断开时发送返回 %v，应为 ErrNotConnected", err)
	}
}
