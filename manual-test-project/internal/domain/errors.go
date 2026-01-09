package domain

import "errors"

// Standardized Domain Errors
var (
	// User Domain Errors
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrRefreshTokenExpired    = errors.New("refresh token expired")
	ErrRefreshTokenRevoked    = errors.New("refresh token revoked")
)
