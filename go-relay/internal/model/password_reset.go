package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// PasswordResetToken represents a temporary token for password resets.
// db tags use camelCase to match Sequelize-created SQLite columns.
type PasswordResetToken struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"userId" json:"userId"`
	TokenHash string    `db:"tokenHash" json:"-"`
	ExpiresAt time.Time `db:"expiresAt" json:"expiresAt"`
	Used      bool      `db:"used" json:"used"`
	CreatedAt time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt" json:"updatedAt"`
}

// PasswordResetTokenStore defines operations on password reset tokens.
type PasswordResetTokenStore interface {
	FindByTokenHash(ctx context.Context, hash string) (*PasswordResetToken, error)
	Create(ctx context.Context, token *PasswordResetToken) error
	MarkUsed(ctx context.Context, id int64) error
	InvalidateForUser(ctx context.Context, userID int64) error
	CleanupExpired(ctx context.Context) (int64, error)
}

// SQLPasswordResetTokenStore implements PasswordResetTokenStore with sqlx.
type SQLPasswordResetTokenStore struct {
	DB     *sqlx.DB
	DBType string
}

func (s *SQLPasswordResetTokenStore) tableName() string {
	if s.DBType == "sqlite" {
		return "PasswordResetTokens"
	}
	return `"PasswordResetTokens"`
}

func (s *SQLPasswordResetTokenStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLPasswordResetTokenStore) FindByTokenHash(ctx context.Context, hash string) (*PasswordResetToken, error) {
	var token PasswordResetToken
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND used = %s AND %s > $2",
		s.tableName(), s.col("token_hash"), s.boolFalse(), s.col("expires_at"))
	err := s.DB.GetContext(ctx, &token, query, hash, time.Now())
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &token, err
}

func (s *SQLPasswordResetTokenStore) Create(ctx context.Context, token *PasswordResetToken) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, used, %s, %s)
		VALUES ($1, $2, $3, %s, $4, $5)`,
		s.tableName(),
		s.col("user_id"), s.col("token_hash"), s.col("expires_at"),
		s.col("created_at"), s.col("updated_at"),
		s.boolFalse())

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			token.UserID, token.TokenHash, token.ExpiresAt, now, now,
		).Scan(&token.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		token.UserID, token.TokenHash, token.ExpiresAt, now, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	token.ID = id
	token.CreatedAt = now
	token.UpdatedAt = now
	return nil
}

func (s *SQLPasswordResetTokenStore) MarkUsed(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE %s SET used = %s, %s = $1 WHERE id = $2",
		s.tableName(), s.boolTrue(), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (s *SQLPasswordResetTokenStore) InvalidateForUser(ctx context.Context, userID int64) error {
	query := fmt.Sprintf("UPDATE %s SET used = %s, %s = $1 WHERE %s = $2 AND used = %s",
		s.tableName(), s.boolTrue(), s.col("updated_at"), s.col("user_id"), s.boolFalse())
	_, err := s.DB.ExecContext(ctx, query, time.Now(), userID)
	return err
}

func (s *SQLPasswordResetTokenStore) CleanupExpired(ctx context.Context) (int64, error) {
	cutoff := time.Now().Add(-24 * time.Hour)
	query := fmt.Sprintf("DELETE FROM %s WHERE (used = %s OR %s < $1) AND %s < $2",
		s.tableName(), s.boolTrue(), s.col("expires_at"), s.col("created_at"))
	result, err := s.DB.ExecContext(ctx, query, time.Now(), cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *SQLPasswordResetTokenStore) boolTrue() string {
	if s.DBType == "sqlite" {
		return "1"
	}
	return "TRUE"
}

func (s *SQLPasswordResetTokenStore) boolFalse() string {
	if s.DBType == "sqlite" {
		return "0"
	}
	return "FALSE"
}
