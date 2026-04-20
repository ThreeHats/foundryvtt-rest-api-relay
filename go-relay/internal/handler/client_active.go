package handler

import (
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
)

// ClientActiveHandler returns whether a clientId currently has an active
// WebSocket connection on any relay instance.
//
// GET /api/clients/{clientId}/active
//
// Public, rate-limited (see middleware.ProbeRateLimiter). Used by the Foundry
// module's init wizard to decide whether to prompt a newly-joining GM for
// pairing — when another GM is already connected, the prompt is suppressed.
//
// Response shape:
//
//	{ "active": true }   // 200 — clientId is registered AND a client is connected
//	{ "active": false }  // 200 — clientId is registered but no client is connected
//	404                  // clientId is not registered to any user
//
// We deliberately do NOT distinguish "clientId doesn't exist" from
// "clientId exists but no user owns it" — both return 404. This avoids
// leaking which clientIds are paired vs unknown.
//
// Memory mode: there's no KnownClients table, so we just return active=true
// if the manager has a live client for this ID, false otherwise. No 404.
func ClientActiveHandler(db *database.DB, manager *ws.ClientManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := chi.URLParam(r, "clientId")
		if clientID == "" {
			helpers.WriteError(w, http.StatusBadRequest, "clientId is required")
			return
		}

		// Look up the clientId in KnownClients (any user). We don't need a
		// per-user filter because the clientId itself is opaque and unique
		// (fvtt_{random16}); enumeration would require guessing 62^16 values.
		// FindAnyByClientID checks across all users.
		ctx := r.Context()
		known, err := db.KnownClientStore().FindAnyByClientID(ctx, clientID)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Lookup failed")
			return
		}
		if known == nil {
			helpers.WriteError(w, http.StatusNotFound, "Unknown clientId")
			return
		}

		// Check local manager first, then Redis for cross-instance.
		active := false
		if manager != nil {
			if manager.GetClient(clientID) != nil {
				active = true
			} else if instanceID, err := manager.GetClientInstanceID(ctx, clientID); err == nil && instanceID != "" {
				active = true
			}
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]any{"active": active})
	}
}
