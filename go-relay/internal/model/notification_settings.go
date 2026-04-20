package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// NotificationSettings represents a user's account-level notification preferences.
//
// Each event type has its own toggle so the user can opt in/out individually.
// Events:
//   - NotifyOnConnect: a Foundry client connected to the relay
//   - NotifyOnDisconnect: a Foundry client disconnected from the relay
//   - NotifyOnMetadataMismatch: a connection's metadata fingerprint differs from the
//     previously-known clientId (possible impersonation)
//   - NotifyOnSettingsChange: a sensitive Foundry module setting was changed
//   - NotifyOnExecuteJs: the execute-js action ran in a Foundry world
//   - NotifyOnMacroExecute: the macro-execute action ran in a Foundry world
//
// Per-scoped-key settings (see ApiKeyNotificationSettings) can fire IN ADDITION to
// these account-level settings — the dispatcher routes events to all matching
// destinations.
type NotificationSettings struct {
	ID                       int64          `db:"id" json:"id"`
	UserID                   int64          `db:"userId" json:"userId"`
	DiscordWebhookURL        sql.NullString `db:"discordWebhookUrl" json:"discordWebhookUrl"`
	NotifyEmail              sql.NullString `db:"notifyEmail" json:"notifyEmail"`
	NotifyOnConnect          bool           `db:"notifyOnConnect" json:"notifyOnConnect"`
	NotifyOnDisconnect       bool           `db:"notifyOnDisconnect" json:"notifyOnDisconnect"`
	NotifyOnMetadataMismatch bool           `db:"notifyOnMetadataMismatch" json:"notifyOnMetadataMismatch"`
	NotifyOnSettingsChange   bool           `db:"notifyOnSettingsChange" json:"notifyOnSettingsChange"`
	NotifyOnExecuteJs              bool           `db:"notifyOnExecuteJs" json:"notifyOnExecuteJs"`
	NotifyOnMacroExecute           bool           `db:"notifyOnMacroExecute" json:"notifyOnMacroExecute"`
	NotifyOnNewClientConnect       bool           `db:"notifyOnNewClientConnect" json:"notifyOnNewClientConnect"`
	NotificationDebounceWindowSecs int            `db:"notificationDebounceWindowSecs" json:"notificationDebounceWindowSecs"`
	RemoteRequestBatchWindowSecs   int            `db:"remoteRequestBatchWindowSecs" json:"remoteRequestBatchWindowSecs"`
	LogCrossWorldRequests          bool           `db:"logCrossWorldRequests" json:"logCrossWorldRequests"`
	CreatedAt                      SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt                      SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// NotificationSettingsStore defines operations on notification settings.
type NotificationSettingsStore interface {
	FindByUser(ctx context.Context, userID int64) (*NotificationSettings, error)
	Upsert(ctx context.Context, settings *NotificationSettings) error
}

// SQLNotificationSettingsStore implements NotificationSettingsStore with sqlx.
type SQLNotificationSettingsStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLNotificationSettingsStore) tableName() string {
	if s.DBType == "sqlite" {
		return "NotificationSettings"
	}
	return `"NotificationSettings"`
}

func (s *SQLNotificationSettingsStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLNotificationSettingsStore) FindByUser(ctx context.Context, userID int64) (*NotificationSettings, error) {
	var ns NotificationSettings
	err := s.DB.GetContext(ctx, &ns, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("user_id")), userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &ns, err
}

func (s *SQLNotificationSettingsStore) Upsert(ctx context.Context, settings *NotificationSettings) error {
	now := time.Now()

	cols := []string{
		s.col("user_id"),
		s.col("discord_webhook_url"),
		s.col("notify_email"),
		s.col("notify_on_connect"),
		s.col("notify_on_disconnect"),
		s.col("notify_on_metadata_mismatch"),
		s.col("notify_on_settings_change"),
		s.col("notify_on_execute_js"),
		s.col("notify_on_macro_execute"),
		s.col("notify_on_new_client_connect"),
		s.col("notification_debounce_window_secs"),
		s.col("remote_request_batch_window_secs"),
		s.col("log_cross_world_requests"),
		s.col("created_at"),
		s.col("updated_at"),
	}
	placeholders := "$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15"
	args := []interface{}{
		settings.UserID,
		settings.DiscordWebhookURL,
		settings.NotifyEmail,
		settings.NotifyOnConnect,
		settings.NotifyOnDisconnect,
		settings.NotifyOnMetadataMismatch,
		settings.NotifyOnSettingsChange,
		settings.NotifyOnExecuteJs,
		settings.NotifyOnMacroExecute,
		settings.NotifyOnNewClientConnect,
		settings.NotificationDebounceWindowSecs,
		settings.RemoteRequestBatchWindowSecs,
		settings.LogCrossWorldRequests,
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

	// Build the UPDATE SET clause for ON CONFLICT
	updateCols := cols[1:14] // skip user_id (conflict column) and created_at; include updated_at separately below
	upsertSet := ""
	keyword := "EXCLUDED"
	if s.DBType == "sqlite" {
		keyword = "excluded"
	}
	for i, c := range updateCols {
		if i > 0 {
			upsertSet += ", "
		}
		upsertSet += fmt.Sprintf("%s=%s.%s", c, keyword, c)
	}
	// Always touch updated_at
	upsertSet += fmt.Sprintf(", %s=$15", s.col("updated_at"))

	if s.DBType == "sqlite" {
		query := fmt.Sprintf(`INSERT INTO %s (%s)
			VALUES (%s)
			ON CONFLICT(%s) DO UPDATE SET %s`,
			s.tableName(), colsList, placeholders, s.col("user_id"), upsertSet)

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
		s.tableName(), colsList, placeholders, s.col("user_id"), upsertSet)

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&settings.ID)
}
