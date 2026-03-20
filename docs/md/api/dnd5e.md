---
tag: dnd5e
---

import ApiTester from '@site/src/components/ApiTester';

# dnd5e

## GET /get-actor-details

Get detailed information for a specific D&D 5e actor. Retrieves comprehensive details about an actor including stats, inventory, spells, features, and other character information based on the requested details array.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor |
| details | array | ✓ | body, query | Array of detail types to retrieve (e.g., ["resources", "items", "spells", "features"]) |

### Returns

**object** - Actor details object containing requested information

### Try It Out

<ApiTester
  method="GET"
  path="/get-actor-details"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"details","type":"array","required":true,"source":"body"}]}
/>

---

## POST /modify-item-charges

Modify the charges for a specific item owned by an actor. Increases or decreases the charges/uses of an item in an actor's inventory. Useful for consumable items like potions, scrolls, or charged magic items.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor who owns the item |
| amount | number | ✓ | body, query | The amount to modify charges by (positive or negative) |
| itemUuid | string |  | body, query | The UUID of the specific item (optional if itemName provided) |
| itemName | string |  | body, query | The name of the item if UUID not provided (optional if itemUuid provided) |

### Returns

**object** - Result of the charge modification operation

### Try It Out

<ApiTester
  method="POST"
  path="/modify-item-charges"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"}]}
/>

---

## POST /use-ability

Use a general ability for an actor. Triggers the use of any ability, feature, spell, or item for an actor. This is a generic endpoint that can handle various types of abilities.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor using the ability |
| abilityUuid | string |  | body, query | The UUID of the specific ability (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the ability (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |

### Returns

**object** - Result of the ability use operation

### Try It Out

<ApiTester
  method="POST"
  path="/use-ability"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"}]}
/>

---

## POST /use-feature

Use a class or racial feature for an actor. Activates class features (like Action Surge, Rage) or racial features (like Dragonborn Breath Weapon) for a character.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor using the feature |
| abilityUuid | string |  | body, query | The UUID of the specific feature (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the feature if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the feature (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |

### Returns

**object** - Result of the feature use operation

### Try It Out

<ApiTester
  method="POST"
  path="/use-feature"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"}]}
/>

---

## POST /use-spell

Cast a spell for an actor. Casts a spell from the actor's spell list, consuming spell slots as appropriate. Handles cantrips, leveled spells, and spell-like abilities.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor casting the spell |
| abilityUuid | string |  | body, query | The UUID of the specific spell (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the spell if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the spell (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |

### Returns

**object** - Result of the spell casting operation

### Try It Out

<ApiTester
  method="POST"
  path="/use-spell"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"}]}
/>

---

## POST /use-item

Use an item for an actor. Activates an item from the actor's inventory, such as drinking a potion, using a magic item, or activating equipment with special properties.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor using the item |
| abilityUuid | string |  | body, query | The UUID of the specific item (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the item if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the item (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |

### Returns

**object** - Result of the item use operation

### Try It Out

<ApiTester
  method="POST"
  path="/use-item"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"}]}
/>

---

## POST /modify-experience

Modify the experience points for a specific actor. Adds or removes experience points from an actor.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| amount | number | ✓ | body, query | The amount of experience to add (can be negative) |
| actorUuid | string |  | body, query | UUID of the actor to modify |
| selected | boolean |  | body, query | Modify the selected token's actor |

### Returns

**object** - Result of the experience modification operation

### Try It Out

<ApiTester
  method="POST"
  path="/modify-experience"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"body"}]}
/>

