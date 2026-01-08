package main

import (
	"go.uber.org/fx"

	"manual-test-project/internal/infrastructure/database"
	"manual-test-project/internal/infrastructure/server"
	"manual-test-project/pkg/logger"
)

func main() {
	fx.New(
		// Register all modules
		logger.Module,
		database.Module,
		server.Module,
	).Run()
}
