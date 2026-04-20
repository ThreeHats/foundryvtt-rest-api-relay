---
id: websocket
title: WebSocket Integration
sidebar_label: WebSocket
sidebar_position: 8
---

# WebSocket Integration

The relay supports a persistent WebSocket connection for external integrations that want real-time bidirectional communication with Foundry VTT — the same operations as the REST API, plus live event subscriptions.

## Two WebSocket endpoints

The relay exposes two distinct WebSocket endpoints with different audiences and authentication:

| Endpoint | Audience | Auth | Purpose |
|----------|---------|------|---------|
| `/relay` | Foundry module ONLY | Connection token (auth-via-first-message) | The Foundry module's relay connection |
| `/ws/api` | External integrations | Scoped API key OR connection token | Persistent WS sessions for apps that want streaming events |

This page documents `/ws/api`. The Foundry module's `/relay` endpoint is internal to the module — you don't connect to it directly from your code.

## Connection (`/ws/api`)

Connect to the client API WebSocket endpoint and authenticate with the **first message** after the socket opens.

```
ws://<host>/ws/api?clientId=<clientId>
```

After the WebSocket opens, send your auth payload as the first message:

```json
{
  "type": "auth",
  "token": "YOUR_SCOPED_API_KEY"
}
```

The relay responds with `{ "type": "auth-success" }` on success, or closes the connection with code `4002` on failure. After successful auth, you can send any of the message types listed below.

`clientId` is auto-resolved when omitted: if your scoped key is bound to one client (`scopedClientId`) it's used automatically; if not, the relay uses the unique connected client under your account; if multiple clients are connected, you must specify which one.

## Message Format

All messages are JSON objects with a `type` field. Request messages must also include a `requestId` for correlation.

### Request

```json
{
  "type": "search",
  "requestId": "my-unique-id",
  "query": "dragon"
}
```

### Response

```json
{
  "type": "search-result",
  "requestId": "my-unique-id",
  "clientId": "abc123",
  "results": [...]
}
```

## Event Subscriptions

Subscribe to real-time events from Foundry:

```json
{
  "type": "subscribe",
  "channel": "chat-events",
  "filters": { "speaker": "GM" }
}
```

Available channels: `chat-events`, `roll-events`

To unsubscribe:

```json
{
  "type": "unsubscribe",
  "channel": "chat-events"
}
```

## Supported Message Types

The full per-message reference — parameters, descriptions, and interactive testers — is in the [WebSocket API Reference](/api/websocket).

The AsyncAPI specification is also available at `/asyncapi.json`.
