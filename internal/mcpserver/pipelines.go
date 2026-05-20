package mcpserver

import (
	"context"
	"errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"net/url"
	"strconv"
)

func registerPipelineTools(srv *server.MCPServer, provider ClientProvider, defaultProject string) {
	pipelinesListTool := newTool("pipelines_list", "列出 Pipelines", pipelineListArgs(), mcp.WithOutputSchema[ListResult]())
	srv.AddTool(pipelinesListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args PipelineListArgs) (ListResult, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return ListResult{}, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return ListResult{}, err
		}
		q := url.Values{}
		setString(q, "ref", args.Ref)
		setString(q, "status", args.Status)
		setString(q, "scope", args.Scope)
		setString(q, "source", args.Source)
		setString(q, "username", args.Username)
		setString(q, "updated_after", args.UpdatedAfter)
		setString(q, "updated_before", args.UpdatedBefore)
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		items, err := client.ListPipelines(project, q, args.All)
		if err != nil {
			return ListResult{}, err
		}
		return listResult(items), nil
	}))

	pipelineGetTool := newTool("pipeline_get", "获取单个 Pipeline", pipelineGetArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(pipelineGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args PipelineGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if args.ID <= 0 {
			return nil, errors.New("id is required")
		}
		return client.GetPipeline(project, strconv.Itoa(args.ID))
	}))

	pipelineRetryTool := newTool("pipeline_retry", "重试 Pipeline", pipelineGetArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(pipelineRetryTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args PipelineGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if args.ID <= 0 {
			return nil, errors.New("id is required")
		}
		return client.RetryPipeline(project, strconv.Itoa(args.ID))
	}))

}
