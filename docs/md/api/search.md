---
tag: search
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


# search

## GET /search

Search entities This endpoint allows searching for entities in the Foundry world based on a query string. Requires Quick Insert module to be installed and enabled.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| query | string | ✓ | query | Search query string |
| filter | string |  | query | Filter to apply (simple: filter="Actor", property-based: filter="key:value,key2:value2") |

### Returns

**object** - Search results containing matching entities

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/search';
const params = {
  clientId: 'your-client-id',
  query: 'test-item',
  filter: 'documentType:Item,subType:base'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

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
curl -X GET 'http://localhost:3010/search?clientId=your-client-id&query=test-item&filter=documentType%3AItem%2CsubType%3Abase' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/search'
params = {
    'clientId': 'your-client-id',
    'query': 'test-item',
    'filter': 'documentType:Item,subType:base'
}
url = f'{base_url}{path}'

response = requests.get(
    url,
    params=params,
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
  const path = '/search';
  const params = {
    clientId: 'your-client-id',
    query: 'test-item',
    filter: 'documentType:Item,subType:base'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

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
  🔤/search🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤query=test-item🔤 ➡️ query
  🔤filter=documentType:Item,subType:base🔤 ➡️ filter
  🔤?🧲clientId🧲&🧲query🧲&🧲filter🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /search🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "search_1765635979813",
  "clientId": "your-client-id",
  "type": "search-result",
  "query": "test-item",
  "filter": "documentType:Item,subType:base",
  "results": [
    {
      "documentType": "Item",
      "id": "ifzt7D41CkQeFgm4",
      "name": "test-item",
      "subType": "base",
      "uuid": "Item.ifzt7D41CkQeFgm4",
      "icon": "<i class=\"fas fa-suitcase entity-icon\"></i>",
      "journalLink": "@UUID[Item.ifzt7D41CkQeFgm4]{test-item}",
      "tagline": "Items Directory",
      "formattedMatch": "<strong>test-item</strong>",
      "resultType": "EntitySearchItem"
    }
  ]
}
```

