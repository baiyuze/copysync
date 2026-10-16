package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/pion/transport/v5"
	"github.com/pion/transport/v5/stdnet"
)

// Net only exists for a relay-only PeerConnection. It is not a system proxy.
type Net struct {
	transport.Net
	client   *Client
	session  string
	turn     *net.UDPAddr
	ctx      context.Context
	OnClosed func(error)
}

func (c *Client) Network(ctx context.Context, session, turnAddress string) (*Net, error) {
	remote, err := net.ResolveUDPAddr("udp4", turnAddress)
	if err != nil {
		return nil, err
	}
	standard, err := stdnet.NewNet()
	if err != nil {
		return nil, err
	}
	return &Net{Net: standard, client: c, session: session, turn: remote, ctx: ctx}, nil
}
func (n *Net) ListenPacket(network, address string) (net.PacketConn, error) {
	if network != "udp4" && network != "udp" {
		return nil, errors.New("official relay supports UDP4 only")
	}
	return n.client.openPacket(n.ctx, n.session, n.turn, n.OnClosed)
}
func (n *Net) DialUDP(network string, local, remote *net.UDPAddr) (transport.UDPConn, error) {
	return nil, errors.New("official relay does not permit direct UDP dialing")
}
func (n *Net) ListenUDP(network string, local *net.UDPAddr) (transport.UDPConn, error) {
	return nil, errors.New("official relay uses a restricted packet transport")
}
func (n *Net) Dial(network, address string) (net.Conn, error) {
	return nil, errors.New("official relay does not permit arbitrary dialing")
}
func (n *Net) DialTCP(network string, local, remote *net.TCPAddr) (transport.TCPConn, error) {
	return nil, errors.New("official relay supports UDP TURN only")
}

type packetConn struct {
	client                      *Client
	session                     string
	ws                          *websocket.Conn
	local, remote               *net.UDPAddr
	ctx                         context.Context
	cancel                      context.CancelFunc
	packets                     chan []byte
	mu                          sync.Mutex
	readDeadline, writeDeadline time.Time
	changed                     chan struct{}
	once                        sync.Once
	onClosed                    func(error)
}

