package mcpserver

type ListResult struct {
	Items []map[string]any `json:"items" jsonschema_description:"列表数据"`
	Count int              `json:"count" jsonschema_description:"当前返回数量"`
}

type ProjectArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
}

type IssueGetArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	IID     int    `json:"iid" jsonschema:"required,minimum=1" jsonschema_description:"Issue IID"`
}

type IssueListArgs struct {
	Project      string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	State        string `json:"state,omitempty" jsonschema_description:"opened|closed|all"`
	Labels       string `json:"labels,omitempty" jsonschema_description:"逗号分隔 labels"`
	AuthorID     int    `json:"author_id,omitempty" jsonschema_description:"作者 ID"`
	AssigneeID   int    `json:"assignee_id,omitempty" jsonschema_description:"负责人 ID"`
	CreatedAfter string `json:"created_after,omitempty" jsonschema_description:"ISO 时间，如 2024-01-01T00:00:00Z"`
	UpdatedAfter string `json:"updated_after,omitempty" jsonschema_description:"ISO 时间，如 2024-01-01T00:00:00Z"`
	Search       string `json:"search,omitempty" jsonschema_description:"搜索关键词"`
	PerPage      int    `json:"per_page,omitempty" jsonschema:"minimum=1,maximum=100" jsonschema_description:"分页大小"`
	Page         int    `json:"page,omitempty" jsonschema:"minimum=1" jsonschema_description:"页码"`
	All          bool   `json:"all,omitempty" jsonschema_description:"是否拉取全部分页"`
}

type IssueCreateArgs struct {
	Project      string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	Title        string `json:"title" jsonschema:"required,minLength=1" jsonschema_description:"Issue 标题"`
	Description  string `json:"description,omitempty" jsonschema_description:"Issue 描述（Markdown）"`
	Labels       string `json:"labels,omitempty" jsonschema_description:"逗号分隔 labels"`
	AssigneeID   int    `json:"assignee_id,omitempty" jsonschema_description:"负责人 ID"`
	MilestoneID  int    `json:"milestone_id,omitempty" jsonschema_description:"里程碑 ID"`
	DueDate      string `json:"due_date,omitempty" jsonschema_description:"截止日期，格式 YYYY-MM-DD"`
	IssueType    string `json:"issue_type,omitempty" jsonschema_description:"issue|incident|test_case|task"`
	Confidential bool   `json:"confidential,omitempty" jsonschema_description:"是否保密 Issue"`
}

type IssueUpdateArgs struct {
	Project          string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	IID              int    `json:"iid" jsonschema:"required,minimum=1" jsonschema_description:"Issue IID"`
	Title            string `json:"title,omitempty" jsonschema_description:"Issue 标题"`
	Description      string `json:"description,omitempty" jsonschema_description:"Issue 描述（Markdown）"`
	Labels           string `json:"labels,omitempty" jsonschema_description:"逗号分隔 labels"`
	AddLabels        string `json:"add_labels,omitempty" jsonschema_description:"追加 labels（逗号分隔）"`
	RemoveLabels     string `json:"remove_labels,omitempty" jsonschema_description:"移除 labels（逗号分隔）"`
	AssigneeID       int    `json:"assignee_id,omitempty" jsonschema_description:"负责人 ID"`
	MilestoneID      int    `json:"milestone_id,omitempty" jsonschema_description:"里程碑 ID"`
	DueDate          string `json:"due_date,omitempty" jsonschema_description:"截止日期，格式 YYYY-MM-DD"`
	StateEvent       string `json:"state_event,omitempty" jsonschema_description:"close|reopen"`
	Confidential     *bool  `json:"confidential,omitempty" jsonschema_description:"是否保密"`
	DiscussionLocked *bool  `json:"discussion_locked,omitempty" jsonschema_description:"是否锁定讨论"`
}

type NoteListArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	IID     int    `json:"iid" jsonschema:"required,minimum=1" jsonschema_description:"IID"`
	PerPage int    `json:"per_page,omitempty" jsonschema:"minimum=1,maximum=100" jsonschema_description:"分页大小"`
	Page    int    `json:"page,omitempty" jsonschema:"minimum=1" jsonschema_description:"页码"`
	All     bool   `json:"all,omitempty" jsonschema_description:"是否拉取全部分页"`
}

type NoteCreateArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	IID     int    `json:"iid" jsonschema:"required,minimum=1" jsonschema_description:"Issue/MR IID"`
	Body    string `json:"body" jsonschema:"required,minLength=1" jsonschema_description:"评论内容（Markdown）"`
}

type MRGetArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	IID     int    `json:"iid" jsonschema:"required,minimum=1" jsonschema_description:"MR IID"`
}

