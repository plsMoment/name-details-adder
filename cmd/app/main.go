package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"name-details-adder/config"
	"name-details-adder/internal/db"
	"name-details-adder/transport/rest"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		e := fmt.Errorf("cannot parse config: %w", err)
		log.Fatal(e)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("Startup app")

	pool, err := db.New(cfg)
	if err != nil {
		logger.Error("connect to database failed: ", err)
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("Connections pool to database initialized")

	if err = migrateUp(cfg); err != nil {
		logger.Error("migrate up failed: ", err)
		return
	}
	logger.Info("Database migrated successfully")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/user", rest.CreateUser(logger, pool))
	router.Delete("/user/{userId}", rest.DeleteUser(logger, pool))
	router.Patch("/user/{userId}", rest.UpdateUser(logger, pool))
	router.Get("/user/*", rest.GetUsers(logger, pool))

	logger.Info("Router initialized successfully")
	logger.Info("Starting server", slog.String("address", cfg.ServerHostAddress))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.ServerHostAddress,
		Handler: router,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("Server stopped")
			} else {
				logger.Error("error during server shutdown", err)
			}
		}
	}()

	logger.Info("Server started")

	<-done
	logger.Info("Stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", err)
		return
	}
}

func migrateUp(cfg *config.Config) error {
	connStr := fmt.Sprintf(
		"pgx5://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
