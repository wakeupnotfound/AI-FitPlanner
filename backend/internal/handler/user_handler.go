package handler

import (
	"time"

	"github.com/ai-fitness-planner/backend/internal/api/request"
	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
// Requirements: 2.1, 2.2, 2.3, 2.4, 2.5
type UserHandler struct {
	*BaseHandler
	userService service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(),
		userService: userService,
	}
}

// GetProfile handles GET /api/v1/user/profile
// Requirements: 2.1
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserProfileResponse "User profile retrieved successfully"
// @Failure 401 {object} response.BaseResponse "Unauthorized"
// @Failure 404 {object} response.BaseResponse "User not found"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.UserProfileResponse{
		User: response.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	if user.Nickname != nil {
		resp.User.Nickname = *user.Nickname
	}
	if user.Phone != nil {
		resp.User.Phone = *user.Phone
	}
	if user.Avatar != nil {
		resp.User.Avatar = *user.Avatar
	}

	h.Success(c, resp)
}

// UpdateProfile handles PUT /api/v1/user/profile
// Requirements: 2.2
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.UpdateUserRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Convert to service request
	serviceReq := &service.UpdateProfileRequest{}
	if req.Phone != "" {
		serviceReq.Phone = &req.Phone
	}
	if req.Nickname != "" {
		serviceReq.Nickname = &req.Nickname
	}
	if req.Avatar != "" {
		serviceReq.Avatar = &req.Avatar
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, serviceReq)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	if user.Nickname != nil {
		resp.Nickname = *user.Nickname
	}
	if user.Phone != nil {
		resp.Phone = *user.Phone
	}
	if user.Avatar != nil {
		resp.Avatar = *user.Avatar
	}

	h.Success(c, resp)
}

// AddBodyData handles POST /api/v1/user/body-data
// Requirements: 2.3
func (h *UserHandler) AddBodyData(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.AddBodyDataRequest
	if !h.BindJSON(c, &req) {
		return
	}

	measurementDate, err := time.ParseInLocation("2006-01-02", req.MeasurementDate, time.Local)
	if err != nil {
		h.BadRequest(c, "无效的日期格式")
		return
	}

	// Convert to service request
	serviceReq := &service.BodyDataRequest{
		Age:               req.Age,
		Gender:            req.Gender,
		Height:            req.Height,
		Weight:            req.Weight,
		BodyFatPercentage: req.BodyFatPercentage,
		MusclePercentage:  req.MusclePercentage,
		MeasurementDate:   measurementDate,
	}

	bodyData, err := h.userService.AddBodyData(c.Request.Context(), userID, serviceReq)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.BodyDataInfo{
		ID:              bodyData.ID,
		Age:             bodyData.Age,
		Gender:          bodyData.Gender,
		Height:          bodyData.Height,
		Weight:          bodyData.Weight,
		MeasurementDate: bodyData.MeasurementDate.Format("2006-01-02"),
		CreatedAt:       bodyData.CreatedAt.Format(time.RFC3339),
	}

	if bodyData.BodyFatPercentage != nil {
		resp.BodyFatPercentage = *bodyData.BodyFatPercentage
	}
	if bodyData.MusclePercentage != nil {
		resp.MusclePercentage = *bodyData.MusclePercentage
	}

	h.Created(c, resp)
}

