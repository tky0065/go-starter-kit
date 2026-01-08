package logger

import (
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// Module provides the logger dependency via fx
var Module = fx.Module("logger",
	fx.Provide(NewLogger),
)

// NewLogger creates a new zerolog logger instance
func NewLogger() zerolog.Logger {
	// Use JSON format in production, console format in development
	env := os.Getenv("APP_ENV")

	var logger zerolog.Logger
	if env == "production" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	return logger
}
