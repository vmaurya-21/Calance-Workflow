package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Helper methods for PR creation

// getDefaultBranch gets the default branch of the repository
func (s *WorkflowService) getDefaultBranch(ctx context.Context, token, owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get repository info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var repoInfo struct {
		DefaultBranch string `json:"default_branch"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &repoInfo); err != nil {
		return "", err
	}

	return repoInfo.DefaultBranch, nil
}

// getBranchSHA gets the SHA of the latest commit on a branch
func (s *WorkflowService) getBranchSHA(ctx context.Context, token, owner, repo, branch string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/%s", owner, repo, branch)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get branch SHA: status %d, body: %s", resp.StatusCode, string(body))
	}

	var refInfo struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &refInfo); err != nil {
		return "", err
	}

	return refInfo.Object.SHA, nil
}

// createBranch creates a new branch from a base SHA
func (s *WorkflowService) createBranch(ctx context.Context, token, owner, repo, branchName, baseSHA string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs", owner, repo)

	requestBody := map[string]interface{}{
		"ref": fmt.Sprintf("refs/heads/%s", branchName),
		"sha": baseSHA,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create branch: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// createPullRequest creates a pull request
func (s *WorkflowService) createPullRequest(ctx context.Context, token, owner, repo, headBranch, baseBranch, workflowName string) (string, int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)

	requestBody := map[string]interface{}{
		"title": fmt.Sprintf("Add workflow: %s", workflowName),
		"body":  fmt.Sprintf("This PR adds the GitHub Actions workflow for `%s`.\n\nGenerated automatically by Calance Workflow Manager.", workflowName),
		"head":  headBranch,
		"base":  baseBranch,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", 0, fmt.Errorf("failed to create pull request: status %d, body: %s", resp.StatusCode, string(body))
	}

	var prInfo struct {
		HTMLURL string `json:"html_url"`
		Number  int    `json:"number"`
	}

	if err := json.Unmarshal(body, &prInfo); err != nil {
		return "", 0, err
	}

	return prInfo.HTMLURL, prInfo.Number, nil
}
