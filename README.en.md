# gitlab-mcp

[中文文档](./README.md)

A GitLab Model Context Protocol (MCP) server for AI coding agents.

`gitlab-mcp` exposes common GitLab workflows as MCP tools, including issues, merge requests, pipelines, branches, commits, and repository files. It supports both local `stdio` usage and remote Streamable HTTP deployment, so it can run as a personal local tool or as a shared server for a team.

## Features

- 25 GitLab tools for project, issue, MR, pipeline, branch, commit, and repository file workflows
- Local `stdio` transport for Codex, Claude Code, Cursor, and other MCP clients
- Remote Streamable HTTP transport for server deployment
- Per-user GitLab token isolation in HTTP mode via `X-GitLab-Token`
- Server-level HTTP access protection via `Authorization: Bearer <MCP_AUTH_TOKEN>`
- Environment-variable based configuration with no config file or embedded secret required

## Use Cases

- Let an AI coding agent inspect GitLab issues, merge requests, and pipeline status
- Create issues, comments, branches, and merge requests from an MCP client
- Deploy one shared MCP endpoint while keeping each user's GitLab permissions separate
- Use a self-hosted GitLab instance without relying on GitLab's official MCP availability

## Tool Coverage

Currently supports 25 tools:

- Project: get project metadata
- Issue: list, get, create, update, list notes, create notes
- Merge Request: list, get, create, update, approve, merge, get changes, list notes, create notes
- Pipeline: list, get, retry
- Repository: list/create branches, list/get commits, get/update repository files

## Requirements

- Go 1.25+
- A reachable GitLab instance
- A GitLab token with `api` scope

## Configuration

Common optional environment variables:

- `GITLAB_BASE_URL`: GitLab base URL, defaults to `https://gitlab.example.com`
- `GITLAB_PROJECT`: default project path, such as `group/project`; tool-level `project` arguments override it
- `GITLAB_INSECURE`: skip TLS certificate verification, defaults to `true`
- `MCP_TRANSPORT`: `stdio` or `http`, defaults to `stdio`

### Local stdio

Required:

- `GITLAB_TOKEN`: GitLab access token

Example:

```bash
export MCP_TRANSPORT='stdio'
export GITLAB_TOKEN='your_gitlab_token'
export GITLAB_BASE_URL='https://gitlab.example.com'
export GITLAB_PROJECT='group/project'
export GITLAB_INSECURE='true'
```

Codex stdio configuration:

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

### Remote HTTP

Server required:

- `MCP_AUTH_TOKEN`: access token that protects the MCP server itself

Server optional:

- `MCP_HTTP_ADDR`: HTTP listen address, defaults to `:8080`
- `MCP_HTTP_PATH`: MCP endpoint path, defaults to `/mcp`

Server startup example:

```bash
export MCP_TRANSPORT='http'
export MCP_AUTH_TOKEN='server_access_token'
export GITLAB_BASE_URL='https://gitlab.example.com'
export GITLAB_PROJECT='group/project'
export GITLAB_INSECURE='true'
./bin/server
```

Codex remote HTTP configuration:

```toml
[mcp_servers.gitlab]
url = "https://example.com/mcp"
enabled = true

[mcp_servers.gitlab.http_headers]
Authorization = "Bearer server_access_token"
X-GitLab-Token = "your_personal_gitlab_token"
X-GitLab-Project = "group/project"
```

Trae remote HTTP configuration:

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

Cursor remote HTTP configuration:

Create `.cursor/mcp.json` in a project, or `~/.cursor/mcp.json` for global configuration:

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

Claude Code remote HTTP configuration:

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

You can also put the same JSON under a project-level `.mcp.json` for team-shared configuration. Do not commit real tokens to the repository; for team-shared config, commit only the URL and let each user manage personal tokens locally.

In HTTP mode, the server does not read `GITLAB_TOKEN` from each user's local environment. Each user must send their own GitLab token through `X-GitLab-Token`, keeping GitLab permissions and audit logs separated.

