package database

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Module provides the database dependency via fx
var Module = fx.Module("database",
	fx.Provide(NewDatabase),
	fx.Invoke(registerHooks),
)

// NewDatabase creates a new GORM database connection
func NewDatabase(logger zerolog.Logger) (*gorm.DB, error) {
	// Build DSN from environment variables
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "manual-test-project"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info().Msg("Successfully connected to database")

	// AutoMigrate database schemas
	// Add your domain models here when ready
	// Example: db.AutoMigrate(&models.User{})

	logger.Info().Msg("Database migrations completed")

	return db, nil
}

// registerHooks registers lifecycle hooks for graceful shutdown
func registerHooks(lifecycle fx.Lifecycle, db *gorm.DB, logger zerolog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Closing database connection")
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
