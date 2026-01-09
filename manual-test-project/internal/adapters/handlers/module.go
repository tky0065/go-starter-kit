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
	// Create API group hierarchy for versioning
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register domain-specific routes
	RegisterAuthRoutes(v1, authHandler)
	RegisterUserRoutes(v1, userHandler, authMiddleware)
}

// RegisterAuthRoutes registers authentication-related routes (public)
func RegisterAuthRoutes(v1 fiber.Router, authHandler *AuthHandler) {
	authGroup := v1.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.Refresh)
}

// RegisterUserRoutes registers user-related routes (protected)
func RegisterUserRoutes(v1 fiber.Router, userHandler *UserHandler, authMiddleware fiber.Handler) {
	userGroup := v1.Group("/users", authMiddleware)
	userGroup.Get("/me", userHandler.GetMe)
	userGroup.Get("", userHandler.GetAllUsers)
	userGroup.Put("/:id", userHandler.UpdateUser)
	userGroup.Delete("/:id", userHandler.DeleteUser)
}
