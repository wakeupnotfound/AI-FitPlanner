package service

import (
	"context"
	"sync"
	"time"

	"github.com/ai-fitness-planner/backend/internal/errors"
	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/repository"
	"github.com/google/uuid"
)

// TrainingService defines the interface for training operations
type TrainingService interface {
	// GeneratePlan generates a training plan asynchronously and returns a task ID
	GeneratePlan(ctx context.Context, userID int64, req *GeneratePlanRequest) (*TaskResponse, error)
	// GetPlanStatus retrieves the status of a plan generation task
	GetPlanStatus(ctx context.Context, taskID string) (*TaskStatus, error)
	// ListPlans retrieves training plans for a user with optional status filter
	ListPlans(ctx context.Context, userID int64, status string) ([]*model.TrainingPlan, error)
	// GetPlanDetail retrieves a specific training plan
	GetPlanDetail(ctx context.Context, planID int64, userID int64) (*model.TrainingPlan, error)
	// GetTodayTraining retrieves today's training schedule
	GetTodayTraining(ctx context.Context, userID int64) (*model.DayPlan, error)
	// RecordTraining records a training session with validation
	RecordTraining(ctx context.Context, userID int64, record *model.TrainingRecord) error
}

// GeneratePlanRequest holds parameters for plan generation request
type GeneratePlanRequest struct {
	PlanName        string `json:"plan_name" validate:"required,min=1,max=200"`
	DurationWeeks   int    `json:"duration_weeks" validate:"required,min=1,max=52"`
	Goal            string `json:"goal" validate:"required,max=100"`
	DifficultyLevel string `json:"difficulty_level" validate:"required,oneof=easy medium hard extreme"`
	AIAPIID         *int64 `json:"ai_api_id"` // Optional, uses default if not provided
}

