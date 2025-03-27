package swagger

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupSwaggerRoutes configures the routes for Swagger documentation
func SetupSwaggerRoutes(router *gin.Engine) {
	// Don't use Swagger in production unless explicitly enabled
	if os.Getenv("ENV") == "production" && os.Getenv("ENABLE_SWAGGER") != "true" {
		return
	}

	// Log that we're setting up swagger
	log.Println("Setting up Swagger documentation routes")

	// Add a basic redirect from /swagger to /swagger/index.html for convenience
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	// Configure swagger with explicit JSON URL to avoid potential path issues
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"), // The URL should be relative to the server
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	log.Println("Swagger documentation available at /swagger/index.html")
}
