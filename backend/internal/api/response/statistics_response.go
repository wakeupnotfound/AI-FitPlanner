package response

// TrainingStatsResponse represents training statistics response
type TrainingStatsResponse struct {
	Period            string           `json:"period"`
	StartDate         string           `json:"start_date"`
	EndDate           string           `json:"end_date"`
	TotalWorkouts     int64            `json:"total_workouts"`
	TotalDuration     int64            `json:"total_duration_minutes"`
	TotalCalories     int64            `json:"total_calories"`
	AverageRating     float64          `json:"average_rating"`
	WorkoutsByType    map[string]int64 `json:"workouts_by_type,omitempty"`
	AverageDuration   float64          `json:"average_duration_minutes"`
	HasSufficientData bool             `json:"has_sufficient_data"`
	Message           string           `json:"message,omitempty"`
}

// ProgressReportResponse represents progress report response
type ProgressReportResponse struct {
	CurrentPeriod     *PeriodSummaryInfo  `json:"current_period"`
	PreviousPeriod    *PeriodSummaryInfo  `json:"previous_period"`
	BodyProgress      *BodyProgressInfo   `json:"body_progress,omitempty"`
	WorkoutComparison *WorkoutCompareInfo `json:"workout_comparison"`
	HasSufficientData bool                `json:"has_sufficient_data"`
	Message           string              `json:"message,omitempty"`
}

// PeriodSummaryInfo represents summary data for a time period
type PeriodSummaryInfo struct {
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	TotalWorkouts int64   `json:"total_workouts"`
	TotalDuration int64   `json:"total_duration_minutes"`
	TotalCalories int64   `json:"total_calories"`
	AverageRating float64 `json:"average_rating"`
}

// BodyProgressInfo represents body measurement progress
type BodyProgressInfo struct {
	CurrentWeight   *float64 `json:"current_weight,omitempty"`
	PreviousWeight  *float64 `json:"previous_weight,omitempty"`
	WeightChange    *float64 `json:"weight_change,omitempty"`
	CurrentBodyFat  *float64 `json:"current_body_fat,omitempty"`
	PreviousBodyFat *float64 `json:"previous_body_fat,omitempty"`
	BodyFatChange   *float64 `json:"body_fat_change,omitempty"`
}

// WorkoutCompareInfo represents workout comparison between periods
type WorkoutCompareInfo struct {
	WorkoutCountChange  int64   `json:"workout_count_change"`
	DurationChange      int64   `json:"duration_change_minutes"`
	CaloriesChange      int64   `json:"calories_change"`
	WorkoutCountPercent float64 `json:"workout_count_percent_change"`
	DurationPercent     float64 `json:"duration_percent_change"`
	CaloriesPercent     float64 `json:"calories_percent_change"`
}

// TrendsReportResponse represents trends report response
type TrendsReportResponse struct {
	Period            string           `json:"period"`
	DataPoints        []TrendPointInfo `json:"data_points"`
	HasSufficientData bool             `json:"has_sufficient_data"`
	Message           string           `json:"message,omitempty"`
}

// TrendPointInfo represents a single data point in the trend
type TrendPointInfo struct {
	PeriodLabel   string  `json:"period_label"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	TotalWorkouts int64   `json:"total_workouts"`
	TotalDuration int64   `json:"total_duration_minutes"`
	TotalCalories int64   `json:"total_calories"`
	AverageRating float64 `json:"average_rating"`
}
