---
title: Session
sidebar_label: Session
---

# Session

## DELETE /end-session

Ends a headless Foundry session.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| sessionId | string | No | query | The session to end |

### Returns

**object** - OK

---

## GET /session

Gets the currently active headless Foundry session.

### Returns

**object** - OK

---

## POST /session-handshake

Creates a temporary, one-time-use, token that can be used to create a headless Foundry session.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| x-foundry-url | string | No | header | The url to your foundry game |
| x-password | string | No | header | The password to log in with |
| x-username | string | No | header | The username to log in with (eg. "Gamemaster") |
| x-world-name | string | No | header | (Optional) The name of the world as it appears in foundry if the world is not already loaded. |

### Returns

**object** - OK

---

## POST /start-session

Starts a headless Foundry session. Must provide a handshake token and the encrypted password.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| body | object | No | body |  |

### Returns

**object** - OK
