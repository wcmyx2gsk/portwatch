package alert

import (
	"strings"
	"testing"

	"bytes"

	"github.com/user/portwatch/internal/ports"
)

func makeAnomaly(kind ports.AnomalyKind, proto string, port uint16) ports.Anomaly {
	return ports.Anomaly{
		Kind: kind,
		Port: ports.Port{
			Protocol:  proto,
			LocalPort: port,
			State:     "LISTEN",
		},
		Message: "test anomaly",
	}
}

func TestAnomalyNotifier_NoAnomalies(t *testing.T) {
	var buf bytes.Buffer
	n := NewAnomalyNotifier(&buf)
	n.Notify(nil)

	out := buf.String()
	if !strings.Contains(out, "no anomalies") {
		t.Errorf("expected 'no anomalies' in output, got: %s", out)
	}
}

func TestAnomalyNotifier_WithAnomalies(t *testing.T) {
	var buf bytes.Buffer
	n := NewAnomalyNotifier(&buf)
	anomalies := []ports.Anomaly{
		makeAnomaly(ports.AnomalyNewPort, "tcp", 9000),
		makeAnomaly(ports.AnomalyRarePort, "udp", 5353),
	}
	n.Notify(anomalies)

	out := buf.String()
	if !strings.Contains(out, "2 anomaly(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "9000") {
		t.Errorf("expected port 9000 in output, got: %s", out)
	}
	if !strings.Contains(out, "5353") {
		t.Errorf("expected port 5353 in output, got: %s", out)
	}
}

func TestAnomalyNotifier_NotifyIfCritical_ReturnsFalseWhenEmpty(t *testing.T) {
	var buf bytes.Buffer
	n := NewAnomalyNotifier(&buf)
	result := n.NotifyIfCritical(nil)

	if result {
		t.Error("expected false when no anomalies")
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestAnomalyNotifier_NotifyIfCritical_ReturnsTrueWhenPresent(t *testing.T) {
	var buf bytes.Buffer
	n := NewAnomalyNotifier(&buf)
	anomalies := []ports.Anomaly{makeAnomaly(ports.AnomalyNewPort, "tcp", 1234)}
	result := n.NotifyIfCritical(anomalies)

	if !result {
		t.Error("expected true when anomalies present")
	}
	if !strings.Contains(buf.String(), "1234") {
		t.Errorf("expected port in output: %s", buf.String())
	}
}

func TestNewAnomalyNotifier_NilWriterUsesStdout(t *testing.T) {
	n := NewAnomalyNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.w == nil {
		t.Error("expected writer to be set")
	}
}
