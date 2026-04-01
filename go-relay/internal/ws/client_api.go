package ws

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// WSEventSub represents a WebSocket event subscription tracked per connection.
type WSEventSub struct {
	ClientID string
	Channel  string // "chat-events" or "roll-events"
	SendFunc func(data interface{}) bool
	Remove   func() // Cleanup function to call on unsubscribe/disconnect
}

// SSEManagerInterface allows the WS layer to manage event subscriptions without importing handler/helpers.
type SSEManagerInterface interface {
	AddWSEventFunc(clientID, channel string, sendFunc func(data interface{}) bool) (remove func(), ok bool)
}

// ClientAPIConfig holds configuration for the client-facing WebSocket API.
type ClientAPIConfig struct {
	PingInterval   time.Duration
	ValidateAPIKey func(token string) (*APIKeyValidation, error)
	TrackUsage     func(apiKey string) (bool, string)
	// AutoStart attempts to start a headless session for a scoped key with stored credentials.
	// Args: masterAPIKey, scopedClientID, scopedUserID. Returns the new clientID or empty string.
	AutoStart      func(masterAPIKey, scopedClientID, scopedUserID string) string
	SSEManager     SSEManagerInterface
	SheetSessions  *SheetSessionManager
}

// clientWSState tracks the state of a client API WebSocket connection.
type clientWSState struct {
	mu               sync.Mutex
	apiKey           string
	masterAPIKey     string
	clientID         string
	scopedUserID     string
	conn              *websocket.Conn
	pendingRequestIDs map[string]string // internalID -> clientRequestID
	subscriptions     []*WSEventSub
	done              chan struct{}
}

// HandleClientAPIConnection handles WebSocket connections on /ws/api.
func HandleClientAPIConnection(manager *ClientManager, pending *PendingRequests, cfg *ClientAPIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("url", r.URL.String()).Msg("WS /ws/api connection attempt")

		query := r.URL.Query()
		token := query.Get("token")
		clientID := query.Get("clientId")

		if token == "" {
			log.Warn().Msg("WS /ws/api rejected: missing token")
			http.Error(w, "Missing token query parameter", http.StatusBadRequest)
			return
		}

		truncatedToken := token
		if len(truncatedToken) > 8 {
			truncatedToken = truncatedToken[:8] + "..."
		}

		// Validate API key
		var matchKey string
		var scopedUserID string
		var validation *APIKeyValidation
		if cfg.ValidateAPIKey != nil {
			var err error
			validation, err = cfg.ValidateAPIKey(token)
			if err != nil || !validation.Valid {
				log.Warn().Str("token", truncatedToken).Msg("WS /ws/api rejected: invalid API key")
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}
			matchKey = validation.MasterAPIKey
			if matchKey == "" {
				matchKey = token
			}
			scopedUserID = validation.ScopedUserID

			log.Debug().Str("token", truncatedToken).Str("matchKey", matchKey[:8]+"...").Str("scopedClientId", validation.ScopedClientID).Msg("WS API key validated")

			// Apply scoped clientId constraint
			if validation.ScopedClientID != "" {
				clientID = validation.ScopedClientID
			}
		} else {
			matchKey = token
		}

		// Auto-resolve clientId
		if clientID == "" {
			clients := manager.GetConnectedClients(matchKey)
			log.Debug().Str("matchKey", matchKey[:8]+"...").Int("connectedClients", len(clients)).Msg("WS auto-resolving clientId")
			switch len(clients) {
			case 1:
				clientID = clients[0]
				log.Info().Str("clientId", clientID).Msg("WS auto-resolved clientId")
			case 0:
				// Try auto-start headless session for scoped keys with stored credentials
				if cfg.AutoStart != nil && validation != nil {
					if autoID := cfg.AutoStart(matchKey, validation.ScopedClientID, validation.ScopedUserID); autoID != "" {
						clientID = autoID
						log.Info().Str("clientId", clientID).Msg("WS auto-started headless session")
						break
					}
				}
				log.Warn().Str("matchKey", matchKey[:8]+"...").Msg("WS /ws/api rejected: no connected clients")
				http.Error(w, "No connected Foundry client found", http.StatusNotFound)
				return
			default:
				log.Warn().Int("count", len(clients)).Msg("WS /ws/api rejected: multiple clients")
				http.Error(w, "Multiple clients connected. Please specify clientId.", http.StatusBadRequest)
				return
			}
		}

		// Verify client exists and belongs to this API key
		foundryClient := manager.GetClient(clientID)
		if foundryClient == nil {
			log.Warn().Str("clientId", clientID).Msg("WS /ws/api rejected: client not found")
			http.Error(w, "Invalid clientId", http.StatusNotFound)
			return
		}
		if foundryClient.APIKey() != matchKey {
			log.Warn().Str("clientId", clientID).Str("clientKey", foundryClient.APIKey()[:8]+"...").Str("matchKey", matchKey[:8]+"...").Msg("WS /ws/api rejected: key mismatch")
			http.Error(w, "API key does not match the specified clientId", http.StatusForbidden)
			return
		}

		// Upgrade to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("Client API WebSocket upgrade failed")
			return
		}

		state := &clientWSState{
			apiKey:           token,
			masterAPIKey:     matchKey,
			clientID:         clientID,
			scopedUserID:     scopedUserID,
			conn:             conn,
			pendingRequestIDs: make(map[string]string),
			done:             make(chan struct{}),
		}

		truncated := token
		if len(truncated) > 8 {
			truncated = truncated[:8] + "..."
		}
		log.Info().Str("clientId", clientID).Str("apiKey", truncated).Msg("Client WS connected")

		// Send welcome
		sendWSJSON(conn, map[string]interface{}{
			"type":           "connected",
			"clientId":       clientID,
			"supportedTypes": pendingRequestTypesList(),
			"eventChannels":  []string{"chat-events", "roll-events"},
		})

		// Ping keepalive
		go func() {
			ticker := time.NewTicker(cfg.PingInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
						return
					}
				case <-state.done:
					return
				}
			}
		}()

		// Read pump
		go func() {
			defer func() {
				close(state.done)
				conn.Close()
				cleanupClientWSState(state, pending, cfg, manager)
				log.Info().Str("clientId", clientID).Msg("Client WS disconnected")
			}()

			for {
				_, messageBytes, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						log.Error().Err(err).Str("clientId", clientID).Msg("Client WS read error")
					}
					return
				}

				var msg map[string]interface{}
				if err := json.Unmarshal(messageBytes, &msg); err != nil {
					sendWSJSON(conn, map[string]interface{}{"type": "error", "error": "Invalid JSON"})
					continue
				}

				// Track usage
				if cfg.TrackUsage != nil {
					allowed, errMsg := cfg.TrackUsage(state.masterAPIKey)
					if !allowed {
						requestID, _ := msg["requestId"].(string)
						sendWSJSON(conn, map[string]interface{}{"type": "error", "error": errMsg, "requestId": requestID})
						continue
					}
				}

				handleClientWSMessage(state, manager, pending, cfg, msg)
			}
		}()
	}
}

