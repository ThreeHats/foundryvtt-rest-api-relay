package model

import "strings"

// snakeToCamel maps all known snake_case column names to their Sequelize camelCase equivalents.
var snakeToCamel = map[string]string{
	// Users table
	"api_key":              "apiKey",
	"requests_this_month":  "requestsThisMonth",
	"requests_today":       "requestsToday",
	"last_request_date":    "lastRequestDate",
	"stripe_customer_id":   "stripeCustomerId",
	"subscription_status":  "subscriptionStatus",
	"subscription_id":      "subscriptionId",
	"subscription_ends_at": "subscriptionEndsAt",
	"max_headless_sessions": "maxHeadlessSessions",
	"created_at":           "createdAt",
	"updated_at":           "updatedAt",
	// PasswordResetTokens table
	"user_id":    "userId",
	"token_hash": "tokenHash",
	"expires_at": "expiresAt",
	// ApiKeys table
	"master_api_key":            "masterApiKey",
	"scoped_client_id":          "scopedClientId",
	"scoped_user_id":            "scopedUserId",
	"daily_limit":               "dailyLimit",
	"foundry_url":               "foundryUrl",
	"foundry_username":          "foundryUsername",
	"foundry_password":          "foundryPassword",
	"encrypted_foundry_password": "encryptedFoundryPassword",
	"password_iv":               "passwordIv",
	"password_auth_tag":         "passwordAuthTag",
}

// Col returns the correct column name for a given database type.
// For SQLite (Sequelize-created tables), maps snake_case to quoted camelCase.
// For PostgreSQL, returns the snake_case name as-is.
func Col(dbType, name string) string {
	if dbType != "sqlite" {
		return name
	}
	if camel, ok := snakeToCamel[name]; ok {
		return `"` + camel + `"`
	}
	return name // id, email, password, name, used — same in both
}

// NormalizeColumnName strips underscores and lowercases for sqlx mapper matching.
// This allows db:"api_key" tags to match both "api_key" (Go-created) and "apiKey" (Sequelize-created) columns.
func NormalizeColumnName(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "_", ""))
}
