---
sidebar_position: 99
---

# Migrating to v3 Authentication

Version 3 introduces scoped API keys and a device-flow key request system. This guide covers what changed and how to update your integration.

## What Changed

| | v2 | v3 |
|---|---|---|
| Base URL | `foundryvtt-rest-api-relay.fly.dev` | `foundryrestapi.com` |
| API key type | Single master key, full access | Scoped key with explicit permissions |
| Key distribution | Copy/paste from dashboard | Device-flow: app requests, user approves in browser |
| Header | `x-api-key: <master-key>` | `x-api-key: <scoped-key>` (same header) |

The `x-api-key` header is unchanged — only the key value and how you obtain it are different.

---

## The New Key Request Flow

Instead of asking users to manually copy a key from the dashboard, your app initiates a device-flow request and the user approves it in their browser. This is similar to OAuth device authorization.

### Step 1 — Create a key request

```http
POST https://foundryrestapi.com/auth/key-request
Content-Type: application/json

{
  "appName": "My App",
  "appDescription": "Brief description of what your app does",
  "appUrl": "https://your-app-website.com",
  "scopes": ["entity:read", "roll:execute", "clients:read"],
  "clientIds": ["optional-client-id-to-scope-to"]
}
```

**Response:**
```json
{
  "code": "abc123",
  "approvalUrl": "https://foundryrestapi.com/approve/abc123",
  "expiresIn": 600,
  "expiresAt": "2026-01-01T12:10:00Z"
}
```

- `appName` and `scopes` are required. Everything else is optional.
- `clientIds` pre-scopes the key to specific Foundry clients. Users can adjust this at approval time.
- The request expires in **10 minutes**.
- No authentication is required for this endpoint.

### Step 2 — Open the approval URL

Open `approvalUrl` in the user's browser. The user will be prompted to log in to their relay account (if not already) and then approve or deny the request. They can adjust the scopes and client bindings before approving.

### Step 3 — Poll for the result

Poll `GET /auth/key-request/{code}/status` every 3–5 seconds until the status changes.

```http
GET https://foundryrestapi.com/auth/key-request/abc123/status
```

**Possible responses:**

| Status | Meaning |
|--------|---------|
| `pending` | User hasn't acted yet — keep polling |
| `approved` | User approved — response includes `apiKey` and `scopes` |
| `denied` | User denied the request |
| `expired` | Request timed out (10 min) |
| `exchanged` | Key was already retrieved — don't poll again |

**Approved response:**
```json
{
  "status": "approved",
  "apiKey": "abc123...64hexchars",
  "scopes": ["entity:read", "roll:execute", "clients:read"],
  "clientIds": ["world-abc123"]
}
```

Store the `apiKey` securely. It will not be returned again.

### Step 4 — Use the key

Pass the scoped key in the `x-api-key` header, exactly as before:

```http
GET https://foundryrestapi.com/get?clientId=world-abc123&type=Actor
x-api-key: abc123...64hexchars
```

---

## Available Scopes

Request only the scopes your app actually needs.

| Scope | What it allows |
|-------|----------------|
| `entity:read` | Read actors, items, and other entities |
| `entity:write` | Create, update, delete entities |
| `roll:read` | Read recent roll results |
| `roll:execute` | Execute rolls |
| `chat:read` | Read chat messages |
| `chat:write` | Send chat messages |
| `encounter:read` | Read combat encounter state |
| `encounter:manage` | Start, end, advance encounters |
| `macro:list` | List available macros |
| `macro:execute` | Execute macros |
| `macro:write` | Create, update, delete macros |
| `scene:read` | Read scene data and screenshots |
| `scene:write` | Modify scenes |
| `canvas:read` | Read canvas/token state |
| `canvas:write` | Modify canvas/tokens |
| `effects:read` | Read active effects |
| `effects:write` | Modify active effects |
| `user:read` | Read Foundry user list |
| `user:write` | Modify Foundry users |
| `file:read` | Download files from Foundry |
| `file:write` | Upload files to Foundry |
| `playlist:control` | Control playlists and sounds |
| `world:info` | Read world/system metadata |
| `clients:read` | List connected Foundry clients |
| `sheet:read` | Read and screenshot character sheets |
| `events:subscribe` | Subscribe to roll/chat SSE streams |
| `session:manage` | Manage headless browser sessions |
| `search` | Search entities via QuickInsert |
| `structure:read` | Read folder/compendium structure |
| `structure:write` | Modify folder structure |
| `dnd5e` | D&D 5e system-specific endpoints |
| `execute-js` | Execute arbitrary JavaScript in Foundry |

:::caution
`execute-js` grants full arbitrary code execution inside Foundry. Only request it if absolutely necessary, and users should understand what they are approving.
:::

---

## Example Implementation (Python)

This mirrors the pattern used in the Foundry REST API MIDI Integration.

```python
import requests
import time
import webbrowser

RELAY_URL = "https://foundryrestapi.com"
APP_NAME = "My App"
APP_SCOPES = ["entity:read", "roll:execute", "clients:read"]

def request_api_key():
    # Step 1: Create key request
    resp = requests.post(f"{RELAY_URL}/auth/key-request", json={
        "appName": APP_NAME,
        "appDescription": "Reads entities and executes rolls",
        "scopes": APP_SCOPES,
    })
    resp.raise_for_status()
    data = resp.json()

    code = data["code"]
    approval_url = data["approvalUrl"]
    expires_in = data.get("expiresIn", 600)

    # Step 2: Open browser
    print(f"Opening browser for approval (code: {code})")
    webbrowser.open(approval_url)

    # Step 3: Poll for result
    start = time.time()
    while time.time() - start < expires_in:
        time.sleep(3)
        status_resp = requests.get(f"{RELAY_URL}/auth/key-request/{code}/status")
        status_data = status_resp.json()
        status = status_data.get("status")

        if status == "approved":
            api_key = status_data["apiKey"]
            print(f"Approved! Scopes: {status_data['scopes']}")
            return api_key
        elif status == "denied":
            raise Exception("Key request was denied by the user")
        elif status in ("expired", "exchanged"):
            raise Exception(f"Key request ended with status: {status}")
        # status == "pending" → keep polling

    raise Exception("Key request timed out")
```

---

## Web App Flow (Callback URL)

If your app runs in a browser or has a server that can receive a redirect, pass a `callbackUrl` in the initial request. After the user approves, the relay will redirect to:

```
https://your-app.com/callback?code=<exchange_code>
```

Exchange the code for the API key:

```http
POST https://foundryrestapi.com/auth/key-request/exchange
Content-Type: application/json

{
  "code": "<exchange_code>"
}
```

**Response:**
```json
{
  "apiKey": "abc123...64hexchars",
  "scopes": ["entity:read"],
  "clientIds": []
}
```

The exchange code is single-use.

---

## Self-Hosted Relays

Everything above applies equally to self-hosted relays. Replace `foundryrestapi.com` with your relay's URL. Users approve requests against their own account on that relay.
