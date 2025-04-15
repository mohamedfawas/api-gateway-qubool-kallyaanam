// Package middleware provides HTTP middleware components for the API Gateway
package middleware

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/errors"
)

// RateLimiterMiddleware creates a Gin middleware for rate limiting requests using Redis
func RateLimiterMiddleware(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	// Check if rate limiting is enabled in the config file
	// If it's not enabled, just return a dummy middleware that does nothing
	if !cfg.RateLimiting.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// If rate limiting is enabled, we first create a Redis client to connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RateLimiting.RedisAddress, // Redis server address
	})

	// Use a context for Redis operations (required by Redis client)
	ctx := context.Background()

	// Try pinging Redis to check if the connection is successful
	if err := redisClient.Ping(ctx).Err(); err != nil {
		// Log the failure and disable rate limiting
		logger.Error("Redis connection failed, rate limiting disabled",
			zap.Error(err),
			zap.String("redisAddress", cfg.RateLimiting.RedisAddress),
		)
		// Return a dummy middleware since Redis is not available
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// If Redis is reachable, create a new rate limiter using the Redis client
	limiter := redis_rate.NewLimiter(redisClient)

	// Calculate how many requests per second are allowed
	// Example: if the config says 60 requests per 1 minute, it becomes 1 request per second
	rps := int(float64(cfg.RateLimiting.Limit) / cfg.RateLimiting.Window.Seconds())

	// Log the rate limiting configuration
	logger.Info("Redis rate limiter initialized",
		zap.Int("requestsPerSecond", rps),
		zap.Int("limit", cfg.RateLimiting.Limit),
		zap.Duration("window", cfg.RateLimiting.Window),
	)

	// Return the actual middleware function
	return func(c *gin.Context) {
		clientIP := c.ClientIP() // Get the IP address of the client
		key := "rl:" + clientIP  // Generate a Redis key using the IP (e.g., "rl:192.168.1.1")

		// Ask the limiter if this client/IP can make a request now
		res, err := limiter.Allow(c.Request.Context(), key, redis_rate.PerSecond(rps))

		if err != nil {
			// If Redis has an error, just allow the request but log it
			logger.Error("Rate limiter Redis error",
				zap.Error(err),
				zap.String("clientIP", clientIP),
			)
			c.Next()
			return
		}

		// Add rate limit headers to response with proper values
		// These headers help clients understand how many requests they have left and when the limit will reset
		c.Header("X-RateLimit-Limit", strconv.Itoa(rps))                                    // Total allowed requests per second
		c.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))                      // How many requests are left
		c.Header("X-RateLimit-Reset", strconv.FormatInt(res.ResetAfter.Milliseconds(), 10)) // Time left to reset limit

		// If not allowed (i.e., rate limit exceeded), block the request
		if res.Allowed <= 0 {
			// Log that the client has exceeded the limit
			logger.Warn("Rate limit exceeded",
				zap.String("clientIP", clientIP),
				zap.String("path", c.Request.URL.Path),
			)

			// Create a rate-limited error using the custom error package
			apiErr := errors.RateLimitedError("Rate limit exceeded")
			c.Error(apiErr) // Attach the error to the Gin context
			c.Abort()       // Stop further processing of the request
			return
		}

		// If allowed, continue to the next handler in the middleware chain
		c.Next()
	}
}
