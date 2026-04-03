package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	appmw "github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server holds all dependencies and provides the HTTP router.
type Server struct {
	cfg     *config.Config
	db      *database.DB
	redis   *config.RedisClient
	version string
	router  *chi.Mux

	// WebSocket infrastructure
	ClientManager   *ws.ClientManager
	PendingReqs     *ws.PendingRequests
	SSEManager      *helpers.SSEManager
	SheetSessions   *ws.SheetSessionManager
	Headless        *worker.HeadlessManager
}

// New creates a new Server with routes configured.
func New(cfg *config.Config, db *database.DB, redis *config.RedisClient, version string) *Server {
	s := &Server{
		cfg:           cfg,
		db:            db,
		redis:         redis,
		version:       version,
		ClientManager: ws.NewClientManager(redis, cfg.InstanceID()),
		PendingReqs:    ws.NewPendingRequests(),
		SSEManager:     helpers.NewSSEManager(),
		SheetSessions:  ws.NewSheetSessionManager(cfg.MaxSheetSessionsPerKey),
	}
	if cfg.AllowHeadless {
		s.Headless = worker.NewHeadlessManager(s.ClientManager, redis, cfg)
	}

	// Clean up when a Foundry client disconnects
	s.ClientManager.OnClientRemoved = func(clientID string) {
		s.SheetSessions.TerminateSessionsForClient(clientID)
		if s.Headless != nil {
			s.Headless.OnClientDisconnected(clientID)
		}
	}

	// Register WS message handlers with SSE fan-out
	fanout := &ws.EventFanout{
		OnChatEvent: func(clientID string, data map[string]interface{}) {
			s.fanoutChatEvent(clientID, data)
		},
		OnRollData: func(clientID string, data map[string]interface{}) {
			s.fanoutRollData(clientID, data)
		},
	}
	ws.RegisterMessageHandlers(s.ClientManager, s.PendingReqs, fanout, s.SheetSessions)

	// Start background cleanup loops
	ctx := context.Background()
	s.ClientManager.StartCleanupLoop(ctx, time.Duration(cfg.ClientCleanupIntervalMs)*time.Millisecond)
	s.PendingReqs.StartCleanupLoop(ctx, 30*time.Second, 30*time.Second)
	s.SheetSessions.StartCleanupLoop(ctx)
	if s.Headless != nil {
		s.Headless.StartCleanupLoop(ctx)
	}

	// Set up auto-start for scoped keys with stored credentials
	if s.Headless != nil {
		var autoStartMu sync.Mutex
		autoStartCooldowns := make(map[int64]time.Time) // scopedKeyID -> cooldown until

		helpers.AutoStartFunc = func(reqCtx *helpers.RequestContext) string {
			if reqCtx == nil || reqCtx.ScopedKey == nil || reqCtx.MasterAPIKey == "" {
				return ""
			}

			keyID := reqCtx.ScopedKey.ID

			// Check cooldown
			autoStartMu.Lock()
			if until, ok := autoStartCooldowns[keyID]; ok && time.Now().Before(until) {
				autoStartMu.Unlock()
				return ""
			}
			autoStartMu.Unlock()

			clientID, err := s.Headless.StartHeadlessWithStoredCredentials(
				reqCtx.ScopedKey.ID, reqCtx.MasterAPIKey, s.db, s.cfg, "")
			if err != nil {
				log.Warn().Err(err).Int64("scopedKeyId", keyID).Msg("Auto-start headless failed")
				// Cooldown: don't retry for 60 seconds
				autoStartMu.Lock()
				autoStartCooldowns[keyID] = time.Now().Add(60 * time.Second)
				autoStartMu.Unlock()
				return ""
			}
			return clientID
		}
	}

	s.router = s.setupRouter()
	return s
}

