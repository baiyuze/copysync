package main

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/identity"
	"github.com/baiyuze/copysync/client-core/internal/transport/lan"
	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	"github.com/baiyuze/copysync/client-core/internal/transport/signed"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// runLANPeer 不经任何服务器，只靠局域网通道（mDNS 发现 + UDP 信令）连接对端。
// 两端的身份由约定的种子派生，相当于已经配对过。
//
//	natlab lanpeer -seed a -peer-seed b [-mdns] [-addr 192.168.20.2:8788]
func runLANPeer(args []string) error {
	fs := flag.NewFlagSet("lanpeer", flag.ExitOnError)
	seed := fs.String("seed", "", "本端身份的种子")
	peerSeed := fs.String("peer-seed", "", "对端身份的种子")
	useMDNS := fs.Bool("mdns", false, "用 mDNS 在局域网里发现对端")
	addr := fs.String("addr", "", "记得的对端地址（相当于配对时输入过、或地址簿里存着），逗号分隔")
	timeout := fs.Duration("timeout", 20*time.Second, "等待连接的最长时间")
	linger := fs.Duration("linger", 4*time.Second, "报完结果后保持连接的时间")
	_ = fs.Parse(args)

	self := identity.FromSeed(seedBytes(*seed))
	peer := identity.FromSeed(seedBytes(*peerSeed))
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ctx, cancel := context.WithTimeout(context.Background(), *timeout+10*time.Second)
	defer cancel()
	start := time.Now()

	var (
		tr     *lan.Transport
		m      *p2p.Manager
		states = make(chan p2p.ConnState, 16)
	)
	m = p2p.NewManager(p2p.ManagerOptions{
		SelfID: self.DeviceID,
		Logger: log,
		SendSignal: func(_ context.Context, to string, msg *pb.SignalMessage) error {
			s, err := signed.Sign(self, to, msg)
			if err != nil {
				return err
			}
			return tr.Send(to, s)
		},
		OnState: func(_ string, s p2p.ConnState) {
			select {
			case states <- s:
			default:
			}
		},
	})
	defer m.Close()

	var discovery lan.Discovery
	if *useMDNS {
		discovery = lan.NewMDNS()
	} else {
		discovery = lan.NewMemoryNetwork().Discovery() // 谁也发现不了：只靠记得的地址
	}
	book := &memBook{addrs: map[string][]string{}}
	reachable := make(chan struct{}, 1)
	tr = lan.New(lan.Options{
		Identity:  self,
		Logger:    log,
		Discovery: discovery,
		Verifier: signed.NewVerifier(self.DeviceID, func(id string) ed25519.PublicKey {
			if id == peer.DeviceID {
				return peer.PublicKey
			}
			return nil
		}),
		Book:       book,
		Peers:      func() []string { return []string{peer.DeviceID} },
		Connected:  func(id string) bool { s := m.State(id); return s == p2p.StateDirect || s == p2p.StateRelay },
		DeviceInfo: func() (string, string) { return *seed, "linux" },
		OnSignal: func(from string, msg *pb.SignalMessage) {
			if err := m.HandleSignal(ctx, from, msg); err != nil {
				log.Warn("处理信令失败", "err", err)
			}
		},
		OnReachable: func(id string, ok bool) {
			if ok {
				select {
				case reachable <- struct{}{}:
				default:
				}
			}
		},
	})
	if err := tr.Start(); err != nil {
		return err
	}
	defer tr.Stop()
	if *addr != "" {
		tr.Remember(peer.DeviceID, splitList(*addr))
		tr.PeersChanged()
	}

	// 在局域网里找到对端后发起连接；谁发 offer 由 p2p 按 ID 决定，两端都调用即可
	go func() {
		select {
		case <-reachable:
			_ = m.Connect(ctx, peer.DeviceID)
		case <-ctx.Done():
		}
	}()

	final := waitFinal(states, m, peer.DeviceID, *timeout, 0)
	err := json.NewEncoder(os.Stdout).Encode(result{
		ID: *seed, State: final.String(), Elapsed: time.Since(start).Round(time.Millisecond).String(),
	})
	time.Sleep(*linger)
	return err
}

func seedBytes(s string) []byte {
	sum := sha256.Sum256([]byte("natlab|" + s))
	return sum[:]
}

// memBook 是只在内存里的地址簿。
type memBook struct {
	mu    sync.Mutex
	addrs map[string][]string
}

func (b *memBook) DeviceAddrs(id string) ([]string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]string(nil), b.addrs[id]...), nil
}

func (b *memBook) PutDeviceAddr(id, addr string, _ time.Time) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, a := range b.addrs[id] {
		if a == addr {
			return nil
		}
	}
	b.addrs[id] = append(b.addrs[id], addr)
	return nil
}
