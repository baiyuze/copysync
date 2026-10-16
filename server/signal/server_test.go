package signal

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/baiyuze/pasteferry/proto/deviceid"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	"github.com/coder/websocket"
)

func TestAuthenticatedHalfCloseReleasesSocket(t *testing.T) {
	s := NewServer(Options{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))})
	server := httptest.NewServer(s)
	defer server.Close()
	var raw *net.TCPConn
	transport := &http.Transport{DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		conn, err := (&net.Dialer{}).DialContext(ctx, network, address)
		if err == nil {
			raw = conn.(*net.TCPConn)
		}
		return conn, err
	}}
	defer transport.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ws, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http"), &websocket.DialOptions{HTTPClient: &http.Client{Transport: transport}})
	if err != nil {
		t.Fatal(err)
	}
	defer ws.CloseNow()
	var hello pb.ServerEnvelope
	if err = readProto(ctx, ws, &hello); err != nil {
		t.Fatal(err)
	}
	pub, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	err = writeProto(ctx, ws, &pb.ClientEnvelope{Payload: &pb.ClientEnvelope_Hello{Hello: &pb.ClientHello{
		DeviceId: deviceid.For(pub), PublicKey: pub, Signature: ed25519.Sign(private, hello.GetHello().GetChallenge()),
	}}})
	if err != nil {
		t.Fatal(err)
	}
	for s.Hub().Count() != 1 {
		select {
		case <-ctx.Done():
			t.Fatal("authenticated client never registered")
		case <-time.After(time.Millisecond):
		}
	}
	if err = raw.CloseWrite(); err != nil {
		t.Fatal(err)
	}
	_ = raw.SetReadDeadline(time.Now().Add(time.Second))
	if _, err = io.Copy(io.Discard, raw); err != nil {
		t.Fatalf("server retained a disconnected authenticated socket: %v", err)
	}
	if s.Hub().Count() != 0 {
		t.Fatal("disconnected client is still registered")
	}
}
