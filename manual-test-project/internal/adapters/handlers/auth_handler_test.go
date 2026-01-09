package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"manual-test-project/internal/adapters/handlers"
	"manual-test-project/internal/domain"
	"manual-test-project/internal/domain/user"
)

type MockService struct {
	ShouldReturnError error
}

func (m *MockService) Register(ctx context.Context, email, password string) (*user.User, error) {
	if m.ShouldReturnError != nil {
		return nil, m.ShouldReturnError
	}
	return &user.User{ID: 1, Email: email, CreatedAt: time.Now()}, nil
}

func (m *MockService) Authenticate(ctx context.Context, email, password string) (*user.AuthResponse, error) {
	if m.ShouldReturnError != nil {
		return nil, m.ShouldReturnError
	}
	return &user.AuthResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresIn:    900,
	}, nil
}

func (m *MockService) RefreshToken(ctx context.Context, oldToken string) (*user.AuthResponse, error) {
	if m.ShouldReturnError != nil {
		return nil, m.ShouldReturnError
	}
	return &user.AuthResponse{
		AccessToken:  "test-new-access-token",
		RefreshToken: "test-new-refresh-token",
		ExpiresIn:    900,
	}, nil
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]interface{}
		serviceError   error
		expectedStatus int
		description    string
	}{
		{
			name: "Success",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			serviceError:   nil,
			expectedStatus: http.StatusCreated,
			description:    "Valid registration should return 201",
		},
		{
			name: "Duplicate Email",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			serviceError:   domain.ErrEmailAlreadyRegistered,
			expectedStatus: http.StatusConflict,
			description:    "Duplicate email should return 409",
		},
		{
			name: "Invalid Email",
			body: map[string]interface{}{
				"email":    "notanemail",
				"password": "password123",
			},
			serviceError:   nil,
			expectedStatus: http.StatusBadRequest,
			description:    "Invalid email format should return 400",
		},
		{
			name: "Empty Email",
			body: map[string]interface{}{
				"email":    "",
				"password": "password123",
			},
			serviceError:   nil,
			expectedStatus: http.StatusBadRequest,
			description:    "Empty email should return 400",
		},
		{
			name: "Empty Password",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "",
			},
			serviceError:   nil,
			expectedStatus: http.StatusBadRequest,
			description:    "Empty password should return 400",
		},
		{
			name: "Short Password",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "short",
			},
			serviceError:   nil,
			expectedStatus: http.StatusBadRequest,
			description:    "Password too short should return 400",
		},
		{
			name: "Malformed JSON",
			body: map[string]interface{}{
				"email": 12345, // Wrong type
			},
			serviceError:   nil,
			expectedStatus: http.StatusBadRequest,
			description:    "Malformed JSON should return 400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockSvc := &MockService{ShouldReturnError: tt.serviceError}
			h := handlers.NewAuthHandler(mockSvc)
			app.Post("/api/v1/auth/register", h.Register)

			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Test request failed: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.description, tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
