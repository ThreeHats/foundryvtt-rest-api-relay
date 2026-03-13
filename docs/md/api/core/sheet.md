---
title: Sheet
sidebar_label: Sheet
---

# Sheet

## GET /sheet

Returns raw HTML (or a string in a JSON response) for an entity If returning HTML there are options for scale, tab to open (if available), and darkMode. If returning JSON the HTML is untouched.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actor | boolean | No | query | (Optional) Return dark mode HTML. Default = false |
| clientId | string | No | query | Auth token to connect to specific Foundry world |
| scale | number | No | query |  |
| selected | boolean | No | query | (Optional) Tab to open if available |

### Returns

**string** - OK
