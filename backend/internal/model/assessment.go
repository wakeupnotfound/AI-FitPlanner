package model

import (
	"time"
)

// FitnessAssessment represents a user's fitness assessment
type FitnessAssessment struct {
	ID                    int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                int64     `gorm:"not null;index" json:"user_id" validate:"required"`
	ExperienceLevel       string    `gorm:"type:enum('beginner','intermediate','advanced');not null" json:"experience_level" validate:"required,oneof=beginner intermediate advanced"`
	WeeklyAvailableDays   int       `gorm:"not null" json:"weekly_available_days" validate:"required,min=1,max=7"`
	DailyAvailableMinutes int       `gorm:"not null" json:"daily_available_minutes" validate:"required,min=10,max=480"`
	ActivityType          *string   `gorm:"size:50" json:"activity_type" validate:"omitempty,max=50"`
	InjuryHistory         *string   `gorm:"type:text" json:"injury_history"`
	HealthConditions      *string   `gorm:"type:text" json:"health_conditions"`
	PreferredDays         JSONSlice `gorm:"type:json" json:"preferred_days"`
	EquipmentAvailable    JSONSlice `gorm:"type:json" json:"equipment_available"`
	AssessmentDate        time.Time `gorm:"type:date;not null" json:"assessment_date" validate:"required"`
	CreatedAt             time.Time `json:"created_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (FitnessAssessment) TableName() string {
	return "fitness_assessments"
}

// ExperienceLevel constants
type ExperienceLevel string

const (
	ExperienceLevelBeginner     ExperienceLevel = "beginner"
	ExperienceLevelIntermediate ExperienceLevel = "intermediate"
	ExperienceLevelAdvanced     ExperienceLevel = "advanced"
)
