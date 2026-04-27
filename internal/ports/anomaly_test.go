package ports

import (
	"strings"
	"testing"
)

func basePort(proto string, port uint16, pid int) Port {
	return Port{
		Protocol:  proto,
		LocalPort: port,
		State:     "LISTEN",
		PID:       pid,
	}
}

func baseTrend(proto string, port uint16, occ int) TrendSummary {
	return TrendSummary{
		Protocol:    proto,
		LocalPort:   port,
		Occurrences: occ,
	}
}

func TestDetectAnomalies_NewPort(t *testing.T) {
	current := []Port{basePort("tcp", 9999, 1234)}
	trend := []TrendSummary{}
	opts := DefaultAnomalyOptions()

	result := DetectAnomalies(current, trend, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(result))
	}
	if result[0].Kind != AnomalyNewPort {
		t.Errorf("expected AnomalyNewPort, got %s", result[0].Kind)
	}
}

func TestDetectAnomalies_RarePort(t *testing.T) {
	current := []Port{basePort("tcp", 8080, 1234)}
	trend := []TrendSummary{baseTrend("tcp", 8080, 1)}
	opts := DefaultAnomalyOptions()
	opts.MinOccurrences = 3

	result := DetectAnomalies(current, trend, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(result))
	}
	if result[0].Kind != AnomalyRarePort {
		t.Errorf("expected AnomalyRarePort, got %s", result[0].Kind)
	}
	if !strings.Contains(result[0].Message, "1 time(s)") {
		t.Errorf("message missing occurrence count: %s", result[0].Message)
	}
}

func TestDetectAnomalies_UnknownProc_PrivilegedPort(t *testing.T) {
	current := []Port{basePort("tcp", 80, 0)}
	trend := []TrendSummary{baseTrend("tcp", 80, 5)}
	opts := DefaultAnomalyOptions()

	result := DetectAnomalies(current, trend, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(result))
	}
	if result[0].Kind != AnomalyUnknownProc {
		t.Errorf("expected AnomalyUnknownProc, got %s", result[0].Kind)
	}
}

func TestDetectAnomalies_NoAnomalies(t *testing.T) {
	current := []Port{basePort("tcp", 443, 500)}
	trend := []TrendSummary{baseTrend("tcp", 443, 10)}
	opts := DefaultAnomalyOptions()

	result := DetectAnomalies(current, trend, opts)
	if len(result) != 0 {
		t.Errorf("expected no anomalies, got %d", len(result))
	}
}

func TestAnomaly_String(t *testing.T) {
	a := Anomaly{
		Kind:    AnomalyNewPort,
		Port:    basePort("tcp", 4444, 0),
		Message: "test message",
	}
	s := a.String()
	if !strings.Contains(s, "new_port") {
		t.Errorf("String() missing kind: %s", s)
	}
	if !strings.Contains(s, "4444") {
		t.Errorf("String() missing port: %s", s)
	}
}
