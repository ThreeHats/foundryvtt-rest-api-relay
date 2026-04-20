package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/alerts"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// getEnvInt reads an integer from an environment variable or returns a fallback.
func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

// keyState holds a token-bucket limiter and a last-seen timestamp for cleanup.
type keyState struct {
	lim      *rate.Limiter
	lastSeen time.Time
}

// RateLimiter provides IP- or API-key-based rate limiting using a token bucket
// per key. The token bucket algorithm is O(1) per request: a map lookup plus
// a single atomic compare-and-swap inside rate.Limiter.Allow().
type RateLimiter struct {
	mu      sync.Mutex
	keys    map[string]*keyState
	r       rate.Limit // tokens per second
	burst   int        // max burst (= window max)
	message string
}

// NewRateLimiter creates a new rate limiter.
//
// window and max define the effective rate: max requests per window.
// Internally this is converted to a token-bucket with rate = max/window and
// burst = max so the steady-state throughput matches the sliding-window intent.
func NewRateLimiter(ctx context.Context, window time.Duration, max int, message string) *RateLimiter {
	var r rate.Limit
	if max <= 0 {
		r = rate.Inf // unlimited — Middleware/APIKeyMiddleware short-circuit anyway
	} else {
		r = rate.Every(window / time.Duration(max))
	}
	rl := &RateLimiter{
		keys:    make(map[string]*keyState),
		r:       r,
		burst:   max,
		message: message,
	}
	go func() {
		// Evict stale entries every window so the map doesn't grow unbounded.
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rl.cleanup(window)
			case <-ctx.Done():
				return
			}
		}
	}()
	return rl
}

// allow checks and consumes one token for key. It creates a new limiter on
// first access. The lock is held only for the map lookup + Allow() call.
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	ks, ok := rl.keys[key]
	if !ok {
		ks = &keyState{lim: rate.NewLimiter(rl.r, rl.burst)}
		rl.keys[key] = ks
	}
	ks.lastSeen = time.Now()
	allowed := ks.lim.Allow()
	rl.mu.Unlock()
	return allowed
}

func (rl *RateLimiter) cleanup(window time.Duration) {
	cutoff := time.Now().Add(-window)
	rl.mu.Lock()
	for key, ks := range rl.keys {
		if ks.lastSeen.Before(cutoff) {
			delete(rl.keys, key)
		}
	}
	rl.mu.Unlock()
}

// Middleware returns an HTTP middleware that enforces the rate limit per remote IP.
// A max of 0 or less disables the limiter (passes all requests through).
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	if rl.burst <= 0 {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !rl.allow(ip) {
			log.Warn().Str("ip", ip).Str("path", r.URL.Path).Msg("Rate limit exceeded")
			helpers.WriteJSON(w, http.StatusTooManyRequests, map[string]string{
				"error": rl.message,
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// APIKeyMiddleware returns an HTTP middleware that enforces rate limits per API key.
// A max of 0 or less disables the limiter (passes all requests through).
func (rl *RateLimiter) APIKeyMiddleware(next http.Handler) http.Handler {
	if rl.burst <= 0 {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = r.Header.Get("x-api-key")
		}
		if apiKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !rl.allow(apiKey) {
			log.Warn().Str("path", r.URL.Path).Msg("API key rate limit exceeded")
			if alerts.Track("ratelimit:"+apiKey, 10, 5*time.Minute, 30*time.Minute) {
				alerts.Fire(alerts.Event{
					Type:     alerts.TypeRateLimitBurst,
					Severity: "warning",
					Message:  "API key hit rate limit 10+ times in 5 minutes",
					Details:  map[string]interface{}{"path": r.URL.Path},
				})
			}
			helpers.WriteJSON(w, http.StatusTooManyRequests, map[string]string{
				"error": rl.message,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Pre-configured rate limiters matching the Node.js implementation.
// Limits are configurable via environment variables for testing.
// These are initialized with defaults and re-initialized by InitRateLimiters()
// after .env is loaded in main().
var (
	AuthRateLimiter              *RateLimiter
	PasswordResetRateLimiter     *RateLimiter
	AccountManagementRateLimiter *RateLimiter
	APIKeyRateLimiter            *RateLimiter
	PairingRateLimiter           *RateLimiter
	KeyRequestRateLimiter        *RateLimiter
	// ProbeRateLimiter throttles the public /api/clients/{id}/active probe.
	// 60 requests per minute per IP is well above any legitimate need (the
	// Foundry module makes this call at most once per page load) but blocks
	// scripted clientId enumeration / liveness sweeps.
	ProbeRateLimiter *RateLimiter
)

func init() {
	// Default initialization (before .env is loaded)
	ctx := context.Background()
	AuthRateLimiter = NewRateLimiter(ctx, 15*time.Minute, 5,
		"Too many authentication attempts from this IP, please try again after 15 minutes")
	PasswordResetRateLimiter = NewRateLimiter(ctx, 1*time.Hour, 3,
		"Too many password reset attempts from this IP, please try again later")
	AccountManagementRateLimiter = NewRateLimiter(ctx, 1*time.Hour, 10,
		"Too many account management requests from this IP, please try again later")
	APIKeyRateLimiter = NewRateLimiter(ctx, 1*time.Minute, 300,
		"Too many API requests, please slow down")
	PairingRateLimiter = NewRateLimiter(ctx, 15*time.Minute, 10,
		"Too many pairing attempts from this IP, please try again later")
	KeyRequestRateLimiter = NewRateLimiter(ctx, 15*time.Minute, 20,
		"Too many key request attempts from this IP, please try again later")
	ProbeRateLimiter = NewRateLimiter(ctx, 1*time.Minute, 60,
		"Too many active-connection probes from this IP, please slow down")
}

// InitRateLimiters re-initializes rate limiters after .env has been loaded.
// Call this from main() after godotenv.Load().
func InitRateLimiters(ctx context.Context) {
	AuthRateLimiter = NewRateLimiter(ctx, 15*time.Minute, getEnvInt("AUTH_RATE_LIMIT", 5),
		"Too many authentication attempts from this IP, please try again after 15 minutes")
	PasswordResetRateLimiter = NewRateLimiter(ctx, 1*time.Hour, getEnvInt("PASSWORD_RESET_RATE_LIMIT", 3),
		"Too many password reset attempts from this IP, please try again later")
	AccountManagementRateLimiter = NewRateLimiter(ctx, 1*time.Hour, getEnvInt("ACCOUNT_MGMT_RATE_LIMIT", 10),
		"Too many account management requests from this IP, please try again later")
	APIKeyRateLimiter = NewRateLimiter(ctx, 1*time.Minute, getEnvIntWithFallback("PER_MINUTE_REQUEST_LIMIT", "API_KEY_RATE_LIMIT", 300),
		"Too many API requests, please slow down")
	PairingRateLimiter = NewRateLimiter(ctx, 15*time.Minute, getEnvInt("PAIRING_RATE_LIMIT", 10),
		"Too many pairing attempts from this IP, please try again later")
	KeyRequestRateLimiter = NewRateLimiter(ctx, 15*time.Minute, getEnvInt("KEY_REQUEST_RATE_LIMIT", 20),
		"Too many key request attempts from this IP, please try again later")
	ProbeRateLimiter = NewRateLimiter(ctx, 1*time.Minute, getEnvInt("PROBE_RATE_LIMIT", 60),
		"Too many active-connection probes from this IP, please slow down")
}

// getEnvIntWithFallback tries the primary key, then the fallback key, then the default.
func getEnvIntWithFallback(key, fallbackKey string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	if val := os.Getenv(fallbackKey); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
