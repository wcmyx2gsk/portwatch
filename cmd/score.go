package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"portwatch/internal/ports"
)

var scoreMinLevel string

func runScore(cmd *cobra.Command, args []string) {
	scanned, err := ports.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		os.Exit(1)
	}

	opts := ports.DefaultScoreOptions()
	scores := ports.ScorePorts(scanned, opts)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PROTO\tPORT\tADDRESS\tSCORE\tLEVEL\tREASONS")

	filterLevel := ports.RiskLevel(scoreMinLevel)
	printed := 0
	for _, s := range scores {
		if filterLevel != "" && !meetsMinLevel(s.Level, filterLevel) {
			continue
		}
		proc := ""
		if s.Port.Process != nil {
			proc = s.Port.Process.Name
		}
		_ = proc
		fmt.Fprintf(w, "%s\t%d\t%s\t%d\t%s\t%v\n",
			s.Port.Protocol,
			s.Port.Port,
			s.Port.Address,
			s.Score,
			s.Level,
			s.Reasons,
		)
		printed++
	}
	w.Flush()

	if printed == 0 {
		fmt.Println("no ports matched the specified risk level filter")
	}
}

func meetsMinLevel(level, min ports.RiskLevel) bool {
	order := map[ports.RiskLevel]int{
		ports.RiskLow:      0,
		ports.RiskMedium:   1,
		ports.RiskHigh:     2,
		ports.RiskCritical: 3,
	}
	return order[level] >= order[min]
}

func init() {
	scoreCmd := &cobra.Command{
		Use:   "score",
		Short: "Score open ports by risk level",
		Long:  "Scans open ports and assigns a risk score based on heuristics such as privileged ports, known risky services, unknown processes, and external reachability.",
		Run:   runScore,
	}
	scoreCmd.Flags().StringVar(&scoreMinLevel, "min-level", "", "minimum risk level to display (low, medium, high, critical)")
	rootCmd.AddCommand(scoreCmd)
}
