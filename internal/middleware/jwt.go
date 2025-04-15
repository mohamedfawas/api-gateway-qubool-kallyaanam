// Package middleware provides HTTP middleware components for the API Gateway
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
		// Get the Authorization header
		authHeader := c.GetHeader(constants.HeaderAuthorization)
		if authHeader == "" {
			logger.Debug("Missing authorization header")
			c.Error(apiErrors.New(apiErrors.ErrorTypeUnauthorized, "Missing authorization header", nil))
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Debug("Invalid authorization header format")
			c.Error(apiErrors.New(apiErrors.ErrorTypeUnauthorized, "Invalid authorization header format", nil))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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
