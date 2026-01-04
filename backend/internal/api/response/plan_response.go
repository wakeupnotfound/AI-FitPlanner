package response

// 通用响应类型
type TaskResponse struct {
	TaskID         string `json:"task_id"`
	Status         string `json:"status"`
	Progress       int    `json:"progress"`
	EstimatedTime  int    `json:"estimated_time"`
	Result         interface{} `json:"result,omitempty"`
	ErrorMessage   string `json:"error_message,omitempty"`
}

type PlanListResponse struct {
	Plans      []PlanInfo `json:"plans"`
	Pagination PaginationInfo `json:"pagination"`
}

type PlanInfo struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	TotalWeeks     int    `json:"total_weeks"`
	DifficultyLevel string `json:"difficulty_level"`
	Status         string `json:"status"`
}

type TodayTrainingResponse struct {
	Schedule TodaySchedule `json:"schedule"`
}

type TodaySchedule struct {
	Date             string          `json:"date"`
	Type             string          `json:"type"`
	FocusArea        string          `json:"focus_area"`
	Exercises        []ExerciseInfo  `json:"exercises"`
	Duration         int             `json:"duration"`
	IsCompleted      bool            `json:"is_completed"`
	CompletedExercises int           `json:"completed_exercises"`
	TotalExercises   int             `json:"total_exercises"`
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
	Plan   NutritionPlanInfo   `json:"plan"`
	Meals  map[string]MealInfo `json:"meals"`
}

type NutritionPlanInfo struct {
	TargetCalories float64 `json:"target_calories"`
	TargetProtein  float64 `json:"target_protein"`
	TargetCarbs    float64 `json:"target_carbs"`
	TargetFat      float64 `json:"target_fat"`
}

type MealInfo struct {
	Time          string           `json:"time"`
	Foods         []NutritionFoodItem `json:"foods"`
	TotalCalories float64          `json:"total_calories"`
}
