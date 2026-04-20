package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// KeyRequest status constants.
const (
	KeyRequestStatusPending   = "pending"
	KeyRequestStatusApproved  = "approved"
	KeyRequestStatusDenied    = "denied"
	KeyRequestStatusExchanged = "exchanged"
	KeyRequestStatusExpired   = "expired"
)

// KeyRequest represents a third-party application's request for an API key.
type KeyRequest struct {
	ID                  int64          `db:"id" json:"id"`
	Code                string         `db:"code" json:"code"`
	AppName             string         `db:"appName" json:"appName"`
	AppDescription      string         `db:"appDescription" json:"appDescription"`
	AppURL              sql.NullString `db:"appUrl" json:"appUrl"`
	RequestedScopes     string         `db:"requestedScopes" json:"requestedScopes"`
	RequestedClientIDs  string         `db:"requestedClientIds" json:"requestedClientIds"`
	CallbackURL         sql.NullString `db:"callbackUrl" json:"callbackUrl"`
	SuggestedMonthlyLimit sql.NullInt64  `db:"suggestedMonthlyLimit" json:"suggestedMonthlyLimit"`
	SuggestedExpiry     sql.NullString `db:"suggestedExpiry" json:"suggestedExpiry"`
	Status              string         `db:"status" json:"status"`
	ApprovedKeyID       sql.NullInt64  `db:"approvedKeyId" json:"approvedKeyId"`
	ApprovedByID        sql.NullInt64  `db:"approvedById" json:"approvedById"`
	ExchangeCode        sql.NullString `db:"exchangeCode" json:"exchangeCode"`
	ExpiresAt           SQLiteTime     `db:"expiresAt" json:"expiresAt"`
	CreatedAt           SQLiteTime     `db:"createdAt" json:"createdAt"`
	UpdatedAt           SQLiteTime     `db:"updatedAt" json:"updatedAt"`
}

// KeyRequestStore defines operations on key requests.
type KeyRequestStore interface {
	FindByCode(ctx context.Context, code string) (*KeyRequest, error)
	Create(ctx context.Context, req *KeyRequest) error
	UpdateStatus(ctx context.Context, id int64, status string, updates map[string]interface{}) error
	CleanupExpired(ctx context.Context) (int64, error)
}

// SQLKeyRequestStore implements KeyRequestStore with sqlx.
type SQLKeyRequestStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLKeyRequestStore) tableName() string {
	if s.DBType == "sqlite" {
		return "KeyRequests"
	}
	return `"KeyRequests"`
}

func (s *SQLKeyRequestStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLKeyRequestStore) FindByCode(ctx context.Context, code string) (*KeyRequest, error) {
	var r KeyRequest
	err := s.DB.GetContext(ctx, &r, fmt.Sprintf("SELECT * FROM %s WHERE code = $1", s.tableName()), code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &r, err
}

func (s *SQLKeyRequestStore) Create(ctx context.Context, req *KeyRequest) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (code, %s, %s, %s, %s, %s, %s, %s, %s, status, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		s.tableName(),
		s.col("app_name"), s.col("app_description"), s.col("app_url"),
		s.col("requested_scopes"), s.col("requested_client_ids"), s.col("callback_url"),
		s.col("suggested_monthly_limit"), s.col("suggested_expiry"),
		s.col("expires_at"), s.col("created_at"), s.col("updated_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			req.Code, req.AppName, req.AppDescription, req.AppURL,
			req.RequestedScopes, req.RequestedClientIDs, req.CallbackURL,
			req.SuggestedMonthlyLimit, req.SuggestedExpiry, req.Status,
			req.ExpiresAt, now, now,
		).Scan(&req.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		req.Code, req.AppName, req.AppDescription, req.AppURL,
		req.RequestedScopes, req.RequestedClientIDs, req.CallbackURL,
		req.SuggestedMonthlyLimit, req.SuggestedExpiry, req.Status,
		req.ExpiresAt, now, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	req.ID = id
	req.CreatedAt = NewSQLiteTime(now)
	req.UpdatedAt = NewSQLiteTime(now)
	return nil
}

func (s *SQLKeyRequestStore) UpdateStatus(ctx context.Context, id int64, status string, updates map[string]interface{}) error {
	now := time.Now()

	// Build SET clause dynamically from updates map
	setClauses := fmt.Sprintf("status=$1, %s=$2", s.col("updated_at"))
	args := []interface{}{status, now}
	paramIdx := 3

	for key, val := range updates {
		setClauses += fmt.Sprintf(", %s=$%d", s.col(key), paramIdx)
		args = append(args, val)
		paramIdx++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", s.tableName(), setClauses, paramIdx)
	_, err := s.DB.ExecContext(ctx, query, args...)
	return err
}

func (s *SQLKeyRequestStore) CleanupExpired(ctx context.Context) (int64, error) {
	// Use SQLiteTime for the expires_at comparison so the parameter is formatted
	// with the same string format as the stored value. See pairing_code.go for
	// the same fix and explanation.
	nowSQL := NewSQLiteTime(time.Now())
	now := time.Now()
	query := fmt.Sprintf(`UPDATE %s SET status=$1, %s=$2 WHERE status=$3 AND %s < $4`,
		s.tableName(), s.col("updated_at"), s.col("expires_at"))
	result, err := s.DB.ExecContext(ctx, query, KeyRequestStatusExpired, now, KeyRequestStatusPending, nowSQL)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
