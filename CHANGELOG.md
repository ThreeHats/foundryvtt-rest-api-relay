# Changelog

## [3.1.1] - 2026-05-06

- Updated deploy workflow version bump step to update docgen version number

## [3.1.0] - 2026-05-04

### Added
- Integration test runner now samples CPU, RAM, GPU utilization, and VRAM every 3 seconds during the test run and prints a formatted resource usage report (with bar charts and peak/average/min stats) after all tests complete, alongside total test suite and test pass/fail counts
- `POST /dnd5e/modify-currency` — add or remove currency from a single actor (use a negative `amount` to deduct)
- `POST /dnd5e/prepare-spell` — toggle a spell's prepared state; only applies to spellcasters that prepare spells
- `ownedByUserId` query parameter on entity search — filter results to documents the specified Foundry user owns (by ID or username)
- `regions` canvas document type — full CRUD for Scene Regions (Foundry v12+); `regions` is the replacement for `templates` on Foundry v14+

### Changed
- Foundry module `module.json` `verified` updated from `"13"` to `"14"`
- `GET /clients` now includes offline known-client records alongside live connections; each entry has an `isOnline` field
- `CREDENTIALS_ENCRYPTION_KEY` and `ADMIN_JWT_SECRET` are no longer required at startup — both are auto-generated on first run and persisted to `<DATA_DIR>/.secrets.env` (mode 0600). Environment variables still take precedence when set. If the encryption key is lost, users will need to re-enter any stored Foundry credentials.
- Docker Compose files (`docker-compose.local.yml`, `docker-compose.postgres.yml`) no longer require a `.env` file to start
- NVIDIA headless rendering now uses ANGLE Vulkan + native NVIDIA Vulkan ICD (`--use-angle=vulkan --use-vulkan=native --ozone-platform=headless`) — provides true hardware-accelerated WebGL; VirtualGL and all its dependencies have been removed from the image
- `NVIDIA_VISIBLE_DEVICES=all` and `NVIDIA_DRIVER_CAPABILITIES=graphics,utility` are now baked into the Dockerfile `ENV` — NVIDIA Option A users no longer need to add these variables to their compose file; the NVIDIA container runtime picks them up automatically at runtime and they have no effect when NVIDIA is not active
- `go-relay/Dockerfile` now explicitly installs `libvulkan1`, `mesa-vulkan-drivers`, `xvfb`, `libgl1-mesa-dri`, `libegl-mesa0`, `libegl1`, and `libglu1-mesa` to cover all rendering paths (NVIDIA Vulkan, Intel/AMD Mesa DRI, Xvfb software)

### Fixed
- `CONFIG.statusEffects` is now read correctly on Foundry v14, which changed the value from an Array to an Object keyed by ID; the `kill` entity action now correctly applies the dead/defeated status effect on v14
- Reconnects from the same IP address (e.g. after a network drop) no longer receive a `DuplicateConnection` (4004) rejection; the stale connection is silently replaced
- Headless session startup now polls for world-list readiness instead of checking once — fixes intermittent launch failures on Foundry v14 where the world list is rendered by JavaScript after initial page load
- Stale `isOnline` flags left by a crash or unclean shutdown are now cleared at startup on single-instance deployments
- Email verification page now renders styled HTML with distinct messaging for each state (expired link, invalid token, already verified, success)
- Admin users panel now shows each account's email verification status
- Dashboard Connections page memory leak: `tokenRow` snippet was defined inside the `{#each clients}` loop, causing a new closure to be created per client on every re-render; moved to component scope with an explicit `activeTokenId` parameter
- Dashboard Connections page: Map and Set reactive state (`worldNotifSettings`, `worldNotifOpen`, `worldNotifSaving`, `worldNotifTesting`, `worldNotifMsg`, `selectedTokenIds`, `expandedInactive`) now mutated in-place via Svelte 5 reactive proxy methods (`.set()`, `.add()`, `.delete()`, `.clear()`) instead of create-copy-assign, eliminating intermediate object allocations on every state change
- Xvfb virtual display now starts at 24-bit color depth (was 16-bit); fixes OpenGL/GLX context creation failures that caused canvas operations to silently produce blank output
- `libegl1` (GLVND EGL dispatcher) and `libglu1-mesa` are now explicitly installed in the image, fixing NVIDIA Vulkan ICD initialization

