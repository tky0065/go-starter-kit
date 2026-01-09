package user_test

import (
	"context"
	"testing"
	"time"

	"manual-test-project/internal/domain"
	"manual-test-project/internal/domain/user"
)

// mockUserRepository is a mock implementation of interfaces.UserRepository for testing
type mockUserRepository struct {
	users         map[string]*user.User
	refreshTokens map[string]*user.RefreshToken
	nextUserID    uint
	nextTokenID   uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:         make(map[string]*user.User),
		refreshTokens: make(map[string]*user.RefreshToken),
		nextUserID:    1,
		nextTokenID:   1,
	}
}

func (m *mockUserRepository) CreateUser(ctx context.Context, u *user.User) error {
	u.ID = m.nextUserID
	m.nextUserID++
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	m.users[u.Email] = u
	return nil
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	u, exists := m.users[email]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	rt := &user.RefreshToken{
		ID:        m.nextTokenID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.nextTokenID++
	m.refreshTokens[token] = rt
	return nil
}

func (m *mockUserRepository) GetRefreshToken(ctx context.Context, token string) (*user.RefreshToken, error) {
	rt, exists := m.refreshTokens[token]
	if !exists {
		return nil, nil
	}
	return rt, nil
}

func (m *mockUserRepository) RevokeRefreshToken(ctx context.Context, tokenID uint) error {
	for _, rt := range m.refreshTokens {
		if rt.ID == tokenID {
			rt.Revoked = true
			return nil
		}
	}
	return nil
}

func (m *mockUserRepository) RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *user.RefreshToken) error {
	// 1. Check/Revoke old token
	var foundOld bool
	for _, rt := range m.refreshTokens {
		if rt.ID == oldTokenID {
			if rt.Revoked {
				return domain.ErrRefreshTokenRevoked // Simulate DB check
			}
			rt.Revoked = true
			foundOld = true
			break
		}
	}
	if !foundOld {
		// In a real DB, if ID not found, rows affected = 0.
		// For mock, we can just say revoked (as per implementation logic if rows=0) or ignore
		return domain.ErrRefreshTokenRevoked
	}

	// 2. Create new token
	newToken.ID = m.nextTokenID
	m.nextTokenID++
	newToken.CreatedAt = time.Now()
	newToken.UpdatedAt = time.Now()
	m.refreshTokens[newToken.Token] = newToken
	return nil
}

