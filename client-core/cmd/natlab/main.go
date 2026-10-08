// natlab 是 NAT 打洞实验室的驱动程序，配合 tools/natlab/run.sh 在 Linux 网络命名空间里
// 复现真实网络拓扑（多出口、对称 NAT 等），验证两台设备最终走直连还是中转。
//
//	natlab infra -http 10.200.0.2:9000 -stun 10.200.1.1,10.200.1.2 -turn 10.200.0.3
//	natlab peer -id a -peer b -infra http://10.200.0.2:9000 -stun ... -probe ... -turn ...
//
// infra 扮演「互联网」上的服务：信令转发、若干台 STUN 服务器、一台 TURN 中转。
// peer 用的就是 CopySync 真实的 P2P 连接代码，只把信令换成了这里的 HTTP 转发。
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pion/stun/v4"
	"google.golang.org/protobuf/proto"

	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
	turnrelay "github.com/baiyuze/copysync/server/relay"
)

const turnSecret = "natlab"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "用法: natlab <infra|peer> [参数]")
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "infra":
		err = runInfra(os.Args[2:])
	case "peer":
		err = runPeer(os.Args[2:])
	default:
		err = fmt.Errorf("未知的子命令 %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "natlab:", err)
		os.Exit(1)
	}
}

// ─────────────────────────── infra ───────────────────────────

func runInfra(args []string) error {
	fs := flag.NewFlagSet("infra", flag.ExitOnError)
	httpAddr := fs.String("http", "", "信令转发的 HTTP 监听地址")
	stunIPs := fs.String("stun", "", "STUN 服务器的 IP，逗号分隔，每个 IP 各起一台")
	turnIP := fs.String("turn", "", "TURN 中转的 IP")
	_ = fs.Parse(args)

	for _, ip := range splitList(*stunIPs) {
		if err := serveSTUN(ip); err != nil {
			return err
		}
	}
	turn, err := turnrelay.New(turnrelay.Options{
		PublicIP: *turnIP, ListenIP: *turnIP, Port: 3478, Secret: turnSecret,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	if err != nil {
		return fmt.Errorf("启动 TURN: %w", err)
	}

	relay := &signalRelay{queues: map[string]chan envelope{}}
	mux := http.NewServeMux()
	mux.HandleFunc("/send", relay.send)
	mux.HandleFunc("/recv", relay.recv)
	mux.HandleFunc("/turn", func(w http.ResponseWriter, r *http.Request) {
		c := turn.Issue(r.URL.Query().Get("id"))
		_ = json.NewEncoder(w).Encode(map[string]any{
			"urls": c.URLs, "username": c.Username, "credential": c.Password,
		})
	})
	return http.ListenAndServe(*httpAddr, mux)
}

// serveSTUN 在指定 IP 上起一台最简的 STUN 服务器：只回答 Binding 请求。
// 绑定具体 IP 而不是 0.0.0.0：回包的源地址必须与客户端发往的地址一致。
func serveSTUN(ip string) error {
	conn, err := net.ListenPacket("udp4", net.JoinHostPort(ip, "3478"))
	if err != nil {
		return fmt.Errorf("STUN %s: %w", ip, err)
	}
	go func() {
		buf := make([]byte, 1500)
		for {
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				return
			}
			req := &stun.Message{Raw: append([]byte(nil), buf[:n]...)}
			if req.Decode() != nil || req.Type != stun.BindingRequest {
				continue
			}
			ua := addr.(*net.UDPAddr)
			resp, err := stun.Build(stun.NewTransactionIDSetter(req.TransactionID), stun.BindingSuccess,
				&stun.XORMappedAddress{IP: ua.IP, Port: ua.Port}, stun.Fingerprint)
			if err == nil {
				_, _ = conn.WriteTo(resp.Raw, addr)
			}
		}
	}()
	return nil
}

type envelope struct {
	from string
	body []byte
}

// signalRelay 是最简的信令转发：按收件人排队，收件人长轮询取走。
type signalRelay struct {
	mu     sync.Mutex
	queues map[string]chan envelope
}

func (r *signalRelay) queue(id string) chan envelope {
	r.mu.Lock()
	defer r.mu.Unlock()
	q, ok := r.queues[id]
	if !ok {
		q = make(chan envelope, 256)
		r.queues[id] = q
	}
	return q
}

func (r *signalRelay) send(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)
	q := req.URL.Query()
	r.queue(q.Get("to")) <- envelope{from: q.Get("from"), body: body}
}

func (r *signalRelay) recv(w http.ResponseWriter, req *http.Request) {
	select {
	case e := <-r.queue(req.URL.Query().Get("id")):
		w.Header().Set("X-From", e.from)
		_, _ = w.Write(e.body)
	case <-time.After(20 * time.Second):
		w.WriteHeader(http.StatusNoContent)
	case <-req.Context().Done():
	}
}

// ─────────────────────────── peer ───────────────────────────

type result struct {
	ID      string   `json:"id"`
	State   string   `json:"state"`
	Egress  []string `json:"egress"`
	Elapsed string   `json:"elapsed"`
}

