package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"gitlab-mcp/internal/config"
	"gitlab-mcp/internal/gitlab"
	"gitlab-mcp/internal/mcpserver"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg, err := config.FromEnv()
	exitOnErr(err)
	mcpserver.SetDebug(cfg.Debug)

	switch cfg.Transport {
	case "stdio":
		exitOnErr(serveStdio(cfg))
	case "http":
		exitOnErr(serveHTTP(cfg))
	default:
		exitOnErr(fmt.Errorf("unsupported MCP_TRANSPORT: %s", cfg.Transport))
	}
}

func serveStdio(cfg config.Config) error {
	client, err := gitlab.NewClient(cfg.BaseURL, cfg.Token, cfg.Insecure)
	if err != nil {
		return err
	}
	provider := mcpserver.StaticClientProvider{ClientValue: client}
	return server.ServeStdio(mcpserver.New(provider, cfg.Project))
}

func serveHTTP(cfg config.Config) error {
	provider := mcpserver.HeaderClientProvider{
		BaseURL:  cfg.BaseURL,
		Insecure: cfg.Insecure,
	}
	mcpHTTP := server.NewStreamableHTTPServer(
		mcpserver.New(provider, cfg.Project),
		server.WithEndpointPath(cfg.HTTPPath),
		server.WithStateLess(true),
	)
	return (&http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: authMiddleware(cfg.MCPAuthToken, cfg.HTTPPath, mcpHTTP),
	}).ListenAndServe()
}

func authMiddleware(authToken, path string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		if !validBearer(r.Header.Get("Authorization"), authToken) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func validBearer(header, token string) bool {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return false
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix)) == token
}

func exitOnErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
