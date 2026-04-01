---
tag: WebSocket
---

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
| `entity` | Get entity details This endpoint retrieves the details of a specific entity. | — |
| `create` | Create a new entity This endpoint creates a new entity in the Foundry world. | `entityType`, `data` |
| `update` | Update an entity This endpoint updates an existing entity in the Foundry world. | `data` |
| `delete` | Delete an entity This endpoint deletes an entity from the Foundry world. | — |
| `give` | Give an item to another entity Transfers an item from one entity to another. | — |
| `remove` | Remove an item from an entity Removes an item from an entity's inventory. | — |
| `decrease` | Decrease an attribute Decreases a numeric attribute of an entity by the specified amount. | `attribute`, `amount` |
| `increase` | Increase an attribute Increases a numeric attribute of an entity by the specified amount. | `attribute`, `amount` |
| `kill` | Kill an entity Sets the entity's HP to 0. | — |
| `search` | Search entities This endpoint allows searching for entities in the Foundry world based on a query string. Requires Quick Insert module to be installed and enabled. | `query` |
| `rolls` | Get recent rolls Retrieves a list of up to 20 recent rolls made in the Foundry world. | — |
| `last-roll` | Get the last roll Retrieves the most recent roll made in the Foundry world. | — |
| `roll` | Make a roll Executes a roll with the specified formula. | `formula` |
| `encounters` | Get all active encounters Retrieves a list of all currently active encounters in the Foundry world. | — |
| `start-encounter` | Start a new encounter Initiates a new encounter in the Foundry world. | — |
| `next-turn` | Advance to the next turn in the encounter Moves the encounter to the next turn. | — |
| `next-round` | Advance to the next round in the encounter Moves the encounter to the next round. | — |
| `last-turn` | Go back to the last turn in the encounter Moves the encounter back to the last turn. | — |
| `last-round` | Go back to the last round in the encounter Moves the encounter back to the last round. | — |
| `end-encounter` | End an encounter Ends the current encounter in the Foundry world. | — |
| `add-to-encounter` | Add tokens to an encounter Adds selected tokens or specified UUIDs to the current encounter. | — |
| `remove-from-encounter` | Remove tokens from an encounter Removes selected tokens or specified UUIDs from the current encounter. | — |
| `macros` | Get all macros Retrieves a list of all macros available in the Foundry world. | — |
| `macro-execute` | Execute a macro by UUID Executes a specific macro in the Foundry world by its UUID. | `uuid` |
| `structure` | Get the structure of the Foundry world Retrieves the folder and compendium structure for the specified Foundry world. | — |
| `get-folder` | Get a specific folder by name | `name` |
| `create-folder` | Create a new folder | `name`, `folderType` |
| `delete-folder` | Delete a folder | `folderId` |
| `select` | Select token(s) Selects one or more tokens in the Foundry VTT client. | — |
| `selected` | Get selected token(s) Retrieves the currently selected token(s) in the Foundry VTT client. | — |
| `players` | Get players/users Retrieves a list of all users configured in the Foundry VTT world. Useful for discovering valid userId values for permission-scoped API calls. | — |
| `execute-js` | Execute JavaScript Executes a JavaScript script in the Foundry VTT client. | — |
| `get-scene` | Get scene(s) Retrieves one or more scenes by ID, name, active status, viewed status, or all. | — |
| `create-scene` | Create a new scene | `data` |
| `update-scene` | Update an existing scene | `data` |
| `delete-scene` | Delete a scene | — |
| `switch-scene` | Switch the active scene | — |
| `get-canvas-documents` | Get canvas embedded documents | `documentType` |
| `create-canvas-document` | Create canvas embedded document(s) | `documentType`, `data` |
| `update-canvas-document` | Update a canvas embedded document | `documentType`, `documentId`, `data` |
| `delete-canvas-document` | Delete a canvas embedded document | `documentType`, `documentId` |
| `chat-messages` | Get chat messages Retrieves chat messages from the Foundry world with optional pagination and filtering. | — |
| `chat-send` | Send a chat message Creates a new chat message in the Foundry world. | `content` |
| `chat-delete` | Delete a specific chat message Deletes a chat message by its ID. Only the message author or a GM can delete messages. | `messageId` |
| `chat-flush` | Clear all chat messages Flushes all chat message history. Only GMs can perform this action. | — |
| `get-effects` | Get all active effects on an actor or token Returns the collection of ActiveEffect documents currently applied to the specified actor or token. | `uuid` |
| `add-effect` | Add an active effect to an actor or token Adds a status condition (by statusId) or a custom ActiveEffect (via effectData) to the specified actor or token. | `uuid` |
| `remove-effect` | Remove an active effect from an actor or token Removes an effect by its document ID (effectId) or by status condition identifier (statusId). | `uuid` |
| `file-system` | Get file system structure | — |
| `download-file` | Download a file from Foundry's file system | — |
| `get-actor-details` | Get detailed information for a specific D&D 5e actor Retrieves comprehensive details about an actor including stats, inventory, spells, features, and other character information based on the requested details array. | `actorUuid`, `details` |
| `modify-item-charges` | Modify the charges for a specific item owned by an actor Increases or decreases the charges/uses of an item in an actor's inventory. Useful for consumable items like potions, scrolls, or charged magic items. | `actorUuid`, `amount` |
| `short-rest` | Perform a short rest for an actor Triggers the D&D 5e short rest workflow including hit dice recovery, class feature resets, and HP recovery. | — |
| `long-rest` | Perform a long rest for an actor Triggers the D&D 5e long rest workflow including full HP recovery, spell slot restoration, hit dice recovery, and feature resets. | — |
| `skill-check` | Roll a skill check for an actor Rolls a D&D 5e skill check with all applicable modifiers including proficiency, expertise, Jack of All Trades, and conditional bonuses. | `actorUuid`, `skill` |
| `ability-save` | Roll an ability saving throw for an actor Rolls a D&D 5e ability saving throw with all applicable modifiers. | `actorUuid`, `ability` |
| `ability-check` | Roll an ability check for an actor Rolls a D&D 5e ability check (raw ability test, not a skill check) with all applicable modifiers. | `actorUuid`, `ability` |
| `death-save` | Roll a death saving throw for an actor Rolls a D&D 5e death saving throw, handling DC 10 CON save, three successes/failures tracking, nat 20 healing, and nat 1 double failure. | `actorUuid` |
| `modify-experience` | Modify the experience points for a specific actor Adds or removes experience points from an actor. | `amount` |
| `use-ability` | Use an ability Activates a specific ability for an actor, optionally targeting another entity | `actorUuid` |
| `use-feature` | Use a feature Activates a specific feature for an actor, optionally targeting another entity | `actorUuid` |
| `use-spell` | Use a spell Casts a specific spell for an actor, optionally targeting another entity | `actorUuid` |
| `use-item` | Use an item Uses a specific item for an actor, optionally targeting another entity | `actorUuid` |

## AsyncAPI Spec

The full AsyncAPI specification is available at `/asyncapi.json`.
