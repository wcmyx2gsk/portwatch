package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/ports"
)

func TestNotify_WritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	a := alert.Alert{
		Level:   alert.LevelAlert,
		Message: "test message",
	}
	n.Notify(a)

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "test message") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestNotifyDiff_Added(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	diff := baseline.DiffResult{
		Added: []ports.Port{
			{Proto: "tcp", Addr: "0.0.0.0", Port: 8080, PID: 42},
		},
	}
	n.NotifyDiff(diff)

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level for added port, got: %s", out)
	}
	if !strings.Contains(out, "unexpected listener") {
		t.Errorf("expected 'unexpected listener' in output, got: %s", out)
	}
}

func TestNotifyDiff_Removed(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	diff := baseline.DiffResult{
		Removed: []ports.Port{
			{Proto: "tcp", Addr: "127.0.0.1", Port: 9090, PID: 7},
		},
	}
	n.NotifyDiff(diff)

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level for removed port, got: %s", out)
	}
	if !strings.Contains(out, "previously known listener gone") {
		t.Errorf("expected removal message in output, got: %s", out)
	}
}

func TestNotifyDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	n.NotifyDiff(baseline.DiffResult{})

	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO level when no changes, got: %s", out)
	}
	if !strings.Contains(out, "no changes detected") {
		t.Errorf("expected 'no changes detected' in output, got: %s", out)
	}
}

func TestNewNotifier_DefaultsToStdout(t *testing.T) {
	n := alert.NewNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
