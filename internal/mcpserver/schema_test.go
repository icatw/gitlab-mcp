package mcpserver

import (
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestToolSchemasExposeProject(t *testing.T) {
	srv := New(StaticClientProvider{}, "")
	tools := srv.ListTools()
	if len(tools) != 26 {
		t.Fatalf("expected 26 tools, got %d", len(tools))
	}

	for name, entry := range tools {
		props := entry.Tool.InputSchema.Properties
		if name == "user_current" {
			continue
		}
		if len(props) == 0 {
			t.Fatalf("%s input schema properties is empty", name)
		}
		if _, ok := props["project"]; !ok {
			t.Fatalf("%s input schema does not expose project: %#v", name, props)
		}
	}

	mrsList := tools["mrs_list"]
	if mrsList == nil {
		t.Fatal("mrs_list tool is missing")
	}
	for _, field := range []string{"project", "state", "per_page", "page", "all"} {
		if _, ok := mrsList.Tool.InputSchema.Properties[field]; !ok {
			t.Fatalf("mrs_list input schema does not expose %s: %#v", field, mrsList.Tool.InputSchema.Properties)
		}
	}
}

func TestProjectOrDefaultPriority(t *testing.T) {
	provider := HeaderClientProvider{}
	req := mcp.CallToolRequest{Header: http.Header{}}
	req.Header.Set("X-GitLab-Project", "header/project")

	project, err := projectOrDefault(req, provider, "arg/project", "env/project")
	if err != nil {
		t.Fatal(err)
	}
	if project != "arg/project" {
		t.Fatalf("expected arg project, got %q", project)
	}

	project, err = projectOrDefault(req, provider, "", "env/project")
	if err != nil {
		t.Fatal(err)
	}
	if project != "header/project" {
		t.Fatalf("expected header project, got %q", project)
	}

	req.Header.Del("X-GitLab-Project")
	project, err = projectOrDefault(req, provider, "", "env/project")
	if err != nil {
		t.Fatal(err)
	}
	if project != "env/project" {
		t.Fatalf("expected env project, got %q", project)
	}

	_, err = projectOrDefault(req, provider, "", "")
	if err == nil {
		t.Fatal("expected missing project error")
	}
}
