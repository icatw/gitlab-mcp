package mcpserver

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerProjectTools(srv *server.MCPServer, provider ClientProvider, defaultProject string) {
	projectTool := newTool("project_get", "获取项目信息", []mcp.ToolOption{projectArg()}, mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(projectTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args ProjectArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		return client.GetProject(project)
	}))
}
