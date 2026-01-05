package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-fitness-planner/backend/internal/api/request"
	"github.com/ai-fitness-planner/backend/internal/api/response"
	apperrors "github.com/ai-fitness-planner/backend/internal/errors"
	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/pkg/crypto"
	"github.com/ai-fitness-planner/backend/internal/repository"
)

// AIAPIService defines the interface for AI API management operations
type AIAPIService interface {
	// AddAPI adds a new AI API configuration with encrypted API key
	AddAPI(ctx context.Context, userID int64, req *request.AddAIAPIRequest) (*response.AIAPIInfo, error)
	// ListAPIs returns all AI API configurations for a user without exposing encrypted keys
	ListAPIs(ctx context.Context, userID int64) (*response.AIAPIListResponse, error)
	// GetAPI returns a single AI API configuration by ID
	GetAPI(ctx context.Context, userID int64, apiID int64) (*response.AIAPIInfo, error)
	// UpdateAPI updates an existing AI API configuration
	UpdateAPI(ctx context.Context, userID int64, apiID int64, req *request.UpdateAIAPIRequest) (*response.AIAPIInfo, error)
	// TestAPI tests the connection to an AI API
	TestAPI(ctx context.Context, userID int64, apiID int64) (*response.TestAPIResponse, error)
	// SetDefault sets an AI API as the default for a user
	SetDefault(ctx context.Context, userID int64, apiID int64) error
	// DeleteAPI deletes an AI API configuration
	DeleteAPI(ctx context.Context, userID int64, apiID int64) error
}

// aiAPIService implements AIAPIService interface
type aiAPIService struct {
	aiAPIRepo repository.AIAPIRepository
	encryptor crypto.Encryptor
}

// NewAIAPIService creates a new instance of AIAPIService
func NewAIAPIService(
	aiAPIRepo repository.AIAPIRepository,
	encryptor crypto.Encryptor,
) AIAPIService {
	return &aiAPIService{
		aiAPIRepo: aiAPIRepo,
		encryptor: encryptor,
	}
}

// AddAPI adds a new AI API configuration with encrypted API key
// Requirements: 3.1 - Encrypt API key using AES-256 before storage
func (s *aiAPIService) AddAPI(ctx context.Context, userID int64, req *request.AddAIAPIRequest) (*response.AIAPIInfo, error) {
	// Encrypt the API key before storage
	encryptedKey, err := s.encryptor.Encrypt(req.APIKey)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternalServer, "failed to encrypt API key")
	}

	// Create the model
	api := &model.AIAPI{
		UserID:          userID,
		Provider:        req.Provider,
		Name:            req.Name,
		APIEndpoint:     req.APIEndpoint,
		APIKeyEncrypted: encryptedKey,
		Status:          1, // Active by default
	}

	// Set optional fields
	if req.Model != "" {
		api.Model = &req.Model
	}
	if req.MaxTokens != nil {
		api.MaxTokens = req.MaxTokens
	}
	if req.Temperature != nil {
		temp := float32(*req.Temperature)
		api.Temperature = &temp
	}

	// Handle is_default flag
	if req.IsDefault != nil && *req.IsDefault {
		// If setting as default, first unset other defaults via SetDefault after creation
		api.IsDefault = true
	}

	// Create the API configuration
	if err := s.aiAPIRepo.Create(ctx, api); err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "failed to create AI API configuration")
	}

	// If this API should be default, ensure only this one is default
	if api.IsDefault {
		if err := s.aiAPIRepo.SetDefault(ctx, userID, api.ID); err != nil {
			return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "failed to set API as default")
		}
	}

	return s.modelToAPIInfo(api), nil
}

// ListAPIs returns all AI API configurations for a user without exposing encrypted keys
// Requirements: 3.2 - Return configurations without exposing the encrypted keys
func (s *aiAPIService) ListAPIs(ctx context.Context, userID int64) (*response.AIAPIListResponse, error) {
	apis, err := s.aiAPIRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "failed to list AI APIs")
	}

	apiInfos := make([]response.AIAPIInfo, 0, len(apis))
	for _, api := range apis {
		apiInfos = append(apiInfos, *s.modelToAPIInfo(api))
	}

	return &response.AIAPIListResponse{
		APIs: apiInfos,
	}, nil
}

// GetAPI returns a single AI API configuration by ID
func (s *aiAPIService) GetAPI(ctx context.Context, userID int64, apiID int64) (*response.AIAPIInfo, error) {
	api, err := s.aiAPIRepo.GetByID(ctx, apiID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "failed to get AI API")
	}
	if api == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "AI API not found")
	}

	// Verify ownership
	if api.UserID != userID {
		return nil, apperrors.New(apperrors.ErrForbidden, "unauthorized access to AI API")
	}

	return s.modelToAPIInfo(api), nil
}

