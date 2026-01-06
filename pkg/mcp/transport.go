package mcp

import (
	"context"
	"fmt"

	"github.com/jneo8/skeleton-mcp-server/pkg/config"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

type TransportStarter func(ctx context.Context, cfg *config.Config, mcpServer *server.MCPServer) error

func GetTransportStarter(transportType string) (TransportStarter, error) {
	switch transportType {
	case "stdio":
		return startStdio, nil
	case "http":
		return startHTTP, nil
	default:
		return nil, fmt.Errorf("unsupported transport type: %s", transportType)
	}
}

func startStdio(ctx context.Context, cfg *config.Config, mcpServer *server.MCPServer) error {
	log.Info().Msg("Starting stdio transport")
	errChan := make(chan error, 1)
	go func() {
		if err := server.ServeStdio(mcpServer); err != nil {
			errChan <- fmt.Errorf("stdio server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("Context cancelled, shutting down stdio server")
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func startHTTP(ctx context.Context, cfg *config.Config, mcpServer *server.MCPServer) error {
	log.Info().
		Str("host", cfg.Server.Host).
		Int("port", cfg.Server.Port).
		Msg("Starting HTTP streamable transport")

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer := server.NewStreamableHTTPServer(mcpServer)

	errChan := make(chan error, 1)
	go func() {
		log.Info().Str("address", addr).Msg("HTTP server listening")
		if err := httpServer.Start(addr); err != nil {
			errChan <- fmt.Errorf("http server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("Context cancelled, shutting down HTTP server")
		if err := httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("error shutting down HTTP server: %w", err)
		}
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
