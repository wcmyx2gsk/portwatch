package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"portwatch/internal/ports"
)

var classifyCmd = &cobra.Command{
	Use:   "classify",
	Short: "Classify open ports by risk category",
	Long:  `Scans open ports and classifies each into a risk category such as web, database, remote-access, or system.`,
	RunE:  runClassify,
}

var classifyFilterCategory string

func runClassify(cmd *cobra.Command, args []string) error {
	scanned, err := ports.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	results := ports.ClassifyAll(scanned)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PORT\tPROTOCOL\tCATEGORY\tREASON")
	fmt.Fprintln(w, "----\t--------\t--------\t------")

	printed := 0
	for _, r := range results {
		if classifyFilterCategory != "" && string(r.Category) != classifyFilterCategory {
			continue
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			r.Port.Port,
			r.Port.Protocol,
			r.Category,
			r.Reason,
		)
		printed++
	}
	w.Flush()

	if printed == 0 {
		fmt.Println("No ports matched the filter.")
	}
	return nil
}

func init() {
	classifyCmd.Flags().StringVarP(
		&classifyFilterCategory, "category", "c", "",
		"Filter output to a specific category (web, database, remote-access, system, unknown, messaging)",
	)
	rootCmd.AddCommand(classifyCmd)
}
