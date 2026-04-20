package model

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"time"
)

// HashAPIKey returns the SHA-256 hex digest of an API key. The relay stores
// only this hash in the Users table; the plaintext key exists momentarily in
// the registration/regeneration response and never again. Auth middleware
// hashes incoming x-api-key header values and looks them up by hash.
//
// SHA-256 (vs bcrypt) is appropriate here because the input is already a
// 256-bit cryptographically random token. There's nothing to brute-force; we
// just need a one-way function fast enough to run on every authenticated
// request.
func HashAPIKey(plaintext string) string {
	if plaintext == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}

// SQLiteTime handles both time.Time (PostgreSQL) and string (SQLite) scanning.
type SQLiteTime struct {
	Time  time.Time
	Valid bool
}

func (t *SQLiteTime) Scan(src interface{}) error {
	if src == nil {
		t.Valid = false
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		t.Time = v
		t.Valid = true
	case string:
		// Parse as UTC. SQLiteTime.Value() always formats in UTC, so when reading
		// back from a SQLite TEXT column we know the bare "2006-01-02 15:04:05"
		// string represents UTC. Using time.Parse() without ParseInLocation defaults
		// to UTC, which is correct for our format. We previously had a timezone bug
		// where Value() formatted in local time but Scan() parsed as UTC, causing
		// timestamps to shift by the local UTC offset on every write→read roundtrip.
		parsed, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			parsed, err = time.Parse(time.RFC3339, v)
		}
		if err != nil {
			parsed, err = time.Parse("2006-01-02", v)
		}
		if err != nil {
			t.Valid = false
			return nil
		}
		t.Time = parsed
		t.Valid = true
	default:
		t.Valid = false
	}
	return nil
}

// Value implements driver.Valuer for database writes.
//
// IMPORTANT: We force UTC here. The Time.Format() method formats in whatever
// timezone the time.Time is in, but Scan() parses bare timestamp strings as
// UTC. If we stored local time and read it back as UTC the timestamp would
// silently shift by the local UTC offset (causing pairing codes to look
// "expired" immediately for users in non-UTC timezones).
func (t SQLiteTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.UTC().Format("2006-01-02 15:04:05"), nil
}

