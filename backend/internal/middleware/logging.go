package middleware

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/ai-fitness-planner/backend/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Sensitive field patterns to mask in logs
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)"password"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"api_key"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"apikey"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"secret"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"access_token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"refresh_token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"authorization"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"credit_card"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"card_number"\s*:\s*"[^"]*"`),
}

// Sensitive headers to mask
var sensitiveHeaders = map[string]bool{
	"authorization": true,
	"x-api-key":     true,
	"cookie":        true,
	"set-cookie":    true,
}

// LoggingConfig holds logging middleware configuration
type LoggingConfig struct {
	// Skip logging for these paths
	SkipPaths []string
	// Log request body (be careful with sensitive data)
	LogRequestBody bool
	// Log response body
	LogResponseBody bool
	// Maximum body size to log (in bytes)
	MaxBodyLogSize int
}

// DefaultLoggingConfig returns default logging configuration
func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		SkipPaths:       []string{"/health", "/metrics"},
		LogRequestBody:  false,
		LogResponseBody: false,
		MaxBodyLogSize:  1024,
	}
}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware creates logging middleware for request/response logging
func LoggingMiddleware(config *LoggingConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultLoggingConfig()
	}

	return func(c *gin.Context) {
		// Skip logging for certain paths
		path := c.Request.URL.Path
		for _, skipPath := range config.SkipPaths {
			if path == skipPath || strings.HasPrefix(path, skipPath) {
				c.Next()
				return
			}
		}

		// Start timer
		start := time.Now()

		// Get request ID if exists
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)

		// Read request body if configured
		var requestBody string
		if config.LogRequestBody && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// Restore body for handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				if len(bodyBytes) <= config.MaxBodyLogSize {
					requestBody = maskSensitiveData(string(bodyBytes))
				} else {
					requestBody = "[body too large]"
				}
			}
		}

		// Wrap response writer if logging response body
		var blw *responseWriter
		if config.LogResponseBody {
			blw = &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user ID if authenticated
		userID, _ := GetUserID(c)

		// Build log fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", maskSensitiveData(c.Request.URL.RawQuery)),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("response_size", c.Writer.Size()),
		}

		// Add user ID if authenticated
		if userID > 0 {
			fields = append(fields, zap.Int64("user_id", userID))
		}

		// Add request body if configured
		if config.LogRequestBody && requestBody != "" {
			fields = append(fields, zap.String("request_body", requestBody))
		}

		// Add response body if configured
		if config.LogResponseBody && blw != nil {
			responseBody := blw.body.String()
			if len(responseBody) <= config.MaxBodyLogSize {
				fields = append(fields, zap.String("response_body", maskSensitiveData(responseBody)))
			}
		}

		// Add error if exists
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// Log based on status code
		status := c.Writer.Status()
		switch {
		case status >= 500:
			logger.Error("HTTP请求", fields...)
		case status >= 400:
			logger.Warn("HTTP请求", fields...)
		default:
			logger.Info("HTTP请求", fields...)
		}
	}
}

// maskSensitiveData masks sensitive information in the input string
func maskSensitiveData(input string) string {
	result := input

	// Mask sensitive JSON fields
	for _, pattern := range sensitivePatterns {
		result = pattern.ReplaceAllStringFunc(result, func(match string) string {
			// Find the field name and replace value with [MASKED]
			parts := strings.SplitN(match, ":", 2)
			if len(parts) == 2 {
				return parts[0] + `: "[MASKED]"`
			}
			return match
		})
	}

	return result
}

// MaskHeaders masks sensitive headers for logging
func MaskHeaders(headers map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if sensitiveHeaders[lowerKey] {
			result[key] = "[MASKED]"
		} else if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
