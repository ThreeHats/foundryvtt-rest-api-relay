package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/rs/zerolog/log"
)

// authCache caches API key lookups to avoid hitting SQLite on every request.
type authCacheEntry struct {
	user      *model.User
	scopedKey *model.ApiKey
	cachedAt  time.Time
}

var (
	authCacheMu    sync.RWMutex
	authCacheMap   = make(map[string]*authCacheEntry)
	authCacheTTL   = 5 * time.Second // Short TTL — just enough to avoid repeated lookups in burst traffic
)

func getCachedAuth(apiKey string) (*authCacheEntry, bool) {
	authCacheMu.RLock()
	defer authCacheMu.RUnlock()
	entry, ok := authCacheMap[apiKey]
	if !ok || time.Since(entry.cachedAt) > authCacheTTL {
		return nil, false
	}
	return entry, true
}

func setCachedAuth(apiKey string, user *model.User, scopedKey *model.ApiKey) {
	authCacheMu.Lock()
	authCacheMap[apiKey] = &authCacheEntry{user: user, scopedKey: scopedKey, cachedAt: time.Now()}
	authCacheMu.Unlock()
}

// AuthMiddleware validates API keys (master + scoped) and sets request context.
func AuthMiddleware(db *database.DB, manager *ws.ClientManager, isMemoryStore bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Memory store mode: bypass auth
			if isMemoryStore {
				reqCtx := &helpers.RequestContext{
					MasterAPIKey:       "local-dev",
					SubscriptionStatus: "active",
				}
				ctx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				apiKey = r.Header.Get("x-api-key")
			}
			if apiKey == "" {
				helpers.WriteError(w, http.StatusUnauthorized, "API key is required")
				return
			}

			ctx := r.Context()

			// Check cache first to avoid SQLite on every request
			if cached, ok := getCachedAuth(apiKey); ok {
				// Validate clientId even on cache hit
				clientID := r.URL.Query().Get("clientId")
				if clientID != "" && manager != nil {
					matchKey := apiKey
					if cached.user != nil {
						matchKey = cached.user.APIKey
					}
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

				if cached.user != nil && cached.scopedKey == nil {
					reqCtx := &helpers.RequestContext{
						User: cached.user, MasterAPIKey: cached.user.APIKey,
						SubscriptionStatus: cached.user.GetSubscriptionStatus(),
					}
					rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
					next.ServeHTTP(w, r.WithContext(rCtx))
					return
				} else if cached.scopedKey != nil && cached.user != nil {
					reqCtx := &helpers.RequestContext{
						User: cached.user, MasterAPIKey: cached.user.APIKey,
						SubscriptionStatus: cached.user.GetSubscriptionStatus(),
						ScopedKey: &helpers.ScopedKeyInfo{
							ID: cached.scopedKey.ID,
							ScopedClientID: cached.scopedKey.ScopedClientID.String,
							ScopedUserID: cached.scopedKey.ScopedUserID.String,
						},
					}
					rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
					next.ServeHTTP(w, r.WithContext(rCtx))
					return
				}
			}

			// 1. Try master key lookup
			user, err := db.UserStore().FindByAPIKey(ctx, apiKey)
			if err != nil {
				log.Error().Err(err).Msg("Auth error during user lookup")
				helpers.WriteError(w, http.StatusInternalServerError, "Authentication error")
				return
			}

			if user != nil {
				// Master key auth
				clientID := r.URL.Query().Get("clientId")
				if clientID != "" {
					client := manager.GetClient(clientID)
					if client == nil {
						helpers.WriteError(w, http.StatusNotFound, "Invalid client ID")
						return
					}
					if client.APIKey() != apiKey {
						helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key for this client ID")
						return
					}
				}

				setCachedAuth(apiKey, user, nil) // Cache master key lookup
				reqCtx := &helpers.RequestContext{
					User:               user,
					MasterAPIKey:       user.APIKey,
					SubscriptionStatus: user.GetSubscriptionStatus(),
				}
				rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
				next.ServeHTTP(w, r.WithContext(rCtx))
				return
			}

			// 2. Try scoped API key lookup
			scopedKey, err := db.ApiKeyStore().FindByKey(ctx, apiKey)
			if err != nil {
				log.Error().Err(err).Msg("Auth error during scoped key lookup")
				helpers.WriteError(w, http.StatusInternalServerError, "Authentication error")
				return
			}

			if scopedKey == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key")
				return
			}

			if !scopedKey.Enabled {
				helpers.WriteError(w, http.StatusUnauthorized, "API key is disabled")
				return
			}

			if scopedKey.IsExpired() {
				helpers.WriteError(w, http.StatusUnauthorized, "API key has expired")
				return
			}

			// Look up parent user
			parentUser, err := db.UserStore().FindByID(ctx, scopedKey.UserID)
			if err != nil || parentUser == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key")
				return
			}

			// Validate clientId against parent's master key
			clientID := r.URL.Query().Get("clientId")
			if clientID != "" {
				client := manager.GetClient(clientID)
				if client == nil {
					helpers.WriteError(w, http.StatusNotFound, "Invalid client ID")
					return
				}
				if client.APIKey() != parentUser.APIKey {
					helpers.WriteError(w, http.StatusUnauthorized, "Invalid API key for this client ID")
					return
				}
			}

			setCachedAuth(apiKey, parentUser, scopedKey) // Cache scoped key lookup
			reqCtx := &helpers.RequestContext{
				User:               parentUser,
				MasterAPIKey:       parentUser.APIKey,
				SubscriptionStatus: parentUser.GetSubscriptionStatus(),
				ScopedKey: &helpers.ScopedKeyInfo{
					ID:             scopedKey.ID,
					ScopedClientID: scopedKey.ScopedClientID.String,
					ScopedUserID:   scopedKey.ScopedUserID.String,
				},
			}

			if scopedKey.DailyLimit.Valid {
				// Store daily limit for later use (rate limiting in trackApiUsage)
				// We pass it through ScopedKeyInfo
			}

			rCtx := context.WithValue(r.Context(), helpers.RequestContextKey, reqCtx)
			next.ServeHTTP(w, r.WithContext(rCtx))
		})
	}
}

// TrackAPIUsage increments request counters and enforces rate limits.
func TrackAPIUsage(db *database.DB, freeTierLimit, dailyLimit int, isMemoryStore bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isMemoryStore {
				next.ServeHTTP(w, r)
				return
			}

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
			if user.RequestsToday >= dailyLimit {
				tomorrow := time.Now().AddDate(0, 0, 1)
				tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
				helpers.WriteJSON(w, http.StatusTooManyRequests, map[string]interface{}{
					"error":           "Daily API request limit reached",
					"dailyLimit":      dailyLimit,
					"currentRequests": user.RequestsToday,
					"message":         "You have reached the daily limit. Please try again tomorrow.",
					"resetsAt":        tomorrow.Format(time.RFC3339),
				})
				return
			}

			// Monthly limit for free tier
			if reqCtx.SubscriptionStatus != "active" && user.RequestsThisMonth >= freeTierLimit {
				helpers.WriteJSON(w, http.StatusTooManyRequests, map[string]interface{}{
					"error":      "Monthly API request limit reached",
					"limit":      freeTierLimit,
					"message":    "Please upgrade to a paid subscription for unlimited monthly API access",
					"upgradeUrl": "/api/subscriptions/create-checkout-session",
				})
				return
			}

			// Increment counters asynchronously — never blocks the request
			userID := user.ID
			go func() {
				if err := db.UserStore().IncrementRequests(context.Background(), userID); err != nil {
					log.Warn().Err(err).Msg("Failed to increment request count")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
