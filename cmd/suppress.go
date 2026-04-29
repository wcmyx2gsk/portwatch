package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"portwatch/internal/ports"
)

var suppressFile string

var suppressCmd = &cobra.Command{
	Use:   "suppress",
	Short: "Manage port alert suppression rules",
}

var suppressAddCmd = &cobra.Command{
	Use:   "add <port> <protocol>",
	Short: "Add a suppression rule for a port/protocol",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var port int
		if _, err := fmt.Sscanf(args[0], "%d", &port); err != nil {
			return fmt.Errorf("invalid port: %s", args[0])
		}
		protocol := args[1]

		reason, _ := cmd.Flags().GetString("reason")
		ttl, _ := cmd.Flags().GetDuration("ttl")

		sl, err := ports.LoadSuppressList(suppressFile)
		if err != nil {
			return fmt.Errorf("loading suppress list: %w", err)
		}

		rule := ports.SuppressRule{
			Port:     port,
			Protocol: protocol,
			Reason:   reason,
		}
		if ttl > 0 {
			rule.Expires = time.Now().Add(ttl)
		}
		sl.Rules = append(sl.Rules, rule)

		if err := ports.SaveSuppressList(suppressFile, sl); err != nil {
			return fmt.Errorf("saving suppress list: %w", err)
		}
		fmt.Fprintf(os.Stdout, "suppression rule added: %d/%s\n", port, protocol)
		return nil
	},
}

var suppressListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active suppression rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		sl, err := ports.LoadSuppressList(suppressFile)
		if err != nil {
			return fmt.Errorf("loading suppress list: %w", err)
		}
		if len(sl.Rules) == 0 {
			fmt.Println("no suppression rules defined")
			return nil
		}
		fmt.Printf("%-8s %-10s %-30s %s\n", "PORT", "PROTOCOL", "REASON", "EXPIRES")
		for _, r := range sl.Rules {
			expires := "never"
			if !r.Expires.IsZero() {
				expires = r.Expires.Format(time.RFC3339)
			}
			fmt.Printf("%-8d %-10s %-30s %s\n", r.Port, r.Protocol, r.Reason, expires)
		}
		return nil
	},
}

func init() {
	suppressCmd.PersistentFlags().StringVar(&suppressFile, "file", "suppress.json", "path to suppression rules file")
	suppressAddCmd.Flags().String("reason", "", "reason for suppression")
	suppressAddCmd.Flags().Duration("ttl", 0, "duration before rule expires (e.g. 24h)")
	suppressCmd.AddCommand(suppressAddCmd)
	suppressCmd.AddCommand(suppressListCmd)
	rootCmd.AddCommand(suppressCmd)
}
