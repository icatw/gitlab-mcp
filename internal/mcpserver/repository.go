package mcpserver

import (
	"context"
	"errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"net/url"
	"strings"
)

func registerRepositoryTools(srv *server.MCPServer, provider ClientProvider, defaultProject string) {
	branchesListTool := newTool("branches_list", "列出分支", branchListArgs(), mcp.WithOutputSchema[ListResult]())
	srv.AddTool(branchesListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args BranchListArgs) (ListResult, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return ListResult{}, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return ListResult{}, err
		}
		q := url.Values{}
		setString(q, "search", args.Search)
		setString(q, "regex", args.Regex)
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		items, err := client.ListBranches(project, q, args.All)
		if err != nil {
			return ListResult{}, err
		}
		return listResult(items), nil
	}))

	branchesCreateTool := newTool("branches_create", "创建分支", branchCreateArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(branchesCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args BranchCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.Branch) == "" || strings.TrimSpace(args.Ref) == "" {
			return nil, errors.New("branch and ref are required")
		}
		body := map[string]any{
			"branch": args.Branch,
			"ref":    args.Ref,
		}
		return client.CreateBranch(project, body)
	}))

	commitsListTool := newTool("commits_list", "列出提交记录", commitListArgs(), mcp.WithOutputSchema[ListResult]())
	srv.AddTool(commitsListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args CommitListArgs) (ListResult, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return ListResult{}, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return ListResult{}, err
		}
		q := url.Values{}
		setString(q, "ref_name", args.RefName)
		setString(q, "path", args.Path)
		setString(q, "since", args.Since)
		setString(q, "until", args.Until)
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		items, err := client.ListCommits(project, q, args.All)
		if err != nil {
			return ListResult{}, err
		}
		return listResult(items), nil
	}))

	commitGetTool := newTool("commit_get", "获取单个提交", commitGetArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(commitGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args CommitGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		sha := strings.TrimSpace(args.SHA)
		if sha == "" {
			return nil, errors.New("sha is required")
		}
		return client.GetCommit(project, sha)
	}))

	repoFileGetTool := newTool("repository_file_get", "获取仓库文件（包含 base64 内容）", repoFileGetArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(repoFileGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args RepoFileGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.FilePath) == "" || strings.TrimSpace(args.Ref) == "" {
			return nil, errors.New("file_path and ref are required")
		}
		return client.GetRepositoryFile(project, args.FilePath, args.Ref)
	}))

	repoFileUpdateTool := newTool("repository_file_update", "更新仓库文件", repoFileUpdateArgs(), mcp.WithOutputSchema[map[string]any]())
	srv.AddTool(repoFileUpdateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args RepoFileUpdateArgs) (map[string]any, error) {
		project, err := projectOrDefault(req, provider, args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		client, err := provider.Client(req)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.FilePath) == "" || strings.TrimSpace(args.Branch) == "" || strings.TrimSpace(args.CommitMessage) == "" {
			return nil, errors.New("file_path, branch and commit_message are required")
		}
		body := map[string]any{
			"branch":         args.Branch,
			"content":        args.Content,
			"commit_message": args.CommitMessage,
		}
		setStringMap(body, "author_email", args.AuthorEmail)
		setStringMap(body, "author_name", args.AuthorName)
		setStringMap(body, "encoding", args.Encoding)
		setBoolPtrMap(body, "execute_filemode", args.ExecuteFilemode)
		setStringMap(body, "last_commit_id", args.LastCommitID)
		return client.UpdateRepositoryFile(project, args.FilePath, body)
	}))
}
