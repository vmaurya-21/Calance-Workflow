package workflow

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
	"github.com/vmaurya-21/Calance-Workflow/internal/pkg/logger"
)

// List retrieves all workflow files from a repository
// GET /api/workflows/:owner/:repo
func (h *Handler) List(c *gin.Context) {
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

	logger.Info().Str("owner", owner).Str("repo", repo).Msg("Fetching workflows")

	// Get workflows from repository
	workflows, err := h.workflowService.GetWorkflows(c.Request.Context(), accessToken, owner, repo)
	if err != nil {
		logger.Error().Err(err).Str("owner", owner).Str("repo", repo).Msg("Failed to fetch workflows")
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch workflows", err)
		return
	}

	logger.Info().Str("owner", owner).Str("repo", repo).Int("count", len(workflows)).Msg("Workflows fetched successfully")

	pkghttp.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Successfully retrieved %d workflow(s)", len(workflows)), gin.H{
		"owner":      owner,
		"repository": repo,
		"workflows":  workflows,
		"count":      len(workflows),
	})
}
