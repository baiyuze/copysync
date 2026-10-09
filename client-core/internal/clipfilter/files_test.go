package clipfilter

import (
	"reflect"
	"testing"
)

func TestRemotePlaceholder(t *testing.T) {
	for _, name := range []string{
		".uuremote_aeawv665203311565082",
		"/Users/me/Library/Application Support/com.netease.uuremote.server/Clipboard/.uuremote_aeawv1",
		`C:\Users\me\AppData\Local\Temp\.uuremote_aeawv123`,
	} {
		if !RemotePlaceholder(name) {
			t.Errorf("missed placeholder %q", name)
		}
	}
	for _, name := range []string{
		"", ".gitignore", "empty.txt", ".uuremote_", ".uuremote_aeawv",
		".uuremote_aeawv123.backup", "uuremote_aeawv123", ".uuremote_aeawvabc",
		"/tmp/.uuremote_aeawv123/real.txt", "/Clipboard/real.txt",
	} {
		if RemotePlaceholder(name) {
			t.Errorf("ordinary file incorrectly filtered: %q", name)
		}
	}
	paths := []string{"/tmp/a.txt", "/tmp/.uuremote_aeawv123", "/tmp/.gitignore"}
	want := []string{"/tmp/a.txt", "/tmp/.gitignore"}
	if got := UserFiles(paths); !reflect.DeepEqual(got, want) {
		t.Fatalf("mixed selection = %v, want %v", got, want)
	}
	if paths[1] != "/tmp/.uuremote_aeawv123" {
		t.Fatal("modified input slice")
	}
}
