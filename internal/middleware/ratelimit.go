// Package middleware provides HTTP middleware components for the API Gateway
package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/errors"
)

// RateLimiter interface defines the rate limit behavior
type RateLimiter interface {
	Allow(key string) bool
	CleanUp()
}

// MemoryRateLimiter implements RateLimiter using memory storage
type MemoryRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	limit    rate.Limit
	burst    int
	interval time.Duration
}

// NewMemoryRateLimiter creates a new memory-based rate limiter
func NewMemoryRateLimiter(rps int, burst int, window time.Duration) *MemoryRateLimiter {
	return &MemoryRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		limit:    rate.Limit(float64(rps) / window.Seconds()),
		burst:    burst,
		interval: window,
	}
}

// Allow checks if a request is allowed based on the key
func (l *MemoryRateLimiter) Allow(key string) bool {
	l.mu.RLock()
	limiter, exists := l.limiters[key]
	l.mu.RUnlock()

	if !exists {
		l.mu.Lock()
		// Double check if limiter was created while we were waiting for the write lock
		limiter, exists = l.limiters[key]
		if !exists {
			limiter = rate.NewLimiter(l.limit, l.burst)
			l.limiters[key] = limiter
		}
		l.mu.Unlock()
	}

	return limiter.Allow()
}

// CleanUp removes stale limiters
func (l *MemoryRateLimiter) CleanUp() {
	// Implementation would remove old limiters
}

// RateLimiterMiddleware creates middleware for rate limiting requests
func RateLimiterMiddleware(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	if !cfg.RateLimiting.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Create limiter based on configuration
	var limiter RateLimiter

	switch cfg.RateLimiting.StoreType {
	case "memory":
		limiter = NewMemoryRateLimiter(
			cfg.RateLimiting.Limit,
			cfg.RateLimiting.Burst,
			cfg.RateLimiting.Window,
		)
	// Support for Redis could be added here when needed
	default:
		limiter = NewMemoryRateLimiter(
			cfg.RateLimiting.Limit,
			cfg.RateLimiting.Burst,
			cfg.RateLimiting.Window,
		)
	}

	logger.Info("Rate limiting middleware initialized",
		zap.Bool("enabled", cfg.RateLimiting.Enabled),
		zap.Int("limit", cfg.RateLimiting.Limit),
		zap.Int("burst", cfg.RateLimiting.Burst),
		zap.Duration("window", cfg.RateLimiting.Window),
		zap.String("storeType", cfg.RateLimiting.StoreType),
	)

	return func(c *gin.Context) {
		// Use client IP as the rate limit key
		clientIP := c.ClientIP()

		// Check rate limit
		if !limiter.Allow(clientIP) {
			logger.Warn("Rate limit exceeded",
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
