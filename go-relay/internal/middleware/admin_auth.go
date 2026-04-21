package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	// AdminCookieNameSecure is used when serving over HTTPS (production).
	// The __Host- prefix gives browsers extra guarantees but requires Secure=true and Path=/.
	AdminCookieNameSecure   = "__Host-admin_token"
	// AdminCookieNameInsecure is used in development (HTTP). The __Host- prefix
	// is silently rejected by browsers when Secure is false, so we use a plain name.
	AdminCookieNameInsecure = "admin_token"
	AdminCSRFHeader         = "X-CSRF-Token"
	adminCookiePath         = "/" // __Host- requires Path=/
	DefaultAccessTTL        = 15 * time.Minute
	defaultRefreshWindow    = 7*time.Minute + 30*time.Second
)

// adminCookieName picks the right cookie name based on whether secure cookies are enabled.
func adminCookieName(secure bool) string {
	if secure {
		return AdminCookieNameSecure
	}
	return AdminCookieNameInsecure
}

// AdminCookieName returns the cookie name appropriate for the given config.
// Used by the test suite and any other code that needs to read the cookie.
func AdminCookieName(cfg *AdminAuthConfig) string {
	return adminCookieName(cfg.SecureCookies)
}

// AdminClaims is the JWT claim payload for admin sessions.
type AdminClaims struct {
	Email        string `json:"email"`
	Role         string `json:"role"`
	SessionStart int64  `json:"session_start"` // unix seconds when session began
	CSRFToken    string `json:"csrf"`
	jwt.RegisteredClaims
}

// adminCtxKey is the context key for the authenticated admin user.
type adminCtxKey string

const AdminContextKey adminCtxKey = "adminUser"

// GetAdminUser extracts the authenticated admin from request context.
func GetAdminUser(r *http.Request) *model.User {
	if v := r.Context().Value(AdminContextKey); v != nil {
		if u, ok := v.(*model.User); ok {
			return u
		}
	}
	return nil
}

// AdminAuthConfig holds configuration for admin auth middleware.
type AdminAuthConfig struct {
	JWTSecret     []byte
	SessionMaxHrs int
	TokenTTL      time.Duration // JWT lifetime per issuance; shorter in production for sliding-refresh security
	IsDevelopment bool          // true in dev (APP_ENV != production): longer token TTL, no max-session cutoff on sliding refresh
	SecureCookies bool          // true: Secure cookie flag + __Host- prefix (requires HTTPS); set via ADMIN_SECURE_COOKIES
}

// generateJTI returns a random JWT ID.
func generateJTI() string {
	return uuid.New().String()
}

