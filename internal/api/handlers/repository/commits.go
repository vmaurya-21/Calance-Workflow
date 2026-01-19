package repository

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// GetCommits returns the latest commits for a specific branch
// GET /api/repositories/:owner/:repo/branches/:branch/commits
func (h *Handler) GetCommits(c *gin.Context) {
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
	branch := c.Param("branch")

	if owner == "" || repo == "" || branch == "" {
		pkghttp.BadRequestResponse(c, "Owner, repository name, and branch are required")
		return
	}

	perPage := 30
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		var parsed int
		if _, err := fmt.Sscanf(perPageStr, "%d", &parsed); err == nil {
			if parsed > 0 && parsed <= 100 {
				perPage = parsed
			}
		}
	}

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	commits, err := h.repositoryService.GetCommits(c.Request.Context(), accessToken, owner, repo, branch, perPage)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch commits", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Commits fetched successfully", gin.H{
		"owner":        owner,
		"repository":   repo,
		"branch":       branch,
		"commits":      commits,
		"commit_count": len(commits),
	})
}
