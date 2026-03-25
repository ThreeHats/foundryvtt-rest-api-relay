---
id: scoped-keys
title: Scoped API Keys
sidebar_position: 8
---

# Scoped API Keys

Scoped API keys let you create restricted sub-keys under your master API key. Each scoped key can be locked to a specific Foundry client, user, and daily request limit.

## When to Use Scoped Keys

- **Discord bots** — give the bot a key locked to one Foundry world
- **Player apps** — create per-player keys scoped to their Foundry user ID
- **Third-party integrations** — limit access to specific endpoints and rate limits
- **Headless sessions** — store Foundry credentials on the key so sessions start with zero extra params

## Creating a Scoped Key

### Via Dashboard

1. Log into the web dashboard
2. Click the **API Keys** tab
3. Click **+ Create Scoped Key**
4. Fill in the form (only Name is required)
5. Copy the full key when it's displayed — it won't be shown again

### Via API

```bash
curl -X POST https://your-relay.fly.dev/auth/api-keys \
  -H "x-api-key: YOUR_MASTER_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Discord Bot",
    "scopedClientId": "foundry-abc123",
    "scopedUserId": "PlayerOne",
    "dailyLimit": 500,
    "expiresAt": "2026-12-31T00:00:00Z"
  }'
```

The response includes the full `key` value — save it immediately.

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
curl -X POST https://your-relay.fly.dev/auth/api-keys \
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
curl -X POST https://your-relay.fly.dev/start-session \
  -H "x-api-key: SCOPED_KEY_WITH_CREDENTIALS"
```

No handshake or encrypted password needed.

## Rate Limits

Scoped keys have two layers of rate limiting:

1. **User-level limits** — all requests (master + scoped keys) count against the parent user's daily and monthly quotas
2. **Per-key limits** — if `dailyLimit` is set, the key has its own daily counter that resets at midnight UTC

Both limits must pass for a request to succeed.

## Managing Scoped Keys

### List Keys
```bash
GET /auth/api-keys
```

### Update a Key
```bash
PATCH /auth/api-keys/:id
```

### Delete a Key
```bash
DELETE /auth/api-keys/:id
```

### Enable/Disable
```bash
PATCH /auth/api-keys/:id
{ "enabled": false }
```

## Cascade Behaviors

- **Master key regeneration** (`POST /auth/regenerate-key`): deletes all scoped keys
- **Account deletion** (`DELETE /auth/account`): deletes all scoped keys
- **Data export** (`GET /auth/export-data`): includes scoped key metadata (not key values or credentials)

## WebSocket Connections

Scoped keys can be used for WebSocket connections (e.g., the Foundry module connecting with a scoped key). The connection is registered under the parent user's master key in the client registry. If `scopedClientId` is set, the connecting client's ID must match it.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CREDENTIALS_ENCRYPTION_KEY` | (required for credential storage) | 32-byte key for AES-256-GCM, hex or base64 encoded |
| `MAX_HEADLESS_SESSIONS` | `1` | Default max concurrent headless sessions per paid user |