// Router returns the chi router.
func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key", "x-api-key"},
		ExposedHeaders:   []string{"Content-Disposition", "Content-Type", "X-Image-Width", "X-Image-Height"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Request forwarder for multi-instance routing (Fly.io)
	if s.redis != nil && s.redis.IsConnected() {
		forwarder := appmw.NewRequestForwarder(s.redis, s.cfg)
		r.Use(forwarder.Middleware)
	}

	// Status endpoint
	r.Get("/api/status", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":    "ok",
			"version":   s.version,
			"websocket": "/relay",
		})
	})

	// Health endpoint
	r.Get("/api/health", s.healthHandler)

	// Auth routes
	r.Mount("/auth", handler.AuthRouter(s.db, s.cfg))

	// API routes — auth middleware needed by Stripe and API routes
	isMemory := s.cfg.DBType == "memory"
	authMw := appmw.AuthMiddleware(s.db, s.ClientManager, isMemory)
	usageMw := appmw.TrackAPIUsage(s.db, s.cfg.FreeAPIRequestsLimit, s.cfg.DailyRequestLimit, isMemory)

	// Stripe routes — subscriptions require auth, webhooks do not
	r.Route("/api/subscriptions", func(sub chi.Router) {
		sub.Use(authMw)
		sub.Mount("/", handler.StripeRouter(s.db, s.cfg))
	})
	r.Mount("/api/webhooks", handler.WebhookRouter(s.db, s.cfg))

	// WebSocket routes — registered before API route group to avoid catch-all conflict
	// Use isMemory (not IsLocalMode) so sqlite still validates API keys
	relayCfg := &ws.RelayConfig{
		PingInterval:    time.Duration(s.cfg.WSPingIntervalMs) * time.Millisecond,
		CleanupInterval: time.Duration(s.cfg.ClientCleanupIntervalMs) * time.Millisecond,
		ValidateAPIKey:   service.MakeWSValidateAPIKey(s.db, isMemory),
		ValidateHeadless: func(clientID, token string) (bool, error) {
			if s.Headless == nil {
				return false, nil
			}
			return s.Headless.ValidateHeadlessSession(clientID, token)
		},
	}
	relayHandler := ws.HandleRelayConnection(s.ClientManager, relayCfg)
	r.HandleFunc("/relay", relayHandler)
	r.HandleFunc("/relay/", relayHandler)

	// Client API WebSocket
	clientWSHandler := ws.HandleClientAPIConnection(s.ClientManager, s.PendingReqs, &ws.ClientAPIConfig{
		PingInterval:   time.Duration(s.cfg.WSPingIntervalMs) * time.Millisecond,
		ValidateAPIKey: service.MakeWSValidateAPIKey(s.db, isMemory),
		TrackUsage: func(apiKey string) (bool, string) {
			return service.TrackWSAPIUsage(context.Background(), s.db, apiKey, s.cfg.DailyRequestLimit, isMemory)
		},
		AutoStart: func(masterAPIKey, scopedClientID, scopedUserID string) string {
			if helpers.AutoStartFunc == nil {
				return ""
			}
			reqCtx := &helpers.RequestContext{
				MasterAPIKey: masterAPIKey,
				ScopedKey: &helpers.ScopedKeyInfo{
					ScopedClientID: scopedClientID,
					ScopedUserID:   scopedUserID,
				},
			}
			return helpers.AutoStartFunc(reqCtx)
		},
		SSEManager:     s.SSEManager,
		SheetSessions:  s.SheetSessions,
	})
	r.HandleFunc("/ws/api", clientWSHandler)

	// Authenticated API routes
	handler.RegisterAPIRoutes(r, s.ClientManager, s.PendingReqs, s.cfg, s.db, s.SSEManager, s.Headless, authMw, usageMw)

	// Static file serving
	baseDir := os.Getenv("STATIC_DIR")
	if baseDir == "" {
		baseDir = ".." // Relative to go-relay directory (local dev)
	}

	// API spec endpoints with dynamic URL injection
	if p := findStaticFile(baseDir, "public/openapi.json"); p != "" {
		r.Get("/openapi.json", handler.OpenAPIHandler(p))
	}
	if p := findStaticFile(baseDir, "public/asyncapi.json"); p != "" {
		r.Get("/asyncapi.json", handler.AsyncAPIHandler(p))
	}
	if p := findStaticFile(baseDir, "public/api-docs.json"); p != "" {
		r.Get("/api/docs", handler.APIDocsHandler(p))
	}

	// Static assets
	if d := findStaticDir(baseDir, "public"); d != "" {
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(d))))
	}
	if d := findStaticDir(baseDir, "public-dist/_astro"); d != "" {
		r.Handle("/_astro/*", http.StripPrefix("/_astro/", http.FileServer(http.Dir(d))))
	}
	if p := findStaticFile(baseDir, "public/default-token.png"); p != "" {
		r.Get("/default-token.png", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, p) })
	}

	// Astro frontend pages
	r.Get("/privacy", handler.ServeStaticPage(
		filepath.Join(baseDir, "public-dist/privacy/index.html"),
		filepath.Join(baseDir, "public/privacy.html"),
	))
	r.Get("/subscription-success", handler.ServeStaticPage(
		filepath.Join(baseDir, "public-dist/subscription-success/index.html"),
		filepath.Join(baseDir, "public/subscription-success.html"),
	))
	r.Get("/subscription-cancel", handler.ServeStaticPage(
		filepath.Join(baseDir, "public-dist/subscription-cancel/index.html"),
		filepath.Join(baseDir, "public/subscription-cancel.html"),
	))

	// Docusaurus documentation — check multiple locations
	docsDir := ""
	for _, candidate := range []string{
		filepath.Join(baseDir, "docs/build"),                             // standard layout (main repo)
		filepath.Join(baseDir, "..", "..", "..", "docs/build"),           // worktree (.claude/worktrees/X/) → main repo
		filepath.Join(baseDir, "..", "..", "..", "..", "docs/build"),     // deeper nesting
	} {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			docsDir = candidate
			break
		}
	}
	if docsDir != "" {
		r.Get("/docs", handler.DocsFileServer(docsDir))
		r.Get("/docs/*", handler.DocsFileServer(docsDir))
	}

	// Asset proxy
	r.Get("/proxy-asset/*", handler.ProxyAssetHandler())

	// Root — WebSocket upgrade or Astro frontend
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {
			relayHandler(w, r)
			return
		}
		// Serve Astro frontend
		for _, p := range []string{
			filepath.Join(baseDir, "public-dist/index.html"),
			filepath.Join(baseDir, "public/index.html"),
		} {
			if _, err := os.Stat(p); err == nil {
				http.ServeFile(w, r, p)
				return
			}
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":     "ok",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"instanceId": s.cfg.InstanceID(),
	}

	if s.redis != nil {
		status, _ := s.redis.CheckHealth(r.Context())
		health["redis"] = status
	}

	writeJSON(w, http.StatusOK, health)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// fanoutChatEvent sends chat events to all SSE and WS event subscribers for the client.
