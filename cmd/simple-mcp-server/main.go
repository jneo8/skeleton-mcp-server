package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jneo8/skeleton-mcp-server/pkg/app"
	"github.com/jneo8/skeleton-mcp-server/pkg/cli"
	"github.com/jneo8/skeleton-mcp-server/pkg/config"
	"github.com/jneo8/skeleton-mcp-server/pkg/mcp"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

// ServerApp implements the app.App interface
type ServerApp struct {
	cfg  *config.Config
	once sync.Once
}

func NewServerApp(cfg *config.Config) *ServerApp {
	return &ServerApp{cfg: cfg}
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
type PingHandler struct{}

func (h *PingHandler) AddTool(mcpServer *server.MCPServer, readOnly bool) error {
	pingTool := mcpgo.NewTool(
		"ping",
		mcpgo.WithDescription("A simple tool that responds with pong."),
	)
	pingHandler := func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		return mcpgo.NewToolResultText("pong"), nil
	}
	mcpServer.AddTool(pingTool, pingHandler)
	return nil
}

func (a *ServerApp) GetHandlers() ([]mcp.Handler, error) {
	return []mcp.Handler{&PingHandler{}}, nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "A skeleton MCP server",
		Long:  `A skeleton implementation of a server that uses the Model Context Protocol.`,
	}

	appFactory := func(cfg *config.Config) app.App {
		return NewServerApp(cfg)
	}

	opts := cli.ConfigOptions{
		ConfigName:  "simple-mcp-server",
		ConfigPaths: []string{".", "/etc/simple-mcp-server"},
		EnvPrefix:   "MCP",
	}

	rootCmd.AddCommand(cli.NewServeCommand(appFactory, opts))
	rootCmd.AddCommand(cli.NewVersionCommand(version, commit, buildDate))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
