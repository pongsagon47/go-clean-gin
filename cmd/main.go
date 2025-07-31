package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-clean-gin/config"
	"go-clean-gin/internal/container"
	"go-clean-gin/internal/router"
	"go-clean-gin/pkg/database"
	"go-clean-gin/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("env", cfg.Env),
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// // Run migrations
	// if err := database.RunMigrations(db); err != nil {
	// 	logger.Fatal("Failed to run migrations", zap.Error(err))
	// }

	// // Seed data (only in development)
	// if cfg.Env == "development" {
	// 	if err := database.SeedData(db, ""); err != nil {
	// 		logger.Warn("Failed to seed data", zap.Error(err))
	// 	}
	// }

	// Initialize dependency injection container
	containerInstance := container.NewContainer(cfg, db)

	// Setup routes
	routerInstance := router.SetupRouter(containerInstance)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      routerInstance,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting",
			zap.String("address", server.Addr),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started successfully", zap.String("address", server.Addr))

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		}
	}

	logger.Info("Server exited")
}
