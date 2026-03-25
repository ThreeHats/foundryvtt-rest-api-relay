---
tag: WebSocket
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# WebSocket API

The WebSocket API provides bidirectional communication with Foundry VTT through the relay server. It supports the same operations as the REST API, plus real-time event subscriptions.

## Connection

Connect to the WebSocket endpoint with your API key and target Foundry client ID:

```
ws://<host>/ws/api?token=<apiKey>&clientId=<clientId>
```

On successful connection, you will receive a `connected` message listing all supported message types and event channels.

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

## Supported Message Types

| Type | Description | Required Params |
|------|-------------|-----------------|
| `json` | Get file system structure | — |
| `get-canvas-documents` | Get canvas embedded documents | `documentType` |
| `create-canvas-document` | Create canvas embedded document(s) | `documentType` |
| `update-canvas-document` | Update a canvas embedded document | `documentType`, `documentId`, `data` |
| `delete-canvas-document` | Delete a canvas embedded document | `documentType`, `documentId` |
| `chat-messages` | Get chat messages Retrieves chat messages from the Foundry world with optional p | — |
| `chat-send` | Send a chat message Creates a new chat message in the Foundry world. | `content` |
| `chat-delete` | Delete a specific chat message Deletes a chat message by its ID. Only the messag | `messageId` |
| `chat-flush` | Clear all chat messages Flushes all chat message history. Only GMs can perform t | — |
| `get-actor-details` | Get detailed information for a specific D&D 5e actor. Retrieves comprehensive de | `actorUuid`, `details` |
| `modify-item-charges` | Modify the charges for a specific item owned by an actor. Increases or decreases | `actorUuid`, `amount` |
| `use-ability` | Use a general ability for an actor. Triggers the use of any ability, feature, sp | `actorUuid` |
| `use-feature` | Use a class or racial feature for an actor. Activates class features (like Actio | `actorUuid` |
| `use-spell` | Cast a spell for an actor. Casts a spell from the actor's spell list, consuming  | `actorUuid` |
| `use-item` | Use an item for an actor. Activates an item from the actor's inventory, such as  | `actorUuid` |
| `short-rest` | Perform a short rest for an actor. Triggers the D&D 5e short rest workflow inclu | — |
| `long-rest` | Perform a long rest for an actor. Triggers the D&D 5e long rest workflow includi | — |
| `skill-check` | Roll a skill check for an actor. Rolls a D&D 5e skill check with all applicable  | `actorUuid`, `skill` |
| `ability-save` | Roll an ability saving throw for an actor. Rolls a D&D 5e ability saving throw w | `actorUuid`, `ability` |
| `ability-check` | Roll an ability check for an actor. Rolls a D&D 5e ability check (raw ability te | `actorUuid`, `ability` |
| `death-save` | Roll a death saving throw for an actor. Rolls a D&D 5e death saving throw, handl | `actorUuid` |
| `modify-experience` | Modify the experience points for a specific actor. Adds or removes experience po | `amount` |
| `get-effects` | Get all active effects on an actor or token. Returns the collection of ActiveEff | `uuid` |
| `add-effect` | Add an active effect to an actor or token. Adds a status condition (by statusId) | `uuid` |
| `remove-effect` | Remove an active effect from an actor or token. Removes an effect by its documen | `uuid` |
| `encounters` | Get all active encounters Retrieves a list of all currently active encounters in | — |
| `start-encounter` | Start a new encounter Initiates a new encounter in the Foundry world. | — |
| `next-turn` | Advance to the next turn in the encounter Moves the encounter to the next turn. | — |
| `next-round` | Advance to the next round in the encounter Moves the encounter to the next round | — |
| `last-turn` | Advance to the last turn in the encounter Moves the encounter to the last turn. | — |
| `last-round` | Advance to the last round in the encounter Moves the encounter to the last round | — |
| `end-encounter` | End an encounter Ends the current encounter in the Foundry world. | — |
| `add-to-encounter` | Add tokens to an encounter Adds selected tokens or specified UUIDs to the curren | — |
| `remove-from-encounter` | Remove tokens from an encounter Removes selected tokens or specified UUIDs from  | — |
| `entity` | Get entity details This endpoint retrieves the details of a specific entity. | — |
| `create` | Create a new entity This endpoint creates a new entity in the Foundry world. | `entityType`, `data` |
| `update` | Update an existing entity This endpoint updates an existing entity in the Foundr | `data` |
| `delete` | Delete an entity This endpoint deletes an entity from the Foundry world. | — |
| `give` | Give an item to an entity This endpoint gives an item to a specified entity. Opt | — |
| `remove` | Remove an item from an entity This endpoint removes an item from a specified ent | — |
| `decrease` | Decrease an attribute This endpoint decreases an attribute of a specified entity | `attribute`, `amount` |
| `increase` | Increase an attribute This endpoint increases an attribute of a specified entity | `attribute`, `amount` |
| `kill` | Kill an entity Marks an entity as killed in the combat tracker, gives it the "de | — |
| `macros` | Get all macros Retrieves a list of all macros available in the Foundry world. | — |
| `macro-execute` | Execute a macro by UUID Executes a specific macro in the Foundry world by its UU | `uuid` |
| `rolls` | Get recent rolls Retrieves a list of up to 20 recent rolls made in the Foundry w | — |
| `last-roll` | Get the last roll Retrieves the most recent roll made in the Foundry world. | — |
| `roll` | Make a roll Executes a roll with the specified formula | `formula` |
| `get-scene` | Get scene(s) Retrieves one or more scenes by ID, name, active status, viewed sta | — |
| `create-scene` | Create a new scene | `data` |
| `update-scene` | Update an existing scene | `data` |
| `delete-scene` | Delete a scene | — |
| `switch-scene` | Switch the active scene | — |
| `search` | Search entities This endpoint allows searching for entities in the Foundry world | `query` |
| `structure` | Get the structure of the Foundry world Retrieves the folder and compendium struc | — |
| `get-folder` | Get a specific folder by name | `name` |
| `create-folder` | Create a new folder | `name`, `folderType` |
| `delete-folder` | Delete a folder | `folderId` |
| `select` | Select token(s) Selects one or more tokens in the Foundry VTT client. | — |
| `selected` | Get selected token(s) Retrieves the currently selected token(s) in the Foundry V | — |
| `players` | Get players/users Retrieves a list of all users configured in the Foundry VTT wo | — |
| `execute-js` | Execute JavaScript Executes a JavaScript script in the Foundry VTT client. | — |

## AsyncAPI Spec

The full AsyncAPI specification is available at [`/asyncapi.json`](/asyncapi.json).

