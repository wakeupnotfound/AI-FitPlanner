package model

import (
	"time"
)

type FeedbackRecord struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"not null;index" json:"user_id"`
	PlanType     string    `gorm:"type:enum('training','nutrition');index" json:"plan_type"`
	PlanID       *int64    `json:"plan_id"`
	FeedbackType string    `gorm:"size:50" json:"feedback_type"`
	FeedbackData JSONMap   `gorm:"type:json" json:"feedback_data"`
	Satisfaction *int      `json:"satisfaction"`
	AIResponse   *string   `gorm:"type:text" json:"ai_response"`
	CreatedAt    time.Time `json:"created_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (FeedbackRecord) TableName() string {
	return "feedback_records"
}

type FeedbackType string

const (
	FeedbackTypeGeneral   FeedbackType = "general"
	FeedbackTypeDifficulty FeedbackType = "difficulty"
	FeedbackEffectiveness FeedbackType = "effectiveness"
	FeedbackSafety        FeedbackType = "safety"
)

type FeedbackData struct {
	PlanQuality   int    `json:"plan_quality"`
	ExerciseVariety int  `json:"exercise_variety"`
	DifficultyFit   bool   `json:"difficulty_fit"`
	TimeReasonable  bool   `json:"time_reasonable"`
	Comments      string `json:"comments"`
	Suggestions   string `json:"suggestions"`
}
