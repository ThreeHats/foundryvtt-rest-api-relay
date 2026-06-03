---
tag: structure
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Structure

## GET /structure

Get the structure of the Foundry world

Retrieves the folder and compendium structure for the specified Foundry world.

**Required scope:** `structure:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| includeEntityData | boolean |  | query | Whether to include full entity data or just UUIDs and names |
| path | string |  | query | Path to read structure from (null = root) |
| recursive | boolean |  | query | Whether to read down the folder tree |
| recursiveDepth | number |  | query | Depth to recurse into folders (default 5) |
| types | string |  | query | Types to return (Scene/Actor/Item/JournalEntry/RollTable/Cards/Macro/Playlist), can be comma-separated or JSON array |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The folder and compendium structure

### Try It Out

<ApiTester
  method="GET"
  path="/structure"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"includeEntityData","type":"boolean","required":false,"source":"query"},{"name":"path","type":"string","required":false,"source":"query"},{"name":"recursive","type":"boolean","required":false,"source":"query"},{"name":"recursiveDepth","type":"number","required":false,"source":"query"},{"name":"types","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /get-folder

Get a specific folder by name

**Required scope:** `structure:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| name | string | ✓ | body, query | Name of the folder to retrieve |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The folder information and its contents

### Try It Out

<ApiTester
  method="GET"
  path="/get-folder"
  parameters={[{"name":"name","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /create-folder

Create a new folder

**Required scope:** `structure:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| name | string | ✓ | body, query | Name of the new folder |
| folderType | string | ✓ | body, query | Type of folder (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist) |
| clientId | string |  | query | Client ID for the Foundry world |
| parentFolderId | string |  | body, query | ID of the parent folder (optional for root level) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The created folder information

### Try It Out

<ApiTester
  method="POST"
  path="/create-folder"
  parameters={[{"name":"name","type":"string","required":true,"source":"body"},{"name":"folderType","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"parentFolderId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## DELETE /delete-folder

Delete a folder

**Required scope:** `structure:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| folderId | string | ✓ | body, query | ID of the folder to delete |
| clientId | string |  | query | Client ID for the Foundry world |
| deleteAll | boolean |  | body, query | Whether to delete all entities in the folder |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Confirmation of deletion

### Try It Out

<ApiTester
  method="DELETE"
  path="/delete-folder"
  parameters={[{"name":"folderId","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"deleteAll","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /contents/:path

This route is deprecated

Use /structure with the path query parameter instead.

**Required scope:** `structure:read`

### Returns

**object** - Error message directing to use /structure endpoint

### Try It Out

<ApiTester
  method="GET"
  path="/contents/:path"
  parameters={[]}
/>

