package lan

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"sync"

	"github.com/pion/logging"
	"github.com/pion/mdns/v2"
	"golang.org/x/net/ipv4"
)

// mdnsDiscovery 用 pion/mdns 在局域网里广播和浏览。
//
// Windows 上广播改由系统的 DNS-SD 服务负责（见 mdns_windows.go）：mDNS 查询是别人主动发来的包，
// 后台服务没有防火墙放行规则时收不到，也就答不了；系统服务有现成的放行规则。
// 浏览不受影响——本机先发出查询，对方的回答属于回包，防火墙会放行。
type mdnsDiscovery struct {
	mu         sync.Mutex
	conn       *mdns.Conn
	name       string
	port       int
	registered bool
	sys        *systemAdvertiser // 仅 Windows
}

// NewMDNS 返回默认的 mDNS 发现实现。
func NewMDNS() Discovery { return &mdnsDiscovery{} }

func (d *mdnsDiscovery) Start(ctx context.Context, name string, port int, text map[string]string, found func(Instance)) error {
	addr, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		return err
	}
	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("监听 mDNS 端口: %w", err)
	}
	logger := logging.NewDefaultLoggerFactory()
	logger.DefaultLogLevel = logging.LogLevelError
	conn, err := mdns.NewServer(ipv4.NewPacketConn(l), nil, mdns.WithLoggerFactory(logger))
	if err != nil {
		_ = l.Close()
		return fmt.Errorf("启动 mDNS: %w", err)
	}

	d.mu.Lock()
	d.conn, d.name, d.port = conn, name, port
	d.mu.Unlock()

	if err := d.advertise(text); err != nil {
		d.Stop()
		return err
	}

	conn.OnServiceDiscovered(func(ev mdns.ServiceEvent) {
		if ev.Instance.Service != ServiceType || ev.Instance.Instance == name {
			return
		}
		ip := ev.Addr.Unmap().WithZone("")
		if !ip.Is4() {
			return
		}
		found(Instance{
			Name: ev.Instance.Instance,
			Addr: netip.AddrPortFrom(ip, ev.Instance.Port),
			Text: txtMap(ev.Instance.Text),
		})
	})
	if err := conn.Browse(ctx, ServiceType); err != nil {
		d.Stop()
		return fmt.Errorf("浏览 mDNS: %w", err)
	}
	return nil
}

func (d *mdnsDiscovery) Update(text map[string]string) error { return d.advertise(text) }

func (d *mdnsDiscovery) advertise(text map[string]string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.conn == nil {
		return nil
	}
	if useSystemAdvertiser {
		if d.sys == nil {
			d.sys = &systemAdvertiser{}
		}
		return d.sys.register(d.name, d.port, text)
	}

	entries := txtEntries(text)
	if d.registered {
		return d.conn.UpdateTXT(d.name, ServiceType, entries)
	}
	err := d.conn.Register(mdns.ServiceInstance{
		Instance: d.name,
		Service:  ServiceType,
		Domain:   "local",
		Host:     d.name + ".local.",
		Port:     uint16(d.port),
		Text:     entries,
	})
	d.registered = err == nil
	return err
}

func (d *mdnsDiscovery) Stop() {
	d.mu.Lock()
	conn, sys := d.conn, d.sys
	d.conn, d.sys, d.registered = nil, nil, false
	d.mu.Unlock()
	if sys != nil {
		sys.unregister()
	}
	if conn != nil {
		_ = conn.Close()
	}
}

func txtEntries(text map[string]string) []mdns.TXTEntry {
	out := make([]mdns.TXTEntry, 0, len(text))
	for _, k := range sortedKeys(text) {
		out = append(out, mdns.TXTEntry{Key: k, Value: []byte(text[k])})
	}
	return out
}

func txtMap(entries []mdns.TXTEntry) map[string]string {
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		out[e.Key] = string(e.Value)
	}
	return out
}
