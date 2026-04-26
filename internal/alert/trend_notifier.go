package alert

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// TrendNotifier writes trend summaries to a writer.
type TrendNotifier struct {
	w io.Writer
}

// NewTrendNotifier creates a TrendNotifier writing to the given writer.
// If w is nil, os.Stdout is used.
func NewTrendNotifier(w io.Writer) *TrendNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &TrendNotifier{w: w}
}

// Notify prints a formatted trend table to the underlying writer.
func (tn *TrendNotifier) Notify(summaries []ports.TrendSummary) {
	if len(summaries) == 0 {
		fmt.Fprintln(tn.w, "[trend] no trend data available")
		return
	}

	tw := tabwriter.NewWriter(tn.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROTO\tPORT\tADDRESS\tFIRST SEEN\tLAST SEEN\tOCCURRENCES\tMISSED")
	for _, s := range summaries {
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\t%d\t%d\n",
			s.Port.Protocol,
			s.Port.LocalPort,
			s.Port.LocalAddress,
			s.FirstSeen.Format(time.RFC3339),
			s.LastSeen.Format(time.RFC3339),
			s.Occurrences,
			s.Missed,
		)
	}
	tw.Flush()
}

// NotifyUnstable prints only ports that were not seen in every snapshot.
func (tn *TrendNotifier) NotifyUnstable(summaries []ports.TrendSummary) {
	unstable := make([]ports.TrendSummary, 0)
	for _, s := range summaries {
		if s.Missed > 0 {
			unstable = append(unstable, s)
		}
	}
	if len(unstable) == 0 {
		fmt.Fprintln(tn.w, "[trend] all ports are stable across snapshots")
		return
	}
	fmt.Fprintf(tn.w, "[trend] %d unstable port(s) detected:\n", len(unstable))
	tn.Notify(unstable)
}
