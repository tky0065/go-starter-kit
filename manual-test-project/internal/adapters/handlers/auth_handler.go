package handlers

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"manual-test-project/internal/domain"
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
		return domain.NewBadRequestError("Invalid JSON format", "INVALID_JSON", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		details := make(map[string]string)
		for _, fieldErr := range validationErrors {
			details[fieldErr.Field()] = formatValidationError(fieldErr)
		}
		return domain.NewBadRequestError("Validation failed", "VALIDATION_FAILED", details)
	}

	// Register user
	newUser, err := h.service.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		return err // Handled by middleware
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": RegisterResponse{
			ID:        newUser.ID,
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
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
		return domain.NewBadRequestError("Invalid JSON format", "INVALID_JSON", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		details := make(map[string]string)
		for _, fieldErr := range validationErrors {
			details[fieldErr.Field()] = formatValidationError(fieldErr)
		}
		return domain.NewBadRequestError("Validation failed", "VALIDATION_FAILED", details)
	}

	// Authenticate user
	authResponse, err := h.service.Authenticate(c.Context(), req.Email, req.Password)
	if err != nil {
		return err // Handled by middleware
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
		return domain.NewBadRequestError("Invalid JSON format", "INVALID_JSON", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return domain.NewBadRequestError("Refresh token is required", "VALIDATION_FAILED", nil)
	}

	authResponse, err := h.service.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return err // Handled by middleware
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
		return domain.NewUnauthorizedError("Unable to extract user information", "UNAUTHORIZED")
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
