package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS middleware configuration
type CORSConfig struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from
	// Use "*" to allow all origins (not recommended for production)
	AllowedOrigins []string

	// AllowedMethods is a list of methods the client is allowed to use
	AllowedMethods []string

	// AllowedHeaders is a list of headers the client is allowed to use
	AllowedHeaders []string

	// ExposedHeaders is a list of headers that are safe to expose to the API
	ExposedHeaders []string

	// AllowCredentials indicates whether the request can include user credentials
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached
	MaxAge int
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Authorization",
			"X-Request-ID",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
			"Retry-After",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// ProductionCORSConfig returns CORS configuration suitable for production
func ProductionCORSConfig(allowedOrigins []string) *CORSConfig {
	config := DefaultCORSConfig()
	config.AllowedOrigins = allowedOrigins
	return config
}

// CORSMiddleware creates CORS middleware with the given configuration
func CORSMiddleware(config *CORSConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed
		allowedOrigin := ""
		if origin != "" {
			allowedOrigin = getAllowedOrigin(origin, config.AllowedOrigins)
		}

		// Set CORS headers
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)

			if config.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			// Expose headers
			if len(config.ExposedHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
			}
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			// Set preflight response headers
			if len(config.AllowedMethods) > 0 {
				c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			}

			if len(config.AllowedHeaders) > 0 {
				c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
			}

			if config.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
			}

			// Respond to preflight request
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// getAllowedOrigin checks if the origin is allowed and returns the appropriate value
func getAllowedOrigin(origin string, allowedOrigins []string) string {
	for _, allowed := range allowedOrigins {
		// Wildcard allows all origins
		if allowed == "*" {
			return origin
		}

		// Exact match
		if allowed == origin {
			return origin
		}

		// Pattern matching (e.g., "*.example.com")
		if strings.HasPrefix(allowed, "*.") {
			suffix := allowed[1:] // Remove the "*"
			if strings.HasSuffix(origin, suffix) {
				return origin
			}
		}
	}

	return ""
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	return getAllowedOrigin(origin, allowedOrigins) != ""
}
