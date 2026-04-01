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
	sheetPendingTimeout  = 30 * time.Second
	sheetInactiveTimeout = 5 * time.Minute
	sheetCleanupInterval = 15 * time.Second
)

// SheetSession represents an active sheet streaming session between a consumer WS and a Foundry client.
type SheetSession struct {
	SessionID    string
	ClientID     string
	APIKey       string
	ConsumerConn *websocket.Conn
	State        string // "pending", "active", "closed"
	CreatedAt    time.Time
	LastActivity time.Time
	Metadata     SheetSessionMetadata
}

// SheetSessionMetadata holds optional parameters for the sheet session.
type SheetSessionMetadata struct {
	UUID    string
	Quality float64
	Scale   float64
}

// SheetSessionManager manages concurrent sheet viewing sessions.
type SheetSessionManager struct {
	mu        sync.RWMutex
	sessions  map[string]*SheetSession // sessionID -> session
	maxPerKey int
}

// NewSheetSessionManager creates a new manager.
func NewSheetSessionManager(maxPerKey int) *SheetSessionManager {
	if maxPerKey <= 0 {
		maxPerKey = 3
	}
	return &SheetSessionManager{
		sessions:  make(map[string]*SheetSession),
		maxPerKey: maxPerKey,
	}
}

// CreateSession creates a new sheet session. Returns error string if limit reached.
func (m *SheetSessionManager) CreateSession(clientID, apiKey string, consumerConn *websocket.Conn, metadata SheetSessionMetadata) (string, error) {
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
		return "", fmt.Errorf("maximum sheet sessions (%d) reached for this API key", m.maxPerKey)
	}

	sessionID := fmt.Sprintf("ss_%d_%s", time.Now().UnixMilli(), randomStr(6))
	now := time.Now()

	m.sessions[sessionID] = &SheetSession{
		SessionID:    sessionID,
		ClientID:     clientID,
		APIKey:       apiKey,
		ConsumerConn: consumerConn,
		State:        "pending",
		CreatedAt:    now,
		LastActivity: now,
		Metadata:     metadata,
	}

	return sessionID, nil
}

// GetSession returns a session by ID.
func (m *SheetSessionManager) GetSession(sessionID string) *SheetSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[sessionID]
}

// ActivateSession transitions a session from pending to active.
func (m *SheetSessionManager) ActivateSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok && s.State == "pending" {
		s.State = "active"
		s.LastActivity = time.Now()
	}
}

// UpdateActivity refreshes the session's last activity timestamp.
func (m *SheetSessionManager) UpdateActivity(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		s.LastActivity = time.Now()
	}
}

// EndSession marks a session as closed and removes it.
func (m *SheetSessionManager) EndSession(sessionID string) bool {
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
func (m *SheetSessionManager) TerminateSessionsForClient(clientID string) {
	m.mu.Lock()
	var toNotify []*SheetSession
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
			"type":      "sheet-session-ended",
			"sessionId": s.SessionID,
			"reason":    "foundry-disconnected",
		})
		s.ConsumerConn.WriteMessage(websocket.TextMessage, msg)
	}
}

// TerminateSessionsForConsumer ends all sessions for a consumer WS (when it disconnects).
// Returns session IDs that should be forwarded as end messages to Foundry.
func (m *SheetSessionManager) TerminateSessionsForConsumer(consumerConn *websocket.Conn) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var sessionIDs []string
	for id, s := range m.sessions {
		if s.ConsumerConn == consumerConn && s.State != "closed" {
			s.State = "closed"
			sessionIDs = append(sessionIDs, id)
			delete(m.sessions, id)
		}
	}
	return sessionIDs
}

// StartCleanupLoop starts a goroutine that cleans up timed-out sessions.
func (m *SheetSessionManager) StartCleanupLoop(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(sheetCleanupInterval)
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

func (m *SheetSessionManager) cleanup() {
	m.mu.Lock()
	now := time.Now()
	var timedOut []*SheetSession

	for id, s := range m.sessions {
		if s.State == "closed" {
			delete(m.sessions, id)
			continue
		}
		if s.State == "pending" && now.Sub(s.CreatedAt) > sheetPendingTimeout {
			s.State = "closed"
			timedOut = append(timedOut, s)
			delete(m.sessions, id)
			continue
		}
		if s.State == "active" && now.Sub(s.LastActivity) > sheetInactiveTimeout {
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
			"type":      "sheet-session-ended",
			"sessionId": s.SessionID,
			"reason":    reason,
		})
		s.ConsumerConn.WriteMessage(websocket.TextMessage, msg)
		log.Info().Str("sessionId", s.SessionID).Str("reason", reason).Msg("Sheet session timed out")
	}
}

// unused import guards
var _ = json.Marshal
var _ = rand.Int
