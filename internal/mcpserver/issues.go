package mcpserver

import (
	"context"
	"errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"net/url"
	"strconv"
	"strings"
)

func registerIssueTools(srv *server.MCPServer, provider ClientProvider, defaultProject string) {
	issuesListTool := newTool("issues_list", "列出 Issue", issueListArgs(), mcp.WithOutputSchema[ListResult]())
	srv.AddTool(issuesListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args IssueListArgs) (ListResult, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return ListResult{}, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return ListResult{}, err
		}
		q := url.Values{}
		setString(q, "state", args.State)
		setString(q, "labels", args.Labels)
		setInt(q, "author_id", args.AuthorID)
		setInt(q, "assignee_id", args.AssigneeID)
		setString(q, "created_after", args.CreatedAfter)
		setString(q, "updated_after", args.UpdatedAfter)
		setString(q, "search", args.Search)
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		items, err := client.ListIssues(project, q, args.All)
		if err != nil {
			return ListResult{}, err
		}
		return listResult(items), nil
	}))

	issueGetTool := newTool("issues_get", "获取单个 Issue", issueGetArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(issueGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args IssueGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if args.IID <= 0 {
			return nil, errors.New("iid is required")
		}
		return client.GetIssue(project, strconv.Itoa(args.IID))
	}))

	issuesCreateTool := newTool("issues_create", "创建 Issue", issueCreateArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(issuesCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args IssueCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(args.Title)
		if title == "" {
			return nil, errors.New("title is required")
		}
		if args.DueDate != "" && !isYYYYMMDD(args.DueDate) {
			return nil, errors.New("due_date must be in YYYY-MM-DD format")
		}

		body := map[string]any{
			"title": title,
		}
		setStringMap(body, "description", args.Description)
		setStringMap(body, "labels", args.Labels)
		setIntMap(body, "assignee_id", args.AssigneeID)
		setIntMap(body, "milestone_id", args.MilestoneID)
		setStringMap(body, "due_date", args.DueDate)
		setStringMap(body, "issue_type", args.IssueType)
		if args.Confidential {
			body["confidential"] = true
		}
		return client.CreateIssue(project, body)
	}))

	issuesUpdateTool := newTool("issues_update", "更新 Issue（状态、标签、指派、描述等）", issueUpdateArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(issuesUpdateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args IssueUpdateArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if args.IID <= 0 {
			return nil, errors.New("iid is required")
		}
		if args.DueDate != "" && !isYYYYMMDD(args.DueDate) {
			return nil, errors.New("due_date must be in YYYY-MM-DD format")
		}
		body := map[string]any{}
		setStringMap(body, "title", args.Title)
		setStringMap(body, "description", args.Description)
		setStringMap(body, "labels", args.Labels)
		setStringMap(body, "add_labels", args.AddLabels)
		setStringMap(body, "remove_labels", args.RemoveLabels)
		setIntMap(body, "assignee_id", args.AssigneeID)
		setIntMap(body, "milestone_id", args.MilestoneID)
		setStringMap(body, "due_date", args.DueDate)
		setStringMap(body, "state_event", args.StateEvent)
		setBoolPtrMap(body, "confidential", args.Confidential)
		setBoolPtrMap(body, "discussion_locked", args.DiscussionLocked)
		if len(body) == 0 {
			return nil, errors.New("at least one updatable field is required")
		}
		return client.UpdateIssue(project, strconv.Itoa(args.IID), body)
	}))

	issueNotesTool := newTool("issues_notes", "获取 Issue 评论", noteListArgs(), mcp.WithOutputSchema[ListResult]())
	srv.AddTool(issueNotesTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args NoteListArgs) (ListResult, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return ListResult{}, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return ListResult{}, err
		}
		if args.IID <= 0 {
			return ListResult{}, errors.New("iid is required")
		}
		q := url.Values{}
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		items, err := client.GetIssueNotes(project, strconv.Itoa(args.IID), q, args.All)
		if err != nil {
			return ListResult{}, err
		}
		return listResult(items), nil
	}))

	issueNotesCreateTool := newTool("issues_notes_create", "创建 Issue 评论", noteCreateArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(issueNotesCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args NoteCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if args.IID <= 0 {
			return nil, errors.New("iid is required")
		}
		body := strings.TrimSpace(args.Body)
		if body == "" {
			return nil, errors.New("body is required")
		}
		return client.CreateIssueNote(project, strconv.Itoa(args.IID), map[string]any{"body": body})
	}))

}
