package handler

import (
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/go-chi/chi/v5"
)

// AdminSessionsRouter exposes admin endpoints for headless ChromeDP sessions.
func AdminSessionsRouter(db *database.DB, headless *worker.HeadlessManager) chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/headless-sessions
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		if headless == nil {
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"total":    0,
				"sessions": []interface{}{},
			})
			return
		}
		sessions := headless.ListSessions()
		// PII-safe view: don't include the Foundry username (it's a credential identifier).
		safeSessions := make([]map[string]interface{}, 0, len(sessions))
		for _, s := range sessions {
			safeSessions = append(safeSessions, map[string]interface{}{
				"sessionId":    s.SessionID,
				"clientId":     s.ClientID,
				"foundryUrl":   s.FoundryURL,
				"worldName":    s.WorldName,
				"startedAt":    s.StartedAt,
				"lastActivity": s.LastActivity,
			})
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"total":    len(safeSessions),
			"sessions": safeSessions,
		})
	})

	// DELETE /admin/api/headless-sessions/{id}
	r.Delete("/{id}", func(w http.ResponseWriter, req *http.Request) {
		if headless == nil {
			helpers.WriteError(w, http.StatusServiceUnavailable, "Headless sessions disabled")
			return
		}
		id := chi.URLParam(req, "id")
		if err := headless.EndSession(id); err != nil {
			helpers.WriteError(w, http.StatusNotFound, "Session not found or could not be ended")
			return
		}
		auditAdmin(req, db, "session.kill", "headless_session", id, "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Session ended"})
	})

	return r
}
