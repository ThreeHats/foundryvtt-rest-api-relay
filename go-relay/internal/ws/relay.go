package ws

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// fastExtractType extracts the "type" field from the top-level JSON object.
// It tracks brace depth and skips strings properly to avoid matching a "type"
// key inside a nested object (e.g. {"data":{"type":"fake"},"type":"real"}).
// Returns empty string if not found.
func fastExtractType(data []byte) string {
	depth := 0
	i := 0
	for i < len(data) {
		b := data[i]
		switch {
		case b == '{':
			depth++
			i++
		case b == '}':
			depth--
			i++
		case b == '"':
			// Scan string, tracking start/end
			i++ // skip opening quote
			keyStart := i
			for i < len(data) {
				if data[i] == '\\' {
					i += 2
					continue
				}
				if data[i] == '"' {
					break
				}
				i++
			}
			keyEnd := i
			i++ // skip closing quote

			if depth != 1 {
				continue
			}

			// Skip whitespace + colon between key and value
			for i < len(data) && (data[i] == ' ' || data[i] == ':' || data[i] == '\t' || data[i] == '\n') {
				i++
			}

			if string(data[keyStart:keyEnd]) == "type" && i < len(data) && data[i] == '"' {
				i++ // skip opening quote of value
				valStart := i
				for i < len(data) {
					if data[i] == '\\' {
						i += 2
						continue
					}
					if data[i] == '"' {
						break
					}
					i++
				}
				return string(data[valStart:i])
			}
		default:
			i++
		}
	}
	return ""
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  32 * 1024,
	WriteBufferSize: 32 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (CORS handled elsewhere)
	},
}

// RelayConfig holds configuration for the WebSocket relay.
type RelayConfig struct {
	PingInterval            time.Duration
	CleanupInterval         time.Duration
	ValidateAPIKey          func(token string) (*APIKeyValidation, error)
	ValidateConnectionToken func(token string) (masterAPIKey string, allowedIPs string, tokenName string, tokenID int64, err error)
	ValidateHeadless        func(clientID, token string) (bool, error)
	// OnClientConnected is called after a Foundry module client successfully
	// connects and registers. Used to sync current subscription state to the module.
	OnClientConnected func(clientID string)
}

// isIPAllowed checks whether remoteAddr is permitted by a comma-separated
// list of IP addresses and/or CIDR ranges. Returns true if the list is empty.
func isIPAllowed(remoteAddr, allowedIPs string) bool {
	if allowedIPs == "" {
		return true
	}

	// Extract just the IP from host:port
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// remoteAddr might already be a bare IP
		host = remoteAddr
	}
	connectingIP := net.ParseIP(host)
	if connectingIP == nil {
		return false
	}

	for _, entry := range strings.Split(allowedIPs, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		// Try CIDR first
		_, cidr, err := net.ParseCIDR(entry)
		if err == nil {
			if cidr.Contains(connectingIP) {
				return true
			}
			continue
		}

		// Fall back to exact IP match
		if ip := net.ParseIP(entry); ip != nil && ip.Equal(connectingIP) {
			return true
		}
	}

	return false
}

// APIKeyValidation holds the result of API key validation.
type APIKeyValidation struct {
	Valid          bool
	UserID         int64 // relay DB user ID; 0 if unknown (e.g. connection token auth)
	MasterAPIKey   string
	ScopedClientID string
	ScopedUserID   string
}

// relayValidateToken validates a token for the relay (Foundry module) connection.
// Only connection tokens are accepted — scoped API keys cannot authenticate
// to /relay (they're for the consumer-facing /ws/api endpoint instead).
//
// Returns (registrationToken, allowedIPs, tokenName, tokenID, validation, error).
// allowedIPs is the token's IP allowlist (may be empty for unrestricted tokens).
// tokenName is the human-readable name of the connection token (e.g. "Firefox on Linux").
// tokenID is the database ID of the connection token (0 if not found).
func relayValidateToken(cfg *RelayConfig, token, id string) (string, string, string, int64, *APIKeyValidation, error) {
	// Memory mode — no validators configured, accept any token
	if cfg.ValidateConnectionToken == nil && cfg.ValidateAPIKey == nil {
		return token, "", "", 0, &APIKeyValidation{Valid: true}, nil
	}

	// Connection token is the ONLY accepted auth method for /relay
	if cfg.ValidateConnectionToken == nil {
		return "", "", "", 0, nil, fmt.Errorf("connection token validation not configured")
	}

	masterAPIKey, allowedIPs, tokenName, tokenID, err := cfg.ValidateConnectionToken(token)
	if err != nil || masterAPIKey == "" {
		return "", "", "", 0, nil, fmt.Errorf("invalid connection token")
	}

	log.Debug().Str("clientId", id).Int64("tokenId", tokenID).Msg("Relay auth via connection token")
	return masterAPIKey, allowedIPs, tokenName, tokenID, &APIKeyValidation{Valid: true, MasterAPIKey: masterAPIKey}, nil
}

