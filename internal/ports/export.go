package ports

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// ExportFormat defines the output format for port data exports.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// ExportRecord wraps a Port with an export timestamp.
type ExportRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Port      Port      `json:"port"`
}

// ExportJSON writes ports as a JSON array to w.
func ExportJSON(w io.Writer, ports []Port) error {
	records := make([]ExportRecord, len(ports))
	now := time.Now().UTC()
	for i, p := range ports {
		records[i] = ExportRecord{Timestamp: now, Port: p}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

// ExportCSV writes ports as CSV rows to w.
func ExportCSV(w io.Writer, ports []Port) error {
	csvW := csv.NewWriter(w)
	if err := csvW.Write([]string{"timestamp", "protocol", "local_addr", "local_port", "state", "pid", "process"}); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	for _, p := range ports {
		row := []string{
			now,
			p.Protocol,
			p.LocalAddr,
			strconv.Itoa(int(p.LocalPort)),
			p.State,
			strconv.Itoa(p.PID),
			p.ProcessName,
		}
		if err := csvW.Write(row); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}
	csvW.Flush()
	return csvW.Error()
}

// Export dispatches to ExportJSON or ExportCSV based on format.
func Export(w io.Writer, ports []Port, format ExportFormat) error {
	switch format {
	case FormatJSON:
		return ExportJSON(w, ports)
	case FormatCSV:
		return ExportCSV(w, ports)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}
