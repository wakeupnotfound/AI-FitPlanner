package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ai-fitness-planner/backend/internal/model"
)

// AIClient defines the interface for AI service providers
type AIClient interface {
	// Call sends a prompt to the AI service and returns the response
	Call(ctx context.Context, prompt string, config *AIClientConfig) (string, error)
	// TestConnection tests the connectivity to the AI service
	TestConnection(ctx context.Context, config *AIClientConfig) error
}

// AIClientConfig holds the configuration for an AI client
type AIClientConfig struct {
	APIEndpoint string
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float32
}

// NewAIClientFromModel creates an AIClientConfig from a model.AIAPI
func NewAIClientFromModel(api *model.AIAPI, decryptedKey string) *AIClientConfig {
	config := &AIClientConfig{
		APIEndpoint: api.APIEndpoint,
		APIKey:      decryptedKey,
	}

	if api.Model != nil {
		config.Model = *api.Model
	}
	if api.MaxTokens != nil {
		config.MaxTokens = *api.MaxTokens
	}
	if api.Temperature != nil {
		config.Temperature = *api.Temperature
	}

	return config
}

// GetAIClient returns the appropriate AI client based on the provider
func GetAIClient(provider string) (AIClient, error) {
	switch provider {
	case "openai":
		return &OpenAIClient{}, nil
	case "wenxin":
		return &WenxinClient{}, nil
	case "tongyi":
		return &TongyiClient{}, nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", provider)
	}
}

// OpenAIClient implements AIClient for OpenAI API
type OpenAIClient struct{}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response structure from OpenAI API
type OpenAIResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choice  `json:"choices"`
	Usage   Usage     `json:"usage"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice represents a response choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// APIError represents an API error response
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// Call sends a request to OpenAI API
func (c *OpenAIClient) Call(ctx context.Context, prompt string, config *AIClientConfig) (string, error) {
	// Set defaults
	model := config.Model
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2000
	}
	temperature := config.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	reqBody := OpenAIRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := config.APIEndpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1"
	}
	url := fmt.Sprintf("%s/chat/completions", endpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// TestConnection tests the connection to OpenAI API
func (c *OpenAIClient) TestConnection(ctx context.Context, config *AIClientConfig) error {
	_, err := c.Call(ctx, "Hello, this is a test message.", config)
	return err
}

// WenxinClient implements AIClient for Baidu Wenxin API
type WenxinClient struct{}

// WenxinRequest represents the request structure for Wenxin API
type WenxinRequest struct {
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	TopP        float32   `json:"top_p,omitempty"`
}

// WenxinResponse represents the response structure from Wenxin API
type WenxinResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	Result           string `json:"result"`
	IsTruncated      bool   `json:"is_truncated"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            Usage  `json:"usage"`
	ErrorCode        int    `json:"error_code,omitempty"`
	ErrorMsg         string `json:"error_msg,omitempty"`
}

// Call sends a request to Wenxin API
func (c *WenxinClient) Call(ctx context.Context, prompt string, config *AIClientConfig) (string, error) {
	temperature := config.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	reqBody := WenxinRequest{
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Temperature: temperature,
		TopP:        0.8,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Wenxin API requires access_token in URL
	url := fmt.Sprintf("%s?access_token=%s", config.APIEndpoint, config.APIKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var wenxinResp WenxinResponse
	if err := json.Unmarshal(body, &wenxinResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if wenxinResp.ErrorCode != 0 {
		return "", fmt.Errorf("Wenxin API error: %s", wenxinResp.ErrorMsg)
	}

	return wenxinResp.Result, nil
}

// TestConnection tests the connection to Wenxin API
func (c *WenxinClient) TestConnection(ctx context.Context, config *AIClientConfig) error {
	_, err := c.Call(ctx, "你好，这是一条测试消息。", config)
	return err
}

// TongyiClient implements AIClient for Alibaba Tongyi API using OpenAI-compatible format
type TongyiClient struct{}

// TongyiRequest represents the OpenAI-compatible request structure for Tongyi API
type TongyiRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// TongyiResponse represents the OpenAI-compatible response structure from Tongyi API
type TongyiResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// Call sends a request to Tongyi API using OpenAI-compatible format
func (c *TongyiClient) Call(ctx context.Context, prompt string, config *AIClientConfig) (string, error) {
	model := config.Model
	if model == "" {
		model = "qwen-turbo"
	}
	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2000
	}
	temperature := config.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	// Use OpenAI-compatible format
	var reqBody TongyiRequest
	reqBody.Model = model
	reqBody.Messages = []Message{
		{Role: "user", Content: prompt},
	}
	reqBody.Temperature = temperature
	reqBody.MaxTokens = maxTokens

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Ensure the endpoint uses the OpenAI-compatible format
	endpoint := strings.TrimSpace(config.APIEndpoint)
	defaultEndpoint := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	if endpoint == "" {
		endpoint = defaultEndpoint
	} else if strings.Contains(endpoint, "compatible-mode") {
		if !strings.Contains(endpoint, "/chat/completions") {
			endpoint = strings.TrimRight(endpoint, "/")
			if strings.HasSuffix(endpoint, "compatible-mode") || strings.HasSuffix(endpoint, "compatible-mode/v1") {
				endpoint = endpoint + "/chat/completions"
			} else {
				endpoint = endpoint + "/v1/chat/completions"
			}
		}
	} else {
		// If user provided a custom endpoint without compatible-mode, use the standard one
		endpoint = defaultEndpoint
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Add debug logging for response
	fmt.Printf("Tongyi API Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Tongyi API Response Body: %s\n", string(body))
	fmt.Printf("Tongyi API Response Headers: %v\n", resp.Header)

	statusCode := resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		if len(body) == 0 {
			return "", fmt.Errorf("Tongyi API error: status %d, empty body", statusCode)
		}
		var tongyiResp TongyiResponse
		if err := json.Unmarshal(body, &tongyiResp); err == nil && tongyiResp.Error != nil {
			return "", fmt.Errorf("Tongyi API error: %s (type: %s, code: %s)",
				tongyiResp.Error.Message, tongyiResp.Error.Type, tongyiResp.Error.Code)
		}
		return "", fmt.Errorf("Tongyi API error: status %d, body: %s", statusCode, string(body))
	}

	// Check if response is empty
	if len(body) == 0 {
		return "", fmt.Errorf("empty response from Tongyi API (status %d)", statusCode)
	}

	var tongyiResp TongyiResponse
	if err := json.Unmarshal(body, &tongyiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w, response body: %s", err, string(body))
	}

	// Check for API errors
	if tongyiResp.Error != nil {
		return "", fmt.Errorf("Tongyi API error: %s (type: %s, code: %s)", 
			tongyiResp.Error.Message, tongyiResp.Error.Type, tongyiResp.Error.Code)
	}

	// Check if we have choices
	if len(tongyiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in Tongyi API response")
	}

	return tongyiResp.Choices[0].Message.Content, nil
}

// TestConnection tests the connection to Tongyi API
func (c *TongyiClient) TestConnection(ctx context.Context, config *AIClientConfig) error {
	_, err := c.Call(ctx, "你好，这是一条测试消息。", config)
	return err
}
