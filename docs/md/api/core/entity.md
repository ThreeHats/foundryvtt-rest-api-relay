---
title: Entity
sidebar_label: Entity
---

# Entity

## POST /create

Creates a new entity in Foundry with the given JSON

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| body | object | No | body |  |

### Returns

**object** - Created

---

## POST /decrease

Decrease an attribute

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| selected | boolean | No | query |  |
| body | object | No | body |  |

### Returns

**object** - OK

---

## DELETE /delete

Deletes an entity from Foundry

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query |  |
| selected | boolean | No | query |  |

### Returns

**object** - OK

---

## GET /get

Returns JSON data for entity

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actor | boolean | No | query | Use selected entity's actor |
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| selected | boolean | No | query | Use selected entity |

### Returns

**object** - OK

---

## POST /give

Give an item to an actor - Item uuid required - Selected actor or actor uuid - Optionally take from another actor - Optionally specify quantity

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| body | object | No | body |  |

### Returns

**object** - OK

---

## POST /increase

Increase an attribute

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| uuid | string | No | query |  |
| body | object | No | body |  |

### Returns

**object** - OK

---

## POST /kill

Reduce to 0hp and mark dead and defeated

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| selected | boolean | No | query |  |

### Returns

**object** - OK

---

## PUT /update

Updates and entity with the given JSON props

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| uuid | string | No | query |  |
| body | object | No | body |  |

### Returns

**object** - OK
