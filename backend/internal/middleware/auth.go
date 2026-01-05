package middleware

import (
	"net/http"
	"strings"

	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/pkg/jwt"
	"github.com/ai-fitness-planner/backend/internal/pkg/logger"
	"github.com/ai-fitness-planner/backend/internal/pkg/session"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Context keys for user information
const (
	ContextKeyUserID    = "user_id"
	ContextKeyUsername  = "username"
	ContextKeySessionID = "session_id"
)

// AuthMiddleware creates authentication middleware with JWT validation and session verification
func AuthMiddleware(jwtManager jwt.JWTManager, sessionManager session.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.UnauthorizedError("缺少认证令牌"))
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.UnauthorizedError("无效的认证格式"))
			return
		}

		tokenString := parts[1]

		// Validate JWT token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			logger.Warn("JWT验证失败",
				zap.Error(err),
				zap.String("ip", c.ClientIP()),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.UnauthorizedError("无效或过期的令牌"))
			return
		}

		// Verify it's an access token
		if claims.Type != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.UnauthorizedError("无效的令牌类型"))
			return
		}

		// Verify session exists in Redis
		sess, err := sessionManager.GetSession(c.Request.Context(), claims.SessionID)
		if err != nil {
			logger.Error("获取会话失败",
				zap.Error(err),
				zap.String("session_id", claims.SessionID),
			)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerError("会话验证失败"))
			return
		}

		if sess == nil {
			logger.Warn("会话不存在或已过期",
				zap.String("session_id", claims.SessionID),
				zap.Int64("user_id", claims.UserID),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.UnauthorizedError("会话不存在或已过期"))
			return
		}

		// Verify session belongs to the same user
		if sess.UserID != claims.UserID {
			logger.Warn("会话用户不匹配",
				zap.Int64("token_user_id", claims.UserID),
				zap.Int64("session_user_id", sess.UserID),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.UnauthorizedError("会话验证失败"))
			return
		}

		// Set user info in context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeySessionID, claims.SessionID)

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

// GetUsername extracts username from context
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

// GetSessionID extracts session ID from context
func GetSessionID(c *gin.Context) (string, bool) {
	sessionID, exists := c.Get(ContextKeySessionID)
	if !exists {
		return "", false
	}
	id, ok := sessionID.(string)
	return id, ok
}

// OptionalAuthMiddleware creates middleware that extracts user info if token is present but doesn't require it
func OptionalAuthMiddleware(jwtManager jwt.JWTManager, sessionManager session.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		if claims.Type != "access" {
			c.Next()
			return
		}

		sess, err := sessionManager.GetSession(c.Request.Context(), claims.SessionID)
		if err != nil || sess == nil || sess.UserID != claims.UserID {
			c.Next()
			return
		}

		// Set user info in context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeySessionID, claims.SessionID)

		c.Next()
	}
}
