---
tag: scene
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Scene

## GET /scene

Get scene(s)

Retrieves one or more scenes by ID, name, active status, viewed status, or all.

**Required scope:** `scene:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | query | ID of a specific scene to retrieve |
| name | string |  | query | Name of the scene to retrieve |
| active | boolean |  | query | Set to true to get the currently active scene |
| viewed | boolean |  | query | Set to true to get the currently viewed scene |
| all | boolean |  | query | Set to true to get all scenes |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Scene data

### Try It Out

<ApiTester
  method="GET"
  path="/scene"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"active","type":"boolean","required":false,"source":"query"},{"name":"viewed","type":"boolean","required":false,"source":"query"},{"name":"all","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /scene/image/raw

Get the raw background image of a scene

Returns the scene's background image file without any tokens, lights, or other canvas elements rendered on it.

**Required scope:** `scene:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | Scene ID (defaults to viewed/active scene) |
| active | boolean |  | body, query | If true, explicitly use the player-facing active scene instead of the viewed scene |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**binary** - The raw scene background image

### Try It Out

<ApiTester
  method="GET"
  path="/scene/image/raw"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"active","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /scene

Create a new scene

**Required scope:** `scene:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| data | object | ✓ | body | Scene data object (name, width, height, grid, etc.) |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Created scene data

### Try It Out

<ApiTester
  method="POST"
  path="/scene"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## PUT /scene

Update an existing scene

**Required scope:** `scene:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| data | object | ✓ | body | Object containing the scene fields to update |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | ID of the scene to update |
| name | string |  | body, query | Name of the scene to update (alternative to sceneId) |
| active | boolean |  | body, query | Set to true to target the active scene |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Updated scene data

### Try It Out

<ApiTester
  method="PUT"
  path="/scene"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"active","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## DELETE /scene

Delete a scene

**Required scope:** `scene:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | query | ID of the scene to delete |
| name | string |  | query | Name of the scene to delete (alternative to sceneId) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Deletion result

### Try It Out

<ApiTester
  method="DELETE"
  path="/scene"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /switch-scene

Switch the active scene

**Required scope:** `scene:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | ID of the scene to activate |
| name | string |  | body, query | Name of the scene to activate (alternative to sceneId) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the scene switch

### Try It Out

<ApiTester
  method="POST"
  path="/switch-scene"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /scene/image

Get a rendered screenshot of a scene

Captures the full rendered canvas of a scene including all visible layers (tokens, lights, walls, etc.) as an image. The scene can be specified by ID or defaults to the active scene.

**Required scope:** `scene:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| sceneId | string |  | query | Scene ID (defaults to viewed/active scene) |
| active | boolean |  | query | If true, explicitly use the player-facing active scene instead of the viewed scene |
| clientId | string |  | query | Client ID for the Foundry world |
| format | string |  | query | Image format: png or jpeg (default: png) |
| quality | number |  | query | Image quality 0-1 for JPEG (default: 0.9) |
| viewport | boolean |  | query | If true, capture exactly what the browser currently shows instead of the full scene |
| width | number |  | query | Output image width in pixels (default: scene width) |
| height | number |  | query | Output image height in pixels (default: scene height) |
| showGrid | boolean |  | query | Include grid lines in capture (default: false) |
| hideOverlays | boolean |  | query | Hide fog of war, weather, vision, and UI overlays (default: false) |
| userId | string |  | query | Foundry user ID or username |

### Returns

**binary** - The scene screenshot as an image

### Try It Out

<ApiTester
  method="GET"
  path="/scene/image"
  parameters={[{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"active","type":"boolean","required":false,"source":"query"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"format","type":"string","required":false,"source":"query"},{"name":"quality","type":"number","required":false,"source":"query"},{"name":"viewport","type":"boolean","required":false,"source":"query"},{"name":"width","type":"number","required":false,"source":"query"},{"name":"height","type":"number","required":false,"source":"query"},{"name":"showGrid","type":"boolean","required":false,"source":"query"},{"name":"hideOverlays","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

