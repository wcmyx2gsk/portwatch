package ports

import (
	"testing"
)

func groupPorts() []Port {
	return []Port{
		{Port: 80, Protocol: "tcp"},
		{Port: 8080, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
		{Port: 53, Protocol: "udp"},
		{Port: 9999, Protocol: "tcp"},
	}
}

func fakeResolve(port uint16, proto string) string {
	switch {
	case port == 80 && proto == "tcp":
		return "http"
	case port == 8080 && proto == "tcp":
		return "http"
	case port == 443 && proto == "tcp":
		return "https"
	case port == 53 && proto == "udp":
		return "dns"
	default:
		return ""
	}
}

func TestGroupByService_ReturnsCorrectGroupCount(t *testing.T) {
	groups := GroupByService(groupPorts(), fakeResolve)
	// expected groups: tcp/http, tcp/https, tcp/unknown, udp/dns
	if len(groups) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(groups))
	}
}

func TestGroupByService_HTTPGroupHasTwoPorts(t *testing.T) {
	groups := GroupByService(groupPorts(), fakeResolve)
	for _, g := range groups {
		if g.Protocol == "tcp" && g.Service == "http" {
			if len(g.Ports) != 2 {
				t.Fatalf("expected 2 ports in tcp/http group, got %d", len(g.Ports))
			}
			return
		}
	}
	t.Fatal("tcp/http group not found")
}

func TestGroupByService_UnknownGroup_CatchesUnresolved(t *testing.T) {
	groups := GroupByService(groupPorts(), fakeResolve)
	for _, g := range groups {
		if g.Service == "unknown" {
			if len(g.Ports) != 1 {
				t.Fatalf("expected 1 port in unknown group, got %d", len(g.Ports))
			}
			if g.Ports[0].Port != 9999 {
				t.Fatalf("expected port 9999 in unknown group, got %d", g.Ports[0].Port)
			}
			return
		}
	}
	t.Fatal("unknown group not found")
}

func TestGroupByService_SortedByProtocolThenService(t *testing.T) {
	groups := GroupByService(groupPorts(), fakeResolve)
	for i := 1; i < len(groups); i++ {
		prev, curr := groups[i-1], groups[i]
		if prev.Protocol > curr.Protocol {
			t.Fatalf("groups not sorted by protocol: %s > %s", prev.Protocol, curr.Protocol)
		}
		if prev.Protocol == curr.Protocol && prev.Service > curr.Service {
			t.Fatalf("groups not sorted by service within protocol: %s > %s", prev.Service, curr.Service)
		}
	}
}

func TestGroupByService_NilResolver_AllUnknown(t *testing.T) {
	groups := GroupByService(groupPorts(), nil)
	if len(groups) != 2 {
		// two protocols: tcp and udp, both unknown
		t.Fatalf("expected 2 groups (one per protocol), got %d", len(groups))
	}
	for _, g := range groups {
		if g.Service != "unknown" {
			t.Fatalf("expected service 'unknown', got %q", g.Service)
		}
	}
}
