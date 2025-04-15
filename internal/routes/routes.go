package routes

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/middleware"
)

const requestTimeout = 10 * time.Second // Set a timeout for HTTP requests to 10 seconds

// RouteGroup represents a group of routes with optional authentication.
// It stores the route group, configuration, and a logger for use within the routes.
type RouteGroup struct {
	Router *gin.RouterGroup
	Config *config.Config
	Logger *zap.Logger
}

// RegisterRoutes sets up all API routes for the gateway
func RegisterRoutes(router *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// Create API version group. This groups all routes under the /api/v1 prefix.
	apiV1 := router.Group("/api/v1")

	// Register a health-check endpoint that can be used to check if the API gateway is running.
	router.GET("/health", createHealthHandler(cfg, logger))

	// Register service routes
	registerAuthRoutes(apiV1.Group("/auth"), cfg, logger)
	registerUserRoutes(apiV1.Group("/users"), cfg, logger)
	registerAdminRoutes(apiV1.Group("/admin"), cfg, logger)
}

// newPublicGroup creates a route group without authentication
func newPublicGroup(router *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) *RouteGroup {
	return &RouteGroup{
		Router: router,
		Config: cfg,
		Logger: logger,
	}
}

// newProtectedGroup creates a route group with authentication
func newProtectedGroup(router *gin.RouterGroup, cfg *config.Config, logger *zap.Logger, roles []string) *RouteGroup {
	// Apply JWT middleware to this group
	router.Use(middleware.JWTAuthMiddleware(cfg, logger))

	// If roles are specified, apply role middleware
	if len(roles) > 0 {
		router.Use(middleware.RoleAuthMiddleware(roles, logger))
	}

	return &RouteGroup{
		Router: router,
		Config: cfg,
		Logger: logger,
	}
}

// registerAuthRoutes sets up all authentication-related routes
func registerAuthRoutes(router *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	authGroup := newPublicGroup(router, cfg, logger)

	// Health check endpoint
	authGroup.Router.GET("/health", createProxyHandler(cfg.Services.AuthServiceURL+"/health", http.MethodGet, logger))

	// Auth endpoints
	authGroup.Router.POST("/register", createProxyHandler(cfg.Services.AuthServiceURL+"/auth/register", http.MethodPost, logger))
	authGroup.Router.POST("/verify-email", createProxyHandler(cfg.Services.AuthServiceURL+"/auth/verify-email", http.MethodPost, logger))
	authGroup.Router.POST("/login", createProxyHandler(cfg.Services.AuthServiceURL+"/auth/login", http.MethodPost, logger))
}

// registerUserRoutes sets up all user-related routes
func registerUserRoutes(router *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	// Public user routes
	publicGroup := newPublicGroup(router, cfg, logger)
	publicGroup.Router.GET("/health", createProxyHandler(cfg.Services.UserServiceURL+"/health", http.MethodGet, logger))

	// Protected user routes
	protectedGroup := newProtectedGroup(router, cfg, logger, []string{constants.RoleUser})
	protectedGroup.Router.POST("/profile", createProxyHandler(cfg.Services.UserServiceURL+"/api/v1/user/profile", http.MethodPost, logger))
	protectedGroup.Router.GET("/profile", createProxyHandler(cfg.Services.UserServiceURL+"/api/v1/user/profile", http.MethodGet, logger))
}

// registerAdminRoutes sets up all admin-related routes
func registerAdminRoutes(router *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	// Public admin routes
	publicGroup := newPublicGroup(router, cfg, logger)
	publicGroup.Router.GET("/health", createProxyHandler(cfg.Services.AdminServiceURL+"/health", http.MethodGet, logger))

	// Protected admin routes
	// protectedGroup := newProtectedGroup(router, cfg, logger, []string{constants.RoleAdmin})
	// Add protected admin routes as needed
}

