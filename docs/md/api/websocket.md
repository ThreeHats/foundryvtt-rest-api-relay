---
tag: WebSocket
---

import WsTester from '@site/src/components/WsTester';
import WsMessageTester, { WsConnectionBar } from '@site/src/components/WsMessageTester';

# WebSocket API

The WebSocket API provides bidirectional communication with Foundry VTT through the relay server. It supports the same operations as the REST API, plus real-time event subscriptions.

## Connection

Connect to `/ws/api` and authenticate with the **first message** after the socket opens.

```
ws://<host>/ws/api?clientId=<clientId>
```

After the WebSocket opens, send your auth payload as the first message:

```json
{
  "type": "auth",
  "token": "YOUR_SCOPED_API_KEY"
}
```

The relay responds with `{ "type": "auth-success" }` on success, or closes the connection with code `4002` on failure.

`clientId` is auto-resolved when omitted: if your scoped key is bound to one client it is used automatically; if multiple clients are connected, you must specify which one.

## Message Format

All messages are JSON objects with a `type` field. Request messages must also include a `requestId` for correlation.

### Request

```json
{
  "type": "search",
  "requestId": "my-unique-id",
  "query": "dragon"
}
```

### Response

```json
{
  "type": "search-result",
  "requestId": "my-unique-id",
  "clientId": "abc123",
  "results": [...]
}
```

## Event Subscriptions

Subscribe to real-time events from Foundry:

```json
{
  "type": "subscribe",
  "channel": "chat-events",
  "filters": { "speaker": "GM" }
}
```

Available channels: `chat-events`, `roll-events`

To unsubscribe:

```json
{
  "type": "unsubscribe",
  "channel": "chat-events"
}
```

## Try It Out

Use the connection bar below to connect once — all message testers on this page share the same connection.

<WsConnectionBar />

---

## Entity

### `entity`

Get entity details

