package handlers

import (
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
)

func TestProtectedRoute_GetCurrentUser_WithoutToken(t *testing.T) {
	// Set JWT_SECRET
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	// Setup app with protected route
	app := fiber.New()
	handler := &AuthHandler{
		service:  nil, // Not needed for this test
		validate: nil,
	}

	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetCurrentUser)

	// Test without token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	// Verify error format
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body["status"] != "error" {
		t.Errorf("Expected status 'error', got '%v'", body["status"])
	}
}

func TestProtectedRoute_GetCurrentUser_WithInvalidToken(t *testing.T) {
	// Set JWT_SECRET
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	// Setup app with protected route
	app := fiber.New()
	handler := &AuthHandler{
		service:  nil,
		validate: nil,
	}

	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetCurrentUser)

	// Test with invalid token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestProtectedRoute_GetCurrentUser_WithValidToken(t *testing.T) {
	// Set JWT_SECRET
	secret := "test-secret"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Create valid token
	userID := uint(789)
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

	// Setup app with protected route
	app := fiber.New()
	handler := &AuthHandler{
		service:  nil,
		validate: nil,
	}

	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetCurrentUser)

	// Test with valid token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 200 OK
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify response contains user_id
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

	if userIDFloat, ok := data["user_id"].(float64); !ok || uint(userIDFloat) != userID {
		t.Errorf("Expected user_id %d, got %v", userID, data["user_id"])
	}
}

func TestProtectedRoute_GetCurrentUser_WithExpiredToken(t *testing.T) {
	// Set JWT_SECRET
	secret := "test-secret"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Create expired token
	userID := uint(789)
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Setup app with protected route
	app := fiber.New()
	handler := &AuthHandler{
		service:  nil,
		validate: nil,
	}

	protected := app.Group("/api/v1", middleware.NewAuthMiddleware())
	protected.Get("/users/me", handler.GetCurrentUser)

	// Test with expired token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}
