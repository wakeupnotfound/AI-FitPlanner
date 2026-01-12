package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/pkg/crypto"
	"github.com/ai-fitness-planner/backend/internal/repository"
)

// AIService defines the interface for AI integration operations
type AIService interface {
	// GenerateTrainingPlan generates a training plan using AI
	GenerateTrainingPlan(ctx context.Context, params *TrainingPlanParams) (*model.TrainingPlan, error)
	// GenerateNutritionPlan generates a nutrition plan using AI
	GenerateNutritionPlan(ctx context.Context, params *NutritionPlanParams) (*model.NutritionPlan, error)
	// TestConnection tests the connection to an AI API
	TestConnection(ctx context.Context, apiID int64, userID int64) error
}

// aiService implements AIService interface
type aiService struct {
	aiAPIRepo  repository.AIAPIRepository
	encryptor  crypto.Encryptor
	maxRetries int
	retryDelay time.Duration
}

// NewAIService creates a new instance of AIService
func NewAIService(
	aiAPIRepo repository.AIAPIRepository,
	encryptor crypto.Encryptor,
	maxRetries int,
	retryDelay time.Duration,
) AIService {
	return &aiService{
		aiAPIRepo:  aiAPIRepo,
		encryptor:  encryptor,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// TrainingPlanParams holds parameters for training plan generation
type TrainingPlanParams struct {
	UserID          int64
	PlanName        string
	DurationWeeks   int
	Goal            string
	DifficultyLevel string
	AIAPIID         int64
	Assessment      *model.FitnessAssessment
	BodyData        *model.UserBodyData
	FitnessGoals    []*model.FitnessGoal
}

// NutritionPlanParams holds parameters for nutrition plan generation
type NutritionPlanParams struct {
	UserID              int64
	PlanName            string
	DurationDays        int
	DailyCalories       float64
	ProteinRatio        float64
	CarbRatio           float64
	FatRatio            float64
	DietaryRestrictions []string
	Preferences         []string
	AIAPIID             int64
	BodyData            *model.UserBodyData
	FitnessGoals        []*model.FitnessGoal
}

// GenerateTrainingPlan generates a training plan using AI with retry logic
func (s *aiService) GenerateTrainingPlan(ctx context.Context, params *TrainingPlanParams) (*model.TrainingPlan, error) {
	// Get AI API configuration
	aiAPI, err := s.aiAPIRepo.GetByID(ctx, params.AIAPIID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI API: %w", err)
	}
	if aiAPI == nil {
		return nil, fmt.Errorf("AI API not found")
	}

	// Decrypt API key
	apiKey, err := s.encryptor.Decrypt(aiAPI.APIKeyEncrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}

	// Get AI client
	client, err := GetAIClient(aiAPI.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client: %w", err)
	}

	// Build prompt
	prompt := s.buildTrainingPlanPrompt(params)

	// Create client config
	config := NewAIClientFromModel(aiAPI, apiKey)

	// Call AI with retry logic (including parse errors)
	var lastErr error
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * s.retryDelay
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		response, err := client.Call(ctx, prompt, config)
		if err != nil {
			lastErr = err
			continue
		}

		planData, err := s.parseTrainingPlanResponse(response)
		if err != nil {
			lastErr = err
			continue
		}

		// Create training plan model
		startDate := time.Now()
		endDate := startDate.AddDate(0, 0, params.DurationWeeks*7)

		trainingPlan := &model.TrainingPlan{
			UserID:          params.UserID,
			PlanName:        params.PlanName,
			StartDate:       startDate,
			EndDate:         endDate,
			TotalWeeks:      params.DurationWeeks,
			DifficultyLevel: params.DifficultyLevel,
			TrainingPurpose: &params.Goal,
			AIAPIID:         params.AIAPIID,
			PlanData:        planData,
			Status:          "active",
		}

		return trainingPlan, nil
	}

	return nil, fmt.Errorf("failed to generate training plan after %d attempts: %w", s.maxRetries+1, lastErr)
}

// GenerateNutritionPlan generates a nutrition plan using AI with retry logic
func (s *aiService) GenerateNutritionPlan(ctx context.Context, params *NutritionPlanParams) (*model.NutritionPlan, error) {
	// Get AI API configuration
	aiAPI, err := s.aiAPIRepo.GetByID(ctx, params.AIAPIID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI API: %w", err)
	}
	if aiAPI == nil {
		return nil, fmt.Errorf("AI API not found")
	}

	// Decrypt API key
	apiKey, err := s.encryptor.Decrypt(aiAPI.APIKeyEncrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}

	// Get AI client
	client, err := GetAIClient(aiAPI.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client: %w", err)
	}

	// Build prompt
	prompt := s.buildNutritionPlanPrompt(params)

	// Create client config
	config := NewAIClientFromModel(aiAPI, apiKey)

	// Call AI with retry logic (including parse errors)
	var lastErr error
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * s.retryDelay
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		response, err := client.Call(ctx, prompt, config)
		if err != nil {
			lastErr = err
			continue
		}

		planData, err := s.parseNutritionPlanResponse(response)
		if err != nil {
			lastErr = err
			continue
		}

		// Create nutrition plan model
		startDate := time.Now()
		endDate := startDate.AddDate(0, 0, params.DurationDays)

		nutritionPlan := &model.NutritionPlan{
			UserID:              params.UserID,
			PlanName:            params.PlanName,
			StartDate:           startDate,
			EndDate:             endDate,
			DailyCalories:       params.DailyCalories,
			ProteinRatio:        params.ProteinRatio,
			CarbRatio:           params.CarbRatio,
			FatRatio:            params.FatRatio,
			DietaryRestrictions: model.JSONSlice(interfaceSlice(params.DietaryRestrictions)),
			Preferences:         model.JSONSlice(interfaceSlice(params.Preferences)),
			PlanData:            planData,
			AIAPIID:             params.AIAPIID,
			Status:              "active",
		}

		return nutritionPlan, nil
	}

	return nil, fmt.Errorf("failed to generate nutrition plan after %d attempts: %w", s.maxRetries+1, lastErr)
}

