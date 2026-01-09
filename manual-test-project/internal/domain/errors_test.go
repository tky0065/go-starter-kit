package domain

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := &AppError{
		Code:    "TEST_ERROR",
		Message: "This is a test error",
		Status:  400,
		Details: nil,
	}

	assert.Equal(t, "This is a test error", err.Error())
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("User not found", "USER_NOT_FOUND")

	assert.NotNil(t, err)
	assert.Equal(t, "USER_NOT_FOUND", err.Code)
	assert.Equal(t, "User not found", err.Message)
	assert.Equal(t, fiber.StatusNotFound, err.Status)
	assert.Nil(t, err.Details)
}

func TestNewBadRequestError(t *testing.T) {
	details := map[string]string{"field": "email", "issue": "invalid format"}
	err := NewBadRequestError("Invalid input", "INVALID_INPUT", details)

	assert.NotNil(t, err)
	assert.Equal(t, "INVALID_INPUT", err.Code)
	assert.Equal(t, "Invalid input", err.Message)
	assert.Equal(t, fiber.StatusBadRequest, err.Status)
	assert.Equal(t, details, err.Details)
}

func TestNewInternalError(t *testing.T) {
	err := NewInternalError("Something went wrong", "INTERNAL_ERROR")

	assert.NotNil(t, err)
	assert.Equal(t, "INTERNAL_ERROR", err.Code)
	assert.Equal(t, "Something went wrong", err.Message)
	assert.Equal(t, fiber.StatusInternalServerError, err.Status)
	assert.Nil(t, err.Details)
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("Access denied", "UNAUTHORIZED")

	assert.NotNil(t, err)
	assert.Equal(t, "UNAUTHORIZED", err.Code)
	assert.Equal(t, "Access denied", err.Message)
	assert.Equal(t, fiber.StatusUnauthorized, err.Status)
	assert.Nil(t, err.Details)
}

func TestNewForbiddenError(t *testing.T) {
	err := NewForbiddenError("Forbidden resource", "FORBIDDEN")

	assert.NotNil(t, err)
	assert.Equal(t, "FORBIDDEN", err.Code)
	assert.Equal(t, "Forbidden resource", err.Message)
	assert.Equal(t, fiber.StatusForbidden, err.Status)
	assert.Nil(t, err.Details)
}

func TestNewConflictError(t *testing.T) {
	err := NewConflictError("Resource already exists", "CONFLICT")

	assert.NotNil(t, err)
	assert.Equal(t, "CONFLICT", err.Code)
	assert.Equal(t, "Resource already exists", err.Message)
	assert.Equal(t, fiber.StatusConflict, err.Status)
	assert.Nil(t, err.Details)
}
