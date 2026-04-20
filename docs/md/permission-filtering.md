---
title: Permission Filtering
sidebar_position: 7
---

# Permission Filtering

All API endpoints support an optional `userId` parameter that scopes requests to what a specific Foundry VTT user can see and do. When omitted, the API operates with full GM-level access.

## How It Works

Pass `userId` as a query parameter or in the request body:

```
GET /entity/get?clientId=fvtt_3a9f1c2e4b7d8e0f&uuid=Actor.xyz&userId=playerUserId
```

The module resolves the user and checks Foundry's built-in document permission system before returning data or performing operations.

## Auto-injection via Scoped API Keys

Both `clientId` and `userId` can be bound to a scoped API key so callers never need to pass them explicitly:

- **Client scoping** — a key bound to one or more `clientIds` auto-resolves `clientId` on every request. A single-client key always targets that world; no `clientId` parameter needed.
- **User scoping** — a key can carry a `scopedUserId` (global, applied to every client) or per-client user IDs. When set, the relay injects `userId` into the request after `clientId` resolves. Callers send no `userId` at all.

This is the recommended approach for player-facing integrations: create a key scoped to the target world and the player's Foundry user ID. The integration never deals with `clientId` or `userId` bookkeeping.

```bash
# Key is scoped to fvtt_3a9f1c2e4b7d8e0f and userId=playerOne — no params needed
curl "https://your-relay/entity/get?uuid=Actor.xyz" \
  -H "x-api-key: PLAYER_SCOPED_KEY"
```

## User Resolution

The `userId` parameter (whether passed explicitly or injected by the key) accepts either:

1. **Foundry User ID** (tried first) - e.g., `userId=abc123def456`
2. **Username** (fallback, case-insensitive) - e.g., `userId=PlayerOne`

If the user is not found, the API returns an error: `"User not found: <userId>"`.

## Permission Levels

Foundry VTT has four document permission levels. The API respects these when `userId` is provided:

| Level | Read Behavior | Write Behavior |
|---|---|---|
| **NONE** | Document excluded from results | Operation denied |
| **LIMITED** | Minimal data returned (name, uuid, type, img) | Operation denied |
| **OBSERVER** | Full document data returned | Operation denied |
| **OWNER** | Full document data returned | Operation allowed |

### Read Endpoints

For read operations (GET entity, search, structure, etc.):
- **Single document**: Returns full data, limited data, or a permission error depending on the user's permission level
- **Collections**: Documents are filtered — NONE-permission documents are excluded, LIMITED documents return minimal data, OBSERVER+ documents return full data

### Write Endpoints

For write operations (update, delete, create, give, remove, kill, etc.):
- Requires **OWNER** permission on the target document
- Returns a descriptive error if the user lacks permission

## GM-Only Operations

Some operations require the user to be a GM regardless of document permissions:

| Endpoint | Reason |
|---|---|
| `POST /execute-js` | Arbitrary JavaScript execution |
| `POST /select`, `GET /selected` | Canvas token control |
| `POST /start-encounter`, `POST /end-encounter` | Combat management |
| `POST /next-turn`, `POST /next-round`, etc. | Combat navigation |
| `POST /create-folder`, `DELETE /delete-folder` | Folder management |

## File System Permissions

File system endpoints check Foundry's user permissions:

| Endpoint | Required Permission |
|---|---|
| `GET /file-system` | `FILES_BROWSE` |
| `GET /download` | `FILES_BROWSE` |
| `POST /upload` | `FILES_UPLOAD` |

## Examples

### Reading an entity as a player

```bash
# Player has OBSERVER permission on this actor — returns full data
curl "https://your-api/api/entity/get?clientId=abc&uuid=Actor.xyz&userId=playerOne"

# Player has LIMITED permission — returns only name, uuid, type, img
curl "https://your-api/api/entity/get?clientId=abc&uuid=Actor.npc123&userId=playerOne"

# Player has no permission — returns permission error
curl "https://your-api/api/entity/get?clientId=abc&uuid=Actor.secret&userId=playerOne"
```

### Writing as a player

```bash
# Player owns this actor — update succeeds
curl -X PUT "https://your-api/api/entity/update?clientId=abc&uuid=Actor.xyz&userId=playerOne" \
  -H "Content-Type: application/json" \
  -d '{"data": {"name": "New Name"}}'

# Player doesn't own this actor — returns permission error
curl -X PUT "https://your-api/api/entity/update?clientId=abc&uuid=Actor.npc123&userId=playerOne" \
  -H "Content-Type: application/json" \
  -d '{"data": {"name": "New Name"}}'
```

### GM-only operation with non-GM user

```bash
# Non-GM user tries to execute JavaScript — returns error
curl -X POST "https://your-api/api/utility/execute-js?clientId=abc&userId=playerOne" \
  -H "Content-Type: application/json" \
  -d '{"script": "return game.actors.size"}'
# Error: "User 'PlayerOne' must be a GM to execute JavaScript"
```

### Backward compatibility (no userId)

```bash
# No userId — full GM-level access (same as before)
curl "https://your-api/api/entity/get?clientId=abc&uuid=Actor.secret"
# Returns full data regardless of permissions
```
