package interfaces

import (
	"context"

	"manual-test-project/internal/domain/user"
)

// AuthService defines the interface for authentication operations.
// Implemented by internal/domain/user/Service.
type AuthService interface {
	Register(ctx context.Context, email, password string) (*user.User, error)
	Authenticate(ctx context.Context, email, password string) (*user.AuthResponse, error)
	RefreshToken(ctx context.Context, oldToken string) (*user.AuthResponse, error)
}

// UserService defines the business logic operations for user management.
// Implemented by internal/domain/user/Service.
type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*user.User, error)
}

// TokenService defines the interface for token generation.
// Implemented by pkg/auth/JWTService.
type TokenService interface {
	GenerateTokens(userID uint) (accessToken string, refreshToken string, expiresIn int64, err error)
}