`X-GitLab-Project` is the client-level default project. It is useful for clients such as Trae, Cursor, and Claude Code because users do not need to repeat the project in every tool call. Project resolution priority is: tool argument `project` > request header `X-GitLab-Project` > server environment variable `GITLAB_PROJECT`.

The HTTP server uses Streamable HTTP stateless mode. Every request should include both `Authorization` and `X-GitLab-Token` headers.

### Docker Deployment

Build the image:

```bash
docker build -t gitlab-mcp:local .
```

Run directly:

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

Use Docker Compose:

```bash
export MCP_AUTH_TOKEN='server_access_token'
export GITLAB_BASE_URL='https://gitlab.example.com'
export GITLAB_PROJECT='group/project'
export GITLAB_INSECURE='true'

docker compose up -d --build
```

`GITLAB_PROJECT` is only the default project. Every tool supports a `project` argument, such as `{"project":"other-group/other-project"}`. For multi-project teams, you can leave `GITLAB_PROJECT` empty and pass `project` explicitly on each tool call. For teams that mostly work in one project, setting `GITLAB_PROJECT` reduces repetitive input.

## Run

```bash
go run ./cmd/mcp-gitlab
```

After changing code, restart `go run` or rebuild the binary so the MCP client loads the latest implementation.

Recommended build path:

```bash
mkdir -p bin
go build -o ./bin/server ./cmd/mcp-gitlab
./bin/server
```

## Tool List

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

## Minimal JSON Arguments

### Create MR: `mrs_create`

```json
{
  "source_branch": "feature/demo",
  "target_branch": "main",
  "title": "feat: add demo",
  "description": "Add demo feature",
  "remove_source_branch": true,
  "squash": true
}
```

### Update MR: `mrs_update`

```json
{
  "iid": 123,
  "title": "feat: add demo (updated)",
  "add_labels": "backend,review",
  "state_event": "reopen"
}
```

### Approve MR: `mrs_approve`

```json
{
  "iid": 123
}
```

### Merge MR: `mrs_merge`

```json
{
  "iid": 123,
  "merge_when_pipeline_succeeds": true,
  "should_remove_source_branch": true,
  "squash": true
}
```

### Create MR note: `mrs_notes_create`

```json
{
  "iid": 123,
  "body": "Checked and ready to merge."
}
```

### Create Issue note: `issues_notes_create`

```json
{
  "iid": 456,
  "body": "Reproduced and working on a fix."
}
```

### Update Issue: `issues_update`

```json
{
  "iid": 456,
  "state_event": "close",
  "add_labels": "done",
  "assignee_id": 2364
}
```

### List pipelines: `pipelines_list`

```json
{
  "ref": "main",
  "status": "failed",
  "per_page": 20
}
```

### Get pipeline: `pipeline_get`

```json
{
  "id": 78901
}
```

### Retry pipeline: `pipeline_retry`

```json
{
  "id": 78901
}
```

### Create branch: `branches_create`

```json
{
  "branch": "feature/new-api",
  "ref": "main"
}
```

### List commits: `commits_list`

```json
{
  "ref_name": "main",
  "per_page": 20
}
```

### Get commit: `commit_get`

```json
{
  "sha": "a1b2c3d4e5f6"
}
```

### Get repository file: `repository_file_get`

```json
{
  "file_path": "README.md",
  "ref": "main"
}
```

The returned `content` is base64 encoded.

### Update repository file: `repository_file_update`

```json
{
  "file_path": "README.md",
  "branch": "feature/update-readme",
  "content": "# new content",
  "commit_message": "docs: update readme",
  "encoding": "text"
}
```

## FAQ

### 401 `invalid_token`

Usually caused by an expired GitLab token, missing `api` scope, missing `GITLAB_TOKEN` in stdio mode, or missing `X-GitLab-Token` in HTTP mode.

### Network status code 000

Check whether local proxy settings intercept your GitLab domain. Configure `NO_PROXY` when needed.

### `project is required`

Pass `project` in the tool call or set `GITLAB_PROJECT`.

## Development

```bash
gofmt -w cmd/mcp-gitlab/main.go internal/**/*.go
go build ./...
```
