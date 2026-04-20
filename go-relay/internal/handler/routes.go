package handler

import (
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	appmw "github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
)

// RegisterAPIRoutes registers all authenticated API routes on the given router.
// All routes in this group require auth + usage tracking middleware.
func RegisterAPIRoutes(r chi.Router, mgr *ws.ClientManager, pending *ws.PendingRequests, cfg *config.Config, db *database.DB, sseMgr *helpers.SSEManager, headless *worker.HeadlessManager, disp *service.Dispatcher, authMw, usageMw func(http.Handler) http.Handler) {
	scope := appmw.RequireScope

	r.Group(func(r chi.Router) {
		r.Use(authMw)
		r.Use(appmw.APIKeyRateLimiter.APIKeyMiddleware)
		r.Use(usageMw)

		// Entity read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEntityRead))
			r.Get("/get", h(mgr, pending, entityGet))
		})

		// Entity write
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEntityWrite))
			r.Post("/create", h(mgr, pending, entityCreate))
			r.Put("/update", h(mgr, pending, entityUpdate))
			r.Delete("/delete", h(mgr, pending, entityDelete))
			r.Post("/give", h(mgr, pending, entityGive))
			r.Post("/remove", h(mgr, pending, entityRemove))
			r.Post("/decrease", h(mgr, pending, entityDecrease))
			r.Post("/increase", h(mgr, pending, entityIncrease))
			r.Post("/kill", h(mgr, pending, entityKill))
		})

		// Search
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeSearch))
			r.Get("/search", h(mgr, pending, searchGet))
		})

		// Roll read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeRollRead))
			r.Get("/rolls", h(mgr, pending, rollsGet))
			r.Get("/lastroll", h(mgr, pending, lastRollGet))
		})

		// Roll execute
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeRollExecute))
			r.Post("/roll", h(mgr, pending, rollPost))
		})

		// Roll subscribe (events)
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEventsSubscribe))
			r.Get("/rolls/subscribe", sseRollSubscribe(mgr, sseMgr))
		})

		// Encounter read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEncounterRead))
			r.Get("/encounters", h(mgr, pending, encountersGet))
		})

		// Encounter manage
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEncounterManage))
			r.Post("/start-encounter", h(mgr, pending, startEncounter))
			r.Post("/next-turn", h(mgr, pending, nextTurn))
			r.Post("/next-round", h(mgr, pending, nextRound))
			r.Post("/last-turn", h(mgr, pending, lastTurn))
			r.Post("/last-round", h(mgr, pending, lastRound))
			r.Post("/end-encounter", h(mgr, pending, endEncounter))
			r.Post("/add-to-encounter", h(mgr, pending, addToEncounter))
			r.Post("/remove-from-encounter", h(mgr, pending, removeFromEncounter))
		})

		// Macro list
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeMacroList))
			r.Get("/macros", h(mgr, pending, macrosGet))
		})

		// Macro execute
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeMacroExecute))
			r.Post("/macro/{uuid}/execute", helpers.CreateAPIRoute(mgr, pending, macroExecute(disp, db)))
		})

		// Structure read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeStructureRead))
			r.Get("/structure", h(mgr, pending, structureGet))
			r.Get("/contents/*", contentsDeprecated)
			r.Get("/get-folder", h(mgr, pending, getFolderRoute))
		})

		// Structure write
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeStructureWrite))
			r.Post("/create-folder", h(mgr, pending, createFolderRoute))
			r.Delete("/delete-folder", h(mgr, pending, deleteFolderRoute))
		})

		// Utility (entity-level operations that use entity scopes)
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEntityRead))
			r.Post("/select", h(mgr, pending, selectPost))
			r.Get("/selected", h(mgr, pending, selectedGet))
			r.Get("/players", h(mgr, pending, playersGet))
		})

		// execute-js — opt-in only, never in defaults
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeExecuteJS))
			r.Post("/execute-js", helpers.CreateAPIRoute(mgr, pending, executeJsPost(disp, db)))
		})

		// User read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeUserRead))
			r.Get("/users", h(mgr, pending, usersGet))
			r.Get("/user", h(mgr, pending, userGet))
		})

		// User write
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeUserWrite))
			r.Post("/user", h(mgr, pending, userCreate))
			r.Put("/user", h(mgr, pending, userUpdate))
			r.Delete("/user", h(mgr, pending, userDelete))
		})

		// Scene read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeSceneRead))
			r.Get("/scene", h(mgr, pending, sceneGet))
			r.Get("/scene/image", sceneImageHandler(mgr, pending))
			r.Get("/scene/image/raw", h(mgr, pending, sceneRawImage))
		})

		// Scene write
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeSceneWrite))
			r.Post("/scene", h(mgr, pending, sceneCreate))
			r.Put("/scene", h(mgr, pending, sceneUpdate))
			r.Delete("/scene", h(mgr, pending, sceneDelete))
			r.Post("/switch-scene", h(mgr, pending, switchScene))
		})

		// Canvas read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeCanvasRead))
			r.Get("/canvas/{documentType}", h(mgr, pending, canvasGet))
			r.Get("/measure-distance", h(mgr, pending, measureDistance))
		})

		// Canvas write
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeCanvasWrite))
			r.Post("/canvas/{documentType}", h(mgr, pending, canvasCreate))
			r.Put("/canvas/{documentType}", h(mgr, pending, canvasUpdate))
			r.Delete("/canvas/{documentType}", h(mgr, pending, canvasDelete))
			r.Post("/move-token", h(mgr, pending, tokenMove))
		})

		// Chat read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeChatRead))
			r.Get("/chat", h(mgr, pending, chatGet))
		})

		// Chat write
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeChatWrite))
			r.Post("/chat", h(mgr, pending, chatSend))
			r.Delete("/chat/{messageId}", h(mgr, pending, chatDeleteMsg))
			r.Delete("/chat", h(mgr, pending, chatFlush))
		})

		// Chat subscribe (events)
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEventsSubscribe))
			r.Get("/chat/subscribe", sseChatSubscribe(mgr, sseMgr))
		})

		// Effects read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEffectsRead))
			r.Get("/effects", h(mgr, pending, effectsGet))
			r.Get("/effects/list", h(mgr, pending, statusEffectsGet))
		})

		// Effects write
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEffectsWrite))
			r.Post("/effects", h(mgr, pending, effectsAdd))
			r.Delete("/effects", h(mgr, pending, effectsRemove))
		})

		// World info
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeWorldInfo))
			r.Get("/world-info", h(mgr, pending, worldInfoGet))
		})

		// Playlist control
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopePlaylistControl))
			r.Get("/playlists", h(mgr, pending, playlistsGet))
			r.Post("/playlist/play", h(mgr, pending, playlistPlay))
			r.Post("/playlist/stop", h(mgr, pending, playlistStop))
			r.Post("/playlist/next", h(mgr, pending, playlistNext))
			r.Post("/playlist/volume", h(mgr, pending, playlistVolume))
			r.Post("/play-sound", h(mgr, pending, soundPlay))
			r.Post("/stop-sound", h(mgr, pending, soundStop))
		})

		// Clients
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeClientsRead))
			r.Get("/clients", clientsHandler(mgr, cfg, db))
		})

		// File read
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeFileRead))
			r.Get("/file-system", h(mgr, pending, fileSystemGet))
			r.Get("/download", h(mgr, pending, downloadFileGet))
		})

		// File write
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeFileWrite))
			r.Post("/upload", uploadHandler(mgr, pending))
		})

		// Sheet
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeSheetRead))
			r.Get("/sheet", sheetHandler(mgr, pending))
		})

		// Session management
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeSessionManage))
			r.Post("/session-handshake", sessionHandshakeHandler(db, cfg))
			r.Post("/start-session", sessionStartHandler(db, cfg, headless))
			r.Delete("/end-session", sessionEndHandler(headless))
			r.Get("/session", sessionListHandler(headless, mgr))
		})

		// Event Subscriptions (SSE)
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeEventsSubscribe))
			r.Get("/hooks/subscribe", sseEventSubscribe(mgr, sseMgr, "hooks"))
			r.Get("/encounters/subscribe", sseEventSubscribe(mgr, sseMgr, "combat"))
			r.Get("/actor/subscribe", sseEventSubscribe(mgr, sseMgr, "actor"))
			r.Get("/scene/subscribe", sseEventSubscribe(mgr, sseMgr, "scene"))
		})

		// D&D 5e
		r.Group(func(r chi.Router) {
			r.Use(scope(model.ScopeDnd5e))
			r.Mount("/dnd5e", Dnd5eRouter(mgr, pending))
		})
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

// Subscribe to real-time events via Server-Sent Events (SSE)
//
// Opens a persistent SSE connection for the specified event type.
// Supported event types: hooks (all Foundry hooks), combat, actor, scene.
// @tag Events
// @param {string} clientId [query] Client ID for the Foundry world
// @param {string} hooks [query] Comma-separated hook names to filter (hooks type only)
// @param {string} encounterId [query] Filter by encounter ID (combat type only)
// @param {string} actorUuid [query] Actor UUID to subscribe to (actor type, required)
// @param {string} sceneId [query] Scene ID to filter (scene type only)
// @returns SSE event stream
func sseEventSubscribe(mgr *ws.ClientManager, sseMgr *helpers.SSEManager, eventType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sseSubscribeHandler(w, r, mgr, sseMgr, eventType)
	}
}
