package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"portwatch/internal/ports"
)

var snapshotFile string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Capture or compare port snapshots",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a snapshot of current open ports to a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanned, err := ports.Scan()
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		snap := ports.NewSnapshot(scanned)
		if err := ports.SaveSnapshot(snapshotFile, snap); err != nil {
			return fmt.Errorf("save snapshot: %w", err)
		}
		fmt.Printf("Snapshot saved to %s (%d ports at %s)\n",
			snapshotFile, len(snap.Ports), snap.Timestamp.Format(time.RFC3339))
		return nil
	},
}

var snapshotDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare a saved snapshot against current open ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		old, err := ports.LoadSnapshot(snapshotFile)
		if err != nil {
			return fmt.Errorf("load snapshot: %w", err)
		}
		scanned, err := ports.Scan()
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		newSnap := ports.NewSnapshot(scanned)
		added, removed := ports.Compare(old, newSnap)

		if len(added) == 0 && len(removed) == 0 {
			fmt.Println("No changes since snapshot.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, p := range added {
			fmt.Fprintf(w, "+ ADDED\t%s\t%s:%d\n", p.Protocol, p.LocalAddr, p.LocalPort)
		}
		for _, p := range removed {
			fmt.Fprintf(w, "- REMOVED\t%s\t%s:%d\n", p.Protocol, p.LocalAddr, p.LocalPort)
		}
		w.Flush()
		return nil
	},
}

func init() {
	snapshotSaveCmd.Flags().StringVarP(&snapshotFile, "output", "o", "portwatch-snapshot.json", "Path to snapshot file")
	snapshotDiffCmd.Flags().StringVarP(&snapshotFile, "file", "f", "portwatch-snapshot.json", "Path to snapshot file")
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotDiffCmd)
	rootCmd.AddCommand(snapshotCmd)
}
