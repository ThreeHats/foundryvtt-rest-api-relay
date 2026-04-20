package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/metrics"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/rs/zerolog/log"
)

// APIRouteConfig defines how to create a standardized API route handler.
type APIRouteConfig struct {
	Type           string
	RequiredParams []ParamDef
	OptionalParams []ParamDef
	Timeout        time.Duration // Default: 10s

	// ValidateParams performs custom validation. Return error message or empty string.
	ValidateParams func(params Params, r *http.Request) (map[string]interface{}, bool)

	// BuildPayload customizes the WS payload. If nil, params are used directly.
	BuildPayload func(params Params, r *http.Request) map[string]interface{}

	// BuildPendingRequest adds extra fields to the pending request.
	BuildPendingRequest func(params Params) *ws.PendingRequest

	// AfterForward is called in a goroutine after the WS message is successfully
	// sent to Foundry. Use for side-effects like notification dispatch or event
	// logging. reqCtx and params are not mutated after this point and are safe to read.
	AfterForward func(reqCtx *RequestContext, clientID string, params Params)
}

// RequestContext holds auth context set by middleware.
type RequestContext struct {
	User           interface{} // *model.User
	MasterAPIKey   string
	ScopedKey      *ScopedKeyInfo
	SubscriptionStatus string
}

// ScopedKeyInfo holds scoped key details from auth middleware.
type ScopedKeyInfo struct {
	ID              int64
	KeyPrefix       string            // first 8 chars of the raw key, safe to log
	ScopedClientID  string            // single allowed client; use ScopedClientIDs for multiple
	ScopedClientIDs []string          // multi-client scoping
	ScopedUserID    string            // userId applied to all clients; use ScopedUserIDs for per-client
	ScopedUserIDs   map[string]string // per-client userId: clientId → userId
	Scopes          []string          // parsed scope list (nil = no scopes, deny all)
	MonthlyLimit      int64           // 0 when MonthlyLimitSet is false
	MonthlyLimitSet   bool            // true if this key has a MonthlyLimit configured
	RequestsThisMonth int             // snapshot from auth time; used for approach alerts
}

// HasScope returns true if this scoped key has the given scope.
// An empty scope list means no scopes — all requests are denied.
// Master keys bypass this check (ScopedKey == nil in request context).
func (s *ScopedKeyInfo) HasScope(scope string) bool {
	if len(s.Scopes) == 0 {
		return false
	}
	return slices.Contains(s.Scopes, scope)
}

// CanAccessClient returns true if this key can access the given clientId.
// Nil/empty ScopedClientIDs means unrestricted.
func (s *ScopedKeyInfo) CanAccessClient(clientID string) bool {
	if len(s.ScopedClientIDs) == 0 {
		return true
	}
	return slices.Contains(s.ScopedClientIDs, clientID)
}

// GetUser returns the authenticated User from the request context, if present.
func (c *RequestContext) GetUser() (*model.User, bool) {
	if c == nil {
		return nil, false
	}
	u, ok := c.User.(*model.User)
	return u, ok && u != nil
}

// AutoStartFunc is called when no client is found, to attempt headless auto-start.
// targetClientID is the world to start; empty string means no specific target (no-op).
// Set by the server at startup. Returns clientID or empty string.
var AutoStartFunc func(reqCtx *RequestContext, targetClientID string) string

type contextKey string

const RequestContextKey contextKey = "requestContext"

// GetRequestContext extracts auth context from the request.
func GetRequestContext(r *http.Request) *RequestContext {
	ctx, _ := r.Context().Value(RequestContextKey).(*RequestContext)
	return ctx
}