// TestConnection tests the connection to an AI API
func (s *aiService) TestConnection(ctx context.Context, apiID int64, userID int64) error {
	// Get AI API configuration
	aiAPI, err := s.aiAPIRepo.GetByID(ctx, apiID)
	if err != nil {
		return fmt.Errorf("failed to get AI API: %w", err)
	}
	if aiAPI == nil {
		return fmt.Errorf("AI API not found")
	}

	// Verify ownership
	if aiAPI.UserID != userID {
		return fmt.Errorf("unauthorized access to AI API")
	}

	// Decrypt API key
	apiKey, err := s.encryptor.Decrypt(aiAPI.APIKeyEncrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt API key: %w", err)
	}

	// Get AI client
	client, err := GetAIClient(aiAPI.Provider)
	if err != nil {
		return fmt.Errorf("failed to get AI client: %w", err)
	}

	// Create client config
	config := NewAIClientFromModel(aiAPI, apiKey)

	// Test connection
	return client.TestConnection(ctx, config)
}

// buildTrainingPlanPrompt builds the prompt for training plan generation
func (s *aiService) buildTrainingPlanPrompt(params *TrainingPlanParams) string {
	prompt := fmt.Sprintf(`Generate a detailed %d-week training plan with the following specifications:

Goal: %s
Difficulty Level: %s
Plan Name: %s

`, params.DurationWeeks, params.Goal, params.DifficultyLevel, params.PlanName)

	// Add assessment information
	if params.Assessment != nil {
		prompt += fmt.Sprintf(`User Assessment:
- Experience Level: %s
- Weekly Available Days: %d
- Daily Available Minutes: %d
`, params.Assessment.ExperienceLevel, params.Assessment.WeeklyAvailableDays, params.Assessment.DailyAvailableMinutes)

		if params.Assessment.InjuryHistory != nil && *params.Assessment.InjuryHistory != "" {
			prompt += fmt.Sprintf("- Injury History: %s\n", *params.Assessment.InjuryHistory)
		}
		if params.Assessment.HealthConditions != nil && *params.Assessment.HealthConditions != "" {
			prompt += fmt.Sprintf("- Health Conditions: %s\n", *params.Assessment.HealthConditions)
		}
		if len(params.Assessment.EquipmentAvailable) > 0 {
			prompt += fmt.Sprintf("- Equipment Available: %v\n", params.Assessment.EquipmentAvailable)
		}
	}

	// Add body data
	if params.BodyData != nil {
		prompt += fmt.Sprintf(`
User Body Data:
- Age: %d
- Gender: %s
- Height: %.2f cm
- Weight: %.2f kg
`, params.BodyData.Age, params.BodyData.Gender, params.BodyData.Height, params.BodyData.Weight)

		if params.BodyData.BodyFatPercentage != nil {
			prompt += fmt.Sprintf("- Body Fat: %.2f%%\n", *params.BodyData.BodyFatPercentage)
		}
	}

	// Add fitness goals
	if len(params.FitnessGoals) > 0 {
		prompt += "\nFitness Goals:\n"
		for _, goal := range params.FitnessGoals {
			prompt += fmt.Sprintf("- %s", goal.GoalType)
			if goal.GoalDescription != nil {
				prompt += fmt.Sprintf(": %s", *goal.GoalDescription)
			}
			prompt += "\n"
		}
	}

	prompt += `
Please generate a comprehensive training plan in JSON format with the following structure:
{
  "weeks": [
    {
      "week": 1,
      "days": [
        {
          "day": 1,
          "date": "YYYY-MM-DD",
          "type": "strength|cardio|rest",
          "focus_area": "upper_body|lower_body|full_body|cardio",
          "exercises": [
            {
              "name": "中文动作名称",
              "sets": 4,
              "reps": "8-10",
              "weight": "70kg or bodyweight",
              "rest": "90s",
              "difficulty": "easy|medium|hard",
              "safety_notes": "标准姿势与注意事项（中文，简洁）"
            }
          ],
          "duration": 60,
          "estimated_calories": 350
        }
      ]
    }
  ]
}

Ensure the plan:
1. Progressively increases in difficulty
2. Includes proper rest days
3. Balances different muscle groups
4. Considers any injuries or health conditions
5. Fits within the user's available time
6. Includes safety notes for complex exercises
7. Uses Chinese exercise names and Chinese safety notes

Return ONLY the JSON object, no additional text.
The response must start with "{" and end with "}".
If you cannot generate the full plan, return {"weeks": []}.`

	return prompt
}

