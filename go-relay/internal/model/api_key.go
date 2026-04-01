package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// ApiKey represents a scoped API key with optional per-client/per-user binding.
// db tags use camelCase to match Sequelize-created SQLite columns.
type ApiKey struct {
	ID                       int64          `db:"id" json:"id"`
	UserID                   int64          `db:"userId" json:"userId"`
	Key                      string         `db:"key" json:"key"`
	Name                     string         `db:"name" json:"name"`
	ScopedClientID           sql.NullString `db:"scopedClientId" json:"scopedClientId"`
	ScopedUserID             sql.NullString `db:"scopedUserId" json:"scopedUserId"`
	DailyLimit               sql.NullInt64  `db:"dailyLimit" json:"dailyLimit"`
	RequestsToday            int            `db:"requestsToday" json:"requestsToday"`
	LastRequestDate          *SQLiteTime    `db:"lastRequestDate" json:"-"`
	FoundryURL               sql.NullString `db:"foundryUrl" json:"foundryUrl"`
	FoundryUsername           sql.NullString `db:"foundryUsername" json:"foundryUsername"`
	EncryptedFoundryPassword sql.NullString `db:"encryptedFoundryPassword" json:"-"`
	PasswordIV               sql.NullString `db:"passwordIv" json:"-"`
	PasswordAuthTag          sql.NullString `db:"passwordAuthTag" json:"-"`
	ExpiresAt                *SQLiteTime    `db:"expiresAt" json:"expiresAt"`
	Enabled                  bool           `db:"enabled" json:"enabled"`
	CreatedAt                SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt                SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// IsExpired returns true if the key has an expiration date and it has passed.
func (k *ApiKey) IsExpired() bool {
	return k.ExpiresAt != nil && k.ExpiresAt.Valid && k.ExpiresAt.Time.Before(time.Now())
}

// HasStoredCredentials returns true if the key has encrypted Foundry credentials.
func (k *ApiKey) HasStoredCredentials() bool {
	return k.EncryptedFoundryPassword.Valid && k.EncryptedFoundryPassword.String != "" &&
		k.FoundryURL.Valid && k.FoundryURL.String != "" &&
		k.FoundryUsername.Valid && k.FoundryUsername.String != ""
}

// ApiKeyStore defines operations on scoped API keys.
type ApiKeyStore interface {
	FindByKey(ctx context.Context, key string) (*ApiKey, error)
	FindByID(ctx context.Context, id int64) (*ApiKey, error)
	FindAllByUser(ctx context.Context, userID int64) ([]*ApiKey, error)
	Create(ctx context.Context, key *ApiKey) error
	Update(ctx context.Context, key *ApiKey) error
	Delete(ctx context.Context, id int64) error
	DeleteAllByUser(ctx context.Context, userID int64) (int64, error)
	ResetDailyCounters(ctx context.Context) error
	IncrementRequests(ctx context.Context, id int64) error
}

// SQLApiKeyStore implements ApiKeyStore with sqlx.
type SQLApiKeyStore struct {
	DB     *sqlx.DB
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
	err := s.DB.GetContext(ctx, &k, fmt.Sprintf("SELECT * FROM %s WHERE key = $1", s.tableName()), key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &k, err
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

func (s *SQLApiKeyStore) Create(ctx context.Context, key *ApiKey) error {
	now := time.Now()
	if key.Key == "" {
		key.Key = GenerateAPIKey()
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s, key, name, %s, %s,
		%s, %s, %s, %s, %s,
		%s, %s, %s, enabled, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		s.tableName(),
		s.col("user_id"), s.col("scoped_client_id"), s.col("scoped_user_id"),
		s.col("daily_limit"), s.col("requests_today"), s.col("foundry_url"), s.col("foundry_username"),
		s.col("encrypted_foundry_password"), s.col("password_iv"), s.col("password_auth_tag"),
		s.col("expires_at"), s.col("created_at"), s.col("updated_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			key.UserID, key.Key, key.Name, key.ScopedClientID, key.ScopedUserID,
			key.DailyLimit, key.RequestsToday, key.FoundryURL, key.FoundryUsername,
			key.EncryptedFoundryPassword, key.PasswordIV, key.PasswordAuthTag,
			key.ExpiresAt, key.Enabled, now, now,
		).Scan(&key.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		key.UserID, key.Key, key.Name, key.ScopedClientID, key.ScopedUserID,
		key.DailyLimit, key.RequestsToday, key.FoundryURL, key.FoundryUsername,
		key.EncryptedFoundryPassword, key.PasswordIV, key.PasswordAuthTag,
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
		%s=$8, %s=$9, %s=$10,
		%s=$11, %s=$12, enabled=$13, %s=$14
		WHERE id=$15`,
		s.tableName(),
		s.col("scoped_client_id"), s.col("scoped_user_id"),
		s.col("daily_limit"), s.col("requests_today"), s.col("last_request_date"), s.col("foundry_url"),
		s.col("foundry_username"), s.col("encrypted_foundry_password"), s.col("password_iv"),
		s.col("password_auth_tag"), s.col("expires_at"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query,
		key.Name, key.ScopedClientID, key.ScopedUserID,
		key.DailyLimit, key.RequestsToday, key.LastRequestDate, key.FoundryURL,
		key.FoundryUsername, key.EncryptedFoundryPassword, key.PasswordIV,
		key.PasswordAuthTag, key.ExpiresAt, key.Enabled, key.UpdatedAt,
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

func (s *SQLApiKeyStore) ResetDailyCounters(ctx context.Context) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = 0, %s = NULL, %s = $1`,
		s.tableName(), s.col("requests_today"), s.col("last_request_date"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, time.Now())
	return err
}

func (s *SQLApiKeyStore) IncrementRequests(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = %s + 1, %s = $1, %s = $2 WHERE id = $3`,
		s.tableName(),
		s.col("requests_today"), s.col("requests_today"),
		s.col("last_request_date"), s.col("updated_at"))
	now := time.Now()
	_, err := s.DB.ExecContext(ctx, query, now, now, id)
	return err
}
