---
tag: encounter
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Encounter

## GET /encounters

Get all active encounters

Retrieves a list of all currently active encounters in the Foundry world.

**Required scope:** `encounter:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - An array of active encounters with details

### Try It Out

<ApiTester
  method="GET"
  path="/encounters"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /start-encounter

Start a new encounter

Initiates a new encounter in the Foundry world.

**Required scope:** `encounter:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| tokens | array |  | body | Array of token UUIDs to include in the encounter |
| startWithSelected | boolean |  | body | Whether to start with selected tokens |
| startWithPlayers | boolean |  | body | Whether to start with players |
| rollNPC | boolean |  | body | Whether to roll for NPCs |
| rollAll | boolean |  | body | Whether to roll for all tokens |
| name | string |  | body | The name of the encounter |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Details of the started encounter

### Try It Out

<ApiTester
  method="POST"
  path="/start-encounter"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"tokens","type":"array","required":false,"source":"body"},{"name":"startWithSelected","type":"boolean","required":false,"source":"body"},{"name":"startWithPlayers","type":"boolean","required":false,"source":"body"},{"name":"rollNPC","type":"boolean","required":false,"source":"body"},{"name":"rollAll","type":"boolean","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /next-turn

Advance to the next turn in the encounter

Moves the encounter to the next turn.

**Required scope:** `encounter:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| encounter | string |  | body, query | The ID of the encounter (optional, defaults to current encounter) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Details of the next turn

### Try It Out

<ApiTester
  method="POST"
  path="/next-turn"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"encounter","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /next-round

Advance to the next round in the encounter

Moves the encounter to the next round.

**Required scope:** `encounter:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| encounter | string |  | body, query | The ID of the encounter (optional, defaults to current encounter) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Details of the next round

### Try It Out

<ApiTester
  method="POST"
  path="/next-round"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"encounter","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /last-turn

Go back to the last turn in the encounter

Moves the encounter back to the last turn.

**Required scope:** `encounter:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| encounter | string |  | body, query | The ID of the encounter (optional, defaults to current encounter) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Details of the last turn

### Try It Out

<ApiTester
  method="POST"
  path="/last-turn"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"encounter","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /last-round

Go back to the last round in the encounter

Moves the encounter back to the last round.

**Required scope:** `encounter:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| encounter | string |  | body, query | The ID of the encounter (optional, defaults to current encounter) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Details of the last round

### Try It Out

<ApiTester
  method="POST"
  path="/last-round"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"encounter","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /end-encounter

End an encounter

Ends the current encounter in the Foundry world.

**Required scope:** `encounter:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| encounter | string |  | body, query | The ID of the encounter (optional, defaults to current encounter) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Details of the ended encounter

### Try It Out

<ApiTester
  method="POST"
  path="/end-encounter"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"encounter","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /add-to-encounter

Add tokens to an encounter

Adds selected tokens or specified UUIDs to the current encounter.

**Required scope:** `encounter:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| encounter | string |  | body, query | The ID of the encounter (optional, defaults to current encounter) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| uuids | array |  | body | The UUIDs of the tokens to add |
| rollInitiative | boolean |  | body | Whether to roll initiative for the added tokens |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Details of the updated encounter

### Try It Out

<ApiTester
  method="POST"
  path="/add-to-encounter"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"encounter","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"uuids","type":"array","required":false,"source":"body"},{"name":"rollInitiative","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /remove-from-encounter

Remove tokens from an encounter

Removes selected tokens or specified UUIDs from the current encounter.

**Required scope:** `encounter:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| encounter | string |  | body, query | The ID of the encounter (optional, defaults to current encounter) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| uuids | array |  | body | The UUIDs of the tokens to remove |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Details of the updated encounter

### Try It Out

<ApiTester
  method="POST"
  path="/remove-from-encounter"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"encounter","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"uuids","type":"array","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

