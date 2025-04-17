package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
	apiErrors "github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/errors"
)

// UserClaims represents the custom claims for JWT tokens
type UserClaims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTAuthMiddleware creates a middleware for JWT authentication
func JWTAuthMiddleware(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if we have user info from gateway in headers
		userID := c.GetHeader(constants.HeaderUserID)
		userRole := c.GetHeader(constants.HeaderUserRole)

		// If headers are present and we trust the API gateway, we can use these values
		if userID != "" && userRole != "" {
			// Store the user info in the context
			logger.Debug("User authenticated via gateway headers",
				zap.String("user_id", userID),
				zap.String("role", userRole))

			// Create simplified claims object from headers
			c.Set("user", &UserClaims{
				UserID: userID,
				Roles:  []string{userRole},
			})
			c.Next()
			return
		}
		// Otherwise, fall back to direct JWT validation
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Debug("Missing authorization header or gateway headers")
			c.JSON(401, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Debug("Invalid authorization header format")
			c.JSON(401, gin.H{"error": "Invalid authentication format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse the token with improved validation
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing algorithm (explicit check for security)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Check the 'kid' header if present (for future key rotation)
			if kid, ok := token.Header["kid"].(string); ok {
				logger.Debug("Token key ID", zap.String("kid", kid))
				// We currently use just one key, but this can be extended for key rotation
			}

			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			logger.Debug("Failed to parse token", zap.Error(err))
			c.Error(apiErrors.New(apiErrors.ErrorTypeUnauthorized, "Invalid token", err))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			// Check if the token is expired - using the new JWT v5 approach
			expirationTime, err := claims.GetExpirationTime()
			if err != nil || expirationTime == nil || expirationTime.Before(time.Now()) {
				logger.Debug("Token expired")
				c.Error(apiErrors.New(apiErrors.ErrorTypeUnauthorized, "Token expired", nil))
				c.Abort()
				return
			}

			// Store the claims in the context for later use
			c.Set("user", claims)
			logger.Debug("Authenticated user",
				zap.String("user_id", claims.UserID),
				zap.String("email", claims.Email),
				zap.Strings("roles", claims.Roles))
		} else {
			logger.Debug("Invalid token claims")
			c.Error(apiErrors.New(apiErrors.ErrorTypeUnauthorized, "Invalid token claims", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleAuthMiddleware creates a middleware to check user roles
func RoleAuthMiddleware(requiredRoles []string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user claims from context
		userValue, exists := c.Get("user")
		if !exists {
			logger.Debug("User claims not found in context")
			c.Error(apiErrors.New(apiErrors.ErrorTypeUnauthorized, "User claims not found", nil))
			c.Abort()
			return
		}

		claims, ok := userValue.(*UserClaims)
		if !ok {
			logger.Debug("Invalid user claims type")
			c.Error(apiErrors.New(apiErrors.ErrorTypeUnauthorized, "Invalid user claims", nil))
			c.Abort()
			return
		}

		// If no roles are required, just continue
		if len(requiredRoles) == 0 {
			c.Next()
			return
		}

		// Check if the user has any of the required roles
		hasRole := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range claims.Roles {
				if requiredRole == userRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			logger.Debug("User does not have required role",
				zap.String("user_id", claims.UserID),
				zap.Strings("required_roles", requiredRoles))
			c.Error(apiErrors.New(apiErrors.ErrorTypeForbidden, "Insufficient permissions", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
