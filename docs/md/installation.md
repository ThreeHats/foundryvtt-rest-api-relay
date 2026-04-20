---
id: installation
title: Installation
sidebar_position: 2
---

# Installation

There are two primary ways to install the FoundryVTT REST API Relay server: using Docker (recommended for ease of use and deployment) or manually.

## Recommended: Docker Installation

Using Docker and Docker Compose is the simplest way to get the relay server running.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/JustAnotherIdea/foundryvtt-rest-api-relay.git
    cd foundryvtt-rest-api-relay
    ```

Instead of cloning the repository, you can also download the `docker-compose.local.yml` file from the repo and use that directly.

Or create your own compose file using the image:
```yaml
services:
  relay:
    # For production stability, pin to a specific version tag instead of 'latest'
    # See available tags at: https://github.com/ThreeHats/foundryvtt-rest-api-relay/tags
    image: threehats/foundryvtt-rest-api-relay:latest
    container_name: foundryvtt-rest-api-relay
    ports:
      - "3010:3010"
    environment:
      - APP_ENV=production
      - PORT=3010
      - DB_TYPE=sqlite
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

:::tip Version Pinning for Production
For production deployments, replace `latest` with a specific version tag (e.g., `threehats/foundryvtt-rest-api-relay:3.0.0`) to avoid unexpected breaking changes from updates.
:::

:::info Request limits
The Docker image has **no request limits by default** (`MONTHLY_REQUEST_LIMIT=0` — unlimited). If you want to enforce a monthly quota on your self-hosted users, add it to your compose file's `environment` section:
```yaml
- MONTHLY_REQUEST_LIMIT=5000
```
Paid subscribers (Stripe `active` status) are always unlimited regardless of this value.
:::

:::info Billing / subscription UI
Subscription-related UI (plan badges, upgrade buttons) is automatically **hidden** when `STRIPE_SECRET_KEY` is not set. Self-hosted deployments show a clean dashboard with no billing elements.
:::

:::info Headless sessions
Headless browser sessions (automated GM login via Chromium) are **enabled by default** on self-hosted instances (`ALLOW_HEADLESS=true`). Key tuning variables:
- `HEADLESS_SESSION_TIMEOUT` — inactivity timeout in seconds before a session is stopped (default: `600`). Set to `0` to never time out.
- `MAX_HEADLESS_SESSIONS` — max concurrent headless sessions (default: `1`).
- `PUPPETEER_EXECUTABLE_PATH` — path to Chrome/Chromium if not auto-detected.

See [Server Configuration](./configuration) for the full list of variables.
:::

2.  **Start the server:**
    ```bash
    docker compose -f docker-compose.local.yml up -d
    ```
    This will pull the latest Docker image and start the relay server in the background. The server will be available at `http://localhost:3010`.

3.  **Create Your Account:**
    The default Docker setup uses an SQLite database for persistence, stored in the `data` directory.
    - Open `http://localhost:3010` in your browser.
    - Click **Sign Up** and create an account.
    - **Your master API key will be displayed exactly once** in a one-time modal after registration. Copy it and save it in a password manager. Never share this key with anyone, and do not create applications using this key.
    - The dashboard never displays the master key again. For routine HTTP API calls, create a **scoped API key** with narrow scopes via the API Keys page. See [Authentication](./authentication) for the full credential model.

4.  **Stopping the server:**
    ```bash
    docker compose -f docker-compose.local.yml down
    ```

4.  **Updating the server:**
-   **[Updating the docker image](./update-docker-image):** Commands to update your docker image.

### Using PostgreSQL
If you prefer to use PostgreSQL for your database, you can use the provided `docker-compose.postgres.yml` file. See the [PostgreSQL Setup Guide](/postgres-setup) for more details.

### Relay + Foundry + duckDNS
For an in depth guide for a full setup using duckDNS see [Relay + App + DNS Example](/relay-app-duckdns-example)

## Manual Installation

If you prefer not to use Docker, you can build and run the Go server directly.

1.  **Prerequisites:**
    - Go 1.22 or later
    - Node.js v18+ and pnpm (only needed for frontend build and tests)
    - Chromium/Chrome (only needed for headless session features)

2.  **Clone the repository:**
    ```bash
    git clone https://github.com/ThreeHats/foundryvtt-rest-api-relay.git
    cd foundryvtt-rest-api-relay
    ```

3.  **Build and run the Go server:**
    ```bash
    # Start relay with SQLite (persists to data/relay.db)
    pnpm run local:sqlite

    # Or build the binary directly
    cd go-relay
    go build -o relay ./cmd/server/
    DB_TYPE=sqlite PORT=3010 ./relay
    ```

4.  **Build the frontend (optional):**
    ```bash
    pnpm run frontend:build
    ```

The server will be running at `http://localhost:3010`.

## Accessing the Relay From Other Devices on Your LAN

The relay binds to all network interfaces (`0.0.0.0`) by default, so a self-hosted instance is automatically reachable from other devices on the same local network - phones, tablets, another PC, your Foundry host, etc. On startup, the server logs every LAN URL it can be reached at, e.g.:

```
INF Server listening on all interfaces (0.0.0.0) port=3010
INF Local URL url=http://localhost:3010
INF LAN URL (reachable from other devices on your network) url=http://192.168.1.42:3010
```

To use it from another device, just hit `http://<that-ip>:3010` instead of `localhost`.

A few things to check if it doesn't work:

- **Host firewall** — make sure inbound TCP on port `3010` is allowed (ufw, firewalld, Windows Defender Firewall, etc.).
- **Docker** — the example `docker-compose.local.yml` already publishes `3010:3010`, which is all you need. If you're running with `docker run`, include `-p 3010:3010`.
- **Foundry module** — when configuring the FoundryVTT REST API module, point its relay URL at the LAN IP (`http://192.168.1.42:3010`), not `localhost`, unless Foundry is running on the exact same machine as the relay.
- **CORS** — already permissive (`*`), so browser-based clients on other LAN devices work without extra config.

If you only want the relay reachable on the local machine (e.g., you're exposing it via a reverse proxy and don't want it directly bound to the LAN), bind Docker's port mapping to loopback only: `127.0.0.1:3010:3010`.
