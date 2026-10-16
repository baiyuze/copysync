package p2p

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/baiyuze/pasteferry/client-core/internal/gateway"
	"net"
	"strings"
	"sync"
	"time"

	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	"github.com/pion/webrtc/v4"
)

const officialDirectWait = 8 * time.Second

type officialAttempt struct {
	id          string
	ctx         context.Context
	peer        *Peer
	cancel      context.CancelFunc
	pending     []*pb.SignalMessage
	failedUntil time.Time
	failReason  string
	mu          sync.Mutex
}

func (m *Manager) officialEnabled() bool {
	return m.gateway != nil && m.official != nil && m.official()
}
func (m *Manager) selectedLocked(id string) *Peer {
	p := m.peers[id]
	if p != nil && p.State() == StateDirect {
		return p
	}
	if a := m.relays[id]; a != nil && a.peer != nil && a.peer.State() == StateRelay {
		return a.peer
	}
	return p
}
func (m *Manager) scheduleFallback(id string) {
	if !m.officialEnabled() || m.selfID > id {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if p := m.peers[id]; p == nil || !p.DirectAttemptReady() {
		return
	}
	if !m.capable[id] || m.fallbackTimers[id] != nil {
		return
	}
	m.fallbackTimers[id] = time.AfterFunc(officialDirectWait, func() {
		m.mu.Lock()
		delete(m.fallbackTimers, id)
		p := m.peers[id]
		existing := m.relays[id]
		keepExisting := existing != nil && (time.Now().Before(existing.failedUntil) ||
			existing.ctx != nil && existing.ctx.Err() == nil)
		m.mu.Unlock()
		if p == nil || p.State() == StateDirect || keepExisting {
			return
		}
		if m.officialOnline != nil && !m.officialOnline(id) {
			return
		}
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			return
		}
		m.startOfficial(id, base64.RawURLEncoding.EncodeToString(b), nil, true)
	})
}
func (m *Manager) handleOfficialSignal(from string, msg *pb.SignalMessage) error {
	id := msg.GetConnectionAttemptId()
	b, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil || len(b) != 32 || !msg.GetGatewayCapable() {
		return errors.New("invalid official connection attempt")
	}
	if m.officialOnline != nil && !m.officialOnline(from) {
		return errors.New("peer is not online on official signaling")
	}
	if reason := msg.GetRelayDenialCode(); reason != "" {
		m.mu.Lock()
		a := m.relays[from]
		if reason == "quota_available" {
			if a != nil && a.id == id && a.failReason == "peer_quota_exhausted" {
				a.failedUntil = time.Time{}
				a.failReason = ""
			}
			m.mu.Unlock()
			m.scheduleFallback(from)
			return nil
		}
		m.mu.Unlock()
		if a == nil || a.id != id {
			return nil
		}
		switch reason {
		case "quota_exhausted":
			reason = "peer_quota_exhausted"
		case "global_exhausted", "global_paused", "relay_blocked", "device_banned":
		default:
			return nil
		}
		m.officialFailed(from, a, &gateway.Error{Code: reason})
		return nil
	}
	m.mu.Lock()
	a := m.relays[from]
	direct := m.peers[from]
	if direct != nil && direct.State() == StateDirect {
		m.mu.Unlock()
		return nil
	}
	if a != nil && a.id == id {
		p := a.peer
		if p == nil {
			if len(a.pending) < 256 {
				a.pending = append(a.pending, msg)
			}
			m.mu.Unlock()
			return nil
		}
		m.mu.Unlock()
		return p.HandleSignal(context.Background(), msg)
	}
	m.mu.Unlock()
	// Only the deterministic offerer may initiate a relay attempt.
	if from > m.selfID || msg.GetOffer() == nil {
		return nil
	}
	m.startOfficial(from, id, msg, false)
	return nil
}
func (m *Manager) startOfficial(deviceID, id string, initial *pb.SignalMessage, offerer bool) {
	if !m.officialEnabled() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	a := &officialAttempt{id: id, ctx: ctx, cancel: cancel}
	if initial != nil {
		a.pending = append(a.pending, initial)
	}
	m.mu.Lock()
	old := m.relays[deviceID]
	if old != nil && old.id == id {
		m.mu.Unlock()
		cancel()
		return
	}
	m.relays[deviceID] = a
	m.mu.Unlock()
	if old != nil {
		old.cancel()
		m.gateway.ReleaseSession(old.id)
		if old.peer != nil {
			old.peer.Close()
		}
	}
	go func() {
		authCtx, stop := context.WithTimeout(ctx, 15*time.Second)
		_, err := m.gateway.Authorize(authCtx, deviceID, id)
		stop()
		if err != nil {
			m.officialFailed(deviceID, a, err)
			return
		}
		m.mu.RLock()
		servers := append([]webrtc.ICEServer(nil), m.iceServers...)
		m.mu.RUnlock()
		servers, addr := officialTURN(servers)
		if addr == "" {
			m.officialFailed(deviceID, a, errors.New("TURN configuration unavailable"))
			return
		}
		network, err := m.gateway.Network(ctx, id, addr)
		if err != nil {
			m.officialFailed(deviceID, a, err)
			return
		}
		p, err := NewPeer(PeerOptions{DeviceID: deviceID, SelfID: m.selfID, Logger: m.log, ICEServers: servers, ForceRelay: true, Net: network,
			SendSignal: func(ctx context.Context, msg *pb.SignalMessage) error {
				msg.GatewayCapable = true
				msg.OfficialRelay = true
				msg.ConnectionAttemptId = id
				return m.sendSignal(ctx, deviceID, msg)
			},
			OnState: func(s ConnState) {
				m.mu.RLock()
				current := m.relays[deviceID] == a
				direct := m.peers[deviceID]
				m.mu.RUnlock()
				if !current || direct != nil && direct.State() == StateDirect {
					return
				}
				if s == StateOffline && ctx.Err() == nil {
					go m.officialFailed(deviceID, a, errors.New("relay connection closed"))
					return
				}
				if m.onState != nil {
					m.onState(deviceID, s)
				}
			},
			OnMessage: func(msg *pb.PeerMessage) {
				if m.onMessage != nil {
					m.onMessage(deviceID, msg)
				}
			}, OnStream: func(label string, stream *Stream) {
				if m.onStream != nil {
					m.onStream(deviceID, label, stream)
				}
			},
		})
		if err != nil {
			m.officialFailed(deviceID, a, err)
			return
		}
		m.mu.Lock()
		if m.relays[deviceID] != a || ctx.Err() != nil {
			m.mu.Unlock()
			p.Close()
			return
		}
		a.peer = p
		pending := a.pending
		a.pending = nil
		m.mu.Unlock()
		if offerer {
			if err = p.Connect(ctx); err != nil {
				m.officialFailed(deviceID, a, err)
				return
			}
		}
		for _, msg := range pending {
			if err = p.HandleSignal(ctx, msg); err != nil {
				m.officialFailed(deviceID, a, err)
				return
			}
		}
	}()
}
func officialTURN(servers []webrtc.ICEServer) ([]webrtc.ICEServer, string) {
	for _, s := range servers {
		for _, u := range s.URLs {
			if !strings.HasPrefix(u, "turn:") || strings.Contains(u, "transport=tcp") {
				continue
			}
			addr := strings.SplitN(strings.TrimPrefix(u, "turn:"), "?", 2)[0]
			host, _, err := net.SplitHostPort(addr)
			if err != nil || net.ParseIP(host).To4() == nil {
				continue
			}
			s.URLs = []string{u}
			return []webrtc.ICEServer{s}, addr
		}
	}
	return nil, ""
}
func directICEServers(servers []webrtc.ICEServer) []webrtc.ICEServer {
	var out []webrtc.ICEServer
	for _, s := range servers {
		var urls []string
		for _, u := range s.URLs {
			if strings.HasPrefix(u, "stun:") || strings.HasPrefix(u, "stuns:") {
				urls = append(urls, u)
			}
		}
		if len(urls) > 0 {
			s.URLs = urls
			out = append(out, s)
		}
	}
	return out
}
func (m *Manager) officialFailed(id string, a *officialAttempt, err error) {
	m.mu.Lock()
	if m.relays[id] == a {
		var ge *gateway.Error
		if errors.As(err, &ge) {
			a.failReason = ge.Code
		}
		u := m.gateway.Usage()
		a.failedUntil = time.Now().Add(time.Minute)
		if u.ResetAt > time.Now().Unix() && (u.Reason == "quota_exhausted" || u.Reason == "global_exhausted" || u.Reason == "peer_quota_exhausted") {
			a.failedUntil = time.Unix(u.ResetAt, 0)
		}
	}
	p := a.peer
	m.mu.Unlock()
	a.cancel()
	m.gateway.ReleaseSession(a.id)
	if p != nil {
		p.Close()
	}
	m.gateway.ReportError(err)
	var denied *gateway.Error
	if errors.As(err, &denied) {
		switch denied.Code {
		case "quota_exhausted", "global_exhausted", "global_paused", "relay_blocked", "device_banned":
			ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
			defer stop()
			_ = m.sendSignal(ctx, id, &pb.SignalMessage{GatewayCapable: true, OfficialRelay: true, ConnectionAttemptId: a.id, RelayDenialCode: denied.Code})
		}
	}
}

