package handler

import (
	"crypto"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/rs/zerolog/log"
)

// PendingHandshake stores RSA key pair info for a session handshake.
type PendingHandshake struct {
	PrivateKey   *rsa.PrivateKey
	PublicKeyPEM string
	Nonce        string
	FoundryURL   string
	Username     string
	WorldName    string
	APIKey       string
	InstanceID   string
	ExpiresAt    time.Time
}

var (
	handshakeMu sync.RWMutex
	handshakes  = make(map[string]*PendingHandshake)
)

// Create a handshake token for secure authentication
//
// @tag Session
// @param {string} x-api-key [header,required] API key header
// @param {string} x-foundry-url [header,required] Foundry URL header
// @param {string} x-world-name [header] World name header
// @param {string} x-username [header,required] Username header
// @returns Handshake token and encryption details
func sessionHandshakeHandler(db *database.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := helpers.GetRequestContext(r)
		if reqCtx == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		foundryURL := r.Header.Get("x-foundry-url")
		username := r.Header.Get("x-username")
		worldName := r.Header.Get("x-world-name")

		if foundryURL == "" || username == "" {
			helpers.WriteError(w, http.StatusBadRequest, "x-foundry-url and x-username headers are required")
			return
		}

		privateKey, err := rsa.GenerateKey(cryptorand.Reader, 2048)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate RSA key pair")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate handshake credentials")
			return
		}

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to marshal public key")
			return
		}

		publicKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes}))

		nonceBytes := make([]byte, 16)
		cryptorand.Read(nonceBytes)
		nonce := fmt.Sprintf("%x", nonceBytes)

		tokenBytes := make([]byte, 32)
		cryptorand.Read(tokenBytes)
		handshakeToken := fmt.Sprintf("%x", tokenBytes)

		hs := &PendingHandshake{
			PrivateKey: privateKey, PublicKeyPEM: publicKeyPEM, Nonce: nonce,
			FoundryURL: foundryURL, Username: username, WorldName: worldName,
			APIKey: reqCtx.MasterAPIKey, InstanceID: cfg.InstanceID(),
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		handshakeMu.Lock()
		handshakes[handshakeToken] = hs
		handshakeMu.Unlock()

		log.Info().Str("foundryUrl", foundryURL).Str("username", username).Msg("Session handshake created")

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"token": handshakeToken, "publicKey": publicKeyPEM, "nonce": nonce,
			"instanceId": cfg.InstanceID(), "foundryUrl": foundryURL, "username": username,
			"expires": hs.ExpiresAt.Format(time.RFC3339),
		})
	}
}

