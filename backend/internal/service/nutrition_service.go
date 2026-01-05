package service

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/ai-fitness-planner/backend/internal/errors"
	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/repository"
	"github.com/google/uuid"
)

// NutritionService defines the interface for nutrition operations
type NutritionService interface {
	// GeneratePlan generates a nutrition plan asynchronously and returns a task ID
	GeneratePlan(ctx context.Context, userID int64, req *GenerateNutritionPlanRequest) (*TaskResponse, error)
	// GetPlanStatus retrieves the status of a plan generation task
	GetPlanStatus(ctx context.Context, taskID string) (*NutritionTaskStatus, error)
	// ListPlans retrieves nutrition plans for a user with optional status filter
	ListPlans(ctx context.Context, userID int64, status string) ([]*model.NutritionPlan, error)
	// GetPlanDetail retrieves a specific nutrition plan
	GetPlanDetail(ctx context.Context, planID int64, userID int64) (*model.NutritionPlan, error)
	// GetTodayMeals retrieves today's meal plan
	GetTodayMeals(ctx context.Context, userID int64) ([]model.NutritionPlanMeal, error)
	// RecordMeal records a meal with nutrition calculation
	RecordMeal(ctx context.Context, userID int64, record *model.NutritionRecord) error
	// GetDailySummary retrieves aggregated nutrition data for a specific day
	GetDailySummary(ctx context.Context, userID int64, date time.Time) (*repository.DailyNutritionSummary, error)
	// GetNutritionHistory retrieves nutrition records for a user
	GetNutritionHistory(ctx context.Context, userID int64, startDate, endDate *time.Time) ([]*model.NutritionRecord, error)
}

// GenerateNutritionPlanRequest holds parameters for nutrition plan generation request
type GenerateNutritionPlanRequest struct {
	PlanName            string   `json:"plan_name" validate:"required,min=1,max=200"`
	DurationDays        int      `json:"duration_days" validate:"required,min=1,max=365"`
	DailyCalories       *float64 `json:"daily_calories"` // Optional, calculated if not provided
	ProteinRatio        float64  `json:"protein_ratio" validate:"required,min=0,max=1"`
	CarbRatio           float64  `json:"carb_ratio" validate:"required,min=0,max=1"`
	FatRatio            float64  `json:"fat_ratio" validate:"required,min=0,max=1"`
	DietaryRestrictions []string `json:"dietary_restrictions"`
	Preferences         []string `json:"preferences"`
	AIAPIID             *int64   `json:"ai_api_id"` // Optional, uses default if not provided
}

