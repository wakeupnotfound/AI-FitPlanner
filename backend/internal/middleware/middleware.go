// Package middleware provides HTTP middleware components for the AI Fitness Planner API.
//
// This package includes:
//   - Authentication middleware (JWT validation and session verification)
//   - Rate limiting middleware (token bucket algorithm with Redis)
//   - Security middleware (input sanitization, XSS/SQL injection prevention, security headers)
//   - Logging middleware (request/response logging with sensitive data masking)
//   - CORS middleware (cross-origin resource sharing configuration)
//   - Recovery middleware (panic recovery with stack trace logging)
//
// Usage example:
//
//	router := gin.New()
//
//	// Add middleware stack
//	router.Use(middleware.RecoveryMiddleware(nil))
//	router.Use(middleware.LoggingMiddleware(nil))
//	router.Use(middleware.CORSMiddleware(nil))
//	router.Use(middleware.SecurityMiddleware(nil))
//
//	// Protected routes
//	protected := router.Group("/api/v1")
//	protected.Use(middleware.AuthMiddleware(jwtManager, sessionManager))
//	protected.Use(rateLimiter.RateLimitMiddleware())
package middleware
