# Foundry REST API

**A REST API bridge for Foundry Virtual Tabletop** — enabling external tools, automations, and integrations to interact with your Foundry VTT worlds.

[![Discord](https://img.shields.io/discord/1234567890?label=Discord&logo=discord&logoColor=white)](https://discord.gg/7SdGxJmKkE)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## What Is This?

Foundry VTT is a self-hosted virtual tabletop platform, but it lacks a native API for external integrations. This project fills that gap by providing a **WebSocket relay server** that exposes your Foundry world data through a clean REST API.

**Use cases:**
- Build custom dashboards or companion apps for your game
- Automate game master tasks (e.g., trigger macros, roll dice remotely)
- Integrate with MIDI controllers, stream decks, or voice assistants
- Create Discord bots that interact with your world data
- Run automated integration tests
- Develop tools that read/write actors, items, scenes, and more

## How It Works

```
┌─────────────────┐      WebSocket       ┌─────────────────┐      REST API      ┌─────────────────┐
│   Foundry VTT   │ ◄──────────────────► │  Relay Server   │ ◄────────────────► │  Your App/Tool  │
│   + Module      │                      │  (this repo)    │                    │                 │
└─────────────────┘                      └─────────────────┘                    └─────────────────┘
```

1. **Foundry Module** connects to the relay server via WebSocket
2. **Relay Server** exposes REST endpoints and routes requests to your Foundry world
3. **Your application** calls the REST API to search, read, create, and modify entities

## Project Components

| Component | Description |
|-----------|-------------|
| [**Relay Server**](https://github.com/ThreeHats/foundryvtt-rest-api-relay) (this repo) | Node.js/Express server with WebSocket relay and REST API |
| [**Foundry Module**](https://github.com/ThreeHats/foundryvtt-rest-api) | Foundry VTT module that connects your world to the relay |

---

## Quick Start

You can either use the **public relay server** (no setup required) or **self-host** your own instance.

### Option A: Use the Public Relay (Easiest)

1. Go to **[https://foundryvtt-rest-api-relay.fly.dev](https://foundryvtt-rest-api-relay.fly.dev)**
2. Create an account and copy your API key from the dashboard
3. Install the [Foundry module](#3-install-the-foundry-module) and set the relay URL to `wss://foundryvtt-rest-api-relay.fly.dev`

> The public relay has rate limits (100 requests/month free, 1000/day). For unlimited usage, self-host or [subscribe](https://foundryvtt-rest-api-relay.fly.dev) for $5/month.

### Option B: Self-Host with Docker

```bash
# Download and run with Docker Compose
curl -O https://raw.githubusercontent.com/ThreeHats/foundryvtt-rest-api-relay/main/docker-compose.yml
docker-compose up -d
```

The server runs at `http://localhost:3010`.

### 2. Get Your API Key

1. Open your relay server (`http://localhost:3010` or the [public URL](https://foundryvtt-rest-api-relay.fly.dev)) in your browser
2. **Create an account** using the Sign Up form
3. After logging in, your API key is displayed on the dashboard — click **Copy** to copy it

### 3. Install the Foundry Module

Install via manifest URL in Foundry VTT:
```
https://github.com/ThreeHats/foundryvtt-rest-api/releases/latest/download/module.json
```

Enable the module in your world and configure the relay URL in the module settings:
- **Public relay:** `wss://foundryvtt-rest-api-relay.fly.dev`
- **Self-hosted:** `ws://localhost:3010` (or `wss://` if using HTTPS)

### 4. Make Your First API Call

```bash
curl -X GET "http://localhost:3010/clients" \
  -H "x-api-key: YOUR_API_KEY"
```

You'll see your connected Foundry worlds listed in the response.

---

## API Overview

The REST API provides endpoints across these resource groups:

| Endpoint Group | Description |
|---------------|-------------|
| `/clients` | List connected Foundry worlds |
| `/get`, `/create`, `/update`, `/delete` | CRUD operations for entities (actors, items, scenes, etc.) |
| `/search` | Search your Foundry database |
| `/roll` | Perform dice rolls |
| `/macro` | Execute macros remotely |
| `/encounter` | Combat encounter management |
| `/dnd5e` | D&D 5th Edition specific operations |
| `/structure` | World structure information |
| `/session` | Headless Foundry sessions via Puppeteer |
| `and more` | Read the [api reference](https://foundryvtt-rest-api-relay.fly.dev/docs/api) for more details |

**Authentication:** Include your API key in the `x-api-key` header for all requests.

📖 **[Full API Documentation](https://foundryvtt-rest-api-relay.fly.dev/docs)**

> *The Postman Collection is deprecated — the interactive API documentation now includes code examples in multiple languages.*

---

## Installation Options

### Docker Compose (Recommended)

```yaml
services:
  relay:
    # For production, pin to a specific version (e.g., 2.1.0) instead of 'latest'
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
    command: pnpm local:sqlite
    restart: unless-stopped
```

```bash
docker-compose up -d
```

> 💡 **Tip:** For production stability, pin to a [specific release tag](https://github.com/ThreeHats/foundryvtt-rest-api-relay/tags) (e.g., `threehats/foundryvtt-rest-api-relay:2.1.0`) instead of `latest`.

### Manual Installation

**Prerequisites:** Node.js v18+ (v20 recommended), pnpm

```bash
git clone https://github.com/ThreeHats/foundryvtt-rest-api-relay.git
cd foundryvtt-rest-api-relay

pnpm install

# Build SQLite native module (path may vary with sqlite3 version)
cd node_modules/.pnpm/sqlite3@*/node_modules/sqlite3 && npm run install && cd -

# Run with SQLite (recommended for local development)
pnpm run local:sqlite
```

> **Note:** On Windows, the SQLite build step may require additional configuration. See the [full installation guide](https://foundryvtt-rest-api-relay.fly.dev/docs/installation) for platform-specific instructions.

---

## Configuration

Key environment variables (see [full configuration guide](https://foundryvtt-rest-api-relay.fly.dev/docs/configuration) for details):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3010` | Server port |
| `DB_TYPE` | `sqlite` | Database type (`sqlite`, `postgres`, `memory`) |
| `DATABASE_URL` | — | PostgreSQL connection string (required for `postgres` DB_TYPE) |
| `WEBSOCKET_PING_INTERVAL_MS` | `20000` | WebSocket ping interval for keep-alive |
| `CLIENT_CLEANUP_INTERVAL_MS` | `15000` | Interval for removing inactive clients |
| `REDIS_URL` | — | Redis URL for session storage in multi-instance deployments |
| `FREE_API_REQUESTS_LIMIT` | `100` | Monthly API request limit per user |
| `DAILY_REQUEST_LIMIT` | `1000` | Daily API request limit per user |

> **Self-hosting tip:** Set `FREE_API_REQUESTS_LIMIT=999999999` and `DAILY_REQUEST_LIMIT=999999999` to effectively disable rate limits for local use.

---

## Tech Stack

**Relay Server:**
- **Runtime:** Node.js + TypeScript
- **Framework:** Express.js
- **Real-time:** WebSocket (ws)
- **Database:** SQLite (default), PostgreSQL + Redis (production), or in-memory
- **Headless Browser:** Puppeteer (for unattended Foundry sessions)
- **Docs:** Docusaurus + TypeDoc (auto-generated API reference)
- **Testing:** Jest + TypeDoc (partially auto-generated tests)

**Foundry Module:**
- TypeScript module for Foundry VTT
- Integrates with QuickInsert for search capabilities

---

## Documentation and Testing Development

```bash
pnpm docs:install    # Install doc dependencies
pnpm docs:generate   # Generate API docs from TypeScript comments
pnpm test:generate   # Generate bare bones Jest tests
pnpm test            # Run full integration tests while capturing examples
pnpm docs:examples   # Update API docs with test examples
pnpm docs:dev        # Start docs server at localhost:3000
```

---

## Links

- 📖 [Documentation](https://foundryvtt-rest-api-relay.fly.dev/docs)
- 🤝 [Contributing Guide](https://foundryvtt-rest-api-relay.fly.dev/docs/contributing)
- 💬 [Discord Community](https://discord.gg/7SdGxJmKkE)

---

## License

MIT © [ThreeHats](https://github.com/ThreeHats)

