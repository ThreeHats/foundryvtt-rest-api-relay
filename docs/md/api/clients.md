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
const baseUrl = 'http://localhost:3011';
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
curl -X GET 'http://localhost:3011/clients' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
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
  const baseUrl = 'http://localhost:3011';
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
  3011 ➡️ port
  🔤/clients🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /clients HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
      "clientId": "qsl-integration-test",
      "instanceId": "local",
      "lastSeen": 1782956903871,
      "connectedSince": 1782956903870,
      "worldId": "test-world",
      "worldTitle": "Test World",
      "foundryVersion": "14.363",
      "systemId": "dnd5e",
      "systemTitle": "Dungeons & Dragons Fifth Edition",
      "systemVersion": "5.0.4",
      "publicUrl": "http://foundry:30000",
      "ipAddress": "172.24.0.3:37168",
      "tokenName": "headless session 2026-07-02 01:48",
      "isOnline": true
    }
  ],
  "total": 1
}
```


