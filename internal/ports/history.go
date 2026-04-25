package ports

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// HistoryEntry records a snapshot of open ports at a point in time.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []Port    `json:"ports"`
}

// History holds a series of port snapshots.
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// AppendEntry adds a new snapshot to the history and persists it to path.
func AppendEntry(path string, ports []Port) error {
	h, err := LoadHistory(path)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}

	h.Entries = append(h.Entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	})

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling history: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing history file: %w", err)
	}
	return nil
}

// LoadHistory reads the history file at path. Returns an empty History if the
// file does not exist.
func LoadHistory(path string) (*History, error) {
	h := &History{}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading history file: %w", err)
	}

	if err := json.Unmarshal(data, h); err != nil {
		return nil, fmt.Errorf("unmarshalling history: %w", err)
	}
	return h, nil
}

// Latest returns the most recent HistoryEntry, or nil if there are none.
func (h *History) Latest() *HistoryEntry {
	if len(h.Entries) == 0 {
		return nil
	}
	return &h.Entries[len(h.Entries)-1]
}
