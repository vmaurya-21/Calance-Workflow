package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vmaurya-21/Calance-Workflow/internal/config"

	// "github.com/vmaurya-21/Calance-Workflow/internal/repositories"
	"github.com/vmaurya-21/Calance-Workflow/internal/services"
	"github.com/vmaurya-21/Calance-Workflow/internal/utils"
)

type AuthController struct {
	oauthService   *services.GitHubOAuthService
	userRepository interface{} // TODO: Database support
	config         *config.Config
}

// NewAuthController creates a new auth controller
func NewAuthController(
	oauthService *services.GitHubOAuthService,
	userRepository interface{},
	cfg *config.Config,
) *AuthController {
	return &AuthController{
		oauthService:   oauthService,
		userRepository: userRepository,
		config:         cfg,
	}
}

// GitHubLogin redirects the user to GitHub OAuth authorization page
// GET /api/auth/github
func (ac *AuthController) GitHubLogin(c *gin.Context) {
	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate state", err)
		return
	}

	// Store state in session or cookie (simplified: using cookie)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true) // 5 minutes expiry

	// Get GitHub OAuth URL
	authURL := ac.oauthService.GetAuthURL(state)

	// Redirect to GitHub
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback handles the OAuth callback from GitHub
// GET /api/auth/github/callback
func (ac *AuthController) GitHubCallback(c *gin.Context) {
	// Verify state for CSRF protection
	state := c.Query("state")
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || state == "" || state != cookieState {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid state parameter", fmt.Errorf("CSRF validation failed"))
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Authorization code not found", nil)
		return
	}

	// Exchange code for token
	token, err := ac.oauthService.ExchangeCode(c.Request.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to exchange authorization code", err)
		return
	}

	// Get GitHub user information
	githubUser, err := ac.oauthService.GetGitHubUser(c.Request.Context(), token)
	if err != nil {
		log.Printf("Failed to get GitHub user: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch user information from GitHub", err)
		return
	}

	// Convert to User model
	user := ac.oauthService.ConvertToUser(githubUser)

	// TODO: Create or update user in database
	// if err := ac.userRepository.CreateOrUpdate(user); err != nil {
	// 	log.Printf("Failed to create/update user: %v", err)
	// 	utils.InternalServerErrorResponse(c, "Failed to save user", err)
	// 	return
	// }

	// Generate JWT token with GitHub access token
	jwtToken, err := utils.GenerateToken(user.ID, githubUser.ID, githubUser.Login, githubUser.Email, token.AccessToken)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to generate authentication token", err)
		return
	}

	// Log successful GitHub OAuth
	log.Printf("GitHub OAuth successful for user: %s", githubUser.Login)

	// Redirect to frontend with token
	frontendURL := fmt.Sprintf("%s/auth/callback?token=%s", ac.config.Frontend.URL, jwtToken)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// GetCurrentUser returns the current authenticated user
// GET /api/auth/me
// TODO: Implement when database is available
func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "User endpoint - database required", nil)
}

// Logout logs out the user (client-side token removal)
// POST /api/auth/logout
func (ac *AuthController) Logout(c *gin.Context) {
	// JWT is stateless, so logout is handled client-side by removing the token
	utils.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}

// GetUserOrganizations returns all GitHub organizations for the authenticated user
// GET /api/auth/organizations
func (ac *AuthController) GetUserOrganizations(c *gin.Context) {
	// Get the authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Authorization header is required")
		return
	}

	// Extract token from header
	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid authorization header format")
		return
	}

	// Validate JWT token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	log.Printf("DEBUG - Fetching organizations for user: %s", claims.Username)

	// Get organizations using the GitHub token stored in JWT
	organizations, err := ac.oauthService.GetUserOrganizations(c.Request.Context(), claims.GitHubToken)
	if err != nil {
		log.Printf("ERROR - Failed to fetch organizations: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch organizations", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Organizations fetched successfully", gin.H{
		"user":               claims.Username,
		"organizations":      organizations,
		"organization_count": len(organizations),
	})
}

// GetOrganizationRepositories returns all repositories for a specific organization
// GET /api/auth/organizations/:org/repositories
func (ac *AuthController) GetOrganizationRepositories(c *gin.Context) {
	// Get organization name from URL parameter
	orgName := c.Param("org")
	if orgName == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Organization name is required", nil)
		return
	}

	// Get the authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Authorization header is required")
		return
	}

	// Extract token from header
	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid authorization header format")
		return
	}

	// Validate JWT token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	log.Printf("DEBUG - Fetching repositories for organization: %s", orgName)

	// Get repositories for the organization
	repositories, err := ac.oauthService.GetOrganizationRepositories(c.Request.Context(), claims.GitHubToken, orgName)
	if err != nil {
		log.Printf("ERROR - Failed to fetch repositories: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch repositories", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Repositories fetched successfully", gin.H{
		"user":             claims.Username,
		"organization":     orgName,
		"repositories":     repositories,
		"repository_count": len(repositories),
	})
}

// GetUserRepositories returns all repositories accessible to the authenticated user from their organizations
// GET /api/auth/repositories
func (ac *AuthController) GetUserRepositories(c *gin.Context) {
	// Get the authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Authorization header is required")
		return
	}

	// Extract token from header
	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid authorization header format")
		return
	}

	// Validate JWT token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	log.Printf("DEBUG - Fetching all accessible repositories for user: %s", claims.Username)

	// Get repositories from all organizations
	repositoriesByOrg, err := ac.oauthService.GetUserRepositories(c.Request.Context(), claims.GitHubToken)
	if err != nil {
		log.Printf("ERROR - Failed to fetch repositories: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to fetch repositories", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Repositories fetched successfully", gin.H{
		"user":                claims.Username,
		"repositories_by_org": repositoriesByOrg,
		"organization_count":  len(repositoriesByOrg),
	})
}

// generateRandomState generates a random state string for CSRF protection
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
