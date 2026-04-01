package ws

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// fastExtractType extracts the "type" field from JSON without full parsing.
// Returns empty string if not found.
func fastExtractType(data []byte) string {
	// Look for "type":" pattern
	needle := []byte(`"type":"`)
	idx := bytes.Index(data, needle)
	if idx < 0 {
		return ""
	}
	start := idx + len(needle)
	end := bytes.IndexByte(data[start:], '"')
	if end < 0 {
		return ""
	}
	return string(data[start : start+end])
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (CORS handled elsewhere)
	},
}

// RelayConfig holds configuration for the WebSocket relay.
type RelayConfig struct {
	PingInterval      time.Duration
	CleanupInterval   time.Duration
	ValidateAPIKey    func(token string) (*APIKeyValidation, error)
	ValidateHeadless  func(clientID, token string) (bool, error)
}

// APIKeyValidation holds the result of API key validation.
type APIKeyValidation struct {
	Valid          bool
	MasterAPIKey   string
	ScopedClientID string
	ScopedUserID   string
}

// HandleRelayConnection handles a WebSocket upgrade for the Foundry module relay.
func HandleRelayConnection(manager *ClientManager, cfg *RelayConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		query := r.URL.Query()
		id := query.Get("id")
		token := query.Get("token")
		worldID := query.Get("worldId")
		worldTitle := query.Get("worldTitle")
		foundryVersion := query.Get("foundryVersion")
		systemID := query.Get("systemId")
		systemTitle := query.Get("systemTitle")
		systemVersion := query.Get("systemVersion")
		customName := query.Get("customName")

		if id == "" || token == "" {
			log.Warn().Msg("Rejecting WebSocket connection: missing id or token")
			http.Error(w, "Missing client ID or token", http.StatusBadRequest)
			return
		}

		// Validate API key
		var registrationToken string
		if cfg.ValidateAPIKey != nil {
			validation, err := cfg.ValidateAPIKey(token)
			if err != nil || !validation.Valid {
				log.Warn().Str("clientId", id).Msg("Rejecting WebSocket: invalid API key")
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			registrationToken = token
			if validation.MasterAPIKey != "" && validation.MasterAPIKey != token {
				// Scoped key — validate clientId constraint
				if validation.ScopedClientID != "" && validation.ScopedClientID != id {
					log.Warn().
						Str("expected", validation.ScopedClientID).
						Str("got", id).
						Msg("Rejecting WebSocket: scoped key clientId mismatch")
					http.Error(w, "Client ID does not match scoped API key", http.StatusForbidden)
					return
				}
				// Register under parent's master key
				registrationToken = validation.MasterAPIKey
			}
		} else {
			registrationToken = token
		}

		// Validate headless session
		if cfg.ValidateHeadless != nil {
			valid, err := cfg.ValidateHeadless(id, token)
			if err != nil || !valid {
				log.Warn().Str("clientId", id).Msg("Rejecting invalid headless client")
				http.Error(w, "Invalid headless session", http.StatusForbidden)
				return
			}
		}

		// Upgrade to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("WebSocket upgrade failed")
			return
		}

		// Register client
		client, err := manager.AddClient(conn, id, registrationToken, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName)
		if err != nil {
			return
		}

		// Set up ping/pong
		conn.SetPongHandler(func(appData string) error {
			client.UpdateLastSeen()
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

			for {
				_, messageBytes, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
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
					// For raw handlers, we can pass nil for the parsed map
					// and let the handler use rawBytes directly.
					// For regular handlers, we need to parse.
					manager.HandleIncomingMessageFast(id, msgType, messageBytes)
				}
			}
		}()
	}
}
