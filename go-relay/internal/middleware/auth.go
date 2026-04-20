package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/alerts"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/counter"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/metrics"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/rs/zerolog/log"
)

// RequestCounts is the package-level batch counter for API request tracking.
// main() wires a goroutine that periodically flushes it to the DB.
var RequestCounts = counter.New()

// keyPrefix returns the first 8 characters of a raw API key, safe to log.
func keyPrefix(key string) string {
	if len(key) <= 8 {
		return key
	}
	return key[:8]
}

// authCache caches API key lookups to avoid hitting SQLite on every request.
type authCacheEntry struct {
	user      *model.User
	scopedKey *model.ApiKey
	expiresAt time.Time
}

var (
	authCacheMu  sync.RWMutex
	authCacheMap = make(map[string]*authCacheEntry)
	authCacheTTL = 60 * time.Second
)

// SetAuthCacheTTL configures the auth cache TTL. Call once during server init.
func SetAuthCacheTTL(d time.Duration) {
	authCacheMu.Lock()
	authCacheTTL = d
	authCacheMu.Unlock()
}

func getCachedAuth(apiKey string) (*authCacheEntry, bool) {
	authCacheMu.RLock()
	entry, ok := authCacheMap[apiKey]
	authCacheMu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry, true
}

func setCachedAuth(apiKey string, user *model.User, scopedKey *model.ApiKey) {
	authCacheMu.Lock()
	authCacheMap[apiKey] = &authCacheEntry{user: user, scopedKey: scopedKey, expiresAt: time.Now().Add(authCacheTTL)}
	authCacheMu.Unlock()
}

// InvalidateCachedAuth removes a specific API key from the cache.
// Call this after key rotation, deletion, or any other event that should
// invalidate cached auth for a given key.
func InvalidateCachedAuth(apiKey string) {
	authCacheMu.Lock()
	delete(authCacheMap, apiKey)
	authCacheMu.Unlock()
}

// InvalidateCachedAuthForUser removes all cache entries that map to a given user.
// Slower (linear scan) but useful when the API key value is no longer known.
func InvalidateCachedAuthForUser(userID int64) {
	authCacheMu.Lock()
	for k, entry := range authCacheMap {
		if entry.user != nil && entry.user.ID == userID {
			delete(authCacheMap, k)
		}
	}
	authCacheMu.Unlock()
}

