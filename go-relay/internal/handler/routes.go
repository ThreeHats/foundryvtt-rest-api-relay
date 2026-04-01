package handler

import (
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
)

// RegisterAPIRoutes registers all authenticated API routes on the given router.
// All routes in this group require auth + usage tracking middleware.
func RegisterAPIRoutes(r chi.Router, mgr *ws.ClientManager, pending *ws.PendingRequests, cfg *config.Config, db *database.DB, sseMgr *helpers.SSEManager, headless *worker.HeadlessManager, authMw, usageMw func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(authMw)
		r.Use(usageMw)

		// Entity CRUD
		r.Get("/get", h(mgr, pending, entityGet))
		r.Post("/create", h(mgr, pending, entityCreate))
		r.Put("/update", h(mgr, pending, entityUpdate))
		r.Delete("/delete", h(mgr, pending, entityDelete))
		r.Post("/give", h(mgr, pending, entityGive))
		r.Post("/remove", h(mgr, pending, entityRemove))
		r.Post("/decrease", h(mgr, pending, entityDecrease))
		r.Post("/increase", h(mgr, pending, entityIncrease))
		r.Post("/kill", h(mgr, pending, entityKill))

		// Search
		r.Get("/search", h(mgr, pending, searchGet))

		// Rolls
		r.Get("/rolls", h(mgr, pending, rollsGet))
		r.Get("/lastroll", h(mgr, pending, lastRollGet))
		r.Post("/roll", h(mgr, pending, rollPost))
		r.Get("/rolls/subscribe", sseRollSubscribe(mgr, sseMgr))

		// Encounters
		r.Get("/encounters", h(mgr, pending, encountersGet))
		r.Post("/start-encounter", h(mgr, pending, startEncounter))
		r.Post("/next-turn", h(mgr, pending, nextTurn))
		r.Post("/next-round", h(mgr, pending, nextRound))
		r.Post("/last-turn", h(mgr, pending, lastTurn))
		r.Post("/last-round", h(mgr, pending, lastRound))
		r.Post("/end-encounter", h(mgr, pending, endEncounter))
		r.Post("/add-to-encounter", h(mgr, pending, addToEncounter))
		r.Post("/remove-from-encounter", h(mgr, pending, removeFromEncounter))

		// Macros
		r.Get("/macros", h(mgr, pending, macrosGet))
		r.Post("/macro/{uuid}/execute", h(mgr, pending, macroExecute))

		// Structure
		r.Get("/structure", h(mgr, pending, structureGet))
		r.Get("/contents/*", contentsDeprecated)
		r.Get("/get-folder", h(mgr, pending, getFolderRoute))
		r.Post("/create-folder", h(mgr, pending, createFolderRoute))
		r.Delete("/delete-folder", h(mgr, pending, deleteFolderRoute))

		// Utility
		r.Post("/select", h(mgr, pending, selectPost))
		r.Get("/selected", h(mgr, pending, selectedGet))
		r.Get("/players", h(mgr, pending, playersGet))
		r.Post("/execute-js", h(mgr, pending, executeJsPost))

		// Scenes
		r.Get("/scene", h(mgr, pending, sceneGet))
		r.Post("/scene", h(mgr, pending, sceneCreate))
		r.Put("/scene", h(mgr, pending, sceneUpdate))
		r.Delete("/scene", h(mgr, pending, sceneDelete))
		r.Post("/switch-scene", h(mgr, pending, switchScene))

		// Canvas
		r.Get("/canvas/{documentType}", h(mgr, pending, canvasGet))
		r.Post("/canvas/{documentType}", h(mgr, pending, canvasCreate))
		r.Put("/canvas/{documentType}", h(mgr, pending, canvasUpdate))
		r.Delete("/canvas/{documentType}", h(mgr, pending, canvasDelete))

		// Chat
		r.Get("/chat", h(mgr, pending, chatGet))
		r.Post("/chat", h(mgr, pending, chatSend))
		r.Delete("/chat/{messageId}", h(mgr, pending, chatDeleteMsg))
		r.Delete("/chat", h(mgr, pending, chatFlush))
		r.Get("/chat/subscribe", sseChatSubscribe(mgr, sseMgr))

		// Effects
		r.Get("/effects", h(mgr, pending, effectsGet))
		r.Post("/effects", h(mgr, pending, effectsAdd))
		r.Delete("/effects", h(mgr, pending, effectsRemove))

		// Clients
		r.Get("/clients", clientsHandler(mgr, cfg))

		// File System
		r.Get("/file-system", h(mgr, pending, fileSystemGet))
		r.Post("/upload", uploadHandler(mgr, pending))
		r.Get("/download", h(mgr, pending, downloadFileGet))

		// Sheet
		r.Get("/sheet", sheetGetHandler(mgr, pending))
		r.Get("/sheet/image", sheetImageHandler(mgr, pending))

		// Session
		r.Post("/session-handshake", sessionHandshakeHandler(db, cfg))
		r.Post("/start-session", sessionStartHandler(db, cfg, headless))
		r.Delete("/end-session", sessionEndHandler(headless))
		r.Get("/session", sessionListHandler(headless, mgr))

		// D&D 5e
		r.Mount("/dnd5e", Dnd5eRouter(mgr, pending))
	})
}

// h wraps a route config builder into an http.HandlerFunc via CreateAPIRoute.
func h(mgr *ws.ClientManager, pending *ws.PendingRequests, builder func() helpers.APIRouteConfig) http.HandlerFunc {
	return helpers.CreateAPIRoute(mgr, pending, builder())
}

// Subscribe to real-time roll events via Server-Sent Events (SSE)
//
// Opens a persistent SSE connection that streams roll events as they occur.
// @tag Roll
// @param {string} clientId [query] Client ID for the Foundry world
// @param {string} userId [query] Foundry user ID or username to scope permissions (omit for GM-level access)
// @returns SSE event stream
func sseRollSubscribe(mgr *ws.ClientManager, sseMgr *helpers.SSEManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sseSubscribeHandler(w, r, mgr, sseMgr, "roll")
	}
}

// Subscribe to real-time chat events via Server-Sent Events (SSE)
//
// Opens a persistent SSE connection that streams chat events (create, update, delete)
// as they occur in the Foundry world.
// @tag Chat
// @param {string} clientId [query] Client ID for the Foundry world
// @param {string} speaker [query] Filter events by speaker name or actor ID
// @param {number} type [query] Filter events by chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll)
// @param {boolean} whisperOnly [query] Only receive whispered messages
// @param {string} userId [query] Foundry user ID or username to scope permissions (omit for GM-level access)
// @returns SSE event stream
func sseChatSubscribe(mgr *ws.ClientManager, sseMgr *helpers.SSEManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sseSubscribeHandler(w, r, mgr, sseMgr, "chat")
	}
}
