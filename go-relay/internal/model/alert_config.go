package model

import (
	"context"
	"fmt"
	"time"
)

// AlertConfig holds the global alert destinations (one row, id=1).
type AlertConfig struct {
	ID                int64      `db:"id"`
	DiscordWebhookURL string     `db:"discordWebhookUrl"`
	EmailDestination  string     `db:"emailDestination"`
	UpdatedAt         SQLiteTime `db:"updatedAt"`
}

// AlertConfigStore manages the single AlertConfig row.
type AlertConfigStore interface {
	Get(ctx context.Context) (*AlertConfig, error)
	Save(ctx context.Context, cfg *AlertConfig) error
}

// AlertSubscription represents a per-(alertType, channel) enabled state.
type AlertSubscription struct {
	ID          int64      `db:"id" json:"id"`
	AlertType   string     `db:"alertType" json:"alertType"`
	Channel     string     `db:"channel" json:"channel"` // "discord" | "email"
	Destination string     `db:"destination" json:"-"`   // unused; kept for schema compat
	Enabled     bool       `db:"enabled" json:"enabled"`
	CreatedAt   SQLiteTime `db:"createdAt" json:"createdAt"`
	UpdatedAt   SQLiteTime `db:"updatedAt" json:"updatedAt"`
}

// AlertSubscriptionStore defines operations on alert subscriptions.
type AlertSubscriptionStore interface {
	Create(ctx context.Context, sub *AlertSubscription) error
	FindAll(ctx context.Context) ([]*AlertSubscription, error)
	FindByAlertType(ctx context.Context, alertType string) ([]*AlertSubscription, error)
	FindByID(ctx context.Context, id int64) (*AlertSubscription, error)
	Delete(ctx context.Context, id int64) error
	BulkReplace(ctx context.Context, subs []*AlertSubscription) error
}

// SQLAlertSubscriptionStore implements AlertSubscriptionStore with sqlx.
type SQLAlertSubscriptionStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLAlertSubscriptionStore) tableName() string {
	if s.DBType == "sqlite" {
		return "AlertSubscriptions"
	}
	return `"AlertSubscriptions"`
}

func (s *SQLAlertSubscriptionStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLAlertSubscriptionStore) Create(ctx context.Context, sub *AlertSubscription) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (%s, channel, destination, enabled, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		s.tableName(), s.col("alert_type"), s.col("created_at"), s.col("updated_at"))
	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query, sub.AlertType, sub.Channel, sub.Destination, sub.Enabled, now, now).Scan(&sub.ID)
	}
	res, err := s.DB.ExecContext(ctx, query, sub.AlertType, sub.Channel, sub.Destination, sub.Enabled, now, now)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	sub.ID = id
	return nil
}

func (s *SQLAlertSubscriptionStore) FindAll(ctx context.Context) ([]*AlertSubscription, error) {
	subs := []*AlertSubscription{}
	err := s.DB.SelectContext(ctx, &subs, fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC", s.tableName()))
	return subs, err
}

func (s *SQLAlertSubscriptionStore) FindByAlertType(ctx context.Context, alertType string) ([]*AlertSubscription, error) {
	subs := []*AlertSubscription{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND enabled = TRUE", s.tableName(), s.col("alert_type"))
	if s.DBType == "sqlite" {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND enabled = 1", s.tableName(), s.col("alert_type"))
	}
	err := s.DB.SelectContext(ctx, &subs, query, alertType)
	return subs, err
}

func (s *SQLAlertSubscriptionStore) FindByID(ctx context.Context, id int64) (*AlertSubscription, error) {
	var sub AlertSubscription
	err := s.DB.GetContext(ctx, &sub, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", s.tableName()), id)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *SQLAlertSubscriptionStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName()), id)
	return err
}

func (s *SQLAlertSubscriptionStore) BulkReplace(ctx context.Context, subs []*AlertSubscription) error {
	if _, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", s.tableName())); err != nil {
		return err
	}
	now := time.Now()
	for _, sub := range subs {
		query := fmt.Sprintf(`INSERT INTO %s (%s, channel, destination, enabled, %s, %s)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			s.tableName(), s.col("alert_type"), s.col("created_at"), s.col("updated_at"))
		if s.DBType != "sqlite" {
			query += " RETURNING id"
			if err := s.DB.QueryRowContext(ctx, query, sub.AlertType, sub.Channel, "", true, now, now).Scan(&sub.ID); err != nil {
				return err
			}
		} else {
			res, err := s.DB.ExecContext(ctx, query, sub.AlertType, sub.Channel, "", true, now, now)
			if err != nil {
				return err
			}
			sub.ID, _ = res.LastInsertId()
		}
	}
	return nil
}

// SQLAlertConfigStore implements AlertConfigStore with sqlx.
type SQLAlertConfigStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLAlertConfigStore) tableName() string {
	if s.DBType == "sqlite" {
		return "AlertConfig"
	}
	return `"AlertConfig"`
}

func (s *SQLAlertConfigStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLAlertConfigStore) Get(ctx context.Context) (*AlertConfig, error) {
	// Ensure the single row exists.
	if s.DBType == "sqlite" {
		s.DB.ExecContext(ctx, fmt.Sprintf( //nolint:errcheck
			`INSERT OR IGNORE INTO %s (id, %s, %s, %s) VALUES (1, '', '', datetime('now'))`,
			s.tableName(), s.col("discord_webhook_url"), s.col("email_destination"), s.col("updated_at"),
		))
	} else {
		s.DB.ExecContext(ctx, fmt.Sprintf( //nolint:errcheck
			`INSERT INTO %s (id, %s, %s, %s) VALUES (1, '', '', NOW()) ON CONFLICT (id) DO NOTHING`,
			s.tableName(), s.col("discord_webhook_url"), s.col("email_destination"), s.col("updated_at"),
		))
	}
	var cfg AlertConfig
	err := s.DB.GetContext(ctx, &cfg, fmt.Sprintf("SELECT * FROM %s WHERE id = 1", s.tableName()))
	if err != nil {
		return &AlertConfig{}, nil
	}
	return &cfg, nil
}

func (s *SQLAlertConfigStore) Save(ctx context.Context, cfg *AlertConfig) error {
	cfg.UpdatedAt = NewSQLiteTime(time.Now())
	if s.DBType == "sqlite" {
		_, err := s.DB.ExecContext(ctx, fmt.Sprintf(
			`INSERT OR REPLACE INTO %s (id, %s, %s, %s) VALUES (1, $1, $2, $3)`,
			s.tableName(), s.col("discord_webhook_url"), s.col("email_destination"), s.col("updated_at"),
		), cfg.DiscordWebhookURL, cfg.EmailDestination, cfg.UpdatedAt)
		return err
	}
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, %s, %s, %s) VALUES (1, $1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET %s=$1, %s=$2, %s=$3`,
		s.tableName(),
		s.col("discord_webhook_url"), s.col("email_destination"), s.col("updated_at"),
		s.col("discord_webhook_url"), s.col("email_destination"), s.col("updated_at"),
	), cfg.DiscordWebhookURL, cfg.EmailDestination, cfg.UpdatedAt)
	return err
}
