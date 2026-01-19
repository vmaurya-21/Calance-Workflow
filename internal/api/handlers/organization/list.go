package organization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// List returns all GitHub organizations for the authenticated user
// GET /api/organizations
func (h *Handler) List(c *gin.Context) {
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

	organizations, err := h.organizationService.GetUserOrganizations(c.Request.Context(), accessToken)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch organizations", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Organizations fetched successfully", gin.H{
		"organizations":      organizations,
		"organization_count": len(organizations),
	})
}
