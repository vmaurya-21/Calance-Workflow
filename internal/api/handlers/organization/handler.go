package organization

import (
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/organization"
	database "github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/database/repositories"
)

// Handler handles organization-related HTTP requests
type Handler struct {
	organizationService *organization.Service
	tokenRepository     *database.TokenRepository
}

// NewHandler creates a new organization handler
func NewHandler(
	organizationService *organization.Service,
	tokenRepo *database.TokenRepository,
) *Handler {
	return &Handler{
		organizationService: organizationService,
		tokenRepository:     tokenRepo,
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
