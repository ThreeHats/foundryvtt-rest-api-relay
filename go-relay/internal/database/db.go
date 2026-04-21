package database

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	sqlDB  *sqlx.DB
	dbType string
}

// New creates a new database connection based on configuration.
func New(cfg *config.Config) (*DB, error) {
	db := &DB{dbType: cfg.DBType}

	switch cfg.DBType {
	case "sqlite":
		log.Info().Msg("Using SQLite database")
		dbPath := os.Getenv("SQLITE_PATH")
		if dbPath == "" {
			// Project root data/ directory (go-relay runs from go-relay/, so ../data/)
			dataDir := filepath.Join("..", "data")
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				return nil, fmt.Errorf("create data dir: %w", err)
			}
			dbPath = filepath.Join(dataDir, "relay.db")
		}
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
		sqlDB = sqlDB.Unsafe()
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
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)
		// Sequelize-created Postgres uses camelCase columns; normalize for matching.
		sqlDB = sqlDB.Unsafe()
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
			"apiKeyHash" TEXT NOT NULL UNIQUE,
			"requestsThisMonth" INTEGER DEFAULT 0,
			"requestsToday" INTEGER DEFAULT 0,
			"lastRequestDate" TEXT,
			"stripeCustomerId" TEXT,
			"subscriptionStatus" TEXT DEFAULT 'free',
			"subscriptionId" TEXT,
			"subscriptionEndsAt" TEXT,
			"maxHeadlessSessions" INTEGER,
			role TEXT DEFAULT 'user',
			disabled INTEGER DEFAULT 0,
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
			"monthlyLimit" INTEGER,
			"requestsThisMonth" INTEGER DEFAULT 0,
			"lastResetDate" TEXT,
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
		`CREATE TABLE IF NOT EXISTS ConnectionTokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"tokenHash" TEXT NOT NULL UNIQUE,
			name TEXT DEFAULT '',
			"allowedIps" TEXT DEFAULT '',
			"allowedTargetClients" TEXT DEFAULT '',
			"remoteScopes" TEXT DEFAULT '',
			"remoteRequestsPerHour" INTEGER DEFAULT 0,
			source TEXT DEFAULT 'dashboard',
			"lastUsedAt" TEXT,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS PairingCodes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			code TEXT NOT NULL UNIQUE,
			"clientId" TEXT,
			"allowedTargetClients" TEXT DEFAULT '',
			"remoteScopes" TEXT DEFAULT '',
			"expiresAt" TEXT NOT NULL,
			used INTEGER DEFAULT 0,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS ConnectionLogs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"clientId" TEXT NOT NULL,
			"tokenName" TEXT,
			"ipAddress" TEXT,
			"userAgent" TEXT,
			"worldId" TEXT,
			"worldTitle" TEXT,
			"systemId" TEXT,
			"foundryVersion" TEXT,
			"metadataMatch" INTEGER DEFAULT 1,
			flagged INTEGER DEFAULT 0,
			"flagReason" TEXT,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS Credentials (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			name TEXT NOT NULL,
			"foundryUrl" TEXT NOT NULL,
			"foundryUsername" TEXT NOT NULL,
			"encryptedFoundryPassword" TEXT NOT NULL,
			"passwordIv" TEXT NOT NULL,
			"passwordAuthTag" TEXT NOT NULL,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS KnownClients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"clientId" TEXT NOT NULL,
			"worldId" TEXT,
			"worldTitle" TEXT,
			"systemId" TEXT,
			"systemTitle" TEXT,
			"systemVersion" TEXT,
			"foundryVersion" TEXT,
			"customName" TEXT,
			"lastSeenAt" TEXT,
			"isOnline" INTEGER DEFAULT 0,
			"autoStartOnRemoteRequest" INTEGER DEFAULT 0,
			"credentialId" INTEGER,
			"allowedTargetClients" TEXT DEFAULT '',
			"remoteScopes" TEXT DEFAULT '',
			"remoteRequestsPerHour" INTEGER DEFAULT 0,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now')),
			UNIQUE("userId", "clientId")
		)`,
		`CREATE TABLE IF NOT EXISTS KeyRequests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL UNIQUE,
			"appName" TEXT NOT NULL,
			"appDescription" TEXT DEFAULT '',
			"appUrl" TEXT DEFAULT '',
			"requestedScopes" TEXT NOT NULL,
			"requestedClientIds" TEXT DEFAULT '',
			"callbackUrl" TEXT,
			"suggestedMonthlyLimit" INTEGER,
			"suggestedExpiry" TEXT,
			status TEXT DEFAULT 'pending',
			"approvedKeyId" INTEGER,
			"approvedById" INTEGER,
			"exchangeCode" TEXT,
			"expiresAt" TEXT NOT NULL,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS NotificationSettings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL UNIQUE,
			"notifyOnConnect" INTEGER DEFAULT 1,
			"discordWebhookUrl" TEXT,
			"notifyEmail" TEXT,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS AdminAuditLogs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"adminUserId" INTEGER NOT NULL,
			action TEXT NOT NULL,
			"targetType" TEXT NOT NULL,
			"targetId" TEXT,
			details TEXT,
			"ipAddress" TEXT,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_audit_created ON AdminAuditLogs("createdAt")`,
		`CREATE INDEX IF NOT EXISTS idx_admin_audit_action ON AdminAuditLogs(action)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_audit_admin_user ON AdminAuditLogs("adminUserId")`,
		`CREATE TABLE IF NOT EXISTS JWTDenylist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"tokenJti" TEXT NOT NULL UNIQUE,
			"expiresAt" TEXT NOT NULL,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_jwt_denylist_jti ON JWTDenylist("tokenJti")`,
		`CREATE INDEX IF NOT EXISTS idx_jwt_denylist_expires ON JWTDenylist("expiresAt")`,
		`CREATE TABLE IF NOT EXISTS AlertSubscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"alertType" TEXT NOT NULL,
			channel TEXT NOT NULL,
			destination TEXT NOT NULL,
			enabled INTEGER DEFAULT 1,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS AlertConfig (
			id INTEGER PRIMARY KEY DEFAULT 1,
			"discordWebhookUrl" TEXT DEFAULT '',
			"emailDestination" TEXT DEFAULT '',
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS Sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"tokenHash" TEXT NOT NULL UNIQUE,
			"userAgent" TEXT,
			"ipAddress" TEXT,
			"expiresAt" TEXT NOT NULL,
			"lastUsedAt" TEXT,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON Sessions("tokenHash")`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON Sessions("userId")`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires ON Sessions("expiresAt")`,
		`CREATE TABLE IF NOT EXISTS RemoteRequestLogs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"sourceClientId" TEXT NOT NULL,
			"sourceTokenId" INTEGER NOT NULL,
			"targetClientId" TEXT NOT NULL,
			action TEXT NOT NULL,
			success INTEGER NOT NULL,
			"errorMessage" TEXT,
			"sourceIp" TEXT,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_remote_req_user ON RemoteRequestLogs("userId")`,
		`CREATE INDEX IF NOT EXISTS idx_remote_req_created ON RemoteRequestLogs("createdAt")`,
		`CREATE TABLE IF NOT EXISTS PairRequests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL UNIQUE,
			"worldId" TEXT NOT NULL,
			"worldTitle" TEXT NOT NULL DEFAULT '',
			"systemId" TEXT NOT NULL DEFAULT '',
			"systemTitle" TEXT NOT NULL DEFAULT '',
			"systemVersion" TEXT NOT NULL DEFAULT '',
			"foundryVersion" TEXT NOT NULL DEFAULT '',
			"requestedRemoteScopes" TEXT,
			"requestedTargetClients" TEXT,
			"upgradeOnly" INTEGER DEFAULT 0,
			status TEXT DEFAULT 'pending',
			"pairingCode" TEXT,
			"userId" INTEGER,
			"expiresAt" TEXT NOT NULL,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS KnownUsers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"knownClientId" INTEGER NOT NULL,
			"userId" TEXT NOT NULL,
			name TEXT NOT NULL,
			role INTEGER NOT NULL DEFAULT 1,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now')),
			UNIQUE("knownClientId", "userId"),
			FOREIGN KEY("knownClientId") REFERENCES KnownClients(id) ON DELETE CASCADE
		)`,
	}

	for _, m := range migrations {
		if _, err := db.sqlDB.ExecContext(ctx, m); err != nil {
			return fmt.Errorf("sqlite migration: %w", err)
		}
	}

	// Add columns that may be missing from existing tables
	sqliteAlterMigrations := []string{
		`ALTER TABLE ApiKeys ADD COLUMN scopes TEXT DEFAULT ''`,
		`ALTER TABLE ApiKeys ADD COLUMN "scopedClientIds" TEXT DEFAULT ''`,
		`ALTER TABLE Users ADD COLUMN "emailVerified" INTEGER DEFAULT 1`,
		`ALTER TABLE Users ADD COLUMN "verificationTokenHash" TEXT`,
		`ALTER TABLE Users ADD COLUMN "verificationTokenExpiresAt" TEXT`,
		`ALTER TABLE Users ADD COLUMN role TEXT DEFAULT 'user'`,
		`ALTER TABLE Users ADD COLUMN disabled INTEGER DEFAULT 0`,
		`ALTER TABLE Users ADD COLUMN "apiKeyRotationRequired" INTEGER DEFAULT 0`,
		`ALTER TABLE Users ADD COLUMN "apiKeyHash" TEXT`,
		// Drop the plaintext apiKey column. Users with apiKeyRotationRequired=1
		// will regenerate via the password-based regenerate-key endpoint and get
		// a fresh key whose hash is stored in apiKeyHash.
		`ALTER TABLE Users DROP COLUMN "apiKey"`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_api_key_hash ON Users("apiKeyHash") WHERE "apiKeyHash" IS NOT NULL`,
		// Connection token permissions for cross-world tunneling
		`ALTER TABLE ConnectionTokens ADD COLUMN "allowedTargetClients" TEXT DEFAULT ''`,
		`ALTER TABLE ConnectionTokens ADD COLUMN "remoteScopes" TEXT DEFAULT ''`,
		`ALTER TABLE ConnectionTokens ADD COLUMN "remoteRequestsPerHour" INTEGER DEFAULT 0`,
		`ALTER TABLE ConnectionTokens ADD COLUMN source TEXT DEFAULT 'dashboard'`,
		// Link each token to the KnownClient (world) it was paired for.
		// Enables the dashboard to group browsers under their world.
		`ALTER TABLE ConnectionTokens ADD COLUMN "clientId" TEXT DEFAULT ''`,
		// Pairing code can be bound to an existing clientId for "add browser" flows
		`ALTER TABLE PairingCodes ADD COLUMN "clientId" TEXT`,
		`ALTER TABLE PairingCodes ADD COLUMN "allowedTargetClients" TEXT DEFAULT ''`,
		`ALTER TABLE PairingCodes ADD COLUMN "remoteScopes" TEXT DEFAULT ''`,
		// KnownClients can opt-in to auto-start headless on incoming remote-request
		`ALTER TABLE KnownClients ADD COLUMN "autoStartOnRemoteRequest" INTEGER DEFAULT 0`,
		// Optional explicit link to a stored Credential. When set, the headless
		// auto-start flow uses this credential. When NULL, falls back to the
		// user's first credential (works for single-Foundry-server deployments).
		`ALTER TABLE KnownClients ADD COLUMN "credentialId" INTEGER`,
		// Cross-world tunneling settings: moved from per-token to per-world (KnownClient).
		// All tokens for a world inherit these permissions; enforcement reads from KnownClient.
		`ALTER TABLE KnownClients ADD COLUMN "allowedTargetClients" TEXT DEFAULT ''`,
		`ALTER TABLE KnownClients ADD COLUMN "remoteScopes" TEXT DEFAULT ''`,
		`ALTER TABLE KnownClients ADD COLUMN "remoteRequestsPerHour" INTEGER DEFAULT 0`,
		// Connection log browser name — the name of the connection token used to authenticate.
		`ALTER TABLE ConnectionLogs ADD COLUMN "tokenName" TEXT`,
		// PairRequests cross-world fields
		`ALTER TABLE PairRequests ADD COLUMN "requestedRemoteScopes" TEXT`,
		`ALTER TABLE PairRequests ADD COLUMN "requestedTargetClients" TEXT`,
		`ALTER TABLE PairRequests ADD COLUMN "upgradeOnly" INTEGER DEFAULT 0`,
		// Extended notification toggles for account-level NotificationSettings
		`ALTER TABLE NotificationSettings ADD COLUMN "notifyOnDisconnect" INTEGER DEFAULT 1`,
		`ALTER TABLE NotificationSettings ADD COLUMN "notifyOnMetadataMismatch" INTEGER DEFAULT 1`,
		`ALTER TABLE NotificationSettings ADD COLUMN "notifyOnSettingsChange" INTEGER DEFAULT 1`,
		`ALTER TABLE NotificationSettings ADD COLUMN "notifyOnExecuteJs" INTEGER DEFAULT 1`,
		`ALTER TABLE NotificationSettings ADD COLUMN "notifyOnMacroExecute" INTEGER DEFAULT 0`,
		// Per-key notification settings table
		`CREATE TABLE IF NOT EXISTS ApiKeyNotificationSettings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"apiKeyId" INTEGER NOT NULL UNIQUE,
			"discordWebhookUrl" TEXT,
			"notifyEmail" TEXT,
			"notifyOnExecuteJs" INTEGER DEFAULT 0,
			"notifyOnMacroExecute" INTEGER DEFAULT 0,
			"notifyOnRateLimit" INTEGER DEFAULT 0,
			"notifyOnError" INTEGER DEFAULT 0,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		// Dashboard session tokens — replaces master-key-as-bearer for the dashboard
		`CREATE TABLE IF NOT EXISTS Sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"tokenHash" TEXT NOT NULL UNIQUE,
			"userAgent" TEXT,
			"ipAddress" TEXT,
			"expiresAt" TEXT NOT NULL,
			"lastUsedAt" TEXT,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON Sessions("tokenHash")`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON Sessions("userId")`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires ON Sessions("expiresAt")`,
		// Audit trail for cross-world remote-request operations
		`CREATE TABLE IF NOT EXISTS RemoteRequestLogs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"sourceClientId" TEXT NOT NULL,
			"sourceTokenId" INTEGER NOT NULL,
			"targetClientId" TEXT NOT NULL,
			action TEXT NOT NULL,
			success INTEGER NOT NULL,
			"errorMessage" TEXT,
			"sourceIp" TEXT,
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_remote_req_user ON RemoteRequestLogs("userId")`,
		`CREATE INDEX IF NOT EXISTS idx_remote_req_created ON RemoteRequestLogs("createdAt")`,
		// D7: hash scoped API keys — add keyHash column and index
		`ALTER TABLE ApiKeys ADD COLUMN "keyHash" TEXT`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_hash ON ApiKeys("keyHash") WHERE "keyHash" IS NOT NULL`,
		// v3.0: extended notification settings
		`ALTER TABLE NotificationSettings ADD COLUMN "notifyOnNewClientConnect" INTEGER DEFAULT 1`,
		`ALTER TABLE NotificationSettings ADD COLUMN "notificationDebounceWindowSecs" INTEGER DEFAULT 0`,
		`ALTER TABLE NotificationSettings ADD COLUMN "remoteRequestBatchWindowSecs" INTEGER DEFAULT 300`,
		// Per-world notification settings table
		`CREATE TABLE IF NOT EXISTS KnownClientNotificationSettings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"knownClientId" INTEGER NOT NULL UNIQUE,
			"userId" INTEGER NOT NULL,
			"discordWebhookUrl" TEXT,
			"notifyEmail" TEXT,
			"notifyOnConnect" INTEGER NOT NULL DEFAULT 0,
			"notifyOnDisconnect" INTEGER NOT NULL DEFAULT 0,
			"notifyOnExecuteJs" INTEGER NOT NULL DEFAULT 0,
			"notifyOnMacroExecute" INTEGER NOT NULL DEFAULT 0,
			"createdAt" TEXT DEFAULT (datetime('now')),
			"updatedAt" TEXT DEFAULT (datetime('now'))
		)`,
		// Module event log (execute-js, macro-execute, settings-change persistence)
		`CREATE TABLE IF NOT EXISTS ModuleEventLogs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"userId" INTEGER NOT NULL,
			"clientId" TEXT NOT NULL,
			"worldTitle" TEXT NOT NULL DEFAULT '',
			"eventType" TEXT NOT NULL,
			actor TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			"createdAt" TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_module_event_user ON ModuleEventLogs("userId")`,
		`CREATE INDEX IF NOT EXISTS idx_module_event_created ON ModuleEventLogs("createdAt")`,
		// Cumulative metrics snapshot (single-row, persists endpoint/user tallies across restarts)
		`CREATE TABLE IF NOT EXISTS MetricsSnapshots (
			id INTEGER PRIMARY KEY DEFAULT 1 CHECK (id = 1),
			"byEndpoint"  TEXT NOT NULL DEFAULT '{}',
			"byUser"      TEXT NOT NULL DEFAULT '{}',
			"errorsTotal" INTEGER NOT NULL DEFAULT 0,
			"savedAt"     TEXT DEFAULT (datetime('now'))
		)`,
		// Log cross-world request toggle
		`ALTER TABLE NotificationSettings ADD COLUMN "logCrossWorldRequests" INTEGER NOT NULL DEFAULT 1`,
		// Drop credential columns from ApiKeys — credentials moved to KnownClients
		`ALTER TABLE ApiKeys DROP COLUMN "foundryUrl"`,
		`ALTER TABLE ApiKeys DROP COLUMN "foundryUsername"`,
		`ALTER TABLE ApiKeys DROP COLUMN "encryptedFoundryPassword"`,
		`ALTER TABLE ApiKeys DROP COLUMN "passwordIv"`,
		`ALTER TABLE ApiKeys DROP COLUMN "passwordAuthTag"`,
		`ALTER TABLE ApiKeys DROP COLUMN "credentialId"`,
		// Per-client user scoping: JSON map of clientId → userId
		`ALTER TABLE ApiKeys ADD COLUMN "scopedUserIds" TEXT`,
		// Rename dailyLimit → monthlyLimit (and related columns) to match updated model.
		// These are no-ops if the column doesn't exist or was already renamed.
		`ALTER TABLE ApiKeys RENAME COLUMN "dailyLimit" TO "monthlyLimit"`,
		`ALTER TABLE ApiKeys RENAME COLUMN "requestsToday" TO "requestsThisMonth"`,
		`ALTER TABLE ApiKeys RENAME COLUMN "lastRequestDate" TO "lastResetDate"`,
		`ALTER TABLE KeyRequests RENAME COLUMN "suggestedDailyLimit" TO "suggestedMonthlyLimit"`,
	}
	for _, m := range sqliteAlterMigrations {
		_, _ = db.sqlDB.ExecContext(ctx, m) // ignore "duplicate column" errors
	}

	// One-time force-rotation migration: mark ALL existing users as requiring rotation.
	// This runs once on first upgrade — we detect "first run" by checking if any user
	// exists without apiKeyRotationRequired set. We use a dedicated migration flag row.
	forceRotationIfNeeded(ctx, db.sqlDB, "sqlite")

	// One-time data fix: an earlier bug in KnownClients.Upsert had `is_online=$11`
	// where $11 was a timestamp value, so the isOnline column was corrupted with
	// timestamp strings. sqlx fails to scan those into Go bool. We clear the table
	// since the rows have no metadata (created during pairing) — they'll be
	// recreated on the next WS connect.
	cleanupCorruptedKnownClients(ctx, db.sqlDB, "sqlite")

	// One-time schema cleanup: remove the legacy plaintext apiKey column from
	// Users. SQLite 3.37.x can't DROP COLUMN on a UNIQUE column, so we recreate
	// the table. This is idempotent: if the column is already gone, the migration
	// is skipped (and the SchemaMigrations check ensures it runs at most once).
	removeUsersAPIKeyColumn(ctx, db.sqlDB)

	// Backfill keyHash for scoped API keys that pre-date the hash migration.
	backfillAPIKeyHashes(ctx, db.sqlDB, "sqlite")

	// One-time data fix: a prior INSERT column-order bug stored ScopedClientIDs
	// in the scopes column and Scopes in scopedClientIds. Rows affected have
	// NULL in scopes (which sqlx cannot scan into a non-nullable string).
	fixNullScopesInApiKeys(ctx, db.sqlDB, "sqlite")

	// v3.0: purge unknown-world rows and enforce UNIQUE(userId, worldId).
	addWorldIdUniqueConstraint(ctx, db.sqlDB, "sqlite")

	// v3.1: grandfather existing accounts — mark all unverified users as verified
	// since email verification was introduced after these accounts were created.
	backfillEmailVerified(ctx, db.sqlDB, "sqlite")

	return nil
}

// cleanupCorruptedKnownClients deletes any KnownClients rows where isOnline is
// not a valid boolean (0 or 1). Called once on startup, idempotent via
// SchemaMigrations tracking.
func cleanupCorruptedKnownClients(ctx context.Context, sqlDB *sqlx.DB, dbType string) {
	migrationName := "cleanup_corrupted_known_clients_2026"
	tableName := "SchemaMigrations"
	if dbType != "sqlite" {
		tableName = `"SchemaMigrations"`
	}

	var count int
	_ = sqlDB.GetContext(ctx, &count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE name = $1`, tableName), migrationName)
	if count > 0 {
		return // Already applied
	}

	knownClientsTable := "KnownClients"
	if dbType != "sqlite" {
		knownClientsTable = `"KnownClients"`
	}

	// Delete rows where isOnline isn't 0 or 1. SQLite is loose about types so we
	// use CAST to compare. PostgreSQL uses BOOLEAN so this query is a no-op there.
	var deleteQuery string
	if dbType == "sqlite" {
		deleteQuery = fmt.Sprintf(`DELETE FROM %s WHERE CAST("isOnline" AS TEXT) NOT IN ('0', '1')`, knownClientsTable)
	} else {
		// On Postgres, isOnline is BOOLEAN — corruption isn't possible since the
		// driver enforces type. Skip for safety.
		_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1)`, tableName), migrationName)
		return
	}

	result, err := sqlDB.ExecContext(ctx, deleteQuery)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to cleanup corrupted KnownClients")
	} else {
		affected, _ := result.RowsAffected()
		if affected > 0 {
			log.Info().Int64("rowsDeleted", affected).Msg("Deleted KnownClients rows with corrupted isOnline column")
		}
	}

	// Also reset isOnline=0 on remaining rows in case any stuck "online" through restart
	_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`UPDATE %s SET "isOnline" = 0`, knownClientsTable))

	// Mark migration as applied
	_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1)`, tableName), migrationName)
}

