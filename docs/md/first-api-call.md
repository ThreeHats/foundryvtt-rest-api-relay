---
id: first-api-call
title: Your First API Call
sidebar_position: 4
---

# Your First API Call

Once you have the relay server running and the Foundry VTT module installed and configured, you're ready to make your first API call.

This guide will walk you through retrieving a list of connected Foundry VTT worlds (clients).

## Prerequisites

1. **Relay Server is running.** (See [Installation](./installation))
2. **Foundry VTT is running** with a world loaded.
3. **The `foundryvtt-rest-api` module is installed, enabled, and paired with the relay** via the Connection menu — see [Foundry Module Setup](./foundry-module). Pairing is per-browser; the connection token lives in localStorage and is never broadcast to other Foundry users.
4. **You have a scoped API key** for HTTP calls. Create one in the dashboard's **API Keys** page with `clients:read` scope (and any other scopes you'll use). **Do not use your master API key for routine API calls** — that's the recovery credential, not an everyday-use key.

When the module connects successfully, you'll see a log message in the Foundry VTT console (F12) and the world appears in the dashboard's Connected Clients list.

## Finding Connected Clients

The `/clients` endpoint returns a list of all Foundry VTT instances currently connected to the relay server.

### Request

This is a `GET` request and does not require any parameters in the body or query string. You only need to provide your scoped API key in the header.

**Endpoint:** `GET /clients`

**Header:**
- `x-api-key`: Your scoped API key (with `clients:read` scope)

### Example using `curl`

Replace `YOUR_SCOPED_KEY_HERE` with the value from the dashboard's API Keys page.

```bash
curl -X GET http://localhost:3010/clients \
  -H "x-api-key: YOUR_SCOPED_KEY_HERE"
```

### Expected Response

If successful, you will receive a JSON response with a `clients` array. Each object in the array represents a connected Foundry VTT world with detailed information about the world and system.

```json
{
  "total": 2,
  "clients": [
      {
          "clientId": "fvtt_3a9f1c2e4b7d8e0f",
          "instanceId": "01JQKX4Y2N",
          "lastSeen": 1765376805728,
          "connectedSince": 1765376805728,
          "worldId": "rest-api",
          "worldTitle": "rest-api",
          "foundryVersion": "12.331",
          "systemId": "dnd5e",
          "systemTitle": "Dungeons & Dragons Fifth Edition",
          "systemVersion": "4.3.8",
          "customName": "v12-test",
          "isOnline": true
      },
      {
          "clientId": "fvtt_b2e7a05c9d1f3e4a",
          "worldId": "testing",
          "worldTitle": "testing",
          "foundryVersion": "13.348",
          "systemId": "dnd5e",
          "systemTitle": "Dungeons & Dragons Fifth Edition",
          "systemVersion": "5.0.4",
          "customName": "",
          "isOnline": false
      }
  ]
}
```

**Response Fields:**

| Field | Description |
|-------|-------------|
| `total` | Total number of clients (online and offline) |
| `clientId` | Unique client identifier used in API calls |
| `isOnline` | Whether the client is currently connected |
| `instanceId` | Relay instance handling the connection (omitted when offline) |
| `lastSeen` | Unix timestamp (ms) of last activity (omitted when offline) |
| `connectedSince` | Unix timestamp (ms) when client connected (omitted when offline) |
| `worldId` | Foundry world ID |
| `worldTitle` | Human-readable world name |
| `foundryVersion` | Foundry VTT version (e.g., "13.348") |
| `systemId` | Game system ID (e.g., "dnd5e", "pf2e") |
| `systemTitle` | Full game system name |
| `systemVersion` | Game system version |
| `customName` | Optional custom name set in module settings |

### Using the `clientId`

The `clientId` identifies which connected Foundry world to interact with. Include it as a query parameter in API calls.

**`clientId` is optional** in two cases — you can omit it and the API auto-resolves it:

1. **Only one client is online** — auto-selected regardless of key type.
2. **Your scoped key is bound to specific clients** — the key's allowed client list is used. If only one of the allowed clients is online, it auto-selects. If exactly one client is configured on the key, every request goes to that client automatically with no `clientId` needed.

When multiple clients are connected and the key is not client-scoped, you must specify which one to use.

```bash
# Explicit clientId
curl -X GET "http://localhost:3010/structure?clientId=fvtt_3a9f1c2e4b7d8e0f&types=Actor" \
  -H "x-api-key: YOUR_API_KEY_HERE"

# Auto-resolved (one client online, or key is bound to one client)
curl -X GET "http://localhost:3010/structure?types=Actor" \
  -H "x-api-key: YOUR_API_KEY_HERE"
```

If you omit `clientId` with multiple unfiltered clients connected, the API returns a 400 error listing the connected clients so you can choose one.

Congratulations! You've made your first successful API call. From here:

- Explore the other endpoints in the [API Reference](/api).
- Learn about [Scoped API Keys](./scoped-keys) for production-grade integrations with narrow scopes and per-key rate limits.
- If you're building a Foundry module that needs to talk to OTHER Foundry worlds (not just call the relay's HTTP API), see [Building Cross-World Foundry Modules](./cross-world-modules) — modules NEVER hold HTTP API keys; they use the WS tunnel pattern instead.
- For the full credential model and threat reasoning, read [Authentication & Security Model](./authentication).
