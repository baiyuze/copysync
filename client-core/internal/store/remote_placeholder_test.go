package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUpgradeRemovesOnlyPlaceholderRecords(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.db")
	source := filepath.Join(dir, ".uuremote_aeawv123")
	if err := os.WriteFile(source, nil, 0600); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	fixtures := []struct {
		id     string
		kind   Kind
		names  []string
		remove bool
	}{
		{"local", KindFile, []string{".uuremote_aeawv123"}, true},
		{"received", KindFile, []string{".uuremote_aeawv456", ".uuremote_aeawv789"}, true},
		{"mixed", KindFile, []string{".uuremote_aeawv123", "real.txt"}, false},
		{"mixed-reversed", KindFile, []string{"real.txt", ".uuremote_aeawv123"}, false},
		{"empty", KindFile, []string{"empty.txt"}, false},
		{"hidden", KindFile, []string{".gitignore"}, false},
		{"similar", KindFile, []string{".uuremote_aeawv123.backup"}, false},
		{"no-items", KindFile, nil, false},
		{"text", KindText, []string{".uuremote_aeawv123"}, false},
	}
	for _, f := range fixtures {
		c := Clip{ID: f.id, Kind: f.kind, Outgoing: f.id == "local", SourcePaths: []string{source},
			CreatedAt: time.Now(), ExpiresAt: time.Now().Add(time.Hour)}
		for _, name := range f.names {
			c.Items = append(c.Items, Item{Name: name})
		}
		if err := s.PutClip(c); err != nil {
			t.Fatal(err)
		}
	}
	s.Close()
	for range 2 { // cleanup is idempotent on subsequent launches
		s, err = Open(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range fixtures {
			_, err := s.GetClip(f.id)
			if f.remove && !errors.Is(err, ErrNotFound) || !f.remove && err != nil {
				t.Errorf("%s: remove=%v, err=%v", f.id, f.remove, err)
			}
		}
		var count int
		if err := s.db.QueryRow("SELECT count(*) FROM clip_items WHERE clip_id IN ('local','received')").Scan(&count); err != nil || count != 0 {
			t.Errorf("orphan items: %d, %v", count, err)
		}
		s.Close()
	}
	if _, err := os.Stat(source); err != nil {
		t.Fatalf("UU source file touched: %v", err)
	}
}