// AuthMiddleware validates API keys (master + scoped) and sets request context.
//
// Security note: Master API keys are stored as SHA-256 hashes in Users.apiKeyHash
// (UserStore.FindByAPIKey hashes the input before lookup). Connection tokens are
// stored as SHA-256 hashes in ConnectionTokens.tokenHash. Scoped API keys are
// stored as SHA-256 hashes in ApiKeys.keyHash (ApiKeyStore.FindByKey hashes the
// input). The plaintext key is never stored at rest after the initial create/return.
func AuthMiddleware(db *database.DB, manager *ws.ClientManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try session-token auth first (Authorization: Bearer <token>).
			// Dashboard requests use this path; programmatic API users use
			// x-api-key.
			if user, ok := tryBearerAuth(r, db); ok {
				if user.Disabled {
					helpers.WriteError(w, http.StatusForbidden, "Account disabled")
					return
				}
				if !user.EmailVerified {
					helpers.WriteError(w, http.StatusForbidden, "Please verify your email address before using the API.")
					return
				}
				if user.APIKeyRotationRequired {
					helpers.WriteError(w, http.StatusForbidden, "Your API key must be rotated for security. Please log in to the dashboard and regenerate your master API key.")
					return
				}
				reqCtx := &helpers.RequestContext{
					User:               user,
					MasterAPIKey:       user.APIKeyHash.String,
					SubscriptionStatus: user.GetSubscriptionStatus(),
				}
				backfillAccessLog(r, user.ID, "")
				rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
				next.ServeHTTP(w, r.WithContext(rCtx))
				return
			}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				apiKey = r.Header.Get("x-api-key")
			}
			if apiKey == "" {
				metrics.AuthFailuresTotal.WithLabelValues("missing_key").Inc()
				helpers.WriteError(w, http.StatusUnauthorized, "API key is required")
				return
			}

			ctx := r.Context()

			// Check cache first to avoid SQLite on every request
			if cached, ok := getCachedAuth(apiKey); ok {
				// MasterAPIKey in the request context is now the user's
				// SHA-256 hash (an opaque per-account identifier), NOT the
				// plaintext master key. The plaintext master key never lives
				// outside the registration/regeneration response payload.
				matchKey := ""
				if cached.user != nil && cached.user.APIKeyHash.Valid {
					matchKey = cached.user.APIKeyHash.String
				}

				// Validate clientId even on cache hit
				clientID := r.URL.Query().Get("clientId")
				if clientID != "" && manager != nil {
					client := manager.GetClient(clientID)
					if client == nil {
						helpers.WriteError(w, http.StatusNotFound, "Invalid client ID")
						return
					}
					if client.APIKey() != matchKey {
						helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key for this client ID")
						return
					}
				}

				if cached.user != nil && !cached.user.EmailVerified {
					helpers.WriteError(w, http.StatusForbidden, "Please verify your email address before using the API.")
					return
				}

				if cached.user != nil && cached.user.Disabled {
					helpers.WriteError(w, http.StatusForbidden, "Account disabled")
					return
				}

				if cached.user != nil && cached.user.APIKeyRotationRequired {
					helpers.WriteError(w, http.StatusForbidden, "Your API key must be rotated for security. Please log in to the dashboard and regenerate your master API key.")
					return
				}

				if cached.user != nil && cached.scopedKey == nil {
					reqCtx := &helpers.RequestContext{
						User: cached.user, MasterAPIKey: matchKey,
						SubscriptionStatus: cached.user.GetSubscriptionStatus(),
					}
					backfillAccessLog(r, cached.user.ID, keyPrefix(apiKey))
					rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
					next.ServeHTTP(w, r.WithContext(rCtx))
					return
				} else if cached.scopedKey != nil && cached.user != nil {
					kp := keyPrefix(cached.scopedKey.Key)
					reqCtx := &helpers.RequestContext{
						User: cached.user, MasterAPIKey: matchKey,
						SubscriptionStatus: cached.user.GetSubscriptionStatus(),
						ScopedKey: &helpers.ScopedKeyInfo{
							ID:                cached.scopedKey.ID,
							KeyPrefix:         kp,
							ScopedClientID:    cached.scopedKey.ScopedClientID.String,
							ScopedClientIDs:   cached.scopedKey.GetScopedClientIDs(),
							ScopedUserID:      cached.scopedKey.ScopedUserID.String,
							ScopedUserIDs:     cached.scopedKey.GetScopedUserIDs(),
							Scopes:            cached.scopedKey.GetScopes(),
							MonthlyLimit:      cached.scopedKey.MonthlyLimit.Int64,
							MonthlyLimitSet:   cached.scopedKey.MonthlyLimit.Valid,
							RequestsThisMonth: cached.scopedKey.RequestsThisMonth,
						},
					}
					backfillAccessLog(r, cached.user.ID, kp)
					rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
					next.ServeHTTP(w, r.WithContext(rCtx))
					return
				}
			}

			// 1. Try master key lookup
			user, err := db.UserStore().FindByAPIKey(ctx, apiKey)
			if err != nil {
				log.Error().Err(err).Str("method", r.Method).Str("path", r.URL.Path).Msg("Auth error during user lookup")
				helpers.WriteError(w, http.StatusInternalServerError, "Authentication error")
				return
			}

			if user != nil {
				// Check email verification
				if !user.EmailVerified {
					helpers.WriteError(w, http.StatusForbidden, "Please verify your email address before using the API.")
					return
				}

				if user.Disabled {
					helpers.WriteError(w, http.StatusForbidden, "Account disabled")
					return
				}

				if user.APIKeyRotationRequired {
					helpers.WriteError(w, http.StatusForbidden, "Your API key must be rotated for security. Please log in to the dashboard and regenerate your master API key.")
					return
				}

				// Master key auth — clients are registered in ClientManager
				// under the user's apiKeyHash (the per-account identifier).
				masterIdentifier := user.APIKeyHash.String
				clientID := r.URL.Query().Get("clientId")
				if clientID != "" {
					client := manager.GetClient(clientID)
					if client == nil {
						helpers.WriteError(w, http.StatusNotFound, "Invalid client ID")
						return
					}
					if client.APIKey() != masterIdentifier {
						helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key for this client ID")
						return
					}
				}

				setCachedAuth(apiKey, user, nil) // Cache master key lookup
				reqCtx := &helpers.RequestContext{
					User:               user,
					MasterAPIKey:       masterIdentifier,
					SubscriptionStatus: user.GetSubscriptionStatus(),
				}
				backfillAccessLog(r, user.ID, keyPrefix(apiKey))
				rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
				next.ServeHTTP(w, r.WithContext(rCtx))
				return
			}

			// 2. Try scoped API key lookup
			scopedKey, err := db.ApiKeyStore().FindByKey(ctx, apiKey)
			if err != nil {
				log.Error().Err(err).Str("method", r.Method).Str("path", r.URL.Path).Msg("Auth error during scoped key lookup")
				helpers.WriteError(w, http.StatusInternalServerError, "Authentication error")
				return
			}

			if scopedKey == nil {
				metrics.AuthFailuresTotal.WithLabelValues("bad_key").Inc()
				if alerts.Track("invalid_token", 10, 5*time.Minute, 30*time.Minute) {
					alerts.Fire(alerts.Event{
						Type:     alerts.TypeInvalidTokenSpike,
						Severity: "warning",
						Message:  "10+ requests with invalid/unknown API keys in 5 minutes",
						Details:  map[string]interface{}{"path": r.URL.Path},
					})
				}
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key")
				return
			}

			if !scopedKey.Enabled {
				metrics.AuthFailuresTotal.WithLabelValues("disabled").Inc()
				helpers.WriteError(w, http.StatusUnauthorized, "API key is disabled")
				return
			}

			if scopedKey.IsExpired() {
				metrics.AuthFailuresTotal.WithLabelValues("expired").Inc()
				helpers.WriteError(w, http.StatusUnauthorized, "API key has expired")
				return
			}

			// Look up parent user
			parentUser, err := db.UserStore().FindByID(ctx, scopedKey.UserID)
			if err != nil || parentUser == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key")
				return
			}

			// Check parent user's email verification
			if !parentUser.EmailVerified {
				helpers.WriteError(w, http.StatusForbidden, "Account email not verified. Please verify your email address.")
				return
			}

			if parentUser.APIKeyRotationRequired {
				helpers.WriteError(w, http.StatusForbidden, "Account requires master API key rotation. Please log in to the dashboard and regenerate the master API key.")
				return
			}

			if parentUser.Disabled {
				helpers.WriteError(w, http.StatusForbidden, "Account disabled")
				return
			}

			// Validate clientId against parent account identifier (apiKeyHash)
			parentIdentifier := parentUser.APIKeyHash.String
			clientID := r.URL.Query().Get("clientId")
			if clientID != "" {
				client := manager.GetClient(clientID)
				if client == nil {
					helpers.WriteError(w, http.StatusNotFound, "Invalid client ID")
					return
				}
				if client.APIKey() != parentIdentifier {
					helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key for this client ID")
					return
				}
			}

			kp := keyPrefix(scopedKey.Key)
			setCachedAuth(apiKey, parentUser, scopedKey) // Cache scoped key lookup
			reqCtx := &helpers.RequestContext{
				User:               parentUser,
				MasterAPIKey:       parentIdentifier,
				SubscriptionStatus: parentUser.GetSubscriptionStatus(),
				ScopedKey: &helpers.ScopedKeyInfo{
					ID:                scopedKey.ID,
					KeyPrefix:         kp,
					ScopedClientID:    scopedKey.ScopedClientID.String,
					ScopedClientIDs:   scopedKey.GetScopedClientIDs(),
					ScopedUserID:      scopedKey.ScopedUserID.String,
					ScopedUserIDs:     scopedKey.GetScopedUserIDs(),
					Scopes:            scopedKey.GetScopes(),
					MonthlyLimit:      scopedKey.MonthlyLimit.Int64,
					MonthlyLimitSet:   scopedKey.MonthlyLimit.Valid,
					RequestsThisMonth: scopedKey.RequestsThisMonth,
				},
			}
			backfillAccessLog(r, parentUser.ID, kp)

			if scopedKey.MonthlyLimit.Valid && int64(scopedKey.RequestsThisMonth) >= scopedKey.MonthlyLimit.Int64 {
				helpers.WriteError(w, http.StatusTooManyRequests, "Scoped key monthly request limit exceeded")
				return
			}

			rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
			next.ServeHTTP(w, r.WithContext(rCtx))
		})
	}
}

