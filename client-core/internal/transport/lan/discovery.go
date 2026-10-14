package lan

import (
	"context"
	"net/netip"
	"sync"
)

// ServiceType 是 PasteFerry 在 mDNS（DNS-SD）里用的服务类型。
const ServiceType = "_pasteferry._udp"

// Instance 是在局域网里发现的一个 PasteFerry 实例。
type Instance struct {
	Name string
	Addr netip.AddrPort
	Text map[string]string
}

// Discovery 负责在局域网里广播本机、发现其他实例。
//
// 默认实现是 mDNS（见 mdns.go）；测试和 natlab 用 MemoryNetwork。
type Discovery interface {
	// Start 开始广播本机并浏览其他实例。found 对每个发现或更新的实例调用，
	// 同一个实例可能反复出现，由调用方去重；不会对本机自己调用。
	Start(ctx context.Context, name string, port int, text map[string]string, found func(Instance)) error
	// Update 更换本机广播的 TXT 记录。
	Update(text map[string]string) error
	Stop()
}

// MemoryNetwork 是进程内的「局域网」，供测试使用：
// 在同一个 MemoryNetwork 上 Start 的实例互相可见，地址都是 127.0.0.1。
type MemoryNetwork struct {
	mu      sync.Mutex
	members map[*memoryDiscovery]bool
}

func NewMemoryNetwork() *MemoryNetwork {
	return &MemoryNetwork{members: make(map[*memoryDiscovery]bool)}
}

// Discovery 返回一个挂在这个网络上的 Discovery。
func (n *MemoryNetwork) Discovery() Discovery { return &memoryDiscovery{net: n} }

type memoryDiscovery struct {
	net   *MemoryNetwork
	self  Instance
	found func(Instance)
}

func (d *memoryDiscovery) Start(_ context.Context, name string, port int, text map[string]string, found func(Instance)) error {
	d.net.mu.Lock()
	d.self = Instance{
		Name: name,
		Addr: netip.AddrPortFrom(netip.MustParseAddr("127.0.0.1"), uint16(port)),
		Text: copyText(text),
	}
	d.found = found
	d.net.members[d] = true
	others := d.othersLocked()
	d.net.mu.Unlock()

	for _, o := range others {
		found(o.self)
		o.found(d.self)
	}
	return nil
}

func (d *memoryDiscovery) Update(text map[string]string) error {
	d.net.mu.Lock()
	d.self.Text = copyText(text)
	self := d.self
	others := d.othersLocked()
	d.net.mu.Unlock()

	for _, o := range others {
		o.found(self)
	}
	return nil
}

func (d *memoryDiscovery) Stop() {
	d.net.mu.Lock()
	delete(d.net.members, d)
	d.net.mu.Unlock()
}

func (d *memoryDiscovery) othersLocked() []*memoryDiscovery {
	var out []*memoryDiscovery
	for m := range d.net.members {
		if m != d {
			out = append(out, m)
		}
	}
	return out
}

func copyText(text map[string]string) map[string]string {
	out := make(map[string]string, len(text))
	for k, v := range text {
		out[k] = v
	}
	return out
}
