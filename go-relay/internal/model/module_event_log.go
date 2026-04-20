package model

import (
	"context"
	"fmt"
	"time"
)

// ModuleEventLog records a Foundry module event (execute-js, macro-execute,
// settings-change) that was reported to the relay via the module-notify WS message.
type ModuleEventLog struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"userId" json:"userId"`
	ClientID    string     `db:"clientId" json:"clientId"`
	WorldTitle  string     `db:"worldTitle" json:"worldTitle"`
	EventType   string     `db:"eventType" json:"eventType"` // "execute-js" | "macro-execute" | "settings-change"
	Actor       string     `db:"actor" json:"actor"`
	Description string     `db:"description" json:"description"`
	CreatedAt   SQLiteTime `db:"createdAt" json:"createdAt"`
}

// ModuleEventLogStore defines operations on the module event log.
type ModuleEventLogStore interface {
	Create(ctx context.Context, log *ModuleEventLog) error
	FindFiltered(ctx context.Context, filters ActivityFilters) ([]*ModuleEventLog, int64, error)
	CleanupOlderThan(ctx context.Context, cutoff time.Time) (int64, error)
}

// SQLModuleEventLogStore implements ModuleEventLogStore with sqlx.
type SQLModuleEventLogStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLModuleEventLogStore) tableName() string {
	if s.DBType == "sqlite" {
		return "ModuleEventLogs"
	}
	return `"ModuleEventLogs"`
}

func (s *SQLModuleEventLogStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLModuleEventLogStore) Create(ctx context.Context, l *ModuleEventLog) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, actor, description, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		s.tableName(),
		s.col("user_id"),
		s.col("client_id"),
		s.col("world_title"),
		s.col("event_type"),
		s.col("created_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			l.UserID, l.ClientID, l.WorldTitle, l.EventType, l.Actor, l.Description, now,
		).Scan(&l.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		l.UserID, l.ClientID, l.WorldTitle, l.EventType, l.Actor, l.Description, now)
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

func (s *SQLModuleEventLogStore) FindFiltered(ctx context.Context, filters ActivityFilters) ([]*ModuleEventLog, int64, error) {
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
	if filters.Action != "" {
		addAnd(fmt.Sprintf("%s = $%d", s.col("event_type"), argN), filters.Action)
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

	var logs []*ModuleEventLog
	dataQ := fmt.Sprintf("SELECT * FROM %s%s ORDER BY id DESC LIMIT $%d",
		s.tableName(), where, argN)
	args = append(args, fetchN)
	if err := s.DB.SelectContext(ctx, &logs, dataQ, args...); err != nil {
		return nil, 0, err
	}
	if logs == nil {
		logs = []*ModuleEventLog{}
	}
	return logs, total, nil
}

func (s *SQLModuleEventLogStore) CleanupOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	now := NewSQLiteTime(cutoff)
	result, err := s.DB.ExecContext(ctx,
		fmt.Sprintf(`DELETE FROM %s WHERE %s < $1`, s.tableName(), s.col("created_at")),
		now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
