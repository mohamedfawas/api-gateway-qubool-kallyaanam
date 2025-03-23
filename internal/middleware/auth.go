package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// JWTAuth validates the JWT token in the Authorization header
func JWTAuth(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(
				http.StatusUnauthorized,
				"Authorization required",
				nil,
			))
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(
				http.StatusUnauthorized,
				"Invalid authorization format",
				nil,
			))
			return
		}

		// Parse and validate token using v5 of the jwt package
		token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Auth.JWTSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(
				http.StatusUnauthorized,
				"Invalid or expired token",
				nil,
			))
			return
		}

		// Extract user details from token
		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(
				http.StatusUnauthorized,
				"Invalid or expired token",
				nil,
			))
			return
		}

		// Add user ID to context for downstream services
		c.Set("userID", claims.UserID)
		// Add user ID to headers for microservices
		c.Request.Header.Set("X-User-ID", claims.UserID)

		c.Next()
	}
}