// TaskResponse represents the response for async task creation
type TaskResponse struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// TaskStatus represents the status of an async task
type TaskStatus struct {
	TaskID    string              `json:"task_id"`
	Status    string              `json:"status"` // pending, processing, completed, failed
	Progress  int                 `json:"progress"`
	Message   string              `json:"message,omitempty"`
	Error     string              `json:"error,omitempty"`
	Result    *model.TrainingPlan `json:"result,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// Task status constants
const (
	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
)

// trainingService implements TrainingService interface
type trainingService struct {
	planRepo        repository.TrainingPlanRepository
	recordRepo      repository.TrainingRecordRepository
	aiAPIRepo       repository.AIAPIRepository
	assessmentRepo  repository.AssessmentRepository
	bodyDataRepo    repository.BodyDataRepository
	fitnessGoalRepo repository.FitnessGoalRepository
	aiService       AIService

	// In-memory task storage (in production, use Redis)
	tasks      map[string]*TaskStatus
	tasksMutex sync.RWMutex
}

// NewTrainingService creates a new instance of TrainingService
func NewTrainingService(
	planRepo repository.TrainingPlanRepository,
	recordRepo repository.TrainingRecordRepository,
	aiAPIRepo repository.AIAPIRepository,
	assessmentRepo repository.AssessmentRepository,
	bodyDataRepo repository.BodyDataRepository,
	fitnessGoalRepo repository.FitnessGoalRepository,
	aiService AIService,
) TrainingService {
	return &trainingService{
		planRepo:        planRepo,
		recordRepo:      recordRepo,
		aiAPIRepo:       aiAPIRepo,
		assessmentRepo:  assessmentRepo,
		bodyDataRepo:    bodyDataRepo,
		fitnessGoalRepo: fitnessGoalRepo,
		aiService:       aiService,
		tasks:           make(map[string]*TaskStatus),
	}
}

// GeneratePlan generates a training plan asynchronously
// Requirements: 5.1, 5.2, 5.4
func (s *trainingService) GeneratePlan(ctx context.Context, userID int64, req *GeneratePlanRequest) (*TaskResponse, error) {
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
	task := &TaskStatus{
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
		Message: "训练计划生成任务已创建",
	}, nil
}

// processGeneratePlan handles the async plan generation
func (s *trainingService) processGeneratePlan(userID int64, req *GeneratePlanRequest, aiAPIID int64, taskID string) {
	ctx := context.Background()

	// Update task status to processing
	s.updateTaskStatus(taskID, TaskStatusProcessing, 10, "正在收集用户数据...", "", nil)

	// Get user's latest assessment
	assessment, err := s.assessmentRepo.GetLatest(ctx, userID)
	if err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "获取用户评估数据失败: "+err.Error(), nil)
		return
	}

	s.updateTaskStatus(taskID, TaskStatusProcessing, 20, "正在获取身体数据...", "", nil)

	// Get user's latest body data
	bodyData, err := s.bodyDataRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "获取身体数据失败: "+err.Error(), nil)
		return
	}

	s.updateTaskStatus(taskID, TaskStatusProcessing, 30, "正在获取健身目标...", "", nil)

	// Get user's fitness goals
	fitnessGoals, err := s.fitnessGoalRepo.GetByUserID(ctx, userID, "active")
	if err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "获取健身目标失败: "+err.Error(), nil)
		return
	}

	s.updateTaskStatus(taskID, TaskStatusProcessing, 50, "正在调用AI生成训练计划...", "", nil)

	// Build AI params
	params := &TrainingPlanParams{
		UserID:          userID,
		PlanName:        req.PlanName,
		DurationWeeks:   req.DurationWeeks,
		Goal:            req.Goal,
		DifficultyLevel: req.DifficultyLevel,
		AIAPIID:         aiAPIID,
		Assessment:      assessment,
		BodyData:        bodyData,
		FitnessGoals:    fitnessGoals,
	}

	// Generate plan using AI service
	plan, err := s.aiService.GenerateTrainingPlan(ctx, params)
	if err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "AI生成计划失败: "+err.Error(), nil)
		return
	}

	s.updateTaskStatus(taskID, TaskStatusProcessing, 80, "正在保存训练计划...", "", nil)

	// Save the plan to database
	if err := s.planRepo.Create(ctx, plan); err != nil {
		s.updateTaskStatus(taskID, TaskStatusFailed, 0, "", "保存计划失败: "+err.Error(), nil)
		return
	}

	// Update task status to completed
	s.updateTaskStatus(taskID, TaskStatusCompleted, 100, "训练计划生成完成", "", plan)
}

// updateTaskStatus updates the status of a task
func (s *trainingService) updateTaskStatus(taskID, status string, progress int, message, errMsg string, result *model.TrainingPlan) {
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
func (s *trainingService) GetPlanStatus(ctx context.Context, taskID string) (*TaskStatus, error) {
	s.tasksMutex.RLock()
	defer s.tasksMutex.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, errors.New(errors.ErrNotFound, "任务不存在")
	}

	return task, nil
}

// ListPlans retrieves training plans for a user with optional status filter
// Requirements: 5.5
func (s *trainingService) ListPlans(ctx context.Context, userID int64, status string) ([]*model.TrainingPlan, error) {
	plans, err := s.planRepo.ListByUser(ctx, userID, status)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取训练计划列表失败")
	}
	return plans, nil
}

// GetPlanDetail retrieves a specific training plan
// Requirements: 5.4
func (s *trainingService) GetPlanDetail(ctx context.Context, planID int64, userID int64) (*model.TrainingPlan, error) {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取训练计划失败")
	}
	if plan == nil {
		return nil, errors.New(errors.ErrPlanNotFound, "训练计划不存在")
	}

	// Verify ownership
	if plan.UserID != userID {
		return nil, errors.New(errors.ErrForbidden, "无权访问此训练计划")
	}

	return plan, nil
}

// GetTodayTraining retrieves today's training schedule
// Requirements: 5.6
func (s *trainingService) GetTodayTraining(ctx context.Context, userID int64) (*model.DayPlan, error) {
	today := time.Now()

	dayPlan, err := s.planRepo.GetTodaySchedule(ctx, userID, today)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取今日训练失败")
	}

	// Return nil if no training scheduled for today (not an error)
	return dayPlan, nil
}

// RecordTraining records a training session with validation
// Requirements: 7.1, 7.2
func (s *trainingService) RecordTraining(ctx context.Context, userID int64, record *model.TrainingRecord) error {
	// Validate that workout date is not in the future
	// Property 12: Future Date Rejection
	now := time.Now()
	// Truncate to start of day for comparison
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	localWorkout := record.WorkoutDate.In(time.Local)
	workoutDate := time.Date(localWorkout.Year(), localWorkout.Month(), localWorkout.Day(), 0, 0, 0, 0, time.Local)

	if workoutDate.After(todayStart) {
		return errors.New(errors.ErrInvalidParam, "训练日期不能是未来日期")
	}

	// Set user ID
	record.UserID = userID

	// Validate plan ID if provided
	if record.PlanID != nil {
		plan, err := s.planRepo.GetByID(ctx, *record.PlanID)
		if err != nil {
			return errors.Wrap(err, errors.ErrDatabase, "验证训练计划失败")
		}
		if plan == nil || plan.UserID != userID {
			return errors.New(errors.ErrNotFound, "训练计划不存在")
		}
	}

	// Create the record
	if err := s.recordRepo.Create(ctx, record); err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "保存训练记录失败")
	}

	return nil
}

// GetTrainingHistory retrieves training records for a user
// Requirements: 7.4
func (s *trainingService) GetTrainingHistory(ctx context.Context, userID int64, startDate, endDate *time.Time) ([]*model.TrainingRecord, error) {
	records, err := s.recordRepo.ListByUser(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取训练记录失败")
	}
	return records, nil
}

// GetTrainingStatistics retrieves aggregated training statistics
// Requirements: 7.5
func (s *trainingService) GetTrainingStatistics(ctx context.Context, userID int64, startDate, endDate time.Time) (*repository.TrainingStatistics, error) {
	stats, err := s.recordRepo.GetStatistics(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "获取训练统计失败")
	}
	return stats, nil
}
