---
id: first-api-call
title: Your First API Call
sidebar_position: 6
---

# Your First API Call

Once you have the relay server running and the Foundry VTT module installed and configured, you're ready to make your first API call.

This guide will walk you through retrieving a list of connected Foundry VTT worlds (clients).

## Prerequisites

1.  **Relay Server is running.** (See [Installation](./installation))
2.  **Foundry VTT is running** with a world loaded.
3.  **The `foundryvtt-rest-api` module is installed and enabled** in your Foundry world.
4.  The module is **configured with the correct relay server URL and your API key**. (See [Foundry Module Setup](./foundry-module))

When the module connects successfully, you should see a log message in the Foundry VTT console (F12).

## Finding Connected Clients

The `/clients` endpoint returns a list of all Foundry VTT instances currently connected to the relay server.

### Request

This is a `GET` request and does not require any parameters in the body or query string. You only need to provide your API key in the header.

**Endpoint:** `GET /clients`

**Header:**
- `x-api-key`: Your API Key

### Example using `curl`

Replace `YOUR_API_KEY_HERE` with the key you obtained from the relay server's web interface.

```bash
curl -X GET http://localhost:3010/clients \
  -H "x-api-key: YOUR_API_KEY_HERE"
```

### Expected Response

If successful, you will receive a JSON response with a `clients` array. Each object in the array represents a connected Foundry VTT world with detailed information about the world and system.

```json
{
  "total": 2,
  "clients": [
      {
          "id": "foundry-vXdyKYLgpmko3kHx",
          "instanceId": "default",
          "lastSeen": 1765376805728,
          "connectedSince": 1765376805728,
          "worldId": "rest-api",
          "worldTitle": "rest-api",
          "foundryVersion": "12.331",
          "systemId": "dnd5e",
          "systemTitle": "Dungeons & Dragons Fifth Edition",
          "systemVersion": "4.3.8",
          "customName": "v12-test"
      },
      {
          "id": "foundry-5ypAoBvOiyjDKiaZ",
          "instanceId": "default",
          "lastSeen": 1765376807068,
          "connectedSince": 1765376807068,
          "worldId": "testing",
          "worldTitle": "testing",
          "foundryVersion": "13.348",
          "systemId": "dnd5e",
          "systemTitle": "Dungeons & Dragons Fifth Edition",
          "systemVersion": "5.0.4",
          "customName": ""
      }
  ]
}
```

**Response Fields:**

| Field | Description |
|-------|-------------|
| `total` | Total number of connected clients |
| `id` | Unique client identifier used in API calls |
| `instanceId` | Server instance handling this connection |
| `lastSeen` | Timestamp of last activity |
| `connectedSince` | Timestamp when client connected |
| `worldId` | Foundry world ID |
| `worldTitle` | Human-readable world name |
| `foundryVersion` | Foundry VTT version (e.g., "13.348") |
| `systemId` | Game system ID (e.g., "dnd5e", "pf2e") |
| `systemTitle` | Full game system name |
| `systemVersion` | Game system version |
| `customName` | Optional custom name set in module settings |
```

### Using the `clientId`

The `clientId` is crucial. You will need to include it as a parameter in most other API calls to specify which connected world you want to interact with.

For example, to get a list of actors from this world, you would make a request to `/structure` and include `clientId=aBcDeFgHiJkLmNoP` and `types=Actor` as query parameters.

Congratulations! You've made your first successful API call. You can now explore the other endpoints in the [API Reference](/api).
