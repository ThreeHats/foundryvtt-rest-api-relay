---
tag: canvas
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Canvas

## GET /canvas/:documentType

Get canvas embedded documents

**Required scope:** `canvas:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | query | Scene ID to query (defaults to the active scene) |
| documentId | string |  | query | Specific document ID to retrieve |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of embedded documents

### Try It Out

<ApiTester
  method="GET"
  path="/canvas/:documentType"
  parameters={[{"name":"documentType","type":"string","required":true,"source":"params"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"documentId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /measure-distance

Measure the distance between two points or tokens

Calculates the distance between two positions on the canvas, respecting the grid type and measurement rules. Points can be specified as coordinates or by referencing tokens by UUID or name.

**Required scope:** `canvas:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| originX | number |  | body, query | Origin x coordinate (optional if originUuid/originName provided) |
| originY | number |  | body, query | Origin y coordinate |
| targetX | number |  | body, query | Target x coordinate (optional if targetUuid/targetName provided) |
| targetY | number |  | body, query | Target y coordinate |
| originUuid | string |  | body, query | UUID of the origin token |
| originName | string |  | body, query | Name of the origin token |
| targetUuid | string |  | body, query | UUID of the target token |
| targetName | string |  | body, query | Name of the target token |
| sceneId | string |  | body, query | Scene ID (defaults to active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Distance measurement including units and grid spaces

### Try It Out

<ApiTester
  method="GET"
  path="/measure-distance"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"originX","type":"number","required":false,"source":"body"},{"name":"originY","type":"number","required":false,"source":"body"},{"name":"targetX","type":"number","required":false,"source":"body"},{"name":"targetY","type":"number","required":false,"source":"body"},{"name":"originUuid","type":"string","required":false,"source":"body"},{"name":"originName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /canvas/:documentType

Create canvas embedded document(s)

**Required scope:** `canvas:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| data | object | ✓ | body | Document data object or array of objects to create |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | Scene ID to create in (defaults to the active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Created document(s)

### Try It Out

<ApiTester
  method="POST"
  path="/canvas/:documentType"
  parameters={[{"name":"documentType","type":"string","required":true,"source":"params"},{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## PUT /canvas/:documentType

Update a canvas embedded document

**Required scope:** `canvas:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| documentId | string | ✓ | body, query | ID of the document to update |
| data | object | ✓ | body | Object containing the fields to update |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | Scene ID containing the document (defaults to the active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Updated document

### Try It Out

<ApiTester
  method="PUT"
  path="/canvas/:documentType"
  parameters={[{"name":"documentType","type":"string","required":true,"source":"params"},{"name":"documentId","type":"string","required":true,"source":"body"},{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## DELETE /canvas/:documentType

Delete a canvas embedded document

**Required scope:** `canvas:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| documentId | string | ✓ | query | ID of the document to delete |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | query | Scene ID containing the document (defaults to the active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Deletion result

### Try It Out

<ApiTester
  method="DELETE"
  path="/canvas/:documentType"
  parameters={[{"name":"documentType","type":"string","required":true,"source":"params"},{"name":"documentId","type":"string","required":true,"source":"query"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /move-token

Move a token to specific coordinates

Moves a token on the canvas to the specified x,y position, optionally animating through waypoints. Token can be identified by UUID or name.

**Required scope:** `canvas:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| x | number | ✓ | body, query | Target x coordinate |
| y | number | ✓ | body, query | Target y coordinate |
| clientId | string |  | query | Client ID for the Foundry world |
| uuid | string |  | body, query | UUID of the token to move (optional if name provided) |
| name | string |  | body, query | Name of the token to move (optional if uuid provided) |
| waypoints | array |  | body, query | Array of waypoint objects with x and y coordinates to animate through before reaching final position |
| animate | boolean |  | body, query | Whether to animate the movement (default: true) |
| sceneId | string |  | body, query | Scene ID (defaults to active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the token movement including new position

### Try It Out

<ApiTester
  method="POST"
  path="/move-token"
  parameters={[{"name":"x","type":"number","required":true,"source":"body"},{"name":"y","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"waypoints","type":"array","required":false,"source":"body"},{"name":"animate","type":"boolean","required":false,"source":"body"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

