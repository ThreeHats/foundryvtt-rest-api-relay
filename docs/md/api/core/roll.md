---
title: Roll
sidebar_label: Roll
---

# Roll

## GET /lastroll

Returns the last roll made

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK

---

## POST /roll

Makes a new roll in Foundry

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| formula | string | Yes | body | Dice formula to execute. |
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| createChatMessage | boolean | No | body |  |
| flavor | string | No | body |  |
| itemUuid | string | No | body |  |
| speaker | string | No | body |  |
| target | string | No | body |  |
| whisper | array | No | body |  |

### Returns

**object** - OK

---

## GET /rolls

Returns up to the last 20 rolls

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| limit | integer | No | query | (Optional) Max number of rolls to return. Max 20. Default 20. |

### Returns

**object** - OK
