package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/jneo8/skeleton-mcp-server/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ServeOptions contains configuration for the serve command
type ServeOptions struct {
	ConfigOptions ConfigOptions
	// CustomFlagSetup allows users to add their own flags
	CustomFlagSetup func(*cobra.Command) error
}

// NewServeCommand creates a serve command that combines default config with user config
func NewServeCommand(application app.App, opts ServeOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config (combines file, env vars, and flags)
			cfg, err := LoadConfig(opts.ConfigOptions)
			if err != nil {
				return err
			}

			// Update the app's config with the loaded configuration
			appCfg := application.GetConfig()
			*appCfg = *cfg

			return app.Run(application)
		},
	}

	// Add default flags
	addDefaultFlags(cmd)

	// Add custom flags if provided
	if opts.CustomFlagSetup != nil {
		if err := opts.CustomFlagSetup(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up custom flags: %v\n", err)
			os.Exit(1)
		}
	}

	return cmd
}

// addDefaultFlags adds the default MCP server flags
func addDefaultFlags(cmd *cobra.Command) {
	// Logging flags
	cmd.Flags().String("log-level", "info", "log level (debug, info, warn, error)")
	cmd.Flags().String("log-format", "text", "log format (text, json)")

	// Server flags
	cmd.Flags().String("transport", "stdio", "transport type (stdio, http)")
	cmd.Flags().Int("port", 8080, "port for http transport")
	cmd.Flags().String("host", "localhost", "host for http transport")
	cmd.Flags().Duration("transport-timeout", 30*time.Second, "transport timeout")
	cmd.Flags().Bool("read-only", false, "enable read-only mode")

	// Bind flags to viper
	flagBindings := map[string]string{
		"logging.level":            "log-level",
		"logging.format":           "log-format",
		"server.transport_type":    "transport",
		"server.port":              "port",
		"server.host":              "host",
		"server.transport_timeout": "transport-timeout",
		"server.read_only":         "read-only",
	}
	for key, flag := range flagBindings {
		if err := viper.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
			fmt.Fprintf(os.Stderr, "Error binding flag %s: %v\n", flag, err)
			os.Exit(1)
		}
	}
}
