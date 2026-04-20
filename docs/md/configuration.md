---
id: configuration
title: Server Configuration
sidebar_position: 10
---

# Server Configuration

All configuration is done via environment variables. Copy `.env.example` to `.env` and adjust for your setup.

---

## Core

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3010` | Port the server listens on |
| `APP_ENV` | `development` | Environment mode. Set to `production` for live deployments — enables SSL for PostgreSQL and tightens security defaults |
| `DB_TYPE` | `sqlite` | Database backend: `sqlite` (default, self-hosted) or `postgres` (production/scaled) |
| `DATABASE_URL` | — | PostgreSQL connection string. Required when `DB_TYPE=postgres`. Example: `postgresql://user:pass@host:5432/db` |
| `LOG_LEVEL` | `info` | Zerolog log level: `debug`, `info`, `warn`, `error` |
| `DATA_DIR` | `data` | Directory for persistent data (SQLite DB, browser logs). Resolved to an absolute path at startup |
| `SQLITE_PATH` | `<DATA_DIR>/relay.db` | Override the SQLite file path. Only used when `DB_TYPE=sqlite` |

---

## Auth & Registration

| Variable | Default | Description |
|----------|---------|-------------|
| `ADMIN_EMAIL` | — | If set, creates this admin account on startup if it doesn't exist |
| `ADMIN_PASSWORD` | — | Password for the auto-created admin account |
| `DISABLE_REGISTRATION` | `false` | When set, prevents new user self-registration |
| `CREDENTIALS_ENCRYPTION_KEY` | — | AES-256 key for encrypting Foundry credentials stored on scoped API keys. **Required for all backends.** Generate: `openssl rand -hex 32` |
| `RETURN_RESET_TOKEN` | `false` | **Testing only.** Returns the raw password reset token in the API response instead of emailing it. Never enable in production |

---

## Admin Dashboard

| Variable | Default | Description |
|----------|---------|-------------|
| `ADMIN_JWT_SECRET` | — | JWT signing secret for admin sessions. **Required for all deployments** (sqlite and postgres). Generate: `openssl rand -base64 32` |
| `ADMIN_ALLOWED_IPS` | — (allow all) | Comma-separated IP allowlist for all `/admin` routes. Example: `127.0.0.1,10.0.0.5` |
| `ADMIN_INTERNAL_PORT` | `0` (disabled) | If set, serves the admin panel on a separate internal port (useful for Fly.io private networking) |
| `ADMIN_SESSION_MAX_HOURS` | `8` | Maximum admin session duration in hours. In dev mode the full window is used as the JWT TTL — no re-login needed within this window |
| `ADMIN_AUDIT_LOG_RETENTION_DAYS` | `90` | How many days to keep admin audit log entries |
| `ADMIN_LOGIN_RATE_LIMIT` | `10` | Max admin login attempts per IP per 15-minute window |
| `ADMIN_LOCKOUT_THRESHOLD` | `5` | Failed login attempts before account lockout |
| `ADMIN_LOCKOUT_MINUTES` | `15` | How long (minutes) a locked account stays locked |
| `ADMIN_SECURE_COOKIES` | `false` | Set to `true` to enable the `Secure` cookie flag and `__Host-` prefix for the admin session cookie. Only works when the admin panel is served over HTTPS. Default `false` so plain HTTP self-hosted deployments work out of the box |

---

## API Usage Limits

:::info Self-hosted deployments
The Docker image defaults the monthly limit to **`0` (unlimited)**. You only need to set this if you want to enforce a quota on your own instance. Subscription UI (plan badge, upgrade button) is automatically hidden when `STRIPE_SECRET_KEY` is not set — no extra config required.
:::

| Variable | Default (Docker) | Description |
|----------|---------|-------------|
| `MONTHLY_REQUEST_LIMIT` | `0` (unlimited) | Monthly request quota for free-tier users. `0` = unlimited. Paid subscribers (Stripe active) are always unlimited. The commercial hosted relay uses `5000` |

---

## Rate Limiting

All rate limiters are per-IP unless noted.

