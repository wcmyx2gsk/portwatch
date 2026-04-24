package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/portwatch/internal/ports"
)

var rootCmd = &cobra.Command{
	Use:   "portwatch",
	Short: "Monitor open ports and alert on unexpected listeners",
	RunE:  runScan,
}

func runScan(cmd *cobra.Command, args []string) error {
	openPorts, err := ports.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}
	if len(openPorts) == 0 {
		fmt.Println("No listening ports found.")
		return nil
	}
	fmt.Printf("%-6s %-21s %s\n", "PROTO", "ADDRESS", "STATE")
	fmt.Println("----------------------------------------------")
	for _, p := range openPorts {
		fmt.Println(p.Display())
	}
	return nil
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
