package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	BaseURL      string
	Project      string
	Token        string
	Insecure     bool
	Transport    string
	HTTPAddr     string
	HTTPPath     string
	MCPAuthToken string
	Debug        bool
}

func FromEnv() (Config, error) {
	cfg := Config{
		BaseURL:      strings.TrimSpace(os.Getenv("GITLAB_BASE_URL")),
		Project:      strings.TrimSpace(os.Getenv("GITLAB_PROJECT")),
		Token:        strings.TrimSpace(os.Getenv("GITLAB_TOKEN")),
		Insecure:     parseBoolEnv(os.Getenv("GITLAB_INSECURE"), true),
		Transport:    strings.TrimSpace(strings.ToLower(os.Getenv("MCP_TRANSPORT"))),
		HTTPAddr:     strings.TrimSpace(os.Getenv("MCP_HTTP_ADDR")),
		HTTPPath:     strings.TrimSpace(os.Getenv("MCP_HTTP_PATH")),
		MCPAuthToken: strings.TrimSpace(os.Getenv("MCP_AUTH_TOKEN")),
		Debug:        parseBoolEnv(os.Getenv("MCP_DEBUG"), false),
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://gitlab.example.com"
	}
	if cfg.Transport == "" {
		cfg.Transport = "stdio"
	}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if cfg.HTTPPath == "" {
		cfg.HTTPPath = "/mcp"
	}
	switch cfg.Transport {
	case "stdio":
		if cfg.Token == "" {
			return Config{}, errors.New("GITLAB_TOKEN is required")
		}
	case "http":
		if cfg.MCPAuthToken == "" {
			return Config{}, errors.New("MCP_AUTH_TOKEN is required when MCP_TRANSPORT=http")
		}
	default:
		return Config{}, errors.New("MCP_TRANSPORT must be stdio or http")
	}
	return cfg, nil
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
