package ports

import "time"

// RiskLevel represents the severity of a port's risk score.
type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

// PortScore holds a computed risk score for a single port entry.
type PortScore struct {
	Port      Port
	Score     int
	Level     RiskLevel
	Reasons   []string
}

// ScoreOptions controls which heuristics are applied during scoring.
type ScoreOptions struct {
	PrivilegedThreshold uint16
	KnownRiskyPorts     []uint16
	MaxAgeForBonus      time.Duration
}

// DefaultScoreOptions returns sensible defaults for risk scoring.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		PrivilegedThreshold: 1024,
		KnownRiskyPorts:     []uint16{21, 23, 135, 139, 445, 3389, 5900},
		MaxAgeForBonus:      30 * time.Minute,
	}
}

// ScorePorts computes a risk score for each port in the slice.
func ScorePorts(ports []Port, opts ScoreOptions) []PortScore {
	results := make([]PortScore, 0, len(ports))
	for _, p := range ports {
		results = append(results, scorePort(p, opts))
	}
	return results
}

func scorePort(p Port, opts ScoreOptions) PortScore {
	var score int
	var reasons []string

	if p.Port < opts.PrivilegedThreshold && p.Port > 0 {
		score += 10
		reasons = append(reasons, "privileged port (<1024)")
	}

	for _, risky := range opts.KnownRiskyPorts {
		if p.Port == risky {
			score += 30
			reasons = append(reasons, "known risky port")
			break
		}
	}

	if p.Process == nil || p.Process.Name == "" {
		score += 20
		reasons = append(reasons, "unknown process")
	}

	if !isLoopback(p.Address) {
		score += 15
		reasons = append(reasons, "externally reachable")
	}

	return PortScore{
		Port:    p,
		Score:   score,
		Level:   riskLevel(score),
		Reasons: reasons,
	}
}

func riskLevel(score int) RiskLevel {
	switch {
	case score >= 60:
		return RiskCritical
	case score >= 40:
		return RiskHigh
	case score >= 20:
		return RiskMedium
	default:
		return RiskLow
	}
}
