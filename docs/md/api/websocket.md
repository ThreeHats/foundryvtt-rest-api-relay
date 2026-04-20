---
tag: WebSocket
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


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
| `documentType` | string | **yes** | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls) |
| `sceneId` | string | no | Scene ID to query (defaults to the active scene) |
| `documentId` | string | no | Specific document ID to retrieve |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="get-canvas-documents" parameters={[{"name":"documentType","type":"string","required":true,"description":"Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID to query (defaults to the active scene)"},{"name":"documentId","type":"string","required":false,"description":"Specific document ID to retrieve"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

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
| `documentType` | string | **yes** | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls) |
| `data` | object | **yes** | Document data object or array of objects to create |
| `sceneId` | string | no | Scene ID to create in (defaults to the active scene) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="create-canvas-document" parameters={[{"name":"documentType","type":"string","required":true,"description":"Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)"},{"name":"data","type":"object","required":true,"description":"Document data object or array of objects to create"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID to create in (defaults to the active scene)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `update-canvas-document`

Update a canvas embedded document

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `documentType` | string | **yes** | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls) |
| `documentId` | string | **yes** | ID of the document to update |
| `data` | object | **yes** | Object containing the fields to update |
| `sceneId` | string | no | Scene ID containing the document (defaults to the active scene) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="update-canvas-document" parameters={[{"name":"documentType","type":"string","required":true,"description":"Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)"},{"name":"documentId","type":"string","required":true,"description":"ID of the document to update"},{"name":"data","type":"object","required":true,"description":"Object containing the fields to update"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID containing the document (defaults to the active scene)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

---

### `delete-canvas-document`

Delete a canvas embedded document

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `documentType` | string | **yes** | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls) |
| `documentId` | string | **yes** | ID of the document to delete |
| `sceneId` | string | no | Scene ID containing the document (defaults to the active scene) |
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="delete-canvas-document" parameters={[{"name":"documentType","type":"string","required":true,"description":"Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)"},{"name":"documentId","type":"string","required":true,"description":"ID of the document to delete"},{"name":"sceneId","type":"string","required":false,"description":"Scene ID containing the document (defaults to the active scene)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

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
| `userId` | string | no | Foundry user ID or username to scope permissions (omit for GM-level access) |

<WsMessageTester messageType="search" parameters={[{"name":"query","type":"string","required":false,"description":"Search query string (omit to browse all entities matching filter)"},{"name":"filter","type":"string","required":false,"description":"Filter string — simple: filter=\"Actor\"; compound: filter=\"documentType:Item,subType:weapon\". Supported keys: documentType, subType, folder, package, resultType"},{"name":"excludeCompendiums","type":"boolean","required":false,"description":"Exclude compendium entries from results (default: false — compendiums are included by default)"},{"name":"limit","type":"number","required":false,"description":"Maximum number of results to return (default: 200, max: 500)"},{"name":"minified","type":"boolean","required":false,"description":"Return minimal fields only — uuid, id, name, img, documentType (default: false)"},{"name":"userId","type":"string","required":false,"description":"Foundry user ID or username to scope permissions (omit for GM-level access)"}]} />

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

## Code Examples

### search

Search for entities

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const ws = new WebSocket('ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID');

ws.onopen = () => {
  // Send auth message first — token must not be in the URL
  ws.send(JSON.stringify({ type: 'auth', token: 'YOUR_API_KEY' }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  if (data.type === 'connected') {
    // Now send your request
    ws.send(JSON.stringify({
      "query": "test",
      "type": "search",
      "requestId": "unique-id"
    }));
  }
  if (data.type === 'search-result') {
    console.log(data);
  }
};
```

</TabItem>
<TabItem value="python" label="Python">

```python
import asyncio
import websockets
import json

async def main():
    uri = 'ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID'
    async with websockets.connect(uri) as ws:
        # Send auth message first — token must not be in the URL
        await ws.send(json.dumps({'type': 'auth', 'token': 'YOUR_API_KEY'}))
        connected = json.loads(await ws.recv())
        if connected.get('type') != 'connected':
            raise Exception('Auth failed')
        await ws.send(json.dumps({"query":"test","type":"search","requestId":"unique-id"}))
        response = await ws.recv()
        data = json.loads(response)
        print(data)

asyncio.run(main())
```

</TabItem>
<TabItem value="typescript" label="TypeScript">

```typescript
import WebSocket from 'ws';

const ws = new WebSocket('ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID');

ws.on('open', () => {
  // Send auth message first — token must not be in the URL
  ws.send(JSON.stringify({ type: 'auth', token: 'YOUR_API_KEY' }));
});

ws.on('message', (raw: string) => {
  const data = JSON.parse(raw);
  if (data.type === 'connected') {
    // Now send your request
    ws.send(JSON.stringify({
      "query": "test",
      "type": "search",
      "requestId": "unique-id"
    }));
  }
  if (data.type === 'search-result') {
    console.log(data);
  }
}):
```

</TabItem>
</Tabs>

#### Response

```json
{
  "clientId": "fvtt_099ad17ea199e7e3",
  "query": "test",
  "requestId": "test_1776657997410_tmbzni",
  "results": [
    {
      "documentType": "Scene",
      "folder": null,
      "formattedMatch": "<strong>test</strong>",
      "icon": "worlds/testing/assets/scenes/fWuYVDZAnXQiKXM5-thumb.webp",
      "id": "fWuYVDZAnXQiKXM5",
      "journalLink": "@UUID[Scene.fWuYVDZAnXQiKXM5]{test}",
      "name": "test",
      "package": null,
      "packageName": null,
      "resultType": "WorldEntity",
      "subType": "",
      "tagline": "Scenes Directory",
      "uuid": "Scene.fWuYVDZAnXQiKXM5"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "<strong>test</strong>-journalentry",
      "icon": "",
      "id": "t4YBhtHLJvfzUxLy",
      "journalLink": "@UUID[JournalEntry.t4YBhtHLJvfzUxLy]{test-journalentry}",
      "name": "test-journalentry",
      "package": null,
      "packageName": null,
      "resultType": "WorldEntity",
      "subType": "",
      "tagline": "Journal Directory",
      "uuid": "JournalEntry.t4YBhtHLJvfzUxLy"
    },
    {
      "documentType": "Macro",
      "folder": null,
      "formattedMatch": "<strong>test</strong>-macro",
      "icon": "icons/svg/dice-target.svg",
      "id": "AkcLmoRwvkrPvjyA",
      "journalLink": "@UUID[Macro.AkcLmoRwvkrPvjyA]{test-macro}",
      "name": "test-macro",
      "package": null,
      "packageName": null,
      "resultType": "WorldEntity",
      "subType": "script",
      "tagline": "Macros Directory",
      "uuid": "Macro.AkcLmoRwvkrPvjyA"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "formattedMatch": "<strong>test</strong>-perrin (halfling monk)",
      "icon": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
      "id": "w5STPCwE3YTDztRk",
      "journalLink": "@UUID[Actor.w5STPCwE3YTDztRk]{test-perrin (halfling monk)}",
      "name": "test-perrin (halfling monk)",
      "package": null,
      "packageName": null,
      "resultType": "WorldEntity",
      "subType": "character",
      "tagline": "Actors Directory",
      "uuid": "Actor.w5STPCwE3YTDztRk"
    },
    {
      "documentType": "Scene",
      "folder": null,
      "formattedMatch": "<strong>test</strong>-scene-updated",
      "icon": "",
      "id": "r36nfimJGHYGUGQX",
      "journalLink": "@UUID[Scene.r36nfimJGHYGUGQX]{test-scene-updated}",
      "name": "test-scene-updated",
      "package": null,
      "packageName": null,
      "resultType": "WorldEntity",
      "subType": "",
      "tagline": "Scenes Directory",
      "uuid": "Scene.r36nfimJGHYGUGQX"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "<strong>test</strong>-studded leather armor +3",
      "icon": "icons/equipment/chest/breastplate-rivited-red.webp",
      "id": "0rkF7xi8VrKkFvzZ",
      "journalLink": "@UUID[Item.0rkF7xi8VrKkFvzZ]{test-studded leather armor +3}",
      "name": "test-studded leather armor +3",
      "package": null,
      "packageName": null,
      "resultType": "WorldEntity",
      "subType": "equipment",
      "tagline": "Items Directory",
      "uuid": "Item.0rkF7xi8VrKkFvzZ"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "D20 <strong>Test</strong>s",
      "icon": "",
      "id": "phbD20Tests00000",
      "journalLink": "@UUID[Compendium.dnd5e.content24.JournalEntry.phbD20Tests00000]{D20 Tests}",
      "name": "D20 Tests",
      "package": "dnd5e.content24",
      "packageName": "Rules",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules",
      "uuid": "Compendium.dnd5e.content24.JournalEntry.phbD20Tests00000"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "formattedMatch": "Updated <strong>Test</strong> Actor",
      "icon": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
      "id": "q9uWyfdPwTlzbpxb",
      "journalLink": "@UUID[Actor.q9uWyfdPwTlzbpxb]{Updated Test Actor}",
      "name": "Updated Test Actor",
      "package": null,
      "packageName": null,
      "resultType": "WorldEntity",
      "subType": "character",
      "tagline": "Actors Directory",
      "uuid": "Actor.q9uWyfdPwTlzbpxb"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Adamantine Breastplate",
      "icon": "icons/equipment/chest/breastplate-collared-steel-grey.webp",
      "id": "DevmObXWP9MfwE2c",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.DevmObXWP9MfwE2c]{Adamantine Breastplate}",
      "name": "Adamantine Breastplate",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.DevmObXWP9MfwE2c"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Adamantine Chain Shirt",
      "icon": "icons/equipment/chest/breastplate-scale-grey.webp",
      "id": "kjTPoUeomTPWJ9h3",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.kjTPoUeomTPWJ9h3]{Adamantine Chain Shirt}",
      "name": "Adamantine Chain Shirt",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.kjTPoUeomTPWJ9h3"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Adamantine Splint Armor",
      "icon": "icons/equipment/chest/breastplate-layered-steel.webp",
      "id": "LDuqUcosOK8Bf76S",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.LDuqUcosOK8Bf76S]{Adamantine Splint Armor}",
      "name": "Adamantine Splint Armor",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.LDuqUcosOK8Bf76S"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Amulet of the Planes Destination",
      "icon": "icons/equipment/neck/amulet-carved-stone-purple.webp",
      "id": "dmgAmuletOfThePl",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.dmgAmuletOfThePl]{Amulet of the Planes Destination}",
      "name": "Amulet of the Planes Destination",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.dmgAmuletOfThePl"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Animated Object Catalysts",
      "icon": "icons/commodities/treasure/stone-cracked-lightning-blue.webp",
      "id": "mmAnimatedObject",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.mmAnimatedObject]{Animated Object Catalysts}",
      "name": "Animated Object Catalysts",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.mmAnimatedObject"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "formattedMatch": "Animated Rug of Smothering",
      "icon": "systems/dnd5e/tokens/construct/RugOfSmothering.webp",
      "id": "mmAnimatedRugOfS",
      "journalLink": "@UUID[Compendium.dnd5e.actors24.Actor.mmAnimatedRugOfS]{Animated Rug of Smothering}",
      "name": "Animated Rug of Smothering",
      "package": "dnd5e.actors24",
      "packageName": "Actors",
      "resultType": "CompendiumEntity",
      "subType": "npc",
      "tagline": "Actors",
      "uuid": "Compendium.dnd5e.actors24.Actor.mmAnimatedRugOfS"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Belt of Fire Giant Strength",
      "icon": "icons/equipment/waist/belt-coiled-leather-steel.webp",
      "id": "bq9YKwEHLQ7p7ric",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.bq9YKwEHLQ7p7ric]{Belt of Fire Giant Strength}",
      "name": "Belt of Fire Giant Strength",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.bq9YKwEHLQ7p7ric"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Belt of Giant Strength (frost)",
      "icon": "icons/equipment/waist/belt-buckle-gold-blue.webp",
      "id": "dmgfroBeltofGian",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgfroBeltofGian]{Belt of Giant Strength (frost)}",
      "name": "Belt of Giant Strength (frost)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgfroBeltofGian"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Belt of Giant Strength (stone)",
      "icon": "icons/equipment/waist/belt-armored-steel.webp",
      "id": "dmgstoBeltofGian",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgstoBeltofGian]{Belt of Giant Strength (stone)}",
      "name": "Belt of Giant Strength (stone)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgstoBeltofGian"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Belt of Giant Strength (storm)",
      "icon": "icons/equipment/waist/belt-thick-gemmed-gold-blue.webp",
      "id": "dmgstmBeltofGian",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgstmBeltofGian]{Belt of Giant Strength (storm)}",
      "name": "Belt of Giant Strength (storm)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgstmBeltofGian"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Belt of Stone Giant Strength",
      "icon": "icons/equipment/waist/belt-armored-steel.webp",
      "id": "fCUZ7h8YYrs16UhX",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.fCUZ7h8YYrs16UhX]{Belt of Stone Giant Strength}",
      "name": "Belt of Stone Giant Strength",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.fCUZ7h8YYrs16UhX"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Blessed Strikes: Divine Strike",
      "icon": "icons/magic/holy/prayer-hands-glowing-yellow.webp",
      "id": "phbDivineStrike0",
      "journalLink": "@UUID[Compendium.dnd5e.classes24.Item.phbDivineStrike0]{Blessed Strikes: Divine Strike}",
      "name": "Blessed Strikes: Divine Strike",
      "package": "dnd5e.classes24",
      "packageName": "Character Classes",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Classes",
      "uuid": "Compendium.dnd5e.classes24.Item.phbDivineStrike0"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Blessed Strikes: Potent Spellcasting",
      "icon": "icons/magic/life/cross-yellow-green.webp",
      "id": "phbClcPotentSpel",
      "journalLink": "@UUID[Compendium.dnd5e.classes24.Item.phbClcPotentSpel]{Blessed Strikes: Potent Spellcasting}",
      "name": "Blessed Strikes: Potent Spellcasting",
      "package": "dnd5e.classes24",
      "packageName": "Character Classes",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Classes",
      "uuid": "Compendium.dnd5e.classes24.Item.phbClcPotentSpel"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Boon of the Night Spirit",
      "icon": "icons/magic/unholy/silhouette-robe-evil-power.webp",
      "id": "phbBoonoftheNigh",
      "journalLink": "@UUID[Compendium.dnd5e.feats24.Item.phbBoonoftheNigh]{Boon of the Night Spirit}",
      "name": "Boon of the Night Spirit",
      "package": "dnd5e.feats24",
      "packageName": "Feats",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Feats",
      "uuid": "Compendium.dnd5e.feats24.Item.phbBoonoftheNigh"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Boon of Truesight",
      "icon": "icons/creatures/eyes/humanoid-single-blind.webp",
      "id": "phbBoonofTruesig",
      "journalLink": "@UUID[Compendium.dnd5e.feats24.Item.phbBoonofTruesig]{Boon of Truesight}",
      "name": "Boon of Truesight",
      "package": "dnd5e.feats24",
      "packageName": "Feats",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Feats",
      "uuid": "Compendium.dnd5e.feats24.Item.phbBoonofTruesig"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Breastplate Armor of Resistance",
      "icon": "icons/equipment/chest/breastplate-collared-steel-grey.webp",
      "id": "lccm5AjIk91aIHbi",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.lccm5AjIk91aIHbi]{Breastplate Armor of Resistance}",
      "name": "Breastplate Armor of Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.lccm5AjIk91aIHbi"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Candle of Invocation: Outer Plane Destination",
      "icon": "icons/sundries/lights/candle-lit-yellow.webp",
      "id": "dmgCandleOfInvoc",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.dmgCandleOfInvoc]{Candle of Invocation: Outer Plane Destination}",
      "name": "Candle of Invocation: Outer Plane Destination",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.dmgCandleOfInvoc"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Carpenter's Tools",
      "icon": "icons/tools/hand/saw-steel-grey.webp",
      "id": "8NS6MSOdXtUqD7Ib",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.8NS6MSOdXtUqD7Ib]{Carpenter's Tools}",
      "name": "Carpenter's Tools",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.8NS6MSOdXtUqD7Ib"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Carpenter's Tools",
      "icon": "icons/tools/hand/saw-steel-grey.webp",
      "id": "phbtulCarpenters",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.phbtulCarpenters]{Carpenter's Tools}",
      "name": "Carpenter's Tools",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.phbtulCarpenters"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Cartographer's Tools",
      "icon": "icons/tools/navigation/map-chart-tan.webp",
      "id": "fC0lFK8P4RuhpfaU",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.fC0lFK8P4RuhpfaU]{Cartographer's Tools}",
      "name": "Cartographer's Tools",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.fC0lFK8P4RuhpfaU"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Cartographer's Tools",
      "icon": "icons/tools/navigation/map-chart-tan.webp",
      "id": "phbtulCartograph",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.phbtulCartograph]{Cartographer's Tools}",
      "name": "Cartographer's Tools",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.phbtulCartograph"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Chain Shirt Armor of Resistance",
      "icon": "icons/equipment/chest/breastplate-scale-grey.webp",
      "id": "HF32aZSVw4P0MR4K",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.HF32aZSVw4P0MR4K]{Chain Shirt Armor of Resistance}",
      "name": "Chain Shirt Armor of Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.HF32aZSVw4P0MR4K"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "Chapter 1: Beyond 1st Level",
      "icon": "",
      "id": "5LoAJLkfIYBAgWTW",
      "journalLink": "@UUID[Compendium.dnd5e.rules.JournalEntry.5LoAJLkfIYBAgWTW]{Chapter 1: Beyond 1st Level}",
      "name": "Chapter 1: Beyond 1st Level",
      "package": "dnd5e.rules",
      "packageName": "Rules (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules (SRD)",
      "uuid": "Compendium.dnd5e.rules.JournalEntry.5LoAJLkfIYBAgWTW"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "Chapter 10: Spellcasting",
      "icon": "",
      "id": "QvPDSUsAiEn3hD8s",
      "journalLink": "@UUID[Compendium.dnd5e.rules.JournalEntry.QvPDSUsAiEn3hD8s]{Chapter 10: Spellcasting}",
      "name": "Chapter 10: Spellcasting",
      "package": "dnd5e.rules",
      "packageName": "Rules (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules (SRD)",
      "uuid": "Compendium.dnd5e.rules.JournalEntry.QvPDSUsAiEn3hD8s"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "Chapter 4: Personality and Background",
      "icon": "",
      "id": "kWXplnmp5JXCo84x",
      "journalLink": "@UUID[Compendium.dnd5e.rules.JournalEntry.kWXplnmp5JXCo84x]{Chapter 4: Personality and Background}",
      "name": "Chapter 4: Personality and Background",
      "package": "dnd5e.rules",
      "packageName": "Rules (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules (SRD)",
      "uuid": "Compendium.dnd5e.rules.JournalEntry.kWXplnmp5JXCo84x"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "Chapter 6: Customization Options",
      "icon": "",
      "id": "hgHJdp8lTiJ5TpN9",
      "journalLink": "@UUID[Compendium.dnd5e.rules.JournalEntry.hgHJdp8lTiJ5TpN9]{Chapter 6: Customization Options}",
      "name": "Chapter 6: Customization Options",
      "package": "dnd5e.rules",
      "packageName": "Rules (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules (SRD)",
      "uuid": "Compendium.dnd5e.rules.JournalEntry.hgHJdp8lTiJ5TpN9"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "Chapter 7: Using Ability Scores",
      "icon": "",
      "id": "0AGfrwZRzSG0vNKb",
      "journalLink": "@UUID[Compendium.dnd5e.rules.JournalEntry.0AGfrwZRzSG0vNKb]{Chapter 7: Using Ability Scores}",
      "name": "Chapter 7: Using Ability Scores",
      "package": "dnd5e.rules",
      "packageName": "Rules (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules (SRD)",
      "uuid": "Compendium.dnd5e.rules.JournalEntry.0AGfrwZRzSG0vNKb"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Clothes, Traveler's",
      "icon": "icons/equipment/chest/shirt-collared-yellow.webp",
      "id": "phbagClothesTrav",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.phbagClothesTrav]{Clothes, Traveler's}",
      "name": "Clothes, Traveler's",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.phbagClothesTrav"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Create or Destroy Water",
      "icon": "icons/magic/air/wind-swirl-purple-blue.webp",
      "id": "a3XtAO5n2GrqiAh5",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.a3XtAO5n2GrqiAh5]{Create or Destroy Water}",
      "name": "Create or Destroy Water",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.a3XtAO5n2GrqiAh5"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Create or Destroy Water",
      "icon": "icons/magic/water/water-hand.webp",
      "id": "phbsplCreateorDe",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplCreateorDe]{Create or Destroy Water}",
      "name": "Create or Destroy Water",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplCreateorDe"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Create Specter",
      "icon": "icons/creatures/abilities/dragon-breath-purple.webp",
      "id": "SlAF2AE4ZKoUvQql",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.SlAF2AE4ZKoUvQql]{Create Specter}",
      "name": "Create Specter",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.SlAF2AE4ZKoUvQql"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Decanter of Endless Water",
      "icon": "icons/consumables/potions/potion-flask-corked-blue.webp",
      "id": "qXcUKfCVxEvV3VU8",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.qXcUKfCVxEvV3VU8]{Decanter of Endless Water}",
      "name": "Decanter of Endless Water",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.qXcUKfCVxEvV3VU8"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Decanter of Endless Water",
      "icon": "icons/consumables/potions/potion-flask-corked-blue.webp",
      "id": "dmgDecanterOfEnd",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgDecanterOfEnd]{Decanter of Endless Water}",
      "name": "Decanter of Endless Water",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgDecanterOfEnd"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Dominate Beast",
      "icon": "icons/magic/air/air-burst-spiral-large-teal-green.webp",
      "id": "LrPvWHBPmiMQQsKB",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.LrPvWHBPmiMQQsKB]{Dominate Beast}",
      "name": "Dominate Beast",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.LrPvWHBPmiMQQsKB"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Dominate Beast",
      "icon": "icons/creatures/mammals/bull-horns-eyes-glowin-orange.webp",
      "id": "phbsplDominateBe",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplDominateBe]{Dominate Beast}",
      "name": "Dominate Beast",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplDominateBe"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Dominate Monster",
      "icon": "icons/magic/air/air-burst-spiral-large-pink.webp",
      "id": "eEpy1ONlXumKS1mp",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.eEpy1ONlXumKS1mp]{Dominate Monster}",
      "name": "Dominate Monster",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.eEpy1ONlXumKS1mp"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Dominate Monster",
      "icon": "icons/magic/control/hypnosis-mesmerism-watch.webp",
      "id": "phbsplDominateMo",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplDominateMo]{Dominate Monster}",
      "name": "Dominate Monster",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplDominateMo"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Elemental Fury: Potent Spellcasting",
      "icon": "icons/magic/nature/leaf-rune-glow-green.webp",
      "id": "phbFuryPotentSpe",
      "journalLink": "@UUID[Compendium.dnd5e.classes24.Item.phbFuryPotentSpe]{Elemental Fury: Potent Spellcasting}",
      "name": "Elemental Fury: Potent Spellcasting",
      "package": "dnd5e.classes24",
      "packageName": "Character Classes",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Classes",
      "uuid": "Compendium.dnd5e.classes24.Item.phbFuryPotentSpe"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ethereal Sight",
      "icon": "icons/magic/perception/eye-tendrils-web-purple.webp",
      "id": "We6R4thWKYDRYlEc",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.We6R4thWKYDRYlEc]{Ethereal Sight}",
      "name": "Ethereal Sight",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.We6R4thWKYDRYlEc"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ethereal Stride",
      "icon": "icons/magic/lightning/orb-ball-purple.webp",
      "id": "NfTCXq8eRrqjhvAo",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.NfTCXq8eRrqjhvAo]{Ethereal Stride}",
      "name": "Ethereal Stride",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.NfTCXq8eRrqjhvAo"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Faultless Tracker",
      "icon": "icons/creatures/mammals/wolf-howl-moon-forest-blue.webp",
      "id": "E8SiDA7Z3Ybd6wt0",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.E8SiDA7Z3Ybd6wt0]{Faultless Tracker}",
      "name": "Faultless Tracker",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.E8SiDA7Z3Ybd6wt0"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Feather Token (Swan Boat)",
      "icon": "icons/commodities/materials/feather-blue-grey.webp",
      "id": "dmgSwanBoatQuaal",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgSwanBoatQuaal]{Feather Token (Swan Boat)}",
      "name": "Feather Token (Swan Boat)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgSwanBoatQuaal"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Feather Token Swan Boat",
      "icon": "icons/commodities/materials/feather-blue-grey.webp",
      "id": "UgnUJhu0tW1tLt7g",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.UgnUJhu0tW1tLt7g]{Feather Token Swan Boat}",
      "name": "Feather Token Swan Boat",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.UgnUJhu0tW1tLt7g"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Flame Tongue Scimitar",
      "icon": "icons/weapons/swords/sword-hooked-engraved.webp",
      "id": "qVHCzgVvOZAtuk4N",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.qVHCzgVvOZAtuk4N]{Flame Tongue Scimitar}",
      "name": "Flame Tongue Scimitar",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.qVHCzgVvOZAtuk4N"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Flame Tongue Shortsword",
      "icon": "icons/weapons/swords/sword-guard-red-jewel.webp",
      "id": "Z9FBwEoMi6daDGRj",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.Z9FBwEoMi6daDGRj]{Flame Tongue Shortsword}",
      "name": "Flame Tongue Shortsword",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.Z9FBwEoMi6daDGRj"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Giant Ancestry",
      "icon": "icons/creatures/magical/construct-iron-stomping-yellow.webp",
      "id": "phbsptGiantAnces",
      "journalLink": "@UUID[Compendium.dnd5e.origins24.Item.phbsptGiantAnces]{Giant Ancestry}",
      "name": "Giant Ancestry",
      "package": "dnd5e.origins24",
      "packageName": "Character Origins",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Origins",
      "uuid": "Compendium.dnd5e.origins24.Item.phbsptGiantAnces"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Giant Slayer Scimitar",
      "icon": "icons/weapons/swords/scimitar-broad.webp",
      "id": "ZLpj1bpnWlAFUEHE",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.ZLpj1bpnWlAFUEHE]{Giant Slayer Scimitar}",
      "name": "Giant Slayer Scimitar",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.ZLpj1bpnWlAFUEHE"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Giant Slayer Shortsword",
      "icon": "icons/weapons/swords/sword-guard-red.webp",
      "id": "tTqixDDmzAfs995G",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.tTqixDDmzAfs995G]{Giant Slayer Shortsword}",
      "name": "Giant Slayer Shortsword",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.tTqixDDmzAfs995G"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Greater Invisibility",
      "icon": "icons/magic/air/fog-gas-smoke-swirling-gray.webp",
      "id": "tEpDmYZNGc9f5OhJ",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.tEpDmYZNGc9f5OhJ]{Greater Invisibility}",
      "name": "Greater Invisibility",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.tEpDmYZNGc9f5OhJ"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Greater Invisibility",
      "icon": "icons/creatures/magical/spirit-undead-ghost-blue.webp",
      "id": "phbsplGreaterInv",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplGreaterInv]{Greater Invisibility}",
      "name": "Greater Invisibility",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplGreaterInv"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Greater Restoration",
      "icon": "icons/magic/life/heart-cross-strong-flame-blue.webp",
      "id": "WzvJ7G3cqvIubsLk",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.WzvJ7G3cqvIubsLk]{Greater Restoration}",
      "name": "Greater Restoration",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.WzvJ7G3cqvIubsLk"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Greater Restoration",
      "icon": "icons/magic/life/heart-hand-gold-green-light.webp",
      "id": "phbsplGreaterRes",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplGreaterRes]{Greater Restoration}",
      "name": "Greater Restoration",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplGreaterRes"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Greatsword of Life Stealing",
      "icon": "icons/weapons/swords/greatsword-crossguard-blue.webp",
      "id": "sdHSbitJxgTX6aDG",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.sdHSbitJxgTX6aDG]{Greatsword of Life Stealing}",
      "name": "Greatsword of Life Stealing",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.sdHSbitJxgTX6aDG"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Half Plate Armor of Resistance",
      "icon": "icons/equipment/chest/breastplate-cuirass-steel-grey.webp",
      "id": "lN1VbnGFo3HNZXNb",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.lN1VbnGFo3HNZXNb]{Half Plate Armor of Resistance}",
      "name": "Half Plate Armor of Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.lN1VbnGFo3HNZXNb"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "How to Use a Monster",
      "icon": "",
      "id": "mmMonsterManual1",
      "journalLink": "@UUID[Compendium.dnd5e.content24.JournalEntry.mmMonsterManual1]{How to Use a Monster}",
      "name": "How to Use a Monster",
      "package": "dnd5e.content24",
      "packageName": "Rules",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules",
      "uuid": "Compendium.dnd5e.content24.JournalEntry.mmMonsterManual1"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Innate Spellcasting",
      "icon": "icons/magic/light/projectiles-star-purple.webp",
      "id": "hkmTEk6klT6QL4K4",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.hkmTEk6klT6QL4K4]{Innate Spellcasting}",
      "name": "Innate Spellcasting",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.hkmTEk6klT6QL4K4"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Investment of the Chain Master",
      "icon": "icons/magic/control/debuff-chains-orb-movement-blue.webp",
      "id": "phbinvInvestment",
      "journalLink": "@UUID[Compendium.dnd5e.classes24.Item.phbinvInvestment]{Investment of the Chain Master}",
      "name": "Investment of the Chain Master",
      "package": "dnd5e.classes24",
      "packageName": "Character Classes",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Classes",
      "uuid": "Compendium.dnd5e.classes24.Item.phbinvInvestment"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Absorption",
      "icon": "icons/commodities/gems/gem-rough-ball-purple.webp",
      "id": "NGVEouqK0I6J6jV5",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.NGVEouqK0I6J6jV5]{Ioun Stone of Absorption}",
      "name": "Ioun Stone of Absorption",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.NGVEouqK0I6J6jV5"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Absorption",
      "icon": "icons/commodities/gems/gem-rough-ball-purple.webp",
      "id": "dmgAbsorptionIou",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgAbsorptionIou]{Ioun Stone of Absorption}",
      "name": "Ioun Stone of Absorption",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgAbsorptionIou"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Greater Absorption",
      "icon": "icons/commodities/stone/ore-pile-green.webp",
      "id": "7FEcfqz1piPHN1tV",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.7FEcfqz1piPHN1tV]{Ioun Stone of Greater Absorption}",
      "name": "Ioun Stone of Greater Absorption",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.7FEcfqz1piPHN1tV"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Greater Absorption",
      "icon": "icons/commodities/stone/ore-pile-green.webp",
      "id": "dmgGreaterAbsorp",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgGreaterAbsorp]{Ioun Stone of Greater Absorption}",
      "name": "Ioun Stone of Greater Absorption",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgGreaterAbsorp"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Insight",
      "icon": "icons/commodities/stone/ore-pile-teal.webp",
      "id": "9jMQEm99q1ttAV1Q",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.9jMQEm99q1ttAV1Q]{Ioun Stone of Insight}",
      "name": "Ioun Stone of Insight",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.9jMQEm99q1ttAV1Q"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Insight",
      "icon": "icons/commodities/stone/ore-pile-teal.webp",
      "id": "dmgInsightIounSt",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgInsightIounSt]{Ioun Stone of Insight}",
      "name": "Ioun Stone of Insight",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgInsightIounSt"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Mastery",
      "icon": "icons/commodities/gems/gem-rough-cushion-green.webp",
      "id": "nk2MH16KcZmKp7FQ",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.nk2MH16KcZmKp7FQ]{Ioun Stone of Mastery}",
      "name": "Ioun Stone of Mastery",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.nk2MH16KcZmKp7FQ"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Mastery",
      "icon": "icons/commodities/gems/gem-rough-cushion-green.webp",
      "id": "dmgMasteryIounSt",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgMasteryIounSt]{Ioun Stone of Mastery}",
      "name": "Ioun Stone of Mastery",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgMasteryIounSt"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Strength",
      "icon": "icons/commodities/gems/gem-rough-cushion-blue.webp",
      "id": "0G5LSgbb5NTV4XC7",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.0G5LSgbb5NTV4XC7]{Ioun Stone of Strength}",
      "name": "Ioun Stone of Strength",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.0G5LSgbb5NTV4XC7"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Strength",
      "icon": "icons/commodities/gems/gem-rough-cushion-blue.webp",
      "id": "dmgStrengthIounS",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgStrengthIounS]{Ioun Stone of Strength}",
      "name": "Ioun Stone of Strength",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgStrengthIounS"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Sustenance",
      "icon": "icons/commodities/stone/geode-raw-brown.webp",
      "id": "6MDTnMG4Hcw7qZsy",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.6MDTnMG4Hcw7qZsy]{Ioun Stone of Sustenance}",
      "name": "Ioun Stone of Sustenance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.6MDTnMG4Hcw7qZsy"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ioun Stone of Sustenance",
      "icon": "icons/commodities/stone/geode-raw-brown.webp",
      "id": "dmgSustenanceIou",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgSustenanceIou]{Ioun Stone of Sustenance}",
      "name": "Ioun Stone of Sustenance",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgSustenanceIou"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Leather Armor of Resistance",
      "icon": "icons/equipment/chest/breastplate-scale-leather.webp",
      "id": "dRtb9Tg34NKX9mGF",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.dRtb9Tg34NKX9mGF]{Leather Armor of Resistance}",
      "name": "Leather Armor of Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.dRtb9Tg34NKX9mGF"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Leatherworker's Tools",
      "icon": "icons/commodities/leather/leather-buckle-steel-tan.webp",
      "id": "PUMfwyVUbtyxgYbD",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.PUMfwyVUbtyxgYbD]{Leatherworker's Tools}",
      "name": "Leatherworker's Tools",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.PUMfwyVUbtyxgYbD"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Leatherworker's Tools",
      "icon": "icons/commodities/leather/leather-buckle-steel-tan.webp",
      "id": "phbtulLeatherwor",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.phbtulLeatherwor]{Leatherworker's Tools}",
      "name": "Leatherworker's Tools",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.phbtulLeatherwor"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Lessons of the First Ones",
      "icon": "icons/creatures/magical/humanoid-giant-forest-blue.webp",
      "id": "phbinvLessonsoft",
      "journalLink": "@UUID[Compendium.dnd5e.classes24.Item.phbinvLessonsoft]{Lessons of the First Ones}",
      "name": "Lessons of the First Ones",
      "package": "dnd5e.classes24",
      "packageName": "Character Classes",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Classes",
      "uuid": "Compendium.dnd5e.classes24.Item.phbinvLessonsoft"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Light Sensitivity",
      "icon": "icons/magic/time/day-night-sunset-sunrise.webp",
      "id": "2l557y06401lwsqs",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.2l557y06401lwsqs]{Light Sensitivity}",
      "name": "Light Sensitivity",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.2l557y06401lwsqs"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Locate Animals or Plants",
      "icon": "icons/magic/nature/leaf-glow-triple-orange-purple.webp",
      "id": "Iv2qqSAT7OkXKPFx",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.Iv2qqSAT7OkXKPFx]{Locate Animals or Plants}",
      "name": "Locate Animals or Plants",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.Iv2qqSAT7OkXKPFx"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Locate Animals or Plants",
      "icon": "icons/magic/nature/leaf-juggle-humanoid-green.webp",
      "id": "phbsplLocateAnim",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplLocateAnim]{Locate Animals or Plants}",
      "name": "Locate Animals or Plants",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplLocateAnim"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "Magic Item Lists",
      "icon": "",
      "id": "dmgMagicItemList",
      "journalLink": "@UUID[Compendium.dnd5e.content24.JournalEntry.dmgMagicItemList]{Magic Item Lists}",
      "name": "Magic Item Lists",
      "package": "dnd5e.content24",
      "packageName": "Rules",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules",
      "uuid": "Compendium.dnd5e.content24.JournalEntry.dmgMagicItemList"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Mantle of Spell Resistance",
      "icon": "icons/equipment/back/cape-layered-violet-white-swirl.webp",
      "id": "oxzUb5j1TMsccGW4",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.oxzUb5j1TMsccGW4]{Mantle of Spell Resistance}",
      "name": "Mantle of Spell Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.oxzUb5j1TMsccGW4"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Mantle of Spell Resistance",
      "icon": "icons/equipment/back/cape-layered-violet-white-swirl.webp",
      "id": "dmgMantleOfSpell",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgMantleOfSpell]{Mantle of Spell Resistance}",
      "name": "Mantle of Spell Resistance",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgMantleOfSpell"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Manual of Golems: Type, Time, and Cost",
      "icon": "icons/sundries/books/book-eye-purple.webp",
      "id": "dmgManualOfGolem",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.dmgManualOfGolem]{Manual of Golems: Type, Time, and Cost}",
      "name": "Manual of Golems: Type, Time, and Cost",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.dmgManualOfGolem"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Mithral Breastplate",
      "icon": "icons/equipment/chest/breastplate-collared-steel-grey.webp",
      "id": "CcTGZzQHejxEVLK1",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.CcTGZzQHejxEVLK1]{Mithral Breastplate}",
      "name": "Mithral Breastplate",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.CcTGZzQHejxEVLK1"
    },
    {
      "documentType": "JournalEntry",
      "folder": null,
      "formattedMatch": "Monsters A to Z",
      "icon": "",
      "id": "mmMonstersAtoZ00",
      "journalLink": "@UUID[Compendium.dnd5e.content24.JournalEntry.mmMonstersAtoZ00]{Monsters A to Z}",
      "name": "Monsters A to Z",
      "package": "dnd5e.content24",
      "packageName": "Rules",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Rules",
      "uuid": "Compendium.dnd5e.content24.JournalEntry.mmMonstersAtoZ00"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Nature's Sanctuary",
      "icon": "icons/magic/control/encase-creature-humanoid-hold.webp",
      "id": "EuX1kJNIw1F68yus",
      "journalLink": "@UUID[Compendium.dnd5e.classfeatures.Item.EuX1kJNIw1F68yus]{Nature's Sanctuary}",
      "name": "Nature's Sanctuary",
      "package": "dnd5e.classfeatures",
      "packageName": "Class & Subclass Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Class & Subclass Features (SRD)",
      "uuid": "Compendium.dnd5e.classfeatures.Item.EuX1kJNIw1F68yus"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Nature's Sanctuary",
      "icon": "icons/magic/nature/vines-thorned-curled-glow-green.webp",
      "id": "phbdrdNaturesSan",
      "journalLink": "@UUID[Compendium.dnd5e.classes24.Item.phbdrdNaturesSan]{Nature's Sanctuary}",
      "name": "Nature's Sanctuary",
      "package": "dnd5e.classes24",
      "packageName": "Character Classes",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Classes",
      "uuid": "Compendium.dnd5e.classes24.Item.phbdrdNaturesSan"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Nine Lives Stealer Scimitar",
      "icon": "icons/weapons/swords/scimitar-guard.webp",
      "id": "9Mdes2tKt0cqsNTw",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.9Mdes2tKt0cqsNTw]{Nine Lives Stealer Scimitar}",
      "name": "Nine Lives Stealer Scimitar",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.9Mdes2tKt0cqsNTw"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Nine Lives Stealer Shortsword",
      "icon": "icons/weapons/swords/shortsword-winged.webp",
      "id": "2Lkub0qIwucWEfp3",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.2Lkub0qIwucWEfp3]{Nine Lives Stealer Shortsword}",
      "name": "Nine Lives Stealer Shortsword",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.2Lkub0qIwucWEfp3"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "formattedMatch": "Otherworldly Steed",
      "icon": "",
      "id": "phbmobOtherworld",
      "journalLink": "@UUID[Compendium.dnd5e.actors24.Actor.phbmobOtherworld]{Otherworldly Steed}",
      "name": "Otherworldly Steed",
      "package": "dnd5e.actors24",
      "packageName": "Actors",
      "resultType": "CompendiumEntity",
      "subType": "npc",
      "tagline": "Actors",
      "uuid": "Compendium.dnd5e.actors24.Actor.phbmobOtherworld"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Plate Armor of Resistance",
      "icon": "icons/equipment/chest/breastplate-collared-steel.webp",
      "id": "azxwKFHrNmG3HpVy",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.azxwKFHrNmG3HpVy]{Plate Armor of Resistance}",
      "name": "Plate Armor of Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.azxwKFHrNmG3HpVy"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Acid Resistance",
      "icon": "icons/consumables/potions/bottle-bulb-corked-green.webp",
      "id": "zgZkJAyFAfYmyn11",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.zgZkJAyFAfYmyn11]{Potion of Acid Resistance}",
      "name": "Potion of Acid Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.zgZkJAyFAfYmyn11"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Cold Resistance",
      "icon": "icons/consumables/potions/bottle-bulb-corked-labeled-blue.webp",
      "id": "34YKlIJVVWLeBv7R",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.34YKlIJVVWLeBv7R]{Potion of Cold Resistance}",
      "name": "Potion of Cold Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.34YKlIJVVWLeBv7R"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Fire Giant Strength",
      "icon": "icons/consumables/potions/bottle-round-corked-yellow.webp",
      "id": "bEZOY6uvHRweMM56",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.bEZOY6uvHRweMM56]{Potion of Fire Giant Strength}",
      "name": "Potion of Fire Giant Strength",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.bEZOY6uvHRweMM56"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Fire Resistance",
      "icon": "icons/consumables/potions/bottle-bulb-corked-glowing-red.webp",
      "id": "Jj4iFQQGvckx8Wsj",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.Jj4iFQQGvckx8Wsj]{Potion of Fire Resistance}",
      "name": "Potion of Fire Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.Jj4iFQQGvckx8Wsj"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Force Resistance",
      "icon": "icons/consumables/potions/bottle-bulb-corked-labeled-blue.webp",
      "id": "kKGJjVVlJVoakWgQ",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.kKGJjVVlJVoakWgQ]{Potion of Force Resistance}",
      "name": "Potion of Force Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.kKGJjVVlJVoakWgQ"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Giant Strength (Frost)",
      "icon": "icons/consumables/potions/bottle-round-corked-blue.webp",
      "id": "dmgFrostPotionOf",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgFrostPotionOf]{Potion of Giant Strength (Frost)}",
      "name": "Potion of Giant Strength (Frost)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgFrostPotionOf"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Giant Strength (Stone)",
      "icon": "icons/consumables/potions/bottle-bulb-corked-green.webp",
      "id": "dmgStonePotionOf",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgStonePotionOf]{Potion of Giant Strength (Stone)}",
      "name": "Potion of Giant Strength (Stone)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgStonePotionOf"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Giant Strength (Storm)",
      "icon": "icons/consumables/potions/bottle-bulb-corked-labeled-blue.webp",
      "id": "dmgStormPotionOf",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgStormPotionOf]{Potion of Giant Strength (Storm)}",
      "name": "Potion of Giant Strength (Storm)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgStormPotionOf"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Lightning Resistance",
      "icon": "icons/consumables/potions/bottle-round-corked-yellow.webp",
      "id": "8MPnSrvEeZhPhtTi",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.8MPnSrvEeZhPhtTi]{Potion of Lightning Resistance}",
      "name": "Potion of Lightning Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.8MPnSrvEeZhPhtTi"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Necrotic Resistance",
      "icon": "icons/consumables/potions/bottle-round-corked-pink.webp",
      "id": "xw99pcqPBVwtMOLw",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.xw99pcqPBVwtMOLw]{Potion of Necrotic Resistance}",
      "name": "Potion of Necrotic Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.xw99pcqPBVwtMOLw"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Poison Resistance",
      "icon": "icons/consumables/potions/bottle-bulb-corked-green.webp",
      "id": "f5chGcpQCi1HYPQw",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.f5chGcpQCi1HYPQw]{Potion of Poison Resistance}",
      "name": "Potion of Poison Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.f5chGcpQCi1HYPQw"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Psychic Resistance",
      "icon": "icons/consumables/potions/bottle-round-corked-pink.webp",
      "id": "c0luemOP0iW8L23R",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.c0luemOP0iW8L23R]{Potion of Psychic Resistance}",
      "name": "Potion of Psychic Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.c0luemOP0iW8L23R"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Radiant Resistance",
      "icon": "icons/consumables/potions/bottle-round-corked-yellow.webp",
      "id": "LBQWNqX6hZOKhQ8a",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.LBQWNqX6hZOKhQ8a]{Potion of Radiant Resistance}",
      "name": "Potion of Radiant Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.LBQWNqX6hZOKhQ8a"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Potion of Resistance",
      "icon": "icons/consumables/potions/bottle-conical-fumes-green.webp",
      "id": "JzLOE4IxcmxjLLuz",
      "journalLink": "@UUID[Compendium.dnd5e.tables.RollTable.JzLOE4IxcmxjLLuz]{Potion of Resistance}",
      "name": "Potion of Resistance",
      "package": "dnd5e.tables",
      "packageName": "Tables (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Tables (SRD)",
      "uuid": "Compendium.dnd5e.tables.RollTable.JzLOE4IxcmxjLLuz"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Resistance",
      "icon": "icons/consumables/potions/potion-bottle-labeled-medicine-capped-red-black.webp",
      "id": "dmgPotionOfResis",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgPotionOfResis]{Potion of Resistance}",
      "name": "Potion of Resistance",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgPotionOfResis"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Potion of Resistance Type",
      "icon": "icons/consumables/potions/potion-bottle-labeled-medicine-capped-red-black.webp",
      "id": "dmgPotionOfResis",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.dmgPotionOfResis]{Potion of Resistance Type}",
      "name": "Potion of Resistance Type",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.dmgPotionOfResis"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Stone Giant Strength",
      "icon": "icons/consumables/potions/bottle-bulb-corked-green.webp",
      "id": "4ZiJsDTRA1GgcWKP",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.4ZiJsDTRA1GgcWKP]{Potion of Stone Giant Strength}",
      "name": "Potion of Stone Giant Strength",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.4ZiJsDTRA1GgcWKP"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potion of Thunder Resistance",
      "icon": "icons/consumables/potions/bottle-round-corked-yellow.webp",
      "id": "zBX8LLC2CjC89Dzl",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.zBX8LLC2CjC89Dzl]{Potion of Thunder Resistance}",
      "name": "Potion of Thunder Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.zBX8LLC2CjC89Dzl"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potter's Tools",
      "icon": "icons/containers/kitchenware/vase-bottle-brown.webp",
      "id": "hJS8yEVkqgJjwfWa",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.hJS8yEVkqgJjwfWa]{Potter's Tools}",
      "name": "Potter's Tools",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.hJS8yEVkqgJjwfWa"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Potter's Tools",
      "icon": "icons/containers/kitchenware/vase-bottle-brown.webp",
      "id": "phbtulPottersToo",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.phbtulPottersToo]{Potter's Tools}",
      "name": "Potter's Tools",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.phbtulPottersToo"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Private Sanctum",
      "icon": "icons/magic/defensive/shield-barrier-flaming-diamond-orange.webp",
      "id": "NJgxf7pmSsBArIG7",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.NJgxf7pmSsBArIG7]{Private Sanctum}",
      "name": "Private Sanctum",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.NJgxf7pmSsBArIG7"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Private Sanctum",
      "icon": "icons/environment/wilderness/cave-entrance-dwarven-hill.webp",
      "id": "phbPrivateSanctu",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbPrivateSanctu]{Private Sanctum}",
      "name": "Private Sanctum",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbPrivateSanctu"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Quarterstaff",
      "icon": "icons/weapons/staves/staff-simple.webp",
      "id": "g2dWN7PQiMRYWzyk",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.g2dWN7PQiMRYWzyk]{Quarterstaff}",
      "name": "Quarterstaff",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.g2dWN7PQiMRYWzyk"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Quarterstaff",
      "icon": "icons/weapons/staves/staff-simple.webp",
      "id": "phbwepQuartersta",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.phbwepQuartersta]{Quarterstaff}",
      "name": "Quarterstaff",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.phbwepQuartersta"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Quarterstaff +1",
      "icon": "icons/weapons/staves/staff-simple-carved.webp",
      "id": "t8L7B0JWamsvxhui",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.t8L7B0JWamsvxhui]{Quarterstaff +1}",
      "name": "Quarterstaff +1",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.t8L7B0JWamsvxhui"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Quarterstaff +2",
      "icon": "icons/weapons/staves/staff-orb-feather.webp",
      "id": "7kVZo4DLBq22406E",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.7kVZo4DLBq22406E]{Quarterstaff +2}",
      "name": "Quarterstaff +2",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.7kVZo4DLBq22406E"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Quarterstaff +3",
      "icon": "icons/weapons/staves/staff-ornate.webp",
      "id": "BmWnprrj0QWQ1BL3",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.BmWnprrj0QWQ1BL3]{Quarterstaff +3}",
      "name": "Quarterstaff +3",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.BmWnprrj0QWQ1BL3"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Quarterstaff of the Acrobat",
      "icon": "icons/weapons/staves/staff-simple-wrapped.webp",
      "id": "dmgQuarterstaffO",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgQuarterstaffO]{Quarterstaff of the Acrobat}",
      "name": "Quarterstaff of the Acrobat",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgQuarterstaffO"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Relentless Hunter",
      "icon": "icons/skills/wounds/injury-pain-body-orange.webp",
      "id": "phbrgrRelentless",
      "journalLink": "@UUID[Compendium.dnd5e.classes24.Item.phbrgrRelentless]{Relentless Hunter}",
      "name": "Relentless Hunter",
      "package": "dnd5e.classes24",
      "packageName": "Character Classes",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Classes",
      "uuid": "Compendium.dnd5e.classes24.Item.phbrgrRelentless"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ring of Lightning Resistance",
      "icon": "icons/equipment/finger/ring-cabochon-engraved-gold-orange.webp",
      "id": "XJ8CG4UvLELCmOi2",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.XJ8CG4UvLELCmOi2]{Ring of Lightning Resistance}",
      "name": "Ring of Lightning Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.XJ8CG4UvLELCmOi2"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ring of Necrotic Resistance",
      "icon": "icons/equipment/finger/ring-faceted-grey.webp",
      "id": "qMGkmzfLHfXd7DiJ",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.qMGkmzfLHfXd7DiJ]{Ring of Necrotic Resistance}",
      "name": "Ring of Necrotic Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.qMGkmzfLHfXd7DiJ"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ring of Radiant Resistance",
      "icon": "icons/equipment/finger/ring-cabochon-thin-gold-orange.webp",
      "id": "IrC5LPbWNxlAQoK7",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.IrC5LPbWNxlAQoK7]{Ring of Radiant Resistance}",
      "name": "Ring of Radiant Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.IrC5LPbWNxlAQoK7"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Ring of Resistance: Damage Type and Gemstone",
      "icon": "icons/equipment/finger/ring-ball-silver.webp",
      "id": "dmgRingOfResista",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.dmgRingOfResista]{Ring of Resistance: Damage Type and Gemstone}",
      "name": "Ring of Resistance: Damage Type and Gemstone",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.dmgRingOfResista"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Ring of Thunder Resistance",
      "icon": "icons/equipment/finger/ring-faceted-silver-orange.webp",
      "id": "IpBBqr0r7JanyVn0",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.IpBBqr0r7JanyVn0]{Ring of Thunder Resistance}",
      "name": "Ring of Thunder Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.IpBBqr0r7JanyVn0"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Scimitar of Life Stealing",
      "icon": "icons/weapons/swords/sword-hooked-worn.webp",
      "id": "sfegfmo59MHJg2YC",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.sfegfmo59MHJg2YC]{Scimitar of Life Stealing}",
      "name": "Scimitar of Life Stealing",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.sfegfmo59MHJg2YC"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Secret Chest",
      "icon": "icons/magic/perception/eye-ringed-glow-angry-teal.webp",
      "id": "8sgwRh8NUNkn9Vi0",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.8sgwRh8NUNkn9Vi0]{Secret Chest}",
      "name": "Secret Chest",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.8sgwRh8NUNkn9Vi0"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Secret Chest",
      "icon": "icons/containers/chest/chest-simple-box-gold-brown.webp",
      "id": "phbsplLeomundsSe",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplLeomundsSe]{Secret Chest}",
      "name": "Secret Chest",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplLeomundsSe"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "formattedMatch": "Secret Chest",
      "icon": "icons/containers/chest/chest-simple-box-gold-brown.webp",
      "id": "phbsplSecretChes",
      "journalLink": "@UUID[Compendium.dnd5e.actors24.Actor.phbsplSecretChes]{Secret Chest}",
      "name": "Secret Chest",
      "package": "dnd5e.actors24",
      "packageName": "Actors",
      "resultType": "CompendiumEntity",
      "subType": "npc",
      "tagline": "Actors",
      "uuid": "Compendium.dnd5e.actors24.Actor.phbsplSecretChes"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Sentient Magic Items Alignment",
      "icon": "icons/weapons/polearms/spear-flared-silver-pink.webp",
      "id": "NdjHJMSSVWw5fHsL",
      "journalLink": "@UUID[Compendium.dnd5e.tables.RollTable.NdjHJMSSVWw5fHsL]{Sentient Magic Items Alignment}",
      "name": "Sentient Magic Items Alignment",
      "package": "dnd5e.tables",
      "packageName": "Tables (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Tables (SRD)",
      "uuid": "Compendium.dnd5e.tables.RollTable.NdjHJMSSVWw5fHsL"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Sentient Magic Items Communication",
      "icon": "icons/weapons/polearms/spear-flared-silver-pink.webp",
      "id": "BHckoLKDwoL9d5p3",
      "journalLink": "@UUID[Compendium.dnd5e.tables.RollTable.BHckoLKDwoL9d5p3]{Sentient Magic Items Communication}",
      "name": "Sentient Magic Items Communication",
      "package": "dnd5e.tables",
      "packageName": "Tables (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Tables (SRD)",
      "uuid": "Compendium.dnd5e.tables.RollTable.BHckoLKDwoL9d5p3"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Shortsword of Life Stealing",
      "icon": "icons/weapons/swords/sword-guard-worn.webp",
      "id": "902yxeFDwavpm6cv",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.902yxeFDwavpm6cv]{Shortsword of Life Stealing}",
      "name": "Shortsword of Life Stealing",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.902yxeFDwavpm6cv"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Speak with Beasts and Plants",
      "icon": "icons/commodities/currency/coin-engraved-sun-smile-copper.webp",
      "id": "59DUUDZet1J4PIlA",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.59DUUDZet1J4PIlA]{Speak with Beasts and Plants}",
      "name": "Speak with Beasts and Plants",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.59DUUDZet1J4PIlA"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Sphere of Annihilation Interaction Results",
      "icon": "icons/magic/unholy/orb-glowing-purple.webp",
      "id": "dmgSphereOfAnnih",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.dmgSphereOfAnnih]{Sphere of Annihilation Interaction Results}",
      "name": "Sphere of Annihilation Interaction Results",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.dmgSphereOfAnnih"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Splint Armor of Resistance",
      "icon": "icons/equipment/chest/breastplate-layered-steel.webp",
      "id": "JNkjtTxYmEC7W34O",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.JNkjtTxYmEC7W34O]{Splint Armor of Resistance}",
      "name": "Splint Armor of Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.JNkjtTxYmEC7W34O"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Steadfast",
      "icon": "icons/magic/defensive/illusion-evasion-echo-purple.webp",
      "id": "4N7S29kDROQ932pG",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.4N7S29kDROQ932pG]{Steadfast}",
      "name": "Steadfast",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.4N7S29kDROQ932pG"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Stirge Roosts",
      "icon": "icons/environment/settlement/city-night.webp",
      "id": "mmStirgeRoosts00",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.mmStirgeRoosts00]{Stirge Roosts}",
      "name": "Stirge Roosts",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.mmStirgeRoosts00"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Stone of Good Luck (Luckstone)",
      "icon": "icons/commodities/gems/gem-rough-rectangle-red.webp",
      "id": "296Zgo9RhltWShE1",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.296Zgo9RhltWShE1]{Stone of Good Luck (Luckstone)}",
      "name": "Stone of Good Luck (Luckstone)",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.296Zgo9RhltWShE1"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Stone of Good Luck (Luckstone)",
      "icon": "icons/commodities/gems/gem-rough-rectangle-red.webp",
      "id": "dmgStoneOfGoodLu",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgStoneOfGoodLu]{Stone of Good Luck (Luckstone)}",
      "name": "Stone of Good Luck (Luckstone)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgStoneOfGoodLu"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Studded Leather Armor of Resistance",
      "icon": "icons/equipment/chest/breastplate-rivited-red.webp",
      "id": "W1kDsFekjroIywuz",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.W1kDsFekjroIywuz]{Studded Leather Armor of Resistance}",
      "name": "Studded Leather Armor of Resistance",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.W1kDsFekjroIywuz"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Sunlight Sensitivity",
      "icon": "icons/magic/light/explosion-star-glow-blue-purple.webp",
      "id": "F14aW2Ke3I5ZtSg4",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.F14aW2Ke3I5ZtSg4]{Sunlight Sensitivity}",
      "name": "Sunlight Sensitivity",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.F14aW2Ke3I5ZtSg4"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Thieves' Cant",
      "icon": "icons/sundries/documents/document-torn-diagram-tan.webp",
      "id": "ohwfuwnvuoBWlSQr",
      "journalLink": "@UUID[Compendium.dnd5e.classfeatures.Item.ohwfuwnvuoBWlSQr]{Thieves' Cant}",
      "name": "Thieves' Cant",
      "package": "dnd5e.classfeatures",
      "packageName": "Class & Subclass Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Class & Subclass Features (SRD)",
      "uuid": "Compendium.dnd5e.classfeatures.Item.ohwfuwnvuoBWlSQr"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Thieves' Cant",
      "icon": "icons/sundries/documents/document-symbol-eye.webp",
      "id": "phbrgeThievesCan",
      "journalLink": "@UUID[Compendium.dnd5e.classes24.Item.phbrgeThievesCan]{Thieves' Cant}",
      "name": "Thieves' Cant",
      "package": "dnd5e.classes24",
      "packageName": "Character Classes",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Character Classes",
      "uuid": "Compendium.dnd5e.classes24.Item.phbrgeThievesCan"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Thieves' Tools",
      "icon": "icons/tools/hand/lockpicks-steel-grey.webp",
      "id": "woWZ1sO5IUVGzo58",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.woWZ1sO5IUVGzo58]{Thieves' Tools}",
      "name": "Thieves' Tools",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.woWZ1sO5IUVGzo58"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Thieves' Tools",
      "icon": "icons/tools/hand/lockpicks-steel-grey.webp",
      "id": "phbtulThievesToo",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.phbtulThievesToo]{Thieves' Tools}",
      "name": "Thieves' Tools",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.phbtulThievesToo"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Thunderous Greatclub",
      "icon": "icons/weapons/clubs/club-spiked-glowing.webp",
      "id": "dmgThunderousGre",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgThunderousGre]{Thunderous Greatclub}",
      "name": "Thunderous Greatclub",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgThunderousGre"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Time Stop",
      "icon": "icons/magic/time/clock-stopwatch-white-blue.webp",
      "id": "JYuRBwxpoFhXduvD",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.JYuRBwxpoFhXduvD]{Time Stop}",
      "name": "Time Stop",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.JYuRBwxpoFhXduvD"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Time Stop",
      "icon": "icons/magic/time/hourglass-yellow-green.webp",
      "id": "phbsplTimeStop00",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplTimeStop00]{Time Stop}",
      "name": "Time Stop",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplTimeStop00"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Tinker's Tools",
      "icon": "icons/commodities/cloth/thread-spindle-white-needle.webp",
      "id": "0d08g1i5WXnNrCNA",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.0d08g1i5WXnNrCNA]{Tinker's Tools}",
      "name": "Tinker's Tools",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.0d08g1i5WXnNrCNA"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Tinker's Tools",
      "icon": "icons/commodities/cloth/thread-spindle-white-needle.webp",
      "id": "phbtulTinkersToo",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.phbtulTinkersToo]{Tinker's Tools}",
      "name": "Tinker's Tools",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "tool",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.phbtulTinkersToo"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Tome of Understanding",
      "icon": "icons/sundries/books/book-turquoise-moon.webp",
      "id": "WnKWD1FuAFUE7f4v",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.WnKWD1FuAFUE7f4v]{Tome of Understanding}",
      "name": "Tome of Understanding",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.WnKWD1FuAFUE7f4v"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Tome of Understanding",
      "icon": "icons/sundries/books/book-turquoise-moon.webp",
      "id": "dmgTomeOfUnderst",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgTomeOfUnderst]{Tome of Understanding}",
      "name": "Tome of Understanding",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgTomeOfUnderst"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Traveler's Clothes",
      "icon": "icons/equipment/back/cloak-brown-collared-fur-white-tied.webp",
      "id": "SsAmWV6YBqeOFihT",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.SsAmWV6YBqeOFihT]{Traveler's Clothes}",
      "name": "Traveler's Clothes",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "equipment",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.SsAmWV6YBqeOFihT"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Tree Stride",
      "icon": "icons/magic/nature/leaf-glow-maple-orange-purple.webp",
      "id": "DUBgwHPakcLDkB6W",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.DUBgwHPakcLDkB6W]{Tree Stride}",
      "name": "Tree Stride",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.DUBgwHPakcLDkB6W"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Tree Stride",
      "icon": "icons/magic/nature/tree-spirit-green.webp",
      "id": "g4V02wJbEstUpwi9",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.g4V02wJbEstUpwi9]{Tree Stride}",
      "name": "Tree Stride",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.g4V02wJbEstUpwi9"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Tree Stride",
      "icon": "icons/magic/nature/tree-spirit-blue.webp",
      "id": "phbsplTreeStride",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplTreeStride]{Tree Stride}",
      "name": "Tree Stride",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplTreeStride"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "True Resurrection",
      "icon": "icons/magic/life/heart-cross-strong-flame-blue.webp",
      "id": "qLeEXZDbW5y4bmLY",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.qLeEXZDbW5y4bmLY]{True Resurrection}",
      "name": "True Resurrection",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.qLeEXZDbW5y4bmLY"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "True Resurrection",
      "icon": "icons/magic/life/ankh-gold-blue.webp",
      "id": "phbsplTrueResurr",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplTrueResurr]{True Resurrection}",
      "name": "True Resurrection",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplTrueResurr"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "True Strike",
      "icon": "icons/magic/fire/dagger-rune-enchant-blue-gray.webp",
      "id": "mGGlcLdggHwcL7MG",
      "journalLink": "@UUID[Compendium.dnd5e.spells.Item.mGGlcLdggHwcL7MG]{True Strike}",
      "name": "True Strike",
      "package": "dnd5e.spells",
      "packageName": "Spells (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells (SRD)",
      "uuid": "Compendium.dnd5e.spells.Item.mGGlcLdggHwcL7MG"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "True Strike",
      "icon": "icons/magic/unholy/hand-weapon-glow-black-green.webp",
      "id": "phbsplTrueStrike",
      "journalLink": "@UUID[Compendium.dnd5e.spells24.Item.phbsplTrueStrike]{True Strike}",
      "name": "True Strike",
      "package": "dnd5e.spells24",
      "packageName": "Spells",
      "resultType": "CompendiumEntity",
      "subType": "spell",
      "tagline": "Spells",
      "uuid": "Compendium.dnd5e.spells24.Item.phbsplTrueStrike"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Turn Resistance",
      "icon": "icons/magic/fire/flame-burning-creature-skeleton.webp",
      "id": "r9aMLZ7F3gSRLgRr",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.r9aMLZ7F3gSRLgRr]{Turn Resistance}",
      "name": "Turn Resistance",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.r9aMLZ7F3gSRLgRr"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Two-Person Tent",
      "icon": "icons/environment/wilderness/camp-improvised.webp",
      "id": "PanSr5EbqlfpSvwK",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.PanSr5EbqlfpSvwK]{Two-Person Tent}",
      "name": "Two-Person Tent",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "loot",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.PanSr5EbqlfpSvwK"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Vicious Quarterstaff",
      "icon": "icons/weapons/maces/mace-round-steel.webp",
      "id": "Z7xno2zMzRtqqUIQ",
      "journalLink": "@UUID[Compendium.dnd5e.items.Item.Z7xno2zMzRtqqUIQ]{Vicious Quarterstaff}",
      "name": "Vicious Quarterstaff",
      "package": "dnd5e.items",
      "packageName": "Items (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "weapon",
      "tagline": "Items (SRD)",
      "uuid": "Compendium.dnd5e.items.Item.Z7xno2zMzRtqqUIQ"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Voice of the Chain Master",
      "icon": "icons/creatures/mammals/humanoid-cat-skulking-teal.webp",
      "id": "k5M8gsl7MMcdjOjs",
      "journalLink": "@UUID[Compendium.dnd5e.classfeatures.Item.k5M8gsl7MMcdjOjs]{Voice of the Chain Master}",
      "name": "Voice of the Chain Master",
      "package": "dnd5e.classfeatures",
      "packageName": "Class & Subclass Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Class & Subclass Features (SRD)",
      "uuid": "Compendium.dnd5e.classfeatures.Item.k5M8gsl7MMcdjOjs"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Water Susceptibility",
      "icon": "icons/magic/water/pseudopod-swirl-blue.webp",
      "id": "5V7SCABXvIbnk2Zn",
      "journalLink": "@UUID[Compendium.dnd5e.monsterfeatures.Item.5V7SCABXvIbnk2Zn]{Water Susceptibility}",
      "name": "Water Susceptibility",
      "package": "dnd5e.monsterfeatures",
      "packageName": "Monster Features (SRD)",
      "resultType": "CompendiumEntity",
      "subType": "feat",
      "tagline": "Monster Features (SRD)",
      "uuid": "Compendium.dnd5e.monsterfeatures.Item.5V7SCABXvIbnk2Zn"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Water, fresh (Pint)",
      "icon": "icons/magic/water/water-drop-swirl-blue.webp",
      "id": "dmgspWaterfresh0",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgspWaterfresh0]{Water, fresh (Pint)}",
      "name": "Water, fresh (Pint)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgspWaterfresh0"
    },
    {
      "documentType": "Item",
      "folder": null,
      "formattedMatch": "Water, salt (Pint)",
      "icon": "icons/sundries/survival/waterskin-leather-brown.webp",
      "id": "dmgspWatersalt00",
      "journalLink": "@UUID[Compendium.dnd5e.equipment24.Item.dmgspWatersalt00]{Water, salt (Pint)}",
      "name": "Water, salt (Pint)",
      "package": "dnd5e.equipment24",
      "packageName": "Equipment",
      "resultType": "CompendiumEntity",
      "subType": "consumable",
      "tagline": "Equipment",
      "uuid": "Compendium.dnd5e.equipment24.Item.dmgspWatersalt00"
    },
    {
      "documentType": "RollTable",
      "folder": null,
      "formattedMatch": "Wraith Manifestations",
      "icon": "icons/magic/death/skull-energy-light-purple.webp",
      "id": "mmWraithManifest",
      "journalLink": "@UUID[Compendium.dnd5e.tables24.RollTable.mmWraithManifest]{Wraith Manifestations}",
      "name": "Wraith Manifestations",
      "package": "dnd5e.tables24",
      "packageName": "Roll Tables",
      "resultType": "CompendiumEntity",
      "subType": "",
      "tagline": "Roll Tables",
      "uuid": "Compendium.dnd5e.tables24.RollTable.mmWraithManifest"
    }
  ],
  "type": "search-result"
}
```

### chat-event

Received chat-event via WS subscription

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const ws = new WebSocket('ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID');

ws.onopen = () => {
  // Send auth message first — token must not be in the URL
  ws.send(JSON.stringify({ type: 'auth', token: 'YOUR_API_KEY' }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  if (data.type === 'connected') {
    // Now send your request
    ws.send(JSON.stringify({
      "type": "chat-event",
      "requestId": "unique-id"
    }));
  }
  if (data.type === 'chat-event-result') {
    console.log(data);
  }
};
```

