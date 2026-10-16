//go:build windows

package notifications

import (
	"context"
	"encoding/base64"
	"golang.org/x/sys/windows/registry"
	"os"
	"strings"
	"testing"
	"time"
)

func TestWindowsNativeToastXML(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	markup := base64.StdEncoding.EncodeToString([]byte(toastXML(Message{ID: "test", Title: "Receive", Body: `<file & name>`, Action: "receive", ActionLabel: "Receive"})))
	output, err := powershell(ctx, toastPrelude+`$xml=New-Object Windows.Data.Xml.Dom.XmlDocument; $xml.LoadXml([Text.Encoding]::UTF8.GetString([Convert]::FromBase64String('`+markup+`'))); [Windows.UI.Notifications.ToastNotification]::new($xml).Content.GetXml()`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "pasteferry-notify:") {
		t.Fatalf("native toast lost its activation URI: %s", output)
	}
}

func TestWindowsDesktopNotificationSmoke(t *testing.T) {
	if os.Getenv("PASTEFERRY_DESKTOP_SMOKE") != "1" {
		t.Skip("requires an interactive Windows desktop")
	}
	const appID = "PasteFerry.Validation.LargeFile"
	keyPath := `Software\Classes\AppUserModelId\` + appID
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	if err := key.SetStringValue("DisplayName", "PasteFerry validation"); err != nil {
		t.Fatal(err)
	}
	key.Close()
	defer registry.DeleteKey(registry.CURRENT_USER, keyPath)
	b := windowsBackend{appID: appID}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	status, err := b.Status(ctx)
	if err != nil || !status.Supported {
		t.Fatalf("native status: %v %v", status, err)
	}
	t.Logf("native notification status: %v", status)
	for _, n := range []Message{
		{ID: "desktop-smoke", Title: "PasteFerry test: Receive to paste", Body: "Office PC · archive.zip · 80 MB\nReceive this item before pasting.", Action: "receive", ActionLabel: "Receive"},
		{ID: "desktop-smoke", Title: "PasteFerry test: Ready to paste", Body: "archive.zip\nReceived and placed on the clipboard.", Action: "open", ActionLabel: "History"},
	} {
		if err := b.Show(ctx, n); err != nil {
			t.Fatal(err)
		}
		t.Log(n.Title)
		time.Sleep(8 * time.Second)
	}
	_, err = powershell(ctx, toastPrelude+`[Windows.UI.Notifications.ToastNotificationManager]::History.Clear('`+appID+`')`)
	if err != nil {
		t.Fatal(err)
	}
}
