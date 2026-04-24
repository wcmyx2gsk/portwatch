package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/ports"
)

func samplePorts() []ports.Port {
	return []ports.Port{
		{Proto: "tcp", LocalAddr: "0.0.0.0", LocalPort: 22, State: "LISTEN"},
		{Proto: "tcp", LocalAddr: "0.0.0.0", LocalPort: 80, State: "LISTEN"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	if err := baseline.Save(path, samplePorts()); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(b.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(b.Ports))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestDiff_Added(t *testing.T) {
	base := &baseline.Baseline{Ports: samplePorts()}
	current := append(samplePorts(), ports.Port{Proto: "tcp", LocalAddr: "0.0.0.0", LocalPort: 443, State: "LISTEN"})

	added, removed := baseline.Diff(base, current)
	if len(added) != 1 || added[0].LocalPort != 443 {
		t.Errorf("expected 1 added port (443), got %v", added)
	}
	if len(removed) != 0 {
		t.Errorf("expected 0 removed ports, got %v", removed)
	}
}

func TestDiff_Removed(t *testing.T) {
	base := &baseline.Baseline{Ports: samplePorts()}
	current := []ports.Port{samplePorts()[0]} // only port 22

	added, removed := baseline.Diff(base, current)
	if len(added) != 0 {
		t.Errorf("expected 0 added ports, got %v", added)
	}
	if len(removed) != 1 || removed[0].LocalPort != 80 {
		t.Errorf("expected 1 removed port (80), got %v", removed)
	}
}

func TestDiff_NoChange(t *testing.T) {
	base := &baseline.Baseline{Ports: samplePorts()}
	added, removed := baseline.Diff(base, samplePorts())
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := baseline.Save("/nonexistent_dir/baseline.json", samplePorts())
	if err == nil {
		t.Fatal("expected error writing to invalid path")
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
