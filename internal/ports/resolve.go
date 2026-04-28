package ports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ServiceName maps a port number and protocol to a well-known service name.
type ServiceName struct {
	Port     int
	Protocol string
	Name     string
}

// ResolveService attempts to resolve a port/protocol pair to a known service
// name by reading /etc/services. Falls back to a small built-in table if the
// file is unavailable.
func ResolveService(port int, protocol string) string {
	name, err := lookupEtcServices(port, protocol)
	if err == nil && name != "" {
		return name
	}
	return builtinLookup(port, protocol)
}

func lookupEtcServices(port int, protocol string) (string, error) {
	f, err := os.Open("/etc/services")
	if err != nil {
		return "", err
	}
	defer f.Close()

	proto := strings.ToLower(protocol)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		portProto := fields[1] // e.g. "80/tcp"
		parts := strings.SplitN(portProto, "/", 2)
		if len(parts) != 2 {
			continue
		}
		p, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		if p == port && strings.ToLower(parts[1]) == proto {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("not found")
}

// builtinLookup provides a minimal fallback table for common ports.
func builtinLookup(port int, protocol string) string {
	type key struct {
		port  int
		proto string
	}
	table := map[key]string{
		{22, "tcp"}:   "ssh",
		{25, "tcp"}:   "smtp",
		{53, "tcp"}:   "domain",
		{53, "udp"}:   "domain",
		{80, "tcp"}:   "http",
		{443, "tcp"}:  "https",
		{3306, "tcp"}: "mysql",
		{5432, "tcp"}: "postgresql",
		{6379, "tcp"}: "redis",
		{8080, "tcp"}: "http-alt",
	}
	if name, ok := table[key{port, strings.ToLower(protocol)}]; ok {
		return name
	}
	return ""
}
