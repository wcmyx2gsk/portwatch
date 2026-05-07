package ports

import "sort"

// GroupKey identifies a logical group of ports.
type GroupKey struct {
	Protocol string
	Service  string
}

// PortGroup is a collection of ports sharing the same protocol and service name.
type PortGroup struct {
	Protocol string
	Service  string
	Ports    []Port
}

// GroupByService organises a slice of Port entries into groups keyed by
// (Protocol, Service). Ports with no resolved service name are placed under
// the label "unknown". Groups are returned sorted by Protocol then Service.
func GroupByService(ports []Port, resolve func(port uint16, proto string) string) []PortGroup {
	if resolve == nil {
		resolve = func(port uint16, proto string) string { return "" }
	}

	index := make(map[GroupKey]*PortGroup)

	for _, p := range ports {
		svc := resolve(p.Port, p.Protocol)
		if svc == "" {
			svc = "unknown"
		}
		key := GroupKey{Protocol: p.Protocol, Service: svc}
		if _, ok := index[key]; !ok {
			index[key] = &PortGroup{
				Protocol: p.Protocol,
				Service:  svc,
			}
		}
		index[key].Ports = append(index[key].Ports, p)
	}

	groups := make([]PortGroup, 0, len(index))
	for _, g := range index {
		groups = append(groups, *g)
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Protocol != groups[j].Protocol {
			return groups[i].Protocol < groups[j].Protocol
		}
		return groups[i].Service < groups[j].Service
	})

	return groups
}
