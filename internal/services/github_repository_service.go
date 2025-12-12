package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

// GitHubTag represents a Git tag
type GitHubTag struct {
	Tag     string `json:"tag"`
	SHA     string `json:"sha"`
	URL     string `json:"url"`
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

	log.Printf("DEBUG - Fetching repositories for organization: %s", orgName)

	var allRepos []GitHubRepository
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?page=%d&per_page=%d&type=all&visibility=all", orgName, page, perPage)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Printf("ERROR - Failed to create request: %v", err)
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add authorization header with access token
		req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR - Failed to fetch repositories: %v", err)
			return nil, fmt.Errorf("failed to fetch repositories: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("ERROR - Repositories endpoint returned status %d: %s", resp.StatusCode, string(body))
			return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
		}

		var repos []GitHubRepository
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			resp.Body.Close()
			log.Printf("ERROR - Failed to decode repositories: %v", err)
			return nil, fmt.Errorf("failed to decode repositories: %w", err)
		}
		resp.Body.Close()

		if len(repos) == 0 {
			break // No more pages
		}

		allRepos = append(allRepos, repos...)
		page++
	}

	log.Printf("DEBUG - Successfully fetched %d repositories for organization %s", len(allRepos), orgName)
	return allRepos, nil
}

// GetRepositoryBranches fetches all branches for a specific repository
func (s *GitHubRepositoryService) GetRepositoryBranches(ctx context.Context, accessToken, owner, repo string) ([]GitHubBranch, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	log.Printf("DEBUG - Fetching branches for repository: %s/%s", owner, repo)

	var allBranches []GitHubBranch
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches?page=%d&per_page=%d", owner, repo, page, perPage)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Printf("ERROR - Failed to create request: %v", err)
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add authorization header with access token
		req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR - Failed to fetch branches: %v", err)
			return nil, fmt.Errorf("failed to fetch branches: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("ERROR - Branches endpoint returned status %d: %s", resp.StatusCode, string(body))
			return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
		}

		var branches []GitHubBranch
		if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
			resp.Body.Close()
			log.Printf("ERROR - Failed to decode branches: %v", err)
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

	log.Printf("DEBUG - Successfully fetched %d branches for repository %s/%s", len(allBranches), owner, repo)
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

	log.Printf("DEBUG - Fetching %d commits for branch %s in repository: %s/%s", perPage, branch, owner, repo)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?sha=%s&per_page=%d&page=1", owner, repo, branch, perPage)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("ERROR - Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with access token
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR - Failed to fetch commits: %v", err)
		return nil, fmt.Errorf("failed to fetch commits: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		// Handle empty repository (409 Conflict)
		if resp.StatusCode == http.StatusConflict {
			log.Printf("INFO - Repository %s/%s branch %s is empty (no commits yet)", owner, repo, branch)
			return []GitHubCommit{}, nil // Return empty array instead of error
		}
		log.Printf("ERROR - Commits endpoint returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var commits []GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		log.Printf("ERROR - Failed to decode commits: %v", err)
		return nil, fmt.Errorf("failed to decode commits: %w", err)
	}

	log.Printf("DEBUG - Successfully fetched %d commits for branch %s", len(commits), branch)
	return commits, nil
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

	log.Printf("DEBUG - Creating %s tag '%s' for commit %s in %s/%s", req.TagType, req.TagName, req.CommitSHA, owner, req.Repo)

	// Verify the commit exists
	if err := s.verifyCommit(ctx, client, accessToken, owner, req.Repo, req.CommitSHA); err != nil {
		log.Printf("ERROR - Commit verification failed: %v", err)
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

	log.Printf("DEBUG - Successfully created tag '%s' at ref %s", req.TagName, refResponse.Ref)
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
		log.Printf("ERROR - Failed to marshal tag payload: %v", err)
		return nil, fmt.Errorf("failed to marshal tag payload: %w", err)
	}

	log.Printf("DEBUG - Creating tag object at URL: %s", tagURL)
	log.Printf("DEBUG - Tag payload: %s", string(tagBody))

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
			log.Printf("ERROR - Repository check failed with status %d: %s", repoResp.StatusCode, string(repoBody))
			return nil, fmt.Errorf("cannot access repository %s/%s (status %d). Make sure the OAuth app has access to this repository/organization", owner, req.Repo, repoResp.StatusCode)
		}
		log.Printf("DEBUG - Repository %s/%s is accessible", owner, req.Repo)
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
		log.Printf("ERROR - Failed to create tag object: %v", err)
		return nil, fmt.Errorf("failed to create tag object: %w", err)
	}
	defer tagResp.Body.Close()

	if tagResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(tagResp.Body)
		log.Printf("ERROR - Tag creation returned status %d: %s", tagResp.StatusCode, string(body))
		return nil, fmt.Errorf("GitHub API returned status %d: %s", tagResp.StatusCode, string(body))
	}

	var tagObj GitHubTag
	if err := json.NewDecoder(tagResp.Body).Decode(&tagObj); err != nil {
		log.Printf("ERROR - Failed to decode tag object: %v", err)
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
		log.Printf("ERROR - Failed to marshal reference payload: %v", err)
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
		log.Printf("ERROR - Failed to create reference: %v", err)
		return nil, fmt.Errorf("failed to create reference: %w", err)
	}
	defer refResp.Body.Close()

	if refResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(refResp.Body)
		log.Printf("ERROR - Reference creation returned status %d: %s", refResp.StatusCode, string(body))
		return nil, fmt.Errorf("GitHub API returned status %d: %s", refResp.StatusCode, string(body))
	}

	var reference GitHubReference
	if err := json.NewDecoder(refResp.Body).Decode(&reference); err != nil {
		log.Printf("ERROR - Failed to decode reference: %v", err)
		return nil, fmt.Errorf("failed to decode reference: %w", err)
	}

	return &reference, nil
}

// verifyCommit checks if a commit exists in the repository
func (s *GitHubRepositoryService) verifyCommit(ctx context.Context, client *http.Client, accessToken, owner, repo, commitSHA string) error {
	commitURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/commits/%s", owner, repo, commitSHA)

	log.Printf("DEBUG - Verifying commit %s exists in %s/%s", commitSHA, owner, repo)

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
		log.Printf("ERROR - Commit %s not found in %s/%s: %s", commitSHA, owner, repo, string(body))
		return fmt.Errorf("commit %s not found in repository %s/%s", commitSHA, owner, repo)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR - Commit verification returned status %d: %s", resp.StatusCode, string(body))
		return fmt.Errorf("failed to verify commit, GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("DEBUG - Commit %s verified successfully", commitSHA)
	return nil
}