// CreateAPIRoute creates a standardized HTTP handler that:
// 1. Parses and validates parameters
// 2. Enforces scoped key constraints
// 3. Auto-resolves clientId
// 4. Sends WS message to Foundry and waits for response
func CreateAPIRoute(manager *ws.ClientManager, pending *ws.PendingRequests, cfg APIRouteConfig) http.HandlerFunc {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	allDefs := append(cfg.RequiredParams, cfg.OptionalParams...)

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse body for POST/PUT/PATCH/DELETE
		var body map[string]interface{}
		if r.Method != http.MethodGet && r.Body != nil {
			json.NewDecoder(r.Body).Decode(&body)
		}

		// Extract parameters
		params, err := ExtractParams(r, body, allDefs)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Validate required parameters
		if err := ValidateRequired(params, cfg.RequiredParams); err != nil {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Custom validation
		if cfg.ValidateParams != nil {
			if errData, hasError := cfg.ValidateParams(params, r); hasError {
				WriteJSON(w, http.StatusBadRequest, errData)
				return
			}
		}

		// Scoped key enforcement
		reqCtx := GetRequestContext(r)
		if reqCtx != nil && reqCtx.ScopedKey != nil {
			// Multi-client scoping: validate or auto-select allowed client
			if len(reqCtx.ScopedKey.ScopedClientIDs) > 0 {
				requestedClient := params.GetString("clientId")
				if requestedClient != "" {
					if !reqCtx.ScopedKey.CanAccessClient(requestedClient) {
						WriteError(w, http.StatusForbidden, "API key is not authorized for this client")
						return
					}
				} else if len(reqCtx.ScopedKey.ScopedClientIDs) == 1 {
					// Auto-select if only one client allowed
					params["clientId"] = reqCtx.ScopedKey.ScopedClientIDs[0]
				}
				// else: multiple allowed, none specified — fall through to auto-resolve
				// which will be filtered below
			} else if reqCtx.ScopedKey.ScopedClientID != "" {
				params["clientId"] = reqCtx.ScopedKey.ScopedClientID
			}
		}

		// Auto-resolve clientId
		clientID := params.GetString("clientId")
		if clientID == "" {
			masterKey := ""
			if reqCtx != nil {
				masterKey = reqCtx.MasterAPIKey
			}
			if masterKey != "" {
				clients := manager.GetConnectedClients(masterKey)

				// Filter to allowed clients when multi-client scoping is active
				if reqCtx != nil && reqCtx.ScopedKey != nil && len(reqCtx.ScopedKey.ScopedClientIDs) > 0 {
					filtered := make([]string, 0, len(clients))
					for _, c := range clients {
						if reqCtx.ScopedKey.CanAccessClient(c) {
							filtered = append(filtered, c)
						}
					}
					clients = filtered
				}

				switch len(clients) {
				case 1:
					clientID = clients[0]
					params["clientId"] = clientID
				case 0:
					// Try auto-start (no specific clientId known at this point; no-op for unrestricted keys)
					if AutoStartFunc != nil && reqCtx != nil {
						if autoID := AutoStartFunc(reqCtx, ""); autoID != "" {
							clientID = autoID
							params["clientId"] = clientID
							break
						}
					}
					WriteError(w, http.StatusNotFound, "No connected Foundry clients found")
					return
				default:
					WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
						"error":            "Multiple clients connected. Please specify clientId.",
						"connectedClients": clients,
					})
					return
				}
			} else {
				WriteError(w, http.StatusBadRequest, "'clientId' is required")
				return
			}
		}

		// Inject userId after clientId is resolved (per-client or global fallback)
		if reqCtx != nil && reqCtx.ScopedKey != nil {
			if len(reqCtx.ScopedKey.ScopedUserIDs) > 0 {
				if uid := reqCtx.ScopedKey.ScopedUserIDs[clientID]; uid != "" {
					params["userId"] = uid
				}
			} else if reqCtx.ScopedKey.ScopedUserID != "" {
				params["userId"] = reqCtx.ScopedKey.ScopedUserID
			}
		}

		// Get client
		client := manager.GetClient(clientID)
		if client == nil {
			// Debug: log the state
			hasScopedKey := reqCtx != nil && reqCtx.ScopedKey != nil
			hasAutoStart := AutoStartFunc != nil
			ev := log.Warn().
				Str("clientId", clientID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Bool("hasScopedKey", hasScopedKey).
				Bool("hasAutoStart", hasAutoStart)
			if u, ok := reqCtx.GetUser(); ok {
				ev = ev.Int64("userId", u.ID)
			}
			ev.Msg("Client not found, checking auto-start")

			// Try auto-start using the world's stored credential
			if AutoStartFunc != nil && reqCtx != nil && reqCtx.ScopedKey != nil {
				if autoID := AutoStartFunc(reqCtx, clientID); autoID != "" {
					clientID = autoID
					params["clientId"] = clientID
					client = manager.GetClient(clientID)
				}
			}
		}
		if client == nil {
			WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"error":    "Invalid client ID",
				"clientId": clientID,
				"message":  "No Foundry client connected with this ID.",
			})
			return
		}

		// Build request ID
		requestID := fmt.Sprintf("%s_%d", cfg.Type, time.Now().UnixMilli())

		// Register pending request
		responseCh := make(chan *ws.WSResponse, 1)
		pendingReq := &ws.PendingRequest{
			ResponseCh: responseCh,
			Type:       cfg.Type,
			ClientID:   clientID,
			Timestamp:  time.Now(),
		}
		if cfg.BuildPendingRequest != nil {
			extra := cfg.BuildPendingRequest(params)
			if extra != nil {
				pendingReq.Format = extra.Format
				pendingReq.ActiveTab = extra.ActiveTab
				pendingReq.DarkMode = extra.DarkMode
				pendingReq.InitScale = extra.InitScale
				pendingReq.UUID = extra.UUID
				pendingReq.Path = extra.Path
				pendingReq.Query = extra.Query
				pendingReq.Filter = extra.Filter
			}
		}
		pending.Store(requestID, pendingReq)

		// Build payload
		payload := make(map[string]interface{})
		if cfg.BuildPayload != nil {
			payload = cfg.BuildPayload(params, r)
		} else {
			for k, v := range params {
				if k != "clientId" && k != "type" {
					payload[k] = v
				}
			}
		}

		// Add type and requestId
		payload["type"] = cfg.Type
		payload["requestId"] = requestID

		// Inject scopes for scoped-key requests so the module can enforce
		// entity-type-level restrictions (e.g. macro:write) after UUID lookup.
		// Master keys (ScopedKey == nil) get no _scopes field; the module treats
		// absence as "all scopes allowed". The relay constructs this from the
		// authenticated key — it is never populated from attacker-controlled input.
		if reqCtx != nil && reqCtx.ScopedKey != nil {
			payload["_scopes"] = reqCtx.ScopedKey.Scopes
		}

		// Ensure data sub-object exists
		if _, ok := payload["data"]; !ok {
			payload["data"] = map[string]interface{}{}
		}

		// Send to Foundry client
		wsSendStart := time.Now()
		if !client.Send(payload) {
			pending.Delete(requestID)
			WriteError(w, http.StatusInternalServerError, "Failed to send request to Foundry client")
			return
		}

		if cfg.AfterForward != nil {
			go cfg.AfterForward(reqCtx, clientID, params)
		}

		// Wait for response with timeout
		timer := time.NewTimer(cfg.Timeout)
		defer timer.Stop()
		select {
		case resp := <-responseCh:
			wsRoundTrip := time.Since(wsSendStart)
			if wsRoundTrip > 500*time.Millisecond {
				ev := log.Warn().
					Str("type", cfg.Type).
					Str("requestId", requestID).
					Dur("roundTrip", wsRoundTrip).
					Str("clientId", clientID).
					Str("method", r.Method).
					Str("path", r.URL.Path)
				if u, ok := reqCtx.GetUser(); ok {
					ev = ev.Int64("userId", u.ID)
				}
				if reqCtx != nil && reqCtx.ScopedKey != nil {
					ev = ev.Str("keyPrefix", reqCtx.ScopedKey.KeyPrefix)
				}
				ev.Msg("Slow WS round-trip")
			}
			if resp == nil {
				metrics.WSRoundTripTimeouts.Inc()
				WriteError(w, http.StatusGatewayTimeout, "Request timed out")
				return
			}
			if resp.RawData != nil {
				// Check if this is binary data with metadata headers, or raw JSON pass-through
				if resp.Data != nil {
					if ct, ok := resp.Data["_contentType"].(string); ok {
						// Binary response (file download, screenshot)
						w.Header().Set("Content-Type", ct)
						if cd, ok := resp.Data["_contentDisposition"].(string); ok {
							w.Header().Set("Content-Disposition", cd)
						}
						if cl, ok := resp.Data["_contentLength"].(int); ok {
							w.Header().Set("Content-Length", fmt.Sprintf("%d", cl))
						}
						if width := resp.Data["width"]; width != nil {
							w.Header().Set("X-Image-Width", fmt.Sprintf("%v", width))
						}
						if height := resp.Data["height"]; height != nil {
							w.Header().Set("X-Image-Height", fmt.Sprintf("%v", height))
						}
						w.WriteHeader(resp.StatusCode)
						w.Write(resp.RawData)
					} else {
						// Raw JSON pass-through
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(resp.StatusCode)
						w.Write(resp.RawData)
					}
				} else {
					// Raw JSON pass-through (no metadata)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(resp.StatusCode)
					w.Write(resp.RawData)
				}
			} else if resp.Data != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(resp.StatusCode)
				json.NewEncoder(w).Encode(resp.Data)
			}
		case <-timer.C:
			pending.Delete(requestID)
			metrics.WSRoundTripTimeouts.Inc()
			WriteError(w, http.StatusRequestTimeout, "Request timed out")
		case <-r.Context().Done():
			pending.Delete(requestID)
		}
	}
}
