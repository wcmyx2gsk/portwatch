package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"portwatch/internal/ports"
)

var processCmd = &cobra.Command{
	Use:   "process <inode>",
	Short: "Look up the process owning a socket inode",
	Long: `Resolve a socket inode number to the owning process.

The inode can be obtained from the 'scan' or 'snapshot' output.
Example:
  portwatch process 123456`,
	Args: cobra.ExactArgs(1),
	RunE: runProcess,
}

func runProcess(cmd *cobra.Command, args []string) error {
	inodeVal, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid inode %q: %w", args[0], err)
	}

	info, err := ports.LookupProcess(inodeVal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		fmt.Println("No matching process found.")
		return nil
	}

	fmt.Printf("PID  : %d\n", info.PID)
	fmt.Printf("Name : %s\n", info.Name)
	fmt.Printf("Exe  : %s\n", info.Exe)
	return nil
}

func init() {
	rootCmd.AddCommand(processCmd)
}
