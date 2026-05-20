# gitlab-mcp

[English](./README.en.md)

`gitlab-mcp` 是一个面向 AI 编程代理的 GitLab MCP（Model Context Protocol）服务。它把常见 GitLab 工作流封装成 MCP tools，覆盖 Issue、Merge Request、Pipeline、分支、提交和仓库文件等操作。

项目同时支持本地 `stdio` 和远程 Streamable HTTP 两种传输方式：可以作为个人本地 MCP 工具使用，也可以部署成团队共享服务。HTTP 模式下通过 `X-GitLab-Token` 传递每个使用者自己的 GitLab token，避免多人共用同一个 GitLab 身份。

## 特性

- 提供 25 个 GitLab MCP tools，覆盖项目、Issue、MR、Pipeline、分支、提交和仓库文件
- 支持本地 `stdio` 传输，适配 Codex、Claude Code、Cursor 等 MCP 客户端
- 支持远程 Streamable HTTP 传输，方便部署到服务器
- HTTP 模式通过 `X-GitLab-Token` 隔离每个使用者的 GitLab 权限
- HTTP 模式通过 `Authorization: Bearer <MCP_AUTH_TOKEN>` 保护 MCP 服务入口
- 全部配置使用环境变量，不需要配置文件，也不内置任何密钥

## 适用场景

- 让 AI 编程代理读取 GitLab Issue、MR 和 Pipeline 状态
- 从 MCP 客户端创建 Issue、评论、分支和 MR
- 部署一个团队共享 MCP endpoint，同时保留每个用户自己的 GitLab 权限和审计记录
- 在官方 GitLab MCP 不可用或不满足需求时，对接自托管 GitLab 实例

## 功能概览

当前支持 25 个工具，覆盖：

- 项目：查询项目信息
- Issue：查询、创建、更新、评论查询、评论创建
- Merge Request：查询、创建、更新、审批、合并、变更查询、评论查询、评论创建
- Pipeline：列表、详情、重试
- 仓库：分支列表/创建、提交列表/详情、文件读取/更新

## 环境要求

- Go 1.25+
- 可访问的 GitLab 实例
- 具备 `api` scope 的 GitLab Token

## 配置

通用可选环境变量：

- `GITLAB_BASE_URL`：GitLab 根地址，默认 `https://gitlab.example.com`
- `GITLAB_PROJECT`：默认项目路径（`group/project`），工具入参里的 `project` 可覆盖该默认值
- `GITLAB_INSECURE`：是否跳过 TLS 证书校验，默认 `true`
- `MCP_TRANSPORT`：`stdio` 或 `http`，默认 `stdio`

### 本地 stdio

必填：

- `GITLAB_TOKEN`：GitLab 访问令牌

示例：

```bash
export MCP_TRANSPORT='stdio'
export GITLAB_TOKEN='your_gitlab_token'
export GITLAB_BASE_URL='https://gitlab.example.com'
export GITLAB_PROJECT='group/project'
export GITLAB_INSECURE='true'
```

Codex stdio 配置示例：

```toml
[mcp_servers.gitlab]
type = "stdio"
command = "/path/to/gitlab-mcp/bin/server"
enabled = true

[mcp_servers.gitlab.env]
MCP_TRANSPORT = "stdio"
GITLAB_TOKEN = "your_gitlab_token"
GITLAB_BASE_URL = "https://gitlab.example.com"
GITLAB_PROJECT = "group/project"
GITLAB_INSECURE = "true"
```

### 远程 HTTP

服务端必填：

- `MCP_AUTH_TOKEN`：保护 MCP 服务本身的访问令牌

服务端可选：

- `MCP_HTTP_ADDR`：HTTP 监听地址，默认 `:8080`
- `MCP_HTTP_PATH`：MCP endpoint，默认 `/mcp`

服务端启动示例：

```bash
export MCP_TRANSPORT='http'
export MCP_AUTH_TOKEN='server_access_token'
export GITLAB_BASE_URL='https://gitlab.example.com'
export GITLAB_PROJECT='group/project'
export GITLAB_INSECURE='true'
./bin/server
```

Codex 远程 HTTP 配置示例：

