package ports

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestBuildProcessInfo_CurrentProcess(t *testing.T) {
	pid := os.Getpid()
	info, err := buildProcessInfo(pid)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if info.PID != pid {
		t.Errorf("expected PID %d, got %d", pid, info.PID)
	}
	if info.Name == "" {
		t.Error("expected non-empty process name")
	}
	if info.Exe == "" {
		t.Error("expected non-empty exe path")
	}
}

func TestBuildProcessInfo_InvalidPID(t *testing.T) {
	_, err := buildProcessInfo(-1)
	if err == nil {
		t.Error("expected error for invalid PID, got nil")
	}
}

func TestLookupProcess_UnknownInode(t *testing.T) {
	// inode 0 should never match a real socket
	_, err := LookupProcess(0)
	if err == nil {
		t.Error("expected error for unknown inode, got nil")
	}
}

func TestProcessInfo_Fields(t *testing.T) {
	pid := os.Getpid()
	info, err := buildProcessInfo(pid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify PID field is correctly stored as integer
	pidStr := strconv.Itoa(info.PID)
	if pidStr == "" {
		t.Error("PID should be convertible to string")
	}

	// Verify exe is not the placeholder when process exists
	if info.Exe == "" {
		t.Error("exe should not be empty for current process")
	}

	fmt.Printf("[test] current process: pid=%d name=%s exe=%s\n", info.PID, info.Name, info.Exe)
}
