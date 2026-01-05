package repository

import (
	"context"
	"errors"

	"github.com/ai-fitness-planner/backend/internal/model"
	"gorm.io/gorm"
)

// AssessmentRepository defines the interface for fitness assessment operations
type AssessmentRepository interface {
	Create(ctx context.Context, assessment *model.FitnessAssessment) error
	GetByID(ctx context.Context, id int64) (*model.FitnessAssessment, error)
	GetLatest(ctx context.Context, userID int64) (*model.FitnessAssessment, error)
	ListByUser(ctx context.Context, userID int64) ([]*model.FitnessAssessment, error)
}

// assessmentRepository implements AssessmentRepository interface
type assessmentRepository struct {
	db *gorm.DB
}

// NewAssessmentRepository creates a new instance of AssessmentRepository
func NewAssessmentRepository(db *gorm.DB) AssessmentRepository {
	return &assessmentRepository{db: db}
}

// Create creates a new fitness assessment
func (r *assessmentRepository) Create(ctx context.Context, assessment *model.FitnessAssessment) error {
	if err := r.db.WithContext(ctx).Create(assessment).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a fitness assessment by ID
func (r *assessmentRepository) GetByID(ctx context.Context, id int64) (*model.FitnessAssessment, error) {
	var assessment model.FitnessAssessment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&assessment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &assessment, nil
}

// GetLatest retrieves the most recent fitness assessment for a user
func (r *assessmentRepository) GetLatest(ctx context.Context, userID int64) (*model.FitnessAssessment, error) {
	var assessment model.FitnessAssessment
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("assessment_date DESC, created_at DESC").
		First(&assessment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &assessment, nil
}

// ListByUser retrieves all fitness assessments for a user
func (r *assessmentRepository) ListByUser(ctx context.Context, userID int64) ([]*model.FitnessAssessment, error) {
	var assessments []*model.FitnessAssessment
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("assessment_date DESC, created_at DESC").
		Find(&assessments).Error; err != nil {
		return nil, err
	}
	return assessments, nil
}
