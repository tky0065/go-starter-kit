package repository_test

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"manual-test-project/internal/adapters/repository"
	"manual-test-project/internal/domain/user"
	"manual-test-project/internal/interfaces"
)

func TestUserRepository_ImplementsInterface(t *testing.T) {
	var _ interfaces.UserRepository = (*repository.UserRepository)(nil)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Auto-migrate tables
	err = db.AutoMigrate(&user.User{}, &user.RefreshToken{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestUserRepository_RevokeRefreshToken(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	// Create a test user first
	testUser := &user.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	err := repo.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Save a refresh token
	tokenString := "test-refresh-token"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = repo.SaveRefreshToken(ctx, testUser.ID, tokenString, expiresAt)
	if err != nil {
		t.Fatalf("failed to save refresh token: %v", err)
	}

	// Retrieve the token to get its ID
	rt, err := repo.GetRefreshToken(ctx, tokenString)
	if err != nil {
		t.Fatalf("failed to get refresh token: %v", err)
	}
	if rt == nil {
		t.Fatal("expected refresh token to exist")
	}

	// Verify token is not revoked initially
	if rt.Revoked {
		t.Error("expected token to not be revoked initially")
	}

	// Revoke the token
	err = repo.RevokeRefreshToken(ctx, rt.ID)
	if err != nil {
		t.Fatalf("failed to revoke refresh token: %v", err)
	}

	// Retrieve again and verify it's revoked
	revokedToken, err := repo.GetRefreshToken(ctx, tokenString)
	if err != nil {
		t.Fatalf("failed to get refresh token after revocation: %v", err)
	}
	if revokedToken == nil {
		t.Fatal("expected revoked token to still exist")
	}
	if !revokedToken.Revoked {
		t.Error("expected token to be revoked after RevokeRefreshToken call")
	}
}