func (s *Server) fanoutChatEvent(clientID string, data map[string]interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	jsonStr := string(jsonBytes)

	eventType := "chat-create"
	if et, ok := data["eventType"].(string); ok {
		eventType = "chat-" + et
	}

	// Fan out to SSE connections
	for _, conn := range s.SSEManager.GetChatSSE(clientID) {
		// Apply filters
		if conn.Filters.Speaker != "" {
			if speaker, ok := data["speaker"].(string); ok && speaker != conn.Filters.Speaker {
				continue
			}
		}
		if conn.Filters.Type != nil {
			if msgType, ok := data["type"].(float64); ok && int(msgType) != *conn.Filters.Type {
				continue
			}
		}

		select {
		case <-conn.Done:
			continue
		default:
			fmt.Fprintf(conn.W, "event: %s\ndata: %s\n\n", eventType, jsonStr)
			conn.Flusher.Flush()
		}
	}

	// Fan out to WS event connections
	for _, conn := range s.SSEManager.GetWSEvents(clientID) {
		if conn.Channel != "chat-events" {
			continue
		}
		if conn.Filters.Speaker != "" {
			if speaker, ok := data["speaker"].(string); ok && speaker != conn.Filters.Speaker {
				continue
			}
		}
		if conn.Filters.Type != nil {
			if msgType, ok := data["type"].(float64); ok && int(msgType) != *conn.Filters.Type {
				continue
			}
		}
		conn.SendFunc(map[string]interface{}{"type": "chat-event", "event": eventType, "data": data})
	}
}

// fanoutRollData sends roll events to all SSE and WS event subscribers for the client.
func (s *Server) fanoutRollData(clientID string, data map[string]interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	jsonStr := string(jsonBytes)

	for _, conn := range s.SSEManager.GetRollSSE(clientID) {
		// Apply userId filter
		if conn.Filters.UserID != "" {
			if uid, ok := data["userId"].(string); ok && uid != conn.Filters.UserID {
				continue
			}
		}

		select {
		case <-conn.Done:
			continue
		default:
			fmt.Fprintf(conn.W, "event: roll\ndata: %s\n\n", jsonStr)
			conn.Flusher.Flush()
		}
	}

	// Fan out to WS event connections
	for _, conn := range s.SSEManager.GetWSEvents(clientID) {
		if conn.Channel != "roll-events" {
			continue
		}
		if conn.Filters.UserID != "" {
			if uid, ok := data["userId"].(string); ok && uid != conn.Filters.UserID {
				continue
			}
		}
		conn.SendFunc(map[string]interface{}{"type": "roll-event", "data": data})
	}
}

func findStaticFile(baseDir, relPath string) string {
	p := filepath.Join(baseDir, relPath)
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return ""
}

func findStaticDir(baseDir, relPath string) string {
	p := filepath.Join(baseDir, relPath)
	if info, err := os.Stat(p); err == nil && info.IsDir() {
		return p
	}
	return ""
}

// import guard
var _ = fmt.Sprintf
