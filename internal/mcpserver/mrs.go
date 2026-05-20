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

func registerMRTools(srv *server.MCPServer, provider ClientProvider, defaultProject string) {
	mrsListTool := newTool("mrs_list", "列出 Merge Request", mrListArgs(), mcp.WithOutputSchema[ListResult]())
	srv.AddTool(mrsListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRListArgs) (ListResult, error) {
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
		items, err := client.ListMergeRequests(project, q, args.All)
		if err != nil {
			debugLog(provider, req, "mrs_list failed", "project", project, "state", args.State, "per_page", args.PerPage, "error", err)
			return ListResult{}, err
		}
		debugLog(provider, req, "mrs_list ok", "project", project, "state", args.State, "per_page", args.PerPage, "count", len(items))
		return listResult(items), nil
	}))

	mrGetTool := newTool("mrs_get", "获取单个 Merge Request", mrGetArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(mrGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRGetArgs) (map[string]any, error) {
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
		return client.GetMergeRequest(project, strconv.Itoa(args.IID))
	}))

	mrCreateTool := newTool("mrs_create", "创建 Merge Request", mrCreateArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(mrCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.SourceBranch) == "" || strings.TrimSpace(args.TargetBranch) == "" || strings.TrimSpace(args.Title) == "" {
			return nil, errors.New("source_branch, target_branch and title are required")
		}
		body := map[string]any{
			"source_branch": args.SourceBranch,
			"target_branch": args.TargetBranch,
			"title":         args.Title,
		}
		setStringMap(body, "description", args.Description)
		setIntMap(body, "assignee_id", args.AssigneeID)
		setStringMap(body, "labels", args.Labels)
		if args.RemoveSourceBranch {
			body["remove_source_branch"] = true
		}
		if args.Squash {
			body["squash"] = true
		}
		if args.Draft {
			body["draft"] = true
		}
		setBoolPtrMap(body, "allow_collaboration", args.AllowCollaboration)
		return client.CreateMergeRequest(project, body)
	}))

	mrUpdateTool := newTool("mrs_update", "更新 Merge Request（状态、标签、目标分支、描述等）", mrUpdateArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(mrUpdateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRUpdateArgs) (map[string]any, error) {
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
		body := map[string]any{}
		setStringMap(body, "target_branch", args.TargetBranch)
		setStringMap(body, "title", args.Title)
		setStringMap(body, "description", args.Description)
		setIntMap(body, "assignee_id", args.AssigneeID)
		setStringMap(body, "labels", args.Labels)
		setStringMap(body, "add_labels", args.AddLabels)
		setStringMap(body, "remove_labels", args.RemoveLabels)
		setStringMap(body, "state_event", args.StateEvent)
		setBoolPtrMap(body, "squash", args.Squash)
		setBoolPtrMap(body, "discussion_locked", args.DiscussionLocked)
		setBoolPtrMap(body, "remove_source_branch", args.RemoveSourceBranch)
		if len(body) == 0 {
			return nil, errors.New("at least one updatable field is required")
		}
		return client.UpdateMergeRequest(project, strconv.Itoa(args.IID), body)
	}))

	mrApproveTool := newTool("mrs_approve", "审批 Merge Request", mrApproveArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(mrApproveTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRApproveArgs) (map[string]any, error) {
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
		body := map[string]any{}
		setStringMap(body, "sha", args.Sha)
		setStringMap(body, "approval_password", args.ApprovalPwd)
		return client.ApproveMergeRequest(project, strconv.Itoa(args.IID), body)
	}))

	mrMergeTool := newTool("mrs_merge", "合并 Merge Request", mrMergeArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(mrMergeTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRMergeArgs) (map[string]any, error) {
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
		body := map[string]any{}
		setStringMap(body, "merge_commit_message", args.MergeCommitMessage)
		setStringMap(body, "squash_commit_message", args.SquashCommitMessage)
		setBoolPtrMap(body, "should_remove_source_branch", args.ShouldRemoveSourceBranch)
		if args.MergeWhenPipelineSucceeds {
			body["merge_when_pipeline_succeeds"] = true
		}
		setStringMap(body, "sha", args.Sha)
		setBoolPtrMap(body, "squash", args.Squash)
		return client.MergeMergeRequest(project, strconv.Itoa(args.IID), body)
	}))

	mrChangesTool := newTool("mrs_changes", "获取 MR 变更详情", mrGetArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(mrChangesTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRGetArgs) (map[string]any, error) {
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
		return client.GetMergeRequestChanges(project, strconv.Itoa(args.IID))
	}))

	mrNotesTool := newTool("mrs_notes", "获取 MR 评论", noteListArgs(), mcp.WithOutputSchema[ListResult]())
	srv.AddTool(mrNotesTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args NoteListArgs) (ListResult, error) {
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
		items, err := client.GetMergeRequestNotes(project, strconv.Itoa(args.IID), q, args.All)
		if err != nil {
			return ListResult{}, err
		}
		return listResult(items), nil
	}))

	mrNotesCreateTool := newTool("mrs_notes_create", "创建 MR 评论", noteCreateArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(mrNotesCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args NoteCreateArgs) (map[string]any, error) {
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
		return client.CreateMergeRequestNote(project, strconv.Itoa(args.IID), map[string]any{"body": body})
	}))

}
