package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/models"
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// FindByUserID finds a token by User ID
func (r *TokenRepository) FindByUserID(userID uuid.UUID) (*models.Token, error) {
	var token models.Token
	err := r.db.Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// Create creates a new token
func (r *TokenRepository) Create(token *models.Token) error {
	return r.db.Create(token).Error
}

// Update updates an existing token
func (r *TokenRepository) Update(token *models.Token) error {
	return r.db.Save(token).Error
}

// CreateOrUpdate creates a new token or updates if exists for the user
func (r *TokenRepository) CreateOrUpdate(token *models.Token) error {
	existingToken, err := r.FindByUserID(token.UserID)
	if err != nil {
		return err
	}

	if existingToken != nil {
		// Update existing token
		token.ID = existingToken.ID
		token.CreatedAt = existingToken.CreatedAt
		return r.Update(token)
	}

	// Create new token
	return r.Create(token)
}

// Delete soft deletes a token
func (r *TokenRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Token{}, id).Error
}
