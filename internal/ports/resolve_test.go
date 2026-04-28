package ports

import (
	"testing"
)

func TestBuiltinLookup_KnownPorts(t *testing.T) {
	cases := []struct {
		port     int
		protocol string
		want     string
	}{
		{22, "tcp", "ssh"},
		{80, "tcp", "http"},
		{443, "tcp", "https"},
		{53, "udp", "domain"},
		{3306, "tcp", "mysql"},
		{6379, "tcp", "redis"},
	}
	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			got := builtinLookup(tc.port, tc.protocol)
			if got != tc.want {
				t.Errorf("builtinLookup(%d, %q) = %q; want %q", tc.port, tc.protocol, got, tc.want)
			}
		})
	}
}

func TestBuiltinLookup_UnknownPort(t *testing.T) {
	got := builtinLookup(9999, "tcp")
	if got != "" {
		t.Errorf("expected empty string for unknown port, got %q", got)
	}
}

func TestBuiltinLookup_WrongProtocol(t *testing.T) {
	// port 80 is only registered for tcp in builtin table
	got := builtinLookup(80, "udp")
	if got != "" {
		t.Errorf("expected empty for 80/udp, got %q", got)
	}
}

func TestResolveService_FallsBackToBuiltin(t *testing.T) {
	// Even if /etc/services exists, common ports should resolve.
	got := ResolveService(22, "tcp")
	if got == "" {
		t.Error("expected non-empty service name for port 22/tcp")
	}
}

func TestResolveService_Unknown_ReturnsEmpty(t *testing.T) {
	got := ResolveService(19999, "tcp")
	if got != "" {
		t.Errorf("expected empty for obscure port, got %q", got)
	}
}

func TestResolveService_CaseInsensitiveProtocol(t *testing.T) {
	got1 := ResolveService(443, "TCP")
	got2 := ResolveService(443, "tcp")
	if got1 != got2 {
		t.Errorf("case should not matter: %q vs %q", got1, got2)
	}
}
