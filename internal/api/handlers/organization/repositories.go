package organization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// GetRepositories returns all repositories for a specific organization
// GET /api/organizations/:org/repositories
func (h *Handler) GetRepositories(c *gin.Context) {
	orgName := c.Param("org")
	if orgName == "" {
		pkghttp.BadRequestResponse(c, "Organization name is required")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	repositories, err := h.organizationService.GetOrganizationRepositories(c.Request.Context(), accessToken, orgName)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch repositories", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Repositories fetched successfully", gin.H{
		"organization":     orgName,
		"repositories":     repositories,
		"repository_count": len(repositories),
	})
}

// GetUserRepositories returns all repositories accessible to the authenticated user from their organizations
// GET /api/repositories
func (h *Handler) GetUserRepositories(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	repositoriesByOrg, err := h.organizationService.GetUserRepositories(c.Request.Context(), accessToken)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch repositories", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Repositories fetched successfully", gin.H{
		"repositories_by_org": repositoriesByOrg,
		"organization_count":  len(repositoriesByOrg),
	})
}