// TrackAPIUsage increments request counters and enforces monthly limits.
func TrackAPIUsage(freeTierLimit int, db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Get user from context
			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Check limits using the counts we already have from auth lookup
			// (avoids extra DB round-trip on every request)

			// Fire approaching-limit alert (once per 24h window per user).
			// A freeTierLimit of 0 means unlimited — skip all checks.
			if freeTierLimit > 0 && reqCtx.SubscriptionStatus != "active" {
				monthlyThreshold := int(float64(freeTierLimit) * 0.8)
				if monthlyThreshold > 0 && user.RequestsThisMonth >= monthlyThreshold {
					if alerts.Track("monthly:"+fmt.Sprintf("%d", user.ID), 1, 24*time.Hour, 24*time.Hour) {
						alerts.Fire(alerts.Event{
							Type:     alerts.TypeUserMonthlyLimitApproach,
							Severity: "info",
							Message:  "Free-tier user approaching 80% of monthly request limit",
							Details:  map[string]interface{}{"userId": user.ID, "requestsThisMonth": user.RequestsThisMonth, "monthlyLimit": freeTierLimit},
						})
					}
				}
			}

			// Scoped key monthly approach alert — fires at most once per 24h per key.
			if reqCtx.ScopedKey != nil && reqCtx.ScopedKey.MonthlyLimitSet {
				threshold := int(float64(reqCtx.ScopedKey.MonthlyLimit) * 0.8)
				if threshold > 0 && reqCtx.ScopedKey.RequestsThisMonth >= threshold {
					if alerts.Track("scoped_monthly:"+fmt.Sprintf("%d", reqCtx.ScopedKey.ID), 1, 24*time.Hour, 24*time.Hour) {
						alerts.Fire(alerts.Event{
							Type:     alerts.TypeScopedKeyMonthlyLimitApproach,
							Severity: "info",
							Message:  "Scoped API key approaching 80% of monthly request limit",
							Details:  map[string]interface{}{"userId": user.ID, "scopedKeyId": reqCtx.ScopedKey.ID, "keyPrefix": reqCtx.ScopedKey.KeyPrefix, "requestsThisMonth": reqCtx.ScopedKey.RequestsThisMonth, "monthlyLimit": reqCtx.ScopedKey.MonthlyLimit},
						})
					}
				}
			}

			// Monthly limit for free tier. A freeTierLimit of 0 means unlimited.
			if freeTierLimit > 0 && reqCtx.SubscriptionStatus != "active" && user.RequestsThisMonth >= freeTierLimit {
				helpers.WriteJSON(w, http.StatusTooManyRequests, map[string]interface{}{
					"error":      "Monthly API request limit reached",
					"limit":      freeTierLimit,
					"message":    "Please upgrade to a paid subscription for unlimited monthly API access",
					"upgradeUrl": "/api/subscriptions/create-checkout-session",
				})
				return
			}

			// Batch the increment — flushed to DB every 500ms by the counter goroutine in main()
			RequestCounts.Add(user.ID)

			// Increment scoped key monthly counter if this request used a scoped key.
			if reqCtx.ScopedKey != nil && db != nil {
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()
					_ = db.ApiKeyStore().IncrementMonthlyRequests(ctx, reqCtx.ScopedKey.ID)
				}()
			}

			next.ServeHTTP(w, r)
		})
	}
}

