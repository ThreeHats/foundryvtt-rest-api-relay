package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/metrics"
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
	// Receives clientID and the API key the client was registered under.
	OnClientRemoved func(clientID, apiKey string)
	// OnClientConnected is called when a new Foundry client connects. Set by server after init.
	OnClientConnected func(clientID, apiKey string, info ClientInfo)
	// OnDuplicateConnectionRejected is called when an incoming WS connection
	// is rejected because the clientId is already held by another live client.
	// newIP is the incoming connection's address; existingIP is the address of
	// the already-connected client. Callers can compare the two to distinguish
	// a same-machine multi-tab reload (suppress) from a remote takeover attempt
	// (alert).
	OnDuplicateConnectionRejected func(clientID, apiKey, newIP, existingIP string)
}

const clientExpiry = 2 * time.Hour

// disconnectChannel is the Redis pub/sub channel used to broadcast forced
// disconnect commands across all relay instances. Every instance running
// StartDisconnectSubscriber listens on this channel and acts locally when
// a matching client or connection token is found.
const disconnectChannel = "relay:disconnect"

// disconnectCmd is the payload format for disconnect broadcasts.
// Exactly one of ClientID / ConnectionTokenID / APIKey should be set.
type disconnectCmd struct {
	Kind              string `json:"kind"` // "client" | "token" | "apikey"
	ClientID          string `json:"clientId,omitempty"`
	ConnectionTokenID int64  `json:"tokenId,omitempty"`
	APIKey            string `json:"apiKey,omitempty"`
	Reason            string `json:"reason,omitempty"`
	CloseCode         int    `json:"closeCode,omitempty"` // WS close code, default 4002
	Origin            string `json:"origin,omitempty"`    // instance ID of the publisher
}

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
func (m *ClientManager) AddClient(conn *websocket.Conn, id, token, tokenName, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName string) (*Client, error) {
	m.mu.Lock()

	// Reject duplicate connections
	if existing, exists := m.clients[id]; exists {
		existingIP := existing.IPAddress() // capture before unlock
		m.mu.Unlock()
		newIP := ""
		if conn.RemoteAddr() != nil {
			newIP = conn.RemoteAddr().String()
		}
		log.Warn().Str("clientId", id).Str("newIp", newIP).Str("existingIp", existingIP).Msg("Client already exists, rejecting connection")
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(4004, "Client ID already connected"))
		conn.Close()
		if m.OnDuplicateConnectionRejected != nil {
			go m.OnDuplicateConnectionRejected(id, token, newIP, existingIP)
		}
		return nil, fmt.Errorf("duplicate client ID: %s", id)
	}

	client := NewClient(conn, id, token, tokenName, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName)
	m.clients[id] = client

	// Add to token group
	if m.tokenGroups[token] == nil {
		m.tokenGroups[token] = make(map[string]bool)
	}
	m.tokenGroups[token][id] = true

	m.mu.Unlock()

	// Store in Redis (non-blocking, best effort) — pipeline all writes into one round-trip.
	if m.redis != nil && m.redis.IsConnected() {
		ctx := context.Background()
		truncatedToken := token[:8] + "..."
		pipe := m.redis.Pipeline()
		if pipe != nil {
			pipe.Set(ctx, fmt.Sprintf("apikey:%s:instance", token), m.instanceID, clientExpiry)
			pipe.Set(ctx, fmt.Sprintf("client:%s:instance", id), m.instanceID, clientExpiry)
			pipe.Set(ctx, fmt.Sprintf("client:%s:apikey", id), token, clientExpiry)
			if worldID != "" {
				pipe.Set(ctx, fmt.Sprintf("client:%s:worldId", id), worldID, clientExpiry)
			}
			if worldTitle != "" {
				pipe.Set(ctx, fmt.Sprintf("client:%s:worldTitle", id), worldTitle, clientExpiry)
			}
			if foundryVersion != "" {
				pipe.Set(ctx, fmt.Sprintf("client:%s:foundryVersion", id), foundryVersion, clientExpiry)
			}
			if systemID != "" {
				pipe.Set(ctx, fmt.Sprintf("client:%s:systemId", id), systemID, clientExpiry)
			}
			if systemTitle != "" {
				pipe.Set(ctx, fmt.Sprintf("client:%s:systemTitle", id), systemTitle, clientExpiry)
			}
			if systemVersion != "" {
				pipe.Set(ctx, fmt.Sprintf("client:%s:systemVersion", id), systemVersion, clientExpiry)
			}
			if customName != "" {
				pipe.Set(ctx, fmt.Sprintf("client:%s:customName", id), customName, clientExpiry)
			}
			pipe.SAdd(ctx, fmt.Sprintf("apikey:%s:clients", token), id)
			pipe.Expire(ctx, fmt.Sprintf("apikey:%s:clients", token), clientExpiry)
			if _, err := pipe.Exec(ctx); err != nil {
				log.Warn().Err(err).Str("clientId", id).Msg("Redis pipeline exec failed")
			} else {
				log.Info().Str("clientId", id).Str("token", truncatedToken).Msg("Client registered in Redis")
			}
		}
	}

	log.Info().
		Str("clientId", id).
		Str("worldId", client.WorldID()).
		Str("foundryVersion", client.FoundryVersion()).
		Str("systemId", client.SystemID()).
		Str("tokenName", client.TokenName()).
		Msg("Client connected")

	// Notify listeners (e.g., upsert known client, log connection, send notification)
	if m.OnClientConnected != nil {
		m.OnClientConnected(id, token, client.Info(m.instanceID))
	}

	metrics.WSConnectionsActive.Inc()
	return client, nil
}

