package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// KnownClient represents a Foundry VTT client that has connected to the relay.
type KnownClient struct {
	ID             int64          `db:"id" json:"id"`
	UserID         int64          `db:"userId" json:"userId"`
	ClientID       string         `db:"clientId" json:"clientId"`
	WorldID        sql.NullString `db:"worldId" json:"worldId"`
	WorldTitle     sql.NullString `db:"worldTitle" json:"worldTitle"`
	SystemID       sql.NullString `db:"systemId" json:"systemId"`
	SystemTitle    sql.NullString `db:"systemTitle" json:"systemTitle"`
	SystemVersion  sql.NullString `db:"systemVersion" json:"systemVersion"`
	FoundryVersion sql.NullString `db:"foundryVersion" json:"foundryVersion"`
	CustomName     sql.NullString `db:"customName" json:"customName"`
	LastSeenAt     *SQLiteTime    `db:"lastSeenAt" json:"lastSeenAt"`
	// IsOnline uses LooseBool so corrupted SQLite data (timestamps, strings)
	// from earlier param-binding bugs doesn't break the entire endpoint.
	IsOnline LooseBool `db:"isOnline" json:"isOnline"`
	// AutoStartOnRemoteRequest, when true, lets the relay spawn a headless
	// session for this clientId in response to an incoming remote-request from
	// a sibling client (when this client is currently offline). Default false:
	// resource consumption is opt-in per client.
	AutoStartOnRemoteRequest LooseBool `db:"autoStartOnRemoteRequest" json:"autoStartOnRemoteRequest"`
	// CredentialID is the optional explicit link to a stored Credential. When
	// set, the headless auto-start flow uses this credential to log in. When
	// NULL, the auto-start falls back to the user's first credential — works
	// for the common single-Foundry-server deployment.
	CredentialID sql.NullInt64 `db:"credentialId" json:"credentialId"`
	// ServerFingerprint is a random stable ID generated once per world on first
	// pair and stored as a world-scoped Foundry setting. It lets the relay
	// distinguish "same server re-pairing" (reuse clientId) from "different
	// server with the same worldId slug" (mint a new clientId).
	ServerFingerprint sql.NullString `db:"serverFingerprint" json:"-"`
	// PublicUrl is the browser Origin header captured at WebSocket connection time.
	// Populated when the Foundry module connects to /relay.
	PublicUrl sql.NullString `db:"publicUrl" json:"publicUrl"`
	// AllowedTargetClients lists the clientIds this world may invoke
	// remote-request operations against. CSV. Empty = no cross-world access.
	AllowedTargetClients sql.NullString `db:"allowedTargetClients" json:"allowedTargetClients"`
	// RemoteScopes lists the scope strings this world holds for remote-request
	// operations on allowed targets. CSV. Empty = no cross-world access.
	RemoteScopes sql.NullString `db:"remoteScopes" json:"remoteScopes"`
	// RemoteRequestsPerHour is the per-world rate limit for cross-world
	// remote-request operations. 0 = unlimited.
	RemoteRequestsPerHour int        `db:"remoteRequestsPerHour" json:"remoteRequestsPerHour"`
	CreatedAt             SQLiteTime `db:"createdAt" json:"createdAt"`
	UpdatedAt             SQLiteTime `db:"updatedAt" json:"updatedAt"`
}

