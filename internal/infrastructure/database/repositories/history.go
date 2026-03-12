package database

import (
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/workflow"
	"gorm.io/gorm"
)

// HistoryRepository handles workflow history database operations
type HistoryRepository struct {
	db *gorm.DB
}

// NewHistoryRepository creates a new history repository
func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

// Create creates a new history record
func (r *HistoryRepository) Create(history *workflow.History) error {
	return r.db.Create(history).Error
}

// ListByRepository retrieves history records for a specific repository
func (r *HistoryRepository) ListByRepository(owner, repository string) ([]workflow.History, error) {
	var history []workflow.History
	err := r.db.Where("owner = ? AND repository = ?", owner, repository).
		Order("created_at desc").
		Find(&history).Error
	return history, err
}

// ListByUser retrieves history records for a specific user
func (r *HistoryRepository) ListByUser(userID string) ([]workflow.History, error) {
	var history []workflow.History
	err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&history).Error
	return history, err
}
