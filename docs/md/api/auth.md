---
tag: auth
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Auth

## POST /auth/key-request

Request a scoped API key

Initiates the device-flow key request. The response contains an `approvalUrl` the user must visit to review and approve the requested scopes. Poll `GET /auth/key-request/{code}/status` until `status` is `approved` or `denied`.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| appName | string | ✓ | body | Name of the application requesting access |
| scopes | array | ✓ | body | List of permission scopes the key requires |
| appDescription | string |  | body | Short description of what the application does |
| appUrl | string |  | body | Homepage or docs URL for the application |
| callbackUrl | string |  | body | URL to redirect to after approval (web flow) |
| clientIds | array |  | body | Foundry client IDs to restrict the key to |
| suggestedMonthlyLimit | number |  | body | Suggested per-key monthly request cap (user may override) |
| suggestedExpiry | string |  | body | Suggested expiry date ISO 8601 (user may override) |

### Returns

**object** - `code`, `approvalUrl`, `expiresIn`, `expiresAt`

### Try It Out

<ApiTester
  method="POST"
  path="/auth/key-request"
  parameters={[{"name":"appName","type":"string","required":true,"source":"body"},{"name":"scopes","type":"array","required":true,"source":"body"},{"name":"appDescription","type":"string","required":false,"source":"body"},{"name":"appUrl","type":"string","required":false,"source":"body"},{"name":"callbackUrl","type":"string","required":false,"source":"body"},{"name":"clientIds","type":"array","required":false,"source":"body"},{"name":"suggestedMonthlyLimit","type":"number","required":false,"source":"body"},{"name":"suggestedExpiry","type":"string","required":false,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/key-request';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "appName": "Test Discord Bot",
      "appDescription": "A test integration",
      "scopes": [
        "entity:read",
        "roll:read",
        "chat:read"
      ]
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/key-request' \
  -H "Content-Type: application/json" \
  -d '{"appName":"Test Discord Bot","appDescription":"A test integration","scopes":["entity:read","roll:read","chat:read"]}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/key-request'
url = f'{base_url}{path}'

response = requests.post(
    url,
    json={
      "appName": "Test Discord Bot",
      "appDescription": "A test integration",
      "scopes": [
        "entity:read",
        "roll:read",
        "chat:read"
      ]
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
  const path = '/auth/key-request';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "appName": "Test Discord Bot",
        "appDescription": "A test integration",
        "scopes": [
          "entity:read",
          "roll:read",
          "chat:read"
        ]
      }
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
  🔤/auth/key-request🔤 ➡️ path

  💭 Request body
  🔤{"appName":"Test Discord Bot","appDescription":"A test integration","scopes":["entity:read","roll:read","chat:read"]}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/key-request HTTP/1.1❌r❌nHost: localhost:3010❌r❌nContent-Type: application/json❌r❌nContent-Length: 117❌r❌n❌r❌n{"appName":"Test Discord Bot","appDescription":"A test integration","scopes":["entity:read","roll:read","chat:read"]}🔤 ➡️ request

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

**Status:** 201

```json
{
  "approvalUrl": "http://localhost:3010/approve/F4JQJX",
  "code": "F4JQJX",
  "expiresAt": "2026-05-14T15:28:09-05:00",
  "expiresIn": 600
}
```


---

## GET /auth/key-request/:code/status

Poll key request status

Returns the current status of a pending key request. When `status` is `approved`, the response includes the newly created `key`. Once the key has been retrieved, the code is invalidated.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| code | string | ✓ | params | The code returned by POST /auth/key-request |

### Returns

**object** - `status` (`pending` | `approved` | `denied` | `expired`), `key` (when approved)

### Try It Out

<ApiTester
  method="GET"
  path="/auth/key-request/:code/status"
  parameters={[{"name":"code","type":"string","required":true,"source":"params"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/key-request/F4JQJX/status';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'GET'
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X GET 'http://localhost:3010/auth/key-request/F4JQJX/status'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/key-request/F4JQJX/status'
url = f'{base_url}{path}'

response = requests.get(
    url
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
  const path = '/auth/key-request/F4JQJX/status';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'get',
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
  🔤/auth/key-request/F4JQJX/status🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /auth/key-request/F4JQJX/status HTTP/1.1❌r❌nHost: localhost:3010❌r❌n❌r❌n🔤 ➡️ request

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
  "apiKey": "your-api-key-here",
  "clientIds": null,
  "scopes": [
    "entity:read",
    "roll:read"
  ],
  "status": "approved"
}
```


