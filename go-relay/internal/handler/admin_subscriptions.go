package handler

import (
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/go-chi/chi/v5"
)

// AdminSubscriptionsRouter exposes Stripe subscription overview for admins.
// PII safety: returns aggregate counts and the same admin-safe user view used elsewhere
// (no Stripe customer IDs, no payment method details).
func AdminSubscriptionsRouter(db *database.DB) chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/subscriptions — overview: counts by status + paginated subscriber list
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		offset, limit := parsePagination(req)
		users, total, err := db.UserStore().FindAllPaginated(req.Context(), offset, limit)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to query subscriptions")
			return
		}

		// Aggregate counts (over the page only — for full counts, the dashboard
		// can iterate or we can add a dedicated count query later)
		statusCounts := map[string]int{}
		subscribers := make([]adminUserView, 0, len(users))
		for _, u := range users {
			status := "free"
			if u.SubscriptionStatus.Valid && u.SubscriptionStatus.String != "" {
				status = u.SubscriptionStatus.String
			}
			statusCounts[status]++
			if status != "free" {
				subscribers = append(subscribers, newAdminUserView(u))
			}
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"total":        total,
			"offset":       offset,
			"limit":        limit,
			"statusCounts": statusCounts,
			"subscribers":  subscribers,
		})
	})

	return r
}
