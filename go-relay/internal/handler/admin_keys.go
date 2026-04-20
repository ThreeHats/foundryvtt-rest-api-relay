package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// adminApiKeyView is the PII-safe view of an API key for admin endpoints.
// Never includes the full key value — only a prefix.
type adminApiKeyView struct {
	ID              int64    `json:"id"`
	UserID          int64    `json:"userId"`
	Name            string   `json:"name"`
	KeyPrefix       string   `json:"keyPrefix"`
	Enabled         bool     `json:"enabled"`
	IsExpired       bool     `json:"isExpired"`
	Scopes          []string `json:"scopes"`
	ScopedClientIDs []string `json:"scopedClientIds"`
	ScopedUserID    string   `json:"scopedUserId"`
	MonthlyLimit       int64    `json:"monthlyLimit"`
	RequestsThisMonth  int      `json:"requestsThisMonth"`
	ExpiresAt       string   `json:"expiresAt,omitempty"`
	CreatedAt       string   `json:"createdAt"`
}

func newAdminApiKeyView(k *model.ApiKey) adminApiKeyView {
	prefix := k.Key
	if len(prefix) > 8 {
		prefix = prefix[:8] + "..."
	}
	v := adminApiKeyView{
		ID:             k.ID,
		UserID:         k.UserID,
		Name:           k.Name,
		KeyPrefix:      prefix,
		Enabled:        k.Enabled,
		IsExpired:      k.IsExpired(),
		Scopes:             k.GetScopes(),
		RequestsThisMonth:  k.RequestsThisMonth,
	}
	if ids := k.GetScopedClientIDs(); ids != nil {
		v.ScopedClientIDs = ids
	} else {
		v.ScopedClientIDs = []string{}
	}
	if k.ScopedUserID.Valid {
		v.ScopedUserID = k.ScopedUserID.String
	}
	if k.MonthlyLimit.Valid {
		v.MonthlyLimit = k.MonthlyLimit.Int64
	}
	if k.ExpiresAt != nil && k.ExpiresAt.Valid {
		v.ExpiresAt = k.ExpiresAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	if k.CreatedAt.Valid {
		v.CreatedAt = k.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	return v
}

// AdminKeysRouter exposes admin endpoints for managing API keys across all users.
func AdminKeysRouter(db *database.DB) chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/keys — paginated list across all users
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		offset, limit := parsePagination(r)
		keys, total, err := db.ApiKeyStore().FindAllPaginated(r.Context(), offset, limit)
		if err != nil {
			log.Error().Err(err).Msg("admin keys: list failed")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to list keys")
			return
		}
		views := make([]adminApiKeyView, 0, len(keys))
		for _, k := range keys {
			views = append(views, newAdminApiKeyView(k))
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"keys":   views,
			"total":  total,
			"offset": offset,
			"limit":  limit,
		})
	})

	// GET /admin/api/keys/{id}
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		key, err := db.ApiKeyStore().FindByID(r.Context(), id)
		if err != nil || key == nil {
			helpers.WriteError(w, http.StatusNotFound, "Key not found")
			return
		}
		helpers.WriteJSON(w, http.StatusOK, newAdminApiKeyView(key))
	})

	// PATCH /admin/api/keys/{id} — enable/disable, change scopes
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
		key, err := db.ApiKeyStore().FindByID(r.Context(), id)
		if err != nil || key == nil {
			helpers.WriteError(w, http.StatusNotFound, "Key not found")
			return
		}
		changed := false
		if v, ok := body["enabled"].(bool); ok {
			key.Enabled = v
			changed = true
		}
		if v, ok := body["scopes"].([]interface{}); ok {
			scopes := make([]string, 0, len(v))
			for _, s := range v {
				if str, ok := s.(string); ok {
					scopes = append(scopes, str)
				}
			}
			key.Scopes = strings.Join(scopes, ",")
			changed = true
		}
		if !changed {
			helpers.WriteError(w, http.StatusBadRequest, "No supported fields to update")
			return
		}
		if err := db.ApiKeyStore().Update(r.Context(), key); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to update key")
			return
		}
		auditAdmin(r, db, "key.update", "apiKey", strconv.FormatInt(id, 10), "")
		helpers.WriteJSON(w, http.StatusOK, newAdminApiKeyView(key))
	})

	// DELETE /admin/api/keys/{id} — revoke key
	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInt64Param(w, r, "id")
		if !ok {
			return
		}
		if err := db.ApiKeyStore().Delete(r.Context(), id); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to revoke key")
			return
		}
		auditAdmin(r, db, "key.revoke", "apiKey", strconv.FormatInt(id, 10), "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Key revoked"})
	})

	return r
}
