package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"manual-test-project/internal/domain"
	"manual-test-project/internal/domain/user"
)

// UserRepository implements user data persistence using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// GetUserByEmail retrieves a user by email, returns nil if not found
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByID retrieves a user by ID, returns nil if not found
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// SaveRefreshToken saves a refresh token for the given user
func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	refreshToken := &user.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}
	return r.db.WithContext(ctx).Create(refreshToken).Error
}

// GetRefreshToken retrieves a refresh token by token string
func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (*user.RefreshToken, error) {
	var rt user.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func (r *UserRepository) RevokeRefreshToken(ctx context.Context, tokenID uint) error {
	return r.db.WithContext(ctx).Model(&user.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked", true).Error
}

// RotateRefreshToken performs atomic token rotation: revocation of old and creation of new
func (r *UserRepository) RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *user.RefreshToken) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Revoke old token with optimistic locking check
		// We only update if revoked is currently false
		result := tx.Model(&user.RefreshToken{}).
			Where("id = ? AND revoked = ?", oldTokenID, false).
			Update("revoked", true)

		if result.Error != nil {
			return result.Error
		}

		// If no rows affected, it means the token was already revoked (race condition)
		// or doesn't exist
		if result.RowsAffected == 0 {
			return domain.ErrRefreshTokenRevoked
		}

		// 2. Create new token
		if err := tx.Create(newToken).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindAll retrieves users with pagination
func (r *UserRepository) FindAll(ctx context.Context, page, limit int) ([]*user.User, int64, error) {
	var users []*user.User
	var total int64

	// Get total count (ignoring soft deleted)
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated users
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

// Delete performs a soft delete on the user (sets deleted_at)
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, id).Error
}