---

## [3.0.3] - 2026-04-20

### Added
- Cross-instance HTTP request forwarding: requests arriving at a Fly.io instance that doesn't have the target client connected are automatically forwarded to the correct instance

### Fixed
- SSE connections now correctly re-establish across multi-instance Fly.io deployments

---

## [3.0.2] - 2026-04-20

### Fixed
- Additional database stability fixes

---

## [3.0.1] - 2026-04-19

### Fixed
- Database schema migration errors on fresh 3.0.0 installs
- Email verification token hashing and validation corrected; resend flow fixed
- Existing accounts without a verification record are automatically marked verified on upgrade when SMTP is not configured

---

## [3.0.0] - 2026-04-19

### New Features

#### Scoped API Key Permission System
Every API endpoint is now gated by a named scope. Master keys pass all scope checks automatically; scoped keys must explicitly include the required scope.
- 32 scope constants defined: `entity:read`, `entity:write`, `roll:read`, `roll:execute`, `chat:read`, `chat:write`, `encounter:read`, `encounter:manage`, `macro:list`, `macro:execute`, `macro:write`, `scene:read`, `scene:write`, `canvas:read`, `canvas:write`, `effects:read`, `effects:write`, `user:read`, `user:write`, `file:read`, `file:write`, `playlist:control`, `world:info`, `clients:read`, `sheet:read`, `session:manage`, `events:subscribe`, `execute-js`, `search`, `dnd5e`, `structure:read`, `structure:write`
- `DefaultScopes` — safe read-only set assigned to new keys
- `ActionToScopeRequired` map — all WS action types mapped to their required scope
- All routes wrapped in `scope(model.Scope...)` groups
- `RequireScope` middleware enforces scopes on HTTP requests
- Scoped key form updated with full scope checklist and danger warnings

#### World Pairing & Connection Tokens
A new device-flow pairing model replaces direct API key sharing.
- **Pair Request flow**: Foundry module initiates pairing without credentials; user approves in dashboard
  - `POST /auth/pair-request` — module sends world metadata and requested scopes
  - `GET /auth/pair-request/{code}/status` — module polls for approval
  - `POST /auth/pair-request/{code}/approve` — user grants access, optionally with cross-world settings
  - `POST /auth/pair-request/{code}/deny` — user rejects
  - Supports `upgradeOnly` mode for updating existing world permissions
- **Connection Tokens** (replaces raw API key in WebSocket URL):
  - `POST /auth/connection-tokens` — generate token from pairing code
  - `GET /auth/connection-tokens` — list all tokens
  - `PATCH /auth/connection-tokens/{id}` — update cross-world settings, rate limits
  - `DELETE /auth/connection-tokens/{id}` — revoke and force-disconnect live sessions
  - Tokens stored hashed (never in plaintext); bound to a specific client/world
- **Self-unpair**: `POST /api/self-unpair` — module can disconnect itself and clean up its token
- **PairApprovalPage** in dashboard for users to review and approve pairing requests with cross-world configuration

#### Cross-World Remote Request Tunnel
Worlds in the same account can now execute actions against each other over the relay.
- New `remote-request` WebSocket action: source world sends a request targeting another world
- Permission checks: target must be in `AllowedTargetClients`, scope must be in `RemoteScopes`
- Per-token rate limiting (`RemoteRequestsPerHour`)
- Auto-start headless session on target world if it is offline and `AutoStartOnRemoteRequest` is enabled
- Full audit log (`RemoteRequestLog`) with source token, target client, action, success/failure
- `RemoteRequestsPage` in dashboard shows all cross-world activity

