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
	fx.Invoke(func(h *AuthHandler, app *fiber.App) {
		h.RegisterRoutes(app)
	}),
)
