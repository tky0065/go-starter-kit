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

func TestUserRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	// Test empty database
	users, _, err := repo.FindAll(ctx, 1, 10)
	if err != nil {
		t.Fatalf("FindAll failed on empty database: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}

	// Create test users
	user1 := &user.User{
		Email:        "user1@example.com",
		PasswordHash: "hash1",
	}
	user2 := &user.User{
		Email:        "user2@example.com",
		PasswordHash: "hash2",
	}

	err = repo.CreateUser(ctx, user1)
	if err != nil {
		t.Fatalf("failed to create user1: %v", err)
	}

	err = repo.CreateUser(ctx, user2)
	if err != nil {
		t.Fatalf("failed to create user2: %v", err)
	}

	// Test FindAll returns all users
	users, _, err = repo.FindAll(ctx, 1, 10)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}

	// Verify emails
	emails := make(map[string]bool)
	for _, u := range users {
		emails[u.Email] = true
	}

	if !emails["user1@example.com"] || !emails["user2@example.com"] {
		t.Error("expected both user emails to be present")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	// Create a test user
	testUser := &user.User{
		Email:        "original@example.com",
		PasswordHash: "originalhash",
	}

	err := repo.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Update the user
	testUser.Email = "updated@example.com"
	err = repo.Update(ctx, testUser)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Retrieve and verify update
	updated, err := repo.FindByID(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("failed to find updated user: %v", err)
	}

	if updated == nil {
		t.Fatal("expected user to exist after update")
	}

	if updated.Email != "updated@example.com" {
		t.Errorf("expected email to be 'updated@example.com', got '%s'", updated.Email)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	// Create a test user
	testUser := &user.User{
		Email:        "todelete@example.com",
		PasswordHash: "hash",
	}

	err := repo.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	userID := testUser.ID

	// Delete the user (soft delete)
	err = repo.Delete(ctx, userID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify user is soft deleted (not returned in normal queries)
	found, err := repo.FindByID(ctx, userID)
	if err != nil {
		t.Fatalf("FindByID failed after delete: %v", err)
	}

	if found != nil {
		t.Error("expected user to not be found after soft delete")
	}

	// Verify user still exists in database with deleted_at set
	var deletedUser user.User
	err = db.Unscoped().Where("id = ?", userID).First(&deletedUser).Error
	if err != nil {
		t.Fatalf("failed to find deleted user with Unscoped: %v", err)
	}

	if deletedUser.DeletedAt.Time.IsZero() {
		t.Error("expected DeletedAt to be set after soft delete")
	}
}
