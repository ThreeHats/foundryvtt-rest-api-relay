package handler

import (
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
)

// AdminClientsRouter exposes admin endpoints for viewing/managing connected WebSocket clients.
func AdminClientsRouter(db *database.DB, manager *ws.ClientManager) chi.Router {
	r := chi.NewRouter()

	// GET /admin/api/clients — list all connected clients across all users
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		infos := manager.GetAllClientInfos()
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"total":   len(infos),
			"clients": infos,
		})
	})

	// GET /admin/api/clients/{id}
	r.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		client := manager.GetClient(id)
		if client == nil {
			helpers.WriteError(w, http.StatusNotFound, "Client not found")
			return
		}
		helpers.WriteJSON(w, http.StatusOK, client.Info(""))
	})

	// POST /admin/api/clients/{id}/disconnect — force disconnect a client
	r.Post("/{id}/disconnect", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		client := manager.GetClient(id)
		if client == nil {
			helpers.WriteError(w, http.StatusNotFound, "Client not found")
			return
		}
		client.Disconnect()
		manager.RemoveClient(id)
		auditAdmin(req, db, "client.disconnect", "client", id, "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Client disconnected"})
	})

	return r
}
