package ports

import (
	"os"
	"testing"
)

func TestParseHexAddr_Valid(t *testing.T) {
	addr, port, err := parseHexAddr("0100007F:0050")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port != 80 {
		t.Errorf("expected port 80, got %d", port)
	}
	if addr != "0100007F" {
		t.Errorf("unexpected address: %s", addr)
	}
}

func TestParseHexAddr_Invalid(t *testing.T) {
	_, _, err := parseHexAddr("invalid")
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestParseProcNet_MissingFile(t *testing.T) {
	_, err := parseProcNet("/nonexistent/path", "tcp")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseProcNet_ValidFile(t *testing.T) {
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0051 00000000:0000 02 00000000:00000000 00:00000000 00000000     0        0 12346 1 0000000000000000 100 0 0 10 0
`
	tmpFile, err := os.CreateTemp("", "proc_net_tcp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	listeners, err := parseProcNet(tmpFile.Name(), "tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only state 0A (LISTEN) should be included
	if len(listeners) != 1 {
		t.Errorf("expected 1 listener, got %d", len(listeners))
	}
	if listeners[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", listeners[0].Port)
	}
}
