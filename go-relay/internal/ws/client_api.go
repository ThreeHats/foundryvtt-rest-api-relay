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
	PingInterval            time.Duration
	ValidateAPIKey          func(token string) (*APIKeyValidation, error)
	ValidateConnectionToken func(token string) (masterAPIKey string, allowedIPs string, tokenName string, tokenID int64, err error)
	TrackUsage              func(apiKey string) (bool, string)
	// AutoStart attempts to start a headless session for a scoped key with stored credentials.
	// Args: masterAPIKey, scopedClientID, scopedUserID. Returns the new clientID or empty string.
	AutoStart      func(masterAPIKey, scopedClientID, scopedUserID string) string
	SSEManager     SSEManagerInterface
	InteractiveSessions  *InteractiveSessionManager
}

// clientWSState tracks the state of a client API WebSocket connection.
type clientWSState struct {
	mu               sync.Mutex
	apiKey           string
	masterAPIKey     string
	userID           int64  // relay DB user ID; 0 for connection-token auth or unknown
	clientID         string
	scopedUserID     string
	conn              *websocket.Conn
	consumerID        string            // stable identity for interactive session cleanup
	pendingRequestIDs map[string]string // internalID -> clientRequestID
	subscriptions     []*WSEventSub
	done              chan struct{}
	semaphore         chan struct{} // limits concurrent goroutines spawned per connection
}

// clientAPIValidateToken validates a token for the client API connection, trying connection token first, then API key.
// Returns (matchKey, allowedIPs, scopedUserID, validation, error).
// allowedIPs is non-empty only for connection token auth.
func clientAPIValidateToken(cfg *ClientAPIConfig, token string) (matchKey string, allowedIPs string, scopedUserID string, validation *APIKeyValidation, err error) {
	// Try connection token first
	if cfg.ValidateConnectionToken != nil {
		masterAPIKey, ips, _, _, ctErr := cfg.ValidateConnectionToken(token)
		if ctErr == nil && masterAPIKey != "" {
			return masterAPIKey, ips, "", &APIKeyValidation{Valid: true, MasterAPIKey: masterAPIKey}, nil
		}
		// Fall through to API key validation
	}

	// Fall back to API key validation
	if cfg.ValidateAPIKey != nil {
		validation, err = cfg.ValidateAPIKey(token)
		if err != nil || !validation.Valid {
			return "", "", "", nil, fmt.Errorf("invalid API key")
		}
		matchKey = validation.MasterAPIKey
		if matchKey == "" {
			matchKey = token
		}
		return matchKey, "", validation.ScopedUserID, validation, nil
	}

	return "", "", "", nil, fmt.Errorf("no validators configured")
}

