package cmd

import (
	"fmt"
	"os"

	"github.com/portwatch/portwatch/internal/ports"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "portwatch",
	Short: "Monitor open ports and alert on unexpected listeners",
	Long: `portwatch is a lightweight CLI daemon that scans active
network listeners and reports unexpected open ports based on
a configurable allowlist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runScan()
	},
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform a one-shot scan of open ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runScan()
	},
}

func runScan() error {
	listeners, err := ports.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if len(listeners) == 0 {
		fmt.Println("No active listeners found.")
		return nil
	}

	fmt.Printf("%-8s %-20s %s\n", "PROTO", "ADDRESS", "PORT")
	fmt.Println("---------------------------------------")
	for _, l := range listeners {
		fmt.Printf("%-8s %-20s %d\n", l.Protocol, l.Address, l.Port)
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
