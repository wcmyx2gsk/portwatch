package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"portwatch/internal/ports"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Show tag information for known ports",
	Long:  `Display well-known service tags for open ports from the current scan.`,
	RunE:  runTag,
}

func runTag(cmd *cobra.Command, args []string) error {
	scanned, err := ports.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	registry := ports.DefaultTagRegistry()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "PROTO\tPORT\tADDR\tTAGS\tREASON")
	fmt.Fprintln(w, "-----\t----\t----\t----\t------")

	for _, p := range scanned {
		tags := registry.Lookup(p.Port, p.Protocol)
		tagName := "(untagged)"
		reason := ""
		if len(tags) > 0 {
			tagName = tags[0].Name
			reason = tags[0].Reason
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\n",
			p.Protocol, p.Port, p.LocalAddr, tagName, reason)
	}

	return w.Flush()
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