// clientAPIResolveAndValidate validates the token and resolves the clientID.
// Returns the validated state fields or an error string suitable for HTTP response.
// allowedIPs is non-empty only for connection token auth and must be checked by the caller.
func clientAPIResolveAndValidate(manager *ClientManager, cfg *ClientAPIConfig, token, clientID string) (matchKey, allowedIPs, resolvedClientID, scopedUserID string, validation *APIKeyValidation, httpErr string, httpStatus int) {
	truncatedToken := token
	if len(truncatedToken) > 8 {
		truncatedToken = truncatedToken[:8] + "..."
	}

	matchKey, allowedIPs, scopedUserID, validation, err := clientAPIValidateToken(cfg, token)
	if err != nil {
		log.Warn().Str("token", truncatedToken).Msg("WS /ws/api rejected: invalid API key")
		return "", "", "", "", nil, "Invalid API key", http.StatusUnauthorized
	}

	truncatedMatchKey := matchKey
	if len(truncatedMatchKey) > 8 {
		truncatedMatchKey = truncatedMatchKey[:8] + "..."
	}
	log.Debug().Str("token", truncatedToken).Str("matchKey", truncatedMatchKey).Msg("WS API key validated")

	// Apply scoped clientId constraint
	if validation != nil && validation.ScopedClientID != "" {
		clientID = validation.ScopedClientID
	}

	// Auto-resolve clientId
	if clientID == "" {
		clients := manager.GetConnectedClients(matchKey)
		log.Debug().Str("matchKey", matchKey[:8]+"...").Int("connectedClients", len(clients)).Msg("WS auto-resolving clientId")
		switch len(clients) {
		case 1:
			clientID = clients[0]
			log.Debug().Str("clientId", clientID).Msg("WS auto-resolved clientId")
		case 0:
			if cfg.AutoStart != nil && validation != nil {
				if autoID := cfg.AutoStart(matchKey, validation.ScopedClientID, validation.ScopedUserID); autoID != "" {
					clientID = autoID
					log.Info().Str("clientId", clientID).Msg("WS auto-started headless session")
					break
				}
			}
			log.Warn().Str("matchKey", matchKey[:8]+"...").Msg("WS /ws/api rejected: no connected clients")
			return "", "", "", "", nil, "No connected Foundry client found", http.StatusNotFound
		default:
			log.Warn().Int("count", len(clients)).Msg("WS /ws/api rejected: multiple clients")
			return "", "", "", "", nil, "Multiple clients connected. Please specify clientId.", http.StatusBadRequest
		}
	}

	// Verify client exists and belongs to this API key
	foundryClient := manager.GetClient(clientID)
	if foundryClient == nil {
		log.Warn().Str("clientId", clientID).Msg("WS /ws/api rejected: client not found")
		return "", "", "", "", nil, "Invalid clientId", http.StatusNotFound
	}
	if foundryClient.APIKey() != matchKey {
		log.Warn().Str("clientId", clientID).Msg("WS /ws/api rejected: key mismatch")
		return "", "", "", "", nil, "API key does not match the specified clientId", http.StatusForbidden
	}

	return matchKey, allowedIPs, clientID, scopedUserID, validation, "", 0
}

