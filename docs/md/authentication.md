---
id: authentication
title: Authentication & Security Model
sidebar_position: 5
---

# Authentication & Security Model

The Foundry REST API uses **four distinct credential types**, each with a different lifetime, audience, and trust boundary. Understanding which credential to use is the most important security decision you'll make.

## TL;DR

| If you are… | Use this credential | Sent as… |
|---|---|---|
| An external app (Discord bot, Obsidian plugin, script) | **Scoped API key** | `x-api-key: <scoped-key>` header |
| The Foundry REST API module | **Connection token** | First WebSocket message after connect |
| Testing the API by hand | **Scoped API key** with the scopes you need | `x-api-key: <scoped-key>` header |

## The three credentials

| Credential | What it auths | Where it lives | How to get one |
|---|---|---|---|
| **Dashboard session** | Dashboard access | Browser (managed automatically) | Email + password login |
| **Connection token** | One Foundry browser's WebSocket relay connection | The Foundry browser's client-scope localStorage (per-device) | 6-char pairing code from dashboard |
| **Scoped API key** | External integrations calling HTTP endpoints | The integration's secrets store | Created in the dashboard |

---

## Account access

Your account is protected by **email + password**. When you log in, the dashboard mints a short-lived session token stored in your browser — you never handle it directly. All dashboard operations run under this session.

There is no "master API key" to copy, store, or protect. Your password is your account credential.

### If your password is compromised

Change it immediately from the dashboard. You can also trigger a **full credential reset** using the **Reset Credentials** button at the bottom of the dashboard, which invalidates all scoped keys, connection tokens, and active sessions at once. Every Foundry browser will need to re-pair.

---

## Connection token

The credential a Foundry module uses to connect to the relay's WebSocket. **This is the only credential a Foundry module ever holds.** It is intentionally WS-only - it cannot authenticate HTTP API calls.

### Storage

- `clientId` and `pairedRelayUrl` live in **world settings** (shared across all GMs). These are non-secret pointers - the clientId is opaque and learning it grants no access.
- `connectionToken` lives in a **client-scope Foundry setting** (browser localStorage, per-device). This is the only Foundry storage that doesn't broadcast to all connected clients, verified against Foundry v13 source at `client/helpers/client-settings.mjs:42-46`.

This split means: each GM pairs each device once. A second GM joining an already-paired world gets prompted to pair their browser using an "Add Browser" code from the dashboard's Known Clients page.

### Why per-browser?

Every other Foundry storage mechanism that claims to hide data from non-GM players (user flags, `scope: "user"` settings, `gmOnly` schema fields, JournalEntry ownership) **broadcasts to every connected client**. Foundry's permission system gates the UI, not the wire data - a player can read any of those values from the browser console. Client-scope settings are the only exception.

### Pairing a browser

1. In the dashboard: **Connection Tokens** → **Generate Pairing Code** (6-char code, expires 10 minutes).
2. In Foundry: **Game Settings** → **REST API Connection** → **Enter Code**.
3. The module stores the token in this browser's localStorage. Reload Foundry and it connects.

Or use **Pair via Browser** - the recommended flow that opens the relay dashboard directly and completes pairing without manually entering a code.

To pair a second GM's browser: the dashboard's **Known Clients** page has an **Add Browser** button per world. The second GM enters that code in their Foundry.

### Cross-world capability (WS tunnel)

Connection tokens have two optional permissions set at pair time:

- `allowedTargetClients` - other clientIds this token may invoke `remote-request` operations against
- `remoteScopes` - scope strings the token holds for those remote operations

A Foundry module with these permissions can call `module.api.remoteRequest(targetClientId, action, payload)` to invoke an action on another Foundry world owned by the same account, without ever holding an HTTP API key.

Default: empty. A GM pairing their browser doesn't accidentally get cross-world powers - they'd need to explicitly enable them when generating the pair code.

See [Building Cross-World Foundry Modules](./cross-world-modules) for the full guide.

---

## Scoped API key

For external integrations (Discord bots, Obsidian plugins, custom scripts) that call the relay's HTTP API. Created with a specific scope set, optional client restriction, optional monthly rate limit, and optional expiry.

Authenticated via `x-api-key: <scoped-key>` on every HTTP request.

### Requesting a key from your integration (preferred)

Integrations can request a scoped API key directly - the user reviews and approves the request in the dashboard, and the key is delivered back to the integration automatically. This is the preferred approach: users don't need to manually create keys and copy them into your integration.

There are two flows depending on whether your integration can receive an HTTP redirect.

#### Device flow (CLI tools, scripts, desktop apps)

Use this when your integration can't receive an inbound HTTP request.

1. **Request a key** - `POST /auth/key-request` with your app info and the scopes you need:

```json
{
  "appName": "My Discord Bot",
  "appDescription": "Rolls dice and looks up characters",
  "appUrl": "https://github.com/you/my-bot",
  "scopes": ["entity:read", "roll:execute"],
  "suggestedMonthlyLimit": 1000
}
```

Response:
```json
{
  "code": "a1b2c3",
  "approvalUrl": "https://foundryrestapi.com/approve/a1b2c3",
  "expiresIn": 600,
  "expiresAt": "2024-01-01T12:10:00Z"
}
```

2. **Direct the user** to the `approvalUrl` (print it, open it in a browser, etc.). They log in to the dashboard, review the requested scopes, and approve.

