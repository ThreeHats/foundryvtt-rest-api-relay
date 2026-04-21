package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	appmw "github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// adminUserView is the PII-safe view of a user returned by admin endpoints.
// It deliberately omits password, full apiKey, stripe customer ID, and verification token hash.
type adminUserView struct {
	ID                     int64  `json:"id"`
	Email                  string `json:"email"`
	Role                   string `json:"role"`
	Disabled               bool   `json:"disabled"`
	EmailVerified          bool   `json:"emailVerified"`
	SubscriptionStatus     string `json:"subscriptionStatus"`
	RequestsToday          int    `json:"requestsToday"`
	RequestsThisMonth      int    `json:"requestsThisMonth"`
	MaxHeadlessSessions    int64  `json:"maxHeadlessSessions"`
	APIKeyRotationRequired bool   `json:"apiKeyRotationRequired"`
	CreatedAt              string `json:"createdAt"`
	UpdatedAt              string `json:"updatedAt"`
}

func newAdminUserView(u *model.User) adminUserView {
	v := adminUserView{
		ID:                     u.ID,
		Email:                  u.Email,
		Role:                   u.Role,
		Disabled:               u.Disabled,
		EmailVerified:          u.EmailVerified,
		RequestsToday:          u.RequestsToday,
		RequestsThisMonth:      u.RequestsThisMonth,
		APIKeyRotationRequired: u.APIKeyRotationRequired,
	}
	if u.SubscriptionStatus.Valid {
		v.SubscriptionStatus = u.SubscriptionStatus.String
	} else {
		v.SubscriptionStatus = "free"
	}
	if u.MaxHeadlessSessions.Valid {
		v.MaxHeadlessSessions = u.MaxHeadlessSessions.Int64
	}
	if u.CreatedAt.Valid {
		v.CreatedAt = u.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	if u.UpdatedAt.Valid {
		v.UpdatedAt = u.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	return v
}

// AdminUsersRouter exposes admin endpoints for managing users.
// Mounted under /admin/api/users (RequireAdmin already applied at parent).
func AdminUsersRouter(db *database.DB) chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/users — paginated list
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		offset, limit := parsePagination(r)
		users, total, err := db.UserStore().FindAllPaginated(r.Context(), offset, limit)
		if err != nil {
			log.Error().Err(err).Msg("admin users: list failed")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to list users")
			return
		}
		views := make([]adminUserView, 0, len(users))
		for _, u := range users {
			views = append(views, newAdminUserView(u))
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"users":  views,
			"total":  total,
			"offset": offset,
			"limit":  limit,
		})
	})

	// GET /admin/api/users/{id}
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		user, err := db.UserStore().FindByID(r.Context(), id)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		helpers.WriteJSON(w, http.StatusOK, newAdminUserView(user))
	})

	// PATCH /admin/api/users/{id} — update role / limits / max headless sessions
	r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		body, err := parseBody(r)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := db.UserStore().FindByID(r.Context(), id)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		changed := false
		if v, ok := body["role"].(string); ok && v != "" {
			if v != "user" && v != "admin" {
				helpers.WriteError(w, http.StatusBadRequest, "role must be 'user' or 'admin'")
				return
			}
			user.Role = v
			changed = true
		}
		if v, ok := body["maxHeadlessSessions"].(float64); ok {
			user.MaxHeadlessSessions = sql.NullInt64{Int64: int64(v), Valid: true}
			changed = true
		}
		if v, ok := body["emailVerified"].(bool); ok {
			user.EmailVerified = v
			changed = true
		}
		if !changed {
			helpers.WriteError(w, http.StatusBadRequest, "No supported fields to update")
			return
		}
		if err := db.UserStore().Update(r.Context(), user); err != nil {
			log.Error().Err(err).Msg("admin users: update failed")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}
		auditAdmin(r, db, "user.update", "user", strconv.FormatInt(id, 10), "")
		helpers.WriteJSON(w, http.StatusOK, newAdminUserView(user))
	})

	// POST /admin/api/users/{id}/disable
	r.Post("/{id}/disable", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		// Don't let admin disable themselves
		if me := appmw.GetAdminUser(r); me != nil && me.ID == id {
			helpers.WriteError(w, http.StatusBadRequest, "Cannot disable your own account")
			return
		}
		if err := db.UserStore().SetDisabled(r.Context(), id, true); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to disable user")
			return
		}
		appmw.InvalidateCachedAuthForUser(id)
		auditAdmin(r, db, "user.disable", "user", strconv.FormatInt(id, 10), "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "User disabled"})
	})

	// POST /admin/api/users/{id}/enable
	r.Post("/{id}/enable", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		if err := db.UserStore().SetDisabled(r.Context(), id, false); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to enable user")
			return
		}
		auditAdmin(r, db, "user.enable", "user", strconv.FormatInt(id, 10), "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "User enabled"})
	})

	// POST /admin/api/users/{id}/force-rotation — flag user as requiring master key rotation
	r.Post("/{id}/force-rotation", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		user, err := db.UserStore().FindByID(r.Context(), id)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		user.APIKeyRotationRequired = true
		if err := db.UserStore().Update(r.Context(), user); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to flag rotation")
			return
		}
		appmw.InvalidateCachedAuthForUser(id)
		auditAdmin(r, db, "user.flag_rotation", "user", strconv.FormatInt(id, 10), "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Rotation required flag set"})
	})

	// POST /admin/api/users/{id}/rotate-key
	r.Post("/{id}/rotate-key", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		user, err := db.UserStore().FindByID(r.Context(), id)
		if err != nil || user == nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		newKey, err := model.GenerateAPIKey()
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate key")
			return
		}
		user.APIKey = newKey
		user.APIKeyRotationRequired = false
		if err := db.UserStore().Update(r.Context(), user); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to rotate key")
			return
		}
		auditAdmin(r, db, "user.rotate_key", "user", strconv.FormatInt(id, 10), "")
		// Don't return the new key in the audit-log–only context.
		// Admin can view it via /admin/api/users/{id} only if we add that — for safety,
		// we return only the prefix here.
		prefix := newKey
		if len(prefix) > 8 {
			prefix = prefix[:8] + "..."
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message":   "API key rotated",
			"keyPrefix": prefix,
		})
	})

	// DELETE /admin/api/users/{id}
	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		// Don't let admin delete themselves
		if me := appmw.GetAdminUser(r); me != nil && me.ID == id {
			helpers.WriteError(w, http.StatusBadRequest, "Cannot delete your own account")
			return
		}
		if _, err := db.ApiKeyStore().DeleteAllByUser(r.Context(), id); err != nil {
			log.Warn().Err(err).Msg("admin users: failed to delete user's API keys")
		}
		if err := db.UserStore().Delete(r.Context(), id); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to delete user")
			return
		}
		auditAdmin(r, db, "user.delete", "user", strconv.FormatInt(id, 10), "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "User deleted"})
	})

	return r
}

