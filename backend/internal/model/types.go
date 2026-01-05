package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONMap is a custom type for JSON object fields
type JSONMap map[string]interface{}

// Scan implements the sql.Scanner interface for JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for JSONMap
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// JSONSlice is a custom type for JSON array fields
type JSONSlice []interface{}

// Scan implements the sql.Scanner interface for JSONSlice
func (j *JSONSlice) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONSlice, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for JSONSlice
func (j JSONSlice) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// User model represents a registered user in the system
type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`
	Email        string    `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email,max=100"`
	Phone        *string   `gorm:"size:20" json:"phone" validate:"omitempty,max=20"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Avatar       *string   `gorm:"size:255" json:"avatar" validate:"omitempty,url,max=255"`
	Status       int8      `gorm:"default:1" json:"status" validate:"oneof=0 1"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// AIAPI model represents user's AI service configuration
type AIAPI struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64     `gorm:"not null;index" json:"user_id" validate:"required"`
	Provider        string    `gorm:"size:50;not null" json:"provider" validate:"required,oneof=openai wenxin tongyi"`
	Name            string    `gorm:"size:100;not null" json:"name" validate:"required,min=1,max=100"`
	APIEndpoint     string    `gorm:"size:500;not null" json:"api_endpoint" validate:"required,url,max=500"`
	APIKeyEncrypted string    `gorm:"type:text;not null" json:"-"`
	Model           *string   `gorm:"size:100" json:"model" validate:"omitempty,max=100"`
	MaxTokens       *int      `json:"max_tokens" validate:"omitempty,min=1,max=32000"`
	Temperature     *float32  `gorm:"type:decimal(3,2)" json:"temperature" validate:"omitempty,min=0,max=2"`
	IsDefault       bool      `gorm:"default:false" json:"is_default"`
	Status          int8      `gorm:"default:1" json:"status" validate:"oneof=0 1"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (AIAPI) TableName() string {
	return "ai_apis"
}

// TrainingPlan model represents a user's training plan
type TrainingPlan struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64     `gorm:"not null;index" json:"user_id" validate:"required"`
	PlanName        string    `gorm:"size:200;not null" json:"plan_name" validate:"required,min=1,max=200"`
	StartDate       time.Time `gorm:"type:date;not null" json:"start_date" validate:"required"`
	EndDate         time.Time `gorm:"type:date;not null" json:"end_date" validate:"required,gtfield=StartDate"`
	TotalWeeks      int       `gorm:"not null" json:"total_weeks" validate:"required,min=1,max=52"`
	DifficultyLevel string    `gorm:"type:enum('easy','medium','hard','extreme')" json:"difficulty_level" validate:"oneof=easy medium hard extreme"`
	TrainingPurpose *string   `gorm:"size:100" json:"training_purpose" validate:"omitempty,max=100"`
	AIAPIID         int64     `gorm:"not null;index" json:"ai_api_id" validate:"required"`
	PlanData        JSONMap   `gorm:"type:json;not null" json:"plan_data"`
	Status          string    `gorm:"size:20;default:'active'" json:"status" validate:"oneof=active inactive completed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (TrainingPlan) TableName() string {
	return "training_plans"
}

// PlanData represents the structure of training plan data stored in JSON
type PlanData struct {
	Weeks []WeekPlan `json:"weeks"`
}

// WeekPlan represents a week in the training plan
type WeekPlan struct {
	Week int       `json:"week"`
	Days []DayPlan `json:"days"`
}

// DayPlan represents a single day's training schedule
type DayPlan struct {
	Day               int        `json:"day"`
	Date              string     `json:"date"`
	Type              string     `json:"type"` // strength/cardio/rest
	FocusArea         string     `json:"focus_area"`
	Exercises         []Exercise `json:"exercises"`
	Duration          int        `json:"duration"` // minutes
	EstimatedCalories int        `json:"estimated_calories"`
}

// Exercise represents a single exercise in a training plan
type Exercise struct {
	Name        string `json:"name"`
	Sets        int    `json:"sets"`
	Reps        string `json:"reps"` // "8-10" or "12"
	Weight      string `json:"weight"`
	Rest        string `json:"rest"` // "90s"
	Difficulty  string `json:"difficulty"`
	SafetyNotes string `json:"safety_notes"`
}
