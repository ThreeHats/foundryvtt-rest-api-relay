---
tag: search
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Search

## GET /search

Search entities

This endpoint allows searching for entities in the Foundry world based on a query string. Search world entities and compendiums using the native built-in search engine. No third-party modules required. Results are ranked by relevance: exact match, prefix match, substring match, word-prefix match, and subsequence match.

**Required scope:** `search`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| query | string |  | query | Search query string (omit to browse all entities matching filter) |
| filter | string |  | query | Filter string — simple: filter="Actor"; compound: filter="documentType:Item,subType:weapon". Supported keys: documentType, subType, folder, package, resultType |
| excludeCompendiums | boolean |  | query | Exclude compendium entries from results (default: false — compendiums are included by default) |
| limit | number |  | query | Maximum number of results to return (default: 200, max: 500) |
| minified | boolean |  | query | Return minimal fields only — uuid, id, name, img, documentType (default: false) |
| ownedByUserId | string |  | query | Filter results to only documents the specified Foundry user (ID or username) has Owner permission on |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Search results containing matching entities

### Try It Out

<ApiTester
  method="GET"
  path="/search"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"query","type":"string","required":false,"source":"query"},{"name":"filter","type":"string","required":false,"source":"query"},{"name":"excludeCompendiums","type":"boolean","required":false,"source":"query"},{"name":"limit","type":"number","required":false,"source":"query"},{"name":"minified","type":"boolean","required":false,"source":"query"},{"name":"ownedByUserId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

