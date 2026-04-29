package ports

import (
	"encoding/json"
	"os"
	"time"
)

// SuppressRule defines a rule to silence alerts for a specific port/protocol.
type SuppressRule struct {
	Port     int       `json:"port"`
	Protocol string    `json:"protocol"`
	Reason   string    `json:"reason,omitempty"`
	Expires  time.Time `json:"expires,omitempty"`
}

// SuppressList holds a collection of suppression rules.
type SuppressList struct {
	Rules []SuppressRule `json:"rules"`
}

// IsSuppressed returns true if the given port+protocol matches an active rule.
func (sl *SuppressList) IsSuppressed(port int, protocol string) bool {
	now := time.Now()
	for _, r := range sl.Rules {
		if r.Port == port && r.Protocol == protocol {
			if r.Expires.IsZero() || r.Expires.After(now) {
				return true
			}
		}
	}
	return false
}

// FilterSuppressed removes ports that match active suppression rules.
func (sl *SuppressList) FilterSuppressed(ports []Port) []Port {
	var result []Port
	for _, p := range ports {
		if !sl.IsSuppressed(p.Port, p.Protocol) {
			result = append(result, p)
		}
	}
	return result
}

// LoadSuppressList reads a suppression list from a JSON file.
// Returns an empty list if the file does not exist.
func LoadSuppressList(path string) (*SuppressList, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &SuppressList{}, nil
	}
	if err != nil {
		return nil, err
	}
	var sl SuppressList
	if err := json.Unmarshal(data, &sl); err != nil {
		return nil, err
	}
	return &sl, nil
}

// SaveSuppressList writes the suppression list to a JSON file.
func SaveSuppressList(path string, sl *SuppressList) error {
	data, err := json.MarshalIndent(sl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
