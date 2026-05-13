---
tag: macro
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Macro

## GET /macros

Get all macros

Retrieves a list of all macros available in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - An array of macros with details

### Try It Out

<ApiTester
  method="GET"
  path="/macros"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/macros';
const params = {
  clientId: 'fvtt_71dbc81bd608978a'
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
curl -X GET 'http://localhost:3010/macros?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/macros'
params = {
    'clientId': 'fvtt_71dbc81bd608978a'
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
    clientId: 'fvtt_71dbc81bd608978a'
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
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
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
  "type": "macros-result",
  "requestId": "macros_1777996626130",
  "macros": [
    {
      "uuid": "Macro.bRLMIJXvTs22CRK0",
      "id": "bRLMIJXvTs22CRK0",
      "name": "test-macro",
      "type": "script",
      "author": "tester",
      "command": "// Example macro that uses parameters\nfunction myMacro(args) {\n  const targetName = args.targetName || \"Target\";\n  const damage = args.damage || 0;\n  const effect = args.effect || \"none\";\n  \n  // Use the parameters\n  console.log(`Attacking ${targetName} for ${damage} ${effect} damage`);\n  \n  // Return a value (can be any data type)\n  return {\n    success: true,\n    damageDealt: damage,\n    target: targetName\n  };\n}\n\n// Don't forget to return the result of your function\nreturn myMacro(args);",
      "img": "icons/svg/dice-target.svg",
      "scope": "global",
      "canExecute": true
    },
    {
      "uuid": "Macro.vMGQmIwYBNK6DiJx",
      "id": "vMGQmIwYBNK6DiJx",
      "name": "Nuke",
      "type": "script",
      "author": "Gamemaster",
      "command": "async function cleanSlate() {\n    // List of document types to wipe\n    const collections = [\n        game.scenes,\n        game.actors,\n        game.items,\n        game.journal,\n        game.tables,\n        game.playlists,\n        game.cards,\n        // game.macros // Uncomment this line if you want to delete all macros too\n    ];\n\n    for (let collection of collections) {\n        const ids = collection.map(doc => doc.id);\n        if (ids.length > 0) {\n            console.log(`Deleting ${ids.length} documents from ${collection.name}...`);\n            await collection.documentClass.deleteDocuments(ids);\n        }\n    }\n\n    ui.notifications.info(\"World cleanup complete. A fresh start awaits!\");\n}\n\n// Confirmation Dialog\nnew Dialog({\n    title: \"Nuclear Option: Clear World Data\",\n    content: `\n        <div style=\"text-align: center;\">\n            <p><i class=\"fas fa-exclamation-triangle fa-3x\" style=\"color: #ff6b6b;\"></i></p>\n            <p>This will <strong>permanently delete</strong> all Scenes, Actors, Items, Journals, and more.</p>\n            <p><em>Are you absolutely sure?</em></p>\n        </div>`,\n    buttons: {\n        confirm: {\n            icon: '<i class=\"fas fa-trash\"></i>',\n            label: \"Delete Everything\",\n            callback: () => cleanSlate()\n        },\n        cancel: {\n            icon: '<i class=\"fas fa-times\"></i>',\n            label: \"Cancel\"\n        }\n    },\n    default: \"cancel\"\n}).render(true);",
      "img": "icons/svg/poison.svg",
      "scope": "global",
      "canExecute": true
    },
    {
      "uuid": "Macro.QsVroAGCs54xe5aX",
      "id": "QsVroAGCs54xe5aX",
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

Execute a macro by UUID

Executes a specific macro in the Foundry world by its UUID.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | params | UUID of the macro to execute |
| clientId | string |  | query | Client ID for the Foundry world |
| args | object |  | body | Optional arguments to pass to the macro execution |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the macro execution

### Try It Out

<ApiTester
  method="POST"
  path="/macro/:uuid/execute"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"params"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"args","type":"object","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/macro/Macro.QsVroAGCs54xe5aX/execute';
const params = {
  clientId: 'fvtt_71dbc81bd608978a'
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
curl -X POST 'http://localhost:3010/macro/Macro.QsVroAGCs54xe5aX/execute?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"args":{"targetName":"Goblin","damage":100000,"effect":"poison"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/macro/Macro.QsVroAGCs54xe5aX/execute'
params = {
    'clientId': 'fvtt_71dbc81bd608978a'
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
  const path = '/macro/Macro.QsVroAGCs54xe5aX/execute';
  const params = {
    clientId: 'fvtt_71dbc81bd608978a'
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
  🔤/macro/Macro.QsVroAGCs54xe5aX/execute🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"args":{"targetName":"Goblin","damage":100000,"effect":"poison"}}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /macro/Macro.QsVroAGCs54xe5aX/execute🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 66❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "macro-execute-result",
  "requestId": "macro-execute_1777996626133",
  "uuid": "Macro.QsVroAGCs54xe5aX",
  "success": true,
  "result": {
    "success": true,
    "damageDealt": 100000,
    "target": "Goblin"
  }
}
```


