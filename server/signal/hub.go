// Package signal 是信令服务器：设备注册、在线状态广播、消息转发。
//
// 服务器不存任何用户数据，也读不懂设备间的内容——它只做三件事：
// 校验「连上来的设备确实持有其声称公钥的私钥」、告诉设备它关心的对端是否在线、
// 把带签名的信令原样转发。设备间的互信完全由配对时交换的公钥承担。
package signal

import (
	"sync"
)

// Conn 是一条已完成认证的设备连接。
type Conn struct {
	DeviceID   string
	DeviceName string
	Platform   string
	PublicKey  []byte

	// send 是该连接的出站队列。写满说明对端消费不过来，
	// 此时丢弃消息而不是阻塞 hub——信令都是可重试的。
	send chan []byte

	// watching 是此设备关心在线状态的对端集合（即它已配对的设备）
	watching map[string]struct{}
}

func (c *Conn) Send(data []byte) bool {
	select {
	case c.send <- data:
		return true
	default:
		return false
	}
}

func (c *Conn) Outbound() <-chan []byte { return c.send }

type Hub struct {
	mu    sync.RWMutex
	conns map[string]*Conn
}

func NewHub() *Hub {
	return &Hub{conns: make(map[string]*Conn)}
}

// Register 接纳一条新连接。同一设备重复连接时，旧连接被顶替并关闭——
// 客户端重连（网络切换、进程重启）时不应留下僵尸连接。
func (h *Hub) Register(c *Conn) (replaced *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if old, ok := h.conns[c.DeviceID]; ok {
		delete(h.conns, c.DeviceID)
		replaced = old
	}
	h.conns[c.DeviceID] = c
	return replaced
}

// Unregister 移除连接。只有当前注册的那条才会被移除——
// 避免被顶替的旧连接在收尾时误删新连接。
func (h *Hub) Unregister(c *Conn) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if cur, ok := h.conns[c.DeviceID]; ok && cur == c {
		delete(h.conns, c.DeviceID)
		return true
	}
	return false
}

func (h *Hub) Get(deviceID string) (*Conn, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.conns[deviceID]
	return c, ok
}

func (h *Hub) Online(deviceID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.conns[deviceID]
	return ok
}

// SetWatching 记录某设备关心哪些对端的在线状态。
func (h *Hub) SetWatching(deviceID string, peers []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	c, ok := h.conns[deviceID]
	if !ok {
		return
	}
	c.watching = make(map[string]struct{}, len(peers))
	for _, p := range peers {
		c.watching[p] = struct{}{}
	}
}

// WatchersOf 返回正关注 deviceID 的所有在线连接。
//
// 只通知声明过关心该设备的连接：设备间是否配对只有客户端知道，
// 服务器不维护配对关系，靠 ClientHello 里带上来的列表来路由。
func (h *Hub) WatchersOf(deviceID string) []*Conn {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var out []*Conn
	for id, c := range h.conns {
		if id == deviceID {
			continue
		}
		if _, ok := c.watching[deviceID]; ok {
			out = append(out, c)
		}
	}
	return out
}

// PresenceFor 返回 deviceID 所关注的对端中，当前在线的那些。
func (h *Hub) PresenceFor(deviceID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	c, ok := h.conns[deviceID]
	if !ok {
		return nil
	}
	var online []string
	for peer := range c.watching {
		if _, ok := h.conns[peer]; ok {
			online = append(online, peer)
		}
	}
	return online
}

func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.conns)
}
