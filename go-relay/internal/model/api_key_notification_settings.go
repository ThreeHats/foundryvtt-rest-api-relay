package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ApiKeyNotificationSettings holds per-scoped-key notification preferences.
//
// These are independent from account-level NotificationSettings — when a relevant
// event fires, the dispatcher checks both the account settings AND any matching
// per-key settings, sending notifications to all configured destinations.
//
// Per-key notifications are most useful for events triggered by a specific
// integration: "I want to know when this Discord bot calls execute-js" or "alert
// me when this Obsidian plugin hits its rate limit". Account-level settings
// cover broader concerns like connect/disconnect.
type ApiKeyNotificationSettings struct {
	ID                   int64          `db:"id" json:"id"`
	ApiKeyID             int64          `db:"apiKeyId" json:"apiKeyId"`
	DiscordWebhookURL    sql.NullString `db:"discordWebhookUrl" json:"discordWebhookUrl"`
	NotifyEmail          sql.NullString `db:"notifyEmail" json:"notifyEmail"`
	NotifyOnExecuteJs    bool           `db:"notifyOnExecuteJs" json:"notifyOnExecuteJs"`
	NotifyOnMacroExecute bool           `db:"notifyOnMacroExecute" json:"notifyOnMacroExecute"`
	NotifyOnRateLimit    bool           `db:"notifyOnRateLimit" json:"notifyOnRateLimit"`
	NotifyOnError        bool           `db:"notifyOnError" json:"notifyOnError"`
	CreatedAt            SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt            SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// ApiKeyNotificationSettingsStore defines operations on per-key notification settings.
type ApiKeyNotificationSettingsStore interface {
	FindByApiKey(ctx context.Context, apiKeyID int64) (*ApiKeyNotificationSettings, error)
	Upsert(ctx context.Context, settings *ApiKeyNotificationSettings) error
	Delete(ctx context.Context, apiKeyID int64) error
}

// SQLApiKeyNotificationSettingsStore implements the store with sqlx.
type SQLApiKeyNotificationSettingsStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLApiKeyNotificationSettingsStore) tableName() string {
	if s.DBType == "sqlite" {
		return "ApiKeyNotificationSettings"
	}
	return `"ApiKeyNotificationSettings"`
}

func (s *SQLApiKeyNotificationSettingsStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLApiKeyNotificationSettingsStore) FindByApiKey(ctx context.Context, apiKeyID int64) (*ApiKeyNotificationSettings, error) {
	var ns ApiKeyNotificationSettings
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("api_key_id"))
	err := s.DB.GetContext(ctx, &ns, query, apiKeyID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &ns, err
}

func (s *SQLApiKeyNotificationSettingsStore) Upsert(ctx context.Context, settings *ApiKeyNotificationSettings) error {
	now := time.Now()

	cols := []string{
		s.col("api_key_id"),
		s.col("discord_webhook_url"),
		s.col("notify_email"),
		s.col("notify_on_execute_js"),
		s.col("notify_on_macro_execute"),
		s.col("notify_on_rate_limit"),
		s.col("notify_on_error"),
		s.col("created_at"),
		s.col("updated_at"),
	}
	placeholders := "$1, $2, $3, $4, $5, $6, $7, $8, $9"
	args := []interface{}{
		settings.ApiKeyID,
		settings.DiscordWebhookURL,
		settings.NotifyEmail,
		settings.NotifyOnExecuteJs,
		settings.NotifyOnMacroExecute,
		settings.NotifyOnRateLimit,
		settings.NotifyOnError,
		now,
		now,
	}

	colsList := ""
	for i, c := range cols {
		if i > 0 {
			colsList += ", "
		}
		colsList += c
	}

	keyword := "EXCLUDED"
	if s.DBType == "sqlite" {
		keyword = "excluded"
	}

	updateCols := cols[1:8] // skip api_key_id and created_at
	upsertSet := ""
	for i, c := range updateCols {
		if i > 0 {
			upsertSet += ", "
		}
		upsertSet += fmt.Sprintf("%s=%s.%s", c, keyword, c)
	}
	upsertSet += fmt.Sprintf(", %s=$9", s.col("updated_at"))

	if s.DBType == "sqlite" {
		query := fmt.Sprintf(`INSERT INTO %s (%s)
			VALUES (%s)
			ON CONFLICT(%s) DO UPDATE SET %s`,
			s.tableName(), colsList, placeholders, s.col("api_key_id"), upsertSet)

		result, err := s.DB.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		if id > 0 {
			settings.ID = id
		}
		settings.CreatedAt = NewSQLiteTime(now)
		settings.UpdatedAt = NewSQLiteTime(now)
		return nil
	}

	// PostgreSQL
	query := fmt.Sprintf(`INSERT INTO %s (%s)
		VALUES (%s)
		ON CONFLICT (%s) DO UPDATE SET %s
		RETURNING id`,
		s.tableName(), colsList, placeholders, s.col("api_key_id"), upsertSet)

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&settings.ID)
}

func (s *SQLApiKeyNotificationSettingsStore) Delete(ctx context.Context, apiKeyID int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", s.tableName(), s.col("api_key_id"))
	_, err := s.DB.ExecContext(ctx, query, apiKeyID)
	return err
}
