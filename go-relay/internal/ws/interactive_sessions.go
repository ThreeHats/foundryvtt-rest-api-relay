package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	interactivePendingTimeout  = 30 * time.Second
	interactiveInactiveTimeout = 5 * time.Minute
	interactiveCleanupInterval = 15 * time.Second
)

// InteractiveSession represents an active interactive streaming session between a consumer WS and a Foundry client.
type InteractiveSession struct {
	SessionID    string
	ClientID     string
	APIKey       string
	ConsumerConn *websocket.Conn
	ConsumerID   string // stable string identity for the consumer connection; used by TerminateSessionsForConsumer
	State        string // "pending", "active", "closed"
	CreatedAt    time.Time
	LastActivity time.Time
	Metadata     InteractiveSessionMetadata
}

// InteractiveSessionMetadata holds optional parameters for the interactive session.
type InteractiveSessionMetadata struct {
	UUID    string
	Quality float64
	Scale   float64
}

// InteractiveSessionManager manages concurrent interactive viewing sessions.
type InteractiveSessionManager struct {
	mu        sync.RWMutex
	sessions  map[string]*InteractiveSession // sessionID -> session
	maxPerKey int
}

// NewInteractiveSessionManager creates a new manager.
func NewInteractiveSessionManager(maxPerKey int) *InteractiveSessionManager {
	if maxPerKey <= 0 {
		maxPerKey = 3
	}
	return &InteractiveSessionManager{
		sessions:  make(map[string]*InteractiveSession),
		maxPerKey: maxPerKey,
	}
}

// CreateSession creates a new interactive session. Returns error string if limit reached.
// consumerID is a stable string identity for the consumer connection (e.g. randomStr(8))
// so that TerminateSessionsForConsumer can match by ID rather than pointer.
func (m *InteractiveSessionManager) CreateSession(clientID, apiKey, consumerID string, consumerConn *websocket.Conn, metadata InteractiveSessionMetadata) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Count sessions for this API key
	count := 0
	for _, s := range m.sessions {
		if s.APIKey == apiKey && s.State != "closed" {
			count++
		}
	}
	if count >= m.maxPerKey {
		return "", fmt.Errorf("maximum interactive sessions (%d) reached for this API key", m.maxPerKey)
	}

	sessionID := fmt.Sprintf("is_%d_%s", time.Now().UnixMilli(), randomStr(6))
	now := time.Now()

	m.sessions[sessionID] = &InteractiveSession{
		SessionID:    sessionID,
		ClientID:     clientID,
		APIKey:       apiKey,
		ConsumerConn: consumerConn,
		ConsumerID:   consumerID,
		State:        "pending",
		CreatedAt:    now,
		LastActivity: now,
		Metadata:     metadata,
	}

	return sessionID, nil
}

// GetSession returns a session by ID.
func (m *InteractiveSessionManager) GetSession(sessionID string) *InteractiveSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[sessionID]
}

// ActivateSession transitions a session from pending to active.
func (m *InteractiveSessionManager) ActivateSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok && s.State == "pending" {
		s.State = "active"
		s.LastActivity = time.Now()
	}
}

// UpdateActivity refreshes the session's last activity timestamp.
func (m *InteractiveSessionManager) UpdateActivity(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		s.LastActivity = time.Now()
	}
}

// EndSession marks a session as closed and removes it.
func (m *InteractiveSessionManager) EndSession(sessionID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		s.State = "closed"
		delete(m.sessions, sessionID)
		return true
	}
	return false
}

// TerminateSessionsForClient ends all sessions for a Foundry client (when it disconnects).
func (m *InteractiveSessionManager) TerminateSessionsForClient(clientID string) {
	m.mu.Lock()
	var toNotify []*InteractiveSession
	for id, s := range m.sessions {
		if s.ClientID == clientID && s.State != "closed" {
			s.State = "closed"
			toNotify = append(toNotify, s)
			delete(m.sessions, id)
		}
	}
	m.mu.Unlock()

	// Notify consumers
	for _, s := range toNotify {
		msg, _ := json.Marshal(map[string]interface{}{
			"type":      "interactive-session-ended",
			"sessionId": s.SessionID,
			"reason":    "foundry-disconnected",
		})
		s.ConsumerConn.WriteMessage(websocket.TextMessage, msg)
	}
}

// TerminateSessionsForConsumer ends all sessions for a consumer WS (when it disconnects).
// consumerID is the stable string identity assigned at CreateSession time.
// Returns session IDs that should be forwarded as end messages to Foundry.
func (m *InteractiveSessionManager) TerminateSessionsForConsumer(consumerID string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var sessionIDs []string
	for id, s := range m.sessions {
		if s.ConsumerID == consumerID && s.State != "closed" {
			s.State = "closed"
			sessionIDs = append(sessionIDs, id)
			delete(m.sessions, id)
		}
	}
	return sessionIDs
}

// StartCleanupLoop starts a goroutine that cleans up timed-out sessions.
func (m *InteractiveSessionManager) StartCleanupLoop(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(interactiveCleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.cleanup()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m *InteractiveSessionManager) cleanup() {
	m.mu.Lock()
	now := time.Now()
	var timedOut []*InteractiveSession

	for id, s := range m.sessions {
		if s.State == "closed" {
			delete(m.sessions, id)
			continue
		}
		if s.State == "pending" && now.Sub(s.CreatedAt) > interactivePendingTimeout {
			s.State = "closed"
			timedOut = append(timedOut, s)
			delete(m.sessions, id)
			continue
		}
		if s.State == "active" && now.Sub(s.LastActivity) > interactiveInactiveTimeout {
			s.State = "closed"
			timedOut = append(timedOut, s)
			delete(m.sessions, id)
		}
	}
	m.mu.Unlock()

	// Notify timed-out consumers
	for _, s := range timedOut {
		reason := "inactivity-timeout"
		if s.State == "pending" {
			reason = "pending-timeout"
		}
		msg, _ := json.Marshal(map[string]interface{}{
			"type":      "interactive-session-ended",
			"sessionId": s.SessionID,
			"reason":    reason,
		})
		s.ConsumerConn.WriteMessage(websocket.TextMessage, msg)
		log.Info().Str("sessionId", s.SessionID).Str("reason", reason).Msg("Interactive session timed out")
	}
}

// unused import guards
var _ = json.Marshal
var _ = rand.Int
