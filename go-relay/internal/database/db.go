package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// DB wraps the database connection and provides model operations.
type DB struct {
	sqlDB    *sqlx.DB
	memory   *MemoryStore
	isMemory bool
	dbType   string
}

// New creates a new database connection based on configuration.
func New(cfg *config.Config) (*DB, error) {
	db := &DB{dbType: cfg.DBType}

	switch cfg.DBType {
	case "memory":
		log.Info().Msg("Using in-memory database")
		db.memory = NewMemoryStore()
		db.isMemory = true
		return db, nil

	case "sqlite":
		log.Info().Msg("Using SQLite database")
		dataDir := filepath.Join(".", "data")
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("create data dir: %w", err)
		}
		dbPath := filepath.Join(dataDir, "relay.db")
		log.Info().Str("path", dbPath).Msg("SQLite database path")

		sqlDB, err := sqlx.Connect("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL")
		if err != nil {
			return nil, fmt.Errorf("connect sqlite: %w", err)
		}
		// SQLite needs single writer — limit connections to avoid SQLITE_BUSY
		sqlDB.SetMaxOpenConns(1)
		// Use Unsafe so sqlx doesn't fail on column/struct mismatches.
		// Sequelize-created SQLite uses camelCase columns; Go struct tags use snake_case.
		// MapperFunc normalizes both to lowercase for matching.
		sqlDB.Unsafe()
		sqlDB.MapperFunc(func(s string) string {
			return strings.ToLower(strings.ReplaceAll(s, "_", ""))
		})
		db.sqlDB = sqlDB
		return db, nil

	default: // postgres
		if cfg.DBUrl == "" {
			return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
		}
		log.Info().Msg("Using PostgreSQL database")

		sqlDB, err := sqlx.Connect("pgx", cfg.DBUrl)
		if err != nil {
			return nil, fmt.Errorf("connect postgres: %w", err)
		}
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		// Sequelize-created Postgres uses camelCase columns; normalize for matching.
		sqlDB.Unsafe()
		sqlDB.MapperFunc(func(s string) string {
			return strings.ToLower(strings.ReplaceAll(s, "_", ""))
		})
		db.sqlDB = sqlDB
		return db, nil
	}
}

// Close closes the database connection.
func (db *DB) Close() error {
	if db.sqlDB != nil {
		return db.sqlDB.Close()
	}
	return nil
}

// Migrate creates tables if they don't exist.
func (db *DB) Migrate() error {
	if db.isMemory {
		return nil
	}

	ctx := context.Background()

	if db.dbType == "sqlite" {
		return db.migrateSQLite(ctx)
	}
	return db.migratePostgres(ctx)
}

func (db *DB) migrateSQLite(ctx context.Context) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS Users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			"apiKey" TEXT NOT NULL UNIQUE,
			"requestsThisMonth" INTEGER DEFAULT 0,
			"requestsToday" INTEGER DEFAULT 0,
			"lastRequestDate" TEXT,
			"stripeCustomerId" TEXT,
			"subscriptionStatus" TEXT DEFAULT 'free',
			"subscriptionId" TEXT,
			"subscriptionEndsAt" TEXT,
			"maxHeadlessSessions" INTEGER,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS PasswordResetTokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"tokenHash" TEXT NOT NULL UNIQUE,
			"expiresAt" TEXT NOT NULL,
			used INTEGER DEFAULT 0,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS ApiKeys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			key TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			"scopedClientId" TEXT,
			"scopedUserId" TEXT,
			"dailyLimit" INTEGER,
			"requestsToday" INTEGER DEFAULT 0,
			"lastRequestDate" TEXT,
			"foundryUrl" TEXT,
			"foundryUsername" TEXT,
			"encryptedFoundryPassword" TEXT,
			"passwordIv" TEXT,
			"passwordAuthTag" TEXT,
			"expiresAt" TEXT,
			enabled INTEGER DEFAULT 1,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
	}

	for _, m := range migrations {
		if _, err := db.sqlDB.ExecContext(ctx, m); err != nil {
			return fmt.Errorf("sqlite migration: %w", err)
		}
	}
	return nil
}

