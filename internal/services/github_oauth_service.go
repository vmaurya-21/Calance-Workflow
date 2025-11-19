package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vmaurya-21/Calance-Workflow/internal/config"
	"github.com/vmaurya-21/Calance-Workflow/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubUser represents the GitHub user response
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Company   string `json:"company"`
}

// GitHubOAuthService handles GitHub OAuth operations
type GitHubOAuthService struct {
	config *oauth2.Config
}

// NewGitHubOAuthService creates a new GitHub OAuth service
func NewGitHubOAuthService(cfg *config.Config) *GitHubOAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GitHub.ClientID,
		ClientSecret: cfg.GitHub.ClientSecret,
		RedirectURL:  cfg.GitHub.RedirectURL,
		Scopes:       []string{"user:email", "read:user", "read:org"},
		Endpoint:     github.Endpoint,
	}

	return &GitHubOAuthService{
		config: oauthConfig,
	}
}

// GetAuthURL returns the GitHub OAuth authorization URL
func (s *GitHubOAuthService) GetAuthURL(state string) string {
	authURL := s.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	log.Printf("DEBUG - Generated Auth URL: %s", authURL)
	log.Printf("DEBUG - OAuth Config ClientID: %s", s.config.ClientID)
	log.Printf("DEBUG - OAuth Config RedirectURL: %s", s.config.RedirectURL)
	return authURL
}

// ExchangeCode exchanges the authorization code for an access token
func (s *GitHubOAuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

// GetGitHubUser fetches the GitHub user information using the access token
func (s *GitHubOAuthService) GetGitHubUser(ctx context.Context, token *oauth2.Token) (*GitHubUser, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := s.config.Client(ctx, token)

	log.Printf("DEBUG - Fetching user info from GitHub API with token: %v", token.AccessToken[:20]+"...")

	// Fetch user data from GitHub API
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Printf("ERROR - GitHub API request failed: %v (type: %T)", err, err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR - GitHub API returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		log.Printf("ERROR - Failed to decode GitHub user response: %v", err)
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	log.Printf("DEBUG - Successfully fetched GitHub user: %s (ID: %d)", githubUser.Login, githubUser.ID)

	// If email is not public, fetch it from emails endpoint
	if githubUser.Email == "" {
		log.Printf("DEBUG - Email not public, fetching from emails endpoint")
		email, err := s.getPrimaryEmail(client)
		if err == nil {
			githubUser.Email = email
			log.Printf("DEBUG - Fetched email: %s", githubUser.Email)
		} else {
			log.Printf("WARNING - Could not fetch email: %v", err)
		}
	}

	return &githubUser, nil
}

// getPrimaryEmail fetches the user's primary email from GitHub
func (s *GitHubOAuthService) getPrimaryEmail(client *http.Client) (string, error) {
	log.Printf("DEBUG - Fetching emails from GitHub API")
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		log.Printf("ERROR - Failed to fetch emails: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR - Emails endpoint returned status %d", resp.StatusCode)
		return "", fmt.Errorf("failed to fetch emails: status %d", resp.StatusCode)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	// Find primary verified email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// Fallback to first verified email
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}

// GitHubOrganization represents a GitHub organization
type GitHubOrganization struct {
	ID     int64  `json:"id"`
	Login  string `json:"login"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Avatar string `json:"avatar_url"`
}

// GetUserOrganizations fetches all organizations for the authenticated user using access token string
func (s *GitHubOAuthService) GetUserOrganizations(ctx context.Context, accessToken string) ([]GitHubOrganization, error) {
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

// ConvertToUser converts GitHubUser to models.User
func (s *GitHubOAuthService) ConvertToUser(githubUser *GitHubUser) *models.User {
	return &models.User{
		GitHubID:  githubUser.ID,
		Username:  githubUser.Login,
		Email:     githubUser.Email,
		AvatarURL: githubUser.AvatarURL,
		Name:      githubUser.Name,
		Bio:       githubUser.Bio,
		Location:  githubUser.Location,
		Company:   githubUser.Company,
	}
}
