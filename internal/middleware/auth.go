package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/pkg/auth"
)

// JWTAuth validates JWT tokens and adds user info to the request context
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(
				http.StatusUnauthorized,
				"Authorization header is required",
				nil,
			))
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(
				http.StatusUnauthorized,
				"Invalid token format",
				nil,
			))
			return
		}

		claims, err := auth.ValidateToken(tokenParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(
				http.StatusUnauthorized,
				"Invalid or expired token",
				err.Error(),
			))
			return
		}

		// Add user info to context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RBACMiddleware restricts access based on user role
func RBACMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, models.NewErrorResponse(
				http.StatusForbidden,
				"Role information not found",
				nil,
			))
			return
		}

		role := userRole.(string)
		allowed := false
		for _, r := range allowedRoles {
			if r == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, models.NewErrorResponse(
				http.StatusForbidden,
				"Insufficient permissions",
				nil,
			))
			return
		}

		c.Next()
	}
}
