package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ConnectionLog records metadata about each WebSocket connection for security auditing.
// db tags use camelCase to match Sequelize-created SQLite columns.
type ConnectionLog struct {
	ID             int64          `db:"id" json:"id"`
	UserID         int64          `db:"userId" json:"userId"`
	ClientID       string         `db:"clientId" json:"clientId"`
	TokenName      sql.NullString `db:"tokenName" json:"tokenName"`
	IPAddress      sql.NullString `db:"ipAddress" json:"ipAddress"`
	UserAgent      sql.NullString `db:"userAgent" json:"userAgent"`
	WorldID        sql.NullString `db:"worldId" json:"worldId"`
	WorldTitle     sql.NullString `db:"worldTitle" json:"worldTitle"`
	SystemID       sql.NullString `db:"systemId" json:"systemId"`
	FoundryVersion sql.NullString `db:"foundryVersion" json:"foundryVersion"`
	MetadataMatch  bool           `db:"metadataMatch" json:"metadataMatch"`
	Flagged        bool           `db:"flagged" json:"flagged"`
	FlagReason     sql.NullString `db:"flagReason" json:"flagReason"`
	CreatedAt      SQLiteTime     `db:"createdAt" json:"createdAt"`
}

// ConnectionLogStore defines operations on connection logs.
type ConnectionLogStore interface {
	Create(ctx context.Context, log *ConnectionLog) error
	FindByUser(ctx context.Context, userID int64, limit, offset int) ([]*ConnectionLog, error)
	FindFiltered(ctx context.Context, filters ActivityFilters) ([]*ConnectionLog, int64, error)
	CleanupOlderThan(ctx context.Context, days int) (int64, error)
}

// SQLConnectionLogStore implements ConnectionLogStore with sqlx.
type SQLConnectionLogStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLConnectionLogStore) tableName() string {
	if s.DBType == "sqlite" {
		return "ConnectionLogs"
	}
	return `"ConnectionLogs"`
}

func (s *SQLConnectionLogStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLConnectionLogStore) boolTrue() string {
	if s.DBType == "sqlite" {
		return "1"
	}
	return "TRUE"
}

func (s *SQLConnectionLogStore) boolFalse() string {
	if s.DBType == "sqlite" {
		return "0"
	}
	return "FALSE"
}

func (s *SQLConnectionLogStore) Create(ctx context.Context, connLog *ConnectionLog) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, flagged, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		s.tableName(),
		s.col("user_id"), s.col("client_id"), s.col("token_name"), s.col("ip_address"), s.col("user_agent"),
		s.col("world_id"), s.col("world_title"), s.col("system_id"), s.col("foundry_version"),
		s.col("metadata_match"), s.col("flag_reason"), s.col("created_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			connLog.UserID, connLog.ClientID, connLog.TokenName, connLog.IPAddress, connLog.UserAgent,
			connLog.WorldID, connLog.WorldTitle, connLog.SystemID, connLog.FoundryVersion,
			connLog.MetadataMatch, connLog.Flagged, connLog.FlagReason, now,
		).Scan(&connLog.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		connLog.UserID, connLog.ClientID, connLog.TokenName, connLog.IPAddress, connLog.UserAgent,
		connLog.WorldID, connLog.WorldTitle, connLog.SystemID, connLog.FoundryVersion,
		connLog.MetadataMatch, connLog.Flagged, connLog.FlagReason, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	connLog.ID = id
	connLog.CreatedAt = now
	return nil
}

func (s *SQLConnectionLogStore) FindByUser(ctx context.Context, userID int64, limit, offset int) ([]*ConnectionLog, error) {
	var logs []*ConnectionLog
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 ORDER BY %s DESC LIMIT $2 OFFSET $3",
		s.tableName(), s.col("user_id"), s.col("created_at"))
	err := s.DB.SelectContext(ctx, &logs, query, userID, limit, offset)
	if logs == nil {
		logs = []*ConnectionLog{}
	}
	return logs, err
}

func (s *SQLConnectionLogStore) FindFiltered(ctx context.Context, filters ActivityFilters) ([]*ConnectionLog, int64, error) {
	limit := filters.Limit
	if limit <= 0 {
		limit = 50
	}
	fetchN := filters.Offset + limit

	args := []interface{}{}
	where := ""
	argN := 1

	addAnd := func(clause string, vals ...interface{}) {
		if where == "" {
			where = " WHERE "
		} else {
			where += " AND "
		}
		where += clause
		args = append(args, vals...)
		argN += len(vals)
	}

	if filters.UserID != 0 {
		addAnd(fmt.Sprintf("%s = $%d", s.col("user_id"), argN), filters.UserID)
	}
	if filters.World != "" {
		addAnd(fmt.Sprintf("%s = $%d", s.col("client_id"), argN), filters.World)
	}
	if !filters.Since.IsZero() {
		addAnd(fmt.Sprintf("%s >= $%d", s.col("created_at"), argN), NewSQLiteTime(filters.Since))
	}
	if !filters.Until.IsZero() {
		addAnd(fmt.Sprintf("%s <= $%d", s.col("created_at"), argN), NewSQLiteTime(filters.Until))
	}

	var total int64
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", s.tableName(), where)
	if err := s.DB.GetContext(ctx, &total, countQ, args...); err != nil {
		return nil, 0, err
	}

	var logs []*ConnectionLog
	dataQ := fmt.Sprintf("SELECT * FROM %s%s ORDER BY %s DESC LIMIT $%d",
		s.tableName(), where, s.col("created_at"), argN)
	args = append(args, fetchN)
	if err := s.DB.SelectContext(ctx, &logs, dataQ, args...); err != nil {
		return nil, 0, err
	}
	if logs == nil {
		logs = []*ConnectionLog{}
	}
	return logs, total, nil
}

func (s *SQLConnectionLogStore) CleanupOlderThan(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s < $1",
		s.tableName(), s.col("created_at"))
	result, err := s.DB.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
