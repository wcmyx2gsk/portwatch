package watch_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/watch"
)

func writeTempBaseline(t *testing.T, ps []ports.Port) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	if err := baseline.Save(path, ps); err != nil {
		t.Fatalf("writeTempBaseline: %v", err)
	}
	return path
}

func TestWatcher_Stop(t *testing.T) {
	var buf bytes.Buffer
	notifier := alert.NewNotifier(&buf)

	cfg := config.Default()
	cfg.BaselinePath = writeTempBaseline(t, []ports.Port{})
	cfg.IntervalSeconds = 1

	w := watch.New(cfg, notifier)

	done := make(chan struct{})
	go func() {
		w.Start()
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	w.Stop()

	select {
	case <-done:
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("watcher did not stop in time")
	}
}

func TestWatcher_New_NotNil(t *testing.T) {
	var buf bytes.Buffer
	notifier := alert.NewNotifier(&buf)
	cfg := config.Default()
	w := watch.New(cfg, notifier)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

// Ensure baseline path is honoured: a missing baseline causes no panic.
func TestWatcher_MissingBaseline_NocrashOnCycle(t *testing.T) {
	var buf bytes.Buffer
	notifier := alert.NewNotifier(&buf)
	cfg := config.Default()
	cfg.BaselinePath = filepath.Join(t.TempDir(), "nonexistent.json")
	cfg.IntervalSeconds = 60 // long interval so we don't actually tick

	// Verify the config is serialisable (sanity check).
	if _, err := json.Marshal(cfg); err != nil {
		t.Fatalf("config marshal: %v", err)
	}

	w := watch.New(cfg, notifier)
	w.Stop() // stop immediately; just ensure construction is fine
	_ = os.Remove(cfg.BaselinePath)
}