// HandleRelayConnection handles a WebSocket upgrade for the Foundry module relay.
func HandleRelayConnection(manager *ClientManager, cfg *RelayConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters (metadata only — token is sent in the first WS message)
		query := r.URL.Query()
		id := query.Get("id")
		worldID := query.Get("worldId")
		worldTitle := query.Get("worldTitle")
		foundryVersion := query.Get("foundryVersion")
		systemID := query.Get("systemId")
		systemTitle := query.Get("systemTitle")
		systemVersion := query.Get("systemVersion")
		customName := query.Get("customName")

		if id == "" {
			log.Warn().Str("ip", r.RemoteAddr).Msg("Rejecting WebSocket connection: missing id")
			http.Error(w, "Missing client ID", http.StatusBadRequest)
			return
		}

		// Token URL parameter is not supported — authenticate using the first WebSocket message.
		if query.Get("token") != "" {
			http.Error(w, "Token URL parameter is not supported. Send auth as the first WebSocket message.", http.StatusBadRequest)
			return
		}

		// Auth via first message after upgrade
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warn().Err(err).Str("ip", r.RemoteAddr).Msg("WebSocket upgrade failed")
			return
		}

		// Set a 10-second deadline for the auth message
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Warn().Str("clientId", id).Err(err).Msg("Auth message read failed")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(4002, "Authentication timeout"))
			conn.Close()
			return
		}

		// Parse auth message
		var authMsg struct {
			Type  string `json:"type"`
			Token string `json:"token"`
		}
		if err := json.Unmarshal(messageBytes, &authMsg); err != nil || authMsg.Type != "auth" || authMsg.Token == "" {
			log.Warn().Str("clientId", id).Msg("Invalid auth message")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(4002, "Invalid auth message"))
			conn.Close()
			return
		}

		registrationToken, allowedIPs, tokenName, tokenID, _, err := relayValidateToken(cfg, authMsg.Token, id)
		if err != nil {
			log.Warn().Str("clientId", id).Err(err).Msg("Rejecting WebSocket: auth failed")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(1008, "Authentication failed"))
			conn.Close()
			return
		}

		// Validate connection token IP allowlist
		if allowedIPs != "" && !isIPAllowed(r.RemoteAddr, allowedIPs) {
			log.Warn().Str("clientId", id).Str("remoteAddr", r.RemoteAddr).Msg("Rejecting WebSocket: IP not in allowlist")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(4002, "IP address not allowed"))
			conn.Close()
			return
		}

		// Validate headless session
		if cfg.ValidateHeadless != nil {
			valid, err := cfg.ValidateHeadless(id, authMsg.Token)
			if err != nil || !valid {
				log.Warn().Str("clientId", id).Msg("Rejecting invalid headless client")
				conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(4002, "Invalid headless session"))
				conn.Close()
				return
			}
		}

		// Clear the read deadline
		conn.SetReadDeadline(time.Time{})

		// All validation passed — send auth-success before starting client pumps
		sendWSJSON(conn, map[string]interface{}{"type": "auth-success"})

		startRelayClient(manager, cfg, conn, id, registrationToken, authMsg.Token, tokenName, tokenID, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName)
	}
}

// startRelayClient registers the client and starts ping/read pumps.
// tokenID is the connection token's database ID (0 if connected via legacy API key);
// it is recorded on the Client so DELETE /auth/connection-tokens/:id can target this connection.
// tokenName is the human-readable name of the connection token (e.g. "Firefox on Linux"),
// stored on the client for inclusion in connection log entries.
func startRelayClient(manager *ClientManager, cfg *RelayConfig, conn *websocket.Conn, id, registrationToken, rawToken, tokenName string, tokenID int64, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName string) {
	// Register client
	client, err := manager.AddClient(conn, id, registrationToken, tokenName, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName)
	if err != nil {
		return
	}
	if tokenID != 0 {
		// Go through the manager so the token ID is also mirrored into Redis
		// for cross-instance lookups (used by the broadcast disconnect path).
		manager.SetClientConnectionTokenID(client.ID(), tokenID)
	}

	// Notify server so it can sync current subscription state to this module.
	if cfg.OnClientConnected != nil {
		cfg.OnClientConnected(client.ID())
	}

	// Set up ping/pong — extend read deadline on each pong so a stalled
	// connection is detected within 2× the ping interval.
	conn.SetPongHandler(func(appData string) error {
		client.UpdateLastSeen()
		conn.SetReadDeadline(time.Now().Add(2 * cfg.PingInterval))
		return nil
	})

	// Ping goroutine
	go func() {
		ticker := time.NewTicker(cfg.PingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if !client.IsAlive() {
					return
				}
				if err := conn.WriteControl(websocket.PingMessage, []byte("keepalive"), time.Now().Add(5*time.Second)); err != nil {
					log.Debug().Err(err).Str("clientId", id).Msg("Ping failed")
					return
				}
			case <-client.done:
				return
			}
		}
	}()

	// Read pump goroutine
	go func() {
		defer func() {
			manager.RemoveClient(id)
			conn.Close()
		}()

		conn.SetReadDeadline(time.Now().Add(2 * cfg.PingInterval))

		for {
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					log.Error().Err(err).Str("clientId", id).Msg("WebSocket read error")
				}
				return
			}

			// Fast path: extract only "type" without full JSON parse.
			// Full parse is deferred to handlers that need it.
			msgType := fastExtractType(messageBytes)
			if msgType == "" {
				// Fallback to full parse for malformed messages
				var message map[string]interface{}
				if err := json.Unmarshal(messageBytes, &message); err != nil {
					log.Error().Err(err).Str("clientId", id).Msg("Invalid JSON message")
					continue
				}
				msgType, _ = message["type"].(string)
				manager.HandleIncomingMessage(id, message, messageBytes)
			} else {
				manager.HandleIncomingMessageFast(id, msgType, messageBytes)
			}
		}
	}()
}
