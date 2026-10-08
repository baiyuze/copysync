package sync

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	gosync "sync"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/store"
	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

func TestBacklogKeepsNewest(t *testing.T) {
	var b backlog
	for i := range backlogMax + 5 {
		b.add("p", fmt.Sprint(i))
	}
	b.add("p", "10") // 再次加入的挪到最后，不重复

	got := b.take("p")
	if len(got) != backlogMax {
		t.Fatalf("留下 %d 条，应为 %d", len(got), backlogMax)
	}
	if got[0] != "5" || got[len(got)-1] != "10" {
		t.Errorf("首尾 = %s…%s，应为 5…10", got[0], got[len(got)-1])
	}
	if again := b.take("p"); len(again) != 0 {
		t.Errorf("取走后还剩 %v", again)
	}
}

type fakePeers struct {
	mu     gosync.Mutex
	up     map[string]bool
	offers map[string][]*pb.ClipOffer
}

func (f *fakePeers) send(peer string, msg *pb.PeerMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.up[peer] {
		return p2p.ErrNotConnected
	}
	if o := msg.GetOffer(); o != nil {
		f.offers[peer] = append(f.offers[peer], o)
	}
	return nil
}

func (f *fakePeers) received(peer string) []*pb.ClipOffer {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.offers[peer]
}

func newBacklogEngine(t *testing.T, peers *fakePeers) *Engine {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "copysync.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	cfg := config.Default()
	return &Engine{
		log:          slog.New(slog.NewTextHandler(io.Discard, nil)),
		store:        db,
		loadConfig:   func() config.Config { return cfg },
		send:         peers.send,
		pairedPeers:  func() []string { return []string{"office", "home"} },
		quiet:        map[string]bool{},
		lastProgress: map[string]time.Time{},
	}
}

func putText(t *testing.T, e *Engine, id, text string, at time.Time) store.Clip {
	t.Helper()
	clip := store.Clip{
		ID: id, Kind: store.KindText, Status: store.StatusReady, Outgoing: true,
		TextContent: text, TextPreview: text, TotalSize: int64(len(text)),
		CreatedAt: at, ExpiresAt: at.Add(72 * time.Hour),
	}
	if err := e.store.PutClip(clip); err != nil {
		t.Fatal(err)
	}
	return clip
}

// 复制时对端连不上的内容，连接恢复后补发，并标为 backlog；已送达的不重复发。
func TestMissedOfferIsReplayedOnReconnect(t *testing.T) {
	peers := &fakePeers{up: map[string]bool{"office": true}, offers: map[string][]*pb.ClipOffer{}}
	e := newBacklogEngine(t, peers)
	ctx := context.Background()

	now := time.Now()
	for i, text := range []string{"一", "二", "三"} {
		clip := putText(t, e, fmt.Sprintf("clip-%d", i), text, now.Add(time.Duration(i)*time.Second))
		e.offerToPeers(ctx, clip, e.loadConfig())
	}
	if n := len(peers.received("office")); n != 3 {
		t.Fatalf("在线的设备收到 %d 条，应为 3", n)
	}
	if n := len(peers.received("home")); n != 0 {
		t.Fatalf("连不上的设备收到 %d 条，应为 0", n)
	}

	peers.mu.Lock()
	peers.up["home"] = true
	peers.mu.Unlock()
	e.PeerReady(ctx, "home")

	got := peers.received("home")
	if len(got) != 3 {
		t.Fatalf("补发了 %d 条，应为 3", len(got))
	}
	for i, o := range got {
		if !o.GetBacklog() {
			t.Errorf("第 %d 条没有标为补发", i)
		}
		if want := fmt.Sprintf("clip-%d", i); o.GetClipId() != want {
			t.Errorf("第 %d 条是 %s，应按复制先后补发 %s", i, o.GetClipId(), want)
		}
	}

	e.PeerReady(ctx, "office")
	e.PeerReady(ctx, "home")
	if n := len(peers.received("office")); n != 3 {
		t.Errorf("已送达的设备又收到补发，共 %d 条", n)
	}
	if n := len(peers.received("home")); n != 3 {
		t.Errorf("补发过的又发了一次，共 %d 条", n)
	}
}

// 太旧的、已删除的不补发；补发途中又断开，剩下的留到下次。
func TestReplaySkipsStaleAndRequeuesOnFailure(t *testing.T) {
	peers := &fakePeers{up: map[string]bool{}, offers: map[string][]*pb.ClipOffer{}}
	e := newBacklogEngine(t, peers)
	ctx := context.Background()

	old := putText(t, e, "old", "旧", time.Now().Add(-2*backlogMaxAge))
	e.offerToPeers(ctx, old, e.loadConfig())
	e.backlog.add("home", "deleted")
	fresh := putText(t, e, "fresh", "新", time.Now())
	e.offerToPeers(ctx, fresh, e.loadConfig())

	// 还没连上就调用：太旧的、已删除的丢掉，发不出去的留到下次
	e.PeerReady(ctx, "home")
	left := e.backlog.take("home")
	if len(left) != 1 || left[0] != "fresh" {
		t.Fatalf("留下 %v，应只留下 fresh", left)
	}
	e.backlog.add("home", "fresh")

	peers.mu.Lock()
	peers.up["home"] = true
	peers.mu.Unlock()
	e.PeerReady(ctx, "home")
	got := peers.received("home")
	if len(got) != 1 || got[0].GetClipId() != "fresh" {
		t.Fatalf("补发了 %v，应只补发 fresh", got)
	}
}

// 补发来的文件收完后不写进剪贴板：对端此刻剪贴板里的内容更新。
func TestBacklogOfferIsMarkedQuiet(t *testing.T) {
	peers := &fakePeers{up: map[string]bool{}, offers: map[string][]*pb.ClipOffer{}}
	e := newBacklogEngine(t, peers)
	e.onRecord = func(*pb.ClipRecord) {}

	for _, tc := range []struct {
		id      string
		backlog bool
	}{{"live", false}, {"late", true}} {
		e.handleOffer(context.Background(), "office", &pb.ClipOffer{
			ClipId: tc.id, Kind: pb.ClipKind(store.KindFile), TotalSize: 10, WillPush: true,
			Items: []*pb.ClipItem{{Name: "a.txt", Size: 10}}, Backlog: tc.backlog,
			CreatedAtUnix: time.Now().Unix(),
		})
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.quiet["live"] || !e.quiet["late"] {
		t.Errorf("quiet = %v，应只有补发的 late", e.quiet)
	}
}
