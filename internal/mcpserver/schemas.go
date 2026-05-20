package mcpserver

import "github.com/mark3labs/mcp-go/mcp"

func newTool(name, description string, input []mcp.ToolOption, extra ...mcp.ToolOption) mcp.Tool {
	opts := []mcp.ToolOption{mcp.WithDescription(description)}
	opts = append(opts, input...)
	opts = append(opts, extra...)
	return mcp.NewTool(name, opts...)
}

func noArgs() []mcp.ToolOption {
	return []mcp.ToolOption{}
}

func projectArg() mcp.ToolOption {
	return mcp.WithString("project", mcp.Description("项目路径，例如 group/project；优先级高于 X-GitLab-Project 和服务端 GITLAB_PROJECT"))
}

func iidArg(desc string) mcp.ToolOption {
	return mcp.WithInteger("iid", mcp.Required(), mcp.Min(1), mcp.Description(desc))
}

func paginationArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithInteger("per_page", mcp.Min(1), mcp.Max(100), mcp.Description("分页大小，1-100")),
		mcp.WithInteger("page", mcp.Min(1), mcp.Description("页码")),
		mcp.WithBoolean("all", mcp.Description("是否拉取全部分页")),
	}
}

func noteListArgs() []mcp.ToolOption {
	return append([]mcp.ToolOption{
		projectArg(),
		iidArg("Issue/MR IID"),
	}, paginationArgs()...)
}

func noteCreateArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		iidArg("Issue/MR IID"),
		mcp.WithString("body", mcp.Required(), mcp.MinLength(1), mcp.Description("评论内容（Markdown）")),
	}
}

func issueListArgs() []mcp.ToolOption {
	return append([]mcp.ToolOption{
		projectArg(),
		mcp.WithString("state", mcp.Enum("opened", "closed", "all"), mcp.Description("Issue 状态")),
		mcp.WithString("labels", mcp.Description("逗号分隔 labels")),
		mcp.WithInteger("author_id", mcp.Description("作者 ID")),
		mcp.WithInteger("assignee_id", mcp.Description("负责人 ID")),
		mcp.WithString("created_after", mcp.Description("ISO 时间，如 2024-01-01T00:00:00Z")),
		mcp.WithString("updated_after", mcp.Description("ISO 时间，如 2024-01-01T00:00:00Z")),
		mcp.WithString("search", mcp.Description("搜索关键词")),
	}, paginationArgs()...)
}

func issueGetArgs() []mcp.ToolOption {
	return []mcp.ToolOption{projectArg(), iidArg("Issue IID")}
}

func issueCreateArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		mcp.WithString("title", mcp.Required(), mcp.MinLength(1), mcp.Description("Issue 标题")),
		mcp.WithString("description", mcp.Description("Issue 描述（Markdown）")),
		mcp.WithString("labels", mcp.Description("逗号分隔 labels")),
		mcp.WithInteger("assignee_id", mcp.Description("负责人 ID")),
		mcp.WithInteger("milestone_id", mcp.Description("里程碑 ID")),
		mcp.WithString("due_date", mcp.Description("截止日期，格式 YYYY-MM-DD")),
		mcp.WithString("issue_type", mcp.Enum("issue", "incident", "test_case", "task"), mcp.Description("Issue 类型")),
		mcp.WithBoolean("confidential", mcp.Description("是否保密 Issue")),
	}
}

func issueUpdateArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		iidArg("Issue IID"),
		mcp.WithString("title", mcp.Description("Issue 标题")),
		mcp.WithString("description", mcp.Description("Issue 描述（Markdown）")),
		mcp.WithString("labels", mcp.Description("逗号分隔 labels")),
		mcp.WithString("add_labels", mcp.Description("追加 labels（逗号分隔）")),
		mcp.WithString("remove_labels", mcp.Description("移除 labels（逗号分隔）")),
		mcp.WithInteger("assignee_id", mcp.Description("负责人 ID")),
		mcp.WithInteger("milestone_id", mcp.Description("里程碑 ID")),
		mcp.WithString("due_date", mcp.Description("截止日期，格式 YYYY-MM-DD")),
		mcp.WithString("state_event", mcp.Enum("close", "reopen"), mcp.Description("状态流转")),
		mcp.WithBoolean("confidential", mcp.Description("是否保密")),
		mcp.WithBoolean("discussion_locked", mcp.Description("是否锁定讨论")),
	}
}

