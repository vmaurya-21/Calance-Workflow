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

// GitHubOrganization represents a GitHub organization
type GitHubOrganization struct {
	ID     int64  `json:"id"`
	Login  string `json:"login"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Avatar string `json:"avatar_url"`
}

// GitHubOrganizationService handles GitHub organization operations
type GitHubOrganizationService struct {
	repositoryService *GitHubRepositoryService
}

// NewGitHubOrganizationService creates a new GitHub organization service
func NewGitHubOrganizationService(repositoryService *GitHubRepositoryService) *GitHubOrganizationService {
	return &GitHubOrganizationService{
		repositoryService: repositoryService,
	}
}

// GetUserOrganizations fetches all organizations for the authenticated user using access token string
func (s *GitHubOrganizationService) GetUserOrganizations(ctx context.Context, accessToken string) ([]GitHubOrganization, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	log.Printf("DEBUG - Fetching organizations from GitHub API with token")

	var allOrgs []GitHubOrganization
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/user/orgs?page=%d&per_page=%d", page, perPage)
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
			log.Printf("ERROR - Failed to fetch organizations: %v", err)
			return nil, fmt.Errorf("failed to fetch organizations: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("ERROR - Organizations endpoint returned status %d: %s", resp.StatusCode, string(body))
			return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
		}

		var orgs []GitHubOrganization
		if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
			resp.Body.Close()
			log.Printf("ERROR - Failed to decode organizations: %v", err)
			return nil, fmt.Errorf("failed to decode organizations: %w", err)
		}
		resp.Body.Close()

		if len(orgs) == 0 {
			break // No more pages
		}

		allOrgs = append(allOrgs, orgs...)
		page++
	}

	log.Printf("DEBUG - Successfully fetched %d organizations", len(allOrgs))
	return allOrgs, nil
}

// GetOrganizationRepositories fetches all repositories for a specific GitHub organization
func (s *GitHubOrganizationService) GetOrganizationRepositories(ctx context.Context, accessToken, orgName string) ([]GitHubRepository, error) {
	return s.repositoryService.GetOrganizationRepositories(ctx, accessToken, orgName)
}

// GetUserRepositories fetches all repositories accessible to the authenticated user from their organizations
func (s *GitHubOrganizationService) GetUserRepositories(ctx context.Context, accessToken string) (map[string][]GitHubRepository, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	log.Printf("DEBUG - Fetching all accessible repositories for user")

	// First, get all user organizations
	organizations, err := s.GetUserOrganizations(ctx, accessToken)
	if err != nil {
		log.Printf("ERROR - Failed to fetch organizations: %v", err)
		return nil, err
	}

	log.Printf("DEBUG - Found %d organizations for user", len(organizations))

	// Fetch repositories for each organization
	reposByOrg := make(map[string][]GitHubRepository)

	for _, org := range organizations {
		repos, err := s.repositoryService.GetOrganizationRepositories(ctx, accessToken, org.Login)
		if err != nil {
			log.Printf("WARNING - Failed to fetch repositories for organization %s: %v", org.Login, err)
			// Continue with other organizations instead of failing
			continue
		}
		reposByOrg[org.Login] = repos
		log.Printf("DEBUG - Fetched %d repositories for organization %s", len(repos), org.Login)
	}

	return reposByOrg, nil
}