| Variable | Default | Window | Applies to |
|----------|---------|--------|------------|
| `PER_MINUTE_REQUEST_LIMIT` | `300` | 1 min | Per-API-key burst limit on all API endpoints. Set to `0` to disable |
| `AUTH_RATE_LIMIT` | `5` | 15 min | Register, login, regenerate-key endpoints |
| `PASSWORD_RESET_RATE_LIMIT` | `3` | 1 hour | Forgot-password / reset-password endpoints |
| `ACCOUNT_MGMT_RATE_LIMIT` | `10` | 1 hour | Change-password, delete-account, export-data endpoints |
| `PAIRING_RATE_LIMIT` | `10` | 15 min | Relay-to-Foundry pairing attempts |
| `KEY_REQUEST_RATE_LIMIT` | `20` | 15 min | Scoped key approval flow submissions |
| `PROBE_RATE_LIMIT` | `60` | 1 min | Active-connection probe requests |

---

## Redis

Redis is optional. When enabled it allows multiple relay instances to route requests to the correct Foundry client across instances. Automatically disabled when `DB_TYPE=sqlite`.

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_URL` | — | Redis connection URL. Example: `redis://default:password@host:6379` |
| `REDIS_ENABLED` | `true` | Set to `false` to disable Redis even when `REDIS_URL` is set |

---

## Email (Password Reset & Notifications)

If `SMTP_HOST` is not configured, password reset tokens are logged to the console instead of emailed.

| Variable | Default | Description |
|----------|---------|-------------|
| `SMTP_HOST` | — | SMTP server hostname |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_USER` | — | SMTP username |
| `SMTP_PASS` | — | SMTP password |
| `SMTP_FROM` | `noreply@foundryvtt-relay.com` | From address for outgoing emails |
| `SMTP_SECURE` | `false` | Use TLS for SMTP connection |

---

## Stripe (Subscription Billing)

Stripe is fully disabled when `STRIPE_SECRET_KEY` is not set.

| Variable | Default | Description |
|----------|---------|-------------|
| `STRIPE_SECRET_KEY` | — | Stripe API secret key (`sk_test_...` for dev, `sk_live_...` for production) |
| `STRIPE_WEBHOOK_SECRET` | — | Stripe webhook signing secret (`whsec_...`) |
| `STRIPE_PRICE_ID` | — | Stripe Price ID for the subscription product |
| `STRIPE_PORTAL_URL` | — | Stripe Customer Portal URL for self-service subscription management |
| `FRONTEND_URL` | `https://foundryrestapi.com` | Base URL for Stripe checkout redirect URLs |

---

## WebSocket & Connection Tuning

| Variable | Default | Description |
|----------|---------|-------------|
| `WEBSOCKET_PING_INTERVAL_MS` | `20000` | Interval (ms) for WebSocket keepalive pings |
| `CLIENT_CLEANUP_INTERVAL_MS` | `15000` | Interval (ms) to check for and remove dead clients |
| `AUTH_CACHE_TTL_SECONDS` | `30` | How long to cache API key lookups in memory to reduce DB load |

---

## Headless Sessions

Headless sessions use Chromium to automate Foundry logins.

