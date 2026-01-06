package config

import (
	"fmt"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// Validate performs comprehensive configuration validation
func (c *Config) Validate() error {
	var errors ValidationErrors

	// Validate Logging configuration
	errors = append(errors, c.validateLogging()...)

	// Validate Server configuration
	errors = append(errors, c.validateServer()...)

	if c.ShutdownTimeout < 0 {
		errors = append(errors, ValidationError{
			Field:   "shutdown_timeout",
			Message: "must be greater than or equal to 0",
		})
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// validateLogging validates logging configuration
func (c *Config) validateLogging() []ValidationError {
	var errors []ValidationError

	// Validate log level
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if c.Logging.Level != "" && !validLevels[c.Logging.Level] {
		errors = append(errors, ValidationError{
			Field:   "logging.level",
			Message: "must be one of: debug, info, warn, error",
		})
	}

	// Validate log format
	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	if c.Logging.Format != "" && !validFormats[c.Logging.Format] {
		errors = append(errors, ValidationError{
			Field:   "logging.format",
			Message: "must be one of: json, text",
		})
	}

	return errors
}

// validateServer validates MCP configuration
func (c *Config) validateServer() []ValidationError {
	var errors []ValidationError

	// Validate transport type
	validTypes := map[string]bool{
		"stdio": true,
		"http":  true,
	}
	if c.Server.TransportType != "" && !validTypes[c.Server.TransportType] {
		errors = append(errors, ValidationError{
			Field:   "server.transport_type",
			Message: "must be one of: stdio, http",
		})
	}

	// Validate HTTP-specific settings
	if c.Server.TransportType == "http" {
		if c.Server.Port != 0 && (c.Server.Port <= 0 || c.Server.Port > 65535) {
			errors = append(errors, ValidationError{
				Field:   "server.port",
				Message: "port must be between 1 and 65535",
			})
		}

		if c.Server.Host == "" {
			errors = append(errors, ValidationError{
				Field:   "server.host",
				Message: "host is required for HTTP transport",
			})
		}
	}

	return errors
}
