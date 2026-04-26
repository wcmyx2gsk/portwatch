package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func makeSummary(proto string, port uint16, occ, missed int) ports.TrendSummary {
	now := time.Now()
	return ports.TrendSummary{
		Port: ports.Port{
			Protocol:     proto,
			LocalAddress: "0.0.0.0",
			LocalPort:    port,
			State:        "LISTEN",
		},
		FirstSeen:   now.Add(-time.Hour),
		LastSeen:    now,
		Occurrences: occ,
		Missed:      missed,
	}
}

func TestTrendNotifier_Empty(t *testing.T) {
	var buf bytes.Buffer
	tn := NewTrendNotifier(&buf)
	tn.Notify(nil)
	if !strings.Contains(buf.String(), "no trend data") {
		t.Errorf("expected 'no trend data' message, got: %s", buf.String())
	}
}

func TestTrendNotifier_Notify_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	tn := NewTrendNotifier(&buf)
	tn.Notify([]ports.TrendSummary{makeSummary("tcp", 80, 3, 0)})
	out := buf.String()
	for _, header := range []string{"PROTO", "PORT", "ADDRESS", "OCCURRENCES", "MISSED"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestTrendNotifier_Notify_ContainsPortData(t *testing.T) {
	var buf bytes.Buffer
	tn := NewTrendNotifier(&buf)
	tn.Notify([]ports.TrendSummary{makeSummary("tcp", 8080, 5, 2)})
	out := buf.String()
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "tcp") {
		t.Errorf("expected protocol tcp in output")
	}
}

func TestTrendNotifier_NotifyUnstable_AllStable(t *testing.T) {
	var buf bytes.Buffer
	tn := NewTrendNotifier(&buf)
	tn.NotifyUnstable([]ports.TrendSummary{
		makeSummary("tcp", 80, 3, 0),
		makeSummary("tcp", 443, 3, 0),
	})
	if !strings.Contains(buf.String(), "stable") {
		t.Errorf("expected stable message, got: %s", buf.String())
	}
}

func TestTrendNotifier_NotifyUnstable_ShowsUnstableOnly(t *testing.T) {
	var buf bytes.Buffer
	tn := NewTrendNotifier(&buf)
	tn.NotifyUnstable([]ports.TrendSummary{
		makeSummary("tcp", 80, 3, 0),
		makeSummary("tcp", 9999, 1, 2),
	})
	out := buf.String()
	if !strings.Contains(out, "9999") {
		t.Errorf("expected unstable port 9999 in output")
	}
	if strings.Contains(out, "  80 ") {
		t.Errorf("stable port 80 should not appear in unstable output")
	}
}