// backfillAccessLog populates the access-log holder (injected by MetricsMiddleware)
// with the authenticated user ID and key prefix so the log line shows real values.
func backfillAccessLog(r *http.Request, userID int64, kp string) {
	if f := GetAccessLogFields(r); f != nil {
		f.UserID = userID
		f.KeyPrefix = kp
	}
}

// tryBearerAuth attempts to authenticate the request via Authorization: Bearer
// <session-token>. Returns (user, true) on success, (nil, false) if no Bearer
// header is present or the token is invalid/expired. Used by AuthMiddleware to
// support session-based dashboard auth alongside the legacy x-api-key path.
func tryBearerAuth(r *http.Request, db *database.DB) (*model.User, bool) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, false
	}
	const prefix = "Bearer "
	if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
		return nil, false
	}
	rawToken := header[len(prefix):]
	if rawToken == "" {
		return nil, false
	}

	hash := model.HashSessionToken(rawToken)
	ctx := r.Context()
	session, err := db.SessionStore().FindByTokenHash(ctx, hash)
	if err != nil || session == nil {
		return nil, false
	}

	user, err := db.UserStore().FindByID(ctx, session.UserID)
	if err != nil || user == nil {
		// Session row outlived its user — clean up and reject.
		_ = db.SessionStore().Delete(ctx, session.ID)
		return nil, false
	}

	// Bump last-used asynchronously (best effort)
	go func() {
		_ = db.SessionStore().UpdateLastUsed(context.Background(), session.ID)
	}()

	return user, true
}
