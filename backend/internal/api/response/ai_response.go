package response

// AI API配置响应

type AIAPIListResponse struct {
	APIs []AIAPIInfo `json:"apis"`
}

type AIAPIInfo struct {
	ID          int64   `json:"id"`
	Provider    string  `json:"provider"`
	Name        string  `json:"name"`
	APIEndpoint string  `json:"api_endpoint"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	IsDefault   bool    `json:"is_default"`
	Status      bool    `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

type AIAPIDetailResponse struct {
	API AIAPIInfo `json:"api"`
}

type TestAPIResponse struct {
	TestResult APITestResult `json:"test_result"`
}

type APITestResult struct {
	Status       string    `json:"status"`
	ResponseTime int       `json:"response_time"`
	ModelInfo    ModelInfo `json:"model_info"`
	Message      string    `json:"message,omitempty"`
}

type ModelInfo struct {
	Name      string `json:"name"`
	MaxTokens int    `json:"max_tokens"`
	Version   string `json:"version,omitempty"`
}
