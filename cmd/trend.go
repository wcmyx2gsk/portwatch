package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ports"
)

var (
	trendHistoryFile string
	trendUnstableOnly bool
)

var trendCmd = &cobra.Command{
	Use:   "trend",
	Short: "Analyse port trend data from history",
	Long:  `Reads the port scan history file and prints a trend summary showing how often each port was observed and whether it was intermittently missing.`,
	RunE:  runTrend,
}

func runTrend(cmd *cobra.Command, args []string) error {
	entries, err := ports.LoadHistory(trendHistoryFile)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "[trend] history file is empty — run 'portwatch watch' to collect data")
		return nil
	}

	summaries := ports.BuildTrend(entries)
	tn := alert.NewTrendNotifier(os.Stdout)

	if trendUnstableOnly {
		tn.NotifyUnstable(summaries)
	} else {
		tn.Notify(summaries)
	}

	return nil
}

func init() {
	trendCmd.Flags().StringVarP(
		&trendHistoryFile, "history", "H", "portwatch_history.yaml",
		"Path to the history file",
	)
	trendCmd.Flags().BoolVarP(
		&trendUnstableOnly, "unstable", "u", false,
		"Show only ports that were not seen in every snapshot",
	)
	rootCmd.AddCommand(trendCmd)
}
