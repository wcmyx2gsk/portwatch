package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/user/portwatch/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage portwatch configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the active configuration as YAML",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		out, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(os.Stdout, string(out))
		return err
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Write a default config file to disk",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		if cfgPath == "" {
			cfgPath = "portwatch.yaml"
		}
		if _, err := os.Stat(cfgPath); err == nil {
			return fmt.Errorf("config file already exists: %s", cfgPath)
		}
		if err := config.Save(cfgPath, config.Default()); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}
		fmt.Printf("Config written to %s\n", cfgPath)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.PersistentFlags().String("config", "", "path to config file")
	rootCmd.AddCommand(configCmd)
}
