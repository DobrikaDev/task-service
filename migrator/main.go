package main

import (
	"DobrikaDev/task-service/di"
	"DobrikaDev/task-service/utils/config"
	"DobrikaDev/task-service/utils/logger"
	"context"
	"os"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoadConfigFromFile("deployments/config.yaml")
	logger, _ := logger.NewLogger()
	defer logger.Sync()
	container := di.NewContainer(ctx, cfg, logger)

	goose.SetTableName("migrations")

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("Error setting dialect:", zap.Error(err))
		os.Exit(1)
	}

	if err := goose.Up(container.GetDB().DB, "migrations/postgres"); err != nil {
		logger.Error("Error running migrations:", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Migrations completed successfully")
}
