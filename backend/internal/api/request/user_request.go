package request

import "time"

// 添加身体数据请求
type AddBodyDataRequest struct {
	Age               int       `json:"age" binding:"required,min=1,max=150"`
	Gender            string    `json:"gender" binding:"required,oneof=male female other"`
	Height            float64   `json:"height" binding:"required,min=50,max=300"`
	Weight            float64   `json:"weight" binding:"required,min=20,max=500"`
	BodyFatPercentage *float64  `json:"body_fat_percentage" binding:"omitempty,min=0,max=80"`
	MusclePercentage  *float64  `json:"muscle_percentage" binding:"omitempty,min=0,max=100"`
	MeasurementDate   time.Time `json:"measurement_date" binding:"required"`
}

// 添加健身目标请求
type AddGoalRequest struct {
	GoalType        string   `json:"goal_type" binding:"required"`
	GoalDescription string   `json:"goal_description" binding:"required"`
	TargetWeight    *float64 `json:"target_weight" binding:"omitempty,min=20,max=500"`
	Deadline        *string  `json:"deadline" binding:"omitempty,datetime=2006-01-02"`
	Priority        *int     `json:"priority" binding:"omitempty,min=1,max=10"`
}

// 更新目标请求
type UpdateGoalRequest struct {
	GoalType        string   `json:"goal_type" binding:"omitempty"`
	GoalDescription string   `json:"goal_description" binding:"omitempty"`
	TargetWeight    *float64 `json:"target_weight" binding:"omitempty,min=20,max=500"`
	Deadline        *string  `json:"deadline" binding:"omitempty,datetime=2006-01-02"`
	Priority        *int     `json:"priority" binding:"omitempty,min=1,max=10"`
	Status          *string  `json:"status" binding:"omitempty,oneof=active completed cancelled"`
}

// 目标ID参数
type GoalIDParam struct {
	GoalID int64 `uri:"id" binding:"required,min=1"`
}

// 分页查询参数
type PaginationParams struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1,max=100"`
}

// 日期范围查询参数
type DateRangeParams struct {
	StartDate string `form:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"required,datetime=2006-01-02"`
}
