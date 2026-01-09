package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"manual-test-project/internal/domain/user"
)

var Module = fx.Module("handlers",
	fx.Provide(func(s *user.Service) *AuthHandler {
		return NewAuthHandler(s)
	}),
	fx.Provide(func(s *user.Service) *UserHandler {
		return NewUserHandler(s)
	}),
	fx.Invoke(RegisterAllRoutes),
)

// RegisterAllRoutes registers all application routes with public and protected groups
func RegisterAllRoutes(authHandler *AuthHandler, userHandler *UserHandler, app *fiber.App, authMiddleware fiber.Handler) {
	// Public routes (no authentication required)
	public := app.Group("/api/v1")
	public.Post("/auth/register", authHandler.Register)
	public.Post("/auth/login", authHandler.Login)
	public.Post("/auth/refresh", authHandler.Refresh)

	// Protected routes (authentication required)
	protected := app.Group("/api/v1", authMiddleware)
	protected.Get("/users/me", userHandler.GetMe)
}
