package service

import (
	"context"
	"time"

	"github.com/ai-fitness-planner/backend/internal/errors"
	"github.com/ai-fitness-planner/backend/internal/repository"
)

// StatisticsService defines the interface for statistics operations
// Requirements: 10.1, 10.2, 10.3, 10.4
type StatisticsService interface {
	// GetTrainingStatistics calculates total workouts, duration, and calories
	// Requirements: 10.1
	GetTrainingStatistics(ctx context.Context, userID int64, period string) (*TrainingStats, error)
	// GetProgressReport compares current data with historical data
	// Requirements: 10.2
	GetProgressReport(ctx context.Context, userID int64) (*ProgressReport, error)
	// CalculateTrends aggregates data by week or month
	// Requirements: 10.3
	CalculateTrends(ctx context.Context, userID int64, period string, count int) (*TrendsReport, error)
}

// TrainingStats represents aggregated training statistics
// Requirements: 10.1
type TrainingStats struct {
	Period            string           `json:"period"`
	StartDate         time.Time        `json:"start_date"`
	EndDate           time.Time        `json:"end_date"`
	TotalWorkouts     int64            `json:"total_workouts"`
	TotalDuration     int64            `json:"total_duration_minutes"`
	TotalCalories     int64            `json:"total_calories"`
	AverageRating     float64          `json:"average_rating"`
	WorkoutsByType    map[string]int64 `json:"workouts_by_type"`
	AverageDuration   float64          `json:"average_duration_minutes"`
	HasSufficientData bool             `json:"has_sufficient_data"`
	Message           string           `json:"message,omitempty"`
}

// ProgressReport represents a comparison of current vs historical data
// Requirements: 10.2
type ProgressReport struct {
	CurrentPeriod     *PeriodSummary     `json:"current_period"`
	PreviousPeriod    *PeriodSummary     `json:"previous_period"`
	BodyProgress      *BodyProgressData  `json:"body_progress,omitempty"`
	WorkoutComparison *WorkoutComparison `json:"workout_comparison"`
	HasSufficientData bool               `json:"has_sufficient_data"`
	Message           string             `json:"message,omitempty"`
}

// PeriodSummary represents summary data for a time period
type PeriodSummary struct {
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TotalWorkouts int64     `json:"total_workouts"`
	TotalDuration int64     `json:"total_duration_minutes"`
	TotalCalories int64     `json:"total_calories"`
	AverageRating float64   `json:"average_rating"`
}

// BodyProgressData represents body measurement progress
type BodyProgressData struct {
	CurrentWeight   *float64 `json:"current_weight,omitempty"`
	PreviousWeight  *float64 `json:"previous_weight,omitempty"`
	WeightChange    *float64 `json:"weight_change,omitempty"`
	CurrentBodyFat  *float64 `json:"current_body_fat,omitempty"`
	PreviousBodyFat *float64 `json:"previous_body_fat,omitempty"`
	BodyFatChange   *float64 `json:"body_fat_change,omitempty"`
}

// WorkoutComparison represents workout comparison between periods
type WorkoutComparison struct {
	WorkoutCountChange  int64   `json:"workout_count_change"`
	DurationChange      int64   `json:"duration_change_minutes"`
	CaloriesChange      int64   `json:"calories_change"`
	WorkoutCountPercent float64 `json:"workout_count_percent_change"`
	DurationPercent     float64 `json:"duration_percent_change"`
	CaloriesPercent     float64 `json:"calories_percent_change"`
}

// TrendsReport represents trend data over multiple periods
// Requirements: 10.3
type TrendsReport struct {
	Period            string       `json:"period"` // "week" or "month"
	DataPoints        []TrendPoint `json:"data_points"`
	HasSufficientData bool         `json:"has_sufficient_data"`
	Message           string       `json:"message,omitempty"`
}

