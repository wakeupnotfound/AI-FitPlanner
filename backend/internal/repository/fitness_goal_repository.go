package repository

import (
	"context"
	"errors"

	"github.com/ai-fitness-planner/backend/internal/model"
	"gorm.io/gorm"
)

// FitnessGoalRepository defines the interface for fitness goal operations
type FitnessGoalRepository interface {
	Create(ctx context.Context, goal *model.FitnessGoal) error
	GetByUserID(ctx context.Context, userID int64, status string) ([]*model.FitnessGoal, error)
	Update(ctx context.Context, goal *model.FitnessGoal) error
	Delete(ctx context.Context, id int64) error
}

// fitnessGoalRepository implements FitnessGoalRepository interface
type fitnessGoalRepository struct {
	db *gorm.DB
}

// NewFitnessGoalRepository creates a new instance of FitnessGoalRepository
func NewFitnessGoalRepository(db *gorm.DB) FitnessGoalRepository {
	return &fitnessGoalRepository{db: db}
}

// Create creates a new fitness goal
func (r *fitnessGoalRepository) Create(ctx context.Context, goal *model.FitnessGoal) error {
	if err := r.db.WithContext(ctx).Create(goal).Error; err != nil {
		return err
	}
	return nil
}

// GetByUserID retrieves fitness goals for a user, optionally filtered by status
func (r *fitnessGoalRepository) GetByUserID(ctx context.Context, userID int64, status string) ([]*model.FitnessGoal, error) {
	var goals []*model.FitnessGoal
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("priority DESC, created_at DESC").Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

// Update updates an existing fitness goal
func (r *fitnessGoalRepository) Update(ctx context.Context, goal *model.FitnessGoal) error {
	if err := r.db.WithContext(ctx).Save(goal).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a fitness goal
func (r *fitnessGoalRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.FitnessGoal{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}
