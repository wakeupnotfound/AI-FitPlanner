package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	// Per-user limits
	UserRequestsPerMinute int64
	UserRequestsPerHour   int64

	// Per-IP limits
	IPRequestsPerMinute int64

	// AI generation endpoint limits (stricter)
	AIGenerationPerMinute int64
}

// DefaultRateLimitConfig returns default rate limit configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		UserRequestsPerMinute: 60,
		UserRequestsPerHour:   1000,
		IPRequestsPerMinute:   100,
		AIGenerationPerMinute: 2,
	}
}

// RateLimiter handles rate limiting using Redis token bucket algorithm
type RateLimiter struct {
	client *redis.Client
	config *RateLimitConfig
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(client *redis.Client, config *RateLimitConfig) *RateLimiter {
	if config == nil {
		config = DefaultRateLimitConfig()
	}
	return &RateLimiter{
		client: client,
		config: config,
	}
}

// RateLimitMiddleware creates rate limiting middleware for general API endpoints
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Check per-IP rate limit
		clientIP := c.ClientIP()
		ipKey := fmt.Sprintf("ratelimit:ip:%s:minute", clientIP)

		allowed, retryAfter, err := rl.checkRateLimit(ctx, ipKey, rl.config.IPRequestsPerMinute, time.Minute)
		if err != nil {
			logger.Error("IP限流检查失败", zap.Error(err), zap.String("ip", clientIP))
			// Allow request on error to avoid blocking legitimate users
			c.Next()
			return
		}

		if !allowed {
			c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Error(4290, "请求过于频繁，请稍后再试"))
			return
		}

		// Check per-user rate limit if authenticated
		userID, exists := GetUserID(c)
		if exists {
			// Per-minute limit
			userMinuteKey := fmt.Sprintf("ratelimit:user:%d:minute", userID)
			allowed, retryAfter, err = rl.checkRateLimit(ctx, userMinuteKey, rl.config.UserRequestsPerMinute, time.Minute)
			if err != nil {
				logger.Error("用户分钟限流检查失败", zap.Error(err), zap.Int64("user_id", userID))
				c.Next()
				return
			}

			if !allowed {
				c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
				c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Error(4290, "请求过于频繁，请稍后再试"))
				return
			}

			// Per-hour limit
			userHourKey := fmt.Sprintf("ratelimit:user:%d:hour", userID)
			allowed, retryAfter, err = rl.checkRateLimit(ctx, userHourKey, rl.config.UserRequestsPerHour, time.Hour)
			if err != nil {
				logger.Error("用户小时限流检查失败", zap.Error(err), zap.Int64("user_id", userID))
				c.Next()
				return
			}

			if !allowed {
				c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
				c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Error(4290, "请求过于频繁，请稍后再试"))
				return
			}
		}

		c.Next()
	}
}

// AIGenerationRateLimitMiddleware creates stricter rate limiting for AI generation endpoints
func (rl *RateLimiter) AIGenerationRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userID, exists := GetUserID(c)
		if !exists {
			// Should not happen as AI endpoints require authentication
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.UnauthorizedError("需要认证"))
			return
		}

		// Stricter per-minute limit for AI generation
		aiKey := fmt.Sprintf("ratelimit:ai:%d:minute", userID)
		allowed, retryAfter, err := rl.checkRateLimit(ctx, aiKey, rl.config.AIGenerationPerMinute, time.Minute)
		if err != nil {
			logger.Error("AI生成限流检查失败", zap.Error(err), zap.Int64("user_id", userID))
			c.Next()
			return
		}

		if !allowed {
			c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Error(4290, "AI生成请求过于频繁，请稍后再试"))
			return
		}

		c.Next()
	}
}

// checkRateLimit implements token bucket algorithm using Redis
// Returns (allowed, retryAfterSeconds, error)
func (rl *RateLimiter) checkRateLimit(ctx context.Context, key string, limit int64, window time.Duration) (bool, int64, error) {
	// Use Redis INCR with EXPIRE for simple rate limiting
	// This is a sliding window counter approach

	pipe := rl.client.Pipeline()

	// Increment counter
	incrCmd := pipe.Incr(ctx, key)

	// Set expiration only if key is new (NX flag equivalent via checking TTL)
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, fmt.Errorf("failed to execute rate limit pipeline: %w", err)
	}

	count := incrCmd.Val()
	ttl := ttlCmd.Val()

	// If TTL is -1 (no expiration) or -2 (key doesn't exist), set expiration
	if ttl < 0 {
		if err := rl.client.Expire(ctx, key, window).Err(); err != nil {
			logger.Warn("设置限流键过期时间失败", zap.Error(err), zap.String("key", key))
		}
		ttl = window
	}

	if count > limit {
		// Calculate retry-after in seconds
		retryAfter := int64(ttl.Seconds())
		if retryAfter < 1 {
			retryAfter = 1
		}
		return false, retryAfter, nil
	}

	return true, 0, nil
}

// GetRemainingRequests returns the number of remaining requests for a key
func (rl *RateLimiter) GetRemainingRequests(ctx context.Context, key string, limit int64) (int64, error) {
	count, err := rl.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return limit, nil
	}
	if err != nil {
		return 0, err
	}

	remaining := limit - count
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}
