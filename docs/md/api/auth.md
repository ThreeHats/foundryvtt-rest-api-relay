---
tag: auth
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Auth

## POST /auth/key-request

Request a scoped API key

Initiates a key request that a user must approve in the dashboard. Supports two flows:

**Device flow (no `callbackUrl`)** — for CLI tools, scripts, and desktop apps. The response contains an `approvalUrl` to direct the user to. Poll `GET /auth/key-request/{code}/status` until `status` is `approved`, then read `apiKey` from the response.

**Web flow (`callbackUrl` provided)** — for web apps that can receive an HTTP redirect. After the user approves, the dashboard redirects to your `callbackUrl` with a `code` query parameter containing the exchange code. POST that code to `POST /auth/key-request/exchange` to retrieve the API key.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| appName | string | ✓ | body | Name of the application requesting access |
| scopes | array | ✓ | body | List of permission scopes the key requires |
| appDescription | string |  | body | Short description of what the application does |
| appUrl | string |  | body | Homepage or docs URL for the application |
| callbackUrl | string |  | body | Enables web flow: URL the dashboard redirects to after approval, with a `code` query parameter containing the exchange code |
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
  "approvalUrl": "http://localhost:3010/approve/R76CTM",
  "code": "R76CTM",
  "expiresAt": "2026-05-15T21:04:05-05:00",
  "expiresIn": 600
}
```


---

## GET /auth/key-request/:code/status

Poll key request status

Returns the current status of a pending key request. When `status` is `approved`, the response includes the newly created `apiKey`. Once the key has been retrieved, the status becomes `exchanged` and the key is no longer returned.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| code | string | ✓ | params | The code returned by POST /auth/key-request |

### Returns

**object** - `status` (`pending` | `approved` | `denied` | `expired` | `exchanged`), `apiKey`, `scopes`, `clientIds` (when approved)

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
const path = '/auth/key-request/R76CTM/status';
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
curl -X GET 'http://localhost:3010/auth/key-request/R76CTM/status'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/key-request/R76CTM/status'
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
  const path = '/auth/key-request/R76CTM/status';
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
  🔤/auth/key-request/R76CTM/status🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /auth/key-request/R76CTM/status HTTP/1.1❌r❌nHost: localhost:3010❌r❌n❌r❌n🔤 ➡️ request

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


---

## POST /auth/key-request/exchange

Exchange approval code for API key (web flow)

Web flow only. After the user approves the request in the dashboard, the relay redirects to the `callbackUrl` with a one-time `code` query parameter. POST that code here to receive the API key. The code is single-use — a second call returns 410.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| code | string | ✓ | body | The exchange code delivered to your callbackUrl |

### Returns

**object** - `apiKey`, `scopes`, `clientIds`

### Try It Out

<ApiTester
  method="POST"
  path="/auth/key-request/exchange"
  parameters={[{"name":"code","type":"string","required":true,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/key-request/exchange';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "code": "your-exchange-code-here"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/key-request/exchange' \
  -H "Content-Type: application/json" \
  -d '{"code": "your-exchange-code-here"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/key-request/exchange'
url = f'{base_url}{path}'

response = requests.post(
    url,
    json={
      "code": "your-exchange-code-here"
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
  const path = '/auth/key-request/exchange';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "code": "your-exchange-code-here"
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
  🔤/auth/key-request/exchange🔤 ➡️ path

  💭 Request body
  🔤{"code": "your-exchange-code-here"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/key-request/exchange HTTP/1.1❌r❌nHost: localhost:3010❌r❌nContent-Type: application/json❌r❌nContent-Length: 43❌r❌n❌r❌n{"code": "your-exchange-code-here"}🔤 ➡️ request

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
    "entity:read"
  ]
}
```


