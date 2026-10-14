package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateLegacyMovesOldDataDir(t *testing.T) {
	base := t.TempDir()
	old := filepath.Join(base, "CopySync")
	if err := os.MkdirAll(filepath.Join(old, "cache"), 0o700); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"identity.json": "id", "copysync.db": "db", "copysync.db-wal": "wal"} {
		if err := os.WriteFile(filepath.Join(old, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	p := PathsUnder(filepath.Join(base, "PasteFerry"))
	moved, err := MigrateLegacy(p)
	if err != nil || !moved {
		t.Fatalf("MigrateLegacy = %v, %v; want true, nil", moved, err)
	}
	for path, want := range map[string]string{
		filepath.Join(p.Root, "identity.json"): "id",
		p.DB:                                   "db",
		p.DB + "-wal":                          "wal",
	} {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != want {
			t.Errorf("%s = %q, %v; want %q", path, got, err, want)
		}
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Errorf("old dir still exists: %v", err)
	}
}

func TestMigrateLegacyLeavesExistingDataDirAlone(t *testing.T) {
	base := t.TempDir()
	for _, d := range []string{"CopySync", "PasteFerry"} {
		if err := os.MkdirAll(filepath.Join(base, d), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	moved, err := MigrateLegacy(PathsUnder(filepath.Join(base, "PasteFerry")))
	if err != nil || moved {
		t.Fatalf("MigrateLegacy = %v, %v; want false, nil", moved, err)
	}
	if _, err := os.Stat(filepath.Join(base, "CopySync")); err != nil {
		t.Errorf("old dir should be untouched: %v", err)
	}
}

func TestMigrateLegacyWithoutOldDir(t *testing.T) {
	moved, err := MigrateLegacy(PathsUnder(filepath.Join(t.TempDir(), "pasteferry")))
	if err != nil || moved {
		t.Fatalf("MigrateLegacy = %v, %v; want false, nil", moved, err)
	}
}
