package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

// Handler defines the interface for registering MCP tools
type Handler interface {
	AddTool(mcpServer *server.MCPServer, readOnly bool) error
}
