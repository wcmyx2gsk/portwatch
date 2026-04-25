package ports

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleSnapshot() Snapshot {
	return Snapshot{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Ports: []Port{
			{Protocol: "tcp", LocalAddr: "0.0.0.0", LocalPort: 22, State: "LISTEN", PID: 100},
			{Protocol: "tcp", LocalAddr: "0.0.0.0", LocalPort: 80, State: "LISTEN", PID: 200},
		},
	}
}

func TestNewSnapshot_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	snap := NewSnapshot(sampleSnapshot().Ports)
	after := time.Now().UTC()

	if snap.Timestamp.Before(before) || snap.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", snap.Timestamp, before, after)
	}
	if len(snap.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(snap.Ports))
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := sampleSnapshot()
	if err := SaveSnapshot(path, orig); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if !loaded.Timestamp.Equal(orig.Timestamp) {
		t.Errorf("timestamp mismatch: got %v want %v", loaded.Timestamp, orig.Timestamp)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("port count mismatch: got %d want %d", len(loaded.Ports), len(orig.Ports))
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

func TestCompare_AddedAndRemoved(t *testing.T) {
	old := Snapshot{Ports: []Port{
		{Protocol: "tcp", LocalPort: 22},
		{Protocol: "tcp", LocalPort: 80},
	}}
	new := Snapshot{Ports: []Port{
		{Protocol: "tcp", LocalPort: 22},
		{Protocol: "tcp", LocalPort: 443},
	}}

	added, removed := Compare(old, new)

	if len(added) != 1 || added[0].LocalPort != 443 {
		t.Errorf("expected added port 443, got %+v", added)
	}
	if len(removed) != 1 || removed[0].LocalPort != 80 {
		t.Errorf("expected removed port 80, got %+v", removed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	snap := sampleSnapshot()
	added, removed := Compare(snap, snap)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no changes, got added=%v removed=%v", added, removed)
	}
}
