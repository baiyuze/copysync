//go:build linux

package notifications

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
)

type desktopNotice struct {
	replaces uint32
	body     string
	actions  []string
}
type desktopFixture struct{ calls chan desktopNotice }

func (*desktopFixture) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"actions", "body-markup"}, nil
}
func (s *desktopFixture) Notify(_ string, replaces uint32, _ string, _ string, body string, actions []string, _ map[string]dbus.Variant, _ int32) (uint32, *dbus.Error) {
	s.calls <- desktopNotice{replaces, body, actions}
	return 17, nil
}

func TestLinuxDesktopNotificationActionsAndReplacement(t *testing.T) {
	if os.Getenv("DBUS_SESSION_BUS_ADDRESS") == "" {
		t.Skip("run under dbus-run-session for native D-Bus coverage")
	}
	bus, err := dbus.ConnectSessionBus()
	if err != nil {
		t.Fatal(err)
	}
	defer bus.Close()
	reply, err := bus.RequestName(desktopNotifications, dbus.NameFlagDoNotQueue)
	if err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		t.Skip("a real desktop notification service owns this session; use dbus-run-session")
	}
	fixture := &desktopFixture{calls: make(chan desktopNotice, 4)}
	if err := bus.Export(fixture, desktopNotificationPath, desktopNotifications); err != nil {
		t.Fatal(err)
	}
	actions := make(chan string, 1)
	b, err := newBackend(nil, func(id, action string) { actions <- id + ":" + action })
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	status, err := b.Status(ctx)
	if err != nil || !status.Supported || !status.ActionsSupported {
		t.Fatal(status, err)
	}
	n := Message{ID: "clip", Title: "Receive to paste", Body: `<b>file & name</b>`, Action: "receive", ActionLabel: "Receive"}
	if err := b.Show(ctx, n); err != nil {
		t.Fatal(err)
	}
	first := <-fixture.calls
	if first.replaces != 0 || len(first.actions) != 4 || strings.Contains(first.body, "<b>") {
		t.Fatal(first)
	}
	if err := bus.Emit(desktopNotificationPath, desktopNotifications+".ActionInvoked", uint32(17), "receive"); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-actions:
		if got != "clip:receive" {
			t.Fatal(got)
		}
	case <-ctx.Done():
		t.Fatal("desktop action not delivered")
	}
	if err := b.Show(ctx, n); err != nil {
		t.Fatal(err)
	}
	if next := <-fixture.calls; next.replaces != 17 {
		t.Fatal("notification was not replaced", next)
	}
}

func TestLinuxDesktopNotificationSmoke(t *testing.T) {
	if os.Getenv("PASTEFERRY_DESKTOP_SMOKE") != "1" {
		t.Skip("requires an interactive Linux desktop")
	}
	actions := make(chan string, 4)
	b, err := newBackend(nil, func(id, action string) { actions <- id + ":" + action })
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	status, err := b.Status(ctx)
	if err != nil || !status.Supported {
		t.Fatal(status, err)
	}
	t.Logf("native notification status: %v", status)
	if err := b.Show(ctx, Message{ID: "desktop-smoke", Title: "PasteFerry test: Receive to paste", Body: "Office PC · archive.zip · 80 MB\nReceive before pasting. History also offers Receive.", Action: "receive", ActionLabel: "Receive"}); err != nil {
		t.Fatal(err)
	}
	t.Log("waiting for the desktop Receive action")
	select {
	case action := <-actions:
		if action != "desktop-smoke:receive" {
			t.Fatal(action)
		}
	case <-ctx.Done():
		t.Fatal("native notification action was not clicked")
	}
	if err := b.Show(ctx, Message{ID: "desktop-smoke", Title: "PasteFerry test: Ready to paste", Body: "Native Receive action delivered successfully.", Action: "open", ActionLabel: "History"}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	if linux, ok := b.(*linuxBackend); ok {
		for _, id := range linux.ids {
			_ = linux.bus.Object(desktopNotifications, desktopNotificationPath).CallWithContext(ctx, desktopNotifications+".CloseNotification", 0, id).Err
		}
	}
}
