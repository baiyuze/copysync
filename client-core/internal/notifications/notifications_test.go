package notifications

import (
	"context"
	"encoding/xml"
	"errors"
	"strings"
	"testing"
	"time"

	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

type fakeBackend struct{ messages chan Message }

func (*fakeBackend) Status(context.Context) (*pb.NotificationStatus, error) {
	return &pb.NotificationStatus{Supported: true, ActionsSupported: true, Permission: "allowed"}, nil
}
func (b *fakeBackend) Request(ctx context.Context) (*pb.NotificationStatus, error) {
	return b.Status(ctx)
}
func (b *fakeBackend) Show(_ context.Context, n Message) error { b.messages <- n; return nil }
func (*fakeBackend) Close()                                    {}

func nextMessage(t *testing.T, b *fakeBackend) Message {
	t.Helper()
	select {
	case n := <-b.messages:
		return n
	case <-time.After(time.Second):
		t.Fatal("notification not delivered")
		return Message{}
	}
}

func TestNotificationActionsAndClipboardOutcomes(t *testing.T) {
	r := &pb.ClipRecord{Id: "clip", OriginDeviceName: "Office PC", Items: []*pb.ClipItem{{Name: "file.zip"}}, TotalSize: 80 << 20}
	var receives, applies int
	m := New(context.Background(), Options{Language: func() string { return "en" }, Record: func(string) (*pb.ClipRecord, error) { return r, nil },
		Fetch: func(context.Context, string) error { receives++; return nil }, Apply: func(context.Context, string) error { applies++; return errors.New("clipboard busy") }})
	b := &fakeBackend{messages: make(chan Message, 16)}
	m.backend = b
	t.Cleanup(m.Close)
	m.Pending(r)
	n := nextMessage(t, b)
	if n.Action != "receive" || !strings.Contains(n.Body, "80.0 MB") || !strings.Contains(n.Body, "Office PC") {
		t.Fatal(n)
	}
	if err := m.Activate(context.Background(), r.Id, "receive"); err != nil || receives != 1 {
		t.Fatalf("receive: %d %v", receives, err)
	}
	r.Status = pb.ClipStatus_CLIP_STATUS_READY
	if err := m.Activate(context.Background(), r.Id, "receive"); err != nil {
		t.Fatal(err)
	}
	n = nextMessage(t, b)
	if n.Action != "apply" || receives != 1 || applies != 0 {
		t.Fatal("stale receive toast overwrote clipboard", n)
	}
	if err := m.Activate(context.Background(), r.Id, "apply"); err != nil {
		t.Fatal(err)
	}
	n = nextMessage(t, b)
	if n.Action != "apply" || strings.Contains(n.Title, "Ready to paste") || !strings.Contains(n.Body, "clipboard busy") {
		t.Fatal("write failure reported as paste-ready", n)
	}
	m.Completed(r, &pb.ClipboardWriteResult{Outcome: pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_APPLIED})
	n = nextMessage(t, b)
	if n.Title != "Ready to paste" {
		t.Fatal(n)
	}
	if err := m.Activate(context.Background(), r.Id, "unknown"); err == nil {
		t.Fatal("unknown action accepted")
	}
}

func TestNativeCatalogHasEveryLocaleAndPlaceholder(t *testing.T) {
	for _, lang := range []string{"en", "zh-Hans", "ja", "zh-Hant-TW", "zh-Hant-HK"} {
		for key := range catalog["en"] {
			if catalog[lang][key] == "" {
				t.Fatalf("missing %s/%s", lang, key)
			}
		}
		for _, key := range []string{"notifyPendingBody", "notifyReceivingBody", "notifyReadyBody", "notifyReceivedBody", "notifyFailedBody", "notifyClipboardFailedBody"} {
			text := strings.NewReplacer("{device}", "PC", "{name}", "file.bin", "{size}", "80 MB", "{error}", "offline").Replace(translate(lang, key))
			if strings.ContainsAny(text, "{}") || strings.Contains(text, `\n`) {
				t.Fatalf("unexpanded %s/%s: %q", lang, key, text)
			}
		}
	}
}

func TestWindowsToastEscapesFilenamesAndKeepsFullActionIdentity(t *testing.T) {
	n := Message{ID: "full-id/with spaces?&", Title: `<receive & paste>`, Body: `a"b & <file>`, Action: "receive", ActionLabel: `Receive "file"`}
	markup := toastXML(n)
	var parsed struct {
		Launch  string `xml:"launch,attr"`
		Actions []struct {
			Arguments string `xml:"arguments,attr"`
		} `xml:"actions>action"`
	}
	if err := xml.Unmarshal([]byte(markup), &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed.Launch != actionURI(n.ID, "open") || len(parsed.Actions) != 1 || parsed.Actions[0].Arguments != actionURI(n.ID, "receive") {
		t.Fatal(markup)
	}
	if len(notificationTag(n.ID)) > 16 || notificationTag(n.ID) == notificationTag("other") {
		t.Fatal("invalid toast identity")
	}
}
