package p2p

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/gateway"
	"github.com/baiyuze/pasteferry/client-core/internal/identity"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	turnrelay "github.com/baiyuze/pasteferry/server/relay"
)

func officialPair(t *testing.T, base string, backend *turnrelay.Server, dropDirect bool) (*Manager, *Manager, <-chan *pb.PeerMessage) {
	t.Helper()
	a := identity.FromSeed(make([]byte, 32))
	seed := sha256.Sum256([]byte("official p2p b"))
	b := identity.FromSeed(seed[:])
	ac, _ := gateway.NewClient(a, base, nil)
	bc, _ := gateway.NewClient(b, base, nil)
	t.Cleanup(ac.Close)
	t.Cleanup(bc.Close)
	var am, bm *Manager
	received := make(chan *pb.PeerMessage, 16)
	blackhole, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { blackhole.Close() })
	blackholePort := blackhole.LocalAddr().(*net.UDPAddr).Port
	deliver := func(from string, target **Manager) func(context.Context, string, *pb.SignalMessage) error {
		return func(ctx context.Context, to string, msg *pb.SignalMessage) error {
			if dropDirect && !msg.GetOfficialRelay() {
				if msg.GetCandidate() != nil {
					msg.GetCandidate().Candidate = fmt.Sprintf("candidate:1 1 udp 2130706431 127.0.0.1 %d typ host", blackholePort)
				}
				if offer := msg.GetOffer(); offer != nil {
					offer.Sdp = stripCandidates(offer.Sdp)
				}
				if answer := msg.GetAnswer(); answer != nil {
					answer.Sdp = stripCandidates(answer.Sdp)
				}
			}
			go func() { _ = (*target).HandleSignal(context.Background(), from, msg) }()
			return nil
		}
	}
	opts := func(id string, c *gateway.Client, send func(context.Context, string, *pb.SignalMessage) error) ManagerOptions {
		return ManagerOptions{SelfID: id, Gateway: c, Official: func() bool { return true }, OfficialOnline: func(string) bool { return true }, DisableEgressProbe: true, SendSignal: send, OnMessage: func(_ string, msg *pb.PeerMessage) {
			select {
			case received <- msg:
			default:
			}
		}}
	}
	am = NewManager(opts(a.DeviceID, ac, deliver(a.DeviceID, &bm)))
	bm = NewManager(opts(b.DeviceID, bc, deliver(b.DeviceID, &am)))
	t.Cleanup(am.Close)
	t.Cleanup(bm.Close)
	if backend != nil {
		for _, x := range []struct {
			m  *Manager
			id string
		}{{am, a.DeviceID}, {bm, b.DeviceID}} {
			cred := backend.Issue(x.id)
			x.m.SetICEServers([]*pb.IceServer{{Urls: cred.URLs, Username: cred.Username, Credential: cred.Password}})
		}
	}
	return am, bm, received
}
func stripCandidates(s string) string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "a=candidate:") {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}
func waitPair(t *testing.T, a, b *Manager, want ConnState, timeout time.Duration) {
	t.Helper()
	until := time.Now().Add(timeout)
	for time.Now().Before(until) {
		if a.State(b.selfID) == want && b.State(a.selfID) == want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("pair states %v/%v; expected %v", a.State(b.selfID), b.State(a.selfID), want)
}
func TestOfficialDirectNeverRequestsGateway(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { requests.Add(1); w.WriteHeader(500) }))
	defer server.Close()
	a, b, _ := officialPair(t, server.URL, nil, false)
	if a.selfID < b.selfID {
		_ = a.Connect(context.Background(), b.selfID)
	} else {
		_ = b.Connect(context.Background(), a.selfID)
	}
	waitPair(t, a, b, StateDirect, 10*time.Second)
	if requests.Load() != 0 {
		t.Fatal("direct path contacted relay authorization")
	}
}

