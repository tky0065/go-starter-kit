package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"manual-test-project/internal/domain/user"
)

// mockAuthService is a mock for UserService that supports authentication
type mockAuthService struct {
	registerFunc     func(ctx context.Context, email, password string) (*user.User, error)
	authenticateFunc func(ctx context.Context, email, password string) (*user.AuthResponse, error)
}

func (m *mockAuthService) Register(ctx context.Context, email, password string) (*user.User, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, email, password)
	}
	return nil, nil
}

func (m *mockAuthService) Authenticate(ctx context.Context, email, password string) (*user.AuthResponse, error) {
	if m.authenticateFunc != nil {
		return m.authenticateFunc(ctx, email, password)
	}
	return nil, nil
}

func (m *mockAuthService) RefreshToken(ctx context.Context, oldToken string) (*user.AuthResponse, error) {
	return nil, nil
}

func TestAuthHandler_Login_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		authenticateFunc: func(ctx context.Context, email, password string) (*user.AuthResponse, error) {
			return &user.AuthResponse{
				AccessToken:  "test-access-token",
				RefreshToken: "test-refresh-token",
				ExpiresIn:    900,
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)
	app.Post("/api/v1/auth/login", handler.Login)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, fiber.StatusOK)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &response)

	if response["status"] != "success" {
		t.Errorf("status = %v, want success", response["status"])
	}

	data := response["data"].(map[string]interface{})
	if data["access_token"] != "test-access-token" {
		t.Errorf("access_token = %v, want test-access-token", data["access_token"])
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		authenticateFunc: func(ctx context.Context, email, password string) (*user.AuthResponse, error) {
			return nil, user.ErrInvalidCredentials
		},
	}

	handler := NewAuthHandler(mockService)
	app.Post("/api/v1/auth/login", handler.Login)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, fiber.StatusUnauthorized)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &response)

	if response["error"] == nil {
		t.Error("expected error message in response")
	}
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)
	app.Post("/api/v1/auth/login", handler.Login)

	tests := []struct {
		name    string
		reqBody map[string]string
	}{
		{
			name: "missing email",
			reqBody: map[string]string{
				"password": "password123",
			},
		},
		{
			name: "missing password",
			reqBody: map[string]string{
				"email": "test@example.com",
			},
		},
		{
			name: "invalid email format",
			reqBody: map[string]string{
				"email":    "not-an-email",
				"password": "password123",
			},
		},
		{
			name: "password too short",
			reqBody: map[string]string{
				"email":    "test@example.com",
				"password": "short",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() error: %v", err)
			}

			if resp.StatusCode != fiber.StatusBadRequest {
				t.Errorf("StatusCode = %v, want %v", resp.StatusCode, fiber.StatusBadRequest)
			}
		})
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)
	app.Post("/api/v1/auth/login", handler.Login)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, fiber.StatusBadRequest)
	}
}

func TestAuthHandler_Login_ServiceError(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		authenticateFunc: func(ctx context.Context, email, password string) (*user.AuthResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewAuthHandler(mockService)
	app.Post("/api/v1/auth/login", handler.Login)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error: %v", err)
	}

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, fiber.StatusInternalServerError)
	}
}
