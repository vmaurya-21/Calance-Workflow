package database

import (
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/auth"
	"gorm.io/gorm"
)

// UserRepository handles user data access
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uuid.UUID) (*auth.User, error) {
	var user auth.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByGitHubID finds a user by GitHub ID
func (r *UserRepository) FindByGitHubID(githubID int64) (*auth.User, error) {
	var user auth.User
	if err := r.db.Where("github_id = ?", githubID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateOrUpdate creates or updates a user
func (r *UserRepository) CreateOrUpdate(user *auth.User) error {
	var existing auth.User
	err := r.db.Where("github_id = ?", user.GitHubID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(user).Error
	}

	if err != nil {
		return err
	}

	user.ID = existing.ID
	return r.db.Save(user).Error
}
