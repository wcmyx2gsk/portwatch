package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/config"
	"portwatch/internal/ports"
)

var (
	flagProtocols   []string
	flagExcludePorts []uint16
	flagLocalOnly   bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan open ports and print results",
	Long:  "Performs a one-shot scan of open ports and prints them to stdout, respecting filter flags.",
	RunE:  runFilteredScan,
}

func runFilteredScan(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load config: %v\n", err)
		cfg = config.Default()
	}

	scanned, err := ports.Scan(cfg.ProcNetPath)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	f := ports.DefaultFilterOptions()
	if len(flagProtocols) > 0 {
		f.Protocols = flagProtocols
	}
	if flagLocalOnly {
		f.LocalOnly = true
	}
	if len(flagExcludePorts) > 0 {
		for _, p := range flagExcludePorts {
			f.ExcludePorts[p] = struct{}{}
		}
	}

	result := f.Apply(scanned)

	if len(result) == 0 {
		fmt.Println("No open ports found matching the given filters.")
		return nil
	}

	fmt.Printf("%-8s %-20s %s\n", "PROTO", "ADDRESS", "PORT")
	for _, p := range result {
		fmt.Printf("%-8s %-20s %d\n", p.Protocol, p.Address, p.Port)
	}
	return nil
}

func init() {
	scanCmd.Flags().StringSliceVar(&flagProtocols, "protocol", []string{}, "filter by protocol (tcp, udp)")
	scanCmd.Flags().Uint16SliceVar(&flagExcludePorts, "exclude-port", []uint16{}, "ports to exclude from output")
	scanCmd.Flags().BoolVar(&flagLocalOnly, "local-only", false, "only show loopback/wildcard bound ports")
	rootCmd.AddCommand(scanCmd)
}
