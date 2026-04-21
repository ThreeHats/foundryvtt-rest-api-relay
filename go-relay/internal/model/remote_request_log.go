package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// RemoteRequestLog records a single cross-world remote-request invocation
// from one connected Foundry module to another. The relay writes one row per
// remote-request, regardless of whether it succeeded.
//
// This is the audit trail for the WS tunnel pattern: every cross-world action
// is attributable to a specific source token, target client, and timestamp,
// and is queryable from the dashboard's "Remote Requests" view.
type RemoteRequestLog struct {
	ID             int64          `db:"id" json:"id"`
	UserID         int64          `db:"userId" json:"userId"`
	SourceClientID string         `db:"sourceClientId" json:"sourceClientId"`
	SourceTokenID  int64          `db:"sourceTokenId" json:"sourceTokenId"`
	TargetClientID string         `db:"targetClientId" json:"targetClientId"`
	Action         string         `db:"action" json:"action"`
	Success        LooseBool      `db:"success" json:"success"`
	ErrorMessage   sql.NullString `db:"errorMessage" json:"errorMessage"`
	SourceIP       sql.NullString `db:"sourceIp" json:"sourceIp"`
	CreatedAt      SQLiteTime     `db:"createdAt" json:"createdAt"`
}

// RemoteRequestLogStore defines operations on the cross-world audit log.
type RemoteRequestLogStore interface {
	Create(ctx context.Context, log *RemoteRequestLog) error
	FindRecentByUser(ctx context.Context, userID int64, limit, offset int) ([]*RemoteRequestLog, int64, error)
	FindFiltered(ctx context.Context, filters ActivityFilters) ([]*RemoteRequestLog, int64, error)
	CleanupOlderThan(ctx context.Context, cutoff time.Time) (int64, error)
}

// SQLRemoteRequestLogStore implements RemoteRequestLogStore with sqlx.
type SQLRemoteRequestLogStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLRemoteRequestLogStore) tableName() string {
	if s.DBType == "sqlite" {
		return "RemoteRequestLogs"
	}
	return `"RemoteRequestLogs"`
}

func (s *SQLRemoteRequestLogStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLRemoteRequestLogStore) Create(ctx context.Context, l *RemoteRequestLog) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`INSERT INTO %s
		(%s, %s, %s, %s, action, success, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		s.tableName(),
		s.col("user_id"),
		s.col("source_client_id"),
		s.col("source_token_id"),
		s.col("target_client_id"),
		s.col("error_message"),
		s.col("source_ip"),
		s.col("created_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			l.UserID, l.SourceClientID, l.SourceTokenID, l.TargetClientID,
			l.Action, bool(l.Success), l.ErrorMessage, l.SourceIP, now,
		).Scan(&l.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		l.UserID, l.SourceClientID, l.SourceTokenID, l.TargetClientID,
		l.Action, l.Success, l.ErrorMessage, l.SourceIP, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	l.ID = id
	l.CreatedAt = now
	return nil
}

func (s *SQLRemoteRequestLogStore) FindRecentByUser(ctx context.Context, userID int64, limit, offset int) ([]*RemoteRequestLog, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	logs := []*RemoteRequestLog{}
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`,
		s.tableName(), s.col("user_id"))
	if err := s.DB.SelectContext(ctx, &logs, query, userID, limit, offset); err != nil {
		return nil, 0, err
	}
	var total int64
	if err := s.DB.GetContext(ctx, &total,
		fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s = $1`, s.tableName(), s.col("user_id")),
		userID); err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (s *SQLRemoteRequestLogStore) FindFiltered(ctx context.Context, filters ActivityFilters) ([]*RemoteRequestLog, int64, error) {
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
		addAnd(fmt.Sprintf("(%s = $%d OR %s = $%d)", s.col("source_client_id"), argN, s.col("target_client_id"), argN+1), filters.World, filters.World)
	}
	if filters.Action != "" {
		addAnd(fmt.Sprintf("action = $%d", argN), filters.Action)
	}
	if filters.Success != nil {
		v := 0
		if *filters.Success {
			v = 1
		}
		addAnd(fmt.Sprintf("success = $%d", argN), v)
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

	var logs []*RemoteRequestLog
	dataQ := fmt.Sprintf("SELECT * FROM %s%s ORDER BY id DESC LIMIT $%d",
		s.tableName(), where, argN)
	args = append(args, fetchN)
	if err := s.DB.SelectContext(ctx, &logs, dataQ, args...); err != nil {
		return nil, 0, err
	}
	if logs == nil {
		logs = []*RemoteRequestLog{}
	}
	return logs, total, nil
}

func (s *SQLRemoteRequestLogStore) CleanupOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	now := NewSQLiteTime(cutoff)
	result, err := s.DB.ExecContext(ctx,
		fmt.Sprintf(`DELETE FROM %s WHERE %s < $1`, s.tableName(), s.col("created_at")),
		now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
