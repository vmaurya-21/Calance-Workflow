package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByGitHubID finds a user by GitHub ID
func (r *UserRepository) FindByGitHubID(githubID int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("git_hub_id = ?", githubID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by UUID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// CreateOrUpdate creates a new user or updates if exists
func (r *UserRepository) CreateOrUpdate(user *models.User) error {
	existingUser, err := r.FindByGitHubID(user.GitHubID)
	if err != nil {
		return err
	}

	if existingUser != nil {
		// Update existing user
		user.ID = existingUser.ID
		user.CreatedAt = existingUser.CreatedAt
		return r.Update(user)
	}

	// Create new user
	return r.Create(user)
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List returns all users with pagination
func (r *UserRepository) List(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Count returns total number of users
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}
