package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/alerts"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	appmw "github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server holds all dependencies and provides the HTTP router.
type Server struct {
	cfg     *config.Config
	db      *database.DB
	redis   *config.RedisClient
	version string
	router  *chi.Mux
	ctx     context.Context // shutdown context passed in from main; stops background goroutines

	// WebSocket infrastructure
	ClientManager   *ws.ClientManager
	PendingReqs     *ws.PendingRequests
	SSEManager      *helpers.SSEManager
	InteractiveSessions   *ws.InteractiveSessionManager
	Headless        *worker.HeadlessManager

	// Notification dispatcher (account + per-key routing)
	Dispatcher *service.Dispatcher
	// RemoteRequestBatcher batches cross-world activity notifications per user.
	RemoteRequestBatcher *service.RemoteRequestBatcher

	// headlessAutoStart is set during NewServer if a headless manager is configured.
	// Used by setupRouter to wire the WebSocket client API auto-start path.
	headlessAutoStart func(masterAPIKey, clientID string) string
}

// New creates a new Server with routes configured.
// ctx should be the application shutdown context so that background goroutines
// (cleanup loops, Redis subscriber) stop cleanly on SIGTERM.
func New(ctx context.Context, cfg *config.Config, db *database.DB, redis *config.RedisClient, version string) *Server {
	s := &Server{
		cfg:           cfg,
		db:            db,
		redis:         redis,
		version:       version,
		ctx:           ctx,
		ClientManager: ws.NewClientManager(redis, cfg.InstanceID()),
		PendingReqs:    ws.NewPendingRequests(),
		SSEManager:     helpers.NewSSEManager(), // callback wired below
		InteractiveSessions:  ws.NewInteractiveSessionManager(cfg.MaxInteractiveSessionsPerKey),
	}
	// Notify the Foundry module whenever subscriber counts change so it can
	// enable/disable event hooks on demand rather than always streaming everything.
	s.SSEManager.OnSubscriberCountChanged = s.notifyModuleChannelCount

	if cfg.AllowHeadless {
		s.Headless = worker.NewHeadlessManager(s.ClientManager, redis, cfg)
		// Wire DB + cfg for the AutoStartForKnownClient path used by the
		// remote-request handler. Done after construction to avoid an
		// import cycle in the worker package.
		s.Headless.SetDeps(&worker.HeadlessDeps{DB: db, Cfg: cfg})
	}

	// Build the unified notification dispatcher
	smtpCfg := &service.NotificationConfig{
		SMTPHost:    cfg.SMTPHost,
		SMTPPort:    cfg.SMTPPort,
		SMTPUser:    cfg.SMTPUser,
		SMTPPass:    cfg.SMTPPass,
		SMTPFrom:    cfg.SMTPFrom,
		FrontendURL: cfg.FrontendURL,
	}
	s.Dispatcher = service.NewDispatcher(
		service.NotificationStores{
			NotificationSettings:            db.NotificationSettingsStore(),
			ApiKeyNotificationSettings:      db.ApiKeyNotificationSettingsStore(),
			KnownClientNotificationSettings: db.KnownClientNotificationSettingsStore(),
			KnownClients:                    db.KnownClientStore(),
		},
		smtpCfg,
		cfg.FrontendURL,
	)

	// Initialise the alerts package with DB and SMTP config so Fire/Track
	// can be called from middleware and handlers without threading deps through.
	alerts.Init(db, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom, cfg.FrontendURL)

	// Track connection/disconnection events

	// Clean up when a Foundry client disconnects
	s.ClientManager.OnClientRemoved = func(clientID, apiKey string) {
		s.InteractiveSessions.TerminateSessionsForClient(clientID)
		if s.Headless != nil {
			s.Headless.OnClientDisconnected(clientID)
		}

		if alerts.Track("client_disconnect", 5, 2*time.Minute, 15*time.Minute) {
			alerts.Fire(alerts.Event{
				Type:     alerts.TypeClientDisconnectSpike,
				Severity: "warning",
				Message:  "5+ Foundry clients disconnected within 2 minutes",
				Details:  map[string]interface{}{"clientId": clientID},
			})
		}

		// Set known client offline + dispatch disconnect notification
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// `apiKey` is the per-account identifier (apiKeyHash); use
			// the hash variant to skip re-hashing.
			user, _ := db.UserStore().FindByAPIKeyHash(ctx, apiKey)
			if user == nil {
				return
			}

			// Pull the known client metadata BEFORE marking offline so we can include it
			known, _ := db.KnownClientStore().FindByClientID(ctx, user.ID, clientID)

			db.KnownClientStore().SetOffline(ctx, user.ID, clientID)

			nc := service.NotificationContext{
				Event:    service.EventDisconnect,
				UserID:   user.ID,
				ClientID: clientID,
			}
			if known != nil {
				nc.WorldID = known.WorldID.String
				nc.WorldTitle = known.WorldTitle.String
				nc.SystemID = known.SystemID.String
			}
			s.Dispatcher.Dispatch(nc)
		}()
	}

	// Silent-alarm notification when a duplicate-connection attempt is rejected.
	// `apiKey` is the per-account identifier (apiKeyHash) of the user who
	// owns the legitimate connection; that's who gets notified.
	s.ClientManager.OnDuplicateConnectionRejected = func(clientID, apiKey, newIP, existingIP string) {
		go func() {
			// Extract bare hosts (without port) for comparison.
			newHost, _, err := net.SplitHostPort(newIP)
			if err != nil {
				newHost = newIP
			}
			existingHost, _, err := net.SplitHostPort(existingIP)
			if err != nil {
				existingHost = existingIP
			}
			sameHost := newHost == existingHost && newHost != ""

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			user, _ := db.UserStore().FindByAPIKeyHash(ctx, apiKey)
			if user == nil {
				return
			}

			// Only fire the spike alert and user notification for cross-IP duplicates.
			// Same-host duplicates (e.g. two browser tabs on the same machine) are
			// expected behaviour and should not generate alerts.
			if !sameHost {
				if alerts.Track("dup_conn:"+fmt.Sprintf("%d", user.ID), 3, 10*time.Minute, 30*time.Minute) {
					alerts.Fire(alerts.Event{
						Type:     alerts.TypeDuplicateConnectionSpike,
						Severity: "warning",
						Message:  "3+ duplicate WebSocket connection rejections for the same user in 10 minutes",
						Details:  map[string]interface{}{"userId": user.ID, "clientId": clientID, "newIp": newIP, "existingIp": existingIP},
					})
				}
				s.Dispatcher.Dispatch(service.NotificationContext{
					Event:     service.EventDuplicateConnectionRejected,
					UserID:    user.ID,
					ClientID:  clientID,
					IPAddress: newIP,
					Severity:  "alert",
					Reason:    "Another client tried to connect with this clientId from a different IP while the legitimate client was already connected. The attempt was rejected.",
				})
			}
		}()
	}

	// Track when a Foundry client connects
	s.ClientManager.OnClientConnected = func(clientID, apiKey string, metadata ws.ClientInfo) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// `apiKey` here is the per-account identifier (apiKeyHash) under
			// which the WS client was registered. Look up by hash directly.
			user, _ := db.UserStore().FindByAPIKeyHash(ctx, apiKey)
			if user == nil {
				return
			}

			// Check for metadata mismatch against previously known client
			var metadataMatch bool = true
			var flagged bool
			var flagReason string

			existing, _ := db.KnownClientStore().FindByClientID(ctx, user.ID, clientID)
			if existing != nil {
				var mismatches []string
				if existing.WorldID.Valid && metadata.WorldID != "" && existing.WorldID.String != metadata.WorldID {
					mismatches = append(mismatches, fmt.Sprintf("worldId changed from %q to %q", existing.WorldID.String, metadata.WorldID))
				}
				if existing.SystemID.Valid && metadata.SystemID != "" && existing.SystemID.String != metadata.SystemID {
					mismatches = append(mismatches, fmt.Sprintf("systemId changed from %q to %q", existing.SystemID.String, metadata.SystemID))
				}
				if existing.FoundryVersion.Valid && metadata.FoundryVersion != "" && existing.FoundryVersion.String != metadata.FoundryVersion {
					mismatches = append(mismatches, fmt.Sprintf("foundryVersion changed from %q to %q", existing.FoundryVersion.String, metadata.FoundryVersion))
				}
				if len(mismatches) > 0 {
					metadataMatch = false
					flagged = true
					flagReason = "metadata mismatch: " + strings.Join(mismatches, "; ")
					log.Warn().Str("clientId", clientID).Str("reason", flagReason).Msg("Connection metadata mismatch detected")
				}
			}

			db.KnownClientStore().Upsert(ctx, &model.KnownClient{
				UserID:         user.ID,
				ClientID:       clientID,
				WorldID:        sql.NullString{String: metadata.WorldID, Valid: metadata.WorldID != ""},
				WorldTitle:     sql.NullString{String: metadata.WorldTitle, Valid: metadata.WorldTitle != ""},
				SystemID:       sql.NullString{String: metadata.SystemID, Valid: metadata.SystemID != ""},
				SystemTitle:    sql.NullString{String: metadata.SystemTitle, Valid: metadata.SystemTitle != ""},
				SystemVersion:  sql.NullString{String: metadata.SystemVersion, Valid: metadata.SystemVersion != ""},
				FoundryVersion: sql.NullString{String: metadata.FoundryVersion, Valid: metadata.FoundryVersion != ""},
				CustomName:     sql.NullString{String: metadata.CustomName, Valid: metadata.CustomName != ""},
				IsOnline:       true,
			})

			// Detect same worldId connecting under a different user account (multi-account abuse signal).
			if metadata.WorldID != "" {
				if conflict, err := db.KnownClientStore().FindByWorldIDCrossUser(ctx, metadata.WorldID, user.ID); err == nil && conflict != nil {
					if alerts.Track("world_cross:"+metadata.WorldID, 1, 24*time.Hour, 24*time.Hour) {
						alerts.Fire(alerts.Event{
							Type:     alerts.TypeWorldIDCrossAccount,
							Severity: "warning",
							Message:  "Foundry world connected under a different user account than previously seen",
							Details:  map[string]interface{}{"worldId": metadata.WorldID, "currentUserId": user.ID, "previousUserId": conflict.UserID},
						})
					}
				}
			}

			// Create connection log with metadata fingerprint result
			db.ConnectionLogStore().Create(ctx, &model.ConnectionLog{
				UserID:         user.ID,
				ClientID:       clientID,
				TokenName:      sql.NullString{String: metadata.TokenName, Valid: metadata.TokenName != ""},
				IPAddress:      sql.NullString{String: metadata.IPAddress, Valid: metadata.IPAddress != ""},
				WorldID:        sql.NullString{String: metadata.WorldID, Valid: metadata.WorldID != ""},
				WorldTitle:     sql.NullString{String: metadata.WorldTitle, Valid: metadata.WorldTitle != ""},
				SystemID:       sql.NullString{String: metadata.SystemID, Valid: metadata.SystemID != ""},
				FoundryVersion: sql.NullString{String: metadata.FoundryVersion, Valid: metadata.FoundryVersion != ""},
				MetadataMatch:  metadataMatch,
				Flagged:        flagged,
				FlagReason:     sql.NullString{String: flagReason, Valid: flagReason != ""},
			})

			if flagged {
				alerts.Fire(alerts.Event{
					Type:     alerts.TypeFlaggedConnection,
					Severity: "warning",
					Message:  "Connection flagged due to metadata mismatch",
					Details:  map[string]interface{}{"userId": user.ID, "clientId": clientID, "reason": flagReason},
				})
				if alerts.Track("meta:"+clientID, 3, 10*time.Minute, 1*time.Hour) {
					alerts.Fire(alerts.Event{
						Type:     alerts.TypeMetadataMismatchSpike,
						Severity: "warning",
						Message:  "3+ metadata mismatches from the same client in 10 minutes",
						Details:  map[string]interface{}{"userId": user.ID, "clientId": clientID},
					})
				}
			}

			// Load user notification settings to apply per-user batcher window.
			if notifSettings, err := db.NotificationSettingsStore().FindByUser(ctx, user.ID); err == nil && notifSettings != nil {
				if notifSettings.RemoteRequestBatchWindowSecs > 0 {
					s.RemoteRequestBatcher.SetUserWindow(user.ID, time.Duration(notifSettings.RemoteRequestBatchWindowSecs)*time.Second)
				}
			}

			// If this client has never been seen before (no existing KnownClients row),
			// dispatch a "new world connected" event in addition to the regular connect.
			if existing == nil {
				s.Dispatcher.Dispatch(service.NotificationContext{
					Event:      service.EventNewClientConnect,
					UserID:     user.ID,
					ClientID:   clientID,
					WorldID:    metadata.WorldID,
					WorldTitle: metadata.WorldTitle,
					SystemID:   metadata.SystemID,
					IPAddress:  metadata.IPAddress,
				})
			}

			// Dispatch connect notification via unified dispatcher
			s.Dispatcher.Dispatch(service.NotificationContext{
				Event:      service.EventConnect,
				UserID:     user.ID,
				ClientID:   clientID,
				WorldID:    metadata.WorldID,
				WorldTitle: metadata.WorldTitle,
				SystemID:   metadata.SystemID,
				IPAddress:  metadata.IPAddress,
			})

			// If metadata mismatch was detected, also fire a metadata-mismatch event
			// (separate notification so users can opt-in/out independently of regular connects)
			if flagged {
				s.Dispatcher.Dispatch(service.NotificationContext{
					Event:      service.EventMetadataMismatch,
					UserID:     user.ID,
					ClientID:   clientID,
					WorldID:    metadata.WorldID,
					WorldTitle: metadata.WorldTitle,
					SystemID:   metadata.SystemID,
					IPAddress:  metadata.IPAddress,
					Reason:     flagReason,
					Severity:   "alert",
				})
			}
		}()
	}

	// Register WS message handlers with SSE fan-out
	fanout := &ws.EventFanout{
		OnChatEvent: func(clientID string, data map[string]interface{}) {
			s.fanoutChatEvent(clientID, data)
		},
		OnRollData: func(clientID string, data map[string]interface{}) {
			s.fanoutRollData(clientID, data)
		},
		OnHookEvent: func(clientID string, data map[string]interface{}) {
			s.fanoutHookEvent(clientID, data)
		},
		OnCombatEvent: func(clientID string, data map[string]interface{}) {
			s.fanoutCombatEvent(clientID, data)
		},
		OnActorEvent: func(clientID string, data map[string]interface{}) {
			s.fanoutActorEvent(clientID, data)
		},
		OnSceneEvent: func(clientID string, data map[string]interface{}) {
			s.fanoutSceneEvent(clientID, data)
		},
	}
	ws.RegisterMessageHandlers(s.ClientManager, s.PendingReqs, fanout, s.InteractiveSessions)

	// Cross-world tunnel: source-module remote-request → relay → target client.
	// Gated by source connection token's allowedTargetClients + remoteScopes.
	// The Headless field is passed via interface so the handler can ask for
	// auto-start when a target is offline (currently stubbed; see Task #85).
	var headlessAutoStarter ws.HeadlessAutoStarter
	if s.Headless != nil {
		headlessAutoStarter = s.Headless
	}
	// Notification batcher for cross-world remote-request events. Batches
	// multiple calls within a 5-minute window into one summary notification.
	// Stored on the Server so OnClientConnected can update per-user windows.
	remoteRequestBatcher := service.NewRemoteRequestBatcher(s.Dispatcher, 5*time.Minute)
	s.RemoteRequestBatcher = remoteRequestBatcher

	ws.RegisterRemoteRequestHandler(ws.RemoteRequestConfig{
		Manager:  s.ClientManager,
		Pending:  s.PendingReqs,
		DB:       s.db,
		Headless: headlessAutoStarter,
		Batcher:  remoteRequestBatcher,
	})

	// Known-clients: lets a connected Foundry module query all worlds on its
	// relay account (online + offline) to populate the transfer target picker.
	ws.RegisterKnownClientsHandler(ws.KnownClientsConfig{
		Manager: s.ClientManager,
		DB:      s.db,
	})

	// Module-notify: Foundry module reports in-Foundry events (settings change,
	// execute-js, macro-execute) that the relay can't observe directly. The
	// callback resolves the user via API key and routes through the unified
	// notification dispatcher.
	ws.RegisterModuleNotifyHandler(s.ClientManager, func(event ws.ModuleNotifyEvent) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// event.APIKey is the WS Client's registration token, which is
			// the user's apiKeyHash. Look up by hash directly (don't re-hash).
			user, _ := db.UserStore().FindByAPIKeyHash(ctx, event.APIKey)
			if user == nil {
				return
			}

			// Map module event name to dispatcher event constant.
			// execute-js and macro-execute are now dispatched directly from the
			// HTTP handler (with full ApiKeyID context). Only settings-change
			// is genuinely unobservable from the HTTP layer.
			var dispEvent service.NotificationEvent
			switch event.EventName {
			case "settings-change":
				dispEvent = service.EventSettingsChange
			default:
				return
			}

			desc := event.Details
			if event.Actor != "" && desc != "" {
				desc = "**Triggered by:** " + event.Actor + "\n" + desc
			} else if event.Actor != "" {
				desc = "**Triggered by:** " + event.Actor
			}

			s.Dispatcher.Dispatch(service.NotificationContext{
				Event:       dispEvent,
				UserID:      user.ID,
				ClientID:    event.ClientID,
				WorldTitle:  event.World,
				Description: desc,
			})

			// Persist the event to the module event log for the activity feed.
			_ = db.ModuleEventLogStore().Create(ctx, &model.ModuleEventLog{
				UserID:      user.ID,
				ClientID:    event.ClientID,
				WorldTitle:  event.World,
				EventType:   event.EventName,
				Actor:       event.Actor,
				Description: desc,
			})
		}()
	})

	// player-list: Foundry module pushes its full user list on auth-success and on
	// user CRUD hooks. Store the list in KnownUsers so scoped-key per-client user
	// selection works even when the world is offline.
	s.ClientManager.OnMessageType("player-list", func(client *ws.Client, data map[string]interface{}) {
		raw, _ := json.Marshal(data["users"])
		var incoming []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Role int    `json:"role"`
		}
		if err := json.Unmarshal(raw, &incoming); err != nil || len(incoming) == 0 {
			return
		}
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			user, err := db.UserStore().FindByAPIKeyHash(ctx, client.APIKey())
			if err != nil || user == nil {
				return
			}
			known, err := db.KnownClientStore().FindByClientID(ctx, user.ID, client.ID())
			if err != nil || known == nil {
				return
			}
			users := make([]*model.KnownUser, 0, len(incoming))
			for _, u := range incoming {
				users = append(users, &model.KnownUser{
					KnownClientID: known.ID,
					UserID:        u.ID,
					Name:          u.Name,
					Role:          u.Role,
				})
			}
			if err := db.KnownUserStore().UpsertAll(ctx, known.ID, users); err != nil {
				log.Warn().Err(err).Str("clientId", client.ID()).Msg("Failed to upsert player list")
			} else {
				log.Debug().Int("count", len(users)).Str("clientId", client.ID()).Msg("Stored player list")
			}
		}()
	})

	// Start background cleanup loops using the application shutdown context so
	// all goroutines stop cleanly when the server receives SIGTERM/SIGINT.
	s.ClientManager.StartCleanupLoop(s.ctx, time.Duration(cfg.ClientCleanupIntervalMs)*time.Millisecond)
	// Subscribe to the cross-instance disconnect channel so that token
	// revocations and known-client deletes kill live sessions pinned to
	// any Fly.io instance, not just the one handling the HTTP request.
	s.ClientManager.StartDisconnectSubscriber(s.ctx)
	// Safety net: periodically walk local clients and validate their
	// credentials against the database. Catches any disconnect broadcast
	// that didn't deliver (e.g., transient Redis outage) or any bug in the
	// broadcast path. The DB is the source of truth for credential state.
	s.startReconciliationLoop(s.ctx, 30*time.Second)
	s.PendingReqs.StartCleanupLoop(s.ctx, 30*time.Second, 30*time.Second)
	s.InteractiveSessions.StartCleanupLoop(s.ctx)
	if s.Headless != nil {
		s.Headless.StartCleanupLoop(s.ctx)
		go s.Headless.WarmUpBrowser()
	}

	// Set up auto-start using world-side credentials (KnownClient.credentialId)
	if s.Headless != nil {
		var autoStartMu sync.Mutex
		autoStartCooldowns := make(map[string]time.Time) // clientID -> cooldown until

		// Evict expired cooldown entries so the map doesn't grow unbounded.
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					autoStartMu.Lock()
					now := time.Now()
					for id, until := range autoStartCooldowns {
						if now.After(until) {
							delete(autoStartCooldowns, id)
						}
					}
					autoStartMu.Unlock()
				case <-s.ctx.Done():
					return
				}
			}
		}()

		autoStart := func(userID int64, targetClientID string) string {
			if targetClientID == "" {
				return ""
			}
			if s.Headless.IsLaunching(targetClientID) {
				log.Info().Str("clientId", targetClientID).Msg("Headless session launching, queuing request")
				clientID, err := s.Headless.WaitForLaunch(targetClientID, 5*time.Minute)
				if err != nil {
					log.Warn().Err(err).Str("clientId", targetClientID).Msg("Queued request: launch wait failed")
					return ""
				}
				return clientID
			}
			autoStartMu.Lock()
			if until, ok := autoStartCooldowns[targetClientID]; ok && time.Now().Before(until) {
				autoStartMu.Unlock()
				return ""
			}
			autoStartMu.Unlock()
			s.Headless.RegisterLaunch(targetClientID)
			startCtx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			defer cancel()
			clientID, err := s.Headless.AutoStartForKnownClient(startCtx, userID, targetClientID)
			if err != nil {
				log.Warn().Err(err).Str("clientId", targetClientID).Msg("Auto-start headless failed")
				s.Headless.FailLaunch(targetClientID, err)
				autoStartMu.Lock()
				autoStartCooldowns[targetClientID] = time.Now().Add(60 * time.Second)
				autoStartMu.Unlock()
				return ""
			}
			s.Headless.CompleteLaunch(targetClientID, clientID)
			return clientID
		}

		helpers.AutoStartFunc = func(reqCtx *helpers.RequestContext, targetClientID string) string {
			if reqCtx == nil || targetClientID == "" {
				return ""
			}
			user, ok := reqCtx.GetUser()
			if !ok {
				return ""
			}
			return autoStart(user.ID, targetClientID)
		}

		// Expose the same logic to the WebSocket client API path.
		s.headlessAutoStart = func(masterAPIKey, scopedClientID string) string {
			if scopedClientID == "" {
				return ""
			}
			user, err := s.db.UserStore().FindByAPIKeyHash(context.Background(), masterAPIKey)
			if err != nil || user == nil {
				return ""
			}
			return autoStart(user.ID, scopedClientID)
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

	// Resolve static file base directory early so it's available throughout
	// this function (e.g. for serving the admin SPA page inside the admin group).
	baseDir := os.Getenv("STATIC_DIR")
	if baseDir == "" {
		baseDir = ".." // Relative to go-relay directory (local dev)
	}

	// Global middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(handler.MaintenanceModeMiddleware)
	r.Use(appmw.MetricsMiddleware)
	r.Use(appmw.BodyLimit(int64(s.cfg.MaxUploadSizeMB) * 1024 * 1024))
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

	// Prometheus metrics endpoint — gated by admin IP allowlist (no admin JWT,
	// because Prometheus scrapers don't support cookies). In production, restrict
	// via ADMIN_ALLOWED_IPS or run on Fly internal network.
	r.With(appmw.AdminIPAllowlist(s.cfg.AdminAllowedIPs)).Handle("/metrics", promhttp.Handler())

	// Status endpoint — public, no auth. Also used by the frontend to detect
	// self-hosted mode (billingEnabled = false when STRIPE_SECRET_KEY is unset).
	r.Get("/api/status", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":          "ok",
			"version":         s.version,
			"websocket":       "/relay",
			"billingEnabled":  s.cfg.StripeSecretKey != "",
			"headlessEnabled": s.cfg.AllowHeadless,
		})
	})

	// Health endpoint
	r.Get("/api/health", s.healthHandler)

	// Active-connection probe — public, rate-limited, returns only a boolean.
	// Used by the Foundry module's init wizard to detect "is another GM
	// already holding the slot for this clientId?" before deciding whether to
	// prompt for pairing. The clientId is non-secret (it's stored in world
	// settings, readable by every Foundry player), so this leaks no
	// confidential data — just a single bit of liveness state.
	//
	// Returns 404 if the clientId doesn't exist in KnownClients (we don't
	// distinguish "not present" from "exists but offline" beyond that, to
	// minimize information leakage about which clientIds exist).
	r.With(appmw.ProbeRateLimiter.Middleware).Get("/api/clients/{clientId}/active",
		handler.ClientActiveHandler(s.db, s.ClientManager))

	// Self-unpair — called by the Foundry module when the GM clicks Unpair.
	// Authenticates via the connection token in the request body (not the
	// master API key), deletes the token, and cleans up KnownClients if no
	// other tokens remain for that world.
	r.Post("/api/self-unpair", handler.SelfUnpairHandler(s.db, s.ClientManager))

	// Auth routes
	r.Mount("/auth", handler.AuthRouter(s.db, s.cfg, s.ClientManager))

	// Admin dashboard routes — gated by IP allowlist; auth/CSRF via JWT cookies inside.
	// AdminCSP is applied per-sub-route inside AdminRouter so that the HTML page
	// response is not subject to it (Astro's inline hydration scripts must run).
	r.Group(func(adm chi.Router) {
		adm.Use(appmw.AdminIPAllowlist(s.cfg.AdminAllowedIPs))
		adm.Mount("/admin", handler.AdminRouter(&handler.AdminDeps{
			DB:            s.db,
			Cfg:           s.cfg,
			ClientManager: s.ClientManager,
			PendingReqs:   s.PendingReqs,
			Headless:      s.Headless,
			Redis:         s.redis,
			Version:       s.version,
			// Serve the admin SPA at GET /admin — handled inside AdminRouter at "/"
			// because Chi's Mount intercepts the exact /admin path before an outer
			// Get("/admin") can match it in the trie.
			AdminPageHandler: func(w http.ResponseWriter, r *http.Request) {
				p := filepath.Join(baseDir, "public-dist/admin/index.html")
				if _, err := os.Stat(p); err == nil {
					http.ServeFile(w, r, p)
					return
				}
				http.NotFound(w, r)
			},
		}))
	})

	// API routes — auth middleware needed by Stripe and API routes
	authMw := appmw.AuthMiddleware(s.db, s.ClientManager)
	usageMw := appmw.TrackAPIUsage(s.cfg.MonthlyRequestLimit, s.db)

	// Stripe routes — subscriptions require auth, webhooks do not
	r.Route("/api/subscriptions", func(sub chi.Router) {
		sub.Use(authMw)
		sub.Mount("/", handler.StripeRouter(s.db, s.cfg))
	})
	r.Mount("/api/webhooks", handler.WebhookRouter(s.db, s.cfg))

	// WebSocket routes — registered before API route group to avoid catch-all conflict
	relayCfg := &ws.RelayConfig{
		PingInterval:             time.Duration(s.cfg.WSPingIntervalMs) * time.Millisecond,
		CleanupInterval:          time.Duration(s.cfg.ClientCleanupIntervalMs) * time.Millisecond,
		ValidateAPIKey:           service.MakeWSValidateAPIKey(s.db),
		ValidateConnectionToken:  service.MakeWSValidateConnectionToken(s.db),
		ValidateHeadless: func(clientID, token string) (bool, error) {
			if s.Headless == nil {
				return false, nil
			}
			return s.Headless.ValidateHeadlessSession(clientID, token)
		},
		// On (re)connect, send the module the current subscription counts so it
		// can re-enable any channels that already have active subscribers.
		OnClientConnected: func(clientID string) {
			channels := []string{"chat-events", "roll-events", "hooks", "combat-events", "actor-events", "scene-events"}
			for _, ch := range channels {
				count := s.SSEManager.TotalForChannel(clientID, ch)
				if count > 0 {
					s.notifyModuleChannelCount(clientID, ch, count)
				}
			}
		},
	}
	relayHandler := ws.HandleRelayConnection(s.ClientManager, relayCfg)
	r.HandleFunc("/relay", relayHandler)
	r.HandleFunc("/relay/", relayHandler)

	// Client API WebSocket
	clientWSHandler := ws.HandleClientAPIConnection(s.ClientManager, s.PendingReqs, &ws.ClientAPIConfig{
		PingInterval:            time.Duration(s.cfg.WSPingIntervalMs) * time.Millisecond,
		ValidateAPIKey:          service.MakeWSValidateAPIKey(s.db),
		ValidateConnectionToken: service.MakeWSValidateConnectionToken(s.db),
		TrackUsage: func(apiKey string) (bool, string) {
			return service.TrackWSAPIUsage(context.Background(), s.db, apiKey)
		},
		AutoStart: func(masterAPIKey, scopedClientID, _ string) string {
			if s.headlessAutoStart == nil {
				return ""
			}
			return s.headlessAutoStart(masterAPIKey, scopedClientID)
		},
		SSEManager:     s.SSEManager,
		InteractiveSessions:  s.InteractiveSessions,
	})
	r.HandleFunc("/ws/api", clientWSHandler)

	// Authenticated API routes
	handler.RegisterAPIRoutes(r, s.ClientManager, s.PendingReqs, s.cfg, s.db, s.SSEManager, s.Headless, s.Dispatcher, authMw, usageMw)

	// Static file serving

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

	// Static assets (Astro build output)
	if d := findStaticDir(baseDir, "public-dist/_astro"); d != "" {
		r.Handle("/_astro/*", cacheControl("public, max-age=31536000, immutable", http.StripPrefix("/_astro/", http.FileServer(http.Dir(d)))))
	}

	// Astro frontend pages
	r.Get("/privacy", handler.ServeStaticPage(
		filepath.Join(baseDir, "public-dist/privacy/index.html"),
	))
	r.Get("/subscription-success", handler.ServeStaticPage(
		filepath.Join(baseDir, "public-dist/subscription-success/index.html"),
	))
	r.Get("/subscription-cancel", handler.ServeStaticPage(
		filepath.Join(baseDir, "public-dist/subscription-cancel/index.html"),
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
		docsHandler := handler.DocsFileServer(docsDir)
		r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=3600")
			docsHandler(w, r)
		})
		r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=3600")
			docsHandler(w, r)
		})
	}

	// Asset proxy
	r.Get("/proxy-asset/*", handler.ProxyAssetHandler())

	// SPA deep-link routes — serve index.html so the Astro frontend can
	// read window.location.pathname and route to the right view.
	spaHandler := func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Join(baseDir, "public-dist/index.html")
		if _, err := os.Stat(p); err == nil {
			http.ServeFile(w, r, p)
			return
		}
		http.NotFound(w, r)
	}
	r.Get("/approve/{code}", spaHandler)
	r.Get("/pair/{code}", spaHandler)

	// Root — WebSocket upgrade or Astro frontend
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {
			relayHandler(w, r)
			return
		}
		// Serve Astro frontend
		p := filepath.Join(baseDir, "public-dist/index.html")
		if _, err := os.Stat(p); err == nil {
			http.ServeFile(w, r, p)
			return
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

// notifyModuleChannelCount sends an event-subscription-update message to the
// Foundry module so it can enable or disable hooks for the given channel.
func (s *Server) notifyModuleChannelCount(clientID, channel string, count int) {
	client := s.ClientManager.GetClient(clientID)
	if client == nil {
		return
	}
	client.Send(map[string]interface{}{
		"type":    "event-subscription-update",
		"channel": channel,
		"count":   count,
	})
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

	// Fan out to SSE connections — each write in its own goroutine so a
	// slow subscriber doesn't block delivery to other subscribers.
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

		conn := conn
		go func() {
			select {
			case <-conn.Done:
				return
			default:
				fmt.Fprintf(conn.W, "event: %s\ndata: %s\n\n", eventType, jsonStr)
				conn.Flusher.Flush()
			}
		}()
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

		conn := conn
		go func() {
			select {
			case <-conn.Done:
				return
			default:
				fmt.Fprintf(conn.W, "event: roll\ndata: %s\n\n", jsonStr)
				conn.Flusher.Flush()
			}
		}()
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

// fanoutHookEvent sends generic Foundry hook events to hooks SSE and WS subscribers.
func (s *Server) fanoutHookEvent(clientID string, data map[string]interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	jsonStr := string(jsonBytes)

	hookName, _ := data["hook"].(string)
	if hookName == "" {
		if innerData, ok := data["data"].(map[string]interface{}); ok {
			hookName, _ = innerData["hook"].(string)
		}
	}

	// Fan out to generic SSE connections
	for _, conn := range s.SSEManager.GetGenericSSE(clientID) {
		if len(conn.HookFilters) > 0 {
			matched := false
			for _, f := range conn.HookFilters {
				if f == hookName {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		conn := conn
		go func() {
			select {
			case <-conn.Done:
				return
			default:
				fmt.Fprintf(conn.W, "event: %s\ndata: %s\n\n", hookName, jsonStr)
				conn.Flusher.Flush()
			}
		}()
	}

	// Fan out to WS event connections subscribed to "hooks"
	for _, conn := range s.SSEManager.GetWSEvents(clientID) {
		if conn.Channel != "hooks" {
			continue
		}
		if len(conn.HookFilters) > 0 {
			matched := false
			for _, f := range conn.HookFilters {
				if f == hookName {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		conn.SendFunc(map[string]interface{}{"type": "hook-event", "hook": hookName, "data": data})
	}
}

// fanoutCombatEvent sends combat events to combat SSE and WS subscribers.
func (s *Server) fanoutCombatEvent(clientID string, data map[string]interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	jsonStr := string(jsonBytes)

	eventType, _ := data["eventType"].(string)
	encounterId, _ := data["encounterId"].(string)

	for _, conn := range s.SSEManager.GetCombatSSE(clientID) {
		if conn.EncounterID != "" && conn.EncounterID != encounterId {
			continue
		}
		conn := conn
		go func() {
			select {
			case <-conn.Done:
				return
			default:
				fmt.Fprintf(conn.W, "event: combat-%s\ndata: %s\n\n", eventType, jsonStr)
				conn.Flusher.Flush()
			}
		}()
	}

	for _, conn := range s.SSEManager.GetWSEvents(clientID) {
		if conn.Channel != "combat-events" {
			continue
		}
		conn.SendFunc(map[string]interface{}{"type": "combat-event", "event": "combat-" + eventType, "data": data})
	}
}

// fanoutActorEvent sends actor events to actor SSE and WS subscribers.
func (s *Server) fanoutActorEvent(clientID string, data map[string]interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	jsonStr := string(jsonBytes)

	actorUUID, _ := data["actorUuid"].(string)
	eventType, _ := data["eventType"].(string)

	for _, conn := range s.SSEManager.GetActorSSE(clientID) {
		if conn.ActorUUID != "" && conn.ActorUUID != actorUUID {
			continue
		}
		conn := conn
		go func() {
			select {
			case <-conn.Done:
				return
			default:
				fmt.Fprintf(conn.W, "event: actor-%s\ndata: %s\n\n", eventType, jsonStr)
				conn.Flusher.Flush()
			}
		}()
	}

	for _, conn := range s.SSEManager.GetWSEvents(clientID) {
		if conn.Channel != "actor-events" {
			continue
		}
		if conn.ActorUUID != "" && conn.ActorUUID != actorUUID {
			continue
		}
		conn.SendFunc(map[string]interface{}{"type": "actor-event", "event": "actor-" + eventType, "data": data})
	}
}

// fanoutSceneEvent sends scene events to scene SSE and WS subscribers.
func (s *Server) fanoutSceneEvent(clientID string, data map[string]interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	jsonStr := string(jsonBytes)

	sceneID, _ := data["sceneId"].(string)
	eventType, _ := data["eventType"].(string)

	for _, conn := range s.SSEManager.GetSceneSSE(clientID) {
		if conn.SceneID != "" && conn.SceneID != sceneID {
			continue
		}
		conn := conn
		go func() {
			select {
			case <-conn.Done:
				return
			default:
				fmt.Fprintf(conn.W, "event: scene-%s\ndata: %s\n\n", eventType, jsonStr)
				conn.Flusher.Flush()
			}
		}()
	}

	for _, conn := range s.SSEManager.GetWSEvents(clientID) {
		if conn.Channel != "scene-events" {
			continue
		}
		if conn.SceneID != "" && conn.SceneID != sceneID {
			continue
		}
		conn.SendFunc(map[string]interface{}{"type": "scene-event", "event": "scene-" + eventType, "data": data})
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

// cacheControl wraps an http.Handler to add Cache-Control headers.
func cacheControl(maxAge string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", maxAge)
		handler.ServeHTTP(w, r)
	})
}

// startReconciliationLoop periodically validates every locally-connected
// WebSocket client against the database. Any client whose master API key
// or connection token has been revoked is force-disconnected with a 4002
// close frame.
//
// This is the safety net behind the Redis pub/sub disconnect broadcast.
// Pub/sub is fire-and-forget: if Redis is briefly unavailable on the
// subscriber side, messages are lost. The DB is the source of truth, so
// reconciliation catches anything the broadcast missed.
//
// Each instance only reconciles clients it owns locally — sharded by
// connection, not by user — so the combined work is O(connected clients)
// across the whole fleet regardless of instance count.
func (s *Server) startReconciliationLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.reconcileConnections(ctx)
			}
		}
	}()
}

func (s *Server) reconcileConnections(ctx context.Context) {
	clients := s.ClientManager.SnapshotLocalClients()
	if len(clients) == 0 {
		return
	}

	disconnected := 0
	for _, c := range clients {
		if !c.IsAlive() {
			continue
		}
		// Use a short per-client context so one slow DB query doesn't
		// block the whole reconciliation pass.
		lookupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		// Check 1: master API key still valid (catches rotation + account deletion).
		// c.APIKey() returns the registration token (apiKeyHash), so we use
		// the hash variant to avoid double-hashing.
		user, err := s.db.UserStore().FindByAPIKeyHash(lookupCtx, c.APIKey())
		if err != nil {
			cancel()
			log.Warn().Err(err).Str("clientId", c.ID()).Msg("Reconciliation: FindByAPIKey failed, skipping client")
			continue
		}
		if user == nil {
			cancel()
			log.Info().Str("clientId", c.ID()).Msg("Reconciliation: master API key no longer valid, disconnecting")
			s.ClientManager.ForceDisconnectLocal(c, 4002, "Master API key revoked")
			s.ClientManager.RemoveClient(c.ID())
			disconnected++
			continue
		}

		// Check 2: if the client authenticated with a connection token,
		// that token must still exist in the DB.
		if tokenID := c.ConnectionTokenID(); tokenID != 0 {
			token, err := s.db.ConnectionTokenStore().FindByID(lookupCtx, tokenID)
			if err != nil {
				cancel()
				log.Warn().Err(err).Str("clientId", c.ID()).Int64("tokenId", tokenID).Msg("Reconciliation: FindByID failed, skipping client")
				continue
			}
			if token == nil || token.UserID != user.ID {
				cancel()
				log.Info().Str("clientId", c.ID()).Int64("tokenId", tokenID).Msg("Reconciliation: connection token revoked, disconnecting")
				s.ClientManager.ForceDisconnectLocal(c, 4002, "Connection token revoked")
				s.ClientManager.RemoveClient(c.ID())
				disconnected++
				continue
			}
		}

		cancel()
	}

	if disconnected > 0 {
		log.Info().Int("disconnected", disconnected).Int("scanned", len(clients)).Msg("Reconciliation pass complete")
	}
}

// import guard
var _ = fmt.Sprintf
