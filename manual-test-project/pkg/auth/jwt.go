package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTService creates a new JWTService instance.
func NewJWTService(secret string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secret:        secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateTokens generates both access and refresh tokens for a user.
// Returns accessToken, refreshToken, expiresIn (seconds), and error.
func (s *JWTService) GenerateTokens(UserID uint) (string, string, int64, error) {
	// Generate access token (JWT)
	now := time.Now()
	expiresAt := now.Add(s.accessExpiry)

	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", UserID),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token (opaque string)
	refreshToken, err := generateOpaqueToken()
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresIn := int64(s.accessExpiry.Seconds())

	return accessToken, refreshToken, expiresIn, nil
}

// ValidateAccessToken validates an access token and returns the user ID.
func (s *JWTService) ValidateAccessToken(tokenString string) (uint, error) {
	if tokenString == "" {
		return 0, fmt.Errorf("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, fmt.Errorf("failed to parse claims")
	}

	// Parse user ID from subject claim
	var userID uint
	_, err = fmt.Sscanf(claims.Subject, "%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user ID from subject: %w", err)
	}

	return userID, nil
}

// generateOpaqueToken generates a cryptographically secure random token.
func generateOpaqueToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
