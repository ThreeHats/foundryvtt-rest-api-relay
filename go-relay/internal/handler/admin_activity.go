package handler

import (
	"net/http"
	"strconv"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/go-chi/chi/v5"
)

// AdminActivityRouter mounts GET /admin/api/activity.
// Auth is handled by the outer RequireAdmin middleware in AdminRouter.
func AdminActivityRouter(db *database.DB) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse filters — same as user endpoint but with optional userId override
		filters, limit, offset := parseActivityParams(r, 0 /* 0 = all users */)

		// Allow narrowing to a specific user
		if uid := r.URL.Query().Get("userId"); uid != "" {
			if v, err := strconv.ParseInt(uid, 10, 64); err == nil && v > 0 {
				filters.UserID = v
			}
		}

		events, total, err := mergeActivityLogs(r.Context(), db, filters, limit, offset)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to fetch activity log")
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"events": events,
			"total":  total,
			"offset": offset,
			"limit":  limit,
		})
	})

	return r
}
