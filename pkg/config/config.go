package config

import "time"

// Config represents the complete application configuration
type Config struct {
	MCP     MCPConfig     `mapstructure:"mcp"`
	Logging LoggingConfig `mapstructure:"logging"`
}

// MCPConfig contains MCP server settings
type MCPConfig struct {
	// Transport configuration
	Transport TransportConfig `mapstructure:"transport"`

	// Server behavior
	ServerName    string `mapstructure:"server_name"`
	ServerVersion string `mapstructure:"server_version"`

	// Feature flag
	ReadOnly bool `mapstructure:"read_only"`
}

// TransportConfig defines how MCP communicates
type TransportConfig struct {
	Type string `mapstructure:"type"` // "stdio" or "http"

	// HTTP-specific settings
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`

	// Connection timeout
	Timeout time.Duration `mapstructure:"timeout"`
}

// LoggingConfig controls application logging
type LoggingConfig struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
}
