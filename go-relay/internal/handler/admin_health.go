package handler

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
)

// AdminHealthRouter exposes detailed system health for the admin dashboard.
func AdminHealthRouter(
	cfg *config.Config,
	manager *ws.ClientManager,
	pending *ws.PendingRequests,
	headless *worker.HeadlessManager,
	redis *config.RedisClient,
	version string,
) chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/system/health
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		health := map[string]interface{}{
			"status":            "ok",
			"version":           version,
			"timestamp":         time.Now().UTC().Format(time.RFC3339),
			"instanceId":        cfg.InstanceID(),
			"goroutines":        runtime.NumGoroutine(),
			"memoryAllocBytes":  memStats.Alloc,
			"memoryHeapBytes":   memStats.HeapAlloc,
			"memorySysBytes":    memStats.Sys,
			"gcCount":           memStats.NumGC,
			"wsConnectionCount": manager.CountConnectedClients(),
		}

		if pending != nil {
			health["pendingRequestCount"] = pending.Count()
		}
		if headless != nil {
			health["headlessSessionCount"] = len(headless.ListSessions())
		}

		// Redis
		if redis != nil && redis.IsConnected() {
			ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
			defer cancel()
			if status, err := redis.CheckHealth(ctx); err == nil {
				health["redisStatus"] = status
			} else {
				health["redisStatus"] = "error"
				health["redisError"] = err.Error()
				health["status"] = "degraded"
			}
		} else {
			health["redisStatus"] = "disabled"
		}

		helpers.WriteJSON(w, http.StatusOK, health)
	})

	return r
}
