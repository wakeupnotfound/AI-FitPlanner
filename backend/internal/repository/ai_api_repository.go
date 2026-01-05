package repository

import (
	"context"
	"errors"

	"github.com/ai-fitness-planner/backend/internal/model"
	"gorm.io/gorm"
)

// AIAPIRepository defines the interface for AI API configuration operations
type AIAPIRepository interface {
	Create(ctx context.Context, api *model.AIAPI) error
	GetByID(ctx context.Context, id int64) (*model.AIAPI, error)
	ListByUser(ctx context.Context, userID int64) ([]*model.AIAPI, error)
	Update(ctx context.Context, api *model.AIAPI) error
	Delete(ctx context.Context, id int64) error
	GetDefaultByUser(ctx context.Context, userID int64) (*model.AIAPI, error)
	SetDefault(ctx context.Context, userID int64, apiID int64) error
}

// aiAPIRepository implements AIAPIRepository interface
type aiAPIRepository struct {
	db *gorm.DB
}

// NewAIAPIRepository creates a new instance of AIAPIRepository
func NewAIAPIRepository(db *gorm.DB) AIAPIRepository {
	return &aiAPIRepository{db: db}
}

// Create creates a new AI API configuration
func (r *aiAPIRepository) Create(ctx context.Context, api *model.AIAPI) error {
	if err := r.db.WithContext(ctx).Create(api).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves an AI API configuration by ID
func (r *aiAPIRepository) GetByID(ctx context.Context, id int64) (*model.AIAPI, error) {
	var api model.AIAPI
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&api).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

// ListByUser retrieves all AI API configurations for a user
func (r *aiAPIRepository) ListByUser(ctx context.Context, userID int64) ([]*model.AIAPI, error) {
	var apis []*model.AIAPI
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

// Update updates an existing AI API configuration
func (r *aiAPIRepository) Update(ctx context.Context, api *model.AIAPI) error {
	if err := r.db.WithContext(ctx).Save(api).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes an AI API configuration
func (r *aiAPIRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.AIAPI{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetDefaultByUser retrieves the default AI API configuration for a user
func (r *aiAPIRepository) GetDefaultByUser(ctx context.Context, userID int64) (*model.AIAPI, error) {
	var api model.AIAPI
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&api).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

// SetDefault sets an AI API as the default for a user
// This operation uses a transaction to ensure only one API is marked as default
func (r *aiAPIRepository) SetDefault(ctx context.Context, userID int64, apiID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, unset all defaults for this user
		if err := tx.Model(&model.AIAPI{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Then, set the specified API as default
		if err := tx.Model(&model.AIAPI{}).
			Where("id = ? AND user_id = ?", apiID, userID).
			Update("is_default", true).Error; err != nil {
			return err
		}

		return nil
	})
}