// TrendPoint represents a single data point in the trend
type TrendPoint struct {
	PeriodLabel   string    `json:"period_label"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TotalWorkouts int64     `json:"total_workouts"`
	TotalDuration int64     `json:"total_duration_minutes"`
	TotalCalories int64     `json:"total_calories"`
	AverageRating float64   `json:"average_rating"`
}

// statisticsService implements StatisticsService interface
type statisticsService struct {
	trainingRecordRepo repository.TrainingRecordRepository
	bodyDataRepo       repository.BodyDataRepository
}

// NewStatisticsService creates a new instance of StatisticsService
func NewStatisticsService(
	trainingRecordRepo repository.TrainingRecordRepository,
	bodyDataRepo repository.BodyDataRepository,
) StatisticsService {
	return &statisticsService{
		trainingRecordRepo: trainingRecordRepo,
		bodyDataRepo:       bodyDataRepo,
	}
}

// GetTrainingStatistics calculates total workouts, duration, and calories
// Requirements: 10.1
// Property 16: Training Statistics Accuracy - calculated total duration should equal sum of individual record durations
func (s *statisticsService) GetTrainingStatistics(ctx context.Context, userID int64, period string) (*TrainingStats, error) {
	startDate, endDate, err := s.calculateDateRange(period)
	if err != nil {
		return nil, err
	}

	stats, err := s.trainingRecordRepo.GetStatistics(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取训练统计失败")
	}

	result := &TrainingStats{
		Period:         period,
		StartDate:      startDate,
		EndDate:        endDate,
		TotalWorkouts:  stats.TotalWorkouts,
		TotalDuration:  stats.TotalDuration,
		TotalCalories:  stats.TotalCalories,
		AverageRating:  stats.AverageRating,
		WorkoutsByType: stats.WorkoutsByType,
	}

	// Calculate average duration
	if stats.TotalWorkouts > 0 {
		result.AverageDuration = float64(stats.TotalDuration) / float64(stats.TotalWorkouts)
		result.HasSufficientData = true
	} else {
		result.HasSufficientData = false
		result.Message = "该时间段内没有训练记录"
	}

	return result, nil
}

// GetProgressReport compares current data with historical data
// Requirements: 10.2
func (s *statisticsService) GetProgressReport(ctx context.Context, userID int64) (*ProgressReport, error) {
	now := time.Now()

	// Current period: last 30 days
	currentEnd := now
	currentStart := now.AddDate(0, 0, -30)

	// Previous period: 30 days before current period
	previousEnd := currentStart.AddDate(0, 0, -1)
	previousStart := previousEnd.AddDate(0, 0, -30)

	// Get current period stats
	currentStats, err := s.trainingRecordRepo.GetStatistics(ctx, userID, currentStart, currentEnd)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取当前周期统计失败")
	}

	// Get previous period stats
	previousStats, err := s.trainingRecordRepo.GetStatistics(ctx, userID, previousStart, previousEnd)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取上一周期统计失败")
	}

	report := &ProgressReport{
		CurrentPeriod: &PeriodSummary{
			StartDate:     currentStart,
			EndDate:       currentEnd,
			TotalWorkouts: currentStats.TotalWorkouts,
			TotalDuration: currentStats.TotalDuration,
			TotalCalories: currentStats.TotalCalories,
			AverageRating: currentStats.AverageRating,
		},
		PreviousPeriod: &PeriodSummary{
			StartDate:     previousStart,
			EndDate:       previousEnd,
			TotalWorkouts: previousStats.TotalWorkouts,
			TotalDuration: previousStats.TotalDuration,
			TotalCalories: previousStats.TotalCalories,
			AverageRating: previousStats.AverageRating,
		},
	}

	// Calculate workout comparison
	report.WorkoutComparison = s.calculateWorkoutComparison(currentStats, previousStats)

	// Get body progress data
	bodyProgress, err := s.getBodyProgress(ctx, userID)
	if err == nil && bodyProgress != nil {
		report.BodyProgress = bodyProgress
	}

	// Check if we have sufficient data
	// Requirements: 10.4 - handle insufficient data cases
	if currentStats.TotalWorkouts == 0 && previousStats.TotalWorkouts == 0 {
		report.HasSufficientData = false
		report.Message = "没有足够的训练数据来生成进度报告"
	} else if currentStats.TotalWorkouts == 0 {
		report.HasSufficientData = true
		report.Message = "当前周期没有训练记录"
	} else if previousStats.TotalWorkouts == 0 {
		report.HasSufficientData = true
		report.Message = "上一周期没有训练记录，无法进行对比"
	} else {
		report.HasSufficientData = true
	}

	return report, nil
}

// CalculateTrends aggregates data by week or month
// Requirements: 10.3
func (s *statisticsService) CalculateTrends(ctx context.Context, userID int64, period string, count int) (*TrendsReport, error) {
	if period != "week" && period != "month" {
		return nil, errors.New(errors.ErrInvalidParam, "period必须是'week'或'month'")
	}

	if count <= 0 || count > 52 {
		count = 12 // Default to 12 periods
	}

	now := time.Now()
	dataPoints := make([]TrendPoint, 0, count)

	for i := count - 1; i >= 0; i-- {
		var startDate, endDate time.Time
		var label string

		if period == "week" {
			// Calculate week boundaries
			endDate = now.AddDate(0, 0, -7*i)
			startDate = endDate.AddDate(0, 0, -6)
			label = startDate.Format("01/02") + " - " + endDate.Format("01/02")
		} else {
			// Calculate month boundaries
			targetMonth := now.AddDate(0, -i, 0)
			startDate = time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, targetMonth.Location())
			endDate = startDate.AddDate(0, 1, -1)
			label = startDate.Format("2006-01")
		}

		stats, err := s.trainingRecordRepo.GetStatistics(ctx, userID, startDate, endDate)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrDatabase, "获取趋势数据失败")
		}

		dataPoints = append(dataPoints, TrendPoint{
			PeriodLabel:   label,
			StartDate:     startDate,
			EndDate:       endDate,
			TotalWorkouts: stats.TotalWorkouts,
			TotalDuration: stats.TotalDuration,
			TotalCalories: stats.TotalCalories,
			AverageRating: stats.AverageRating,
		})
	}

	report := &TrendsReport{
		Period:     period,
		DataPoints: dataPoints,
	}

	// Check if we have sufficient data
	// Requirements: 10.4 - handle insufficient data cases
	hasData := false
	for _, dp := range dataPoints {
		if dp.TotalWorkouts > 0 {
			hasData = true
			break
		}
	}

	if !hasData {
		report.HasSufficientData = false
		report.Message = "没有足够的训练数据来生成趋势报告"
	} else {
		report.HasSufficientData = true
	}

	return report, nil
}

// calculateDateRange calculates start and end dates based on period string
func (s *statisticsService) calculateDateRange(period string) (time.Time, time.Time, error) {
	now := time.Now()
	endDate := now

	var startDate time.Time

	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	case "all":
		// Use a very old date for "all time"
		startDate = time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location())
	default:
		return time.Time{}, time.Time{}, errors.New(errors.ErrInvalidParam, "无效的时间周期，支持: week, month, quarter, year, all")
	}

	return startDate, endDate, nil
}

// calculateWorkoutComparison calculates the comparison between two periods
func (s *statisticsService) calculateWorkoutComparison(current, previous *repository.TrainingStatistics) *WorkoutComparison {
	comparison := &WorkoutComparison{
		WorkoutCountChange: current.TotalWorkouts - previous.TotalWorkouts,
		DurationChange:     current.TotalDuration - previous.TotalDuration,
		CaloriesChange:     current.TotalCalories - previous.TotalCalories,
	}

	// Calculate percentage changes
	if previous.TotalWorkouts > 0 {
		comparison.WorkoutCountPercent = float64(comparison.WorkoutCountChange) / float64(previous.TotalWorkouts) * 100
	}
	if previous.TotalDuration > 0 {
		comparison.DurationPercent = float64(comparison.DurationChange) / float64(previous.TotalDuration) * 100
	}
	if previous.TotalCalories > 0 {
		comparison.CaloriesPercent = float64(comparison.CaloriesChange) / float64(previous.TotalCalories) * 100
	}

	return comparison
}

// getBodyProgress retrieves body measurement progress
func (s *statisticsService) getBodyProgress(ctx context.Context, userID int64) (*BodyProgressData, error) {
	bodyDataList, err := s.bodyDataRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(bodyDataList) == 0 {
		return nil, nil
	}

	progress := &BodyProgressData{}

	// Get current (most recent) body data
	current := bodyDataList[0]
	progress.CurrentWeight = &current.Weight
	if current.BodyFatPercentage != nil {
		progress.CurrentBodyFat = current.BodyFatPercentage
	}

	// Find previous body data (at least 7 days old)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	for _, bd := range bodyDataList[1:] {
		if bd.MeasurementDate.Before(sevenDaysAgo) {
			progress.PreviousWeight = &bd.Weight
			if bd.BodyFatPercentage != nil {
				progress.PreviousBodyFat = bd.BodyFatPercentage
			}
			break
		}
	}

	// Calculate changes
	if progress.CurrentWeight != nil && progress.PreviousWeight != nil {
		change := *progress.CurrentWeight - *progress.PreviousWeight
		progress.WeightChange = &change
	}
	if progress.CurrentBodyFat != nil && progress.PreviousBodyFat != nil {
		change := *progress.CurrentBodyFat - *progress.PreviousBodyFat
		progress.BodyFatChange = &change
	}

	return progress, nil
}
