package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// MessageHandler is called when a message of a registered type arrives.
type MessageHandler func(client *Client, message map[string]interface{})

// RawMessageHandler receives both parsed and raw bytes of the message.
type RawMessageHandler func(client *Client, message map[string]interface{}, rawBytes []byte)

// ClientManager manages all connected Foundry clients.
type ClientManager struct {
	mu              sync.RWMutex
	clients         map[string]*Client           // clientID -> Client
	tokenGroups     map[string]map[string]bool    // apiKey -> set of clientIDs
	messageHandlers    map[string][]MessageHandler    // messageType -> handlers
	rawMessageHandlers map[string][]RawMessageHandler // messageType -> raw handlers
	redis           *config.RedisClient
	instanceID      string
	// OnClientRemoved is called when a client disconnects. Set by server after init.
	OnClientRemoved func(clientID string)
}

const clientExpiry = 2 * time.Hour

// NewClientManager creates a new ClientManager.
func NewClientManager(redis *config.RedisClient, instanceID string) *ClientManager {
	return &ClientManager{
		clients:            make(map[string]*Client),
		tokenGroups:        make(map[string]map[string]bool),
		messageHandlers:    make(map[string][]MessageHandler),
		rawMessageHandlers: make(map[string][]RawMessageHandler),
		redis:           redis,
		instanceID:      instanceID,
	}
}

// AddClient registers a new Foundry client connection.
func (m *ClientManager) AddClient(conn *websocket.Conn, id, token, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName string) (*Client, error) {
	m.mu.Lock()

	// Reject duplicate connections
	if _, exists := m.clients[id]; exists {
		m.mu.Unlock()
		log.Warn().Str("clientId", id).Msg("Client already exists, rejecting connection")
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(4004, "Client ID already connected"))
		conn.Close()
		return nil, fmt.Errorf("duplicate client ID: %s", id)
	}

	client := NewClient(conn, id, token, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName)
	m.clients[id] = client

	// Add to token group
	if m.tokenGroups[token] == nil {
		m.tokenGroups[token] = make(map[string]bool)
	}
	m.tokenGroups[token][id] = true

	m.mu.Unlock()

	// Store in Redis (non-blocking, best effort)
	if m.redis != nil && m.redis.IsConnected() {
		ctx := context.Background()
		truncatedToken := token[:8] + "..."

		m.redis.SafeSet(ctx, fmt.Sprintf("apikey:%s:instance", token), m.instanceID, clientExpiry)
		m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:instance", id), m.instanceID, clientExpiry)
		m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:apikey", id), token, clientExpiry)

		if worldID != "" {
			m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:worldId", id), worldID, clientExpiry)
		}
		if worldTitle != "" {
			m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:worldTitle", id), worldTitle, clientExpiry)
		}
		if foundryVersion != "" {
			m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:foundryVersion", id), foundryVersion, clientExpiry)
		}
		if systemID != "" {
			m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:systemId", id), systemID, clientExpiry)
		}
		if systemTitle != "" {
			m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:systemTitle", id), systemTitle, clientExpiry)
		}
		if systemVersion != "" {
			m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:systemVersion", id), systemVersion, clientExpiry)
		}
		if customName != "" {
			m.redis.SafeSet(ctx, fmt.Sprintf("client:%s:customName", id), customName, clientExpiry)
		}

		m.redis.SafeSAdd(ctx, fmt.Sprintf("apikey:%s:clients", token), id)
		m.redis.SafeExpire(ctx, fmt.Sprintf("apikey:%s:clients", token), clientExpiry)

		log.Info().Str("clientId", id).Str("token", truncatedToken).Msg("Client registered in Redis")
	}

	truncatedToken := token
	if len(token) > 8 {
		truncatedToken = token[:8] + "..."
	}
	log.Info().Str("clientId", id).Str("token", truncatedToken).Msg("Client connected")
	return client, nil
}

