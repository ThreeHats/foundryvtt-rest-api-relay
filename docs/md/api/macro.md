---
tag: macro
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Macro

## GET /macros

Get all macros

Retrieves a list of all macros available in the Foundry world.

**Required scope:** `macro:list`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - An array of macros with details

### Try It Out

<ApiTester
  method="GET"
  path="/macros"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /macro/:uuid/execute

Execute a macro by UUID

Executes a specific macro in the Foundry world by its UUID.

**Required scope:** `macro:execute`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | params | UUID of the macro to execute |
| clientId | string |  | query | Client ID for the Foundry world |
| args | object |  | body | Optional arguments to pass to the macro execution |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the macro execution

### Try It Out

<ApiTester
  method="POST"
  path="/macro/:uuid/execute"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"params"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"args","type":"object","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

