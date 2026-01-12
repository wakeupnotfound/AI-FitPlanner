package request

// TrainingStatsParams represents query parameters for training statistics
type TrainingStatsParams struct {
	Period    string `form:"period" binding:"omitempty,oneof=week month quarter year all"`
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

// TrendsParams represents query parameters for trends
type TrendsParams struct {
	Period    string `form:"period" binding:"omitempty,oneof=week month"`
	Count     int    `form:"count" binding:"omitempty,min=1,max=52"`
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
}
