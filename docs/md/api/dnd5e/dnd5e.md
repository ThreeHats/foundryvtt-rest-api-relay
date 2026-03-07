---
title: DnD5e
sidebar_label: DnD5e
---

# DnD5e

## GET /get-actor-details

get-actor-details

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | No | query |  |
| clientId | string | No | query |  |
| details | string | No | query |  |

### Returns

**object** - Successful response

---

## POST /modify-experience

modify-experience

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | No | query |  |
| amount | integer | No | query |  |
| clientId | string | No | query |  |

### Returns

**object** - Successful response

---

## POST /modify-item-charges

modify-item-charges

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | No | query |  |
| amount | string | No | query |  |
| clientId | string | No | query |  |
| itemName | string | No | query |  |

### Returns

**object** - Successful response

---

## POST /use-ability

use-item

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| abilityName | string | No | query |  |
| actorUuid | string | No | query |  |
| clientId | string | No | query |  |
| details | string | No | header |  |

### Returns

**object** - Successful response