// A fresh owner snapshot can re-enable relay after a policy change or reset.
// A paired peer receives only a signed availability hint; the gateway still
// independently verifies all new grants and never controls pairing trust.
func (m *Manager) OfficialAllowanceChanged() {
	if !m.officialEnabled() {
		return
	}
	m.mu.Lock()
	ready := map[string]string{}
	for id, a := range m.relays {
		switch a.failReason {
		case "quota_exhausted", "global_exhausted", "global_paused", "relay_blocked":
			a.failedUntil = time.Time{}
			a.failReason = ""
			ready[id] = a.id
		}
	}
	m.mu.Unlock()
	for id, session := range ready {
		m.scheduleFallback(id)
		go func(id, session string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = m.sendSignal(ctx, id, &pb.SignalMessage{GatewayCapable: true, OfficialRelay: true, ConnectionAttemptId: session, RelayDenialCode: "quota_available"})
		}(id, session)
	}
}
func (m *Manager) stopOfficial(id string) {
	m.mu.Lock()
	a := m.relays[id]
	delete(m.relays, id)
	if t := m.fallbackTimers[id]; t != nil {
		t.Stop()
		delete(m.fallbackTimers, id)
	}
	m.mu.Unlock()
	if a != nil {
		a.cancel()
		m.gateway.ReleaseSession(a.id)
		if a.peer != nil {
			a.peer.Close()
		}
	}
}
func (m *Manager) stopAllOfficial() {
	m.mu.RLock()
	ids := make([]string, 0, len(m.relays)+len(m.fallbackTimers))
	for id := range m.relays {
		ids = append(ids, id)
	}
	for id := range m.fallbackTimers {
		ids = append(ids, id)
	}
	m.mu.RUnlock()
	for _, id := range ids {
		m.stopOfficial(id)
	}
}
func (m *Manager) officialRetryLoop() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-m.stop:
			return
		case <-t.C:
			if !m.officialEnabled() {
				continue
			}
			m.mu.RLock()
			ps := map[string]*Peer{}
			for id, a := range m.relays {
				if a.peer != nil && a.peer.State() == StateRelay && !a.peer.busy() && m.selfID < id {
					ps[id] = m.peers[id]
				}
			}
			m.mu.RUnlock()
			for _, p := range ps {
				if p != nil && p.State() != StateDirect {
					_ = p.Connect(context.Background())
				}
			}
		}
	}
}
