package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	GitHubID  int64          `gorm:"column:github_id;uniqueIndex;not null" json:"github_id"`
	Username  string         `gorm:"not null" json:"username"`
	Email     string         `json:"email"`
	AvatarURL string         `json:"avatar_url"`
	Name      string         `json:"name"`
	Bio       string         `json:"bio"`
	Location  string         `json:"location"`
	Company   string         `json:"company"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// ToResponse converts User to response format
func (u *User) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"github_id":  u.GitHubID,
		"username":   u.Username,
		"email":      u.Email,
		"avatar_url": u.AvatarURL,
		"name":       u.Name,
		"bio":        u.Bio,
		"location":   u.Location,
		"company":    u.Company,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}

// Token represents an access token
type Token struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	AccessToken string         `gorm:"not null" json:"access_token"`
	TokenType   string         `json:"token_type"`
	Scope       string         `json:"scope"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook
func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
