package repository

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// GetBranches returns all branches for a specific repository
// GET /api/repositories/:owner/:repo/branches
func (h *Handler) GetBranches(c *gin.Context) {
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

	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		pkghttp.BadRequestResponse(c, "Owner and repository name are required")
		return
	}

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	branches, err := h.repositoryService.GetBranches(c.Request.Context(), accessToken, owner, repo)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch branches", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Branches fetched successfully", gin.H{
		"owner":        owner,
		"repository":   repo,
		"branches":     branches,
		"branch_count": len(branches),
	})
}
