---
title: Structure
sidebar_label: Structure
---

# Structure

## GET /contents/{path}

Returns the contents of a folder or compendium

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| path | string | Yes | path | Folder or compendium to return |
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK

---

## GET /structure

Returns the folders and compendiums in the world

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK
