package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	domainAuth "github.com/vmaurya-21/Calance-Workflow/internal/domain/auth"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// Callback handles the OAuth callback from GitHub
// GET /api/auth/github/callback
func (h *Handler) Callback(c *gin.Context) {
	// Verify state for CSRF protection
	state := c.Query("state")
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || state == "" || state != cookieState {
		pkghttp.BadRequestResponse(c, "Invalid state parameter")
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		pkghttp.BadRequestResponse(c, "Authorization code not found")
		return
	}

	// Exchange code for token
	token, err := h.authService.ExchangeCode(c.Request.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		pkghttp.InternalServerErrorResponse(c, "Failed to exchange authorization code", err)
		return
	}

	// Get GitHub user information
	user, err := h.authService.GetGitHubUser(c.Request.Context(), token.AccessToken)
	if err != nil {
		log.Printf("Failed to get GitHub user: %v", err)
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch user information from GitHub", err)
		return
	}

	// Create or update user in database
	dbUser := &domainAuth.User{
		GitHubID:  user.GitHubID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Name:      user.Name,
		Bio:       user.Bio,
		Location:  user.Location,
		Company:   user.Company,
	}

	if err := h.userRepository.CreateOrUpdate(dbUser); err != nil {
		log.Printf("Failed to create/update user: %v", err)
		pkghttp.InternalServerErrorResponse(c, "Failed to save user", err)
		return
	}

	// Fetch the updated user to get the correct UUID
	savedUser, err := h.userRepository.FindByGitHubID(user.GitHubID)
	if err != nil || savedUser == nil {
		log.Printf("Failed to fetch saved user: %v", err)
		pkghttp.InternalServerErrorResponse(c, "Failed to retrieve saved user", err)
		return
	}

	// Create or update token in database
	expiry := token.Expiry
	tokenModel := &domainAuth.Token{
		UserID:      savedUser.ID,
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		Scope:       "user:email,read:user,read:org,repo,workflow",
		ExpiresAt:   &expiry,
	}

	if err := h.tokenRepository.CreateOrUpdate(tokenModel); err != nil {
		log.Printf("Failed to create/update token: %v", err)
		pkghttp.InternalServerErrorResponse(c, "Failed to save access token", err)
		return
	}

	// Generate JWT token
	jwtToken, err := h.authService.GenerateJWT(savedUser.ID, savedUser.Username)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		pkghttp.InternalServerErrorResponse(c, "Failed to generate authentication token", err)
		return
	}

	log.Printf("GitHub OAuth successful for user: %s", savedUser.Username)

	// Redirect to frontend with token
	frontendURL := fmt.Sprintf("%s/auth/callback?token=%s", "http://localhost:3000", jwtToken)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}