func mrListArgs() []mcp.ToolOption {
	return append([]mcp.ToolOption{
		projectArg(),
		mcp.WithString("state", mcp.Enum("opened", "closed", "merged", "all"), mcp.Description("MR 状态")),
		mcp.WithString("labels", mcp.Description("逗号分隔 labels")),
		mcp.WithInteger("author_id", mcp.Description("作者 ID")),
		mcp.WithInteger("assignee_id", mcp.Description("负责人 ID")),
		mcp.WithString("created_after", mcp.Description("ISO 时间，如 2024-01-01T00:00:00Z")),
		mcp.WithString("updated_after", mcp.Description("ISO 时间，如 2024-01-01T00:00:00Z")),
		mcp.WithString("search", mcp.Description("搜索关键词")),
	}, paginationArgs()...)
}

func mrGetArgs() []mcp.ToolOption {
	return []mcp.ToolOption{projectArg(), iidArg("MR IID")}
}

func mrCreateArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		mcp.WithString("source_branch", mcp.Required(), mcp.MinLength(1), mcp.Description("源分支")),
		mcp.WithString("target_branch", mcp.Required(), mcp.MinLength(1), mcp.Description("目标分支")),
		mcp.WithString("title", mcp.Required(), mcp.MinLength(1), mcp.Description("MR 标题")),
		mcp.WithString("description", mcp.Description("MR 描述（Markdown）")),
		mcp.WithInteger("assignee_id", mcp.Description("负责人 ID")),
		mcp.WithString("labels", mcp.Description("逗号分隔 labels")),
		mcp.WithBoolean("remove_source_branch", mcp.Description("合并后删除源分支")),
		mcp.WithBoolean("squash", mcp.Description("是否 squash")),
		mcp.WithBoolean("draft", mcp.Description("是否草稿")),
		mcp.WithBoolean("allow_collaboration", mcp.Description("是否允许协作者 push")),
	}
}

func mrUpdateArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		iidArg("MR IID"),
		mcp.WithString("target_branch", mcp.Description("目标分支")),
		mcp.WithString("title", mcp.Description("MR 标题")),
		mcp.WithString("description", mcp.Description("MR 描述（Markdown）")),
		mcp.WithInteger("assignee_id", mcp.Description("负责人 ID")),
		mcp.WithString("labels", mcp.Description("逗号分隔 labels")),
		mcp.WithString("add_labels", mcp.Description("追加 labels（逗号分隔）")),
		mcp.WithString("remove_labels", mcp.Description("移除 labels（逗号分隔）")),
		mcp.WithString("state_event", mcp.Enum("close", "reopen"), mcp.Description("状态流转")),
		mcp.WithBoolean("squash", mcp.Description("是否 squash")),
		mcp.WithBoolean("discussion_locked", mcp.Description("是否锁定讨论")),
		mcp.WithBoolean("remove_source_branch", mcp.Description("合并后删除源分支")),
	}
}

func mrApproveArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		iidArg("MR IID"),
		mcp.WithString("sha", mcp.Description("指定当前 HEAD SHA，防止过期审批")),
		mcp.WithString("approval_password", mcp.Description("审批密码（若启用）")),
	}
}

func mrMergeArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		iidArg("MR IID"),
		mcp.WithString("merge_commit_message", mcp.Description("合并提交信息")),
		mcp.WithString("squash_commit_message", mcp.Description("squash 提交信息")),
		mcp.WithBoolean("should_remove_source_branch", mcp.Description("是否删除源分支")),
		mcp.WithBoolean("merge_when_pipeline_succeeds", mcp.Description("流水线通过后自动合并")),
		mcp.WithString("sha", mcp.Description("指定当前 HEAD SHA，防止过期合并")),
		mcp.WithBoolean("squash", mcp.Description("是否 squash")),
	}
}

