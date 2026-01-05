package request

// GenerateNutritionPlanRequest represents the request to generate a nutrition plan
type GenerateNutritionPlanRequest struct {
	PlanName            string   `json:"plan_name" binding:"required,min=1,max=200"`
	DurationDays        int      `json:"duration_days" binding:"required,min=1,max=365"`
	DailyCalories       *float64 `json:"daily_calories" binding:"omitempty,min=500,max=10000"`
	ProteinRatio        float64  `json:"protein_ratio" binding:"required,min=0,max=1,macro_ratio"`
	CarbRatio           float64  `json:"carb_ratio" binding:"required,min=0,max=1,macro_ratio"`
	FatRatio            float64  `json:"fat_ratio" binding:"required,min=0,max=1,macro_ratio"`
	DietaryRestrictions []string `json:"dietary_restrictions" binding:"omitempty,dive,min=1,max=100"`
	Preferences         []string `json:"preferences" binding:"omitempty,dive,min=1,max=100"`
	AIAPIID             *int64   `json:"ai_api_id" binding:"omitempty,min=1"`
}

// RecordMealRequest represents the request to record a meal
type RecordMealRequest struct {
	PlanID   *int64                 `json:"plan_id" binding:"omitempty,min=1"`
	MealDate string                 `json:"meal_date" binding:"required,datetime=2006-01-02,future_date"`
	MealType string                 `json:"meal_type" binding:"required,oneof=breakfast lunch dinner snack"`
	Calories float64                `json:"calories" binding:"omitempty,min=0,max=10000"`
	Protein  float64                `json:"protein" binding:"omitempty,min=0,max=1000"`
	Carbs    float64                `json:"carbs" binding:"omitempty,min=0,max=1000"`
	Fat      float64                `json:"fat" binding:"omitempty,min=0,max=1000"`
	Fiber    float64                `json:"fiber" binding:"omitempty,min=0,max=500"`
	Foods    map[string]interface{} `json:"foods" binding:"required"`
	Notes    *string                `json:"notes" binding:"omitempty,max=1000"`
}

// NutritionPlanListParams represents query parameters for listing nutrition plans
type NutritionPlanListParams struct {
	Status string `form:"status" binding:"omitempty,oneof=active completed cancelled"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// NutritionRecordListParams represents query parameters for listing nutrition records
type NutritionRecordListParams struct {
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// DailySummaryParams represents query parameters for daily nutrition summary
type DailySummaryParams struct {
	Date string `form:"date" binding:"required,datetime=2006-01-02"`
}
