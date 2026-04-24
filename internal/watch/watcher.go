package watch

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ports"
)

// Watcher periodically scans open ports and alerts on changes from baseline.
type Watcher struct {
	cfg      *config.Config
	notifier *alert.Notifier
	stopCh   chan struct{}
}

// New creates a new Watcher using the provided config and notifier.
func New(cfg *config.Config, notifier *alert.Notifier) *Watcher {
	return &Watcher{
		cfg:      cfg,
		notifier: notifier,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the watch loop, scanning at the configured interval.
// It blocks until Stop is called.
func (w *Watcher) Start() {
	interval := time.Duration(w.cfg.IntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("portwatch: starting watcher (interval=%s, baseline=%s)",
		interval, w.cfg.BaselinePath)

	for {
		select {
		case <-ticker.C:
			w.runCycle()
		case <-w.stopCh:
			log.Println("portwatch: watcher stopped")
			return
		}
	}
}

// Stop signals the watcher to cease operation.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

func (w *Watcher) runCycle() {
	current, err := ports.Scan()
	if err != nil {
		log.Printf("portwatch: scan error: %v", err)
		return
	}

	saved, err := baseline.Load(w.cfg.BaselinePath)
	if err != nil {
		log.Printf("portwatch: baseline load error: %v", err)
		return
	}

	diff := baseline.Diff(saved, current)
	if err := w.notifier.NotifyDiff(diff); err != nil {
		log.Printf("portwatch: notify error: %v", err)
	}
}
