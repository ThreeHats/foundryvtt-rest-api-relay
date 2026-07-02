---
tag: roll
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';
import SseTester from '@site/src/components/SseTester';

# Roll

## GET /rolls

Get recent rolls

Retrieves a list of up to 20 recent rolls made in the Foundry world.

**Required scope:** `roll:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| limit | number |  | query | Optional limit on the number of rolls to return (default is 20) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - An array of recent rolls with details

### Try It Out

<ApiTester
  method="GET"
  path="/rolls"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"limit","type":"number","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /lastroll

Get the last roll

Retrieves the most recent roll made in the Foundry world.

**Required scope:** `roll:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The most recent roll with details

### Try It Out

<ApiTester
  method="GET"
  path="/lastroll"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /roll

Make a roll

Executes a roll with the specified formula.

**Required scope:** `roll:execute`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| formula | string | ✓ | body | The roll formula to evaluate (e.g., "1d20 + 5") |
| clientId | string |  | query | Client ID for the Foundry world |
| flavor | string |  | body | Optional flavor text for the roll |
| createChatMessage | boolean |  | body | Whether to create a chat message for the roll |
| speaker | string |  | body | The speaker for the roll |
| whisper | array |  | body | Users to whisper the roll result to |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the roll operation

### Try It Out

<ApiTester
  method="POST"
  path="/roll"
  parameters={[{"name":"formula","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"flavor","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"},{"name":"speaker","type":"string","required":false,"source":"body"},{"name":"whisper","type":"array","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /rolls/subscribe

Subscribe to real-time roll events via Server-Sent Events (SSE)

Opens a persistent SSE connection that streams roll events as they occur.

**Required scope:** `events:subscribe`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**SSE stream** - SSE event stream

### Try It Out

<SseTester
  path="/rolls/subscribe"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

