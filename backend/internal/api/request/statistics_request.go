package request

// TrainingStatsParams represents query parameters for training statistics
type TrainingStatsParams struct {
	Period string `form:"period" binding:"required,oneof=week month quarter year all"`
}

// TrendsParams represents query parameters for trends
type TrendsParams struct {
	Period string `form:"period" binding:"required,oneof=week month"`
	Count  int    `form:"count" binding:"omitempty,min=1,max=52"`
}
