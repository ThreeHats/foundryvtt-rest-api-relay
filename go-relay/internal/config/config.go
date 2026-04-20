package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Core
	Port   int
	Env    string // "development" or "production"
	DBType string // "sqlite", "postgres"
	DBUrl  string

	// Redis
	RedisURL     string
	RedisEnabled bool

	// Auth
	AdminEmail               string
	AdminPassword            string
	CredentialsEncryptionKey string
	DisableRegistration      bool

	// Limits
	MonthlyRequestLimit          int
	AllowHeadless                bool
	MaxHeadlessSessions          int
	MaxInteractiveSessionsPerKey int
	HeadlessInactiveTimeoutSecs  int // env: HEADLESS_SESSION_TIMEOUT (alias: HEADLESS_INACTIVE_TIMEOUT), default 600; 0 = never time out
	MaxJSONBodySizeMB            int // env: MAX_JSON_BODY_SIZE_MB, default 250
	MaxUploadSizeMB              int // env: MAX_UPLOAD_SIZE_MB, default 250

	// Cache
	AuthCacheTTLSeconds int

	// WebSocket
	WSPingIntervalMs        int
	ClientCleanupIntervalMs int

	// Email
	SMTPHost   string
	SMTPPort   int
	SMTPUser   string
	SMTPPass   string
	SMTPFrom   string
	SMTPSecure bool

	// Stripe
	StripeSecretKey     string
	StripeWebhookSecret string
	StripePriceID       string
	StripePortalURL     string
	FrontendURL         string

	// Infrastructure
	FlyAllocID      string
	FlyInternalPort string
	AppName         string

	// Browser/Logging
	ChromePath              string
	LogLevel                string
	CaptureBrowserConsole   string
	BrowserLogRetentionDays int    // days to keep browser log files (0 = use default 3)
	DataDir                 string // absolute path to data directory (e.g. browser logs, sqlite)
	ChromeUserDataDir       string // persistent Chrome profile dir for V8 bytecode + HTTP cache
	ChromeJSHeapMB          int    // max V8 old-space heap in MB
	ChromeWindowWidth       int    // headless viewport width
	ChromeWindowHeight      int    // headless viewport height (≥768 for Foundry)
	ChromeEnableSHM         bool   // allow Chrome to use /dev/shm (needs ≥256MB shm)
	ChromeGPUMode           string // rendering backend: auto|gpu|xvfb|swiftshader

	// Admin dashboard
	AdminJWTSecret             string
	AdminAllowedIPs            string // comma-separated CIDRs/IPs; empty = allow all
	AdminInternalPort          int    // if >0, serve admin on this internal port (Fly private network)
	AdminSessionMaxHours       int
	AdminAuditLogRetentionDays int
	AdminSecureCookies         bool   // env: ADMIN_SECURE_COOKIES; enables Secure flag + __Host- prefix (requires HTTPS)

	// Log retention (days)
	ConnectionLogRetentionDays    int // env: CONNECTION_LOG_RETENTION_DAYS, default 90
	RemoteRequestLogRetentionDays int // env: REMOTE_REQUEST_LOG_RETENTION_DAYS, default 30
	ModuleEventLogRetentionDays   int // env: MODULE_EVENT_LOG_RETENTION_DAYS, default 7
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	cfg := &Config{
		Port:   getEnvInt("PORT", 3010),
		Env:    getEnv("APP_ENV", getEnv("NODE_ENV", "development")),
		DBType: getEnv("DB_TYPE", "sqlite"),
		DBUrl:  getEnv("DATABASE_URL", ""),

		RedisURL:     getEnv("REDIS_URL", ""),
		RedisEnabled: getEnvBool("REDIS_ENABLED", true),

		AdminEmail:               getEnv("ADMIN_EMAIL", ""),
		AdminPassword:            getEnv("ADMIN_PASSWORD", ""),
		CredentialsEncryptionKey: getEnv("CREDENTIALS_ENCRYPTION_KEY", ""),
		DisableRegistration:      getEnvBool("DISABLE_REGISTRATION", false),

		MonthlyRequestLimit:          getEnvIntWithFallback("MONTHLY_REQUEST_LIMIT", "FREE_API_REQUESTS_LIMIT", 0),
		AllowHeadless:                getEnvBool("ALLOW_HEADLESS", true),
		MaxHeadlessSessions:          getEnvInt("MAX_HEADLESS_SESSIONS", 1),
		MaxInteractiveSessionsPerKey: getEnvInt("MAX_INTERACTIVE_SESSIONS_PER_KEY", 3),
		HeadlessInactiveTimeoutSecs:  getEnvIntWithFallback("HEADLESS_SESSION_TIMEOUT", "HEADLESS_INACTIVE_TIMEOUT", 600),
		MaxJSONBodySizeMB:            getEnvInt("MAX_JSON_BODY_SIZE_MB", 250),
		MaxUploadSizeMB:              getEnvInt("MAX_UPLOAD_SIZE_MB", 250),

		AuthCacheTTLSeconds: getEnvInt("AUTH_CACHE_TTL_SECONDS", 30),

		WSPingIntervalMs:        getEnvInt("WEBSOCKET_PING_INTERVAL_MS", 20000),
		ClientCleanupIntervalMs: getEnvInt("CLIENT_CLEANUP_INTERVAL_MS", 15000),

		SMTPHost:   getEnv("SMTP_HOST", ""),
		SMTPPort:   getEnvInt("SMTP_PORT", 587),
		SMTPUser:   getEnv("SMTP_USER", ""),
		SMTPPass:   getEnv("SMTP_PASS", ""),
		SMTPFrom:   getEnv("SMTP_FROM", "noreply@foundryvtt-relay.com"),
		SMTPSecure: getEnvBool("SMTP_SECURE", false),

		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		StripePriceID:       getEnv("STRIPE_PRICE_ID", ""),
		StripePortalURL:     getEnv("STRIPE_PORTAL_URL", ""),
		FrontendURL:         getEnv("FRONTEND_URL", "https://foundryrestapi.com"),

		FlyAllocID:      getEnv("FLY_ALLOC_ID", "local"),
		FlyInternalPort: getEnv("FLY_INTERNAL_PORT", ""),
		AppName:         getEnv("APP_NAME", ""),

		ChromePath:              getEnv("PUPPETEER_EXECUTABLE_PATH", ""),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
		CaptureBrowserConsole:   getEnv("CAPTURE_BROWSER_CONSOLE", ""),
		BrowserLogRetentionDays: getEnvInt("BROWSER_LOG_RETENTION_DAYS", 3),
		ChromeUserDataDir:       getEnv("CHROME_USER_DATA_DIR", ""),
		ChromeJSHeapMB:          getEnvInt("CHROME_JS_HEAP_MB", 2048),
		ChromeWindowWidth:       getEnvInt("CHROME_WINDOW_WIDTH", 1280),
		ChromeWindowHeight:      getEnvInt("CHROME_WINDOW_HEIGHT", 800),
		ChromeEnableSHM:         getEnvBool("CHROME_ENABLE_SHM", false),
		ChromeGPUMode:           getEnv("CHROME_GPU_MODE", "auto"),

		AdminJWTSecret:             getEnv("ADMIN_JWT_SECRET", ""),
		AdminAllowedIPs:            getEnv("ADMIN_ALLOWED_IPS", ""),
		AdminInternalPort:          getEnvInt("ADMIN_INTERNAL_PORT", 0),
		AdminSessionMaxHours:       getEnvInt("ADMIN_SESSION_MAX_HOURS", 8),
		AdminAuditLogRetentionDays: getEnvInt("ADMIN_AUDIT_LOG_RETENTION_DAYS", 90),
		AdminSecureCookies:         getEnvBool("ADMIN_SECURE_COOKIES", false),

		ConnectionLogRetentionDays:    getEnvInt("CONNECTION_LOG_RETENTION_DAYS", 90),
		RemoteRequestLogRetentionDays: getEnvInt("REMOTE_REQUEST_LOG_RETENTION_DAYS", 30),
		ModuleEventLogRetentionDays:   getEnvInt("MODULE_EVENT_LOG_RETENTION_DAYS", 7),
	}

	// Resolve DataDir to an absolute path so it's consistent regardless of CWD.
	// Default is "../data" so that when the server runs from go-relay/ the data
	// directory is at the project root. Docker sets DATA_DIR=/app/data explicitly.
	dataDir, _ := filepath.Abs(getEnv("DATA_DIR", "../data"))
	cfg.DataDir = dataDir

	// Disable Redis for local database modes
	if cfg.DBType == "sqlite" {
		cfg.RedisEnabled = false
	}

	return cfg
}

