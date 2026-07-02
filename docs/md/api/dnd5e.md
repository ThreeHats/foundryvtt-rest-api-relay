---
tag: dnd5e
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Dnd5e

## GET /dnd5e/get-actor-details

Get detailed information for a specific D&D 5e actor

Retrieves comprehensive details about an actor including stats, inventory, spells, features, and other character information based on the requested details array.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| details | array | ✓ | body, query | Array of detail types to retrieve (e.g., ["resources", "items", "spells", "features"]) |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Actor details object containing requested information

### Try It Out

<ApiTester
  method="GET"
  path="/dnd5e/get-actor-details"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"details","type":"array","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/modify-item-charges

Modify the charges for a specific item owned by an actor

Increases or decreases the charges/uses of an item in an actor's inventory. Useful for consumable items like potions, scrolls, or charged magic items.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| amount | number | ✓ | body, query | The amount to modify charges by (positive or negative) |
| clientId | string |  | query | Client ID for the Foundry world |
| itemUuid | string |  | body, query | The UUID of the specific item (optional if itemName provided) |
| itemName | string |  | body, query | The name of the item if UUID not provided (optional if itemUuid provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the charge modification operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/modify-item-charges"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/short-rest

Perform a short rest for an actor

Triggers the D&D 5e short rest workflow including hit dice recovery, class feature resets, and HP recovery.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor (optional if selected is true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| autoHD | boolean |  | body, query | Automatically spend hit dice during short rest |
| autoHDThreshold | number |  | body, query | HP threshold below which to auto-spend hit dice (0-1 as fraction of max HP) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the short rest operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/short-rest"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"autoHD","type":"boolean","required":false,"source":"body"},{"name":"autoHDThreshold","type":"number","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/long-rest

Perform a long rest for an actor

Triggers the D&D 5e long rest workflow including full HP recovery, spell slot restoration, hit dice recovery, and feature resets.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor (optional if selected is true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| newDay | boolean |  | body, query | Whether the long rest marks a new day (default: true) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the long rest operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/long-rest"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"newDay","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/skill-check

Roll a skill check for an actor

Rolls a D&D 5e skill check with all applicable modifiers including proficiency, expertise, Jack of All Trades, and conditional bonuses.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| skill | string | ✓ | body, query | Skill abbreviation (e.g., "acr", "ath", "ste", "prc") |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |
| advantage | boolean |  | body, query | Roll with advantage |
| disadvantage | boolean |  | body, query | Roll with disadvantage |
| bonus | string |  | body, query | Extra bonus formula to add (e.g., "1d4", "+2") |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |

### Returns

**object** - Result of the skill check roll

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/skill-check"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"skill","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"disadvantage","type":"boolean","required":false,"source":"body"},{"name":"bonus","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"}]}
/>

---

## POST /dnd5e/ability-save

Roll an ability saving throw for an actor

Rolls a D&D 5e ability saving throw with all applicable modifiers.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| ability | string | ✓ | body, query | Ability abbreviation (e.g., "str", "dex", "con", "int", "wis", "cha") |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |
| advantage | boolean |  | body, query | Roll with advantage |
| disadvantage | boolean |  | body, query | Roll with disadvantage |
| bonus | string |  | body, query | Extra bonus formula to add (e.g., "1d4", "+2") |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |

### Returns

**object** - Result of the saving throw roll

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/ability-save"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"ability","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"disadvantage","type":"boolean","required":false,"source":"body"},{"name":"bonus","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"}]}
/>

---

## POST /dnd5e/ability-check

Roll an ability check for an actor

Rolls a D&D 5e ability check (raw ability test, not a skill check) with all applicable modifiers.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| ability | string | ✓ | body, query | Ability abbreviation (e.g., "str", "dex", "con", "int", "wis", "cha") |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |
| advantage | boolean |  | body, query | Roll with advantage |
| disadvantage | boolean |  | body, query | Roll with disadvantage |
| bonus | string |  | body, query | Extra bonus formula to add (e.g., "1d4", "+2") |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |

### Returns

**object** - Result of the ability check roll

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/ability-check"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"ability","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"disadvantage","type":"boolean","required":false,"source":"body"},{"name":"bonus","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"}]}
/>

---

## POST /dnd5e/death-save

Roll a death saving throw for an actor

Rolls a D&D 5e death saving throw, handling DC 10 CON save, three successes/failures tracking, nat 20 healing, and nat 1 double failure.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| clientId | string |  | query | Client ID for the Foundry world |
| advantage | boolean |  | body, query | Roll with advantage |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the death saving throw

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/death-save"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/modify-experience

Modify the experience points for a specific actor

