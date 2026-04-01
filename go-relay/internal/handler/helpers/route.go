package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
	ID             int64
	ScopedClientID string
	ScopedUserID   string
}

// AutoStartFunc is called when no client is found and the request has a scoped key with stored credentials.
// Set by the server at startup. Returns clientID or empty string.
var AutoStartFunc func(reqCtx *RequestContext) string

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
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil && len(bodyBytes) > 0 {
				json.Unmarshal(bodyBytes, &body)
			}
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
			if reqCtx.ScopedKey.ScopedClientID != "" {
				params["clientId"] = reqCtx.ScopedKey.ScopedClientID
			}
			if reqCtx.ScopedKey.ScopedUserID != "" {
				params["userId"] = reqCtx.ScopedKey.ScopedUserID
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
				switch len(clients) {
				case 1:
					clientID = clients[0]
					params["clientId"] = clientID
				case 0:
					// Try auto-start for scoped keys with stored credentials
					if AutoStartFunc != nil && reqCtx != nil {
						if autoID := AutoStartFunc(reqCtx); autoID != "" {
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

		// Get client
		client := manager.GetClient(clientID)
		if client == nil {
			// Debug: log the state
			hasScopedKey := reqCtx != nil && reqCtx.ScopedKey != nil
			hasAutoStart := AutoStartFunc != nil
			log.Warn().Str("clientId", clientID).Bool("hasScopedKey", hasScopedKey).Bool("hasAutoStart", hasAutoStart).Msg("Client not found, checking auto-start")

			// Try auto-start for scoped keys with stored credentials
			if AutoStartFunc != nil && reqCtx != nil && reqCtx.ScopedKey != nil {
				if autoID := AutoStartFunc(reqCtx); autoID != "" {
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
				"message":  "No Foundry client connected with this ID. If using a scoped key with stored credentials, check server logs for auto-start errors.",
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

		// Wait for response with timeout
		select {
		case resp := <-responseCh:
			wsRoundTrip := time.Since(wsSendStart)
			if wsRoundTrip > 500*time.Millisecond {
				log.Warn().Str("type", cfg.Type).Str("requestId", requestID).Dur("roundTrip", wsRoundTrip).Msg("Slow WS round-trip")
			}
			if resp == nil {
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
		case <-time.After(cfg.Timeout):
			pending.Delete(requestID)
			WriteError(w, http.StatusRequestTimeout, "Request timed out")
		case <-r.Context().Done():
			pending.Delete(requestID)
		}
	}
}