// buildNutritionPlanPrompt builds the prompt for nutrition plan generation
func (s *aiService) buildNutritionPlanPrompt(params *NutritionPlanParams) string {
	prompt := fmt.Sprintf(`Generate a detailed %d-day nutrition plan with the following specifications:

Plan Name: %s
Daily Calories: %.0f kcal
Macronutrient Ratios:
- Protein: %.0f%%
- Carbohydrates: %.0f%%
- Fat: %.0f%%

`, params.DurationDays, params.PlanName, params.DailyCalories,
		params.ProteinRatio*100, params.CarbRatio*100, params.FatRatio*100)

	// Add dietary restrictions
	if len(params.DietaryRestrictions) > 0 {
		prompt += fmt.Sprintf("Dietary Restrictions: %v\n", params.DietaryRestrictions)
	}

	// Add preferences
	if len(params.Preferences) > 0 {
		prompt += fmt.Sprintf("Preferences: %v\n", params.Preferences)
	}

	// Add body data
	if params.BodyData != nil {
		prompt += fmt.Sprintf(`
User Body Data:
- Age: %d
- Gender: %s
- Height: %.2f cm
- Weight: %.2f kg
`, params.BodyData.Age, params.BodyData.Gender, params.BodyData.Height, params.BodyData.Weight)
	}

	// Add fitness goals
	if len(params.FitnessGoals) > 0 {
		prompt += "\nFitness Goals:\n"
		for _, goal := range params.FitnessGoals {
			prompt += fmt.Sprintf("- %s", goal.GoalType)
			if goal.GoalDescription != nil {
				prompt += fmt.Sprintf(": %s", *goal.GoalDescription)
			}
			prompt += "\n"
		}
	}

	prompt += `
Please generate a comprehensive nutrition plan in JSON format with the following structure:
{
  "days": [
    {
      "day": 1,
      "date": "YYYY-MM-DD",
      "meals": {
        "breakfast": {
          "time": "07:00-08:00",
          "foods": [
            {
              "name": "Food name",
              "amount": "100g",
              "calories": 200,
              "protein": 10,
              "carbs": 25,
              "fat": 5,
              "fiber": 3
            }
          ],
          "total_calories": 450
        },
        "lunch": { ... },
        "dinner": { ... },
        "snacks": { ... }
      },
      "daily_totals": {
        "calories": 2000,
        "protein": 150,
        "carbs": 200,
        "fat": 67
      }
    }
  ]
}

Ensure the plan:
1. Meets the specified calorie and macro targets
2. Respects all dietary restrictions
3. Includes variety across days
4. Provides balanced nutrition
5. Includes meal timing suggestions
6. Lists specific portion sizes

Return ONLY the JSON object, no additional text.
The response must start with "{" and end with "}".
If you cannot generate the full plan, return {"days": []}.`

	return prompt
}

