---
id: cross-world-modules
title: Building Cross-World Foundry Modules
sidebar_position: 9
---

# Building Cross-World Foundry Modules

If you're building a Foundry module that needs to interact with **other Foundry worlds** (transferring entities, syncing folders, creating users, pushing economy updates, etc.), this is the guide for you.

Putting an HTTP API key in your module's settings defeats the relay's security model in three ways:

1. The key has to live somewhere in Foundry, and every Foundry-side storage mechanism (world settings, user flags, document flags with restricted ownership) is **broadcast to every connected client including non-GM players**. The key leaks the moment any player opens the browser console.
2. A single HTTP API key gives access to ALL of the user's worlds at once. Stealing it from one Foundry world pivots to every world the user has paired.
3. The key has no per-action gating - an integration accidentally given a key with `execute-js` scope can run arbitrary code in your world.

The relay provides a better way: the **WS tunnel via `module.api.remoteRequest`**.

## The architectural rule

> Foundry modules never hold HTTP API credentials. They use their existing connection token (WS-only, client-scope localStorage) and ask the relay to invoke actions on other worlds via WebSocket messages. Cross-world authority is bounded by per-token allowed-targets and per-token scopes.

## How it works

1. The user pairs each of their Foundry worlds with the relay normally. Each pairing produces a `clientId` and a connection token. The connection token lives in that browser's `localStorage`.
2. When generating the pair code in the dashboard, the user can optionally enable cross-world capabilities:
   - **Allowed target clients**: which other worlds this browser may interact with
   - **Remote scopes**: which actions are permitted (`entity:write`, `user:write`, `file:write`, etc.)
3. From a Foundry module, you call `module.api.remoteRequest(targetClientId, action, payload)`.
4. The relay validates that your connection token's permissions cover the action + target, then forwards the action to the target world's WebSocket connection.
5. The target world processes the action (using its existing handlers - same code as if it had been called via HTTP) and returns a response.
6. The relay routes the response back to your module as a `remote-response` message, which resolves your Promise.

If the target world is offline AND its `KnownClients.autoStartOnRemoteRequest` toggle is enabled, the relay can spawn a headless ChromeDP session for the target on the fly.

## Calling the API

```javascript
const restApi = game.modules.get("foundry-rest-api");
if (!restApi?.active) {
  ui.notifications.error("Foundry REST API module is not installed or active");
  return;
}

try {
  // Create an actor on another world
  const actorResult = await restApi.api.remoteRequest(
    "fvtt_target_world_id",
    "create",
    {
      entityType: "Actor",
      data: { name: "Alice", type: "character", /* ... */ },
    }
  );

  // Create a user with a specific password on another world
  const userResult = await restApi.api.remoteRequest(
    "fvtt_target_world_id",
    "create-user",
    {
      name: "Alice",
      role: 1, // 1 = Player, 2 = Trusted, 3 = Assistant GM, 4 = GM
      password: "secret-password",
    }
  );

  // Upload a file to another world
  const uploadResult = await restApi.api.remoteRequest(
    "fvtt_target_world_id",
    "upload",
    {
      path: "worlds/myworld/assets",
      filename: "portrait.png",
      fileData: "data:image/png;base64,iVBORw0KGgo...",
    }
  );

  // Get a folder structure
  const structure = await restApi.api.remoteRequest(
    "fvtt_target_world_id",
    "structure",
    { types: ["Actor", "Item"] }
  );
} catch (err) {
  ui.notifications.error(`Cross-world request failed: ${err.message}`);
}
```

## API reference

### `module.api.remoteRequest(targetClientId, action, payload?, opts?)`

**Parameters:**