</TabItem>
<TabItem value="python" label="Python">

```python
import asyncio
import websockets
import json

async def main():
    uri = 'ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID'
    async with websockets.connect(uri) as ws:
        # Send auth message first — token must not be in the URL
        await ws.send(json.dumps({'type': 'auth', 'token': 'YOUR_API_KEY'}))
        connected = json.loads(await ws.recv())
        if connected.get('type') != 'connected':
            raise Exception('Auth failed')
        await ws.send(json.dumps({"type":"chat-event","requestId":"unique-id"}))
        response = await ws.recv()
        data = json.loads(response)
        print(data)

asyncio.run(main())
```

</TabItem>
<TabItem value="typescript" label="TypeScript">

```typescript
import WebSocket from 'ws';

const ws = new WebSocket('ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID');

ws.on('open', () => {
  // Send auth message first — token must not be in the URL
  ws.send(JSON.stringify({ type: 'auth', token: 'YOUR_API_KEY' }));
});

ws.on('message', (raw: string) => {
  const data = JSON.parse(raw);
  if (data.type === 'connected') {
    // Now send your request
    ws.send(JSON.stringify({
      "type": "chat-event",
      "requestId": "unique-id"
    }));
  }
  if (data.type === 'chat-event-result') {
    console.log(data);
  }
}):
```

