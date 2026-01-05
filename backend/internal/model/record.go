package model

import (
	"time"
)

type TrainingRecord struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64     `gorm:"not null;index;index:user_date" json:"user_id" validate:"required"`
	PlanID          *int64    `gorm:"index;index:user_date" json:"plan_id"`
	WorkoutDate     time.Time `gorm:"type:date;not null;index:user_date" json:"workout_date" validate:"required"`
	WorkoutType     string    `gorm:"size:100;not null" json:"workout_type" validate:"required,max=100"`
	DurationMinutes *int      `json:"duration_minutes" validate:"omitempty,min=0"`
	Exercises       JSONMap   `gorm:"type:json" json:"exercises"`
	PerformanceData JSONMap   `gorm:"type:json" json:"performance_data"`
	Notes           *string   `gorm:"type:text" json:"notes"`
	Rating          *int      `json:"rating" validate:"omitempty,min=1,max=5"`
	InjuryReport    *string   `gorm:"type:text" json:"injury_report"`
	CreatedAt       time.Time `json:"created_at"`

	// 关联关系
	User         User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TrainingPlan *TrainingPlan `gorm:"foreignKey:PlanID" json:"training_plan,omitempty"`
}

func (TrainingRecord) TableName() string {
	return "training_records"
}

type ExerciseRecord struct {
	ExerciseName string    `json:"exercise_name"`
	Sets         int       `json:"sets"`
	RepsPerSet   []int     `json:"reps_per_set"`
	WeightUsed   []float64 `json:"weight_used"`
	Notes        string    `json:"notes"`
	Difficulty   string    `json:"difficulty"`
}

type PerformanceData struct {
	TotalVolume       float64 `json:"total_volume"`
	EstimatedCalories int     `json:"estimated_calories"`
	AvgHeartRate      *int    `json:"avg_heart_rate"`
	MaxHeartRate      *int    `json:"max_heart_rate"`
}
