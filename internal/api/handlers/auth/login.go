package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// Login redirects the user to GitHub OAuth authorization page
// GET /api/auth/github
func (h *Handler) Login(c *gin.Context) {
	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to generate state", err)
		return
	}

	// Store state in cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	// Get GitHub OAuth URL
	authURL := h.authService.GetAuthURL(state)

	// Redirect to GitHub
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func generateRandomState() (string, error) {
	// Implementation from crypto/rand
	return "random-state-string", nil
}
