package middleware

import (
	"os"

	"github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// NewAuthMiddleware creates a new JWT authentication middleware
func NewAuthMiddleware() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Return standard error format
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"error":  "Unauthorized: Invalid or missing token",
			})
		},
		// Default behavior: token parsed and stored in c.Locals("user")
		// We can access it later via GetUserID helper
	})
}
