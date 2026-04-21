package ws

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/rs/zerolog/log"
)

// HeadlessAutoStarter is the contract the remote-request handler needs from
// the headless worker. Defined here as an interface to avoid pulling the full
// worker package into the ws package (and the import cycle that would
// create). The server wires the actual *worker.HeadlessManager when
// RegisterRemoteRequestHandler is called.
type HeadlessAutoStarter interface {
	// AutoStartForKnownClient launches a headless session for the given
	// userId + clientId pair, using the relay's stored Foundry credentials,
	// and blocks until the resulting client is connected (or times out).
	// Returns the connected clientId on success.
	AutoStartForKnownClient(ctx context.Context, userID int64, clientID string) (string, error)
}

// rateBucket tracks per-token remote-request counts for hourly rate limiting.
type rateBucket struct {
	count       int
	windowStart time.Time
}

var (
	rateBuckets   = make(map[int64]*rateBucket) // tokenID → bucket
	rateBucketsMu sync.Mutex
)

// checkRemoteRequestRate returns true if the request is allowed, false if the
// rate limit has been exceeded. Increments the counter on success.
func checkRemoteRequestRate(tokenID int64, limitPerHour int) bool {
	rateBucketsMu.Lock()
	defer rateBucketsMu.Unlock()

	bucket, ok := rateBuckets[tokenID]
	if !ok || time.Since(bucket.windowStart) > time.Hour {
		// New window
		rateBuckets[tokenID] = &rateBucket{count: 1, windowStart: time.Now()}
		return true
	}
	if bucket.count >= limitPerHour {
		return false
	}
	bucket.count++
	return true
}

// RemoteRequestBatcherInterface is the subset of service.RemoteRequestBatcher
// the handler needs. Defined here as an interface to avoid importing the
// service package (and the cycle that would create).
type RemoteRequestBatcherInterface interface {
	Add(userID int64, sourceTokenID int64, targetClientID, action, tokenName string, success bool)
}

// RemoteRequestConfig wires the remote-request handler with its dependencies.
type RemoteRequestConfig struct {
	Manager  *ClientManager
	Pending  *PendingRequests
	DB       *database.DB
	Headless HeadlessAutoStarter           // may be nil if headless is disabled
	Batcher  RemoteRequestBatcherInterface // may be nil; falls back to no notification
	// ForwardToInstance proxies an action to a client on a different relay instance.
	// Nil on single-instance deployments or when Redis is unavailable.
	ForwardToInstance func(ctx context.Context, instanceID, targetClientID, action string, payload map[string]interface{}) (map[string]interface{}, error)
}

// RegisterRemoteRequestHandler installs the WS message handler for the
// `remote-request` action. This is the cross-world tunnel: a connected
// Foundry module (source) asks the relay to invoke an action on another
// Foundry module (target) owned by the same account, gated by the source
// connection token's allowedTargetClients + remoteScopes.
//
// @ws-type    remote-request
// @ws-result  remote-response
// @ws-summary Invoke any supported action on a remote Foundry world via the relay tunnel
// @ws-description The source connection token must list the target in allowedTargetClients
// @ws-description and hold the required scope in remoteScopes. Configure these in the
// @ws-description dashboard → Connections → Edit browser.
// @ws-note    Foundry module connection token required; API key clients are rejected
// @ws-send targetClientId     {string}  required  Client ID of the target world (must be in allowedTargetClients)
// @ws-send action             {string}  required  Action to invoke on the target (e.g. create-user, entity, search)
// @ws-send payload            {object}  optional  Action payload forwarded verbatim to the target module
// @ws-send autoStartIfOffline {boolean} optional  Start a headless session if the target is offline (requires autoStartOnRemoteRequest enabled for that world)
// @ws-recv success            {boolean} required  Whether the action succeeded
// @ws-recv data               {object}  optional  Response data from the target module (present on success)
// @ws-recv error              {string}  optional  Error message if the action failed or was rejected
//
// Wire format:
//
//	source → relay
//	{
//	  "type": "remote-request",
//	  "requestId": "src_1234",
//	  "targetClientId": "fvtt_target_world",
//	  "action": "create-user",
//	  "payload": { ... },
//	  "autoStartIfOffline": true
//	}
//
//	relay → source (eventually)
//	{
//	  "type": "remote-response",
//	  "requestId": "src_1234",
//	  "success": true,
//	  "data": { ... }
//	}
func RegisterRemoteRequestHandler(cfg RemoteRequestConfig) {
	if cfg.Manager == nil || cfg.Pending == nil || cfg.DB == nil {
		log.Warn().Msg("RegisterRemoteRequestHandler called with incomplete config; skipping")
		return
	}

	cfg.Manager.OnMessageType("remote-request", func(sourceClient *Client, data map[string]interface{}) {
		go handleRemoteRequest(cfg, sourceClient, data)
	})

	log.Info().Msg("Registered remote-request handler")
}

