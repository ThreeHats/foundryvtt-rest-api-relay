package model

import (
	"context"
	"fmt"
	"time"
)

// KnownUser represents a Foundry VTT user stored per-world so scoped API keys
// can reference a specific player even when the world is offline.
type KnownUser struct {
	ID            int64      `db:"id" json:"id"`
	KnownClientID int64      `db:"knownClientId" json:"knownClientId"`
	UserID        string     `db:"userId" json:"userId"`
	Name          string     `db:"name" json:"name"`
	Role          int        `db:"role" json:"role"`
	CreatedAt     SQLiteTime `db:"createdAt" json:"createdAt"`
	UpdatedAt     SQLiteTime `db:"updatedAt" json:"updatedAt"`
}

// KnownUserStore defines operations on stored Foundry users per-world.
type KnownUserStore interface {
	// UpsertAll replaces the full user list for the given knownClientId.
	UpsertAll(ctx context.Context, knownClientID int64, users []*KnownUser) error
	FindAllByKnownClient(ctx context.Context, knownClientID int64) ([]*KnownUser, error)
}

// SQLKnownUserStore implements KnownUserStore with sqlx.
type SQLKnownUserStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLKnownUserStore) tableName() string {
	if s.DBType == "sqlite" {
		return "KnownUsers"
	}
	return `"KnownUsers"`
}

func (s *SQLKnownUserStore) col(name string) string {
	return Col(s.DBType, name)
}

// UpsertAll deletes all existing users for the given knownClientId and inserts
// the new list. This is a full replace — the module always sends the complete list.
func (s *SQLKnownUserStore) UpsertAll(ctx context.Context, knownClientID int64, users []*KnownUser) error {
	if _, err := s.DB.ExecContext(ctx,
		fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, s.tableName(), s.col("known_client_id")),
		knownClientID,
	); err != nil {
		return fmt.Errorf("delete known users: %w", err)
	}

	now := NewSQLiteTime(time.Now())
	for _, u := range users {
		query := fmt.Sprintf(
			`INSERT INTO %s (%s, %s, name, role, %s, %s) VALUES ($1, $2, $3, $4, $5, $6)`,
			s.tableName(),
			s.col("known_client_id"), s.col("user_id"),
			s.col("created_at"), s.col("updated_at"),
		)
		if _, err := s.DB.ExecContext(ctx, query, knownClientID, u.UserID, u.Name, u.Role, now, now); err != nil {
			return fmt.Errorf("insert known user %s: %w", u.UserID, err)
		}
	}
	return nil
}

func (s *SQLKnownUserStore) FindAllByKnownClient(ctx context.Context, knownClientID int64) ([]*KnownUser, error) {
	var users []*KnownUser
	err := s.DB.SelectContext(ctx, &users,
		fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1 ORDER BY name ASC`, s.tableName(), s.col("known_client_id")),
		knownClientID,
	)
	if users == nil {
		users = []*KnownUser{}
	}
	return users, err
}
