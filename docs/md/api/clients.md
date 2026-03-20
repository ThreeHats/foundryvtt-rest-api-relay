---
tag: clients
---

import ApiTester from '@site/src/components/ApiTester';

# clients

## GET /clients

Get all connected clients for the authenticated API key Returns a list of all currently connected Foundry VTT clients associated with the provided API key, including their connection details and world information.

### Returns

**object** - Object containing total count and array of connected client details

### Try It Out

<ApiTester
  method="GET"
  path="/clients"
  parameters={[]}
/>

