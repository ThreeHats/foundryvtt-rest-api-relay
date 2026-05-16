---
tag: filesystem
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# FileSystem

## GET /file-system

Get file system structure

**Required scope:** `file:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| path | string |  | query | The path to retrieve (relative to source) |
| source | string |  | query | The source directory to use (data, systems, modules, etc.) |
| recursive | boolean |  | query | Whether to recursively list all subdirectories |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - File system structure with files and directories

### Try It Out

<ApiTester
  method="GET"
  path="/file-system"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"path","type":"string","required":false,"source":"query"},{"name":"source","type":"string","required":false,"source":"query"},{"name":"recursive","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /download

Download a file from Foundry's file system

**Required scope:** `file:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| path | string |  | query | The full path to the file to download |
| source | string |  | query | The source directory to use (data, systems, modules, etc.) |
| format | string |  | query | The format to return the file in (binary, base64) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - File contents in the requested format

### Try It Out

<ApiTester
  method="GET"
  path="/download"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"path","type":"string","required":false,"source":"query"},{"name":"source","type":"string","required":false,"source":"query"},{"name":"format","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /upload

Upload a file to Foundry's file system (handles both base64 and binary data)

**Required scope:** `file:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| path | string |  | query, body | The directory path to upload to |
| filename | string |  | query, body | The filename to save as |
| source | string |  | query, body | The source directory to use (data, systems, modules, etc.) |
| mimeType | string |  | query, body | The MIME type of the file |
| overwrite | boolean |  | query, body | Whether to overwrite an existing file |
| fileData | string |  | body | Base64 encoded file data (if sending as JSON) 250MB limit |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the file upload operation

### Try It Out

<ApiTester
  method="POST"
  path="/upload"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"path","type":"string","required":false,"source":"query"},{"name":"filename","type":"string","required":false,"source":"query"},{"name":"source","type":"string","required":false,"source":"query"},{"name":"mimeType","type":"string","required":false,"source":"query"},{"name":"overwrite","type":"boolean","required":false,"source":"query"},{"name":"fileData","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

