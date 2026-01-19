package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vmaurya-21/Calance-Workflow/internal/config"
	"github.com/vmaurya-21/Calance-Workflow/internal/logger"
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
		Scopes:       []string{"user:email", "read:user", "read:org", "repo"},
		Endpoint:     github.Endpoint,
	}

	return &GitHubOAuthService{
		config: oauthConfig,
	}
}

// GetAuthURL returns the GitHub OAuth authorization URL
func (s *GitHubOAuthService) GetAuthURL(state string) string {
	authURL := s.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	logger.Debug().
		Str("auth_url", authURL).
		Str("client_id", s.config.ClientID).
		Str("redirect_url", s.config.RedirectURL).
		Msg("Generated OAuth authorization URL")
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

	logger.Debug().Str("token_prefix", token.AccessToken[:20]+"...").Msg("Fetching user info from GitHub API")

	// Fetch user data from GitHub API
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		logger.Error().Err(err).Msg("GitHub API request failed")
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error().Int("status_code", resp.StatusCode).Str("response", string(body)).Msg("GitHub API returned error status")
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		logger.Error().Err(err).Msg("Failed to decode GitHub user response")
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	logger.Debug().Str("login", githubUser.Login).Int64("id", githubUser.ID).Msg("Successfully fetched GitHub user")

	// If email is not public, fetch it from emails endpoint
	if githubUser.Email == "" {
		logger.Debug().Msg("Email not public, fetching from emails endpoint")
		email, err := s.getPrimaryEmail(client)
		if err == nil {
			githubUser.Email = email
			logger.Debug().Str("email", githubUser.Email).Msg("Fetched email")
		} else {
			logger.Warn().Err(err).Msg("Could not fetch email")
		}
	}

	return &githubUser, nil
}

// getPrimaryEmail fetches the user's primary email from GitHub
func (s *GitHubOAuthService) getPrimaryEmail(client *http.Client) (string, error) {
	logger.Debug().Msg("Fetching emails from GitHub API")
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch emails")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status_code", resp.StatusCode).Msg("Emails endpoint returned error status")
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
	user := &models.User{
		GitHubID:  githubUser.ID,
		Username:  githubUser.Login,
		AvatarURL: githubUser.AvatarURL,
	}

	// Convert empty strings to nil pointers for optional fields
	if githubUser.Email != "" {
		user.Email = &githubUser.Email
	}
	if githubUser.Name != "" {
		user.Name = &githubUser.Name
	}
	if githubUser.Bio != "" {
		user.Bio = &githubUser.Bio
	}
	if githubUser.Location != "" {
		user.Location = &githubUser.Location
	}
	if githubUser.Company != "" {
		user.Company = &githubUser.Company
	}

	return user
}
