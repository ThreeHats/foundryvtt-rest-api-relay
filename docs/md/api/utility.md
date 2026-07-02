---
tag: utility
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Utility

## POST /select

Select token(s)

Selects one or more tokens in the Foundry VTT client.

**Required scope:** `entity:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| uuids | array |  | body | Array of UUIDs to select |
| name | string |  | body | Name of the token(s) to select |
| data | object |  | body | Data to match for selection (e.g., "data.attributes.hp.value": 20) |
| overwrite | boolean |  | body | Whether to overwrite existing selection |
| all | boolean |  | body | Whether to select all tokens on the canvas |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The selected token(s)

### Try It Out

<ApiTester
  method="POST"
  path="/select"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuids","type":"array","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"data","type":"object","required":false,"source":"body"},{"name":"overwrite","type":"boolean","required":false,"source":"body"},{"name":"all","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3011';
const path = '/select';
const params = {
  clientId: 'qsl-integration-test'
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
      "all": true,
      "overwrite": true
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3011/select?clientId=qsl-integration-test' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"all":true,"overwrite":true}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/select'
params = {
    'clientId': 'qsl-integration-test'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "all": True,
      "overwrite": True
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
  const path = '/select';
  const params = {
    clientId: 'qsl-integration-test'
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
        "all": true,
        "overwrite": true
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
  3011 ➡️ port
  🔤/select🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"all":true,"overwrite":true}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /select🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 29❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "select-result",
  "requestId": "select_1782956929084",
  "success": true,
  "count": 1,
  "message": "1 entities selected",
  "selected": [
    "Scene.SM7AhDv5JgZh6IvK.Token.h982FKt6QfxjSNkM"
  ]
}
```


---

## GET /selected

Get selected token(s)

Retrieves the currently selected token(s) in the Foundry VTT client.

**Required scope:** `entity:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The selected token(s)

### Try It Out

<ApiTester
  method="GET"
  path="/selected"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3011';
const path = '/selected';
const params = {
  clientId: 'qsl-integration-test'
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
curl -X GET 'http://localhost:3011/selected?clientId=qsl-integration-test' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/selected'
params = {
    'clientId': 'qsl-integration-test'
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
  const baseUrl = 'http://localhost:3011';
  const path = '/selected';
  const params = {
    clientId: 'qsl-integration-test'
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
  3011 ➡️ port
  🔤/selected🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /selected🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "selected-result",
  "requestId": "selected_1782956929088",
  "success": true,
  "selected": [
    {
      "tokenUuid": "Scene.SM7AhDv5JgZh6IvK.Token.h982FKt6QfxjSNkM",
      "actorUuid": "Scene.SM7AhDv5JgZh6IvK.Token.h982FKt6QfxjSNkM.Actor.V5OF1QXHjaIy6iO8"
    }
  ]
}
```


---

## GET /players

Get players/users

Retrieves a list of all users configured in the Foundry VTT world. Useful for discovering valid userId values for permission-scoped API calls.

**Required scope:** `entity:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - List of users with their IDs, names, roles, and active status

### Try It Out

<ApiTester
  method="GET"
  path="/players"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3011';
const path = '/players';
const params = {
  clientId: 'qsl-integration-test'
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
curl -X GET 'http://localhost:3011/players?clientId=qsl-integration-test' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/players'
params = {
    'clientId': 'qsl-integration-test'
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
  const baseUrl = 'http://localhost:3011';
  const path = '/players';
  const params = {
    clientId: 'qsl-integration-test'
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
  3011 ➡️ port
  🔤/players🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /players🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "players-result",
  "requestId": "players_1782956929091",
  "users": [
    {
      "id": "cpvaGKk3hgoBCzCS",
      "name": "Gamemaster",
      "role": 4,
      "isGM": true,
      "active": true,
      "color": "#161068",
      "avatar": "icons/svg/mystery-man.svg"
    },
    {
      "id": "zi1MHwh4aJs4L2Mn",
      "name": "test",
      "role": 1,
      "isGM": false,
      "active": false,
      "color": "#ccad28",
      "avatar": "icons/svg/mystery-man.svg"
    }
  ]
}
```


---

## GET /world-info

Get comprehensive world information

Returns a single object with world name, game system, Foundry version, all modules (with active status), all users (with online status), and the active scene. Useful for API clients to discover the world state.

**Required scope:** `world:info`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - World information object

### Try It Out

<ApiTester
  method="GET"
  path="/world-info"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3011';
const path = '/world-info';
const params = {
  clientId: 'qsl-integration-test'
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
curl -X GET 'http://localhost:3011/world-info?clientId=qsl-integration-test' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/world-info'
params = {
    'clientId': 'qsl-integration-test'
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
  const baseUrl = 'http://localhost:3011';
  const path = '/world-info';
  const params = {
    clientId: 'qsl-integration-test'
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
  3011 ➡️ port
  🔤/world-info🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /world-info🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "world-info-result",
  "requestId": "world-info_1782956929095",
  "data": {
    "world": {
      "id": "test-world",
      "title": "Test World"
    },
    "system": {
      "id": "dnd5e",
      "title": "Dungeons & Dragons Fifth Edition",
      "version": "5.0.4"
    },
    "foundryVersion": "14.363",
    "modules": [
      {
        "id": "foundry-rest-api",
        "title": "Foundry REST API",
        "version": "3.2.3",
        "active": true
      }
    ],
    "users": [
      {
        "id": "cpvaGKk3hgoBCzCS",
        "name": "Gamemaster",
        "role": 4,
        "isGM": true,
        "active": true,
        "color": "#161068",
        "avatar": "icons/svg/mystery-man.svg"
      },
      {
        "id": "zi1MHwh4aJs4L2Mn",
        "name": "test",
        "role": 1,
        "isGM": false,
        "active": false,
        "color": "#ccad28",
        "avatar": "icons/svg/mystery-man.svg"
      }
    ],
    "activeScene": {
      "id": "SM7AhDv5JgZh6IvK",
      "name": "test-scene-updated"
    }
  }
}
```


---

## POST /execute-js

Execute JavaScript

Executes a JavaScript script in the Foundry VTT client.

**Required scope:** `execute-js`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| script | string |  | body | JavaScript script to execute |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The result of the executed script

### Try It Out

<ApiTester
  method="POST"
  path="/execute-js"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"script","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3011';
const path = '/execute-js';
const params = {
  clientId: 'qsl-integration-test'
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
      "script": "const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3011/execute-js?clientId=qsl-integration-test' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"script":"const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/execute-js'
params = {
    'clientId': 'qsl-integration-test'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "script": "const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;"
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
  const path = '/execute-js';
  const params = {
    clientId: 'qsl-integration-test'
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
        "script": "const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;"
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
  3011 ➡️ port
  🔤/execute-js🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"script":"const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /execute-js🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 104❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "execute-js-result",
  "requestId": "execute-js_1782956929093",
  "success": true,
  "result": "ws://relay:3010"
}
```


