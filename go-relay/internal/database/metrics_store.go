package database

import (
	"context"
	"database/sql"
	"encoding/json"
)

// SaveMetricsSnapshot upserts the single-row cumulative metrics snapshot.
// Uses INSERT OR REPLACE (SQLite) or INSERT … ON CONFLICT DO UPDATE (Postgres)
// so the table always holds exactly one row.
func (db *DB) SaveMetricsSnapshot(byEndpoint map[string]int, byUser map[int64]int, errorsTotal int) error {
	ctx := context.Background()

	epJSON, err := json.Marshal(byEndpoint)
	if err != nil {
		return err
	}
	usrJSON, err := json.Marshal(byUser)
	if err != nil {
		return err
	}

	if db.dbType == "sqlite" {
		_, err = db.sqlDB.ExecContext(ctx,
			`INSERT OR REPLACE INTO MetricsSnapshots (id, "byEndpoint", "byUser", "errorsTotal", "savedAt")
			 VALUES (1, ?, ?, ?, datetime('now'))`,
			string(epJSON), string(usrJSON), errorsTotal,
		)
	} else {
		_, err = db.sqlDB.ExecContext(ctx,
			`INSERT INTO "MetricsSnapshots" (id, "byEndpoint", "byUser", "errorsTotal", "savedAt")
			 VALUES (1, $1, $2, $3, NOW())
			 ON CONFLICT (id) DO UPDATE SET
			   "byEndpoint" = EXCLUDED."byEndpoint",
			   "byUser"     = EXCLUDED."byUser",
			   "errorsTotal"= EXCLUDED."errorsTotal",
			   "savedAt"    = NOW()`,
			string(epJSON), string(usrJSON), errorsTotal,
		)
	}
	return err
}

// LoadMetricsSnapshot reads the persisted cumulative metrics snapshot.
// Returns empty maps (not an error) when no snapshot row exists yet.
func (db *DB) LoadMetricsSnapshot() (byEndpoint map[string]int, byUser map[int64]int, errorsTotal int, err error) {
	byEndpoint = make(map[string]int)
	byUser = make(map[int64]int)

	ctx := context.Background()
	var epJSON, usrJSON string

	tableName := `MetricsSnapshots`
	if db.dbType != "sqlite" {
		tableName = `"MetricsSnapshots"`
	}

	row := db.sqlDB.QueryRowContext(ctx,
		`SELECT "byEndpoint", "byUser", "errorsTotal" FROM `+tableName+` WHERE id = 1`,
	)
	if scanErr := row.Scan(&epJSON, &usrJSON, &errorsTotal); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return byEndpoint, byUser, 0, nil // first startup, no snapshot yet
		}
		return byEndpoint, byUser, 0, scanErr
	}

	// Ignore unmarshal errors — malformed JSON just results in empty maps.
	_ = json.Unmarshal([]byte(epJSON), &byEndpoint)
	_ = json.Unmarshal([]byte(usrJSON), &byUser)
	return byEndpoint, byUser, errorsTotal, nil
}
