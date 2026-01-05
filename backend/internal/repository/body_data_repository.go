package repository

import (
	"context"
	"errors"

	"github.com/ai-fitness-planner/backend/internal/model"
	"gorm.io/gorm"
)

// BodyDataRepository defines the interface for body data operations
type BodyDataRepository interface {
	Create(ctx context.Context, bodyData *model.UserBodyData) error
	GetByUserID(ctx context.Context, userID int64) ([]*model.UserBodyData, error)
	GetLatestByUserID(ctx context.Context, userID int64) (*model.UserBodyData, error)
}

// bodyDataRepository implements BodyDataRepository interface
type bodyDataRepository struct {
	db *gorm.DB
}

// NewBodyDataRepository creates a new instance of BodyDataRepository
func NewBodyDataRepository(db *gorm.DB) BodyDataRepository {
	return &bodyDataRepository{db: db}
}

// Create creates a new body data record
func (r *bodyDataRepository) Create(ctx context.Context, bodyData *model.UserBodyData) error {
	if err := r.db.WithContext(ctx).Create(bodyData).Error; err != nil {
		return err
	}
	return nil
}

// GetByUserID retrieves all body data records for a user, ordered by measurement date descending
func (r *bodyDataRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.UserBodyData, error) {
	var bodyDataList []*model.UserBodyData
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("measurement_date DESC").
		Find(&bodyDataList).Error; err != nil {
		return nil, err
	}
	return bodyDataList, nil
}

// GetLatestByUserID retrieves the most recent body data record for a user
func (r *bodyDataRepository) GetLatestByUserID(ctx context.Context, userID int64) (*model.UserBodyData, error) {
	var bodyData model.UserBodyData
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("measurement_date DESC").
		First(&bodyData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bodyData, nil
}
