package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryConfig holds recovery middleware configuration
type RecoveryConfig struct {
	// EnableStackTrace enables logging of stack traces
	EnableStackTrace bool
	// StackTraceSize is the maximum size of stack trace to log
	StackTraceSize int
}

// DefaultRecoveryConfig returns default recovery configuration
func DefaultRecoveryConfig() *RecoveryConfig {
	return &RecoveryConfig{
		EnableStackTrace: true,
		StackTraceSize:   4096,
	}
}

// RecoveryMiddleware creates recovery middleware that catches panics
func RecoveryMiddleware(config *RecoveryConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultRecoveryConfig()
	}

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				var stackTrace string
				if config.EnableStackTrace {
					stack := debug.Stack()
					if len(stack) > config.StackTraceSize {
						stack = stack[:config.StackTraceSize]
					}
					stackTrace = string(stack)
				}

				// Get request info for logging
				requestID := c.GetHeader("X-Request-ID")
				userID, _ := GetUserID(c)

				// Log the panic
				fields := []zap.Field{
					zap.Any("panic", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.String("user_agent", c.Request.UserAgent()),
				}

				if requestID != "" {
					fields = append(fields, zap.String("request_id", requestID))
				}

				if userID > 0 {
					fields = append(fields, zap.Int64("user_id", userID))
				}

				if stackTrace != "" {
					fields = append(fields, zap.String("stack_trace", stackTrace))
				}

				logger.Error("服务器发生panic", fields...)

				// Abort with 500 error
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerError("服务器内部错误"))
			}
		}()

		c.Next()
	}
}

// RecoveryWithWriter creates recovery middleware with custom error writer
func RecoveryWithWriter(config *RecoveryConfig, errorHandler func(c *gin.Context, err interface{})) gin.HandlerFunc {
	if config == nil {
		config = DefaultRecoveryConfig()
	}

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				var stackTrace string
				if config.EnableStackTrace {
					stack := debug.Stack()
					if len(stack) > config.StackTraceSize {
						stack = stack[:config.StackTraceSize]
					}
					stackTrace = string(stack)
				}

				// Log the panic
				logger.Error("服务器发生panic",
					zap.Any("panic", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("stack_trace", stackTrace),
				)

				// Call custom error handler if provided
				if errorHandler != nil {
					errorHandler(c, err)
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerError("服务器内部错误"))
				}
			}
		}()

		c.Next()
	}
}