// Start a headless Foundry session using puppeteer
//
// @tag Session
// @param {string} handshakeToken [body,required] The token received from session-handshake
// @param {string} encryptedPassword [body,required] Password encrypted with the public key
// @param {string} captureBrowserConsole [body] Log level for browser console capture ("error", "warn", or "debug")
// @param {string} x-api-key [header,required] API key header
// @returns Session information including sessionId and clientId
func sessionStartHandler(db *database.DB, cfg *config.Config, headless *worker.HeadlessManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := parseBody(r)
		handshakeToken := bodyStr(body, "handshakeToken")
		encryptedPassword := bodyStr(body, "encryptedPassword")

		if handshakeToken == "" || encryptedPassword == "" {
			helpers.WriteError(w, http.StatusBadRequest, "handshakeToken and encryptedPassword are required")
			return
		}

		handshakeMu.Lock()
		hs, exists := handshakes[handshakeToken]
		if exists {
			delete(handshakes, handshakeToken)
		}
		handshakeMu.Unlock()

		if !exists || hs.ExpiresAt.Before(time.Now()) {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid or expired handshake token")
			return
		}

		// Decrypt password using RSA private key
		encBytes, err := base64.StdEncoding.DecodeString(encryptedPassword)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid encrypted password encoding")
			return
		}

		// Node.js RSA_PKCS1_OAEP_PADDING defaults to SHA-1
		decrypted, err := rsa.DecryptOAEP(sha1.New(), cryptorand.Reader, hs.PrivateKey, encBytes, nil)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Failed to decrypt password")
			return
		}

		// Parse decrypted JSON: { "password": "...", "nonce": "..." }
		var creds struct {
			Password string `json:"password"`
			Nonce    string `json:"nonce"`
		}
		if err := json.Unmarshal(decrypted, &creds); err != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid decrypted credential format")
			return
		}

		// Verify nonce
		if creds.Nonce != hs.Nonce {
			helpers.WriteError(w, http.StatusUnauthorized, "Nonce mismatch")
			return
		}

		worldName := hs.WorldName
		if wn := bodyStr(body, "worldName"); wn != "" {
			worldName = wn
		}

		if headless == nil {
			helpers.WriteError(w, http.StatusServiceUnavailable, "Headless sessions not available")
			return
		}

		// Launch headless session (this blocks until client connects or timeout)
		sessionID, clientID, err := headless.LaunchSession(hs.APIKey, hs.FoundryURL, hs.Username, creds.Password, worldName)
		if err != nil {
			log.Error().Err(err).Msg("Headless session launch failed")
			helpers.WriteJSON(w, http.StatusRequestTimeout, map[string]interface{}{
				"error":   err.Error(),
				"message": "Failed to start headless session. This may happen if another GM is already connected.",
			})
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"success":   true,
			"sessionId": sessionID,
			"clientId":  clientID,
		})
	}
}

// Stop a headless Foundry session
//
// @tag Session
// @param {string} sessionId [query,required] The ID of the session to end
// @param {string} x-api-key [header,required] API key header
// @returns Status of the operation
func sessionEndHandler(headless *worker.HeadlessManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("sessionId")
		if sessionID == "" {
			body := parseBody(r)
			sessionID = bodyStr(body, "sessionId")
		}

		if sessionID == "" {
			helpers.WriteError(w, http.StatusBadRequest, "sessionId is required")
			return
		}

		if headless == nil {
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
			return
		}

		if err := headless.EndSession(sessionID); err != nil {
			helpers.WriteError(w, http.StatusNotFound, err.Error())
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "message": "Session ended"})
	}
}

// Get all active headless Foundry sessions
//
// @tag Session
// @param {string} x-api-key [header,required] API key header
// @returns List of active sessions for the current API key
func sessionListHandler(headless *worker.HeadlessManager, mgr *ws.ClientManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if headless == nil {
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"activeSessions": []interface{}{}})
			return
		}

		sessions := headless.ListSessions()
		// Enrich with client metadata
		var activeSessions []map[string]interface{}
		for _, s := range sessions {
			info := map[string]interface{}{
				"sessionId":    s.SessionID,
				"clientId":     s.ClientID,
				"foundryUrl":   s.FoundryURL,
				"username":     s.Username,
				"worldName":    s.WorldName,
				"startedAt":    s.StartedAt,
				"lastActivity": s.LastActivity,
			}
			// Add client metadata if connected
			if client := mgr.GetClient(s.ClientID); client != nil {
				info["worldId"] = client.WorldID()
				info["worldTitle"] = client.WorldTitle()
				info["foundryVersion"] = client.FoundryVersion()
				info["systemId"] = client.SystemID()
				info["systemTitle"] = client.SystemTitle()
				info["systemVersion"] = client.SystemVersion()
			}
			activeSessions = append(activeSessions, info)
		}
		if activeSessions == nil {
			activeSessions = []map[string]interface{}{}
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"activeSessions": activeSessions})
	}
}

// CleanupExpiredHandshakes removes expired handshake tokens.
func CleanupExpiredHandshakes() {
	handshakeMu.Lock()
	defer handshakeMu.Unlock()
	now := time.Now()
	for token, hs := range handshakes {
		if hs.ExpiresAt.Before(now) {
			delete(handshakes, token)
		}
	}
}

// unused import guards
var _ = crypto.SHA256
