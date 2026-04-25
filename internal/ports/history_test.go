package ports

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleHistoryPorts() []Port {
	return []Port{
		{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 8080, State: "LISTEN", PID: 42},
		{Protocol: "tcp", LocalAddress: "127.0.0.1", LocalPort: 5432, State: "LISTEN", PID: 99},
	}
}

func TestLoadHistory_MissingFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	h, err := LoadHistory(filepath.Join(dir, "history.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(h.Entries))
	}
}

func TestAppendEntry_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	if err := AppendEntry(path, sampleHistoryPorts()); err != nil {
		t.Fatalf("AppendEntry failed: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("history file was not created")
	}
}

func TestAppendEntry_AccumulatesEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	for i := 0; i < 3; i++ {
		if err := AppendEntry(path, sampleHistoryPorts()); err != nil {
			t.Fatalf("AppendEntry iteration %d failed: %v", i, err)
		}
	}

	h, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}
	if len(h.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(h.Entries))
	}
}

func TestHistory_Latest_ReturnsLastEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	_ = AppendEntry(path, sampleHistoryPorts())
	time.Sleep(2 * time.Millisecond)
	_ = AppendEntry(path, sampleHistoryPorts()[:1])

	h, _ := LoadHistory(path)
	latest := h.Latest()
	if latest == nil {
		t.Fatal("expected a latest entry, got nil")
	}
	if len(latest.Ports) != 1 {
		t.Errorf("expected latest entry to have 1 port, got %d", len(latest.Ports))
	}
}

func TestHistory_Latest_EmptyHistory(t *testing.T) {
	h := &History{}
	if h.Latest() != nil {
		t.Error("expected nil for empty history")
	}
}