// startClientWSPumps sets up the client WS state, sends welcome, and starts ping/read pumps.
func startClientWSPumps(conn *websocket.Conn, manager *ClientManager, pending *PendingRequests, cfg *ClientAPIConfig, token, matchKey, clientID, scopedUserID string, userID int64) {
	state := &clientWSState{
		apiKey:            token,
		masterAPIKey:      matchKey,
		userID:            userID,
		clientID:          clientID,
		scopedUserID:      scopedUserID,
		conn:              conn,
		consumerID:        randomStr(8),
		pendingRequestIDs: make(map[string]string),
		done:              make(chan struct{}),
		semaphore:         make(chan struct{}, 100),
	}

	kp := token
	if len(kp) > 8 {
		kp = kp[:8]
	}
	log.Info().Str("clientId", clientID).Int64("userId", userID).Str("keyPrefix", kp).Msg("Client WS connected")

	// Send welcome
	sendWSJSON(conn, map[string]interface{}{
		"type":           "connected",
		"clientId":       clientID,
		"supportedTypes": pendingRequestTypesList(),
		"eventChannels":  []string{"chat-events", "roll-events", "hooks", "combat-events", "actor-events", "scene-events"},
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
			dkp := state.apiKey
			if len(dkp) > 8 {
				dkp = dkp[:8]
			}
			log.Info().Str("clientId", clientID).Int64("userId", state.userID).Str("keyPrefix", dkp).Msg("Client WS disconnected")
		}()

		for {
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
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

// HandleClientAPIConnection handles WebSocket connections on /ws/api.
//
// Auth-via-first-message: after the WebSocket upgrade, the client must send
// a JSON message {"type":"auth","token":"<key>"} within 10 seconds. Both
// connection tokens and API keys (master or scoped) are accepted.
func HandleClientAPIConnection(manager *ClientManager, pending *PendingRequests, cfg *ClientAPIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("url", r.URL.String()).Msg("WS /ws/api connection attempt")

		query := r.URL.Query()
		clientID := query.Get("clientId")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("Client API WebSocket upgrade failed")
			return
		}

		// Auth-via-first-message ONLY. Tokens must never be passed in URL params
		// (they appear in server access logs and are visible to proxies).
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("WS /ws/api auth message read failed")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(4002, "Authentication timeout"))
			conn.Close()
			return
		}

		var authMsg struct {
			Type     string `json:"type"`
			Token    string `json:"token"`
			ClientID string `json:"clientId"`
		}
		if err := json.Unmarshal(messageBytes, &authMsg); err != nil || authMsg.Type != "auth" || authMsg.Token == "" {
			log.Warn().Msg("WS /ws/api invalid auth message")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(4002, "Invalid auth message"))
			conn.Close()
			return
		}

		token := authMsg.Token
		if authMsg.ClientID != "" {
			clientID = authMsg.ClientID
		}

		matchKey, allowedIPs, resolvedClientID, scopedUserID, validation, wsErr, _ := clientAPIResolveAndValidate(manager, cfg, token, clientID)
		if wsErr != "" {
			log.Warn().Str("error", wsErr).Msg("WS /ws/api auth rejected")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(4002, wsErr))
			conn.Close()
			return
		}

		// Validate connection token IP allowlist
		if allowedIPs != "" && !isIPAllowed(r.RemoteAddr, allowedIPs) {
			log.Warn().Str("remoteAddr", r.RemoteAddr).Msg("WS /ws/api rejected: IP not in allowlist")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(4002, "IP address not allowed"))
			conn.Close()
			return
		}

		var userID int64
		if validation != nil {
			userID = validation.UserID
		}

		conn.SetReadDeadline(time.Time{})
		startClientWSPumps(conn, manager, pending, cfg, token, matchKey, resolvedClientID, scopedUserID, userID)
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

	// Interactive session messages
	if msgType == "interactive-session-start" {
		handleInteractiveSessionStart(state, manager, cfg, msg)
		return
	}
	if msgType == "interactive-input" {
		handleInteractiveInput(state, manager, cfg, msg)
		return
	}
	if msgType == "interactive-session-end" {
		handleInteractiveSessionEnd(state, manager, cfg, msg)
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

	// Acquire semaphore slot — cap concurrent response goroutines per connection to 100.
	// If the slot isn't available immediately, reject the in-flight request rather than
	// spawning an unbounded number of goroutines.
	select {
	case state.semaphore <- struct{}{}:
	default:
		pending.Delete(internalID)
		state.mu.Lock()
		delete(state.pendingRequestIDs, internalID)
		state.mu.Unlock()
		sendWSJSON(state.conn, map[string]interface{}{
			"type":      fmt.Sprintf("%s-result", msgType),
			"requestId": requestID,
			"error":     "Too many concurrent requests; try again shortly",
		})
		return
	}

	// Wait for response in a goroutine
	go func() {
		defer func() { <-state.semaphore }()
		select {
		case resp := <-responseCh:
			state.mu.Lock()
			clientReqID := state.pendingRequestIDs[internalID]
			delete(state.pendingRequestIDs, internalID)
			state.mu.Unlock()

			if resp != nil {
				var responseData map[string]interface{}
				if resp.Data != nil {
					responseData = resp.Data
				} else if resp.RawData != nil {
					if err := json.Unmarshal(resp.RawData, &responseData); err != nil {
						responseData = map[string]interface{}{"error": "failed to parse response"}
					}
				}
				if responseData != nil {
					responseData["type"] = fmt.Sprintf("%s-result", msgType)
					responseData["requestId"] = clientReqID
					responseData["clientId"] = state.clientID
					sendWSJSON(state.conn, responseData)
				}
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

	validChannels := map[string]bool{
		"chat-events": true, "roll-events": true,
		"hooks": true, "combat-events": true, "actor-events": true, "scene-events": true,
	}
	if !validChannels[channel] {
		sendWSJSON(state.conn, map[string]interface{}{
			"type":      "error",
			"error":     fmt.Sprintf("Invalid channel: %q. Supported: chat-events, roll-events, hooks, combat-events, actor-events, scene-events", channel),
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
	log.Debug().Str("clientId", state.clientID).Str("channel", channel).Msg("Client WS subscribed")
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
	log.Debug().Str("clientId", state.clientID).Str("channel", channel).Int("removed", removed).Msg("Client WS unsubscribed")
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

func handleInteractiveSessionStart(state *clientWSState, manager *ClientManager, cfg *ClientAPIConfig, msg map[string]interface{}) {
	if cfg.InteractiveSessions == nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "interactive-session-error", "error": "Interactive sessions not available"})
		return
	}

	foundryClient := manager.GetClient(state.clientID)
	if foundryClient == nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "interactive-session-error", "error": "Foundry client is no longer connected"})
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

	sessionID, err := cfg.InteractiveSessions.CreateSession(state.clientID, state.apiKey, state.consumerID, state.conn, InteractiveSessionMetadata{
		UUID: uuid, Quality: quality, Scale: scale,
	})
	if err != nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "interactive-session-error", "error": err.Error()})
		return
	}

	userId := state.scopedUserID
	if userId == "" {
		if u, ok := msg["userId"].(string); ok {
			userId = u
		}
	}

	payload := map[string]interface{}{
		"type": "interactive-session-start", "sessionId": sessionID,
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
		cfg.InteractiveSessions.EndSession(sessionID)
		sendWSJSON(state.conn, map[string]interface{}{"type": "interactive-session-error", "error": "Failed to send session start to Foundry"})
	}
}

