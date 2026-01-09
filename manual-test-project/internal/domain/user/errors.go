package user

import "errors"

// Sentinel errors for user domain operations
var (
	// ErrEmailAlreadyRegistered is returned when attempting to register with an existing email
	ErrEmailAlreadyRegistered = errors.New("email already registered")

	// ErrInvalidCredentials is returned when login credentials are incorrect
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound is returned when a user cannot be found
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidRefreshToken is returned when a refresh token is invalid or not found
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// ErrRefreshTokenExpired is returned when a refresh token has expired
	ErrRefreshTokenExpired = errors.New("refresh token expired")

	// ErrRefreshTokenRevoked is returned when a refresh token has been revoked
	ErrRefreshTokenRevoked = errors.New("refresh token revoked")
)
