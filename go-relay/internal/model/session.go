package model

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// Session represents a dashboard authentication session. The browser holds an
// opaque random token; the server stores only its SHA-256 hash. Lookups hash
// the incoming Authorization: Bearer <token> header value before querying.
//
// This is the replacement for "dashboard sends master API key as x-api-key on
// every request" — the master key now lives only in the user's password
// manager and is shown exactly once at registration/regeneration. The
// dashboard authenticates via these short-lived sessions instead.
type Session struct {
	ID         int64          `db:"id" json:"id"`
	UserID     int64          `db:"userId" json:"userId"`
	TokenHash  string         `db:"tokenHash" json:"-"`
	UserAgent  sql.NullString `db:"userAgent" json:"userAgent"`
	IPAddress  sql.NullString `db:"ipAddress" json:"ipAddress"`
	ExpiresAt  SQLiteTime     `db:"expiresAt" json:"expiresAt"`
	LastUsedAt *SQLiteTime    `db:"lastUsedAt" json:"lastUsedAt"`
	CreatedAt  SQLiteTime     `db:"createdAt" json:"createdAt"`
}

// GenerateSessionToken creates a 32-byte random token, returns the raw hex
// (sent to the browser exactly once in the login response) and its SHA-256
// hash (what gets stored in the DB).
func GenerateSessionToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate session token: %w", err)
	}
	raw = hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(sum[:])
	return raw, hash, nil
}

// HashSessionToken returns the SHA-256 hex digest of a raw session token.
// Used by the auth middleware to look up sessions by hash.
func HashSessionToken(raw string) string {
	if raw == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// SessionStore defines operations on dashboard sessions.
type SessionStore interface {
	Create(ctx context.Context, session *Session) error
	FindByTokenHash(ctx context.Context, hash string) (*Session, error)
	UpdateLastUsed(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	DeleteByTokenHash(ctx context.Context, hash string) error
	DeleteAllByUser(ctx context.Context, userID int64) (int64, error)
	CleanupExpired(ctx context.Context) (int64, error)
}

// SQLSessionStore implements SessionStore with sqlx.
type SQLSessionStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLSessionStore) tableName() string {
	if s.DBType == "sqlite" {
		return "Sessions"
	}
	return `"Sessions"`
}

func (s *SQLSessionStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLSessionStore) Create(ctx context.Context, session *Session) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		s.tableName(),
		s.col("user_id"), s.col("token_hash"), s.col("user_agent"),
		s.col("ip_address"), s.col("expires_at"), s.col("created_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			session.UserID, session.TokenHash, session.UserAgent,
			session.IPAddress, session.ExpiresAt, now,
		).Scan(&session.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		session.UserID, session.TokenHash, session.UserAgent,
		session.IPAddress, session.ExpiresAt, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	session.ID = id
	session.CreatedAt = now
	return nil
}

func (s *SQLSessionStore) FindByTokenHash(ctx context.Context, hash string) (*Session, error) {
	if hash == "" {
		return nil, nil
	}
	var session Session
	now := NewSQLiteTime(time.Now())
	// Only return non-expired sessions. Expired ones are cleaned up by the
	// background cleanup loop, but we double-check here so an expired session
	// is treated as missing instantly.
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s > $2",
		s.tableName(), s.col("token_hash"), s.col("expires_at"))
	err := s.DB.GetContext(ctx, &session, query, hash, now)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &session, err
}

func (s *SQLSessionStore) UpdateLastUsed(ctx context.Context, id int64) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = $2",
		s.tableName(), s.col("last_used_at"))
	_, err := s.DB.ExecContext(ctx, query, now, id)
	return err
}

func (s *SQLSessionStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName()), id)
	return err
}

func (s *SQLSessionStore) DeleteByTokenHash(ctx context.Context, hash string) error {
	_, err := s.DB.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE %s = $1", s.tableName(), s.col("token_hash")),
		hash)
	return err
}

func (s *SQLSessionStore) DeleteAllByUser(ctx context.Context, userID int64) (int64, error) {
	result, err := s.DB.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE %s = $1", s.tableName(), s.col("user_id")),
		userID)
	if err != nil {
		return 0, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (s *SQLSessionStore) CleanupExpired(ctx context.Context) (int64, error) {
	now := NewSQLiteTime(time.Now())
	result, err := s.DB.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE %s < $1", s.tableName(), s.col("expires_at")),
		now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
