package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PairingCode represents a short-lived code used to pair a Foundry module with a user account.
// db tags use camelCase to match Sequelize-created SQLite columns.
type PairingCode struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"userId" json:"userId"`
	Code      string     `db:"code" json:"code"`
	// ClientID, when non-null, makes this pairing code reuse an existing
	// known clientId instead of minting a fresh one. Used by the "add this
	// browser" flow where a second GM joins an already-paired world.
	ClientID sql.NullString `db:"clientId" json:"clientId"`
	// AllowedTargetClients is a CSV of clientIds the resulting connection
	// token will be allowed to invoke remote-request operations against.
	// Empty = no cross-world access (default).
	AllowedTargetClients sql.NullString `db:"allowedTargetClients" json:"allowedTargetClients"`
	// RemoteScopes is a CSV of scope strings the resulting connection token
	// will hold for cross-world operations. Empty = no cross-world access.
	RemoteScopes sql.NullString `db:"remoteScopes" json:"remoteScopes"`
	ExpiresAt    SQLiteTime     `db:"expiresAt" json:"expiresAt"`
	Used         bool           `db:"used" json:"used"`
	CreatedAt    SQLiteTime     `db:"createdAt" json:"createdAt"`
}

// PairingCodeStore defines operations on pairing codes.
type PairingCodeStore interface {
	FindByCode(ctx context.Context, code string) (*PairingCode, error)
	Create(ctx context.Context, code *PairingCode) error
	MarkUsed(ctx context.Context, id int64) error
	CleanupExpired(ctx context.Context) (int64, error)
}

// SQLPairingCodeStore implements PairingCodeStore with sqlx.
type SQLPairingCodeStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLPairingCodeStore) tableName() string {
	if s.DBType == "sqlite" {
		return "PairingCodes"
	}
	return `"PairingCodes"`
}

func (s *SQLPairingCodeStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLPairingCodeStore) boolTrue() string {
	if s.DBType == "sqlite" {
		return "1"
	}
	return "TRUE"
}

func (s *SQLPairingCodeStore) boolFalse() string {
	if s.DBType == "sqlite" {
		return "0"
	}
	return "FALSE"
}

func (s *SQLPairingCodeStore) FindByCode(ctx context.Context, code string) (*PairingCode, error) {
	var pc PairingCode
	// Note: SQLiteTime.Value() formats as "2006-01-02 15:04:05" so we must pass
	// a SQLiteTime value here (NOT a raw time.Time) so the parameter formatting
	// matches what's stored in the column. Using time.Now() directly produces a
	// different string format and the comparison silently fails (everything looks
	// "expired" because the string compare is wrong).
	now := NewSQLiteTime(time.Now())
	// Find by code regardless of used status so callers can distinguish 404 (not found)
	// from 410 (found but already used). Expired codes still return nil → 404.
	query := fmt.Sprintf("SELECT * FROM %s WHERE code = $1 AND %s > $2",
		s.tableName(), s.col("expires_at"))
	err := s.DB.GetContext(ctx, &pc, query, code, now)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &pc, err
}

func (s *SQLPairingCodeStore) Create(ctx context.Context, code *PairingCode) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (%s, code, %s, %s, %s, %s, used, %s)
		VALUES ($1, $2, $3, $4, $5, $6, %s, $7)`,
		s.tableName(),
		s.col("user_id"),
		s.col("client_id"),
		s.col("allowed_target_clients"),
		s.col("remote_scopes"),
		s.col("expires_at"),
		s.col("created_at"),
		s.boolFalse())

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			code.UserID, code.Code, code.ClientID, code.AllowedTargetClients,
			code.RemoteScopes, code.ExpiresAt, now,
		).Scan(&code.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		code.UserID, code.Code, code.ClientID, code.AllowedTargetClients,
		code.RemoteScopes, code.ExpiresAt, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	code.ID = id
	code.CreatedAt = NewSQLiteTime(now)
	return nil
}

func (s *SQLPairingCodeStore) MarkUsed(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE %s SET used = %s WHERE id = $1",
		s.tableName(), s.boolTrue())
	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}

func (s *SQLPairingCodeStore) CleanupExpired(ctx context.Context) (int64, error) {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf("DELETE FROM %s WHERE used = %s OR %s < $1",
		s.tableName(), s.boolTrue(), s.col("expires_at"))
	result, err := s.DB.ExecContext(ctx, query, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
