package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"manual-test-project/internal/domain"
	"manual-test-project/internal/domain/user"
)

type mockRepo struct {
	users map[string]*user.User
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		users: make(map[string]*user.User),
	}
}

func (m *mockRepo) CreateUser(ctx context.Context, u *user.User) error {
	if _, exists := m.users[u.Email]; exists {
		return errors.New("conflict")
	}
	u.ID = uint(len(m.users) + 1)
	m.users[u.Email] = u
	return nil
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	if u, exists := m.users[email]; exists {
		return u, nil
	}
	return nil, nil // Return nil, nil if not found for this mock
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*user.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	return nil
}

func (m *mockRepo) GetRefreshToken(ctx context.Context, token string) (*user.RefreshToken, error) {
	return nil, nil
}

func (m *mockRepo) RevokeRefreshToken(ctx context.Context, tokenID uint) error {
	return nil
}

func (m *mockRepo) RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *user.RefreshToken) error {
	return nil
}

func TestService_Register(t *testing.T) {
	repo := newMockRepo()
	svc := user.NewService(repo)
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"

	// Test Success
	createdUser, err := svc.Register(ctx, email, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdUser.Email != email {
		t.Errorf("expected email %s, got %s", email, createdUser.Email)
	}

	if createdUser.PasswordHash == "" {
		t.Error("expected password hash to be set")
	}

	if createdUser.PasswordHash == password {
		t.Error("expected password to be hashed, got plain text")
	}

	// Verify hash
	err = bcrypt.CompareHashAndPassword([]byte(createdUser.PasswordHash), []byte(password))
	if err != nil {
		t.Errorf("password hash verification failed: %v", err)
	}

	// Test Duplicate
	_, err = svc.Register(ctx, email, "newpassword")
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
	if !errors.Is(err, domain.ErrEmailAlreadyRegistered) {
		t.Errorf("expected ErrEmailAlreadyRegistered, got %v", err)
	}
}
