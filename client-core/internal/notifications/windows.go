//go:build windows

package notifications

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/baiyuze/pasteferry/client-core/internal/clipboard"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const windowsAppID = "PasteFerry.Desktop"

type windowsBackend struct{ appID string }

func newBackend(_ *clipboard.MainLoop, _ func(string, string)) (Backend, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	// Per-user registration works for the portable distribution too. Actions
	// launch a short-lived client; it authenticates to the existing daemon.
	entries := []struct{ path, name, value string }{
		{`Software\Classes\` + protocolName, "", "PasteFerry notification action"},
		{`Software\Classes\` + protocolName, "URL Protocol", ""},
		{`Software\Classes\` + protocolName + `\shell\open\command`, "", `"` + exe + `" -notification-action "%1"`},
		{`Software\Classes\AppUserModelId\` + windowsAppID, "DisplayName", "PasteFerry"},
	}
	for _, entry := range entries {
		key, _, err := registry.CreateKey(registry.CURRENT_USER, entry.path, registry.SET_VALUE)
		if err != nil {
			return nil, err
		}
		err = key.SetStringValue(entry.name, entry.value)
		key.Close()
		if err != nil {
			return nil, err
		}
	}
	return windowsBackend{appID: windowsAppID}, nil
}

func (b windowsBackend) Status(ctx context.Context) (*pb.NotificationStatus, error) {
	output, err := powershell(ctx, toastPrelude+`$n=[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('`+b.appID+`'); try{$setting=$n.get_Setting(); if($null -eq $setting){'unknown'}else{$setting.ToString()}}catch{'unknown'}`)
	if err != nil {
		return &pb.NotificationStatus{Permission: "unavailable"}, nil
	}
	permission := "unknown"
	switch string(bytes.TrimSpace(output)) {
	case "Enabled":
		permission = "allowed"
	case "DisabledForApplication", "DisabledForUser", "DisabledByGroupPolicy", "DisabledByManifest":
		permission = "denied"
	}
	return &pb.NotificationStatus{Supported: true, ActionsSupported: true, Permission: permission}, nil
}
func (b windowsBackend) Request(ctx context.Context) (*pb.NotificationStatus, error) {
	cmd := exec.Command("explorer.exe", "ms-settings:notifications")
	_ = cmd.Start()
	if cmd.Process != nil {
		go cmd.Wait()
	}
	return b.Status(ctx)
}

const toastPrelude = `$ErrorActionPreference='Stop'; [Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType=WindowsRuntime] > $null; [Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType=WindowsRuntime] > $null; [Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType=WindowsRuntime] > $null; `

func (b windowsBackend) Show(ctx context.Context, n Message) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(toastXML(n)))
	// All dynamic XML is data, never PowerShell source.
	script := toastPrelude + `$xml=New-Object Windows.Data.Xml.Dom.XmlDocument; $xml.LoadXml([Text.Encoding]::UTF8.GetString([Convert]::FromBase64String('` + encoded + `'))); $toast=[Windows.UI.Notifications.ToastNotification]::new($xml); $toast.Tag='` + notificationTag(n.ID) + `'; $toast.Group='receive'; [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('` + b.appID + `').Show($toast)`
	_, err := powershell(ctx, script)
	return err
}
func (windowsBackend) Close() {}
func powershell(ctx context.Context, script string) ([]byte, error) {
	var buf bytes.Buffer
	for _, value := range utf16.Encode([]rune(script)) {
		_ = binary.Write(&buf, binary.LittleEndian, value)
	}
	cmd := exec.CommandContext(ctx, "powershell.exe", "-NoLogo", "-NoProfile", "-NonInteractive", "-EncodedCommand", base64.StdEncoding.EncodeToString(buf.Bytes()))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("Windows notification: %w", ctx.Err())
		}
		return nil, fmt.Errorf("Windows notification: %w: %s", err, bytes.TrimSpace(output))
	}
	return output, nil
}
func systemLanguage() string {
	var name [85]uint16
	proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetUserDefaultLocaleName")
	count, _, _ := proc.Call(uintptr(unsafe.Pointer(&name[0])), uintptr(len(name)))
	if count == 0 {
		return "en"
	}
	return windows.UTF16ToString(name[:])
}
