package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/rs/zerolog/log"
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

// RateLimiter provides IP-based rate limiting.
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	window   time.Duration
	max      int
	message  string
}

// NewRateLimiter creates a new IP-based rate limiter.
func NewRateLimiter(window time.Duration, max int, message string) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		max:      max,
		message:  message,
	}
	// Cleanup goroutine
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()
	return rl
}

// Middleware returns an HTTP middleware that enforces the rate limit.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)

		// Filter old entries
		var recent []time.Time
		for _, t := range rl.requests[ip] {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}

		if len(recent) >= rl.max {
			rl.mu.Unlock()
			log.Warn().Str("ip", ip).Str("path", r.URL.Path).Msg("Rate limit exceeded")
			helpers.WriteJSON(w, http.StatusTooManyRequests, map[string]string{
				"error": rl.message,
			})
			return
		}

		recent = append(recent, now)
		rl.requests[ip] = recent
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	cutoff := time.Now().Add(-rl.window)
	for ip, times := range rl.requests {
		var recent []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}
		if len(recent) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = recent
		}
	}
}

// Pre-configured rate limiters matching the Node.js implementation.
// Limits are configurable via environment variables for testing.
// These are initialized with defaults and re-initialized by InitRateLimiters()
// after .env is loaded in main().
var (
	AuthRateLimiter             *RateLimiter
	PasswordResetRateLimiter    *RateLimiter
	AccountManagementRateLimiter *RateLimiter
)

func init() {
	// Default initialization (before .env is loaded)
	AuthRateLimiter = NewRateLimiter(15*time.Minute, 5,
		"Too many authentication attempts from this IP, please try again after 15 minutes")
	PasswordResetRateLimiter = NewRateLimiter(1*time.Hour, 3,
		"Too many password reset attempts from this IP, please try again later")
	AccountManagementRateLimiter = NewRateLimiter(1*time.Hour, 10,
		"Too many account management requests from this IP, please try again later")
}

// InitRateLimiters re-initializes rate limiters after .env has been loaded.
// Call this from main() after godotenv.Load().
func InitRateLimiters() {
	AuthRateLimiter = NewRateLimiter(15*time.Minute, getEnvInt("AUTH_RATE_LIMIT", 5),
		"Too many authentication attempts from this IP, please try again after 15 minutes")
	PasswordResetRateLimiter = NewRateLimiter(1*time.Hour, getEnvInt("PASSWORD_RESET_RATE_LIMIT", 3),
		"Too many password reset attempts from this IP, please try again later")
	AccountManagementRateLimiter = NewRateLimiter(1*time.Hour, getEnvInt("ACCOUNT_MGMT_RATE_LIMIT", 10),
		"Too many account management requests from this IP, please try again later")
}
