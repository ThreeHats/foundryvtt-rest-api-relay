package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/cron"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/server"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const version = "2.2.1"

func main() {
	// Load .env file (silent fail if not found)
	godotenv.Load()      // .env in current directory
	godotenv.Load("../.env") // .env in parent (worktree root)

	// Re-initialize rate limiters now that .env is loaded
	middleware.InitRateLimiters()

	// Load configuration
	cfg := config.Load()

	// Setup logging
	setupLogging(cfg.LogLevel)

	log.Info().Str("version", version).Msg("Starting FoundryVTT REST API Relay")

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Create admin user if configured
	if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
		if err := db.CreateAdminUser(cfg.AdminEmail, cfg.AdminPassword); err != nil {
			log.Warn().Err(err).Msg("Failed to create admin user")
		}
	}

	// Initialize Redis (optional)
	var redisClient *config.RedisClient
	if cfg.RedisEnabled {
		redisClient, err = config.NewRedisClient(cfg)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to initialize Redis, continuing without it")
		} else {
			defer redisClient.Close()
		}
	}

	// Start cron jobs
	scheduler := cron.NewScheduler(db, redisClient)
	scheduler.Start()
	defer scheduler.Stop()

	// Create and start HTTP server
	srv := server.New(cfg, db, redisClient, version)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      srv.Router(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second, // Long for SSE/file uploads
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		log.Info().Int("port", cfg.Port).Msg("Server listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server stopped")
}

func setupLogging(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