func (c *Client) openPacket(ctx context.Context, session string, remote *net.UDPAddr, onClosed ...func(error)) (*packetConn, error) {
	c.mu.Lock()
	a, ok := c.access[session]
	enabled, root := c.enabled, c.root
	c.mu.Unlock()
	if !ok || !enabled {
		return nil, errors.New("relay is not authorized")
	}
	ctx, cancelRoot := context.WithCancel(ctx)
	stopRoot := context.AfterFunc(root, cancelRoot)
	owned := false
	defer func() {
		if !owned {
			stopRoot()
			cancelRoot()
		}
	}()
	u := "ws" + strings.TrimPrefix(c.base, "http") + "/gateway/v1/tunnel"
	header := make(http.Header)
	header.Set("Authorization", "Bearer "+a.Token)
	var ws *websocket.Conn
	var err error
	for attempt := 0; attempt < 30; attempt++ {
		var resp *http.Response
		ws, resp, err = websocket.Dial(ctx, u, &websocket.DialOptions{HTTPHeader: header, HTTPClient: c.http})
		if err == nil {
			break
		}
		if resp == nil || resp.StatusCode != 409 {
			return nil, errors.New("official relay tunnel unavailable")
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(300 * time.Millisecond):
		}
	}
	if err != nil {
		return nil, err
	}
	ws.SetReadLimit(65536)
	helloCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	_, data, err := ws.Read(helloCtx)
	cancel()
	if err != nil {
		ws.CloseNow()
		return nil, err
	}
	var hello struct {
		Local string `json:"local_address"`
		TURN  string `json:"turn_address"`
	}
	if json.Unmarshal(data, &hello) != nil || hello.TURN != remote.String() {
		ws.CloseNow()
		return nil, errors.New("unexpected relay endpoint")
	}
	local, err := net.ResolveUDPAddr("udp4", hello.Local)
	if err != nil {
		ws.CloseNow()
		return nil, err
	}
	connCtx, stop := context.WithCancel(ctx)
	p := &packetConn{client: c, session: session, ws: ws, local: local, remote: remote, ctx: connCtx, cancel: func() { stop(); stopRoot(); cancelRoot() }, packets: make(chan []byte, 128), changed: make(chan struct{})}
	if len(onClosed) > 0 {
		p.onClosed = onClosed[0]
	}
	c.mu.Lock()
	if !c.enabled || root.Err() != nil {
		c.mu.Unlock()
		p.cancel()
		ws.CloseNow()
		return nil, errors.New("official service disabled")
	}
	c.conns[p] = true
	c.mu.Unlock()
	owned = true
	go p.receive()
	return p, nil
}
func (p *packetConn) receive() {
	defer p.Close()
	for {
		kind, b, err := p.ws.Read(p.ctx)
		if err != nil {
			if p.ctx.Err() != nil {
				return
			}
			reason := "service_unavailable"
			status := websocket.CloseStatus(err)
			if status == websocket.StatusPolicyViolation {
				u := p.client.Usage()
				u.Reason = "service_unavailable"
				var ce websocket.CloseError
				if errors.As(err, &ce) {
					switch ce.Reason {
					case "quota_exhausted", "peer_quota_exhausted", "global_exhausted", "global_paused", "relay_blocked", "device_banned":
						u.Reason = ce.Reason
					}
				}
				p.client.update(u)
				reason = u.Reason
			}
			// ICE may otherwise keep a closed relay selected for many seconds.
			// Notify the owning connection so new offers are queued for retry.
			if p.onClosed != nil {
				go p.onClosed(&Error{Code: reason})
			}
			return
		}
		if kind != websocket.MessageBinary {
			continue
		}
		select {
		case p.packets <- b:
		case <-p.ctx.Done():
			return
		default: /* bounded UDP-style receive queue */
		}
	}
}
func (p *packetConn) ReadFrom(b []byte) (int, net.Addr, error) {
	for {
		p.mu.Lock()
		deadline, changed := p.readDeadline, p.changed
		p.mu.Unlock()
		var timer *time.Timer
		var timeout <-chan time.Time
		if !deadline.IsZero() {
			left := time.Until(deadline)
			if left <= 0 {
				return 0, nil, os.ErrDeadlineExceeded
			}
			timer = time.NewTimer(left)
			timeout = timer.C
		}
		select {
		case data := <-p.packets:
			if timer != nil {
				timer.Stop()
			}
			n := copy(b, data)
			if n < len(data) {
				return n, p.remote, io.ErrShortBuffer
			}
			return n, p.remote, nil
		case <-p.ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			return 0, nil, net.ErrClosed
		case <-changed:
			if timer != nil {
				timer.Stop()
			}
			continue
		case <-timeout:
			return 0, nil, os.ErrDeadlineExceeded
		}
	}
}
func (p *packetConn) WriteTo(b []byte, a net.Addr) (int, error) {
	if a == nil || a.String() != p.remote.String() {
		return 0, errors.New("relay destination is not permitted")
	}
	if len(b) > 65536 {
		return 0, errors.New("datagram too large")
	}
	p.mu.Lock()
	deadline := p.writeDeadline
	p.mu.Unlock()
	ctx := p.ctx
	cancel := func() {}
	if !deadline.IsZero() {
		ctx, cancel = context.WithDeadline(ctx, deadline)
	}
	defer cancel()
	if err := p.ws.Write(ctx, websocket.MessageBinary, b); err != nil {
		return 0, err
	}
	return len(b), nil
}
func (p *packetConn) Close() error {
	p.once.Do(func() {
		p.cancel()
		p.ws.CloseNow()
		p.client.mu.Lock()
		delete(p.client.conns, p)
		p.client.mu.Unlock()
	})
	return nil
}
func (p *packetConn) LocalAddr() net.Addr { return p.local }
func (p *packetConn) SetDeadline(t time.Time) error {
	_ = p.SetReadDeadline(t)
	return p.SetWriteDeadline(t)
}
func (p *packetConn) SetReadDeadline(t time.Time) error {
	p.mu.Lock()
	p.readDeadline = t
	close(p.changed)
	p.changed = make(chan struct{})
	p.mu.Unlock()
	return nil
}
func (p *packetConn) SetWriteDeadline(t time.Time) error {
	p.mu.Lock()
	p.writeDeadline = t
	p.mu.Unlock()
	return nil
}
func (p *packetConn) refresh(token string) {
	ctx, cancel := context.WithTimeout(p.ctx, 5*time.Second)
	defer cancel()
	b, _ := json.Marshal(map[string]string{"token": token})
	_ = p.ws.Write(ctx, websocket.MessageText, b)
}
