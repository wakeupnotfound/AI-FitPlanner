package middleware

import (
	"html"
	"net/http"
	"regexp"
	"strings"

	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SQL injection patterns to detect
var sqlInjectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(\b(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|ALTER|CREATE|TRUNCATE|EXEC|EXECUTE)\b)`),
	regexp.MustCompile(`(?i)(--|#|/\*|\*/)`),
	regexp.MustCompile(`(?i)(\bOR\b\s+\d+\s*=\s*\d+)`),
	regexp.MustCompile(`(?i)(\bAND\b\s+\d+\s*=\s*\d+)`),
	regexp.MustCompile(`(?i)(;\s*(SELECT|INSERT|UPDATE|DELETE|DROP))`),
	regexp.MustCompile(`(?i)('|\"|;|\\x00|\\n|\\r|\\x1a)`),
}

// XSS patterns to detect
var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
	regexp.MustCompile(`(?i)<script[^>]*>`),
	regexp.MustCompile(`(?i)javascript:`),
	regexp.MustCompile(`(?i)on\w+\s*=`),
	regexp.MustCompile(`(?i)<iframe[^>]*>`),
	regexp.MustCompile(`(?i)<object[^>]*>`),
	regexp.MustCompile(`(?i)<embed[^>]*>`),
	regexp.MustCompile(`(?i)<link[^>]*>`),
	regexp.MustCompile(`(?i)<meta[^>]*>`),
}

// SecurityConfig holds security middleware configuration
type SecurityConfig struct {
	// Enable SQL injection detection
	EnableSQLInjectionDetection bool
	// Enable XSS detection
	EnableXSSDetection bool
	// Enable security headers
	EnableSecurityHeaders bool
	// Allowed content types
	AllowedContentTypes []string
}

// DefaultSecurityConfig returns default security configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		EnableSQLInjectionDetection: true,
		EnableXSSDetection:          true,
		EnableSecurityHeaders:       true,
		AllowedContentTypes: []string{
			"application/json",
			"application/x-www-form-urlencoded",
			"multipart/form-data",
		},
	}
}

// SecurityMiddleware creates security middleware for input sanitization and security headers
func SecurityMiddleware(config *SecurityConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSecurityConfig()
	}

	return func(c *gin.Context) {
		// Add security headers
		if config.EnableSecurityHeaders {
			addSecurityHeaders(c)
		}

		// Skip input validation for GET and OPTIONS requests
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Validate content type for POST/PUT/PATCH requests
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			contentType := c.GetHeader("Content-Type")
			if contentType != "" && !isAllowedContentType(contentType, config.AllowedContentTypes) {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, response.Error(4150, "不支持的内容类型"))
				return
			}
		}

		// Check query parameters for SQL injection and XSS
		for key, values := range c.Request.URL.Query() {
			for _, value := range values {
				if config.EnableSQLInjectionDetection && containsSQLInjection(value) {
					logger.Warn("检测到SQL注入尝试",
						zap.String("ip", c.ClientIP()),
						zap.String("param", key),
						zap.String("path", c.Request.URL.Path),
					)
					c.AbortWithStatusJSON(http.StatusBadRequest, response.BadRequestError("检测到非法字符"))
					return
				}

				if config.EnableXSSDetection && containsXSS(value) {
					logger.Warn("检测到XSS尝试",
						zap.String("ip", c.ClientIP()),
						zap.String("param", key),
						zap.String("path", c.Request.URL.Path),
					)
					c.AbortWithStatusJSON(http.StatusBadRequest, response.BadRequestError("检测到非法字符"))
					return
				}
			}
		}

		c.Next()
	}
}

// addSecurityHeaders adds security-related HTTP headers
func addSecurityHeaders(c *gin.Context) {
	// Prevent MIME type sniffing
	c.Header("X-Content-Type-Options", "nosniff")

	// Enable XSS filter in browsers
	c.Header("X-XSS-Protection", "1; mode=block")

	// Prevent clickjacking
	c.Header("X-Frame-Options", "DENY")

	// Control referrer information
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

	// Content Security Policy
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")

	// Strict Transport Security (for HTTPS)
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

	// Permissions Policy
	c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
}

// containsSQLInjection checks if input contains SQL injection patterns
func containsSQLInjection(input string) bool {
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

// containsXSS checks if input contains XSS patterns
func containsXSS(input string) bool {
	for _, pattern := range xssPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

// isAllowedContentType checks if content type is allowed
func isAllowedContentType(contentType string, allowed []string) bool {
	// Extract main content type (ignore charset and other parameters)
	mainType := strings.Split(contentType, ";")[0]
	mainType = strings.TrimSpace(strings.ToLower(mainType))

	for _, allowedType := range allowed {
		if strings.ToLower(allowedType) == mainType {
			return true
		}
	}
	return false
}

// SanitizeInput sanitizes user input by escaping HTML characters
func SanitizeInput(input string) string {
	return html.EscapeString(input)
}

// SanitizeMap sanitizes all string values in a map
func SanitizeMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = SanitizeInput(v)
		case map[string]interface{}:
			result[key] = SanitizeMap(v)
		default:
			result[key] = value
		}
	}
	return result
}

// ValidateInput validates input against SQL injection and XSS patterns
// Returns true if input is safe, false otherwise
func ValidateInput(input string) bool {
	return !containsSQLInjection(input) && !containsXSS(input)
}

// StripTags removes HTML tags from input
func StripTags(input string) string {
	// Simple tag stripping - removes anything between < and >
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}