type MRListArgs struct {
	Project      string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	State        string `json:"state,omitempty" jsonschema_description:"opened|closed|merged|all"`
	Labels       string `json:"labels,omitempty" jsonschema_description:"逗号分隔 labels"`
	AuthorID     int    `json:"author_id,omitempty" jsonschema_description:"作者 ID"`
	AssigneeID   int    `json:"assignee_id,omitempty" jsonschema_description:"负责人 ID"`
	CreatedAfter string `json:"created_after,omitempty" jsonschema_description:"ISO 时间，如 2024-01-01T00:00:00Z"`
	UpdatedAfter string `json:"updated_after,omitempty" jsonschema_description:"ISO 时间，如 2024-01-01T00:00:00Z"`
	Search       string `json:"search,omitempty" jsonschema_description:"搜索关键词"`
	PerPage      int    `json:"per_page,omitempty" jsonschema:"minimum=1,maximum=100" jsonschema_description:"分页大小"`
	Page         int    `json:"page,omitempty" jsonschema:"minimum=1" jsonschema_description:"页码"`
	All          bool   `json:"all,omitempty" jsonschema_description:"是否拉取全部分页"`
}

type MRCreateArgs struct {
	Project            string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	SourceBranch       string `json:"source_branch" jsonschema:"required,minLength=1" jsonschema_description:"源分支"`
	TargetBranch       string `json:"target_branch" jsonschema:"required,minLength=1" jsonschema_description:"目标分支"`
	Title              string `json:"title" jsonschema:"required,minLength=1" jsonschema_description:"MR 标题"`
	Description        string `json:"description,omitempty" jsonschema_description:"MR 描述（Markdown）"`
	AssigneeID         int    `json:"assignee_id,omitempty" jsonschema_description:"负责人 ID"`
	Labels             string `json:"labels,omitempty" jsonschema_description:"逗号分隔 labels"`
	RemoveSourceBranch bool   `json:"remove_source_branch,omitempty" jsonschema_description:"合并后删除源分支"`
	Squash             bool   `json:"squash,omitempty" jsonschema_description:"是否 squash"`
	Draft              bool   `json:"draft,omitempty" jsonschema_description:"是否草稿"`
	AllowCollaboration *bool  `json:"allow_collaboration,omitempty" jsonschema_description:"是否允许协作者 push"`
}

type MRUpdateArgs struct {
	Project            string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	IID                int    `json:"iid" jsonschema:"required,minimum=1" jsonschema_description:"MR IID"`
	TargetBranch       string `json:"target_branch,omitempty" jsonschema_description:"目标分支"`
	Title              string `json:"title,omitempty" jsonschema_description:"MR 标题"`
	Description        string `json:"description,omitempty" jsonschema_description:"MR 描述（Markdown）"`
	AssigneeID         int    `json:"assignee_id,omitempty" jsonschema_description:"负责人 ID"`
	Labels             string `json:"labels,omitempty" jsonschema_description:"逗号分隔 labels"`
	AddLabels          string `json:"add_labels,omitempty" jsonschema_description:"追加 labels（逗号分隔）"`
	RemoveLabels       string `json:"remove_labels,omitempty" jsonschema_description:"移除 labels（逗号分隔）"`
	StateEvent         string `json:"state_event,omitempty" jsonschema_description:"close|reopen"`
	Squash             *bool  `json:"squash,omitempty" jsonschema_description:"是否 squash"`
	DiscussionLocked   *bool  `json:"discussion_locked,omitempty" jsonschema_description:"是否锁定讨论"`
	RemoveSourceBranch *bool  `json:"remove_source_branch,omitempty" jsonschema_description:"合并后删除源分支"`
}

type MRApproveArgs struct {
	Project     string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	IID         int    `json:"iid" jsonschema:"required,minimum=1" jsonschema_description:"MR IID"`
	Sha         string `json:"sha,omitempty" jsonschema_description:"指定当前 HEAD SHA，防止过期审批"`
	ApprovalPwd string `json:"approval_password,omitempty" jsonschema_description:"审批密码（若启用）"`
}

type MRMergeArgs struct {
	Project                   string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	IID                       int    `json:"iid" jsonschema:"required,minimum=1" jsonschema_description:"MR IID"`
	MergeCommitMessage        string `json:"merge_commit_message,omitempty" jsonschema_description:"合并提交信息"`
	SquashCommitMessage       string `json:"squash_commit_message,omitempty" jsonschema_description:"squash 提交信息"`
	ShouldRemoveSourceBranch  *bool  `json:"should_remove_source_branch,omitempty" jsonschema_description:"是否删除源分支"`
	MergeWhenPipelineSucceeds bool   `json:"merge_when_pipeline_succeeds,omitempty" jsonschema_description:"流水线通过后自动合并"`
	Sha                       string `json:"sha,omitempty" jsonschema_description:"指定当前 HEAD SHA，防止过期合并"`
	Squash                    *bool  `json:"squash,omitempty" jsonschema_description:"是否 squash"`
}

