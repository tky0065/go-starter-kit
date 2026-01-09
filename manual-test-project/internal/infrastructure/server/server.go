package server

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/gorm"

	httphandlers "manual-test-project/internal/adapters/http"
	"manual-test-project/internal/adapters/middleware"
)

// Module provides the Fiber server dependency via fx
var Module = fx.Module("server",
	fx.Provide(NewServer),
	fx.Provide(middleware.NewAuthMiddleware),
	fx.Invoke(registerHooks),
)

// NewServer creates and configures a new Fiber application
func NewServer(logger zerolog.Logger, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "manual-test-project",
		ErrorHandler: middleware.ErrorHandler,
	})

	logger.Info().Msg("Fiber server initialized with centralized error handler")

	// Register routes
	httphandlers.RegisterHealthRoutes(app)

	return app
}

// registerHooks registers lifecycle hooks for server startup and shutdown
func registerHooks(lifecycle fx.Lifecycle, app *fiber.App, logger zerolog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := getEnv("APP_PORT", "3000")
			logger.Info().Str("port", port).Msg("Starting Fiber server")

			go func() {
				if err := app.Listen(":" + port); err != nil {
					logger.Fatal().Err(err).Msg("Failed to start server")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Shutting down Fiber server gracefully")
			return app.Shutdown()
		},
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