// removeUsersAPIKeyColumn recreates the Users table without the legacy plaintext
// apiKey column. SQLite 3.37.x cannot DROP COLUMN on a UNIQUE-constrained column,
// so we use the table-recreation pattern instead. Idempotent via SchemaMigrations.
func removeUsersAPIKeyColumn(ctx context.Context, sqlDB *sqlx.DB) {
	const migrationName = "remove_users_api_key_column_2026"

	var count int
	_ = sqlDB.GetContext(ctx, &count, `SELECT COUNT(*) FROM SchemaMigrations WHERE name = $1`, migrationName)
	if count > 0 {
		return // Already applied
	}

	// Check whether the apiKey column still exists before doing anything.
	var hasCol int
	_ = sqlDB.GetContext(ctx, &hasCol,
		`SELECT COUNT(*) FROM pragma_table_info('Users') WHERE name = 'apiKey'`)
	if hasCol == 0 {
		// Column already gone — mark done and return.
		_, _ = sqlDB.ExecContext(ctx, `INSERT INTO SchemaMigrations (name) VALUES ($1)`, migrationName)
		return
	}

	// Recreate Users without the apiKey column using individual statements.
	// (SQLite drivers don't support multi-statement ExecContext.)
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction for Users apiKey removal")
		return
	}
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS Users_migrate_tmp (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			"requestsThisMonth" INTEGER DEFAULT 0,
			"stripeCustomerId" VARCHAR(255),
			"subscriptionStatus" VARCHAR(255) DEFAULT 'free',
			"subscriptionId" VARCHAR(255),
			"subscriptionEndsAt" DATETIME,
			"createdAt" DATETIME NOT NULL,
			"updatedAt" DATETIME NOT NULL,
			"requestsToday" INTEGER DEFAULT 0,
			"lastRequestDate" DATE,
			"maxHeadlessSessions" INTEGER,
			"emailVerified" INTEGER DEFAULT 1,
			"verificationTokenHash" TEXT,
			"verificationTokenExpiresAt" TEXT,
			role TEXT DEFAULT 'user',
			disabled INTEGER DEFAULT 0,
			"apiKeyRotationRequired" INTEGER DEFAULT 0,
			"apiKeyHash" TEXT
		)`,
		`INSERT INTO Users_migrate_tmp
			SELECT id, email, password, "requestsThisMonth", "stripeCustomerId",
			       "subscriptionStatus", "subscriptionId", "subscriptionEndsAt",
			       "createdAt", "updatedAt", "requestsToday", "lastRequestDate",
			       "maxHeadlessSessions", "emailVerified", "verificationTokenHash",
			       "verificationTokenExpiresAt", role, disabled, "apiKeyRotationRequired",
			       "apiKeyHash"
			FROM Users`,
		`DROP TABLE Users`,
		`ALTER TABLE Users_migrate_tmp RENAME TO Users`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_api_key_hash ON Users("apiKeyHash") WHERE "apiKeyHash" IS NOT NULL`,
	}
	for _, stmt := range stmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			_ = tx.Rollback()
			log.Error().Err(err).Str("stmt", stmt[:min(50, len(stmt))]).Msg("Failed to remove apiKey column from Users table")
			return
		}
	}
	// Record migration completion using a parameterized query (not string concatenation).
	if _, err := tx.ExecContext(ctx, `INSERT INTO SchemaMigrations (name) VALUES ($1)`, migrationName); err != nil {
		_ = tx.Rollback()
		log.Error().Err(err).Msg("Failed to record migration completion")
		return
	}
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit Users apiKey removal migration")
		return
	}
	log.Info().Msg("Removed legacy apiKey column from Users table")
}

