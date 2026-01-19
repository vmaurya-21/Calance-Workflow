package workflow

import (
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/workflow"
	database "github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/database/repositories"
)

// Handler handles workflow-related HTTP requests
type Handler struct {
	workflowService *workflow.Service
	tokenRepository *database.TokenRepository
}

// NewHandler creates a new workflow handler
func NewHandler(
	workflowService *workflow.Service,
	tokenRepo *database.TokenRepository,
) *Handler {
	return &Handler{
		workflowService: workflowService,
		tokenRepository: tokenRepo,
	}
}

// getUserID extracts user ID from context
func (h *Handler) getUserID(c interface{}) (uuid.UUID, error) {
	return uuid.Nil, nil
}

// getAccessToken retrieves access token for user
func (h *Handler) getAccessToken(userID uuid.UUID) (string, error) {
	token, err := h.tokenRepository.FindByUserID(userID)
	if err != nil || token == nil {
		return "", err
	}
	return token.AccessToken, nil
}