func (db *DB) migratePostgres(ctx context.Context) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS "Users" (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			api_key VARCHAR(64) NOT NULL UNIQUE,
			requests_this_month INTEGER DEFAULT 0,
			requests_today INTEGER DEFAULT 0,
			last_request_date DATE,
			stripe_customer_id VARCHAR(255),
			subscription_status VARCHAR(50) DEFAULT 'free',
			subscription_id VARCHAR(255),
			subscription_ends_at TIMESTAMPTZ,
			max_headless_sessions INTEGER,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "PasswordResetTokens" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			token_hash VARCHAR(255) NOT NULL UNIQUE,
			expires_at TIMESTAMPTZ NOT NULL,
			used BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "ApiKeys" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			key VARCHAR(64) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			scoped_client_id VARCHAR(255),
			scoped_user_id VARCHAR(255),
			daily_limit INTEGER,
			requests_today INTEGER DEFAULT 0,
			last_request_date DATE,
			foundry_url TEXT,
			foundry_username VARCHAR(255),
			encrypted_foundry_password TEXT,
			password_iv VARCHAR(255),
			password_auth_tag VARCHAR(255),
			expires_at TIMESTAMPTZ,
			enabled BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
	}

	for _, m := range migrations {
		if _, err := db.sqlDB.ExecContext(ctx, m); err != nil {
			return fmt.Errorf("postgres migration: %w", err)
		}
	}

	// Add columns that may be missing from Sequelize-created tables
	alterMigrations := []string{
		`ALTER TABLE "Users" ADD COLUMN IF NOT EXISTS "maxHeadlessSessions" INTEGER`,
	}
	for _, m := range alterMigrations {
		_, _ = db.sqlDB.ExecContext(ctx, m)
	}

	// Rename snake_case columns to camelCase to match Sequelize convention.
	// These are safe to run repeatedly — they no-op if already renamed.
	renames := []struct{ table, from, to string }{
		// ApiKeys table
		{"ApiKeys", "user_id", "userId"},
		{"ApiKeys", "scoped_client_id", "scopedClientId"},
		{"ApiKeys", "scoped_user_id", "scopedUserId"},
		{"ApiKeys", "daily_limit", "dailyLimit"},
		{"ApiKeys", "requests_today", "requestsToday"},
		{"ApiKeys", "last_request_date", "lastRequestDate"},
		{"ApiKeys", "foundry_url", "foundryUrl"},
		{"ApiKeys", "foundry_username", "foundryUsername"},
		{"ApiKeys", "encrypted_foundry_password", "encryptedFoundryPassword"},
		{"ApiKeys", "password_iv", "passwordIv"},
		{"ApiKeys", "password_auth_tag", "passwordAuthTag"},
		{"ApiKeys", "expires_at", "expiresAt"},
		{"ApiKeys", "created_at", "createdAt"},
		{"ApiKeys", "updated_at", "updatedAt"},
		// PasswordResetTokens table
		{"PasswordResetTokens", "user_id", "userId"},
		{"PasswordResetTokens", "token_hash", "tokenHash"},
		{"PasswordResetTokens", "expires_at", "expiresAt"},
		{"PasswordResetTokens", "created_at", "createdAt"},
		{"PasswordResetTokens", "updated_at", "updatedAt"},
	}
	for _, r := range renames {
		_, _ = db.sqlDB.ExecContext(ctx, fmt.Sprintf(
			`ALTER TABLE "%s" RENAME COLUMN %s TO "%s"`, r.table, r.from, r.to))
	}

	return nil
}

// CreateAdminUser creates the initial admin user if it doesn't exist.
func (db *DB) CreateAdminUser(email, password string) error {
	ctx := context.Background()

	// Check if user exists
	existing, err := db.UserStore().FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		log.Info().Str("email", email).Msg("Admin user already exists")
		return nil
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Email:              email,
		Password:           string(hash),
		APIKey:             model.GenerateAPIKey(),
		RequestsThisMonth:  0,
		RequestsToday:      0,
		SubscriptionStatus: sql.NullString{String: "free", Valid: true},
	}

	if err := db.UserStore().Create(ctx, user); err != nil {
		return err
	}

	log.Info().Str("email", email).Msg("Admin user created")
	return nil
}

// UserStore returns the user store for this database.
func (db *DB) UserStore() model.UserStore {
	if db.isMemory {
		return db.memory.Users()
	}
	return &model.SQLUserStore{DB: db.sqlDB, DBType: db.dbType}
}

// ApiKeyStore returns the API key store for this database.
func (db *DB) ApiKeyStore() model.ApiKeyStore {
	if db.isMemory {
		return db.memory.ApiKeys()
	}
	return &model.SQLApiKeyStore{DB: db.sqlDB, DBType: db.dbType}
}

// PasswordResetTokenStore returns the token store for this database.
func (db *DB) PasswordResetTokenStore() model.PasswordResetTokenStore {
	if db.isMemory {
		return db.memory.PasswordResetTokens()
	}
	return &model.SQLPasswordResetTokenStore{DB: db.sqlDB, DBType: db.dbType}
}
