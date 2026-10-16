package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/identity"
)

type Client struct {
	id          *identity.Identity
	base        string
	http        *http.Client
	mu          sync.Mutex
	access      map[string]Access
	peers       map[string]string
	conns       map[*packetConn]bool
	last        Usage
	usageAccess Access
	onUsage     func(Usage)
	enabled     bool
	root        context.Context
	cancel      context.CancelFunc
}

func (c *Client) ReportError(err error) {
	u := c.Usage()
	var ge *Error
	if errors.As(err, &ge) {
		if ge.Usage != nil {
			u = *ge.Usage
		}
		u.Reason = ge.Code
	} else if u.Reason == "" {
		u.Available = false
		u.Reason = "service_unavailable"
	}
	c.update(u)
}

func NewClient(id *identity.Identity, base string, onUsage func(Usage)) (*Client, error) {
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || u.RawQuery != "" || u.User != nil || u.Fragment != "" {
		return nil, errors.New("invalid gateway URL")
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost")) {
		return nil, errors.New("gateway requires HTTPS")
	}
	root, cancel := context.WithCancel(context.Background())
	return &Client{id: id, base: strings.TrimRight(base, "/"), http: &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("gateway redirects are not allowed") }}, access: map[string]Access{}, peers: map[string]string{}, conns: map[*packetConn]bool{}, onUsage: onUsage, enabled: true, root: root, cancel: cancel}, nil
}

func (c *Client) SetEnabled(enabled bool) {
	if !enabled {
		c.Close()
		return
	}
	c.mu.Lock()
	if !c.enabled {
		c.root, c.cancel = context.WithCancel(context.Background())
		c.enabled = true
	}
	c.mu.Unlock()
}
func (c *Client) request(ctx context.Context, method, path, token string, body, out any) error {
	c.mu.Lock()
	enabled, root := c.enabled, c.root
	c.mu.Unlock()
	if !enabled {
		return &Error{Code: "disabled"}
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stop := context.AfterFunc(root, cancel)
	defer stop()
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+"/gateway/v1/"+path, r)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return errors.New("official service unavailable")
	}
	defer resp.Body.Close()
	d := json.NewDecoder(io.LimitReader(resp.Body, 32768))
	if resp.StatusCode != 200 {
		e := &Error{Code: "service_unavailable"}
		_ = d.Decode(e)
		if e.Usage != nil {
			c.update(*e.Usage)
		}
		return e
	}
	return d.Decode(out)
}
func (c *Client) update(u Usage) {
	c.mu.Lock()
	if !c.enabled {
		c.mu.Unlock()
		return
	}
	c.last = u
	c.mu.Unlock()
	if c.onUsage != nil {
		c.onUsage(u)
	}
}
func (c *Client) Usage() Usage { c.mu.Lock(); defer c.mu.Unlock(); return c.last }
func (c *Client) Authorize(ctx context.Context, peer, session string) (Access, error) {
	var nonce struct {
		Nonce string `json:"nonce"`
	}
	if err := c.request(ctx, "GET", "challenge", "", nil, &nonce); err != nil {
		return Access{}, err
	}
	// The only machine-ID read in the application happens after direct failure.
	p := proof(c.id, "relay", nonce.Nonce, peer, session, MachineDigest())
	var a Access
	if err := c.request(ctx, "POST", "authorize", "", p, &a); err != nil {
		var ge *Error
		if errors.As(err, &ge) && ge.UsageToken != "" && ge.ExpiresAt > time.Now().Unix() {
			c.mu.Lock()
			if c.enabled {
				c.usageAccess = Access{Token: ge.UsageToken, ExpiresAt: ge.ExpiresAt}
			}
			c.mu.Unlock()
		}
		return Access{}, err
	}
	if a.Token == "" || a.ExpiresAt <= time.Now().Unix() {
		return Access{}, errors.New("invalid gateway response")
	}
	c.mu.Lock()
	if !c.enabled {
		c.mu.Unlock()
		return Access{}, &Error{Code: "disabled"}
	}
	c.access[session] = a
	c.usageAccess = a
	c.peers[session] = peer
	conns := make([]*packetConn, 0)
	for conn := range c.conns {
		if conn.session == session {
			conns = append(conns, conn)
		}
	}
	c.mu.Unlock()
	for _, conn := range conns {
		conn.refresh(a.Token)
	}
	c.update(a.Usage)
	return a, nil
}
func (c *Client) Run(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	defer c.Close()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			owner := c.usageAccess
			sessions := make(map[string]Access, len(c.access))
			peers := make(map[string]string, len(c.peers))
			for k, v := range c.access {
				sessions[k] = v
				peers[k] = c.peers[k]
			}
			c.mu.Unlock()
			if len(sessions) == 0 && owner.Token != "" {
				var err error
				if time.Now().Unix() >= owner.ExpiresAt-120 {
					err = c.refreshUsage(ctx)
				} else {
					var u Usage
					err = c.request(ctx, "GET", "usage", owner.Token, nil, &u)
					if err == nil {
						c.update(u)
					}
				}
				if err != nil {
					c.ReportError(err)
				}
			}
			for session, a := range sessions {
				var err error
				if time.Now().Unix() >= a.ExpiresAt-120 {
					_, err = c.Authorize(ctx, peers[session], session)
				} else {
					var u Usage
					err = c.request(ctx, "GET", "usage", a.Token, nil, &u)
					if err == nil {
						c.update(u)
					}
				}
				if err != nil {
					u := c.Usage()
					u.Available = false
					u.Reason = "service_unavailable"
					var ge *Error
					if errors.As(err, &ge) {
						u.Reason = ge.Code
					}
					c.update(u)
				}
			}
		}
	}
}
func (c *Client) Close() {
	c.mu.Lock()
	c.enabled = false
	c.cancel()
	cs := make([]*packetConn, 0, len(c.conns))
	for x := range c.conns {
		cs = append(cs, x)
	}
	c.access = map[string]Access{}
	c.peers = map[string]string{}
	c.usageAccess = Access{}
	c.mu.Unlock()
	for _, x := range cs {
		x.Close()
	}
}

func (c *Client) ReleaseSession(session string) {
	c.mu.Lock()
	delete(c.access, session)
	delete(c.peers, session)
	c.mu.Unlock()
}
func (c *Client) refreshUsage(ctx context.Context) error {
	var nonce struct {
		Nonce string `json:"nonce"`
	}
	if err := c.request(ctx, "GET", "challenge", "", nil, &nonce); err != nil {
		return err
	}
	var a Access
	if err := c.request(ctx, "POST", "authorize", "", proof(c.id, "usage", nonce.Nonce, "", "", ""), &a); err != nil {
		return err
	}
	c.mu.Lock()
	if c.enabled {
		c.usageAccess = a
	}
	c.mu.Unlock()
	c.update(a.Usage)
	return nil
}