// NutritionTaskStatus represents the status of an async nutrition task
type NutritionTaskStatus struct {
	TaskID    string               `json:"task_id"`
	Status    string               `json:"status"` // pending, processing, completed, failed
	Progress  int                  `json:"progress"`
	Message   string               `json:"message,omitempty"`
	Error     string               `json:"error,omitempty"`
	Result    *model.NutritionPlan `json:"result,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// nutritionService implements NutritionService interface
type nutritionService struct {
	planRepo        repository.NutritionPlanRepository
	recordRepo      repository.NutritionRecordRepository
	aiAPIRepo       repository.AIAPIRepository
	bodyDataRepo    repository.BodyDataRepository
	fitnessGoalRepo repository.FitnessGoalRepository
	aiService       AIService

	// In-memory task storage (in production, use Redis)
	tasks      map[string]*NutritionTaskStatus
	tasksMutex sync.RWMutex
}

// NewNutritionService creates a new instance of NutritionService
func NewNutritionService(
	planRepo repository.NutritionPlanRepository,
	recordRepo repository.NutritionRecordRepository,
	aiAPIRepo repository.AIAPIRepository,
	bodyDataRepo repository.BodyDataRepository,
	fitnessGoalRepo repository.FitnessGoalRepository,
	aiService AIService,
) NutritionService {
	return &nutritionService{
		planRepo:        planRepo,
		recordRepo:      recordRepo,
		aiAPIRepo:       aiAPIRepo,
		bodyDataRepo:    bodyDataRepo,
		fitnessGoalRepo: fitnessGoalRepo,
		aiService:       aiService,
		tasks:           make(map[string]*NutritionTaskStatus),
	}
}

// GeneratePlan generates a nutrition plan asynchronously
// Requirements: 6.1, 6.2, 6.3
func (s *nutritionService) GeneratePlan(ctx context.Context, userID int64, req *GenerateNutritionPlanRequest) (*TaskResponse, error) {
	// Validate macro nutrient ratios sum to 1.0 (±0.01 tolerance)
	// Property 11: Macro Nutrient Ratio Sum
	ratioSum := req.ProteinRatio + req.CarbRatio + req.FatRatio
	if math.Abs(ratioSum-1.0) > 0.01 {
		return nil, errors.New(errors.ErrInvalidParam, "宏量营养素比例之和必须等于100%")
	}

	// Determine which AI API to use
	var aiAPIID int64
	if req.AIAPIID != nil {
		aiAPIID = *req.AIAPIID
		// Verify the API exists and belongs to the user
		api, err := s.aiAPIRepo.GetByID(ctx, aiAPIID)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrDatabase, "获取AI API失败")
		}
		if api == nil || api.UserID != userID {
			return nil, errors.New(errors.ErrNotFound, "AI API不存在")
		}
	} else {
		// Use default AI API
		defaultAPI, err := s.aiAPIRepo.GetDefaultByUser(ctx, userID)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrDatabase, "获取默认AI API失败")
		}
		if defaultAPI == nil {
			return nil, errors.ErrNoDefaultAIAPI
		}
		aiAPIID = defaultAPI.ID
	}

	// Create task ID
	taskID := uuid.New().String()

	// Initialize task status
	now := time.Now()
	task := &NutritionTaskStatus{
		TaskID:    taskID,
		Status:    TaskStatusPending,
		Progress:  0,
		Message:   "任务已创建，等待处理",
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.tasksMutex.Lock()
	s.tasks[taskID] = task
	s.tasksMutex.Unlock()

	// Start async generation
	go s.processGeneratePlan(userID, req, aiAPIID, taskID)

	return &TaskResponse{
		TaskID:  taskID,
		Status:  TaskStatusPending,
		Message: "饮食计划生成任务已创建",
	}, nil
}

// processGeneratePlan handles the async plan generation
func (s *nutritionService) processGeneratePlan(userID int64, req *GenerateNutritionPlanRequest, aiAPIID int64, taskID string) {
	ctx := context.Background()

	// Update task status to processing
	s.updateTaskStatus(taskID, TaskStatusProcessing, 10, "正在收集用户数据...", "", nil)

	// Get user's latest body data
	bodyData, err := s.bodyDataRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "获取身体数据失败: "+err.Error(), nil)
		return
	}

	s.updateTaskStatus(taskID, TaskStatusProcessing, 20, "正在获取健身目标...", "", nil)

	// Get user's fitness goals
	fitnessGoals, err := s.fitnessGoalRepo.GetByUserID(ctx, userID, "active")
	if err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "获取健身目标失败: "+err.Error(), nil)
		return
	}

	s.updateTaskStatus(taskID, TaskStatusProcessing, 30, "正在计算每日热量需求...", "", nil)

	// Calculate daily calories if not provided
	// Requirements: 6.1 - Calculate daily calorie needs based on body data
	dailyCalories := req.DailyCalories
	if dailyCalories == nil || *dailyCalories <= 0 {
		calculatedCalories := s.calculateDailyCalories(bodyData, fitnessGoals)
		dailyCalories = &calculatedCalories
	}

	s.updateTaskStatus(taskID, TaskStatusProcessing, 50, "正在调用AI生成饮食计划...", "", nil)

	// Build AI params
	params := &NutritionPlanParams{
		UserID:              userID,
		PlanName:            req.PlanName,
		DurationDays:        req.DurationDays,
		DailyCalories:       *dailyCalories,
		ProteinRatio:        req.ProteinRatio,
		CarbRatio:           req.CarbRatio,
		FatRatio:            req.FatRatio,
		DietaryRestrictions: req.DietaryRestrictions,
		Preferences:         req.Preferences,
		AIAPIID:             aiAPIID,
		BodyData:            bodyData,
		FitnessGoals:        fitnessGoals,
	}

	// Generate plan using AI service
	plan, err := s.aiService.GenerateNutritionPlan(ctx, params)
	if err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "AI生成计划失败: "+err.Error(), nil)
		return
	}

	s.updateTaskStatus(taskID, TaskStatusProcessing, 80, "正在保存饮食计划...", "", nil)

	// Save the plan to database
	if err := s.planRepo.Create(ctx, plan); err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "保存计划失败: "+err.Error(), nil)
		return
	}

	// Update task status to completed
	s.updateTaskStatus(taskID, TaskStatusCompleted, 100, "饮食计划生成完成", "", plan)
}

// calculateDailyCalories calculates daily calorie needs based on body data and goals
// Uses Mifflin-St Jeor equation for BMR calculation
// Requirements: 6.1
func (s *nutritionService) calculateDailyCalories(bodyData *model.UserBodyData, goals []*model.FitnessGoal) float64 {
	if bodyData == nil {
		// Default to 2000 calories if no body data available
		return 2000.0
	}

	// Calculate BMR using Mifflin-St Jeor equation
	var bmr float64
	if bodyData.Gender == "male" {
		// Men: BMR = 10 × weight(kg) + 6.25 × height(cm) - 5 × age(years) + 5
		bmr = 10*bodyData.Weight + 6.25*bodyData.Height - 5*float64(bodyData.Age) + 5
	} else {
		// Women: BMR = 10 × weight(kg) + 6.25 × height(cm) - 5 × age(years) - 161
		bmr = 10*bodyData.Weight + 6.25*bodyData.Height - 5*float64(bodyData.Age) - 161
	}

	// Apply activity multiplier (default to moderate activity: 1.55)
	activityMultiplier := 1.55
	tdee := bmr * activityMultiplier

	// Adjust based on fitness goals
	for _, goal := range goals {
		switch goal.GoalType {
		case "weight_loss", "fat_loss", "减脂", "减重":
			// Create a caloric deficit (15-20%)
			tdee *= 0.85
		case "muscle_gain", "bulk", "增肌":
			// Create a caloric surplus (10-15%)
			tdee *= 1.15
		case "maintenance", "maintain", "维持":
			// Keep TDEE as is
		}
		// Only apply the first matching goal
		break
	}

	// Round to nearest 50
	return math.Round(tdee/50) * 50
}

// updateTaskStatus updates the status of a task
func (s *nutritionService) updateTaskStatus(taskID, status string, progress int, message, errMsg string, result *model.NutritionPlan) {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()

	if task, exists := s.tasks[taskID]; exists {
		task.Status = status
		task.Progress = progress
		task.Message = message
		task.Error = errMsg
		task.Result = result
		task.UpdatedAt = time.Now()
	}
}

// GetPlanStatus retrieves the status of a plan generation task
func (s *nutritionService) GetPlanStatus(ctx context.Context, taskID string) (*NutritionTaskStatus, error) {
	s.tasksMutex.RLock()
	defer s.tasksMutex.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, errors.New(errors.ErrNotFound, "任务不存在")
	}

	return task, nil
}

// ListPlans retrieves nutrition plans for a user with optional status filter
// Requirements: 6.3
func (s *nutritionService) ListPlans(ctx context.Context, userID int64, status string) ([]*model.NutritionPlan, error) {
	plans, err := s.planRepo.ListByUser(ctx, userID, status)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取饮食计划列表失败")
	}
	return plans, nil
}

// GetPlanDetail retrieves a specific nutrition plan
// Requirements: 6.3
func (s *nutritionService) GetPlanDetail(ctx context.Context, planID int64, userID int64) (*model.NutritionPlan, error) {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取饮食计划失败")
	}
	if plan == nil {
		return nil, errors.New(errors.ErrPlanNotFound, "饮食计划不存在")
	}

	// Verify ownership
	if plan.UserID != userID {
		return nil, errors.New(errors.ErrForbidden, "无权访问此饮食计划")
	}

	return plan, nil
}

// GetTodayMeals retrieves today's meal plan
// Requirements: 6.4
func (s *nutritionService) GetTodayMeals(ctx context.Context, userID int64) ([]model.NutritionPlanMeal, error) {
	today := time.Now()

	meals, err := s.planRepo.GetTodayMeals(ctx, userID, today)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取今日餐食失败")
	}

	// Return empty slice if no meals scheduled for today (not an error)
	if meals == nil {
		return []model.NutritionPlanMeal{}, nil
	}

	return meals, nil
}

// RecordMeal records a meal with nutrition calculation
// Requirements: 8.1, 8.2
func (s *nutritionService) RecordMeal(ctx context.Context, userID int64, record *model.NutritionRecord) error {
	// Set user ID
	record.UserID = userID

	// Calculate total nutrition from foods if not already set
	// Requirements: 8.1 - Calculate total calories and macronutrients
	if record.Foods != nil {
		totalCalories, totalProtein, totalCarbs, totalFat, totalFiber := s.calculateNutritionFromFoods(record.Foods)

		// Only override if not explicitly set
		if record.Calories == 0 {
			record.Calories = totalCalories
		}
		if record.Protein == 0 {
			record.Protein = totalProtein
		}
		if record.Carbs == 0 {
			record.Carbs = totalCarbs
		}
		if record.Fat == 0 {
			record.Fat = totalFat
		}
		if record.Fiber == 0 {
			record.Fiber = totalFiber
		}
	}

	// Create the record
	if err := s.recordRepo.Create(ctx, record); err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "保存饮食记录失败")
	}

	return nil
}

// calculateNutritionFromFoods calculates total nutrition values from foods JSON
func (s *nutritionService) calculateNutritionFromFoods(foods model.JSONMap) (calories, protein, carbs, fat, fiber float64) {
	// Try to extract foods array from the JSON map
	foodsInterface, ok := foods["items"]
	if !ok {
		// Try alternative key
		foodsInterface, ok = foods["foods"]
		if !ok {
			return 0, 0, 0, 0, 0
		}
	}

	foodsArray, ok := foodsInterface.([]interface{})
	if !ok {
		return 0, 0, 0, 0, 0
	}

	for _, foodInterface := range foodsArray {
		foodMap, ok := foodInterface.(map[string]interface{})
		if !ok {
			continue
		}

		if cal, ok := foodMap["calories"].(float64); ok {
			calories += cal
		}
		if p, ok := foodMap["protein"].(float64); ok {
			protein += p
		}
		if c, ok := foodMap["carbs"].(float64); ok {
			carbs += c
		}
		if f, ok := foodMap["fat"].(float64); ok {
			fat += f
		}
		if fb, ok := foodMap["fiber"].(float64); ok {
			fiber += fb
		}
	}

	return calories, protein, carbs, fat, fiber
}

// GetDailySummary retrieves aggregated nutrition data for a specific day
// Requirements: 8.2
func (s *nutritionService) GetDailySummary(ctx context.Context, userID int64, date time.Time) (*repository.DailyNutritionSummary, error) {
	summary, err := s.recordRepo.GetDailySummary(ctx, userID, date)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取每日营养摘要失败")
	}

	return summary, nil
}

// GetNutritionHistory retrieves nutrition records for a user
// Requirements: 8.4
func (s *nutritionService) GetNutritionHistory(ctx context.Context, userID int64, startDate, endDate *time.Time) ([]*model.NutritionRecord, error) {
	records, err := s.recordRepo.ListByUser(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取饮食记录失败")
	}
	return records, nil
}
