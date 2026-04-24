package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/watch"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Continuously monitor open ports and alert on changes",
	Long: `Starts a daemon that scans open ports at a regular interval.
Any port that appears or disappears relative to the saved baseline
will be reported to stdout (or the configured log file).`,
	RunE: runWatch,
}

func runWatch(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")

	var cfg *config.Config
	var err error
	if cfgPath != "" {
		cfg, err = config.Load(cfgPath)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
	} else {
		cfg = config.Default()
	}

	notifier := alert.NewNotifier(os.Stdout)
	watcher := watch.New(cfg, notifier)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		watcher.Stop()
	}()

	watcher.Start()
	return nil
}

func init() {
	watchCmd.Flags().StringP("config", "c", "", "path to config file (optional)")
	rootCmd.AddCommand(watchCmd)
}
