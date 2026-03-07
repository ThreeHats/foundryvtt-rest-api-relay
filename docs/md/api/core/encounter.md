---
title: Encounter
sidebar_label: Encounter
---

# Encounter

## POST /add-to-encounter

Add to encounter

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| body | object | No | body |  |

### Returns

**object** - OK

---

## GET /encounters

Returns all encounters

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK

---

## POST /end-encounter

End an encounter

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK

---

## POST /last-round

Move backward 1 round

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK

---

## POST /last-turn

Move backward 1 turn

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK

---

## POST /next-round

Move forward 1 round

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK

---

## POST /next-turn

Move forward 1 turn

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK

---

## POST /remove-from-encounter

Remove from encounter

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| body | object | No | body |  |

### Returns

**object** - OK

---

## POST /start-encounter

Starts an encouter

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| body | object | No | body |  |

### Returns

**object** - Created
