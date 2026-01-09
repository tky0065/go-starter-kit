package handlers

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"manual-test-project/internal/domain/user"
	"manual-test-project/internal/interfaces"
	"manual-test-project/pkg/auth"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	service  interfaces.AuthService
	validate *validator.Validate
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(service interfaces.AuthService) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

// RegisterRequest represents the user registration request payload
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// RegisterResponse represents the user registration response
type RegisterResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		// Parse validation errors for user-friendly messages
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			errorMessages = append(errorMessages, formatValidationError(fieldErr))
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"fields": errorMessages,
		})
	}

	u, err := h.service.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		// Handle specific domain errors with appropriate HTTP status codes
		if errors.Is(err, user.ErrEmailAlreadyRegistered) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already registered",
			})
		}
		// Internal server error for unexpected errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user account",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": RegisterResponse{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		},
	})
}

// LoginRequest represents the user login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// LoginResponse represents the user login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		// Parse validation errors for user-friendly messages
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			errorMessages = append(errorMessages, formatValidationError(fieldErr))
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"fields": errorMessages,
		})
	}

	authResponse, err := h.service.Authenticate(c.Context(), req.Email, req.Password)
	if err != nil {
		// Handle specific domain errors with appropriate HTTP status codes
		if errors.Is(err, user.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		// Internal server error for unexpected errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": LoginResponse{
			AccessToken:  authResponse.AccessToken,
			RefreshToken: authResponse.RefreshToken,
			ExpiresIn:    authResponse.ExpiresIn,
		},
	})
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Refresh godoc
// @Summary Refresh access token
// @Description Exchange a valid refresh token for new access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		// Parse validation errors for user-friendly messages
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			errorMessages = append(errorMessages, formatValidationError(fieldErr))
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"fields": errorMessages,
		})
	}

	authResponse, err := h.service.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		// Handle specific domain errors with appropriate HTTP status codes
		if errors.Is(err, user.ErrInvalidRefreshToken) ||
			errors.Is(err, user.ErrRefreshTokenExpired) ||
			errors.Is(err, user.ErrRefreshTokenRevoked) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired refresh token",
			})
		}
		// Internal server error for unexpected errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token refresh failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": LoginResponse{
			AccessToken:  authResponse.AccessToken,
			RefreshToken: authResponse.RefreshToken,
			ExpiresIn:    authResponse.ExpiresIn,
		},
	})
}

// GetCurrentUser godoc
// @Summary Get current user information
// @Description Get the current authenticated user's information
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/users/me [get]
// @Security BearerAuth
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"error":  "Unable to extract user information",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user_id": userID,
			"message": "You are authenticated",
		},
	})
}

// formatValidationError formats a validator.FieldError into a user-friendly message
func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	switch err.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + err.Param() + " characters"
	case "max":
		return field + " must be at most " + err.Param() + " characters"
	default:
		return field + " is invalid"
	}
}
