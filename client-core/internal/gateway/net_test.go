package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/identity"
	"github.com/coder/websocket"
)

func TestTunnelClosureImmediatelyNotifiesOwningConnection(t *testing.T) {
	closeTunnel := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer ws.CloseNow()
		b, _ := json.Marshal(map[string]string{"local_address": "127.0.0.1:12345", "turn_address": "127.0.0.1:3478"})
		if ws.Write(r.Context(), websocket.MessageText, b) != nil {
			return
		}
		<-closeTunnel
		_ = ws.Close(websocket.StatusPolicyViolation, "quota_exhausted")
	}))
	defer server.Close()
	c, err := NewClient(identity.FromSeed(make([]byte, 32)), server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.access["fixture"] = Access{Token: "test-only-token", ExpiresAt: time.Now().Add(time.Hour).Unix()}
	closed := make(chan error, 1)
	p, err := c.openPacket(context.Background(), "fixture", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 3478}, func(err error) { closed <- err })
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	close(closeTunnel)
	select {
	case err = <-closed:
		var e *Error
		if !errors.As(err, &e) || e.Code != "quota_exhausted" || c.Usage().Reason != "quota_exhausted" {
			t.Fatalf("lost immediate quota closure: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("closed relay waited for ICE timeout")
	}
}