// SetClientConnectionTokenID records the connection token ID that
// authenticated a client, both on the local Client struct and in Redis
// (so cross-instance handlers can look it up later for revocation).
//
// Called by the relay connection handler after AddClient.
func (m *ClientManager) SetClientConnectionTokenID(clientID string, tokenID int64) {
	m.mu.RLock()
	client := m.clients[clientID]
	m.mu.RUnlock()
	if client == nil {
		return
	}
	client.SetConnectionTokenID(tokenID)

	if tokenID != 0 && m.redis != nil && m.redis.IsConnected() {
		ctx := context.Background()
		m.redis.SafeSet(ctx,
			fmt.Sprintf("client:%s:connectionTokenId", clientID),
			fmt.Sprintf("%d", tokenID),
			clientExpiry,
		)
	}
}

// LookupClientConnectionTokenID returns the connection token ID associated
// with a client, checking local memory first and falling back to Redis for
// clients connected to other instances. Returns 0 if no token ID is known.
func (m *ClientManager) LookupClientConnectionTokenID(ctx context.Context, clientID string) int64 {
	m.mu.RLock()
	client := m.clients[clientID]
	m.mu.RUnlock()
	if client != nil {
		if id := client.ConnectionTokenID(); id != 0 {
			return id
		}
	}
	if m.redis != nil && m.redis.IsConnected() {
		val, err := m.redis.SafeGet(ctx, fmt.Sprintf("client:%s:connectionTokenId", clientID))
		if err == nil && val != "" {
			if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
				return parsed
			}
		}
	}
	return 0
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

	// Notify listeners (e.g., terminate interactive sessions)
	if m.OnClientRemoved != nil {
		m.OnClientRemoved(id, token)
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
		m.redis.SafeDel(ctx, fmt.Sprintf("client:%s:connectionTokenId", id))
		m.redis.SafeSRem(ctx, fmt.Sprintf("apikey:%s:clients", token), id)

		remaining, _ := m.redis.SafeSCard(ctx, fmt.Sprintf("apikey:%s:clients", token))
		if remaining == 0 {
			m.redis.SafeDel(ctx, fmt.Sprintf("apikey:%s:instance", token))
			m.redis.SafeDel(ctx, fmt.Sprintf("apikey:%s:clients", token))
		}
	}

	log.Info().
		Str("clientId", id).
		Str("worldId", client.WorldID()).
		Str("foundryVersion", client.FoundryVersion()).
		Str("systemId", client.SystemID()).
		Msg("Client disconnected")
	metrics.WSConnectionsActive.Dec()
}

// SnapshotLocalClients returns a point-in-time copy of all clients connected
// to this instance. Used by the reconciliation loop in server.go to
// independently validate each client's credentials against the database.
//
// The returned slice is safe to iterate outside of any lock — the underlying
// *Client pointers remain valid until their natural Disconnect.
func (m *ClientManager) SnapshotLocalClients() []*Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Client, 0, len(m.clients))
	for _, c := range m.clients {
		out = append(out, c)
	}
	return out
}

// ForceDisconnectLocal closes a specific client on this instance with the
// given close code + reason. Used by the reconciliation loop.
func (m *ClientManager) ForceDisconnectLocal(c *Client, closeCode int, reason string) {
	if c == nil {
		return
	}
	m.closeClients([]*Client{c}, closeCode, reason, "reconciliation")
}

