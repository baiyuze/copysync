package sync

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/cache"
	"github.com/baiyuze/pasteferry/client-core/internal/clipboard"
	"github.com/baiyuze/pasteferry/client-core/internal/config"
	"github.com/baiyuze/pasteferry/client-core/internal/pack"
	"github.com/baiyuze/pasteferry/client-core/internal/store"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

type receiveClipboard struct {
	mu       sync.Mutex
	revision int64
	writes   int
	failure  error
}

func (c *receiveClipboard) Peek() (clipboard.Snapshot, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return clipboard.Snapshot{ChangeCount: c.revision}, nil
}
func (*receiveClipboard) Read() (clipboard.Content, error) { return clipboard.Content{}, nil }
func (c *receiveClipboard) Write(clipboard.Content) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.failure != nil {
		return c.failure
	}
	c.writes++
	c.revision++
	return nil
}
func (*receiveClipboard) Permission() clipboard.Permission { return clipboard.PermissionNotApplicable }
func (*receiveClipboard) OpenPermissionSettings() error    { return nil }

func receiveFixture(t *testing.T, cb *receiveClipboard, cfg config.Config, timeout time.Duration) (*Engine, chan ReceiveEvent) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	blobs, err := cache.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	loop := clipboard.NewMainLoop()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { runtime.LockOSThread(); defer runtime.UnlockOSThread(); defer close(done); loop.Run(ctx) }()
	t.Cleanup(func() { cancel(); <-done })
	events := make(chan ReceiveEvent, 32)
	e := New(Options{Store: db, Cache: blobs, Watcher: clipboard.NewWatcher(clipboard.WatcherOptions{Clipboard: cb, MainLoop: loop}),
		LoadConfig: func() config.Config { return cfg }, Send: func(string, *pb.PeerMessage) error { return nil },
		ReceiveStartTimeout: timeout, OnReceive: func(event ReceiveEvent) { events <- event },
	})
	t.Cleanup(e.Close)
	return e, events
}

func offerLarge(e *Engine, id string, backlog bool) {
	e.handleOffer(context.Background(), "peer", &pb.ClipOffer{ClipId: id, Kind: pb.ClipKind_CLIP_KIND_FILE,
		TotalSize: (50 << 20) + 1, CreatedAtUnix: time.Now().Unix(), Items: []*pb.ClipItem{{Name: "large.bin"}}, Backlog: backlog})
}

func waitReceive(t *testing.T, events <-chan ReceiveEvent, stage ReceiveStage) ReceiveEvent {
	t.Helper()
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()
	for {
		select {
		case event := <-events:
			if event.Stage == stage {
				return event
			}
		case <-timer.C:
			t.Fatalf("no stage %v", stage)
			return ReceiveEvent{}
		}
	}
}

func TestReceiveNotifiesOnlyNewLiveOffers(t *testing.T) {
	e, events := receiveFixture(t, &receiveClipboard{}, config.Default(), 0)
	offerLarge(e, "live", false)
	offerLarge(e, "live", false)
	offerLarge(e, "late", true)
	offerLarge(e, "late", false)
	if event := waitReceive(t, events, ReceivePending); event.Record.Id != "live" {
		t.Fatal(event)
	}
	select {
	case event := <-events:
		t.Fatalf("duplicate/backlog notification: %v", event)
	default:
	}
	// A sender replays an offer after it has been downloaded. Keep its cache.
	if err := e.store.SetClipCache("live", "cached", 123); err != nil {
		t.Fatal(err)
	}
	offerLarge(e, "live", false)
	clip, _ := e.store.GetClip("live")
	if clip.CachePath != "cached" || clip.Status != store.StatusReady {
		t.Fatal("replay erased ready record")
	}
}

func TestReceiveDeduplicatesAndTimesOutBeforeHeader(t *testing.T) {
	e, events := receiveFixture(t, &receiveClipboard{}, config.Default(), time.Second)
	offerLarge(e, "clip", false)
	var sends atomic.Int32
	e.send = func(string, *pb.PeerMessage) error { sends.Add(1); return nil }
	var group sync.WaitGroup
	for range 20 {
		group.Go(func() {
			if err := e.Fetch(context.Background(), "clip"); err != nil {
				t.Error(err)
			}
		})
	}
	group.Wait()
	if sends.Load() != 1 {
		t.Fatalf("sent %d fetches", sends.Load())
	}
	waitReceive(t, events, ReceiveFailed)
	clip, _ := e.store.GetClip("clip")
	if clip.Status != store.StatusFailed {
		t.Fatal(clip.Status)
	}
	if err := e.Fetch(context.Background(), "clip"); err != nil {
		t.Fatal(err)
	}
	if sends.Load() != 2 {
		t.Fatal("retry did not send")
	}
}

