---
id: configuration
title: Server Configuration
sidebar_position: 5
---

# Server Configuration

The FoundryVTT REST API Relay server can be configured using environment variables. This allows you to customize its behavior for different environments (development, production) and database setups.

## Environment Variables

Here are the primary environment variables you can use:

-   `PORT`: The port the server will listen on.
    -   **Default:** `3010`

-   `NODE_ENV`: The environment mode.
    -   **Default:** `development`
    -   **Set to `production` for live deployments.**

-   `DB_TYPE`: The type of database to use.
    -   **Options:** `sqlite`, `postgres`, `memory`
    -   **Default (in `docker-compose.yml`):** `sqlite`
    -   **Default (manual `pnpm local`):** `memory`

-   `DATABASE_URL`: The connection string for the database.
    -   **Required for `postgres` `DB_TYPE`.**
    -   **Example:** `postgres://user:password@host:port/database`

-   `WEBSOCKET_PING_INTERVAL_MS`: Interval in milliseconds for sending WebSocket protocol pings to keep connections alive.
    -   **Default:** `20000` (20 seconds)

-   `CLIENT_CLEANUP_INTERVAL_MS`: Interval in milliseconds for the server to check for and remove inactive or disconnected clients.
    -   **Default:** `15000` (15 seconds)

-   `REDIS_URL`: Connection URL for a Redis instance.
    -   **Optional.** Used for session storage and can help in multi-instance deployments.

-   `FREE_API_REQUESTS_LIMIT`: Number of free requests per month
    -   **Default:** 100

-   `DAILY_REQUEST_LIMIT`: Number of requests per day allowed per API key
    -   **Default** 1000

-   `CAPTURE_BROWSER_CONSOLE`: Log level for capturing browser console output from headless Foundry sessions to files in `data/browser-logs/`.
    -   **Options:** `error`, `warn`, `debug`
    -   Each level includes all levels above it (e.g., `warn` captures both warnings and errors).
    -   **Default:** unset (no browser console logging)

## Docker Configuration

When using Docker, you can set these variables in the `environment` section of your `docker-compose.yml` file.

**Example `docker-compose.yml`:**
```yaml
services:
  relay:
    image: threehats/foundryvtt-rest-api-relay:latest
    container_name: foundryvtt-rest-api-relay
    ports:
      - "3010:3010"
    environment:
      - NODE_ENV=production
      - PORT=3010
      - DB_TYPE=sqlite
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

## Manual Start Configuration

When running the server manually, you can set these variables in your shell or use a `.env` file with a tool like `dotenv-cli`. The provided `package.json` scripts already handle some of this for you.

**Example `package.json` script:**
```json
"scripts": {
  "local:sqlite": "cross-env DB_TYPE=sqlite PORT=3010 tsx watch src/index.ts"
}
```
This script sets `DB_TYPE` to `sqlite` and `PORT` to `3010` when you run `pnpm local:sqlite`.
