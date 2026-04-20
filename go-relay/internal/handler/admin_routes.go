package handler

import (
	"net/http"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	appmw "github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// AdminDeps holds the dependencies needed by the admin router.
type AdminDeps struct {
	DB              *database.DB
	Cfg             *config.Config
	ClientManager   *ws.ClientManager
	PendingReqs     *ws.PendingRequests
	Headless        *worker.HeadlessManager
	Redis           *config.RedisClient
	Version         string
	AdminPageHandler http.HandlerFunc // serves public-dist/admin/index.html at GET /admin
}

// AdminRouter mounts all /admin routes — both unauthenticated (login) and protected.
// Mounted by server.go under /admin with IP allowlist already applied.
func AdminRouter(deps *AdminDeps) chi.Router {
	r := chi.NewRouter()

	// Serve the admin SPA page at the root of this sub-router (i.e. GET /admin).
	// No CSP here — Astro's inline bootstrap scripts and the BaseLayout theme
	// script must run; applying script-src 'self' would block them and prevent
	// Svelte from hydrating (no login, no dark mode).
	if deps.AdminPageHandler != nil {
		r.Get("/", deps.AdminPageHandler)
	}

	// Resolve admin auth config — JWT secret is auto-generated in dev if missing.
	authCfg := buildAdminAuthConfig(deps.Cfg)

	// Apply strict CSP to all API/auth endpoints but NOT to the HTML page above.
	r.Group(func(api chi.Router) {
		api.Use(appmw.AdminCSP)

		// Auth endpoints (login is the entry point, no RequireAdmin)
		api.Mount("/auth", AdminAuthRouter(deps.DB, authCfg))

		// Protected admin API routes (each handler in its own router file)
		api.Route("/api", func(apiRoutes chi.Router) {
			// Test report: no JWT required — still behind the IP allowlist at the server level.
			apiRoutes.Get("/tests/report", TestReportHandler)

			// Everything else requires admin auth.
			apiRoutes.Group(func(protected chi.Router) {
				protected.Use(appmw.RequireAdmin(deps.DB, authCfg))
				protected.Mount("/users", AdminUsersRouter(deps.DB))
				protected.Mount("/keys", AdminKeysRouter(deps.DB))
				protected.Mount("/clients", AdminClientsRouter(deps.DB, deps.ClientManager))
				protected.Mount("/audit-logs", AdminAuditRouter(deps.DB))
				protected.Mount("/activity", AdminActivityRouter(deps.DB))
				protected.Mount("/headless-sessions", AdminSessionsRouter(deps.DB, deps.Headless))
				protected.Mount("/system/health", AdminHealthRouter(deps.Cfg, deps.ClientManager, deps.PendingReqs, deps.Headless, deps.Redis, deps.Version))
				protected.Mount("/ops", AdminOpsRouter(deps.DB, deps.ClientManager))
				protected.Mount("/subscriptions", AdminSubscriptionsRouter(deps.DB))
				protected.Mount("/metrics", AdminMetricsRouter())
				protected.Mount("/tests", AdminTestRunnerRouter(deps.DB, deps.Cfg))
				protected.Mount("/alerts", AdminAlertsRouter(deps.DB))
			})
		})
	})

	return r
}

// buildAdminAuthConfig produces an AdminAuthConfig from the loaded config.
// ADMIN_JWT_SECRET is required for all backends — Validate() enforces this before we get here.
func buildAdminAuthConfig(cfg *config.Config) *appmw.AdminAuthConfig {
	secret := cfg.AdminJWTSecret
	if secret == "" {
		// Validate() should have caught this — last-resort guard.
		log.Fatal().Msg("ADMIN_JWT_SECRET is required — generate one with: openssl rand -base64 32")
	}
	maxHrs := cfg.AdminSessionMaxHours
	if maxHrs <= 0 {
		maxHrs = 8
	}
	isProd := cfg.IsProduction()
	// In production, short-lived tokens with sliding refresh. In dev, long sessions
	// so the operator doesn't need to re-login during a work session.
	tokenTTL := appmw.DefaultAccessTTL
	if !isProd {
		tokenTTL = time.Duration(maxHrs) * time.Hour
	}
	return &appmw.AdminAuthConfig{
		JWTSecret:     []byte(secret),
		SessionMaxHrs: maxHrs,
		TokenTTL:      tokenTTL,
		IsDevelopment: !isProd,
		SecureCookies: cfg.AdminSecureCookies,
	}
}
