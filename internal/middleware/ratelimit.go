// Package middleware provides HTTP middleware components for the API Gateway
package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/errors"
)

// RateLimiterMiddleware creates a Gin middleware for rate limiting requests using Redis
func RateLimiterMiddleware(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	// If rate limiting is disabled, return a no-op middleware
	if !cfg.RateLimiting.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RateLimiting.RedisAddress,
	})

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Warn("Redis connection failed, rate limiting may not work properly",
			zap.Error(err),
			zap.String("redisAddress", cfg.RateLimiting.RedisAddress),
		)
	}

	// Create the Redis-based rate limiter
	limiter := redis_rate.NewLimiter(redisClient)

	// Calculate requests per second from config values
	rps := int(float64(cfg.RateLimiting.Limit) / cfg.RateLimiting.Window.Seconds())

	logger.Info("Redis rate limiter initialized",
		zap.Int("requestsPerSecond", rps),
		zap.Int("limit", cfg.RateLimiting.Limit),
		zap.Duration("window", cfg.RateLimiting.Window),
	)

	// Return the actual middleware function
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "rl:" + clientIP // Prefix helps organize Redis keys

		// Check if request is allowed
		res, err := limiter.Allow(c.Request.Context(), key, redis_rate.PerSecond(rps))
		if err != nil {
			// Redis error - allow the request but log the error
			logger.Error("Rate limiter Redis error",
				zap.Error(err),
				zap.String("clientIP", clientIP),
			)
			c.Next()
			return
		}

		// Add rate limit headers to response
		c.Header("X-RateLimit-Limit", "1")
		c.Header("X-RateLimit-Remaining", "1")
		c.Header("X-RateLimit-Reset", "1")

		if res.Allowed <= 0 {
			logger.Info("Rate limit exceeded",
				zap.String("clientIP", clientIP),
				zap.String("path", c.Request.URL.Path),
			)

			apiErr := errors.New(errors.ErrorTypeRateLimited, "Rate limit exceeded", nil)
			c.Error(apiErr)
			c.Abort()
			return
		}

		c.Next()
	}
}
