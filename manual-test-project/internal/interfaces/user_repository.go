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
	FindByID(ctx context.Context, id uint) (*user.User, error)
	FindAll(ctx context.Context, page, limit int) ([]*user.User, int64, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id uint) error
	SaveRefreshToken(ctx context.Context, UserID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*user.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uint) error
	RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *user.RefreshToken) error
}
