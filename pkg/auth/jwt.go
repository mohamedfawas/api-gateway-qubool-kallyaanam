package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// Claims represents the JWT claims used in the application
type Claims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ValidateToken validates a JWT token and returns the claims if valid
func ValidateToken(tokenString string) (*Claims, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, errors.New("failed to load configuration")
	}

	claims := &Claims{}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		// Return the secret key used for signing
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateToken creates a new JWT token for testing purposes
// In production, tokens should be generated by the auth service
func GenerateToken(userID, role string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", errors.New("failed to load configuration")
	}

	// Create expiration time
	expirationTime := time.Now().Add(time.Duration(cfg.JWT.Expiration) * time.Minute)

	// Create claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "matrimonial-api-gateway",
			Subject:   userID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
