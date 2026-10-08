package sync

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/store"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// 补发：复制时对端连不上（断线、睡眠、连接已失效但还没被发现），
// 这条内容就记下来，连接恢复后再发过去。只记在内存里，后台服务重启后清空。
const (
	backlogMax    = 20             // 每台设备最多补发的条数，只留最新的
	backlogMaxAge = 24 * time.Hour // 再旧的不补
)

type backlog struct {
	mu      sync.Mutex
	pending map[string][]string // 设备 ID → 待补发的记录 ID，按复制先后排列
}

func (b *backlog) add(peer, clipID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.pending == nil {
		b.pending = make(map[string][]string)
	}
	list := slices.DeleteFunc(b.pending[peer], func(id string) bool { return id == clipID })
	list = append(list, clipID)
	if len(list) > backlogMax {
		list = slices.Clone(list[len(list)-backlogMax:])
	}
	b.pending[peer] = list
}

// take 取走某台设备的全部待补发记录。
func (b *backlog) take(peer string) []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	list := b.pending[peer]
	delete(b.pending, peer)
	return list
}

// offerToPeers 把一条本机复制发给每台已配对设备。没送达的记进待补发，
// 连接恢复后再发（见 PeerReady）。
func (e *Engine) offerToPeers(ctx context.Context, clip store.Clip, cfg config.Config) {
	peers := e.pairedPeers()
	if len(peers) == 0 {
		return
	}
	offer := e.buildOffer(clip, cfg)
	msg := &pb.PeerMessage{Payload: &pb.PeerMessage_Offer{Offer: offer}}
	var reached []string
	for _, peer := range peers {
		if err := e.send(peer, msg); err != nil {
			e.backlog.add(peer, clip.ID)
			continue
		}
		reached = append(reached, peer)
	}
	e.log.Info("已广播剪贴板内容", "clip", clip.ID, "kind", clip.Kind, "size", clip.TotalSize,
		"对端数", len(reached), "待补发", len(peers)-len(reached))

	// 小内容随即推送，对端无需等待即可粘贴
	if pushes(offer, clip) {
		for _, peer := range reached {
			go e.pushOrRequeue(context.WithoutCancel(ctx), peer, clip)
		}
	}
}

// PeerReady 在与某台设备的连接建立或恢复时调用，补发连不上期间没送达的内容。
//
// 补发的 offer 带 backlog 标记，对端只记进复制记录、不写进剪贴板：
// 它此刻剪贴板里的内容比这些旧内容更新。
func (e *Engine) PeerReady(ctx context.Context, peer string) {
	ids := e.backlog.take(peer)
	if len(ids) == 0 {
		return
	}
	cfg := e.loadConfig()
	var sent int
	for i, id := range ids {
		clip, err := e.store.GetClip(id)
		if err != nil || time.Since(clip.CreatedAt) > backlogMaxAge {
			continue // 已删除或太旧
		}
		offer := e.buildOffer(clip, cfg)
		offer.Backlog = true
		if err := e.send(peer, &pb.PeerMessage{Payload: &pb.PeerMessage_Offer{Offer: offer}}); err != nil {
			// 刚连上又断了：剩下的放回去，等下次
			for _, rest := range ids[i:] {
				e.backlog.add(peer, rest)
			}
			break
		}
		sent++
		if pushes(offer, clip) {
			// 逐条推、不并发：补发可能有一串，同时开流会互相抢带宽
			e.pushOrRequeue(ctx, peer, clip)
		}
	}
	if sent > 0 {
		e.log.Info("已补发连接中断期间的内容", "to", peer, "条数", sent)
	}
}

// pushOrRequeue 推送内容；中途断了就记进待补发，重连后连同 offer 一起重来，
// 不让对端的记录一直停在「接收中」。
func (e *Engine) pushOrRequeue(ctx context.Context, peer string, clip store.Clip) {
	if err := e.pushTo(ctx, peer, clip); err != nil {
		e.log.Warn("推送内容失败，连接恢复后补发", "clip", clip.ID, "to", peer, "err", err)
		e.backlog.add(peer, clip.ID)
	}
}

func pushes(offer *pb.ClipOffer, clip store.Clip) bool {
	return offer.GetWillPush() && (clip.Kind == store.KindFile || clip.Kind == store.KindImage)
}