// Validate checks that required configuration is present and returns an error if not.
// Call this after Load() to fail fast on misconfiguration.
func (c *Config) Validate() error {
	// Credentials encryption is required for all database backends.
	if c.CredentialsEncryptionKey == "" {
		return fmt.Errorf("CREDENTIALS_ENCRYPTION_KEY must be set — generate one with: openssl rand -hex 32")
	}
	// Admin JWT secret is required for all backends — auto-generation is not supported.
	if c.AdminJWTSecret == "" {
		return fmt.Errorf("ADMIN_JWT_SECRET must be set — generate one with: openssl rand -base64 32")
	}
	return nil
}

// configForLog is an alias used to print Config without triggering the String() method recursively.
type configForLog Config

// String returns a safe representation of the config with secrets redacted.
// This prevents accidental secret exposure in logs or debug output.
func (c Config) String() string {
	redact := func(s string) string {
		if s == "" {
			return ""
		}
		return "[redacted]"
	}
	c.AdminPassword = redact(c.AdminPassword)
	c.CredentialsEncryptionKey = redact(c.CredentialsEncryptionKey)
	c.AdminJWTSecret = redact(c.AdminJWTSecret)
	c.StripeSecretKey = redact(c.StripeSecretKey)
	c.StripeWebhookSecret = redact(c.StripeWebhookSecret)
	c.SMTPPass = redact(c.SMTPPass)
	return fmt.Sprintf("%+v", configForLog(c))
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

// getEnvIntWithFallback tries the primary key, then the fallback key, then the default.
func getEnvIntWithFallback(key, fallbackKey string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	if val := os.Getenv(fallbackKey); val != "" {
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
