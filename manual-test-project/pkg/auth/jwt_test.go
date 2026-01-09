package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateTokens(t *testing.T) {
	secret := "test-secret-key-for-testing-purposes"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 7 * 24 * time.Hour

	jwtService := NewJWTService(secret, accessExpiry, refreshExpiry)

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "valid user ID",
			userID:  123,
			wantErr: false,
		},
		{
			name:    "user ID zero",
			userID:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessToken, refreshToken, expiresIn, err := jwtService.GenerateTokens(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if accessToken == "" {
					t.Error("GenerateTokens() accessToken is empty")
				}
				if refreshToken == "" {
					t.Error("GenerateTokens() refreshToken is empty")
				}
				if expiresIn <= 0 {
					t.Error("GenerateTokens() expiresIn should be positive")
				}
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	secret := "test-secret-key-for-testing-purposes"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 7 * 24 * time.Hour

	jwtService := NewJWTService(secret, accessExpiry, refreshExpiry)

	// Generate a valid token
	userID := uint(456)
	accessToken, _, _, err := jwtService.GenerateTokens(userID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		wantID    uint
		wantValid bool
	}{
		{
			name:      "valid token",
			token:     accessToken,
			wantID:    userID,
			wantValid: true,
		},
		{
			name:      "empty token",
			token:     "",
			wantID:    0,
			wantValid: false,
		},
		{
			name:      "invalid token",
			token:     "invalid.token.string",
			wantID:    0,
			wantValid: false,
		},
		{
			name:      "tampered token",
			token:     accessToken + "tampered",
			wantID:    0,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := jwtService.ValidateAccessToken(tt.token)

			if tt.wantValid && err != nil {
				t.Errorf("ValidateAccessToken() unexpected error = %v", err)
				return
			}

			if !tt.wantValid && err == nil {
				t.Error("ValidateAccessToken() expected error but got nil")
				return
			}

			if tt.wantValid && gotID != tt.wantID {
				t.Errorf("ValidateAccessToken() gotID = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	secret := "test-secret-key-for-testing-purposes"
	accessExpiry := -1 * time.Second // Already expired
	refreshExpiry := 7 * 24 * time.Hour

	jwtService := NewJWTService(secret, accessExpiry, refreshExpiry)

	// Generate an expired token
	userID := uint(789)
	accessToken, _, _, err := jwtService.GenerateTokens(userID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate the expired token
	_, err = jwtService.ValidateAccessToken(accessToken)
	if err == nil {
		t.Error("ValidateAccessToken() expected error for expired token but got nil")
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	secret1 := "secret-one"
	secret2 := "secret-two"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 7 * 24 * time.Hour

	jwtService1 := NewJWTService(secret1, accessExpiry, refreshExpiry)
	jwtService2 := NewJWTService(secret2, accessExpiry, refreshExpiry)

	// Generate token with secret1
	userID := uint(999)
	accessToken, _, _, err := jwtService1.GenerateTokens(userID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with secret2
	_, err = jwtService2.ValidateAccessToken(accessToken)
	if err == nil {
		t.Error("ValidateAccessToken() expected error for wrong secret but got nil")
	}
}

func TestJWTClaims(t *testing.T) {
	secret := "test-secret-key-for-testing-purposes"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 7 * 24 * time.Hour

	jwtService := NewJWTService(secret, accessExpiry, refreshExpiry)

	userID := uint(111)
	accessToken, _, _, err := jwtService.GenerateTokens(userID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Parse token without validation to check claims structure
	token, _, err := jwt.NewParser().ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get claims from token")
	}

	// Check standard claims exist
	if _, ok := claims["sub"]; !ok {
		t.Error("Token missing 'sub' claim")
	}
	if _, ok := claims["exp"]; !ok {
		t.Error("Token missing 'exp' claim")
	}
	if _, ok := claims["iat"]; !ok {
		t.Error("Token missing 'iat' claim")
	}
}
