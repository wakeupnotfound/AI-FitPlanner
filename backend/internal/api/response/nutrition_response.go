package response

// NutritionPlanListResponse represents a list of nutrition plans
type NutritionPlanListResponse struct {
	Plans      []NutritionPlanInfo `json:"plans"`
	Pagination PaginationInfo      `json:"pagination"`
}

// NutritionPlanInfo represents a nutrition plan in responses
type NutritionPlanInfo struct {
	ID                  int64    `json:"id"`
	PlanName            string   `json:"plan_name"`
	StartDate           string   `json:"start_date"`
	EndDate             string   `json:"end_date"`
	DailyCalories       float64  `json:"daily_calories"`
	ProteinRatio        float64  `json:"protein_ratio"`
	CarbRatio           float64  `json:"carb_ratio"`
	FatRatio            float64  `json:"fat_ratio"`
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
	Preferences         []string `json:"preferences,omitempty"`
	Status              string   `json:"status"`
	CreatedAt           string   `json:"created_at"`
}

// NutritionRecordInfo represents a nutrition record in responses
type NutritionRecordInfo struct {
	ID        int64                  `json:"id"`
	MealDate  string                 `json:"meal_date"`
	MealType  string                 `json:"meal_type"`
	Calories  float64                `json:"calories"`
	Protein   float64                `json:"protein"`
	Carbs     float64                `json:"carbs"`
	Fat       float64                `json:"fat"`
	Fiber     float64                `json:"fiber"`
	Foods     map[string]interface{} `json:"foods,omitempty"`
	CreatedAt string                 `json:"created_at"`
}

// NutritionRecordListResponse represents a list of nutrition records
type NutritionRecordListResponse struct {
	Records    []NutritionRecordInfo `json:"records"`
	Pagination PaginationInfo        `json:"pagination"`
}

// DailySummaryResponse represents daily nutrition summary
type DailySummaryResponse struct {
	Date          string  `json:"date"`
	TotalCalories float64 `json:"total_calories"`
	TotalProtein  float64 `json:"total_protein"`
	TotalCarbs    float64 `json:"total_carbs"`
	TotalFat      float64 `json:"total_fat"`
	TotalFiber    float64 `json:"total_fiber"`
	MealCount     int     `json:"meal_count"`
}
