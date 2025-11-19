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
		Scopes:       []string{"user:email", "read:user"},
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
