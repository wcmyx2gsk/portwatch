package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/baseline"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single alert event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Notifier writes alerts to an output destination.
type Notifier struct {
	out io.Writer
}

// NewNotifier creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func NewNotifier(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes an alert to the output.
func (n *Notifier) Notify(a Alert) {
	fmt.Fprintf(n.out, "[%s] %s: %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Message,
	)
}

// NotifyDiff inspects a baseline.DiffResult and emits alerts for any
// added or removed ports.
func (n *Notifier) NotifyDiff(diff baseline.DiffResult) {
	for _, p := range diff.Added {
		n.Notify(Alert{
			Timestamp: time.Now(),
			Level:     LevelAlert,
			Message:   fmt.Sprintf("unexpected listener detected: %s (pid %d)", p.String(), p.PID),
		})
	}
	for _, p := range diff.Removed {
		n.Notify(Alert{
			Timestamp: time.Now(),
			Level:     LevelWarn,
			Message:   fmt.Sprintf("previously known listener gone: %s (pid %d)", p.String(), p.PID),
		})
	}
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		n.Notify(Alert{
			Timestamp: time.Now(),
			Level:     LevelInfo,
			Message:   "no changes detected",
		})
	}
}