3. **Poll for the result** - `GET /auth/key-request/:code/status` every few seconds until status changes:

| `status` | Meaning |
|---|---|
| `pending` | User hasn't acted yet - keep polling |
| `approved` | Approved - response includes `apiKey`, `scopes`, `clientIds` |
| `denied` | User denied the request |
| `expired` | 10-minute window passed without action |
| `exchanged` | Key already retrieved (poll returned `approved` once, then became `exchanged`) |

4. **Store the key** from the `approved` response. The key is returned exactly once - the status becomes `exchanged` immediately after.

#### Web flow (web apps that can receive a redirect)

Use this when your integration has a server that can receive an inbound HTTP request.

1. **Request a key** - same `POST /auth/key-request` as above, but include a `callbackUrl`:

```json
{
  "appName": "My Web App",
  "scopes": ["entity:read", "roll:execute"],
  "callbackUrl": "https://myapp.example.com/foundry/callback"
}
```

2. **Redirect the user** to the `approvalUrl`. They approve in the dashboard.

3. **Receive the callback** - after approval, the dashboard redirects to `callbackUrl?code=<exchangeCode>`.

4. **Exchange the code for the key** - `POST /auth/key-request/exchange`:

```json
{ "code": "<exchangeCode>" }
```

Response:
```json
{
  "apiKey": "...",
  "scopes": ["entity:read", "roll:execute"],
  "clientIds": []
}
```

### Creating a key manually (dashboard)

1. Log in to the dashboard.
2. Go to **API Keys** → **Create New Key**.
3. Pick the scopes the integration needs (start narrow: `entity:read`, `roll:read`, etc.).
4. Optionally bind to specific Foundry clients and set a monthly limit.
5. Save the key - the dashboard shows the value once. Copy it into your integration's secrets store.

Revoke a scoped key at any time from the dashboard. Each integration should have its own key.

---

## Making authenticated requests

Include the `x-api-key` header with your **scoped API key**:

```bash
curl -X GET http://localhost:3010/structure \
  -H "x-api-key: YOUR_SCOPED_KEY_HERE"
```

---

## Unauthenticated endpoints

A small number of endpoints require no authentication:

- `GET /api/health` - health check
- `GET /api/status` - version and configuration info
- `GET /api/clients/:clientId/active` - public, rate-limited; returns only `{ active: bool }`. Used by the Foundry module's init wizard.
- `POST /auth/register` - account creation
- `POST /auth/login` - dashboard login
- `POST /auth/pair` - exchange a pairing code for a connection token (rate-limited)
- `POST /auth/forgot-password` and friends

All data-handling endpoints (`/get`, `/create`, `/upload`, etc.) require `x-api-key`.

---

## The architectural rule

> **Foundry modules never hold HTTP API credentials. The only credential a Foundry module ever stores is its own connection token (WS-only, client-scope localStorage). Anything a module needs from the relay is expressed as a WebSocket message - including operations on other Foundry worlds. Cross-world authority is bounded by per-token allowed-targets and per-token scopes.**

This rule is what makes the storage model work. Because connection tokens are the only thing a module holds, and they live exclusively in client-scope localStorage, there is no place in shared Foundry world data where a secret credential lives.

---

## Defense in depth

Beyond the credential model, the relay applies these defenses:

- **SHA-256 hashing of all secrets at rest** - session tokens, scoped API keys, connection tokens, password reset tokens. A database leak exposes no usable credential.
- **Force reset flag** - an account can be flagged from the admin panel. Until the user logs in and completes a credential reset, all scoped keys and sessions are blocked. The flag clears automatically on successful reset.
- **Connection notifications** - every WebSocket connect, disconnect, metadata-mismatch, and rejected duplicate-connection event fires through the unified notification dispatcher (Discord webhook + email).
- **Per-token IP allowlist** - connection tokens can be restricted to specific IPs/CIDRs, configurable per browser pairing from the dashboard's Connection Tokens page. Comma-separated; accepts individual IPs and CIDR ranges. Leave blank for unrestricted.
- **Scope enforcement** - scoped API keys are gated per-endpoint by scope strings; the same vocabulary gates cross-world `remote-request` operations.
- **Audit trails** - `ConnectionLogs` records every connection attempt with metadata fingerprint validation; `RemoteRequestLogs` records every cross-world action.
- **Cross-instance revocation** - token revocation propagates over Redis pub/sub so all relay instances enforce immediately.
- **Active-connection probe** (`GET /api/clients/:id/active`) - public, rate-limited, returns only a boolean so the Foundry init wizard can check slot occupancy without revealing additional state.

---

## Security warnings

- **Use a strong, unique password.** Your password is your account credential — protect it like any other sensitive login.
- **For testing the API**, create a "personal-test" scoped key with the scopes you need, use it, then delete it.
- **Foundry modules NEVER hold HTTP API credentials.** If a module asks you to paste an API key into its settings, that's a security smell — it should use a connection token via the cross-world tunnel instead.
- **If you suspect your account is compromised**, use the **Reset Credentials** button at the bottom of the dashboard immediately. It purges all scoped keys, connection tokens, and sessions.
- **Run behind HTTPS in production.** Use `wss://` for the WebSocket relay URL — `ws://` to a non-localhost host sends the connection token unencrypted and triggers a warning in the Foundry module.
