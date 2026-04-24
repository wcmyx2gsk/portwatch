package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()

	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.ScanInterval)
	}
	if cfg.BaselineFile != "baseline.json" {
		t.Errorf("unexpected baseline file: %s", cfg.BaselineFile)
	}
	if !cfg.AlertOnNew || !cfg.AlertOnGone {
		t.Error("expected both alert flags to be true by default")
	}
}

func TestLoad_MissingFile_ReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("/nonexistent/portwatch.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected default scan interval, got %v", cfg.ScanInterval)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	orig := config.Default()
	orig.ScanInterval = 10 * time.Second
	orig.BaselineFile = "/tmp/bl.json"
	orig.AlertOnGone = false
	orig.IgnoredPorts = []uint16{22, 80}

	if err := config.Save(path, orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.ScanInterval != orig.ScanInterval {
		t.Errorf("ScanInterval mismatch: got %v", loaded.ScanInterval)
	}
	if loaded.BaselineFile != orig.BaselineFile {
		t.Errorf("BaselineFile mismatch: got %s", loaded.BaselineFile)
	}
	if loaded.AlertOnGone != orig.AlertOnGone {
		t.Errorf("AlertOnGone mismatch")
	}
	if len(loaded.IgnoredPorts) != 2 || loaded.IgnoredPorts[0] != 22 {
		t.Errorf("IgnoredPorts mismatch: %v", loaded.IgnoredPorts)
	}
}

func TestLoad_InvalidYAML_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	_ = os.WriteFile(path, []byte("scan_interval: [not a duration"), 0o644)

	_, err := config.Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}
