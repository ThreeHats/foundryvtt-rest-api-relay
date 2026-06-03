---
tag: entity
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Entity

## GET /get

Get entity details

This endpoint retrieves the details of a specific entity.

**Required scope:** `entity:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to retrieve (optional if selected=true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| actor | boolean |  | query | Return the actor of specified entity |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Entity details object containing requested information

### Try It Out

<ApiTester
  method="GET"
  path="/get"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"actor","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /create

Create a new entity

This endpoint creates a new entity in the Foundry world.

**Required scope:** `entity:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| entityType | string | ✓ | body | Document type of entity to create (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist, ext.) |
| data | object | ✓ | body | Data for the new entity |
| clientId | string |  | query | Client ID for the Foundry world |
| folder | string |  | body | Optional folder UUID to place the new entity in |
| keepId | boolean |  | body, query | If true, preserve the _id from the provided data instead of generating a new one |
| override | boolean |  | body, query | If true and keepId is set, replace any existing entity with the same ID |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the entity creation operation

### Try It Out

<ApiTester
  method="POST"
  path="/create"
  parameters={[{"name":"entityType","type":"string","required":true,"source":"body"},{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"folder","type":"string","required":false,"source":"body"},{"name":"keepId","type":"boolean","required":false,"source":"body"},{"name":"override","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## PUT /update

Update an entity

This endpoint updates an existing entity in the Foundry world.

**Required scope:** `entity:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| data | object | ✓ | body | Object containing the fields to update |
| clientId | string |  | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to retrieve (optional if selected=true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| actor | boolean |  | query | Whether to update the actor of specified entity |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the entity update operation

### Try It Out

<ApiTester
  method="PUT"
  path="/update"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"actor","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## DELETE /delete

Delete an entity

This endpoint deletes an entity from the Foundry world.

**Required scope:** `entity:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to retrieve (optional if selected=true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the entity deletion operation

### Try It Out

<ApiTester
  method="DELETE"
  path="/delete"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /give

Give an item to another entity

Transfers an item from one entity to another.

**Required scope:** `entity:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| fromUuid | string |  | body | UUID of the entity giving the item |
| toUuid | string |  | body | UUID of the entity receiving the item |
| selected | boolean |  | query, body | Whether to get the selected entity |
| itemUuid | string |  | body | UUID of the item to give |
| itemName | string |  | body | Name of the item to give (alternative to itemUuid) |
| quantity | number |  | body | Quantity of items to give |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the give operation

### Try It Out

<ApiTester
  method="POST"
  path="/give"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"fromUuid","type":"string","required":false,"source":"body"},{"name":"toUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"},{"name":"quantity","type":"number","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /remove

Remove an item from an entity

Removes an item from an entity's inventory.

**Required scope:** `entity:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body | UUID of the actor to remove the item from |
| selected | boolean |  | query, body | Whether to get the selected entity |
| itemUuid | string |  | body | UUID of the item to remove |
| itemName | string |  | body | Name of the item to remove (alternative to itemUuid) |
| quantity | number |  | body | Quantity of items to remove |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the remove operation

### Try It Out

<ApiTester
  method="POST"
  path="/remove"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"},{"name":"quantity","type":"number","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /decrease

Decrease an attribute

Decreases a numeric attribute of an entity by the specified amount.

**Required scope:** `entity:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| attribute | string | ✓ | body | The attribute to decrease (e.g., hp.value) |
| amount | number | ✓ | body | The amount to decrease by |
| clientId | string |  | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to retrieve (optional if selected=true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the decrease operation

### Try It Out

<ApiTester
  method="POST"
  path="/decrease"
  parameters={[{"name":"attribute","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /increase

Increase an attribute

Increases a numeric attribute of an entity by the specified amount.

**Required scope:** `entity:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| attribute | string | ✓ | body | The attribute to increase (e.g., hp.value) |
| amount | number | ✓ | body | The amount to increase by |
| clientId | string |  | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to retrieve (optional if selected=true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the increase operation

### Try It Out

<ApiTester
  method="POST"
  path="/increase"
  parameters={[{"name":"attribute","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /kill

Kill an entity

Sets the entity's HP to 0.

**Required scope:** `entity:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to retrieve (optional if selected=true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the kill operation

### Try It Out

<ApiTester
  method="POST"
  path="/kill"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

