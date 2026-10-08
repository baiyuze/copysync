package rpc

import (
	"sync"

	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// Hub 把 daemon 内部发生的事件广播给所有已连接的 UI。
//
// 订阅者是慢速消费者时不能拖垮 daemon：每个订阅者有独立缓冲，
// 满了就丢弃该订阅者的这条事件，而不是阻塞发布方。UI 可以用
// SubscribeRequest.since_unix 重连补齐。
type Hub struct {
	mu     sync.Mutex
	subs   map[int]chan *pb.Event
	nextID int
}

func NewHub() *Hub {
	return &Hub{subs: make(map[int]chan *pb.Event)}
}

// Subscribe 返回事件通道与取消函数。调用方必须在结束时调用取消函数。
func (h *Hub) Subscribe() (<-chan *pb.Event, func()) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := h.nextID
	h.nextID++
	ch := make(chan *pb.Event, 64)
	h.subs[id] = ch

	return ch, func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		if c, ok := h.subs[id]; ok {
			delete(h.subs, id)
			close(c)
		}
	}
}

func (h *Hub) Publish(e *pb.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range h.subs {
		select {
		case ch <- e:
		default:
			// 该订阅者积压严重，丢弃这条而不是阻塞其他订阅者
		}
	}
}

func (h *Hub) SubscriberCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.subs)
}

// 便捷构造：事件类型很多，包一层省得调用点到处写 oneof 包装。

func EventClipAdded(r *pb.ClipRecord) *pb.Event {
	return &pb.Event{Payload: &pb.Event_ClipAdded{ClipAdded: r}}
}

func EventClipUpdated(r *pb.ClipRecord) *pb.Event {
	return &pb.Event{Payload: &pb.Event_ClipUpdated{ClipUpdated: r}}
}

func EventClipRemoved(id string) *pb.Event {
	return &pb.Event{Payload: &pb.Event_ClipRemoved{ClipRemoved: id}}
}

func EventProgress(p *pb.TransferProgress) *pb.Event {
	return &pb.Event{Payload: &pb.Event_Progress{Progress: p}}
}

func EventDeviceChanged(d *pb.Device) *pb.Event {
	return &pb.Event{Payload: &pb.Event_DeviceChanged{DeviceChanged: d}}
}

func EventPermissionChanged(p pb.ClipboardPermission) *pb.Event {
	return &pb.Event{Payload: &pb.Event_PermissionChanged{PermissionChanged: p}}
}

func EventStatusChanged(s *pb.Status) *pb.Event {
	return &pb.Event{Payload: &pb.Event_StatusChanged{StatusChanged: s}}
}
