package domain

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// AppError represents a structured application error with HTTP status and details.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"` // HTTP Status, not serialized in JSON
	Details any    `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// NewNotFoundError creates a 404 error.
func NewNotFoundError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusNotFound,
		Details: nil,
	}
}

// NewBadRequestError creates a 400 error with optional validation details.
func NewBadRequestError(msg string, code string, details any) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusBadRequest,
		Details: details,
	}
}

// NewInternalError creates a 500 error.
func NewInternalError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusInternalServerError,
		Details: nil,
	}
}

// NewUnauthorizedError creates a 401 error.
func NewUnauthorizedError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusUnauthorized,
		Details: nil,
	}
}

// NewForbiddenError creates a 403 error.
func NewForbiddenError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusForbidden,
		Details: nil,
	}
}

// NewConflictError creates a 409 error.
func NewConflictError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusConflict,
		Details: nil,
	}
}

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
