package ports

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of open ports.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []Port    `json:"ports"`
}

// NewSnapshot creates a Snapshot from the given ports with the current time.
func NewSnapshot(ports []Port) Snapshot {
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
}

// SaveSnapshot writes a Snapshot as JSON to the given file path.
func SaveSnapshot(path string, snap Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// LoadSnapshot reads a Snapshot from the given file path.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Compare returns ports added and removed between two snapshots.
func Compare(old, new Snapshot) (added []Port, removed []Port) {
	oldSet := make(map[string]struct{}, len(old.Ports))
	newSet := make(map[string]struct{}, len(new.Ports))

	for _, p := range old.Ports {
		oldSet[p.Key()] = struct{}{}
	}
	for _, p := range new.Ports {
		newSet[p.Key()] = struct{}{}
	}

	for _, p := range new.Ports {
		if _, ok := oldSet[p.Key()]; !ok {
			added = append(added, p)
		}
	}
	for _, p := range old.Ports {
		if _, ok := newSet[p.Key()]; !ok {
			removed = append(removed, p)
		}
	}
	return added, removed
}
