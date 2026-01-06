package cli

import (
	"fmt"
	"strings"

	"github.com/jneo8/skeleton-mcp-server/pkg/config"
	"github.com/spf13/viper"
)

// ConfigOptions defines the options for loading configuration
type ConfigOptions struct {
	ConfigName  string
	ConfigPaths []string
	EnvPrefix   string
}

// initConfig initializes viper configuration
func InitConfig(opts ConfigOptions) error {
	viper.SetConfigName(opts.ConfigName)
	viper.SetConfigType("yaml")

	for _, path := range opts.ConfigPaths {
		viper.AddConfigPath(path)
	}

	viper.SetEnvPrefix(opts.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("reading config file: %w", err)
		}
	}

	return nil
}

// setDefaults sets default values in viper
func setDefaults() {
	viper.SetDefault("server.transport_type", "stdio")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.transport_timeout", "30s")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("shutdown_timeout", "10s")
}

// LoadConfig loads and validates configuration from all sources
func LoadConfig(opts ConfigOptions) (*config.Config, error) {
	if err := InitConfig(opts); err != nil {
		return nil, err
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}
