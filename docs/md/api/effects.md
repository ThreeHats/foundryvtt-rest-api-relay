---
tag: effects
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Effects

## GET /effects

Get all active effects on an actor or token

Returns the collection of ActiveEffect documents currently applied to the specified actor or token.

**Required scope:** `effects:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to query |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of active effects

### Try It Out

<ApiTester
  method="GET"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /effects/list

List all available status effects

Returns all status effects defined by the game system's configuration. Useful for discovering valid statusId values for the add/remove effect endpoints.

**Required scope:** `effects:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of available status effects with id, name, and icon

### Try It Out

<ApiTester
  method="GET"
  path="/effects/list"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /effects

Add an active effect to an actor or token

Adds a status condition (by statusId) or a custom ActiveEffect (via effectData) to the specified actor or token.

**Required scope:** `effects:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to add the effect to |
| clientId | string |  | query | Client ID for the Foundry world |
| statusId | string |  | body, query | Standard status condition ID (e.g., "poisoned", "blinded", "prone") |
| effectData | object |  | body, query | Custom ActiveEffect data object (name, icon, duration, changes) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the add operation

### Try It Out

<ApiTester
  method="POST"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"statusId","type":"string","required":false,"source":"body"},{"name":"effectData","type":"object","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## DELETE /effects

Remove an active effect from an actor or token

Removes an effect by its document ID (effectId) or by status condition identifier (statusId).

**Required scope:** `effects:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to remove the effect from |
| clientId | string |  | query | Client ID for the Foundry world |
| effectId | string |  | body, query | The ActiveEffect document ID to remove |
| statusId | string |  | body, query | Standard status condition ID to remove (e.g., "poisoned") |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the remove operation

### Try It Out

<ApiTester
  method="DELETE"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"effectId","type":"string","required":false,"source":"body"},{"name":"statusId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

