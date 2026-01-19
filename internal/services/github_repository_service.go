package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vmaurya-21/Calance-Workflow/internal/logger"
)

// GitHubRepository represents a GitHub repository
type GitHubRepository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	URL         string `json:"html_url"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Owner       struct {
		Login string `json:"login"`
		Type  string `json:"type"`
	} `json:"owner"`
	Language      string `json:"language"`
	StargazCount  int    `json:"stargazers_count"`
	ForksCount    int    `json:"forks_count"`
	OpenIssues    int    `json:"open_issues_count"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PushedAt      string `json:"pushed_at"`
	DefaultBranch string `json:"default_branch"`
}

// GitHubBranch represents a branch in a GitHub repository
type GitHubBranch struct {
	Name      string          `json:"name"`
	Protected bool            `json:"protected"`
	Commit    GitHubCommitRef `json:"commit"`
}

type GitHubCommitRef struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

// GitHubCommit represents a commit in a GitHub repository
type GitHubCommit struct {
	SHA    string `json:"sha"`
	NodeID string `json:"node_id"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			SHA string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		CommentCount int `json:"comment_count"`
	} `json:"commit"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Author      struct {
		Login     string `json:"login"`
		ID        int64  `json:"id"`
		AvatarURL string `json:"avatar_url"`
		URL       string `json:"url"`
	} `json:"author"`
	Committer struct {
		Login     string `json:"login"`
		ID        int64  `json:"id"`
		AvatarURL string `json:"avatar_url"`
		URL       string `json:"url"`
	} `json:"committer"`
	Parents []GitHubCommitRef `json:"parents"`
}

// GitHubTag represents a Git tag from GitHub API (for listing tags)
type GitHubTag struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	NodeID string `json:"node_id"`
}

// GitHubTagObject represents a Git tag object from GitHub API (for creating tags)
type GitHubTagObject struct {
	SHA     string `json:"sha"`
	URL     string `json:"url"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
	Tagger  struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"tagger"`
	Object struct {
		Type string `json:"type"`
		SHA  string `json:"sha"`
		URL  string `json:"url"`
	} `json:"object"`
}

// GitHubReference represents a Git reference
type GitHubReference struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
	Object struct {
		Type string `json:"type"`
		SHA  string `json:"sha"`
		URL  string `json:"url"`
	} `json:"object"`
}

// CreateTagRequest represents the request to create a tag
type CreateTagRequest struct {
	Owner       string `json:"owner"` // Optional: if not provided, uses authenticated user's username
	Repo        string `json:"repo" binding:"required"`
	CommitSHA   string `json:"commit_sha" binding:"required"`
	TagName     string `json:"tag_name" binding:"required"`
	TagMessage  string `json:"tag_message"`
	TagType     string `json:"tag_type"` // "lightweight" or "annotated"
	TaggerName  string `json:"tagger_name"`
	TaggerEmail string `json:"tagger_email"`
}

// GitHubWorkflowRun represents a GitHub Actions workflow run
type GitHubWorkflowRun struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	DisplayTitle string    `json:"display_title"`
	Status       string    `json:"status"`     // queued, in_progress, completed
	Conclusion   string    `json:"conclusion"` // success, failure, cancelled, skipped, timed_out, action_required
	HeadBranch   string    `json:"head_branch"`
	HeadSHA      string    `json:"head_sha"`
	RunNumber    int       `json:"run_number"`
	Event        string    `json:"event"` // push, pull_request, workflow_dispatch, etc.
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RunStartedAt time.Time `json:"run_started_at"`
	HTMLURL      string    `json:"html_url"`
	Actor        struct {
		Login     string `json:"login"`
		ID        int64  `json:"id"`
		AvatarURL string `json:"avatar_url"`
	} `json:"actor"`
	WorkflowID int64  `json:"workflow_id"`
	Path       string `json:"path"` // Path to workflow file
}

