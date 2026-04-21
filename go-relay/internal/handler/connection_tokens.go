package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// RegisterConnectionTokenRoutes registers routes for connection token management and pairing.
//
// The ClientManager is required so that DELETE /auth/connection-tokens/:id can
// force-disconnect any active WebSocket clients that authenticated using the
// revoked token.
func RegisterConnectionTokenRoutes(r chi.Router, db *database.DB, cfg *config.Config, manager *ws.ClientManager) {
	// Authenticated routes (require master API key)
	r.Route("/connection-tokens", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		// POST /auth/connection-tokens — Generate a pairing code.
		// Body (optional):
		//   { "clientId": "fvtt_...", "allowedTargetClients": [...], "remoteScopes": [...] }
		// - clientId, when provided, must already exist in KnownClients for
		//   this user. Pairing this code reuses that clientId instead of
		//   minting a fresh one (used for the "add this browser" flow).
		// - allowedTargetClients lists the clientIds the resulting connection
		//   token may invoke remote-request operations against. Empty = none.
		// - remoteScopes lists the scope strings (e.g. "entity:write",
		//   "user:write") the resulting token holds for those operations.
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key to manage connection tokens.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			var body struct {
				ClientID string `json:"clientId"`
			}
			// Body is optional — ignore decode errors.
			_ = json.NewDecoder(r.Body).Decode(&body)

			// If clientId is provided, validate it belongs to this user.
			var boundClientID sql.NullString
			if body.ClientID != "" {
				kc, err := db.KnownClientStore().FindByClientID(r.Context(), user.ID, body.ClientID)
				if err != nil {
					helpers.WriteError(w, http.StatusInternalServerError, "Failed to validate clientId")
					return
				}
				if kc == nil {
					helpers.WriteError(w, http.StatusNotFound, "clientId not found in your known clients")
					return
				}
				boundClientID = sql.NullString{String: body.ClientID, Valid: true}
			}

			// Generate 6-char uppercase alphanumeric code
			code, err := generatePairingCode(6)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate pairing code")
				return
			}

			pairingCode := &model.PairingCode{
				UserID:    user.ID,
				Code:      code,
				ClientID:  boundClientID,
				ExpiresAt: model.NewSQLiteTime(time.Now().Add(10 * time.Minute)),
			}

			if err := db.PairingCodeStore().Create(r.Context(), pairingCode); err != nil {
				log.Error().Err(err).Msg("Failed to create pairing code")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to create pairing code")
				return
			}

			helpers.WriteJSON(w, http.StatusCreated, map[string]interface{}{
				"code":      code,
				"expiresAt": pairingCode.ExpiresAt.Time.Format(time.RFC3339),
				"expiresIn": 600, // seconds
			})
		})

		// GET /auth/connection-tokens — List active tokens
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			tokens, err := db.ConnectionTokenStore().FindAllByUser(r.Context(), user.ID)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to list connection tokens")
				return
			}

			// Apply pagination query params (limit default 50, max 200)
			limit := 50
			if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
				if l > 200 {
					l = 200
				}
				limit = l
			}
			offset := 0
			if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o > 0 {
				offset = o
			}
			total := len(tokens)
			if offset < total {
				tokens = tokens[offset:]
			} else {
				tokens = nil
			}
			if limit > 0 && len(tokens) > limit {
				tokens = tokens[:limit]
			}

			var result []map[string]interface{}
			for _, t := range tokens {
				result = append(result, map[string]interface{}{
					"id":                    t.ID,
					"name":                  t.Name,
					"clientId":              t.ClientID,
					"source":                t.Source,
					"allowedIps":            t.AllowedIPs.String,
					"allowedTargetClients":  t.GetAllowedTargets(),
					"remoteScopes":          t.GetRemoteScopes(),
					"remoteRequestsPerHour": t.RemoteRequestsPerHour,
					"lastUsedAt":            t.LastUsedAt,
					"createdAt":             t.CreatedAt,
				})
			}
			if result == nil {
				result = []map[string]interface{}{}
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"tokens": result,
				"total":  total,
			})
		})

		// PATCH /auth/connection-tokens/:id — edit permissions without re-pairing
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}
			tokenID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid token ID")
				return
			}

			// Verify ownership
			existing, err := db.ConnectionTokenStore().FindByID(r.Context(), tokenID)
			if err != nil || existing == nil || existing.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "Token not found")
				return
			}

			var body struct {
				Name       string  `json:"name"`
				AllowedIPs *string `json:"allowedIps"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
				return
			}

			name := body.Name
			if name == "" {
				name = existing.Name
			}

			allowedIPs := existing.AllowedIPs
			if body.AllowedIPs != nil {
				normalized := strings.TrimSpace(*body.AllowedIPs)
				if msg := validateIPAllowlist(normalized); msg != "" {
					helpers.WriteError(w, http.StatusBadRequest, "Invalid allowedIps: "+msg)
					return
				}
				if normalized == "" {
					allowedIPs = sql.NullString{Valid: false}
				} else {
					allowedIPs = sql.NullString{String: normalized, Valid: true}
				}
			}

			if err := db.ConnectionTokenStore().UpdatePermissions(
				r.Context(), tokenID, name,
				allowedIPs,
				existing.AllowedTargetClients,
				existing.RemoteScopes,
				existing.RemoteRequestsPerHour,
			); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update token")
				return
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		})

		// DELETE /auth/connection-tokens/:id — Revoke a token AND force-disconnect any active client
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			tokenID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid token ID")
				return
			}

			if err := db.ConnectionTokenStore().Delete(r.Context(), tokenID); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to delete connection token")
				return
			}

			// Force-disconnect any active WebSocket clients that authenticated with this token.
			// Broadcast via Redis pub/sub so sessions pinned to other Fly instances are killed too.
			disconnected := 0
			if manager != nil {
				disconnected = manager.BroadcastDisconnectByConnectionToken(r.Context(), tokenID, "Connection token revoked")
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"success":             true,
				"message":             "Connection token revoked",
				"disconnectedClients": disconnected,
			})
		})
	})

	// POST /auth/pair — Exchange pairing code for connection token (NO AUTH, RATE LIMITED)
	r.With(middleware.PairingRateLimiter.Middleware).Post("/pair", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Code       string `json:"code"`
			WorldID    string `json:"worldId"`
			WorldTitle string `json:"worldTitle"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Code is required")
			return
		}

		// worldId is required (v3.0+) unless this is an "add browser" pairing
		// where the code was generated with a bound clientId. We validate this
		// after loading the pairing code below.

		ctx := r.Context()
		pairingCode, err := db.PairingCodeStore().FindByCode(ctx, body.Code)
		if err != nil || pairingCode == nil {
			helpers.WriteError(w, http.StatusNotFound, "Invalid pairing code")
			return
		}

		if pairingCode.Used {
			helpers.WriteError(w, http.StatusGone, "Pairing code has already been used")
			return
		}

		if pairingCode.ExpiresAt.Time.Before(time.Now()) {
			helpers.WriteError(w, http.StatusGone, "Pairing code has expired")
			return
		}

		// Enforce worldId requirement for first-pair flows (not "add browser").
		isAddBrowser := pairingCode.ClientID.Valid && pairingCode.ClientID.String != ""
		if !isAddBrowser && body.WorldID == "" {
			helpers.WriteError(w, http.StatusBadRequest, "worldId is required")
			return
		}

		// CRITICAL: Mark pairing code as used FIRST to prevent race conditions.
		// If two concurrent requests hit this endpoint with the same code,
		// only one will succeed at marking it used.
		if err := db.PairingCodeStore().MarkUsed(ctx, pairingCode.ID); err != nil {
			log.Error().Err(err).Msg("Failed to mark pairing code as used")
			helpers.WriteError(w, http.StatusConflict, "Pairing code could not be claimed. It may have been used by another request.")
			return
		}

		// Generate connection token (32 random bytes)
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}
		rawToken := hex.EncodeToString(tokenBytes)

		// Hash for storage
		hash := sha256.Sum256([]byte(rawToken))
		tokenHash := hex.EncodeToString(hash[:])

		// Determine clientId:
		//   - "add browser" flow: reuse the pairing code's bound clientId
		//   - First-pair/new-world: look up by (userId, worldId); reuse if found,
		//     otherwise generate a deterministic clientId from (userId, worldId)
		var clientID string
		if isAddBrowser {
			clientID = pairingCode.ClientID.String
		} else {
			existing, lookupErr := db.KnownClientStore().FindByWorldID(ctx, pairingCode.UserID, body.WorldID)
			if lookupErr != nil {
				log.Warn().Err(lookupErr).Msg("Failed to look up known client by worldId")
			}
			if existing != nil {
				clientID = existing.ClientID
			} else {
				clientID = deterministicClientID(pairingCode.UserID, body.WorldID)
			}
		}

		// Store connection token. Cross-world permissions are world-level (KnownClient)
		// not per-token, so we don't copy them from the pairing code here.
		connToken := &model.ConnectionToken{
			UserID:    pairingCode.UserID,
			TokenHash: tokenHash,
			Name:      parseUAName(r.Header.Get("User-Agent")),
			ClientID:  clientID,
			Source:    model.TokenSourceDashboard,
		}
		if err := db.ConnectionTokenStore().Create(ctx, connToken); err != nil {
			log.Error().Err(err).Msg("Failed to create connection token")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to create connection token")
			return
		}

		// Always upsert KnownClient with world metadata. For "add browser" the
		// row already exists; the upsert updates worldTitle if it changed.
		// For first-pair, this creates the row with worldId set from the start.
		knownClient := &model.KnownClient{
			UserID:   pairingCode.UserID,
			ClientID: clientID,
		}
		if !isAddBrowser {
			knownClient.WorldID = sql.NullString{String: body.WorldID, Valid: true}
			knownClient.WorldTitle = sql.NullString{String: body.WorldTitle, Valid: body.WorldTitle != ""}
		}
		if err := db.KnownClientStore().Upsert(ctx, knownClient); err != nil {
			log.Warn().Err(err).Msg("Failed to upsert known client entry")
		}

		// If the PairingCode carries cross-world settings (set by the pair-request
		// approval handler), apply them to the KnownClient now that it exists.
		if pairingCode.AllowedTargetClients.Valid || pairingCode.RemoteScopes.Valid {
			if kc, err := db.KnownClientStore().FindByClientID(ctx, pairingCode.UserID, clientID); err == nil && kc != nil {
				_ = db.KnownClientStore().SetCrossWorldSettings(ctx, kc.ID,
					pairingCode.AllowedTargetClients,
					pairingCode.RemoteScopes, 0)
			}
		}

		// Determine relay URL — derive from request host so it works on localhost,
		// production, and any custom domain.
		relayURL := buildRelayWSURL(r, cfg.FrontendURL)

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"token":    rawToken,
			"clientId": clientID,
			"relayUrl": relayURL,
		})

		log.Info().Int64("userId", pairingCode.UserID).Str("clientId", clientID).Str("worldId", body.WorldID).Int64("tokenId", connToken.ID).Str("tokenHashPrefix", tokenHash[:8]+"…").Msg("Pairing completed")
	})

	// GET /auth/remote-request-logs — View cross-world audit log (authenticated, master key only)
	r.Route("/remote-request-logs", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			limit := 50
			offset := 0
			if l := r.URL.Query().Get("limit"); l != "" {
				if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 200 {
					limit = v
				}
			}
			if o := r.URL.Query().Get("offset"); o != "" {
				if v, err := strconv.Atoi(o); err == nil && v >= 0 {
					offset = v
				}
			}

			logs, total, err := db.RemoteRequestLogStore().FindRecentByUser(r.Context(), user.ID, limit, offset)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to list remote request logs")
				return
			}

			result := make([]map[string]interface{}, 0, len(logs))
			for _, l := range logs {
				result = append(result, map[string]interface{}{
					"id":             l.ID,
					"sourceClientId": l.SourceClientID,
					"sourceTokenId":  l.SourceTokenID,
					"targetClientId": l.TargetClientID,
					"action":         l.Action,
					"success":        bool(l.Success),
					"errorMessage":   l.ErrorMessage.String,
					"sourceIp":       l.SourceIP.String,
					"createdAt":      l.CreatedAt,
				})
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"logs":   result,
				"limit":  limit,
				"offset": offset,
				"total":  total,
			})
		})
	})

	// GET /auth/connection-logs — View connection audit log (authenticated)
	r.Route("/connection-logs", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			limit := 50
			offset := 0
			if l := r.URL.Query().Get("limit"); l != "" {
				if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 200 {
					limit = v
				}
			}
			if o := r.URL.Query().Get("offset"); o != "" {
				if v, err := strconv.Atoi(o); err == nil && v >= 0 {
					offset = v
				}
			}

			logs, err := db.ConnectionLogStore().FindByUser(r.Context(), user.ID, limit, offset)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to list connection logs")
				return
			}

			result := make([]map[string]interface{}, 0, len(logs))
			for _, log := range logs {
				result = append(result, map[string]interface{}{
					"id":             log.ID,
					"clientId":       log.ClientID,
					"tokenName":      log.TokenName.String,
					"ipAddress":      log.IPAddress.String,
					"userAgent":      log.UserAgent.String,
					"worldId":        log.WorldID.String,
					"worldTitle":     log.WorldTitle.String,
					"systemId":       log.SystemID.String,
					"foundryVersion": log.FoundryVersion.String,
					"metadataMatch":  log.MetadataMatch,
					"flagged":        log.Flagged,
					"flagReason":     log.FlagReason.String,
					"createdAt":      log.CreatedAt,
				})
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"logs":   result,
				"limit":  limit,
				"offset": offset,
				"total":  len(result), // Note: this is the page size, not the absolute total. Pagination should use Has Next.
			})
		})
	})
}