// GetClient returns a client by ID, or nil if not found locally.
// DisconnectByConnectionToken closes any WebSocket clients that authenticated
// with the given connection token ID. Returns the number of clients disconnected.
//
// Called by the DELETE /auth/connection-tokens/:id endpoint to immediately
// terminate any active sessions using the revoked token. Without this, a
// revoked token would still allow its existing connection to function until
// the next reconnect.
func (m *ClientManager) DisconnectByConnectionToken(tokenID int64) int {
	return m.disconnectLocalByToken(tokenID, 4002, "Connection token revoked")
}

// disconnectLocalByToken closes matching clients on this instance only.
func (m *ClientManager) disconnectLocalByToken(tokenID int64, closeCode int, reason string) int {
	if tokenID == 0 {
		return 0
	}
	m.mu.RLock()
	var toDisconnect []*Client
	for _, c := range m.clients {
		if c.ConnectionTokenID() == tokenID {
			toDisconnect = append(toDisconnect, c)
		}
	}
	m.mu.RUnlock()

	return m.closeClients(toDisconnect, closeCode, reason, "connection token revocation")
}

// disconnectLocalByClientID closes the matching client on this instance only.
func (m *ClientManager) disconnectLocalByClientID(clientID string, closeCode int, reason string) int {
	if clientID == "" {
		return 0
	}
	m.mu.RLock()
	c := m.clients[clientID]
	m.mu.RUnlock()
	if c == nil {
		return 0
	}
	return m.closeClients([]*Client{c}, closeCode, reason, "client-id disconnect broadcast")
}

// disconnectLocalByAPIKey closes every client registered under the given
// master API key on this instance only.
func (m *ClientManager) disconnectLocalByAPIKey(apiKey string, closeCode int, reason string) int {
	if apiKey == "" {
		return 0
	}
	m.mu.RLock()
	var toDisconnect []*Client
	if group, ok := m.tokenGroups[apiKey]; ok {
		for id := range group {
			if c, ok := m.clients[id]; ok {
				toDisconnect = append(toDisconnect, c)
			}
		}
	}
	m.mu.RUnlock()
	return m.closeClients(toDisconnect, closeCode, reason, "api-key disconnect broadcast")
}

func (m *ClientManager) closeClients(clients []*Client, closeCode int, reason, logReason string) int {
	if closeCode == 0 {
		closeCode = 4002
	}
	count := 0
	for _, c := range clients {
		log.Info().
			Str("clientId", c.ID()).
			Int("closeCode", closeCode).
			Str("reason", logReason).
			Msg("Force-disconnecting client")
		// Send close frame first so the module gets a concrete code.
		_ = c.Conn().WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(closeCode, reason))
		c.Disconnect()
		count++
	}
	return count
}

// BroadcastDisconnectByConnectionToken disconnects matching clients on THIS
// instance and publishes a Redis pub/sub command so every other instance
// does the same. Returns the number of clients disconnected locally.
//
// Safe to call when Redis is not configured — falls back to local-only.
func (m *ClientManager) BroadcastDisconnectByConnectionToken(ctx context.Context, tokenID int64, reason string) int {
	local := m.disconnectLocalByToken(tokenID, 4002, reason)
	if m.redis != nil && m.redis.IsConnected() {
		payload, err := json.Marshal(disconnectCmd{
			Kind:              "token",
			ConnectionTokenID: tokenID,
			Reason:            reason,
			CloseCode:         4002,
			Origin:            m.instanceID,
		})
		if err == nil {
			if pubErr := m.redis.SafePublish(ctx, disconnectChannel, payload); pubErr != nil {
				log.Warn().Err(pubErr).Msg("Failed to publish disconnect command")
			}
		}
	}
	return local
}

// BroadcastDisconnectByClientID does the same thing, keyed by client ID.
func (m *ClientManager) BroadcastDisconnectByClientID(ctx context.Context, clientID, reason string) int {
	local := m.disconnectLocalByClientID(clientID, 4002, reason)
	if m.redis != nil && m.redis.IsConnected() {
		payload, err := json.Marshal(disconnectCmd{
			Kind:      "client",
			ClientID:  clientID,
			Reason:    reason,
			CloseCode: 4002,
			Origin:    m.instanceID,
		})
		if err == nil {
			if pubErr := m.redis.SafePublish(ctx, disconnectChannel, payload); pubErr != nil {
				log.Warn().Err(pubErr).Msg("Failed to publish disconnect command")
			}
		}
	}
	return local
}