// GitHubWorkflowRunDetail represents detailed workflow run information
type GitHubWorkflowRunDetail struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	DisplayTitle string    `json:"display_title"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	HeadBranch   string    `json:"head_branch"`
	HeadSHA      string    `json:"head_sha"`
	RunNumber    int       `json:"run_number"`
	Event        string    `json:"event"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RunStartedAt time.Time `json:"run_started_at"`
	HTMLURL      string    `json:"html_url"`
	Actor        struct {
		Login     string `json:"login"`
		ID        int64  `json:"id"`
		AvatarURL string `json:"avatar_url"`
	} `json:"actor"`
	WorkflowID int64  `json:"workflow_id"`
	Path       string `json:"path"`
}

// GitHubWorkflowJob represents a job within a workflow run
type GitHubWorkflowJob struct {
	ID          int64                `json:"id"`
	RunID       int64                `json:"run_id"`
	Name        string               `json:"name"`
	Status      string               `json:"status"`
	Conclusion  string               `json:"conclusion"`
	StartedAt   time.Time            `json:"started_at"`
	CompletedAt time.Time            `json:"completed_at"`
	Steps       []GitHubWorkflowStep `json:"steps"`
}

// GitHubWorkflowStep represents a step within a job
type GitHubWorkflowStep struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	Number      int       `json:"number"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	Logs        string    `json:"logs,omitempty"` // Step logs
}

// GitHubWorkflowRunsResponse represents the response from GitHub Actions runs API
type GitHubWorkflowRunsResponse struct {
	TotalCount   int                 `json:"total_count"`
	WorkflowRuns []GitHubWorkflowRun `json:"workflow_runs"`
}

// GitHubWorkflowJobsResponse represents the response from GitHub Actions jobs API
type GitHubWorkflowJobsResponse struct {
	TotalCount int                 `json:"total_count"`
	Jobs       []GitHubWorkflowJob `json:"jobs"`
}

// GitHubRepositoryService handles GitHub repository operations
type GitHubRepositoryService struct{}

// NewGitHubRepositoryService creates a new GitHub repository service
func NewGitHubRepositoryService() *GitHubRepositoryService {
	return &GitHubRepositoryService{}
}

// GetOrganizationRepositories fetches all repositories for a specific GitHub organization
func (s *GitHubRepositoryService) GetOrganizationRepositories(ctx context.Context, accessToken, orgName string) ([]GitHubRepository, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	logger.Debug().Str("organization", orgName).Msg("Fetching repositories")

	var allRepos []GitHubRepository
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?page=%d&per_page=%d&type=all&visibility=all", orgName, page, perPage)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create request")
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add authorization header with access token
		req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := client.Do(req)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to fetch repositories")
			return nil, fmt.Errorf("failed to fetch repositories: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			logger.Error().Int("status_code", resp.StatusCode).Str("response", string(body)).Msg("Repositories endpoint returned error")
			return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
		}

		var repos []GitHubRepository
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			resp.Body.Close()
			logger.Error().Err(err).Msg("Failed to decode repositories")
			return nil, fmt.Errorf("failed to decode repositories: %w", err)
		}
		resp.Body.Close()

		if len(repos) == 0 {
			break // No more pages
		}

		allRepos = append(allRepos, repos...)
		page++
	}

	logger.Debug().Int("count", len(allRepos)).Str("organization", orgName).Msg("Successfully fetched repositories")
	return allRepos, nil
}

// GetRepositoryBranches fetches all branches for a specific repository
func (s *GitHubRepositoryService) GetRepositoryBranches(ctx context.Context, accessToken, owner, repo string) ([]GitHubBranch, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	logger.Debug().Str("owner", owner).Str("repo", repo).Msg("Fetching branches")

	var allBranches []GitHubBranch
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches?page=%d&per_page=%d", owner, repo, page, perPage)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create request")
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add authorization header with access token
		req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := client.Do(req)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to fetch branches")
			return nil, fmt.Errorf("failed to fetch branches: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			logger.Error().Int("status_code", resp.StatusCode).Str("response", string(body)).Msg("Branches endpoint returned error")
			return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
		}

		var branches []GitHubBranch
		if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
			resp.Body.Close()
			logger.Error().Err(err).Msg("Failed to decode branches")
			return nil, fmt.Errorf("failed to decode branches: %w", err)
		}
		resp.Body.Close()

		if len(branches) == 0 {
			break // No more pages
		}

		allBranches = append(allBranches, branches...)

		// GitHub API typically returns all branches in one page for most repos
		// Break if we got less than perPage results
		if len(branches) < perPage {
			break
		}
		page++
	}

	logger.Debug().Int("count", len(allBranches)).Str("owner", owner).Str("repo", repo).Msg("Successfully fetched branches")
	return allBranches, nil
}

// GetBranchCommits fetches the latest commits for a specific branch
func (s *GitHubRepositoryService) GetBranchCommits(ctx context.Context, accessToken, owner, repo, branch string, perPage int) ([]GitHubCommit, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if perPage <= 0 || perPage > 100 {
		perPage = 30 // Default to 30
	}

	logger.Debug().Int("per_page", perPage).Str("branch", branch).Str("owner", owner).Str("repo", repo).Msg("Fetching commits")

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?sha=%s&per_page=%d&page=1", owner, repo, branch, perPage)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with access token
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch commits")
		return nil, fmt.Errorf("failed to fetch commits: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		// Handle empty repository (409 Conflict)
		if resp.StatusCode == http.StatusConflict {
			logger.Info().Str("owner", owner).Str("repo", repo).Str("branch", branch).Msg("Repository branch is empty")
			return []GitHubCommit{}, nil // Return empty array instead of error
		}
		logger.Error().Int("status_code", resp.StatusCode).Str("response", string(body)).Msg("Commits endpoint returned error")
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var commits []GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		logger.Error().Err(err).Msg("Failed to decode commits")
		return nil, fmt.Errorf("failed to decode commits: %w", err)
	}

	logger.Debug().Int("count", len(commits)).Str("branch", branch).Msg("Successfully fetched commits")
	return commits, nil
}

// GetRepositoryTags fetches all tags for a repository
func (s *GitHubRepositoryService) GetRepositoryTags(ctx context.Context, accessToken, owner, repo string) ([]GitHubTag, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	logger.Debug().Str("owner", owner).Str("repo", repo).Msg("Fetching tags")

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with access token
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch tags")
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error().Int("status_code", resp.StatusCode).Str("response", string(body)).Msg("Tags endpoint returned error")
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var tags []GitHubTag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		logger.Error().Err(err).Msg("Failed to decode tags")
		return nil, fmt.Errorf("failed to decode tags: %w", err)
	}

	logger.Debug().Int("count", len(tags)).Str("owner", owner).Str("repo", repo).Msg("Successfully fetched tags")
	return tags, nil
}

// GetWorkflowRunDetail fetches detailed information for a specific workflow run, including jobs and their logs.
func (s *GitHubRepositoryService) GetWorkflowRunDetail(ctx context.Context, accessToken, owner, repo string, runID int64) (*GitHubWorkflowRunDetail, []GitHubWorkflowJob, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second) // Increased timeout for logs
	defer cancel()

	client := &http.Client{
		Timeout: 30 * time.Second, // Increased client timeout
	}

	// 1. Fetch workflow run details
	runURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%d", owner, repo, runID)
	req, err := http.NewRequestWithContext(ctx, "GET", runURL, nil)
	if err != nil {
		logger.Error().Int64("run_id", runID).Err(err).Msg("Failed to create request for workflow run")
		return nil, nil, fmt.Errorf("failed to create request for workflow run: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Int64("run_id", runID).Err(err).Msg("Failed to fetch workflow run")
		return nil, nil, fmt.Errorf("failed to fetch workflow run: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error().Int("status_code", resp.StatusCode).Int64("run_id", runID).Str("response", string(body)).Msg("Workflow run endpoint returned error")
		return nil, nil, fmt.Errorf("GitHub API returned status %d for workflow run %d: %s", resp.StatusCode, runID, string(body))
	}

	var runDetail GitHubWorkflowRunDetail
	if err := json.NewDecoder(resp.Body).Decode(&runDetail); err != nil {
		logger.Error().Int64("run_id", runID).Err(err).Msg("Failed to decode workflow run")
		return nil, nil, fmt.Errorf("failed to decode workflow run: %w", err)
	}
	logger.Debug().Int64("run_id", runID).Msg("Successfully fetched workflow run details")

	// 2. Fetch jobs for the workflow run
	jobsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%d/jobs?per_page=100", owner, repo, runID)
	req, err = http.NewRequestWithContext(ctx, "GET", jobsURL, nil)
	if err != nil {
		logger.Error().Int64("run_id", runID).Err(err).Msg("Failed to create request for jobs")
		return nil, nil, fmt.Errorf("failed to create request for workflow jobs: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err = client.Do(req)
	if err != nil {
		logger.Error().Int64("run_id", runID).Err(err).Msg("Failed to fetch jobs")
		return nil, nil, fmt.Errorf("failed to fetch workflow jobs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error().Int("status_code", resp.StatusCode).Int64("run_id", runID).Str("response", string(body)).Msg("Workflow jobs endpoint returned error")
		return nil, nil, fmt.Errorf("GitHub API returned status %d for workflow jobs of run %d: %s", resp.StatusCode, runID, string(body))
	}

	var jobsResponse GitHubWorkflowJobsResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobsResponse); err != nil {
		logger.Error().Int64("run_id", runID).Err(err).Msg("Failed to decode jobs")
		return nil, nil, fmt.Errorf("failed to decode workflow jobs: %w", err)
	}
	logger.Debug().Int("count", len(jobsResponse.Jobs)).Int64("run_id", runID).Msg("Successfully fetched jobs")
	return &runDetail, jobsResponse.Jobs, nil
}

// GetJobLogs fetches ALL logs for a job in a single response
func (s *GitHubRepositoryService) GetJobLogs(ctx context.Context, accessToken, owner, repo string, jobID int64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	logger.Debug().Int64("job_id", jobID).Str("owner", owner).Str("repo", repo).Msg("Fetching complete job logs")

	// GitHub API endpoint for complete job logs
	logsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/jobs/%d/logs", owner, repo, jobID)
	req, err := http.NewRequestWithContext(ctx, "GET", logsURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for job logs: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch job logs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API returned status %d for job logs: %s", resp.StatusCode, string(body))
	}

	logs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read job logs: %w", err)
	}

	logger.Debug().Int64("job_id", jobID).Int("bytes", len(logs)).Msg("Successfully fetched complete job logs")
	return string(logs), nil
}

// CreateTag creates and pushes a tag for a specific commit
func (s *GitHubRepositoryService) CreateTag(ctx context.Context, accessToken string, owner string, req CreateTagRequest) (*GitHubReference, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Default to annotated tag if not specified
	if req.TagType == "" {
		req.TagType = "annotated"
	}

	logger.Debug().Str("tag_type", req.TagType).Str("tag_name", req.TagName).Str("commit_sha", req.CommitSHA).Str("owner", owner).Str("repo", req.Repo).Msg("Creating tag")

	// Verify the commit exists
	if err := s.verifyCommit(ctx, client, accessToken, owner, req.Repo, req.CommitSHA); err != nil {
		logger.Error().Err(err).Msg("Commit verification failed")
		return nil, fmt.Errorf("commit verification failed: %w", err)
	}

	var refResponse *GitHubReference
	var err error

	if req.TagType == "annotated" {
		refResponse, err = s.createAnnotatedTag(ctx, client, accessToken, owner, req)
	} else {
		refResponse, err = s.createLightweightTag(ctx, client, accessToken, owner, req)
	}

	if err != nil {
		return nil, err
	}

	logger.Debug().Str("tag_name", req.TagName).Str("ref", refResponse.Ref).Msg("Successfully created tag")
	return refResponse, nil
}

// createAnnotatedTag creates an annotated tag (with metadata)
func (s *GitHubRepositoryService) createAnnotatedTag(ctx context.Context, client *http.Client, accessToken string, owner string, req CreateTagRequest) (*GitHubReference, error) {
	// Step 1: Create the tag object
	tagURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/tags", owner, req.Repo)

	tagMessage := req.TagMessage
	if tagMessage == "" {
		tagMessage = fmt.Sprintf("Tag %s", req.TagName)
	}

	tagPayload := map[string]interface{}{
		"tag":     req.TagName,
		"message": tagMessage,
		"object":  req.CommitSHA,
		"type":    "commit",
	}

	// Add tagger information if provided
	if req.TaggerName != "" && req.TaggerEmail != "" {
		tagPayload["tagger"] = map[string]interface{}{
			"name":  req.TaggerName,
			"email": req.TaggerEmail,
			"date":  time.Now().Format(time.RFC3339),
		}
	}

	tagBody, err := json.Marshal(tagPayload)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal tag payload")
		return nil, fmt.Errorf("failed to marshal tag payload: %w", err)
	}

	logger.Debug().Str("url", tagURL).Msg("Creating tag object")
	logger.Debug().Str("payload", string(tagBody)).Msg("Tag payload")

	// First verify repository exists and is accessible
	repoCheckURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, req.Repo)
	repoReq, _ := http.NewRequestWithContext(ctx, "GET", repoCheckURL, nil)
	repoReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	repoReq.Header.Set("Accept", "application/vnd.github.v3+json")
	repoResp, _ := client.Do(repoReq)
	if repoResp != nil {
		defer repoResp.Body.Close()
		if repoResp.StatusCode != http.StatusOK {
			repoBody, _ := io.ReadAll(repoResp.Body)
			logger.Error().Int("status_code", repoResp.StatusCode).Str("response", string(repoBody)).Msg("Repository check failed")
			return nil, fmt.Errorf("cannot access repository %s/%s (status %d). Make sure the OAuth app has access to this repository/organization", owner, req.Repo, repoResp.StatusCode)
		}
		logger.Debug().Str("owner", owner).Str("repo", req.Repo).Msg("Repository is accessible")
	}

	tagReq, err := http.NewRequestWithContext(ctx, "POST", tagURL, bytes.NewReader(tagBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create tag request: %w", err)
	}
	tagReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	tagReq.Header.Set("Accept", "application/vnd.github.v3+json")
	tagReq.Header.Set("Content-Type", "application/json")

	tagResp, err := client.Do(tagReq)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create tag object")
		return nil, fmt.Errorf("failed to create tag object: %w", err)
	}
	defer tagResp.Body.Close()

	if tagResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(tagResp.Body)
		logger.Error().Int("status_code", tagResp.StatusCode).Str("response", string(body)).Msg("Tag creation returned error")
		return nil, fmt.Errorf("GitHub API returned status %d: %s", tagResp.StatusCode, string(body))
	}

	var tagObj GitHubTagObject
	if err := json.NewDecoder(tagResp.Body).Decode(&tagObj); err != nil {
		logger.Error().Err(err).Msg("Failed to decode tag object")
		return nil, fmt.Errorf("failed to decode tag object: %w", err)
	}

	// Step 2: Create the reference (push the tag)
	return s.createReference(ctx, client, accessToken, owner, req.Repo, req.TagName, tagObj.SHA)
}

// createLightweightTag creates a lightweight tag (just a reference)
func (s *GitHubRepositoryService) createLightweightTag(ctx context.Context, client *http.Client, accessToken string, owner string, req CreateTagRequest) (*GitHubReference, error) {
	return s.createReference(ctx, client, accessToken, owner, req.Repo, req.TagName, req.CommitSHA)
}

// createReference creates a Git reference (pushes the tag to remote)
func (s *GitHubRepositoryService) createReference(ctx context.Context, client *http.Client, accessToken, owner, repo, tagName, sha string) (*GitHubReference, error) {
	refURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs", owner, repo)

	refPayload := map[string]interface{}{
		"ref": fmt.Sprintf("refs/tags/%s", tagName),
		"sha": sha,
	}

	refBody, err := json.Marshal(refPayload)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal reference payload")
		return nil, fmt.Errorf("failed to marshal reference payload: %w", err)
	}

	refReq, err := http.NewRequestWithContext(ctx, "POST", refURL, bytes.NewReader(refBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create reference request: %w", err)
	}
	refReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	refReq.Header.Set("Accept", "application/vnd.github.v3+json")
	refReq.Header.Set("Content-Type", "application/json")

	refResp, err := client.Do(refReq)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create reference")
		return nil, fmt.Errorf("failed to create reference: %w", err)
	}
	defer refResp.Body.Close()

	if refResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(refResp.Body)
		logger.Error().Int("status_code", refResp.StatusCode).Str("response", string(body)).Msg("Reference creation returned error")
		return nil, fmt.Errorf("GitHub API returned status %d: %s", refResp.StatusCode, string(body))
	}

	var reference GitHubReference
	if err := json.NewDecoder(refResp.Body).Decode(&reference); err != nil {
		logger.Error().Err(err).Msg("Failed to decode reference")
		return nil, fmt.Errorf("failed to decode reference: %w", err)
	}

	return &reference, nil
}

// verifyCommit checks if a commit exists in the repository
func (s *GitHubRepositoryService) verifyCommit(ctx context.Context, client *http.Client, accessToken, owner, repo, commitSHA string) error {
	commitURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/commits/%s", owner, repo, commitSHA)

	logger.Debug().Str("commit_sha", commitSHA).Str("owner", owner).Str("repo", repo).Msg("Verifying commit exists")

	req, err := http.NewRequestWithContext(ctx, "GET", commitURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create commit verification request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify commit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		logger.Error().Str("commit_sha", commitSHA).Str("owner", owner).Str("repo", repo).Str("response", string(body)).Msg("Commit not found in repository")
		return fmt.Errorf("commit %s not found in repository %s/%s", commitSHA, owner, repo)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error().Int("status_code", resp.StatusCode).Str("response", string(body)).Msg("Commit verification returned error status")
		return fmt.Errorf("failed to verify commit, GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	logger.Debug().Str("commit_sha", commitSHA).Msg("Commit verified successfully")
	return nil
}

// GetRepositoryWorkflowRuns fetches workflow runs for a specific repository
func (s *GitHubRepositoryService) GetRepositoryWorkflowRuns(ctx context.Context, accessToken, owner, repo string, perPage int) ([]GitHubWorkflowRun, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if perPage <= 0 || perPage > 100 {
		perPage = 30 // Default to 30
	}

	logger.Debug().Int("per_page", perPage).Str("owner", owner).Str("repo", repo).Msg("Fetching workflow runs")

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs?per_page=%d&page=1", owner, repo, perPage)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with access token
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch workflow runs")
		return nil, fmt.Errorf("failed to fetch workflow runs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error().Int("status_code", resp.StatusCode).Str("response", string(body)).Msg("Workflow runs endpoint returned error")
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var runsResponse GitHubWorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&runsResponse); err != nil {
		logger.Error().Err(err).Msg("Failed to decode workflow runs")
		return nil, fmt.Errorf("failed to decode workflow runs: %w", err)
	}

	logger.Debug().Int("count", len(runsResponse.WorkflowRuns)).Int("total", runsResponse.TotalCount).Msg("Successfully fetched workflow runs")
	return runsResponse.WorkflowRuns, nil
}