// UpdateAPI updates an existing AI API configuration
func (s *aiAPIService) UpdateAPI(ctx context.Context, userID int64, apiID int64, req *request.UpdateAIAPIRequest) (*response.AIAPIInfo, error) {
	// Get existing API
	api, err := s.aiAPIRepo.GetByID(ctx, apiID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "failed to get AI API")
	}
	if api == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "AI API not found")
	}

	// Verify ownership
	if api.UserID != userID {
		return nil, apperrors.New(apperrors.ErrForbidden, "unauthorized access to AI API")
	}

	// Update fields if provided
	if req.Name != "" {
		api.Name = req.Name
	}
	if req.APIEndpoint != "" {
		api.APIEndpoint = req.APIEndpoint
	}
	if req.APIKey != "" {
		// Encrypt the new API key
		encryptedKey, err := s.encryptor.Encrypt(req.APIKey)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.ErrInternalServer, "failed to encrypt API key")
		}
		api.APIKeyEncrypted = encryptedKey
	}
	if req.Model != "" {
		api.Model = &req.Model
	}
	if req.MaxTokens != nil {
		api.MaxTokens = req.MaxTokens
	}
	if req.Temperature != nil {
		temp := float32(*req.Temperature)
		api.Temperature = &temp
	}
	if req.Status != nil {
		if *req.Status {
			api.Status = 1
		} else {
			api.Status = 0
		}
	}

	// Update the API
	if err := s.aiAPIRepo.Update(ctx, api); err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "failed to update AI API")
	}

	// Handle is_default flag separately to ensure single default invariant
	if req.IsDefault != nil && *req.IsDefault {
		if err := s.aiAPIRepo.SetDefault(ctx, userID, apiID); err != nil {
			return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "failed to set API as default")
		}
		api.IsDefault = true
	}

	return s.modelToAPIInfo(api), nil
}

// TestAPI tests the connection to an AI API
// Requirements: 3.3 - Attempt a connection and return the test result
func (s *aiAPIService) TestAPI(ctx context.Context, userID int64, apiID int64) (*response.TestAPIResponse, error) {
	// Get the API configuration
	api, err := s.aiAPIRepo.GetByID(ctx, apiID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "failed to get AI API")
	}
	if api == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "AI API not found")
	}

	// Verify ownership
	if api.UserID != userID {
		return nil, apperrors.New(apperrors.ErrForbidden, "unauthorized access to AI API")
	}

	// Decrypt API key for testing
	// Requirements: 3.6 - Decrypt the API key for the request
	apiKey, err := s.encryptor.Decrypt(api.APIKeyEncrypted)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternalServer, "failed to decrypt API key")
	}

	// Get the appropriate AI client
	client, err := GetAIClient(api.Provider)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternalServer, fmt.Sprintf("unsupported provider: %s", api.Provider))
	}

	// Create client config
	config := NewAIClientFromModel(api, apiKey)

	// Test the connection and measure response time
	startTime := time.Now()
	testErr := client.TestConnection(ctx, config)
	responseTime := int(time.Since(startTime).Milliseconds())

	// Build response
	testResult := response.APITestResult{
		ResponseTime: responseTime,
	}

	modelName := ""
	if api.Model != nil {
		modelName = *api.Model
	}
	maxTokens := 0
	if api.MaxTokens != nil {
		maxTokens = *api.MaxTokens
	}

	testResult.ModelInfo = response.ModelInfo{
		Name:      modelName,
		MaxTokens: maxTokens,
	}

	if testErr != nil {
		testResult.Status = "failed"
		testResult.Message = testErr.Error()
	} else {
		testResult.Status = "success"
		testResult.Message = "Connection successful"
	}

	return &response.TestAPIResponse{
		TestResult: testResult,
	}, nil
}

// SetDefault sets an AI API as the default for a user
// Requirements: 3.4 - Update the default flag and unset other defaults
func (s *aiAPIService) SetDefault(ctx context.Context, userID int64, apiID int64) error {
	// Verify the API exists and belongs to the user
	api, err := s.aiAPIRepo.GetByID(ctx, apiID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "failed to get AI API")
	}
	if api == nil {
		return apperrors.New(apperrors.ErrNotFound, "AI API not found")
	}

	// Verify ownership
	if api.UserID != userID {
		return apperrors.New(apperrors.ErrForbidden, "unauthorized access to AI API")
	}

	// Set as default (repository handles unsetting other defaults in a transaction)
	if err := s.aiAPIRepo.SetDefault(ctx, userID, apiID); err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "failed to set API as default")
	}

	return nil
}

// DeleteAPI deletes an AI API configuration
// Requirements: 3.5 - Remove the configuration and associated data
func (s *aiAPIService) DeleteAPI(ctx context.Context, userID int64, apiID int64) error {
	// Verify the API exists and belongs to the user
	api, err := s.aiAPIRepo.GetByID(ctx, apiID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "failed to get AI API")
	}
	if api == nil {
		return apperrors.New(apperrors.ErrNotFound, "AI API not found")
	}

	// Verify ownership
	if api.UserID != userID {
		return apperrors.New(apperrors.ErrForbidden, "unauthorized access to AI API")
	}

	// Delete the API
	if err := s.aiAPIRepo.Delete(ctx, apiID); err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "failed to delete AI API")
	}

	return nil
}

// modelToAPIInfo converts a model.AIAPI to response.AIAPIInfo
// This function ensures encrypted keys are never exposed in responses
func (s *aiAPIService) modelToAPIInfo(api *model.AIAPI) *response.AIAPIInfo {
	info := &response.AIAPIInfo{
		ID:          api.ID,
		Provider:    api.Provider,
		Name:        api.Name,
		APIEndpoint: api.APIEndpoint,
		IsDefault:   api.IsDefault,
		Status:      api.Status == 1,
		CreatedAt:   api.CreatedAt.Format(time.RFC3339),
	}

	if api.Model != nil {
		info.Model = *api.Model
	}
	if api.MaxTokens != nil {
		info.MaxTokens = *api.MaxTokens
	}
	if api.Temperature != nil {
		info.Temperature = float64(*api.Temperature)
	}

	return info
}
