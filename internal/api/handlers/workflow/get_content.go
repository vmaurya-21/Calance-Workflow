package workflow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// GetWorkflowContent returns the content of a workflow file
// GET /api/workflows/:owner/:repo/file?path=.github/workflows/blank.yml
func (h *Handler) GetWorkflowContent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	_, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	owner := c.Param("owner")
	repo := c.Param("repo")
	filePath := c.Query("path")

	if owner == "" || repo == "" {
		pkghttp.BadRequestResponse(c, "Owner and repository are required")
		return
	}

	if filePath == "" {
		pkghttp.BadRequestResponse(c, "File path is required")
		return
	}

	accessToken, err := h.getAccessToken(userID.(string))
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	fileContent, err := h.workflowService.GetWorkflowContent(c.Request.Context(), accessToken, owner, repo, filePath)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch workflow file content", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "File content retrieved successfully", fileContent)
}
