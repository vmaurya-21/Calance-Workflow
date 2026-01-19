package github

import (
	"context"
	"fmt"
	"net/http"
)

// RepositoryClient handles GitHub repository operations
type RepositoryClient struct {
	*Client
}

// NewRepositoryClient creates a new repository client
func NewRepositoryClient() *RepositoryClient {
	return &RepositoryClient{
		Client: NewClient(),
	}
}

// GetBranches retrieves all branches for a repository
func (rc *RepositoryClient) GetBranches(ctx context.Context, token, owner, repo string) ([]Branch, error) {
	path := fmt.Sprintf("/repos/%s/%s/branches", owner, repo)
	resp, err := rc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var branches []Branch
	if err := resp.UnmarshalJSON(&branches); err != nil {
		return nil, err
	}

	return branches, nil
}

// GetCommits retrieves commits for a branch
func (rc *RepositoryClient) GetCommits(ctx context.Context, token, owner, repo, branch string, perPage int) ([]interface{}, error) {
	path := fmt.Sprintf("/repos/%s/%s/commits?sha=%s&per_page=%d", owner, repo, branch, perPage)
	resp, err := rc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var commits []interface{}
	if err := resp.UnmarshalJSON(&commits); err != nil {
		return nil, err
	}

	return commits, nil
}

// GetTags retrieves all tags for a repository
func (rc *RepositoryClient) GetTags(ctx context.Context, token, owner, repo string) ([]interface{}, error) {
	path := fmt.Sprintf("/repos/%s/%s/tags", owner, repo)
	resp, err := rc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var tags []interface{}
	if err := resp.UnmarshalJSON(&tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// CreateTag creates a new tag
func (rc *RepositoryClient) CreateTag(ctx context.Context, token, owner, repo, tagName, commitSHA string) (*Ref, error) {
	path := fmt.Sprintf("/repos/%s/%s/git/refs", owner, repo)
	body := map[string]interface{}{
		"ref": fmt.Sprintf("refs/tags/%s", tagName),
		"sha": commitSHA,
	}

	resp, err := rc.doRequest(ctx, token, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create tag: %s", resp.GetErrorMessage())
	}

	var ref Ref
	if err := resp.UnmarshalJSON(&ref); err != nil {
		return nil, err
	}

	return &ref, nil
}

// GetWorkflowRuns retrieves workflow runs for a repository
func (rc *RepositoryClient) GetWorkflowRuns(ctx context.Context, token, owner, repo string, perPage int) ([]interface{}, error) {
	path := fmt.Sprintf("/repos/%s/%s/actions/runs?per_page=%d", owner, repo, perPage)
	resp, err := rc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := resp.UnmarshalJSON(&result); err != nil {
		return nil, err
	}

	runs, ok := result["workflow_runs"].([]interface{})
	if !ok {
		return []interface{}{}, nil
	}

	return runs, nil
}

// GetWorkflowRunDetail retrieves detailed information about a workflow run
func (rc *RepositoryClient) GetWorkflowRunDetail(ctx context.Context, token, owner, repo string, runID int64) (interface{}, []interface{}, error) {
	// Get run details
	runPath := fmt.Sprintf("/repos/%s/%s/actions/runs/%d", owner, repo, runID)
	runResp, err := rc.doRequest(ctx, token, http.MethodGet, runPath, nil)
	if err != nil {
		return nil, nil, err
	}

	if err := checkResponse(runResp); err != nil {
		return nil, nil, err
	}

	var runDetail interface{}
	if err := runResp.UnmarshalJSON(&runDetail); err != nil {
		return nil, nil, err
	}

	// Get jobs
	jobsPath := fmt.Sprintf("/repos/%s/%s/actions/runs/%d/jobs", owner, repo, runID)
	jobsResp, err := rc.doRequest(ctx, token, http.MethodGet, jobsPath, nil)
	if err != nil {
		return runDetail, nil, err
	}

	if err := checkResponse(jobsResp); err != nil {
		return runDetail, nil, err
	}

	var jobsResult map[string]interface{}
	if err := jobsResp.UnmarshalJSON(&jobsResult); err != nil {
		return runDetail, nil, err
	}

	jobs, ok := jobsResult["jobs"].([]interface{})
	if !ok {
		return runDetail, []interface{}{}, nil
	}

	return runDetail, jobs, nil
}

// GetJobLogs retrieves logs for a specific job
func (rc *RepositoryClient) GetJobLogs(ctx context.Context, token, owner, repo string, jobID int64) (string, error) {
	path := fmt.Sprintf("/repos/%s/%s/actions/jobs/%d/logs", owner, repo, jobID)
	resp, err := rc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	if err := checkResponse(resp); err != nil {
		return "", err
	}

	return string(resp.Body), nil
}
