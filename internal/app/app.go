package app

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/handlers"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/middleware"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/swagger"
)

// App represents the application
type App struct {
	config *config.Config
	router *gin.Engine
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config) *App {
	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &App{
		config: cfg,
		router: gin.Default(),
	}
}

// SetupRoutes configures all application routes
func (a *App) SetupRoutes(proxyHandler *handlers.ProxyHandler) {
	// Apply global middleware
	a.router.Use(middleware.CORS())
	a.router.Use(middleware.SimpleErrorHandler())
	// Add rate limiting - 100 requests per second with burst of 200
	a.router.Use(middleware.RateLimit(100, 200))

	// Setup health check endpoint with Swagger annotation
	// @Summary API health check
	// @Description Returns OK status if API is running
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	a.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Setup Swagger documentation routes
	swagger.SetupSwaggerRoutes(a.router)

	// Auth service routes - no auth required
	authRoutes := a.router.Group("/auth")
	authRoutes.Any("/*path", proxyHandler.ProxyRequest("auth"))

	// Protected routes - require auth
	protectedRoutes := a.router.Group("")
	protectedRoutes.Use(middleware.JWTAuth(a.config))

	// User service routes
	protectedRoutes.Any("/user/*path", proxyHandler.ProxyRequest("user"))

	// Admin service routes
	protectedRoutes.Any("/admin/*path", proxyHandler.ProxyRequest("admin"))
}

// SetupServer initializes the HTTP server
func (a *App) SetupServer() *http.Server {
	// Create proxy handler for services
	proxyHandler, err := handlers.NewProxyHandler(a.config.Services) // creates a new proxy handler instance with the services map from the configuration
	if err != nil {
		log.Fatalf("Failed to initialize proxy handler: %v", err)
	}

	// Setup routes
	a.SetupRoutes(proxyHandler) // sets up the routes for the application using the proxy handler

	// Return configured HTTP server
	return &http.Server{
		Addr:    ":" + a.config.Server.Port,
		Handler: a.router,
	}
}