// csvNullString turns a slice of strings into a comma-separated
// sql.NullString. Empty/nil slices return Valid=false so the column gets a
// real NULL instead of an empty string (matters for the empty-vs-default
// distinction in some queries).
func csvNullString(values []string) sql.NullString {
	if len(values) == 0 {
		return sql.NullString{}
	}
	clean := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			clean = append(clean, v)
		}
	}
	if len(clean) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.Join(clean, ","), Valid: true}
}

// generatePairingCode generates a random uppercase alphanumeric code of the given length.
func generatePairingCode(length int) (string, error) {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Removed ambiguous chars: 0OI1
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		code[i] = chars[n.Int64()]
	}
	return string(code), nil
}

// deterministicClientID generates a stable client ID derived from the user ID
// and world ID. The same (userId, worldId) pair always produces the same clientId,
// which prevents duplicate KnownClient rows across re-pairings.
func deterministicClientID(userID int64, worldID string) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d:%s", userID, worldID)))
	return "fvtt_" + hex.EncodeToString(h[:])[:16]
}

// parseUAName extracts a human-readable "Browser on OS" string from a User-Agent
// header value. Used as the default connection token name so users can identify
// which browser/device a pairing belongs to without renaming it.
func parseUAName(ua string) string {
	lower := strings.ToLower(ua)
	// Browser detection — order matters: Edge before Chrome, OPR before Chrome
	browser := ""
	switch {
	case strings.Contains(lower, "edg/") || strings.Contains(lower, "edge/"):
		browser = "Edge"
	case strings.Contains(lower, "opr/") || strings.Contains(lower, "opera"):
		browser = "Opera"
	case strings.Contains(lower, "firefox"):
		browser = "Firefox"
	case strings.Contains(lower, "chrome"):
		browser = "Chrome"
	case strings.Contains(lower, "safari"):
		browser = "Safari"
	}
	// OS detection
	os := ""
	switch {
	case strings.Contains(lower, "windows"):
		os = "Windows"
	case strings.Contains(lower, "android"):
		os = "Android"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad"):
		os = "iOS"
	case strings.Contains(lower, "mac os x") || strings.Contains(lower, "macos"):
		os = "macOS"
	case strings.Contains(lower, "linux"):
		os = "Linux"
	}
	switch {
	case browser != "" && os != "":
		return browser + " on " + os
	case browser != "":
		return browser
	case os != "":
		return "Browser on " + os
	default:
		return fmt.Sprintf("Paired %s", time.Now().Format("2006-01-02 15:04"))
	}
}

