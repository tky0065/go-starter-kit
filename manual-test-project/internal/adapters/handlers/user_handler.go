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

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service  interfaces.UserService
	validate *validator.Validate
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service interfaces.UserService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
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

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a paginated list of all users
// @Tags users
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10, max 100)"
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with data"
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [get]
// @Security BearerAuth
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// Get all users from service
	users, total, err := h.service.GetAll(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to retrieve users",
		})
	}

	// Convert to response format (no password hashes)
	userResponses := make([]ProfileResponse, len(users))
	for i, u := range users {
		userResponses[i] = ProfileResponse{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   userResponses,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// UpdateUser godoc
// @Summary Update user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Update user request"
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with data"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Parse user ID from path
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid user ID",
		})
	}

	// Parse request body
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid request body",
		})
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Validation failed: " + err.Error(),
		})
	}

	// Update user through service
	u, err := h.service.UpdateUser(c.Context(), uint(userID), req.Email)
	if err != nil {
		// Check specific errors
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"error":  "User not found",
			})
		}
		if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "Email already in use",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to update user",
		})
	}

	// Return updated user
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

// DeleteUser godoc
// @Summary Delete user
// @Description Soft delete a user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Parse user ID from path
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid user ID",
		})
	}

	// Delete user through service
	err = h.service.DeleteUser(c.Context(), uint(userID))
	if err != nil {
		// Check specific errors
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"error":  "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to delete user",
		})
	}

	// Return success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
		"meta":    fiber.Map{},
	})
}
