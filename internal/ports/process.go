package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ProcessInfo holds metadata about the process owning a port.
type ProcessInfo struct {
	PID  int
	Name string
	Exe  string
}

// LookupProcess attempts to find the process that owns the given inode
// by scanning /proc/<pid>/fd and /proc/<pid>/net/tcp entries.
func LookupProcess(inode uint64) (*ProcessInfo, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read /proc: %w", err)
	}

	inodeStr := fmt.Sprintf("socket:[%d]", inode)

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}

		fdDir := fmt.Sprintf("/proc/%d/fd", pid)
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if link == inodeStr {
				return buildProcessInfo(pid)
			}
		}
	}
	return nil, fmt.Errorf("no process found for inode %d", inode)
}

func buildProcessInfo(pid int) (*ProcessInfo, error) {
	exePath := fmt.Sprintf("/proc/%d/exe", pid)
	exe, err := os.Readlink(exePath)
	if err != nil {
		exe = "unknown"
	}

	commPath := fmt.Sprintf("/proc/%d/comm", pid)
	nameBytes, err := os.ReadFile(commPath)
	name := "unknown"
	if err == nil {
		name = strings.TrimSpace(string(nameBytes))
	}

	return &ProcessInfo{
		PID:  pid,
		Name: name,
		Exe:  exe,
	}, nil
}
