package cli

import (
	"fmt"
	"strings"

	"github.com/jneo8/skeleton-mcp-server/pkg/config"
	"github.com/spf13/viper"
)

// ConfigOptions defines the options for loading configuration
type ConfigOptions struct {
	EnvPrefix string
}

// LoadConfig loads and validates configuration from flags and environment variables
func LoadConfig(opts ConfigOptions) (*config.Config, error) {
	// Set environment variable support
	if opts.EnvPrefix != "" {
		viper.SetEnvPrefix(opts.EnvPrefix)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		viper.AutomaticEnv()
	}

	// Unmarshal config from viper (flags + env vars)
	// Defaults come from the flag definitions
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}
