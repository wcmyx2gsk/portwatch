package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// Baseline represents a saved snapshot of known-good open ports.
type Baseline struct {
	CreatedAt time.Time    `json:"created_at"`
	Ports     []ports.Port `json:"ports"`
}

// Save writes the current port list to the baseline file.
func Save(path string, openPorts []ports.Port) error {
	b := Baseline{
		CreatedAt: time.Now(),
		Ports:     openPorts,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal baseline: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write baseline file: %w", err)
	}
	return nil
}

// Load reads a previously saved baseline from disk.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no baseline found at %s; run 'portwatch baseline save' first", path)
		}
		return nil, fmt.Errorf("read baseline file: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("parse baseline file: %w", err)
	}
	return &b, nil
}

// Diff compares a current port list against the baseline.
// Returns new ports not in the baseline and ports that have disappeared.
func Diff(base *Baseline, current []ports.Port) (added []ports.Port, removed []ports.Port) {
	baseSet := make(map[string]struct{}, len(base.Ports))
	for _, p := range base.Ports {
		baseSet[p.String()] = struct{}{}
	}
	currentSet := make(map[string]struct{}, len(current))
	for _, p := range current {
		currentSet[p.String()] = struct{}{}
		if _, ok := baseSet[p.String()]; !ok {
			added = append(added, p)
		}
	}
	for _, p := range base.Ports {
		if _, ok := currentSet[p.String()]; !ok {
			removed = append(removed, p)
		}
	}
	return added, removed
}

// Age returns the duration since the baseline was created.
func (b *Baseline) Age() time.Duration {
	return time.Since(b.CreatedAt).Truncate(time.Second)
}
