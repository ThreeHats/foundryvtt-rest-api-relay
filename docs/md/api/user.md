---
tag: user
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# User

## GET /users

List all Foundry users

Retrieves a list of all users configured in the Foundry VTT world, including their roles and online status. This is a GM-only operation.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of user objects with id, name, role, isGM, active, color, avatar, and character fields

### Try It Out

<ApiTester
  method="GET"
  path="/users"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/users';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3'
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
curl -X GET 'http://localhost:3010/users?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/users'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
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
  const path = '/users';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3'
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
  🔤/users🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /users🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "get-users-result",
  "requestId": "get-users_1776657998667",
  "data": [
    {
      "id": "5ypAoBvOiyjDKiaZ",
      "name": "Gamemaster",
      "role": 4,
      "isGM": true,
      "active": false,
      "color": "#28cca2",
      "avatar": "icons/svg/mystery-man.svg",
      "character": null
    },
    {
      "id": "r6bXhB7k9cXa3cif",
      "name": "tester",
      "role": 4,
      "isGM": true,
      "active": true,
      "color": "#cc2829",
      "avatar": "icons/svg/mystery-man.svg",
      "character": null
    },
    {
      "id": "JLfKTeTgCDpAdDfw",
      "name": "some-cool-guy",
      "role": 1,
      "isGM": false,
      "active": false,
      "color": "#6328cc",
      "avatar": "icons/svg/mystery-man.svg",
      "character": null
    }
  ]
}
```


---

## GET /user

Get a single Foundry user

Retrieves a single user by their ID or name. This is a GM-only operation.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| id | string |  | query | ID of the user to retrieve |
| name | string |  | query | Name of the user to retrieve (alternative to id) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - User object with id, name, role, isGM, active, color, avatar, and character fields

### Try It Out

<ApiTester
  method="GET"
  path="/user"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"id","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/user';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  id: 'cpUI52NlPcovB61k'
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
curl -X GET 'http://localhost:3010/user?clientId=fvtt_099ad17ea199e7e3&id=cpUI52NlPcovB61k' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/user'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'id': 'cpUI52NlPcovB61k'
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
  const path = '/user';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    id: 'cpUI52NlPcovB61k'
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
  🔤/user🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤id=cpUI52NlPcovB61k🔤 ➡️ id
  🔤?🧲clientId🧲&🧲id🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /user🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "get-user-result",
  "requestId": "get-user_1776657998677",
  "data": {
    "id": "cpUI52NlPcovB61k",
    "name": "test-api-user",
    "role": 1,
    "isGM": false,
    "active": false,
    "color": "#64cc28",
    "avatar": "icons/svg/mystery-man.svg",
    "character": null
  }
}
```


---

## POST /user

Create a new Foundry user

Creates a new user in the Foundry VTT world with the specified name, role, and optional password. This is a GM-only operation.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| name | string | ✓ | body | Username for the new user |
| clientId | string |  | query | Client ID for the Foundry world |
| role | number |  | body | User role: 0=None, 1=Player, 2=Trusted, 3=Assistant, 4=GM (default: 1) |
| password | string |  | body | Password for the new user |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The created user object

### Try It Out

<ApiTester
  method="POST"
  path="/user"
  parameters={[{"name":"name","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"role","type":"number","required":false,"source":"body"},{"name":"password","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/user';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3'
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
      "name": "test-api-user",
      "role": 1,
      "password": "testpassword123"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/user?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"name":"test-api-user","role":1,"password":"testpassword123"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/user'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "name": "test-api-user",
      "role": 1,
      "password": "testpassword123"
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
  const path = '/user';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3'
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
        "name": "test-api-user",
        "role": 1,
        "password": "testpassword123"
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
  🔤/user🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"name":"test-api-user","role":1,"password":"testpassword123"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /user🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 62❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "create-user-result",
  "requestId": "create-user_1776657998670",
  "data": {
    "id": "cpUI52NlPcovB61k",
    "name": "test-api-user",
    "role": 1,
    "isGM": false,
    "active": false,
    "color": "#64cc28",
    "avatar": "icons/svg/mystery-man.svg",
    "character": null
  }
}
```


---

## PUT /user

Update an existing Foundry user

Updates fields on an existing user. Identify the user by id or name, then pass the fields to update in the data object. Cannot demote the last GM user. This is a GM-only operation.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| data | object | ✓ | body | Object containing user fields to update (name, role, password, color, avatar, character) |
| clientId | string |  | query | Client ID for the Foundry world |
| id | string |  | body, query | ID of the user to update |
| name | string |  | body, query | Name of the user to update (alternative to id) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The updated user object

### Try It Out

<ApiTester
  method="PUT"
  path="/user"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"id","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/user';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'PUT',
  headers: {
    'x-api-key': 'your-api-key-here',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "id": "cpUI52NlPcovB61k",
      "data": {
        "role": 2
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X PUT 'http://localhost:3010/user?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"id":"cpUI52NlPcovB61k","data":{"role":2}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/user'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
}
url = f'{base_url}{path}'

response = requests.put(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "id": "cpUI52NlPcovB61k",
      "data": {
        "role": 2
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
  const path = '/user';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'put',
    headers: {
      'x-api-key': 'your-api-key-here',
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "id": "cpUI52NlPcovB61k",
        "data": {
          "role": 2
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
  🔤/user🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"id":"cpUI52NlPcovB61k","data":{"role":2}}🔤 ➡️ body

  💭 Build HTTP request
  🔤PUT /user🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 43❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "update-user-result",
  "requestId": "update-user_1776657998681",
  "data": {
    "id": "cpUI52NlPcovB61k",
    "name": "test-api-user",
    "role": 2,
    "isGM": false,
    "active": false,
    "color": "#64cc28",
    "avatar": "icons/svg/mystery-man.svg",
    "character": null
  }
}
```


---

## DELETE /user

Delete a Foundry user

Permanently deletes a user from the Foundry VTT world. Cannot delete yourself or the last GM user. This is a GM-only operation.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| id | string |  | query | ID of the user to delete |
| name | string |  | query | Name of the user to delete (alternative to id) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Deletion result

### Try It Out

<ApiTester
  method="DELETE"
  path="/user"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"id","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/user';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  id: 'cpUI52NlPcovB61k'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'DELETE',
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
curl -X DELETE 'http://localhost:3010/user?clientId=fvtt_099ad17ea199e7e3&id=cpUI52NlPcovB61k' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/user'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'id': 'cpUI52NlPcovB61k'
}
url = f'{base_url}{path}'

response = requests.delete(
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
  const path = '/user';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    id: 'cpUI52NlPcovB61k'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'delete',
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
  🔤/user🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤id=cpUI52NlPcovB61k🔤 ➡️ id
  🔤?🧲clientId🧲&🧲id🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /user🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "delete-user-result",
  "requestId": "delete-user_1776657998687",
  "success": true
}
```


