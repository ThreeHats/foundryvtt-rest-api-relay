---
tag: user
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# User

## GET /users

List all Foundry users

Retrieves a list of all users configured in the Foundry VTT world, including their roles and online status. This is a GM-only operation.

**Required scope:** `user:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of user objects with id, name, role, isGM, active, color, avatar, and character fields

### Try It Out

<ApiTester
  method="GET"
  path="/users"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /user

Get a single Foundry user

Retrieves a single user by their ID or name. This is a GM-only operation.

**Required scope:** `user:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| id | string |  | query | ID of the user to retrieve |
| name | string |  | query | Name of the user to retrieve (alternative to id) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - User object with id, name, role, isGM, active, color, avatar, and character fields

### Try It Out

<ApiTester
  method="GET"
  path="/user"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"id","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /user

Create a new Foundry user

Creates a new user in the Foundry VTT world with the specified name, role, and optional password. This is a GM-only operation.

**Required scope:** `user:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| name | string | ✓ | body | Username for the new user |
| clientId | string |  | query | Client ID for the Foundry world |
| role | number |  | body | User role: 0=None, 1=Player, 2=Trusted, 3=Assistant, 4=GM (default: 1) |
| password | string |  | body | Password for the new user |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The created user object

### Try It Out

<ApiTester
  method="POST"
  path="/user"
  parameters={[{"name":"name","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"role","type":"number","required":false,"source":"body"},{"name":"password","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## PUT /user

Update an existing Foundry user

Updates fields on an existing user. Identify the user by id or name, then pass the fields to update in the data object. Cannot demote the last GM user. This is a GM-only operation.

**Required scope:** `user:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| data | object | ✓ | body | Object containing user fields to update (name, role, password, color, avatar, character) |
| clientId | string |  | query | Client ID for the Foundry world |
| id | string |  | body, query | ID of the user to update |
| name | string |  | body, query | Name of the user to update (alternative to id) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The updated user object

### Try It Out

<ApiTester
  method="PUT"
  path="/user"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"id","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## DELETE /user

Delete a Foundry user

Permanently deletes a user from the Foundry VTT world. Cannot delete yourself or the last GM user. This is a GM-only operation.

**Required scope:** `user:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| id | string |  | query | ID of the user to delete |
| name | string |  | query | Name of the user to delete (alternative to id) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Deletion result

### Try It Out

<ApiTester
  method="DELETE"
  path="/user"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"id","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

