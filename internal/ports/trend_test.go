package ports

import (
	"testing"
	"time"
)

func makePorts(proto string, portNum uint16) Port {
	return Port{
		Protocol:     proto,
		LocalAddress: "0.0.0.0",
		LocalPort:    portNum,
		State:        "LISTEN",
	}
}

func TestBuildTrend_Empty(t *testing.T) {
	result := BuildTrend(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestBuildTrend_SingleEntry(t *testing.T) {
	now := time.Now()
	entries := []HistoryEntry{
		{Timestamp: now, Ports: []Port{makePorts("tcp", 80)}},
	}
	result := BuildTrend(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(result))
	}
	if result[0].Occurrences != 1 {
		t.Errorf("expected 1 occurrence, got %d", result[0].Occurrences)
	}
	if result[0].Missed != 0 {
		t.Errorf("expected 0 missed, got %d", result[0].Missed)
	}
}

func TestBuildTrend_MultipleEntries_Occurrences(t *testing.T) {
	t1 := time.Now().Add(-2 * time.Minute)
	t2 := time.Now().Add(-1 * time.Minute)
	t3 := time.Now()
	p80 := makePorts("tcp", 80)
	p443 := makePorts("tcp", 443)

	entries := []HistoryEntry{
		{Timestamp: t1, Ports: []Port{p80}},
		{Timestamp: t2, Ports: []Port{p80, p443}},
		{Timestamp: t3, Ports: []Port{p80}},
	}

	result := BuildTrend(entries)
	if len(result) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(result))
	}

	for _, s := range result {
		if s.Port.LocalPort == 80 {
			if s.Occurrences != 3 {
				t.Errorf("port 80: expected 3 occurrences, got %d", s.Occurrences)
			}
			if s.Missed != 0 {
				t.Errorf("port 80: expected 0 missed, got %d", s.Missed)
			}
		} else if s.Port.LocalPort == 443 {
			if s.Occurrences != 1 {
				t.Errorf("port 443: expected 1 occurrence, got %d", s.Occurrences)
			}
			if s.Missed != 2 {
				t.Errorf("port 443: expected 2 missed, got %d", s.Missed)
			}
		}
	}
}

func TestBuildTrend_SortedByProtocolAndPort(t *testing.T) {
	now := time.Now()
	entries := []HistoryEntry{
		{
			Timestamp: now,
			Ports: []Port{
				makePorts("udp", 53),
				makePorts("tcp", 443),
				makePorts("tcp", 80),
			},
		},
	}
	result := BuildTrend(entries)
	if len(result) != 3 {
		t.Fatalf("expected 3 summaries, got %d", len(result))
	}
	if result[0].Port.Protocol != "tcp" || result[0].Port.LocalPort != 80 {
		t.Errorf("expected first entry tcp:80, got %s:%d", result[0].Port.Protocol, result[0].Port.LocalPort)
	}
	if result[1].Port.Protocol != "tcp" || result[1].Port.LocalPort != 443 {
		t.Errorf("expected second entry tcp:443, got %s:%d", result[1].Port.Protocol, result[1].Port.LocalPort)
	}
	if result[2].Port.Protocol != "udp" || result[2].Port.LocalPort != 53 {
		t.Errorf("expected third entry udp:53, got %s:%d", result[2].Port.Protocol, result[2].Port.LocalPort)
	}
}
