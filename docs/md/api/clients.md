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

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/clients';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'GET',
  headers: {
    'x-api-key': 'your-api-key-here'
  }
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X GET 'http://localhost:3010/clients' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/clients'
url = f'{base_url}{path}'

response = requests.get(
    url,
    headers={
        'x-api-key': 'your-api-key-here'
    }
)
data = response.json()
print(data)
```

</TabItem>
<TabItem value="typescript" label="TypeScript">

```typescript
import axios from 'axios';

(async () => {
  const baseUrl = 'http://localhost:3010';
  const path = '/clients';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'get',
    headers: {
      'x-api-key': 'your-api-key-here'
    },
    url
  });
  const data = response.data;
  console.log(data);
})();
```

</TabItem>
<TabItem value="emojicode" label="Emojicode">

```emojicode
📦 sockets 🏠

💭 Emojicode HTTP Client
💭 Compile: emojicodec example.🍇 -o example
💭 Run: ./example

🏁 🍇
  💭 Connection settings
  🔤localhost🔤 ➡️ host
  3010 ➡️ port
  🔤/clients🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /clients HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

  💭 Connect and send
  🍺 🆕📞 host port❗ ➡️ socket
  🍺 💬 socket 📇 request❗❗
  
  💭 Read and print response
  🍺 👂 socket 4096❗ ➡️ data
  😀 🍺 🔡 data❗❗
  
  💭 Close socket
  🚪 socket❗
🍉
```

</TabItem>
</Tabs>

#### Response

**Status:** 200

```json
{
  "clients": [
    {
      "clientId": "fvtt_d3db244ce1d8af0b",
      "instanceId": "local",
      "lastSeen": 1780520731407,
      "connectedSince": 1780520731405,
      "worldId": "rest-api",
      "worldTitle": "rest-api",
      "foundryVersion": "12.331",
      "systemId": "dnd5e",
      "systemTitle": "Dungeons & Dragons Fifth Edition",
      "systemVersion": "4.3.8",
      "customName": "my-cool-server",
      "publicUrl": "http://localhost:30012",
      "ipAddress": "[::1]:51316",
      "tokenName": "headless session 2026-06-03 16:05",
      "isOnline": true
    },
    {
      "clientId": "fvtt_099ad17ea199e7e3",
      "instanceId": "local",
      "lastSeen": 1780520736804,
      "connectedSince": 1780520734799,
      "worldId": "testing",
      "worldTitle": "testing",
      "foundryVersion": "13.348",
      "systemId": "dnd5e",
      "systemTitle": "Dungeons & Dragons Fifth Edition",
      "systemVersion": "5.0.4",
      "publicUrl": "http://localhost:30013",
      "ipAddress": "[::1]:51324",
      "tokenName": "headless session 2026-06-03 16:05",
      "isOnline": true
    },
    {
      "clientId": "fvtt_20b2e79d7d679516",
      "instanceId": "local",
      "lastSeen": 1780520738273,
      "connectedSince": 1780520738272,
      "worldId": "5e",
      "worldTitle": "5e",
      "foundryVersion": "14.360",
      "systemId": "dnd5e",
      "systemTitle": "Dungeons & Dragons Fifth Edition",
      "systemVersion": "5.3.1",
      "publicUrl": "http://localhost:30014",
      "ipAddress": "[::1]:51332",
      "tokenName": "headless session 2026-06-03 16:05",
      "isOnline": true
    }
  ],
  "total": 3
}
```


