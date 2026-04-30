package ports

import (
	"testing"
)

func baseScorePort(port uint16, addr string, proc *ProcessInfo) Port {
	return Port{
		Port:     port,
		Protocol: "tcp",
		Address:  addr,
		Process:  proc,
	}
}

func TestRiskLevel_Thresholds(t *testing.T) {
	tests := []struct {
		score    int
		expected RiskLevel
	}{
		{0, RiskLow},
		{19, RiskLow},
		{20, RiskMedium},
		{39, RiskMedium},
		{40, RiskHigh},
		{59, RiskHigh},
		{60, RiskCritical},
		{100, RiskCritical},
	}
	for _, tt := range tests {
		got := riskLevel(tt.score)
		if got != tt.expected {
			t.Errorf("riskLevel(%d) = %s, want %s", tt.score, got, tt.expected)
		}
	}
}

func TestScorePort_KnownRiskyPort(t *testing.T) {
	opts := DefaultScoreOptions()
	p := baseScorePort(22, "127.0.0.1", &ProcessInfo{Name: "sshd"})
	// port 22 is not in KnownRiskyPorts by default, but 21 is
	p.Port = 21
	result := scorePort(p, opts)
	found := false
	for _, r := range result.Reasons {
		if r == "known risky port" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'known risky port' reason for port 21")
	}
	if result.Score < 30 {
		t.Errorf("expected score >= 30 for risky port, got %d", result.Score)
	}
}

func TestScorePort_UnknownProcess_AddsScore(t *testing.T) {
	opts := DefaultScoreOptions()
	p := baseScorePort(8080, "0.0.0.0", nil)
	result := scorePort(p, opts)
	found := false
	for _, r := range result.Reasons {
		if r == "unknown process" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'unknown process' reason")
	}
}

func TestScorePort_ExternalAddress_AddsScore(t *testing.T) {
	opts := DefaultScoreOptions()
	p := baseScorePort(9000, "0.0.0.0", &ProcessInfo{Name: "app"})
	result := scorePort(p, opts)
	found := false
	for _, r := range result.Reasons {
		if r == "externally reachable" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'externally reachable' reason for 0.0.0.0")
	}
}

func TestScorePorts_ReturnsOnePerPort(t *testing.T) {
	opts := DefaultScoreOptions()
	ports := []Port{
		baseScorePort(80, "127.0.0.1", &ProcessInfo{Name: "nginx"}),
		baseScorePort(443, "127.0.0.1", &ProcessInfo{Name: "nginx"}),
	}
	results := ScorePorts(ports, opts)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestDefaultScoreOptions_NotEmpty(t *testing.T) {
	opts := DefaultScoreOptions()
	if len(opts.KnownRiskyPorts) == 0 {
		t.Error("expected non-empty KnownRiskyPorts")
	}
	if opts.PrivilegedThreshold == 0 {
		t.Error("expected non-zero PrivilegedThreshold")
	}
}
