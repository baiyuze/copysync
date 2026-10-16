// Package notifications owns desktop notifications in the background process.
// It never starts a GUI merely to deliver a notification.
package notifications

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/clipboard"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

type Message struct {
	ID, Title, Body, Action, ActionLabel string
}

type Backend interface {
	Status(context.Context) (*pb.NotificationStatus, error)
	Request(context.Context) (*pb.NotificationStatus, error)
	Show(context.Context, Message) error
	Close()
}

type Options struct {
	Enabled      bool
	Loop         *clipboard.MainLoop
	Language     func() string
	Fetch, Apply func(context.Context, string) error
	Record       func(string) (*pb.ClipRecord, error)
}

type Manager struct {
	backend Backend
	opts    Options
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.Mutex
	pending map[string]Message
	wake    chan struct{}
	done    chan struct{}
}

func New(ctx context.Context, opts Options) *Manager {
	ctx, cancel := context.WithCancel(ctx)
	m := &Manager{opts: opts, ctx: ctx, cancel: cancel, pending: make(map[string]Message), wake: make(chan struct{}, 1), done: make(chan struct{})}
	m.backend = unavailable{}
	if opts.Enabled {
		backend, err := newBackend(opts.Loop, func(id, action string) {
			go func() {
				if err := m.Activate(ctx, id, action); err != nil {
					slog.Warn("通知操作失败", "err", err)
				}
			}()
		})
		if err != nil {
			slog.Warn("系统通知不可用", "err", err)
		} else {
			m.backend = backend
		}
	}
	go m.run()
	return m
}

func (m *Manager) Close() { m.cancel(); <-m.done; m.backend.Close() }
func (m *Manager) Status(ctx context.Context) (*pb.NotificationStatus, error) {
	return m.backend.Status(ctx)
}
func (m *Manager) Request(ctx context.Context) (*pb.NotificationStatus, error) {
	return m.backend.Request(ctx)
}

func (m *Manager) run() {
	defer close(m.done)
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-m.wake:
			m.mu.Lock()
			batch := m.pending
			m.pending = make(map[string]Message)
			m.mu.Unlock()
			for _, notice := range batch {
				ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
				if err := m.backend.Show(ctx, notice); err != nil {
					slog.Debug("发送系统通知失败", "err", err)
				}
				cancel()
			}
		}
	}
}

func (m *Manager) show(r *pb.ClipRecord, title, body, action, actionLabel string, failure string) {
	lang := ""
	if m.opts.Language != nil {
		lang = m.opts.Language()
	}
	if lang == "" {
		lang = systemLanguage()
	}
	name := r.GetTextPreview()
	if len(r.GetItems()) > 0 {
		name = r.Items[0].Name
		if len(r.Items) > 1 {
			name += fmt.Sprintf(" (+%d)", len(r.Items)-1)
		}
	}
	if len([]rune(name)) > 120 {
		name = string([]rune(name)[:117]) + "…"
	}
	replace := strings.NewReplacer("{name}", name, "{device}", r.GetOriginDeviceName(), "{size}", humanSize(r.GetTotalSize()), "{error}", failure)
	notice := Message{ID: r.Id, Title: translate(lang, title), Body: replace.Replace(translate(lang, body)), Action: action, ActionLabel: translate(lang, actionLabel)}
	m.mu.Lock()
	m.pending[r.Id] = notice
	m.mu.Unlock()
	select {
	case m.wake <- struct{}{}:
	default:
	}
}

func (m *Manager) Pending(r *pb.ClipRecord) {
	m.show(r, "notifyPendingTitle", "notifyPendingBody", "receive", "fetch", "")
}
func (m *Manager) Started(r *pb.ClipRecord) {
	m.show(r, "notifyReceivingTitle", "notifyReceivingBody", "open", "navHistory", "")
}
func (m *Manager) Failed(r *pb.ClipRecord) {
	m.show(r, "notifyFailedTitle", "notifyFailedBody", "receive", "retry", r.Error)
}
func (m *Manager) Completed(r *pb.ClipRecord, result *pb.ClipboardWriteResult) {
	switch result.GetOutcome() {
	case pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_APPLIED:
		m.show(r, "notifyReadyTitle", "notifyReadyBody", "open", "navHistory", "")
	case pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_FAILED:
		m.show(r, "notifyClipboardFailedTitle", "notifyClipboardFailedBody", "apply", "putOnClipboard", result.Error)
	default:
		m.show(r, "notifyReceivedTitle", "notifyReceivedBody", "apply", "putOnClipboard", "")
	}
}

func (m *Manager) Activate(ctx context.Context, id, action string) error {
	if action == "open" {
		return OpenUI()
	}
	r, err := m.opts.Record(id)
	if err != nil {
		return err
	}
	if r.Outgoing {
		return fmt.Errorf("notification does not refer to a received item")
	}
	switch action {
	case "receive":
		if r.Status == pb.ClipStatus_CLIP_STATUS_READY {
			m.Completed(r, &pb.ClipboardWriteResult{Outcome: pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_PRESERVED})
			return nil
		}
		return m.opts.Fetch(ctx, id)
	case "apply":
		result := &pb.ClipboardWriteResult{ClipId: id, Outcome: pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_APPLIED}
		if err := m.opts.Apply(ctx, id); err != nil {
			result.Outcome = pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_FAILED
			result.Error = err.Error()
		}
		m.Completed(r, result)
		return nil
	default:
		return fmt.Errorf("unknown notification action")
	}
}

// OpenUI resolves the sibling installed/portable UI, never a PATH executable.
func OpenUI() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		app := filepath.Clean(filepath.Join(filepath.Dir(exe), "../../../../.."))
		if filepath.Ext(app) != ".app" {
			return fmt.Errorf("UI app bundle not found")
		}
		cmd = exec.Command("/usr/bin/open", app)
	case "windows":
		cmd = exec.Command(filepath.Join(filepath.Dir(exe), "PasteFerry.exe"))
	default:
		cmd = exec.Command(filepath.Join(filepath.Dir(exe), "pasteferry"))
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}

type unavailable struct{}

func (unavailable) Status(context.Context) (*pb.NotificationStatus, error) {
	return &pb.NotificationStatus{Permission: "unavailable"}, nil
}
func (u unavailable) Request(ctx context.Context) (*pb.NotificationStatus, error) {
	return u.Status(ctx)
}
func (unavailable) Show(context.Context, Message) error { return nil }
func (unavailable) Close()                              {}

//go:embed catalog.json
var catalogJSON []byte
var catalog = func() map[string]map[string]string {
	var c map[string]map[string]string
	_ = json.Unmarshal(catalogJSON, &c)
	return c
}()

func translate(lang, key string) string {
	switch {
	case strings.HasPrefix(lang, "zh-Hant-HK") || strings.HasPrefix(lang, "zh-HK"):
		lang = "zh-Hant-HK"
	case strings.HasPrefix(lang, "zh-Hant") || strings.HasPrefix(lang, "zh-TW"):
		lang = "zh-Hant-TW"
	case strings.HasPrefix(lang, "zh"):
		lang = "zh-Hans"
	case strings.HasPrefix(lang, "ja"):
		lang = "ja"
	default:
		lang = "en"
	}
	if value := catalog[lang][key]; value != "" {
		return value
	}
	return catalog["en"][key]
}
func humanSize(n int64) string {
	if n >= 1<<30 {
		return fmt.Sprintf("%.1f GB", float64(n)/(1<<30))
	}
	return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
}
