package auth

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrNoUserInContext is returned when no user is found in the context
	ErrNoUserInContext = errors.New("no user found in context")
	// ErrInvalidUserClaims is returned when the user claims are invalid
	ErrInvalidUserClaims = errors.New("invalid user claims in context")
)

// GetUserID extracts the user ID from the Fiber context.
// The JWT middleware stores the parsed token in c.Locals("user").
func GetUserID(c *fiber.Ctx) (uint, error) {
	user := c.Locals("user")
	if user == nil {
		return 0, ErrNoUserInContext
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return 0, ErrInvalidUserClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidUserClaims
	}

	// Extract "sub" claim which contains the user ID
	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, ErrInvalidUserClaims
	}

	// Parse user ID from subject claim
	var userID uint
	_, err := fmt.Sscanf(sub, "%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user ID from subject: %w", err)
	}

	return userID, nil
}
