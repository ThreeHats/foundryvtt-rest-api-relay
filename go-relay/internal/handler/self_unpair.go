package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/rs/zerolog/log"
)

// SelfUnpairHandler handles POST /api/self-unpair.
//
// Called by the Foundry module when the GM clicks "Unpair" in the connection
// dialog. The request body contains the raw connection token; the handler:
//  1. Hashes the token and looks it up in ConnectionTokens.
//  2. Force-disconnects any live WebSocket authenticated with that token.
//  3. Deletes the ConnectionToken row.
//
// The KnownClients row (world) is intentionally preserved — the world retains
// its cross-world config, credentials, and auto-start settings so the GM can
// re-pair later without losing those settings. The world can be removed
// explicitly from the dashboard.
//
// Always returns 200 — the response is informational only. The module should
// clear its local settings regardless of the outcome.
func SelfUnpairHandler(db *database.DB, manager *ws.ClientManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Token == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Missing token")
			return
		}

		hash := sha256.Sum256([]byte(body.Token))
		tokenHash := hex.EncodeToString(hash[:])

		ctx := r.Context()

		ct, err := db.ConnectionTokenStore().FindByTokenHash(ctx, tokenHash)
		if err != nil {
			log.Warn().Err(err).Msg("SelfUnpair: token lookup error")
			// Return success anyway — the module should still clear its local state.
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Token not found; nothing to clean up"})
			return
		}
		if ct == nil {
			// Already deleted or never existed — idempotent.
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Token not found; nothing to clean up"})
			return
		}

		clientID := ct.ClientID
		userID := ct.UserID

		// Force-disconnect any live WS connection authenticated with this token.
		if manager != nil {
			n := manager.BroadcastDisconnectByConnectionToken(ctx, ct.ID, "Unpaired by owner")
			if n > 0 {
				log.Info().Str("clientId", clientID).Int("connections", n).Msg("SelfUnpair: force-disconnected live connection(s)")
			}
		}

		// Delete the connection token.
		if err := db.ConnectionTokenStore().Delete(ctx, ct.ID); err != nil {
			log.Warn().Err(err).Int64("tokenId", ct.ID).Msg("SelfUnpair: failed to delete token")
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": false, "message": "Failed to delete token"})
			return
		}

		log.Info().Str("clientId", clientID).Int64("userId", userID).Msg("SelfUnpair: connection token deleted")

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Unpaired successfully"})
	}
}