// backfillAPIKeyHashes populates the keyHash column for any scoped API keys
// that were created before the keyHash migration. This runs once at startup;
// subsequent creates already store the hash via Create().
func backfillAPIKeyHashes(ctx context.Context, sqlDB *sqlx.DB, dbType string) {
	tableName := "ApiKeys"
	hashCol := `"keyHash"`
	if dbType != "sqlite" {
		tableName = `"ApiKeys"`
	}

	// Fetch rows where keyHash is still empty.
	rows, err := sqlDB.QueryContext(ctx, fmt.Sprintf(`SELECT id, key FROM %s WHERE %s IS NULL OR %s = ''`, tableName, hashCol, hashCol))
	if err != nil {
		return
	}
	defer rows.Close()

	type row struct {
		ID  int64
		Key string
	}
	var toUpdate []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.ID, &r.Key); err == nil && r.Key != "" {
			toUpdate = append(toUpdate, r)
		}
	}
	_ = rows.Close()

	for _, r := range toUpdate {
		sum := sha256.Sum256([]byte(r.Key))
		hash := hex.EncodeToString(sum[:])
		_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`UPDATE %s SET %s = $1 WHERE id = $2`, tableName, hashCol), hash, r.ID)
	}
	if len(toUpdate) > 0 {
		log.Info().Int("count", len(toUpdate)).Msg("Backfilled keyHash for scoped API keys")
	}
}

