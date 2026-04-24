package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	ScanInterval  time.Duration `yaml:"scan_interval"`
	BaselineFile  string        `yaml:"baseline_file"`
	LogFile       string        `yaml:"log_file"`
	AlertOnNew    bool          `yaml:"alert_on_new"`
	AlertOnGone   bool          `yaml:"alert_on_gone"`
	IgnoredPorts  []uint16      `yaml:"ignored_ports"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		ScanInterval: 30 * time.Second,
		BaselineFile: "baseline.json",
		LogFile:      "",
		AlertOnNew:   true,
		AlertOnGone:  true,
		IgnoredPorts: []uint16{},
	}
}

// Load reads a YAML config file from path and merges it over defaults.
// If the file does not exist, the defaults are returned without error.
func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save writes the Config to a YAML file at path.
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
