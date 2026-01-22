package workflow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domainWorkflow "github.com/vmaurya-21/Calance-Workflow/internal/domain/workflow"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
	"github.com/vmaurya-21/Calance-Workflow/internal/pkg/logger"
)

// Create creates a new workflow file in a GitHub repository
// POST /api/workflows/create
func (h *Handler) Create(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	// Parse request body
	var request domainWorkflow.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error().Err(err).Msg("Invalid workflow request body")
		pkghttp.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.Error().Err(err).Str("deployment_type", string(request.DeploymentType)).Msg("Workflow request validation failed")
		pkghttp.BadRequestResponse(c, "Validation failed: "+err.Error())
		return
	}

	// Fetch token from database
	accessToken, err := h.getAccessToken(userID.(string))
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	logger.Info().
		Str("owner", request.Owner).
		Str("repo", request.Repository).
		Str("workflow_name", request.WorkflowName).
		Str("deployment_type", string(request.DeploymentType)).
		Msg("Creating workflow")

	// Generate workflow YAML
	yamlContent, err := h.workflowService.GenerateWorkflow(&request)
	if err != nil {
		logger.Error().Err(err).Str("workflow_name", request.WorkflowName).Msg("Failed to generate workflow YAML")
		pkghttp.InternalServerErrorResponse(c, "Failed to generate workflow", err)
		return
	}

	// Create workflow file in GitHub repository
	response, err := h.workflowService.CreateWorkflow(
		c.Request.Context(),
		accessToken,
		request.Owner,
		request.Repository,
		request.WorkflowName,
		yamlContent,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Str("owner", request.Owner).
			Str("repo", request.Repository).
			Str("workflow_name", request.WorkflowName).
			Msg("Failed to create workflow file")
		pkghttp.InternalServerErrorResponse(c, "Failed to create workflow file", err)
		return
	}

	logger.Info().
		Str("owner", request.Owner).
		Str("repo", request.Repository).
		Str("workflow_name", request.WorkflowName).
		Str("file_path", response.FilePath).
		Msg("Workflow created successfully")

	pkghttp.SuccessResponse(c, http.StatusCreated, response.Message, response)
}