// fixNullScopesInApiKeys repairs rows where the scopes column is NULL due to a prior
// INSERT column-order bug that swapped the scopes and scopedClientIds bindings.
// Affected rows have NULL in scopes (causing sqlx scan panics) and the actual scope
// string stored in scopedClientIds. This migration is naturally idempotent — once
// scopes is set to a non-NULL value the WHERE condition no longer matches.
func fixNullScopesInApiKeys(ctx context.Context, sqlDB *sqlx.DB, dbType string) {
	tableName := "ApiKeys"
	scopedClientIdsCol := `"scopedClientIds"`
	if dbType != "sqlite" {
		tableName = `"ApiKeys"`
	}
	query := fmt.Sprintf(
		`UPDATE %s SET scopes = COALESCE(%s, ''), %s = NULL WHERE scopes IS NULL`,
		tableName, scopedClientIdsCol, scopedClientIdsCol)
	result, err := sqlDB.ExecContext(ctx, query)
	if err != nil {
		log.Warn().Err(err).Msg("fixNullScopesInApiKeys: failed")
		return
	}
	if n, _ := result.RowsAffected(); n > 0 {
		log.Info().Int64("count", n).Msg("Fixed NULL scopes in ApiKeys")
	}
}

// backfillEmailVerified grandfathers all existing users as email-verified.
// Email verification was introduced in v3.1 but pre-existing accounts had
// emailVerified=false from Sequelize's default. This runs once at startup.
func backfillEmailVerified(ctx context.Context, sqlDB *sqlx.DB, dbType string) {
	const migrationName = "backfill_email_verified_v3_1"
	smTable := "SchemaMigrations"
	usersTable := "Users"
	if dbType != "sqlite" {
		smTable = `"SchemaMigrations"`
		usersTable = `"Users"`
	}

	var count int
	_ = sqlDB.GetContext(ctx, &count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE name = $1`, smTable), migrationName)
	if count > 0 {
		return
	}

	result, err := sqlDB.ExecContext(ctx, fmt.Sprintf(
		`UPDATE %s SET "emailVerified" = 1 WHERE "emailVerified" = 0 OR "emailVerified" IS NULL`, usersTable))
	if dbType != "sqlite" {
		result, err = sqlDB.ExecContext(ctx, fmt.Sprintf(
			`UPDATE %s SET "emailVerified" = TRUE WHERE "emailVerified" IS NOT TRUE`, usersTable))
	}
	if err != nil {
		log.Warn().Err(err).Msg("backfillEmailVerified: failed to update users")
	} else if n, _ := result.RowsAffected(); n > 0 {
		log.Info().Int64("count", n).Msg("Grandfathered pre-v3.1 accounts as email-verified")
	}

	_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1)`, smTable), migrationName)
}

// backfillUserAPIKeyHashes generates a fresh master API key hash for any Users row
// whose apiKeyHash is NULL. This happens for Sequelize-era accounts where the
// plaintext apiKey column was dropped before its SHA-256 hash was stored.
// Without apiKeyHash, WebSocket connection-token validation always fails because
// the user lookup returns an empty identifier.
//
// The generated key is random — the user never knew it and doesn't need to. It
// exists only so the relay can use apiKeyHash as a stable per-account identifier
// in ClientManager and Redis. If the user later needs to use their master API key
// directly, they can regenerate it from the dashboard (which shows it once).
func backfillUserAPIKeyHashes(ctx context.Context, sqlDB *sqlx.DB, dbType string) {
	const migrationName = "backfill_user_api_key_hashes_v3_1"
	smTable := "SchemaMigrations"
	usersTable := "Users"
	if dbType != "sqlite" {
		smTable = `"SchemaMigrations"`
		usersTable = `"Users"`
	}

	var count int
	_ = sqlDB.GetContext(ctx, &count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE name = $1`, smTable), migrationName)
	if count > 0 {
		return
	}

	// Find all users with no apiKeyHash.
	rows, err := sqlDB.QueryContext(ctx, fmt.Sprintf(`SELECT id FROM %s WHERE "apiKeyHash" IS NULL OR "apiKeyHash" = ''`, usersTable))
	if err != nil {
		log.Warn().Err(err).Msg("backfillUserAPIKeyHashes: failed to query users")
		return
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	_ = rows.Close()

	for _, id := range ids {
		// GenerateAPIKey returns 64 random hex chars (32 bytes via crypto/rand).
		rawKey, err := model.GenerateAPIKey()
		if err != nil {
			log.Warn().Err(err).Int64("userId", id).Msg("backfillUserAPIKeyHashes: failed to generate key")
			continue
		}
		sum := sha256.Sum256([]byte(rawKey))
		hash := hex.EncodeToString(sum[:])
		if _, err := sqlDB.ExecContext(ctx,
			fmt.Sprintf(`UPDATE %s SET "apiKeyHash" = $1, "apiKeyRotationRequired" = TRUE WHERE id = $2`, usersTable),
			hash, id,
		); err != nil {
			log.Warn().Err(err).Int64("userId", id).Msg("backfillUserAPIKeyHashes: failed to update user")
		} else {
			log.Info().Int64("userId", id).Msg("backfillUserAPIKeyHashes: generated apiKeyHash for user")
		}
	}

	if len(ids) > 0 {
		log.Info().Int("count", len(ids)).Msg("backfillUserAPIKeyHashes: populated missing apiKeyHash values")
	}
	_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1)`, smTable), migrationName)
}

