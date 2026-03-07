---
title: Utility
sidebar_label: Utility
---

# Utility

## POST /execute-js

#### Executes Javascript in Foundry Accepts the script as a file upload, or as a raw string in the body with the key "script". If included as a raw string in the body excape quotes and backslashes, and remove comments. Returns the result of the code execution.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query |  |
| body | object | No | body |  |

### Returns

**object** - OK

---

## POST /select

**Selects entities in Foundry**

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| body | object | No | body |  |

### Returns

**object** - OK

---

## GET /selected

**Returns the currently selected entities**

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK
