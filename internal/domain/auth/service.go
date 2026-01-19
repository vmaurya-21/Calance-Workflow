package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/github"
	"github.com/vmaurya-21/Calance-Workflow/internal/utils"
	"golang.org/x/oauth2"
)

// Service handles authentication business logic
type Service struct {
	githubAuth *github.AuthClient
}

// NewService creates a new auth service
func NewService(clientID, clientSecret, redirectURL string, scopes []string) *Service {
	return &Service{
		githubAuth: github.NewAuthClient(clientID, clientSecret, redirectURL, scopes),
	}
}

// GetAuthURL returns the GitHub OAuth authorization URL
func (s *Service) GetAuthURL(state string) string {
	return s.githubAuth.GetAuthURL(state)
}

// ExchangeCode exchanges authorization code for access token
func (s *Service) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.githubAuth.ExchangeCode(ctx, code)
}

// GetGitHubUser fetches user information from GitHub
func (s *Service) GetGitHubUser(ctx context.Context, token string) (*User, error) {
	githubUser, err := s.githubAuth.GetUser(ctx, token)
	if err != nil {
		return nil, err
	}

	return &User{
		GitHubID:  githubUser.ID,
		Username:  githubUser.Login,
		Email:     githubUser.Email,
		AvatarURL: githubUser.AvatarURL,
		Name:      githubUser.Name,
		Bio:       githubUser.Bio,
		Location:  githubUser.Location,
		Company:   githubUser.Company,
	}, nil
}

// GenerateJWT generates a JWT token for a user
func (s *Service) GenerateJWT(userID uuid.UUID, username string) (string, error) {
	token, err := utils.GenerateToken(userID, username)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}
	return token, nil
}

// ValidateJWT validates a JWT token
func (s *Service) ValidateJWT(tokenString string) (*utils.JWTClaims, error) {
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to validate JWT: %w", err)
	}
	return claims, nil
}
