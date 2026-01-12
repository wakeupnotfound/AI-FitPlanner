package model

import (
	"time"
)

// UserBodyData represents a user's body measurements
type UserBodyData struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int64     `gorm:"not null;index:user_date" json:"user_id" validate:"required"`
	Age               int       `gorm:"not null" json:"age" validate:"required,min=1,max=150"`
	Gender            string    `gorm:"type:enum('male','female','other');not null" json:"gender" validate:"required,oneof=male female other"`
	Height            float64   `gorm:"type:decimal(5,2);not null" json:"height" validate:"required,min=50,max=300"`
	Weight            float64   `gorm:"type:decimal(5,2);not null" json:"weight" validate:"required,min=20,max=500"`
	BodyFatPercentage *float64  `gorm:"type:decimal(4,2)" json:"body_fat_percentage" validate:"omitempty,min=0,max=100"`
	MusclePercentage  *float64  `gorm:"type:decimal(4,2)" json:"muscle_percentage" validate:"omitempty,min=0,max=100"`
	MeasurementDate   time.Time `gorm:"type:date;not null;index:user_date" json:"measurement_date" validate:"required"`
	CreatedAt         time.Time `json:"created_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserBodyData) TableName() string {
	return "user_body_data"
}

// FitnessGoal represents a user's fitness goal
type FitnessGoal struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64      `gorm:"not null;index:user_status" json:"user_id" validate:"required"`
	GoalType        string     `gorm:"size:100;not null" json:"goal_type" validate:"required,max=100"`
	GoalDescription *string    `gorm:"type:text" json:"goal_description"`
	InitialWeight   *float64   `gorm:"type:decimal(5,2)" json:"initial_weight" validate:"omitempty,min=20,max=500"`
	InitialBodyFat  *float64   `gorm:"type:decimal(4,2)" json:"initial_body_fat" validate:"omitempty,min=0,max=100"`
	InitialMuscle   *float64   `gorm:"type:decimal(4,2)" json:"initial_muscle_mass" validate:"omitempty,min=0,max=100"`
	TargetWeight    *float64   `gorm:"type:decimal(5,2)" json:"target_weight" validate:"omitempty,min=20,max=500"`
	Deadline        *time.Time `gorm:"type:date" json:"deadline"`
	Priority        int        `gorm:"default:1" json:"priority" validate:"min=1,max=10"`
	Status          string     `gorm:"size:20;default:'active';index:user_status" json:"status" validate:"oneof=active completed cancelled"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (FitnessGoal) TableName() string {
	return "fitness_goals"
}

// Gender constants
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// GoalStatus constants
type GoalStatus string

const (
	GoalStatusActive    GoalStatus = "active"
	GoalStatusCompleted GoalStatus = "completed"
	GoalStatusCancelled GoalStatus = "cancelled"
)
