package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"gitlab-mcp/internal/gitlab"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Config struct {
	BaseURL  string `json:"base_url"`
	Project  string `json:"project"`
	Insecure bool   `json:"insecure"`
}

func main() {
	token := strings.TrimSpace(os.Getenv("GITLAB_TOKEN"))
	if token == "" {
		exitOnErr(errors.New("GITLAB_TOKEN is required"))
	}

	baseURL := strings.TrimSpace(os.Getenv("GITLAB_BASE_URL"))
	if baseURL == "" {
		baseURL = "https://gitlab.example.com"
	}
	defaultProject := strings.TrimSpace(os.Getenv("GITLAB_PROJECT"))
	insecure := parseBoolEnv(os.Getenv("GITLAB_INSECURE"), true)

	client, err := gitlab.NewClient(baseURL, token, insecure)
	exitOnErr(err)

	srv := server.NewMCPServer("gitlab-mcp", "0.3.0")
	registerTools(srv, client, defaultProject)

	exitOnErr(server.ServeStdio(srv))
}

func registerTools(srv *server.MCPServer, client *gitlab.Client, defaultProject string) {
	projectTool := mcp.NewTool("project_get",
		mcp.WithDescription("获取项目信息"),
		mcp.WithInputSchema[ProjectArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(projectTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args ProjectArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		return client.GetProject(project)
	}))

	issuesListTool := mcp.NewTool("issues_list",
		mcp.WithDescription("列出 Issue"),
		mcp.WithInputSchema[IssueListArgs](),
		mcp.WithOutputSchema[[]map[string]any](),
	)
	srv.AddTool(issuesListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args IssueListArgs) ([]map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
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
		return client.ListIssues(project, q, args.All)
	}))

	issueGetTool := mcp.NewTool("issues_get",
		mcp.WithDescription("获取单个 Issue"),
		mcp.WithInputSchema[IssueGetArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(issueGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args IssueGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		if args.IID <= 0 {
			return nil, errors.New("iid is required")
		}
		return client.GetIssue(project, strconv.Itoa(args.IID))
	}))

	issuesCreateTool := mcp.NewTool("issues_create",
		mcp.WithDescription("创建 Issue"),
		mcp.WithInputSchema[IssueCreateArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(issuesCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args IssueCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	issuesUpdateTool := mcp.NewTool("issues_update",
		mcp.WithDescription("更新 Issue（状态、标签、指派、描述等）"),
		mcp.WithInputSchema[IssueUpdateArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(issuesUpdateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args IssueUpdateArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	issueNotesTool := mcp.NewTool("issues_notes",
		mcp.WithDescription("获取 Issue 评论"),
		mcp.WithInputSchema[NoteListArgs](),
		mcp.WithOutputSchema[[]map[string]any](),
	)
	srv.AddTool(issueNotesTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args NoteListArgs) ([]map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		if args.IID <= 0 {
			return nil, errors.New("iid is required")
		}
		q := url.Values{}
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		return client.GetIssueNotes(project, strconv.Itoa(args.IID), q, args.All)
	}))

	issueNotesCreateTool := mcp.NewTool("issues_notes_create",
		mcp.WithDescription("创建 Issue 评论"),
		mcp.WithInputSchema[NoteCreateArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(issueNotesCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args NoteCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	mrsListTool := mcp.NewTool("mrs_list",
		mcp.WithDescription("列出 Merge Request"),
		mcp.WithInputSchema[MRListArgs](),
		mcp.WithOutputSchema[[]map[string]any](),
	)
	srv.AddTool(mrsListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRListArgs) ([]map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
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
		return client.ListMergeRequests(project, q, args.All)
	}))

	mrGetTool := mcp.NewTool("mrs_get",
		mcp.WithDescription("获取单个 Merge Request"),
		mcp.WithInputSchema[MRGetArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(mrGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		if args.IID <= 0 {
			return nil, errors.New("iid is required")
		}
		return client.GetMergeRequest(project, strconv.Itoa(args.IID))
	}))

	mrCreateTool := mcp.NewTool("mrs_create",
		mcp.WithDescription("创建 Merge Request"),
		mcp.WithInputSchema[MRCreateArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(mrCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	mrUpdateTool := mcp.NewTool("mrs_update",
		mcp.WithDescription("更新 Merge Request（状态、标签、目标分支、描述等）"),
		mcp.WithInputSchema[MRUpdateArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(mrUpdateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRUpdateArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	mrApproveTool := mcp.NewTool("mrs_approve",
		mcp.WithDescription("审批 Merge Request"),
		mcp.WithInputSchema[MRApproveArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(mrApproveTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRApproveArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	mrMergeTool := mcp.NewTool("mrs_merge",
		mcp.WithDescription("合并 Merge Request"),
		mcp.WithInputSchema[MRMergeArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(mrMergeTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRMergeArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	mrChangesTool := mcp.NewTool("mrs_changes",
		mcp.WithDescription("获取 MR 变更详情"),
		mcp.WithInputSchema[MRGetArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(mrChangesTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args MRGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		if args.IID <= 0 {
			return nil, errors.New("iid is required")
		}
		return client.GetMergeRequestChanges(project, strconv.Itoa(args.IID))
	}))

	mrNotesTool := mcp.NewTool("mrs_notes",
		mcp.WithDescription("获取 MR 评论"),
		mcp.WithInputSchema[NoteListArgs](),
		mcp.WithOutputSchema[[]map[string]any](),
	)
	srv.AddTool(mrNotesTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args NoteListArgs) ([]map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		if args.IID <= 0 {
			return nil, errors.New("iid is required")
		}
		q := url.Values{}
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		return client.GetMergeRequestNotes(project, strconv.Itoa(args.IID), q, args.All)
	}))

	mrNotesCreateTool := mcp.NewTool("mrs_notes_create",
		mcp.WithDescription("创建 MR 评论"),
		mcp.WithInputSchema[NoteCreateArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(mrNotesCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args NoteCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	pipelinesListTool := mcp.NewTool("pipelines_list",
		mcp.WithDescription("列出 Pipelines"),
		mcp.WithInputSchema[PipelineListArgs](),
		mcp.WithOutputSchema[[]map[string]any](),
	)
	srv.AddTool(pipelinesListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args PipelineListArgs) ([]map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
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
		return client.ListPipelines(project, q, args.All)
	}))

	pipelineGetTool := mcp.NewTool("pipeline_get",
		mcp.WithDescription("获取单个 Pipeline"),
		mcp.WithInputSchema[PipelineGetArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(pipelineGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args PipelineGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		if args.ID <= 0 {
			return nil, errors.New("id is required")
		}
		return client.GetPipeline(project, strconv.Itoa(args.ID))
	}))

	pipelineRetryTool := mcp.NewTool("pipeline_retry",
		mcp.WithDescription("重试 Pipeline"),
		mcp.WithInputSchema[PipelineGetArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(pipelineRetryTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args PipelineGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		if args.ID <= 0 {
			return nil, errors.New("id is required")
		}
		return client.RetryPipeline(project, strconv.Itoa(args.ID))
	}))

	branchesListTool := mcp.NewTool("branches_list",
		mcp.WithDescription("列出分支"),
		mcp.WithInputSchema[BranchListArgs](),
		mcp.WithOutputSchema[[]map[string]any](),
	)
	srv.AddTool(branchesListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args BranchListArgs) ([]map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		q := url.Values{}
		setString(q, "search", args.Search)
		setString(q, "regex", args.Regex)
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		return client.ListBranches(project, q, args.All)
	}))

	branchesCreateTool := mcp.NewTool("branches_create",
		mcp.WithDescription("创建分支"),
		mcp.WithInputSchema[BranchCreateArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(branchesCreateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args BranchCreateArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

	commitsListTool := mcp.NewTool("commits_list",
		mcp.WithDescription("列出提交记录"),
		mcp.WithInputSchema[CommitListArgs](),
		mcp.WithOutputSchema[[]map[string]any](),
	)
	srv.AddTool(commitsListTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args CommitListArgs) ([]map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		q := url.Values{}
		setString(q, "ref_name", args.RefName)
		setString(q, "path", args.Path)
		setString(q, "since", args.Since)
		setString(q, "until", args.Until)
		setInt(q, "per_page", args.PerPage)
		setInt(q, "page", args.Page)
		return client.ListCommits(project, q, args.All)
	}))

	commitGetTool := mcp.NewTool("commit_get",
		mcp.WithDescription("获取单个提交"),
		mcp.WithInputSchema[CommitGetArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(commitGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args CommitGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		sha := strings.TrimSpace(args.SHA)
		if sha == "" {
			return nil, errors.New("sha is required")
		}
		return client.GetCommit(project, sha)
	}))

	repoFileGetTool := mcp.NewTool("repository_file_get",
		mcp.WithDescription("获取仓库文件（包含 base64 内容）"),
		mcp.WithInputSchema[RepoFileGetArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(repoFileGetTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args RepoFileGetArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.FilePath) == "" || strings.TrimSpace(args.Ref) == "" {
			return nil, errors.New("file_path and ref are required")
		}
		return client.GetRepositoryFile(project, args.FilePath, args.Ref)
	}))

	repoFileUpdateTool := mcp.NewTool("repository_file_update",
		mcp.WithDescription("更新仓库文件"),
		mcp.WithInputSchema[RepoFileUpdateArgs](),
		mcp.WithOutputSchema[map[string]any](),
	)
	srv.AddTool(repoFileUpdateTool, mcp.NewStructuredToolHandler(func(ctx context.Context, req mcp.CallToolRequest, args RepoFileUpdateArgs) (map[string]any, error) {
		project, err := projectOrDefault(args.Project, defaultProject)
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

func projectOrDefault(project, fallback string) (string, error) {
	project = strings.TrimSpace(project)
	if project != "" {
		return project, nil
	}
	fallback = strings.TrimSpace(fallback)
	if fallback == "" {
		return "", errors.New("project is required")
	}
	return fallback, nil
}

func parseBoolEnv(raw string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}

func setString(q url.Values, key, val string) {
	if strings.TrimSpace(val) == "" {
		return
	}
	q.Set(key, val)
}

func setInt(q url.Values, key string, val int) {
	if val <= 0 {
		return
	}
	q.Set(key, strconv.Itoa(val))
}

func setStringMap(m map[string]any, key, val string) {
	if strings.TrimSpace(val) == "" {
		return
	}
	m[key] = val
}

func setIntMap(m map[string]any, key string, val int) {
	if val <= 0 {
		return
	}
	m[key] = val
}

func setBoolPtrMap(m map[string]any, key string, val *bool) {
	if val == nil {
		return
	}
	m[key] = *val
}

func isYYYYMMDD(s string) bool {
	if len(s) != 10 {
		return false
	}
	if s[4] != '-' || s[7] != '-' {
		return false
	}
	for i, ch := range s {
		if i == 4 || i == 7 {
			continue
		}
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func exitOnErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
