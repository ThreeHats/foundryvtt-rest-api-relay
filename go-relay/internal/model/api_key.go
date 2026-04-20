package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ApiKey represents a scoped API key with optional per-client/per-user binding.
// db tags use camelCase to match Sequelize-created SQLite columns.
type ApiKey struct {
	ID                       int64          `db:"id" json:"id"`
	UserID                   int64          `db:"userId" json:"userId"`
	Key                      string         `db:"key" json:"key"`
	KeyHash                  string         `db:"keyHash" json:"-"` // SHA-256 of Key; used for auth lookups so plaintext isn't needed at runtime
	Name                     string         `db:"name" json:"name"`
	ScopedClientID           sql.NullString `db:"scopedClientId" json:"scopedClientId"`
	ScopedUserID             sql.NullString `db:"scopedUserId" json:"scopedUserId"`
	MonthlyLimit             sql.NullInt64  `db:"monthlyLimit" json:"monthlyLimit"`
	RequestsThisMonth        int            `db:"requestsThisMonth" json:"requestsThisMonth"`
	LastResetDate            *SQLiteTime    `db:"lastResetDate" json:"-"`
	Scopes                   string         `db:"scopes" json:"scopes"`
	ScopedClientIDs          sql.NullString `db:"scopedClientIds" json:"scopedClientIds"`
	ScopedUserIDs            sql.NullString `db:"scopedUserIds" json:"scopedUserIds"`
	ExpiresAt                *SQLiteTime    `db:"expiresAt" json:"expiresAt"`
	Enabled                  bool           `db:"enabled" json:"enabled"`
	CreatedAt                SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt                SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// IsExpired returns true if the key has an expiration date and it has passed.
func (k *ApiKey) IsExpired() bool {
	return k.ExpiresAt != nil && k.ExpiresAt.Valid && k.ExpiresAt.Time.Before(time.Now())
}

// GetScopes returns the parsed scope list for this key.
//
// An empty scopes string means NO scopes — all requests are denied. Scoped
// keys must explicitly declare their scopes via the dashboard or key request flow.
//
// This applies only to scoped API keys. Master keys are represented by a User
// record with no Scopes column — they bypass scope checks because the request
// context's ScopedKey field is nil.
func (k *ApiKey) GetScopes() []string {
	if k.Scopes == "" {
		return nil
	}
	return ParseScopes(k.Scopes)
}

// GetScopedUserIDs returns the per-client userId map (clientId → userId).
// Returns nil if no per-client user scoping is configured.
func (k *ApiKey) GetScopedUserIDs() map[string]string {
	if !k.ScopedUserIDs.Valid || k.ScopedUserIDs.String == "" {
		return nil
	}
	m := map[string]string{}
	_ = json.Unmarshal([]byte(k.ScopedUserIDs.String), &m)
	if len(m) == 0 {
		return nil
	}
	return m
}

// GetScopedClientIDs returns the list of allowed client IDs, or nil if unrestricted.
func (k *ApiKey) GetScopedClientIDs() []string {
	if k.ScopedClientIDs.Valid && k.ScopedClientIDs.String != "" {
		return strings.Split(k.ScopedClientIDs.String, ",")
	}
	if k.ScopedClientID.Valid && k.ScopedClientID.String != "" {
		return []string{k.ScopedClientID.String}
	}
	return nil
}

// ApiKeyStore defines operations on scoped API keys.
type ApiKeyStore interface {
	FindByKey(ctx context.Context, key string) (*ApiKey, error)
	FindByID(ctx context.Context, id int64) (*ApiKey, error)
	FindAllByUser(ctx context.Context, userID int64) ([]*ApiKey, error)
	FindAllPaginated(ctx context.Context, offset, limit int) ([]*ApiKey, int64, error)
	Create(ctx context.Context, key *ApiKey) error
	Update(ctx context.Context, key *ApiKey) error
	Delete(ctx context.Context, id int64) error
	DeleteAllByUser(ctx context.Context, userID int64) (int64, error)
	ResetMonthlyCounters(ctx context.Context) error
	IncrementMonthlyRequests(ctx context.Context, id int64) error
	// RegenerateKey replaces the key value for an existing scoped key. It
	// generates a new random plaintext key, stores only the SHA-256 hash, and
	// returns the plaintext to the caller (shown once, same as Create).
	RegenerateKey(ctx context.Context, id int64) (string, error)
}

// SQLApiKeyStore implements ApiKeyStore with sqlx.
type SQLApiKeyStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLApiKeyStore) tableName() string {
	if s.DBType == "sqlite" {
		return "ApiKeys"
	}
	return `"ApiKeys"`
}

