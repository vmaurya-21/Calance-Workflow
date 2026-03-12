package workflow

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
	"github.com/vmaurya-21/Calance-Workflow/internal/pkg/logger"
)

// GetRepositoryPRs retrieves all pull requests created by this application for a specific repository
// GET /api/workflows/:owner/:repo/prs
func (h *Handler) GetRepositoryPRs(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	// Get owner and repo from URL parameters
	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		pkghttp.BadRequestResponse(c, "Owner and repository name are required")
		return
	}

	// Fetch token from database
	accessToken, err := h.getAccessToken(userID.(string))
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	logger.Info().Str("owner", owner).Str("repo", repo).Msg("Fetching repository PRs")

	// Get PRs from repository
	prs, err := h.workflowService.GetRepositoryPRs(c.Request.Context(), accessToken, owner, repo)
	if err != nil {
		logger.Error().Err(err).Str("owner", owner).Str("repo", repo).Msg("Failed to fetch repository PRs")
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch repository PRs", err)
		return
	}

	logger.Info().Str("owner", owner).Str("repo", repo).Int("count", len(prs)).Msg("Repository PRs fetched successfully")

	pkghttp.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Successfully retrieved %d PR(s)", len(prs)), gin.H{
		"owner":      owner,
		"repository": repo,
		"prs":        prs,
		"count":      len(prs),
	})
}
