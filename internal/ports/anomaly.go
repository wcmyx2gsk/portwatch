package ports

import "fmt"

// AnomalyKind classifies the type of anomaly detected.
type AnomalyKind string

const (
	AnomalyNewPort      AnomalyKind = "new_port"
	AnomalyRarePort     AnomalyKind = "rare_port"
	AnomalyHighPort     AnomalyKind = "high_port"
	AnomalyUnknownProc  AnomalyKind = "unknown_process"
)

// Anomaly represents a suspicious port observation.
type Anomaly struct {
	Kind    AnomalyKind
	Port    Port
	Message string
}

func (a Anomaly) String() string {
	return fmt.Sprintf("[%s] %s:%d (%s) — %s",
		a.Kind, a.Port.Protocol, a.Port.LocalPort, a.Port.State, a.Message)
}

// AnomalyOptions controls which anomaly checks are enabled.
type AnomalyOptions struct {
	HighPortThreshold uint16
	MinOccurrences    int
	FlagUnknownProcs  bool
}

// DefaultAnomalyOptions returns sensible defaults.
func DefaultAnomalyOptions() AnomalyOptions {
	return AnomalyOptions{
		HighPortThreshold: 1024,
		MinOccurrences:    2,
		FlagUnknownProcs:  true,
	}
}

// DetectAnomalies inspects current ports against trend history and returns anomalies.
func DetectAnomalies(current []Port, trend []TrendSummary, opts AnomalyOptions) []Anomaly {
	seen := make(map[string]TrendSummary)
	for _, t := range trend {
		key := fmt.Sprintf("%s:%d", t.Protocol, t.LocalPort)
		seen[key] = t
	}

	var anomalies []Anomaly
	for _, p := range current {
		key := fmt.Sprintf("%s:%d", p.Protocol, p.LocalPort)
		t, exists := seen[key]

		if !exists {
			anomalies = append(anomalies, Anomaly{
				Kind:    AnomalyNewPort,
				Port:    p,
				Message: "port not seen in historical trend",
			})
			continue
		}

		if t.Occurrences < opts.MinOccurrences {
			anomalies = append(anomalies, Anomaly{
				Kind:    AnomalyRarePort,
				Port:    p,
				Message: fmt.Sprintf("only seen %d time(s) in history", t.Occurrences),
			})
		}

		if p.LocalPort < opts.HighPortThreshold && opts.FlagUnknownProcs && p.PID == 0 {
			anomalies = append(anomalies, Anomaly{
				Kind:    AnomalyUnknownProc,
				Port:    p,
				Message: "listening on privileged port with no identifiable process",
			})
		}
	}
	return anomalies
}
