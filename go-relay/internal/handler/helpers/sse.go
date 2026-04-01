package helpers

import (
	"fmt"
	"net/http"
	"sync"
)

const MaxSSEConnectionsPerClient = 10

// SSEConnection represents a Server-Sent Events connection for chat events.
type SSEConnection struct {
	W        http.ResponseWriter
	Flusher  http.Flusher
	ClientID string
	Filters  SSEFilters
	Done     <-chan struct{} // From request context
}

// SSEFilters holds filter criteria for chat event SSE streams.
type SSEFilters struct {
	Speaker    string
	Type       *int
	WhisperOnly bool
	UserID     string
}

// RollSSEConnection represents a Server-Sent Events connection for roll events.
type RollSSEConnection struct {
	W        http.ResponseWriter
	Flusher  http.Flusher
	ClientID string
	Filters  RollSSEFilters
	Done     <-chan struct{}
}

// RollSSEFilters holds filter criteria for roll event SSE streams.
type RollSSEFilters struct {
	UserID string
}

// WSEventConnection represents a WebSocket connection subscribing to events.
type WSEventConnection struct {
	SendFunc func(data interface{}) bool
	ClientID string
	Channel  string // "chat-events" or "roll-events"
	Filters  SSEFilters
}

// SSEManager manages all SSE and WS event connections.
type SSEManager struct {
	mu              sync.RWMutex
	chatSSE         map[string][]*SSEConnection     // clientID -> connections
	rollSSE         map[string][]*RollSSEConnection  // clientID -> connections
	wsEvents        map[string][]*WSEventConnection  // clientID -> connections
}

// NewSSEManager creates a new SSE connection manager.
func NewSSEManager() *SSEManager {
	return &SSEManager{
		chatSSE:  make(map[string][]*SSEConnection),
		rollSSE:  make(map[string][]*RollSSEConnection),
		wsEvents: make(map[string][]*WSEventConnection),
	}
}

// AddChatSSE adds a chat SSE connection. Returns false if max reached.
func (m *SSEManager) AddChatSSE(conn *SSEConnection) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.chatSSE[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		return false
	}
	m.chatSSE[conn.ClientID] = append(m.chatSSE[conn.ClientID], conn)
	return true
}

// RemoveChatSSE removes a chat SSE connection.
func (m *SSEManager) RemoveChatSSE(conn *SSEConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	conns := m.chatSSE[conn.ClientID]
	for i, c := range conns {
		if c == conn {
			m.chatSSE[conn.ClientID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(m.chatSSE[conn.ClientID]) == 0 {
		delete(m.chatSSE, conn.ClientID)
	}
}

// GetChatSSE returns all chat SSE connections for a client.
func (m *SSEManager) GetChatSSE(clientID string) []*SSEConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.chatSSE[clientID]
}

// AddRollSSE adds a roll SSE connection. Returns false if max reached.
func (m *SSEManager) AddRollSSE(conn *RollSSEConnection) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.rollSSE[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		return false
	}
	m.rollSSE[conn.ClientID] = append(m.rollSSE[conn.ClientID], conn)
	return true
}

// RemoveRollSSE removes a roll SSE connection.
func (m *SSEManager) RemoveRollSSE(conn *RollSSEConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	conns := m.rollSSE[conn.ClientID]
	for i, c := range conns {
		if c == conn {
			m.rollSSE[conn.ClientID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(m.rollSSE[conn.ClientID]) == 0 {
		delete(m.rollSSE, conn.ClientID)
	}
}

// GetRollSSE returns all roll SSE connections for a client.
func (m *SSEManager) GetRollSSE(clientID string) []*RollSSEConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rollSSE[clientID]
}

// AddWSEvent adds a WS event connection. Returns false if max reached.
func (m *SSEManager) AddWSEvent(conn *WSEventConnection) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.wsEvents[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		return false
	}
	m.wsEvents[conn.ClientID] = append(m.wsEvents[conn.ClientID], conn)
	return true
}

// RemoveWSEvent removes a WS event connection.
func (m *SSEManager) RemoveWSEvent(conn *WSEventConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	conns := m.wsEvents[conn.ClientID]
	for i, c := range conns {
		if c == conn {
			m.wsEvents[conn.ClientID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(m.wsEvents[conn.ClientID]) == 0 {
		delete(m.wsEvents, conn.ClientID)
	}
}

// GetWSEvents returns all WS event connections for a client.
func (m *SSEManager) GetWSEvents(clientID string) []*WSEventConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.wsEvents[clientID]
}

// AddWSEventFunc creates a WSEventConnection from a send function and returns a remove callback.
// This satisfies the ws.SSEManagerInterface used by the client API WebSocket.
func (m *SSEManager) AddWSEventFunc(clientID, channel string, sendFunc func(data interface{}) bool) (remove func(), ok bool) {
	conn := &WSEventConnection{
		SendFunc: sendFunc,
		ClientID: clientID,
		Channel:  channel,
	}
	if !m.AddWSEvent(conn) {
		return nil, false
	}
	return func() { m.RemoveWSEvent(conn) }, true
}

// SendSSEEvent sends an SSE event to a connection.
func SendSSEEvent(conn *SSEConnection, eventType string, data interface{}) {
	fmt.Fprintf(conn.W, "event: %s\ndata: %v\n\n", eventType, data)
	conn.Flusher.Flush()
}

// SendRollSSEEvent sends an SSE event to a roll connection.
func SendRollSSEEvent(conn *RollSSEConnection, data interface{}) {
	fmt.Fprintf(conn.W, "event: roll\ndata: %v\n\n", data)
	conn.Flusher.Flush()
}
