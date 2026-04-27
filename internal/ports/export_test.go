package ports

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"
)

func sampleExportPorts() []Port {
	return []Port{
		{Protocol: "tcp", LocalAddr: "0.0.0.0", LocalPort: 80, State: "LISTEN", PID: 100, ProcessName: "nginx"},
		{Protocol: "udp", LocalAddr: "127.0.0.1", LocalPort: 53, State: "UNCONN", PID: 200, ProcessName: "dnsmasq"},
	}
}

func TestExportJSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportJSON(&buf, sampleExportPorts()); err != nil {
		t.Fatalf("ExportJSON error: %v", err)
	}
	var records []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
	if records[0].Port.ProcessName != "nginx" {
		t.Errorf("expected nginx, got %s", records[0].Port.ProcessName)
	}
}

func TestExportCSV_HeaderAndRows(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV(&buf, sampleExportPorts()); err != nil {
		t.Fatalf("ExportCSV error: %v", err)
	}
	r := csv.NewReader(&buf)
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("csv read error: %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("expected 3 rows (header+2), got %d", len(rows))
	}
	if rows[0][0] != "timestamp" {
		t.Errorf("expected header row to start with 'timestamp'")
	}
	if rows[1][2] != "0.0.0.0" {
		t.Errorf("expected local_addr 0.0.0.0, got %s", rows[1][2])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, sampleExportPorts(), ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExport_Dispatch_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, sampleExportPorts(), FormatJSON); err != nil {
		t.Fatalf("Export JSON error: %v", err)
	}
	if !strings.Contains(buf.String(), "nginx") {
		t.Error("expected JSON output to contain 'nginx'")
	}
}

func TestExport_Dispatch_CSV(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, sampleExportPorts(), FormatCSV); err != nil {
		t.Fatalf("Export CSV error: %v", err)
	}
	if !strings.Contains(buf.String(), "dnsmasq") {
		t.Error("expected CSV output to contain 'dnsmasq'")
	}
}
