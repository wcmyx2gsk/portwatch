package ports

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleSuppressList() *SuppressList {
	return &SuppressList{
		Rules: []SuppressRule{
			{Port: 8080, Protocol: "tcp", Reason: "dev server"},
			{Port: 9090, Protocol: "tcp", Expires: time.Now().Add(1 * time.Hour)},
			{Port: 5353, Protocol: "udp", Expires: time.Now().Add(-1 * time.Hour)}, // expired
		},
	}
}

func TestIsSuppressed_ActiveRule(t *testing.T) {
	sl := sampleSuppressList()
	if !sl.IsSuppressed(8080, "tcp") {
		t.Error("expected port 8080/tcp to be suppressed")
	}
}

func TestIsSuppressed_NotInList(t *testing.T) {
	sl := sampleSuppressList()
	if sl.IsSuppressed(22, "tcp") {
		t.Error("expected port 22/tcp to not be suppressed")
	}
}

func TestIsSuppressed_ExpiredRule(t *testing.T) {
	sl := sampleSuppressList()
	if sl.IsSuppressed(5353, "udp") {
		t.Error("expected expired rule for 5353/udp to not suppress")
	}
}

func TestIsSuppressed_FutureExpiry(t *testing.T) {
	sl := sampleSuppressList()
	if !sl.IsSuppressed(9090, "tcp") {
		t.Error("expected port 9090/tcp with future expiry to be suppressed")
	}
}

func TestFilterSuppressed_RemovesMatches(t *testing.T) {
	sl := &SuppressList{
		Rules: []SuppressRule{
			{Port: 8080, Protocol: "tcp"},
		},
	}
	ports := []Port{
		{Port: 8080, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}
	result := sl.FilterSuppressed(ports)
	if len(result) != 1 || result[0].Port != 443 {
		t.Errorf("expected only port 443, got %v", result)
	}
}

func TestSaveAndLoadSuppressList_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suppress.json")

	original := sampleSuppressList()
	if err := SaveSuppressList(path, original); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := LoadSuppressList(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded.Rules) != len(original.Rules) {
		t.Errorf("expected %d rules, got %d", len(original.Rules), len(loaded.Rules))
	}
}

func TestLoadSuppressList_MissingFile_ReturnsEmpty(t *testing.T) {
	sl, err := LoadSuppressList("/nonexistent/suppress.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sl.Rules) != 0 {
		t.Errorf("expected empty list, got %d rules", len(sl.Rules))
	}
}

func TestLoadSuppressList_InvalidJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)

	_, err := LoadSuppressList(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	_ = json.Unmarshal // ensure import used
}
