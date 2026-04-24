package ports

import "strings"

// FilterOptions controls which ports are included in scan results.
type FilterOptions struct {
	// Protocols limits results to the given protocols (e.g. ["tcp", "udp"]).
	// An empty slice means all protocols are included.
	Protocols []string

	// ExcludePorts is a set of port numbers to ignore.
	ExcludePorts map[uint16]struct{}

	// LocalOnly, when true, skips ports bound to non-loopback addresses.
	LocalOnly bool
}

// DefaultFilterOptions returns a FilterOptions with no restrictions.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{
		Protocols:    []string{},
		ExcludePorts: map[uint16]struct{}{},
	}
}

// Apply returns the subset of ports that match the filter options.
func (f FilterOptions) Apply(ports []Port) []Port {
	result := make([]Port, 0, len(ports))
	for _, p := range ports {
		if !f.matchProtocol(p) {
			continue
		}
		if f.isExcluded(p) {
			continue
		}
		if f.LocalOnly && !isLoopback(p.Address) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func (f FilterOptions) matchProtocol(p Port) bool {
	if len(f.Protocols) == 0 {
		return true
	}
	for _, proto := range f.Protocols {
		if strings.EqualFold(proto, p.Protocol) {
			return true
		}
	}
	return false
}

func (f FilterOptions) isExcluded(p Port) bool {
	_, excluded := f.ExcludePorts[p.Port]
	return excluded
}

func isLoopback(addr string) bool {
	return strings.HasPrefix(addr, "127.") || addr == "::1" || addr == "0.0.0.0"
}
