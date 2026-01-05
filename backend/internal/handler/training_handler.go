package handler

import (
	"strconv"
	"time"

	"github.com/ai-fitness-planner/backend/internal/api/request"
	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// TrainingHandler handles training-related HTTP requests
// Requirements: 5.1, 5.2, 5.4, 5.5, 5.6, 7.1, 7.2, 7.3, 7.4
type TrainingHandler struct {
	*BaseHandler
	trainingService service.TrainingService
}

// NewTrainingHandler creates a new TrainingHandler instance
func NewTrainingHandler(trainingService service.TrainingService) *TrainingHandler {
	return &TrainingHandler{
		BaseHandler:     NewBaseHandler(),
		trainingService: trainingService,
	}
}

// GeneratePlan handles POST /api/v1/training-plans/generate
// Requirements: 5.1, 5.2
func (h *TrainingHandler) GeneratePlan(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.GenerateTrainingPlanRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Convert to service request
	serviceReq := &service.GeneratePlanRequest{
		PlanName:        req.PlanName,
		DurationWeeks:   req.DurationWeeks,
		Goal:            req.Goal,
		DifficultyLevel: req.DifficultyLevel,
		AIAPIID:         req.AIAPIID,
	}

	taskResp, err := h.trainingService.GeneratePlan(c.Request.Context(), userID, serviceReq)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.TaskResponse{
		TaskID:        taskResp.TaskID,
		Status:        taskResp.Status,
		Progress:      0,
		EstimatedTime: 60, // Estimated 60 seconds
	}

	h.Success(c, resp)
}

// GetPlanStatus handles GET /api/v1/training-plans/tasks/:taskId
func (h *TrainingHandler) GetPlanStatus(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		h.BadRequest(c, "任务ID不能为空")
		return
	}

	taskStatus, err := h.trainingService.GetPlanStatus(c.Request.Context(), taskID)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.TaskResponse{
		TaskID:   taskStatus.TaskID,
		Status:   taskStatus.Status,
		Progress: taskStatus.Progress,
	}

	if taskStatus.Error != "" {
		resp.ErrorMessage = taskStatus.Error
	}

	if taskStatus.Result != nil {
		resp.Result = h.buildPlanInfo(taskStatus.Result)
	}

	h.Success(c, resp)
}

// ListPlans handles GET /api/v1/training-plans
// Requirements: 5.5
func (h *TrainingHandler) ListPlans(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var params request.TrainingPlanListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		// Use defaults if binding fails
		params.Status = ""
	}

	plans, err := h.trainingService.ListPlans(c.Request.Context(), userID, params.Status)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Convert to response format
	planInfos := make([]response.PlanInfo, 0, len(plans))
	for _, plan := range plans {
		planInfos = append(planInfos, h.buildPlanInfo(plan))
	}

	page, limit, _ := h.GetPagination(c)
	resp := response.PlanListResponse{
		Plans:      planInfos,
		Pagination: h.BuildPaginationInfo(page, limit, int64(len(plans))),
	}

	h.Success(c, resp)
}

// GetPlanDetail handles GET /api/v1/training-plans/:id
// Requirements: 5.4
func (h *TrainingHandler) GetPlanDetail(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.BadRequest(c, "无效的计划ID")
		return
	}

	plan, err := h.trainingService.GetPlanDetail(c.Request.Context(), planID, userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	h.Success(c, gin.H{
		"plan": plan,
	})
}

