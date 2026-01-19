package auth

import (
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/auth"
	database "github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/database/repositories"
)

// Handler handles auth-related HTTP requests
type Handler struct {
	authService     *auth.Service
	userRepository  *database.UserRepository
	tokenRepository *database.TokenRepository
}

// NewHandler creates a new auth handler
func NewHandler(
	authService *auth.Service,
	userRepo *database.UserRepository,
	tokenRepo *database.TokenRepository,
) *Handler {
	return &Handler{
		authService:     authService,
		userRepository:  userRepo,
		tokenRepository: tokenRepo,
	}
}

// getUserID extracts user ID from context
func (h *Handler) getUserID(c interface{}) (uuid.UUID, error) {
	// This will be implemented with gin.Context
	return uuid.Nil, nil
}
