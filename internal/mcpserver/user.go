package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerUserTools(srv *server.MCPServer, provider ClientProvider) {
	userCurrentTool := newTool("user_current", "获取当前 GitLab token 对应的用户", noArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(userCurrentTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args ProjectArgs) (map[string]any, error) {
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		out, err := client.CurrentUser()
		if err != nil {
			debugLog(provider, req, "current user failed", "error", err)
			return nil, err
		}
		debugLog(provider, req, "current user ok", "username", out["username"], "id", out["id"])
		return out, nil
	}))
}
