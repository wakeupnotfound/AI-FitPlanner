package handler

import (
	"net/http"
	"time"

	"github.com/ai-fitness-planner/backend/internal/api/request"
	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
// Requirements: 1.1, 1.2, 1.3, 1.4, 1.5
type AuthHandler struct {
	*BaseHandler
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(),
		authService: authService,
	}
}

// Register handles POST /api/v1/auth/register
// Requirements: 1.1, 1.2
// @Summary Register a new user
// @Description Create a new user account with username, email, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Registration details"
// @Success 201 {object} response.AuthResponse "User registered successfully"
// @Failure 400 {object} response.BaseResponse "Invalid input"
// @Failure 409 {object} response.BaseResponse "Username or email already exists"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Convert to service request
	serviceReq := &service.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	authResp, err := h.authService.Register(c.Request.Context(), serviceReq)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Build response
	resp := response.AuthResponse{
		User: response.UserInfo{
			ID:        authResp.User.ID,
			Username:  authResp.User.Username,
			Email:     authResp.User.Email,
			CreatedAt: authResp.User.CreatedAt.Format(time.RFC3339),
		},
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    3600, // 1 hour
	}

	if authResp.User.Phone != nil {
		resp.User.Phone = *authResp.User.Phone
	}
	if authResp.User.Avatar != nil {
		resp.User.Avatar = *authResp.User.Avatar
	}

	h.Created(c, resp)
}

// Login handles POST /api/v1/auth/login
// Requirements: 1.2, 1.3
// @Summary User login
// @Description Authenticate user with username/email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login credentials"
// @Success 200 {object} response.AuthResponse "Login successful"
// @Failure 400 {object} response.BaseResponse "Invalid input"
// @Failure 401 {object} response.BaseResponse "Invalid credentials"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Get client info for session
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Convert to service request
	serviceReq := &service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	authResp, err := h.authService.Login(c.Request.Context(), serviceReq, ipAddress, userAgent)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Build response
	resp := response.AuthResponse{
		User: response.UserInfo{
			ID:        authResp.User.ID,
			Username:  authResp.User.Username,
			Email:     authResp.User.Email,
			CreatedAt: authResp.User.CreatedAt.Format(time.RFC3339),
		},
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    3600, // 1 hour
	}

	if authResp.User.Phone != nil {
		resp.User.Phone = *authResp.User.Phone
	}
	if authResp.User.Avatar != nil {
		resp.User.Avatar = *authResp.User.Avatar
	}

	h.Success(c, resp)
}

// Logout handles POST /api/v1/auth/logout
// Requirements: 1.5
// @Summary User logout
// @Description Invalidate current user session
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse "Logout successful"
// @Failure 401 {object} response.BaseResponse "Unauthorized"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, ok := h.GetSessionID(c)
	if !ok {
		h.Unauthorized(c, "会话不存在")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), sessionID); err != nil {
		h.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"message": "登出成功",
	}))
}

// RefreshToken handles POST /api/v1/auth/refresh
// Requirements: 1.4
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.RefreshTokenResponse "Token refreshed successfully"
// @Failure 400 {object} response.BaseResponse "Invalid input"
// @Failure 401 {object} response.BaseResponse "Invalid refresh token"
// @Failure 500 {object} response.BaseResponse "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if !h.BindJSON(c, &req) {
		return
	}

	tokenResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.RefreshTokenResponse{
		AccessToken: tokenResp.AccessToken,
		ExpiresIn:   3600, // 1 hour
	}

	h.Success(c, resp)
}