func pipelineListArgs() []mcp.ToolOption {
	return append([]mcp.ToolOption{
		projectArg(),
		mcp.WithString("ref", mcp.Description("分支或 tag")),
		mcp.WithString("status", mcp.Enum("created", "waiting_for_resource", "preparing", "pending", "running", "success", "failed", "canceled", "skipped", "manual", "scheduled"), mcp.Description("Pipeline 状态")),
		mcp.WithString("scope", mcp.Enum("running", "pending", "finished", "branches", "tags"), mcp.Description("Pipeline 范围")),
		mcp.WithString("source", mcp.Description("Pipeline 来源")),
		mcp.WithString("username", mcp.Description("触发者用户名")),
		mcp.WithString("updated_after", mcp.Description("ISO 时间，如 2024-01-01T00:00:00Z")),
		mcp.WithString("updated_before", mcp.Description("ISO 时间，如 2024-01-31T23:59:59Z")),
	}, paginationArgs()...)
}

func pipelineGetArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		mcp.WithInteger("id", mcp.Required(), mcp.Min(1), mcp.Description("Pipeline ID")),
	}
}

func branchListArgs() []mcp.ToolOption {
	return append([]mcp.ToolOption{
		projectArg(),
		mcp.WithString("search", mcp.Description("分支搜索关键词")),
		mcp.WithString("regex", mcp.Description("分支正则（RE2）")),
	}, paginationArgs()...)
}

func branchCreateArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		mcp.WithString("branch", mcp.Required(), mcp.MinLength(1), mcp.Description("新分支名")),
		mcp.WithString("ref", mcp.Required(), mcp.MinLength(1), mcp.Description("源 ref（分支/标签/提交）")),
	}
}

func commitListArgs() []mcp.ToolOption {
	return append([]mcp.ToolOption{
		projectArg(),
		mcp.WithString("ref_name", mcp.Description("分支或 tag")),
		mcp.WithString("path", mcp.Description("文件路径过滤")),
		mcp.WithString("since", mcp.Description("ISO 时间，如 2024-01-01T00:00:00Z")),
		mcp.WithString("until", mcp.Description("ISO 时间，如 2024-01-31T23:59:59Z")),
	}, paginationArgs()...)
}

func commitGetArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		mcp.WithString("sha", mcp.Required(), mcp.MinLength(1), mcp.Description("提交 SHA")),
	}
}

func repoFileGetArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		mcp.WithString("file_path", mcp.Required(), mcp.MinLength(1), mcp.Description("仓库文件路径")),
		mcp.WithString("ref", mcp.Required(), mcp.MinLength(1), mcp.Description("分支、tag 或 commit SHA")),
	}
}

func repoFileUpdateArgs() []mcp.ToolOption {
	return []mcp.ToolOption{
		projectArg(),
		mcp.WithString("file_path", mcp.Required(), mcp.MinLength(1), mcp.Description("仓库文件路径")),
		mcp.WithString("branch", mcp.Required(), mcp.MinLength(1), mcp.Description("目标分支")),
		mcp.WithString("content", mcp.Required(), mcp.Description("新的文件内容")),
		mcp.WithString("commit_message", mcp.Required(), mcp.MinLength(1), mcp.Description("提交信息")),
		mcp.WithString("author_email", mcp.Description("提交作者邮箱")),
		mcp.WithString("author_name", mcp.Description("提交作者名")),
		mcp.WithString("encoding", mcp.Enum("text", "base64"), mcp.Description("内容编码")),
		mcp.WithBoolean("execute_filemode", mcp.Description("是否设置可执行位")),
		mcp.WithString("last_commit_id", mcp.Description("并发保护 commit id")),
	}
}