// GetTodayTraining handles GET /api/v1/training-plans/today
// Requirements: 5.6
func (h *TrainingHandler) GetTodayTraining(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	dayPlan, err := h.trainingService.GetTodayTraining(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	if dayPlan == nil {
		h.Success(c, response.TodayTrainingResponse{
			Schedule: response.TodaySchedule{
				Date:        time.Now().Format("2006-01-02"),
				Type:        "rest",
				FocusArea:   "",
				Exercises:   []response.ExerciseInfo{},
				Duration:    0,
				IsCompleted: false,
			},
		})
		return
	}

	// Convert exercises to response format
	exercises := make([]response.ExerciseInfo, 0, len(dayPlan.Exercises))
	for _, ex := range dayPlan.Exercises {
		exercises = append(exercises, response.ExerciseInfo{
			Name:        ex.Name,
			Sets:        ex.Sets,
			Reps:        ex.Reps,
			Weight:      ex.Weight,
			Rest:        ex.Rest,
			Difficulty:  ex.Difficulty,
			SafetyNotes: ex.SafetyNotes,
		})
	}

	resp := response.TodayTrainingResponse{
		Schedule: response.TodaySchedule{
			Date:           dayPlan.Date,
			Type:           dayPlan.Type,
			FocusArea:      dayPlan.FocusArea,
			Exercises:      exercises,
			Duration:       dayPlan.Duration,
			IsCompleted:    false,
			TotalExercises: len(exercises),
		},
	}

	h.Success(c, resp)
}

// RecordTraining handles POST /api/v1/training-records
// Requirements: 7.1, 7.2, 7.3
func (h *TrainingHandler) RecordTraining(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.RecordTrainingRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Parse workout date
	workoutDate, err := time.Parse("2006-01-02", req.WorkoutDate)
	if err != nil {
		h.BadRequest(c, "无效的训练日期格式")
		return
	}

	// Create record model
	record := &model.TrainingRecord{
		UserID:          userID,
		PlanID:          req.PlanID,
		WorkoutDate:     workoutDate,
		WorkoutType:     req.WorkoutType,
		DurationMinutes: req.DurationMinutes,
		Notes:           req.Notes,
		Rating:          req.Rating,
		InjuryReport:    req.InjuryReport,
		CreatedAt:       time.Now(),
	}

	// Convert exercises and performance data to JSONMap
	if req.Exercises != nil {
		record.Exercises = model.JSONMap(req.Exercises)
	}
	if req.PerformanceData != nil {
		record.PerformanceData = model.JSONMap(req.PerformanceData)
	}

	if err := h.trainingService.RecordTraining(c.Request.Context(), userID, record); err != nil {
		h.Error(c, err)
		return
	}

	h.Created(c, gin.H{
		"id":      record.ID,
		"message": "训练记录已保存",
	})
}

// ListTrainingRecords handles GET /api/v1/training-records
// Requirements: 7.4
func (h *TrainingHandler) ListTrainingRecords(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var params request.TrainingRecordListParams
	_ = c.ShouldBindQuery(&params)

	var startDate, endDate *time.Time
	if params.StartDate != "" {
		t, err := time.Parse("2006-01-02", params.StartDate)
		if err == nil {
			startDate = &t
		}
	}
	if params.EndDate != "" {
		t, err := time.Parse("2006-01-02", params.EndDate)
		if err == nil {
			endDate = &t
		}
	}

	// Use the trainingService's GetTrainingHistory method via type assertion
	type historyGetter interface {
		GetTrainingHistory(ctx interface{}, userID int64, startDate, endDate *time.Time) ([]*model.TrainingRecord, error)
	}

	if getter, ok := h.trainingService.(historyGetter); ok {
		records, err := getter.GetTrainingHistory(c.Request.Context(), userID, startDate, endDate)
		if err != nil {
			h.Error(c, err)
			return
		}

		page, limit, _ := h.GetPagination(c)
		h.Success(c, gin.H{
			"records":    records,
			"pagination": h.BuildPaginationInfo(page, limit, int64(len(records))),
		})
		return
	}

	// Fallback if method not available
	h.Success(c, gin.H{
		"records":    []interface{}{},
		"pagination": h.BuildPaginationInfo(1, 20, 0),
	})
}

// buildPlanInfo converts model to response format
func (h *TrainingHandler) buildPlanInfo(plan *model.TrainingPlan) response.PlanInfo {
	return response.PlanInfo{
		ID:              plan.ID,
		Name:            plan.PlanName,
		StartDate:       plan.StartDate.Format("2006-01-02"),
		EndDate:         plan.EndDate.Format("2006-01-02"),
		TotalWeeks:      plan.TotalWeeks,
		DifficultyLevel: plan.DifficultyLevel,
		Status:          plan.Status,
	}
}
