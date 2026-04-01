package model

import (
	"context"
	"crypto/rand"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

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
func (t SQLiteTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.Format("2006-01-02 15:04:05"), nil
}

// NewSQLiteTime creates a SQLiteTime from a time.Time.
func NewSQLiteTime(t time.Time) SQLiteTime {
	return SQLiteTime{Time: t, Valid: true}
}

// User represents a registered user.
// db tags use camelCase to match Sequelize-created SQLite columns.
// The col() helper maps these to snake_case for PostgreSQL queries.
type User struct {
	ID                  int64          `db:"id" json:"id"`
	Email               string         `db:"email" json:"email"`
	Password            string         `db:"password" json:"-"`
	APIKey              string         `db:"apiKey" json:"apiKey"`
	RequestsThisMonth   int            `db:"requestsThisMonth" json:"requestsThisMonth"`
	RequestsToday       int            `db:"requestsToday" json:"requestsToday"`
	LastRequestDate     *SQLiteTime    `db:"lastRequestDate" json:"-"`
	StripeCustomerID    sql.NullString `db:"stripeCustomerId" json:"-"`
	SubscriptionStatus  sql.NullString `db:"subscriptionStatus" json:"subscriptionStatus"`
	SubscriptionID      sql.NullString `db:"subscriptionId" json:"-"`
	SubscriptionEndsAt  *SQLiteTime    `db:"subscriptionEndsAt" json:"-"`
	MaxHeadlessSessions sql.NullInt64  `db:"maxHeadlessSessions" json:"-"`
	CreatedAt           SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt           SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// GetSubscriptionStatus returns the subscription status string, defaulting to "free".
func (u *User) GetSubscriptionStatus() string {
	if u.SubscriptionStatus.Valid {
		return u.SubscriptionStatus.String
	}
	return "free"
}

// GenerateAPIKey creates a new random 64-character hex API key.
func GenerateAPIKey() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to generate API key: %v", err))
	}
	return hex.EncodeToString(b)
}

// UserStore defines operations on users.
type UserStore interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByAPIKey(ctx context.Context, apiKey string) (*User, error)
	FindByStripeCustomerID(ctx context.Context, customerID string) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	IncrementRequests(ctx context.Context, id int64) error
	ResetMonthlyRequests(ctx context.Context) error
	ResetDailyRequests(ctx context.Context) error
}

// SQLUserStore implements UserStore with sqlx.
type SQLUserStore struct {
	DB     *sqlx.DB
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
	var user User
	err := s.DB.GetContext(ctx, &user, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("api_key")), apiKey)
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

func (s *SQLUserStore) FindAll(ctx context.Context) ([]*User, error) {
	var users []*User
	err := s.DB.SelectContext(ctx, &users, fmt.Sprintf("SELECT * FROM %s", s.tableName()))
	return users, err
}

func (s *SQLUserStore) Create(ctx context.Context, user *User) error {
	query := fmt.Sprintf(`INSERT INTO %s (email, password, api_key, requests_this_month, requests_today, subscription_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, s.tableName())

	now := time.Now()
	status := "free"
	if user.SubscriptionStatus.Valid {
		status = user.SubscriptionStatus.String
	}

	if s.DBType == "sqlite" {
		query = fmt.Sprintf(`INSERT INTO Users (email, password, %s, %s, %s, %s, %s, %s)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			s.col("api_key"), s.col("requests_this_month"), s.col("requests_today"),
			s.col("subscription_status"), s.col("created_at"), s.col("updated_at"))
		result, err := s.DB.ExecContext(ctx, query,
			user.Email, user.Password, user.APIKey, user.RequestsThisMonth, user.RequestsToday, status, now, now)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		user.ID = id
		user.CreatedAt = NewSQLiteTime(now)
		user.UpdatedAt = NewSQLiteTime(now)
		return nil
	}

	return s.DB.QueryRowContext(ctx, query,
		user.Email, user.Password, user.APIKey, user.RequestsThisMonth, user.RequestsToday, status, now, now,
	).Scan(&user.ID)
}

func (s *SQLUserStore) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET email=$1, password=$2, %s=$3, %s=$4,
		%s=$5, %s=$6, %s=$7, %s=$8,
		%s=$9, %s=$10, %s=$11, %s=$12
		WHERE id=$13`,
		s.tableName(),
		s.col("api_key"), s.col("requests_this_month"),
		s.col("requests_today"), s.col("last_request_date"),
		s.col("stripe_customer_id"), s.col("subscription_status"),
		s.col("subscription_id"), s.col("subscription_ends_at"),
		s.col("max_headless_sessions"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query,
		user.Email, user.Password, user.APIKey, user.RequestsThisMonth,
		user.RequestsToday, user.LastRequestDate, user.StripeCustomerID, user.SubscriptionStatus,
		user.SubscriptionID, user.SubscriptionEndsAt, user.MaxHeadlessSessions, user.UpdatedAt,
		user.ID)
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
