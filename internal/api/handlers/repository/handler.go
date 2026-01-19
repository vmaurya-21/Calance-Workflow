package repository

import (
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/repository"
	database "github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/database/repositories"
)

// Handler handles repository-related HTTP requests
type Handler struct {
	repositoryService *repository.Service
	tokenRepository   *database.TokenRepository
}

// NewHandler creates a new repository handler
func NewHandler(
	repositoryService *repository.Service,
	tokenRepo *database.TokenRepository,
) *Handler {
	return &Handler{
		repositoryService: repositoryService,
		tokenRepository:   tokenRepo,
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