- `targetClientId` (string): the opaque clientId of the target Foundry world. Get this from `module.api.remoteRequest(yourClientId, "clients", {})` or by inspecting `game.settings.get("foundry-rest-api", "clientId")` on the source world.
- `action` (string): the action to invoke on the target. See [Action reference](#action-reference) below.
- `payload` (object, optional): action-specific parameters. For most actions this is the same shape as the equivalent HTTP endpoint's request body.
- `opts` (object, optional):
  - `autoStartIfOffline` (boolean, default `true`): if the target is offline, ask the relay to attempt a headless auto-start. Requires the target's `autoStartOnRemoteRequest` to be enabled in the dashboard.
  - `timeoutMs` (number, default `60000`): how long to wait for the target's response before rejecting with a timeout error.

**Returns:** a Promise that resolves with the target's response data, or rejects with an `Error` describing what went wrong.

**Throws** (rejection reasons):
- `target {clientId} not in allowed clients` - the source connection token's `allowedTargetClients` doesn't include the target. Edit the token's permissions in the dashboard's Connection Tokens page.
- `scope {scope} not granted to source token` - the source token doesn't have the required scope in its `remoteScopes`. Edit permissions in the dashboard.
- `action {action} is not exposed via remote-request` - the action isn't in the relay's `ActionToScopeRequired` map. Likely a typo in your action string.
- `target offline; auto-start not configured` - the target is offline and either `autoStartIfOffline` was false or the target's `autoStartOnRemoteRequest` is disabled.
- `target offline; headless worker not available on this instance` - `autoStartIfOffline` was true but the relay instance doesn't have headless support enabled.
- `request timed out` - the target didn't respond within `timeoutMs`.
- `target not owned by source account` - defensive check; should not happen unless the relay's user→client mapping is corrupted.

## Action reference

The relay maps each action to a required scope. When generating a pair code with cross-world capability, the user must enable the appropriate scopes for the actions your module needs.

| Action | Required scope | Notes |
|--------|---------------|-------|
| `get` | `entity:read` | Get an entity by UUID/ID |
| `create` | `entity:write` | Create any entity type (Actor, Item, Scene, etc.) |
| `update` | `entity:write` | Update an existing entity |
| `delete` | `entity:write` | Delete an entity |
| `give`, `remove`, `decrease`, `increase`, `kill` | `entity:write` | Various entity manipulation actions |
| `search` | `search` | QuickInsert search |
| `rolls`, `lastroll` | `roll:read` | Roll history |
| `roll` | `roll:execute` | Execute a roll |
| `chat` | `chat:read` | Read chat messages |
| `send-chat` | `chat:write` | Send a chat message |
| `encounters`, `start-encounter`, `next-turn`, `end-encounter`, ... | `encounter:read` / `encounter:manage` | Combat tracker actions |
| `macros` | `macro:list` | List macros |
| `execute-macro` | `macro:execute` | Execute a macro (gated by ALLOW_MACRO_EXECUTE setting) |
| `scenes`, `change-scene` | `scene:read` / `scene:write` | Scene management |
| `users` | `user:read` | List users on the target world |
| `user`, `create-user` | `user:write` | Create a user with name/role/password |
| `file-system`, `download` | `file:read` | File system operations |
| `upload` | `file:write` | Upload a file |
| `create-folder` | `structure:write` | Create a folder |
| `structure` | `structure:read` | Read folder structure |
| `clients` | `clients:read` | List connected clients |
| `sheet`, `sheet-screenshot` | `sheet:read` | Character sheet rendering |
| `scene-screenshot` | `scene:read` | Scene rendering |
| `playlists`, `playlist-play`, `play-sound`, ... | `playlist:control` | Audio control |
| `world-info` | `world:info` | World metadata |
| `execute-js` | `execute-js` | Execute arbitrary JS - DANGEROUS, opt-in only via the target's ALLOW_EXECUTE_JS setting |

The full mapping is in `internal/model/scopes.go` (`ActionToScopeRequired`). If your module needs an action that isn't listed, file an issue.

## Limitations

### Source world must be connected

Cross-world `remote-request` calls are **initiated by the source world's module**, which must be actively connected to the relay for the request to go anywhere. If no GM is currently logged into the source world with the module active, the call never reaches the relay.

In practice this means your module's code runs in the source world's Foundry game session (a GM is present and the module is running). The request flows from there through the relay to the target world. The source side is always online by definition when your module code executes.

If you need to trigger cross-world actions without a live GM in either world, you'll need to make HTTP API calls directly from a backend system using a scoped API key - that approach doesn't require the source module to be connected.

### Target world

The **target world** does not need to have a GM online. If `autoStartOnRemoteRequest` is enabled for the target client in the dashboard and the relay has headless support, the relay will spawn a headless session to handle the request. See [Why headless auto-start matters](#why-headless-auto-start-matters) below.

### One connection per world

The relay allows only one WebSocket connection per world at a time (lowest-ID active full GM). This means cross-world requests from World A always land on exactly one deterministic GM client in World B - there's no risk of duplicate processing.

---

## Configuring permissions

Cross-world capability is opt-in **at pair time**. The user generates a pair code from the dashboard's Known Clients page, choosing:

1. **Which target clients** the resulting connection token may invoke `remote-request` against (multi-select from their other paired worlds)
2. **Which scopes** the token holds for those operations (checkbox grid)

Example: a user building a character transfer flow between World A and World B would pair their World A browser with `allowedTargetClients = ["fvtt_world_b"]` and `remoteScopes = ["entity:read", "entity:write", "user:read", "user:write", "file:write", "structure:read", "structure:write"]`.

## Why headless auto-start matters

The transfer use case is: a GM is in World A, wants to send an actor to World B, but no GM is currently logged into World B. Without auto-start, the transfer fails.

With auto-start enabled on World B's Known Client row, the relay will:

1. Receive the `remote-request` from World A's module
2. Detect that World B is offline
3. Look up the user's stored Foundry credentials for World B
4. Spawn a ChromeDP headless session, log into World B as the configured user
5. Wait for the headless module to connect to the relay
6. Forward the action to the now-connected headless module
7. Return the response to World A's module
8. Tear down the headless session after a configurable idle timeout

This is fully relay-managed. The source module doesn't need to know how to start headless sessions or manage credentials.

## Security implications

A connection token with cross-world capability is more sensitive than a normal connection token, but its blast radius is still bounded:

- **What it can do**: invoke the explicitly-granted scopes against the explicitly-granted target clients. Nothing else.
- **What it cannot do**: be used as an HTTP API key, modify other tokens, regenerate the master key, or access any clients not in its `allowedTargetClients` list.
- **Where it lives**: the same place as a normal connection token - client-scope localStorage in the browser the user paired. Players cannot read it.
- **If it leaks**: revoke the specific connection token from the dashboard. The relay's broadcast-disconnect machinery propagates the revocation across all relay instances and immediately closes any sessions using that token.
- **Audit trail**: every cross-world `remote-request` is logged in `RemoteRequestLogs` (`GET /auth/remote-request-logs` from the dashboard) with source token, target client, action, source IP, and outcome.
