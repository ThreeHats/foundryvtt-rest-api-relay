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
	Channel  string // "chat-events", "roll-events", "combat-events", "actor-events", "scene-events", "hooks"
	Filters  SSEFilters
	// HookFilters limits which hooks are forwarded (empty = all hooks)
	HookFilters []string
	// ActorUUID filters actor events to a specific actor
	ActorUUID string
	// SceneID filters scene events to a specific scene
	SceneID string
}

// GenericSSEConnection represents an SSE connection for the hooks firehose.
type GenericSSEConnection struct {
	W           http.ResponseWriter
	Flusher     http.Flusher
	ClientID    string
	HookFilters []string // empty = all hooks
	Done        <-chan struct{}
}

// CombatSSEConnection represents an SSE connection for combat events.
type CombatSSEConnection struct {
	W           http.ResponseWriter
	Flusher     http.Flusher
	ClientID    string
	EncounterID string // optional filter
	Done        <-chan struct{}
}

// ActorSSEConnection represents an SSE connection for actor events.
type ActorSSEConnection struct {
	W         http.ResponseWriter
	Flusher   http.Flusher
	ClientID  string
	ActorUUID string // required filter
	Done      <-chan struct{}
}

// SceneSSEConnection represents an SSE connection for scene events.
type SceneSSEConnection struct {
	W        http.ResponseWriter
	Flusher  http.Flusher
	ClientID string
	SceneID  string // optional — defaults to active scene events
	Done     <-chan struct{}
}

// SSEManager manages all SSE and WS event connections.
type SSEManager struct {
	mu              sync.RWMutex
	chatSSE         map[string][]*SSEConnection        // clientID -> connections
	rollSSE         map[string][]*RollSSEConnection     // clientID -> connections
	wsEvents        map[string][]*WSEventConnection     // clientID -> connections
	genericSSE      map[string][]*GenericSSEConnection  // clientID -> connections
	combatSSE       map[string][]*CombatSSEConnection   // clientID -> connections
	actorSSE        map[string][]*ActorSSEConnection    // clientID -> connections
	sceneSSE        map[string][]*SceneSSEConnection    // clientID -> connections

	// OnSubscriberCountChanged is called (outside the lock) whenever the total
	// subscriber count for a channel+client changes. Used to notify the Foundry
	// module to enable or disable event hooks on demand.
	OnSubscriberCountChanged func(clientID, channel string, count int)
}

// NewSSEManager creates a new SSE connection manager.
func NewSSEManager() *SSEManager {
	return &SSEManager{
		chatSSE:    make(map[string][]*SSEConnection),
		rollSSE:    make(map[string][]*RollSSEConnection),
		wsEvents:   make(map[string][]*WSEventConnection),
		genericSSE: make(map[string][]*GenericSSEConnection),
		combatSSE:  make(map[string][]*CombatSSEConnection),
		actorSSE:   make(map[string][]*ActorSSEConnection),
		sceneSSE:   make(map[string][]*SceneSSEConnection),
	}
}

// totalForChannel returns the total subscriber count for a channel (SSE + WS).
// Must be called with m.mu held.
func (m *SSEManager) totalForChannel(clientID, channel string) int {
	var n int
	switch channel {
	case "chat-events":
		n = len(m.chatSSE[clientID])
	case "roll-events":
		n = len(m.rollSSE[clientID])
	case "hooks":
		n = len(m.genericSSE[clientID])
	case "combat-events":
		n = len(m.combatSSE[clientID])
	case "actor-events":
		n = len(m.actorSSE[clientID])
	case "scene-events":
		n = len(m.sceneSSE[clientID])
	}
	for _, ws := range m.wsEvents[clientID] {
		if ws.Channel == channel {
			n++
		}
	}
	return n
}

// TotalForChannel returns the total subscriber count for a channel (SSE + WS).
// Safe for concurrent use.
func (m *SSEManager) TotalForChannel(clientID, channel string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalForChannel(clientID, channel)
}

// emitCountChanged calls OnSubscriberCountChanged outside the lock.
func (m *SSEManager) emitCountChanged(clientID, channel string, count int) {
	if m.OnSubscriberCountChanged != nil {
		m.OnSubscriberCountChanged(clientID, channel, count)
	}
}

