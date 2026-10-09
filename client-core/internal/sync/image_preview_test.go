package sync

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/cache"
	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/store"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

func TestImagePreviewPath(t *testing.T) {
	root := t.TempDir()
	db, err := store.Open(filepath.Join(root, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	blobs, err := cache.New(filepath.Join(root, "cache"))
	if err != nil {
		t.Fatal(err)
	}
	// 不设置 watcher：预览不能依赖或写入剪贴板。
	e := &Engine{store: db, cache: blobs}
	rel, dir, err := blobs.Dir("picture")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "中文 图片.png")
	if err := os.WriteFile(path, []byte("unchanged image bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name string
		clip store.Clip
		want string
	}{
		{"received", store.Clip{Kind: store.KindImage, Status: store.StatusReady, CachePath: rel}, path},
		{"outgoing", store.Clip{Kind: store.KindImage, Status: store.StatusReady, SourcePaths: []string{path}}, path},
		{"remote-only", store.Clip{Kind: store.KindImage, Status: store.StatusRemoteOnly, SourcePaths: []string{path}}, ""},
		{"expired", store.Clip{Kind: store.KindImage, Status: store.StatusExpired, SourcePaths: []string{path}}, ""},
		{"not-image", store.Clip{Kind: store.KindFile, Status: store.StatusReady, SourcePaths: []string{path}}, ""},
		{"deleted-file", store.Clip{Kind: store.KindImage, Status: store.StatusReady, SourcePaths: []string{filepath.Join(root, "absent.png")}}, ""},
		{"directory", store.Clip{Kind: store.KindImage, Status: store.StatusReady, SourcePaths: []string{dir}}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.clip.ID = tc.name
			if err := db.PutClip(tc.clip); err != nil {
				t.Fatal(err)
			}
			got, err := e.ImagePreviewPath(context.Background(), tc.name)
			if tc.want == "" {
				if err == nil {
					t.Fatalf("unexpected path: %s", got)
				}
			} else if err != nil || got != tc.want {
				t.Fatalf("got %q, %v; want %q", got, err, tc.want)
			}
		})
	}
	if _, err := e.ImagePreviewPath(context.Background(), "missing-record"); err == nil {
		t.Fatal("missing record accepted")
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "unchanged image bytes" {
		t.Fatal("preview changed image")
	}
}

func TestPreviewDownloadDoesNotApplyClipboard(t *testing.T) {
	root := t.TempDir()
	db, err := store.Open(filepath.Join(root, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	blobs, err := cache.New(filepath.Join(root, "cache"))
	if err != nil {
		t.Fatal(err)
	}
	clip := store.Clip{ID: "preview", Kind: store.KindImage, Status: store.StatusRemoteOnly, OriginDeviceID: "peer"}
	if err := db.PutClip(clip); err != nil {
		t.Fatal(err)
	}
	sent := false
	e := New(Options{Store: db, Cache: blobs,
		LoadConfig: func() config.Config { return config.Default() },
		Send: func(peer string, message *pb.PeerMessage) error {
			sent = peer == "peer" && message.GetFetch().GetClipId() == "preview"
			return nil
		},
	})
	if err := e.FetchForPreview(context.Background(), clip.ID); err != nil {
		t.Fatal(err)
	}
	if !sent {
		t.Fatal("download was not requested")
	}
	rel, dir, err := blobs.Dir(clip.ID)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "image.png")
	if err := os.WriteFile(path, []byte("downloaded image"), 0o600); err != nil {
		t.Fatal(err)
	}
	// 默认 AutoApplyToClipboard=true 且 watcher=nil：误写剪贴板会直接导致测试失败。
	e.finishReceive(&receiver{clipID: clip.ID, streamID: "test-stream", from: "peer", start: time.Now()}, rel, []string{path})
	got, err := e.ImagePreviewPath(context.Background(), clip.ID)
	if err != nil || got != path {
		t.Fatalf("downloaded preview: %q, %v", got, err)
	}
	if e.quiet[clip.ID] {
		t.Fatal("completed preview download left quiet state")
	}
}
