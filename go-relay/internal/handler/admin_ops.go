package handler

import (
	"net/http"
	"sync"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
)

// FeatureFlags holds runtime-toggleable feature flags shared across the app.
// In multi-instance deployments this should be backed by Redis; for now we use
// in-memory state with a global lock — operators can change flags via the admin
// dashboard without restarting the server.
var FeatureFlags = &featureFlagsStore{
	flags: map[string]bool{
		"disable_registration": false,
		"maintenance_mode":     false,
	},
}

type featureFlagsStore struct {
	mu    sync.RWMutex
	flags map[string]bool
}

func (f *featureFlagsStore) Get(name string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.flags[name]
}

func (f *featureFlagsStore) Set(name string, value bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.flags[name] = value
}

func (f *featureFlagsStore) Snapshot() map[string]bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make(map[string]bool, len(f.flags))
	for k, v := range f.flags {
		out[k] = v
	}
	return out
}

// MaintenanceModeMiddleware blocks non-admin requests when maintenance_mode is on.
// Admin routes (/admin/*) are exempt so operators can re-enable the system.
func MaintenanceModeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Exempt /admin and /api/health
		if len(r.URL.Path) >= 6 && r.URL.Path[:6] == "/admin" {
			next.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == "/api/health" || r.URL.Path == "/api/status" {
			next.ServeHTTP(w, r)
			return
		}
		if FeatureFlags.Get("maintenance_mode") {
			helpers.WriteError(w, http.StatusServiceUnavailable, "Service is in maintenance mode")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AdminOpsRouter exposes operational controls: feature flags, broadcasts, force disconnects.
func AdminOpsRouter(db *database.DB, manager *ws.ClientManager) chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/ops/feature-flags
	r.Get("/feature-flags", func(w http.ResponseWriter, req *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, FeatureFlags.Snapshot())
	})

	// POST /admin/api/ops/feature-flags  { flag: "...", value: bool }
	r.Post("/feature-flags", func(w http.ResponseWriter, req *http.Request) {
		body, err := parseBody(req)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		flag := bodyStr(body, "flag")
		val, ok := body["value"].(bool)
		if flag == "" || !ok {
			helpers.WriteError(w, http.StatusBadRequest, "flag (string) and value (bool) required")
			return
		}
		// Whitelist of known flags
		switch flag {
		case "disable_registration", "maintenance_mode":
			// ok
		default:
			helpers.WriteError(w, http.StatusBadRequest, "Unknown flag")
			return
		}
		FeatureFlags.Set(flag, val)
		auditAdmin(req, db, "ops.feature_flag", "system", flag, "")
		helpers.WriteJSON(w, http.StatusOK, FeatureFlags.Snapshot())
	})

	// POST /admin/api/ops/force-disconnect/{id}
	r.Post("/force-disconnect/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		client := manager.GetClient(id)
		if client == nil {
			helpers.WriteError(w, http.StatusNotFound, "Client not found")
			return
		}
		client.Disconnect()
		manager.RemoveClient(id)
		auditAdmin(req, db, "ops.force_disconnect", "client", id, "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Disconnected"})
	})

	// POST /admin/api/ops/rate-limit-override — placeholder, accepts and audit-logs
	// (full per-user override store can be added later; this hook lets the dashboard
	// communicate intent now and persists via audit log.)
	r.Post("/rate-limit-override", func(w http.ResponseWriter, req *http.Request) {
		body, err := parseBody(req)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		auditAdmin(req, db, "ops.rate_limit_override", "user", bodyStr(body, "userId"), "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Override recorded"})
	})

	return r
}
