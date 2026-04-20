package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// AdminAuditRouter exposes the admin audit log query endpoint.
func AdminAuditRouter(db *database.DB) chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/audit-logs — paginated, filterable
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		offset, limit := parsePagination(r)
		filters := model.AuditLogFilters{
			Offset: offset,
			Limit:  limit,
		}
		if v := r.URL.Query().Get("action"); v != "" {
			filters.Action = v
		}
		if v := r.URL.Query().Get("targetType"); v != "" {
			filters.TargetType = v
		}
		if v := r.URL.Query().Get("adminUserId"); v != "" {
			if id, err := strconv.ParseInt(v, 10, 64); err == nil {
				filters.AdminUserID = id
			}
		}
		if v := r.URL.Query().Get("since"); v != "" {
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				filters.Since = &t
			}
		}
		if v := r.URL.Query().Get("until"); v != "" {
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				filters.Until = &t
			}
		}

		entries, total, err := db.AuditLogStore().FindAll(r.Context(), filters)
		if err != nil {
			log.Error().Err(err).Msg("admin audit: query failed")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to query audit logs")
			return
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"entries": entries,
			"total":   total,
			"offset":  offset,
			"limit":   limit,
		})
	})

	return r
}
