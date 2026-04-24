package ports

import "fmt"

// Port represents a single listening network port entry.
type Port struct {
	Proto     string `json:"proto"`
	LocalAddr string `json:"local_addr"`
	LocalPort uint16 `json:"local_port"`
	State     string `json:"state"`
	PID       int    `json:"pid,omitempty"`
}

// String returns a canonical string representation used for comparison.
func (p Port) String() string {
	return fmt.Sprintf("%s:%s:%d", p.Proto, p.LocalAddr, p.LocalPort)
}

// Display returns a human-readable description of the port.
func (p Port) Display() string {
	if p.PID > 0 {
		return fmt.Sprintf("%-6s %s:%-5d  %-10s pid=%d", p.Proto, p.LocalAddr, p.LocalPort, p.State, p.PID)
	}
	return fmt.Sprintf("%-6s %s:%-5d  %s", p.Proto, p.LocalAddr, p.LocalPort, p.State)
}
