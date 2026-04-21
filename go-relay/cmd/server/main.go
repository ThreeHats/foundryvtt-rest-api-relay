package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/cron"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/metrics"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/server"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const version = "3.0.2"

func main() {
	// Load .env file (silent fail if not found)
	godotenv.Load()      // .env in current directory
	godotenv.Load("../.env") // .env in parent (worktree root)

	// Graceful shutdown context — created early so background goroutines can respect it
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Re-initialize rate limiters now that .env is loaded
	middleware.InitRateLimiters(ctx)
	middleware.InitAdminRateLimiters(ctx)
	middleware.InitAdminLockout()

	// Load configuration
	cfg := config.Load()

	// Validate required configuration (fails fast on misconfiguration)
	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid configuration")
	}

	middleware.SetAuthCacheTTL(time.Duration(cfg.AuthCacheTTLSeconds) * time.Second)

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

	// Restore cumulative metrics from the last persisted snapshot.
	// Rolling time-window stats (per-min/hour/day) always reset on restart.
	if ep, usr, errs, snapErr := db.LoadMetricsSnapshot(); snapErr != nil {
		log.Warn().Err(snapErr).Msg("Failed to load metrics snapshot")
	} else if len(ep) > 0 || errs > 0 {
		metrics.Global.Import(ep, usr, errs)
		log.Info().Int("endpoints", len(ep)).Msg("Restored metrics snapshot")
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

	// Flush batched request counters every 500ms instead of one DB write per request.
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				middleware.RequestCounts.Flush(context.Background(), db.UserStore())
			case <-ctx.Done():
				middleware.RequestCounts.Flush(context.Background(), db.UserStore())
				return
			}
		}
	}()

	// Persist cumulative metrics every 5 minutes so restarts don't lose the
	// endpoint breakdown and per-user totals. Rolling time-window stats reset
	// on restart by design and are not included in the snapshot.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ep, usr, errs := metrics.Global.Export()
				if saveErr := db.SaveMetricsSnapshot(ep, usr, errs); saveErr != nil {
					log.Warn().Err(saveErr).Msg("Failed to save metrics snapshot")
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start cron jobs
	scheduler := cron.NewScheduler(db, cfg, redisClient, cfg.DataDir, cfg.BrowserLogRetentionDays)
	scheduler.Start()
	defer scheduler.Stop()

	// Create and start HTTP server
	srv := server.New(ctx, cfg, db, redisClient, version)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      srv.Router(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second, // Long for SSE/file uploads
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Int("port", cfg.Port).Msg("Server listening on all interfaces (0.0.0.0)")
		log.Info().Str("url", fmt.Sprintf("http://localhost:%d", cfg.Port)).Msg("Local URL")
		for _, ip := range lanIPs() {
			log.Info().Str("url", fmt.Sprintf("http://%s:%d", ip, cfg.Port)).Msg("LAN URL (reachable from other devices on your network)")
		}
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Optional second HTTP server bound only to the Fly internal IPv6 address
	// (fly-local-6pn) for private admin access. Only reachable via flyctl proxy
	// or sibling apps on the same Fly organization's private network.
	var internalServer *http.Server
	if cfg.AdminInternalPort > 0 {
		internalServer = &http.Server{
			Addr:         fmt.Sprintf("[::]:%d", cfg.AdminInternalPort),
			Handler:      srv.Router(),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 600 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		go func() {
			log.Info().Int("port", cfg.AdminInternalPort).Msg("Admin internal server listening (Fly private network)")
			if err := internalServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error().Err(err).Msg("Admin internal server failed")
			}
		}()
	}

	<-ctx.Done()
	log.Info().Msg("Shutting down gracefully...")

	// Persist metrics before exit so the next startup can restore them.
	ep, usr, errs := metrics.Global.Export()
	if saveErr := db.SaveMetricsSnapshot(ep, usr, errs); saveErr != nil {
		log.Warn().Err(saveErr).Msg("Failed to save metrics snapshot on shutdown")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
	if internalServer != nil {
		_ = internalServer.Shutdown(shutdownCtx)
	}

	log.Info().Msg("Server stopped")
}

// lanIPs returns non-loopback IPv4 addresses on up interfaces, so operators
// can see at startup which URLs other devices on their network can use.
func lanIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var ips []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip4 := ip.To4()
			if ip4 == nil {
				continue
			}
			ips = append(ips, ip4.String())
		}
	}
	return ips
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
