package handler

import (
	"strconv"

	"github.com/ai-fitness-planner/backend/internal/api/request"
	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AIAPIHandler handles AI API configuration HTTP requests
// Requirements: 3.1, 3.2, 3.3, 3.4, 3.5
type AIAPIHandler struct {
	*BaseHandler
	aiAPIService service.AIAPIService
}

// NewAIAPIHandler creates a new AIAPIHandler instance
func NewAIAPIHandler(aiAPIService service.AIAPIService) *AIAPIHandler {
	return &AIAPIHandler{
		BaseHandler:  NewBaseHandler(),
		aiAPIService: aiAPIService,
	}
}

// AddAPI handles POST /api/v1/ai-apis
// Requirements: 3.1
func (h *AIAPIHandler) AddAPI(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var req request.AddAIAPIRequest
	if !h.BindJSON(c, &req) {
		return
	}

	apiInfo, err := h.aiAPIService.AddAPI(c.Request.Context(), userID, &req)
	if err != nil {
		h.Error(c, err)
		return
	}

	h.Created(c, response.AIAPIDetailResponse{API: *apiInfo})
}

// ListAPIs handles GET /api/v1/ai-apis
// Requirements: 3.2
func (h *AIAPIHandler) ListAPIs(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	listResp, err := h.aiAPIService.ListAPIs(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	h.Success(c, listResp)
}

// GetAPI handles GET /api/v1/ai-apis/:id
func (h *AIAPIHandler) GetAPI(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	apiID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.BadRequest(c, "无效的API ID")
		return
	}

	apiInfo, err := h.aiAPIService.GetAPI(c.Request.Context(), userID, apiID)
	if err != nil {
		h.Error(c, err)
		return
	}

	h.Success(c, response.AIAPIDetailResponse{API: *apiInfo})
}

// UpdateAPI handles PUT /api/v1/ai-apis/:id
func (h *AIAPIHandler) UpdateAPI(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	apiID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.BadRequest(c, "无效的API ID")
		return
	}

	var req request.UpdateAIAPIRequest
	if !h.BindJSON(c, &req) {
		return
	}

	apiInfo, err := h.aiAPIService.UpdateAPI(c.Request.Context(), userID, apiID, &req)
	if err != nil {
		h.Error(c, err)
		return
	}

	h.Success(c, response.AIAPIDetailResponse{API: *apiInfo})
}

// DeleteAPI handles DELETE /api/v1/ai-apis/:id
// Requirements: 3.5
func (h *AIAPIHandler) DeleteAPI(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	apiID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.BadRequest(c, "无效的API ID")
		return
	}

	if err := h.aiAPIService.DeleteAPI(c.Request.Context(), userID, apiID); err != nil {
		h.Error(c, err)
		return
	}

	h.NoContent(c)
}

// TestAPI handles POST /api/v1/ai-apis/:id/test
// Requirements: 3.3
func (h *AIAPIHandler) TestAPI(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	apiID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.BadRequest(c, "无效的API ID")
		return
	}

	testResp, err := h.aiAPIService.TestAPI(c.Request.Context(), userID, apiID)
	if err != nil {
		h.Error(c, err)
		return
	}

	h.Success(c, testResp)
}

// SetDefault handles POST /api/v1/ai-apis/:id/set-default
// Requirements: 3.4
func (h *AIAPIHandler) SetDefault(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	apiID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.BadRequest(c, "无效的API ID")
		return
	}

	if err := h.aiAPIService.SetDefault(c.Request.Context(), userID, apiID); err != nil {
		h.Error(c, err)
		return
	}

	h.Success(c, gin.H{"message": "已设置为默认API"})
}
