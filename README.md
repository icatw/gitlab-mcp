# gitlab-mcp

一个基于 `mcp-go` 的 GitLab MCP Server（stdio），用于在 AI/Agent 中调用 GitLab 常用能力。

## 功能概览

当前支持 25 个工具，覆盖：

- 项目：查询项目信息
- Issue：查询、创建、更新、评论查询、评论创建
- Merge Request：查询、创建、更新、审批、合并、变更查询、评论查询、评论创建
- Pipeline：列表、详情、重试
- 仓库：分支列表/创建、提交列表/详情、文件读取/更新

## 环境要求

- Go 1.22+
- 可访问的 GitLab 实例
- 具备 `api` scope 的 GitLab Token

## 配置（全部使用环境变量）

必填：

- `GITLAB_TOKEN`：GitLab 访问令牌

可选：

- `GITLAB_BASE_URL`：GitLab 根地址，默认 `https://gitlab.example.com`
- `GITLAB_PROJECT`：默认项目路径（`group/project`）
- `GITLAB_INSECURE`：是否跳过 TLS 证书校验，默认 `true`

示例：

```bash
export GITLAB_TOKEN='your_gitlab_token'
export GITLAB_BASE_URL='https://gitlab.example.com'
export GITLAB_PROJECT='group/project'
export GITLAB_INSECURE='true'
```

## 运行

```bash
go run ./cmd/mcp-gitlab
```

说明：若你修改了代码（新增/修改工具、参数、逻辑），需要重新执行 `go run` 或重新编译二进制后再启动，MCP 才会加载最新实现。

推荐先编译到仓库内固定路径再运行（便于 Codex 配置）：

```bash
mkdir -p bin
go build -o ./bin/server ./cmd/mcp-gitlab
./bin/server
```

Codex MCP 配置示例：

```toml
[mcp_servers.gitlab]
type = "stdio"
command = "/Users/cyberserval/GS-go/gitlab-mcp/bin/server"
enabled = true

[mcp_servers.gitlab.env]
GITLAB_TOKEN = "your_gitlab_token"
GITLAB_BASE_URL = "https://gitlab.example.com"
GITLAB_PROJECT = "group/project"
GITLAB_INSECURE = "true"
```

## MCP 工具清单

- `project_get`
- `issues_list`
- `issues_get`
- `issues_create`
- `issues_update`
- `issues_notes`
- `issues_notes_create`
- `mrs_list`
- `mrs_get`
- `mrs_create`
- `mrs_update`
- `mrs_approve`
- `mrs_merge`
- `mrs_changes`
- `mrs_notes`
- `mrs_notes_create`
- `pipelines_list`
- `pipeline_get`
- `pipeline_retry`
- `branches_list`
- `branches_create`
- `commits_list`
- `commit_get`
- `repository_file_get`
- `repository_file_update`

## 最小调用示例（JSON 入参）

### 1) 创建 MR：`mrs_create`

```json
{
  "source_branch": "feature/demo",
  "target_branch": "main",
  "title": "feat: add demo",
  "description": "新增 demo 功能",
  "remove_source_branch": true,
  "squash": true
}
```

### 2) 更新 MR：`mrs_update`

```json
{
  "iid": 123,
  "title": "feat: add demo (updated)",
  "add_labels": "backend,review",
  "state_event": "reopen"
}
```

### 3) 审批 MR：`mrs_approve`

```json
{
  "iid": 123
}
```

### 4) 合并 MR：`mrs_merge`

```json
{
  "iid": 123,
  "merge_when_pipeline_succeeds": true,
  "should_remove_source_branch": true,
  "squash": true
}
```

### 5) 给 MR 写评论：`mrs_notes_create`

```json
{
  "iid": 123,
  "body": "已完成检查，建议合并。"
}
```

### 6) 给 Issue 写评论：`issues_notes_create`

```json
{
  "iid": 456,
  "body": "已复现，正在修复。"
}
```

### 7) 更新 Issue：`issues_update`

```json
{
  "iid": 456,
  "state_event": "close",
  "add_labels": "done",
  "assignee_id": 2364
}
```

### 8) 查询 Pipeline 列表：`pipelines_list`

```json
{
  "ref": "main",
  "status": "failed",
  "per_page": 20
}
```

### 9) 查询单个 Pipeline：`pipeline_get`

```json
{
  "id": 78901
}
```

### 10) 重试 Pipeline：`pipeline_retry`

```json
{
  "id": 78901
}
```

### 11) 创建分支：`branches_create`

```json
{
  "branch": "feature/new-api",
  "ref": "main"
}
```

### 12) 查询提交：`commits_list`

```json
{
  "ref_name": "main",
  "per_page": 20
}
```

### 13) 查询单个提交：`commit_get`

```json
{
  "sha": "a1b2c3d4e5f6"
}
```

### 14) 读取仓库文件：`repository_file_get`

```json
{
  "file_path": "README.md",
  "ref": "main"
}
```

> 返回内容中的 `content` 为 base64 编码。

### 15) 更新仓库文件：`repository_file_update`

```json
{
  "file_path": "README.md",
  "branch": "feature/update-readme",
  "content": "# new content",
  "commit_message": "docs: update readme",
  "encoding": "text"
}
```

## 常见问题

### 1) 返回 401 `invalid_token`

通常是 token 过期、scope 不足或未正确设置 `GITLAB_TOKEN`。

### 2) 内网地址连通失败（状态码 000）

请检查是否被本地代理劫持；必要时对内网域名设置 `NO_PROXY`。

### 3) `project is required`

说明调用时未传 `project`，且也未设置 `GITLAB_PROJECT`。

## 开发

```bash
gofmt -w cmd/mcp-gitlab/main.go internal/gitlab/client.go
go build ./...
```
