package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// JWTDenylistEntry represents a revoked JWT identified by its jti claim.
// Used to invalidate admin sessions before their natural expiry.
type JWTDenylistEntry struct {
	ID        int64      `db:"id" json:"id"`
	TokenJTI  string     `db:"tokenJti" json:"tokenJti"`
	ExpiresAt SQLiteTime `db:"expiresAt" json:"expiresAt"`
	CreatedAt SQLiteTime `db:"createdAt" json:"createdAt"`
}

// JWTDenylistStore defines operations on the JWT denylist.
type JWTDenylistStore interface {
	Add(ctx context.Context, jti string, expiresAt time.Time) error
	IsDenied(ctx context.Context, jti string) (bool, error)
	CleanupExpired(ctx context.Context) (int64, error)
}

// SQLJWTDenylistStore implements JWTDenylistStore with sqlx.
type SQLJWTDenylistStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLJWTDenylistStore) tableName() string {
	if s.DBType == "sqlite" {
		return "JWTDenylist"
	}
	return `"JWTDenylist"`
}

func (s *SQLJWTDenylistStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLJWTDenylistStore) Add(ctx context.Context, jti string, expiresAt time.Time) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)`,
		s.tableName(), s.col("token_jti"), s.col("expires_at"), s.col("created_at"))
	_, err := s.DB.ExecContext(ctx, query, jti, expiresAt, now)
	return err
}

func (s *SQLJWTDenylistStore) IsDenied(ctx context.Context, jti string) (bool, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s = $1`,
		s.tableName(), s.col("token_jti"))
	err := s.DB.GetContext(ctx, &count, query, jti)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return count > 0, err
}

func (s *SQLJWTDenylistStore) CleanupExpired(ctx context.Context) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s < $1`,
		s.tableName(), s.col("expires_at"))
	result, err := s.DB.ExecContext(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