// buildRelayWSURL constructs the WebSocket URL for the /relay endpoint
// based on the actual incoming request.
//
// We always use r.Host (the host the client just connected to). This is
// correct for all environments:
//   - Local dev: r.Host = "localhost:3010" → ws://localhost:3010/relay
//   - Fly.io production: r.Host = "foundryrestapi.com" with
//     X-Forwarded-Proto: https → wss://foundryrestapi.com/relay
//   - Custom deployments: works automatically
//
// We do NOT use cfg.FrontendURL because that's the dashboard URL and may
// differ from the relay URL (e.g., dashboard at foundryrestapi.com but
// relay at relay.foundryrestapi.com). It also defaults to the production
// frontend URL even in local dev, which would produce a wrong URL.
//
// The frontendURL parameter is unused but kept in the signature for callers
// that may want to override in the future.
func buildRelayWSURL(r *http.Request, frontendURL string) string {
	_ = frontendURL // intentionally unused — see doc comment

	host := r.Host
	if host == "" {
		// Should never happen for a real HTTP request, but be defensive
		host = "localhost:3010"
	}

	// Determine the scheme: TLS at the relay, or X-Forwarded-Proto: https
	// from a reverse proxy (Fly, Cloudflare, etc.), means wss://
	scheme := "ws"
	if r.TLS != nil {
		scheme = "wss"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
		scheme = "wss"
	}

	return scheme + "://" + host + "/relay"
}

