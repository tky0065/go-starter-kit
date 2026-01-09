package middleware

import (
	"errors"
	"os"

	"manual-test-project/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ErrorHandler is a centralized error handler for Fiber that formats all errors
// into a consistent JSON structure following the API standardization requirements.
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	code := fiber.StatusInternalServerError
	resp := fiber.Map{
		"status":  "error",
		"code":    "INTERNAL_SERVER_ERROR",
		"message": "Internal server error",
		"details": nil,
	}

	// Flag to check if we should mask the error message (Production)
	isProd := os.Getenv("APP_ENV") == "production"

	// 1. Handle Domain standard errors (map standard errors to AppErrors)
	if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
		err = domain.NewConflictError("Email already registered", "EMAIL_ALREADY_REGISTERED")
	} else if errors.Is(err, domain.ErrInvalidCredentials) {
		err = domain.NewUnauthorizedError("Invalid email or password", "INVALID_CREDENTIALS")
	} else if errors.Is(err, domain.ErrUserNotFound) {
		err = domain.NewNotFoundError("User not found", "USER_NOT_FOUND")
	} else if errors.Is(err, domain.ErrInvalidRefreshToken) || errors.Is(err, domain.ErrRefreshTokenExpired) || errors.Is(err, domain.ErrRefreshTokenRevoked) {
		err = domain.NewUnauthorizedError(err.Error(), "AUTH_TOKEN_ERROR")
	}

	// 2. Handle Fiber Errors (including 404, 405, etc.)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		resp["message"] = fiberErr.Message
		resp["code"] = mapHTTPStatusToCode(code)
	}

	// 3. Handle Domain AppErrors (business logic errors)
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		code = appErr.Status
		resp["message"] = appErr.Message
		resp["code"] = appErr.Code
		resp["details"] = appErr.Details
	}

	// AC3: Mask internal error messages in production
	if code == fiber.StatusInternalServerError && isProd {
		resp["message"] = "Internal server error"
	}

	// Logging with context
	log.Error().
		Err(err).
		Int("status", code).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Msg("API Error")

	return c.Status(code).JSON(resp)
}

// mapHTTPStatusToCode maps HTTP status codes to readable error codes.
func mapHTTPStatusToCode(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return "BAD_REQUEST"
	case fiber.StatusUnauthorized:
		return "UNAUTHORIZED"
	case fiber.StatusForbidden:
		return "FORBIDDEN"
	case fiber.StatusNotFound:
		return "NOT_FOUND"
	case fiber.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case fiber.StatusConflict:
		return "CONFLICT"
	case fiber.StatusUnprocessableEntity:
		return "UNPROCESSABLE_ENTITY"
	case fiber.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	default:
		return "HTTP_ERROR"
	}
}
