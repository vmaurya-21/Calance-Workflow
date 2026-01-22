package workflow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/workflow"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// UpdateWorkflow updates an existing workflow file and creates a PR
// PUT /api/workflows/:owner/:repo/file
func (h *Handler) UpdateWorkflow(c *gin.Context) {
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

	var req workflow.UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate path parameters match request body
	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner != req.Owner || repo != req.Repository {
		pkghttp.BadRequestResponse(c, "Path parameters must match request body")
		return
	}

	accessToken, err := h.getAccessToken(userID.(string))
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	response, err := h.workflowService.UpdateWorkflow(c.Request.Context(), accessToken, &req)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to update workflow", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, response.Message, response)
}