// addWorldIdUniqueConstraint is the v3.0 migration that:
//  1. Deletes all KnownClients rows with NULL worldId (unknown-world orphans from v2.x)
//  2. Creates a UNIQUE index on (userId, worldId) so each world can only have one row per user
//
// Idempotent via SchemaMigrations tracking.
func addWorldIdUniqueConstraint(ctx context.Context, sqlDB *sqlx.DB, dbType string) {
	const migrationName = "add_world_id_unique_constraint_v3"
	smTable := "SchemaMigrations"
	kcTable := "KnownClients"
	if dbType != "sqlite" {
		smTable = `"SchemaMigrations"`
		kcTable = `"KnownClients"`
	}

	var count int
	_ = sqlDB.GetContext(ctx, &count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE name = $1`, smTable), migrationName)
	if count > 0 {
		return
	}

	// Step 1: delete unknown-world rows
	result, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s WHERE "worldId" IS NULL`, kcTable))
	if err != nil {
		log.Warn().Err(err).Msg("addWorldIdUniqueConstraint: failed to delete NULL worldId rows")
	} else if n, _ := result.RowsAffected(); n > 0 {
		log.Info().Int64("count", n).Msg("Deleted unknown-world KnownClients rows (v3.0 migration)")
	}

	// Step 2: create unique index on (userId, worldId)
	var indexSQL string
	if dbType == "sqlite" {
		indexSQL = `CREATE UNIQUE INDEX IF NOT EXISTS idx_kc_user_world ON KnownClients("userId", "worldId")`
	} else {
		indexSQL = `CREATE UNIQUE INDEX IF NOT EXISTS idx_kc_user_world ON "KnownClients"("userId", "worldId")`
	}
	if _, err := sqlDB.ExecContext(ctx, indexSQL); err != nil {
		log.Warn().Err(err).Msg("addWorldIdUniqueConstraint: failed to create unique index")
	}

	_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1)`, smTable), migrationName)
	log.Info().Msg("addWorldIdUniqueConstraint migration applied")
}

// forceRotationIfNeeded sets apiKeyRotationRequired=1 for all existing users on first run.
// It uses the schema_migrations-like approach: a row in a metadata table that marks whether
// this migration has already been applied. To keep it simple, we use a dedicated column flag
// check on an idempotent approach: only run if the migration hasn't been marked complete.
func forceRotationIfNeeded(ctx context.Context, sqlDB *sqlx.DB, dbType string) {
	// Create the migration tracking table if needed
	if dbType == "sqlite" {
		_, _ = sqlDB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS SchemaMigrations (
			name TEXT PRIMARY KEY,
			"appliedAt" TEXT DEFAULT (datetime('now'))
		)`)
	} else {
		_, _ = sqlDB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS "SchemaMigrations" (
			name VARCHAR(255) PRIMARY KEY,
			"appliedAt" TIMESTAMPTZ DEFAULT NOW()
		)`)
	}

	// Check if the force-rotation migration has already been applied
	var count int
	var tableName string
	if dbType == "sqlite" {
		tableName = "SchemaMigrations"
	} else {
		tableName = `"SchemaMigrations"`
	}
	_ = sqlDB.GetContext(ctx, &count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE name = $1`, tableName), "force_master_key_rotation_2026")
	if count > 0 {
		return // Already applied
	}

	// Apply: mark ALL existing users as requiring rotation
	var usersTable string
	if dbType == "sqlite" {
		usersTable = "Users"
	} else {
		usersTable = `"Users"`
	}
	result, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`UPDATE %s SET "apiKeyRotationRequired" = 1`, usersTable))
	if err != nil {
		log.Warn().Err(err).Msg("Failed to apply force_master_key_rotation migration (this is OK on first install)")
		// Still mark as applied to avoid retry loops on fresh installs where the column doesn't exist yet
	} else {
		affected, _ := result.RowsAffected()
		log.Info().Int64("usersAffected", affected).Msg("Force master key rotation applied to existing users")
	}

	// Mark migration as applied
	_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1)`, tableName), "force_master_key_rotation_2026")
}

func (db *DB) migratePostgres(ctx context.Context) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS "Users" (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			api_key_hash VARCHAR(64) NOT NULL UNIQUE,
			requests_this_month INTEGER DEFAULT 0,
			requests_today INTEGER DEFAULT 0,
			last_request_date DATE,
			stripe_customer_id VARCHAR(255),
			subscription_status VARCHAR(50) DEFAULT 'free',
			subscription_id VARCHAR(255),
			subscription_ends_at TIMESTAMPTZ,
			max_headless_sessions INTEGER,
			role VARCHAR(20) DEFAULT 'user',
			disabled BOOLEAN DEFAULT FALSE,
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
			monthly_limit INTEGER,
			requests_this_month INTEGER DEFAULT 0,
			last_reset_date DATE,
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
		`CREATE TABLE IF NOT EXISTS "ConnectionTokens" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			token_hash VARCHAR(255) NOT NULL UNIQUE,
			name VARCHAR(255) DEFAULT '',
			allowed_ips TEXT DEFAULT '',
			allowed_target_clients TEXT DEFAULT '',
			remote_scopes TEXT DEFAULT '',
			remote_requests_per_hour INTEGER DEFAULT 0,
			source VARCHAR(20) DEFAULT 'dashboard',
			last_used_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "PairingCodes" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			code VARCHAR(32) NOT NULL UNIQUE,
			client_id VARCHAR(255),
			allowed_target_clients TEXT DEFAULT '',
			remote_scopes TEXT DEFAULT '',
			expires_at TIMESTAMPTZ NOT NULL,
			used BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "ConnectionLogs" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			client_id VARCHAR(255) NOT NULL,
			token_name VARCHAR(255),
			ip_address VARCHAR(255),
			user_agent TEXT,
			world_id VARCHAR(255),
			world_title VARCHAR(255),
			system_id VARCHAR(255),
			foundry_version VARCHAR(50),
			metadata_match BOOLEAN DEFAULT TRUE,
			flagged BOOLEAN DEFAULT FALSE,
			flag_reason TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "Credentials" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			name VARCHAR(255) NOT NULL,
			foundry_url TEXT NOT NULL,
			foundry_username VARCHAR(255) NOT NULL,
			encrypted_foundry_password TEXT NOT NULL,
			password_iv VARCHAR(255) NOT NULL,
			password_auth_tag VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "KnownClients" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			client_id VARCHAR(255) NOT NULL,
			world_id VARCHAR(255),
			world_title VARCHAR(255),
			system_id VARCHAR(255),
			system_title VARCHAR(255),
			system_version VARCHAR(50),
			foundry_version VARCHAR(50),
			custom_name VARCHAR(255),
			last_seen_at TIMESTAMPTZ,
			is_online BOOLEAN DEFAULT FALSE,
			auto_start_on_remote_request BOOLEAN DEFAULT FALSE,
			credential_id INTEGER,
			allowed_target_clients TEXT DEFAULT '',
			remote_scopes TEXT DEFAULT '',
			remote_requests_per_hour INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(user_id, client_id)
		)`,
		`CREATE TABLE IF NOT EXISTS "KeyRequests" (
			id SERIAL PRIMARY KEY,
			code VARCHAR(255) NOT NULL UNIQUE,
			app_name VARCHAR(255) NOT NULL,
			app_description TEXT DEFAULT '',
			app_url TEXT DEFAULT '',
			requested_scopes TEXT NOT NULL,
			requested_client_ids TEXT DEFAULT '',
			callback_url TEXT,
			suggested_monthly_limit INTEGER,
			suggested_expiry VARCHAR(255),
			status VARCHAR(50) DEFAULT 'pending',
			approved_key_id INTEGER,
			approved_by_id INTEGER,
			exchange_code VARCHAR(255),
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "PairRequests" (
			id SERIAL PRIMARY KEY,
			code VARCHAR(255) NOT NULL UNIQUE,
			world_id VARCHAR(255) NOT NULL,
			world_title VARCHAR(255) NOT NULL DEFAULT '',
			system_id VARCHAR(255) NOT NULL DEFAULT '',
			system_title VARCHAR(255) NOT NULL DEFAULT '',
			system_version VARCHAR(255) NOT NULL DEFAULT '',
			foundry_version VARCHAR(255) NOT NULL DEFAULT '',
			requested_remote_scopes TEXT,
			requested_target_clients TEXT,
			upgrade_only BOOLEAN DEFAULT FALSE,
			status VARCHAR(50) DEFAULT 'pending',
			pairing_code VARCHAR(255),
			user_id INTEGER,
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "NotificationSettings" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL UNIQUE,
			notify_on_connect BOOLEAN DEFAULT TRUE,
			discord_webhook_url TEXT,
			notify_email VARCHAR(255),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "AdminAuditLogs" (
			id SERIAL PRIMARY KEY,
			"adminUserId" INTEGER NOT NULL,
			action VARCHAR(100) NOT NULL,
			"targetType" VARCHAR(50) NOT NULL,
			"targetId" VARCHAR(255),
			details TEXT,
			"ipAddress" VARCHAR(64),
			"createdAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_audit_created ON "AdminAuditLogs"("createdAt")`,
		`CREATE INDEX IF NOT EXISTS idx_admin_audit_action ON "AdminAuditLogs"(action)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_audit_admin_user ON "AdminAuditLogs"("adminUserId")`,
		`CREATE TABLE IF NOT EXISTS "JWTDenylist" (
			id SERIAL PRIMARY KEY,
			"tokenJti" VARCHAR(255) NOT NULL UNIQUE,
			"expiresAt" TIMESTAMPTZ NOT NULL,
			"createdAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_jwt_denylist_jti ON "JWTDenylist"("tokenJti")`,
		`CREATE INDEX IF NOT EXISTS idx_jwt_denylist_expires ON "JWTDenylist"("expiresAt")`,
		`CREATE TABLE IF NOT EXISTS "AlertSubscriptions" (
			id SERIAL PRIMARY KEY,
			"alertType" VARCHAR(64) NOT NULL,
			channel VARCHAR(32) NOT NULL,
			destination TEXT NOT NULL,
			enabled BOOLEAN DEFAULT TRUE,
			"createdAt" TIMESTAMPTZ DEFAULT NOW(),
			"updatedAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "AlertConfig" (
			id INTEGER PRIMARY KEY DEFAULT 1,
			"discordWebhookUrl" TEXT DEFAULT '',
			"emailDestination" TEXT DEFAULT '',
			"updatedAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "Sessions" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			token_hash VARCHAR(64) NOT NULL UNIQUE,
			user_agent TEXT,
			ip_address VARCHAR(64),
			expires_at TIMESTAMPTZ NOT NULL,
			last_used_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "RemoteRequestLogs" (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			source_client_id VARCHAR(255) NOT NULL,
			source_token_id INTEGER NOT NULL,
			target_client_id VARCHAR(255) NOT NULL,
			action VARCHAR(64) NOT NULL,
			success BOOLEAN NOT NULL,
			error_message TEXT,
			source_ip VARCHAR(64),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "KnownUsers" (
			id SERIAL PRIMARY KEY,
			"knownClientId" INTEGER NOT NULL,
			"userId" VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role INTEGER NOT NULL DEFAULT 1,
			"createdAt" TIMESTAMPTZ DEFAULT NOW(),
			"updatedAt" TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE("knownClientId", "userId"),
			FOREIGN KEY("knownClientId") REFERENCES "KnownClients"(id) ON DELETE CASCADE
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
		`ALTER TABLE "ApiKeys" ADD COLUMN IF NOT EXISTS scopes TEXT DEFAULT ''`,
		`ALTER TABLE "ApiKeys" ADD COLUMN IF NOT EXISTS "scopedClientIds" TEXT DEFAULT ''`,
		`ALTER TABLE "Users" ADD COLUMN IF NOT EXISTS "emailVerified" BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE "Users" ADD COLUMN IF NOT EXISTS "verificationTokenHash" VARCHAR(255)`,
		`ALTER TABLE "Users" ADD COLUMN IF NOT EXISTS "verificationTokenExpiresAt" TIMESTAMPTZ`,
		`ALTER TABLE "Users" ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user'`,
		`ALTER TABLE "Users" ADD COLUMN IF NOT EXISTS disabled BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE "Users" ADD COLUMN IF NOT EXISTS "apiKeyRotationRequired" BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE "Users" ADD COLUMN IF NOT EXISTS "apiKeyHash" VARCHAR(64)`,
		// Pre-production cleanup: drop the legacy plaintext apiKey column.
		// See SQLite migration above for the rationale.
		`ALTER TABLE "Users" DROP COLUMN IF EXISTS "apiKey"`,
		`ALTER TABLE "Users" DROP COLUMN IF EXISTS api_key`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_api_key_hash ON "Users"("apiKeyHash") WHERE "apiKeyHash" IS NOT NULL`,
		// Connection token permissions for cross-world tunneling
		`ALTER TABLE "ConnectionTokens" ADD COLUMN IF NOT EXISTS "allowedTargetClients" TEXT DEFAULT ''`,
		`ALTER TABLE "ConnectionTokens" ADD COLUMN IF NOT EXISTS "remoteScopes" TEXT DEFAULT ''`,
		`ALTER TABLE "ConnectionTokens" ADD COLUMN IF NOT EXISTS "remoteRequestsPerHour" INTEGER DEFAULT 0`,
		`ALTER TABLE "ConnectionTokens" ADD COLUMN IF NOT EXISTS source VARCHAR(20) DEFAULT 'dashboard'`,
		`ALTER TABLE "ConnectionTokens" ADD COLUMN IF NOT EXISTS "clientId" VARCHAR(255) DEFAULT ''`,
		`ALTER TABLE "ConnectionTokens" ADD COLUMN IF NOT EXISTS "lastUsedAt" TIMESTAMPTZ`,
		// Pairing code can be bound to an existing clientId for "add browser" flows
		`ALTER TABLE "PairingCodes" ADD COLUMN IF NOT EXISTS "clientId" VARCHAR(255)`,
		`ALTER TABLE "PairingCodes" ADD COLUMN IF NOT EXISTS "allowedTargetClients" TEXT DEFAULT ''`,
		`ALTER TABLE "PairingCodes" ADD COLUMN IF NOT EXISTS "remoteScopes" TEXT DEFAULT ''`,
		// KnownClients can opt-in to auto-start headless on incoming remote-request
		`ALTER TABLE "KnownClients" ADD COLUMN IF NOT EXISTS "autoStartOnRemoteRequest" BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE "KnownClients" ADD COLUMN IF NOT EXISTS "credentialId" INTEGER`,
		// Cross-world tunneling settings moved to per-world (KnownClient)
		`ALTER TABLE "KnownClients" ADD COLUMN IF NOT EXISTS "allowedTargetClients" TEXT DEFAULT ''`,
		`ALTER TABLE "KnownClients" ADD COLUMN IF NOT EXISTS "remoteScopes" TEXT DEFAULT ''`,
		`ALTER TABLE "KnownClients" ADD COLUMN IF NOT EXISTS "remoteRequestsPerHour" INTEGER DEFAULT 0`,
		// Connection log browser name
		`ALTER TABLE "ConnectionLogs" ADD COLUMN IF NOT EXISTS "tokenName" VARCHAR(255)`,
		// Extended notification toggles
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "notifyOnDisconnect" BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "notifyOnMetadataMismatch" BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "notifyOnSettingsChange" BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "notifyOnExecuteJs" BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "notifyOnMacroExecute" BOOLEAN DEFAULT FALSE`,
		`CREATE TABLE IF NOT EXISTS "ApiKeyNotificationSettings" (
			id SERIAL PRIMARY KEY,
			"apiKeyId" INTEGER NOT NULL UNIQUE,
			"discordWebhookUrl" TEXT,
			"notifyEmail" TEXT,
			"notifyOnExecuteJs" BOOLEAN DEFAULT FALSE,
			"notifyOnMacroExecute" BOOLEAN DEFAULT FALSE,
			"notifyOnRateLimit" BOOLEAN DEFAULT FALSE,
			"notifyOnError" BOOLEAN DEFAULT FALSE,
			"createdAt" TIMESTAMPTZ DEFAULT NOW(),
			"updatedAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		// Dashboard session tokens — replaces master-key-as-bearer for the dashboard
		`CREATE TABLE IF NOT EXISTS "Sessions" (
			id SERIAL PRIMARY KEY,
			"userId" INTEGER NOT NULL,
			"tokenHash" VARCHAR(64) NOT NULL UNIQUE,
			"userAgent" TEXT,
			"ipAddress" VARCHAR(64),
			"expiresAt" TIMESTAMPTZ NOT NULL,
			"lastUsedAt" TIMESTAMPTZ,
			"createdAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON "Sessions"("tokenHash")`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON "Sessions"("userId")`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires ON "Sessions"("expiresAt")`,
		// Audit trail for cross-world remote-request operations
		`CREATE TABLE IF NOT EXISTS "RemoteRequestLogs" (
			id SERIAL PRIMARY KEY,
			"userId" INTEGER NOT NULL,
			"sourceClientId" VARCHAR(255) NOT NULL,
			"sourceTokenId" INTEGER NOT NULL,
			"targetClientId" VARCHAR(255) NOT NULL,
			action VARCHAR(64) NOT NULL,
			success BOOLEAN NOT NULL,
			"errorMessage" TEXT,
			"sourceIp" VARCHAR(64),
			"createdAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_remote_req_user ON "RemoteRequestLogs"("userId")`,
		`CREATE INDEX IF NOT EXISTS idx_remote_req_created ON "RemoteRequestLogs"("createdAt")`,
		// D7: hash scoped API keys
		`ALTER TABLE "ApiKeys" ADD COLUMN IF NOT EXISTS "keyHash" VARCHAR(64)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_hash ON "ApiKeys"("keyHash") WHERE "keyHash" IS NOT NULL`,
		// v3.0: extended notification settings
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "notifyOnNewClientConnect" BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "notificationDebounceWindowSecs" INTEGER DEFAULT 0`,
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "remoteRequestBatchWindowSecs" INTEGER DEFAULT 300`,
		// Per-world notification settings table
		`CREATE TABLE IF NOT EXISTS "KnownClientNotificationSettings" (
			id SERIAL PRIMARY KEY,
			"knownClientId" INTEGER NOT NULL UNIQUE,
			"userId" INTEGER NOT NULL,
			"discordWebhookUrl" TEXT,
			"notifyEmail" TEXT,
			"notifyOnConnect" BOOLEAN NOT NULL DEFAULT FALSE,
			"notifyOnDisconnect" BOOLEAN NOT NULL DEFAULT FALSE,
			"notifyOnExecuteJs" BOOLEAN NOT NULL DEFAULT FALSE,
			"notifyOnMacroExecute" BOOLEAN NOT NULL DEFAULT FALSE,
			"createdAt" TIMESTAMPTZ DEFAULT NOW(),
			"updatedAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		// Module event log (execute-js, macro-execute, settings-change persistence)
		`CREATE TABLE IF NOT EXISTS "ModuleEventLogs" (
			id SERIAL PRIMARY KEY,
			"userId" INTEGER NOT NULL,
			"clientId" VARCHAR(255) NOT NULL,
			"worldTitle" VARCHAR(255) NOT NULL DEFAULT '',
			"eventType" VARCHAR(64) NOT NULL,
			actor VARCHAR(255) NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			"createdAt" TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_module_event_user ON "ModuleEventLogs"("userId")`,
		`CREATE INDEX IF NOT EXISTS idx_module_event_created ON "ModuleEventLogs"("createdAt")`,
		// Cumulative metrics snapshot (single-row, persists endpoint/user tallies across restarts)
		`CREATE TABLE IF NOT EXISTS "MetricsSnapshots" (
			id INTEGER PRIMARY KEY DEFAULT 1 CHECK (id = 1),
			"byEndpoint"  TEXT NOT NULL DEFAULT '{}',
			"byUser"      TEXT NOT NULL DEFAULT '{}',
			"errorsTotal" INTEGER NOT NULL DEFAULT 0,
			"savedAt"     TIMESTAMPTZ DEFAULT NOW()
		)`,
		// Log cross-world request toggle
		`ALTER TABLE "NotificationSettings" ADD COLUMN IF NOT EXISTS "logCrossWorldRequests" BOOLEAN NOT NULL DEFAULT TRUE`,
		// Drop credential columns from ApiKeys — credentials moved to KnownClients
		`ALTER TABLE "ApiKeys" DROP COLUMN IF EXISTS "foundryUrl"`,
		`ALTER TABLE "ApiKeys" DROP COLUMN IF EXISTS "foundryUsername"`,
		`ALTER TABLE "ApiKeys" DROP COLUMN IF EXISTS "encryptedFoundryPassword"`,
		`ALTER TABLE "ApiKeys" DROP COLUMN IF EXISTS "passwordIv"`,
		`ALTER TABLE "ApiKeys" DROP COLUMN IF EXISTS "passwordAuthTag"`,
		`ALTER TABLE "ApiKeys" DROP COLUMN IF EXISTS "credentialId"`,
		// Per-client user scoping: JSON map of clientId → userId
		`ALTER TABLE "ApiKeys" ADD COLUMN IF NOT EXISTS "scopedUserIds" TEXT`,
	}
	for _, m := range alterMigrations {
		_, _ = db.sqlDB.ExecContext(ctx, m) // Idempotent: errors expected if column already exists
	}

	// One-time force-rotation migration
	forceRotationIfNeeded(ctx, db.sqlDB, "postgres")

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
		{"ApiKeys", "dailyLimit", "monthlyLimit"},
		{"ApiKeys", "requestsToday", "requestsThisMonth"},
		{"ApiKeys", "lastRequestDate", "lastResetDate"},
		{"ApiKeys", "expires_at", "expiresAt"},
		{"ApiKeys", "created_at", "createdAt"},
		{"ApiKeys", "updated_at", "updatedAt"},
		// PasswordResetTokens table
		{"PasswordResetTokens", "user_id", "userId"},
		{"PasswordResetTokens", "token_hash", "tokenHash"},
		{"PasswordResetTokens", "expires_at", "expiresAt"},
		{"PasswordResetTokens", "created_at", "createdAt"},
		{"PasswordResetTokens", "updated_at", "updatedAt"},
		// ConnectionTokens table
		{"ConnectionTokens", "user_id", "userId"},
		{"ConnectionTokens", "token_hash", "tokenHash"},
		{"ConnectionTokens", "allowed_ips", "allowedIps"},
		{"ConnectionTokens", "allowed_target_clients", "allowedTargetClients"},
		{"ConnectionTokens", "remote_scopes", "remoteScopes"},
		{"ConnectionTokens", "last_used_at", "lastUsedAt"},
		{"ConnectionTokens", "created_at", "createdAt"},
		{"ConnectionTokens", "updated_at", "updatedAt"},
		// PairingCodes table
		{"PairingCodes", "user_id", "userId"},
		{"PairingCodes", "client_id", "clientId"},
		{"PairingCodes", "allowed_target_clients", "allowedTargetClients"},
		{"PairingCodes", "remote_scopes", "remoteScopes"},
		{"PairingCodes", "expires_at", "expiresAt"},
		{"PairingCodes", "created_at", "createdAt"},
		// ConnectionLogs table
		{"ConnectionLogs", "user_id", "userId"},
		{"ConnectionLogs", "client_id", "clientId"},
		{"ConnectionLogs", "ip_address", "ipAddress"},
		{"ConnectionLogs", "user_agent", "userAgent"},
		{"ConnectionLogs", "world_id", "worldId"},
		{"ConnectionLogs", "world_title", "worldTitle"},
		{"ConnectionLogs", "system_id", "systemId"},
		{"ConnectionLogs", "foundry_version", "foundryVersion"},
		{"ConnectionLogs", "metadata_match", "metadataMatch"},
		{"ConnectionLogs", "flag_reason", "flagReason"},
		{"ConnectionLogs", "created_at", "createdAt"},
		// Credentials table
		{"Credentials", "user_id", "userId"},
		{"Credentials", "foundry_url", "foundryUrl"},
		{"Credentials", "foundry_username", "foundryUsername"},
		{"Credentials", "encrypted_foundry_password", "encryptedFoundryPassword"},
		{"Credentials", "password_iv", "passwordIv"},
		{"Credentials", "password_auth_tag", "passwordAuthTag"},
		{"Credentials", "created_at", "createdAt"},
		{"Credentials", "updated_at", "updatedAt"},
		// KnownClients table
		{"KnownClients", "user_id", "userId"},
		{"KnownClients", "client_id", "clientId"},
		{"KnownClients", "world_id", "worldId"},
		{"KnownClients", "world_title", "worldTitle"},
		{"KnownClients", "system_id", "systemId"},
		{"KnownClients", "system_title", "systemTitle"},
		{"KnownClients", "system_version", "systemVersion"},
		{"KnownClients", "foundry_version", "foundryVersion"},
		{"KnownClients", "custom_name", "customName"},
		{"KnownClients", "last_seen_at", "lastSeenAt"},
		{"KnownClients", "is_online", "isOnline"},
		{"KnownClients", "auto_start_on_remote_request", "autoStartOnRemoteRequest"},
		{"KnownClients", "credential_id", "credentialId"},
		{"KnownClients", "created_at", "createdAt"},
		{"KnownClients", "updated_at", "updatedAt"},
		// Users table — apiKeyHash column needs renaming if it was ever created via Sequelize
		{"Users", "api_key_hash", "apiKeyHash"},
		// Sessions table (Postgres might create with snake_case if Sequelize ever touches it)
		{"Sessions", "user_id", "userId"},
		{"Sessions", "token_hash", "tokenHash"},
		{"Sessions", "user_agent", "userAgent"},
		{"Sessions", "ip_address", "ipAddress"},
		{"Sessions", "expires_at", "expiresAt"},
		{"Sessions", "last_used_at", "lastUsedAt"},
		{"Sessions", "created_at", "createdAt"},
		// RemoteRequestLogs table
		{"RemoteRequestLogs", "user_id", "userId"},
		{"RemoteRequestLogs", "source_client_id", "sourceClientId"},
		{"RemoteRequestLogs", "source_token_id", "sourceTokenId"},
		{"RemoteRequestLogs", "target_client_id", "targetClientId"},
		{"RemoteRequestLogs", "error_message", "errorMessage"},
		{"RemoteRequestLogs", "source_ip", "sourceIp"},
		{"RemoteRequestLogs", "created_at", "createdAt"},
		// PairRequests table
		{"PairRequests", "world_id", "worldId"},
		{"PairRequests", "world_title", "worldTitle"},
		{"PairRequests", "system_id", "systemId"},
		{"PairRequests", "system_title", "systemTitle"},
		{"PairRequests", "system_version", "systemVersion"},
		{"PairRequests", "foundry_version", "foundryVersion"},
		{"PairRequests", "requested_remote_scopes", "requestedRemoteScopes"},
		{"PairRequests", "requested_target_clients", "requestedTargetClients"},
		{"PairRequests", "upgrade_only", "upgradeOnly"},
		{"PairRequests", "pairing_code", "pairingCode"},
		{"PairRequests", "user_id", "userId"},
		{"PairRequests", "expires_at", "expiresAt"},
		{"PairRequests", "created_at", "createdAt"},
		// KeyRequests table
		{"KeyRequests", "app_name", "appName"},
		{"KeyRequests", "app_description", "appDescription"},
		{"KeyRequests", "app_url", "appUrl"},
		{"KeyRequests", "requested_scopes", "requestedScopes"},
		{"KeyRequests", "requested_client_ids", "requestedClientIds"},
		{"KeyRequests", "callback_url", "callbackUrl"},
		{"KeyRequests", "suggested_daily_limit", "suggestedDailyLimit"},
		{"KeyRequests", "suggestedDailyLimit", "suggestedMonthlyLimit"},
		{"KeyRequests", "suggested_expiry", "suggestedExpiry"},
		{"KeyRequests", "approved_key_id", "approvedKeyId"},
		{"KeyRequests", "approved_by_id", "approvedById"},
		{"KeyRequests", "exchange_code", "exchangeCode"},
		{"KeyRequests", "expires_at", "expiresAt"},
		{"KeyRequests", "created_at", "createdAt"},
		{"KeyRequests", "updated_at", "updatedAt"},
		// NotificationSettings table
		{"NotificationSettings", "user_id", "userId"},
		{"NotificationSettings", "notify_on_connect", "notifyOnConnect"},
		{"NotificationSettings", "discord_webhook_url", "discordWebhookUrl"},
		{"NotificationSettings", "notify_email", "notifyEmail"},
		{"NotificationSettings", "created_at", "createdAt"},
		{"NotificationSettings", "updated_at", "updatedAt"},
	}
	for _, r := range renames {
		_, _ = db.sqlDB.ExecContext(ctx, fmt.Sprintf(
			`ALTER TABLE "%s" RENAME COLUMN %s TO "%s"`, r.table, r.from, r.to))
	}

	// Backfill keyHash for scoped API keys that pre-date the hash migration.
	backfillAPIKeyHashes(ctx, db.sqlDB, "postgres")

	// One-time data fix: recover NULL scopes caused by a prior INSERT column-order bug.
	fixNullScopesInApiKeys(ctx, db.sqlDB, "postgres")

	// v3.0: purge unknown-world rows and enforce UNIQUE(userId, worldId).
	addWorldIdUniqueConstraint(ctx, db.sqlDB, "postgres")

	// v3.1: grandfather existing accounts — mark all unverified users as verified
	// since email verification was introduced after these accounts were created.
	backfillEmailVerified(ctx, db.sqlDB, "postgres")

	// v3.1: generate apiKeyHash for Sequelize-era accounts that had their
	// plaintext apiKey column dropped before the hash was computed.
	backfillUserAPIKeyHashes(ctx, db.sqlDB, "postgres")

	return nil
}

// CreateAdminUser creates the initial admin user if it doesn't exist,
// or promotes an existing user with this email to the admin role.
func (db *DB) CreateAdminUser(email, password string) error {
	ctx := context.Background()

	// Check if user exists
	existing, err := db.UserStore().FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		if existing.Role != "admin" {
			if err := db.UserStore().SetRole(ctx, existing.ID, "admin"); err != nil {
				return fmt.Errorf("promote existing user to admin: %w", err)
			}
			log.Info().Str("email", email).Msg("Promoted existing user to admin")
		} else {
			log.Info().Str("email", email).Msg("Admin user already exists")
		}
		return nil
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	apiKey, err := model.GenerateAPIKey()
	if err != nil {
		return fmt.Errorf("generate API key: %w", err)
	}

	user := &model.User{
		Email:              email,
		Password:           string(hash),
		APIKey:             apiKey,
		RequestsThisMonth:  0,
		RequestsToday:      0,
		SubscriptionStatus: sql.NullString{String: "free", Valid: true},
		Role:               "admin",
		EmailVerified:      true,
	}

	if err := db.UserStore().Create(ctx, user); err != nil {
		return err
	}

	log.Info().Str("email", email).Msg("Admin user created")
	return nil
}

// UserStore returns the user store for this database.
func (db *DB) UserStore() model.UserStore {
	return &model.SQLUserStore{DB: db.sqlDB, DBType: db.dbType}
}

// ApiKeyStore returns the API key store for this database.
func (db *DB) ApiKeyStore() model.ApiKeyStore {
	return &model.SQLApiKeyStore{DB: db.sqlDB, DBType: db.dbType}
}

// PasswordResetTokenStore returns the token store for this database.
func (db *DB) PasswordResetTokenStore() model.PasswordResetTokenStore {
	return &model.SQLPasswordResetTokenStore{DB: db.sqlDB, DBType: db.dbType}
}

// WithTx runs fn within a database transaction.
func (db *DB) WithTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.sqlDB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// TxUserStore returns a user store that operates within the given transaction.
func (db *DB) TxUserStore(tx *sqlx.Tx) model.UserStore {
	return &model.SQLUserStore{DB: tx, DBType: db.dbType}
}

// TxApiKeyStore returns an API key store that operates within the given transaction.
func (db *DB) TxApiKeyStore(tx *sqlx.Tx) model.ApiKeyStore {
	return &model.SQLApiKeyStore{DB: tx, DBType: db.dbType}
}

// TxPasswordResetTokenStore returns a token store that operates within the given transaction.
func (db *DB) TxPasswordResetTokenStore(tx *sqlx.Tx) model.PasswordResetTokenStore {
	return &model.SQLPasswordResetTokenStore{DB: tx, DBType: db.dbType}
}

// ConnectionTokenStore returns the connection token store for this database.
func (db *DB) ConnectionTokenStore() model.ConnectionTokenStore {
	return &model.SQLConnectionTokenStore{DB: db.sqlDB, DBType: db.dbType}
}

// PairingCodeStore returns the pairing code store for this database.
func (db *DB) PairingCodeStore() model.PairingCodeStore {
	return &model.SQLPairingCodeStore{DB: db.sqlDB, DBType: db.dbType}
}

// ConnectionLogStore returns the connection log store for this database.
func (db *DB) ConnectionLogStore() model.ConnectionLogStore {
	return &model.SQLConnectionLogStore{DB: db.sqlDB, DBType: db.dbType}
}

// TxConnectionTokenStore returns a connection token store that operates within the given transaction.
func (db *DB) TxConnectionTokenStore(tx *sqlx.Tx) model.ConnectionTokenStore {
	return &model.SQLConnectionTokenStore{DB: tx, DBType: db.dbType}
}

// TxPairingCodeStore returns a pairing code store that operates within the given transaction.
func (db *DB) TxPairingCodeStore(tx *sqlx.Tx) model.PairingCodeStore {
	return &model.SQLPairingCodeStore{DB: tx, DBType: db.dbType}
}

// TxConnectionLogStore returns a connection log store that operates within the given transaction.
func (db *DB) TxConnectionLogStore(tx *sqlx.Tx) model.ConnectionLogStore {
	return &model.SQLConnectionLogStore{DB: tx, DBType: db.dbType}
}

// CredentialStore returns the credential store for this database.
func (db *DB) CredentialStore() model.CredentialStore {
	return &model.SQLCredentialStore{DB: db.sqlDB, DBType: db.dbType}
}

// KnownClientStore returns the known client store for this database.
func (db *DB) KnownClientStore() model.KnownClientStore {
	return &model.SQLKnownClientStore{DB: db.sqlDB, DBType: db.dbType}
}

// KeyRequestStore returns the key request store for this database.
func (db *DB) KeyRequestStore() model.KeyRequestStore {
	return &model.SQLKeyRequestStore{DB: db.sqlDB, DBType: db.dbType}
}

// PairRequestStore returns the pair request store for this database.
func (db *DB) PairRequestStore() model.PairRequestStore {
	return &model.SQLPairRequestStore{DB: db.sqlDB, DBType: db.dbType}
}

// NotificationSettingsStore returns the notification settings store for this database.
func (db *DB) NotificationSettingsStore() model.NotificationSettingsStore {
	return &model.SQLNotificationSettingsStore{DB: db.sqlDB, DBType: db.dbType}
}

// ApiKeyNotificationSettingsStore returns the per-key notification settings store.
func (db *DB) ApiKeyNotificationSettingsStore() model.ApiKeyNotificationSettingsStore {
	return &model.SQLApiKeyNotificationSettingsStore{DB: db.sqlDB, DBType: db.dbType}
}

// KnownClientNotificationSettingsStore returns the per-world notification settings store.
func (db *DB) KnownClientNotificationSettingsStore() model.KnownClientNotificationSettingsStore {
	return &model.SQLKnownClientNotificationSettingsStore{DB: db.sqlDB, DBType: db.dbType}
}

// AuditLogStore returns the admin audit log store for this database.
func (db *DB) AuditLogStore() model.AuditLogStore {
	return &model.SQLAuditLogStore{DB: db.sqlDB, DBType: db.dbType}
}

// JWTDenylistStore returns the JWT denylist store for this database.
func (db *DB) JWTDenylistStore() model.JWTDenylistStore {
	return &model.SQLJWTDenylistStore{DB: db.sqlDB, DBType: db.dbType}
}

// AlertSubscriptionStore returns the alert subscription store for this database.
func (db *DB) AlertSubscriptionStore() model.AlertSubscriptionStore {
	return &model.SQLAlertSubscriptionStore{DB: db.sqlDB, DBType: db.dbType}
}

// AlertConfigStore returns the global alert config store for this database.
func (db *DB) AlertConfigStore() model.AlertConfigStore {
	return &model.SQLAlertConfigStore{DB: db.sqlDB, DBType: db.dbType}
}

// SessionStore returns the dashboard session store for this database.
func (db *DB) SessionStore() model.SessionStore {
	return &model.SQLSessionStore{DB: db.sqlDB, DBType: db.dbType}
}

// RemoteRequestLogStore returns the cross-world audit log store.
func (db *DB) RemoteRequestLogStore() model.RemoteRequestLogStore {
	return &model.SQLRemoteRequestLogStore{DB: db.sqlDB, DBType: db.dbType}
}

// ModuleEventLogStore returns the module event log store.
func (db *DB) ModuleEventLogStore() model.ModuleEventLogStore {
	return &model.SQLModuleEventLogStore{DB: db.sqlDB, DBType: db.dbType}
}

// KnownUserStore returns the known user store for this database.
func (db *DB) KnownUserStore() model.KnownUserStore {
	return &model.SQLKnownUserStore{DB: db.sqlDB, DBType: db.dbType}
}