// parseTrainingPlanResponse parses the AI response for training plan
func (s *aiService) parseTrainingPlanResponse(response string) (model.JSONMap, error) {
	// Try to extract JSON from response (AI might add extra text)
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	var planData model.JSONMap
	if err := json.Unmarshal([]byte(jsonStr), &planData); err != nil {
		var weeks []interface{}
		if err := json.Unmarshal([]byte(jsonStr), &weeks); err == nil {
			planData = model.JSONMap{
				"weeks": weeks,
			}
		} else {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	// Validate structure
	if _, ok := planData["weeks"]; !ok {
		return nil, fmt.Errorf("invalid plan structure: missing 'weeks' field")
	}

	return planData, nil
}

// parseNutritionPlanResponse parses the AI response for nutrition plan
func (s *aiService) parseNutritionPlanResponse(response string) (model.JSONMap, error) {
	// Try to extract JSON from response (AI might add extra text)
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	var planData model.JSONMap
	if err := json.Unmarshal([]byte(jsonStr), &planData); err != nil {
		var days []interface{}
		if err := json.Unmarshal([]byte(jsonStr), &days); err == nil {
			planData = model.JSONMap{
				"days": days,
			}
		} else {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	// Validate structure
	if _, ok := planData["days"]; !ok {
		return nil, fmt.Errorf("invalid plan structure: missing 'days' field")
	}

	return planData, nil
}

// extractJSON extracts JSON object from a string that might contain additional text
func extractJSON(s string) string {
	// Find first { and last }
	start := -1
	end := -1
	braceCount := 0

	for i, c := range s {
		if c == '{' {
			if start == -1 {
				start = i
			}
			braceCount++
		} else if c == '}' {
			braceCount--
			if braceCount == 0 && start != -1 {
				end = i + 1
				break
			}
		}
	}

	if start == -1 || end == -1 {
		return extractJSONArray(s)
	}

	return s[start:end]
}

func extractJSONArray(s string) string {
	start := -1
	end := -1
	bracketCount := 0

	for i, c := range s {
		if c == '[' {
			if start == -1 {
				start = i
			}
			bracketCount++
		} else if c == ']' {
			bracketCount--
			if bracketCount == 0 && start != -1 {
				end = i + 1
				break
			}
		}
	}

	if start == -1 || end == -1 {
		return ""
	}

	return s[start:end]
}

// interfaceSlice converts a string slice to an interface slice
func interfaceSlice(strings []string) []interface{} {
	result := make([]interface{}, len(strings))
	for i, s := range strings {
		result[i] = s
	}
	return result
}
