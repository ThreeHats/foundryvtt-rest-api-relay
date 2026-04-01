---
tag: sheet
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Sheet

## GET /sheet

Get actor sheet HTML

This endpoint retrieves the HTML for an actor sheet based on the provided UUID or selected actor. Only works on Foundry version 12.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string |  | query | The UUID of the entity to get the sheet for |
| selected | boolean |  | query | Whether to get the sheet for the selected entity |
| actor | boolean |  | query | Whether to get the sheet for the selected token's actor if selected is true |
| clientId | string |  | query | Client ID for the Foundry world |
| format | string |  | query | The format to return the sheet in (html, json) |
| scale | number |  | query | The initial scale of the sheet |
| tab | number |  | query | The active tab index to open |
| darkMode | boolean |  | query | Whether to use dark mode for the sheet |
| userId | string |  | query | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The sheet HTML or data depending on format requested

### Try It Out

<ApiTester
  method="GET"
  path="/sheet"
  parameters={[{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"actor","type":"boolean","required":false,"source":"query"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"format","type":"string","required":false,"source":"query"},{"name":"scale","type":"number","required":false,"source":"query"},{"name":"tab","type":"number","required":false,"source":"query"},{"name":"darkMode","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/sheet';
const params = {
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
  uuid: 'Actor.pxZTVHItjx6GgPgC',
  format: 'json'
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
curl -X GET 'http://localhost:3010/sheet?clientId=foundry-testing-r6bXhB7k9cXa3cif&uuid=Actor.pxZTVHItjx6GgPgC&format=json' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/sheet'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif',
    'uuid': 'Actor.pxZTVHItjx6GgPgC',
    'format': 'json'
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
  const path = '/sheet';
  const params = {
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
    uuid: 'Actor.pxZTVHItjx6GgPgC',
    format: 'json'
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
  🔤/sheet🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤uuid=Actor.pxZTVHItjx6GgPgC🔤 ➡️ uuid
  🔤format=json🔤 ➡️ format
  🔤?🧲clientId🧲&🧲uuid🧲&🧲format🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /sheet🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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

**Status:** 400

```json
{
  "clientId": "foundry-testing-r6bXhB7k9cXa3cif",
  "error": "This endpoint is only supported in Foundry VTT version 12",
  "requestId": "get-sheet_1775068878248"
}
```


