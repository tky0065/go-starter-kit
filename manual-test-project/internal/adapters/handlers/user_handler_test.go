package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"manual-test-project/internal/adapters/middleware"
	"manual-test-project/internal/domain"
	"manual-test-project/internal/domain/user"
)

// mockUserService mocks the UserService interface for testing
type mockUserService struct {
	getProfileFunc func(ctx context.Context, userID uint) (*user.User, error)
}

func (m *mockUserService) GetProfile(ctx context.Context, userID uint) (*user.User, error) {
	if m.getProfileFunc != nil {
		return m.getProfileFunc(ctx, userID)
	}
	return nil, nil
}

func TestGetMe_Success(t *testing.T) {
	// Set JWT_SECRET for middleware
	secret := "test-secret-key"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Create a valid JWT token
	userID := uint(123)
	email := "test@example.com"
	createdAt := time.Now().Add(-24 * time.Hour) // Created yesterday

	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Create mock service
	mockService := &mockUserService{
		getProfileFunc: func(ctx context.Context, id uint) (*user.User, error) {
			if id != userID {
				t.Errorf("Expected userID %d, got %d", userID, id)
			}
			return &user.User{
				ID:        userID,
				Email:     email,
				CreatedAt: createdAt,
			}, nil
		},
	}

	// Create handler and app
	handler := NewUserHandler(mockService)
	app := fiber.New()
	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetMe)

	// Test request with valid token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify response contains profile data
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", body["status"])
	}

	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'data' field in response")
	}

	// Verify ID
	if userIDFloat, ok := data["id"].(float64); !ok || uint(userIDFloat) != userID {
		t.Errorf("Expected id %d, got %v", userID, data["id"])
	}

	// Verify email
	if gotEmail, ok := data["email"].(string); !ok || gotEmail != email {
		t.Errorf("Expected email %s, got %v", email, data["email"])
	}

	// Verify created_at exists
	if _, ok := data["created_at"].(string); !ok {
		t.Error("Expected created_at field in response")
	}

	// Verify meta field exists (architecture requirement)
	if _, ok := body["meta"].(map[string]interface{}); !ok {
		t.Error("Expected 'meta' field in response")
	}

	// Verify password hash is NOT included
	if _, exists := data["password_hash"]; exists {
		t.Error("Password hash should not be included in response")
	}
}

func TestGetMe_Unauthorized_NoToken(t *testing.T) {
	// Set JWT_SECRET for middleware
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	// Create mock service (won't be called)
	mockService := &mockUserService{}

	// Create handler and app
	handler := NewUserHandler(mockService)
	app := fiber.New()
	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetMe)

	// Test request without token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	// Verify error response format
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body["status"] != "error" {
		t.Errorf("Expected status 'error', got '%v'", body["status"])
	}
}

func TestGetMe_Unauthorized_InvalidToken(t *testing.T) {
	// Set JWT_SECRET for middleware
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	// Create mock service (won't be called)
	mockService := &mockUserService{}

	// Create handler and app
	handler := NewUserHandler(mockService)
	app := fiber.New()
	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetMe)

	// Test request with invalid token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestGetMe_UserNotFound(t *testing.T) {
	// Set JWT_SECRET for middleware
	secret := "test-secret-key"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Create a valid JWT token
	userID := uint(999)
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Create mock service that returns user not found
	mockService := &mockUserService{
		getProfileFunc: func(ctx context.Context, id uint) (*user.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}

	// Create handler and app
	handler := NewUserHandler(mockService)
	app := fiber.New()
	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetMe)

	// Test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 404 Not Found
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	// Verify error response
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body["status"] != "error" {
		t.Errorf("Expected status 'error', got '%v'", body["status"])
	}

	if body["error"] != "User not found" {
		t.Errorf("Expected error 'User not found', got '%v'", body["error"])
	}
}

func TestGetMe_InternalServerError(t *testing.T) {
	// Set JWT_SECRET for middleware
	secret := "test-secret-key"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Create a valid JWT token
	userID := uint(456)
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Create mock service that returns an unexpected error
	mockService := &mockUserService{
		getProfileFunc: func(ctx context.Context, id uint) (*user.User, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	// Create handler and app
	handler := NewUserHandler(mockService)
	app := fiber.New()
	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetMe)

	// Test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 500 Internal Server Error
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}

	// Verify error response
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body["status"] != "error" {
		t.Errorf("Expected status 'error', got '%v'", body["status"])
	}
}
