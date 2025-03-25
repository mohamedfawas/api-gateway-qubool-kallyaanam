package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer"
	// Use a unique context key to avoid collisions
	userIDContextKey   = "userID"
	userIDHeaderKey    = "X-User-ID"
	errorAuthRequired  = "Authorization required"
	errorInvalidFormat = "Invalid authorization format"
	errorInvalidToken  = "Invalid or expired token"
	errorMissingUserID = "User ID not found in token"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func extractToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}
	scheme := parts[0]
	if !strings.EqualFold(scheme, bearerPrefix) {
		return "", errors.New("invalid authorization scheme")
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token is empty")
	}
	return token, nil
}

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, errorAuthRequired, nil))
			return
		}

		tokenString, err := extractToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, errorInvalidFormat, nil))
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.Auth.JWTSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, errorInvalidToken, nil))
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, errorInvalidToken, nil))
			return
		}

		if claims.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, errorMissingUserID, nil))
			return
		}

		c.Set(userIDContextKey, claims.UserID)
		c.Request.Header.Set(userIDHeaderKey, claims.UserID)
		c.Next()
	}
}
