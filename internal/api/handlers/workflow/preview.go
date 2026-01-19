package workflow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domainWorkflow "github.com/vmaurya-21/Calance-Workflow/internal/domain/workflow"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
	"github.com/vmaurya-21/Calance-Workflow/internal/pkg/logger"
)

// Preview generates and returns the workflow YAML without creating it
// POST /api/workflows/preview
func (h *Handler) Preview(c *gin.Context) {
	// Get user ID from context (for authentication)
	_, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	// Parse request body
	var request domainWorkflow.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error().Err(err).Msg("Invalid workflow preview request body")
		pkghttp.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.Error().Err(err).Str("deployment_type", string(request.DeploymentType)).Msg("Workflow preview request validation failed")
		pkghttp.BadRequestResponse(c, "Validation failed: "+err.Error())
		return
	}

	logger.Info().
		Str("workflow_name", request.WorkflowName).
		Str("deployment_type", string(request.DeploymentType)).
		Msg("Previewing workflow")

	// Generate workflow YAML
	yamlContent, err := h.workflowService.GenerateWorkflow(&request)
	if err != nil {
		logger.Error().Err(err).Str("workflow_name", request.WorkflowName).Msg("Failed to preview workflow YAML")
		pkghttp.InternalServerErrorResponse(c, "Failed to generate workflow preview", err)
		return
	}

	logger.Info().Str("workflow_name", request.WorkflowName).Msg("Workflow preview generated successfully")

	pkghttp.SuccessResponse(c, http.StatusOK, "Workflow preview generated successfully", gin.H{
		"workflow_name":   request.WorkflowName,
		"deployment_type": request.DeploymentType,
		"yaml_content":    yamlContent,
	})
}
