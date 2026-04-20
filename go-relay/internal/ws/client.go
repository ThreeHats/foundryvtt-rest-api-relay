package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// Client represents a connected Foundry VTT instance.
type Client struct {
	mu                sync.Mutex
	conn              *websocket.Conn
	id                string
	apiKey            string
	connectionTokenID int64  // 0 if connected via legacy API key (no connection token)
	tokenName         string // human-readable name of the connection token (e.g. "Firefox on Linux")
	lastSeen          time.Time
	connectedSince    time.Time
	connected         bool
	worldID           string
	worldTitle        string
	foundryVersion    string
	systemID          string
	systemTitle       string
	systemVersion     string
	customName        string
	ipAddress         string
	sendCh            chan []byte
	done              chan struct{}
}

// ClientInfo holds metadata about a client for API responses.
type ClientInfo struct {
	ID             string `json:"clientId"`
	InstanceID     string `json:"instanceId,omitempty"`
	LastSeen       int64  `json:"lastSeen"`
	ConnectedSince int64  `json:"connectedSince"`
	WorldID        string `json:"worldId,omitempty"`
	WorldTitle     string `json:"worldTitle,omitempty"`
	FoundryVersion string `json:"foundryVersion,omitempty"`
	SystemID       string `json:"systemId,omitempty"`
	SystemTitle    string `json:"systemTitle,omitempty"`
	SystemVersion  string `json:"systemVersion,omitempty"`
	CustomName     string `json:"customName,omitempty"`
	IPAddress      string `json:"ipAddress,omitempty"`
	TokenName      string `json:"tokenName,omitempty"`
}

// NewClient creates a new Client and starts its write goroutine.
func NewClient(conn *websocket.Conn, id, apiKey, tokenName string, worldID, worldTitle, foundryVersion, systemID, systemTitle, systemVersion, customName string) *Client {
	remoteAddr := ""
	if conn.RemoteAddr() != nil {
		remoteAddr = conn.RemoteAddr().String()
	}

	c := &Client{
		conn:           conn,
		id:             id,
		apiKey:         apiKey,
		tokenName:      tokenName,
		lastSeen:       time.Now(),
		connectedSince: time.Now(),
		connected:      true,
		worldID:        worldID,
		worldTitle:     worldTitle,
		foundryVersion: foundryVersion,
		systemID:       systemID,
		systemTitle:    systemTitle,
		systemVersion:  systemVersion,
		customName:     customName,
		ipAddress:      remoteAddr,
		sendCh:         make(chan []byte, 1024),
		done:           make(chan struct{}),
	}

	// Write goroutine serializes all writes to avoid concurrent write panics
	go c.writePump()

	return c
}

// writePump reads from sendCh and writes to the WebSocket connection.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.sendCh:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Error().Err(err).Str("clientId", c.id).Msg("Write error")
				return
			}
		case <-c.done:
			return
		}
	}
}

// Send sends a JSON message to the client. Thread-safe.
func (c *Client) Send(data interface{}) bool {
	if !c.IsAlive() {
		return false
	}

	msg, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Str("clientId", c.id).Msg("Marshal error")
		return false
	}

	select {
	case c.sendCh <- msg:
		return true
	default:
		log.Warn().Str("clientId", c.id).Msg("Send channel full, dropping message")
		return false
	}
}

// SendRaw sends a raw JSON byte slice. Thread-safe.
func (c *Client) SendRaw(msg []byte) bool {
	if !c.IsAlive() {
		return false
	}
	select {
	case c.sendCh <- msg:
		return true
	default:
		return false
	}
}

// IsAlive returns true if the client connection is still active.
func (c *Client) IsAlive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

// UpdateLastSeen refreshes the client's last activity timestamp.
func (c *Client) UpdateLastSeen() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSeen = time.Now()
}

// Disconnect closes the client connection.
func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connected {
		c.connected = false
		close(c.done)
		c.conn.Close()
	}
}

// MarkDisconnected marks the client as disconnected without closing the socket.
func (c *Client) MarkDisconnected() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = false
}

// SetConnectionTokenID records which connection token authenticated this client.
// Called by the relay handler after successful auth-via-first-message.
func (c *Client) SetConnectionTokenID(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connectionTokenID = id
}

// ConnectionTokenID returns the ID of the connection token used to authenticate
// this client. Returns 0 if the client connected via legacy API key.
func (c *Client) ConnectionTokenID() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connectionTokenID
}

// Getters
func (c *Client) ID() string               { return c.id }
func (c *Client) APIKey() string            { return c.apiKey }
func (c *Client) IPAddress() string         { return c.ipAddress }
func (c *Client) WorldID() string           { return c.worldID }
func (c *Client) WorldTitle() string        { return c.worldTitle }
func (c *Client) FoundryVersion() string    { return c.foundryVersion }
func (c *Client) SystemID() string          { return c.systemID }
func (c *Client) SystemTitle() string       { return c.systemTitle }
func (c *Client) SystemVersion() string     { return c.systemVersion }
func (c *Client) CustomName() string        { return c.customName }
func (c *Client) TokenName() string         { return c.tokenName }
func (c *Client) LastSeen() time.Time       { return c.lastSeen }
func (c *Client) ConnectedSince() time.Time { return c.connectedSince }
func (c *Client) Conn() *websocket.Conn     { return c.conn }

// Info returns client metadata for API responses.
func (c *Client) Info(instanceID string) ClientInfo {
	return ClientInfo{
		ID:             c.id,
		InstanceID:     instanceID,
		LastSeen:       c.lastSeen.UnixMilli(),
		ConnectedSince: c.connectedSince.UnixMilli(),
		WorldID:        c.worldID,
		WorldTitle:     c.worldTitle,
		FoundryVersion: c.foundryVersion,
		SystemID:       c.systemID,
		SystemTitle:    c.systemTitle,
		SystemVersion:  c.systemVersion,
		CustomName:     c.customName,
		IPAddress:      c.ipAddress,
		TokenName:      c.tokenName,
	}
}
