package services

import (
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