func (m *mockUserRepository) FindAll(ctx context.Context, page, limit int) ([]*user.User, int64, error) {
	users := make([]*user.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, int64(len(users)), nil
}

func (m *mockUserRepository) Update(ctx context.Context, u *user.User) error {
	for email, existing := range m.users {
		if existing.ID == u.ID {
			delete(m.users, email)
			m.users[u.Email] = u
			return nil
		}
	}
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uint) error {
	for email, u := range m.users {
		if u.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return nil
}

// mockTokenService is a mock implementation of interfaces.TokenService for testing
type mockTokenService struct {
	accessToken  string
	refreshToken string
	expiresIn    int64
}

func newMockTokenService() *mockTokenService {
	return &mockTokenService{
		accessToken:  "mock-access-token",
		refreshToken: "mock-refresh-token",
		expiresIn:    3600,
	}
}

func (m *mockTokenService) GenerateTokens(userID uint) (string, string, int64, error) {
	return m.accessToken, m.refreshToken, m.expiresIn, nil
}

func TestService_RefreshToken_Success(t *testing.T) {
	repo := newMockUserRepository()
	tokenService := newMockTokenService()
	service := user.NewServiceWithJWT(repo, tokenService)
	ctx := context.Background()

	// Create a test user
	testUser := &user.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	err := repo.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Save a valid refresh token
	oldToken := "old-refresh-token"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = repo.SaveRefreshToken(ctx, testUser.ID, oldToken, expiresAt)
	if err != nil {
		t.Fatalf("failed to save refresh token: %v", err)
	}

	// Call RefreshToken
	authResponse, err := service.RefreshToken(ctx, oldToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	// Verify response
	if authResponse == nil {
		t.Fatal("expected auth response, got nil")
	}
	if authResponse.AccessToken != tokenService.accessToken {
		t.Errorf("expected access token %s, got %s", tokenService.accessToken, authResponse.AccessToken)
	}
	if authResponse.RefreshToken != tokenService.refreshToken {
		t.Errorf("expected refresh token %s, got %s", tokenService.refreshToken, authResponse.RefreshToken)
	}
	if authResponse.ExpiresIn != tokenService.expiresIn {
		t.Errorf("expected expiresIn %d, got %d", tokenService.expiresIn, authResponse.ExpiresIn)
	}

	// Verify old token was revoked
	oldTokenData, err := repo.GetRefreshToken(ctx, oldToken)
	if err != nil {
		t.Fatalf("failed to get old token: %v", err)
	}
	if oldTokenData == nil {
		t.Fatal("expected old token to exist")
	}
	if !oldTokenData.Revoked {
		t.Error("expected old token to be revoked")
	}

	// Verify new token was saved
	newTokenData, err := repo.GetRefreshToken(ctx, tokenService.refreshToken)
	if err != nil {
		t.Fatalf("failed to get new token: %v", err)
	}
	if newTokenData == nil {
		t.Error("expected new token to be saved")
	}
}

func TestService_RefreshToken_InvalidToken(t *testing.T) {
	repo := newMockUserRepository()
	tokenService := newMockTokenService()
	service := user.NewServiceWithJWT(repo, tokenService)
	ctx := context.Background()

	// Try to refresh with non-existent token
	_, err := service.RefreshToken(ctx, "non-existent-token")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
	if err != domain.ErrInvalidRefreshToken {
		t.Errorf("expected ErrInvalidRefreshToken, got %v", err)
	}
}

func TestService_RefreshToken_ExpiredToken(t *testing.T) {
	repo := newMockUserRepository()
	tokenService := newMockTokenService()
	service := user.NewServiceWithJWT(repo, tokenService)
	ctx := context.Background()

	// Create a test user
	testUser := &user.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	err := repo.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Save an expired refresh token
	expiredToken := "expired-refresh-token"
	expiresAt := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	err = repo.SaveRefreshToken(ctx, testUser.ID, expiredToken, expiresAt)
	if err != nil {
		t.Fatalf("failed to save expired token: %v", err)
	}

	// Try to refresh with expired token
	_, err = service.RefreshToken(ctx, expiredToken)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
	if err != domain.ErrRefreshTokenExpired {
		t.Errorf("expected ErrRefreshTokenExpired, got %v", err)
	}
}

func TestService_RefreshToken_RevokedToken(t *testing.T) {
	repo := newMockUserRepository()
	tokenService := newMockTokenService()
	service := user.NewServiceWithJWT(repo, tokenService)
	ctx := context.Background()

	// Create a test user
	testUser := &user.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	err := repo.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Save a refresh token
	revokedToken := "revoked-refresh-token"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = repo.SaveRefreshToken(ctx, testUser.ID, revokedToken, expiresAt)
	if err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Get the token to revoke it
	rt, err := repo.GetRefreshToken(ctx, revokedToken)
	if err != nil {
		t.Fatalf("failed to get token: %v", err)
	}

	// Revoke the token
	err = repo.RevokeRefreshToken(ctx, rt.ID)
	if err != nil {
		t.Fatalf("failed to revoke token: %v", err)
	}

	// Try to refresh with revoked token
	_, err = service.RefreshToken(ctx, revokedToken)
	if err == nil {
		t.Fatal("expected error for revoked token, got nil")
	}
	if err != domain.ErrRefreshTokenRevoked {
		t.Errorf("expected ErrRefreshTokenRevoked, got %v", err)
	}
}
