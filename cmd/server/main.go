package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matveevaolga/request-managing-app/internal/config"
	"github.com/matveevaolga/request-managing-app/internal/logger"
	"github.com/matveevaolga/request-managing-app/internal/repository"
	"github.com/matveevaolga/request-managing-app/internal/service"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler"
	"github.com/matveevaolga/request-managing-app/internal/transport/middleware"
)

func main() {
	cfg := initConfig()
	initLogger(cfg)

	db := initDB(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	projectTypeRepo := repository.NewProjectTypeRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	projectTypeService := service.NewProjectTypeService(projectTypeRepo)
	appService := service.NewApplicationService(appRepo, projectTypeRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	projectTypeHandler := handler.NewProjectTypeHandler(projectTypeService)
	appHandler := handler.NewApplicationHandler(appService)

	mux := registerRoutes(authHandler, projectTypeHandler, appHandler, authService)

	server := initServer(cfg, mux)
	runServer(server)
}

func initConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	return cfg
}

func initLogger(cfg *config.Config) {
	logger.Init(cfg.LogLevel)
	slog.Info("starting request-managing-app", "port", cfg.ServerPort)
}

func initDB(cfg *config.Config) *pgxpool.Pool {
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	return db
}

func registerRoutes(
	authHandler *handler.AuthHandler,
	projectTypeHandler *handler.ProjectTypeHandler,
	appHandler *handler.ApplicationHandler,
	authService *service.AuthService,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("GET /project/type", projectTypeHandler.GetAllProjects)
	mux.HandleFunc("POST /project/application/external", appHandler.Create)

	mux.HandleFunc("GET /project/application/external/list", func(w http.ResponseWriter, r *http.Request) {
		middleware.RequireAdmin(
			http.HandlerFunc(appHandler.GetAllFiltered),
		).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /project/application/external/{applicationId}", func(w http.ResponseWriter, r *http.Request) {
		middleware.RequireAdmin(
			http.HandlerFunc(appHandler.GetByID),
		).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /project/application/external/{applicationId}/accept", func(w http.ResponseWriter, r *http.Request) {
		middleware.RequireAdmin(
			middleware.Auth(authService)(http.HandlerFunc(appHandler.Accept)),
		).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /project/application/external/{applicationId}/reject", func(w http.ResponseWriter, r *http.Request) {
		middleware.RequireAdmin(
			middleware.Auth(authService)(http.HandlerFunc(appHandler.Reject)),
		).ServeHTTP(w, r)
	})

	return mux
}

func initServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      middleware.Logging(handler),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func runServer(server *http.Server) {
	slog.Info("Starting server", "addr", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	slog.Info("Stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server failed to gracefully shutdown", "error", err)
	}
	slog.Info("Server stopped")
}
