package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ai-fitness-planner/backend/internal/model"
	"gorm.io/gorm"
)

// TrainingPlanRepository defines the interface for training plan operations
type TrainingPlanRepository interface {
	Create(ctx context.Context, plan *model.TrainingPlan) error
	GetByID(ctx context.Context, id int64) (*model.TrainingPlan, error)
	ListByUser(ctx context.Context, userID int64, status string) ([]*model.TrainingPlan, error)
	Update(ctx context.Context, plan *model.TrainingPlan) error
	Delete(ctx context.Context, id int64) error
	GetTodaySchedule(ctx context.Context, userID int64, date time.Time) (*model.DayPlan, error)
}

// trainingPlanRepository implements TrainingPlanRepository interface
type trainingPlanRepository struct {
	db *gorm.DB
}

// NewTrainingPlanRepository creates a new instance of TrainingPlanRepository
func NewTrainingPlanRepository(db *gorm.DB) TrainingPlanRepository {
	return &trainingPlanRepository{db: db}
}

// Create creates a new training plan
func (r *trainingPlanRepository) Create(ctx context.Context, plan *model.TrainingPlan) error {
	if err := r.db.WithContext(ctx).Create(plan).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a training plan by ID
func (r *trainingPlanRepository) GetByID(ctx context.Context, id int64) (*model.TrainingPlan, error) {
	var plan model.TrainingPlan
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// ListByUser retrieves all training plans for a user, optionally filtered by status
func (r *trainingPlanRepository) ListByUser(ctx context.Context, userID int64, status string) ([]*model.TrainingPlan, error) {
	var plans []*model.TrainingPlan
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

// Update updates an existing training plan
func (r *trainingPlanRepository) Update(ctx context.Context, plan *model.TrainingPlan) error {
	if err := r.db.WithContext(ctx).Save(plan).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a training plan
func (r *trainingPlanRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.TrainingPlan{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetTodaySchedule retrieves the training schedule for a specific date
func (r *trainingPlanRepository) GetTodaySchedule(ctx context.Context, userID int64, date time.Time) (*model.DayPlan, error) {
	var plan model.TrainingPlan

	// Format date as string for comparison (YYYY-MM-DD)
	dateStr := date.Format("2006-01-02")

	// Find active plan that includes this date
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND start_date <= ? AND end_date >= ?",
			userID, "active", date, date).
		First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Parse plan data to find the schedule for the specific date
	planData, ok := plan.PlanData["weeks"]
	if !ok {
		return nil, nil
	}

	weeks, ok := planData.([]interface{})
	if !ok {
		return nil, nil
	}

	// Search through weeks and days to find matching date
	for _, weekInterface := range weeks {
		week, ok := weekInterface.(map[string]interface{})
		if !ok {
			continue
		}

		days, ok := week["days"].([]interface{})
		if !ok {
			continue
		}

		for _, dayInterface := range days {
			dayMap, ok := dayInterface.(map[string]interface{})
			if !ok {
				continue
			}

			dayDate, ok := dayMap["date"].(string)
			if !ok {
				continue
			}

			if dayDate == dateStr {
				// Convert map to DayPlan struct
				dayPlan := &model.DayPlan{}

				if day, ok := dayMap["day"].(float64); ok {
					dayPlan.Day = int(day)
				}
				dayPlan.Date = dayDate
				if typ, ok := dayMap["type"].(string); ok {
					dayPlan.Type = typ
				}
				if focus, ok := dayMap["focus_area"].(string); ok {
					dayPlan.FocusArea = focus
				}
				if duration, ok := dayMap["duration"].(float64); ok {
					dayPlan.Duration = int(duration)
				}
				if calories, ok := dayMap["estimated_calories"].(float64); ok {
					dayPlan.EstimatedCalories = int(calories)
				}

				// Parse exercises
				if exercisesInterface, ok := dayMap["exercises"].([]interface{}); ok {
					exercises := make([]model.Exercise, 0, len(exercisesInterface))
					for _, exInterface := range exercisesInterface {
						exMap, ok := exInterface.(map[string]interface{})
						if !ok {
							continue
						}

						exercise := model.Exercise{}
						if name, ok := exMap["name"].(string); ok {
							exercise.Name = name
						}
						if sets, ok := exMap["sets"].(float64); ok {
							exercise.Sets = int(sets)
						}
						if reps, ok := exMap["reps"].(string); ok {
							exercise.Reps = reps
						}
						if weight, ok := exMap["weight"].(string); ok {
							exercise.Weight = weight
						}
						if rest, ok := exMap["rest"].(string); ok {
							exercise.Rest = rest
						}
						if difficulty, ok := exMap["difficulty"].(string); ok {
							exercise.Difficulty = difficulty
						}
						if safety, ok := exMap["safety_notes"].(string); ok {
							exercise.SafetyNotes = safety
						}

						exercises = append(exercises, exercise)
					}
					dayPlan.Exercises = exercises
				}

				return dayPlan, nil
			}
		}
	}

	return nil, nil
}
