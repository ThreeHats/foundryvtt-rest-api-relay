---
tag: clients
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Clients

## GET /clients

Get all connected clients for the authenticated API key

Returns a list of all currently connected Foundry VTT clients associated with the provided API key, including their connection details and world information.

**Required scope:** `clients:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| x-api-key | string |  | header | API key for authentication |

### Returns

**array** - Object containing total count and array of connected client details

### Try It Out

<ApiTester
  method="GET"
  path="/clients"
  parameters={[{"name":"x-api-key","type":"string","required":false,"source":"header"}]}
/>

