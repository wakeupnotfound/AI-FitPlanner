package response

import "github.com/ai-fitness-planner/backend/internal/model"

// 通用响应类型
type TaskResponse struct {
	TaskID        string      `json:"task_id"`
	Status        string      `json:"status"`
	Progress      int         `json:"progress"`
	EstimatedTime int         `json:"estimated_time"`
	Result        interface{} `json:"result,omitempty"`
	ErrorMessage  string      `json:"error_message,omitempty"`
}

type PlanListResponse struct {
	Plans      []PlanInfo     `json:"plans"`
	Pagination PaginationInfo `json:"pagination"`
}

type PlanInfo struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	TotalWeeks      int    `json:"total_weeks"`
	DifficultyLevel string `json:"difficulty_level"`
	Status          string `json:"status"`
}

type PlanDetailResponse struct {
	Plan PlanDetailInfo `json:"plan"`
}

type PlanDetailInfo struct {
	ID              int64                  `json:"id"`
	Name            string                 `json:"name"`
	StartDate       string                 `json:"start_date"`
	EndDate         string                 `json:"end_date"`
	TotalWeeks      int                    `json:"total_weeks"`
	DifficultyLevel string                 `json:"difficulty_level"`
	TrainingPurpose string                 `json:"training_purpose,omitempty"`
	PlanData        map[string]interface{} `json:"plan_data"`
	Status          string                 `json:"status"`
	CreatedAt       string                 `json:"created_at"`
}

type TodayTrainingResponse struct {
	Schedule TodaySchedule `json:"schedule"`
}

type TodaySchedule struct {
	Date               string         `json:"date"`
	Type               string         `json:"type"`
	FocusArea          string         `json:"focus_area"`
	Exercises          []ExerciseInfo `json:"exercises"`
	Duration           int            `json:"duration"`
	IsCompleted        bool           `json:"is_completed"`
	CompletedExercises int            `json:"completed_exercises"`
	TotalExercises     int            `json:"total_exercises"`
}

type ExerciseInfo struct {
	Name        string `json:"name"`
	Sets        int    `json:"sets"`
	Reps        string `json:"reps"`
	Weight      string `json:"weight"`
	Rest        string `json:"rest"`
	Difficulty  string `json:"difficulty"`
	SafetyNotes string `json:"safety_notes"`
}

type TodayNutritionResponse struct {
	Plan  TodayNutritionPlanInfo `json:"plan"`
	Meals map[string]MealInfo    `json:"meals"`
}

type TodayNutritionPlanInfo struct {
	TargetCalories float64 `json:"target_calories"`
	TargetProtein  float64 `json:"target_protein"`
	TargetCarbs    float64 `json:"target_carbs"`
	TargetFat      float64 `json:"target_fat"`
}

type MealInfo struct {
	Time          string                    `json:"time"`
	Foods         []model.NutritionFoodItem `json:"foods"`
	TotalCalories float64                   `json:"total_calories"`
}

type TrainingRecordInfo struct {
	ID              int64                  `json:"id"`
	PlanID          int64                  `json:"plan_id,omitempty"`
	WorkoutDate     string                 `json:"workout_date"`
	WorkoutType     string                 `json:"workout_type"`
	DurationMinutes int                    `json:"duration_minutes,omitempty"`
	Exercises       map[string]interface{} `json:"exercises,omitempty"`
	PerformanceData map[string]interface{} `json:"performance_data,omitempty"`
	Notes           string                 `json:"notes,omitempty"`
	Rating          int                    `json:"rating,omitempty"`
	InjuryReport    string                 `json:"injury_report,omitempty"`
	CreatedAt       string                 `json:"created_at"`
}

type TrainingRecordListResponse struct {
	Records    []TrainingRecordInfo `json:"records"`
	Pagination PaginationInfo       `json:"pagination"`
}
