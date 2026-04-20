package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// KnownClientNotificationSettings holds per-world notification preferences.
//
// These fire in addition to account-level and per-key notifications. The
// four events covered here are the ones that are most naturally scoped to a
// specific Foundry world: connect, disconnect, execute-js, and macro-execute.
type KnownClientNotificationSettings struct {
	ID                   int64          `db:"id" json:"id"`
	KnownClientID        int64          `db:"knownClientId" json:"knownClientId"`
	UserID               int64          `db:"userId" json:"userId"`
	DiscordWebhookURL    sql.NullString `db:"discordWebhookUrl" json:"discordWebhookUrl"`
	NotifyEmail          sql.NullString `db:"notifyEmail" json:"notifyEmail"`
	NotifyOnConnect      bool           `db:"notifyOnConnect" json:"notifyOnConnect"`
	NotifyOnDisconnect   bool           `db:"notifyOnDisconnect" json:"notifyOnDisconnect"`
	NotifyOnExecuteJs    bool           `db:"notifyOnExecuteJs" json:"notifyOnExecuteJs"`
	NotifyOnMacroExecute bool           `db:"notifyOnMacroExecute" json:"notifyOnMacroExecute"`
	CreatedAt            SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt            SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// KnownClientNotificationSettingsStore defines operations on per-world notification settings.
type KnownClientNotificationSettingsStore interface {
	FindByKnownClient(ctx context.Context, knownClientID int64) (*KnownClientNotificationSettings, error)
	Upsert(ctx context.Context, s *KnownClientNotificationSettings) error
	Delete(ctx context.Context, knownClientID int64) error
}

// SQLKnownClientNotificationSettingsStore implements the store with sqlx.
type SQLKnownClientNotificationSettingsStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLKnownClientNotificationSettingsStore) tableName() string {
	if s.DBType == "sqlite" {
		return "KnownClientNotificationSettings"
	}
	return `"KnownClientNotificationSettings"`
}

func (s *SQLKnownClientNotificationSettingsStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLKnownClientNotificationSettingsStore) FindByKnownClient(ctx context.Context, knownClientID int64) (*KnownClientNotificationSettings, error) {
	var ns KnownClientNotificationSettings
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("known_client_id"))
	err := s.DB.GetContext(ctx, &ns, query, knownClientID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &ns, err
}

func (s *SQLKnownClientNotificationSettingsStore) Upsert(ctx context.Context, settings *KnownClientNotificationSettings) error {
	now := time.Now()

	cols := []string{
		s.col("known_client_id"),
		s.col("user_id"),
		s.col("discord_webhook_url"),
		s.col("notify_email"),
		s.col("notify_on_connect"),
		s.col("notify_on_disconnect"),
		s.col("notify_on_execute_js"),
		s.col("notify_on_macro_execute"),
		s.col("created_at"),
		s.col("updated_at"),
	}
	placeholders := "$1, $2, $3, $4, $5, $6, $7, $8, $9, $10"
	args := []interface{}{
		settings.KnownClientID,
		settings.UserID,
		settings.DiscordWebhookURL,
		settings.NotifyEmail,
		settings.NotifyOnConnect,
		settings.NotifyOnDisconnect,
		settings.NotifyOnExecuteJs,
		settings.NotifyOnMacroExecute,
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

	// Update columns: skip known_client_id, user_id, and created_at
	updateCols := cols[2:8]
	upsertSet := ""
	for i, c := range updateCols {
		if i > 0 {
			upsertSet += ", "
		}
		upsertSet += fmt.Sprintf("%s=%s.%s", c, keyword, c)
	}
	upsertSet += fmt.Sprintf(", %s=$10", s.col("updated_at"))

	if s.DBType == "sqlite" {
		query := fmt.Sprintf(`INSERT INTO %s (%s)
			VALUES (%s)
			ON CONFLICT(%s) DO UPDATE SET %s`,
			s.tableName(), colsList, placeholders, s.col("known_client_id"), upsertSet)

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
		s.tableName(), colsList, placeholders, s.col("known_client_id"), upsertSet)

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&settings.ID)
}

func (s *SQLKnownClientNotificationSettingsStore) Delete(ctx context.Context, knownClientID int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", s.tableName(), s.col("known_client_id"))
	_, err := s.DB.ExecContext(ctx, query, knownClientID)
	return err
}