func runPeer(args []string) error {
	fs := flag.NewFlagSet("peer", flag.ExitOnError)
	id := fs.String("id", "", "本端设备 ID")
	peerID := fs.String("peer", "", "对端设备 ID")
	infra := fs.String("infra", "", "infra 的 HTTP 地址")
	stunList := fs.String("stun", "", "信令服务器下发的 STUN（ip:port，逗号分隔）")
	probeList := fs.String("probe", "", "额外的探测服务器，对应真实环境里的公共 STUN")
	statePath := fs.String("state", "", "出口历史文件")
	legacy := fs.Bool("legacy", false, "用旧版本的候选收集方式（对照实验）")
	timeout := fs.Duration("timeout", 40*time.Second, "等待连接的最长时间")
	linger := fs.Duration("linger", 4*time.Second, "报完结果后保持连接的时间")
	settle := fs.Duration("settle", 0, "连上中转后再等多久，看能否换成直连")
	_ = fs.Parse(args)

	ctx, cancel := context.WithTimeout(context.Background(), *timeout+*settle+10*time.Second)
	defer cancel()
	start := time.Now()

	states := make(chan p2p.ConnState, 16)
	m := p2p.NewManager(p2p.ManagerOptions{
		SelfID:             *id,
		Logger:             slog.New(slog.NewTextHandler(os.Stderr, nil)),
		ProbeAllowPrivate:  true, // 实验室用内网地址扮演公网
		DisableEgressProbe: *legacy,
		StatePath:          *statePath,
		ProbeServers:       func() []string { return splitList(*probeList) },
		SendSignal: func(ctx context.Context, to string, msg *pb.SignalMessage) error {
			raw, err := proto.Marshal(msg)
			if err != nil {
				return err
			}
			u := fmt.Sprintf("%s/send?from=%s&to=%s", *infra, url.QueryEscape(*id), url.QueryEscape(to))
			resp, err := http.Post(u, "application/octet-stream", strings.NewReader(string(raw)))
			if err == nil {
				resp.Body.Close()
			}
			return err
		},
		OnState: func(_ string, s p2p.ConnState) {
			select {
			case states <- s:
			default:
			}
		},
	})
	defer m.Close()

	servers, err := iceServers(*infra, *id, splitList(*stunList))
	if err != nil {
		return err
	}
	m.SetICEServers(servers)
	go receiveSignals(ctx, m, *infra, *id)

	// 等首轮出口探测结束，免得用一份空结果去握手
	if !*legacy {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) && m.Egress().At.IsZero() {
			time.Sleep(50 * time.Millisecond)
		}
	}
	if err := m.Connect(ctx, *peerID); err != nil {
		return err
	}

	final := waitFinal(states, m, *peerID, *timeout, *settle)
	var egress []string
	for _, e := range m.Egress().Egresses {
		egress = append(egress, e.Addr.String())
	}
	err = json.NewEncoder(os.Stdout).Encode(result{
		ID: *id, State: final.String(), Egress: egress, Elapsed: time.Since(start).Round(time.Millisecond).String(),
	})
	// 报完结果再留一会儿：先报完的一端立即退出会关掉连接，另一端还在等稳定状态，
	// 读到的就成了「已断开」
	time.Sleep(*linger)
	return err
}

// waitFinal 等连接建立，再多等一会儿让两端交换 LinkInfo，取稳定后的状态。
// settle 大于零时，若稳定后是中转，再最多等这么久，看能否换成直连。
func waitFinal(states <-chan p2p.ConnState, m *p2p.Manager, peerID string, timeout, settle time.Duration) p2p.ConnState {
	deadline := time.After(timeout)
	for {
		select {
		case s := <-states:
			if s != p2p.StateDirect && s != p2p.StateRelay {
				continue
			}
			time.Sleep(1500 * time.Millisecond)
			until := time.Now().Add(settle)
			for m.State(peerID) == p2p.StateRelay && time.Now().Before(until) {
				time.Sleep(100 * time.Millisecond)
			}
			return m.State(peerID)
		case <-deadline:
			return m.State(peerID)
		}
	}
}

func iceServers(infra, id string, stunAddrs []string) ([]*pb.IceServer, error) {
	resp, err := http.Get(infra + "/turn?id=" + url.QueryEscape(id))
	if err != nil {
		return nil, fmt.Errorf("获取 TURN 凭证: %w", err)
	}
	defer resp.Body.Close()
	var cred struct {
		URLs       []string `json:"urls"`
		Username   string   `json:"username"`
		Credential string   `json:"credential"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cred); err != nil {
		return nil, err
	}
	servers := []*pb.IceServer{{Urls: cred.URLs, Username: cred.Username, Credential: cred.Credential}}
	if len(stunAddrs) > 0 {
		var urls []string
		for _, a := range stunAddrs {
			urls = append(urls, "stun:"+a)
		}
		servers = append(servers, &pb.IceServer{Urls: urls})
	}
	return servers, nil
}

func receiveSignals(ctx context.Context, m *p2p.Manager, infra, id string) {
	client := &http.Client{Timeout: 30 * time.Second}
	for ctx.Err() == nil {
		resp, err := client.Get(infra + "/recv?id=" + url.QueryEscape(id))
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			continue
		}
		var msg pb.SignalMessage
		if err := proto.Unmarshal(body, &msg); err != nil {
			continue
		}
		if err := m.HandleSignal(ctx, resp.Header.Get("X-From"), &msg); err != nil &&
			!errors.Is(err, context.Canceled) {
			fmt.Fprintln(os.Stderr, "处理信令:", err)
		}
	}
}

func splitList(s string) []string {
	var out []string
	for _, v := range strings.Split(s, ",") {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
	}
	return out
}
