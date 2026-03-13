---
title: Macro
sidebar_label: Macro
---

# Macro

## POST /macro/{uuid}/execute

Executes a macro

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | Yes | path | UUID of the macro to execute |
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| body | object | No | body |  |

### Returns

**object** - OK

---

## GET /macros

Returns all available macros

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string | No | query | Auth token to connect to specific Foundry world |

### Returns

**object** - OK
