package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// AuthClient handles GitHub OAuth operations
type AuthClient struct {
	*Client
	oauthConfig *oauth2.Config
}

// NewAuthClient creates a new GitHub OAuth client
func NewAuthClient(clientID, clientSecret, redirectURL string, scopes []string) *AuthClient {
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     githuboauth.Endpoint,
	}

	return &AuthClient{
		Client:      NewClient(),
		oauthConfig: oauthConfig,
	}
}

// GetAuthURL returns the GitHub OAuth authorization URL
func (ac *AuthClient) GetAuthURL(state string) string {
	return ac.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeCode exchanges the authorization code for an access token
func (ac *AuthClient) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := ac.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

// GetUser fetches the GitHub user information using the access token
func (ac *AuthClient) GetUser(ctx context.Context, token string) (*User, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := ac.doRequest(ctx, token, http.MethodGet, "/user", nil)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var user User
	if err := resp.UnmarshalJSON(&user); err != nil {
		return nil, err
	}

	// If email is not public, fetch it from emails endpoint
	if user.Email == "" {
		email, err := ac.getPrimaryEmail(ctx, token)
		if err == nil {
			user.Email = email
		}
	}

	return &user, nil
}

// getPrimaryEmail fetches the user's primary email from GitHub
func (ac *AuthClient) getPrimaryEmail(ctx context.Context, token string) (string, error) {
	resp, err := ac.doRequest(ctx, token, http.MethodGet, "/user/emails", nil)
	if err != nil {
		return "", err
	}

	if err := checkResponse(resp); err != nil {
		return "", err
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := resp.UnmarshalJSON(&emails); err != nil {
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
