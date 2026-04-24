package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/ports"
)

var baselinePath string

var baselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Manage the port baseline snapshot",
}

var baselineSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current open ports as the baseline",
	RunE: func(cmd *cobra.Command, args []string) error {
		openPorts, err := ports.Scan()
		if err != nil {
			return fmt.Errorf("scan ports: %w", err)
		}
		if err := baseline.Save(baselinePath, openPorts); err != nil {
			return err
		}
		fmt.Printf("Baseline saved to %s (%d ports)\n", baselinePath, len(openPorts))
		return nil
	},
}

var baselineDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare current ports against the saved baseline",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := baseline.Load(baselinePath)
		if err != nil {
			return err
		}
		current, err := ports.Scan()
		if err != nil {
			return fmt.Errorf("scan ports: %w", err)
		}
		added, removed := baseline.Diff(b, current)
		if len(added) == 0 && len(removed) == 0 {
			fmt.Println("No changes detected.")
			return nil
		}
		for _, p := range added {
			fmt.Printf("[+] NEW      %s\n", p.Display())
		}
		for _, p := range removed {
			fmt.Printf("[-] REMOVED  %s\n", p.Display())
		}
		if len(added) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	baselineCmd.PersistentFlags().StringVar(&baselinePath, "file", "/var/lib/portwatch/baseline.json", "Path to baseline file")
	baselineCmd.AddCommand(baselineSaveCmd)
	baselineCmd.AddCommand(baselineDiffCmd)
	rootCmd.AddCommand(baselineCmd)
}
