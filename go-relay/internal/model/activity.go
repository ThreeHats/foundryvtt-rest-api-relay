package model

import "time"

// ActivityFilters are shared filter parameters for merged activity log queries.
type ActivityFilters struct {
	Type    string    // "connection" | "remote_request" | "module_event" | "" = all
	World   string    // clientId / sourceClientId / targetClientId
	Action  string    // remote_request only
	Success *bool     // remote_request only; nil = all
	Since   time.Time // zero = no lower bound
	Until   time.Time // zero = no upper bound
	UserID  int64     // 0 = all users (admin only)
	Limit   int       // 0 → caller sets default
	Offset  int
}

// ActivityEvent is the unified event type returned by the activity log endpoint.
// It flattens ConnectionLog, RemoteRequestLog, and ModuleEventLog into a single
// shape for the frontend to render in a chronological feed.
type ActivityEvent struct {
	ID             int64     `json:"id"`
	Type           string    `json:"type"`                     // "connection" | "remote_request" | "module_event"
	EventSubtype   string    `json:"eventSubtype,omitempty"`   // "connect"/"disconnect" (connection only)
	ClientID       string    `json:"clientId"`
	WorldTitle     string    `json:"worldTitle,omitempty"`
	TargetClientID string    `json:"targetClientId,omitempty"` // remote_request only
	Action         string    `json:"action,omitempty"`         // remote_request / module_event
	Success        *bool     `json:"success,omitempty"`        // remote_request only
	ErrorMessage   string    `json:"errorMessage,omitempty"`
	TokenName      string    `json:"tokenName,omitempty"`
	IPAddress      string    `json:"ipAddress,omitempty"`
	Flagged        bool      `json:"flagged,omitempty"`
	Actor          string    `json:"actor,omitempty"`        // module_event only
	Description    string    `json:"description,omitempty"`  // module_event only
	UserID         int64     `json:"userId,omitempty"`       // admin view only
	CreatedAt      time.Time `json:"createdAt"`
}