// MarshalJSON implements json.Marshaler. Returns the time as an ISO 8601
// string (e.g. "2026-04-07T06:30:00Z") or null if invalid. Without this,
// SQLiteTime would serialize as the default Go struct shape
// {"Time":"...","Valid":true}, which the frontend can't parse with new Date().
func (t SQLiteTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	// RFC 3339 / ISO 8601 — directly parseable by JavaScript's Date constructor
	formatted := t.Time.UTC().Format(time.RFC3339)
	return []byte(`"` + formatted + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler. Accepts ISO 8601 strings or null.
func (t *SQLiteTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == `""` {
		t.Valid = false
		return nil
	}
	// Strip quotes
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// Fall back to other formats
		parsed, err = time.Parse("2006-01-02 15:04:05", s)
	}
	if err != nil {
		t.Valid = false
		return nil
	}
	t.Time = parsed
	t.Valid = true
	return nil
}

// NewSQLiteTime creates a SQLiteTime from a time.Time.
func NewSQLiteTime(t time.Time) SQLiteTime {
	return SQLiteTime{Time: t, Valid: true}
}

// LooseBool is a database/sql Scanner that accepts a wide range of values
// (int64, []byte, string, bool, nil) and coerces them to bool. Any value
// other than 1, "1", "true", or true is treated as false.
//
// This exists to defend against corrupted data in SQLite columns where the
// underlying value might be a string or timestamp due to a parameter binding
// bug. The standard library's bool scanner errors on such values.
type LooseBool bool

// Scan implements sql.Scanner.
func (b *LooseBool) Scan(src interface{}) error {
	if src == nil {
		*b = false
		return nil
	}
	switch v := src.(type) {
	case bool:
		*b = LooseBool(v)
	case int64:
		*b = v != 0
	case float64:
		*b = v != 0
	case []byte:
		s := string(v)
		*b = s == "1" || s == "true" || s == "TRUE" || s == "True"
	case string:
		*b = v == "1" || v == "true" || v == "TRUE" || v == "True"
	default:
		*b = false
	}
	return nil
}

// Value implements driver.Valuer.
func (b LooseBool) Value() (driver.Value, error) {
	if bool(b) {
		return int64(1), nil
	}
	return int64(0), nil
}

// MarshalJSON serializes as a regular JSON boolean.
func (b LooseBool) MarshalJSON() ([]byte, error) {
	if bool(b) {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// UnmarshalJSON accepts a JSON boolean.
func (b *LooseBool) UnmarshalJSON(data []byte) error {
	s := string(data)
	*b = LooseBool(s == "true" || s == "1")
	return nil
}

// User represents a registered user.
// db tags use camelCase to match Sequelize-created SQLite columns.
// The col() helper maps these to snake_case for PostgreSQL queries.
type User struct {
	ID                  int64          `db:"id" json:"id"`
	Email               string         `db:"email" json:"email"`
	Password            string         `db:"password" json:"-"`
	// APIKey is a TRANSIENT field — never persisted to the DB (db:"-"). It
	// holds the plaintext master API key only momentarily during the
	// registration / rotation flows so the response handler can return it to
	// the caller exactly once. After Create() or Update() runs, the plaintext
	// is hashed into APIKeyHash and the plaintext value is meaningless to
	// keep on the struct (but it stays around in memory until the request
	// goroutine returns).
	APIKey              string         `db:"-" json:"-"`
	// APIKeyHash is the SHA-256 hex digest of the master API key. This is the
	// only thing stored in the DB; lookups by API key hash this column.
	APIKeyHash          sql.NullString `db:"apiKeyHash" json:"-"`
	RequestsThisMonth   int            `db:"requestsThisMonth" json:"requestsThisMonth"`
	RequestsToday       int            `db:"requestsToday" json:"requestsToday"`
	LastRequestDate     *SQLiteTime    `db:"lastRequestDate" json:"-"`
	StripeCustomerID    sql.NullString `db:"stripeCustomerId" json:"-"`
	SubscriptionStatus  sql.NullString `db:"subscriptionStatus" json:"subscriptionStatus"`
	SubscriptionID      sql.NullString `db:"subscriptionId" json:"-"`
	SubscriptionEndsAt  *SQLiteTime    `db:"subscriptionEndsAt" json:"-"`
	MaxHeadlessSessions            sql.NullInt64  `db:"maxHeadlessSessions" json:"-"`
	EmailVerified                  bool           `db:"emailVerified" json:"emailVerified"`
	VerificationTokenHash          sql.NullString `db:"verificationTokenHash" json:"-"`
	VerificationTokenExpiresAt     *SQLiteTime    `db:"verificationTokenExpiresAt" json:"-"`
	Role                           string         `db:"role" json:"role"`
	Disabled                       bool           `db:"disabled" json:"disabled"`
	APIKeyRotationRequired         bool           `db:"apiKeyRotationRequired" json:"apiKeyRotationRequired"`
	CreatedAt                      SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt                      SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// IsAdmin returns true if this user has the admin role.
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// GetSubscriptionStatus returns the subscription status string, defaulting to "free".
func (u *User) GetSubscriptionStatus() string {
	if u.SubscriptionStatus.Valid {
		return u.SubscriptionStatus.String
	}
	return "free"
}

// GenerateAPIKey creates a new random 64-character hex API key.
func GenerateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate API key: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// UserStore defines operations on users.
type UserStore interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	// FindByAPIKey hashes the input API key and looks up by hash.
	// Used by HTTP auth middleware.
	FindByAPIKey(ctx context.Context, apiKey string) (*User, error)
	// FindByAPIKeyHash takes a pre-computed hex SHA-256 hash and looks it up
	// directly. Used by callers that already hold the hash (e.g., the WS
	// reconciliation loop, which compares against ClientManager registration
	// tokens that ARE the hash).
	FindByAPIKeyHash(ctx context.Context, hash string) (*User, error)
	FindByStripeCustomerID(ctx context.Context, customerID string) (*User, error)
	FindAllPaginated(ctx context.Context, offset, limit int) ([]*User, int64, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	SetDisabled(ctx context.Context, id int64, disabled bool) error
	SetRole(ctx context.Context, id int64, role string) error
	IncrementRequests(ctx context.Context, id int64) error
	IncrementRequestsBy(ctx context.Context, id int64, count int) error
	ResetMonthlyRequests(ctx context.Context) error
	ResetDailyRequests(ctx context.Context) error
}

// SQLUserStore implements UserStore with sqlx.
type SQLUserStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLUserStore) tableName() string {
	if s.DBType == "sqlite" {
		return "Users"
	}
	return `"Users"`
}

func (s *SQLUserStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLUserStore) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := s.DB.GetContext(ctx, &user, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", s.tableName()), id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (s *SQLUserStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.DB.GetContext(ctx, &user, fmt.Sprintf("SELECT * FROM %s WHERE email = $1", s.tableName()), email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (s *SQLUserStore) FindByAPIKey(ctx context.Context, apiKey string) (*User, error) {
	if apiKey == "" {
		return nil, nil
	}
	return s.FindByAPIKeyHash(ctx, HashAPIKey(apiKey))
}

func (s *SQLUserStore) FindByAPIKeyHash(ctx context.Context, hash string) (*User, error) {
	if hash == "" {
		return nil, nil
	}
	var user User
	err := s.DB.GetContext(ctx, &user,
		fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("api_key_hash")),
		hash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (s *SQLUserStore) FindByStripeCustomerID(ctx context.Context, customerID string) (*User, error) {
	var user User
	err := s.DB.GetContext(ctx, &user, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("stripe_customer_id")), customerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (s *SQLUserStore) Create(ctx context.Context, user *User) error {
	now := time.Now()
	status := "free"
	if user.SubscriptionStatus.Valid {
		status = user.SubscriptionStatus.String
	}

	role := user.Role
	if role == "" {
		role = "user"
	}

	// Always derive and populate the hash on create. The plaintext stays in
	// the (transient) APIKey field momentarily so the response handler can
	// return it once, but it's never written to the DB.
	if user.APIKey != "" {
		user.APIKeyHash = sql.NullString{String: HashAPIKey(user.APIKey), Valid: true}
	}

	if s.DBType == "sqlite" {
		query := fmt.Sprintf(`INSERT INTO Users (email, password, %s, %s, %s, %s, %s, %s, %s, role, disabled, %s, %s, %s)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
			s.col("api_key_hash"), s.col("requests_this_month"), s.col("requests_today"),
			s.col("subscription_status"), s.col("email_verified"),
			s.col("verification_token_hash"), s.col("verification_token_expires_at"),
			s.col("api_key_rotation_required"),
			s.col("created_at"), s.col("updated_at"))
		result, err := s.DB.ExecContext(ctx, query,
			user.Email, user.Password, user.APIKeyHash, user.RequestsThisMonth, user.RequestsToday, status,
			user.EmailVerified, user.VerificationTokenHash, user.VerificationTokenExpiresAt,
			role, user.Disabled, user.APIKeyRotationRequired,
			now, now)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		user.ID = id
		user.Role = role
		user.CreatedAt = NewSQLiteTime(now)
		user.UpdatedAt = NewSQLiteTime(now)
		return nil
	}

	query := fmt.Sprintf(`INSERT INTO %s (email, password, %s, %s, %s, %s, %s, %s, %s, role, disabled, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`,
		s.tableName(),
		s.col("api_key_hash"), s.col("requests_this_month"), s.col("requests_today"),
		s.col("subscription_status"), s.col("email_verified"),
		s.col("verification_token_hash"), s.col("verification_token_expires_at"),
		s.col("api_key_rotation_required"),
		s.col("created_at"), s.col("updated_at"))
	if err := s.DB.QueryRowContext(ctx, query,
		user.Email, user.Password, user.APIKeyHash, user.RequestsThisMonth, user.RequestsToday, status,
		user.EmailVerified, user.VerificationTokenHash, user.VerificationTokenExpiresAt,
		role, user.Disabled, user.APIKeyRotationRequired,
		now, now,
	).Scan(&user.ID); err != nil {
		return err
	}
	user.Role = role
	return nil
}

func (s *SQLUserStore) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = NewSQLiteTime(time.Now())
	if user.Role == "" {
		user.Role = "user"
	}
	// Re-derive the hash from APIKey on every update so that callers who
	// rotate the master key (which writes user.APIKey = newKey) automatically
	// get the matching hash without having to set it explicitly.
	if user.APIKey != "" {
		user.APIKeyHash = sql.NullString{String: HashAPIKey(user.APIKey), Valid: true}
	}
	query := fmt.Sprintf(`UPDATE %s SET email=$1, password=$2, %s=$3, %s=$4,
		%s=$5, %s=$6, %s=$7, %s=$8,
		%s=$9, %s=$10, %s=$11, %s=$12,
		%s=$13, %s=$14, %s=$15,
		role=$16, disabled=$17, %s=$18
		WHERE id=$19`,
		s.tableName(),
		s.col("api_key_hash"), s.col("requests_this_month"),
		s.col("requests_today"), s.col("last_request_date"),
		s.col("stripe_customer_id"), s.col("subscription_status"),
		s.col("subscription_id"), s.col("subscription_ends_at"),
		s.col("max_headless_sessions"), s.col("updated_at"),
		s.col("email_verified"), s.col("verification_token_hash"),
		s.col("verification_token_expires_at"),
		s.col("api_key_rotation_required"))
	_, err := s.DB.ExecContext(ctx, query,
		user.Email, user.Password, user.APIKeyHash, user.RequestsThisMonth,
		user.RequestsToday, user.LastRequestDate, user.StripeCustomerID, user.SubscriptionStatus,
		user.SubscriptionID, user.SubscriptionEndsAt, user.MaxHeadlessSessions, user.UpdatedAt,
		user.EmailVerified, user.VerificationTokenHash, user.VerificationTokenExpiresAt,
		user.Role, user.Disabled, user.APIKeyRotationRequired,
		user.ID)
	return err
}

// FindAllPaginated returns a page of users plus the total count.
func (s *SQLUserStore) FindAllPaginated(ctx context.Context, offset, limit int) ([]*User, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	users := []*User{}
	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY id ASC LIMIT $1 OFFSET $2`, s.tableName())
	if err := s.DB.SelectContext(ctx, &users, query, limit, offset); err != nil {
		return nil, 0, err
	}
	var total int64
	if err := s.DB.GetContext(ctx, &total, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, s.tableName())); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// SetDisabled toggles the disabled flag for a user.
func (s *SQLUserStore) SetDisabled(ctx context.Context, id int64, disabled bool) error {
	query := fmt.Sprintf(`UPDATE %s SET disabled=$1, %s=$2 WHERE id=$3`,
		s.tableName(), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, disabled, time.Now(), id)
	return err
}

// SetRole updates the role for a user.
func (s *SQLUserStore) SetRole(ctx context.Context, id int64, role string) error {
	if role == "" {
		role = "user"
	}
	query := fmt.Sprintf(`UPDATE %s SET role=$1, %s=$2 WHERE id=$3`,
		s.tableName(), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, role, time.Now(), id)
	return err
}

func (s *SQLUserStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName()), id)
	return err
}

func (s *SQLUserStore) IncrementRequests(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = %s + 1,
		%s = %s + 1, %s = $1, %s = $2 WHERE id = $3`,
		s.tableName(),
		s.col("requests_this_month"), s.col("requests_this_month"),
		s.col("requests_today"), s.col("requests_today"),
		s.col("last_request_date"), s.col("updated_at"))
	now := time.Now()
	_, err := s.DB.ExecContext(ctx, query, now, now, id)
	return err
}

func (s *SQLUserStore) IncrementRequestsBy(ctx context.Context, id int64, count int) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = %s + $1,
		%s = %s + $1, %s = $2, %s = $3 WHERE id = $4`,
		s.tableName(),
		s.col("requests_this_month"), s.col("requests_this_month"),
		s.col("requests_today"), s.col("requests_today"),
		s.col("last_request_date"), s.col("updated_at"))
	now := time.Now()
	_, err := s.DB.ExecContext(ctx, query, count, now, now, id)
	return err
}

func (s *SQLUserStore) ResetMonthlyRequests(ctx context.Context) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = 0, %s = 0, %s = $1`,
		s.tableName(), s.col("requests_this_month"), s.col("requests_today"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, time.Now())
	return err
}

func (s *SQLUserStore) ResetDailyRequests(ctx context.Context) error {
	query := fmt.Sprintf(`UPDATE %s SET %s = 0, %s = $1`,
		s.tableName(), s.col("requests_today"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, time.Now())
	return err
}
