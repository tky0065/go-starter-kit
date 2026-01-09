//go:build integration
// +build integration

package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"manual-test-project/internal/adapters/handlers"
	"manual-test-project/internal/adapters/repository"
	"manual-test-project/internal/domain/user"
	"manual-test-project/pkg/auth"
)

func setupIntegrationTest(t *testing.T) (*fiber.App, *user.Service, *gorm.DB) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Auto-migrate
	err = db.AutoMigrate(&user.User{}, &user.RefreshToken{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	// Create repository
	repo := repository.NewUserRepository(db)

	// Create token service
	tokenService := auth.NewJWTService("test-secret-key", 15*time.Minute, 7*24*time.Hour)

	// Create user service
	service := user.NewServiceWithJWT(repo, tokenService)

	// Create handler
	handler := handlers.NewAuthHandler(service)

	// Create Fiber app
	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")
	handlers.RegisterAuthRoutes(v1, handler)

	return app, service, db
}

func TestRefreshToken_Integration_Success(t *testing.T) {
	app, service, _ := setupIntegrationTest(t)
	ctx := context.Background()

	// Register a user first
	_, err := service.Register(ctx, "test@example.com", "password123")
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	// Authenticate to get tokens
	authResp, err := service.Authenticate(ctx, "test@example.com", "password123")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}

	// Now test the refresh endpoint
	refreshReq := map[string]string{
		"refresh_token": authResp.RefreshToken,
	}
	body, _ := json.Marshal(refreshReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("expected status success, got %v", result["status"])
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}

	if _, ok := data["access_token"]; !ok {
		t.Error("expected access_token in response")
	}
	if _, ok := data["refresh_token"]; !ok {
		t.Error("expected refresh_token in response")
	}
	if _, ok := data["expires_in"]; !ok {
		t.Error("expected expires_in in response")
	}

	// Verify new refresh token is different from old one
	newRefreshToken, ok := data["refresh_token"].(string)
	if !ok {
		t.Fatal("refresh_token should be a string")
	}
	if newRefreshToken == authResp.RefreshToken {
		t.Error("expected new refresh token to be different (token rotation)")
	}
}

func TestRefreshToken_Integration_InvalidToken(t *testing.T) {
	app, _, _ := setupIntegrationTest(t)

	refreshReq := map[string]string{
		"refresh_token": "invalid-token",
	}
	body, _ := json.Marshal(refreshReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestRefreshToken_Integration_MissingToken(t *testing.T) {
	app, _, _ := setupIntegrationTest(t)

	refreshReq := map[string]string{}
	body, _ := json.Marshal(refreshReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