// generateCSRFToken returns a random CSRF token.
func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// IssueAdminJWT creates a signed JWT for the given admin user, sets the cookie,
// and returns the CSRF token (which the caller should include in the response body).
func IssueAdminJWT(w http.ResponseWriter, user *model.User, cfg *AdminAuthConfig) (string, error) {
	csrf, err := generateCSRFToken()
	if err != nil {
		return "", err
	}
	now := time.Now()
	ttl := cfg.TokenTTL
	if ttl <= 0 {
		ttl = DefaultAccessTTL
	}
	claims := AdminClaims{
		Email:        user.Email,
		Role:         user.Role,
		SessionStart: now.Unix(),
		CSRFToken:    csrf,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userIDString(user.ID),
			ID:        generateJTI(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(cfg.JWTSecret)
	if err != nil {
		return "", err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     adminCookieName(cfg.SecureCookies),
		Value:    signed,
		Path:     adminCookiePath,
		HttpOnly: true,
		Secure:   cfg.SecureCookies,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(ttl.Seconds()),
	})
	return csrf, nil
}

// ClearAdminCookie removes the admin auth cookie.
func ClearAdminCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     adminCookieName(secure),
		Value:    "",
		Path:     adminCookiePath,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

// RequireAdmin returns middleware that validates the admin JWT cookie,
// loads the user, checks role, validates CSRF on mutating requests,
// and applies sliding refresh.
func RequireAdmin(db *database.DB, cfg *AdminAuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookie, err := r.Cookie(adminCookieName(cfg.SecureCookies))
			if err != nil || cookie.Value == "" {
				helpers.WriteError(w, http.StatusUnauthorized, "Admin authentication required")
				return
			}

			claims := &AdminClaims{}
			parsed, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return cfg.JWTSecret, nil
			})
			if err != nil || !parsed.Valid {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid or expired admin session")
				return
			}

			// Check denylist
			denied, err := db.JWTDenylistStore().IsDenied(ctx, claims.ID)
			if err != nil {
				log.Error().Err(err).Msg("admin auth: denylist check failed")
				helpers.WriteError(w, http.StatusInternalServerError, "Authentication error")
				return
			}
			if denied {
				ClearAdminCookie(w, cfg.SecureCookies)
				helpers.WriteError(w, http.StatusUnauthorized, "Session has been revoked")
				return
			}

			// Check session age against max session window
			sessionStart := time.Unix(claims.SessionStart, 0)
			maxSessionAge := time.Duration(cfg.SessionMaxHrs) * time.Hour
			if time.Since(sessionStart) > maxSessionAge {
				ClearAdminCookie(w, cfg.SecureCookies)
				helpers.WriteError(w, http.StatusUnauthorized, "Session expired, please log in again")
				return
			}

			// Load user
			userID, err := parseUserID(claims.Subject)
			if err != nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid session")
				return
			}
			user, err := db.UserStore().FindByID(ctx, userID)
			if err != nil || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid session")
				return
			}
			if user.Role != "admin" {
				helpers.WriteError(w, http.StatusForbidden, "Admin access required")
				return
			}
			if user.Disabled {
				helpers.WriteError(w, http.StatusForbidden, "Account disabled")
				return
			}

			// CSRF check on mutating requests
			if isMutating(r.Method) {
				headerToken := r.Header.Get(AdminCSRFHeader)
				if headerToken == "" || headerToken != claims.CSRFToken {
					helpers.WriteError(w, http.StatusForbidden, "CSRF token missing or invalid")
					return
				}
			}

			// Sliding refresh: if token has < refresh window remaining
			// AND total session age has not exceeded max, issue a new cookie.
			expiry := claims.ExpiresAt.Time
			if time.Until(expiry) < defaultRefreshWindow && time.Since(sessionStart) < maxSessionAge {
				if err := refreshAdminJWT(w, claims, cfg); err != nil {
					log.Warn().Err(err).Msg("admin auth: sliding refresh failed")
				}
			}

			backfillAccessLog(r, user.ID, "")
			rCtx := context.WithValue(r.Context(), AdminContextKey, user)
			next.ServeHTTP(w, r.WithContext(rCtx))
		})
	}
}

// refreshAdminJWT issues a new JWT preserving the original session_start and csrf token.
// This avoids forcing a CSRF rotation in the middle of a session.
func refreshAdminJWT(w http.ResponseWriter, oldClaims *AdminClaims, cfg *AdminAuthConfig) error {
	ttl := cfg.TokenTTL
	if ttl <= 0 {
		ttl = DefaultAccessTTL
	}
	now := time.Now()
	newClaims := AdminClaims{
		Email:        oldClaims.Email,
		Role:         oldClaims.Role,
		SessionStart: oldClaims.SessionStart,
		CSRFToken:    oldClaims.CSRFToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   oldClaims.Subject,
			ID:        generateJTI(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signed, err := token.SignedString(cfg.JWTSecret)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     adminCookieName(cfg.SecureCookies),
		Value:    signed,
		Path:     adminCookiePath,
		HttpOnly: true,
		Secure:   cfg.SecureCookies,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(ttl.Seconds()),
	})
	return nil
}

// DenyAdminJWT extracts the jti from the cookie and adds it to the denylist.
// Used by logout to invalidate the current session immediately.
func DenyAdminJWT(r *http.Request, db *database.DB, cfg *AdminAuthConfig) error {
	cookie, err := r.Cookie(adminCookieName(cfg.SecureCookies))
	if err != nil || cookie.Value == "" {
		return nil
	}
	claims := &AdminClaims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return cfg.JWTSecret, nil
	})
	if err != nil {
		return nil // can't parse — already invalid
	}
	if claims.ID == "" || claims.ExpiresAt == nil {
		return nil
	}
	return db.JWTDenylistStore().Add(r.Context(), claims.ID, claims.ExpiresAt.Time)
}

// ParseAdminCookie parses and validates the admin JWT cookie without enforcing
// admin role or context loading. Used by /admin/auth/me to fetch session info.
func ParseAdminCookie(r *http.Request, cfg *AdminAuthConfig) (*AdminClaims, error) {
	cookie, err := r.Cookie(adminCookieName(cfg.SecureCookies))
	if err != nil || cookie.Value == "" {
		return nil, errors.New("no admin cookie")
	}
	claims := &AdminClaims{}
	parsed, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return cfg.JWTSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func isMutating(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

func userIDString(id int64) string {
	return strconv.FormatInt(id, 10)
}

func parseUserID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
