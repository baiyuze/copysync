package store

import (
	"path/filepath"
	"testing"
)

func TestRecoverReceivesOnlyChangesInterruptedIncomingTransfers(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	for _, clip := range []Clip{
		{ID: "interrupted", Status: StatusFetching},
		{ID: "ready", Status: StatusReady, CachePath: "downloaded"},
		{ID: "outgoing", Status: StatusFetching, Outgoing: true},
	} {
		if err := s.PutClip(clip); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.RecoverReceives(); err != nil {
		t.Fatal(err)
	}
	for id, want := range map[string]Status{"interrupted": StatusFailed, "ready": StatusReady, "outgoing": StatusFetching} {
		got, err := s.GetClip(id)
		if err != nil || got.Status != want {
			t.Fatalf("%s: %v %v", id, got.Status, err)
		}
	}
}
