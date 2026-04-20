package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// AdminAuditLog records administrative actions taken via the admin dashboard.
// Used for security audit trail and accountability.
type AdminAuditLog struct {
	ID          int64          `db:"id" json:"id"`
	AdminUserID int64          `db:"adminUserId" json:"adminUserId"`
	Action      string         `db:"action" json:"action"`
	TargetType  string         `db:"targetType" json:"targetType"`
	TargetID    sql.NullString `db:"targetId" json:"targetId"`
	Details     sql.NullString `db:"details" json:"details"`
	IPAddress   sql.NullString `db:"ipAddress" json:"ipAddress"`
	CreatedAt   SQLiteTime     `db:"createdAt" json:"createdAt"`
}

// AuditLogFilters narrows audit log queries.
type AuditLogFilters struct {
	AdminUserID int64
	Action      string
	TargetType  string
	Since       *time.Time
	Until       *time.Time
	Offset      int
	Limit       int
}

// AuditLogStore defines operations on admin audit logs.
type AuditLogStore interface {
	Create(ctx context.Context, entry *AdminAuditLog) error
	FindAll(ctx context.Context, filters AuditLogFilters) ([]*AdminAuditLog, int64, error)
	DeleteByAdminUserID(ctx context.Context, adminUserID int64) (int64, error)
	CleanupOlderThan(ctx context.Context, days int) (int64, error)
}

// SQLAuditLogStore implements AuditLogStore with sqlx.
type SQLAuditLogStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLAuditLogStore) tableName() string {
	if s.DBType == "sqlite" {
		return "AdminAuditLogs"
	}
	return `"AdminAuditLogs"`
}

func (s *SQLAuditLogStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLAuditLogStore) Create(ctx context.Context, entry *AdminAuditLog) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (%s, action, %s, %s, details, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		s.tableName(),
		s.col("admin_user_id"), s.col("target_type"), s.col("target_id"),
		s.col("ip_address"), s.col("created_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			entry.AdminUserID, entry.Action, entry.TargetType, entry.TargetID,
			entry.Details, entry.IPAddress, now,
		).Scan(&entry.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		entry.AdminUserID, entry.Action, entry.TargetType, entry.TargetID,
		entry.Details, entry.IPAddress, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	entry.ID = id
	entry.CreatedAt = NewSQLiteTime(now)
	return nil
}

func (s *SQLAuditLogStore) FindAll(ctx context.Context, filters AuditLogFilters) ([]*AdminAuditLog, int64, error) {
	where := []string{}
	args := []interface{}{}
	idx := 1

	if filters.AdminUserID != 0 {
		where = append(where, fmt.Sprintf("%s = $%d", s.col("admin_user_id"), idx))
		args = append(args, filters.AdminUserID)
		idx++
	}
	if filters.Action != "" {
		where = append(where, fmt.Sprintf("action = $%d", idx))
		args = append(args, filters.Action)
		idx++
	}
	if filters.TargetType != "" {
		where = append(where, fmt.Sprintf("%s = $%d", s.col("target_type"), idx))
		args = append(args, filters.TargetType)
		idx++
	}
	if filters.Since != nil {
		where = append(where, fmt.Sprintf("%s >= $%d", s.col("created_at"), idx))
		args = append(args, *filters.Since)
		idx++
	}
	if filters.Until != nil {
		where = append(where, fmt.Sprintf("%s <= $%d", s.col("created_at"), idx))
		args = append(args, *filters.Until)
		idx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	limit := filters.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	listQuery := fmt.Sprintf(`SELECT * FROM %s %s ORDER BY %s DESC LIMIT $%d OFFSET $%d`,
		s.tableName(), whereClause, s.col("created_at"), idx, idx+1)
	listArgs := append(args, limit, offset)

	logs := []*AdminAuditLog{}
	if err := s.DB.SelectContext(ctx, &logs, listQuery, listArgs...); err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, s.tableName(), whereClause)
	var total int64
	if err := s.DB.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (s *SQLAuditLogStore) DeleteByAdminUserID(ctx context.Context, adminUserID int64) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", s.tableName(), s.col("admin_user_id"))
	result, err := s.DB.ExecContext(ctx, query, adminUserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *SQLAuditLogStore) CleanupOlderThan(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s < $1",
		s.tableName(), s.col("created_at"))
	result, err := s.DB.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