Adds or removes experience points from an actor.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| amount | number | ✓ | body, query | The amount of experience to add (can be negative) |
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor (optional if selected is true) |
| selected | boolean |  | query, body | Whether to get the selected entity |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the experience modification operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/modify-experience"
  parameters={[{"name":"amount","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /dnd5e/concentration

Check if an actor is concentrating on a spell

Returns whether the actor currently has a concentration effect active, and if so, what spell they are concentrating on.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor (optional if selected is true) |
| actorName | string |  | body, query | Name of the actor (optional if actorUuid provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Concentration status with effect details and spell name

### Try It Out

<ApiTester
  method="GET"
  path="/dnd5e/concentration"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"actorName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/break-concentration

Break an actor's concentration

Removes the concentration effect from the actor, ending any spell that requires concentration.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor (optional if selected is true) |
| actorName | string |  | body, query | Name of the actor (optional if actorUuid provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Confirmation that concentration was broken

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/break-concentration"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"actorName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/concentration-save

Roll a concentration saving throw

Rolls a Constitution saving throw to maintain concentration after taking damage. The DC is calculated as max(10, floor(damage/2)). Returns the roll result and whether concentration was maintained or broken.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| damage | number | ✓ | body, query | Amount of damage taken (used to calculate DC = max(10, floor(damage/2))) |
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor (optional if selected is true) |
| actorName | string |  | body, query | Name of the actor (optional if actorUuid provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |
| advantage | boolean |  | body, query | Roll with advantage |
| disadvantage | boolean |  | body, query | Roll with disadvantage |
| bonus | string |  | body, query | Extra bonus formula to add (e.g., "1d4", "+2") |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |

### Returns

**object** - Roll result and concentration maintained status

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/concentration-save"
  parameters={[{"name":"damage","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"actorName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"disadvantage","type":"boolean","required":false,"source":"body"},{"name":"bonus","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"}]}
/>

---

## POST /dnd5e/equip-item

Equip or unequip an item

Changes the equipped status of an item in an actor's inventory.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| equipped | boolean | ✓ | body, query | Whether the item should be equipped (true) or unequipped (false) |
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor (optional if selected is true) |
| actorName | string |  | body, query | Name of the actor (optional if actorUuid provided) |
| itemUuid | string |  | body, query | UUID of the item (optional if itemName provided) |
| itemName | string |  | body, query | Name of the item (optional if itemUuid provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Updated equipment status

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/equip-item"
  parameters={[{"name":"equipped","type":"boolean","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"actorName","type":"string","required":false,"source":"body"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/attune-item

Attune or unattune an item

Changes the attunement status of a magic item in an actor's inventory.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| attuned | boolean | ✓ | body, query | Whether the item should be attuned (true) or unattuned (false) |
| clientId | string |  | query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor (optional if selected is true) |
| actorName | string |  | body, query | Name of the actor (optional if actorUuid provided) |
| itemUuid | string |  | body, query | UUID of the item (optional if itemName provided) |
| itemName | string |  | body, query | Name of the item (optional if itemUuid provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Updated attunement status

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/attune-item"
  parameters={[{"name":"attuned","type":"boolean","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"actorName","type":"string","required":false,"source":"body"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/transfer-currency

Transfer currency between actors

Moves currency from one actor to another. Validates that the source actor has sufficient funds before transferring.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| currency | object | ✓ | body, query | Currency amounts to transfer, e.g. pp, gp, ep, sp, cp denomination keys with numeric values |
| clientId | string |  | query | Client ID for the Foundry world |
| sourceActorUuid | string |  | body, query | UUID of the source actor (optional if sourceActorName provided) |
| sourceActorName | string |  | body, query | Name of the source actor |
| targetActorUuid | string |  | body, query | UUID of the target actor (optional if targetActorName provided) |
| targetActorName | string |  | body, query | Name of the target actor |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Transfer result with updated balances

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/transfer-currency"
  parameters={[{"name":"currency","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sourceActorUuid","type":"string","required":false,"source":"body"},{"name":"sourceActorName","type":"string","required":false,"source":"body"},{"name":"targetActorUuid","type":"string","required":false,"source":"body"},{"name":"targetActorName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/modify-currency

Modify currency balance for a single actor (delta-based, not a transfer between actors)

Adds or removes currency from an actor's wallet. Use a negative amount to remove currency.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| currency | string | ✓ | body, query | Currency denomination to modify (pp, gp, ep, sp, cp) |
| amount | number | ✓ | body, query | Amount to add (positive) or remove (negative) |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the currency modification

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/modify-currency"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"currency","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/prepare-spell

Prepare or unprepare a spell for an actor

Toggles a spell's prepared state. Only applicable to spellcaster classes that prepare spells.

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| spellName | string | ✓ | body, query | Name of the spell to prepare or unprepare |
| prepared | boolean | ✓ | body, query | True to prepare the spell, false to unprepare it |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the prepare spell operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/prepare-spell"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"spellName","type":"string","required":true,"source":"body"},{"name":"prepared","type":"boolean","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/use-ability

Use an ability

Activates a specific ability for an actor, optionally targeting another entity

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| clientId | string |  | query | Client ID for the Foundry world |
| abilityUuid | string |  | body, query | The UUID of the specific ability (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the ability (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the use ability operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/use-ability"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/use-feature

Use a feature

Activates a specific feature for an actor, optionally targeting another entity

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| clientId | string |  | query | Client ID for the Foundry world |
| abilityUuid | string |  | body, query | The UUID of the specific ability (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the ability (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the use feature operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/use-feature"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/use-spell

Use a spell

Casts a specific spell for an actor, optionally targeting another entity

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| clientId | string |  | query | Client ID for the Foundry world |
| abilityUuid | string |  | body, query | The UUID of the specific ability (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the ability (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the use spell operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/use-spell"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /dnd5e/use-item

Use an item

Uses a specific item for an actor, optionally targeting another entity

**Required scope:** `dnd5e`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| clientId | string |  | query | Client ID for the Foundry world |
| abilityUuid | string |  | body, query | The UUID of the specific ability (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the ability (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the use item operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/use-item"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

