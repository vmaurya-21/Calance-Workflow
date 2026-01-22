package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// WorkflowClient handles GitHub workflow operations
type WorkflowClient struct {
	*Client
}

// NewWorkflowClient creates a new workflow client
func NewWorkflowClient() *WorkflowClient {
	return &WorkflowClient{
		Client: NewClient(),
	}
}

// VerifyRepository checks if a repository exists and is accessible
func (wc *WorkflowClient) VerifyRepository(ctx context.Context, token, owner, repo string) error {
	path := fmt.Sprintf("/repos/%s/%s", owner, repo)
	resp, err := wc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return checkResponse(resp)
}

// GetDefaultBranch retrieves the default branch of a repository
func (wc *WorkflowClient) GetDefaultBranch(ctx context.Context, token, owner, repo string) (string, error) {
	path := fmt.Sprintf("/repos/%s/%s", owner, repo)
	resp, err := wc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	if err := checkResponse(resp); err != nil {
		return "", err
	}

	var repoInfo Repository
	if err := resp.UnmarshalJSON(&repoInfo); err != nil {
		return "", err
	}

	return repoInfo.DefaultBranch, nil
}

// GetBranchSHA retrieves the SHA of the latest commit on a branch
func (wc *WorkflowClient) GetBranchSHA(ctx context.Context, token, owner, repo, branch string) (string, error) {
	path := fmt.Sprintf("/repos/%s/%s/git/refs/heads/%s", owner, repo, branch)
	resp, err := wc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	if err := checkResponse(resp); err != nil {
		return "", err
	}

	var refInfo Ref
	if err := resp.UnmarshalJSON(&refInfo); err != nil {
		return "", err
	}

	return refInfo.Object.SHA, nil
}

// CreateBranch creates a new branch from a base SHA
func (wc *WorkflowClient) CreateBranch(ctx context.Context, token, owner, repo, branchName, baseSHA string) error {
	path := fmt.Sprintf("/repos/%s/%s/git/refs", owner, repo)
	body := map[string]interface{}{
		"ref": fmt.Sprintf("refs/heads/%s", branchName),
		"sha": baseSHA,
	}

	resp, err := wc.doRequest(ctx, token, http.MethodPost, path, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create branch: %s", resp.GetErrorMessage())
	}

	return nil
}

// CreateFile creates a file in the repository
func (wc *WorkflowClient) CreateFile(ctx context.Context, token, owner, repo, filePath, content, message, branch string) error {
	path := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, filePath)
	body := map[string]interface{}{
		"message": message,
		"content": base64.StdEncoding.EncodeToString([]byte(content)),
		"branch":  branch,
	}

	resp, err := wc.doRequest(ctx, token, http.MethodPut, path, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return checkResponse(resp)
	}

	return nil
}

// CreatePullRequest creates a pull request
func (wc *WorkflowClient) CreatePullRequest(ctx context.Context, token, owner, repo, headBranch, baseBranch, title, body string) (string, int, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)
	requestBody := map[string]interface{}{
		"title": title,
		"body":  body,
		"head":  headBranch,
		"base":  baseBranch,
	}

	resp, err := wc.doRequest(ctx, token, http.MethodPost, path, requestBody)
	if err != nil {
		return "", 0, err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", 0, fmt.Errorf("failed to create pull request: %s", resp.GetErrorMessage())
	}

	var prInfo PullRequest
	if err := resp.UnmarshalJSON(&prInfo); err != nil {
		return "", 0, err
	}

	return prInfo.HTMLURL, prInfo.Number, nil
}

// GetWorkflowFiles retrieves all workflow files from a repository
func (wc *WorkflowClient) GetWorkflowFiles(ctx context.Context, token, owner, repo string) ([]Content, error) {
	path := fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", owner, repo)
	resp, err := wc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		// .github/workflows directory doesn't exist
		return []Content{}, nil
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var files []Content
	if err := resp.UnmarshalJSON(&files); err != nil {
		return nil, err
	}

	// Filter only .yml and .yaml files
	workflows := make([]Content, 0)
	for _, file := range files {
		if file.Type == "file" && (strings.HasSuffix(file.Name, ".yml") || strings.HasSuffix(file.Name, ".yaml")) {
			workflows = append(workflows, file)
		}
	}

	return workflows, nil
}

// GetFileContent retrieves the content of a file from the repository
func (wc *WorkflowClient) GetFileContent(ctx context.Context, token, owner, repo, filePath string) (string, string, error) {
	path := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, filePath)
	resp, err := wc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return "", "", err
	}

	if err := checkResponse(resp); err != nil {
		return "", "", err
	}

	var fileInfo Content
	if err := resp.UnmarshalJSON(&fileInfo); err != nil {
		return "", "", err
	}

	// Decode base64 content
	decodedContent, err := base64.StdEncoding.DecodeString(fileInfo.Content)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode file content: %w", err)
	}

	return string(decodedContent), fileInfo.SHA, nil
}

// UpdateFile updates an existing file in the repository
func (wc *WorkflowClient) UpdateFile(ctx context.Context, token, owner, repo, filePath, content, message, branch, sha string) error {
	path := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, filePath)
	body := map[string]interface{}{
		"message": message,
		"content": base64.StdEncoding.EncodeToString([]byte(content)),
		"branch":  branch,
		"sha":     sha, // Required for updates to prevent conflicts
	}

	resp, err := wc.doRequest(ctx, token, http.MethodPut, path, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return checkResponse(resp)
	}

	return nil
}
