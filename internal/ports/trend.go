package ports

import (
	"sort"
	"time"
)

// TrendEntry represents a single observation of a port at a point in time.
type TrendEntry struct {
	Timestamp time.Time `yaml:"timestamp"`
	Port      Port      `yaml:"port"`
	Seen      bool      `yaml:"seen"`
}

// TrendSummary aggregates trend data for a single port.
type TrendSummary struct {
	Port        Port
	FirstSeen   time.Time
	LastSeen    time.Time
	Occurrences int
	Missed      int
}

// BuildTrend analyses a history and returns per-port trend summaries.
func BuildTrend(entries []HistoryEntry) []TrendSummary {
	type key struct {
		Proto   string
		Address string
		PortNum uint16
	}

	type stats struct {
		port      Port
		firstSeen time.Time
		lastSeen  time.Time
		count     int
		total     int
	}

	tracker := make(map[key]*stats)

	for _, entry := range entries {
		for _, p := range entry.Ports {
			k := key{Proto: p.Protocol, Address: p.LocalAddress, PortNum: p.LocalPort}
			if _, ok := tracker[k]; !ok {
				tracker[k] = &stats{
					port:      p,
					firstSeen: entry.Timestamp,
					lastSeen:  entry.Timestamp,
				}
			}
			s := tracker[k]
			s.count++
			if entry.Timestamp.After(s.lastSeen) {
				s.lastSeen = entry.Timestamp
			}
			if entry.Timestamp.Before(s.firstSeen) {
				s.firstSeen = entry.Timestamp
			}
		}
	}

	totalSnapshots := len(entries)
	result := make([]TrendSummary, 0, len(tracker))
	for _, s := range tracker {
		missed := totalSnapshots - s.count
		if missed < 0 {
			missed = 0
		}
		result = append(result, TrendSummary{
			Port:        s.port,
			FirstSeen:   s.firstSeen,
			LastSeen:    s.lastSeen,
			Occurrences: s.count,
			Missed:      missed,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Port.Protocol != result[j].Port.Protocol {
			return result[i].Port.Protocol < result[j].Port.Protocol
		}
		return result[i].Port.LocalPort < result[j].Port.LocalPort
	})

	return result
}
