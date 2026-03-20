---
tag: macro
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


# macro

## GET /macros

Get all macros Retrieves a list of all macros available in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - An array of macros with details

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/macros';
const params = {
  clientId: 'your-client-id'
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
curl -X GET 'http://localhost:3010/macros?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/macros'
params = {
    'clientId': 'your-client-id'
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
  const path = '/macros';
  const params = {
    clientId: 'your-client-id'
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
  🔤/macros🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /macros🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "macros_1773999634831",
  "clientId": "your-client-id",
  "type": "macros-result",
  "macros": [
    {
      "uuid": "Macro.TLcYwBMu3V4c09Ap",
      "id": "TLcYwBMu3V4c09Ap",
      "name": "test-macro",
      "type": "script",
      "author": "tester",
      "command": "// Example macro that uses parameters\nfunction myMacro(args) {\n  const targetName = args.targetName || \"Target\";\n  const damage = args.damage || 0;\n  const effect = args.effect || \"none\";\n  \n  // Use the parameters\n  console.log(`Attacking ${targetName} for ${damage} ${effect} damage`);\n  \n  // Return a value (can be any data type)\n  return {\n    success: true,\n    damageDealt: damage,\n    target: targetName\n  };\n}\n\n// Don't forget to return the result of your function\nreturn myMacro(args);",
      "img": "icons/svg/dice-target.svg",
      "scope": "global",
      "canExecute": true
    },
    {
      "uuid": "Macro.jxIBRKLXepzR9tmu",
      "id": "jxIBRKLXepzR9tmu",
      "name": "test-macro",
      "type": "script",
      "author": "tester",
      "command": "// Example macro that uses parameters\nfunction myMacro(args) {\n  const targetName = args.targetName || \"Target\";\n  const damage = args.damage || 0;\n  const effect = args.effect || \"none\";\n  \n  // Use the parameters\n  console.log(`Attacking ${targetName} for ${damage} ${effect} damage`);\n  \n  // Return a value (can be any data type)\n  return {\n    success: true,\n    damageDealt: damage,\n    target: targetName\n  };\n}\n\n// Don't forget to return the result of your function\nreturn myMacro(args);",
      "img": "icons/svg/dice-target.svg",
      "scope": "global",
      "canExecute": true
    },
    {
      "uuid": "Macro.u4CgAClWghPlKQRs",
      "id": "u4CgAClWghPlKQRs",
      "name": "test-macro",
      "type": "script",
      "author": "tester",
      "command": "// Example macro that uses parameters\nfunction myMacro(args) {\n  const targetName = args.targetName || \"Target\";\n  const damage = args.damage || 0;\n  const effect = args.effect || \"none\";\n  \n  // Use the parameters\n  console.log(`Attacking ${targetName} for ${damage} ${effect} damage`);\n  \n  // Return a value (can be any data type)\n  return {\n    success: true,\n    damageDealt: damage,\n    target: targetName\n  };\n}\n\n// Don't forget to return the result of your function\nreturn myMacro(args);",
      "img": "icons/svg/dice-target.svg",
      "scope": "global",
      "canExecute": true
    },
    {
      "uuid": "Macro.ZP1Pu3sJ94JbtlgE",
      "id": "ZP1Pu3sJ94JbtlgE",
      "name": "test-macro",
      "type": "script",
      "author": "tester",
      "command": "// Example macro that uses parameters\nfunction myMacro(args) {\n  const targetName = args.targetName || \"Target\";\n  const damage = args.damage || 0;\n  const effect = args.effect || \"none\";\n  \n  // Use the parameters\n  console.log(`Attacking ${targetName} for ${damage} ${effect} damage`);\n  \n  // Return a value (can be any data type)\n  return {\n    success: true,\n    damageDealt: damage,\n    target: targetName\n  };\n}\n\n// Don't forget to return the result of your function\nreturn myMacro(args);",
      "img": "icons/svg/dice-target.svg",
      "scope": "global",
      "canExecute": true
    }
  ]
}
```


---

## POST /macro/:uuid/execute

Execute a macro by UUID Executes a specific macro in the Foundry world by its UUID.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| uuid | string | ✓ | params | UUID of the macro to execute |
| args | object |  | body | Optional arguments to pass to the macro execution |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the macro execution

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/macro/Macro.ZP1Pu3sJ94JbtlgE/execute';
const params = {
  clientId: 'your-client-id'
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
      "args": {
        "targetName": "Goblin",
        "damage": 100000,
        "effect": "poison"
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/macro/Macro.ZP1Pu3sJ94JbtlgE/execute?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"args":{"targetName":"Goblin","damage":100000,"effect":"poison"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/macro/Macro.ZP1Pu3sJ94JbtlgE/execute'
params = {
    'clientId': 'your-client-id'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "args": {
        "targetName": "Goblin",
        "damage": 100000,
        "effect": "poison"
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
  const path = '/macro/Macro.ZP1Pu3sJ94JbtlgE/execute';
  const params = {
    clientId: 'your-client-id'
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
        "args": {
          "targetName": "Goblin",
          "damage": 100000,
          "effect": "poison"
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
  🔤/macro/Macro.ZP1Pu3sJ94JbtlgE/execute🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"args":{"targetName":"Goblin","damage":100000,"effect":"poison"}}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /macro/Macro.ZP1Pu3sJ94JbtlgE/execute🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 66❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "macro-execute_1773999634963",
  "clientId": "your-client-id",
  "type": "macro-execute-result",
  "uuid": "Macro.ZP1Pu3sJ94JbtlgE",
  "success": true,
  "result": {
    "success": true,
    "damageDealt": 100000,
    "target": "Goblin"
  }
}
```


