package main

import (
	"context"
	"embed"
	"log"
	"log/slog"

	"github.com/matveevaolga/request-managing-app/internal/config"
	"github.com/matveevaolga/request-managing-app/internal/logger"
	"github.com/matveevaolga/request-managing-app/internal/repository"
)

//go:embed data/*
var seedFS embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logger.Init(cfg.LogLevel)
	slog.Info("Starting seed")

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	userRepo := repository.NewUserRepository(db)
	typeRepo := repository.NewProjectTypeRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	if err := seedUsers(ctx, userRepo); err != nil {
		slog.Error("Failed to seed users", "error", err)
		log.Fatal(err)
	}

	if err := seedProjectTypes(ctx, typeRepo); err != nil {
		slog.Error("Failed to seed project types", "error", err)
		log.Fatal(err)
	}

	if err := seedApplications(ctx, appRepo, typeRepo); err != nil {
		slog.Error("Failed to seed applications", "error", err)
		log.Fatal(err)
	}

	slog.Info("Seed completed successfully")
}
