package sync

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/clipboard"
	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/store"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

func TestLegacyRemotePlaceholderOffer(t *testing.T) {
	e := newBacklogEngine(t, &fakePeers{})
	emitted := 0
	e.onRecord = func(*pb.ClipRecord) { emitted++ }
	for i := range 100 {
		id := fmt.Sprintf("legacy-%d", i)
		e.handleOffer(context.Background(), "old-peer", &pb.ClipOffer{
			ClipId: id, Kind: pb.ClipKind(store.KindFile), WillPush: true,
			Items: []*pb.ClipItem{{Name: fmt.Sprintf(".uuremote_aeawv%d", i)}},
		})
		e.handleHeader("old-peer", &pb.TransferHeader{ClipId: id})
	}
	clips, err := e.store.ListClips(500, time.Time{}, store.KindUnspecified)
	if err != nil || len(clips) != 0 || emitted != 0 || len(e.incoming) != 0 {
		t.Fatalf("legacy flood: records=%d, events=%d, receivers=%d, err=%v", len(clips), emitted, len(e.incoming), err)
	}
}

func TestRemotePlaceholderFlood(t *testing.T) {
	e := newBacklogEngine(t, &fakePeers{})
	e.deviceName = func() string { return "test" }
	cfg := config.Default()
	dir := t.TempDir()
	for i := range 100 {
		p := filepath.Join(dir, fmt.Sprintf(".uuremote_aeawv%d", i))
		if err := os.WriteFile(p, nil, 0600); err != nil {
			t.Fatal(err)
		}
		_, err := e.recordOutgoing(clipboard.Content{Kind: clipboard.KindFile, Files: []string{p}}, cfg)
		if !errors.Is(err, errRemotePlaceholder) {
			t.Fatalf("placeholder %d recorded: %v", i, err)
		}
	}
	clips, err := e.store.ListClips(500, time.Time{}, store.KindUnspecified)
	if err != nil || len(clips) != 0 {
		t.Fatalf("flood left %d records: %v", len(clips), err)
	}

	// No broad hidden-file, zero-byte, or >50 MB filter: actual files survive.
	for _, size := range []int64{0, 5 * 1024, 50 << 20, (50 << 20) + 1, 1100 << 20} {
		p := filepath.Join(dir, fmt.Sprintf(".real-%d", size))
		f, err := os.Create(p)
		if err != nil {
			t.Fatal(err)
		}
		err = f.Truncate(size) // sparse: validate metadata without transferring 1 GB
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		c, err := e.recordOutgoing(clipboard.Content{Kind: clipboard.KindFile,
			Files: []string{filepath.Join(dir, ".uuremote_aeawv999"), p}}, cfg)
		if err != nil {
			t.Fatal(err)
		}
		if len(c.Items) != 1 || len(c.SourcePaths) != 1 || c.SourcePaths[0] != p || c.TotalSize != size {
			t.Fatalf("mixed selection corrupted: %+v", c)
		}
		offer := e.buildOffer(c, cfg)
		if len(offer.Items) != 1 || offer.TotalSize != size || offer.WillPush != (size <= cfg.AutoSyncThresholdBytes) {
			t.Fatalf("wrong offer: %v", offer)
		}
	}
}
