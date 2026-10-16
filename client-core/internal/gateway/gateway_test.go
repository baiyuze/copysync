package gateway

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/identity"
	"github.com/baiyuze/pasteferry/server/relay"
	"github.com/pion/transport/v5/stdnet"
	"github.com/pion/turn/v5"
	_ "modernc.org/sqlite"
)

func TestMachineIDPolicy(t *testing.T) {
	var paths []string
	read := func(path string) ([]byte, error) {
		paths = append(paths, path)
		if path == "/etc/machine-id" {
			return []byte("invalid"), nil
		}
		if path == "/var/lib/dbus/machine-id" {
			return []byte("1234567890abcdef1234567890abcdef"), nil
		}
		return nil, os.ErrNotExist
	}
	raw, err := readMachineID(context.Background(), "linux", read, func(context.Context) (string, error) { t.Fatal("Linux must not call host fallback"); return "", nil })
	if err != nil || raw == "" || len(paths) != 2 {
		t.Fatal(raw, paths, err)
	}
	for _, bad := range []string{"", "00000000000000000000000000000000", "ffffffffffffffffffffffffffffffff", "unknown", "123"} {
		if _, err := normalizeMachineID(bad); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
}
func TestSignedProof(t *testing.T) {
	id := identity.FromSeed(make([]byte, 32))
	p := proof(id, "relay", "nonce", "peer", "session", "")
	sig := p.Signature
	p.PeerID = "someone-else"
	if p.Signature != sig {
		t.Fatal("test invalid")
	}
	pub := id.PublicKey
	_ = pub
	if ed25519.Verify(pub, p.Message(), id.Sign(proof(id, "relay", "nonce", "peer", "session", "").Message())) {
		t.Fatal("signature did not bind peer")
	}
}

func freeTCP(t *testing.T) string {
	t.Helper()
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	a := l.Addr().String()
	l.Close()
	return a
}

// This black-box test consumes only the published gateway protocol. The
// optional private executable is never imported, linked or shipped here.
func TestGatewayTURNBlackBox(t *testing.T) {
	binary := os.Getenv("PF_GATEWAY_TEST_BINARY")
	if binary == "" {
		t.Skip("set PF_GATEWAY_TEST_BINARY to independently built gateway")
	}
	udp, e := net.ListenPacket("udp4", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	port := udp.LocalAddr().(*net.UDPAddr).Port
	udp.Close()
	backend, e := relay.New(relay.Options{PublicIP: "127.0.0.1", ListenIP: "127.0.0.1", Port: port, Secret: "integration-only-not-a-production-secret"})
	if e != nil {
		t.Fatal(e)
	}
	defer backend.Close()
	addr := freeTCP(t)
	dir := t.TempDir()
	cmd := exec.Command(binary, "serve", "-listen", addr, "-db", filepath.Join(dir, "gateway.db"), "-secret-file", filepath.Join(dir, "secret"), "-signal", "http://127.0.0.1:1", "-turn", fmt.Sprintf("127.0.0.1:%d", port), "-udp-bind", "127.0.0.1", "-stun", "")
	cmd.Stderr = os.Stderr
	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	base := "http://" + addr
	for {
		res, err := http.Get(base + "/healthz")
		if err == nil {
			res.Body.Close()
			break
		}
		select {
		case <-ctx.Done():
			t.Fatal("gateway failed to start")
		case <-time.After(50 * time.Millisecond):
		}
	}
	seedB := sha256.Sum256([]byte("gateway peer b"))
	aID, bID := identity.FromSeed(make([]byte, 32)), identity.FromSeed(seedB[:])
	a, _ := NewClient(aID, base, nil)
	b, _ := NewClient(bID, base, nil)
	defer a.Close()
	defer b.Close()
	session := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	if _, e = a.Authorize(ctx, bID.DeviceID, session); e != nil {
		t.Fatal(e)
	}
	if _, e = b.Authorize(ctx, aID.DeviceID, session); e != nil {
		t.Fatal(e)
	}
	turnAddr := fmt.Sprintf("127.0.0.1:%d", port)
	allocate := func(c *Client, id *identity.Identity) net.PacketConn {
		t.Helper()
		network, err := c.Network(ctx, session, turnAddr)
		if err != nil {
			t.Fatal(err)
		}
		pc, err := network.ListenPacket("udp4", "0.0.0.0:0")
		if err != nil {
			t.Fatal(err)
		}
		cred := backend.Issue(id.DeviceID)
		n, _ := stdnet.NewNet()
		tc, err := turn.NewClient(&turn.ClientConfig{TURNServerAddr: turnAddr, Conn: pc, Username: cred.Username, Password: cred.Password, Realm: "pasteferry", Net: n})
		if err != nil {
			t.Fatal(err)
		}
		if err = tc.Listen(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(tc.Close)
		allocation, err := tc.Allocate()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { allocation.Close() })
		return allocation
	}
	ca, cb := allocate(a, aID), allocate(b, bID)
	// Establish permission in both directions before testing actual forwarding.
	if _, e = cb.WriteTo([]byte("ready"), ca.LocalAddr()); e != nil {
		t.Fatal(e)
	}
	payload := []byte("end-to-end ciphertext fixture")
	if _, e = ca.WriteTo(payload, cb.LocalAddr()); e != nil {
		t.Fatal(e)
	}
	_ = cb.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 2048)
	n, _, e := cb.ReadFrom(buf)
	if e != nil {
		t.Fatal(e)
	}
	if string(buf[:n]) != string(payload) {
		t.Fatalf("unexpected datagram %q", buf[:n])
	}
	var u Usage
	if e = a.request(ctx, "GET", "usage", a.access[session].Token, nil, &u); e != nil {
		t.Fatal(e)
	}
	if !u.Available || u.Upload < int64(len(payload)) || u.Download < int64(len(payload)) {
		t.Fatalf("missing server accounting: %+v", u)
	}
	// A real UDP sink must not receive an arbitrary-target packet, even though
	// the test TURN backend itself would allow that destination without a gate.
	sink, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer sink.Close()
	outsideAllocation := allocate(a, aID)
	_, _ = outsideAllocation.WriteTo([]byte("must-not-forward"), sink.LocalAddr())
	_ = sink.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	if n, _, err := sink.ReadFrom(buf); err == nil {
		t.Fatalf("arbitrary proxy target received %d bytes", n)
	}
	// Lower a live reservation below usage. Existing channels must stop; a
	// previously issued credential and a large held chunk cannot bypass it.
	db, err := sql.Open("sqlite", filepath.Join(dir, "gateway.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err = db.Exec(`UPDATE settings SET body='{"daily_limit":0,"rate_bps":2000000,"global_limit":1000000000,"paused":false}' WHERE id=1`); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1200 * time.Millisecond)
	_, _ = ca.WriteTo([]byte("over-limit"), cb.LocalAddr())
	until := time.Now().Add(5 * time.Second)
	for time.Now().Before(until) && a.Usage().Reason != "quota_exhausted" {
		time.Sleep(20 * time.Millisecond)
	}
	if a.Usage().Reason != "quota_exhausted" {
		t.Fatal("an existing tunnel was not stopped at quota exhaustion")
	}
}

func TestDisabledGatewayDoesNotEnroll(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { requests.Add(1); w.WriteHeader(500) }))
	defer server.Close()
	c, _ := NewClient(identity.FromSeed(make([]byte, 32)), server.URL, nil)
	c.SetEnabled(false)
	if _, err := c.Authorize(context.Background(), "peer", "session"); err == nil {
		t.Fatal("disabled gateway enrolled")
	}
	if requests.Load() != 0 {
		t.Fatal("disabled mode made an official request")
	}
}

func TestDisablingCancelsTunnelHandshake(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(started); <-r.Context().Done() }))
	defer server.Close()
	c, _ := NewClient(identity.FromSeed(make([]byte, 32)), server.URL, nil)
	defer c.Close()
	c.access["session"] = Access{Token: "test-handshake-token"}
	done := make(chan error, 1)
	go func() {
		_, err := c.openPacket(context.Background(), "session", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 3478})
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("tunnel handshake did not start")
	}
	c.SetEnabled(false)
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("disabled tunnel completed")
		}
	case <-time.After(time.Second):
		t.Fatal("disabling did not cancel handshake")
	}
}