type PipelineListArgs struct {
	Project       string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	Ref           string `json:"ref,omitempty" jsonschema_description:"分支或 tag"`
	Status        string `json:"status,omitempty" jsonschema_description:"created|waiting_for_resource|preparing|pending|running|success|failed|canceled|skipped|manual|scheduled"`
	Scope         string `json:"scope,omitempty" jsonschema_description:"running|pending|finished|branches|tags"`
	Source        string `json:"source,omitempty" jsonschema_description:"push|web|trigger|schedule|api|external|pipeline|chat|webide|merge_request_event"`
	Username      string `json:"username,omitempty" jsonschema_description:"触发者用户名"`
	UpdatedAfter  string `json:"updated_after,omitempty" jsonschema_description:"ISO 时间，如 2024-01-01T00:00:00Z"`
	UpdatedBefore string `json:"updated_before,omitempty" jsonschema_description:"ISO 时间，如 2024-01-01T00:00:00Z"`
	PerPage       int    `json:"per_page,omitempty" jsonschema:"minimum=1,maximum=100" jsonschema_description:"分页大小"`
	Page          int    `json:"page,omitempty" jsonschema:"minimum=1" jsonschema_description:"页码"`
	All           bool   `json:"all,omitempty" jsonschema_description:"是否拉取全部分页"`
}

type PipelineGetArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	ID      int    `json:"id" jsonschema:"required,minimum=1" jsonschema_description:"Pipeline ID"`
}

type BranchListArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	Search  string `json:"search,omitempty" jsonschema_description:"分支搜索关键词"`
	Regex   string `json:"regex,omitempty" jsonschema_description:"分支正则（RE2）"`
	PerPage int    `json:"per_page,omitempty" jsonschema:"minimum=1,maximum=100" jsonschema_description:"分页大小"`
	Page    int    `json:"page,omitempty" jsonschema:"minimum=1" jsonschema_description:"页码"`
	All     bool   `json:"all,omitempty" jsonschema_description:"是否拉取全部分页"`
}

type BranchCreateArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	Branch  string `json:"branch" jsonschema:"required,minLength=1" jsonschema_description:"新分支名"`
	Ref     string `json:"ref" jsonschema:"required,minLength=1" jsonschema_description:"源 ref（分支/标签/提交）"`
}

type CommitListArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	RefName string `json:"ref_name,omitempty" jsonschema_description:"分支或 tag"`
	Path    string `json:"path,omitempty" jsonschema_description:"文件路径过滤"`
	Since   string `json:"since,omitempty" jsonschema_description:"ISO 时间，如 2024-01-01T00:00:00Z"`
	Until   string `json:"until,omitempty" jsonschema_description:"ISO 时间，如 2024-01-31T23:59:59Z"`
	PerPage int    `json:"per_page,omitempty" jsonschema:"minimum=1,maximum=100" jsonschema_description:"分页大小"`
	Page    int    `json:"page,omitempty" jsonschema:"minimum=1" jsonschema_description:"页码"`
	All     bool   `json:"all,omitempty" jsonschema_description:"是否拉取全部分页"`
}

type CommitGetArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	SHA     string `json:"sha" jsonschema:"required,minLength=1" jsonschema_description:"提交 SHA"`
}

type RepoFileGetArgs struct {
	Project  string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	FilePath string `json:"file_path" jsonschema:"required,minLength=1" jsonschema_description:"仓库文件路径"`
	Ref      string `json:"ref" jsonschema:"required,minLength=1" jsonschema_description:"分支、tag 或 commit SHA"`
}

type RepoFileUpdateArgs struct {
	Project         string `json:"project,omitempty" jsonschema_description:"项目路径，未传则使用默认配置"`
	FilePath        string `json:"file_path" jsonschema:"required,minLength=1" jsonschema_description:"仓库文件路径"`
	Branch          string `json:"branch" jsonschema:"required,minLength=1" jsonschema_description:"目标分支"`
	Content         string `json:"content" jsonschema:"required" jsonschema_description:"新的文件内容"`
	CommitMessage   string `json:"commit_message" jsonschema:"required,minLength=1" jsonschema_description:"提交信息"`
	AuthorEmail     string `json:"author_email,omitempty" jsonschema_description:"提交作者邮箱"`
	AuthorName      string `json:"author_name,omitempty" jsonschema_description:"提交作者名"`
	Encoding        string `json:"encoding,omitempty" jsonschema_description:"text|base64"`
	ExecuteFilemode *bool  `json:"execute_filemode,omitempty" jsonschema_description:"是否设置可执行位"`
	LastCommitID    string `json:"last_commit_id,omitempty" jsonschema_description:"并发保护 commit id"`
}
