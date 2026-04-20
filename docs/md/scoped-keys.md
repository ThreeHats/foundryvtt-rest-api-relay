---
id: scoped-keys
title: Scoped API Keys
sidebar_position: 6
---

# Scoped API Keys

Scoped API keys let you create restricted sub-keys under your master API key. Each scoped key can be locked to a specific Foundry client, user, daily request limit, and a set of action scopes.

**Use scoped keys for any HTTP integration with the relay.** They're the right credential for Discord bots, custom apps, scripts, Obsidian plugins - anything that calls the relay's REST API. The master API key should never be used for routine HTTP calls; reach for it only for emergency programmatic root access.

## When to use scoped keys

- **Discord bots** - give the bot a key locked to one Foundry world with read-only entity scopes
- **Player apps** - create per-player keys scoped to their Foundry user ID
- **Third-party integrations** - limit access to specific endpoints and rate limits
- **Headless sessions (legacy)** - store Foundry credentials directly on the key so sessions can start with zero extra params

## When NOT to use scoped keys

- **Don't put a scoped key inside a Foundry module.** That's not what scoped keys are for. Foundry modules use connection tokens (WS-only, per-browser) - see [Building Cross-World Foundry Modules](./cross-world-modules) for the cross-world tunnel pattern.
- **Don't reuse the same scoped key across multiple integrations.** Each integration should have its own key. When you decommission one, revoke that key without affecting the others.
- **Don't store a scoped key in a public repo or pass it around in chat.** They're secrets, just lower-blast-radius secrets than the master key.

## Creating a Scoped Key

### Integration-initiated request (preferred for third-party apps)

If you're building an integration for other people's Foundry instances, use the key-request flow. Your integration requests a key with the scopes it needs; the user reviews and approves it in the dashboard; the key is delivered back to your integration automatically. Users never have to touch the dashboard to set up your integration.

See [Authentication & Security Model - Requesting a key from your integration](./authentication#requesting-a-key-from-your-integration-preferred) for the full flow, including device flow (poll-based) and web flow (callback redirect).

### Via Dashboard

1. Log into the web dashboard
2. Click the **API Keys** tab
3. Click **+ Create Scoped Key**
4. Fill in the form (only Name is required)
5. Copy the full key when it's displayed - it won't be shown again

### Via API (master key)

```bash
curl -X POST https://your-relay.com/auth/api-keys \
  -H "x-api-key: YOUR_MASTER_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Discord Bot",
    "scopedClientId": "foundry-abc123",
    "scopedUserId": "PlayerOne",
    "monthlyLimit": 500,
    "expiresAt": "2026-12-31T00:00:00Z"
  }'
```

The response includes the full `key` value - save it immediately.

## Scope Enforcement

Scoped keys enforce restrictions server-side. The caller **cannot override** scoped values.

| Scope | Behavior |
|-------|----------|
| `scopedClientId` | All requests are forced to this client, regardless of what the caller sends |
| `scopedUserId` | The `userId` parameter is always set to this value for permission filtering |
| `dailyLimit` | Returns 429 when the key's per-day counter exceeds this limit |
| `expiresAt` | Returns 401 after this timestamp |

## clientId Auto-Resolution

When `clientId` is omitted from a request:

1. If the key has a `scopedClientId`, that value is used automatically
2. Otherwise, the API checks how many Foundry clients are connected under the master key
3. If exactly one client is connected, it's used automatically
4. If zero clients → 404 error
5. If multiple clients → 400 error listing available clients

This means most scoped keys don't need to specify `clientId` at all.

## Stored Credentials

Scoped keys can store encrypted Foundry login credentials for headless sessions:

```bash
curl -X POST https://your-relay.com/auth/api-keys \
  -H "x-api-key: YOUR_MASTER_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Auto-Session Key",
    "foundryUrl": "https://my-foundry.com",
    "foundryUsername": "Gamemaster",
    "foundryPassword": "secret123"
  }'
```

Passwords are encrypted with AES-256-GCM at rest. The server must have `CREDENTIALS_ENCRYPTION_KEY` configured.

With stored credentials, starting a headless session requires only:

```bash
curl -X POST https://your-relay.com/start-session \
  -H "x-api-key: SCOPED_KEY_WITH_CREDENTIALS"
```

No handshake or encrypted password needed.

## Rate Limits

Scoped keys have two layers of rate limiting:

1. **User-level limits** - all requests (master + scoped keys) count against the parent user's daily and monthly quotas
2. **Per-key limits** - if `dailyLimit` is set, the key has its own daily counter that resets at midnight UTC

Both limits must pass for a request to succeed.

## Managing Scoped Keys

The **API Keys** page in the dashboard is the intended place to manage scoped keys. From there you can see all keys, their usage counters, enabled/disabled state, expiry, and scopes — and rename, toggle, or delete any of them with a click.

If you need to automate key management (CI pipelines, provisioning scripts), the same operations are available via the API using your master key:

| Operation | Endpoint |
|-----------|----------|
| List keys | `GET /auth/api-keys` |
| Update a key | `PATCH /auth/api-keys/:id` |
| Enable / disable | `PATCH /auth/api-keys/:id` with `{ "enabled": false }` |
| Delete a key | `DELETE /auth/api-keys/:id` |

## Cascade Behaviors

- **Master key regeneration** (`POST /auth/regenerate-key`): deletes all scoped keys, deletes all connection tokens (Foundry modules need to re-pair), invalidates all dashboard sessions across every device, and force-disconnects every active WebSocket. Rotation is a panic button - be sure before you trigger it.
- **Account deletion** (`DELETE /auth/account`): deletes all scoped keys, all connection tokens, all credentials, all known clients, all logs.
- **Data export** (`GET /auth/export-data`): includes scoped key metadata (not key values, credentials, or session tokens).

## WebSocket connections

Scoped keys can authenticate the **client API** WebSocket endpoint (`/ws/api`) - used by external integrations that want a persistent WS session for streaming events. They cannot be used for the Foundry module's `/relay` endpoint, which only accepts connection tokens (and requires the auth-via-first-message handshake).

This split is intentional:

- **`/relay`** - Foundry module side. Uses a connection token sent as the first WebSocket message after connect. Connection tokens are per-browser and stored in client-scope localStorage.
- **`/ws/api`** - External integration side. Accepts both connection tokens (for module-style use cases) AND scoped API keys (for HTTP-style use cases via the WS transport).

If you're building an external app that wants to subscribe to live events from a Foundry world, use a scoped key on `/ws/api`. If you're building a Foundry module, you don't need either - the module already has its own connection token from the pairing flow, and `module.api.remoteRequest()` is the cross-world API.

See [WebSocket API](./websocket) for the full connection flow and message reference.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CREDENTIALS_ENCRYPTION_KEY` | (required for credential storage) | 32-byte key for AES-256-GCM, hex or base64 encoded |
| `MAX_HEADLESS_SESSIONS` | `1` | Default max concurrent headless sessions per paid user |