func handleClientWSMessage(state *clientWSState, manager *ClientManager, pending *PendingRequests, cfg *ClientAPIConfig, msg map[string]interface{}) {
	msgType, _ := msg["type"].(string)
	requestID, _ := msg["requestId"].(string)

	if msgType == "" {
		sendWSJSON(state.conn, map[string]interface{}{"type": "error", "error": "Missing \"type\" field", "requestId": requestID})
		return
	}

	// Ping
	if msgType == "ping" {
		sendWSJSON(state.conn, map[string]interface{}{"type": "pong", "requestId": requestID})
		return
	}

	// Subscribe/unsubscribe
	if msgType == "subscribe" {
		handleWSSubscribe(state, manager, msg, cfg)
		return
	}
	if msgType == "unsubscribe" {
		handleWSUnsubscribe(state, msg)
		return
	}

	// Sheet session messages
	if msgType == "sheet-session-start" {
		handleSheetSessionStart(state, manager, cfg, msg)
		return
	}
	if msgType == "sheet-input" {
		handleSheetInput(state, manager, cfg, msg)
		return
	}
	if msgType == "sheet-session-end" {
		handleSheetSessionEnd(state, manager, cfg, msg)
		return
	}

	// Validate message type
	if !PendingRequestTypes[msgType] {
		sendWSJSON(state.conn, map[string]interface{}{
			"type":      "error",
			"error":     fmt.Sprintf("Unknown message type: %q", msgType),
			"requestId": requestID,
		})
		return
	}

	if requestID == "" {
		sendWSJSON(state.conn, map[string]interface{}{"type": "error", "error": "Missing \"requestId\" field for request messages"})
		return
	}

	// Get Foundry client
	foundryClient := manager.GetClient(state.clientID)
	if foundryClient == nil {
		sendWSJSON(state.conn, map[string]interface{}{
			"type":      fmt.Sprintf("%s-result", msgType),
			"requestId": requestID,
			"error":     "Foundry client is no longer connected",
		})
		return
	}

	// Create internal request ID
	internalID := fmt.Sprintf("ws_%s_%d_%s", msgType, time.Now().UnixMilli(), randomStr(6))

	// Register pending request with WS callback
	responseCh := make(chan *WSResponse, 1)
	pending.Store(internalID, &PendingRequest{
		ResponseCh: responseCh,
		Type:       msgType,
		ClientID:   state.clientID,
		Timestamp:  time.Now(),
	})

	state.mu.Lock()
	state.pendingRequestIDs[internalID] = requestID
	state.mu.Unlock()

	// Build payload
	payload := make(map[string]interface{})
	for k, v := range msg {
		if k != "type" && k != "requestId" && k != "format" {
			payload[k] = v
		}
	}

	// Inject scoped userId
	if state.scopedUserID != "" {
		payload["userId"] = state.scopedUserID
	}

	payload["type"] = msgType
	payload["requestId"] = internalID
	if _, ok := payload["data"]; !ok {
		payload["data"] = map[string]interface{}{}
	}

	if !foundryClient.Send(payload) {
		pending.Delete(internalID)
		state.mu.Lock()
		delete(state.pendingRequestIDs, internalID)
		state.mu.Unlock()
		sendWSJSON(state.conn, map[string]interface{}{
			"type":      fmt.Sprintf("%s-result", msgType),
			"requestId": requestID,
			"error":     "Failed to send request to Foundry client",
		})
		return
	}

	// Wait for response in a goroutine
	go func() {
		select {
		case resp := <-responseCh:
			state.mu.Lock()
			clientReqID := state.pendingRequestIDs[internalID]
			delete(state.pendingRequestIDs, internalID)
			state.mu.Unlock()

			if resp != nil && resp.Data != nil {
				resp.Data["type"] = fmt.Sprintf("%s-result", msgType)
				resp.Data["requestId"] = clientReqID
				sendWSJSON(state.conn, resp.Data)
			}
		case <-time.After(30 * time.Second):
			pending.Delete(internalID)
			state.mu.Lock()
			clientReqID := state.pendingRequestIDs[internalID]
			delete(state.pendingRequestIDs, internalID)
			state.mu.Unlock()

			sendWSJSON(state.conn, map[string]interface{}{
				"type":      fmt.Sprintf("%s-result", msgType),
				"requestId": clientReqID,
				"error":     "Request timed out",
			})
		case <-state.done:
			pending.Delete(internalID)
		}
	}()
}

