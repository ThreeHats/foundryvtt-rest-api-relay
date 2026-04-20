package handler

import (
	"net/http"
	"strconv"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/metrics"
	"github.com/go-chi/chi/v5"
)

// AdminMetricsRouter exposes the in-app rolling metrics for the admin dashboard.
// PII-safe: only user IDs (integers) are returned, never emails.
func AdminMetricsRouter() chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/metrics/overview
	r.Get("/overview", func(w http.ResponseWriter, req *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, metrics.Global.Overview())
	})

	// GET /admin/api/metrics/by-endpoint
	r.Get("/by-endpoint", func(w http.ResponseWriter, req *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"endpoints": metrics.Global.ByEndpoint(),
		})
	})

	// GET /admin/api/metrics/top-consumers?limit=N
	r.Get("/top-consumers", func(w http.ResponseWriter, req *http.Request) {
		limit := 10
		if v := req.URL.Query().Get("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
				limit = n
			}
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"users": metrics.Global.TopConsumers(limit),
		})
	})

	// GET /admin/api/metrics/errors
	r.Get("/errors", func(w http.ResponseWriter, req *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"errorsTotal": metrics.Global.Errors(),
		})
	})

	return r
}
