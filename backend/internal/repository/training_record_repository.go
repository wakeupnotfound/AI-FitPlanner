package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ai-fitness-planner/backend/internal/model"
	"gorm.io/gorm"
)

// TrainingRecordRepository defines the interface for training record operations
type TrainingRecordRepository interface {
	Create(ctx context.Context, record *model.TrainingRecord) error
	GetByID(ctx context.Context, id int64) (*model.TrainingRecord, error)
	ListByUser(ctx context.Context, userID int64, startDate, endDate *time.Time) ([]*model.TrainingRecord, error)
	GetStatistics(ctx context.Context, userID int64, startDate, endDate time.Time) (*TrainingStatistics, error)
}

// TrainingStatistics represents aggregated training statistics
type TrainingStatistics struct {
	TotalWorkouts  int64
	TotalDuration  int64
	TotalCalories  int64
	AverageRating  float64
	WorkoutsByType map[string]int64
}

// trainingRecordRepository implements TrainingRecordRepository interface
type trainingRecordRepository struct {
	db *gorm.DB
}

// NewTrainingRecordRepository creates a new instance of TrainingRecordRepository
func NewTrainingRecordRepository(db *gorm.DB) TrainingRecordRepository {
	return &trainingRecordRepository{db: db}
}

// Create creates a new training record with validation
func (r *trainingRecordRepository) Create(ctx context.Context, record *model.TrainingRecord) error {
	// Validate that workout date is not in the future
	if record.WorkoutDate.After(time.Now()) {
		return errors.New("workout date cannot be in the future")
	}

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a training record by ID
func (r *trainingRecordRepository) GetByID(ctx context.Context, id int64) (*model.TrainingRecord, error) {
	var record model.TrainingRecord
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// ListByUser retrieves training records for a user within an optional date range
func (r *trainingRecordRepository) ListByUser(ctx context.Context, userID int64, startDate, endDate *time.Time) ([]*model.TrainingRecord, error) {
	var records []*model.TrainingRecord
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("workout_date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("workout_date <= ?", *endDate)
	}

	if err := query.Order("workout_date DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// GetStatistics calculates aggregated statistics for a user's training records
func (r *trainingRecordRepository) GetStatistics(ctx context.Context, userID int64, startDate, endDate time.Time) (*TrainingStatistics, error) {
	stats := &TrainingStatistics{
		WorkoutsByType: make(map[string]int64),
	}

	// Get total workouts count
	if err := r.db.WithContext(ctx).
		Model(&model.TrainingRecord{}).
		Where("user_id = ? AND workout_date >= ? AND workout_date <= ?", userID, startDate, endDate).
		Count(&stats.TotalWorkouts).Error; err != nil {
		return nil, err
	}

	// Get sum of duration and calculate average rating
	type AggregateResult struct {
		TotalDuration int64
		AvgRating     float64
	}

	var result AggregateResult
	if err := r.db.WithContext(ctx).
		Model(&model.TrainingRecord{}).
		Select("COALESCE(SUM(duration_minutes), 0) as total_duration, COALESCE(AVG(rating), 0) as avg_rating").
		Where("user_id = ? AND workout_date >= ? AND workout_date <= ?", userID, startDate, endDate).
		Scan(&result).Error; err != nil {
		return nil, err
	}

	stats.TotalDuration = result.TotalDuration
	stats.AverageRating = result.AvgRating

	// Calculate total calories from performance_data JSON field
	var records []*model.TrainingRecord
	if err := r.db.WithContext(ctx).
		Select("performance_data").
		Where("user_id = ? AND workout_date >= ? AND workout_date <= ?", userID, startDate, endDate).
		Find(&records).Error; err != nil {
		return nil, err
	}

	var totalCalories int64
	for _, record := range records {
		if record.PerformanceData != nil {
			if calories, ok := record.PerformanceData["estimated_calories"]; ok {
				switch v := calories.(type) {
				case float64:
					totalCalories += int64(v)
				case int64:
					totalCalories += v
				case int:
					totalCalories += int64(v)
				}
			}
		}
	}
	stats.TotalCalories = totalCalories

	// Get workouts grouped by type
	type TypeCount struct {
		WorkoutType string
		Count       int64
	}

	var typeCounts []TypeCount
	if err := r.db.WithContext(ctx).
		Model(&model.TrainingRecord{}).
		Select("workout_type, COUNT(*) as count").
		Where("user_id = ? AND workout_date >= ? AND workout_date <= ?", userID, startDate, endDate).
		Group("workout_type").
		Scan(&typeCounts).Error; err != nil {
		return nil, err
	}

	for _, tc := range typeCounts {
		stats.WorkoutsByType[tc.WorkoutType] = tc.Count
	}

	return stats, nil
}