func handleWSSubscribe(state *clientWSState, manager *ClientManager, msg map[string]interface{}, cfg *ClientAPIConfig) {
	requestID, _ := msg["requestId"].(string)
	channel, _ := msg["channel"].(string)

	if channel != "chat-events" && channel != "roll-events" {
		sendWSJSON(state.conn, map[string]interface{}{
			"type":      "error",
			"error":     fmt.Sprintf("Invalid channel: %q. Supported: chat-events, roll-events", channel),
			"requestId": requestID,
		})
		return
	}

	if cfg.SSEManager == nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "subscribed", "channel": channel, "requestId": requestID})
		return
	}

	sendFunc := func(data interface{}) bool {
		return sendWSJSONSafe(state.conn, state.done, data)
	}

	remove, ok := cfg.SSEManager.AddWSEventFunc(state.clientID, channel, sendFunc)
	if !ok {
		sendWSJSON(state.conn, map[string]interface{}{
			"type":      "error",
			"error":     "Maximum event subscriptions reached for this client",
			"requestId": requestID,
		})
		return
	}

	state.mu.Lock()
	state.subscriptions = append(state.subscriptions, &WSEventSub{
		ClientID: state.clientID,
		Channel:  channel,
		SendFunc: sendFunc,
		Remove:   remove,
	})
	state.mu.Unlock()

	sendWSJSON(state.conn, map[string]interface{}{"type": "subscribed", "channel": channel, "requestId": requestID})
	log.Info().Str("clientId", state.clientID).Str("channel", channel).Msg("Client WS subscribed")
}

func handleWSUnsubscribe(state *clientWSState, msg map[string]interface{}) {
	requestID, _ := msg["requestId"].(string)
	channel, _ := msg["channel"].(string)

	state.mu.Lock()
	var remaining []*WSEventSub
	removed := 0
	for _, sub := range state.subscriptions {
		if channel == "" || sub.Channel == channel {
			if sub.Remove != nil {
				sub.Remove()
			}
			removed++
		} else {
			remaining = append(remaining, sub)
		}
	}
	state.subscriptions = remaining
	state.mu.Unlock()

	if channel == "" {
		channel = "all"
	}
	sendWSJSON(state.conn, map[string]interface{}{"type": "unsubscribed", "channel": channel, "removed": removed, "requestId": requestID})
	log.Info().Str("clientId", state.clientID).Str("channel", channel).Int("removed", removed).Msg("Client WS unsubscribed")
}

func sendWSJSONSafe(conn *websocket.Conn, done chan struct{}, data interface{}) bool {
	select {
	case <-done:
		return false
	default:
	}
	msg, err := json.Marshal(data)
	if err != nil {
		return false
	}
	return conn.WriteMessage(websocket.TextMessage, msg) == nil
}