func handleInteractiveInput(state *clientWSState, manager *ClientManager, cfg *ClientAPIConfig, msg map[string]interface{}) {
	if cfg.InteractiveSessions == nil {
		return
	}

	sessionID, _ := msg["sessionId"].(string)
	session := cfg.InteractiveSessions.GetSession(sessionID)
	if session == nil || session.ConsumerConn != state.conn {
		sendWSJSON(state.conn, map[string]interface{}{"type": "interactive-session-error", "sessionId": sessionID, "error": "Invalid session"})
		return
	}

	cfg.InteractiveSessions.UpdateActivity(sessionID)

	foundryClient := manager.GetClient(state.clientID)
	if foundryClient == nil {
		sendWSJSON(state.conn, map[string]interface{}{"type": "interactive-session-error", "sessionId": sessionID, "error": "Foundry client disconnected"})
		cfg.InteractiveSessions.EndSession(sessionID)
		return
	}

	foundryClient.Send(map[string]interface{}{
		"type": "interactive-input", "sessionId": sessionID,
		"action": msg["action"], "x": msg["x"], "y": msg["y"], "button": msg["button"],
		"deltaX": msg["deltaX"], "deltaY": msg["deltaY"],
		"key": msg["key"], "code": msg["code"], "modifiers": msg["modifiers"],
	})
}

func handleInteractiveSessionEnd(state *clientWSState, manager *ClientManager, cfg *ClientAPIConfig, msg map[string]interface{}) {
	if cfg.InteractiveSessions == nil {
		return
	}

	sessionID, _ := msg["sessionId"].(string)
	cfg.InteractiveSessions.EndSession(sessionID)

	foundryClient := manager.GetClient(state.clientID)
	if foundryClient != nil {
		foundryClient.Send(map[string]interface{}{"type": "interactive-session-end", "sessionId": sessionID})
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

	// Clean up interactive sessions — notify Foundry to close them
	if cfg != nil && cfg.InteractiveSessions != nil {
		sessionIDs := cfg.InteractiveSessions.TerminateSessionsForConsumer(state.consumerID)
		if len(sessionIDs) > 0 {
			foundryClient := manager.GetClient(state.clientID)
			if foundryClient != nil {
				for _, sid := range sessionIDs {
					foundryClient.Send(map[string]interface{}{"type": "interactive-session-end", "sessionId": sid})
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
