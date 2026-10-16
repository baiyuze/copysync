// envcheck exchanges synthetic fixtures between explicitly selected machines.
// It never reads the OS clipboard or an installed application's data directory.
package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/gateway"
	"github.com/baiyuze/pasteferry/client-core/internal/identity"
	"github.com/baiyuze/pasteferry/client-core/internal/transport/p2p"
	"github.com/baiyuze/pasteferry/client-core/internal/transport/signaling"
	"github.com/baiyuze/pasteferry/client-core/internal/transport/signed"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	"github.com/pion/webrtc/v4"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func fixture(kind string) []byte {
	if kind == "image" {
		im := image.NewNRGBA(image.Rect(0, 0, 32, 32))
		for y := 0; y < 32; y++ {
			for x := 0; x < 32; x++ {
				im.SetNRGBA(x, y, color.NRGBA{uint8(x * 7), uint8(y * 7), 101, 255})
			}
		}
		var b bytes.Buffer
		_ = png.Encode(&b, im)
		return b.Bytes()
	}
	b := make([]byte, 1_048_576)
	for i := range b {
		b[i] = byte((i*31 + i/251) % 256)
	}
	return b
}
func sum(b []byte) string { v := sha256.Sum256(b); return hex.EncodeToString(v[:]) }
func checkIdentity(run, label string) *identity.Identity {
	v := sha256.Sum256([]byte("PasteFerry isolated environment fixture v1\n" + run + "\n" + label))
	return identity.FromSeed(v[:])
}
func run() (resultErr error) {
	selfLabel := flag.String("self", "", "synthetic endpoint label")
	peerLabel := flag.String("peer", "", "synthetic peer label")
	runID := flag.String("run", "", "unique test run label; never use production identity seeds")
	mode := flag.String("mode", "direct", "direct, official-relay or custom-relay")
	signalURL := flag.String("signal", "", "explicit test signaling URL")
	base := flag.String("gateway", "", "explicit HTTPS gateway or loopback HTTP forward")
	out := flag.String("output", "", "new test output directory")
	timeout := flag.Duration("timeout", 3*time.Minute, "overall weak-network test timeout")
	expectDenial := flag.String("expect-denial", "", "expected gateway denial instead of a completed transfer")
	flag.Parse()
	if *selfLabel == "" || *peerLabel == "" || *runID == "" || *signalURL == "" || *out == "" {
		return errors.New("self, peer, run, signal and output are required")
	}
	official := *mode == "direct" || *mode == "official-relay"
	if !official && *mode != "custom-relay" {
		return errors.New("invalid test mode")
	}
	if err := os.Mkdir(*out, 0700); err != nil {
		return fmt.Errorf("output must be a new directory: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	self, peer := checkIdentity(*runID, *selfLabel), checkIdentity(*runID, *peerLabel)
	failures := make(chan error, 8)
	fail := func(err error) {
		select {
		case failures <- err:
		default:
		}
	}
	var online atomic.Bool
	var client *gateway.Client
	var err error
	if official {
		client, err = gateway.NewClient(self, *base, func(u gateway.Usage) {
			switch u.Reason {
			case "quota_exhausted", "peer_quota_exhausted", "global_exhausted", "relay_blocked", "global_paused", "device_banned":
				fail(&gateway.Error{Code: u.Reason})
			}
		})
		if err != nil {
			return err
		}
		defer client.Close()
		go client.Run(ctx)
	}
	defer func() {
		if *expectDenial == "" {
			return
		}
		if client == nil || client.Usage().Reason != *expectDenial || resultErr == nil {
			resultErr = fmt.Errorf("expected gateway denial %s; got %v", *expectDenial, resultErr)
			return
		}
		b, _ := json.MarshalIndent(map[string]any{"platform": runtime.GOOS, "mode": *mode, "denial_verified": client.Usage().Reason, "device_id": self.DeviceID}, "", "  ")
		resultErr = os.WriteFile(filepath.Join(*out, "result.json"), append(b, '\n'), 0600)
		fmt.Println(string(b))
	}()
	blackhole, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer blackhole.Close()
	blackholePort := blackhole.LocalAddr().(*net.UDPAddr).Port
	const text = "PasteFerry VM 测试 · 繁體 · 日本語 · emoji 📋\nline two"
	peerText := "fixture-text:" + *peerLabel + ":" + text
	acks := map[string]chan struct{}{"image": make(chan struct{}), "file": make(chan struct{})}
	received := map[string]chan struct{}{"image": make(chan struct{}), "file": make(chan struct{}), "text": make(chan struct{}), "done": make(chan struct{})}
	var once sync.Map
	closeOnce := func(key string, ch chan struct{}) {
		if _, loaded := once.LoadOrStore(key, true); !loaded {
			close(ch)
		}
	}
	var sc *signaling.Client
	var manager *p2p.Manager
	manager = p2p.NewManager(p2p.ManagerOptions{SelfID: self.DeviceID, Gateway: client, Official: func() bool { return official }, OfficialOnline: func(string) bool { return online.Load() }, DisableEgressProbe: true, ForceRelay: *mode == "custom-relay",
		SendSignal: func(ctx context.Context, to string, msg *pb.SignalMessage) error {
			if *mode == "official-relay" && !msg.GetOfficialRelay() {
				if c := msg.GetCandidate(); c != nil {
					c.Candidate = fmt.Sprintf("candidate:1 1 udp 2130706431 127.0.0.1 %d typ host", blackholePort)
				}
				strip := func(s string) string {
					var lines []string
					for _, line := range strings.Split(s, "\n") {
						if !strings.HasPrefix(strings.TrimSpace(line), "a=candidate:") {
							lines = append(lines, line)
						}
					}
					return strings.Join(lines, "\n")
				}
				if v := msg.GetOffer(); v != nil {
					v.Sdp = strip(v.Sdp)
				}
				if v := msg.GetAnswer(); v != nil {
					v.Sdp = strip(v.Sdp)
				}
			}
			proof, e := signed.Sign(self, to, msg)
			if e != nil {
				return e
			}
			return sc.SendSigned(ctx, proof)
		},
		OnMessage: func(_ string, msg *pb.PeerMessage) {
			v := msg.GetFetch().GetClipId()
			if v == peerText {
				closeOnce("recv:text", received["text"])
			}
			if v == "fixture-done" {
				closeOnce("recv:done", received["done"])
			}
			for _, kind := range []string{"image", "file"} {
				if v == "ack:"+kind+":"+sum(fixture(kind)) {
					closeOnce("ack:"+kind, acks[kind])
				}
			}
		},
		OnStream: func(_ string, label string, stream *p2p.Stream) {
			go func() {
				if label != "envcheck/image" && label != "envcheck/file" {
					fail(errors.New("unexpected fixture stream"))
					return
				}
				kind := strings.TrimPrefix(label, "envcheck/")
				expected := fixture(kind)
				var got bytes.Buffer
				for {
					select {
					case b, ok := <-stream.Chunks():
						if !ok {
							fail(fmt.Errorf("%s ended before completion", kind))
							return
						}
						if got.Len()+len(b) > len(expected) {
							fail(errors.New("oversized fixture"))
							return
						}
						got.Write(b)
						if got.Len() == len(expected) {
							if !bytes.Equal(got.Bytes(), expected) {
								fail(fmt.Errorf("%s checksum mismatch", kind))
								return
							}
							name := "文件 📋.bin"
							if kind == "image" {
								name = "图片 📋.png"
							}
							if e := os.WriteFile(filepath.Join(*out, name), got.Bytes(), 0600); e != nil {
								fail(e)
								return
							}
							if e := manager.Send(peer.DeviceID, &pb.PeerMessage{Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: "ack:" + kind + ":" + sum(got.Bytes())}}}); e != nil {
								fail(e)
								return
							}
							closeOnce("recv:"+kind, received[kind])
							return
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		},
	})
	defer manager.Close()
	sc = signaling.New(signaling.Options{Identity: self, URL: *signalURL, PublicKeyOf: func(id string) ed25519.PublicKey {
		if id == peer.DeviceID {
			return peer.PublicKey
		}
		return nil
	}, DeviceInfo: func() (string, string, []string) { return "PasteFerry isolated fixture", "", []string{peer.DeviceID} }, DialHeaders: func(string) http.Header {
		if official {
			return gateway.SignalHeaders(self)
		}
		return nil
	}, Handlers: signaling.Handlers{OnPresence: func(id string, v bool) {
		if id == peer.DeviceID {
			online.Store(v)
		}
	}, OnIceConfig: manager.SetICEServers, OnSignal: func(from string, msg *pb.SignalMessage) {
		if e := manager.HandleSignal(ctx, from, msg); e != nil {
			fail(e)
		}
	}}})
	go sc.Run(ctx)
	wait := func(ch <-chan struct{}) error {
		select {
		case <-ch:
			return nil
		case e := <-failures:
			return e
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	for !online.Load() {
		select {
		case <-time.After(30 * time.Millisecond):
		case <-ctx.Done():
			return errors.New("peer presence timed out")
		}
	}
	started := time.Now()
	fmt.Fprintln(os.Stderr, "envcheck: peer online; starting", *mode)
	if err = manager.Connect(ctx, peer.DeviceID); err != nil {
		return err
	}
	want := p2p.StateDirect
	if *mode != "direct" {
		want = p2p.StateRelay
	}
	for manager.State(peer.DeviceID) != want {
		select {
		case <-time.After(30 * time.Millisecond):
		case e := <-failures:
			return e
		case <-ctx.Done():
			return fmt.Errorf("connection timed out at %s", manager.State(peer.DeviceID))
		}
	}
	connectedAfter := time.Since(started)
	fmt.Fprintln(os.Stderr, "envcheck: connected", want.String(), "after", connectedAfter.Round(time.Millisecond))
	if *mode == "direct" && (client.Usage().UpdatedAt != 0 || client.Usage().Reason != "") {
		return errors.New("direct connection requested gateway authorization")
	}
	if *mode == "official-relay" && connectedAfter < 8*time.Second {
		return errors.New("relay started before direct attempt window")
	}
	if err = manager.Send(peer.DeviceID, &pb.PeerMessage{Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: "fixture-text:" + *selfLabel + ":" + text}}}); err != nil {
		return err
	}
	transferStarted := time.Now()
	for _, kind := range []string{"image", "file"} {
		dc, e := manager.OpenStream(peer.DeviceID, "envcheck/"+kind)
		if e != nil {
			return e
		}
		opened := make(chan struct{})
		var openOnce sync.Once
		dc.OnOpen(func() { openOnce.Do(func() { close(opened) }) })
		if dc.ReadyState() == webrtc.DataChannelStateOpen {
			openOnce.Do(func() { close(opened) })
		}
		if e = wait(opened); e != nil {
			return e
		}
		payload := fixture(kind)
		for offset := 0; offset < len(payload); {
			for dc.BufferedAmount() > 128*1024 {
				select {
				case <-time.After(5 * time.Millisecond):
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			end := min(offset+16*1024, len(payload))
			if e = dc.Send(payload[offset:end]); e != nil {
				return e
			}
			offset = end
		}
		if e = wait(acks[kind]); e != nil {
			return e
		}
		fmt.Fprintln(os.Stderr, "envcheck: peer verified", kind)
		dc.Close()
	}
	for _, kind := range []string{"text", "image", "file"} {
		if err = wait(received[kind]); err != nil {
			return err
		}
	}
	if err = manager.Send(peer.DeviceID, &pb.PeerMessage{Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: "fixture-done"}}}); err != nil {
		return err
	}
	if err = wait(received["done"]); err != nil {
		return err
	}
	result := map[string]any{"platform": runtime.GOOS, "arch": runtime.GOARCH, "mode": *mode, "self": *selfLabel, "device_id": self.DeviceID, "peer_id": peer.DeviceID, "connection": want.String(), "connected_ms": connectedAfter.Milliseconds(), "transfer_ms": time.Since(transferStarted).Milliseconds(), "text_verified": true, "image_sha256": sum(fixture("image")), "file_sha256": sum(fixture("file")), "file_bytes": len(fixture("file"))}
	if official {
		result["gateway_authorized"] = client.Usage().UpdatedAt > 0
	}
	encoded, _ := json.MarshalIndent(result, "", "  ")
	if err = os.WriteFile(filepath.Join(*out, "result.json"), append(encoded, '\n'), 0600); err != nil {
		return err
	}
	fmt.Println(string(encoded))
	time.Sleep(2 * time.Second)
	return nil
}