// GetAllowedTargets returns the parsed list of clientIds this world can target.
func (c *KnownClient) GetAllowedTargets() []string {
	if !c.AllowedTargetClients.Valid || c.AllowedTargetClients.String == "" {
		return nil
	}
	out := []string{}
	for _, s := range strings.Split(c.AllowedTargetClients.String, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// CanTarget returns true if the given clientId is in this world's allow-list.
// The special value "*" grants access to all target clients.
func (c *KnownClient) CanTarget(clientID string) bool {
	if clientID == "" {
		return false
	}
	for _, allowed := range c.GetAllowedTargets() {
		if allowed == "*" || allowed == clientID {
			return true
		}
	}
	return false
}

// GetRemoteScopes returns the parsed list of scopes this world holds for remote-request.
func (c *KnownClient) GetRemoteScopes() []string {
	if !c.RemoteScopes.Valid || c.RemoteScopes.String == "" {
		return nil
	}
	out := []string{}
	for _, s := range strings.Split(c.RemoteScopes.String, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// HasRemoteScope returns true if the given scope is granted to this world.
func (c *KnownClient) HasRemoteScope(scope string) bool {
	if scope == "" {
		return false
	}
	for _, granted := range c.GetRemoteScopes() {
		if granted == scope {
			return true
		}
	}
	return false
}

// KnownClientStore defines operations on known Foundry clients.
type KnownClientStore interface {
	Upsert(ctx context.Context, client *KnownClient) error
	FindByID(ctx context.Context, id int64) (*KnownClient, error)
	FindAllByUser(ctx context.Context, userID int64) ([]*KnownClient, error)
	FindByClientID(ctx context.Context, userID int64, clientID string) (*KnownClient, error)
	// FindAnyByClientID looks up a known client by clientId WITHOUT a user
	// filter. Used by the public active-connection probe endpoint, which
	// needs to verify the clientId exists to return a meaningful answer
	// without revealing which user owns it.
	FindAnyByClientID(ctx context.Context, clientID string) (*KnownClient, error)
	// FindByWorldID looks up a known client by (userId, worldId). Used at
	// pairing time to reuse the existing clientId for a world that has been
	// paired before.
	FindByWorldID(ctx context.Context, userID int64, worldID string) (*KnownClient, error)
	// FindByWorldIDAndPublicUrl looks up a known client by (userId, worldId, publicUrl).
	// More precise than FindByWorldID: differentiates two Foundry instances that share
	// a worldId slug but run on different servers (e.g. localhost:30000 vs :30001).
	FindByWorldIDAndPublicUrl(ctx context.Context, userID int64, worldID, publicUrl string) (*KnownClient, error)
	// FindByServerFingerprint looks up a known client by (userId, serverFingerprint).
	// Used at pairing time to reuse the existing clientId when the same physical
	// Foundry server re-pairs (even if its worldId slug is shared by another server).
	FindByServerFingerprint(ctx context.Context, userID int64, fingerprint string) (*KnownClient, error)
	// FindByWorldIDCrossUser finds any known client with this worldId belonging
	// to a different user. Used to detect multi-account abuse.
	FindByWorldIDCrossUser(ctx context.Context, worldID string, excludeUserID int64) (*KnownClient, error)
	SetOnline(ctx context.Context, userID int64, clientID string) error
	SetOffline(ctx context.Context, userID int64, clientID string) error
	// ResetAllOnline marks every known client as offline. Called at startup
	// on single-instance deployments to clear stale flags left by a crash or
	// unclean shutdown. Not safe to call on a running multi-instance cluster.
	ResetAllOnline(ctx context.Context) error
	SetAutoStartOnRemoteRequest(ctx context.Context, id int64, enabled bool) error
	SetCredentialID(ctx context.Context, id int64, credentialID sql.NullInt64) error
	SetCrossWorldSettings(ctx context.Context, id int64, allowedTargetClients, remoteScopes sql.NullString, remoteRequestsPerHour int) error
	Delete(ctx context.Context, id int64) error
}

// SQLKnownClientStore implements KnownClientStore with sqlx.
type SQLKnownClientStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLKnownClientStore) tableName() string {
	if s.DBType == "sqlite" {
		return "KnownClients"
	}
	return `"KnownClients"`
}

func (s *SQLKnownClientStore) col(name string) string {
	return Col(s.DBType, name)
}

func (s *SQLKnownClientStore) Upsert(ctx context.Context, client *KnownClient) error {
	now := NewSQLiteTime(time.Now())

	// Insert column order MUST stay in lockstep with the args slice below.
	insertCols := []string{
		"user_id", "client_id", "world_id", "world_title", "system_id", "system_title",
		"system_version", "foundry_version", "custom_name", "server_fingerprint",
		"is_online", "created_at", "updated_at", "last_seen_at", "public_url",
	}
	args := []interface{}{
		client.UserID, client.ClientID, client.WorldID, client.WorldTitle, client.SystemID, client.SystemTitle,
		client.SystemVersion, client.FoundryVersion, client.CustomName, client.ServerFingerprint,
		bool(client.IsOnline), now, now, now, client.PublicUrl,
	}

	colList := make([]string, len(insertCols))
	placeholders := make([]string, len(insertCols))
	for i, c := range insertCols {
		colList[i] = s.col(c)
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	keyword := "EXCLUDED"
	if s.DBType == "sqlite" {
		keyword = "excluded"
	}

	// Preserve the stored value when the incoming column is empty: the WS-connect
	// upsert carries no serverFingerprint and "add browser" carries no world/system
	// metadata, so an unconditional overwrite would null them out. Non-empty updates.
	preserveOnEmpty := map[string]bool{
		"world_id": true, "world_title": true, "system_id": true, "system_title": true,
		"system_version": true, "foundry_version": true, "custom_name": true,
		"server_fingerprint": true, "public_url": true,
	}
	setClauses := make([]string, 0, len(insertCols))
	for _, c := range insertCols {
		// Conflict keys and the immutable created_at are never updated.
		if c == "user_id" || c == "client_id" || c == "created_at" {
			continue
		}
		col := s.col(c)
		if preserveOnEmpty[c] {
			// The existing-row reference must be table-qualified: Postgres rejects
			// bare column names in DO UPDATE as ambiguous vs EXCLUDED (42702).
			setClauses = append(setClauses, fmt.Sprintf("%s=COALESCE(NULLIF(%s.%s, ''), %s.%s)", col, keyword, col, s.tableName(), col))
		} else {
			setClauses = append(setClauses, fmt.Sprintf("%s=%s.%s", col, keyword, col))
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s, %s) DO UPDATE SET %s",
		s.tableName(), strings.Join(colList, ", "), strings.Join(placeholders, ", "),
		s.col("user_id"), s.col("client_id"), strings.Join(setClauses, ", "))

	if s.DBType == "sqlite" {
		result, err := s.DB.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		if id > 0 {
			client.ID = id
		}
		client.CreatedAt = now
		client.UpdatedAt = now
		return nil
	}

	// PostgreSQL
	return s.DB.QueryRowContext(ctx, query+" RETURNING id", args...).Scan(&client.ID)
}

func (s *SQLKnownClientStore) FindAllByUser(ctx context.Context, userID int64) ([]*KnownClient, error) {
	var clients []*KnownClient
	err := s.DB.SelectContext(ctx, &clients, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("user_id")), userID)
	if clients == nil {
		clients = []*KnownClient{}
	}
	return clients, err
}

func (s *SQLKnownClientStore) FindByID(ctx context.Context, id int64) (*KnownClient, error) {
	var c KnownClient
	err := s.DB.GetContext(ctx, &c, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", s.tableName()), id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (s *SQLKnownClientStore) FindByClientID(ctx context.Context, userID int64, clientID string) (*KnownClient, error) {
	var c KnownClient
	err := s.DB.GetContext(ctx, &c, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s = $2",
		s.tableName(), s.col("user_id"), s.col("client_id")), userID, clientID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (s *SQLKnownClientStore) FindAnyByClientID(ctx context.Context, clientID string) (*KnownClient, error) {
	var c KnownClient
	err := s.DB.GetContext(ctx, &c, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 LIMIT 1",
		s.tableName(), s.col("client_id")), clientID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (s *SQLKnownClientStore) FindByWorldID(ctx context.Context, userID int64, worldID string) (*KnownClient, error) {
	var c KnownClient
	err := s.DB.GetContext(ctx, &c, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s = $2 LIMIT 1",
		s.tableName(), s.col("user_id"), s.col("world_id")), userID, worldID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (s *SQLKnownClientStore) FindByWorldIDAndPublicUrl(ctx context.Context, userID int64, worldID, publicUrl string) (*KnownClient, error) {
	var c KnownClient
	err := s.DB.GetContext(ctx, &c, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s = $2 AND %s = $3 LIMIT 1",
		s.tableName(), s.col("user_id"), s.col("world_id"), s.col("public_url")), userID, worldID, publicUrl)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (s *SQLKnownClientStore) FindByServerFingerprint(ctx context.Context, userID int64, fingerprint string) (*KnownClient, error) {
	var c KnownClient
	err := s.DB.GetContext(ctx, &c, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s = $2 LIMIT 1",
		s.tableName(), s.col("user_id"), s.col("server_fingerprint")), userID, fingerprint)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (s *SQLKnownClientStore) FindByWorldIDCrossUser(ctx context.Context, worldID string, excludeUserID int64) (*KnownClient, error) {
	var c KnownClient
	err := s.DB.GetContext(ctx, &c, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s != $2 LIMIT 1",
		s.tableName(), s.col("world_id"), s.col("user_id")), worldID, excludeUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (s *SQLKnownClientStore) SetOnline(ctx context.Context, userID int64, clientID string) error {
	// Use SQLiteTime so the value is formatted via SQLiteTime.Value() (UTC ISO).
	// Passing raw time.Time would use the database driver's default formatter,
	// which may include timezone info that doesn't match SQLiteTime.Scan()'s
	// expectations and breaks roundtrips.
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET %s = 1, %s = $1, %s = $2 WHERE %s = $3 AND %s = $4`,
		s.tableName(), s.col("is_online"), s.col("last_seen_at"), s.col("updated_at"),
		s.col("user_id"), s.col("client_id"))
	_, err := s.DB.ExecContext(ctx, query, now, now, userID, clientID)
	return err
}

func (s *SQLKnownClientStore) SetOffline(ctx context.Context, userID int64, clientID string) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET %s = 0, %s = $1, %s = $2 WHERE %s = $3 AND %s = $4`,
		s.tableName(), s.col("is_online"), s.col("last_seen_at"), s.col("updated_at"),
		s.col("user_id"), s.col("client_id"))
	_, err := s.DB.ExecContext(ctx, query, now, now, userID, clientID)
	return err
}

func (s *SQLKnownClientStore) ResetAllOnline(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf(`UPDATE %s SET %s = $1`, s.tableName(), s.col("is_online")), false)
	return err
}

func (s *SQLKnownClientStore) SetCredentialID(ctx context.Context, id int64, credentialID sql.NullInt64) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET %s = $1, %s = $2 WHERE id = $3`,
		s.tableName(), s.col("credential_id"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, credentialID, now, id)
	return err
}

func (s *SQLKnownClientStore) SetAutoStartOnRemoteRequest(ctx context.Context, id int64, enabled bool) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET %s = $1, %s = $2 WHERE id = $3`,
		s.tableName(), s.col("auto_start_on_remote_request"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, enabled, now, id)
	return err
}

func (s *SQLKnownClientStore) SetCrossWorldSettings(ctx context.Context, id int64, allowedTargetClients, remoteScopes sql.NullString, remoteRequestsPerHour int) error {
	now := NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET %s = $1, %s = $2, %s = $3, %s = $4 WHERE id = $5`,
		s.tableName(),
		s.col("allowed_target_clients"), s.col("remote_scopes"),
		s.col("remote_requests_per_hour"), s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query, allowedTargetClients, remoteScopes, remoteRequestsPerHour, now, id)
	return err
}

func (s *SQLKnownClientStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName()), id)
	return err
}
