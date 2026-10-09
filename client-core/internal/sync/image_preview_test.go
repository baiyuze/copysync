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

func TestImagePreview(t *testing.T) {
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
			got, _, _, err := e.ImagePreview(context.Background(), tc.name, 0)
			if tc.want == "" {
				if err == nil {
					t.Fatalf("unexpected path: %s", got)
				}
			} else if err != nil || got != tc.want {
				t.Fatalf("got %q, %v; want %q", got, err, tc.want)
			}
		})
	}
	if _, _, _, err := e.ImagePreview(context.Background(), "missing-record", 0); err == nil {
		t.Fatal("missing record accepted")
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "unchanged image bytes" {
		t.Fatal("preview changed image")
	}
}

// 文件记录里的图片也能预览：只算图片文件，按文件名排序，越界时回到第一张。
func TestImagePreviewOfFiles(t *testing.T) {
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
	e := &Engine{store: db, cache: blobs}
	var paths []string
	for _, name := range []string{"b.png", "a.jpg", "notes.txt"} {
		p := filepath.Join(root, name)
		if err := os.WriteFile(p, []byte(name), 0o600); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, p)
	}
	clip := store.Clip{ID: "files", Kind: store.KindFile, Status: store.StatusReady, SourcePaths: paths,
		Items: []store.Item{{Name: "b.png"}, {Name: "a.jpg"}, {Name: "notes.txt"}}}
	if err := db.PutClip(clip); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		index int
		want  string
	}{{0, "a.jpg"}, {1, "b.png"}, {5, "a.jpg"}} {
		path, name, count, err := e.ImagePreview(context.Background(), clip.ID, tc.index)
		if err != nil || name != tc.want || path != filepath.Join(root, tc.want) || count != 2 {
			t.Fatalf("第 %d 张：%q %q %d %v，应为 %s，共 2 张", tc.index, path, name, count, err, tc.want)
		}
	}

	text := store.Clip{ID: "text-file", Kind: store.KindFile, Status: store.StatusReady,
		SourcePaths: paths[2:], Items: []store.Item{{Name: "notes.txt"}}}
	if err := db.PutClip(text); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := e.ImagePreview(context.Background(), text.ID, 0); err == nil {
		t.Fatal("没有图片的文件记录也返回了预览")
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
	got, _, _, err := e.ImagePreview(context.Background(), clip.ID, 0)
	if err != nil || got != path {
		t.Fatalf("downloaded preview: %q, %v", got, err)
	}
	if e.quiet[clip.ID] {
		t.Fatal("completed preview download left quiet state")
	}
}