// GetBodyDataHistory handles GET /api/v1/user/body-data
// Requirements: 2.4
func (h *UserHandler) GetBodyDataHistory(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	bodyDataList, err := h.userService.GetBodyDataHistory(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Convert to response format
	bodyDataInfos := make([]response.BodyDataInfo, 0, len(bodyDataList))
	for _, bd := range bodyDataList {
		info := response.BodyDataInfo{
			ID:              bd.ID,
			Age:             bd.Age,
			Gender:          bd.Gender,
			Height:          bd.Height,
			Weight:          bd.Weight,
			MeasurementDate: bd.MeasurementDate.Format("2006-01-02"),
			CreatedAt:       bd.CreatedAt.Format(time.RFC3339),
		}
		if bd.BodyFatPercentage != nil {
			info.BodyFatPercentage = *bd.BodyFatPercentage
		}
		if bd.MusclePercentage != nil {
			info.MusclePercentage = *bd.MusclePercentage
		}
		bodyDataInfos = append(bodyDataInfos, info)
	}

	page, limit, _ := h.GetPagination(c)
	resp := response.BodyDataListResponse{
		BodyData:   bodyDataInfos,
		Pagination: h.BuildPaginationInfo(page, limit, int64(len(bodyDataList))),
	}

	h.Success(c, resp)
}

// SetFitnessGoals handles POST /api/v1/user/fitness-goals
// Requirements: 2.5
func (h *UserHandler) SetFitnessGoals(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.AddGoalRequest
	if !h.BindJSON(c, &req) {
		return
	}

	goalDescription := ""
	switch {
	case req.GoalDescription != "":
		goalDescription = req.GoalDescription
	case req.Notes != nil:
		goalDescription = *req.Notes
	}

	var goalDescriptionPtr *string
	if goalDescription != "" {
		goalDescriptionPtr = &goalDescription
	}

	// Convert to service request
	serviceReq := &service.FitnessGoalRequest{
		GoalType:        req.GoalType,
		GoalDescription: goalDescriptionPtr,
		TargetWeight:    req.TargetWeight,
		Priority:        1, // Default priority
	}

	if req.Priority != nil {
		serviceReq.Priority = *req.Priority
	}

	// Parse deadline if provided
	deadlineValue := req.Deadline
	if deadlineValue == nil {
		deadlineValue = req.TargetDate
	}
	if deadlineValue != nil {
		deadline, err := time.ParseInLocation("2006-01-02", *deadlineValue, time.Local)
		if err == nil {
			serviceReq.Deadline = &deadline
		}
	}

	goal, err := h.userService.SetFitnessGoals(c.Request.Context(), userID, serviceReq)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.GoalInfo{
		ID:        goal.ID,
		GoalType:  goal.GoalType,
		Priority:  goal.Priority,
		Status:    goal.Status,
		CreatedAt: goal.CreatedAt.Format(time.RFC3339),
	}

	if goal.GoalDescription != nil {
		resp.GoalDescription = *goal.GoalDescription
		resp.Notes = *goal.GoalDescription
	}
	if goal.InitialWeight != nil {
		resp.InitialWeight = *goal.InitialWeight
	}
	if goal.InitialBodyFat != nil {
		resp.InitialBodyFat = *goal.InitialBodyFat
	}
	if goal.InitialMuscle != nil {
		resp.InitialMuscle = *goal.InitialMuscle
	}
	if goal.TargetWeight != nil {
		resp.TargetWeight = *goal.TargetWeight
	}
	if goal.Deadline != nil {
		resp.Deadline = goal.Deadline.Format("2006-01-02")
		resp.TargetDate = resp.Deadline
	}

	h.Created(c, resp)
}


// GetFitnessGoals handles GET /api/v1/user/fitness-goals
// Requirements: 2.5
// @Summary Get fitness goals
// @Description Get the authenticated user's fitness goals
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse "Fitness goals retrieved successfully"
// @Failure 401 {object} response.BaseResponse "Unauthorized"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /user/fitness-goals [get]
func (h *UserHandler) GetFitnessGoals(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	goals, err := h.userService.GetFitnessGoals(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Convert to response format
	goalInfos := make([]response.GoalInfo, 0, len(goals))
	for _, goal := range goals {
		info := response.GoalInfo{
			ID:        goal.ID,
			GoalType:  goal.GoalType,
			Priority:  goal.Priority,
			Status:    goal.Status,
			CreatedAt: goal.CreatedAt.Format(time.RFC3339),
		}
		if goal.GoalDescription != nil {
			info.GoalDescription = *goal.GoalDescription
			info.Notes = *goal.GoalDescription
		}
		if goal.InitialWeight != nil {
			info.InitialWeight = *goal.InitialWeight
		}
		if goal.InitialBodyFat != nil {
			info.InitialBodyFat = *goal.InitialBodyFat
		}
		if goal.InitialMuscle != nil {
			info.InitialMuscle = *goal.InitialMuscle
		}
		if goal.TargetWeight != nil {
			info.TargetWeight = *goal.TargetWeight
		}
		if goal.Deadline != nil {
			info.Deadline = goal.Deadline.Format("2006-01-02")
			info.TargetDate = info.Deadline
		}
		goalInfos = append(goalInfos, info)
	}

	h.Success(c, map[string]interface{}{
		"goals": goalInfos,
	})
}

// UpdateFitnessGoals handles PUT /api/v1/user/fitness-goals
// Requirements: 2.5
// @Summary Update fitness goals
// @Description Update the authenticated user's fitness goals
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.AddGoalRequest true "Goal data"
// @Success 200 {object} response.BaseResponse "Fitness goals updated successfully"
// @Failure 400 {object} response.BaseResponse "Bad request"
// @Failure 401 {object} response.BaseResponse "Unauthorized"
// @Failure 404 {object} response.BaseResponse "Goal not found"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /user/fitness-goals [put]
func (h *UserHandler) UpdateFitnessGoals(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.AddGoalRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Get existing goals first
	goals, err := h.userService.GetFitnessGoals(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	// If no goals exist, create a new one
	if len(goals) == 0 {
		goalDescription := ""
		switch {
		case req.GoalDescription != "":
			goalDescription = req.GoalDescription
		case req.Notes != nil:
			goalDescription = *req.Notes
		}

		var goalDescriptionPtr *string
		if goalDescription != "" {
			goalDescriptionPtr = &goalDescription
		}

		serviceReq := &service.FitnessGoalRequest{
			GoalType:        req.GoalType,
			GoalDescription: goalDescriptionPtr,
			TargetWeight:    req.TargetWeight,
			Priority:        1,
		}
		if req.Priority != nil {
			serviceReq.Priority = *req.Priority
		}
		deadlineValue := req.Deadline
		if deadlineValue == nil {
			deadlineValue = req.TargetDate
		}
		if deadlineValue != nil {
			deadline, err := time.ParseInLocation("2006-01-02", *deadlineValue, time.Local)
			if err == nil {
				serviceReq.Deadline = &deadline
			}
		}

		goal, err := h.userService.SetFitnessGoals(c.Request.Context(), userID, serviceReq)
		if err != nil {
			h.Error(c, err)
			return
		}

		resp := response.GoalInfo{
			ID:        goal.ID,
			GoalType:  goal.GoalType,
			Priority:  goal.Priority,
			Status:    goal.Status,
			CreatedAt: goal.CreatedAt.Format(time.RFC3339),
		}
		if goal.GoalDescription != nil {
			resp.GoalDescription = *goal.GoalDescription
			resp.Notes = *goal.GoalDescription
		}
		if goal.InitialWeight != nil {
			resp.InitialWeight = *goal.InitialWeight
		}
		if goal.InitialBodyFat != nil {
			resp.InitialBodyFat = *goal.InitialBodyFat
		}
		if goal.InitialMuscle != nil {
			resp.InitialMuscle = *goal.InitialMuscle
		}
		if goal.TargetWeight != nil {
			resp.TargetWeight = *goal.TargetWeight
		}
		if goal.Deadline != nil {
			resp.Deadline = goal.Deadline.Format("2006-01-02")
			resp.TargetDate = resp.Deadline
		}

		h.Success(c, resp)
		return
	}

	// Update the first (most recent) goal
	goalToUpdate := goals[0]
	goalDescription := ""
	switch {
	case req.GoalDescription != "":
		goalDescription = req.GoalDescription
	case req.Notes != nil:
		goalDescription = *req.Notes
	}

	var goalDescriptionPtr *string
	if goalDescription != "" {
		goalDescriptionPtr = &goalDescription
	}

	serviceReq := &service.FitnessGoalRequest{
		GoalType:        req.GoalType,
		GoalDescription: goalDescriptionPtr,
		TargetWeight:    req.TargetWeight,
		Priority:        1,
	}
	if req.Priority != nil {
		serviceReq.Priority = *req.Priority
	}
	deadlineValue := req.Deadline
	if deadlineValue == nil {
		deadlineValue = req.TargetDate
	}
	if deadlineValue != nil {
		deadline, err := time.ParseInLocation("2006-01-02", *deadlineValue, time.Local)
		if err == nil {
			serviceReq.Deadline = &deadline
		}
	}

	goal, err := h.userService.UpdateFitnessGoals(c.Request.Context(), userID, goalToUpdate.ID, serviceReq)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.GoalInfo{
		ID:        goal.ID,
		GoalType:  goal.GoalType,
		Priority:  goal.Priority,
		Status:    goal.Status,
		CreatedAt: goal.CreatedAt.Format(time.RFC3339),
	}
	if goal.GoalDescription != nil {
		resp.GoalDescription = *goal.GoalDescription
		resp.Notes = *goal.GoalDescription
	}
	if goal.InitialWeight != nil {
		resp.InitialWeight = *goal.InitialWeight
	}
	if goal.InitialBodyFat != nil {
		resp.InitialBodyFat = *goal.InitialBodyFat
	}
	if goal.InitialMuscle != nil {
		resp.InitialMuscle = *goal.InitialMuscle
	}
	if goal.TargetWeight != nil {
		resp.TargetWeight = *goal.TargetWeight
	}
	if goal.Deadline != nil {
		resp.Deadline = goal.Deadline.Format("2006-01-02")
		resp.TargetDate = resp.Deadline
	}

	h.Success(c, resp)
}