// RegisterCredentialRoutes registers routes for the credential vault,
// known clients, notification settings, and per-key notification settings.
//
// The ClientManager is required so DELETE /auth/known-clients/:id can
// force-disconnect any active WebSocket connection for that client.
func RegisterCredentialRoutes(r chi.Router, db *database.DB, cfg *config.Config, manager *ws.ClientManager) {
	r.Route("/credentials", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		// POST /auth/credentials — Create credential set
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			var body struct {
				Name            string `json:"name"`
				FoundryURL      string `json:"foundryUrl"`
				FoundryUsername string `json:"foundryUsername"`
				FoundryPassword string `json:"foundryPassword"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}
			if body.Name == "" || body.FoundryURL == "" || body.FoundryUsername == "" || body.FoundryPassword == "" {
				helpers.WriteError(w, http.StatusBadRequest, "name, foundryUrl, foundryUsername, and foundryPassword are all required")
				return
			}

			if !service.IsEncryptionAvailable(cfg.CredentialsEncryptionKey) {
				helpers.WriteError(w, http.StatusBadRequest, "Credential storage is not available. CREDENTIALS_ENCRYPTION_KEY is not configured.")
				return
			}

			encrypted, err := service.Encrypt(body.FoundryPassword, cfg.CredentialsEncryptionKey)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to encrypt credentials")
				return
			}

			cred := &model.Credential{
				UserID:                   user.ID,
				Name:                     body.Name,
				FoundryURL:               body.FoundryURL,
				FoundryUsername:          body.FoundryUsername,
				EncryptedFoundryPassword: encrypted.Ciphertext,
				PasswordIV:               encrypted.IV,
				PasswordAuthTag:          encrypted.AuthTag,
			}

			if err := db.CredentialStore().Create(r.Context(), cred); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to create credential")
				return
			}

			helpers.WriteJSON(w, http.StatusCreated, map[string]interface{}{
				"id":              cred.ID,
				"name":            cred.Name,
				"foundryUrl":      cred.FoundryURL,
				"foundryUsername": cred.FoundryUsername,
				"createdAt":       cred.CreatedAt,
			})
		})

		// GET /auth/credentials — List credentials (never returns passwords)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			creds, err := db.CredentialStore().FindAllByUser(r.Context(), user.ID)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to list credentials")
				return
			}

			var result []map[string]interface{}
			for _, c := range creds {
				result = append(result, map[string]interface{}{
					"id":              c.ID,
					"name":            c.Name,
					"foundryUrl":      c.FoundryURL,
					"foundryUsername": c.FoundryUsername,
					"createdAt":       c.CreatedAt,
					"updatedAt":       c.UpdatedAt,
				})
			}
			if result == nil {
				result = []map[string]interface{}{}
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"credentials": result})
		})

		// PATCH /auth/credentials/:id — Update credential
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			credID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid credential ID")
				return
			}

			ctx := r.Context()
			cred, _ := db.CredentialStore().FindByID(ctx, credID)
			if cred == nil || cred.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "Credential not found")
				return
			}

			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)

			if name, ok := body["name"].(string); ok && name != "" {
				cred.Name = name
			}
			if url, ok := body["foundryUrl"].(string); ok {
				cred.FoundryURL = url
			}
			if username, ok := body["foundryUsername"].(string); ok {
				cred.FoundryUsername = username
			}
			pw, hasPw := body["foundryPassword"].(string)
			if !hasPw || pw == "" {
				helpers.WriteError(w, http.StatusBadRequest, "Password is required when updating a credential.")
				return
			}
			if !service.IsEncryptionAvailable(cfg.CredentialsEncryptionKey) {
				helpers.WriteError(w, http.StatusBadRequest, "Credential storage is not available.")
				return
			}
			encrypted, err := service.Encrypt(pw, cfg.CredentialsEncryptionKey)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to encrypt credentials")
				return
			}
			cred.EncryptedFoundryPassword = encrypted.Ciphertext
			cred.PasswordIV = encrypted.IV
			cred.PasswordAuthTag = encrypted.AuthTag

			if err := db.CredentialStore().Update(ctx, cred); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update credential")
				return
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"id":              cred.ID,
				"name":            cred.Name,
				"foundryUrl":      cred.FoundryURL,
				"foundryUsername": cred.FoundryUsername,
				"updatedAt":       cred.UpdatedAt,
			})
		})

		// DELETE /auth/credentials/:id
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			credID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid credential ID")
				return
			}

			ctx := r.Context()
			cred, _ := db.CredentialStore().FindByID(ctx, credID)
			if cred == nil || cred.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "Credential not found")
				return
			}

			if err := db.CredentialStore().Delete(ctx, credID); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to delete credential")
				return
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "Credential deleted",
			})
		})
	})

	// Known Clients routes
	r.Route("/known-clients", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		// GET /auth/known-clients — List all known clients with online/offline status,
		// their connection tokens, and which token is currently active.
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			clients, err := db.KnownClientStore().FindAllByUser(r.Context(), user.ID)
			if err != nil {
				log.Error().Err(err).Int64("userId", user.ID).Msg("FindAllByUser(KnownClients) failed")
				helpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list known clients: %s", err))
				return
			}

			// Fetch all tokens for this user, then group by clientId.
			allTokens, err := db.ConnectionTokenStore().FindAllByUser(r.Context(), user.ID)
			if err != nil {
				log.Error().Err(err).Int64("userId", user.ID).Msg("FindAllByUser(ConnectionTokens) failed")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to list connection tokens")
				return
			}
			tokensByClient := make(map[string][]*model.ConnectionToken, len(clients))
			for _, t := range allTokens {
				tokensByClient[t.ClientID] = append(tokensByClient[t.ClientID], t)
			}

			result := make([]map[string]interface{}, 0, len(clients))
			for _, c := range clients {
				activeTokenID := manager.LookupClientConnectionTokenID(r.Context(), c.ClientID)

				// Build lightweight token list (no tokenHash, allowedIps as plain string).
				tokens := tokensByClient[c.ClientID]
				if tokens == nil {
					tokens = []*model.ConnectionToken{}
				}
				// Legacy: tokens created before the clientId bug fix have clientId="".
				// If the active token isn't already in the list, check the "" bucket.
				if activeTokenID != 0 {
					found := false
					for _, t := range tokens {
						if t.ID == activeTokenID {
							found = true
							break
						}
					}
					if !found {
						for _, t := range tokensByClient[""] {
							if t.ID == activeTokenID {
								tokens = append(tokens, t)
								break
							}
						}
					}
				}

				// Build lightweight token views with plain-string fields.
				tokenViews := make([]map[string]interface{}, 0, len(tokens))
				for _, t := range tokens {
					tokenViews = append(tokenViews, map[string]interface{}{
						"id":         t.ID,
						"name":       t.Name,
						"clientId":   t.ClientID,
						"source":     t.Source,
						"allowedIps": t.AllowedIPs.String,
						"lastUsedAt": t.LastUsedAt,
						"createdAt":  t.CreatedAt,
					})
				}

				row := map[string]interface{}{
					"id":                       c.ID,
					"clientId":                 c.ClientID,
					"worldId":                  c.WorldID.String,
					"worldTitle":               c.WorldTitle.String,
					"systemId":                 c.SystemID.String,
					"systemTitle":              c.SystemTitle.String,
					"systemVersion":            c.SystemVersion.String,
					"foundryVersion":           c.FoundryVersion.String,
					"customName":               c.CustomName.String,
					"lastSeenAt":               c.LastSeenAt,
					"isOnline":                 manager.IsClientOnlineAnywhere(r.Context(), c.ClientID),
					"autoStartOnRemoteRequest": bool(c.AutoStartOnRemoteRequest),
					"allowedTargetClients":     c.GetAllowedTargets(),
					"remoteScopes":             c.GetRemoteScopes(),
					"remoteRequestsPerHour":    c.RemoteRequestsPerHour,
					"createdAt":                c.CreatedAt,
					"updatedAt":                c.UpdatedAt,
					"tokens":                   tokenViews,
					"activeTokenId":            activeTokenID,
				}
				if c.CredentialID.Valid {
					row["credentialId"] = c.CredentialID.Int64
				} else {
					row["credentialId"] = nil
				}
				result = append(result, row)
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"clients": result})
		})

		// PATCH /auth/known-clients/:id/credential — set the explicit
		// Credential link used by AutoStartForKnownClient. The dashboard
		// "Set credential" dropdown calls this. Required when the user has
		// multiple credentials and wants auto-start enabled — without this
		// link the auto-start path falls back to "the user's first credential"
		// which only works for single-Foundry-server deployments.
		//
		// Body: { "credentialId": <number|null> }
		r.Patch("/{id}/credential", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}
			rowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid client ID")
				return
			}

			known, err := db.KnownClientStore().FindByID(r.Context(), rowID)
			if err != nil || known == nil || known.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "Known client not found")
				return
			}

			var body struct {
				CredentialID *int64 `json:"credentialId"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
				return
			}

			var credID sql.NullInt64
			if body.CredentialID != nil {
				// Validate ownership
				cred, err := db.CredentialStore().FindByID(r.Context(), *body.CredentialID)
				if err != nil || cred == nil || cred.UserID != user.ID {
					helpers.WriteError(w, http.StatusNotFound, "Credential not found")
					return
				}
				credID = sql.NullInt64{Int64: *body.CredentialID, Valid: true}
			}

			if err := db.KnownClientStore().SetCredentialID(r.Context(), rowID, credID); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to set credential")
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		})

		// PATCH /auth/known-clients/:id — Update cross-world tunneling settings for a world.
		// Body: { allowedTargetClients?: string[], remoteScopes?: string[], remoteRequestsPerHour?: number }
		// These settings apply to all browsers (connection tokens) paired to this world.
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}
			rowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid client ID")
				return
			}

			// Verify ownership
			known, err := db.KnownClientStore().FindByID(r.Context(), rowID)
			if err != nil || known == nil || known.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "Known client not found")
				return
			}

			var body struct {
				AllowedTargetClients  []string `json:"allowedTargetClients"`
				RemoteScopes          []string `json:"remoteScopes"`
				RemoteRequestsPerHour *int     `json:"remoteRequestsPerHour"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
				return
			}

			// Validate allowedTargetClients all belong to this user.
			for _, tcid := range body.AllowedTargetClients {
				kc, err := db.KnownClientStore().FindByClientID(r.Context(), user.ID, tcid)
				if err != nil || kc == nil {
					helpers.WriteError(w, http.StatusBadRequest, fmt.Sprintf("allowedTargetClient %q not found in your known clients", tcid))
					return
				}
			}

			// Validate remoteScopes are known scope strings.
			for _, sc := range body.RemoteScopes {
				if !model.IsValidScope(sc) {
					helpers.WriteError(w, http.StatusBadRequest, fmt.Sprintf("unknown scope %q", sc))
					return
				}
			}

			// Only overwrite fields that were explicitly provided.
			allowedTargets := known.AllowedTargetClients
			if body.AllowedTargetClients != nil {
				allowedTargets = csvNullString(body.AllowedTargetClients)
			}
			remoteScopes := known.RemoteScopes
			if body.RemoteScopes != nil {
				remoteScopes = csvNullString(body.RemoteScopes)
			}
			rateLimit := known.RemoteRequestsPerHour
			if body.RemoteRequestsPerHour != nil {
				rateLimit = *body.RemoteRequestsPerHour
			}

			if err := db.KnownClientStore().SetCrossWorldSettings(r.Context(), rowID, allowedTargets, remoteScopes, rateLimit); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update cross-world settings")
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		})

		// PATCH /auth/known-clients/:id/auto-start — toggle whether the relay
		// will spawn a headless session for this client in response to an
		// incoming remote-request from a sibling client (when this client is
		// currently offline).
		r.Patch("/{id}/auto-start", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}
			rowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid client ID")
				return
			}

			// Verify ownership
			known, err := db.KnownClientStore().FindByID(r.Context(), rowID)
			if err != nil || known == nil || known.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "Known client not found")
				return
			}

			var body struct {
				Enabled bool `json:"enabled"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
				return
			}

			if err := db.KnownClientStore().SetAutoStartOnRemoteRequest(r.Context(), rowID, body.Enabled); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update auto-start setting")
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		})

		// DELETE /auth/known-clients/:id
		// Force-disconnects any live WebSocket connection for this client, then
		// deletes the row from the KnownClients table.
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			rowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid client ID")
				return
			}

			// Look up the known client row to get its opaque clientId string
			// (needed to find the live WS connection) and verify ownership.
			ctx := r.Context()
			known, err := db.KnownClientStore().FindByID(ctx, rowID)
			if err != nil {
				log.Error().Err(err).Int64("rowId", rowID).Msg("FindByID(KnownClient) failed")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to look up known client")
				return
			}
			if known == nil || known.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "Known client not found")
				return
			}

			// Force-disconnect any active WebSocket connection for this clientId,
			// then delete ALL connection tokens associated with this world.
			disconnected := false
			if manager != nil {
				tokenID := manager.LookupClientConnectionTokenID(ctx, known.ClientID)
				if tokenID != 0 {
					n := manager.BroadcastDisconnectByConnectionToken(ctx, tokenID, "Client removed by owner")
					disconnected = n > 0
				} else {
					n := manager.BroadcastDisconnectByClientID(ctx, known.ClientID, "Client removed by owner")
					disconnected = n > 0
				}
			}

			// Delete all tokens for this world (not just the active one).
			if _, err := db.ConnectionTokenStore().DeleteAllByClientID(ctx, user.ID, known.ClientID); err != nil {
				log.Warn().Err(err).Str("clientId", known.ClientID).Msg("Failed to delete connection tokens on known-client delete")
			}

			// Delete the DB row
			if err := db.KnownClientStore().Delete(ctx, rowID); err != nil {
				log.Error().Err(err).Int64("rowId", rowID).Msg("Delete(KnownClient) failed")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to delete known client")
				return
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"success":      true,
				"message":      "Known client removed",
				"disconnected": disconnected,
			})
		})
	})

	// Notification Settings routes
	r.Route("/notification-settings", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		// GET /auth/notification-settings
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			settings, err := db.NotificationSettingsStore().FindByUser(r.Context(), user.ID)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to get notification settings")
				return
			}

			if settings == nil {
				// Return defaults
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"discordWebhookUrl":              "",
					"notifyEmail":                    "",
					"notifyOnConnect":                true,
					"notifyOnDisconnect":             true,
					"notifyOnMetadataMismatch":       true,
					"notifyOnSettingsChange":         true,
					"notifyOnExecuteJs":              true,
					"notifyOnMacroExecute":           false,
					"notifyOnNewClientConnect":       true,
					"notificationDebounceWindowSecs": 0,
					"remoteRequestBatchWindowSecs":   300,
					"logCrossWorldRequests":          true,
					"smtpAvailable":                  cfg.SMTPHost != "",
				})
				return
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"discordWebhookUrl":              settings.DiscordWebhookURL.String,
				"notifyEmail":                    settings.NotifyEmail.String,
				"notifyOnConnect":                settings.NotifyOnConnect,
				"notifyOnDisconnect":             settings.NotifyOnDisconnect,
				"notifyOnMetadataMismatch":       settings.NotifyOnMetadataMismatch,
				"notifyOnSettingsChange":         settings.NotifyOnSettingsChange,
				"notifyOnExecuteJs":              settings.NotifyOnExecuteJs,
				"notifyOnMacroExecute":           settings.NotifyOnMacroExecute,
				"notifyOnNewClientConnect":       settings.NotifyOnNewClientConnect,
				"notificationDebounceWindowSecs": settings.NotificationDebounceWindowSecs,
				"remoteRequestBatchWindowSecs":   settings.RemoteRequestBatchWindowSecs,
				"logCrossWorldRequests":          settings.LogCrossWorldRequests,
				"smtpAvailable":                  cfg.SMTPHost != "",
			})
		})

		// PUT /auth/notification-settings
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			var body struct {
				DiscordWebhookURL              string `json:"discordWebhookUrl"`
				NotifyEmail                    string `json:"notifyEmail"`
				NotifyOnConnect                *bool  `json:"notifyOnConnect"`
				NotifyOnDisconnect             *bool  `json:"notifyOnDisconnect"`
				NotifyOnMetadataMismatch       *bool  `json:"notifyOnMetadataMismatch"`
				NotifyOnSettingsChange         *bool  `json:"notifyOnSettingsChange"`
				NotifyOnExecuteJs              *bool  `json:"notifyOnExecuteJs"`
				NotifyOnMacroExecute           *bool  `json:"notifyOnMacroExecute"`
				NotifyOnNewClientConnect       *bool  `json:"notifyOnNewClientConnect"`
				NotificationDebounceWindowSecs *int   `json:"notificationDebounceWindowSecs"`
				RemoteRequestBatchWindowSecs   *int   `json:"remoteRequestBatchWindowSecs"`
				LogCrossWorldRequests          *bool  `json:"logCrossWorldRequests"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}

			boolPtrOr := func(p *bool, def bool) bool {
				if p == nil {
					return def
				}
				return *p
			}
			intPtrOr := func(p *int, def int) int {
				if p == nil {
					return def
				}
				return *p
			}

			settings := &model.NotificationSettings{
				UserID:                         user.ID,
				DiscordWebhookURL:              sql.NullString{String: body.DiscordWebhookURL, Valid: body.DiscordWebhookURL != ""},
				NotifyEmail:                    sql.NullString{String: body.NotifyEmail, Valid: body.NotifyEmail != ""},
				NotifyOnConnect:                boolPtrOr(body.NotifyOnConnect, true),
				NotifyOnDisconnect:             boolPtrOr(body.NotifyOnDisconnect, true),
				NotifyOnMetadataMismatch:       boolPtrOr(body.NotifyOnMetadataMismatch, true),
				NotifyOnSettingsChange:         boolPtrOr(body.NotifyOnSettingsChange, true),
				NotifyOnExecuteJs:              boolPtrOr(body.NotifyOnExecuteJs, true),
				NotifyOnMacroExecute:           boolPtrOr(body.NotifyOnMacroExecute, false),
				NotifyOnNewClientConnect:       boolPtrOr(body.NotifyOnNewClientConnect, true),
				NotificationDebounceWindowSecs: intPtrOr(body.NotificationDebounceWindowSecs, 0),
				RemoteRequestBatchWindowSecs:   intPtrOr(body.RemoteRequestBatchWindowSecs, 300),
				LogCrossWorldRequests:          boolPtrOr(body.LogCrossWorldRequests, true),
			}

			if err := db.NotificationSettingsStore().Upsert(r.Context(), settings); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update notification settings")
				return
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"discordWebhookUrl":              settings.DiscordWebhookURL.String,
				"notifyEmail":                    settings.NotifyEmail.String,
				"notifyOnConnect":                settings.NotifyOnConnect,
				"notifyOnDisconnect":             settings.NotifyOnDisconnect,
				"notifyOnMetadataMismatch":       settings.NotifyOnMetadataMismatch,
				"notifyOnSettingsChange":         settings.NotifyOnSettingsChange,
				"notifyOnExecuteJs":              settings.NotifyOnExecuteJs,
				"notifyOnMacroExecute":           settings.NotifyOnMacroExecute,
				"notifyOnNewClientConnect":       settings.NotifyOnNewClientConnect,
				"notificationDebounceWindowSecs": settings.NotificationDebounceWindowSecs,
				"remoteRequestBatchWindowSecs":   settings.RemoteRequestBatchWindowSecs,
			})
		})

		// POST /auth/notification-settings/test — Send test notification.
		// Accepts an optional JSON body { discordWebhookUrl, notifyEmail } to
		// test with unsaved form values. Falls back to saved settings if the
		// body fields are empty or the body is absent.
		r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			// Decode optional body — non-empty values override saved settings.
			var body struct {
				DiscordWebhookURL string `json:"discordWebhookUrl"`
				NotifyEmail       string `json:"notifyEmail"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body) // ignore parse errors — body is optional

			webhookURL := body.DiscordWebhookURL
			notifyEmail := body.NotifyEmail

			// Fall back to saved settings for any field not supplied in the body.
			if webhookURL == "" || notifyEmail == "" {
				settings, _ := db.NotificationSettingsStore().FindByUser(r.Context(), user.ID)
				if settings != nil {
					if webhookURL == "" {
						webhookURL = settings.DiscordWebhookURL.String
					}
					if notifyEmail == "" {
						notifyEmail = settings.NotifyEmail.String
					}
				}
			}

			if webhookURL == "" && notifyEmail == "" {
				helpers.WriteError(w, http.StatusBadRequest, "No notification destination configured")
				return
			}

			notifCfg := &service.NotificationConfig{
				SMTPHost:    cfg.SMTPHost,
				SMTPPort:    cfg.SMTPPort,
				SMTPUser:    cfg.SMTPUser,
				SMTPPass:    cfg.SMTPPass,
				SMTPFrom:    cfg.SMTPFrom,
				FrontendURL: cfg.FrontendURL,
			}

			if err := service.SendTestNotification(webhookURL, notifyEmail, notifCfg); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Test notification failed: %s", err))
				return
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "Test notification sent",
			})
		})
	})

	// Per-scoped-key notification settings
	// GET    /auth/api-keys/:id/notification-settings
	// PUT    /auth/api-keys/:id/notification-settings
	// DELETE /auth/api-keys/:id/notification-settings
	// POST   /auth/api-keys/:id/notification-settings/test
	r.Route("/api-keys/{id}/notification-settings", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		// Helper to look up the key + verify ownership
		lookupKey := func(w http.ResponseWriter, r *http.Request) (*model.User, *model.ApiKey, bool) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return nil, nil, false
			}
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return nil, nil, false
			}
			keyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid key ID")
				return nil, nil, false
			}
			key, _ := db.ApiKeyStore().FindByID(r.Context(), keyID)
			if key == nil || key.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "API key not found")
				return nil, nil, false
			}
			return user, key, true
		}

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, key, ok := lookupKey(w, r)
			if !ok {
				return
			}
			settings, err := db.ApiKeyNotificationSettingsStore().FindByApiKey(r.Context(), key.ID)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to get key notification settings")
				return
			}
			if settings == nil {
				// Default empty config
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"apiKeyId":             key.ID,
					"discordWebhookUrl":    "",
					"notifyEmail":          "",
					"notifyOnExecuteJs":    false,
					"notifyOnMacroExecute": false,
					"notifyOnRateLimit":    false,
					"notifyOnError":        false,
					"smtpAvailable":        cfg.SMTPHost != "",
				})
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"apiKeyId":             settings.ApiKeyID,
				"discordWebhookUrl":    settings.DiscordWebhookURL.String,
				"notifyEmail":          settings.NotifyEmail.String,
				"notifyOnExecuteJs":    settings.NotifyOnExecuteJs,
				"notifyOnMacroExecute": settings.NotifyOnMacroExecute,
				"notifyOnRateLimit":    settings.NotifyOnRateLimit,
				"notifyOnError":        settings.NotifyOnError,
				"smtpAvailable":        cfg.SMTPHost != "",
			})
		})

		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			_, key, ok := lookupKey(w, r)
			if !ok {
				return
			}
			var body struct {
				DiscordWebhookURL    string `json:"discordWebhookUrl"`
				NotifyEmail          string `json:"notifyEmail"`
				NotifyOnExecuteJs    *bool  `json:"notifyOnExecuteJs"`
				NotifyOnMacroExecute *bool  `json:"notifyOnMacroExecute"`
				NotifyOnRateLimit    *bool  `json:"notifyOnRateLimit"`
				NotifyOnError        *bool  `json:"notifyOnError"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}
			boolPtrOr := func(p *bool, def bool) bool {
				if p == nil {
					return def
				}
				return *p
			}
			settings := &model.ApiKeyNotificationSettings{
				ApiKeyID:             key.ID,
				DiscordWebhookURL:    sql.NullString{String: body.DiscordWebhookURL, Valid: body.DiscordWebhookURL != ""},
				NotifyEmail:          sql.NullString{String: body.NotifyEmail, Valid: body.NotifyEmail != ""},
				NotifyOnExecuteJs:    boolPtrOr(body.NotifyOnExecuteJs, false),
				NotifyOnMacroExecute: boolPtrOr(body.NotifyOnMacroExecute, false),
				NotifyOnRateLimit:    boolPtrOr(body.NotifyOnRateLimit, false),
				NotifyOnError:        boolPtrOr(body.NotifyOnError, false),
			}
			if err := db.ApiKeyNotificationSettingsStore().Upsert(r.Context(), settings); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update key notification settings")
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"apiKeyId":             settings.ApiKeyID,
				"discordWebhookUrl":    settings.DiscordWebhookURL.String,
				"notifyEmail":          settings.NotifyEmail.String,
				"notifyOnExecuteJs":    settings.NotifyOnExecuteJs,
				"notifyOnMacroExecute": settings.NotifyOnMacroExecute,
				"notifyOnRateLimit":    settings.NotifyOnRateLimit,
				"notifyOnError":        settings.NotifyOnError,
			})
		})

		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			_, key, ok := lookupKey(w, r)
			if !ok {
				return
			}
			if err := db.ApiKeyNotificationSettingsStore().Delete(r.Context(), key.ID); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to delete key notification settings")
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		})

		r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
			_, key, ok := lookupKey(w, r)
			if !ok {
				return
			}
			settings, err := db.ApiKeyNotificationSettingsStore().FindByApiKey(r.Context(), key.ID)
			if err != nil || settings == nil {
				helpers.WriteError(w, http.StatusBadRequest, "No notification settings configured for this key")
				return
			}
			notifCfg := &service.NotificationConfig{
				SMTPHost:    cfg.SMTPHost,
				SMTPPort:    cfg.SMTPPort,
				SMTPUser:    cfg.SMTPUser,
				SMTPPass:    cfg.SMTPPass,
				SMTPFrom:    cfg.SMTPFrom,
				FrontendURL: cfg.FrontendURL,
			}
			if err := service.SendTestNotification(
				settings.DiscordWebhookURL.String,
				settings.NotifyEmail.String,
				notifCfg,
			); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Test notification failed: %s", err))
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "Test notification sent",
			})
		})
	})

	// GET /auth/known-clients/{id}/users — stored Foundry users for this world
	r.Route("/known-clients/{id}/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return
			}
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}
			rowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid client ID")
				return
			}
			known, err := db.KnownClientStore().FindByID(r.Context(), rowID)
			if err != nil || known == nil || known.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "Known client not found")
				return
			}
			users, err := db.KnownUserStore().FindAllByKnownClient(r.Context(), rowID)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to list known users")
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"users": users})
		})
	})

	// Per-world notification settings
	// GET    /auth/known-clients/{id}/notification-settings
	// PUT    /auth/known-clients/{id}/notification-settings
	// DELETE /auth/known-clients/{id}/notification-settings
	// POST   /auth/known-clients/{id}/notification-settings/test
	r.Route("/known-clients/{id}/notification-settings", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		lookupWorld := func(w http.ResponseWriter, r *http.Request) (*model.User, *model.KnownClient, bool) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
				return nil, nil, false
			}
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return nil, nil, false
			}
			rowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid known client ID")
				return nil, nil, false
			}
			kc, _ := db.KnownClientStore().FindByID(r.Context(), rowID)
			if kc == nil || kc.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "World not found")
				return nil, nil, false
			}
			return user, kc, true
		}

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, kc, ok := lookupWorld(w, r)
			if !ok {
				return
			}
			settings, err := db.KnownClientNotificationSettingsStore().FindByKnownClient(r.Context(), kc.ID)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to get world notification settings")
				return
			}
			if settings == nil {
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"knownClientId":        kc.ID,
					"discordWebhookUrl":    "",
					"notifyEmail":          "",
					"notifyOnConnect":      false,
					"notifyOnDisconnect":   false,
					"notifyOnExecuteJs":    false,
					"notifyOnMacroExecute": false,
					"smtpAvailable":        cfg.SMTPHost != "",
				})
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"knownClientId":        settings.KnownClientID,
				"discordWebhookUrl":    settings.DiscordWebhookURL.String,
				"notifyEmail":          settings.NotifyEmail.String,
				"notifyOnConnect":      settings.NotifyOnConnect,
				"notifyOnDisconnect":   settings.NotifyOnDisconnect,
				"notifyOnExecuteJs":    settings.NotifyOnExecuteJs,
				"notifyOnMacroExecute": settings.NotifyOnMacroExecute,
				"smtpAvailable":        cfg.SMTPHost != "",
			})
		})

		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			user, kc, ok := lookupWorld(w, r)
			if !ok {
				return
			}
			var body struct {
				DiscordWebhookURL    string `json:"discordWebhookUrl"`
				NotifyEmail          string `json:"notifyEmail"`
				NotifyOnConnect      *bool  `json:"notifyOnConnect"`
				NotifyOnDisconnect   *bool  `json:"notifyOnDisconnect"`
				NotifyOnExecuteJs    *bool  `json:"notifyOnExecuteJs"`
				NotifyOnMacroExecute *bool  `json:"notifyOnMacroExecute"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}
			boolPtrOr := func(p *bool, def bool) bool {
				if p == nil {
					return def
				}
				return *p
			}
			settings := &model.KnownClientNotificationSettings{
				KnownClientID:        kc.ID,
				UserID:               user.ID,
				DiscordWebhookURL:    sql.NullString{String: body.DiscordWebhookURL, Valid: body.DiscordWebhookURL != ""},
				NotifyEmail:          sql.NullString{String: body.NotifyEmail, Valid: body.NotifyEmail != ""},
				NotifyOnConnect:      boolPtrOr(body.NotifyOnConnect, false),
				NotifyOnDisconnect:   boolPtrOr(body.NotifyOnDisconnect, false),
				NotifyOnExecuteJs:    boolPtrOr(body.NotifyOnExecuteJs, false),
				NotifyOnMacroExecute: boolPtrOr(body.NotifyOnMacroExecute, false),
			}
			if err := db.KnownClientNotificationSettingsStore().Upsert(r.Context(), settings); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update world notification settings")
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"knownClientId":        settings.KnownClientID,
				"discordWebhookUrl":    settings.DiscordWebhookURL.String,
				"notifyEmail":          settings.NotifyEmail.String,
				"notifyOnConnect":      settings.NotifyOnConnect,
				"notifyOnDisconnect":   settings.NotifyOnDisconnect,
				"notifyOnExecuteJs":    settings.NotifyOnExecuteJs,
				"notifyOnMacroExecute": settings.NotifyOnMacroExecute,
				"smtpAvailable":        cfg.SMTPHost != "",
			})
		})

		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			_, kc, ok := lookupWorld(w, r)
			if !ok {
				return
			}
			if err := db.KnownClientNotificationSettingsStore().Delete(r.Context(), kc.ID); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to delete world notification settings")
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		})

		r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
			_, kc, ok := lookupWorld(w, r)
			if !ok {
				return
			}
			settings, err := db.KnownClientNotificationSettingsStore().FindByKnownClient(r.Context(), kc.ID)
			if err != nil || settings == nil {
				helpers.WriteError(w, http.StatusBadRequest, "No notification settings configured for this world")
				return
			}
			notifCfg := &service.NotificationConfig{
				SMTPHost:    cfg.SMTPHost,
				SMTPPort:    cfg.SMTPPort,
				SMTPUser:    cfg.SMTPUser,
				SMTPPass:    cfg.SMTPPass,
				SMTPFrom:    cfg.SMTPFrom,
				FrontendURL: cfg.FrontendURL,
			}
			if err := service.SendTestNotification(
				settings.DiscordWebhookURL.String,
				settings.NotifyEmail.String,
				notifCfg,
			); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Test notification failed: %s", err))
				return
			}
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "Test notification sent",
			})
		})
	})
}

// validateIPAllowlist validates a comma-separated list of IP addresses and/or CIDR ranges.
// Returns an error description if any entry is invalid, empty string if all are valid.
func validateIPAllowlist(raw string) string {
	if raw == "" {
		return ""
	}
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if net.ParseIP(entry) != nil {
			continue
		}
		if _, _, err := net.ParseCIDR(entry); err == nil {
			continue
		}
		return fmt.Sprintf("invalid IP or CIDR: %q", entry)
	}
	return ""
}
