package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jneo8/skeleton-mcp-server/pkg/config"
	"github.com/jneo8/skeleton-mcp-server/pkg/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

type App interface {
	GetConfig() *config.Config
	Init() error
	Shutdown(ctx context.Context) error
	GetHandlers() ([]mcp.Handler, error)
}

func Run(app App) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.Init(); err != nil {
		return err
	}
	defer func() {
		if err := app.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("App shutdown error")
		}
	}()

	cfg := app.GetConfig()
	mcpServer := server.NewMCPServer(
		cfg.MCP.ServerName,
		cfg.MCP.ServerVersion,
		server.WithToolCapabilities(true),
	)

	handlers, err := app.GetHandlers()
	if err != nil {
		return fmt.Errorf("getting handlers: %w", err)
	}

	for _, h := range handlers {
		if err := h.RegisterTools(mcpServer, cfg.MCP.ReadOnly); err != nil {
			return fmt.Errorf("registering tools: %w", err)
		}
	}

	transportStarter, err := mcp.GetTransportStarter(cfg.MCP.Transport.Type)
	if err != nil {
		return err
	}

	return transportStarter(ctx, cfg, mcpServer)
}