func TestOfficialQuotaRecoveryAfterTransportClosure(t *testing.T) {
	id := identity.FromSeed(make([]byte, 32))
	c, err := gateway.NewClient(id, "http://127.0.0.1:1", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.ReportError(&gateway.Error{Code: "quota_exhausted", Usage: &gateway.Usage{Available: true, Reason: "quota_exhausted", ResetAt: time.Now().Add(time.Hour).Unix()}})
	m := NewManager(ManagerOptions{SelfID: id.DeviceID, Gateway: c, Official: func() bool { return true }, DisableEgressProbe: true, SendSignal: func(context.Context, string, *pb.SignalMessage) error { return nil }})
	defer m.Close()
	ctx, cancel := context.WithCancel(context.Background())
	a := &officialAttempt{id: "attempt", ctx: ctx, cancel: cancel}
	m.relays["peer"] = a
	m.officialFailed("peer", a, errors.New("relay connection closed"))
	if a.failReason != "quota_exhausted" {
		t.Fatalf("lost closure reason: %q", a.failReason)
	}
	m.OfficialAllowanceChanged()
	if !a.failedUntil.IsZero() || a.failReason != "" {
		t.Fatal("new allowance did not permit a retry")
	}
}

func TestStaleOfficialClosureDoesNotCancelReplacement(t *testing.T) {
	id := identity.FromSeed(make([]byte, 32))
	c, err := gateway.NewClient(id, "http://127.0.0.1:1", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	m := NewManager(ManagerOptions{SelfID: id.DeviceID, Gateway: c, Official: func() bool { return true }, DisableEgressProbe: true})
	defer m.Close()
	oldCtx, oldCancel := context.WithCancel(context.Background())
	defer oldCancel()
	currentCtx, currentCancel := context.WithCancel(context.Background())
	defer currentCancel()
	old := &officialAttempt{id: "old", ctx: oldCtx, cancel: oldCancel}
	m.relays["peer"] = &officialAttempt{id: "replacement", ctx: currentCtx, cancel: currentCancel}
	m.officialFailed("peer", old, &gateway.Error{Code: "quota_exhausted"})
	if currentCtx.Err() != nil || c.Usage().Reason != "" {
		t.Fatal("stale tunnel closure poisoned the replacement")
	}
}
func TestOfficialDirectThenGatewayTURN(t *testing.T) {
	binary := os.Getenv("PF_GATEWAY_TEST_BINARY")
	if binary == "" {
		t.Skip("set PF_GATEWAY_TEST_BINARY for independent gateway interoperability")
	}
	udp, e := net.ListenPacket("udp4", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	port := udp.LocalAddr().(*net.UDPAddr).Port
	udp.Close()
	backend, e := turnrelay.New(turnrelay.Options{PublicIP: "127.0.0.1", ListenIP: "127.0.0.1", Port: port, Secret: "test-only-separate-turn-secret"})
	if e != nil {
		t.Fatal(e)
	}
	defer backend.Close()
	tcp, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	addr := tcp.Addr().String()
	tcp.Close()
	dir := t.TempDir()
	cmd := exec.Command(binary, "serve", "-listen", addr, "-db", filepath.Join(dir, "gateway.db"), "-secret-file", filepath.Join(dir, "secret"), "-signal", "http://127.0.0.1:1", "-turn", fmt.Sprintf("127.0.0.1:%d", port), "-udp-bind", "127.0.0.1", "-stun", "")
	cmd.Stderr = os.Stderr
	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() { cmd.Process.Kill(); cmd.Wait() }()
	base := "http://" + addr
	until := time.Now().Add(5 * time.Second)
	for {
		r, e := http.Get(base + "/healthz")
		if e == nil {
			r.Body.Close()
			break
		}
		if time.Now().After(until) {
			t.Fatal("gateway startup failed")
		}
		time.Sleep(30 * time.Millisecond)
	}
	a, b, messages := officialPair(t, base, backend, true)
	start := time.Now()
	if a.selfID < b.selfID {
		_ = a.Connect(context.Background(), b.selfID)
	} else {
		_ = b.Connect(context.Background(), a.selfID)
	}
	time.Sleep(time.Second)
	if a.gateway.Usage().UpdatedAt != 0 || b.gateway.Usage().UpdatedAt != 0 {
		t.Fatal("relay authorization happened before direct attempt")
	}
	waitPair(t, a, b, StateRelay, 25*time.Second)
	if time.Since(start) < officialDirectWait {
		t.Fatal("relay selected before direct budget")
	}
	if e = a.Send(b.selfID, &pb.PeerMessage{Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: "through-independent-gateway"}}}); e != nil {
		t.Fatal(e)
	}
	select {
	case msg := <-messages:
		if msg.GetFetch().GetClipId() != "through-independent-gateway" {
			t.Fatal("wrong end-to-end message")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("encrypted peer channel did not deliver")
	}
}