#### OAuth-Style Key Request Flow
Third-party apps can request API keys without pre-shared credentials.
- `POST /auth/key-request` — app submits request with requested scopes, client list, suggested limits
- `GET /auth/key-request/{code}/status` — app polls for user approval
- `GET /auth/key-request/{code}` — fetch request details for dashboard approval UI
- `POST /auth/key-request/{code}/approve` — user grants scopes and selects accessible worlds
- `POST /auth/key-request/{code}/deny`
- `POST /auth/key-request/exchange` — web-flow code exchange to retrieve the approved key
- `KeyApprovalPage` in dashboard with scope matrix, danger warnings, client selector, rate limit and expiry configuration

#### Stored Foundry Credentials
- Users can store encrypted Foundry server credentials (URL, username, password) using AES-GCM encryption
- Used by headless sessions for automatic world login
- `GET/POST/PUT/DELETE /auth/credentials` endpoints
- `CredentialsPage` in dashboard to manage saved credentials
- A world can be linked to a specific stored credential for headless auto-start

#### Known Clients Registry
- `KnownClient` records track paired Foundry worlds, their metadata, and cross-world settings
- `KnownUser` records track individual Foundry players per world (for scoped key references)
- `GET /auth/known-clients` — list all paired worlds with live/offline status
- `PATCH /auth/known-clients/{id}` — update cross-world tunneling settings (allowed targets, scopes, rate limit)
- `PATCH /auth/known-clients/{id}/credential` — link a stored credential for headless auto-start
- `PATCH /auth/known-clients/{id}/auto-start` — toggle auto-start on remote-request
- `DELETE /auth/known-clients/{id}` — unpair a world
- `GET /auth/known-clients/{id}/users` — list stored Foundry users for a world
- `GET/PUT/DELETE /auth/known-clients/{id}/notification-settings` — per-world notification config
- `get-known-clients` WebSocket action returns world list with `canAutoStart`, `tokenScopes`, `tokenAllowedTargets`
- `ConnectionsPage` in dashboard for managing paired worlds and cross-world permissions

#### Notification System
Users can receive Discord or email alerts for events in their connected worlds.
- Per-user account-level settings (`NotificationSettings`):
  - Toggle events: new client connect, connect/disconnect, metadata mismatch, settings change, execute-js, macro-execute, remote requests
  - Dual destination: Discord webhook URL and email address
  - Debounce window — suppress duplicate events within N seconds
  - Remote request batch window — group cross-world activity into a single alert per N minutes
- Per-API-key notification settings
- Per-world notification settings
- Cross-world requests are batched over a configurable window into a single summary notification
- Events dispatched to all configured destinations simultaneously
- `NotificationSettingsForm` in dashboard for configuring channels and toggles
- Test notification button to verify webhook/email config

#### Activity Log
- Unified event feed across all activity types: connections, cross-world requests, module events
- Filters: type, world/client ID, action name, failures only, date range, user ID (admin)
- Auto-refresh every 30 seconds, pagination (50 per page)
- `ActivityLog` component in dashboard
- Three underlying log tables: connection events, cross-world requests, module events

#### Module Event Auditing
- `module-notify` WebSocket message type — module sends security-relevant events
- Tracked events: `settings-change`, `execute-js`, `macro-execute`
- Stored in `ModuleEventLog` and routed to the notification system