// createProxyHandler creates a handler function that forwards requests to a service
// It acts as a reverse proxy, handling requests and responses between the client and the service.
func createProxyHandler(serviceURL string, method string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new HTTP client with a timeout to avoid hanging requests.
		client := &http.Client{Timeout: requestTimeout}

		// Variables for the new request and any potential error.
		var req *http.Request
		var err error

		// Handle the request based on its HTTP method.
		// GET and DELETE requests don't have a body.
		if method == http.MethodGet || method == http.MethodDelete {
			// Create a new HTTP request with the given method (GET/DELETE),
			// the target service URL, and a `nil` body because these methods don't require a body.
			req, err = http.NewRequest(method, serviceURL, nil)
		} else {
			// For methods like POST, PUT, or PATCH, the client usually sends data in the request body.
			// Example: JSON data for registering a user or submitting a form.

			// First, we read the body of the incoming request using `io.ReadAll`.
			// This reads all the data sent by the client (like JSON input).
			bodyBytes, err := io.ReadAll(c.Request.Body)

			if err != nil {
				// If reading the body fails (due to a malformed request or read error),
				// we log the error and send a 400 Bad Request response back to the client.
				logger.Error("Failed to read request body", zap.Error(err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			// Now we create a new HTTP request to the backend service using the same method (POST/PUT/etc),
			// the service URL, and attach the read request body using `bytes.NewBuffer(bodyBytes)`.
			// This forwards the client's original data to the destination microservice.
			req, err = http.NewRequest(method, serviceURL, bytes.NewBuffer(bodyBytes))
		}

		// If there was an error creating the request, log it and return an internal server error.
		if err != nil {
			logger.Error("Failed to create request", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Copy all headers from the original request to the new request
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// If there is user information available from previous authentication middleware,
		// add it as headers in the outgoing request for tracing purposes.
		if user, exists := c.Get("user"); exists {
			// Assert the user information to the expected type.
			if userClaims, ok := user.(*middleware.UserClaims); ok {
				req.Header.Set(constants.HeaderUserID, userClaims.UserID)
				if userClaims.Email != "" {
					req.Header.Set(constants.HeaderUsername, userClaims.Email)
				}
				// If the user has any roles, the first role is added as a header.
				if len(userClaims.Roles) > 0 {
					req.Header.Set(constants.HeaderUserRole, userClaims.Roles[0])
				}
			}
		}

		// Send the request to the target service using the HTTP client.
		resp, err := client.Do(req)

		if err != nil {
			// Log the error and return a service unavailable response if the request fails.
			logger.Error("Service request failed", zap.String("url", serviceURL), zap.Error(err))
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
			return
		}
		// Ensure the response body is closed after reading to free resources.
		defer resp.Body.Close()

		// Set the status code of the response to be the same as the service's response.
		c.Status(resp.StatusCode)

		// Copy headers from the response received from the service to the client response.
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		// Copy the body of the service response directly to the client response.
		io.Copy(c.Writer, resp.Body)
	}
}

// createHealthHandler creates a simple health check endpoint that checks the status of all services.
// It makes GET requests to the health endpoint of each service and aggregates the results.
func createHealthHandler(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create an HTTP client with a 5-second timeout for the health checks.
		client := &http.Client{Timeout: 5 * time.Second}

		// Define a map to hold the status of each service
		services := map[string]string{
			"api-gateway":   "up",      // The API gateway itself is assumed to be up.
			"auth-service":  "unknown", // Initial state is unknown for auth service.
			"user-service":  "unknown",
			"admin-service": "unknown",
		}

		// Check the auth service health endpoint. Update the status based on the response.
		if _, err := client.Get(cfg.Services.AuthServiceURL + "/health"); err == nil {
			services["auth-service"] = "up"
		} else {
			services["auth-service"] = "down"
			// Log a warning if the auth service is not reachable.
			logger.Warn("Auth service is down", zap.Error(err))
		}

		// Check the user service health endpoint.
		if _, err := client.Get(cfg.Services.UserServiceURL + "/health"); err == nil {
			services["user-service"] = "up"
		} else {
			services["user-service"] = "down"
			logger.Warn("User service is down", zap.Error(err))
		}

		// Check the admin service health endpoint.
		if _, err := client.Get(cfg.Services.AdminServiceURL + "/health"); err == nil {
			services["admin-service"] = "up"
		} else {
			services["admin-service"] = "down"
			logger.Warn("Admin service is down", zap.Error(err))
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"services":  services,
				"timestamp": time.Now(),
			},
		})
	}
}
