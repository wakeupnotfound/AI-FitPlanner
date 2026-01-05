package request

// AI API配置请求
type AddAIAPIRequest struct {
	Provider    string   `json:"provider" binding:"required,oneof=openai wenxin tongyi"`
	Name        string   `json:"name" binding:"required,min=1,max=100"`
	APIEndpoint string   `json:"api_endpoint" binding:"required,url,max=500"`
	APIKey      string   `json:"api_key" binding:"required,min=1,max=500"`
	Model       string   `json:"model" binding:"required,min=1,max=100"`
	MaxTokens   *int     `json:"max_tokens" binding:"omitempty,min=1,max=100000"`
	Temperature *float64 `json:"temperature" binding:"omitempty,min=0,max=2"`
	IsDefault   *bool    `json:"is_default"`
}

type UpdateAIAPIRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=1,max=100"`
	APIEndpoint string   `json:"api_endpoint" binding:"omitempty,url,max=500"`
	APIKey      string   `json:"api_key" binding:"omitempty,min=1,max=500"`
	Model       string   `json:"model" binding:"omitempty,min=1,max=100"`
	MaxTokens   *int     `json:"max_tokens" binding:"omitempty,min=1,max=100000"`
	Temperature *float64 `json:"temperature" binding:"omitempty,min=0,max=2"`
	Status      *bool    `json:"status"`
	IsDefault   *bool    `json:"is_default"`
}

type AIAPIIDParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
