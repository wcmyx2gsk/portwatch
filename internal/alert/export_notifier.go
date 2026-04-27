package alert

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/portwatch/internal/ports"
)

// ExportNotifier writes a formatted export summary to a writer.
type ExportNotifier struct {
	w      io.Writer
	format ports.ExportFormat
}

// NewExportNotifier creates an ExportNotifier writing to w in the given format.
func NewExportNotifier(w io.Writer, format ports.ExportFormat) *ExportNotifier {
	return &ExportNotifier{w: w, format: format}
}

// Notify exports the given ports with a header comment (JSON/CSV aware).
func (n *ExportNotifier) Notify(ps []ports.Port) error {
	if len(ps) == 0 {
		_, err := fmt.Fprintln(n.w, "# no open ports to export")
		return err
	}
	if n.format == ports.FormatCSV {
		_, err := fmt.Fprintf(n.w, "# portwatch export — %s — %d ports\n",
			time.Now().UTC().Format(time.RFC3339), len(ps))
		if err != nil {
			return err
		}
	}
	return ports.Export(n.w, ps, n.format)
}
