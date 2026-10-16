package cache

import "testing"

func TestReceivingCacheSurvivesCollectionAndReferenceSnapshotRace(t *testing.T) {
	c, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	rel, _, release, err := c.LeaseDir("receiving")
	if err != nil {
		t.Fatal(err)
	}
	unreferenced := func(string) (bool, error) { return false, nil }
	if removed, err := c.RemoveOrphan(rel, unreferenced); err != nil || removed {
		t.Fatal("removed a live receive", removed, err)
	}
	release()
	// The collector's original snapshot had no reference. A just-completed
	// transfer committed one before releasing its lease; recheck it.
	if removed, err := c.RemoveOrphan(rel, func(string) (bool, error) { return true, nil }); err != nil || removed {
		t.Fatal("removed a completed receive", removed, err)
	}
	if removed, err := c.RemoveOrphan(rel, unreferenced); err != nil || !removed {
		t.Fatal("orphan not collected", removed, err)
	}
}
