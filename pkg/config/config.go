package config

import (
	"time"
)

// Config represents the complete application configuration
type Config struct {
	Logging         LoggingConfig `mapstructure:"logging"`
	Server          ServerConfig  `mapstructure:"server"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, text
}

// ServerConfig contains MCP server settings
type ServerConfig struct {
	// Transport configuration
	TransportType string `mapstructure:"transport_type"` // "stdio" or "http"
	// HTTP-specific settings
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	// Connection timeout
	TransportTimeout time.Duration `mapstructure:"transport_timeout"`
	// Feature flag
	ReadOnly bool `mapstructure:"read_only"`
}