func handleSheetSessionStart(state *clientWSState, manager *ClientManager, cfg *ClientAPIConfig, msg map[string]interface{}) {
	if cfg.SheetSessions == nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "sheet-session-error", "error": "Sheet sessions not available"})
		return
	}

	foundryClient := manager.GetClient(state.clientID)
	if foundryClient == nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "sheet-session-error", "error": "Foundry client is no longer connected"})
		return
	}

	uuid, _ := msg["uuid"].(string)
	quality := 0.9
	if q, ok := msg["quality"].(float64); ok {
		quality = q
	}
	scale := 1.0
	if s, ok := msg["scale"].(float64); ok {
		scale = s
	}

	sessionID, err := cfg.SheetSessions.CreateSession(state.clientID, state.apiKey, state.conn, SheetSessionMetadata{
		UUID: uuid, Quality: quality, Scale: scale,
	})
	if err != nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "sheet-session-error", "error": err.Error()})
		return
	}

	userId := state.scopedUserID
	if userId == "" {
		if u, ok := msg["userId"].(string); ok {
			userId = u
		}
	}

	payload := map[string]interface{}{
		"type": "sheet-session-start", "sessionId": sessionID,
		"uuid": uuid, "quality": quality, "scale": scale,
	}
	if selected, ok := msg["selected"].(bool); ok {
		payload["selected"] = selected
	}
	if actor, ok := msg["actor"].(bool); ok {
		payload["actor"] = actor
	}
	if userId != "" {
		payload["userId"] = userId
	}

	if !foundryClient.Send(payload) {
		cfg.SheetSessions.EndSession(sessionID)
		sendWSJSON(state.conn, map[string]interface{}{"type": "sheet-session-error", "error": "Failed to send session start to Foundry"})
	}
}

func handleSheetInput(state *clientWSState, manager *ClientManager, cfg *ClientAPIConfig, msg map[string]interface{}) {
	if cfg.SheetSessions == nil {
		return
	}

	sessionID, _ := msg["sessionId"].(string)
	session := cfg.SheetSessions.GetSession(sessionID)
	if session == nil || session.ConsumerConn != state.conn {
		sendWSJSON(state.conn, map[string]interface{}{"type": "sheet-session-error", "sessionId": sessionID, "error": "Invalid session"})
		return
	}

	cfg.SheetSessions.UpdateActivity(sessionID)

	foundryClient := manager.GetClient(state.clientID)
	if foundryClient == nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "sheet-session-error", "sessionId": sessionID, "error": "Foundry client disconnected"})
		cfg.SheetSessions.EndSession(sessionID)
		return
	}

	foundryClient.Send(map[string]interface{}{
		"type": "sheet-input", "sessionId": sessionID,
		"action": msg["action"], "x": msg["x"], "y": msg["y"], "button": msg["button"],
		"deltaX": msg["deltaX"], "deltaY": msg["deltaY"],
		"key": msg["key"], "code": msg["code"], "modifiers": msg["modifiers"],
	})
}

func handleSheetSessionEnd(state *clientWSState, manager *ClientManager, cfg *ClientAPIConfig, msg map[string]interface{}) {
	if cfg.SheetSessions == nil {
		return
	}

	sessionID, _ := msg["sessionId"].(string)
	cfg.SheetSessions.EndSession(sessionID)

	foundryClient := manager.GetClient(state.clientID)
	if foundryClient != nil {
		foundryClient.Send(map[string]interface{}{"type": "sheet-session-end", "sessionId": sessionID})
	}
}

func cleanupClientWSState(state *clientWSState, pending *PendingRequests, cfg *ClientAPIConfig, manager *ClientManager) {
	state.mu.Lock()
	for internalID := range state.pendingRequestIDs {
		pending.Delete(internalID)
	}
	state.pendingRequestIDs = nil
	// Clean up event subscriptions
	for _, sub := range state.subscriptions {
		if sub.Remove != nil {
			sub.Remove()
		}
	}
	state.subscriptions = nil
	state.mu.Unlock()

	// Clean up sheet sessions — notify Foundry to close them
	if cfg != nil && cfg.SheetSessions != nil {
		sessionIDs := cfg.SheetSessions.TerminateSessionsForConsumer(state.conn)
		if len(sessionIDs) > 0 {
			foundryClient := manager.GetClient(state.clientID)
			if foundryClient != nil {
				for _, sid := range sessionIDs {
					foundryClient.Send(map[string]interface{}{"type": "sheet-session-end", "sessionId": sid})
				}
			}
		}
	}
}

func sendWSJSON(conn *websocket.Conn, data interface{}) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	conn.WriteMessage(websocket.TextMessage, msg)
}

func pendingRequestTypesList() []string {
	types := make([]string, 0, len(PendingRequestTypes))
	for t := range PendingRequestTypes {
		types = append(types, t)
	}
	return types
}

func randomStr(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
