package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// RegisterPairRequestRoutes registers routes for the OAuth device-flow world pairing.
//
// Normal pairing flow:
//  1. Foundry module POSTs /auth/pair-request with world metadata → gets {code, pairUrl}
//  2. Module opens pairUrl in a browser tab
//  3. User logs into the relay web UI, optionally configures cross-world settings, clicks Approve
//  4. Module polls /auth/pair-request/{code}/status until status="approved"
//  5. Module exchanges pairingCode via the existing POST /auth/pair endpoint
//
// Upgrade-only flow (upgradeOnly=true):
//  - Same steps 1–4, but step 5 is skipped — no new token is created.
//  - The approval handler updates the already-paired world's KnownClient directly.
//  - The status poll returns {status:"approved", upgraded:true} without a pairingCode.
func RegisterPairRequestRoutes(r chi.Router, db *database.DB, cfg *config.Config, manager *ws.ClientManager) {
	r.Route("/pair-request", func(r chi.Router) {

		// POST /auth/pair-request — Create pending pair request (NO AUTH, RATE LIMITED)
		r.With(middleware.KeyRequestRateLimiter.Middleware).Post("/", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				WorldID                string   `json:"worldId"`
				WorldTitle             string   `json:"worldTitle"`
				SystemID               string   `json:"systemId"`
				SystemTitle            string   `json:"systemTitle"`
				SystemVersion          string   `json:"systemVersion"`
				FoundryVersion         string   `json:"foundryVersion"`
				RequestedRemoteScopes  []string `json:"requestedRemoteScopes"`
				RequestedTargetClients []string `json:"requestedTargetClients"`
				UpgradeOnly            bool     `json:"upgradeOnly"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}
			if body.WorldID == "" {
				helpers.WriteError(w, http.StatusBadRequest, "worldId is required")
				return
			}

			// Validate any requested scopes against the canonical list.
			if invalid := model.ValidateScopes(body.RequestedRemoteScopes); invalid != "" {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid scope: "+invalid)
				return
			}

			code, err := generateKeyRequestCode(8)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate code")
				return
			}

			expiresAt := time.Now().Add(10 * time.Minute)
			req := &model.PairRequest{
				Code:           code,
				WorldID:        body.WorldID,
				WorldTitle:     body.WorldTitle,
				SystemID:       body.SystemID,
				SystemTitle:    body.SystemTitle,
				SystemVersion:  body.SystemVersion,
				FoundryVersion: body.FoundryVersion,
				UpgradeOnly:    body.UpgradeOnly,
				Status:         model.PairRequestStatusPending,
				ExpiresAt:      model.NewSQLiteTime(expiresAt),
			}
			if len(body.RequestedRemoteScopes) > 0 {
				req.RequestedRemoteScopes = sql.NullString{
					String: model.ScopesString(body.RequestedRemoteScopes), Valid: true,
				}
			}
			if len(body.RequestedTargetClients) > 0 {
				req.RequestedTargetClients = sql.NullString{
					String: model.ScopesString(body.RequestedTargetClients), Valid: true,
				}
			}
			if err := db.PairRequestStore().Create(r.Context(), req); err != nil {
				log.Error().Err(err).Msg("Failed to create pair request")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to create pair request")
				return
			}

			pairURL := cfg.FrontendURL + "/pair/" + code
			helpers.WriteJSON(w, http.StatusCreated, map[string]interface{}{
				"code":        code,
				"pairUrl":     pairURL,
				"upgradeOnly": body.UpgradeOnly,
				"expiresIn":   600,
				"expiresAt":   expiresAt.Format(time.RFC3339),
			})
		})

		// GET /auth/pair-request/{code}/status — Poll for device flow (NO AUTH, RATE LIMITED)
		r.With(middleware.KeyRequestRateLimiter.Middleware).Get("/{code}/status", func(w http.ResponseWriter, r *http.Request) {
			code := chi.URLParam(r, "code")
			ctx := r.Context()

			req, err := db.PairRequestStore().FindByCode(ctx, code)
			if err != nil || req == nil {
				helpers.WriteError(w, http.StatusNotFound, "Pair request not found")
				return
			}

			if req.ExpiresAt.Time.Before(time.Now()) && req.Status == model.PairRequestStatusPending {
				helpers.WriteJSON(w, http.StatusGone, map[string]interface{}{"status": "expired"})
				return
			}

			switch req.Status {
			case model.PairRequestStatusPending:
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "pending"})
			case model.PairRequestStatusApproved:
				// For upgradeOnly: no pairing code — just signal success.
				if req.UpgradeOnly || !req.PairingCode.Valid {
					db.PairRequestStore().SetStatus(ctx, req.ID, model.PairRequestStatusExchanged)
					helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
						"status":   "approved",
						"upgraded": req.UpgradeOnly,
					})
					return
				}
				// Normal pairing: return the code once, then mark exchanged.
				db.PairRequestStore().SetStatus(ctx, req.ID, model.PairRequestStatusExchanged)
				helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{
					"status":      "approved",
					"pairingCode": req.PairingCode.String,
				})
			case model.PairRequestStatusExchanged:
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"status":  "exchanged",
					"message": "Pairing code has already been retrieved",
				})
			case model.PairRequestStatusDenied:
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "denied"})
			case model.PairRequestStatusExpired:
				helpers.WriteJSON(w, http.StatusGone, map[string]interface{}{"status": "expired"})
			default:
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": req.Status})
			}
		})

		// Authenticated routes — require a valid session token
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(db, nil))

			// GET /auth/pair-request/{code} — Fetch details for the approval page.
			// Also returns all known clients so the UI can render target world checkboxes.
			r.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
				reqCtx := helpers.GetRequestContext(r)
				user, _ := reqCtx.User.(*model.User)

				code := chi.URLParam(r, "code")
				req, err := db.PairRequestStore().FindByCode(r.Context(), code)
				if err != nil || req == nil {
					helpers.WriteError(w, http.StatusNotFound, "Pair request not found")
					return
				}

				// Fetch known clients for target-world selection in the UI.
				// Flatten sql.NullString fields so the frontend receives plain strings.
				knownClients := make([]map[string]interface{}, 0)
				if user != nil {
					if clients, err := db.KnownClientStore().FindAllByUser(r.Context(), user.ID); err == nil {
						for _, c := range clients {
							knownClients = append(knownClients, map[string]interface{}{
								"id":          c.ID,
								"clientId":    c.ClientID,
								"worldId":     c.WorldID.String,
								"worldTitle":  c.WorldTitle.String,
								"customName":  c.CustomName.String,
								"isOnline":    manager.IsClientOnlineAnywhere(r.Context(), c.ClientID),
							})
						}
					}
				}

				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"code":                   req.Code,
					"worldId":                req.WorldID,
					"worldTitle":             req.WorldTitle,
					"systemId":               req.SystemID,
					"systemTitle":            req.SystemTitle,
					"systemVersion":          req.SystemVersion,
					"foundryVersion":         req.FoundryVersion,
					"requestedRemoteScopes":  model.ParseScopes(req.RequestedRemoteScopes.String),
					"requestedTargetClients": model.ParseScopes(req.RequestedTargetClients.String),
					"upgradeOnly":            req.UpgradeOnly,
					"status":                 req.Status,
					"expiresAt":              req.ExpiresAt,
					"knownClients":           knownClients,
				})
			})

			// POST /auth/pair-request/{code}/approve
			r.Post("/{code}/approve", func(w http.ResponseWriter, r *http.Request) {
				reqCtx := helpers.GetRequestContext(r)
				if reqCtx == nil || reqCtx.ScopedKey != nil {
					helpers.WriteError(w, http.StatusForbidden, "Use your master API key to approve pair requests.")
					return
				}
				user, ok := reqCtx.User.(*model.User)
				if !ok || user == nil {
					helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
					return
				}

				code := chi.URLParam(r, "code")
				ctx := r.Context()

				pairReq, err := db.PairRequestStore().FindByCode(ctx, code)
				if err != nil || pairReq == nil {
					helpers.WriteError(w, http.StatusNotFound, "Pair request not found")
					return
				}
				if pairReq.Status != model.PairRequestStatusPending {
					helpers.WriteError(w, http.StatusBadRequest, "Pair request is not pending")
					return
				}
				if pairReq.ExpiresAt.Time.Before(time.Now()) {
					helpers.WriteError(w, http.StatusGone, "Pair request has expired")
					return
				}

				// Parse approved cross-world settings from the approval body.
				var approveBody struct {
					RemoteScopes         []string `json:"remoteScopes"`
					AllowedTargetClients []string `json:"allowedTargetClients"`
					RemoteRequestsPerHour int     `json:"remoteRequestsPerHour"`
				}
				json.NewDecoder(r.Body).Decode(&approveBody)

				approvedScopes := model.ScopesString(approveBody.RemoteScopes)
				approvedTargets := model.ScopesString(approveBody.AllowedTargetClients)

				// Look up whether the world is already registered.
				existing, _ := db.KnownClientStore().FindByWorldID(ctx, user.ID, pairReq.WorldID)

				// ── Upgrade-only: world is already paired, just update cross-world settings ──
				if pairReq.UpgradeOnly {
					if existing == nil {
						helpers.WriteError(w, http.StatusBadRequest, "World is not yet paired — cannot upgrade. Pair first.")
						return
					}
					if err := db.KnownClientStore().SetCrossWorldSettings(ctx, existing.ID,
						sql.NullString{String: approvedTargets, Valid: approvedTargets != ""},
						sql.NullString{String: approvedScopes, Valid: approvedScopes != ""},
						approveBody.RemoteRequestsPerHour,
					); err != nil {
						log.Error().Err(err).Msg("Failed to apply cross-world settings during upgrade")
						helpers.WriteError(w, http.StatusInternalServerError, "Failed to apply settings")
						return
					}
					db.PairRequestStore().SetApproved(ctx, pairReq.ID, "", user.ID)
					log.Info().Str("code", code).Str("worldId", pairReq.WorldID).Int64("userId", user.ID).
						Msg("Pair request upgrade approved")
					helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true, "upgraded": true})
					return
				}

				// ── Normal pairing: create a PairingCode the module exchanges ──
				var clientID sql.NullString
				if existing != nil {
					clientID = sql.NullString{String: existing.ClientID, Valid: true}
				}

				pCode, err := generatePairingCode(6)
				if err != nil {
					helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate pairing code")
					return
				}

				// Thread cross-world settings through the PairingCode so that
				// POST /auth/pair can apply them to the KnownClient after creation.
				pairingCode := &model.PairingCode{
					UserID:   user.ID,
					Code:     pCode,
					ClientID: clientID,
					AllowedTargetClients: sql.NullString{
						String: approvedTargets, Valid: approvedTargets != "",
					},
					RemoteScopes: sql.NullString{
						String: approvedScopes, Valid: approvedScopes != "",
					},
					ExpiresAt: model.NewSQLiteTime(time.Now().Add(5 * time.Minute)),
				}
				if err := db.PairingCodeStore().Create(ctx, pairingCode); err != nil {
					log.Error().Err(err).Msg("Failed to create pairing code for pair request")
					helpers.WriteError(w, http.StatusInternalServerError, "Failed to create pairing code")
					return
				}

				// If the world is already known, apply cross-world settings immediately
				// (covers the case where the user is re-pairing and wants to grant access).
				if existing != nil && (approvedScopes != "" || approvedTargets != "") {
					_ = db.KnownClientStore().SetCrossWorldSettings(ctx, existing.ID,
						sql.NullString{String: approvedTargets, Valid: approvedTargets != ""},
						sql.NullString{String: approvedScopes, Valid: approvedScopes != ""},
						approveBody.RemoteRequestsPerHour,
					)
				}

				if err := db.PairRequestStore().SetApproved(ctx, pairReq.ID, pCode, user.ID); err != nil {
					log.Error().Err(err).Msg("Failed to update pair request")
					helpers.WriteError(w, http.StatusInternalServerError, "Failed to approve pair request")
					return
				}

				log.Info().Str("code", code).Str("worldId", pairReq.WorldID).Int64("userId", user.ID).
					Msg("Pair request approved")

				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true})
			})

			// POST /auth/pair-request/{code}/deny
			r.Post("/{code}/deny", func(w http.ResponseWriter, r *http.Request) {
				reqCtx := helpers.GetRequestContext(r)
				if reqCtx == nil || reqCtx.ScopedKey != nil {
					helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
					return
				}
				code := chi.URLParam(r, "code")
				pairReq, err := db.PairRequestStore().FindByCode(r.Context(), code)
				if err != nil || pairReq == nil {
					helpers.WriteError(w, http.StatusNotFound, "Pair request not found")
					return
				}
				if pairReq.Status != model.PairRequestStatusPending {
					helpers.WriteError(w, http.StatusBadRequest, "Pair request is not pending")
					return
				}
				db.PairRequestStore().SetStatus(r.Context(), pairReq.ID, model.PairRequestStatusDenied)
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"success": true,
					"message": "Pair request denied",
				})
			})
		})
	})
}
