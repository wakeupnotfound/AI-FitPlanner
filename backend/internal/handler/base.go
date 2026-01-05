package handler

import (
	"net/http"
	"time"

	"github.com/ai-fitness-planner/backend/internal/api/response"
	apperrors "github.com/ai-fitness-planner/backend/internal/errors"
	"github.com/ai-fitness-planner/backend/internal/middleware"
	"github.com/ai-fitness-planner/backend/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BaseHandler provides common handler utilities
type BaseHandler struct{}

// NewBaseHandler creates a new BaseHandler instance
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// Success sends a successful response with data
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response.Success(data))
}

// SuccessWithMessage sends a successful response with a custom message
func (h *BaseHandler) SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, &response.BaseResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Created sends a 201 Created response
func (h *BaseHandler) Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, response.Success(data))
}

// NoContent sends a 204 No Content response
func (h *BaseHandler) NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error handles error responses based on error type
func (h *BaseHandler) Error(c *gin.Context, err error) {
	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		// Unknown error - log and return generic error
		logger.Error("Unexpected error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.InternalServerError("服务器内部错误"))
		return
	}

	// Map error code to HTTP status
	httpStatus := h.mapErrorCodeToHTTPStatus(appErr.Code)
	c.JSON(httpStatus, response.Error(appErr.Code, appErr.Message))
}

// BadRequest sends a 400 Bad Request response
func (h *BaseHandler) BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, response.BadRequestError(message))
}

// Unauthorized sends a 401 Unauthorized response
func (h *BaseHandler) Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.UnauthorizedError(message))
}

// Forbidden sends a 403 Forbidden response
func (h *BaseHandler) Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, response.ForbiddenError(message))
}

// NotFound sends a 404 Not Found response
func (h *BaseHandler) NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, response.NotFoundError(message))
}

// InternalError sends a 500 Internal Server Error response
func (h *BaseHandler) InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, response.InternalServerError(message))
}

// mapErrorCodeToHTTPStatus maps application error codes to HTTP status codes
func (h *BaseHandler) mapErrorCodeToHTTPStatus(code int) int {
	switch {
	case code >= 4000 && code < 4010:
		return http.StatusBadRequest
	case code >= 4010 && code < 4030:
		return http.StatusUnauthorized
	case code >= 4030 && code < 4040:
		return http.StatusForbidden
	case code >= 4040 && code < 4050:
		return http.StatusNotFound
	case code >= 4050 && code < 4090:
		return http.StatusMethodNotAllowed
	case code >= 4090 && code < 5000:
		return http.StatusConflict
	case code >= 5000 && code < 6000:
		return http.StatusInternalServerError
	case code >= 6000:
		// Business errors - map to appropriate HTTP status
		switch code {
		case apperrors.ErrUserExists:
			return http.StatusConflict
		case apperrors.ErrUserNotFound, apperrors.ErrPlanNotFound:
			return http.StatusNotFound
		case apperrors.ErrWrongPassword, apperrors.ErrInvalidCredentials:
			return http.StatusUnauthorized
		case apperrors.ErrTokenExpired:
			return http.StatusUnauthorized
		case apperrors.ErrAiApiNotConfigured:
			return http.StatusBadRequest
		case apperrors.ErrApiLimitExceeded:
			return http.StatusTooManyRequests
		default:
			return http.StatusBadRequest
		}
	default:
		return http.StatusInternalServerError
	}
}

// GetUserID extracts user ID from context, returns error response if not found
func (h *BaseHandler) GetUserID(c *gin.Context) (int64, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		h.Unauthorized(c, "用户未认证")
		return 0, false
	}
	return userID, true
}

// GetSessionID extracts session ID from context
func (h *BaseHandler) GetSessionID(c *gin.Context) (string, bool) {
	return middleware.GetSessionID(c)
}

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// GetPagination extracts and validates pagination parameters with defaults
func (h *BaseHandler) GetPagination(c *gin.Context) (page, limit, offset int) {
	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params.Page = 1
		params.Limit = 20
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	offset = (params.Page - 1) * params.Limit
	return params.Page, params.Limit, offset
}

// BuildPaginationInfo creates pagination info for response
func (h *BaseHandler) BuildPaginationInfo(page, limit int, total int64) response.PaginationInfo {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return response.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// BindJSON binds JSON request body and handles validation errors
func (h *BaseHandler) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.BadRequest(c, "请求参数无效: "+err.Error())
		return false
	}
	return true
}

// BindQuery binds query parameters and handles validation errors
func (h *BaseHandler) BindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		h.BadRequest(c, "查询参数无效: "+err.Error())
		return false
	}
	return true
}

// BindURI binds URI parameters and handles validation errors
func (h *BaseHandler) BindURI(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		h.BadRequest(c, "路径参数无效: "+err.Error())
		return false
	}
	return true
}

// ValidateMacroRatioSum validates that macro nutrient ratios sum to approximately 1.0
// Validates: Requirements 6.3
func (h *BaseHandler) ValidateMacroRatioSum(c *gin.Context, protein, carb, fat float64) bool {
	sum := protein + carb + fat
	tolerance := 0.01
	if sum < (1.0-tolerance) || sum > (1.0+tolerance) {
		h.BadRequest(c, "宏量营养素比例之和必须等于1.0")
		return false
	}
	return true
}

// ValidateDateRange validates that start date is before or equal to end date
// Validates: Requirements 7.1
func (h *BaseHandler) ValidateDateRange(c *gin.Context, startDate, endDate string) bool {
	if startDate == "" || endDate == "" {
		return true
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		h.BadRequest(c, "开始日期格式无效")
		return false
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		h.BadRequest(c, "结束日期格式无效")
		return false
	}

	if start.After(end) {
		h.BadRequest(c, "开始日期不能晚于结束日期")
		return false
	}

	return true
}
