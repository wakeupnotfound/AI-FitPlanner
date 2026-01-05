package handler

import (
	"time"

	"github.com/ai-fitness-planner/backend/internal/api/request"
	"github.com/ai-fitness-planner/backend/internal/api/response"
	apperrors "github.com/ai-fitness-planner/backend/internal/errors"
	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// AssessmentHandler handles fitness assessment HTTP requests
// Requirements: 4.1, 4.2, 4.3, 4.4
type AssessmentHandler struct {
	*BaseHandler
	assessmentRepo repository.AssessmentRepository
}

// NewAssessmentHandler creates a new AssessmentHandler instance
func NewAssessmentHandler(assessmentRepo repository.AssessmentRepository) *AssessmentHandler {
	return &AssessmentHandler{
		BaseHandler:    NewBaseHandler(),
		assessmentRepo: assessmentRepo,
	}
}

// CreateAssessment handles POST /api/v1/assessments
// Requirements: 4.1, 4.2, 4.4
func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.CreateAssessmentRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Parse assessment date
	assessmentDate, err := time.Parse("2006-01-02", req.AssessmentDate)
	if err != nil {
		h.BadRequest(c, "无效的评估日期格式")
		return
	}

	// Create assessment model
	assessment := &model.FitnessAssessment{
		UserID:                userID,
		ExperienceLevel:       req.ExperienceLevel,
		WeeklyAvailableDays:   req.WeeklyAvailableDays,
		DailyAvailableMinutes: req.DailyAvailableMinutes,
		ActivityType:          req.ActivityType,
		InjuryHistory:         req.InjuryHistory,
		HealthConditions:      req.HealthConditions,
		AssessmentDate:        assessmentDate,
		CreatedAt:             time.Now(),
	}

	// Convert slices to JSONSlice
	if len(req.PreferredDays) > 0 {
		assessment.PreferredDays = make(model.JSONSlice, len(req.PreferredDays))
		for i, day := range req.PreferredDays {
			assessment.PreferredDays[i] = day
		}
	}

	if len(req.EquipmentAvailable) > 0 {
		assessment.EquipmentAvailable = make(model.JSONSlice, len(req.EquipmentAvailable))
		for i, equip := range req.EquipmentAvailable {
			assessment.EquipmentAvailable[i] = equip
		}
	}

	// Save assessment
	if err := h.assessmentRepo.Create(c.Request.Context(), assessment); err != nil {
		h.Error(c, apperrors.Wrap(err, apperrors.ErrDatabase, "保存评估数据失败"))
		return
	}

	// Build response
	resp := h.buildAssessmentInfo(assessment)
	h.Created(c, response.AssessmentDetailResponse{Assessment: resp})
}

// GetLatestAssessment handles GET /api/v1/assessments/latest
// Requirements: 4.3
func (h *AssessmentHandler) GetLatestAssessment(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	assessment, err := h.assessmentRepo.GetLatest(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, apperrors.Wrap(err, apperrors.ErrDatabase, "获取评估数据失败"))
		return
	}

	if assessment == nil {
		h.NotFound(c, "未找到评估记录")
		return
	}

	resp := h.buildAssessmentInfo(assessment)
	h.Success(c, response.AssessmentDetailResponse{Assessment: resp})
}

// buildAssessmentInfo converts model to response format
func (h *AssessmentHandler) buildAssessmentInfo(assessment *model.FitnessAssessment) response.AssessmentInfo {
	info := response.AssessmentInfo{
		ID:                    assessment.ID,
		ExperienceLevel:       assessment.ExperienceLevel,
		WeeklyAvailableDays:   assessment.WeeklyAvailableDays,
		DailyAvailableMinutes: assessment.DailyAvailableMinutes,
		AssessmentDate:        assessment.AssessmentDate.Format("2006-01-02"),
		CreatedAt:             assessment.CreatedAt.Format(time.RFC3339),
	}

	if assessment.ActivityType != nil {
		info.ActivityType = *assessment.ActivityType
	}
	if assessment.InjuryHistory != nil {
		info.InjuryHistory = *assessment.InjuryHistory
	}
	if assessment.HealthConditions != nil {
		info.HealthConditions = *assessment.HealthConditions
	}

	// Convert JSONSlice to string slice
	if len(assessment.PreferredDays) > 0 {
		info.PreferredDays = make([]string, len(assessment.PreferredDays))
		for i, day := range assessment.PreferredDays {
			if s, ok := day.(string); ok {
				info.PreferredDays[i] = s
			}
		}
	}

	if len(assessment.EquipmentAvailable) > 0 {
		info.EquipmentAvailable = make([]string, len(assessment.EquipmentAvailable))
		for i, equip := range assessment.EquipmentAvailable {
			if s, ok := equip.(string); ok {
				info.EquipmentAvailable[i] = s
			}
		}
	}

	return info
}
