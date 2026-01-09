package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"manual-test-project/internal/domain"
	"manual-test-project/internal/interfaces"
	"manual-test-project/pkg/auth"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service interfaces.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service interfaces.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// ProfileResponse represents the user profile response
type ProfileResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with data"
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/me [get]
// @Security BearerAuth
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Extract user ID from JWT token (set by auth middleware)
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"error":  "Unable to extract user information",
		})
	}

	// Get user profile from service
	u, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		// Check if user not found
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"error":  "User not found",
			})
		}
		// Internal server error for unexpected errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to retrieve user profile",
		})
	}

	// Return profile information (no password hash)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": ProfileResponse{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		},
		"meta": fiber.Map{},
	})
}
