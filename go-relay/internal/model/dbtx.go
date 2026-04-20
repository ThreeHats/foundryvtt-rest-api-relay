package model

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DBTX is the common interface between *sqlx.DB and *sqlx.Tx.
// Both types satisfy this interface, allowing stores to operate
// within either a regular connection or a transaction.
type DBTX interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
