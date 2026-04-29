package ports

import (
	"strings"
	"testing"
)

func TestDefaultTagRegistry_NotNil(t *testing.T) {
	r := DefaultTagRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
	if len(r.Rules) == 0 {
		t.Fatal("expected at least one default rule")
	}
}

func TestLookup_KnownPort(t *testing.T) {
	r := DefaultTagRegistry()
	tags := r.Lookup(22, "tcp")
	if len(tags) == 0 {
		t.Fatal("expected tags for port 22/tcp")
	}
	if tags[0].Name != "ssh" {
		t.Errorf("expected tag name 'ssh', got %q", tags[0].Name)
	}
}

func TestLookup_UnknownPort_ReturnsNil(t *testing.T) {
	r := DefaultTagRegistry()
	tags := r.Lookup(9999, "tcp")
	if tags != nil {
		t.Errorf("expected nil for unknown port, got %v", tags)
	}
}

func TestLookup_CaseInsensitiveProtocol(t *testing.T) {
	r := DefaultTagRegistry()
	tags := r.Lookup(80, "TCP")
	if len(tags) == 0 {
		t.Fatal("expected tags for port 80/TCP (case-insensitive)")
	}
}

func TestTagNames_KnownPort(t *testing.T) {
	r := DefaultTagRegistry()
	names := r.TagNames(443, "tcp")
	if names == "" {
		t.Fatal("expected non-empty tag names for 443/tcp")
	}
	if !strings.Contains(names, "https") {
		t.Errorf("expected 'https' in tag names, got %q", names)
	}
}

func TestTagNames_UnknownPort_ReturnsEmpty(t *testing.T) {
	r := DefaultTagRegistry()
	names := r.TagNames(12345, "udp")
	if names != "" {
		t.Errorf("expected empty string for unknown port, got %q", names)
	}
}

func TestTagSummary_KnownPort(t *testing.T) {
	r := DefaultTagRegistry()
	summary := r.TagSummary(6379, "tcp")
	if !strings.Contains(summary, "redis") {
		t.Errorf("expected 'redis' in summary, got %q", summary)
	}
	if !strings.Contains(summary, "tcp/6379") {
		t.Errorf("expected 'tcp/6379' in summary, got %q", summary)
	}
}

func TestTagSummary_UnknownPort_ShowsUntagged(t *testing.T) {
	r := DefaultTagRegistry()
	summary := r.TagSummary(9999, "tcp")
	if !strings.Contains(summary, "(untagged)") {
		t.Errorf("expected '(untagged)' in summary, got %q", summary)
	}
}

func TestTagRegistry_AddCustomRule(t *testing.T) {
	r := DefaultTagRegistry()
	r.Rules = append(r.Rules, TagRule{
		Port:     8080,
		Protocol: "tcp",
		Tags:     []Tag{{Name: "dev-server", Reason: "local dev HTTP"}},
	})
	names := r.TagNames(8080, "tcp")
	if names != "dev-server" {
		t.Errorf("expected 'dev-server', got %q", names)
	}
}
