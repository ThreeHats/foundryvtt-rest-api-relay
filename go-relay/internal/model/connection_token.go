package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// TokenSourceDashboard means the connection token was minted via the dashboard
// pairing flow (a real GM ran through the pair code UI).
const TokenSourceDashboard = "dashboard"

// TokenSourceHeadless means the connection token was minted by the relay's own
// headless worker for a ChromeDP-managed Foundry session. These are typically
// short-lived and revoked at session end.
const TokenSourceHeadless = "headless"

// ConnectionToken represents a token used to authenticate Foundry module connections.
// db tags use camelCase to match Sequelize-created SQLite columns.
type ConnectionToken struct {
	ID         int64          `db:"id" json:"id"`
	UserID     int64          `db:"userId" json:"userId"`
	TokenHash  string         `db:"tokenHash" json:"-"`
	Name       string         `db:"name" json:"name"`
	AllowedIPs sql.NullString `db:"allowedIps" json:"allowedIps"`
	// AllowedTargetClients lists the clientIds this token may invoke
	// remote-request operations against. CSV. Empty = no cross-world access.
	AllowedTargetClients sql.NullString `db:"allowedTargetClients" json:"allowedTargetClients"`
	// RemoteScopes lists the scope strings (from internal/model/scopes.go) that
	// remote-request operations may invoke against allowed targets. CSV.
	// Empty = no cross-world access (matches AllowedTargetClients default).
	RemoteScopes sql.NullString `db:"remoteScopes" json:"remoteScopes"`
	// RemoteRequestsPerHour is the per-token rate limit for cross-world
	// remote-request operations. 0 = unlimited (default).
	RemoteRequestsPerHour int `db:"remoteRequestsPerHour" json:"remoteRequestsPerHour"`
	// ClientID links this token to the KnownClient (world) it was paired for.
	// Empty for tokens created before this field was added.
	ClientID   string         `db:"clientId" json:"clientId"`
	// Source identifies who minted the token: "dashboard" (TokenSourceDashboard)
	// for normal pair flow, "headless" (TokenSourceHeadless) for relay-managed
	// per-session tokens for ChromeDP instances.
	Source     string         `db:"source" json:"source"`
	LastUsedAt *SQLiteTime    `db:"lastUsedAt" json:"lastUsedAt"`
	CreatedAt  SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt  SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// GetAllowedTargets returns the parsed list of clientIds this token can target.
// Empty list means no cross-world capability.
func (t *ConnectionToken) GetAllowedTargets() []string {
	if !t.AllowedTargetClients.Valid || t.AllowedTargetClients.String == "" {
		return nil
	}
	out := []string{}
	for _, s := range strings.Split(t.AllowedTargetClients.String, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// CanTarget returns true if the given clientId is in this token's allow-list.
func (t *ConnectionToken) CanTarget(clientID string) bool {
	if clientID == "" {
		return false
	}
	for _, allowed := range t.GetAllowedTargets() {
		if allowed == clientID {
			return true
		}
	}
	return false
}

// GetRemoteScopes returns the parsed list of scopes this token has for
// remote-request operations on allowed targets.
func (t *ConnectionToken) GetRemoteScopes() []string {
	if !t.RemoteScopes.Valid || t.RemoteScopes.String == "" {
		return nil
	}
	out := []string{}
	for _, s := range strings.Split(t.RemoteScopes.String, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// HasRemoteScope returns true if the given scope string is in this token's
// remote-scopes list. Used by the remote-request handler to gate which actions
// the token can invoke on allowed targets.
func (t *ConnectionToken) HasRemoteScope(scope string) bool {
	if scope == "" {
		return false
	}
	for _, granted := range t.GetRemoteScopes() {
		if granted == scope {
			return true
		}
	}
	return false
}

// ConnectionTokenStore defines operations on connection tokens.
type ConnectionTokenStore interface {
	FindByTokenHash(ctx context.Context, hash string) (*ConnectionToken, error)
	FindByID(ctx context.Context, id int64) (*ConnectionToken, error)
	FindAllByUser(ctx context.Context, userID int64) ([]*ConnectionToken, error)
	Create(ctx context.Context, token *ConnectionToken) error
	UpdatePermissions(ctx context.Context, id int64, name string, allowedIPs, allowedTargetClients, remoteScopes sql.NullString, remoteRequestsPerHour int) error
	Delete(ctx context.Context, id int64) error
	DeleteAllByUser(ctx context.Context, userID int64) (int64, error)
	DeleteAllByClientID(ctx context.Context, userID int64, clientID string) (int64, error)
	UpdateLastUsed(ctx context.Context, id int64) error
}

// SQLConnectionTokenStore implements ConnectionTokenStore with sqlx.
type SQLConnectionTokenStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLConnectionTokenStore) tableName() string {
	if s.DBType == "sqlite" {
		return "ConnectionTokens"
	}
	return `"ConnectionTokens"`
}

func (s *SQLConnectionTokenStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLConnectionTokenStore) FindByTokenHash(ctx context.Context, hash string) (*ConnectionToken, error) {
	var token ConnectionToken
	err := s.DB.GetContext(ctx, &token, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("token_hash")), hash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &token, err
}

func (s *SQLConnectionTokenStore) FindByID(ctx context.Context, id int64) (*ConnectionToken, error) {
	var token ConnectionToken
	err := s.DB.GetContext(ctx, &token, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", s.tableName()), id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &token, err
}

func (s *SQLConnectionTokenStore) FindAllByUser(ctx context.Context, userID int64) ([]*ConnectionToken, error) {
	var tokens []*ConnectionToken
	err := s.DB.SelectContext(ctx, &tokens, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("user_id")), userID)
	if tokens == nil {
		tokens = []*ConnectionToken{}
	}
	return tokens, err
}

func (s *SQLConnectionTokenStore) Create(ctx context.Context, token *ConnectionToken) error {
	now := NewSQLiteTime(time.Now())
	if token.Source == "" {
		token.Source = TokenSourceDashboard
	}
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, name, %s, %s, %s, source, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		s.tableName(),
		s.col("user_id"), s.col("token_hash"), s.col("allowed_ips"),
		s.col("allowed_target_clients"), s.col("remote_scopes"),
		s.col("client_id"), s.col("created_at"), s.col("updated_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			token.UserID, token.TokenHash, token.Name, token.AllowedIPs,
			token.AllowedTargetClients, token.RemoteScopes, token.Source,
			token.ClientID, now, now,
		).Scan(&token.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		token.UserID, token.TokenHash, token.Name, token.AllowedIPs,
		token.AllowedTargetClients, token.RemoteScopes, token.Source,
		token.ClientID, now, now)
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

// UpdatePermissions edits a token's name, IP allowlist, allowed targets, and remote scopes
// without re-pairing. Used by the dashboard PATCH endpoint.
func (s *SQLConnectionTokenStore) UpdatePermissions(ctx context.Context, id int64, name string, allowedIPs, allowedTargetClients, remoteScopes sql.NullString, remoteRequestsPerHour int) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET name = $1, %s = $2, %s = $3, %s = $4, %s = $5, %s = $6 WHERE id = $7`,
		s.tableName(),
		s.col("allowed_ips"),
		s.col("allowed_target_clients"), s.col("remote_scopes"),
		s.col("remote_requests_per_hour"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, name, allowedIPs, allowedTargetClients, remoteScopes, remoteRequestsPerHour, now, id)
	return err
}

func (s *SQLConnectionTokenStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName()), id)
	return err
}

func (s *SQLConnectionTokenStore) DeleteAllByClientID(ctx context.Context, userID int64, clientID string) (int64, error) {
	result, err := s.DB.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE %s = $1 AND %s = $2", s.tableName(), s.col("user_id"), s.col("client_id")),
		userID, clientID)
	if err != nil {
		return 0, err
	}
	n, _ := result.RowsAffected()
	return n, nil
}

func (s *SQLConnectionTokenStore) DeleteAllByUser(ctx context.Context, userID int64) (int64, error) {
	result, err := s.DB.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE %s = $1", s.tableName(), s.col("user_id")),
		userID)
	if err != nil {
		return 0, err
	}
	n, _ := result.RowsAffected()
	return n, nil
}

func (s *SQLConnectionTokenStore) UpdateLastUsed(ctx context.Context, id int64) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf("UPDATE %s SET %s = $1, %s = $2 WHERE id = $3",
		s.tableName(), s.col("last_used_at"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, now, now, id)
	return err
}
