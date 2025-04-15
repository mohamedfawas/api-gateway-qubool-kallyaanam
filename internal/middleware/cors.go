package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// CORSMiddleware creates a middleware for handling CORS
func CORSMiddleware(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	corsConfig := cors.Config{
		// List of allowed origins (domains) that can make requests to this API.
		// For example: http://localhost:3000 or https://your-frontend.com
		AllowOrigins: cfg.CORS.AllowOrigins,

		// List of HTTP methods that are allowed for cross-origin requests.
		// For example: GET, POST, PUT, DELETE
		AllowMethods: cfg.CORS.AllowMethods,

		// List of headers that the client can use when making a request.
		// For example: Content-Type, Authorization
		AllowHeaders: cfg.CORS.AllowHeaders,

		// List of headers that the browser is allowed to access in the response.
		// For example: Content-Length
		ExposeHeaders: cfg.CORS.ExposeHeaders,

		// If true, the server allows credentials like cookies or authorization headers in cross-origin requests.
		AllowCredentials: cfg.CORS.AllowCredentials,

		// How long (in seconds) the results of a preflight request can be cached by the browser.
		// A preflight request is a CORS mechanism to check if the real request is safe to send.
		MaxAge: cfg.CORS.MaxAge,
	}

	// Use the logger to print the current CORS configuration.
	// This is helpful for debugging and verifying the settings.
	logger.Info("Configured CORS middleware",
		zap.Strings("allowOrigins", cfg.CORS.AllowOrigins),
		zap.Strings("allowMethods", cfg.CORS.AllowMethods),
		zap.Strings("allowHeaders", cfg.CORS.AllowHeaders),
		zap.Strings("exposeHeaders", cfg.CORS.ExposeHeaders),
		zap.Bool("allowCredentials", cfg.CORS.AllowCredentials),
		zap.Duration("maxAge", cfg.CORS.MaxAge),
	)

	// Return the actual middleware function
	return cors.New(corsConfig)
}
