package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiting for different clients
type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	// Rate limit configuration
	rate     rate.Limit
	capacity int
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(r rate.Limit, capacity int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     r,
		capacity: capacity,
	}
}

// getLimiter returns the rate limiter for the provided client IP
func (rl *RateLimiter) getLimiter(clientIP string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[clientIP]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.capacity)
		rl.visitors[clientIP] = limiter
	}

	return limiter
}

// RateLimit middleware function
func RateLimit(requestsPerSecond float64, burst int) gin.HandlerFunc {
	rateLimiter := NewRateLimiter(rate.Limit(requestsPerSecond), burst)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := rateLimiter.getLimiter(clientIP)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				models.NewErrorResponse(
					http.StatusTooManyRequests,
					"Rate limit exceeded",
					"Too many requests, please try again later",
				),
			)
			return
		}

		c.Next()
	}
}

// CleanupVisitors periodically removes old visitors from the map
func (rl *RateLimiter) CleanupVisitors() {
	for {
		time.Sleep(time.Hour) // Run cleanup every hour
		rl.mu.Lock()
		for ip := range rl.visitors {
			delete(rl.visitors, ip)
		}
		rl.mu.Unlock()
	}
}
