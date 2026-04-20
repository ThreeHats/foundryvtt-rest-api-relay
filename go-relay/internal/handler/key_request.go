package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// RegisterKeyRequestRoutes registers routes for the OAuth-like key request flow.
func RegisterKeyRequestRoutes(r chi.Router, db *database.DB, cfg *config.Config) {
	r.Route("/key-request", func(r chi.Router) {

		// POST /auth/key-request — Create pending request (NO AUTH, RATE LIMITED)
		r.With(middleware.KeyRequestRateLimiter.Middleware).Post("/", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				AppName            string   `json:"appName"`
				AppDescription     string   `json:"appDescription"`
				AppURL             string   `json:"appUrl"`
				Scopes             []string `json:"scopes"`
				ClientIDs          []string `json:"clientIds"`
				CallbackURL        string   `json:"callbackUrl"`
				SuggestedMonthlyLimit *int64  `json:"suggestedMonthlyLimit"`
				SuggestedExpiry    string   `json:"suggestedExpiry"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}
			if body.AppName == "" || len(body.Scopes) == 0 {
				helpers.WriteError(w, http.StatusBadRequest, "appName and scopes are required")
				return
			}

			// Validate scopes
			if invalid := model.ValidateScopes(body.Scopes); invalid != "" {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid scope: "+invalid)
				return
			}

			code, err := generateKeyRequestCode(6)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate code")
				return
			}

			expiresAt := time.Now().Add(10 * time.Minute)

			req := &model.KeyRequest{
				Code:               code,
				AppName:            body.AppName,
				AppDescription:     body.AppDescription,
				RequestedScopes:    model.ScopesString(body.Scopes),
				Status:             "pending",
				ExpiresAt:          model.NewSQLiteTime(expiresAt),
			}
			if body.AppURL != "" {
				if err := validateWebURL(body.AppURL); err != nil {
					helpers.WriteError(w, http.StatusBadRequest, "Invalid appUrl: "+err.Error())
					return
				}
				req.AppURL = sql.NullString{String: body.AppURL, Valid: true}
			}
			if body.CallbackURL != "" {
				if err := validateWebURL(body.CallbackURL); err != nil {
					helpers.WriteError(w, http.StatusBadRequest, "Invalid callbackUrl: "+err.Error())
					return
				}
				req.CallbackURL = sql.NullString{String: body.CallbackURL, Valid: true}
			}
			if len(body.ClientIDs) > 0 {
				req.RequestedClientIDs = model.ScopesString(body.ClientIDs) // reuse comma-join
			}
			if body.SuggestedMonthlyLimit != nil {
				req.SuggestedMonthlyLimit = sql.NullInt64{Int64: *body.SuggestedMonthlyLimit, Valid: true}
			}
			if body.SuggestedExpiry != "" {
				req.SuggestedExpiry = sql.NullString{String: body.SuggestedExpiry, Valid: true}
			}

			if err := db.KeyRequestStore().Create(r.Context(), req); err != nil {
				log.Error().Err(err).Msg("Failed to create key request")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to create key request")
				return
			}

			approvalURL := cfg.FrontendURL + "/approve/" + code

			helpers.WriteJSON(w, http.StatusCreated, map[string]interface{}{
				"code":        code,
				"approvalUrl": approvalURL,
				"expiresIn":   600,
				"expiresAt":   expiresAt.Format(time.RFC3339),
			})
		})

		// GET /auth/key-request/:code/status — Poll for device flow (NO AUTH, RATE LIMITED)
		r.With(middleware.KeyRequestRateLimiter.Middleware).Get("/{code}/status", func(w http.ResponseWriter, r *http.Request) {
			code := chi.URLParam(r, "code")
			ctx := r.Context()

			req, err := db.KeyRequestStore().FindByCode(ctx, code)
			if err != nil || req == nil {
				helpers.WriteError(w, http.StatusNotFound, "Key request not found")
				return
			}

			if req.ExpiresAt.Time.Before(time.Now()) && req.Status == "pending" {
				helpers.WriteJSON(w, http.StatusGone, map[string]interface{}{
					"status": "expired",
				})
				return
			}

			switch req.Status {
			case "pending":
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "pending"})
			case "approved":
				// Return the key and mark as exchanged
				key, _ := db.ApiKeyStore().FindByID(ctx, req.ApprovedKeyID.Int64)
				if key == nil {
					helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "approved"})
					return
				}

				db.KeyRequestStore().UpdateStatus(ctx, req.ID, "exchanged", nil)

				helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{
					"status":    "approved",
					"apiKey":    key.Key,
					"scopes":    key.GetScopes(),
					"clientIds": key.GetScopedClientIDs(),
				})
			case "exchanged":
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"status":  "exchanged",
					"message": "API key has already been retrieved",
				})
			case "denied":
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "denied"})
			default:
				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"status": req.Status})
			}
		})

		// POST /auth/key-request/exchange — Web flow code exchange (NO AUTH)
		r.With(middleware.KeyRequestRateLimiter.Middleware).Post("/exchange", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Code string `json:"code"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
				helpers.WriteError(w, http.StatusBadRequest, "Code is required")
				return
			}

			ctx := r.Context()

			// Find by exchange code
			// We need to iterate — for now we search by the exchange code
			// This is a simple approach; in production you'd index this
			req, err := db.KeyRequestStore().FindByCode(ctx, body.Code)
			if err != nil || req == nil || !req.ExchangeCode.Valid || req.ExchangeCode.String != body.Code {
				// Try finding by exchange code via status scan
				helpers.WriteError(w, http.StatusNotFound, "Invalid exchange code")
				return
			}

			if req.Status != "approved" {
				helpers.WriteError(w, http.StatusGone, "Exchange code has already been used or request was denied")
				return
			}

			key, _ := db.ApiKeyStore().FindByID(ctx, req.ApprovedKeyID.Int64)
			if key == nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to retrieve API key")
				return
			}

			db.KeyRequestStore().UpdateStatus(ctx, req.ID, "exchanged", nil)

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"apiKey":    key.Key,
				"scopes":   key.GetScopes(),
				"clientIds": key.GetScopedClientIDs(),
			})
		})

		// Authenticated routes for approval
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(db, nil))

			// GET /auth/key-request/:code — Get details for approval page
			r.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
				code := chi.URLParam(r, "code")
				req, err := db.KeyRequestStore().FindByCode(r.Context(), code)
				if err != nil || req == nil {
					helpers.WriteError(w, http.StatusNotFound, "Key request not found")
					return
				}

				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"code":               req.Code,
					"appName":            req.AppName,
					"appDescription":     req.AppDescription,
					"appUrl":             req.AppURL.String,
					"requestedScopes":    model.ParseScopes(req.RequestedScopes),
					"requestedClientIds": model.ParseScopes(req.RequestedClientIDs),
					"callbackUrl":        req.CallbackURL.String,
					"suggestedMonthlyLimit": req.SuggestedMonthlyLimit.Int64,
					"suggestedExpiry":    req.SuggestedExpiry.String,
					"status":            req.Status,
					"expiresAt":         req.ExpiresAt,
				})
			})

			// POST /auth/key-request/:code/approve — Approve with selected scopes/clients
			r.Post("/{code}/approve", func(w http.ResponseWriter, r *http.Request) {
				reqCtx := helpers.GetRequestContext(r)
				if reqCtx == nil || reqCtx.ScopedKey != nil {
					helpers.WriteError(w, http.StatusForbidden, "Use your master API key to approve key requests.")
					return
				}

				user, ok := reqCtx.User.(*model.User)
				if !ok || user == nil {
					helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
					return
				}

				code := chi.URLParam(r, "code")
				ctx := r.Context()

				keyReq, err := db.KeyRequestStore().FindByCode(ctx, code)
				if err != nil || keyReq == nil {
					helpers.WriteError(w, http.StatusNotFound, "Key request not found")
					return
				}

				if keyReq.Status != "pending" {
					helpers.WriteError(w, http.StatusBadRequest, "Key request is not pending")
					return
				}

				if keyReq.ExpiresAt.Time.Before(time.Now()) {
					helpers.WriteError(w, http.StatusGone, "Key request has expired")
					return
				}

				// Parse approval body
				var body struct {
					Scopes     []string `json:"scopes"`
					ClientIDs  []string `json:"clientIds"`
					MonthlyLimit *int64   `json:"monthlyLimit"`
					ExpiresAt  string   `json:"expiresAt"`
				}
				json.NewDecoder(r.Body).Decode(&body)

				// Use requested scopes if none specified in approval
				scopes := body.Scopes
				if len(scopes) == 0 {
					scopes = model.ParseScopes(keyReq.RequestedScopes)
				}

				// Validate scopes — approver can only select from what was requested
				requestedScopes := model.ParseScopes(keyReq.RequestedScopes)
				for _, s := range scopes {
					if !model.HasScope(requestedScopes, s) {
						helpers.WriteError(w, http.StatusBadRequest, "Cannot grant scope not in original request: "+s)
						return
					}
				}

				// Create the scoped API key
				keyStr, err := model.GenerateAPIKey()
				if err != nil {
					helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate API key")
					return
				}

				apiKey := &model.ApiKey{
					UserID:  user.ID,
					Key:     keyStr,
					Name:    keyReq.AppName,
					Scopes:  model.ScopesString(scopes),
					Enabled: true,
				}

				if len(body.ClientIDs) > 0 {
					apiKey.ScopedClientIDs = sql.NullString{
						String: model.ScopesString(body.ClientIDs),
						Valid:  true,
					}
				}

				if body.MonthlyLimit != nil {
					apiKey.MonthlyLimit = sql.NullInt64{Int64: *body.MonthlyLimit, Valid: true}
				}

				if body.ExpiresAt != "" {
					if t, err := time.Parse(time.RFC3339, body.ExpiresAt); err == nil {
						apiKey.ExpiresAt = &model.SQLiteTime{Time: t, Valid: true}
					}
				}

				if err := db.ApiKeyStore().Create(ctx, apiKey); err != nil {
					helpers.WriteError(w, http.StatusInternalServerError, "Failed to create API key")
					return
				}

				// Update key request status
				updates := map[string]interface{}{
					"approvedKeyId": apiKey.ID,
					"approvedById":  user.ID,
				}

				// For web flow (has callback URL), generate exchange code
				if keyReq.CallbackURL.Valid && keyReq.CallbackURL.String != "" {
					exchangeCode, _ := generateExchangeCode()
					updates["exchangeCode"] = exchangeCode
				}

				db.KeyRequestStore().UpdateStatus(ctx, keyReq.ID, "approved", updates)

				response := map[string]interface{}{
					"success": true,
					"keyId":   apiKey.ID,
					"scopes":  apiKey.GetScopes(),
				}

				if ec, ok := updates["exchangeCode"].(string); ok {
					response["exchangeCode"] = ec
				}

				helpers.WriteJSON(w, http.StatusOK, response)

				log.Info().
					Str("code", code).
					Str("appName", keyReq.AppName).
					Int64("userId", user.ID).
					Int64("keyId", apiKey.ID).
					Msg("Key request approved")
			})

			// POST /auth/key-request/:code/deny
			r.Post("/{code}/deny", func(w http.ResponseWriter, r *http.Request) {
				reqCtx := helpers.GetRequestContext(r)
				if reqCtx == nil || reqCtx.ScopedKey != nil {
					helpers.WriteError(w, http.StatusForbidden, "Use your master API key.")
					return
				}

				code := chi.URLParam(r, "code")
				keyReq, err := db.KeyRequestStore().FindByCode(r.Context(), code)
				if err != nil || keyReq == nil {
					helpers.WriteError(w, http.StatusNotFound, "Key request not found")
					return
				}

				if keyReq.Status != "pending" {
					helpers.WriteError(w, http.StatusBadRequest, "Key request is not pending")
					return
				}

				db.KeyRequestStore().UpdateStatus(r.Context(), keyReq.ID, "denied", nil)

				helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
					"success": true,
					"message": "Key request denied",
				})
			})
		})
	})
}

// validateWebURL ensures a URL is a valid http/https URL with a public host.
// Rejects private/loopback IPs to prevent SSRF.
func validateWebURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return fmt.Errorf("invalid URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("URL must have a non-empty host")
	}
	// Reject loopback and RFC1918 private IPs to prevent SSRF
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
			return fmt.Errorf("URL must not target a private or loopback address")
		}
	}
	return nil
}

func generateKeyRequestCode(length int) (string, error) {
	return generatePairingCode(length) // Reuse the same code generator
}

func generateExchangeCode() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

