---
id: foundry-module
title: Foundry VTT Module Setup
sidebar_position: 3
---

# Foundry VTT Module Setup

The Foundry VTT REST API module is the counterpart to the relay server, running inside your Foundry VTT world. It connects to the relay over a WebSocket and makes your world's data available through the API.

## Installation

1. Open the Foundry VTT setup screen.
2. Navigate to the **Add-on Modules** tab.
3. Click **Install Module**.
4. In the **Manifest URL** field, paste:
   ```
   https://github.com/ThreeHats/foundryvtt-rest-api/releases/latest/download/module.json
   ```
5. Click **Install**.

Once installed, enable the module within your desired world.

## First-time pairing

The module uses **per-browser pairing**. Each GM pairs each device once. The connection token lives in that browser's localStorage and is never exposed to non-GM players.

### Steps

1. **Configure the relay URL** (one-time, world setting):
   - For the public relay: `wss://foundryrestapi.com`
   - For local development: `ws://localhost:3010`
   - For self-hosted production: `wss://your-relay-domain.example`
   - The URL is managed via the **REST API Connection** menu (not the regular settings UI).

2. **Open the Connection menu** in Foundry:
   - **Game Settings** → **Configure Settings** → **REST API Connection** → **Manage Connection**.

3. **Pair your browser** — two options:

   **Option A — Pair via Browser (recommended):**
   - Click **Pair via Browser**.
   - A new tab opens to your relay dashboard's pairing page.
   - Log in if prompted, then click **Approve** to grant the connection.
   - Return to Foundry — the module automatically stores the token and prompts to reload.

   **Option B — Enter Code manually:**
   - Log in to the relay dashboard.
   - Go to **Connection Tokens** → **Generate Pairing Code**.
   - A 6-character code appears (case-insensitive, expires in 10 minutes).
   - Back in Foundry, click **Enter Code** and paste the code.

4. **Reload Foundry**. The module connects automatically on every subsequent load from this browser.

### Pairing additional GMs

If a second GM joins the same world, the module detects that the world is paired (the `clientId` is in world settings) but this browser doesn't have a token yet. The init wizard will:

1. Probe the relay's `/api/clients/:id/active` endpoint to see if another GM is currently connected.
2. **If yes** (the slot is held), the module shows a passive info notification and stays silent.
3. **If no** (slot is free), the module opens the Connection dialog asking the GM to add their browser.

To pair the second browser, the account owner generates an **Add Browser** code from the dashboard's **Known Clients** page. The code is bound to the existing `clientId`, so the second browser pairs into the same world rather than creating a new world entry.

### Per-user "don't ask me again"

Each GM can dismiss the pairing prompt independently - clicking **Don't Ask Again** sets a per-user flag that only affects that GM. Other GMs in the same world still see the prompt.

## Configuration reference

The module's settings (Game Settings → Configure Settings → Foundry REST API Module):

- **Custom Client Name** (world setting): Optional human-readable name for this Foundry world in the dashboard's Connected Clients list.
- **Log Level** (world setting): Module log verbosity. `debug`, `info`, `warn`, `error`.
- **Ping Interval** (world setting): WebSocket keepalive interval in seconds (default: `30`).
- **Max Reconnect Attempts** (world setting): Maximum auto-reconnect attempts after a disconnect (default: `20`).
- **Reconnect Base Delay** (world setting): Initial delay before reconnect, doubles with each attempt (default: `1000`ms).
- **Code Execution Permission Level** (world setting): Minimum Foundry user role required for `execute-js` (default: `4` = GM only).
- **Allow Execute JavaScript** (world setting): Enable the API to run arbitrary JavaScript via `POST /execute-js`. Disabled by default.
- **Allow Macro Execution** (world setting): Enable the API to run macros via `POST /macro`. Disabled by default.
- **Allow Macro Creation/Editing** (world setting): Enable the API to create, update, or delete macros. Disabled by default - a malicious caller could plant code for later execution, so enable only if you trust all API key holders.
- **Notify on Execute JS** (world setting): Show an in-Foundry GM chat whisper whenever the API runs execute-js. Default on. Disable if a trusted integration calls this frequently and the notifications become noise. Discord/email notifications are controlled separately in the relay dashboard.
- **Notify on Macro Execute** (world setting): Show an in-Foundry GM chat whisper whenever the API runs a macro. Default on. Same rationale as above.
- **REST API Connection menu**: opens the unified pair/unpair/edit-URL dialog. GM-restricted.

The following are **not in the regular UI** because they're managed by the Connection menu or are credentials:

- **WebSocket Relay URL** (world, hidden): managed via the Connection menu.
- **Client ID** (world, hidden): non-secret opaque ID set by the pairing flow.
- **Paired Relay URL** (world, hidden): the URL the world was paired against; used for URL pinning.
- **Connection Token** (CLIENT scope, browser localStorage, hidden): the secret. Per-device, never broadcast.

## URL pinning

