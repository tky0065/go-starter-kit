package interfaces

import (
	"context"
	"time"

	"manual-test-project/internal/domain/user"
)

// UserRepository defines the interface for user data persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, user *user.User) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	SaveRefreshToken(ctx context.Context, UserID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*user.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uint) error
	RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *user.RefreshToken) error
}