// AddChatSSE adds a chat SSE connection. Returns false if max reached.
func (m *SSEManager) AddChatSSE(conn *SSEConnection) bool {
	m.mu.Lock()
	if len(m.chatSSE[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		m.mu.Unlock()
		return false
	}
	m.chatSSE[conn.ClientID] = append(m.chatSSE[conn.ClientID], conn)
	count := m.totalForChannel(conn.ClientID, "chat-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "chat-events", count)
	return true
}

// RemoveChatSSE removes a chat SSE connection.
func (m *SSEManager) RemoveChatSSE(conn *SSEConnection) {
	m.mu.Lock()
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
	count := m.totalForChannel(conn.ClientID, "chat-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "chat-events", count)
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
	if len(m.rollSSE[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		m.mu.Unlock()
		return false
	}
	m.rollSSE[conn.ClientID] = append(m.rollSSE[conn.ClientID], conn)
	count := m.totalForChannel(conn.ClientID, "roll-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "roll-events", count)
	return true
}

// RemoveRollSSE removes a roll SSE connection.
func (m *SSEManager) RemoveRollSSE(conn *RollSSEConnection) {
	m.mu.Lock()
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
	count := m.totalForChannel(conn.ClientID, "roll-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "roll-events", count)
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
	if len(m.wsEvents[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		m.mu.Unlock()
		return false
	}
	m.wsEvents[conn.ClientID] = append(m.wsEvents[conn.ClientID], conn)
	count := m.totalForChannel(conn.ClientID, conn.Channel)
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, conn.Channel, count)
	return true
}

// RemoveWSEvent removes a WS event connection.
func (m *SSEManager) RemoveWSEvent(conn *WSEventConnection) {
	m.mu.Lock()
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
	count := m.totalForChannel(conn.ClientID, conn.Channel)
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, conn.Channel, count)
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

// --- Generic/Hooks SSE ---

func (m *SSEManager) AddGenericSSE(conn *GenericSSEConnection) bool {
	m.mu.Lock()
	if len(m.genericSSE[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		m.mu.Unlock()
		return false
	}
	m.genericSSE[conn.ClientID] = append(m.genericSSE[conn.ClientID], conn)
	count := m.totalForChannel(conn.ClientID, "hooks")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "hooks", count)
	return true
}

func (m *SSEManager) RemoveGenericSSE(conn *GenericSSEConnection) {
	m.mu.Lock()
	conns := m.genericSSE[conn.ClientID]
	for i, c := range conns {
		if c == conn {
			m.genericSSE[conn.ClientID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(m.genericSSE[conn.ClientID]) == 0 {
		delete(m.genericSSE, conn.ClientID)
	}
	count := m.totalForChannel(conn.ClientID, "hooks")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "hooks", count)
}

func (m *SSEManager) GetGenericSSE(clientID string) []*GenericSSEConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.genericSSE[clientID]
}

// --- Combat SSE ---

func (m *SSEManager) AddCombatSSE(conn *CombatSSEConnection) bool {
	m.mu.Lock()
	if len(m.combatSSE[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		m.mu.Unlock()
		return false
	}
	m.combatSSE[conn.ClientID] = append(m.combatSSE[conn.ClientID], conn)
	count := m.totalForChannel(conn.ClientID, "combat-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "combat-events", count)
	return true
}

func (m *SSEManager) RemoveCombatSSE(conn *CombatSSEConnection) {
	m.mu.Lock()
	conns := m.combatSSE[conn.ClientID]
	for i, c := range conns {
		if c == conn {
			m.combatSSE[conn.ClientID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(m.combatSSE[conn.ClientID]) == 0 {
		delete(m.combatSSE, conn.ClientID)
	}
	count := m.totalForChannel(conn.ClientID, "combat-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "combat-events", count)
}

func (m *SSEManager) GetCombatSSE(clientID string) []*CombatSSEConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.combatSSE[clientID]
}

// --- Actor SSE ---

func (m *SSEManager) AddActorSSE(conn *ActorSSEConnection) bool {
	m.mu.Lock()
	if len(m.actorSSE[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		m.mu.Unlock()
		return false
	}
	m.actorSSE[conn.ClientID] = append(m.actorSSE[conn.ClientID], conn)
	count := m.totalForChannel(conn.ClientID, "actor-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "actor-events", count)
	return true
}

func (m *SSEManager) RemoveActorSSE(conn *ActorSSEConnection) {
	m.mu.Lock()
	conns := m.actorSSE[conn.ClientID]
	for i, c := range conns {
		if c == conn {
			m.actorSSE[conn.ClientID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(m.actorSSE[conn.ClientID]) == 0 {
		delete(m.actorSSE, conn.ClientID)
	}
	count := m.totalForChannel(conn.ClientID, "actor-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "actor-events", count)
}

func (m *SSEManager) GetActorSSE(clientID string) []*ActorSSEConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.actorSSE[clientID]
}

// --- Scene SSE ---

func (m *SSEManager) AddSceneSSE(conn *SceneSSEConnection) bool {
	m.mu.Lock()
	if len(m.sceneSSE[conn.ClientID]) >= MaxSSEConnectionsPerClient {
		m.mu.Unlock()
		return false
	}
	m.sceneSSE[conn.ClientID] = append(m.sceneSSE[conn.ClientID], conn)
	count := m.totalForChannel(conn.ClientID, "scene-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "scene-events", count)
	return true
}

func (m *SSEManager) RemoveSceneSSE(conn *SceneSSEConnection) {
	m.mu.Lock()
	conns := m.sceneSSE[conn.ClientID]
	for i, c := range conns {
		if c == conn {
			m.sceneSSE[conn.ClientID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(m.sceneSSE[conn.ClientID]) == 0 {
		delete(m.sceneSSE, conn.ClientID)
	}
	count := m.totalForChannel(conn.ClientID, "scene-events")
	m.mu.Unlock()
	m.emitCountChanged(conn.ClientID, "scene-events", count)
}

func (m *SSEManager) GetSceneSSE(clientID string) []*SceneSSEConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sceneSSE[clientID]
}