</TabItem>
</Tabs>

#### Response

```json
{
  "data": {
    "data": {
      "data": {
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "content": "WS chat-event test",
        "flags": {},
        "flavor": "",
        "id": "xbxNmfQ3A2PZOl4s",
        "isRoll": false,
        "rolls": [],
        "speaker": {
          "actor": null,
          "scene": null,
          "token": null
        },
        "timestamp": 1776657997470,
        "type": "base",
        "uuid": "ChatMessage.xbxNmfQ3A2PZOl4s",
        "whisper": []
      },
      "eventType": "create"
    },
    "type": "chat-event"
  },
  "event": "chat-create",
  "type": "chat-event"
}
```

### roll-event

Received roll-event via WS subscription

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const ws = new WebSocket('ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID');

ws.onopen = () => {
  // Send auth message first — token must not be in the URL
  ws.send(JSON.stringify({ type: 'auth', token: 'YOUR_API_KEY' }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  if (data.type === 'connected') {
    // Now send your request
    ws.send(JSON.stringify({
      "type": "roll-event",
      "requestId": "unique-id"
    }));
  }
  if (data.type === 'roll-event-result') {
    console.log(data);
  }
};
```

</TabItem>
<TabItem value="python" label="Python">

```python
import asyncio
import websockets
import json

async def main():
    uri = 'ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID'
    async with websockets.connect(uri) as ws:
        # Send auth message first — token must not be in the URL
        await ws.send(json.dumps({'type': 'auth', 'token': 'YOUR_API_KEY'}))
        connected = json.loads(await ws.recv())
        if connected.get('type') != 'connected':
            raise Exception('Auth failed')
        await ws.send(json.dumps({"type":"roll-event","requestId":"unique-id"}))
        response = await ws.recv()
        data = json.loads(response)
        print(data)

asyncio.run(main())
```

</TabItem>
<TabItem value="typescript" label="TypeScript">

```typescript
import WebSocket from 'ws';

const ws = new WebSocket('ws://localhost:3010/ws/api?clientId=YOUR_CLIENT_ID');

ws.on('open', () => {
  // Send auth message first — token must not be in the URL
  ws.send(JSON.stringify({ type: 'auth', token: 'YOUR_API_KEY' }));
});

ws.on('message', (raw: string) => {
  const data = JSON.parse(raw);
  if (data.type === 'connected') {
    // Now send your request
    ws.send(JSON.stringify({
      "type": "roll-event",
      "requestId": "unique-id"
    }));
  }
  if (data.type === 'roll-event-result') {
    console.log(data);
  }
}):
```

</TabItem>
</Tabs>

#### Response

```json
{
  "data": {
    "data": {
      "dice": [
        {
          "faces": 6,
          "results": [
            {
              "active": true,
              "result": 6
            }
          ]
        }
      ],
      "flavor": "WS roll-event test",
      "formula": "1d6",
      "id": "1ba7afBjJTiWWgUB",
      "isCritical": false,
      "isFumble": false,
      "messageId": "1ba7afBjJTiWWgUB",
      "rollTotal": 6,
      "speaker": {
        "actor": null,
        "scene": null,
        "token": null
      },
      "timestamp": 1776657997490,
      "user": {
        "id": "r6bXhB7k9cXa3cif",
        "name": "tester"
      }
    },
    "type": "roll-data"
  },
  "type": "roll-event"
}
```

