package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	config "scs-user/config"
	"scs-user/internal/models"
	"scs-user/internal/server"
	"scs-user/pkg/db"
	kafka_client "scs-user/pkg/kafka"
	"scs-user/pkg/logger"
	"strings"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
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

	// Initialize Kafka producer
	producer := startKafkaProducer("user.created", &cfg, appLogger)

	// Test sending a Kafka message after producer initialization
	ctx := context.Background()
	err = producer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("test-key"),
		Value: []byte("Hello, Kafka!"),
	})
	if err != nil {
		appLogger.Errorf("Failed to send Kafka message: %v", err)
	} else {
		appLogger.Info("Kafka message sent successfully")
	}
	// Initialize the server
	s := server.NewServer(&cfg, psqlDb, appLogger, producer)

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
func startKafkaProducer(topic string, cfg *config.Config, logger *logger.ApiLogger) *kafka_client.Producer {
	// Initialize Kafka producer
	fmt.Println("topic:", topic)
	kafkaCfg := kafka_client.Config{
		Brokers: strings.Split(cfg.Kafka.Brokers, ","),
		Topic:   topic,
	}
	producerCfg := kafka_client.ProducerConfig{
		BatchSize:    1,
		BatchTimeout: 100,   // In milliseconds
		Async:        false, // Set to false for immediate delivery
	}
	producer := kafka_client.NewProducer(&kafkaCfg, &producerCfg)
	return producer
}