func (s *SQLApiKeyStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLApiKeyStore) FindByKey(ctx context.Context, key string) (*ApiKey, error) {
	var k ApiKey
	// Look up by keyHash (SHA-256 of key) so the plaintext is never stored at rest.
	// Falls back to plaintext key for rows that pre-date the keyHash migration.
	hash := HashAPIKey(key)
	err := s.DB.GetContext(ctx, &k, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("key_hash")), hash)
	if err == nil {
		return &k, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	// Fallback: plaintext key (pre-migration rows). If found, backfill the hash.
	err = s.DB.GetContext(ctx, &k, fmt.Sprintf("SELECT * FROM %s WHERE key = $1", s.tableName()), key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// Backfill keyHash for this row so future lookups use the hash path.
	_, _ = s.DB.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = $2", s.tableName(), s.col("key_hash")), hash, k.ID)
	k.KeyHash = hash
	return &k, nil
}

func (s *SQLApiKeyStore) FindByID(ctx context.Context, id int64) (*ApiKey, error) {
	var k ApiKey
	err := s.DB.GetContext(ctx, &k, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", s.tableName()), id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &k, err
}

func (s *SQLApiKeyStore) FindAllByUser(ctx context.Context, userID int64) ([]*ApiKey, error) {
	var keys []*ApiKey
	err := s.DB.SelectContext(ctx, &keys, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("user_id")), userID)
	if keys == nil {
		keys = []*ApiKey{}
	}
	return keys, err
}

// FindAllPaginated returns a page of API keys plus the total count.
func (s *SQLApiKeyStore) FindAllPaginated(ctx context.Context, offset, limit int) ([]*ApiKey, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	keys := []*ApiKey{}
	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY id ASC LIMIT $1 OFFSET $2`, s.tableName())
	if err := s.DB.SelectContext(ctx, &keys, query, limit, offset); err != nil {
		return nil, 0, err
	}
	var total int64
	if err := s.DB.GetContext(ctx, &total, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, s.tableName())); err != nil {
		return nil, 0, err
	}
	return keys, total, nil
}

func (s *SQLApiKeyStore) Create(ctx context.Context, key *ApiKey) error {
	now := time.Now()
	if key.Key == "" {
		generatedKey, err := GenerateAPIKey()
		if err != nil {
			return fmt.Errorf("generate API key: %w", err)
		}
		key.Key = generatedKey
	}
	// Store only the SHA-256 hash; the plaintext Key is returned to the caller
	// once and never persisted in the database.
	key.KeyHash = HashAPIKey(key.Key)

	query := fmt.Sprintf(`INSERT INTO %s (%s, key, %s, name, %s, %s,
		%s, %s, scopes, %s, %s, %s,
		enabled, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		s.tableName(),
		s.col("user_id"), s.col("key_hash"), s.col("scoped_client_id"), s.col("scoped_user_id"),
		s.col("monthly_limit"), s.col("requests_this_month"),
		s.col("scoped_client_ids"), s.col("scoped_user_ids"), s.col("expires_at"),
		s.col("created_at"), s.col("updated_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			key.UserID, key.Key, key.KeyHash, key.Name, key.ScopedClientID, key.ScopedUserID,
			key.MonthlyLimit, key.RequestsThisMonth, key.Scopes, key.ScopedClientIDs, key.ScopedUserIDs,
			key.ExpiresAt, key.Enabled, now, now,
		).Scan(&key.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		key.UserID, key.Key, key.KeyHash, key.Name, key.ScopedClientID, key.ScopedUserID,
		key.MonthlyLimit, key.RequestsThisMonth, key.Scopes, key.ScopedClientIDs, key.ScopedUserIDs,
		key.ExpiresAt, key.Enabled, now, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	key.ID = id
	key.CreatedAt = NewSQLiteTime(now)
	key.UpdatedAt = NewSQLiteTime(now)
	return nil
}

func (s *SQLApiKeyStore) Update(ctx context.Context, key *ApiKey) error {
	key.UpdatedAt = NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET name=$1, %s=$2, %s=$3,
		%s=$4, %s=$5, %s=$6, %s=$7,
		enabled=$8, %s=$9, scopes=$10, %s=$11, %s=$12
		WHERE id=$13`,
		s.tableName(),
		s.col("scoped_client_id"), s.col("scoped_user_id"),
		s.col("monthly_limit"), s.col("requests_this_month"), s.col("last_reset_date"),
		s.col("expires_at"), s.col("updated_at"),
		s.col("scoped_client_ids"), s.col("scoped_user_ids"))
	_, err := s.DB.ExecContext(ctx, query,
		key.Name, key.ScopedClientID, key.ScopedUserID,
		key.MonthlyLimit, key.RequestsThisMonth, key.LastResetDate,
		key.ExpiresAt, key.Enabled, key.UpdatedAt,
		key.Scopes, key.ScopedClientIDs, key.ScopedUserIDs,
		key.ID)
	return err
}

func (s *SQLApiKeyStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName()), id)
	return err
}

func (s *SQLApiKeyStore) DeleteAllByUser(ctx context.Context, userID int64) (int64, error) {
	result, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE %s = $1", s.tableName(), s.col("user_id")), userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *SQLApiKeyStore) ResetMonthlyCounters(ctx context.Context) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = 0, %s = NULL, %s = $1`,
		s.tableName(), s.col("requests_this_month"), s.col("last_reset_date"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, time.Now())
	return err
}

func (s *SQLApiKeyStore) IncrementMonthlyRequests(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = %s + 1, %s = $1, %s = $2 WHERE id = $3`,
		s.tableName(),
		s.col("requests_this_month"), s.col("requests_this_month"),
		s.col("last_reset_date"), s.col("updated_at"))
	now := time.Now()
	_, err := s.DB.ExecContext(ctx, query, now, now, id)
	return err
}

func (s *SQLApiKeyStore) RegenerateKey(ctx context.Context, id int64) (string, error) {
	newKey, err := GenerateAPIKey()
	if err != nil {
		return "", fmt.Errorf("generate API key: %w", err)
	}
	newHash := HashAPIKey(newKey)
	query := fmt.Sprintf(`UPDATE %s SET key = $1, %s = $2, %s = $3 WHERE id = $4`,
		s.tableName(), s.col("key_hash"), s.col("updated_at"))
	_, err = s.DB.ExecContext(ctx, query, newKey, newHash, time.Now(), id)
	if err != nil {
		return "", err
	}
	return newKey, nil
}
