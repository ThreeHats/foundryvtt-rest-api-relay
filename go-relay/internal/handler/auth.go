package handler

import (
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/alerts"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v78"
	stripeSubscription "github.com/stripe/stripe-go/v78/subscription"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

var (
	upperCase = regexp.MustCompile(`[A-Z]`)
	lowerCase = regexp.MustCompile(`[a-z]`)
	digitRe   = regexp.MustCompile(`[0-9]`)
)

func validatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	if !upperCase.MatchString(password) {
		return false, "Password must contain at least one uppercase letter"
	}
	if !lowerCase.MatchString(password) {
		return false, "Password must contain at least one lowercase letter"
	}
	if !digitRe.MatchString(password) {
		return false, "Password must contain at least one number"
	}
	return true, ""
}

func parseBody(r *http.Request) (map[string]interface{}, error) {
	var body map[string]interface{}
	if r.Body != nil {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		if len(data) > 0 {
			if err := json.Unmarshal(data, &body); err != nil {
				return nil, fmt.Errorf("invalid JSON: %w", err)
			}
		}
	}
	return body, nil
}

func bodyStr(body map[string]interface{}, key string) string {
	v, _ := body[key].(string)
	return v
}

// AuthRouter creates the auth route group.
func AuthRouter(db *database.DB, cfg *config.Config, manager *ws.ClientManager) chi.Router {
	r := chi.NewRouter()

	// POST /auth/register
	r.With(middleware.AuthRateLimiter.Middleware).Post("/register", func(w http.ResponseWriter, r *http.Request) {
		if cfg.DisableRegistration {
			helpers.WriteError(w, http.StatusForbidden, "Registration is disabled on this server")
			return
		}

		body, err := parseBody(r)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		email := bodyStr(body, "email")
		password := bodyStr(body, "password")

		if email == "" || password == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		if msg := helpers.ValidateEmailForRegistration(email); msg != "" {
			helpers.WriteError(w, http.StatusBadRequest, msg)
			return
		}

		if valid, errMsg := validatePassword(password); !valid {
			helpers.WriteError(w, http.StatusBadRequest, errMsg)
			return
		}

		ctx := r.Context()
		existing, _ := db.UserStore().FindByEmail(ctx, email)
		if existing != nil {
			helpers.WriteError(w, http.StatusConflict, "User already exists")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Registration failed")
			return
		}

		apiKey, err := model.GenerateAPIKey()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate API key")
			helpers.WriteError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		user := &model.User{
			Email:              email,
			Password:           string(hash),
			APIKey:             apiKey,
			EmailVerified:      true, // default: verified (overridden below if SMTP configured)
			SubscriptionStatus: sql.NullString{String: "free", Valid: true},
		}

		// If SMTP is configured, require email verification
		var rawVerifyToken string
		if cfg.SMTPHost != "" {
			user.EmailVerified = false

			tokenBytes := make([]byte, 32)
			if _, err := crypto_rand.Read(tokenBytes); err != nil {
				log.Error().Err(err).Msg("Failed to generate verification token")
				helpers.WriteError(w, http.StatusInternalServerError, "Registration failed")
				return
			}
			rawVerifyToken = hex.EncodeToString(tokenBytes)
			hashSum := sha256.Sum256([]byte(rawVerifyToken))
			tokenHash := hex.EncodeToString(hashSum[:])

			user.VerificationTokenHash = sql.NullString{String: tokenHash, Valid: true}
			expiresAt := model.NewSQLiteTime(time.Now().Add(24 * time.Hour))
			user.VerificationTokenExpiresAt = &expiresAt
		}

		if err := db.UserStore().Create(ctx, user); err != nil {
			log.Error().Err(err).Msg("Registration error")
			helpers.WriteError(w, http.StatusInternalServerError, "Registration failed")
			return
		}
		alerts.Fire(alerts.Event{
			Type:     alerts.TypeNewUserRegistration,
			Severity: "info",
			Message:  "New user registered",
			Details:  map[string]interface{}{"userId": user.ID},
		})
		if alerts.Track("reg_ip:"+r.RemoteAddr, 3, 1*time.Hour, 4*time.Hour) {
			alerts.Fire(alerts.Event{
				Type:     alerts.TypeRegistrationSpike,
				Severity: "warning",
				Message:  "3+ new accounts registered from the same IP within 1 hour",
				Details:  map[string]interface{}{"ip": r.RemoteAddr},
			})
		}

		// Send verification email if SMTP is configured
		if cfg.SMTPHost != "" && rawVerifyToken != "" {
			service.SendVerificationEmail(cfg, email, rawVerifyToken)
		}

		// Mint a 30-day session token so the frontend is logged in
		// immediately after registration. The plaintext apiKey is also
		// returned EXACTLY ONCE here — the frontend MUST display it in the
		// one-time modal and never store it.
		rawSessionToken, sessionHash, err := model.GenerateSessionToken()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate session token on register")
			helpers.WriteError(w, http.StatusInternalServerError, "Registration failed")
			return
		}
		session := &model.Session{
			UserID:    user.ID,
			TokenHash: sessionHash,
			ExpiresAt: model.NewSQLiteTime(time.Now().Add(30 * 24 * time.Hour)),
		}
		if ua := r.Header.Get("User-Agent"); ua != "" {
			session.UserAgent = sql.NullString{String: ua, Valid: true}
		}
		if r.RemoteAddr != "" {
			session.IPAddress = sql.NullString{String: r.RemoteAddr, Valid: true}
		}
		_ = db.SessionStore().Create(ctx, session)

		// Auth endpoints intentionally return apiKey — use unsanitized
		resp := map[string]interface{}{
			"id":                 user.ID,
			"email":              user.Email,
			"apiKey":             user.APIKey, // ONE-TIME — never returned again
			"emailVerified":      user.EmailVerified,
			"createdAt":          user.CreatedAt,
			"subscriptionStatus": "free",
			"sessionToken":       rawSessionToken,
			"sessionExpiresAt":   session.ExpiresAt,
		}
		if rawVerifyToken != "" {
			resp["verificationToken"] = rawVerifyToken
		}
		helpers.WriteJSONUnsanitized(w, http.StatusCreated, resp)
	})

	// GET /auth/verify
	r.Get("/verify", func(w http.ResponseWriter, r *http.Request) {
		rawToken := r.URL.Query().Get("token")
		if rawToken == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Verification token is required")
			return
		}

		hashSum := sha256.Sum256([]byte(rawToken))
		tokenHash := hex.EncodeToString(hashSum[:])

		ctx := r.Context()

		// Find user by verification token hash
		var user *model.User

		users := db.UserStore()
		sqlStore, ok := users.(*model.SQLUserStore)
		if !ok {
			helpers.WriteError(w, http.StatusInternalServerError, "Verification failed")
			return
		}

		user = &model.User{}
		col := func(name string) string { return model.Col(sqlStore.DBType, name) }
		tableName := `"Users"`
		if sqlStore.DBType == "sqlite" {
			tableName = "Users"
		}
		query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, col("verification_token_hash"))
		err := sqlStore.DB.GetContext(ctx, user, query, tokenHash)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid or expired verification token")
			return
		}

		// Check expiry
		if user.VerificationTokenExpiresAt != nil && user.VerificationTokenExpiresAt.Valid &&
			time.Now().After(user.VerificationTokenExpiresAt.Time) {
			helpers.WriteError(w, http.StatusBadRequest, "Verification token has expired")
			return
		}

		if user.EmailVerified {
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"message":       "Email already verified",
				"emailVerified": true,
			})
			return
		}

		// Mark as verified and clear token
		user.EmailVerified = true
		user.VerificationTokenHash = sql.NullString{}
		user.VerificationTokenExpiresAt = nil
		if err := db.UserStore().Update(ctx, user); err != nil {
			log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to verify email")
			helpers.WriteError(w, http.StatusInternalServerError, "Verification failed")
			return
		}

		log.Info().Int64("userId", user.ID).Msg("Email verified successfully")
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message":       "Email verified successfully",
			"emailVerified": true,
		})
	})

	// POST /auth/resend-verification
	r.With(middleware.AuthRateLimiter.Middleware).Post("/resend-verification", func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
			helpers.WriteError(w, http.StatusUnauthorized, "Bearer token required")
			return
		}
		rawToken := header[len(prefix):]
		hash := model.HashSessionToken(rawToken)

		ctx := r.Context()
		session, err := db.SessionStore().FindByTokenHash(ctx, hash)
		if err != nil || session == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid session")
			return
		}
		user, err := db.UserStore().FindByID(ctx, session.UserID)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}

		if user.EmailVerified {
			helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Email is already verified"})
			return
		}

		if cfg.SMTPHost == "" {
			helpers.WriteError(w, http.StatusServiceUnavailable, "Email service is not configured")
			return
		}

		// Generate new verification token
		tokenBytes := make([]byte, 32)
		if _, err := crypto_rand.Read(tokenBytes); err != nil {
			log.Error().Err(err).Msg("Failed to generate verification token")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to resend verification")
			return
		}
		newVerifyToken := hex.EncodeToString(tokenBytes)
		hashSum := sha256.Sum256([]byte(newVerifyToken))
		tokenHash := hex.EncodeToString(hashSum[:])

		user.VerificationTokenHash = sql.NullString{String: tokenHash, Valid: true}
		expiresAt := model.NewSQLiteTime(time.Now().Add(24 * time.Hour))
		user.VerificationTokenExpiresAt = &expiresAt

		if err := db.UserStore().Update(ctx, user); err != nil {
			log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to update verification token")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to resend verification")
			return
		}

		service.SendVerificationEmail(cfg, user.Email, newVerifyToken)

		log.Info().Int64("userId", user.ID).Msg("Verification email resent")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Verification email sent"})
	})

	// POST /auth/login
	r.With(middleware.AuthRateLimiter.Middleware).Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := parseBody(r)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		email := bodyStr(body, "email")
		password := bodyStr(body, "password")

		if email == "" || password == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		ctx := r.Context()
		user, err := db.UserStore().FindByEmail(ctx, email)
		if err != nil || user == nil {
			if alerts.Track("failed_auth:"+r.RemoteAddr, 5, 10*time.Minute, 30*time.Minute) {
				alerts.Fire(alerts.Event{
					Type:     alerts.TypeFailedAuthSpike,
					Severity: "critical",
					Message:  "5+ failed login attempts from the same IP in 10 minutes",
					Details:  map[string]interface{}{"ip": r.RemoteAddr},
				})
			}
			log.Warn().Str("ip", r.RemoteAddr).Msg("Login failed: invalid credentials")
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			if alerts.Track("failed_auth:"+r.RemoteAddr, 5, 10*time.Minute, 30*time.Minute) {
				alerts.Fire(alerts.Event{
					Type:     alerts.TypeFailedAuthSpike,
					Severity: "critical",
					Message:  "5+ failed login attempts from the same IP in 10 minutes",
					Details:  map[string]interface{}{"ip": r.RemoteAddr},
				})
			}
			log.Warn().Str("ip", r.RemoteAddr).Msg("Login failed: invalid credentials")
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Mint a 30-day session token for dashboard auth.
		rawSessionToken, sessionHash, err := model.GenerateSessionToken()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate session token")
			helpers.WriteError(w, http.StatusInternalServerError, "Login failed")
			return
		}
		session := &model.Session{
			UserID:    user.ID,
			TokenHash: sessionHash,
			ExpiresAt: model.NewSQLiteTime(time.Now().Add(30 * 24 * time.Hour)),
		}
		if ua := r.Header.Get("User-Agent"); ua != "" {
			session.UserAgent = sql.NullString{String: ua, Valid: true}
		}
		if r.RemoteAddr != "" {
			session.IPAddress = sql.NullString{String: r.RemoteAddr, Valid: true}
		}
		if err := db.SessionStore().Create(ctx, session); err != nil {
			log.Error().Err(err).Msg("Failed to persist session")
			helpers.WriteError(w, http.StatusInternalServerError, "Login failed")
			return
		}

		log.Info().Int64("userId", user.ID).Str("ip", r.RemoteAddr).Msg("Login successful")

		// The master API key is intentionally NOT returned. The dashboard
		// authenticates via the sessionToken below. Users see the master key
		// exactly once at registration / regeneration via the one-time modal.
		helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{
			"id":                     user.ID,
			"email":                  user.Email,
			"emailVerified":          user.EmailVerified,
			"apiKeyRotationRequired": user.APIKeyRotationRequired,
			"role":                   user.Role,
			"requestsThisMonth":      user.RequestsThisMonth,
			"createdAt":              user.CreatedAt,
			"sessionToken":           rawSessionToken,
			"sessionExpiresAt":       session.ExpiresAt,
		})
	})

	// POST /auth/logout — invalidates the current session token.
	// Requires Authorization: Bearer <session-token>.
	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
			helpers.WriteError(w, http.StatusUnauthorized, "Bearer token required")
			return
		}
		rawToken := header[len(prefix):]
		hash := model.HashSessionToken(rawToken)
		session, err := db.SessionStore().FindByTokenHash(r.Context(), hash)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Logout failed")
			return
		}
		if session == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		if err := db.SessionStore().Delete(r.Context(), session.ID); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Logout failed")
			return
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true})
	})

	// POST /auth/regenerate-key
	r.With(middleware.AuthRateLimiter.Middleware).Post("/regenerate-key", func(w http.ResponseWriter, r *http.Request) {
		body, err := parseBody(r)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		email := bodyStr(body, "email")
		password := bodyStr(body, "password")

		if email == "" || password == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		ctx := r.Context()
		user, err := db.UserStore().FindByEmail(ctx, email)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		newKey, err := model.GenerateAPIKey()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate API key")
			helpers.WriteError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		// Capture the OLD account identifier (apiKeyHash) BEFORE update so we
		// can disconnect any live WebSocket clients that were authenticated
		// under it. The plaintext old key isn't recoverable from the DB
		// (we never stored it), but the hash is — and clients are registered
		// under the hash in ClientManager.
		oldIdentifier := user.APIKeyHash.String

		user.APIKey = newKey // transient field; Update() hashes it into APIKeyHash
		// Clear the rotation-required flag now that the key has been rotated
		user.APIKeyRotationRequired = false
		if err := db.UserStore().Update(ctx, user); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to regenerate API key")
			return
		}

		// Invalidate ALL auth cache entries for this user. We don't know
		// every plaintext value the cache might be keyed under, so a user-id
		// scan is the only safe way to clear the cache after rotation.
		middleware.InvalidateCachedAuthForUser(user.ID)

		// Invalidate ALL dashboard sessions for this user — rotation is a
		// panic button, every browser the user is logged into on every
		// device should be kicked out.
		_, _ = db.SessionStore().DeleteAllByUser(ctx, user.ID)

		deleted, _ := db.ApiKeyStore().DeleteAllByUser(ctx, user.ID)
		if deleted > 0 {
			log.Info().Int64("deleted", deleted).Int64("userId", user.ID).Msg("Deleted scoped keys due to master key regeneration")
		}

		// Nuke all connection tokens: rotation is a panic button and must
		// invalidate every credential the user owns, including anything stored
		// in Foundry user flags. Users will need to re-pair.
		tokensDeleted, err := db.ConnectionTokenStore().DeleteAllByUser(ctx, user.ID)
		if err != nil {
			log.Warn().Err(err).Int64("userId", user.ID).Msg("Failed to delete connection tokens during key regeneration")
		} else if tokensDeleted > 0 {
			log.Info().Int64("deleted", tokensDeleted).Int64("userId", user.ID).Msg("Deleted connection tokens due to master key regeneration")
		}

		// Force-disconnect every live WebSocket client registered under the
		// old account identifier. This broadcasts via Redis pub/sub so
		// sessions on any instance are killed, not just this one.
		if manager != nil && oldIdentifier != "" {
			n := manager.BroadcastDisconnectByAPIKey(ctx, oldIdentifier, "Master API key regenerated")
			if n > 0 {
				log.Info().Int("count", n).Int64("userId", user.ID).Msg("Disconnected clients due to master key regeneration")
			}
		}

		helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]string{"apiKey": newKey})
	})

	// GET /auth/user-data
	r.Get("/user-data", func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
			helpers.WriteError(w, http.StatusUnauthorized, "Bearer token required")
			return
		}
		rawToken := header[len(prefix):]
		hash := model.HashSessionToken(rawToken)
		session, err := db.SessionStore().FindByTokenHash(r.Context(), hash)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if session == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid or expired session")
			return
		}

		user, err := db.UserStore().FindByID(r.Context(), session.UserID)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}

		if user.Disabled {
			helpers.WriteError(w, http.StatusForbidden, "Account disabled")
			return
		}

		helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{
			"id":                     user.ID,
			"email":                  user.Email,
			"emailVerified":          user.EmailVerified,
			"apiKeyRotationRequired": user.APIKeyRotationRequired,
			"role":                   user.Role,
			"requestsThisMonth":      user.RequestsThisMonth,
			"subscriptionStatus":     user.GetSubscriptionStatus(),
			"limits": map[string]interface{}{
				"monthlyLimit":     cfg.MonthlyRequestLimit,
				"unlimitedMonthly": user.GetSubscriptionStatus() == "active",
			},
		})
	})

	// POST /auth/forgot-password
	r.With(middleware.PasswordResetRateLimiter.Middleware).Post("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		body, err := parseBody(r)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		email := bodyStr(body, "email")
		genericMsg := "If an account with that email exists, a password reset link has been sent."

		if email == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Email is required")
			return
		}

		ctx := r.Context()
		user, _ := db.UserStore().FindByEmail(ctx, email)
		if user == nil {
			helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": genericMsg})
			return
		}

		db.PasswordResetTokenStore().InvalidateForUser(ctx, user.ID)

		rawToken, err := model.GenerateAPIKey()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate reset token")
			helpers.WriteError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		hashBytes := sha256.Sum256([]byte(rawToken))
		tokenHash := hex.EncodeToString(hashBytes[:])

		token := &model.PasswordResetToken{
			UserID:    user.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}
		db.PasswordResetTokenStore().Create(ctx, token)

		if alerts.Track("pwd_reset:"+r.RemoteAddr, 5, 10*time.Minute, 1*time.Hour) {
			alerts.Fire(alerts.Event{
				Type:     alerts.TypePasswordResetFlood,
				Severity: "warning",
				Message:  "5+ password reset requests from the same IP in 10 minutes",
				Details:  map[string]interface{}{"ip": r.RemoteAddr},
			})
		}

		service.SendPasswordResetEmail(cfg, email, rawToken)

		resp := map[string]string{"message": genericMsg}
		if os.Getenv("RETURN_RESET_TOKEN") == "true" && !cfg.IsProduction() {
			resp["token"] = rawToken
			log.Warn().Msg("RETURN_RESET_TOKEN is enabled — this is a dev-only feature and must not be used in production")
		}
		helpers.WriteJSON(w, http.StatusOK, resp)
	})

	// POST /auth/reset-password
	r.With(middleware.PasswordResetRateLimiter.Middleware).Post("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		body, err := parseBody(r)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		rawToken := bodyStr(body, "token")
		password := bodyStr(body, "password")

		if rawToken == "" || password == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Token and password are required")
			return
		}

		if valid, errMsg := validatePassword(password); !valid {
			helpers.WriteError(w, http.StatusBadRequest, errMsg)
			return
		}

		hashBytes := sha256.Sum256([]byte(rawToken))
		tokenHash := hex.EncodeToString(hashBytes[:])

		ctx := r.Context()
		resetToken, _ := db.PasswordResetTokenStore().FindByTokenHash(ctx, tokenHash)
		if resetToken == nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid or expired reset token")
			return
		}

		user, _ := db.UserStore().FindByID(ctx, resetToken.UserID)
		if user == nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid or expired reset token")
			return
		}

		newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
		if err != nil {
			log.Error().Err(err).Msg("Failed to hash new password")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to reset password")
			return
		}
		user.Password = string(newHash)
		if err := db.WithTx(ctx, func(tx *sqlx.Tx) error {
			userStore := db.TxUserStore(tx)
			tokenStore := db.TxPasswordResetTokenStore(tx)
			if err := userStore.Update(ctx, user); err != nil {
				return err
			}
			if err := tokenStore.MarkUsed(ctx, resetToken.ID); err != nil {
				return err
			}
			return tokenStore.InvalidateForUser(ctx, user.ID)
		}); err != nil {
			log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to reset password")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to reset password")
			return
		}

		log.Info().Int64("userId", user.ID).Msg("Password reset successfully")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password has been reset successfully"})
	})

	// GET /auth/validate-reset-token/:token
	r.Get("/validate-reset-token/{token}", func(w http.ResponseWriter, r *http.Request) {
		rawToken := chi.URLParam(r, "token")
		if rawToken == "" {
			helpers.WriteJSON(w, http.StatusOK, map[string]bool{"valid": false})
			return
		}

		hashBytes := sha256.Sum256([]byte(rawToken))
		tokenHash := hex.EncodeToString(hashBytes[:])

		ctx := r.Context()
		resetToken, _ := db.PasswordResetTokenStore().FindByTokenHash(ctx, tokenHash)
		if resetToken == nil {
			helpers.WriteJSON(w, http.StatusOK, map[string]bool{"valid": false})
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]bool{"valid": true})
	})

	// POST /auth/change-password
	r.With(middleware.AccountManagementRateLimiter.Middleware).Post("/change-password", func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			apiKey = r.Header.Get("X-API-Key")
		}
		if apiKey == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "API key is required")
			return
		}

		body, err := parseBody(r)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		currentPassword := bodyStr(body, "currentPassword")
		newPassword := bodyStr(body, "newPassword")

		if currentPassword == "" || newPassword == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Current password and new password are required")
			return
		}

		ctx := r.Context()
		user, err := db.UserStore().FindByAPIKey(ctx, apiKey)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Current password is incorrect")
			return
		}

		if valid, errMsg := validatePassword(newPassword); !valid {
			helpers.WriteError(w, http.StatusBadRequest, errMsg)
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
		user.Password = string(hash)
		db.UserStore().Update(ctx, user)

		alerts.Fire(alerts.Event{
			Type:     alerts.TypeAccountPasswordChange,
			Severity: "info",
			Message:  "User changed their account password",
			Details:  map[string]interface{}{"userId": user.ID},
		})

		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
	})

	// GET /auth/export-data
	r.With(middleware.AccountManagementRateLimiter.Middleware).Get("/export-data", func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			apiKey = r.Header.Get("X-API-Key")
		}
		if apiKey == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "API key is required")
			return
		}

		ctx := r.Context()
		user, err := db.UserStore().FindByAPIKey(ctx, apiKey)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}

		keys, _ := db.ApiKeyStore().FindAllByUser(ctx, user.ID)
		var scopedKeysExport []map[string]interface{}
		for _, k := range keys {
			scopedKeysExport = append(scopedKeysExport, map[string]interface{}{
				"id":                    k.ID,
				"name":                  k.Name,
				"scopedClientId":        k.ScopedClientID.String,
				"scopedUserId":          k.ScopedUserID.String,
				"monthlyLimit":          k.MonthlyLimit.Int64,
				"expiresAt":             k.ExpiresAt,
				"enabled":               k.Enabled,
				"createdAt":             k.CreatedAt,
			})
		}
		if scopedKeysExport == nil {
			scopedKeysExport = []map[string]interface{}{}
		}

		helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{
			"exportDate": time.Now().UTC().Format(time.RFC3339),
			"user": map[string]interface{}{
				"id": user.ID, "email": user.Email,
				"createdAt": user.CreatedAt, "updatedAt": user.UpdatedAt,
			},
			"subscription": map[string]interface{}{
				"status":           user.GetSubscriptionStatus(),
				"stripeCustomerId": user.StripeCustomerID.String,
				"subscriptionId":   user.SubscriptionID.String,
			},
			"usage": map[string]interface{}{
				"requestsThisMonth": user.RequestsThisMonth,
				"requestsToday":     user.RequestsToday,
			},
			"apiAccess":  map[string]interface{}{"hasMasterKey": user.APIKeyHash.Valid},
			"scopedKeys": scopedKeysExport,
		})
	})

	// DELETE /auth/account
	r.With(middleware.AccountManagementRateLimiter.Middleware).Delete("/account", func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			apiKey = r.Header.Get("X-API-Key")
		}
		if apiKey == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "API key is required")
			return
		}

		body, err := parseBody(r)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		confirmEmail := bodyStr(body, "confirmEmail")
		password := bodyStr(body, "password")

		ctx := r.Context()
		user, err := db.UserStore().FindByAPIKey(ctx, apiKey)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}

		if confirmEmail == "" || confirmEmail != user.Email {
			helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error":   "Email confirmation required",
				"message": "Please provide your email address in the confirmEmail field to confirm account deletion",
			})
			return
		}

		if password == "" {
			helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error":   "Password required",
				"message": "Please provide your password to confirm account deletion",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid password")
			return
		}

		// Cancel active Stripe subscription
		if user.SubscriptionStatus.Valid && user.SubscriptionStatus.String == "active" &&
			user.SubscriptionID.Valid && user.SubscriptionID.String != "" {
			stripe.Key = cfg.StripeSecretKey
			_, err := stripeSubscription.Cancel(user.SubscriptionID.String, nil)
			if err != nil {
				log.Warn().Err(err).Int64("userId", user.ID).Msg("Failed to cancel Stripe subscription during account deletion")
			} else {
				log.Info().Int64("userId", user.ID).Msg("Cancelled Stripe subscription during account deletion")
			}
		}

		// Capture the account identifier (apiKeyHash) BEFORE the deletion tx
		// so we can force-disconnect every live WebSocket registered under
		// it. The plaintext master key isn't recoverable from the DB.
		masterIdentifier := user.APIKeyHash.String

		if err := db.WithTx(ctx, func(tx *sqlx.Tx) error {
			if _, err := db.TxApiKeyStore(tx).DeleteAllByUser(ctx, user.ID); err != nil {
				return err
			}
			if err := db.TxPasswordResetTokenStore(tx).InvalidateForUser(ctx, user.ID); err != nil {
				return err
			}
			return db.TxUserStore(tx).Delete(ctx, user.ID)
		}); err != nil {
			log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to delete account")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to delete account")
			return
		}

		// Best-effort cleanup of secondary state. These run outside the tx
		// because a partial failure here (e.g., Redis down) shouldn't block
		// account deletion — the user is already gone from the Users table.
		if _, err := db.ConnectionTokenStore().DeleteAllByUser(ctx, user.ID); err != nil {
			log.Warn().Err(err).Int64("userId", user.ID).Msg("Failed to delete connection tokens during account deletion")
		}
		middleware.InvalidateCachedAuthForUser(user.ID)

		if manager != nil && masterIdentifier != "" {
			n := manager.BroadcastDisconnectByAPIKey(ctx, masterIdentifier, "Account deleted")
			if n > 0 {
				log.Info().Int("count", n).Int64("userId", user.ID).Msg("Disconnected clients due to account deletion")
			}
		}

		log.Info().Int64("userId", user.ID).Msg("Account deleted")
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Account successfully deleted",
		})
	})

	// Connection token routes
	RegisterConnectionTokenRoutes(r, db, cfg, manager)

	// Credential vault, known clients, notification settings, per-key notification settings
	RegisterCredentialRoutes(r, db, cfg, manager)

	// Activity log (merged: connection + cross-world + module events)
	r.Mount("/activity", ActivityRouter(db))

	// Key request routes
	RegisterKeyRequestRoutes(r, db, cfg)

	// World pair request routes (device-flow pairing)
	RegisterPairRequestRoutes(r, db, cfg, manager)

	// Scoped API Key CRUD — requires auth middleware
	r.Route("/api-keys", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil))

		// POST /auth/api-keys
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Scoped keys cannot create other scoped keys. Use your master API key.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			body, err := parseBody(r)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			name := bodyStr(body, "name")
			if name == "" {
				helpers.WriteError(w, http.StatusBadRequest, "Name is required")
				return
			}

			keyStr, err := model.GenerateAPIKey()
			if err != nil {
				log.Error().Err(err).Msg("Failed to generate API key")
				helpers.WriteError(w, http.StatusInternalServerError, "Internal server error")
				return
			}

			// Parse scopes — default to safe read-only set if not specified
			var scopes []string
			if scopesRaw, exists := body["scopes"]; exists {
				if arr, ok := scopesRaw.([]interface{}); ok {
					for _, s := range arr {
						if str, ok := s.(string); ok {
							scopes = append(scopes, str)
						}
					}
				}
			}
			if len(scopes) == 0 {
				helpers.WriteError(w, http.StatusBadRequest, "At least one scope is required")
				return
			}
			if invalid := model.ValidateScopes(scopes); invalid != "" {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid scope: "+invalid)
				return
			}
			scopesStr := model.ScopesString(scopes)

			// Parse scopedClientIds (multi-client)
			scopedClientIDs := sql.NullString{}
			if idsRaw, exists := body["scopedClientIds"]; exists {
				if arr, ok := idsRaw.([]interface{}); ok {
					ids := make([]string, 0, len(arr))
					for _, s := range arr {
						if str, ok := s.(string); ok && str != "" {
							ids = append(ids, str)
						}
					}
					if len(ids) > 0 {
						scopedClientIDs = sql.NullString{String: strings.Join(ids, ","), Valid: true}
					}
				}
			}

			// Parse scopedUserIds (per-client user map)
			scopedUserIDs := sql.NullString{}
			if raw, exists := body["scopedUserIds"]; exists && raw != nil {
				if m, ok := raw.(map[string]interface{}); ok && len(m) > 0 {
					clean := make(map[string]string, len(m))
					for k, v := range m {
						if s, ok := v.(string); ok && s != "" {
							clean[k] = s
						}
					}
					if len(clean) > 0 {
						if b, err := json.Marshal(clean); err == nil {
							scopedUserIDs = sql.NullString{String: string(b), Valid: true}
						}
					}
				}
			}

			key := &model.ApiKey{
				UserID:          user.ID,
				Key:             keyStr,
				Name:            name,
				ScopedClientID:  sql.NullString{String: bodyStr(body, "scopedClientId"), Valid: bodyStr(body, "scopedClientId") != ""},
				ScopedClientIDs: scopedClientIDs,
				ScopedUserID:    sql.NullString{String: bodyStr(body, "scopedUserId"), Valid: bodyStr(body, "scopedUserId") != ""},
				ScopedUserIDs:   scopedUserIDs,
				Scopes:          scopesStr,
				Enabled:         true,
			}

			if dl := bodyStr(body, "monthlyLimit"); dl != "" {
				if v, err := strconv.ParseInt(dl, 10, 64); err == nil {
					key.MonthlyLimit = sql.NullInt64{Int64: v, Valid: true}
				}
			}
			if ea := bodyStr(body, "expiresAt"); ea != "" {
				if t, err := time.Parse(time.RFC3339, ea); err == nil {
					key.ExpiresAt = &model.SQLiteTime{Time: t, Valid: true}
				}
			}

			if err := db.ApiKeyStore().Create(r.Context(), key); err != nil {
				log.Error().Err(err).Msg("Error creating scoped API key")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to create API key")
				return
			}

			// Return full key (only shown on creation)
			helpers.WriteJSONUnsanitized(w, http.StatusCreated, map[string]interface{}{
				"id":                    key.ID,
				"key":                   key.Key,
				"name":                  key.Name,
				"scopes":                key.GetScopes(),
				"scopedClientId":        key.ScopedClientID.String,
				"scopedClientIds":       key.GetScopedClientIDs(),
				"scopedUserId":          key.ScopedUserID.String,
				"scopedUserIds":         key.GetScopedUserIDs(),
				"monthlyLimit":          key.MonthlyLimit.Int64,
				"expiresAt":             key.ExpiresAt,
				"enabled":               true,
				"createdAt":             key.CreatedAt,
			})
		})

		// GET /auth/api-keys
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key to manage scoped keys.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			keys, err := db.ApiKeyStore().FindAllByUser(r.Context(), user.ID)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to list API keys")
				return
			}

			var result []map[string]interface{}
			for _, k := range keys {
				masked := "..."
				if len(k.Key) >= 8 {
					masked = k.Key[:8] + "..."
				}
				result = append(result, map[string]interface{}{
					"id":                    k.ID,
					"key":                   masked,
					"name":                  k.Name,
					"scopes":                k.GetScopes(),
					"scopedClientId":        k.ScopedClientID.String,
					"scopedClientIds":       k.GetScopedClientIDs(),
					"scopedUserId":          k.ScopedUserID.String,
					"scopedUserIds":         k.GetScopedUserIDs(),
					"monthlyLimit":          k.MonthlyLimit.Int64,
					"requestsThisMonth":     k.RequestsThisMonth,
					"expiresAt":             k.ExpiresAt,
					"isExpired":             k.IsExpired(),
					"enabled":               k.Enabled,
					"createdAt":             k.CreatedAt,
					"updatedAt":             k.UpdatedAt,
				})
			}
			if result == nil {
				result = []map[string]interface{}{}
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"keys": result})
		})

		// PATCH /auth/api-keys/{id}
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key to manage scoped keys.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			keyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid key ID")
				return
			}

			ctx := r.Context()
			key, _ := db.ApiKeyStore().FindByID(ctx, keyID)
			if key == nil || key.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "API key not found")
				return
			}

			body, err := parseBody(r)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}

			if name := bodyStr(body, "name"); name != "" {
				key.Name = name
			}
			if scopesRaw, exists := body["scopes"]; exists {
				var patchScopes []string
				if arr, ok := scopesRaw.([]interface{}); ok {
					for _, s := range arr {
						if str, ok := s.(string); ok {
							patchScopes = append(patchScopes, str)
						}
					}
				}
				if len(patchScopes) == 0 {
					helpers.WriteError(w, http.StatusBadRequest, "At least one scope is required")
					return
				}
				if invalid := model.ValidateScopes(patchScopes); invalid != "" {
					helpers.WriteError(w, http.StatusBadRequest, "Invalid scope: "+invalid)
					return
				}
				key.Scopes = model.ScopesString(patchScopes)
			}
			if idsRaw, exists := body["scopedClientIds"]; exists {
				if idsRaw == nil {
					key.ScopedClientIDs = sql.NullString{}
				} else if arr, ok := idsRaw.([]interface{}); ok {
					ids := make([]string, 0, len(arr))
					for _, s := range arr {
						if str, ok := s.(string); ok && str != "" {
							ids = append(ids, str)
						}
					}
					if len(ids) > 0 {
						key.ScopedClientIDs = sql.NullString{String: strings.Join(ids, ","), Valid: true}
					} else {
						key.ScopedClientIDs = sql.NullString{}
					}
				}
			}
			if v, exists := body["scopedClientId"]; exists {
				if v == nil {
					key.ScopedClientID = sql.NullString{}
				} else if s, ok := v.(string); ok {
					key.ScopedClientID = sql.NullString{String: s, Valid: s != ""}
				}
			}
			if v, exists := body["scopedUserId"]; exists {
				if v == nil {
					key.ScopedUserID = sql.NullString{}
				} else if s, ok := v.(string); ok {
					key.ScopedUserID = sql.NullString{String: s, Valid: s != ""}
				}
			}
			if v, exists := body["scopedUserIds"]; exists {
				if v == nil {
					key.ScopedUserIDs = sql.NullString{}
				} else if m, ok := v.(map[string]interface{}); ok {
					clean := make(map[string]string, len(m))
					for k, val := range m {
						if s, ok := val.(string); ok && s != "" {
							clean[k] = s
						}
					}
					if len(clean) > 0 {
						if b, err := json.Marshal(clean); err == nil {
							key.ScopedUserIDs = sql.NullString{String: string(b), Valid: true}
						}
					} else {
						key.ScopedUserIDs = sql.NullString{}
					}
				}
			}
			if v, exists := body["monthlyLimit"]; exists {
				if v == nil {
					key.MonthlyLimit = sql.NullInt64{}
				} else if f, ok := v.(float64); ok {
					key.MonthlyLimit = sql.NullInt64{Int64: int64(f), Valid: true}
				}
			}
			if v, exists := body["expiresAt"]; exists {
				if v == nil {
					key.ExpiresAt = nil
				} else if s, ok := v.(string); ok {
					if t, err := time.Parse(time.RFC3339, s); err == nil {
						key.ExpiresAt = &model.SQLiteTime{Time: t, Valid: true}
					}
				}
			}
			if v, exists := body["enabled"]; exists {
				if b, ok := v.(bool); ok {
					key.Enabled = b
				}
			}
			if err := db.ApiKeyStore().Update(ctx, key); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update API key")
				return
			}

			log.Info().Int64("keyId", keyID).Int64("userId", user.ID).Msg("Updated scoped API key")
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"id": key.ID, "name": key.Name,
				"scopes": key.GetScopes(), "scopedClientIds": key.GetScopedClientIDs(),
				"scopedClientId": key.ScopedClientID.String, "scopedUserId": key.ScopedUserID.String,
				"scopedUserIds": key.GetScopedUserIDs(),
				"monthlyLimit": key.MonthlyLimit.Int64, "expiresAt": key.ExpiresAt,
				"enabled": key.Enabled, "updatedAt": key.UpdatedAt,
			})
		})

		// DELETE /auth/api-keys/{id}
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key to manage scoped keys.")
				return
			}

			keyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid key ID")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			ctx := r.Context()
			key, _ := db.ApiKeyStore().FindByID(ctx, keyID)
			if key == nil || key.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "API key not found")
				return
			}

			if err := db.ApiKeyStore().Delete(ctx, keyID); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to delete API key")
				return
			}

			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "API key deleted"})
		})

		// POST /auth/api-keys/{id}/regenerate
		r.Post("/{id}/regenerate", func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil || reqCtx.ScopedKey != nil {
				helpers.WriteError(w, http.StatusForbidden, "Use your master API key to manage scoped keys.")
				return
			}

			user, ok := reqCtx.User.(*model.User)
			if !ok || user == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
				return
			}

			keyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Invalid key ID")
				return
			}

			ctx := r.Context()
			key, _ := db.ApiKeyStore().FindByID(ctx, keyID)
			if key == nil || key.UserID != user.ID {
				helpers.WriteError(w, http.StatusNotFound, "API key not found")
				return
			}

			newKey, err := db.ApiKeyStore().RegenerateKey(ctx, keyID)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to regenerate API key")
				return
			}

			log.Info().Int64("keyId", keyID).Int64("userId", user.ID).Msg("Regenerated scoped API key")
			helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{"key": newKey})
		})
	})

	return r
}
