package mcpserver

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"

	"gitlab-mcp/internal/gitlab"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	Name    = "gitlab-mcp"
	Version = "0.3.0"
)

type ClientProvider interface {
	Client(req mcp.CallToolRequest) (*gitlab.Client, error)
	TokenFingerprint(req mcp.CallToolRequest) string
	DefaultProject(req mcp.CallToolRequest) string
}

type StaticClientProvider struct {
	ClientValue *gitlab.Client
}

func (p StaticClientProvider) Client(req mcp.CallToolRequest) (*gitlab.Client, error) {
	if p.ClientValue == nil {
		return nil, errors.New("gitlab client is not configured")
	}
	return p.ClientValue, nil
}

func (p StaticClientProvider) TokenFingerprint(req mcp.CallToolRequest) string {
	return "static"
}

func (p StaticClientProvider) DefaultProject(req mcp.CallToolRequest) string {
	return ""
}

type HeaderClientProvider struct {
	BaseURL  string
	Insecure bool
}

func (p HeaderClientProvider) Client(req mcp.CallToolRequest) (*gitlab.Client, error) {
	token := strings.TrimSpace(req.Header.Get("X-GitLab-Token"))
	if token == "" {
		return nil, errors.New("X-GitLab-Token header is required")
	}
	return gitlab.NewClient(p.BaseURL, token, p.Insecure)
}

func (p HeaderClientProvider) TokenFingerprint(req mcp.CallToolRequest) string {
	token := strings.TrimSpace(req.Header.Get("X-GitLab-Token"))
	if token == "" {
		return "missing"
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])[:12]
}

func (p HeaderClientProvider) DefaultProject(req mcp.CallToolRequest) string {
	return strings.TrimSpace(req.Header.Get("X-GitLab-Project"))
}

var debugEnabled bool

func SetDebug(enabled bool) {
	debugEnabled = enabled
}

func New(provider ClientProvider, defaultProject string) *server.MCPServer {
	srv := server.NewMCPServer(Name, Version)
	registerTools(srv, provider, defaultProject)
	return srv
}

func registerTools(srv *server.MCPServer, provider ClientProvider, defaultProject string) {
	registerUserTools(srv, provider)
	registerProjectTools(srv, provider, defaultProject)
	registerIssueTools(srv, provider, defaultProject)
	registerMRTools(srv, provider, defaultProject)
	registerPipelineTools(srv, provider, defaultProject)
	registerRepositoryTools(srv, provider, defaultProject)
}

func debugLog(provider ClientProvider, req mcp.CallToolRequest, msg string, kv ...any) {
	if !debugEnabled {
		return
	}
	fields := []any{"tool", req.Params.Name, "token_fp", provider.TokenFingerprint(req), "msg", msg}
	fields = append(fields, kv...)
	log.Println(fields...)
}

func listResult(items []map[string]any) ListResult {
	return ListResult{Items: items, Count: len(items)}
}

func projectOrDefault(req mcp.CallToolRequest, provider ClientProvider, project, fallback string) (string, error) {
	project = strings.TrimSpace(project)
	if project != "" {
		return project, nil
	}
	project = strings.TrimSpace(provider.DefaultProject(req))
	if project != "" {
		return project, nil
	}
	fallback = strings.TrimSpace(fallback)
	if fallback == "" {
		return "", errors.New("project is required")
	}
	return fallback, nil
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
