package model

import (
	"time"
)

type PromptTemplate struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Category    string    `gorm:"size:50;not null;index" json:"category"`
	Subcategory *string   `gorm:"size:50" json:"subcategory"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Template    string    `gorm:"type:text;not null" json:"template"`
	Variables   JSONSlice `gorm:"type:json" json:"variables"`
	IsDefault   bool      `gorm:"default:false;index" json:"is_default"`
	Description *string   `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PromptTemplate) TableName() string {
	return "prompt_templates"
}

type PromptCategory string

const (
	PromptCategoryTraining   PromptCategory = "training"
	PromptCategoryNutrition  PromptCategory = "nutrition"
	PromptCategoryAssessment PromptCategory = "assessment"
	PromptCategorySafety     PromptCategory = "safety"
)

type TemplateVariable struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Description  string      `json:"description"`
	Required     bool        `json:"required"`
	DefaultValue interface{} `json:"default_value,omitempty"`
}
