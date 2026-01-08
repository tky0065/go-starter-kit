package http

import (
	"github.com/gofiber/fiber/v2"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// RegisterHealthRoutes registers health check routes
func RegisterHealthRoutes(app *fiber.App) {
	app.Get("/health", healthHandler)
}

// healthHandler handles health check requests
func healthHandler(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status: "ok",
	})
}