// RemoveClient removes a client and cleans up state.
func (m *ClientManager) RemoveClient(id string) {
	m.mu.Lock()
	client, exists := m.clients[id]
	if !exists {
		m.mu.Unlock()
		return
	}

	token := client.APIKey()

	// Notify listeners (e.g., terminate sheet sessions)
	if m.OnClientRemoved != nil {
		m.OnClientRemoved(id)
	}

	// Clean up local state
	delete(m.clients, id)
	if group, ok := m.tokenGroups[token]; ok {
		delete(group, id)
		if len(group) == 0 {
			delete(m.tokenGroups, token)
		}
	}
	m.mu.Unlock()

	// Clean up Redis
	if m.redis != nil && m.redis.IsConnected() {
		ctx := context.Background()
		m.redis.SafeDel(ctx, fmt.Sprintf("client:%s:instance", id))
		m.redis.SafeDel(ctx, fmt.Sprintf("client:%s:apikey", id))
		m.redis.SafeSRem(ctx, fmt.Sprintf("apikey:%s:clients", token), id)

		remaining, _ := m.redis.SafeSCard(ctx, fmt.Sprintf("apikey:%s:clients", token))
		if remaining == 0 {
			m.redis.SafeDel(ctx, fmt.Sprintf("apikey:%s:instance", token))
			m.redis.SafeDel(ctx, fmt.Sprintf("apikey:%s:clients", token))
		}
	}

	log.Info().Str("clientId", id).Msg("Client disconnected")
}

// GetClient returns a client by ID, or nil if not found locally.
func (m *ClientManager) GetClient(id string) *Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[id]
}

// GetClientInstanceID checks Redis for which instance hosts a client.
func (m *ClientManager) GetClientInstanceID(ctx context.Context, id string) (string, error) {
	if m.redis == nil || !m.redis.IsConnected() {
		return "", nil
	}
	return m.redis.SafeGet(ctx, fmt.Sprintf("client:%s:instance", id))
}

// GetInstanceForAPIKey checks Redis for which instance serves an API key.
func (m *ClientManager) GetInstanceForAPIKey(ctx context.Context, apiKey string) (string, error) {
	if m.redis == nil || !m.redis.IsConnected() {
		return "", nil
	}
	return m.redis.SafeGet(ctx, fmt.Sprintf("apikey:%s:instance", apiKey))
}

// GetConnectedClients returns client IDs connected with the given API key.
func (m *ClientManager) GetConnectedClients(apiKey string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	group := m.tokenGroups[apiKey]
	if group == nil {
		return nil
	}

	var ids []string
	for id := range group {
		if client, ok := m.clients[id]; ok && client.IsAlive() {
			ids = append(ids, id)
		}
	}
	return ids
}

// GetConnectedClientInfos returns full client info for an API key.
func (m *ClientManager) GetConnectedClientInfos(apiKey string) []ClientInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	group := m.tokenGroups[apiKey]
	if group == nil {
		return nil
	}

	var infos []ClientInfo
	for id := range group {
		if client, ok := m.clients[id]; ok && client.IsAlive() {
			infos = append(infos, client.Info(m.instanceID))
		}
	}
	return infos
}

// UpdateClientLastSeen updates a client's last seen timestamp and refreshes Redis TTLs.
func (m *ClientManager) UpdateClientLastSeen(id string) {
	m.mu.RLock()
	client, exists := m.clients[id]
	m.mu.RUnlock()

	if !exists {
		return
	}

	client.UpdateLastSeen()

	if m.redis != nil && m.redis.IsConnected() {
		ctx := context.Background()
		token := client.APIKey()
		m.redis.SafeExpire(ctx, fmt.Sprintf("client:%s:instance", id), clientExpiry)
		m.redis.SafeExpire(ctx, fmt.Sprintf("client:%s:apikey", id), clientExpiry)
		m.redis.SafeExpire(ctx, fmt.Sprintf("apikey:%s:clients", token), clientExpiry)
		m.redis.SafeExpire(ctx, fmt.Sprintf("apikey:%s:instance", token), clientExpiry)
	}
}

// OnMessageType registers a handler for a specific message type.
func (m *ClientManager) OnMessageType(msgType string, handler MessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messageHandlers[msgType] = append(m.messageHandlers[msgType], handler)
}

// OnRawMessageType registers a handler that also receives the raw message bytes.
func (m *ClientManager) OnRawMessageType(msgType string, handler RawMessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rawMessageHandlers[msgType] = append(m.rawMessageHandlers[msgType], handler)
}

