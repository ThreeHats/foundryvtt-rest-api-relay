---
tag: clients
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


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
      "id": "foundry-testing-r6bXhB7k9cXa3cif",
      "instanceId": "local",
      "lastSeen": 1774367571744,
      "connectedSince": 1774367571742,
      "worldId": "testing",
      "worldTitle": "testing",
      "foundryVersion": "13.348",
      "systemId": "dnd5e",
      "systemTitle": "Dungeons & Dragons Fifth Edition",
      "systemVersion": "5.0.4"
    },
    {
      "id": "foundry-rest-api-fCfNJPT9Atc26yyv",
      "instanceId": "local",
      "lastSeen": 1774367578644,
      "connectedSince": 1774367578615,
      "worldId": "rest-api",
      "worldTitle": "rest-api",
      "foundryVersion": "12.331",
      "systemId": "dnd5e",
      "systemTitle": "Dungeons & Dragons Fifth Edition",
      "systemVersion": "4.3.8",
      "customName": "v12-test"
    }
  ],
  "total": 2
}
```