#### Admin Dashboard
Full operator-only dashboard at `/admin` with JWT cookie authentication (separate from user auth).
- **AdminLogin** — email/password login, httpOnly cookie session, CSRF token
- **AdminUsers** — list/paginate all users; disable, enable, rotate key, delete; audit-logged
- **AdminClients** — live connected client list (auto-refresh 5s); force-disconnect
- **AdminKeys** — browse all API keys system-wide; enable/disable/revoke
- **AdminAuditLogs** — dual log view: user activity feed + admin action log
- **AdminMetrics** — per-minute/hour/day request counts, error totals, top endpoints, top users
- **AdminAlerts** — configure Discord webhook / email notification destination; enable/disable 20+ alert types by channel; view recent alert events
- **AdminHealth** — system status, WebSocket count, pending requests, goroutines, memory, GC stats, Redis status
- **AdminSessions** — headless Foundry session list; kill sessions
- **AdminSubscriptions** — Stripe subscription status summary and subscriber table (when billing enabled)
- **AdminOps** — maintenance mode toggle and other operational flags
- **AdminTestRunner** — run integration test suite or regenerate docs from the UI, with live ANSI-colored output via EventSource; edit `.env.test` variables in-browser

#### Admin Security Middleware
- `RequireAdmin` middleware — validates JWT cookie, applies CSRF, sliding session refresh, checks denylist
- `AdminIPAllowlist` — restrict `/admin` to whitelisted IPs/CIDRs
- `AdminLoginRateLimiter` — per-IP rate limiting on login endpoint
- `AdminLockoutTracker` — account lockout after repeated failures
- `JWTDenylistEntry` — immediate session revocation before natural JWT expiry
- `AdminCSP` — strict Content Security Policy headers for admin panel

#### Prometheus Metrics
- New `metrics/` package exposes a Prometheus-compatible endpoint
- Counters: HTTP requests by method/path/status, rate limit hits, auth failures by reason, admin login outcomes, WS round-trip timeouts
- Histograms: HTTP request duration by method/path (11 buckets: 5ms–10s)
- Gauges: active WS connections, pending requests, headless sessions, SSE connections
- `MetricsMiddleware` emits structured access logs with userId and key prefix

#### Rolling Window Metrics
- In-memory per-minute, per-hour, and per-day request counters
- Per-user and per-endpoint lifetime request counts since startup, surfaced in the admin dashboard

#### Playlist Endpoints
- `GET /playlists` — list all playlists
- `POST /playlist/play` — play a playlist
- `POST /playlist/stop` — stop a playlist
- `POST /playlist/next` — skip to next sound
- `POST /playlist/volume` — set playlist volume
- `POST /play-sound` — play a specific sound
- `POST /stop-sound` — stop a specific sound

#### User Management Endpoints
- `GET /users` — list Foundry VTT users in the connected world
- `GET /user` — fetch a single Foundry user (by query param: `id` or `name`)
- `POST /user` — create a Foundry user
- `PUT /user` — update a Foundry user
- `DELETE /user` — delete a Foundry user

#### Canvas Document Endpoints
Full CRUD access to all 8 Foundry canvas document types:
- `GET /canvas/{documentType}` — get canvas documents
- `POST /canvas/{documentType}` — create a canvas document
- `PUT /canvas/{documentType}` — update a canvas document
- `DELETE /canvas/{documentType}` — delete a canvas document
- Supported `documentType` values: `tokens`, `tiles`, `drawings`, `lights`, `sounds`, `notes`, `templates`, `walls`
- Requires `canvas:read` (GET) or `canvas:write` (POST/PUT/DELETE) scope

#### Scene Image Endpoints
- `GET /scene/image` — captures the rendered scene canvas as PNG/JPEG; supports viewport capture, custom dimensions, grid display, overlay hiding
- `GET /scene/image/raw` — returns the raw background texture without lighting, tokens, or other rendered layers