This endpoint retrieves the details of a specific entity.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `uuid` | string | no | UUID of the entity to retrieve (optional if selected=true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `actor` | boolean | no | Return the actor of specified entity |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="entity" parameters={[{"name":"uuid","type":"string","required":false,"description":"UUID of the entity to retrieve (optional if selected=true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"actor","type":"boolean","required":false,"description":"Return the actor of specified entity"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `create`

Create a new entity

This endpoint creates a new entity in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `entityType` | string | **yes** | Document type of entity to create (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist, ext.) |
| `data` | object | **yes** | Data for the new entity |
| `folder` | string | no | Optional folder UUID to place the new entity in |
| `keepId` | boolean | no | If true, preserve the _id from the provided data instead of generating a new one |
| `override` | boolean | no | If true and keepId is set, replace any existing entity with the same ID |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="create" parameters={[{"name":"entityType","type":"string","required":true,"description":"Document type of entity to create (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist, ext.)"},{"name":"data","type":"object","required":true,"description":"Data for the new entity"},{"name":"folder","type":"string","required":false,"description":"Optional folder UUID to place the new entity in"},{"name":"keepId","type":"boolean","required":false,"description":"If true, preserve the _id from the provided data instead of generating a new one"},{"name":"override","type":"boolean","required":false,"description":"If true and keepId is set, replace any existing entity with the same ID"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `update`

Update an entity

This endpoint updates an existing entity in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `data` | object | **yes** | Object containing the fields to update |
| `uuid` | string | no | UUID of the entity to retrieve (optional if selected=true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `actor` | boolean | no | Whether to update the actor of specified entity |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="update" parameters={[{"name":"data","type":"object","required":true,"description":"Object containing the fields to update"},{"name":"uuid","type":"string","required":false,"description":"UUID of the entity to retrieve (optional if selected=true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"actor","type":"boolean","required":false,"description":"Whether to update the actor of specified entity"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `delete`

Delete an entity

This endpoint deletes an entity from the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `uuid` | string | no | UUID of the entity to retrieve (optional if selected=true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="delete" parameters={[{"name":"uuid","type":"string","required":false,"description":"UUID of the entity to retrieve (optional if selected=true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `give`

Give an item to another entity

Transfers an item from one entity to another.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `fromUuid` | string | no | UUID of the entity giving the item |
| `toUuid` | string | no | UUID of the entity receiving the item |
| `selected` | boolean | no | Whether to get the selected entity |
| `itemUuid` | string | no | UUID of the item to give |
| `itemName` | string | no | Name of the item to give (alternative to itemUuid) |
| `quantity` | number | no | Quantity of items to give |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="give" parameters={[{"name":"fromUuid","type":"string","required":false,"description":"UUID of the entity giving the item"},{"name":"toUuid","type":"string","required":false,"description":"UUID of the entity receiving the item"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"itemUuid","type":"string","required":false,"description":"UUID of the item to give"},{"name":"itemName","type":"string","required":false,"description":"Name of the item to give (alternative to itemUuid)"},{"name":"quantity","type":"number","required":false,"description":"Quantity of items to give"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `remove`

Remove an item from an entity

Removes an item from an entity's inventory.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | no | UUID of the actor to remove the item from |
| `selected` | boolean | no | Whether to get the selected entity |
| `itemUuid` | string | no | UUID of the item to remove |
| `itemName` | string | no | Name of the item to remove (alternative to itemUuid) |
| `quantity` | number | no | Quantity of items to remove |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="remove" parameters={[{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor to remove the item from"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"itemUuid","type":"string","required":false,"description":"UUID of the item to remove"},{"name":"itemName","type":"string","required":false,"description":"Name of the item to remove (alternative to itemUuid)"},{"name":"quantity","type":"number","required":false,"description":"Quantity of items to remove"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `decrease`

Decrease an attribute

Decreases a numeric attribute of an entity by the specified amount.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `attribute` | string | **yes** | The attribute to decrease (e.g., hp.value) |
| `amount` | number | **yes** | The amount to decrease by |
| `uuid` | string | no | UUID of the entity to retrieve (optional if selected=true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="decrease" parameters={[{"name":"attribute","type":"string","required":true,"description":"The attribute to decrease (e.g., hp.value)"},{"name":"amount","type":"number","required":true,"description":"The amount to decrease by"},{"name":"uuid","type":"string","required":false,"description":"UUID of the entity to retrieve (optional if selected=true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `increase`

Increase an attribute

Increases a numeric attribute of an entity by the specified amount.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `attribute` | string | **yes** | The attribute to increase (e.g., hp.value) |
| `amount` | number | **yes** | The amount to increase by |
| `uuid` | string | no | UUID of the entity to retrieve (optional if selected=true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="increase" parameters={[{"name":"attribute","type":"string","required":true,"description":"The attribute to increase (e.g., hp.value)"},{"name":"amount","type":"number","required":true,"description":"The amount to increase by"},{"name":"uuid","type":"string","required":false,"description":"UUID of the entity to retrieve (optional if selected=true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `kill`

Kill an entity

Sets the entity's HP to 0.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `uuid` | string | no | UUID of the entity to retrieve (optional if selected=true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="kill" parameters={[{"name":"uuid","type":"string","required":false,"description":"UUID of the entity to retrieve (optional if selected=true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Structure

### `structure`

Get the structure of the Foundry world

Retrieves the folder and compendium structure for the specified Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `includeEntityData` | boolean | no | Whether to include full entity data or just UUIDs and names |
| `path` | string | no | Path to read structure from (null = root) |
| `recursive` | boolean | no | Whether to read down the folder tree |
| `recursiveDepth` | number | no | Depth to recurse into folders (default 5) |
| `types` | string | no | Types to return (Scene/Actor/Item/JournalEntry/RollTable/Cards/Macro/Playlist), can be comma-separated or JSON array |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="structure" parameters={[{"name":"includeEntityData","type":"boolean","required":false,"description":"Whether to include full entity data or just UUIDs and names"},{"name":"path","type":"string","required":false,"description":"Path to read structure from (null = root)"},{"name":"recursive","type":"boolean","required":false,"description":"Whether to read down the folder tree"},{"name":"recursiveDepth","type":"number","required":false,"description":"Depth to recurse into folders (default 5)"},{"name":"types","type":"string","required":false,"description":"Types to return (Scene/Actor/Item/JournalEntry/RollTable/Cards/Macro/Playlist), can be comma-separated or JSON array"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `get-folder`

Get a specific folder by name

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | **yes** | Name of the folder to retrieve |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-folder" parameters={[{"name":"name","type":"string","required":true,"description":"Name of the folder to retrieve"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `create-folder`

Create a new folder

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | **yes** | Name of the new folder |
| `folderType` | string | **yes** | Type of folder (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist) |
| `parentFolderId` | string | no | ID of the parent folder (optional for root level) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="create-folder" parameters={[{"name":"name","type":"string","required":true,"description":"Name of the new folder"},{"name":"folderType","type":"string","required":true,"description":"Type of folder (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist)"},{"name":"parentFolderId","type":"string","required":false,"description":"ID of the parent folder (optional for root level)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `delete-folder`

Delete a folder

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `folderId` | string | **yes** | ID of the folder to delete |
| `deleteAll` | boolean | no | Whether to delete all entities in the folder |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="delete-folder" parameters={[{"name":"folderId","type":"string","required":true,"description":"ID of the folder to delete"},{"name":"deleteAll","type":"boolean","required":false,"description":"Whether to delete all entities in the folder"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Encounter

### `encounters`

Get all active encounters

Retrieves a list of all currently active encounters in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="encounters" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `start-encounter`

Start a new encounter

Initiates a new encounter in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tokens` | array | no | Array of token UUIDs to include in the encounter |
| `startWithSelected` | boolean | no | Whether to start with selected tokens |
| `startWithPlayers` | boolean | no | Whether to start with players |
| `rollNPC` | boolean | no | Whether to roll for NPCs |
| `rollAll` | boolean | no | Whether to roll for all tokens |
| `name` | string | no | The name of the encounter |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="start-encounter" parameters={[{"name":"tokens","type":"array","required":false,"description":"Array of token UUIDs to include in the encounter"},{"name":"startWithSelected","type":"boolean","required":false,"description":"Whether to start with selected tokens"},{"name":"startWithPlayers","type":"boolean","required":false,"description":"Whether to start with players"},{"name":"rollNPC","type":"boolean","required":false,"description":"Whether to roll for NPCs"},{"name":"rollAll","type":"boolean","required":false,"description":"Whether to roll for all tokens"},{"name":"name","type":"string","required":false,"description":"The name of the encounter"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `next-turn`

Advance to the next turn in the encounter

Moves the encounter to the next turn.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `encounter` | string | no | The ID of the encounter (optional, defaults to current encounter) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="next-turn" parameters={[{"name":"encounter","type":"string","required":false,"description":"The ID of the encounter (optional, defaults to current encounter)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `next-round`

Advance to the next round in the encounter

Moves the encounter to the next round.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `encounter` | string | no | The ID of the encounter (optional, defaults to current encounter) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="next-round" parameters={[{"name":"encounter","type":"string","required":false,"description":"The ID of the encounter (optional, defaults to current encounter)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `last-turn`

Go back to the last turn in the encounter

Moves the encounter back to the last turn.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `encounter` | string | no | The ID of the encounter (optional, defaults to current encounter) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="last-turn" parameters={[{"name":"encounter","type":"string","required":false,"description":"The ID of the encounter (optional, defaults to current encounter)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `last-round`

Go back to the last round in the encounter

Moves the encounter back to the last round.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `encounter` | string | no | The ID of the encounter (optional, defaults to current encounter) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="last-round" parameters={[{"name":"encounter","type":"string","required":false,"description":"The ID of the encounter (optional, defaults to current encounter)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `end-encounter`

End an encounter

Ends the current encounter in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `encounter` | string | no | The ID of the encounter (optional, defaults to current encounter) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="end-encounter" parameters={[{"name":"encounter","type":"string","required":false,"description":"The ID of the encounter (optional, defaults to current encounter)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `add-to-encounter`

Add tokens to an encounter

Adds selected tokens or specified UUIDs to the current encounter.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `encounter` | string | no | The ID of the encounter (optional, defaults to current encounter) |
| `selected` | boolean | no | Whether to get the selected entity |
| `uuids` | array | no | The UUIDs of the tokens to add |
| `rollInitiative` | boolean | no | Whether to roll initiative for the added tokens |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="add-to-encounter" parameters={[{"name":"encounter","type":"string","required":false,"description":"The ID of the encounter (optional, defaults to current encounter)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"uuids","type":"array","required":false,"description":"The UUIDs of the tokens to add"},{"name":"rollInitiative","type":"boolean","required":false,"description":"Whether to roll initiative for the added tokens"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `remove-from-encounter`

Remove tokens from an encounter

Removes selected tokens or specified UUIDs from the current encounter.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `encounter` | string | no | The ID of the encounter (optional, defaults to current encounter) |
| `selected` | boolean | no | Whether to get the selected entity |
| `uuids` | array | no | The UUIDs of the tokens to remove |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="remove-from-encounter" parameters={[{"name":"encounter","type":"string","required":false,"description":"The ID of the encounter (optional, defaults to current encounter)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"uuids","type":"array","required":false,"description":"The UUIDs of the tokens to remove"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Roll

### `rolls`

Get recent rolls

Retrieves a list of up to 20 recent rolls made in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `limit` | number | no | Optional limit on the number of rolls to return (default is 20) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="rolls" parameters={[{"name":"limit","type":"number","required":false,"description":"Optional limit on the number of rolls to return (default is 20)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `last-roll`

Get the last roll

Retrieves the most recent roll made in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="last-roll" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `roll`

Make a roll

Executes a roll with the specified formula.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `formula` | string | **yes** | The roll formula to evaluate (e.g., "1d20 + 5") |
| `flavor` | string | no | Optional flavor text for the roll |
| `createChatMessage` | boolean | no | Whether to create a chat message for the roll |
| `speaker` | string | no | The speaker for the roll |
| `whisper` | array | no | Users to whisper the roll result to |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="roll" parameters={[{"name":"formula","type":"string","required":true,"description":"The roll formula to evaluate (e.g., \"1d20 + 5\")"},{"name":"flavor","type":"string","required":false,"description":"Optional flavor text for the roll"},{"name":"createChatMessage","type":"boolean","required":false,"description":"Whether to create a chat message for the roll"},{"name":"speaker","type":"string","required":false,"description":"The speaker for the roll"},{"name":"whisper","type":"array","required":false,"description":"Users to whisper the roll result to"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Chat

### `chat-messages`

Get chat messages

Retrieves chat messages from the Foundry world with optional pagination and filtering.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `limit` | number | no | Maximum number of messages to return (default: 10) |
| `offset` | number | no | Number of messages to skip for pagination |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |
| `chatType` | number | no | Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field. |
| `speaker` | string | no | Filter messages by speaker name or actor ID |

<WsMessageTester messageType="chat-messages" parameters={[{"name":"limit","type":"number","required":false,"description":"Maximum number of messages to return (default: 10)"},{"name":"offset","type":"number","required":false,"description":"Number of messages to skip for pagination"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"},{"name":"chatType","type":"number","required":false,"description":"Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field."},{"name":"speaker","type":"string","required":false,"description":"Filter messages by speaker name or actor ID"}]} />

---

### `chat-send`

Send a chat message

Creates a new chat message in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `content` | string | **yes** | The message content (supports HTML) |
| `whisper` | array | no | Array of user IDs to whisper the message to |
| `speaker` | string | no | Actor ID to use as the message speaker |
| `alias` | string | no | Display name alias for the speaker |
| `chatType` | number | no | Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field. |
| `flavor` | string | no | Flavor text displayed above the message content |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="chat-send" parameters={[{"name":"content","type":"string","required":true,"description":"The message content (supports HTML)"},{"name":"whisper","type":"array","required":false,"description":"Array of user IDs to whisper the message to"},{"name":"speaker","type":"string","required":false,"description":"Actor ID to use as the message speaker"},{"name":"alias","type":"string","required":false,"description":"Display name alias for the speaker"},{"name":"chatType","type":"number","required":false,"description":"Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field."},{"name":"flavor","type":"string","required":false,"description":"Flavor text displayed above the message content"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `chat-delete`

Delete a specific chat message

Deletes a chat message by its ID. Only the message author or a GM can delete messages.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `messageId` | string | **yes** | ID of the chat message to delete |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="chat-delete" parameters={[{"name":"messageId","type":"string","required":true,"description":"ID of the chat message to delete"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `chat-flush`

Clear all chat messages

Flushes all chat message history. Only GMs can perform this action.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="chat-flush" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Effects

### `get-effects`

Get all active effects on an actor or token

Returns the collection of ActiveEffect documents currently applied to the specified actor or token.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `uuid` | string | **yes** | UUID of the actor or token to query |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-effects" parameters={[{"name":"uuid","type":"string","required":true,"description":"UUID of the actor or token to query"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `get-status-effects`

List all available status effects

Returns all status effects defined by the game system's configuration. Useful for discovering valid statusId values for the add/remove effect endpoints.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-status-effects" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `add-effect`

Add an active effect to an actor or token

Adds a status condition (by statusId) or a custom ActiveEffect (via effectData) to the specified actor or token.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `uuid` | string | **yes** | UUID of the actor or token to add the effect to |
| `statusId` | string | no | Standard status condition ID (e.g., "poisoned", "blinded", "prone") |
| `effectData` | object | no | Custom ActiveEffect data object (name, icon, duration, changes) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="add-effect" parameters={[{"name":"uuid","type":"string","required":true,"description":"UUID of the actor or token to add the effect to"},{"name":"statusId","type":"string","required":false,"description":"Standard status condition ID (e.g., \"poisoned\", \"blinded\", \"prone\")"},{"name":"effectData","type":"object","required":false,"description":"Custom ActiveEffect data object (name, icon, duration, changes)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `remove-effect`

Remove an active effect from an actor or token

Removes an effect by its document ID (effectId) or by status condition identifier (statusId).

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `uuid` | string | **yes** | UUID of the actor or token to remove the effect from |
| `effectId` | string | no | The ActiveEffect document ID to remove |
| `statusId` | string | no | Standard status condition ID to remove (e.g., "poisoned") |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="remove-effect" parameters={[{"name":"uuid","type":"string","required":true,"description":"UUID of the actor or token to remove the effect from"},{"name":"effectId","type":"string","required":false,"description":"The ActiveEffect document ID to remove"},{"name":"statusId","type":"string","required":false,"description":"Standard status condition ID to remove (e.g., \"poisoned\")"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Scene

### `get-scene`

Get scene(s)

Retrieves one or more scenes by ID, name, active status, viewed status, or all.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sceneId` | string | no | ID of a specific scene to retrieve |
| `name` | string | no | Name of the scene to retrieve |
| `active` | boolean | no | Set to true to get the currently active scene |
| `viewed` | boolean | no | Set to true to get the currently viewed scene |
| `all` | boolean | no | Set to true to get all scenes |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-scene" parameters={[{"name":"sceneId","type":"string","required":false,"description":"ID of a specific scene to retrieve"},{"name":"name","type":"string","required":false,"description":"Name of the scene to retrieve"},{"name":"active","type":"boolean","required":false,"description":"Set to true to get the currently active scene"},{"name":"viewed","type":"boolean","required":false,"description":"Set to true to get the currently viewed scene"},{"name":"all","type":"boolean","required":false,"description":"Set to true to get all scenes"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `scene-raw-image`

Get the raw background image of a scene

Returns the scene's background image file without any tokens, lights, or other canvas elements rendered on it.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sceneId` | string | no | Scene ID (defaults to viewed/active scene) |
| `active` | boolean | no | If true, explicitly use the player-facing active scene instead of the viewed scene |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="scene-raw-image" parameters={[{"name":"sceneId","type":"string","required":false,"description":"Scene ID (defaults to viewed/active scene)"},{"name":"active","type":"boolean","required":false,"description":"If true, explicitly use the player-facing active scene instead of the viewed scene"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `create-scene`

Create a new scene

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `data` | object | **yes** | Scene data object (name, width, height, grid, etc.) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="create-scene" parameters={[{"name":"data","type":"object","required":true,"description":"Scene data object (name, width, height, grid, etc.)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `update-scene`

Update an existing scene

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `data` | object | **yes** | Object containing the scene fields to update |
| `sceneId` | string | no | ID of the scene to update |
| `name` | string | no | Name of the scene to update (alternative to sceneId) |
| `active` | boolean | no | Set to true to target the active scene |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="update-scene" parameters={[{"name":"data","type":"object","required":true,"description":"Object containing the scene fields to update"},{"name":"sceneId","type":"string","required":false,"description":"ID of the scene to update"},{"name":"name","type":"string","required":false,"description":"Name of the scene to update (alternative to sceneId)"},{"name":"active","type":"boolean","required":false,"description":"Set to true to target the active scene"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `delete-scene`

Delete a scene

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sceneId` | string | no | ID of the scene to delete |
| `name` | string | no | Name of the scene to delete (alternative to sceneId) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="delete-scene" parameters={[{"name":"sceneId","type":"string","required":false,"description":"ID of the scene to delete"},{"name":"name","type":"string","required":false,"description":"Name of the scene to delete (alternative to sceneId)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `switch-scene`

Switch the active scene

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sceneId` | string | no | ID of the scene to activate |
| `name` | string | no | Name of the scene to activate (alternative to sceneId) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="switch-scene" parameters={[{"name":"sceneId","type":"string","required":false,"description":"ID of the scene to activate"},{"name":"name","type":"string","required":false,"description":"Name of the scene to activate (alternative to sceneId)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Canvas

### `get-canvas-documents`

Get canvas embedded documents

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `documentType` | string | **yes** | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| `sceneId` | string | no | Scene ID to query (defaults to the active scene) |
| `documentId` | string | no | Specific document ID to retrieve |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-canvas-documents" parameters={[{"name":"documentType","type":"string","required":true,"description":"Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions)"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID to query (defaults to the active scene)"},{"name":"documentId","type":"string","required":false,"description":"Specific document ID to retrieve"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `measure-distance`

Measure the distance between two points or tokens

Calculates the distance between two positions on the canvas, respecting the grid type and measurement rules. Points can be specified as coordinates or by referencing tokens by UUID or name.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `originX` | number | no | Origin x coordinate (optional if originUuid/originName provided) |
| `originY` | number | no | Origin y coordinate |
| `targetX` | number | no | Target x coordinate (optional if targetUuid/targetName provided) |
| `targetY` | number | no | Target y coordinate |
| `originUuid` | string | no | UUID of the origin token |
| `originName` | string | no | Name of the origin token |
| `targetUuid` | string | no | UUID of the target token |
| `targetName` | string | no | Name of the target token |
| `sceneId` | string | no | Scene ID (defaults to active scene) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="measure-distance" parameters={[{"name":"originX","type":"number","required":false,"description":"Origin x coordinate (optional if originUuid/originName provided)"},{"name":"originY","type":"number","required":false,"description":"Origin y coordinate"},{"name":"targetX","type":"number","required":false,"description":"Target x coordinate (optional if targetUuid/targetName provided)"},{"name":"targetY","type":"number","required":false,"description":"Target y coordinate"},{"name":"originUuid","type":"string","required":false,"description":"UUID of the origin token"},{"name":"originName","type":"string","required":false,"description":"Name of the origin token"},{"name":"targetUuid","type":"string","required":false,"description":"UUID of the target token"},{"name":"targetName","type":"string","required":false,"description":"Name of the target token"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID (defaults to active scene)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `create-canvas-document`

Create canvas embedded document(s)

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `documentType` | string | **yes** | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| `data` | object | **yes** | Document data object or array of objects to create |
| `sceneId` | string | no | Scene ID to create in (defaults to the active scene) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="create-canvas-document" parameters={[{"name":"documentType","type":"string","required":true,"description":"Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions)"},{"name":"data","type":"object","required":true,"description":"Document data object or array of objects to create"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID to create in (defaults to the active scene)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `update-canvas-document`

Update a canvas embedded document

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `documentType` | string | **yes** | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| `documentId` | string | **yes** | ID of the document to update |
| `data` | object | **yes** | Object containing the fields to update |
| `sceneId` | string | no | Scene ID containing the document (defaults to the active scene) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="update-canvas-document" parameters={[{"name":"documentType","type":"string","required":true,"description":"Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions)"},{"name":"documentId","type":"string","required":true,"description":"ID of the document to update"},{"name":"data","type":"object","required":true,"description":"Object containing the fields to update"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID containing the document (defaults to the active scene)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `delete-canvas-document`

Delete a canvas embedded document

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `documentType` | string | **yes** | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| `documentId` | string | **yes** | ID of the document to delete |
| `sceneId` | string | no | Scene ID containing the document (defaults to the active scene) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="delete-canvas-document" parameters={[{"name":"documentType","type":"string","required":true,"description":"Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions)"},{"name":"documentId","type":"string","required":true,"description":"ID of the document to delete"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID containing the document (defaults to the active scene)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `move-token`

Move a token to specific coordinates

Moves a token on the canvas to the specified x,y position, optionally animating through waypoints. Token can be identified by UUID or name.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `x` | number | **yes** | Target x coordinate |
| `y` | number | **yes** | Target y coordinate |
| `uuid` | string | no | UUID of the token to move (optional if name provided) |
| `name` | string | no | Name of the token to move (optional if uuid provided) |
| `waypoints` | array | no | Array of waypoint objects with x and y coordinates to animate through before reaching final position |
| `animate` | boolean | no | Whether to animate the movement (default: true) |
| `sceneId` | string | no | Scene ID (defaults to active scene) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="move-token" parameters={[{"name":"x","type":"number","required":true,"description":"Target x coordinate"},{"name":"y","type":"number","required":true,"description":"Target y coordinate"},{"name":"uuid","type":"string","required":false,"description":"UUID of the token to move (optional if name provided)"},{"name":"name","type":"string","required":false,"description":"Name of the token to move (optional if uuid provided)"},{"name":"waypoints","type":"array","required":false,"description":"Array of waypoint objects with x and y coordinates to animate through before reaching final position"},{"name":"animate","type":"boolean","required":false,"description":"Whether to animate the movement (default: true)"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID (defaults to active scene)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Playlist

### `get-playlists`

Get all playlists

Returns all playlists in the world with their tracks/sounds, including playing status, mode, and volume information.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-playlists" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `playlist-play`

Play a playlist or specific sound

Starts playback of an entire playlist or a specific sound within it. The playlist can be identified by ID or name. Optionally specify a specific sound/track to play within the playlist.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `playlistId` | string | no | ID of the playlist (optional if playlistName provided) |
| `playlistName` | string | no | Name of the playlist (optional if playlistId provided) |
| `soundId` | string | no | ID of a specific sound to play within the playlist |
| `soundName` | string | no | Name of a specific sound to play (optional if soundId provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="playlist-play" parameters={[{"name":"playlistId","type":"string","required":false,"description":"ID of the playlist (optional if playlistName provided)"},{"name":"playlistName","type":"string","required":false,"description":"Name of the playlist (optional if playlistId provided)"},{"name":"soundId","type":"string","required":false,"description":"ID of a specific sound to play within the playlist"},{"name":"soundName","type":"string","required":false,"description":"Name of a specific sound to play (optional if soundId provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `playlist-stop`

Stop a playlist

Stops playback of the specified playlist.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `playlistId` | string | no | ID of the playlist (optional if playlistName provided) |
| `playlistName` | string | no | Name of the playlist (optional if playlistId provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="playlist-stop" parameters={[{"name":"playlistId","type":"string","required":false,"description":"ID of the playlist (optional if playlistName provided)"},{"name":"playlistName","type":"string","required":false,"description":"Name of the playlist (optional if playlistId provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `playlist-next`

Skip to next track in a playlist

Advances to the next sound/track in the specified playlist.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `playlistId` | string | no | ID of the playlist (optional if playlistName provided) |
| `playlistName` | string | no | Name of the playlist (optional if playlistId provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="playlist-next" parameters={[{"name":"playlistId","type":"string","required":false,"description":"ID of the playlist (optional if playlistName provided)"},{"name":"playlistName","type":"string","required":false,"description":"Name of the playlist (optional if playlistId provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `playlist-volume`

Set volume for a playlist or specific sound

Adjusts the volume of an entire playlist or a specific sound within it. Volume is specified as a float between 0 (silent) and 1 (full volume).

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `volume` | number | **yes** | Volume level from 0.0 (silent) to 1.0 (full volume) |
| `playlistId` | string | no | ID of the playlist (optional if playlistName provided) |
| `playlistName` | string | no | Name of the playlist (optional if playlistId provided) |
| `soundId` | string | no | ID of a specific sound to adjust volume for |
| `soundName` | string | no | Name of a specific sound to adjust volume for (optional if soundId provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="playlist-volume" parameters={[{"name":"volume","type":"number","required":true,"description":"Volume level from 0.0 (silent) to 1.0 (full volume)"},{"name":"playlistId","type":"string","required":false,"description":"ID of the playlist (optional if playlistName provided)"},{"name":"playlistName","type":"string","required":false,"description":"Name of the playlist (optional if playlistId provided)"},{"name":"soundId","type":"string","required":false,"description":"ID of a specific sound to adjust volume for"},{"name":"soundName","type":"string","required":false,"description":"Name of a specific sound to adjust volume for (optional if soundId provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `stop-sound`

Play a one-shot sound effect

Triggers playback of an audio file by its path. Useful for sound effects, ambient sounds, or any audio that should play once without being part of a playlist. Stop a playing sound Stops playback of a currently playing sound by its source path. If no src is provided, stops all currently playing sounds.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `src` | string | no | Path to the audio file to stop (omit to stop all sounds) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="stop-sound" parameters={[{"name":"src","type":"string","required":false,"description":"Path to the audio file to stop (omit to stop all sounds)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Macro

### `macros`

Get all macros

Retrieves a list of all macros available in the Foundry world.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="macros" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `macro-execute`

Execute a macro by UUID

Executes a specific macro in the Foundry world by its UUID.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `uuid` | string | **yes** | UUID of the macro to execute |
| `args` | object | no | Optional arguments to pass to the macro execution |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="macro-execute" parameters={[{"name":"uuid","type":"string","required":true,"description":"UUID of the macro to execute"},{"name":"args","type":"object","required":false,"description":"Optional arguments to pass to the macro execution"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Utility

### `select`

Select token(s)

Selects one or more tokens in the Foundry VTT client.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `uuids` | array | no | Array of UUIDs to select |
| `name` | string | no | Name of the token(s) to select |
| `data` | object | no | Data to match for selection (e.g., "data.attributes.hp.value": 20) |
| `overwrite` | boolean | no | Whether to overwrite existing selection |
| `all` | boolean | no | Whether to select all tokens on the canvas |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="select" parameters={[{"name":"uuids","type":"array","required":false,"description":"Array of UUIDs to select"},{"name":"name","type":"string","required":false,"description":"Name of the token(s) to select"},{"name":"data","type":"object","required":false,"description":"Data to match for selection (e.g., \"data.attributes.hp.value\": 20)"},{"name":"overwrite","type":"boolean","required":false,"description":"Whether to overwrite existing selection"},{"name":"all","type":"boolean","required":false,"description":"Whether to select all tokens on the canvas"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `selected`

Get selected token(s)

Retrieves the currently selected token(s) in the Foundry VTT client.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="selected" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `players`

Get players/users

Retrieves a list of all users configured in the Foundry VTT world. Useful for discovering valid userId values for permission-scoped API calls.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="players" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `world-info`

Get comprehensive world information

Returns a single object with world name, game system, Foundry version, all modules (with active status), all users (with online status), and the active scene. Useful for API clients to discover the world state.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="world-info" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `play-sound`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `src` | string | **yes** | Path to the audio file (e.g., "sounds/effect.mp3") |
| `volume` | number | no | Volume from 0.0 to 1.0 (default: 0.5) |
| `loop` | boolean | no | Whether to loop the sound (default: false) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="play-sound" parameters={[{"name":"src","type":"string","required":true,"description":"Path to the audio file (e.g., \"sounds/effect.mp3\")"},{"name":"volume","type":"number","required":false,"description":"Volume from 0.0 to 1.0 (default: 0.5)"},{"name":"loop","type":"boolean","required":false,"description":"Whether to loop the sound (default: false)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `execute-js`

Execute JavaScript

Executes a JavaScript script in the Foundry VTT client.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `script` | string | no | JavaScript script to execute |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="execute-js" parameters={[{"name":"script","type":"string","required":false,"description":"JavaScript script to execute"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## User

### `get-users`

List all Foundry users

Retrieves a list of all users configured in the Foundry VTT world, including their roles and online status. This is a GM-only operation.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-users" parameters={[{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `get-user`

Get a single Foundry user

Retrieves a single user by their ID or name. This is a GM-only operation.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | no | ID of the user to retrieve |
| `name` | string | no | Name of the user to retrieve (alternative to id) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-user" parameters={[{"name":"id","type":"string","required":false,"description":"ID of the user to retrieve"},{"name":"name","type":"string","required":false,"description":"Name of the user to retrieve (alternative to id)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `create-user`

Create a new Foundry user

Creates a new user in the Foundry VTT world with the specified name, role, and optional password. This is a GM-only operation.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | **yes** | Username for the new user |
| `role` | number | no | User role: 0=None, 1=Player, 2=Trusted, 3=Assistant, 4=GM (default: 1) |
| `password` | string | no | Password for the new user |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="create-user" parameters={[{"name":"name","type":"string","required":true,"description":"Username for the new user"},{"name":"role","type":"number","required":false,"description":"User role: 0=None, 1=Player, 2=Trusted, 3=Assistant, 4=GM (default: 1)"},{"name":"password","type":"string","required":false,"description":"Password for the new user"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `update-user`

Update an existing Foundry user

Updates fields on an existing user. Identify the user by id or name, then pass the fields to update in the data object. Cannot demote the last GM user. This is a GM-only operation.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `data` | object | **yes** | Object containing user fields to update (name, role, password, color, avatar, character) |
| `id` | string | no | ID of the user to update |
| `name` | string | no | Name of the user to update (alternative to id) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="update-user" parameters={[{"name":"data","type":"object","required":true,"description":"Object containing user fields to update (name, role, password, color, avatar, character)"},{"name":"id","type":"string","required":false,"description":"ID of the user to update"},{"name":"name","type":"string","required":false,"description":"Name of the user to update (alternative to id)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `delete-user`

Delete a Foundry user

Permanently deletes a user from the Foundry VTT world. Cannot delete yourself or the last GM user. This is a GM-only operation.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | no | ID of the user to delete |
| `name` | string | no | Name of the user to delete (alternative to id) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="delete-user" parameters={[{"name":"id","type":"string","required":false,"description":"ID of the user to delete"},{"name":"name","type":"string","required":false,"description":"Name of the user to delete (alternative to id)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Search

### `search`

Search entities

This endpoint allows searching for entities in the Foundry world based on a query string. Search world entities and compendiums using the native built-in search engine. No third-party modules required. Results are ranked by relevance: exact match, prefix match, substring match, word-prefix match, and subsequence match.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | no | Search query string (omit to browse all entities matching filter) |
| `filter` | string | no | Filter string — simple: filter="Actor"; compound: filter="documentType:Item,subType:weapon". Supported keys: documentType, subType, folder, package, resultType |
| `excludeCompendiums` | boolean | no | Exclude compendium entries from results (default: false — compendiums are included by default) |
| `limit` | number | no | Maximum number of results to return (default: 200, max: 500) |
| `minified` | boolean | no | Return minimal fields only — uuid, id, name, img, documentType (default: false) |
| `ownedByUserId` | string | no | Filter results to only documents the specified Foundry user (ID or username) has Owner permission on |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="search" parameters={[{"name":"query","type":"string","required":false,"description":"Search query string (omit to browse all entities matching filter)"},{"name":"filter","type":"string","required":false,"description":"Filter string — simple: filter=\"Actor\"; compound: filter=\"documentType:Item,subType:weapon\". Supported keys: documentType, subType, folder, package, resultType"},{"name":"excludeCompendiums","type":"boolean","required":false,"description":"Exclude compendium entries from results (default: false — compendiums are included by default)"},{"name":"limit","type":"number","required":false,"description":"Maximum number of results to return (default: 200, max: 500)"},{"name":"minified","type":"boolean","required":false,"description":"Return minimal fields only — uuid, id, name, img, documentType (default: false)"},{"name":"ownedByUserId","type":"string","required":false,"description":"Filter results to only documents the specified Foundry user (ID or username) has Owner permission on"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## FileSystem

### `file-system`

Get file system structure

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `path` | string | no | The path to retrieve (relative to source) |
| `source` | string | no | The source directory to use (data, systems, modules, etc.) |
| `recursive` | boolean | no | Whether to recursively list all subdirectories |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="file-system" parameters={[{"name":"path","type":"string","required":false,"description":"The path to retrieve (relative to source)"},{"name":"source","type":"string","required":false,"description":"The source directory to use (data, systems, modules, etc.)"},{"name":"recursive","type":"boolean","required":false,"description":"Whether to recursively list all subdirectories"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `download-file`

Download a file from Foundry's file system

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `path` | string | no | The full path to the file to download |
| `source` | string | no | The source directory to use (data, systems, modules, etc.) |
| `format` | string | no | The format to return the file in (binary, base64) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="download-file" parameters={[{"name":"path","type":"string","required":false,"description":"The full path to the file to download"},{"name":"source","type":"string","required":false,"description":"The source directory to use (data, systems, modules, etc.)"},{"name":"format","type":"string","required":false,"description":"The format to return the file in (binary, base64)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Dnd5e

### `get-actor-details`

Get detailed information for a specific D&D 5e actor

Retrieves comprehensive details about an actor including stats, inventory, spells, features, and other character information based on the requested details array.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `details` | array | **yes** | Array of detail types to retrieve (e.g., ["resources", "items", "spells", "features"]) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-actor-details" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"details","type":"array","required":true,"description":"Array of detail types to retrieve (e.g., [\"resources\", \"items\", \"spells\", \"features\"])"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `modify-item-charges`

Modify the charges for a specific item owned by an actor

Increases or decreases the charges/uses of an item in an actor's inventory. Useful for consumable items like potions, scrolls, or charged magic items.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `amount` | number | **yes** | The amount to modify charges by (positive or negative) |
| `itemUuid` | string | no | The UUID of the specific item (optional if itemName provided) |
| `itemName` | string | no | The name of the item if UUID not provided (optional if itemUuid provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="modify-item-charges" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"amount","type":"number","required":true,"description":"The amount to modify charges by (positive or negative)"},{"name":"itemUuid","type":"string","required":false,"description":"The UUID of the specific item (optional if itemName provided)"},{"name":"itemName","type":"string","required":false,"description":"The name of the item if UUID not provided (optional if itemUuid provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `short-rest`

Perform a short rest for an actor

Triggers the D&D 5e short rest workflow including hit dice recovery, class feature resets, and HP recovery.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | no | UUID of the actor (optional if selected is true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `autoHD` | boolean | no | Automatically spend hit dice during short rest |
| `autoHDThreshold` | number | no | HP threshold below which to auto-spend hit dice (0-1 as fraction of max HP) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="short-rest" parameters={[{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor (optional if selected is true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"autoHD","type":"boolean","required":false,"description":"Automatically spend hit dice during short rest"},{"name":"autoHDThreshold","type":"number","required":false,"description":"HP threshold below which to auto-spend hit dice (0-1 as fraction of max HP)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `long-rest`

Perform a long rest for an actor

Triggers the D&D 5e long rest workflow including full HP recovery, spell slot restoration, hit dice recovery, and feature resets.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | no | UUID of the actor (optional if selected is true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `newDay` | boolean | no | Whether the long rest marks a new day (default: true) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="long-rest" parameters={[{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor (optional if selected is true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"newDay","type":"boolean","required":false,"description":"Whether the long rest marks a new day (default: true)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `skill-check`

Roll a skill check for an actor

Rolls a D&D 5e skill check with all applicable modifiers including proficiency, expertise, Jack of All Trades, and conditional bonuses.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `skill` | string | **yes** | Skill abbreviation (e.g., "acr", "ath", "ste", "prc") |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |
| `advantage` | boolean | no | Roll with advantage |
| `disadvantage` | boolean | no | Roll with disadvantage |
| `bonus` | string | no | Extra bonus formula to add (e.g., "1d4", "+2") |
| `createChatMessage` | boolean | no | Whether to post the roll to chat (default: true) |

<WsMessageTester messageType="skill-check" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"skill","type":"string","required":true,"description":"Skill abbreviation (e.g., \"acr\", \"ath\", \"ste\", \"prc\")"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"},{"name":"advantage","type":"boolean","required":false,"description":"Roll with advantage"},{"name":"disadvantage","type":"boolean","required":false,"description":"Roll with disadvantage"},{"name":"bonus","type":"string","required":false,"description":"Extra bonus formula to add (e.g., \"1d4\", \"+2\")"},{"name":"createChatMessage","type":"boolean","required":false,"description":"Whether to post the roll to chat (default: true)"}]} />

---

### `ability-save`

Roll an ability saving throw for an actor

Rolls a D&D 5e ability saving throw with all applicable modifiers.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `ability` | string | **yes** | Ability abbreviation (e.g., "str", "dex", "con", "int", "wis", "cha") |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |
| `advantage` | boolean | no | Roll with advantage |
| `disadvantage` | boolean | no | Roll with disadvantage |
| `bonus` | string | no | Extra bonus formula to add (e.g., "1d4", "+2") |
| `createChatMessage` | boolean | no | Whether to post the roll to chat (default: true) |

<WsMessageTester messageType="ability-save" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"ability","type":"string","required":true,"description":"Ability abbreviation (e.g., \"str\", \"dex\", \"con\", \"int\", \"wis\", \"cha\")"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"},{"name":"advantage","type":"boolean","required":false,"description":"Roll with advantage"},{"name":"disadvantage","type":"boolean","required":false,"description":"Roll with disadvantage"},{"name":"bonus","type":"string","required":false,"description":"Extra bonus formula to add (e.g., \"1d4\", \"+2\")"},{"name":"createChatMessage","type":"boolean","required":false,"description":"Whether to post the roll to chat (default: true)"}]} />

---

### `ability-check`

Roll an ability check for an actor

Rolls a D&D 5e ability check (raw ability test, not a skill check) with all applicable modifiers.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `ability` | string | **yes** | Ability abbreviation (e.g., "str", "dex", "con", "int", "wis", "cha") |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |
| `advantage` | boolean | no | Roll with advantage |
| `disadvantage` | boolean | no | Roll with disadvantage |
| `bonus` | string | no | Extra bonus formula to add (e.g., "1d4", "+2") |
| `createChatMessage` | boolean | no | Whether to post the roll to chat (default: true) |

<WsMessageTester messageType="ability-check" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"ability","type":"string","required":true,"description":"Ability abbreviation (e.g., \"str\", \"dex\", \"con\", \"int\", \"wis\", \"cha\")"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"},{"name":"advantage","type":"boolean","required":false,"description":"Roll with advantage"},{"name":"disadvantage","type":"boolean","required":false,"description":"Roll with disadvantage"},{"name":"bonus","type":"string","required":false,"description":"Extra bonus formula to add (e.g., \"1d4\", \"+2\")"},{"name":"createChatMessage","type":"boolean","required":false,"description":"Whether to post the roll to chat (default: true)"}]} />

---

### `death-save`

Roll a death saving throw for an actor

Rolls a D&D 5e death saving throw, handling DC 10 CON save, three successes/failures tracking, nat 20 healing, and nat 1 double failure.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `advantage` | boolean | no | Roll with advantage |
| `createChatMessage` | boolean | no | Whether to post the roll to chat (default: true) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="death-save" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"advantage","type":"boolean","required":false,"description":"Roll with advantage"},{"name":"createChatMessage","type":"boolean","required":false,"description":"Whether to post the roll to chat (default: true)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `modify-experience`

Modify the experience points for a specific actor

Adds or removes experience points from an actor.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `amount` | number | **yes** | The amount of experience to add (can be negative) |
| `actorUuid` | string | no | UUID of the actor (optional if selected is true) |
| `selected` | boolean | no | Whether to get the selected entity |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="modify-experience" parameters={[{"name":"amount","type":"number","required":true,"description":"The amount of experience to add (can be negative)"},{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor (optional if selected is true)"},{"name":"selected","type":"boolean","required":false,"description":"Whether to get the selected entity"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `get-concentration`

Check if an actor is concentrating on a spell

Returns whether the actor currently has a concentration effect active, and if so, what spell they are concentrating on.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | no | UUID of the actor (optional if selected is true) |
| `actorName` | string | no | Name of the actor (optional if actorUuid provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-concentration" parameters={[{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor (optional if selected is true)"},{"name":"actorName","type":"string","required":false,"description":"Name of the actor (optional if actorUuid provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `break-concentration`

Break an actor's concentration

Removes the concentration effect from the actor, ending any spell that requires concentration.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | no | UUID of the actor (optional if selected is true) |
| `actorName` | string | no | Name of the actor (optional if actorUuid provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="break-concentration" parameters={[{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor (optional if selected is true)"},{"name":"actorName","type":"string","required":false,"description":"Name of the actor (optional if actorUuid provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `concentration-save`

Roll a concentration saving throw

Rolls a Constitution saving throw to maintain concentration after taking damage. The DC is calculated as max(10, floor(damage/2)). Returns the roll result and whether concentration was maintained or broken.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `damage` | number | **yes** | Amount of damage taken (used to calculate DC = max(10, floor(damage/2))) |
| `actorUuid` | string | no | UUID of the actor (optional if selected is true) |
| `actorName` | string | no | Name of the actor (optional if actorUuid provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |
| `advantage` | boolean | no | Roll with advantage |
| `disadvantage` | boolean | no | Roll with disadvantage |
| `bonus` | string | no | Extra bonus formula to add (e.g., "1d4", "+2") |
| `createChatMessage` | boolean | no | Whether to post the roll to chat (default: true) |

<WsMessageTester messageType="concentration-save" parameters={[{"name":"damage","type":"number","required":true,"description":"Amount of damage taken (used to calculate DC = max(10, floor(damage/2)))"},{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor (optional if selected is true)"},{"name":"actorName","type":"string","required":false,"description":"Name of the actor (optional if actorUuid provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"},{"name":"advantage","type":"boolean","required":false,"description":"Roll with advantage"},{"name":"disadvantage","type":"boolean","required":false,"description":"Roll with disadvantage"},{"name":"bonus","type":"string","required":false,"description":"Extra bonus formula to add (e.g., \"1d4\", \"+2\")"},{"name":"createChatMessage","type":"boolean","required":false,"description":"Whether to post the roll to chat (default: true)"}]} />

---

### `equip-item`

Equip or unequip an item

Changes the equipped status of an item in an actor's inventory.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `equipped` | boolean | **yes** | Whether the item should be equipped (true) or unequipped (false) |
| `actorUuid` | string | no | UUID of the actor (optional if selected is true) |
| `actorName` | string | no | Name of the actor (optional if actorUuid provided) |
| `itemUuid` | string | no | UUID of the item (optional if itemName provided) |
| `itemName` | string | no | Name of the item (optional if itemUuid provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="equip-item" parameters={[{"name":"equipped","type":"boolean","required":true,"description":"Whether the item should be equipped (true) or unequipped (false)"},{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor (optional if selected is true)"},{"name":"actorName","type":"string","required":false,"description":"Name of the actor (optional if actorUuid provided)"},{"name":"itemUuid","type":"string","required":false,"description":"UUID of the item (optional if itemName provided)"},{"name":"itemName","type":"string","required":false,"description":"Name of the item (optional if itemUuid provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `attune-item`

Attune or unattune an item

Changes the attunement status of a magic item in an actor's inventory.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `attuned` | boolean | **yes** | Whether the item should be attuned (true) or unattuned (false) |
| `actorUuid` | string | no | UUID of the actor (optional if selected is true) |
| `actorName` | string | no | Name of the actor (optional if actorUuid provided) |
| `itemUuid` | string | no | UUID of the item (optional if itemName provided) |
| `itemName` | string | no | Name of the item (optional if itemUuid provided) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="attune-item" parameters={[{"name":"attuned","type":"boolean","required":true,"description":"Whether the item should be attuned (true) or unattuned (false)"},{"name":"actorUuid","type":"string","required":false,"description":"UUID of the actor (optional if selected is true)"},{"name":"actorName","type":"string","required":false,"description":"Name of the actor (optional if actorUuid provided)"},{"name":"itemUuid","type":"string","required":false,"description":"UUID of the item (optional if itemName provided)"},{"name":"itemName","type":"string","required":false,"description":"Name of the item (optional if itemUuid provided)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `transfer-currency`

Transfer currency between actors

Moves currency from one actor to another. Validates that the source actor has sufficient funds before transferring.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `currency` | object | **yes** | Currency amounts to transfer, e.g. pp, gp, ep, sp, cp denomination keys with numeric values |
| `sourceActorUuid` | string | no | UUID of the source actor (optional if sourceActorName provided) |
| `sourceActorName` | string | no | Name of the source actor |
| `targetActorUuid` | string | no | UUID of the target actor (optional if targetActorName provided) |
| `targetActorName` | string | no | Name of the target actor |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="transfer-currency" parameters={[{"name":"currency","type":"object","required":true,"description":"Currency amounts to transfer, e.g. pp, gp, ep, sp, cp denomination keys with numeric values"},{"name":"sourceActorUuid","type":"string","required":false,"description":"UUID of the source actor (optional if sourceActorName provided)"},{"name":"sourceActorName","type":"string","required":false,"description":"Name of the source actor"},{"name":"targetActorUuid","type":"string","required":false,"description":"UUID of the target actor (optional if targetActorName provided)"},{"name":"targetActorName","type":"string","required":false,"description":"Name of the target actor"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `modify-currency`

Modify currency balance for a single actor (delta-based, not a transfer between actors)

Adds or removes currency from an actor's wallet. Use a negative amount to remove currency.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `currency` | string | **yes** | Currency denomination to modify (pp, gp, ep, sp, cp) |
| `amount` | number | **yes** | Amount to add (positive) or remove (negative) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="modify-currency" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"currency","type":"string","required":true,"description":"Currency denomination to modify (pp, gp, ep, sp, cp)"},{"name":"amount","type":"number","required":true,"description":"Amount to add (positive) or remove (negative)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `prepare-spell`

Prepare or unprepare a spell for an actor

Toggles a spell's prepared state. Only applicable to spellcaster classes that prepare spells.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `spellName` | string | **yes** | Name of the spell to prepare or unprepare |
| `prepared` | boolean | **yes** | True to prepare the spell, false to unprepare it |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="prepare-spell" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"spellName","type":"string","required":true,"description":"Name of the spell to prepare or unprepare"},{"name":"prepared","type":"boolean","required":true,"description":"True to prepare the spell, false to unprepare it"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `use-ability`

Use an ability

Activates a specific ability for an actor, optionally targeting another entity

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `abilityUuid` | string | no | The UUID of the specific ability (optional if abilityName provided) |
| `abilityName` | string | no | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| `targetUuid` | string | no | The UUID of the target for the ability (optional) |
| `targetName` | string | no | The name of the target if UUID not provided (optional) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="use-ability" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"abilityUuid","type":"string","required":false,"description":"The UUID of the specific ability (optional if abilityName provided)"},{"name":"abilityName","type":"string","required":false,"description":"The name of the ability if UUID not provided (optional if abilityUuid provided)"},{"name":"targetUuid","type":"string","required":false,"description":"The UUID of the target for the ability (optional)"},{"name":"targetName","type":"string","required":false,"description":"The name of the target if UUID not provided (optional)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `use-feature`

Use a feature

Activates a specific feature for an actor, optionally targeting another entity

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `abilityUuid` | string | no | The UUID of the specific ability (optional if abilityName provided) |
| `abilityName` | string | no | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| `targetUuid` | string | no | The UUID of the target for the ability (optional) |
| `targetName` | string | no | The name of the target if UUID not provided (optional) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="use-feature" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"abilityUuid","type":"string","required":false,"description":"The UUID of the specific ability (optional if abilityName provided)"},{"name":"abilityName","type":"string","required":false,"description":"The name of the ability if UUID not provided (optional if abilityUuid provided)"},{"name":"targetUuid","type":"string","required":false,"description":"The UUID of the target for the ability (optional)"},{"name":"targetName","type":"string","required":false,"description":"The name of the target if UUID not provided (optional)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `use-spell`

Use a spell

Casts a specific spell for an actor, optionally targeting another entity

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `abilityUuid` | string | no | The UUID of the specific ability (optional if abilityName provided) |
| `abilityName` | string | no | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| `targetUuid` | string | no | The UUID of the target for the ability (optional) |
| `targetName` | string | no | The name of the target if UUID not provided (optional) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="use-spell" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"abilityUuid","type":"string","required":false,"description":"The UUID of the specific ability (optional if abilityName provided)"},{"name":"abilityName","type":"string","required":false,"description":"The name of the ability if UUID not provided (optional if abilityUuid provided)"},{"name":"targetUuid","type":"string","required":false,"description":"The UUID of the target for the ability (optional)"},{"name":"targetName","type":"string","required":false,"description":"The name of the target if UUID not provided (optional)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `use-item`

Use an item

Uses a specific item for an actor, optionally targeting another entity

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `actorUuid` | string | **yes** | UUID of the actor |
| `abilityUuid` | string | no | The UUID of the specific ability (optional if abilityName provided) |
| `abilityName` | string | no | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| `targetUuid` | string | no | The UUID of the target for the ability (optional) |
| `targetName` | string | no | The name of the target if UUID not provided (optional) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="use-item" parameters={[{"name":"actorUuid","type":"string","required":true,"description":"UUID of the actor"},{"name":"abilityUuid","type":"string","required":false,"description":"The UUID of the specific ability (optional if abilityName provided)"},{"name":"abilityName","type":"string","required":false,"description":"The name of the ability if UUID not provided (optional if abilityUuid provided)"},{"name":"targetUuid","type":"string","required":false,"description":"The UUID of the target for the ability (optional)"},{"name":"targetName","type":"string","required":false,"description":"The name of the target if UUID not provided (optional)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

## Cross-World Operations

### `remote-request`

Invoke any supported action on a remote Foundry world via the relay tunnel

:::note
Foundry module connection token required; API key clients are rejected
:::

The source connection token must list the target in allowedTargetClients and hold the required scope in remoteScopes. Configure these in the dashboard → Connections → Edit browser.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `targetClientId` | string | **yes** | Client ID of the target world (must be in allowedTargetClients) |
| `action` | string | **yes** | Action to invoke on the target (e.g. create-user, entity, search) |
| `payload` | object | no | Action payload forwarded verbatim to the target module |
| `autoStartIfOffline` | boolean | no | Start a headless session if the target is offline (requires autoStartOnRemoteRequest enabled for that world) |

<WsMessageTester messageType="remote-request" resultType="remote-response" parameters={[{"name":"targetClientId","type":"string","required":true,"description":"Client ID of the target world (must be in allowedTargetClients)"},{"name":"action","type":"string","required":true,"description":"Action to invoke on the target (e.g. create-user, entity, search)"},{"name":"payload","type":"object","required":false,"description":"Action payload forwarded verbatim to the target module"},{"name":"autoStartIfOffline","type":"boolean","required":false,"description":"Start a headless session if the target is offline (requires autoStartOnRemoteRequest enabled for that world)"}]} />

---

## AsyncAPI Spec

The full AsyncAPI specification is available at `/asyncapi.json`.
