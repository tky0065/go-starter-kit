package middleware

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
	"manual-test-project/pkg/auth"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
	// Set JWT_SECRET for testing
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Create Fiber app
	app := fiber.New()

	// Apply middleware to protected route
	protected := app.Group("/api/v1", NewAuthMiddleware())
	protected.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Test request without token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401 Unauthorized
	if resp.StatusCode != fiber.StatusUnauthorized {
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

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Set JWT_SECRET for testing
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Create Fiber app
	app := fiber.New()

	// Apply middleware to protected route
	protected := app.Group("/api/v1", NewAuthMiddleware())
	protected.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Test request with invalid token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401 Unauthorized
	if resp.StatusCode != fiber.StatusUnauthorized {
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

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Set JWT_SECRET for testing
	secret := "test-secret-key"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Create a valid JWT token
	userID := uint(123)
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

	// Create Fiber app
	app := fiber.New()

	// Apply middleware to protected route
	protected := app.Group("/api/v1", NewAuthMiddleware())
	protected.Get("/protected", func(c *fiber.Ctx) error {
		// Extract user ID using helper
		extractedUserID, err := auth.GetUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "success",
			"user_id": extractedUserID,
		})
	})

	// Test request with valid token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 200 OK
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify response contains user ID
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body["message"] != "success" {
		t.Errorf("Expected message 'success', got '%v'", body["message"])
	}

	// Verify user_id is correctly extracted
	if userIDFloat, ok := body["user_id"].(float64); !ok || uint(userIDFloat) != userID {
		t.Errorf("Expected user_id %d, got %v", userID, body["user_id"])
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// Set JWT_SECRET for testing
	secret := "test-secret-key"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Create an expired JWT token
	userID := uint(123)
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Create Fiber app
	app := fiber.New()

	// Apply middleware to protected route
	protected := app.Group("/api/v1", NewAuthMiddleware())
	protected.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Test request with expired token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401 Unauthorized
	if resp.StatusCode != fiber.StatusUnauthorized {
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

func TestAuthMiddleware_WrongSignature(t *testing.T) {
	// Set JWT_SECRET for testing
	secret := "test-secret-key"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Create a token with wrong secret
	wrongSecret := "wrong-secret-key"
	userID := uint(123)
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(wrongSecret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Create Fiber app
	app := fiber.New()

	// Apply middleware to protected route
	protected := app.Group("/api/v1", NewAuthMiddleware())
	protected.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Test request with wrong signature
	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert 401 Unauthorized
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestGetUserID_Success(t *testing.T) {
	// Set JWT_SECRET for testing
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

	// Create Fiber app
	app := fiber.New()

	// Apply middleware and test GetUserID
	protected := app.Group("/api/v1", NewAuthMiddleware())
	protected.Get("/me", func(c *fiber.Ctx) error {
		extractedUserID, err := auth.GetUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"user_id": extractedUserID,
		})
	})

	// Test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Assert success
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify user ID matches
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if userIDFloat, ok := body["user_id"].(float64); !ok || uint(userIDFloat) != userID {
		t.Errorf("Expected user_id %d, got %v", userID, body["user_id"])
	}
}
