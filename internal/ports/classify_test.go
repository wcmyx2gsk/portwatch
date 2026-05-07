package ports

import (
	"testing"
)

func classifyPort(port int, proto string) PortInfo {
	return PortInfo{Port: port, Protocol: proto}
}

func TestClassifyPort_Web(t *testing.T) {
	r := ClassifyPort(classifyPort(80, "tcp"))
	if r.Category != CategoryWeb {
		t.Errorf("expected web, got %s", r.Category)
	}
}

func TestClassifyPort_Database(t *testing.T) {
	for _, port := range []int{3306, 5432, 27017} {
		r := ClassifyPort(classifyPort(port, "tcp"))
		if r.Category != CategoryDatabase {
			t.Errorf("port %d: expected database, got %s", port, r.Category)
		}
	}
}

func TestClassifyPort_RemoteAccess(t *testing.T) {
	r := ClassifyPort(classifyPort(22, "tcp"))
	if r.Category != CategoryRemote {
		t.Errorf("expected remote-access, got %s", r.Category)
	}
}

func TestClassifyPort_System(t *testing.T) {
	r := ClassifyPort(classifyPort(53, "udp"))
	if r.Category != CategorySystem {
		t.Errorf("expected system, got %s", r.Category)
	}
}

func TestClassifyPort_Unknown(t *testing.T) {
	r := ClassifyPort(classifyPort(9999, "tcp"))
	if r.Category != CategoryUnknown {
		t.Errorf("expected unknown, got %s", r.Category)
	}
	if r.Reason == "" {
		t.Error("expected non-empty reason for unknown port")
	}
}

func TestClassifyPort_ProtocolMismatch(t *testing.T) {
	// Port 22 rule is tcp only; udp should not match remote-access
	r := ClassifyPort(classifyPort(22, "udp"))
	if r.Category == CategoryRemote {
		t.Error("udp port 22 should not classify as remote-access")
	}
}

func TestClassifyAll_ReturnsCorrectCount(t *testing.T) {
	ports := []PortInfo{
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
		{Port: 9999, Protocol: "tcp"},
	}
	results := ClassifyAll(ports)
	if len(results) != len(ports) {
		t.Errorf("expected %d results, got %d", len(ports), len(results))
	}
}

func TestClassifyAll_Empty(t *testing.T) {
	results := ClassifyAll([]PortInfo{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
