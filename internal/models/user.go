package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	GitHubID  int64          `gorm:"uniqueIndex;not null" json:"github_id"`
	Username  string         `gorm:"not null" json:"username"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	AvatarURL string         `json:"avatar_url"`
	Name      string         `json:"name"`
	Bio       string         `json:"bio"`
	Location  string         `json:"location"`
	Company   string         `json:"company"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID before creating user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// UserResponse is the response structure for user data (excludes sensitive fields)
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	GitHubID  int64     `json:"github_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	Location  string    `json:"location"`
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		GitHubID:  u.GitHubID,
		Username:  u.Username,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Name:      u.Name,
		Bio:       u.Bio,
		Location:  u.Location,
		Company:   u.Company,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