// BroadcastDisconnectByAPIKey disconnects all clients registered under the
// given master API key across every instance. Useful for user-level actions
// like key rotation, account deletion, or lockout.
func (m *ClientManager) BroadcastDisconnectByAPIKey(ctx context.Context, apiKey, reason string) int {
	local := m.disconnectLocalByAPIKey(apiKey, 4002, reason)
	if m.redis != nil && m.redis.IsConnected() {
		payload, err := json.Marshal(disconnectCmd{
			Kind:      "apikey",
			APIKey:    apiKey,
			Reason:    reason,
			CloseCode: 4002,
			Origin:    m.instanceID,
		})
		if err == nil {
			if pubErr := m.redis.SafePublish(ctx, disconnectChannel, payload); pubErr != nil {
				log.Warn().Err(pubErr).Msg("Failed to publish disconnect command")
			}
		}
	}
	return local
}

// StartDisconnectSubscriber subscribes to the Redis disconnect channel and
// acts on incoming commands locally. No-op if Redis isn't connected.
//
// Commands published by this instance are ignored (via Origin check) since
// the local disconnect already happened synchronously at publish time.
func (m *ClientManager) StartDisconnectSubscriber(ctx context.Context) {
	if m.redis == nil || !m.redis.IsConnected() {
		log.Info().Msg("Redis not connected; disconnect broadcaster running in single-instance mode")
		return
	}

	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			pubsub := m.redis.Subscribe(ctx, disconnectChannel)
			if pubsub == nil {
				time.Sleep(5 * time.Second)
				continue
			}
			log.Info().Str("channel", disconnectChannel).Msg("Subscribed to disconnect broadcast channel")
			ch := pubsub.Channel()
		recvLoop:
			for {
				select {
				case <-ctx.Done():
					_ = pubsub.Close()
					return
				case msg, ok := <-ch:
					if !ok {
						_ = pubsub.Close()
						break recvLoop
					}
					m.handleDisconnectCommand(msg.Payload)
				}
			}
			log.Warn().Msg("Disconnect subscriber channel closed, reconnecting in 5s")
			time.Sleep(5 * time.Second)
		}
	}()
}

func (m *ClientManager) handleDisconnectCommand(payload string) {
	var cmd disconnectCmd
	if err := json.Unmarshal([]byte(payload), &cmd); err != nil {
		log.Warn().Err(err).Str("payload", payload).Msg("Failed to parse disconnect command")
		return
	}
	// Skip messages we published ourselves — the local disconnect already ran.
	if cmd.Origin == m.instanceID {
		return
	}
	reason := cmd.Reason
	if reason == "" {
		reason = "Revoked by dashboard"
	}
	closeCode := cmd.CloseCode
	if closeCode == 0 {
		closeCode = 4002
	}
	switch cmd.Kind {
	case "token":
		n := m.disconnectLocalByToken(cmd.ConnectionTokenID, closeCode, reason)
		if n > 0 {
			log.Info().Int64("tokenId", cmd.ConnectionTokenID).Int("count", n).Msg("Acted on token disconnect broadcast")
		}
	case "client":
		n := m.disconnectLocalByClientID(cmd.ClientID, closeCode, reason)
		if n > 0 {
			log.Info().Str("clientId", cmd.ClientID).Msg("Acted on client disconnect broadcast")
		}
	case "apikey":
		n := m.disconnectLocalByAPIKey(cmd.APIKey, closeCode, reason)
		if n > 0 {
			log.Info().Int("count", n).Msg("Acted on api-key disconnect broadcast")
		}
	default:
		log.Warn().Str("kind", cmd.Kind).Msg("Unknown disconnect command kind")
	}
}

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

// GetAllClientInfos returns info for ALL connected clients across all API keys.
// Used by admin endpoints to view system-wide WebSocket connections.
func (m *ClientManager) GetAllClientInfos() []ClientInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	infos := make([]ClientInfo, 0, len(m.clients))
	for _, client := range m.clients {
		if client.IsAlive() {
			infos = append(infos, client.Info(m.instanceID))
		}
	}
	return infos
}

// CountConnectedClients returns the total number of alive connected clients.
func (m *ClientManager) CountConnectedClients() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, c := range m.clients {
		if c.IsAlive() {
			count++
		}
	}
	return count
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