When you pair the world, the module records the relay URL it paired against. If the `wsRelayUrl` setting is later changed (e.g., by another module or via console manipulation), the module refuses to connect and prompts to re-pair. This prevents a malicious script from silently redirecting your Foundry world to an attacker-controlled relay.

## Sensitive setting change notifications

The module hooks `updateSetting` and fires a chat-whisper to all active GMs whenever any of these settings is modified:

- `wsRelayUrl`
- `allowExecuteJs`
- `allowMacroExecute`

The module also calls `notifyRelay()` on each change, which the relay dispatches to your configured destinations (Discord webhook or email).

## execute-js / macro-execute gating

By default, `execute-js` (arbitrary JavaScript execution), `execute-macro` (running existing macros), and `allow-macro-write` (creating/editing macros) are all **disabled at the world level**. To enable them:

1. Open Game Settings → Configure Settings → Foundry REST API Module.
2. Toggle **Allow Execute JS**, **Allow Macro Execution**, or **Allow Macro Creation/Editing**.

Two independent notification layers fire on each invocation:
- **Relay** - dispatches to your configured Discord webhook or email based on per-account, per-API-key, and per-world notification settings in the relay dashboard. Always fires regardless of module settings.
- **Module** - whispers active GMs in-Foundry with a script preview or macro name. Controlled by the **Notify on Execute JS** / **Notify on Macro Execute** module settings (both default on). Disable these if a trusted integration runs them frequently and the in-Foundry notifications become noise.

## Single connection per world

The relay enforces **one WebSocket connection per world at a time**. The module uses **first-come, first-serve** logic: the first full GM (role 4, not Assistant GM) whose browser connects to the relay holds the slot. Other GMs who try to connect receive close code `4004 DuplicateConnection` and wait silently - no error, no notification. When the connected GM leaves the world, the next GM whose `userDisconnected` hook fires will try to claim the slot.

This means which GM holds the connection is determined by who was already online and connected when the world was loaded, not by any kind of ranking or sorting. If the slot holder logs out of Foundry and another GM is already in the world, that GM's module will attempt to connect and become the new slot holder automatically.

There is no way to have two browsers from the same world simultaneously connected to the relay. Plan your workflows accordingly - only one GM's browser holds the relay connection at a time.

## What if I see "REST API: connected via another GM's browser"?

This is the expected message when the world is already paired and another GM is currently holding the relay connection slot. Only one connection per world at a time. To take over, the other GM must close their browser or unpair, at which point your browser can be paired and become the active connection.

## What if I see "Relay URL has been changed since pairing"?

URL pinning was triggered. The relay URL setting differs from what was set during pairing. Either:

1. The change was intentional → open the **REST API Connection** menu and re-pair (Pair via Browser or Enter Code). Completing pairing updates the pinned URL to the current one.
2. The change was unauthorized → investigate. A module or script has modified your `wsRelayUrl` setting.

## Using the module API from other modules

The REST API module exposes a public API via the standard Foundry convention (`game.modules.get(id).api`), available to any other module after the `init` hook:

```js
const restApi = game.modules.get("foundry-rest-api");
if (!restApi?.active) return; // not installed or not enabled
const api = restApi.api;
```

### Available methods

#### `api.search(query, filter?)`

Searches world entities using the module's built-in search index. Same data as the `/search` HTTP endpoint.

```js
// Search all entity types
const results = await restApi.api.search("goblin");

// Filter to a specific type
const swords = await restApi.api.search("sword", "type:Item");
```

Returns an array of matching entity objects. The index builds automatically on `ready`; calling `search()` before it's done triggers a build automatically.


#### `api.getWebSocketManager()`

Returns the active `WebSocketManager` instance, or `null` if the module isn't connected to the relay yet. Useful for registering custom message handlers or checking connection state.

```js
const wsm = restApi.api.getWebSocketManager();
if (wsm?.isConnected()) {
  wsm.onMessageType("my-custom-event", (data, ctx) => {
    console.log("Received:", data);
  });
}
```

#### `api.openConnectionDialog()`

Programmatically opens the pairing/connection dialog. Handy if your module detects the REST API isn't connected and wants to prompt the GM to pair.

```js
restApi.api.openConnectionDialog();
```

#### `api.remoteRequest(targetClientId, action, payload?, opts?)`

Invokes an action on **another Foundry world** via the WebSocket tunnel. See [Building Cross-World Foundry Modules](./cross-world-modules) for the full guide - this method requires cross-world permissions configured.

```js
const result = await restApi.api.remoteRequest(
  "fvtt_other_world",
  "get",
  { uuid: "Actor.abc123" }
);
```

## I want to write a module that interacts with OTHER Foundry worlds

See [Building Cross-World Foundry Modules](./cross-world-modules). The short version: never put an HTTP API key in your module. Use `module.api.remoteRequest(targetClientId, action, payload)` instead - it goes through the WebSocket tunnel, has per-token cross-world permissions configured, and supports headless auto-start for offline targets.
