// Package gateway implements the optional, documented official transport.
// It cannot change device trust or execute remotely supplied instructions.
package gateway

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/identity"
)

const (
	OfficialSignalURL = "wss://rtccpwsapi1.pasteferry.cc/signal"
	OfficialBaseURL   = "https://rtccpwsapi1.pasteferry.cc"
	Version           = "pasteferry-gateway-v1"
)

// Proof is intentionally restricted to authentication and resource selection.
// Raw machine IDs, host names, file metadata and OS details are never fields.
type Proof struct {
	PublicKey     string `json:"public_key"`
	DeviceID      string `json:"device_id"`
	Nonce         string `json:"nonce"`
	Purpose       string `json:"purpose"`
	PeerID        string `json:"peer_id"`
	SessionID     string `json:"session_id"`
	MachineDigest string `json:"machine_digest"`
	Signature     string `json:"signature"`
}

func (p Proof) Message() []byte {
	return []byte(strings.Join([]string{Version, p.Purpose, p.Nonce, p.DeviceID, p.PublicKey, p.PeerID, p.SessionID, p.MachineDigest}, "\n"))
}

func proof(id *identity.Identity, purpose, nonce, peer, session, digest string) Proof {
	p := Proof{PublicKey: base64.RawURLEncoding.EncodeToString(id.PublicKey), DeviceID: id.DeviceID, Nonce: nonce, Purpose: purpose, PeerID: peer, SessionID: session, MachineDigest: digest}
	p.Signature = base64.RawURLEncoding.EncodeToString(id.Sign(p.Message()))
	return p
}

// SignalHeaders does not read a machine identifier or request relay access.
func SignalHeaders(id *identity.Identity) http.Header {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return nil
	}
	nonce := fmt.Sprintf("%d.%s", time.Now().Unix(), base64.RawURLEncoding.EncodeToString(b))
	p := proof(id, "signal", nonce, "", "", "")
	h := make(http.Header)
	h.Set("X-PasteFerry-Key", p.PublicKey)
	h.Set("X-PasteFerry-Device", p.DeviceID)
	h.Set("X-PasteFerry-Nonce", p.Nonce)
	h.Set("X-PasteFerry-Signature", p.Signature)
	return h
}

type Usage struct {
	Available bool   `json:"available"`
	Day       string `json:"day"`
	Group     string `json:"group"`
	Upload    int64  `json:"upload"`
	Download  int64  `json:"download"`
	Limit     int64  `json:"limit"`
	Remaining int64  `json:"remaining"`
	ResetAt   int64  `json:"reset_at"`
	UpdatedAt int64  `json:"updated_at"`
	Reason    string `json:"reason"`
}

type Access struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	Usage     Usage  `json:"usage"`
}

type Error struct {
	Code       string `json:"code"`
	Usage      *Usage `json:"usage,omitempty"`
	UsageToken string `json:"usage_token,omitempty"`
	ExpiresAt  int64  `json:"expires_at,omitempty"`
}

func (e *Error) Error() string { return "official service: " + e.Code }
