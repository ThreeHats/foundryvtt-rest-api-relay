---
tag: effects
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Effects

## GET /effects

Get all active effects on an actor or token

Returns the collection of ActiveEffect documents currently applied to the specified actor or token.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to query |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of active effects

### Try It Out

<ApiTester
  method="GET"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/effects';
const params = {
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
  uuid: 'Actor.pxZTVHItjx6GgPgC'
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
curl -X GET 'http://localhost:3010/effects?clientId=foundry-testing-r6bXhB7k9cXa3cif&uuid=Actor.pxZTVHItjx6GgPgC' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/effects'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif',
    'uuid': 'Actor.pxZTVHItjx6GgPgC'
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
  const path = '/effects';
  const params = {
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
    uuid: 'Actor.pxZTVHItjx6GgPgC'
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
  🔤/effects🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤uuid=Actor.pxZTVHItjx6GgPgC🔤 ➡️ uuid
  🔤?🧲clientId🧲&🧲uuid🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /effects🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "get-effects-result",
  "requestId": "get-effects_1775068881867",
  "data": {
    "uuid": "Actor.pxZTVHItjx6GgPgC",
    "effects": []
  }
}
```


---

## POST /effects

Add an active effect to an actor or token

Adds a status condition (by statusId) or a custom ActiveEffect (via effectData) to the specified actor or token.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to add the effect to |
| clientId | string |  | query | Client ID for the Foundry world |
| statusId | string |  | body, query | Standard status condition ID (e.g., "poisoned", "blinded", "prone") |
| effectData | object |  | body, query | Custom ActiveEffect data object (name, icon, duration, changes) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the add operation

### Try It Out

<ApiTester
  method="POST"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"statusId","type":"string","required":false,"source":"body"},{"name":"effectData","type":"object","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/effects';
const params = {
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'x-api-key': 'your-api-key-here',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "uuid": "Actor.pxZTVHItjx6GgPgC",
      "effectData": {
        "name": "Test Effect",
        "icon": "icons/svg/aura.svg",
        "changes": []
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/effects?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"uuid":"Actor.pxZTVHItjx6GgPgC","effectData":{"name":"Test Effect","icon":"icons/svg/aura.svg","changes":[]}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/effects'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here',
        'Content-Type': 'application/json'
    },
    json={
      "uuid": "Actor.pxZTVHItjx6GgPgC",
      "effectData": {
        "name": "Test Effect",
        "icon": "icons/svg/aura.svg",
        "changes": []
      }
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
  const path = '/effects';
  const params = {
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'post',
    headers: {
      'x-api-key': 'your-api-key-here',
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "uuid": "Actor.pxZTVHItjx6GgPgC",
        "effectData": {
          "name": "Test Effect",
          "icon": "icons/svg/aura.svg",
          "changes": []
        }
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
  🔤/effects🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"uuid":"Actor.pxZTVHItjx6GgPgC","effectData":{"name":"Test Effect","icon":"icons/svg/aura.svg","changes":[]}}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /effects🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 110❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "add-effect-result",
  "requestId": "add-effect_1775068881875",
  "data": {
    "uuid": "Actor.pxZTVHItjx6GgPgC",
    "effect": {
      "id": "enhJy4tchkKsw51Y",
      "uuid": "Actor.pxZTVHItjx6GgPgC.ActiveEffect.enhJy4tchkKsw51Y",
      "name": "Test Effect",
      "icon": "icons/svg/aura.svg",
      "statuses": []
    }
  }
}
```


---

## DELETE /effects

Remove an active effect from an actor or token

Removes an effect by its document ID (effectId) or by status condition identifier (statusId).

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to remove the effect from |
| clientId | string |  | query | Client ID for the Foundry world |
| effectId | string |  | body, query | The ActiveEffect document ID to remove |
| statusId | string |  | body, query | Standard status condition ID to remove (e.g., "poisoned") |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the remove operation

### Try It Out

<ApiTester
  method="DELETE"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"effectId","type":"string","required":false,"source":"body"},{"name":"statusId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/effects';
const params = {
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'DELETE',
  headers: {
    'x-api-key': 'your-api-key-here',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "uuid": "Actor.pxZTVHItjx6GgPgC",
      "effectId": "enhJy4tchkKsw51Y"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X DELETE 'http://localhost:3010/effects?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"uuid":"Actor.pxZTVHItjx6GgPgC","effectId":"enhJy4tchkKsw51Y"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/effects'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif'
}
url = f'{base_url}{path}'

response = requests.delete(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here',
        'Content-Type': 'application/json'
    },
    json={
      "uuid": "Actor.pxZTVHItjx6GgPgC",
      "effectId": "enhJy4tchkKsw51Y"
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
  const path = '/effects';
  const params = {
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'delete',
    headers: {
      'x-api-key': 'your-api-key-here',
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "uuid": "Actor.pxZTVHItjx6GgPgC",
        "effectId": "enhJy4tchkKsw51Y"
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
  🔤/effects🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"uuid":"Actor.pxZTVHItjx6GgPgC","effectId":"enhJy4tchkKsw51Y"}🔤 ➡️ body

  💭 Build HTTP request
  🔤DELETE /effects🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 63❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "remove-effect-result",
  "requestId": "remove-effect_1775068881979",
  "data": {
    "uuid": "Actor.pxZTVHItjx6GgPgC",
    "removedEffectId": "enhJy4tchkKsw51Y"
  }
}
```


