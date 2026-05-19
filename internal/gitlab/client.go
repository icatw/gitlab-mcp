package gitlab

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string, insecure bool) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("base url is required")
	}
	if token == "" {
		return nil, errors.New("token is required")
	}
	baseURL = strings.TrimRight(baseURL, "/")
	transport := &http.Transport{}
	if insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second, Transport: transport},
	}, nil
}

func (c *Client) GetProject(path string) (map[string]any, error) {
	var out map[string]any
	_, err := c.do("GET", "/projects/"+url.PathEscape(path), nil, &out)
	return out, err
}

func (c *Client) GetIssue(projectPath, iid string) (map[string]any, error) {
	var out map[string]any
	_, err := c.do("GET", "/projects/"+url.PathEscape(projectPath)+"/issues/"+iid, nil, &out)
	return out, err
}

func (c *Client) CreateIssue(projectPath string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("POST", "/projects/"+url.PathEscape(projectPath)+"/issues", nil, payload, &out)
	return out, err
}

func (c *Client) UpdateIssue(projectPath, iid string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("PUT", "/projects/"+url.PathEscape(projectPath)+"/issues/"+iid, nil, payload, &out)
	return out, err
}

func (c *Client) ListIssues(projectPath string, q url.Values, all bool) ([]map[string]any, error) {
	return c.getList("/projects/"+url.PathEscape(projectPath)+"/issues", q, all)
}

func (c *Client) GetIssueNotes(projectPath, iid string, q url.Values, all bool) ([]map[string]any, error) {
	return c.getList("/projects/"+url.PathEscape(projectPath)+"/issues/"+iid+"/notes", q, all)
}

func (c *Client) GetMergeRequest(projectPath, iid string) (map[string]any, error) {
	var out map[string]any
	_, err := c.do("GET", "/projects/"+url.PathEscape(projectPath)+"/merge_requests/"+iid, nil, &out)
	return out, err
}

func (c *Client) CreateMergeRequest(projectPath string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("POST", "/projects/"+url.PathEscape(projectPath)+"/merge_requests", nil, payload, &out)
	return out, err
}

func (c *Client) UpdateMergeRequest(projectPath, iid string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("PUT", "/projects/"+url.PathEscape(projectPath)+"/merge_requests/"+iid, nil, payload, &out)
	return out, err
}

func (c *Client) ApproveMergeRequest(projectPath, iid string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("POST", "/projects/"+url.PathEscape(projectPath)+"/merge_requests/"+iid+"/approve", nil, payload, &out)
	return out, err
}

func (c *Client) MergeMergeRequest(projectPath, iid string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("PUT", "/projects/"+url.PathEscape(projectPath)+"/merge_requests/"+iid+"/merge", nil, payload, &out)
	return out, err
}

func (c *Client) GetMergeRequestChanges(projectPath, iid string) (map[string]any, error) {
	var out map[string]any
	_, err := c.do("GET", "/projects/"+url.PathEscape(projectPath)+"/merge_requests/"+iid+"/changes", nil, &out)
	return out, err
}

func (c *Client) ListMergeRequests(projectPath string, q url.Values, all bool) ([]map[string]any, error) {
	return c.getList("/projects/"+url.PathEscape(projectPath)+"/merge_requests", q, all)
}

func (c *Client) GetMergeRequestNotes(projectPath, iid string, q url.Values, all bool) ([]map[string]any, error) {
	return c.getList("/projects/"+url.PathEscape(projectPath)+"/merge_requests/"+iid+"/notes", q, all)
}

func (c *Client) CreateIssueNote(projectPath, iid string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("POST", "/projects/"+url.PathEscape(projectPath)+"/issues/"+iid+"/notes", nil, payload, &out)
	return out, err
}

func (c *Client) CreateMergeRequestNote(projectPath, iid string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("POST", "/projects/"+url.PathEscape(projectPath)+"/merge_requests/"+iid+"/notes", nil, payload, &out)
	return out, err
}

func (c *Client) ListPipelines(projectPath string, q url.Values, all bool) ([]map[string]any, error) {
	return c.getList("/projects/"+url.PathEscape(projectPath)+"/pipelines", q, all)
}

func (c *Client) GetPipeline(projectPath, id string) (map[string]any, error) {
	var out map[string]any
	_, err := c.do("GET", "/projects/"+url.PathEscape(projectPath)+"/pipelines/"+id, nil, &out)
	return out, err
}

func (c *Client) RetryPipeline(projectPath, id string) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("POST", "/projects/"+url.PathEscape(projectPath)+"/pipelines/"+id+"/retry", nil, nil, &out)
	return out, err
}

func (c *Client) ListBranches(projectPath string, q url.Values, all bool) ([]map[string]any, error) {
	return c.getList("/projects/"+url.PathEscape(projectPath)+"/repository/branches", q, all)
}

func (c *Client) CreateBranch(projectPath string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("POST", "/projects/"+url.PathEscape(projectPath)+"/repository/branches", nil, payload, &out)
	return out, err
}

func (c *Client) ListCommits(projectPath string, q url.Values, all bool) ([]map[string]any, error) {
	return c.getList("/projects/"+url.PathEscape(projectPath)+"/repository/commits", q, all)
}

func (c *Client) GetCommit(projectPath, sha string) (map[string]any, error) {
	var out map[string]any
	_, err := c.do("GET", "/projects/"+url.PathEscape(projectPath)+"/repository/commits/"+url.PathEscape(sha), nil, &out)
	return out, err
}

func (c *Client) GetRepositoryFile(projectPath, filePath, ref string) (map[string]any, error) {
	var out map[string]any
	q := url.Values{}
	q.Set("ref", ref)
	_, err := c.do("GET", "/projects/"+url.PathEscape(projectPath)+"/repository/files/"+url.PathEscape(filePath), q, &out)
	return out, err
}

func (c *Client) UpdateRepositoryFile(projectPath, filePath string, payload map[string]any) (map[string]any, error) {
	var out map[string]any
	_, err := c.doJSON("PUT", "/projects/"+url.PathEscape(projectPath)+"/repository/files/"+url.PathEscape(filePath), nil, payload, &out)
	return out, err
}

func (c *Client) do(method, path string, query url.Values, out any) (http.Header, error) {
	return c.doJSON(method, path, query, nil, out)
}

func (c *Client) doJSON(method, path string, query url.Values, payload any, out any) (http.Header, error) {
	u := c.baseURL + "/api/v4" + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var body io.Reader
	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return resp.Header, fmt.Errorf("gitlab api error: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if out == nil {
		return resp.Header, nil
	}
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	return resp.Header, dec.Decode(out)
}

func (c *Client) getList(path string, q url.Values, all bool) ([]map[string]any, error) {
	if q == nil {
		q = url.Values{}
	}
	if q.Get("per_page") == "" {
		q.Set("per_page", "100")
	}
	if q.Get("page") == "" {
		q.Set("page", "1")
	}
	var allItems []map[string]any
	for {
		var items []map[string]any
		hdr, err := c.do("GET", path, q, &items)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, items...)
		if !all {
			return allItems, nil
		}
		next := hdr.Get("X-Next-Page")
		if next == "" {
			break
		}
		q.Set("page", next)
	}
	return allItems, nil
}