| Variable | Default | Description |
|----------|---------|-------------|
| `ALLOW_HEADLESS` | `true` | Set to `false` to completely disable headless session creation. **Note:** headless sessions are disabled on the public hosted relay — self-host to use this feature |
| `MAX_HEADLESS_SESSIONS` | `1` | Max concurrent headless browser sessions per instance. When Redis is configured, this cap is enforced globally across all relay instances |
| `HEADLESS_SESSION_TIMEOUT` | `600` | Seconds of inactivity before a headless session is automatically stopped. Set to `0` to never time out due to inactivity |
| `MAX_INTERACTIVE_SESSIONS_PER_KEY` | `3` | Max interactive sheet sessions per scoped API key |
| `PUPPETEER_EXECUTABLE_PATH` | auto-detected | Path to Chrome/Chromium executable |
| `CAPTURE_BROWSER_CONSOLE` | — (disabled) | Log level for browser console output: `error`, `warn`, or `debug`. Logs saved to `<DATA_DIR>/browser-logs/` |
| `REQUIRE_PAID_HEADLESS` | `false` | Require an active paid subscription to create headless sessions |
| `CHROME_USER_DATA_DIR` | `<DATA_DIR>/chrome-profile` | Persistent Chrome profile directory. Enables V8 bytecode cache and HTTP cache across restarts — significantly reduces warm-session JS execution time |
| `CHROME_JS_HEAP_MB` | `2048` | Max V8 old-space heap in MB. Increase if running heavy module loadouts to reduce GC pause spikes |
| `CHROME_WINDOW_WIDTH` | `1280` | Headless browser viewport width in pixels |
| `CHROME_WINDOW_HEIGHT` | `800` | Headless browser viewport height in pixels. Must be ≥768 (Foundry minimum). Also affects canvas/entity screenshot resolution |
| `CHROME_ENABLE_SHM` | `false` | Allow Chrome to use `/dev/shm` for IPC. Enable when the host/container has ≥256 MB of shared memory (`--shm-size=256m` in Docker) |
| `CHROME_GPU_MODE` | `auto` | Chrome rendering backend. `auto` detects the best available option. `gpu` = Ozone headless + ANGLE GL (real GPU via Mesa/DRI on Linux — no display server, no windows). `xvfb` = Mesa software rendering via Xvfb. `swiftshader` = software WebGL bundled in Chrome (always works, used as fallback). |

---

## Fly.io / Infrastructure

These are typically set automatically by Fly.io and do not need to be configured manually.

| Variable | Default | Description |
|----------|---------|-------------|
| `FLY_ALLOC_ID` | `local` | Instance identifier for multi-instance request routing |
| `FLY_INTERNAL_PORT` | — | Internal port on the Fly.io private network |
| `APP_NAME` | — | App name for Fly.io private networking |

---

## Request Body Limits

| Variable | Default | Description |
|----------|---------|-------------|
| `MAX_JSON_BODY_SIZE_MB` | `250` | Maximum request body size in MB for API calls. Applies globally |
| `MAX_UPLOAD_SIZE_MB` | `250` | Maximum file upload size in MB. Applies to file upload endpoints |

---

## Deprecated / Legacy Names

These names still work as fallbacks but should not be used in new deployments.

| Old name | Current name |
|----------|-------------|
| `FREE_API_REQUESTS_LIMIT` | `MONTHLY_REQUEST_LIMIT` |
| `NODE_ENV` | `APP_ENV` |
| `HEADLESS_INACTIVE_TIMEOUT` | `HEADLESS_SESSION_TIMEOUT` |

---

## Documentation Site

These variables are only needed when building or deploying the Docusaurus documentation site (`docs/`). They are not used by the relay server.

When both `ALGOLIA_APP_ID` and `ALGOLIA_SEARCH_API_KEY` are set, the docs site enables Algolia DocSearch. If either is absent, the search widget is omitted from the build.

| Variable | Default | Description |
|----------|---------|-------------|
| `ALGOLIA_APP_ID` | — | Algolia application ID. Required to enable doc search |
| `ALGOLIA_SEARCH_API_KEY` | — | Algolia search-only API key (public, safe to expose in the browser) |
| `ALGOLIA_INDEX_NAME` | `foundryvtt-rest-api-relay` | Algolia index name to query |

---

## Docker Configuration

Two compose files are included in the repository:

- **`docker-compose.local.yml`** — SQLite, no billing, unlimited requests. The recommended starting point for self-hosted setups.
- **`docker-compose.postgres.yml`** — Postgres + optional Redis. For self-hosted setups that prefer Postgres or need horizontal scaling.

```yaml
# docker-compose.local.yml (self-hosted, SQLite)
# Copy .env.example to .env and fill in ADMIN_EMAIL, ADMIN_PASSWORD, etc.
services:
  relay:
    image: threehats/foundryvtt-rest-api-relay:latest
    container_name: foundryvtt-rest-api-relay
    env_file: .env
    environment:
      - PORT=3010
      - DB_TYPE=sqlite
      - APP_ENV=production
    volumes:
      - ./data:/app/data
    ports:
      - "3010:3010"
    restart: unless-stopped
```

Start it with:

```bash
docker compose -f docker-compose.local.yml up -d
```
