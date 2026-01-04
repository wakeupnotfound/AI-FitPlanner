package model

import (
	"time"
)

type NutritionPlan struct {
	ID                int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int64           `gorm:"not null;index" json:"user_id"`
	PlanName          string          `gorm:"size:200;not null" json:"plan_name"`
	StartDate         time.Time       `gorm:"type:date;not null" json:"start_date"`
	EndDate           time.Time       `gorm:"type:date;not null" json:"end_date"`
	DailyCalories     float64         `gorm:"type:decimal(7,2)" json:"daily_calories"`
	ProteinRatio      float64         `gorm:"type:decimal(3,2)" json:"protein_ratio"`
	CarbRatio         float64         `gorm:"type:decimal(3,2)" json:"carb_ratio"`
	FatRatio          float64         `gorm:"type:decimal(3,2)" json:"fat_ratio"`
	DietaryRestrictions JSONSlice      `gorm:"type:json" json:"dietary_restrictions"`
	Preferences        JSONSlice      `gorm:"type:json" json:"preferences"`
	PlanData          JSONMap         `gorm:"type:json;not null" json:"plan_data"`
	AIAPIID           int64           `gorm:"not null;index" json:"ai_api_id"`
	Status            string          `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`

	// 关联关系
	User  User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AIAPI AIAPI  `gorm:"foreignKey:AIAPIID" json:"ai_api,omitempty"`
}

func (NutritionPlan) TableName() string {
	return "nutrition_plans"
}

type NutritionRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index;index:user_date" json:"user_id"`
	MealDate  time.Time `gorm:"type:date;not null;index:user_date" json:"meal_date"`
	MealTime  string    `gorm:"type:enum('breakfast','lunch','dinner','snack')" json:"meal_time"`
	Foods     JSONMap   `gorm:"type:json;not null" json:"foods"`
	Calories  float64   `gorm:"type:decimal(7,2)" json:"calories"`
	Protein   float64   `gorm:"type:decimal(6,2)" json:"protein"`
	Carbs     float64   `gorm:"type:decimal(6,2)" json:"carbs"`
	Fat       float64   `gorm:"type:decimal(6,2)" json:"fat"`
	Fiber     float64   `gorm:"type:decimal(6,2)" json:"fiber"`
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (NutritionRecord) TableName() string {
	return "nutrition_records"
}

type MealTime string

const (
	MealTimeBreakfast MealTime = "breakfast"
	MealTimeLunch     MealTime = "lunch"
	MealTimeDinner    MealTime = "dinner"
	MealTimeSnack     MealTime = "snack"
)
type NutritionPlanMeal struct {
	Time         string              `json:"time"`
	Foods        []NutritionFoodItem `json:"foods"`
	TotalCalories float64            `json:"total_calories"`
}

type NutritionFoodItem struct {
	Name     string  `json:"name"`
	Amount   string  `json:"amount"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber,omitempty"`
}
