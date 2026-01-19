package database

import (
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/auth"
	"gorm.io/gorm"
)

// TokenRepository handles token data access
type TokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// FindByUserID finds a token by user ID
func (r *TokenRepository) FindByUserID(userID uuid.UUID) (*auth.Token, error) {
	var token auth.Token
	if err := r.db.Where("user_id = ?", userID).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// CreateOrUpdate creates or updates a token
func (r *TokenRepository) CreateOrUpdate(token *auth.Token) error {
	var existing auth.Token
	err := r.db.Where("user_id = ?", token.UserID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(token).Error
	}

	if err != nil {
		return err
	}

	token.ID = existing.ID
	return r.db.Save(token).Error
}
