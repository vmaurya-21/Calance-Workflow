package database

import (
	"github.com/vmaurya-21/Calance-Workflow/internal/domain/workflow"
	"gorm.io/gorm"
)

// TemplateRepository handles workflow template database operations
type TemplateRepository struct {
	db *gorm.DB
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

// Create creates a new template
func (r *TemplateRepository) Create(template *workflow.WorkflowTemplate) error {
	return r.db.Create(template).Error
}

// List retrieves all templates as summaries
func (r *TemplateRepository) List() ([]workflow.TemplateSummary, error) {
	var summaries []workflow.TemplateSummary
	err := r.db.Model(&workflow.WorkflowTemplate{}).
		Select("template_id, name, version, description").
		Order("name asc").
		Find(&summaries).Error
	return summaries, err
}

// GetByID retrieves a full template by its template ID
func (r *TemplateRepository) GetByID(templateID string) (*workflow.WorkflowTemplate, error) {
	var template workflow.WorkflowTemplate
	err := r.db.Where("template_id = ?", templateID).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update updates an existing template
func (r *TemplateRepository) Update(template *workflow.WorkflowTemplate) error {
	return r.db.Save(template).Error
}
