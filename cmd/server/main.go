package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	config "scs-user/config"
	"scs-user/internal/models"
	"scs-user/internal/server"
	"scs-user/pkg/db"
	"scs-user/pkg/logger"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func main() {
	// Load configuration from config file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	//Init logger
	appLogger := logger.GetLogger()
	appLogger.InitLogger(&cfg)
	appLogger.Infof("LogLevel: %s, Mode: %s", cfg.Logger.Level, cfg.Server.Mode)

	//Init db
	psqlDb, err := db.NewGormDB(&cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Info("Postgres connected")
	}

	// Auto-migrate models
	err = psqlDb.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		appLogger.Fatalf("Database migration failed: %s", err)
	}
	// Initialize the server
	s := server.NewServer(&cfg, psqlDb, appLogger)

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := s.Run(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatalf("Error starting server: %v", err)
		}
	}()

	// Block until a signal is received
	<-quit

	appLogger.Info("Shutting down the server...")

	// Create a separate, timeout context for the server shutdown
	serverShutdownCtx, serverShutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer serverShutdownCancel()

	// Shut down the Echo server
	if err := s.Shutdown(serverShutdownCtx); err != nil {
		appLogger.Errorf("Server shutdown failed: %v", err)
	}

	appLogger.Info("Server and consumer stopped.")
}