func TestImmediateRejectionAndOfflineFailure(t *testing.T) {
	for _, rejected := range []bool{false, true} {
		t.Run(map[bool]string{false: "offline", true: "rejected"}[rejected], func(t *testing.T) {
			e, events := receiveFixture(t, &receiveClipboard{}, config.Default(), 0)
			offerLarge(e, "clip", false)
			e.send = func(string, *pb.PeerMessage) error {
				if rejected {
					e.handleReject(&pb.FetchReject{ClipId: "clip", Reason: "source missing"})
					return nil
				}
				return errors.New("offline")
			}
			_ = e.Fetch(context.Background(), "clip")
			waitReceive(t, events, ReceiveFailed)
			clip, _ := e.store.GetClip("clip")
			if clip.Status != store.StatusFailed {
				t.Fatal("fast failure overwritten by fetching")
			}
		})
	}
}

func TestReceiveCompletionProtectsClipboardAndRetainsDownload(t *testing.T) {
	for _, scenario := range []string{"unchanged", "new-copy", "write-error", "auto-off", "preview"} {
		t.Run(scenario, func(t *testing.T) {
			cb := &receiveClipboard{revision: 10}
			cfg := config.Default()
			if scenario == "auto-off" {
				cfg.AutoApplyToClipboard = false
			}
			e, events := receiveFixture(t, cb, cfg, 0)
			offerLarge(e, "clip", false)
			if err := e.fetch(context.Background(), "clip", scenario == "preview"); err != nil {
				t.Fatal(err)
			}
			cb.mu.Lock()
			if scenario == "new-copy" {
				cb.revision++
			}
			if scenario == "write-error" {
				cb.failure = errors.New("locked")
			}
			cb.mu.Unlock()
			rel, dir, err := e.cache.Dir("clip")
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "large.bin")
			if err := os.WriteFile(path, []byte("downloaded"), 0600); err != nil {
				t.Fatal(err)
			}
			if !e.finishReceive(&receiver{clipID: "clip", total: 10}, rel, []string{path}) {
				t.Fatal("receive rejected")
			}
			clip, _ := e.store.GetClip("clip")
			if clip.Status != store.StatusReady || clip.CachePath == "" {
				t.Fatal("download was not retained")
			}
			cb.mu.Lock()
			writes := cb.writes
			cb.mu.Unlock()
			if scenario == "preview" {
				if writes != 0 {
					t.Fatal("preview wrote clipboard")
				}
				return
			}
			event := waitReceive(t, events, ReceiveCompleted)
			want := pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_PRESERVED
			if scenario == "unchanged" {
				want = pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_APPLIED
				if writes != 1 {
					t.Fatal(writes)
				}
			} else if writes != 0 {
				t.Fatal("overwrote clipboard")
			}
			if scenario == "write-error" {
				want = pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_FAILED
			}
			if event.Clipboard.Outcome != want {
				t.Fatalf("outcome %v, want %v", event.Clipboard.Outcome, want)
			}
		})
	}
}

func TestUnpackFailurePublishesTerminalStateAndRetryUsesFreshCache(t *testing.T) {
	e, events := receiveFixture(t, &receiveClipboard{}, config.Default(), 0)
	offerLarge(e, "clip", false)
	if err := e.FetchForPreview(context.Background(), "clip"); err != nil {
		t.Fatal(err)
	}
	e.handleHeader("peer", &pb.TransferHeader{ClipId: "clip", StreamId: "broken"})
	e.mu.Lock()
	r := e.incoming["broken"]
	e.mu.Unlock()
	_, _ = r.pipeW.Write([]byte("invalid zstd stream"))
	r.pipeW.Close()
	select {
	case <-r.done:
	case <-time.After(time.Second):
		t.Fatal("unpacker hung")
	}
	clip, _ := e.store.GetClip("clip")
	if clip.Status != store.StatusFailed {
		t.Fatal(clip.Status)
	}
	e.mu.Lock()
	quiet := e.quiet["clip"]
	e.mu.Unlock()
	if quiet {
		t.Fatal("preview state leaked")
	}
	if err := e.Fetch(context.Background(), "clip"); err != nil {
		t.Fatal(err)
	}
	e.handleHeader("peer", &pb.TransferHeader{ClipId: "clip", StreamId: "retry"})
	e.mu.Lock()
	r = e.incoming["retry"]
	e.mu.Unlock()
	path := filepath.Join(t.TempDir(), "large.bin")
	if err := os.WriteFile(path, []byte("whole file"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := pack.Pack(r.pipeW, []string{path}); err != nil {
		t.Fatal(err)
	}
	r.pipeW.Close()
	waitReceive(t, events, ReceiveCompleted)
	entries, err := e.cache.Entries(func() string { c, _ := e.store.GetClip("clip"); return c.CachePath }())
	if err != nil || len(entries) != 1 {
		t.Fatalf("cache: %v %v", entries, err)
	}
}
