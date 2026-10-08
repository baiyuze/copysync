//go:build fakeclip

package main

// 端到端测试用的「文件剪贴板」，只在 -tags fakeclip 构建时编进来，正式版里没有。
//
// 让同一台电脑上的测试实例不碰真实的剪贴板：往 $COPYSYNC_FAKE_CLIPBOARD/copy.json
// 写一份内容就算「复制」了一次；本端要写剪贴板时（收到对端的内容），写到 pasted.json。
//
//	{"kind": "text", "text": "…"}
//	{"kind": "html", "html": "<b>…</b>", "text": "…"}
//	{"kind": "image", "image": "<PNG 的 base64>"}
//	{"kind": "file", "files": ["/绝对路径", …]}

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/baiyuze/copysync/client-core/internal/clipboard"
)

type fakeContent struct {
	Kind  string   `json:"kind"`
	Text  string   `json:"text,omitempty"`
	HTML  string   `json:"html,omitempty"`
	Image []byte   `json:"image,omitempty"`
	Files []string `json:"files,omitempty"`
}

var fakeKinds = map[string]clipboard.Kind{
	"text": clipboard.KindText, "html": clipboard.KindHTML,
	"image": clipboard.KindImage, "file": clipboard.KindFile,
}

type fakeClipboard struct {
	dir   string
	count int64 // 每写一次 pasted.json 加一，与 copy.json 的修改时间一起构成 ChangeCount
}

func newClipboard() clipboard.Clipboard {
	dir := os.Getenv("COPYSYNC_FAKE_CLIPBOARD")
	if dir == "" {
		panic("fakeclip 构建需要设置 COPYSYNC_FAKE_CLIPBOARD")
	}
	_ = os.MkdirAll(dir, 0o700)
	return &fakeClipboard{dir: dir}
}

func (f *fakeClipboard) read() (fakeContent, int64, error) {
	path := filepath.Join(f.dir, "copy.json")
	info, err := os.Stat(path)
	if err != nil {
		return fakeContent{}, f.count, nil // 还没「复制」过
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return fakeContent{}, 0, err
	}
	var c fakeContent
	if err := json.Unmarshal(raw, &c); err != nil {
		return fakeContent{}, 0, err
	}
	return c, info.ModTime().UnixNano() + f.count, nil
}

func (f *fakeClipboard) Peek() (clipboard.Snapshot, error) {
	c, n, err := f.read()
	if err != nil {
		return clipboard.Snapshot{}, err
	}
	return clipboard.Snapshot{ChangeCount: n, Kind: fakeKinds[c.Kind]}, nil
}

func (f *fakeClipboard) Read() (clipboard.Content, error) {
	c, _, err := f.read()
	if err != nil {
		return clipboard.Content{}, err
	}
	kind, ok := fakeKinds[c.Kind]
	if !ok {
		return clipboard.Content{}, errors.New("copy.json 里的 kind 不认识")
	}
	return clipboard.Content{Kind: kind, Text: c.Text, HTML: c.HTML, Image: c.Image, Files: c.Files}, nil
}

func (f *fakeClipboard) Write(c clipboard.Content) error {
	out := fakeContent{Text: c.Text, HTML: c.HTML, Image: c.Image, Files: c.Files}
	for name, k := range fakeKinds {
		if k == c.Kind {
			out.Kind = name
		}
	}
	raw, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	f.count++
	return os.WriteFile(filepath.Join(f.dir, "pasted.json"), raw, 0o600)
}

func (f *fakeClipboard) Permission() clipboard.Permission { return clipboard.PermissionNotApplicable }

func (f *fakeClipboard) OpenPermissionSettings() error { return nil }
