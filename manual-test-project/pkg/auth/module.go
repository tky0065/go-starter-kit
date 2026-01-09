package auth

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/fx"
)

// Module provides the JWT service via fx
var Module = fx.Module("auth",
	fx.Provide(NewJWTServiceFromEnv),
)

// NewJWTServiceFromEnv creates a new JWTService from environment variables
func NewJWTServiceFromEnv() (*JWTService, error) {
	secret := getEnv("JWT_SECRET", "")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	accessExpiryStr := getEnv("JWT_ACCESS_EXPIRY", "15m")
	accessExpiry, err := time.ParseDuration(accessExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRY: %w", err)
	}

	refreshExpiryStr := getEnv("JWT_REFRESH_EXPIRY", "168h")
	refreshExpiry, err := time.ParseDuration(refreshExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRY: %w", err)
	}

	return NewJWTService(secret, accessExpiry, refreshExpiry), nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