// auditAdmin writes an audit log entry for an admin action.
// Best-effort: failures are logged but do not block the response.
func auditAdmin(r *http.Request, db *database.DB, action, targetType, targetID, details string) {
	admin := appmw.GetAdminUser(r)
	if admin == nil {
		return
	}
	entry := &model.AdminAuditLog{
		AdminUserID: admin.ID,
		Action:      action,
		TargetType:  targetType,
		IPAddress:   sql.NullString{String: r.RemoteAddr, Valid: true},
	}
	if targetID != "" {
		entry.TargetID = sql.NullString{String: targetID, Valid: true}
	}
	if details != "" {
		entry.Details = sql.NullString{String: details, Valid: true}
	}
	if err := db.AuditLogStore().Create(r.Context(), entry); err != nil {
		log.Warn().Err(err).Str("action", action).Msg("admin: audit log write failed")
	}
}

// parsePagination extracts offset/limit from query params with safe defaults.
func parsePagination(r *http.Request) (offset, limit int) {
	limit = 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return offset, limit
}

// parseInt64Param extracts a numeric chi URL param.
func parseInt64Param(w http.ResponseWriter, r *http.Request, name string) (int64, bool) {
	raw := chi.URLParam(r, name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "Invalid "+name)
		return 0, false
	}
	return id, true
}
