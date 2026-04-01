package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Core
	Port     int
	Env      string // "development" or "production"
	DBType   string // "memory", "sqlite", "postgres"
	DBUrl    string

	// Redis
	RedisURL     string
	RedisEnabled bool

	// Auth
	AdminEmail              string
	AdminPassword           string
	CredentialsEncryptionKey string
	DisableRegistration     bool

	// Limits
	FreeAPIRequestsLimit int
	DailyRequestLimit    int
	MaxHeadlessSessions  int
	MaxSheetSessionsPerKey int

	// WebSocket
	WSPingIntervalMs      int
	ClientCleanupIntervalMs int

	// Email
	SMTPHost   string
	SMTPPort   int
	SMTPUser   string
	SMTPPass   string
	SMTPFrom   string
	SMTPSecure bool

	// Stripe
	StripeSecretKey      string
	StripeWebhookSecret  string
	StripePriceID        string
	StripePortalURL      string
	FrontendURL          string

	// Infrastructure
	FlyAllocID      string
	FlyInternalPort string
	AppName         string

	// Browser/Logging
	ChromePath            string
	LogLevel              string
	CaptureBrowserConsole string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	cfg := &Config{
		Port:     getEnvInt("PORT", 3010),
		Env:      getEnv("NODE_ENV", "development"),
		DBType:   getEnv("DB_TYPE", "postgres"),
		DBUrl:    getEnv("DATABASE_URL", ""),

		RedisURL:     getEnv("REDIS_URL", ""),
		RedisEnabled: getEnvBool("REDIS_ENABLED", true),

		AdminEmail:              getEnv("ADMIN_EMAIL", ""),
		AdminPassword:           getEnv("ADMIN_PASSWORD", ""),
		CredentialsEncryptionKey: getEnv("CREDENTIALS_ENCRYPTION_KEY", ""),
		DisableRegistration:     getEnvBool("DISABLE_REGISTRATION", false),

		FreeAPIRequestsLimit: getEnvInt("FREE_API_REQUESTS_LIMIT", 100),
		DailyRequestLimit:    getEnvInt("DAILY_REQUEST_LIMIT", 1000),
		MaxHeadlessSessions:  getEnvInt("MAX_HEADLESS_SESSIONS", 1),
		MaxSheetSessionsPerKey: getEnvInt("MAX_SHEET_SESSIONS_PER_KEY", 3),

		WSPingIntervalMs:      getEnvInt("WEBSOCKET_PING_INTERVAL_MS", 20000),
		ClientCleanupIntervalMs: getEnvInt("CLIENT_CLEANUP_INTERVAL_MS", 15000),

		SMTPHost:   getEnv("SMTP_HOST", ""),
		SMTPPort:   getEnvInt("SMTP_PORT", 587),
		SMTPUser:   getEnv("SMTP_USER", ""),
		SMTPPass:   getEnv("SMTP_PASS", ""),
		SMTPFrom:   getEnv("SMTP_FROM", "noreply@foundryvtt-relay.com"),
		SMTPSecure: getEnvBool("SMTP_SECURE", false),

		StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		StripePriceID:       getEnv("STRIPE_PRICE_ID", ""),
		StripePortalURL:     getEnv("STRIPE_PORTAL_URL", ""),
		FrontendURL:         getEnv("FRONTEND_URL", ""),

		FlyAllocID:      getEnv("FLY_ALLOC_ID", "local"),
		FlyInternalPort: getEnv("FLY_INTERNAL_PORT", ""),
		AppName:         getEnv("APP_NAME", ""),

		ChromePath:            getEnv("PUPPETEER_EXECUTABLE_PATH", ""),
		LogLevel:              getEnv("LOG_LEVEL", "info"),
		CaptureBrowserConsole: getEnv("CAPTURE_BROWSER_CONSOLE", ""),
	}

	// Disable Redis for local database modes
	if cfg.DBType == "memory" || cfg.DBType == "sqlite" {
		cfg.RedisEnabled = false
	}

	return cfg
}

// IsLocalMode returns true if using memory or sqlite database.
func (c *Config) IsLocalMode() bool {
	return c.DBType == "memory" || c.DBType == "sqlite"
}

// IsProduction returns true if running in production.
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// InstanceID returns the instance identifier for multi-instance deployments.
func (c *Config) InstanceID() string {
	return c.FlyAllocID
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		return strings.EqualFold(val, "true") || val == "1"
	}
	return fallback
}
