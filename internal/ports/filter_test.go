package ports

import (
	"testing"
)

func samplePorts() []Port {
	return []Port{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 22},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 80},
		{Protocol: "udp", Address: "0.0.0.0", Port: 53},
		{Protocol: "tcp", Address: "192.168.1.5", Port: 8080},
		{Protocol: "udp", Address: "127.0.0.1", Port: 5353},
	}
}

func TestApply_NoFilter_ReturnsAll(t *testing.T) {
	f := DefaultFilterOptions()
	result := f.Apply(samplePorts())
	if len(result) != len(samplePorts()) {
		t.Errorf("expected %d ports, got %d", len(samplePorts()), len(result))
	}
}

func TestApply_ProtocolFilter_TCPOnly(t *testing.T) {
	f := DefaultFilterOptions()
	f.Protocols = []string{"tcp"}
	result := f.Apply(samplePorts())
	for _, p := range result {
		if p.Protocol != "tcp" {
			t.Errorf("expected only tcp, got %s", p.Protocol)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 tcp ports, got %d", len(result))
	}
}

func TestApply_ExcludePort(t *testing.T) {
	f := DefaultFilterOptions()
	f.ExcludePorts = map[uint16]struct{}{22: {}, 80: {}}
	result := f.Apply(samplePorts())
	for _, p := range result {
		if p.Port == 22 || p.Port == 80 {
			t.Errorf("excluded port %d appeared in results", p.Port)
		}
	}
}

func TestApply_LocalOnly(t *testing.T) {
	f := DefaultFilterOptions()
	f.LocalOnly = true
	result := f.Apply(samplePorts())
	for _, p := range result {
		if p.Address == "192.168.1.5" {
			t.Errorf("non-loopback address %s should have been filtered", p.Address)
		}
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	f := DefaultFilterOptions()
	f.Protocols = []string{"tcp"}
	f.LocalOnly = true
	f.ExcludePorts = map[uint16]struct{}{22: {}}
	result := f.Apply(samplePorts())
	// tcp + loopback/0.0.0.0 only: ports 22 and 80 qualify, but 22 is excluded => only 80
	if len(result) != 1 {
		t.Errorf("expected 1 port, got %d", len(result))
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80, got %d", result[0].Port)
	}
}