#### D&D 5e System Endpoints
New 5e-specific endpoints:
- `GET /dnd5e/concentration` — check whether an actor is currently concentrating
- `POST /dnd5e/break-concentration` — remove an actor's concentration effect
- `POST /dnd5e/concentration-save` — roll a Constitution save to maintain concentration after taking damage
- `POST /dnd5e/equip-item` — toggle an item's equipped state
- `POST /dnd5e/attune-item` — toggle attunement on a magic item
- `POST /dnd5e/transfer-currency` — move currency from one actor to another
- `POST /dnd5e/modify-item-charges` — adjust a consumable item's remaining uses
- `POST /dnd5e/modify-experience` — grant or remove XP from an actor
- `POST /dnd5e/use-spell`, `POST /dnd5e/use-ability`, `POST /dnd5e/use-feature`, `POST /dnd5e/use-item` — activate items and abilities

#### Entity & Search Improvements
- `minified=true` query parameter — returns only `uuid`, `id`, `name`, `img`, `documentType` (smaller payloads)
- `excludeCompendiums=true` query parameter — excludes compendium entries from search results
- Macro and execute-js script validation blocks forbidden patterns before forwarding to Foundry

#### Client Active Check
- `GET /api/clients/{clientId}/active` — returns `{"active": true/false}`
- Public endpoint; used by the Foundry module to detect whether a world is currently connected

#### Email Verification on Registration
- When SMTP is configured, new accounts require email verification before their key is active
- Verification link expires after 24 hours
- Accounts on servers without SMTP are automatically verified immediately

#### Interactive Session Improvements
- Headless manager expanded: auto-start on remote request, credential injection, multi-instance coordination via Redis, configurable inactivity timeout

#### WebSocket Client API Expansion
- Connection tokens authenticate via first WebSocket message (avoids token exposure in URL)
- New WS actions for connection tokens, known clients, remote requests, and module notifications

#### Body Size Limits
- Configurable maximum request body size enforced on all endpoints (default 250 MB)

---

### Changed

#### Authentication
- Password hashing strengthened (bcrypt cost raised from 10 → 12)
- Email format validated at registration
- Session-based auth (cookie) added for the dashboard alongside API key auth

#### Scoped Key Form
- Complete rebuild with full scope checklist
- Danger warnings for high-privilege scopes (`execute-js`, `macro:execute`, `macro:write`)
- World/client selector for per-key access restriction

#### Dashboard
- Updated to show scoped key summary, usage stats, and connection status

#### Docker Compose
- `docker-compose.yml` renamed to `docker-compose.local.yml`
- `docker-compose.postgres.yml` updated

#### Dockerfile
- `go-relay/Dockerfile` updated for image size reduction and build optimizations

#### Fly.io Config
- `fly.toml` and `fly.staging.toml` updated

#### Documentation
- All API doc pages rewritten to match current endpoint behavior
- New pages: `events.md`, `playlist.md`, `user.md`, `cross-world-modules.md`, `security-model.md`
- `configuration.md` updated with all new environment variables
- Three new interactive docs components: `SseTester`, `WsMessageTester`, `WsTester`
- `ApiTester` component updated

---

### Removed
- Old static frontend (`public/index.html`, CSS, JS) replaced by Astro/Svelte build
- Legacy static pages (`privacy.html`, `sheet-test.html`, etc.) consolidated into Astro pages

---

### Environment Variables Added
| Variable | Description | Default |
|---|---|---|
| `ADMIN_JWT_SECRET` | Signing secret for admin sessions (required) | — |
| `ADMIN_SECURE_COOKIES` | Enable `Secure` + `__Host-` prefix for admin cookie | `false` |
| `ADMIN_IP_ALLOWLIST` | Comma-separated IP/CIDR allowlist for `/admin` (empty = open) | — |
| `HEADLESS_SESSION_TIMEOUT` | Session inactivity timeout in seconds (`0` = never) | `600` |
| `MAX_INTERACTIVE_SESSIONS_PER_KEY` | Per-key limit on interactive sheet sessions | `3` |
| `REMOTE_REQUEST_BATCH_WINDOW_SECS` | Notification batch window for cross-world activity | — |
| `LOG_CROSS_WORLD_REQUESTS` | Toggle remote request audit logging per user | — |
