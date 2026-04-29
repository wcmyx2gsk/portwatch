package ports

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a human-readable label attached to a port entry.
type Tag struct {
	Name   string `yaml:"name" json:"name"`
	Reason string `yaml:"reason" json:"reason"`
}

// TagRule maps a port/protocol pair to a set of tags.
type TagRule struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Tags     []Tag  `yaml:"tags"`
}

// TagRegistry holds all known tag rules.
type TagRegistry struct {
	Rules []TagRule
}

// DefaultTagRegistry returns a registry pre-populated with well-known service tags.
func DefaultTagRegistry() *TagRegistry {
	return &TagRegistry{
		Rules: []TagRule{
			{Port: 22, Protocol: "tcp", Tags: []Tag{{Name: "ssh", Reason: "Secure Shell"}}},
			{Port: 80, Protocol: "tcp", Tags: []Tag{{Name: "http", Reason: "HTTP web server"}}},
			{Port: 443, Protocol: "tcp", Tags: []Tag{{Name: "https", Reason: "HTTPS web server"}}},
			{Port: 3306, Protocol: "tcp", Tags: []Tag{{Name: "mysql", Reason: "MySQL database"}}},
			{Port: 5432, Protocol: "tcp", Tags: []Tag{{Name: "postgres", Reason: "PostgreSQL database"}}},
			{Port: 6379, Protocol: "tcp", Tags: []Tag{{Name: "redis", Reason: "Redis cache"}}},
			{Port: 27017, Protocol: "tcp", Tags: []Tag{{Name: "mongodb", Reason: "MongoDB database"}}},
		},
	}
}

// Lookup returns all tags for the given port and protocol.
func (r *TagRegistry) Lookup(port int, protocol string) []Tag {
	proto := strings.ToLower(protocol)
	for _, rule := range r.Rules {
		if rule.Port == port && strings.ToLower(rule.Protocol) == proto {
			return rule.Tags
		}
	}
	return nil
}

// TagNames returns a comma-separated string of tag names for a port/protocol pair.
func (r *TagRegistry) TagNames(port int, protocol string) string {
	tags := r.Lookup(port, protocol)
	if len(tags) == 0 {
		return ""
	}
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

// TagSummary returns a formatted summary string for a port/protocol pair.
func (r *TagRegistry) TagSummary(port int, protocol string) string {
	tags := r.Lookup(port, protocol)
	if len(tags) == 0 {
		return fmt.Sprintf("%s/%d: (untagged)", strings.ToLower(protocol), port)
	}
	names := r.TagNames(port, protocol)
	return fmt.Sprintf("%s/%d: [%s]", strings.ToLower(protocol), port, names)
}
