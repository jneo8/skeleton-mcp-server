package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/jneo8/skeleton-mcp-server/pkg/app"
	"github.com/jneo8/skeleton-mcp-server/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppFactory is a function that creates an app.App
type AppFactory func(cfg *config.Config) app.App

// ServeRunE returns a cobra.RunE function that loads config, creates an app, and runs it.
func ServeRunE(factory AppFactory, opts ConfigOptions) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfig(opts)
		if err != nil {
			return err
		}
		application := factory(cfg)
		return app.Run(application)
	}
}

func NewServeCommand(factory AppFactory, opts ConfigOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server",
		RunE:  ServeRunE(factory, opts),
	}

	AddPersistentFlags(cmd)

	cmd.Flags().String("transport", "stdio", "transport type (stdio, http)")
	cmd.Flags().Int("port", 8080, "port for http transport")
	cmd.Flags().String("host", "localhost", "host for http transport")
	cmd.Flags().Duration("transport-timeout", 30*time.Second, "transport timeout")

	flagBindings := map[string]string{
		"server.transport_type":    "transport",
		"server.port":              "port",
		"server.host":              "host",
		"server.transport_timeout": "transport-timeout",
	}
	for key, flag := range flagBindings {
		if err := viper.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
			fmt.Fprintf(os.Stderr, "Error binding flag %s: %v\n", flag, err)
			os.Exit(1)
		}
	}

	return cmd
}
