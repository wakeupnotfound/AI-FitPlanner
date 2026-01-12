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

// NutritionHandler handles nutrition-related HTTP requests
// Requirements: 6.1, 6.2, 6.3, 6.4, 8.1, 8.2, 8.3, 8.4
type NutritionHandler struct {
	*BaseHandler
	nutritionService service.NutritionService
}

// NewNutritionHandler creates a new NutritionHandler instance
func NewNutritionHandler(nutritionService service.NutritionService) *NutritionHandler {
	return &NutritionHandler{
		BaseHandler:      NewBaseHandler(),
		nutritionService: nutritionService,
	}
}

// GeneratePlan handles POST /api/v1/nutrition-plans/generate
// Requirements: 6.1, 6.2
func (h *NutritionHandler) GeneratePlan(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.GenerateNutritionPlanRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Validate macro ratio sum (Requirements 6.3)
	if !h.ValidateMacroRatioSum(c, req.ProteinRatio, req.CarbRatio, req.FatRatio) {
		return
	}

	// Convert to service request
	serviceReq := &service.GenerateNutritionPlanRequest{
		PlanName:            req.PlanName,
		DurationDays:        req.DurationDays,
		DailyCalories:       req.DailyCalories,
		ProteinRatio:        req.ProteinRatio,
		CarbRatio:           req.CarbRatio,
		FatRatio:            req.FatRatio,
		DietaryRestrictions: req.DietaryRestrictions,
		Preferences:         req.Preferences,
		AIAPIID:             req.AIAPIID,
	}

	taskResp, err := h.nutritionService.GeneratePlan(c.Request.Context(), userID, serviceReq)
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

// GetPlanStatus handles GET /api/v1/nutrition-plans/tasks/:taskId
// Requirements: 6.2
func (h *NutritionHandler) GetPlanStatus(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		h.BadRequest(c, "任务ID不能为空")
		return
	}

	taskStatus, err := h.nutritionService.GetPlanStatus(c.Request.Context(), taskID)
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

// ListPlans handles GET /api/v1/nutrition-plans
// Requirements: 6.3
func (h *NutritionHandler) ListPlans(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var params request.NutritionPlanListParams
	_ = c.ShouldBindQuery(&params)

	plans, err := h.nutritionService.ListPlans(c.Request.Context(), userID, params.Status)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Convert to response format
	planInfos := make([]response.NutritionPlanInfo, 0, len(plans))
	for _, plan := range plans {
		planInfos = append(planInfos, h.buildPlanInfo(plan))
	}

	page, limit, _ := h.GetPagination(c)
	resp := response.NutritionPlanListResponse{
		Plans:      planInfos,
		Pagination: h.BuildPaginationInfo(page, limit, int64(len(plans))),
	}

	h.Success(c, resp)
}

// GetPlanDetail handles GET /api/v1/nutrition-plans/:id
// Requirements: 6.3
func (h *NutritionHandler) GetPlanDetail(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.BadRequest(c, "无效的计划ID")
		return
	}

	plan, err := h.nutritionService.GetPlanDetail(c.Request.Context(), planID, userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	h.Success(c, gin.H{
		"plan": plan,
	})
}

// GetTodayMeals handles GET /api/v1/nutrition-plans/today
// Requirements: 6.4
func (h *NutritionHandler) GetTodayMeals(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	meals, err := h.nutritionService.GetTodayMeals(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Convert to response format
	mealMap := make(map[string]response.MealInfo)
	for _, meal := range meals {
		mealMap[meal.Time] = response.MealInfo{
			Time:          meal.Time,
			Foods:         meal.Foods,
			TotalCalories: meal.TotalCalories,
		}
	}

	resp := response.TodayNutritionResponse{
		Plan:  response.TodayNutritionPlanInfo{},
		Meals: mealMap,
	}

	h.Success(c, resp)
}

// RecordMeal handles POST /api/v1/nutrition-records
// Requirements: 8.1
func (h *NutritionHandler) RecordMeal(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.RecordMealRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Parse meal date
	mealDate, err := time.Parse("2006-01-02", req.MealDate)
	if err != nil {
		h.BadRequest(c, "无效的日期格式")
		return
	}

	// Create record model
	record := &model.NutritionRecord{
		UserID:    userID,
		MealDate:  mealDate,
		MealTime:  req.MealType,
		Calories:  req.Calories,
		Protein:   req.Protein,
		Carbs:     req.Carbs,
		Fat:       req.Fat,
		Fiber:     req.Fiber,
		CreatedAt: time.Now(),
	}

	// Convert foods to JSONMap
	if req.Foods != nil {
		record.Foods = model.JSONMap(req.Foods)
	}

	if err := h.nutritionService.RecordMeal(c.Request.Context(), userID, record); err != nil {
		h.Error(c, err)
		return
	}

	h.Created(c, gin.H{
		"id":      record.ID,
		"message": "饮食记录已保存",
	})
}

// ListNutritionRecords handles GET /api/v1/nutrition-records
// Requirements: 8.4
func (h *NutritionHandler) ListNutritionRecords(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var params request.NutritionRecordListParams
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

	records, err := h.nutritionService.GetNutritionHistory(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Convert to response format
	recordInfos := make([]response.NutritionRecordInfo, 0, len(records))
	for _, record := range records {
		recordInfos = append(recordInfos, h.buildRecordInfo(record))
	}

	page, limit, _ := h.GetPagination(c)
	resp := response.NutritionRecordListResponse{
		Records:    recordInfos,
		Pagination: h.BuildPaginationInfo(page, limit, int64(len(records))),
	}

	h.Success(c, resp)
}

// GetDailySummary handles GET /api/v1/nutrition-records/daily-summary
// Requirements: 8.2
func (h *NutritionHandler) GetDailySummary(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var params request.DailySummaryParams
	if !h.BindQuery(c, &params) {
		return
	}

	date, err := time.Parse("2006-01-02", params.Date)
	if err != nil {
		h.BadRequest(c, "无效的日期格式")
		return
	}

	summary, err := h.nutritionService.GetDailySummary(c.Request.Context(), userID, date)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.DailySummaryResponse{
		Date:          params.Date,
		TotalCalories: summary.TotalCalories,
		TotalProtein:  summary.TotalProtein,
		TotalCarbs:    summary.TotalCarbs,
		TotalFat:      summary.TotalFat,
		TotalFiber:    summary.TotalFiber,
		MealCount:     int(summary.MealCount),
	}

	h.Success(c, resp)
}

// buildPlanInfo converts model to response format
func (h *NutritionHandler) buildPlanInfo(plan *model.NutritionPlan) response.NutritionPlanInfo {
	info := response.NutritionPlanInfo{
		ID:            plan.ID,
		PlanName:      plan.PlanName,
		StartDate:     plan.StartDate.Format("2006-01-02"),
		EndDate:       plan.EndDate.Format("2006-01-02"),
		DailyCalories: plan.DailyCalories,
		ProteinRatio:  plan.ProteinRatio,
		CarbRatio:     plan.CarbRatio,
		FatRatio:      plan.FatRatio,
		Status:        plan.Status,
		CreatedAt:     plan.CreatedAt.Format(time.RFC3339),
	}

	// Convert JSONSlice to string slice
	if len(plan.DietaryRestrictions) > 0 {
		info.DietaryRestrictions = make([]string, 0, len(plan.DietaryRestrictions))
		for _, r := range plan.DietaryRestrictions {
			if s, ok := r.(string); ok {
				info.DietaryRestrictions = append(info.DietaryRestrictions, s)
			}
		}
	}

	if len(plan.Preferences) > 0 {
		info.Preferences = make([]string, 0, len(plan.Preferences))
		for _, p := range plan.Preferences {
			if s, ok := p.(string); ok {
				info.Preferences = append(info.Preferences, s)
			}
		}
	}

	return info
}

// buildRecordInfo converts model to response format
func (h *NutritionHandler) buildRecordInfo(record *model.NutritionRecord) response.NutritionRecordInfo {
	info := response.NutritionRecordInfo{
		ID:        record.ID,
		MealDate:  record.MealDate.Format("2006-01-02"),
		MealType:  record.MealTime,
		Calories:  record.Calories,
		Protein:   record.Protein,
		Carbs:     record.Carbs,
		Fat:       record.Fat,
		Fiber:     record.Fiber,
		CreatedAt: record.CreatedAt.Format(time.RFC3339),
	}

	if record.Foods != nil {
		info.Foods = map[string]interface{}(record.Foods)
	}

	return info
}