```toml
[mcp_servers.gitlab]
url = "https://example.com/mcp"
enabled = true

[mcp_servers.gitlab.http_headers]
Authorization = "Bearer server_access_token"
X-GitLab-Token = "your_personal_gitlab_token"
X-GitLab-Project = "group/project"
```

Trae 远程 HTTP 配置示例：

```json
{
  "mcpServers": {
    "gitlab-http": {
      "url": "https://example.com/mcp",
      "headers": {
        "Authorization": "Bearer your_mcp_auth_token",
        "X-GitLab-Token": "your_personal_gitlab_token",
        "X-GitLab-Project": "group/project"
      }
    }
  }
}
```

Cursor 远程 HTTP 配置示例：

在项目内创建 `.cursor/mcp.json`，或在用户目录创建 `~/.cursor/mcp.json`：

```json
{
  "mcpServers": {
    "gitlab-http": {
      "type": "http",
      "url": "https://example.com/mcp",
      "headers": {
        "Authorization": "Bearer your_mcp_auth_token",
        "X-GitLab-Token": "your_personal_gitlab_token",
        "X-GitLab-Project": "group/project"
      }
    }
  }
}
```

Claude Code 远程 HTTP 配置示例：

```bash
claude mcp add-json gitlab-http '{
  "type": "http",
  "url": "https://example.com/mcp",
  "headers": {
    "Authorization": "Bearer your_mcp_auth_token",
    "X-GitLab-Token": "your_personal_gitlab_token",
    "X-GitLab-Project": "group/project"
  }
}'
```

也可以把同样的 JSON 放到项目级 `.mcp.json` 中，便于团队共享配置。注意不要把真实 token 提交到仓库；团队共享配置建议只提交 URL，个人 token 由每个使用者在自己的本地配置里维护。

HTTP 模式下，服务端不读取使用者本机环境变量里的 `GITLAB_TOKEN`。每个使用者必须通过 `X-GitLab-Token` 请求头传自己的 GitLab token，这样多人使用时权限和 GitLab 审计不会混在一起。

`X-GitLab-Project` 是客户端级默认项目，适合 Trae、Cursor、Claude Code 这类客户端减少每次调用时重复传项目。项目选择优先级为：工具入参 `project` > 请求头 `X-GitLab-Project` > 服务端环境变量 `GITLAB_PROJECT`。

HTTP 服务使用 Streamable HTTP 的 stateless 模式，每次请求都应携带 `Authorization` 和 `X-GitLab-Token` 请求头。

### Docker 部署

构建镜像：

```bash
docker build -t gitlab-mcp:local .
```

直接运行：

```bash
docker run -d --name gitlab-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT=http \
  -e MCP_HTTP_ADDR=0.0.0.0:8080 \
  -e MCP_HTTP_PATH=/mcp \
  -e MCP_AUTH_TOKEN='server_access_token' \
  -e GITLAB_BASE_URL='https://gitlab.example.com' \
  -e GITLAB_PROJECT='group/project' \
  -e GITLAB_INSECURE='true' \
  gitlab-mcp:local
```

使用 Docker Compose：

```bash
export MCP_AUTH_TOKEN='server_access_token'
export GITLAB_BASE_URL='https://gitlab.example.com'
export GITLAB_PROJECT='group/project'
export GITLAB_INSECURE='true'

docker compose up -d --build
```

`GITLAB_PROJECT` 只是默认项目。所有工具都支持传入 `project` 参数，例如 `{"project":"other-group/other-project"}`。如果团队会访问多个项目，可以把 `GITLAB_PROJECT` 留空，让每次工具调用都显式传 `project`；如果团队主要围绕一个项目工作，则建议设置 `GITLAB_PROJECT` 作为默认值，减少日常调用成本。

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

通常是 GitLab token 过期、scope 不足、stdio 模式未正确设置 `GITLAB_TOKEN`，或 HTTP 模式未正确传 `X-GitLab-Token`。

### 2) 内网地址连通失败（状态码 000）

请检查是否被本地代理劫持；必要时对内网域名设置 `NO_PROXY`。

### 3) `project is required`

说明调用时未传 `project`，且也未设置 `GITLAB_PROJECT`。

## 开发

```bash
gofmt -w cmd/mcp-gitlab/main.go internal/**/*.go
go build ./...
```
