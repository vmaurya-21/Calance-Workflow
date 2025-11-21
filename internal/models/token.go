package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Token represents an OAuth access token for a user
// One-to-one relationship with User
type Token struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"` // Unique constraint ensures one-to-one
	AccessToken string         `gorm:"type:text;not null" json:"-"`                   // Don't expose in JSON responses
	TokenType   string         `gorm:"type:varchar(50);default:'Bearer'" json:"token_type"`
	Scope       string         `gorm:"type:text" json:"scope"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship to User
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// BeforeCreate hook to generate UUID before creating token
func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TokenResponse is the response structure for token data (excludes sensitive fields)
type TokenResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenType string     `json:"token_type"`
	Scope     string     `json:"scope"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse converts Token to TokenResponse (excludes access_token)
func (t *Token) ToResponse() TokenResponse {
	return TokenResponse{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenType: t.TokenType,
		Scope:     t.Scope,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
