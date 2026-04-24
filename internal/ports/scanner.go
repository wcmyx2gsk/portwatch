package ports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Listener represents an open port with its associated process info.
type Listener struct {
	Protocol string
	Address  string
	Port     uint16
	PID      int
	Process  string
}

// Scan reads active TCP and UDP listeners from /proc/net.
func Scan() ([]Listener, error) {
	var listeners []Listener

	for _, proto := range []string{"tcp", "tcp6", "udp", "udp6"} {
		path := fmt.Sprintf("/proc/net/%s", proto)
		entries, err := parseProcNet(path, proto)
		if err != nil {
			// Non-fatal: some systems may not have all files
			continue
		}
		listeners = append(listeners, entries...)
	}

	return listeners, nil
}

func parseProcNet(path, proto string) ([]Listener, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var listeners []Listener
	scanner := bufio.NewScanner(f)

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		// State field: 0A = TCP_LISTEN, 07 = UDP (stateless)
		state := fields[3]
		if strings.HasPrefix(proto, "tcp") && state != "0A" {
			continue
		}

		addr, port, err := parseHexAddr(fields[1])
		if err != nil {
			continue
		}

		listeners = append(listeners, Listener{
			Protocol: proto,
			Address:  addr,
			Port:     port,
		})
	}

	return listeners, scanner.Err()
}

func parseHexAddr(hexAddr string) (string, uint16, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address: %s", hexAddr)
	}

	portVal, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return "", 0, err
	}

	return parts[0], uint16(portVal), nil
}
