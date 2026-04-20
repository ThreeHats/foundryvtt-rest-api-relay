package handler

import (
	"database/sql"
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/alerts"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/metrics"
	appmw "github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// AdminAuthRouter exposes admin login/logout/me endpoints under /admin/auth.
// These routes are NOT protected by RequireAdmin (login is the entry point) but DO
// have IP allowlist + rate limiting applied at the parent router level.
func AdminAuthRouter(db *database.DB, authCfg *appmw.AdminAuthConfig) chi.Router {
	r := chi.NewRouter()

	// POST /admin/auth/login
	r.With(appmw.AdminLoginRateLimiter.Middleware).Post("/login", func(w http.ResponseWriter, r *http.Request) {
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

		// Check account lockout
		if appmw.AdminLockout.IsLocked(email) {
			helpers.WriteError(w, http.StatusTooManyRequests, "Account temporarily locked due to repeated failed attempts. Try again later.")
			return
		}

		ctx := r.Context()
		user, err := db.UserStore().FindByEmail(ctx, email)
		if err != nil || user == nil {
			appmw.AdminLockout.RecordFailure(email)
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Check disabled before password validation to avoid leaking role info via error messages.
		// All failure paths return the same generic "Invalid credentials" message.
		if user.Disabled {
			appmw.AdminLockout.RecordFailure(email)
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			appmw.AdminLockout.RecordFailure(email)
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		if user.Role != "admin" {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		appmw.AdminLockout.RecordSuccess(email)

		csrf, err := appmw.IssueAdminJWT(w, user, authCfg)
		if err != nil {
			log.Error().Err(err).Msg("admin: failed to issue JWT")
			metrics.AdminLoginsTotal.WithLabelValues("failure").Inc()
			helpers.WriteError(w, http.StatusInternalServerError, "Login failed")
			return
		}

		metrics.AdminLoginsTotal.WithLabelValues("success").Inc()

		// Audit log
		_ = db.AuditLogStore().Create(ctx, &model.AdminAuditLog{
			AdminUserID: user.ID,
			Action:      "admin.login",
			TargetType:  "session",
			TargetID:    sql.NullString{String: user.Email, Valid: true},
			IPAddress:   sql.NullString{String: r.RemoteAddr, Valid: true},
		})

		// Fire admin login alert (PII safe — uses user ID only, no email)
		alerts.Fire(alerts.Event{
			Type:     alerts.TypeAdminLogin,
			Severity: "info",
			Message:  "Admin login successful",
			Details: map[string]interface{}{
				"userId": user.ID,
				"ip":     r.RemoteAddr,
			},
		})

		w.Header().Set(appmw.AdminCSRFHeader, csrf)
		helpers.WriteJSONUnsanitized(w, http.StatusOK, map[string]interface{}{
			"email":     user.Email,
			"role":      user.Role,
			"csrfToken": csrf,
		})
	})

	// POST /admin/auth/logout
	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		// We don't fully RequireAdmin here because logout should still work even
		// if the cookie is partially invalid (clear it on best-effort basis).
		// But we DO require a valid JWT + CSRF — otherwise this becomes a CSRF vector.
		claims, err := appmw.ParseAdminCookie(r, authCfg)
		if err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Not authenticated")
			return
		}
		if r.Header.Get(appmw.AdminCSRFHeader) != claims.CSRFToken {
			helpers.WriteError(w, http.StatusForbidden, "CSRF token missing or invalid")
			return
		}

		// Add JTI to denylist so the cookie can't be reused
		if err := appmw.DenyAdminJWT(r, db, authCfg); err != nil {
			log.Warn().Err(err).Msg("admin: failed to deny JWT on logout")
		}
		appmw.ClearAdminCookie(w, authCfg.SecureCookies)

		// Audit log
		if userID, perr := parseAdminSubject(claims.Subject); perr == nil {
			_ = db.AuditLogStore().Create(r.Context(), &model.AdminAuditLog{
				AdminUserID: userID,
				Action:      "admin.logout",
				TargetType:  "session",
				IPAddress:   sql.NullString{String: r.RemoteAddr, Valid: true},
			})
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
	})

	// GET /admin/auth/me — returns admin session info if logged in.
	// Also re-issues the CSRF token header so mutations still work after a page reload
	// (the CSRF token is held in memory only, so it's lost on reload without this).
	r.With(appmw.RequireAdmin(db, authCfg)).Get("/me", func(w http.ResponseWriter, r *http.Request) {
		user := appmw.GetAdminUser(r)
		if user == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Not authenticated")
			return
		}
		// Re-issue CSRF token so the frontend can restore it after a page reload.
		if claims, err := appmw.ParseAdminCookie(r, authCfg); err == nil {
			w.Header().Set(appmw.AdminCSRFHeader, claims.CSRFToken)
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		})
	})

	return r
}
