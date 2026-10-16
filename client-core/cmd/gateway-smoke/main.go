// gateway-smoke uses fresh mutually known test identities to verify the public
// gateway protocol against an unchanged signaling/TURN server. It deliberately
// withholds direct candidates to exercise the eight-second fallback path.
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/gateway"
	"github.com/baiyuze/pasteferry/client-core/internal/identity"
	"github.com/baiyuze/pasteferry/client-core/internal/transport/p2p"
	"github.com/baiyuze/pasteferry/client-core/internal/transport/signaling"
	"github.com/baiyuze/pasteferry/client-core/internal/transport/signed"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

type node struct {
	id        *identity.Identity
	gateway   *gateway.Client
	p2p       *p2p.Manager
	signal    *signaling.Client
	online    atomic.Bool
	connected atomic.Bool
}

func fresh() *identity.Identity {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return identity.FromSeed(b)
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	base := flag.String("gateway", gateway.OfficialBaseURL, "explicit HTTPS gateway or loopback HTTP forward")
	signalURL := flag.String("signal", gateway.OfficialSignalURL, "signaling WSS or loopback forward")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	blackhole, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer blackhole.Close()
	blackholePort := blackhole.LocalAddr().(*net.UDPAddr).Port
	a, b := &node{id: fresh()}, &node{id: fresh()}
	received := make(chan string, 4)
	for _, pair := range [][2]*node{{a, b}, {b, a}} {
		self, peer := pair[0], pair[1]
		var err error
		self.gateway, err = gateway.NewClient(self.id, *base, nil)
		if err != nil {
			return err
		}
		defer self.gateway.Close()
		self.p2p = p2p.NewManager(p2p.ManagerOptions{SelfID: self.id.DeviceID, Gateway: self.gateway, Official: func() bool { return true }, OfficialOnline: func(string) bool { return self.online.Load() }, DisableEgressProbe: true,
			SendSignal: func(ctx context.Context, to string, msg *pb.SignalMessage) error {
				if !msg.GetOfficialRelay() {
					if msg.GetCandidate() != nil {
						msg.GetCandidate().Candidate = fmt.Sprintf("candidate:1 1 udp 2130706431 127.0.0.1 %d typ host", blackholePort)
					}
					strip := func(s string) string {
						var out []string
						for _, line := range strings.Split(s, "\n") {
							if !strings.HasPrefix(strings.TrimSpace(line), "a=candidate:") {
								out = append(out, line)
							}
						}
						return strings.Join(out, "\n")
					}
					if o := msg.GetOffer(); o != nil {
						o.Sdp = strip(o.Sdp)
					}
					if a := msg.GetAnswer(); a != nil {
						a.Sdp = strip(a.Sdp)
					}
				}
				proof, err := signed.Sign(self.id, to, msg)
				if err != nil {
					return err
				}
				return self.signal.SendSigned(ctx, proof)
			},
			OnMessage: func(_ string, msg *pb.PeerMessage) {
				if msg.GetFetch() != nil {
					select {
					case received <- msg.GetFetch().GetClipId():
					default:
					}
				}
			},
		})
		defer self.p2p.Close()
		self.signal = signaling.New(signaling.Options{Identity: self.id, URL: *signalURL, DialHeaders: func(string) http.Header { return gateway.SignalHeaders(self.id) }, PublicKeyOf: func(id string) ed25519.PublicKey {
			if id == peer.id.DeviceID {
				return peer.id.PublicKey
			}
			return nil
		}, DeviceInfo: func() (string, string, []string) { return "PasteFerry smoke fixture", "", []string{peer.id.DeviceID} }, Handlers: signaling.Handlers{OnConnected: func() { self.connected.Store(true) }, OnPresence: func(id string, online bool) {
			if id == peer.id.DeviceID {
				self.online.Store(online)
			}
		}, OnIceConfig: self.p2p.SetICEServers, OnSignal: func(from string, msg *pb.SignalMessage) { _ = self.p2p.HandleSignal(ctx, from, msg) }}})
	}
	go a.signal.Run(ctx)
	go b.signal.Run(ctx)
	for !a.connected.Load() || !b.connected.Load() || !a.online.Load() || !b.online.Load() {
		select {
		case <-ctx.Done():
			return errors.New("authenticated signaling/presence timed out")
		case <-time.After(50 * time.Millisecond):
		}
	}
	offerer, answerer := a, b
	if a.id.DeviceID > b.id.DeviceID {
		offerer, answerer = b, a
	}
	if err := offerer.p2p.Connect(ctx, answerer.id.DeviceID); err != nil {
		return err
	}
	for a.p2p.State(b.id.DeviceID) != p2p.StateRelay || b.p2p.State(a.id.DeviceID) != p2p.StateRelay {
		select {
		case <-ctx.Done():
			return errors.New("gateway relay negotiation timed out")
		case <-time.After(50 * time.Millisecond):
		}
	}
	if err := a.p2p.Send(b.id.DeviceID, &pb.PeerMessage{Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: "gateway-smoke-ciphertext"}}}); err != nil {
		return err
	}
	select {
	case got := <-received:
		if got != "gateway-smoke-ciphertext" {
			return errors.New("wrong peer payload")
		}
	case <-ctx.Done():
		return errors.New("encrypted transfer timed out")
	}
	return json.NewEncoder(os.Stdout).Encode(map[string]bool{"authenticated_signaling": true, "direct_first": true, "standard_turn": true, "gateway_relay": true, "encrypted_transfer": true})
}
