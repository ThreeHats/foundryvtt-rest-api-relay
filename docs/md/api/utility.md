---
tag: utility
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Utility

## POST /select

Select token(s)

Selects one or more tokens in the Foundry VTT client.

**Required scope:** `entity:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| uuids | array |  | body | Array of UUIDs to select |
| name | string |  | body | Name of the token(s) to select |
| data | object |  | body | Data to match for selection (e.g., "data.attributes.hp.value": 20) |
| overwrite | boolean |  | body | Whether to overwrite existing selection |
| all | boolean |  | body | Whether to select all tokens on the canvas |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The selected token(s)

### Try It Out

<ApiTester
  method="POST"
  path="/select"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuids","type":"array","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"data","type":"object","required":false,"source":"body"},{"name":"overwrite","type":"boolean","required":false,"source":"body"},{"name":"all","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /selected

Get selected token(s)

Retrieves the currently selected token(s) in the Foundry VTT client.

**Required scope:** `entity:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The selected token(s)

### Try It Out

<ApiTester
  method="GET"
  path="/selected"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /players

Get players/users

Retrieves a list of all users configured in the Foundry VTT world. Useful for discovering valid userId values for permission-scoped API calls.

**Required scope:** `entity:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - List of users with their IDs, names, roles, and active status

### Try It Out

<ApiTester
  method="GET"
  path="/players"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /world-info

Get comprehensive world information

Returns a single object with world name, game system, Foundry version, all modules (with active status), all users (with online status), and the active scene. Useful for API clients to discover the world state.

**Required scope:** `world:info`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - World information object

### Try It Out

<ApiTester
  method="GET"
  path="/world-info"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /execute-js

Execute JavaScript

Executes a JavaScript script in the Foundry VTT client.

**Required scope:** `execute-js`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| script | string |  | body | JavaScript script to execute |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The result of the executed script

### Try It Out

<ApiTester
  method="POST"
  path="/execute-js"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"script","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

