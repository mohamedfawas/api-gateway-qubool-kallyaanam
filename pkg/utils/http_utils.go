package utils

import (
	"errors"
	"strings"
)

// ExtractBearerToken gets the token from an Authorization header
func ExtractBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}
	scheme := parts[0]
	if !strings.EqualFold(scheme, "Bearer") {
		return "", errors.New("invalid authorization scheme")
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token is empty")
	}
	return token, nil
}

// TrimPath removes leading slash if present
func TrimPath(path string) string {
	return strings.TrimPrefix(path, "/")
}
