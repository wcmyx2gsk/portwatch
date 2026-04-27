package alert

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/portwatch/internal/ports"
)

func sampleExportNotifierPorts() []ports.Port {
	return []ports.Port{
		{Protocol: "tcp", LocalAddr: "0.0.0.0", LocalPort: 443, State: "LISTEN", PID: 42, ProcessName: "caddy"},
		{Protocol: "udp", LocalAddr: "127.0.0.1", LocalPort: 5353, State: "UNCONN", PID: 99, ProcessName: "avahi"},
	}
}

func TestExportNotifier_Empty(t *testing.T) {
	var buf bytes.Buffer
	n := NewExportNotifier(&buf, ports.FormatJSON)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no open ports") {
		t.Errorf("expected empty notice, got: %s", buf.String())
	}
}

func TestExportNotifier_JSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	n := NewExportNotifier(&buf, ports.FormatJSON)
	if err := n.Notify(sampleExportNotifierPorts()); err != nil {
		t.Fatalf("Notify error: %v", err)
	}
	var records []ports.ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}

func TestExportNotifier_CSV_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	n := NewExportNotifier(&buf, ports.FormatCSV)
	if err := n.Notify(sampleExportNotifierPorts()); err != nil {
		t.Fatalf("Notify error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "portwatch export") {
		t.Errorf("expected export comment header, got: %s", out)
	}
	if !strings.Contains(out, "caddy") {
		t.Errorf("expected port data in CSV, got: %s", out)
	}
}

func TestExportNotifier_CSV_RowCount(t *testing.T) {
	var buf bytes.Buffer
	n := NewExportNotifier(&buf, ports.FormatCSV)
	if err := n.Notify(sampleExportNotifierPorts()); err != nil {
		t.Fatalf("Notify error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// comment line + header + 2 data rows = 4
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d: %v", len(lines), lines)
	}
}
