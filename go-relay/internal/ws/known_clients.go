package ws

import (
	"context"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/rs/zerolog/log"
)

// KnownClientsConfig wires the get-known-clients handler with its dependencies.
type KnownClientsConfig struct {
	Manager *ClientManager
	DB      *database.DB
}

// RegisterKnownClientsHandler installs the WS message handler for the
// `get-known-clients` action. It returns all Foundry worlds (KnownClients)
// belonging to the authenticated connection token's account, along with
// live online/offline status and headless auto-start capability.
//
// @ws-type    get-known-clients
// @ws-result  known-clients-result
// @ws-summary List all worlds known to this relay account (online and offline)
// @ws-description Requires connection-token authentication. Returns all
// @ws-description KnownClients for the account, including offline worlds that
// @ws-description can be auto-started via remote-request.
// @ws-note    API-key-only connections receive an error response.
// @ws-send requestId {string} required  Correlation ID echoed back in the result
// @ws-recv clients   {array}  required  Array of KnownClient objects (may be empty)
// @ws-recv error     {string} optional  Error message if the request failed
func RegisterKnownClientsHandler(cfg KnownClientsConfig) {
	if cfg.Manager == nil || cfg.DB == nil {
		log.Warn().Msg("RegisterKnownClientsHandler called with incomplete config; skipping")
		return
	}

	cfg.Manager.OnMessageType("get-known-clients", func(client *Client, data map[string]interface{}) {
		requestID, _ := data["requestId"].(string)

		sendError := func(msg string) {
			client.Send(map[string]interface{}{
				"type":      "known-clients-result",
				"requestId": requestID,
				"error":     msg,
				"clients":   []interface{}{},
			})
		}

		// Require connection-token authentication
		tokenID := client.ConnectionTokenID()
		if tokenID == 0 {
			sendError("get-known-clients requires connection-token authentication")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Resolve user account via connection token
		token, err := cfg.DB.ConnectionTokenStore().FindByID(ctx, tokenID)
		if err != nil || token == nil {
			sendError("connection token not found or revoked")
			return
		}

		// Fetch all known clients for this user account
		knownClients, err := cfg.DB.KnownClientStore().FindAllByUser(ctx, token.UserID)
		if err != nil {
			sendError("failed to retrieve known clients")
			return
		}

		// Build response list, augmenting each record with live online status
		clientList := make([]map[string]interface{}, 0, len(knownClients))
		for _, kc := range knownClients {
			isOnline := cfg.Manager.GetClient(kc.ClientID) != nil

			worldTitle := ""
			if kc.WorldTitle.Valid {
				worldTitle = kc.WorldTitle.String
			}
			customName := ""
			if kc.CustomName.Valid {
				customName = kc.CustomName.String
			}
			systemTitle := ""
			if kc.SystemTitle.Valid {
				systemTitle = kc.SystemTitle.String
			}
			systemID := ""
			if kc.SystemID.Valid {
				systemID = kc.SystemID.String
			}

			clientList = append(clientList, map[string]interface{}{
				"clientId":     kc.ClientID,
				"worldTitle":   worldTitle,
				"customName":   customName,
				"systemTitle":  systemTitle,
				"systemId":     systemID,
				"isOnline":     isOnline,
				"canAutoStart": bool(kc.AutoStartOnRemoteRequest),
			})
		}

		// Cross-world permissions are world-level (KnownClient), not per-token.
		// The upgradeOnly flow updates KnownClient scopes, so read from there.
		// Use client.ID() (the connected world's relay clientId) rather than
		// token.ClientID, because headless tokens are created without ClientID set.
		currentClientID := client.ID()
		var tokenScopes []string
		var tokenAllowedTargets []string
		for _, kc := range knownClients {
			if kc.ClientID == currentClientID {
				tokenScopes = kc.GetRemoteScopes()
				tokenAllowedTargets = kc.GetAllowedTargets()
				break
			}
		}
		if tokenScopes == nil {
			tokenScopes = []string{}
		}
		if tokenAllowedTargets == nil {
			tokenAllowedTargets = []string{}
		}

		client.Send(map[string]interface{}{
			"type":                "known-clients-result",
			"requestId":           requestID,
			"clients":             clientList,
			"tokenScopes":         tokenScopes,
			"tokenAllowedTargets": tokenAllowedTargets,
		})

		log.Debug().
			Str("clientId", client.ID()).
			Int("count", len(clientList)).
			Msg("Sent known-clients-result")
	})

	log.Info().Msg("Registered get-known-clients handler")
}
