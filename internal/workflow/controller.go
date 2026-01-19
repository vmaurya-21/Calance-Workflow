package workflow

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/logger"
	"github.com/vmaurya-21/Calance-Workflow/internal/repositories"
)

// WorkflowController handles workflow-related HTTP requests
type WorkflowController struct {
	workflowService *WorkflowService
	tokenRepository *repositories.TokenRepository
}

// NewWorkflowController creates a new workflow controller
func NewWorkflowController(
	workflowService *WorkflowService,
	tokenRepository *repositories.TokenRepository,
) *WorkflowController {
	return &WorkflowController{
		workflowService: workflowService,
		tokenRepository: tokenRepository,
	}
}

// CreateWorkflow creates a new workflow file in a GitHub repository
// POST /api/workflows/create
func (wc *WorkflowController) CreateWorkflow(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found in context",
		})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID format",
			"error":   err.Error(),
		})
		return
	}

	// Parse request body
	var request WorkflowRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error().
			Err(err).
			Msg("Invalid workflow request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.Error().
			Err(err).
			Str("deployment_type", string(request.DeploymentType)).
			Msg("Workflow request validation failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	// Fetch token from database
	token, err := wc.tokenRepository.FindByUserID(userUUID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userUUID.String()).
			Msg("Failed to fetch access token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch access token",
			"error":   err.Error(),
		})
		return
	}
	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Access token not found. Please login again.",
		})
		return
	}

	logger.Info().
		Str("owner", request.Owner).
		Str("repo", request.Repository).
		Str("workflow_name", request.WorkflowName).
		Str("deployment_type", string(request.DeploymentType)).
		Msg("Creating workflow")

	// Generate workflow YAML
	yamlContent, err := wc.workflowService.GenerateWorkflowYAML(&request)
	if err != nil {
		logger.Error().
			Err(err).
			Str("workflow_name", request.WorkflowName).
			Msg("Failed to generate workflow YAML")

		// Check if it's a validation error
		statusCode := http.StatusInternalServerError
		message := "Failed to generate workflow"

		if err == ErrInvalidWorkflowName {
			statusCode = http.StatusBadRequest
			message = "Invalid workflow name. Use only alphanumeric characters, hyphens, and underscores."
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	// Create workflow file in GitHub repository
	response, err := wc.workflowService.CreateWorkflowFile(
		c.Request.Context(),
		token.AccessToken,
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

		// Handle specific errors
		statusCode := http.StatusInternalServerError
		message := "Failed to create workflow file"
		errMsg := err.Error()

		// Check for specific error types
		if err == ErrInsufficientPermissions {
			statusCode = http.StatusForbidden
			message = "Insufficient permissions to create workflow. Ensure your GitHub token has 'workflow' scope."
		} else if err == ErrWorkflowAlreadyExists {
			statusCode = http.StatusConflict
			message = "Workflow file already exists"
		} else if strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "Not Found") {
			// GitHub returns 404 for both missing repos AND missing workflow permissions
			statusCode = http.StatusNotFound
			message = fmt.Sprintf("Repository '%s/%s' not found or your GitHub token lacks 'workflow' scope. Please verify: 1) The repository exists, 2) Your token has access to it, 3) Your token has 'workflow' scope enabled.", request.Owner, request.Repository)
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	logger.Info().
		Str("owner", request.Owner).
		Str("repo", request.Repository).
		Str("workflow_name", request.WorkflowName).
		Str("file_path", response.FilePath).
		Msg("Workflow created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": response.Message,
		"data":    response,
	})
}

// PreviewWorkflow generates and returns the workflow YAML without creating it
// POST /api/workflows/preview
func (wc *WorkflowController) PreviewWorkflow(c *gin.Context) {
	// Get user ID from context (for authentication)
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found in context",
		})
		return
	}

	// Parse request body
	var request WorkflowRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error().
			Err(err).
			Msg("Invalid workflow preview request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.Error().
			Err(err).
			Str("deployment_type", string(request.DeploymentType)).
			Msg("Workflow preview request validation failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	logger.Info().
		Str("workflow_name", request.WorkflowName).
		Str("deployment_type", string(request.DeploymentType)).
		Msg("Previewing workflow")

	// Generate workflow YAML
	yamlContent, err := wc.workflowService.PreviewWorkflowYAML(&request)
	if err != nil {
		logger.Error().
			Err(err).
			Str("workflow_name", request.WorkflowName).
			Msg("Failed to preview workflow YAML")

		// Check if it's a validation error
		statusCode := http.StatusInternalServerError
		message := "Failed to generate workflow preview"

		if err == ErrInvalidWorkflowName {
			statusCode = http.StatusBadRequest
			message = "Invalid workflow name. Use only alphanumeric characters, hyphens, and underscores."
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	logger.Info().
		Str("workflow_name", request.WorkflowName).
		Msg("Workflow preview generated successfully")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Workflow preview generated successfully",
		"data": gin.H{
			"workflow_name":   request.WorkflowName,
			"deployment_type": request.DeploymentType,
			"yaml_content":    yamlContent,
		},
	})
}

// GetWorkflows retrieves all workflow files from a repository
// GET /api/workflows/:owner/:repo
func (wc *WorkflowController) GetWorkflows(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found in context",
		})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID format",
			"error":   err.Error(),
		})
		return
	}

	// Get owner and repo from URL parameters
	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Owner and repository name are required",
		})
		return
	}

	// Fetch token from database
	token, err := wc.tokenRepository.FindByUserID(userUUID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userUUID.String()).
			Msg("Failed to fetch access token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch access token",
			"error":   err.Error(),
		})
		return
	}
	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Access token not found. Please login again.",
		})
		return
	}

	logger.Info().
		Str("owner", owner).
		Str("repo", repo).
		Msg("Fetching workflows")

	// Get workflows from repository
	workflows, err := wc.workflowService.GetWorkflows(c.Request.Context(), token.AccessToken, owner, repo)
	if err != nil {
		logger.Error().
			Err(err).
			Str("owner", owner).
			Str("repo", repo).
			Msg("Failed to fetch workflows")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch workflows",
			"error":   err.Error(),
		})
		return
	}

	logger.Info().
		Str("owner", owner).
		Str("repo", repo).
		Int("count", len(workflows)).
		Msg("Workflows fetched successfully")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Successfully retrieved %d workflow(s)", len(workflows)),
		"data": gin.H{
			"owner":      owner,
			"repository": repo,
			"workflows":  workflows,
			"count":      len(workflows),
		},
	})
}
