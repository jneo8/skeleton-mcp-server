package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jneo8/skeleton-mcp-server/pkg/cli"
	"github.com/jneo8/skeleton-mcp-server/pkg/config"
	"github.com/jneo8/skeleton-mcp-server/pkg/mcp"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

// AppConfig contains application-specific configuration
type AppConfig struct {
	MaxPingCount int `mapstructure:"max_ping_count"`
}

// ServerApp implements the app.App interface
type ServerApp struct {
	cfg    *config.Config
	appCfg *AppConfig
	once   sync.Once
}

func NewServerApp(cfg *config.Config) *ServerApp {
	return &ServerApp{
		cfg:    cfg,
		appCfg: &AppConfig{},
	}
}

func (a *ServerApp) GetConfig() *config.Config {
	return a.cfg
}

func (a *ServerApp) Init() error {
	var err error
	a.once.Do(func() {
		if a.cfg == nil {
			err = fmt.Errorf("configuration is nil")
			return
		}

		config.SetupLogger(a.cfg.Logging.Level, a.cfg.Logging.Format)

		if err = a.cfg.Validate(); err != nil {
			err = fmt.Errorf("configuration validation failed: %w", err)
			return
		}

		// TODO: Setup logging based on a.cfg.Logging
		log.Info().Str("log_level", a.cfg.Logging.Level).Msg("Logging initialized")
	})
	return err
}

func (a *ServerApp) Shutdown(ctx context.Context) error {
	log.Info().Msg("ServerApp shutting down...")
	// Perform any necessary cleanup here
	return nil
}

// PingHandler is a simple handler that provides a "ping" tool.
type PingHandler struct {
	maxCount int
	counter  int
}

func NewPingHandler(maxCount int) *PingHandler {
	return &PingHandler{maxCount: maxCount}
}

func (h *PingHandler) AddTool(mcpServer *server.MCPServer, readOnly bool) error {
	pingTool := mcpgo.NewTool(
		"ping",
		mcpgo.WithDescription("A simple tool that responds with pong."),
	)
	pingHandler := func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		h.counter++
		if h.maxCount > 0 && h.counter > h.maxCount {
			return mcpgo.NewToolResultError(fmt.Sprintf("ping limit reached (max: %d)", h.maxCount)), nil
		}
		msg := fmt.Sprintf("pong (count: %d", h.counter)
		if h.maxCount > 0 {
			msg += fmt.Sprintf("/%d)", h.maxCount)
		} else {
			msg += ", unlimited)"
		}
		return mcpgo.NewToolResultText(msg), nil
	}
	mcpServer.AddTool(pingTool, pingHandler)
	return nil
}

func (a *ServerApp) GetHandlers() ([]mcp.Handler, error) {
	return []mcp.Handler{NewPingHandler(a.appCfg.MaxPingCount)}, nil
}

func main() {
	// Create app instance with initial config
	application := NewServerApp(&config.Config{})

	rootCmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "A skeleton MCP server",
		Long:  `A skeleton implementation of a server that uses the Model Context Protocol.`,
	}

	// Create serve command with combined default + custom configuration
	serveOpts := cli.ServeOptions{
		ConfigOptions: cli.ConfigOptions{
			EnvPrefix: "SIMPLE_MCP",
		},
		// Add custom flags for this application
		CustomFlagSetup: func(cmd *cobra.Command) error {
			// Add application-specific flag
			cmd.Flags().Int("max-ping-count", 0, "maximum number of ping requests (0 = unlimited)")

			// Bind flag to viper so it can be loaded into appCfg
			if err := viper.BindPFlag("max_ping_count", cmd.Flags().Lookup("max-ping-count")); err != nil {
				return err
			}

			// Load app-specific config from viper after all flags are bound
			cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
				return viper.Unmarshal(application.appCfg)
			}

			return nil
		},
	}

	rootCmd.AddCommand(cli.NewServeCommand(application, serveOpts))
	rootCmd.AddCommand(cli.NewVersionCommand(version, commit, buildDate))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
