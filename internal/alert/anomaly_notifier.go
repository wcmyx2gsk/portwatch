package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// AnomalyNotifier writes anomaly reports to a writer.
type AnomalyNotifier struct {
	w io.Writer
}

// NewAnomalyNotifier creates an AnomalyNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewAnomalyNotifier(w io.Writer) *AnomalyNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &AnomalyNotifier{w: w}
}

// Notify writes a formatted anomaly report. Returns early with a message if
// no anomalies are present.
func (n *AnomalyNotifier) Notify(anomalies []ports.Anomaly) {
	ts := time.Now().Format(time.RFC3339)
	if len(anomalies) == 0 {
		fmt.Fprintf(n.w, "%s [anomaly] no anomalies detected\n", ts)
		return
	}

	fmt.Fprintf(n.w, "%s [anomaly] %d anomaly(s) detected\n", ts, len(anomalies))
	for _, a := range anomalies {
		fmt.Fprintf(n.w, "  %s\n", a.String())
	}
}

// NotifyIfCritical writes the report only when anomalies exist.
func (n *AnomalyNotifier) NotifyIfCritical(anomalies []ports.Anomaly) bool {
	if len(anomalies) == 0 {
		return false
	}
	n.Notify(anomalies)
	return true
}
