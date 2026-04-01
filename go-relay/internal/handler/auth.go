package handler

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v78"
	stripeSubscription "github.com/stripe/stripe-go/v78/subscription"
	"golang.org/x/crypto/bcrypt"
)

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

func parseBody(r *http.Request) map[string]interface{} {
	var body map[string]interface{}
	if r.Body != nil {
		data, err := io.ReadAll(r.Body)
		if err == nil && len(data) > 0 {
			json.Unmarshal(data, &body)
		}
	}
	return body
}

func bodyStr(body map[string]interface{}, key string) string {
	v, _ := body[key].(string)
	return v
}

// AuthRouter creates the auth route group.
func AuthRouter(db *database.DB, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	// POST /auth/register
	r.With(middleware.AuthRateLimiter.Middleware).Post("/register", func(w http.ResponseWriter, r *http.Request) {
		if cfg.DisableRegistration {
			helpers.WriteError(w, http.StatusForbidden, "Registration is disabled on this server")
			return
		}

		body := parseBody(r)
		email := bodyStr(body, "email")
		password := bodyStr(body, "password")

		if email == "" || password == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Email and password are required")
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

		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Registration failed")
			return
		}

		user := &model.User{
			Email:              email,
			Password:           string(hash),
			APIKey:             model.GenerateAPIKey(),
			SubscriptionStatus: sql.NullString{String: "free", Valid: true},
		}

		if err := db.UserStore().Create(ctx, user); err != nil {
			log.Error().Err(err).Msg("Registration error")
			helpers.WriteError(w, http.StatusInternalServerError, "Registration failed")
			return
		}

		// Auth endpoints intentionally return apiKey — use unsanitized
		helpers.WriteJSONUnsanitized(w, http.StatusCreated, map[string]interface{}{
			"id":                 user.ID,
			"email":              user.Email,
			"apiKey":             user.APIKey,
			"createdAt":          user.CreatedAt,
			"subscriptionStatus": "free",
		})
	})

	// POST /auth/login
	r.With(middleware.AuthRateLimiter.Middleware).Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body := parseBody(r)
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

		helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{
			"id":                user.ID,
			"email":             user.Email,
			"apiKey":            user.APIKey,
			"requestsThisMonth": user.RequestsThisMonth,
			"createdAt":         user.CreatedAt,
		})
	})

	// POST /auth/regenerate-key
	r.With(middleware.AuthRateLimiter.Middleware).Post("/regenerate-key", func(w http.ResponseWriter, r *http.Request) {
		body := parseBody(r)
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

		newKey := model.GenerateAPIKey()
		user.APIKey = newKey
		if err := db.UserStore().Update(ctx, user); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to regenerate API key")
			return
		}

		deleted, _ := db.ApiKeyStore().DeleteAllByUser(ctx, user.ID)
		if deleted > 0 {
			log.Info().Int64("deleted", deleted).Int64("userId", user.ID).Msg("Deleted scoped keys due to master key regeneration")
		}

		helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]string{"apiKey": newKey})
	})

	// GET /auth/user-data
	r.Get("/user-data", func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			apiKey = r.Header.Get("X-API-Key")
		}
		if apiKey == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "API key is required")
			return
		}

		user, err := db.UserStore().FindByAPIKey(r.Context(), apiKey)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}

		helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{
			"id":                 user.ID,
			"email":              user.Email,
			"apiKey":             user.APIKey,
			"requestsThisMonth":  user.RequestsThisMonth,
			"requestsToday":     user.RequestsToday,
			"subscriptionStatus": user.GetSubscriptionStatus(),
			"limits": map[string]interface{}{
				"dailyLimit":       cfg.DailyRequestLimit,
				"monthlyLimit":     cfg.FreeAPIRequestsLimit,
				"unlimitedMonthly": user.GetSubscriptionStatus() == "active",
			},
		})
	})

	// POST /auth/forgot-password
	r.With(middleware.PasswordResetRateLimiter.Middleware).Post("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		body := parseBody(r)
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

		rawToken := model.GenerateAPIKey()
		hashBytes := sha256.Sum256([]byte(rawToken))
		tokenHash := hex.EncodeToString(hashBytes[:])

		token := &model.PasswordResetToken{
			UserID:    user.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}
		db.PasswordResetTokenStore().Create(ctx, token)

		service.SendPasswordResetEmail(cfg, email, rawToken)

		resp := map[string]string{"message": genericMsg}
		if os.Getenv("RETURN_RESET_TOKEN") == "true" {
			resp["token"] = rawToken
		}
		helpers.WriteJSON(w, http.StatusOK, resp)
	})

	// POST /auth/reset-password
	r.With(middleware.PasswordResetRateLimiter.Middleware).Post("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		body := parseBody(r)
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

		newHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			log.Error().Err(err).Msg("Failed to hash new password")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to reset password")
			return
		}
		user.Password = string(newHash)
		if err := db.UserStore().Update(ctx, user); err != nil {
			log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to update user password")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to reset password")
			return
		}

		db.PasswordResetTokenStore().MarkUsed(ctx, resetToken.ID)
		db.PasswordResetTokenStore().InvalidateForUser(ctx, user.ID)

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

		body := parseBody(r)
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

		hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
		user.Password = string(hash)
		db.UserStore().Update(ctx, user)

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
				"dailyLimit":            k.DailyLimit.Int64,
				"expiresAt":             k.ExpiresAt,
				"enabled":               k.Enabled,
				"hasFoundryCredentials": k.HasStoredCredentials(),
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
			"apiAccess":  map[string]interface{}{"apiKey": user.APIKey},
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

		body := parseBody(r)
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

		db.ApiKeyStore().DeleteAllByUser(ctx, user.ID)
		db.PasswordResetTokenStore().InvalidateForUser(ctx, user.ID)
		db.UserStore().Delete(ctx, user.ID)

		log.Info().Int64("userId", user.ID).Msg("Account deleted")
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Account successfully deleted",
		})
	})

	// Scoped API Key CRUD — requires auth middleware
	r.Route("/api-keys", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(db, nil, cfg.DBType == "memory"))

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

			body := parseBody(r)
			name := bodyStr(body, "name")
			if name == "" {
				helpers.WriteError(w, http.StatusBadRequest, "Name is required")
				return
			}

			key := &model.ApiKey{
				UserID:         user.ID,
				Key:            model.GenerateAPIKey(),
				Name:           name,
				ScopedClientID: sql.NullString{String: bodyStr(body, "scopedClientId"), Valid: bodyStr(body, "scopedClientId") != ""},
				ScopedUserID:   sql.NullString{String: bodyStr(body, "scopedUserId"), Valid: bodyStr(body, "scopedUserId") != ""},
				Enabled:        true,
			}

			if dl := bodyStr(body, "dailyLimit"); dl != "" {
				if v, err := strconv.ParseInt(dl, 10, 64); err == nil {
					key.DailyLimit = sql.NullInt64{Int64: v, Valid: true}
				}
			}
			if ea := bodyStr(body, "expiresAt"); ea != "" {
				if t, err := time.Parse(time.RFC3339, ea); err == nil {
					key.ExpiresAt = &model.SQLiteTime{Time: t, Valid: true}
				}
			}

			foundryPassword := bodyStr(body, "foundryPassword")
			if foundryPassword != "" {
				if !service.IsEncryptionAvailable(cfg.CredentialsEncryptionKey) {
					helpers.WriteError(w, http.StatusBadRequest, "Credential storage is not available. CREDENTIALS_ENCRYPTION_KEY is not configured.")
					return
				}
				encrypted, err := service.Encrypt(foundryPassword, cfg.CredentialsEncryptionKey)
				if err != nil {
					helpers.WriteError(w, http.StatusInternalServerError, "Failed to encrypt credentials")
					return
				}
				key.EncryptedFoundryPassword = sql.NullString{String: encrypted.Ciphertext, Valid: true}
				key.PasswordIV = sql.NullString{String: encrypted.IV, Valid: true}
				key.PasswordAuthTag = sql.NullString{String: encrypted.AuthTag, Valid: true}
			}

			key.FoundryURL = sql.NullString{String: bodyStr(body, "foundryUrl"), Valid: bodyStr(body, "foundryUrl") != ""}
			key.FoundryUsername = sql.NullString{String: bodyStr(body, "foundryUsername"), Valid: bodyStr(body, "foundryUsername") != ""}

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
				"scopedClientId":        key.ScopedClientID.String,
				"scopedUserId":          key.ScopedUserID.String,
				"dailyLimit":            key.DailyLimit.Int64,
				"expiresAt":             key.ExpiresAt,
				"hasFoundryCredentials": foundryPassword != "",
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
				masked := k.Key[:8] + "..."
				result = append(result, map[string]interface{}{
					"id":                    k.ID,
					"key":                   masked,
					"name":                  k.Name,
					"scopedClientId":        k.ScopedClientID.String,
					"scopedUserId":          k.ScopedUserID.String,
					"dailyLimit":            k.DailyLimit.Int64,
					"requestsToday":         k.RequestsToday,
					"expiresAt":             k.ExpiresAt,
					"isExpired":             k.IsExpired(),
					"enabled":               k.Enabled,
					"hasFoundryCredentials": k.HasStoredCredentials(),
					"foundryUrl":            k.FoundryURL.String,
					"foundryUsername":        k.FoundryUsername.String,
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

			body := parseBody(r)

			if name := bodyStr(body, "name"); name != "" {
				key.Name = name
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
			if v, exists := body["dailyLimit"]; exists {
				if v == nil {
					key.DailyLimit = sql.NullInt64{}
				} else if f, ok := v.(float64); ok {
					key.DailyLimit = sql.NullInt64{Int64: int64(f), Valid: true}
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
			if v, exists := body["foundryUrl"]; exists {
				if v == nil {
					key.FoundryURL = sql.NullString{}
				} else if s, ok := v.(string); ok {
					key.FoundryURL = sql.NullString{String: s, Valid: s != ""}
				}
			}
			if v, exists := body["foundryUsername"]; exists {
				if v == nil {
					key.FoundryUsername = sql.NullString{}
				} else if s, ok := v.(string); ok {
					key.FoundryUsername = sql.NullString{String: s, Valid: s != ""}
				}
			}
			if v, exists := body["foundryPassword"]; exists {
				if v == nil {
					key.EncryptedFoundryPassword = sql.NullString{}
					key.PasswordIV = sql.NullString{}
					key.PasswordAuthTag = sql.NullString{}
				} else if pw, ok := v.(string); ok && pw != "" {
					if !service.IsEncryptionAvailable(cfg.CredentialsEncryptionKey) {
						helpers.WriteError(w, http.StatusBadRequest, "Credential storage is not available.")
						return
					}
					encrypted, err := service.Encrypt(pw, cfg.CredentialsEncryptionKey)
					if err != nil {
						helpers.WriteError(w, http.StatusInternalServerError, "Failed to encrypt credentials")
						return
					}
					key.EncryptedFoundryPassword = sql.NullString{String: encrypted.Ciphertext, Valid: true}
					key.PasswordIV = sql.NullString{String: encrypted.IV, Valid: true}
					key.PasswordAuthTag = sql.NullString{String: encrypted.AuthTag, Valid: true}
				}
			}

			if err := db.ApiKeyStore().Update(ctx, key); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to update API key")
				return
			}

			log.Info().Int64("keyId", keyID).Int64("userId", user.ID).Msg("Updated scoped API key")
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"id": key.ID, "name": key.Name,
				"scopedClientId": key.ScopedClientID.String, "scopedUserId": key.ScopedUserID.String,
				"dailyLimit": key.DailyLimit.Int64, "expiresAt": key.ExpiresAt,
				"enabled": key.Enabled, "hasFoundryCredentials": key.HasStoredCredentials(),
				"foundryUrl": key.FoundryURL.String, "foundryUsername": key.FoundryUsername.String,
				"updatedAt": key.UpdatedAt,
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
	})

	return r
}
