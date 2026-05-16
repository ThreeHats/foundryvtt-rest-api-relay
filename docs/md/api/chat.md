---
tag: chat
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';
import SseTester from '@site/src/components/SseTester';

# Chat

## GET /chat

Get chat messages

Retrieves chat messages from the Foundry world with optional pagination and filtering.

**Required scope:** `chat:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| limit | number |  | query | Maximum number of messages to return (default: 10) |
| offset | number |  | query | Number of messages to skip for pagination |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |
| chatType | number |  | query | Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field. |
| speaker | string |  | query | Filter messages by speaker name or actor ID |

### Returns

**object** - Paginated list of chat messages

### Try It Out

<ApiTester
  method="GET"
  path="/chat"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"limit","type":"number","required":false,"source":"query"},{"name":"offset","type":"number","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"},{"name":"chatType","type":"number","required":false,"source":"query"},{"name":"speaker","type":"string","required":false,"source":"query"}]}
/>

---

## POST /chat

Send a chat message

Creates a new chat message in the Foundry world.

**Required scope:** `chat:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| content | string | ✓ | body | The message content (supports HTML) |
| clientId | string |  | query | Client ID for the Foundry world |
| whisper | array |  | body | Array of user IDs to whisper the message to |
| speaker | string |  | body | Actor ID to use as the message speaker |
| alias | string |  | body | Display name alias for the speaker |
| chatType | number |  | body | Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field. |
| flavor | string |  | body | Flavor text displayed above the message content |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The created chat message

### Try It Out

<ApiTester
  method="POST"
  path="/chat"
  parameters={[{"name":"content","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"whisper","type":"array","required":false,"source":"body"},{"name":"speaker","type":"string","required":false,"source":"body"},{"name":"alias","type":"string","required":false,"source":"body"},{"name":"chatType","type":"number","required":false,"source":"body"},{"name":"flavor","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## DELETE /chat/:messageId

Delete a specific chat message

Deletes a chat message by its ID. Only the message author or a GM can delete messages.

**Required scope:** `chat:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| messageId | string | ✓ | params | ID of the chat message to delete |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Success confirmation

### Try It Out

<ApiTester
  method="DELETE"
  path="/chat/:messageId"
  parameters={[{"name":"messageId","type":"string","required":true,"source":"params"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## DELETE /chat

Clear all chat messages

Flushes all chat message history. Only GMs can perform this action.

**Required scope:** `chat:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Success confirmation

### Try It Out

<ApiTester
  method="DELETE"
  path="/chat"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /chat/subscribe

Subscribe to real-time chat events via Server-Sent Events (SSE)

Opens a persistent SSE connection that streams chat events (create, update, delete) as they occur in the Foundry world.

**Required scope:** `events:subscribe`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| speaker | string |  | query | Filter events by speaker name or actor ID |
| type | number |  | query | Filter events by chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll) |
| whisperOnly | boolean |  | query | Only receive whispered messages |
| userId | string |  | query | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**SSE stream** - SSE event stream

### Try It Out

<SseTester
  path="/chat/subscribe"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"speaker","type":"string","required":false,"source":"query"},{"name":"type","type":"number","required":false,"source":"query"},{"name":"whisperOnly","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

