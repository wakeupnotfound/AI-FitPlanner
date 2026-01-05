package request

// GenerateTrainingPlanRequest represents the request to generate a training plan
type GenerateTrainingPlanRequest struct {
	PlanName        string `json:"plan_name" binding:"required,min=1,max=200"`
	DurationWeeks   int    `json:"duration_weeks" binding:"required,min=1,max=52"`
	Goal            string `json:"goal" binding:"required,min=1,max=100"`
	DifficultyLevel string `json:"difficulty_level" binding:"required,oneof=easy medium hard extreme"`
	AIAPIID         *int64 `json:"ai_api_id" binding:"omitempty,min=1"`
}

// RecordTrainingRequest represents the request to record a training session
type RecordTrainingRequest struct {
	PlanID          *int64                 `json:"plan_id" binding:"omitempty,min=1"`
	WorkoutDate     string                 `json:"workout_date" binding:"required,datetime=2006-01-02,future_date"`
	WorkoutType     string                 `json:"workout_type" binding:"required,min=1,max=100"`
	DurationMinutes *int                   `json:"duration_minutes" binding:"omitempty,min=0,max=1440"`
	Exercises       map[string]interface{} `json:"exercises" binding:"required"`
	PerformanceData map[string]interface{} `json:"performance_data"`
	Notes           *string                `json:"notes" binding:"omitempty,max=1000"`
	Rating          *int                   `json:"rating" binding:"omitempty,min=1,max=5"`
	InjuryReport    *string                `json:"injury_report" binding:"omitempty,max=1000"`
}

// TrainingPlanListParams represents query parameters for listing training plans
type TrainingPlanListParams struct {
	Status string `form:"status" binding:"omitempty,oneof=active completed cancelled"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// TrainingRecordListParams represents query parameters for listing training records
type TrainingRecordListParams struct {
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
}
