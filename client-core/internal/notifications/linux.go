//go:build linux

package notifications

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/baiyuze/pasteferry/client-core/internal/clipboard"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	"github.com/godbus/dbus/v5"
)

const desktopNotifications = "org.freedesktop.Notifications"
const desktopNotificationPath = dbus.ObjectPath("/org/freedesktop/Notifications")

type linuxBackend struct {
	bus     *dbus.Conn
	mu      sync.Mutex
	ids     map[string]uint32
	records map[uint32]string
	action  func(string, string)
	signals chan *dbus.Signal
	done    chan struct{}
}

func newBackend(_ *clipboard.MainLoop, action func(string, string)) (Backend, error) {
	bus, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}
	b := &linuxBackend{bus: bus, ids: make(map[string]uint32), records: make(map[uint32]string), action: action, signals: make(chan *dbus.Signal, 32), done: make(chan struct{})}
	if err := bus.AddMatchSignal(dbus.WithMatchObjectPath(desktopNotificationPath), dbus.WithMatchInterface(desktopNotifications)); err != nil {
		bus.Close()
		return nil, err
	}
	bus.Signal(b.signals)
	go func() {
		for {
			select {
			case <-b.done:
				return
			case signal := <-b.signals:
				if signal == nil || signal.Name != desktopNotifications+".ActionInvoked" || len(signal.Body) != 2 {
					continue
				}
				id, ok := signal.Body[0].(uint32)
				if !ok {
					continue
				}
				verb, ok := signal.Body[1].(string)
				if !ok {
					continue
				}
				b.mu.Lock()
				record := b.records[id]
				b.mu.Unlock()
				if record != "" {
					if verb == "default" {
						verb = "open"
					}
					b.action(record, verb)
				}
			}
		}
	}()
	return b, nil
}
func (b *linuxBackend) Status(ctx context.Context) (*pb.NotificationStatus, error) {
	var capabilities []string
	err := b.bus.Object(desktopNotifications, desktopNotificationPath).CallWithContext(ctx, desktopNotifications+".GetCapabilities", 0).Store(&capabilities)
	if err != nil {
		return &pb.NotificationStatus{Permission: "unavailable"}, nil
	}
	status := &pb.NotificationStatus{Supported: true, Permission: "unknown"}
	for _, capability := range capabilities {
		if capability == "actions" {
			status.ActionsSupported = true
		}
	}
	return status, nil // Linux has no portable per-app permission/DND query.
}
func (b *linuxBackend) Request(ctx context.Context) (*pb.NotificationStatus, error) {
	return b.Status(ctx)
}
func (b *linuxBackend) Show(ctx context.Context, n Message) error {
	status, err := b.Status(ctx)
	if err != nil || !status.Supported {
		return err
	}
	actions := []string{}
	if status.ActionsSupported {
		actions = []string{"default", "", n.Action, n.ActionLabel}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	var id uint32
	err = b.bus.Object(desktopNotifications, desktopNotificationPath).CallWithContext(ctx, desktopNotifications+".Notify", 0,
		"PasteFerry", b.ids[n.ID], "com.pasteferry.app", n.Title, escapeMarkup(n.Body), actions,
		map[string]dbus.Variant{"desktop-entry": dbus.MakeVariant("com.pasteferry.app"), "suppress-sound": dbus.MakeVariant(true)}, int32(-1)).Store(&id)
	if err == nil {
		delete(b.records, b.ids[n.ID])
		b.ids[n.ID] = id
		b.records[id] = n.ID
	}
	return err
}
func (b *linuxBackend) Close() { close(b.done); b.bus.RemoveSignal(b.signals); b.bus.Close() }
func systemLanguage() string {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if value := os.Getenv(key); value != "" {
			return strings.ReplaceAll(strings.Split(value, ".")[0], "_", "-")
		}
	}
	return "en"
}
func escapeMarkup(value string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(value)
}