// HandleIncomingMessageFast routes a message using only the type string and raw bytes.
// Full JSON parse is deferred until a non-raw handler needs it.
func (m *ClientManager) HandleIncomingMessageFast(clientID string, msgType string, rawBytes []byte) {
	m.mu.RLock()
	client, exists := m.clients[clientID]
	if !exists {
		m.mu.RUnlock()
		return
	}
	m.mu.RUnlock()

	client.UpdateLastSeen()

	if msgType == "ping" {
		client.Send(map[string]string{"type": "pong"})
		return
	}

	// Check for raw handlers first (no parse needed)
	m.mu.RLock()
	rawHandlers := m.rawMessageHandlers[msgType]
	handlers := m.messageHandlers[msgType]
	m.mu.RUnlock()

	if len(rawHandlers) > 0 {
		// Fast extract requestId for raw handlers
		var partialMsg struct {
			RequestID string `json:"requestId"`
			Error     *string `json:"error"`
		}
		json.Unmarshal(rawBytes, &partialMsg) // Partial parse — only extracts 2 fields

		data := map[string]interface{}{
			"requestId": partialMsg.RequestID,
		}
		if partialMsg.Error != nil {
			data["error"] = *partialMsg.Error
		}

		for _, handler := range rawHandlers {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error().Str("clientId", clientID).Str("messageType", msgType).Interface("panic", r).Msg("Panic in raw handler")
					}
				}()
				handler(client, data, rawBytes)
			}()
		}
		return
	}

	// Fall back to full parse for regular handlers
	if len(handlers) > 0 {
		var message map[string]interface{}
		if err := json.Unmarshal(rawBytes, &message); err != nil {
			return
		}
		for _, handler := range handlers {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error().Str("clientId", clientID).Str("messageType", msgType).Interface("panic", r).Msg("Panic in handler")
					}
				}()
				handler(client, message)
			}()
		}
		return
	}

	// Broadcast unhandled — needs full parse
	var message map[string]interface{}
	if err := json.Unmarshal(rawBytes, &message); err != nil {
		return
	}
	m.BroadcastToGroup(clientID, message)
}

// HandleIncomingMessage routes an incoming message to registered handlers (full parse already done).
func (m *ClientManager) HandleIncomingMessage(clientID string, message map[string]interface{}, rawBytes []byte) {
	m.mu.RLock()
	client, exists := m.clients[clientID]
	if !exists {
		m.mu.RUnlock()
		return
	}
	m.mu.RUnlock()

	client.UpdateLastSeen()

	msgType, _ := message["type"].(string)

	// Handle ping
	if msgType == "ping" {
		client.Send(map[string]string{"type": "pong"})
		return
	}

	// Dispatch to registered handlers
	m.mu.RLock()
	handlers := m.messageHandlers[msgType]
	rawHandlers := m.rawMessageHandlers[msgType]
	m.mu.RUnlock()

	if len(rawHandlers) > 0 {
		for _, handler := range rawHandlers {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error().Str("clientId", clientID).Str("messageType", msgType).Interface("panic", r).Msg("Panic in raw handler")
					}
				}()
				handler(client, message, rawBytes)
			}()
		}
		return
	}

	if len(handlers) > 0 {
		for _, handler := range handlers {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error().Str("clientId", clientID).Str("messageType", msgType).Interface("panic", r).Msg("Panic in handler")
					}
				}()
				handler(client, message)
			}()
		}
		return
	}

	// Broadcast unhandled messages to the group
	m.BroadcastToGroup(clientID, message)
}

// BroadcastToGroup sends a message to all clients sharing the same API key.
func (m *ClientManager) BroadcastToGroup(senderID string, message interface{}) {
	m.mu.RLock()
	sender, exists := m.clients[senderID]
	if !exists {
		m.mu.RUnlock()
		return
	}

	token := sender.APIKey()
	group := m.tokenGroups[token]
	if group == nil {
		m.mu.RUnlock()
		return
	}

	// Collect targets while holding the lock
	var targets []*Client
	for id := range group {
		if id != senderID {
			if client, ok := m.clients[id]; ok && client.IsAlive() {
				targets = append(targets, client)
			}
		}
	}
	m.mu.RUnlock()

	for _, client := range targets {
		client.Send(message)
	}
}

// CleanupInactiveClients removes clients that are no longer alive.
func (m *ClientManager) CleanupInactiveClients() {
	m.mu.RLock()
	var toRemove []string
	for id, client := range m.clients {
		if !client.IsAlive() {
			toRemove = append(toRemove, id)
		}
	}
	m.mu.RUnlock()

	for _, id := range toRemove {
		log.Info().Str("clientId", id).Msg("Removing inactive client")
		m.RemoveClient(id)
	}
}

// StartCleanupLoop starts a goroutine that periodically cleans up inactive clients.
func (m *ClientManager) StartCleanupLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.CleanupInactiveClients()
			case <-ctx.Done():
				return
			}
		}
	}()
}
