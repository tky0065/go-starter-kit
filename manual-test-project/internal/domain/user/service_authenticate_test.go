package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"manual-test-project/internal/domain"
)

// mockJWTService is a mock implementation of JWTService for testing
type mockJWTService struct {
	generateTokensFunc func(userID uint) (string, string, int64, error)
}

func (m *mockJWTService) GenerateTokens(userID uint) (string, string, int64, error) {
	if m.generateTokensFunc != nil {
		return m.generateTokensFunc(userID)
	}
	return "access-token", "refresh-token", 900, nil
}

// mockRepository is extended to support refresh token operations
type mockRepositoryForAuth struct {
	getUserByEmailFunc   func(ctx context.Context, email string) (*User, error)
	saveRefreshTokenFunc func(ctx context.Context, userID uint, token string, expiresAt time.Time) error
}

func (m *mockRepositoryForAuth) CreateUser(ctx context.Context, user *User) error {
	return nil
}

func (m *mockRepositoryForAuth) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *mockRepositoryForAuth) FindByID(ctx context.Context, id uint) (*User, error) {
	// Mock implementation for FindByID
	return nil, nil
}

func (m *mockRepositoryForAuth) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	if m.saveRefreshTokenFunc != nil {
		return m.saveRefreshTokenFunc(ctx, userID, token, expiresAt)
	}
	return nil
}

func (m *mockRepositoryForAuth) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	return nil, nil
}

func (m *mockRepositoryForAuth) RevokeRefreshToken(ctx context.Context, tokenID uint) error {
	return nil
}

func (m *mockRepositoryForAuth) RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *RefreshToken) error {
	return nil
}

func (m *mockRepositoryForAuth) FindAll(ctx context.Context, page, limit int) ([]*User, int64, error) {
	return nil, 0, nil
}

func (m *mockRepositoryForAuth) Update(ctx context.Context, user *User) error {
	return nil
}

func (m *mockRepositoryForAuth) Delete(ctx context.Context, id uint) error {
	return nil
}

func TestService_Authenticate_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo := &mockRepositoryForAuth{
		getUserByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return &User{
				ID:           1,
				Email:        email,
				PasswordHash: string(hashedPassword),
			}, nil
		},
		saveRefreshTokenFunc: func(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
			return nil
		},
	}

	mockJWT := &mockJWTService{
		generateTokensFunc: func(userID uint) (string, string, int64, error) {
			return "access-token-123", "refresh-token-456", 900, nil
		},
	}

	service := NewServiceWithJWT(mockRepo, mockJWT)

	response, err := service.Authenticate(context.Background(), "test@example.com", "password123")

	if err != nil {
		t.Fatalf("Authenticate() unexpected error: %v", err)
	}

	if response == nil {
		t.Fatal("Authenticate() response is nil")
	}

	if response.AccessToken != "access-token-123" {
		t.Errorf("AccessToken = %v, want %v", response.AccessToken, "access-token-123")
	}

	if response.RefreshToken != "refresh-token-456" {
		t.Errorf("RefreshToken = %v, want %v", response.RefreshToken, "refresh-token-456")
	}

	if response.ExpiresIn != 900 {
		t.Errorf("ExpiresIn = %v, want %v", response.ExpiresIn, 900)
	}
}

func TestService_Authenticate_UserNotFound(t *testing.T) {
	mockRepo := &mockRepositoryForAuth{
		getUserByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return nil, nil // User not found
		},
	}

	mockJWT := &mockJWTService{}
	service := NewServiceWithJWT(mockRepo, mockJWT)

	_, err := service.Authenticate(context.Background(), "nonexistent@example.com", "password123")

	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Authenticate() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestService_Authenticate_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockRepo := &mockRepositoryForAuth{
		getUserByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return &User{
				ID:           1,
				Email:        email,
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	mockJWT := &mockJWTService{}
	service := NewServiceWithJWT(mockRepo, mockJWT)

	_, err := service.Authenticate(context.Background(), "test@example.com", "wrongpassword")

	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Authenticate() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestService_Authenticate_RepositoryError(t *testing.T) {
	mockRepo := &mockRepositoryForAuth{
		getUserByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return nil, errors.New("database connection failed")
		},
	}

	mockJWT := &mockJWTService{}
	service := NewServiceWithJWT(mockRepo, mockJWT)

	_, err := service.Authenticate(context.Background(), "test@example.com", "password123")

	if err == nil {
		t.Error("Authenticate() expected error but got nil")
	}

	if errors.Is(err, domain.ErrInvalidCredentials) {
		t.Error("Authenticate() should not return ErrInvalidCredentials for infrastructure errors")
	}
}

func TestService_Authenticate_JWTGenerationError(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo := &mockRepositoryForAuth{
		getUserByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return &User{
				ID:           1,
				Email:        email,
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	mockJWT := &mockJWTService{
		generateTokensFunc: func(userID uint) (string, string, int64, error) {
			return "", "", 0, errors.New("JWT generation failed")
		},
	}

	service := NewServiceWithJWT(mockRepo, mockJWT)

	_, err := service.Authenticate(context.Background(), "test@example.com", "password123")

	if err == nil {
		t.Error("Authenticate() expected error but got nil")
	}
}

func TestService_Authenticate_SaveRefreshTokenError(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo := &mockRepositoryForAuth{
		getUserByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return &User{
				ID:           1,
				Email:        email,
				PasswordHash: string(hashedPassword),
			}, nil
		},
		saveRefreshTokenFunc: func(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
			return errors.New("failed to save refresh token")
		},
	}

	mockJWT := &mockJWTService{}
	service := NewServiceWithJWT(mockRepo, mockJWT)

	_, err := service.Authenticate(context.Background(), "test@example.com", "password123")

	if err == nil {
		t.Error("Authenticate() expected error but got nil")
	}
}
