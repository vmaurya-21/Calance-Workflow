package workflow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
	"github.com/vmaurya-21/Calance-Workflow/internal/pkg/logger"
)

// ListTemplates retrieves all available workflow templates
// GET /api/workflows/templates
func (h *Handler) ListTemplates(c *gin.Context) {
	templates, err := h.workflowService.ListTemplates(c.Request.Context())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list workflow templates")
		pkghttp.InternalServerErrorResponse(c, "Failed to list templates", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Templates retrieved successfully", gin.H{
		"templates": templates,
		"count":     len(templates),
	})
}

// GetTemplate retrieves a full workflow template by ID
// GET /api/workflows/templates/:id
func (h *Handler) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		pkghttp.BadRequestResponse(c, "Template ID is required")
		return
	}

	template, err := h.workflowService.GetTemplate(c.Request.Context(), templateID)
	if err != nil {
		logger.Error().Err(err).Str("template_id", templateID).Msg("Failed to fetch workflow template")
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch template", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Template retrieved successfully", template)
}