func handleRemoteRequest(cfg RemoteRequestConfig, source *Client, data map[string]interface{}) {
	sourceRequestID, _ := data["requestId"].(string)
	targetClientID, _ := data["targetClientId"].(string)
	action, _ := data["action"].(string)
	payload, _ := data["payload"].(map[string]interface{})
	autoStart, _ := data["autoStartIfOffline"].(bool)

	respond := func(success bool, errMsg string, result map[string]interface{}) {
		out := map[string]interface{}{
			"type":      "remote-response",
			"requestId": sourceRequestID,
			"success":   success,
		}
		if errMsg != "" {
			out["error"] = errMsg
		}
		if result != nil {
			out["data"] = result
		}
		source.Send(out)
	}

	logEntry := &model.RemoteRequestLog{
		SourceClientID: source.ID(),
		SourceTokenID:  source.ConnectionTokenID(),
		TargetClientID: targetClientID,
		Action:         action,
	}
	if ip := source.IPAddress(); ip != "" {
		logEntry.SourceIP = sql.NullString{String: ip, Valid: true}
	}
	// tokenName is populated later (after we look up the source token).
	// We capture it here so finishLog can include it in the batched notification.
	var tokenName string

	finishLog := func(success bool, errMsg string) {
		logEntry.Success = model.LooseBool(success)
		if errMsg != "" {
			logEntry.ErrorMessage = sql.NullString{String: errMsg, Valid: true}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Only persist if the user hasn't opted out of cross-world request logging.
		// If settings are nil (never configured), default to logging = on.
		if logEntry.UserID != 0 {
			settings, _ := cfg.DB.NotificationSettingsStore().FindByUser(ctx, logEntry.UserID)
			if settings == nil || settings.LogCrossWorldRequests {
				_ = cfg.DB.RemoteRequestLogStore().Create(ctx, logEntry)
			}
		} else {
			_ = cfg.DB.RemoteRequestLogStore().Create(ctx, logEntry)
		}

		// Feed into the notification batcher (5-minute window). The batcher
		// accumulates events and sends ONE summary notification instead of N.
		if cfg.Batcher != nil && logEntry.UserID != 0 {
			cfg.Batcher.Add(logEntry.UserID, logEntry.SourceTokenID, targetClientID, action, tokenName, success)
		}
	}

	if sourceRequestID == "" || targetClientID == "" || action == "" {
		respond(false, "missing required fields: requestId, targetClientId, action", nil)
		finishLog(false, "missing required fields")
		return
	}

	// Look up the source's connection token to get the user ID and token name.
	sourceTokenID := source.ConnectionTokenID()
	if sourceTokenID == 0 {
		respond(false, "remote-request requires a connection-token-authenticated source", nil)
		finishLog(false, "no connection token")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sourceToken, err := cfg.DB.ConnectionTokenStore().FindByID(ctx, sourceTokenID)
	if err != nil || sourceToken == nil {
		respond(false, "source token not found", nil)
		finishLog(false, "source token not found")
		return
	}
	logEntry.UserID = sourceToken.UserID
	tokenName = sourceToken.Name

	// Look up the source world's KnownClient to read cross-world permissions.
	// Permissions are world-level, not per-token, so we read from KnownClient.
	sourceKC, kcErr := cfg.DB.KnownClientStore().FindByClientID(ctx, sourceToken.UserID, source.ID())
	if kcErr != nil || sourceKC == nil {
		respond(false, "source world not registered", nil)
		finishLog(false, "source world not found")
		return
	}

	// Permission check 1: target must be in allowedTargetClients.
	if !sourceKC.CanTarget(targetClientID) {
		respond(false, fmt.Sprintf("target %q not in allowed clients for this world", targetClientID), nil)
		finishLog(false, "target not in allow list")
		return
	}

	// Permission check 2: action's required scope must be in remoteScopes.
	requiredScope, known := model.ScopeForAction(action)
	if !known {
		respond(false, fmt.Sprintf("action %q is not exposed via remote-request", action), nil)
		finishLog(false, "unknown action")
		return
	}
	if !sourceKC.HasRemoteScope(requiredScope) {
		respond(false, fmt.Sprintf("scope %q not granted to source world", requiredScope), nil)
		finishLog(false, "scope not granted")
		return
	}

	// Rate limit check: per-world remote-request rate limiting.
	if sourceKC.RemoteRequestsPerHour > 0 {
		if !checkRemoteRequestRate(sourceKC.ID, sourceKC.RemoteRequestsPerHour) {
			respond(false, fmt.Sprintf("remote-request rate limit exceeded (%d/hour)", sourceKC.RemoteRequestsPerHour), nil)
			finishLog(false, "rate limited")
			return
		}
	}

	// Permission check 3: target must belong to the same account.
	// Re-fetch the source token here to catch any revocation that happened
	// between the initial fetch and this point (e.g., token deleted via dashboard).
	freshToken, ftErr := cfg.DB.ConnectionTokenStore().FindByID(ctx, sourceTokenID)
	if ftErr != nil || freshToken == nil {
		respond(false, "source token not found or revoked", nil)
		finishLog(false, "source token revoked")
		return
	}
	known2, err := cfg.DB.KnownClientStore().FindByClientID(ctx, freshToken.UserID, targetClientID)
	if err != nil || known2 == nil {
		respond(false, "target not owned by source account", nil)
		finishLog(false, "target not owned")
		return
	}

	// Find or auto-start the target client.
	target := cfg.Manager.GetClient(targetClientID)
	if target == nil {
		// Check if the target is online on a different relay instance.
		if remoteInstance := cfg.Manager.GetClientRemoteInstance(ctx, targetClientID); remoteInstance != "" && cfg.ForwardToInstance != nil {
			log.Info().
				Str("sourceClient", source.ID()).
				Str("targetClient", targetClientID).
				Str("remoteInstance", remoteInstance).
				Msg("Forwarding remote-request to peer instance")
			fwdCtx, fwdCancel := context.WithTimeout(context.Background(), 60*time.Second)
			result, err := cfg.ForwardToInstance(fwdCtx, remoteInstance, targetClientID, action, payload)
			fwdCancel()
			if err != nil {
				respond(false, fmt.Sprintf("cross-instance forward failed: %s", err), nil)
				finishLog(false, "cross-instance forward failed")
			} else {
				if errMsg, ok := result["error"].(string); ok && errMsg != "" {
					respond(false, errMsg, nil)
					finishLog(false, errMsg)
				} else {
					respond(true, "", result)
					finishLog(true, "")
				}
			}
			return
		}

		if !autoStart {
			respond(false, "target offline; set autoStartIfOffline:true to start a headless session", nil)
			finishLog(false, "target offline")
			return
		}
		if !bool(known2.AutoStartOnRemoteRequest) {
			respond(false, "target offline; auto-start not enabled for this clientId", nil)
			finishLog(false, "auto-start disabled")
			return
		}
		if cfg.Headless == nil {
			respond(false, "target offline; headless worker not available on this instance", nil)
			finishLog(false, "headless unavailable")
			return
		}

		log.Info().
			Str("sourceClient", source.ID()).
			Str("targetClient", targetClientID).
			Msg("Auto-starting headless for remote-request")

		startCtx, startCancel := context.WithTimeout(context.Background(), 60*time.Second)
		_, err := cfg.Headless.AutoStartForKnownClient(startCtx, sourceToken.UserID, targetClientID)
		startCancel()
		if err != nil {
			respond(false, fmt.Sprintf("auto-start failed: %s", err), nil)
			finishLog(false, "auto-start failed")
			return
		}

		// Re-fetch the now-connected client.
		target = cfg.Manager.GetClient(targetClientID)
		if target == nil {
			respond(false, "auto-start succeeded but target still not connected", nil)
			finishLog(false, "auto-start race")
			return
		}
	}

	// Forward the action message to the target. Generate a fresh internal
	// requestId so the source's id never clashes with anything else.
	relayRequestID := fmt.Sprintf("rr_%d_%s", time.Now().UnixNano(), action)

	responseCh := make(chan *WSResponse, 1)
	cfg.Pending.Store(relayRequestID, &PendingRequest{
		ResponseCh: responseCh,
		Type:       action,
		ClientID:   target.ID(),
		Timestamp:  time.Now(),
		MaxAge:     90 * time.Second,
	})

	// Build the action message: just the payload + type + requestId. Anything
	// in the source's `payload` map is forwarded verbatim to the target.
	msg := make(map[string]interface{}, len(payload)+2)
	for k, v := range payload {
		msg[k] = v
	}
	msg["type"] = action
	msg["requestId"] = relayRequestID

	if !target.Send(msg) {
		cfg.Pending.Delete(relayRequestID)
		respond(false, "failed to send action to target", nil)
		finishLog(false, "send failed")
		return
	}

	// Wait for the target's *-result message to land in our channel.
	select {
	case resp := <-responseCh:
		if resp == nil {
			respond(false, "no response from target", nil)
			finishLog(false, "no response")
			return
		}
		// Parse the raw response if present (cheaper than re-marshaling).
		var result map[string]interface{}
		if resp.RawData != nil {
			if err := json.Unmarshal(resp.RawData, &result); err != nil {
				respond(false, fmt.Sprintf("failed to parse target response: %s", err), nil)
				finishLog(false, "parse error")
				return
			}
		} else if resp.Data != nil {
			result = resp.Data
		}

		if errMsg, ok := result["error"].(string); ok && errMsg != "" {
			respond(false, errMsg, nil)
			finishLog(false, errMsg)
			return
		}

		// Strip the relay-internal requestId from the response so the
		// source can correlate against its own.
		delete(result, "requestId")
		respond(true, "", result)
		finishLog(true, "")

	case <-time.After(60 * time.Second):
		cfg.Pending.Delete(relayRequestID)
		respond(false, "request timed out", nil)
		finishLog(false, "timeout")
	}
}
