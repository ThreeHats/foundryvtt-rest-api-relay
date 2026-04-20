package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PairRequest status constants.
const (
	PairRequestStatusPending   = "pending"
	PairRequestStatusApproved  = "approved"
	PairRequestStatusDenied    = "denied"
	PairRequestStatusExchanged = "exchanged"
	PairRequestStatusExpired   = "expired"
)

// PairRequest represents a Foundry module's request to be paired with a relay account.
// The module initiates the request without auth; the relay user approves via the web UI;
// the module polls for the resulting pairing code to complete the connection.
type PairRequest struct {
	ID             int64          `db:"id" json:"id"`
	Code           string         `db:"code" json:"code"`
	WorldID        string         `db:"worldId" json:"worldId"`
	WorldTitle     string         `db:"worldTitle" json:"worldTitle"`
	SystemID       string         `db:"systemId" json:"systemId"`
	SystemTitle    string         `db:"systemTitle" json:"systemTitle"`
	SystemVersion  string         `db:"systemVersion" json:"systemVersion"`
	FoundryVersion string         `db:"foundryVersion" json:"foundryVersion"`
	// RequestedRemoteScopes and RequestedTargetClients are optional hints from
	// the module about what cross-world permissions it needs. The approving user
	// sees these pre-selected but can modify them before approving.
	RequestedRemoteScopes  sql.NullString `db:"requestedRemoteScopes" json:"requestedRemoteScopes"`
	RequestedTargetClients sql.NullString `db:"requestedTargetClients" json:"requestedTargetClients"`
	// UpgradeOnly, when true, means the world is already paired and the request
	// is only to grant/update cross-world permissions — no new pairing token is created.
	UpgradeOnly bool   `db:"upgradeOnly" json:"upgradeOnly"`
	Status      string `db:"status" json:"status"`
	// PairingCode is the 6-char code set on approval. The module exchanges it via POST /auth/pair.
	// Empty for upgradeOnly requests.
	PairingCode sql.NullString `db:"pairingCode" json:"-"`
	// UserID is the relay user who approved the request.
	UserID    sql.NullInt64 `db:"userId" json:"-"`
	ExpiresAt SQLiteTime    `db:"expiresAt" json:"expiresAt"`
	CreatedAt SQLiteTime    `db:"createdAt" json:"createdAt"`
}

// PairRequestStore defines operations on pair requests.
type PairRequestStore interface {
	FindByCode(ctx context.Context, code string) (*PairRequest, error)
	Create(ctx context.Context, req *PairRequest) error
	SetApproved(ctx context.Context, id int64, pairingCode string, userID int64) error
	SetStatus(ctx context.Context, id int64, status string) error
	CleanupExpired(ctx context.Context) (int64, error)
}

// SQLPairRequestStore implements PairRequestStore with sqlx.
type SQLPairRequestStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLPairRequestStore) tableName() string {
	if s.DBType == "sqlite" {
		return "PairRequests"
	}
	return `"PairRequests"`
}

func (s *SQLPairRequestStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLPairRequestStore) FindByCode(ctx context.Context, code string) (*PairRequest, error) {
	var r PairRequest
	err := s.DB.GetContext(ctx, &r,
		fmt.Sprintf("SELECT * FROM %s WHERE code = $1", s.tableName()), code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &r, err
}

func (s *SQLPairRequestStore) Create(ctx context.Context, req *PairRequest) error {
	now := time.Now()
	query := fmt.Sprintf(
		`INSERT INTO %s (code, %s, %s, %s, %s, %s, %s, %s, %s, %s, status, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		s.tableName(),
		s.col("world_id"), s.col("world_title"),
		s.col("system_id"), s.col("system_title"),
		s.col("system_version"), s.col("foundry_version"),
		s.col("requested_remote_scopes"), s.col("requested_target_clients"),
		s.col("upgrade_only"),
		s.col("expires_at"), s.col("created_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			req.Code, req.WorldID, req.WorldTitle,
			req.SystemID, req.SystemTitle, req.SystemVersion, req.FoundryVersion,
			req.RequestedRemoteScopes, req.RequestedTargetClients, req.UpgradeOnly,
			req.Status, req.ExpiresAt, now,
		).Scan(&req.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		req.Code, req.WorldID, req.WorldTitle,
		req.SystemID, req.SystemTitle, req.SystemVersion, req.FoundryVersion,
		req.RequestedRemoteScopes, req.RequestedTargetClients, req.UpgradeOnly,
		req.Status, req.ExpiresAt, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	req.ID = id
	req.CreatedAt = NewSQLiteTime(now)
	return nil
}

func (s *SQLPairRequestStore) SetApproved(ctx context.Context, id int64, pairingCode string, userID int64) error {
	query := fmt.Sprintf(`UPDATE %s SET status=$1, %s=$2, %s=$3 WHERE id=$4`,
		s.tableName(), s.col("pairing_code"), s.col("user_id"))
	_, err := s.DB.ExecContext(ctx, query, PairRequestStatusApproved, pairingCode, userID, id)
	return err
}

func (s *SQLPairRequestStore) SetStatus(ctx context.Context, id int64, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE id=$2`, s.tableName())
	_, err := s.DB.ExecContext(ctx, query, status, id)
	return err
}

func (s *SQLPairRequestStore) CleanupExpired(ctx context.Context) (int64, error) {
	nowSQL := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE status=$2 AND %s < $3`,
		s.tableName(), s.col("expires_at"))
	result, err := s.DB.ExecContext(ctx, query, PairRequestStatusExpired, PairRequestStatusPending, nowSQL)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
