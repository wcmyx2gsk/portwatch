package cmd

import (
	"fmt"
	"os"

	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/ports"
	"github.com/spf13/cobra"
)

var (
	exportFormat  string
	exportOutFile string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export current open ports to JSON or CSV",
	RunE:  runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	scanned, err := ports.Scan(cfg.ProcNetTCPPath, cfg.ProcNetUDPPath)
	if err != nil {
		return fmt.Errorf("scan ports: %w", err)
	}

	w := os.Stdout
	if exportOutFile != "" {
		f, ferr := os.Create(exportOutFile)
		if ferr != nil {
			return fmt.Errorf("create output file: %w", ferr)
		}
		defer f.Close()
		w = f
	}

	fmt2 := ports.ExportFormat(exportFormat)
	if err := ports.Export(w, scanned, fmt2); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	if exportOutFile != "" {
		_, _ = fmt.Fprintf(os.Stderr, "exported %d ports to %s\n", len(scanned), exportOutFile)
	}
	return nil
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "output format: json or csv")
	exportCmd.Flags().StringVarP(&exportOutFile, "out", "o", "", "output file path (default: stdout)")
	rootCmd.AddCommand(exportCmd)
}
