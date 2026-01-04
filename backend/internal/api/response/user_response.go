package response

import "time"

type UserProfileResponse struct {
	User     UserInfo      `json:"user"`
	BodyData []BodyDataInfo `json:"body_data,omitempty"`
	Goals    []GoalInfo    `json:"goals,omitempty"`
}

type BodyDataInfo struct {
	ID                int64     `json:"id"`
	Age               int       `json:"age"`
	Gender            string    `json:"gender"`
	Height            float64   `json:"height"`
	Weight            float64   `json:"weight"`
	BodyFatPercentage float64   `json:"body_fat_percentage,omitempty"`
	MusclePercentage  float64   `json:"muscle_percentage,omitempty"`
	MeasurementDate   string    `json:"measurement_date"`
	CreatedAt         string    `json:"created_at"`
}

type GoalInfo struct {
	ID              int64    `json:"id"`
	GoalType        string   `json:"goal_type"`
	GoalDescription string   `json:"goal_description"`
	TargetWeight    float64  `json:"target_weight,omitempty"`
	Deadline        string   `json:"deadline,omitempty"`
	Priority        int      `json:"priority"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"created_at"`
}

type BodyDataListResponse struct {
	BodyData []BodyDataInfo `json:"body_data"`
	Pagination PaginationInfo `json:"pagination"`
}

type GoalListResponse struct {
	Goals      []GoalInfo     `json:"goals"`
	Pagination PaginationInfo `json:"pagination"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
