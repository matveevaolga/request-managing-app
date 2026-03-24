package main

import (
	"log/slog"
	"os"

	"github.com/matveevaolga/request-managing-app/internal/config"
	"github.com/matveevaolga/request-managing-app/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Init(cfg.LogLevel)
	slog.Info("Starting request-managing-app", "port", cfg.ServerPort)
}
