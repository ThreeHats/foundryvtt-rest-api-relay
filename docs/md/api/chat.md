---
tag: chat
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';
import SseTester from '@site/src/components/SseTester';

# Chat

## GET /chat

Get chat messages

Retrieves chat messages from the Foundry world with optional pagination and filtering.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| limit | number |  | query | Maximum number of messages to return (default: 10) |
| offset | number |  | query | Number of messages to skip for pagination |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |
| chatType | number |  | query | Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field. |
| speaker | string |  | query | Filter messages by speaker name or actor ID |

### Returns

**object** - Paginated list of chat messages

### Try It Out

<ApiTester
  method="GET"
  path="/chat"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"limit","type":"number","required":false,"source":"query"},{"name":"offset","type":"number","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"},{"name":"chatType","type":"number","required":false,"source":"query"},{"name":"speaker","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/chat';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  limit: '10'
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
curl -X GET 'http://localhost:3010/chat?clientId=fvtt_099ad17ea199e7e3&limit=10' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/chat'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'limit': '10'
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
  const path = '/chat';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    limit: '10'
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
  🔤/chat🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤limit=10🔤 ➡️ limit
  🔤?🧲clientId🧲&🧲limit🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /chat🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "chat-messages-result",
  "requestId": "chat-messages_1776657991264",
  "success": true,
  "data": {
    "messages": [
      {
        "id": "1rNCZ9ooCI1R5SMv",
        "uuid": "ChatMessage.1rNCZ9ooCI1R5SMv",
        "content": "This is a whispered test message",
        "speaker": {
          "scene": null,
          "actor": null,
          "token": null,
          "alias": "API Test Bot"
        },
        "timestamp": 1776657991259,
        "whisper": [],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "",
        "isRoll": false,
        "rolls": [],
        "flags": {}
      },
      {
        "id": "I5Pyye6zZD9bQVY5",
        "uuid": "ChatMessage.I5Pyye6zZD9bQVY5",
        "content": "Hello from the REST API test suite!",
        "speaker": {
          "scene": null,
          "actor": null,
          "token": null
        },
        "timestamp": 1776657991252,
        "whisper": [],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "Test Message",
        "isRoll": false,
        "rolls": [],
        "flags": {}
      },
      {
        "id": "WCkilrDZ2hWVHXG4",
        "uuid": "ChatMessage.WCkilrDZ2hWVHXG4",
        "content": "19",
        "speaker": {
          "scene": "r36nfimJGHYGUGQX",
          "actor": "q9uWyfdPwTlzbpxb",
          "token": "Q9y6lc2dPYs2WRQT",
          "alias": "Updated Test Actor"
        },
        "timestamp": 1776657990727,
        "whisper": [],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "Updated Test Actor rolls for Initiative!",
        "isRoll": true,
        "rolls": [
          {
            "formula": "1d20 + 3 + 0",
            "total": 19,
            "isCritical": false,
            "isFumble": false,
            "dice": [
              {
                "faces": 20,
                "results": [
                  {
                    "result": 16,
                    "active": true
                  }
                ]
              }
            ]
          }
        ],
        "flags": {
          "core": {
            "initiativeRoll": true
          }
        }
      },
      {
        "id": "7GyAvLpLupFBCUgU",
        "uuid": "ChatMessage.7GyAvLpLupFBCUgU",
        "content": "17",
        "speaker": {
          "scene": "r36nfimJGHYGUGQX",
          "actor": "q9uWyfdPwTlzbpxb",
          "token": "Q9y6lc2dPYs2WRQT",
          "alias": "Updated Test Actor"
        },
        "timestamp": 1776657990531,
        "whisper": [],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "Updated Test Actor rolls for Initiative!",
        "isRoll": true,
        "rolls": [
          {
            "formula": "1d20 + 3 + 0",
            "total": 17,
            "isCritical": false,
            "isFumble": false,
            "dice": [
              {
                "faces": 20,
                "results": [
                  {
                    "result": 14,
                    "active": true
                  }
                ]
              }
            ]
          }
        ],
        "flags": {
          "core": {
            "initiativeRoll": true
          }
        }
      },
      {
        "id": "dV9tVKJTzFwOka7L",
        "uuid": "ChatMessage.dV9tVKJTzFwOka7L",
        "content": "<b>⚠ REST API execute-js:</b> <code>const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;</code>",
        "speaker": {
          "scene": null,
          "actor": null,
          "token": null,
          "alias": "REST API Module"
        },
        "timestamp": 1776657990095,
        "whisper": [
          "r6bXhB7k9cXa3cif"
        ],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "",
        "isRoll": false,
        "rolls": [],
        "flags": {}
      },
      {
        "id": "HqaHnfvFgKeNjS9I",
        "uuid": "ChatMessage.HqaHnfvFgKeNjS9I",
        "content": "<b>⚠ REST API macro-execute:</b> <code>test-macro</code> (Macro.AkcLmoRwvkrPvjyA)",
        "speaker": {
          "scene": null,
          "actor": null,
          "token": null,
          "alias": "REST API Module"
        },
        "timestamp": 1776657989845,
        "whisper": [
          "r6bXhB7k9cXa3cif"
        ],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "",
        "isRoll": false,
        "rolls": [],
        "flags": {}
      },
      {
        "id": "29dKiKxs0o74LgU8",
        "uuid": "ChatMessage.29dKiKxs0o74LgU8",
        "content": "15",
        "speaker": {
          "scene": null,
          "actor": null,
          "token": null
        },
        "timestamp": 1776657970180,
        "whisper": [],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "SSE Test Roll",
        "isRoll": true,
        "rolls": [
          {
            "formula": "1d20",
            "total": 15,
            "isCritical": false,
            "isFumble": false,
            "dice": [
              {
                "faces": 20,
                "results": [
                  {
                    "result": 15,
                    "active": true
                  }
                ]
              }
            ]
          }
        ],
        "flags": {}
      },
      {
        "id": "LQiRbHhHpBIOgXsL",
        "uuid": "ChatMessage.LQiRbHhHpBIOgXsL",
        "content": "13",
        "speaker": {
          "scene": null,
          "actor": null,
          "token": null
        },
        "timestamp": 1776657969661,
        "whisper": [],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "Test Roll",
        "isRoll": true,
        "rolls": [
          {
            "formula": "2d20kh",
            "total": 13,
            "isCritical": false,
            "isFumble": false,
            "dice": [
              {
                "faces": 20,
                "results": [
                  {
                    "result": 13,
                    "active": true
                  },
                  {
                    "result": 8,
                    "active": false
                  }
                ]
              }
            ]
          }
        ],
        "flags": {}
      },
      {
        "id": "j6g16FombklYBmBJ",
        "uuid": "ChatMessage.j6g16FombklYBmBJ",
        "content": "<i class=\"fas fa-link\"></i> REST API connected to <code>localhost:3010</code>",
        "speaker": {
          "scene": null,
          "actor": null,
          "token": null,
          "alias": "REST API Module"
        },
        "timestamp": 1776657947586,
        "whisper": [
          "r6bXhB7k9cXa3cif"
        ],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "",
        "isRoll": false,
        "rolls": [],
        "flags": {}
      },
      {
        "id": "wb1WdTavW1uHF4Jg",
        "uuid": "ChatMessage.wb1WdTavW1uHF4Jg",
        "content": "",
        "speaker": {
          "scene": "Fwzj1xLpvHBiXYla",
          "actor": "pqHYOrBnvc1S4C3y",
          "token": "OBnzhG0BJxttlZIF",
          "alias": "Updated Test Actor"
        },
        "timestamp": 1776656664133,
        "whisper": [],
        "type": "base",
        "author": {
          "id": "r6bXhB7k9cXa3cif",
          "name": "tester"
        },
        "flavor": "Constitution Saving Throw",
        "isRoll": true,
        "rolls": [
          {
            "formula": "1d20 + 1",
            "total": 21,
            "isCritical": true,
            "isFumble": false,
            "dice": [
              {
                "faces": 20,
                "results": [
                  {
                    "result": 20,
                    "active": true
                  }
                ]
              }
            ]
          }
        ],
        "flags": {
          "dnd5e": {
            "messageType": "roll",
            "roll": {
              "ability": "con",
              "type": "save"
            }
          }
        }
      }
    ],
    "total": 22,
    "offset": 0,
    "limit": 10
  }
}
```


---

## POST /chat

Send a chat message

Creates a new chat message in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| content | string | ✓ | body | The message content (supports HTML) |
| clientId | string |  | query | Client ID for the Foundry world |
| whisper | array |  | body | Array of user IDs to whisper the message to |
| speaker | string |  | body | Actor ID to use as the message speaker |
| alias | string |  | body | Display name alias for the speaker |
| chatType | number |  | body | Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field. |
| flavor | string |  | body | Flavor text displayed above the message content |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The created chat message

### Try It Out

<ApiTester
  method="POST"
  path="/chat"
  parameters={[{"name":"content","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"whisper","type":"array","required":false,"source":"body"},{"name":"speaker","type":"string","required":false,"source":"body"},{"name":"alias","type":"string","required":false,"source":"body"},{"name":"chatType","type":"number","required":false,"source":"body"},{"name":"flavor","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/chat';
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
      "content": "Hello from the REST API test suite!",
      "flavor": "Test Message"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/chat?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello from the REST API test suite!","flavor":"Test Message"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/chat'
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
      "content": "Hello from the REST API test suite!",
      "flavor": "Test Message"
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
  const path = '/chat';
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
        "content": "Hello from the REST API test suite!",
        "flavor": "Test Message"
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
  🔤/chat🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"content":"Hello from the REST API test suite!","flavor":"Test Message"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /chat🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 73❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "chat-send-result",
  "requestId": "chat-send_1776657991251",
  "success": true,
  "data": {
    "id": "I5Pyye6zZD9bQVY5",
    "uuid": "ChatMessage.I5Pyye6zZD9bQVY5",
    "content": "Hello from the REST API test suite!",
    "speaker": {
      "scene": null,
      "actor": null,
      "token": null
    },
    "timestamp": 1776657991252,
    "whisper": [],
    "type": "base",
    "author": {
      "id": "r6bXhB7k9cXa3cif",
      "name": "tester"
    },
    "flavor": "Test Message",
    "isRoll": false,
    "rolls": [],
    "flags": {}
  }
}
```


---

## DELETE /chat/:messageId

Delete a specific chat message

Deletes a chat message by its ID. Only the message author or a GM can delete messages.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| messageId | string | ✓ | params | ID of the chat message to delete |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Success confirmation

### Try It Out

<ApiTester
  method="DELETE"
  path="/chat/:messageId"
  parameters={[{"name":"messageId","type":"string","required":true,"source":"params"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/chat/I5Pyye6zZD9bQVY5';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3'
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
curl -X DELETE 'http://localhost:3010/chat/I5Pyye6zZD9bQVY5?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/chat/I5Pyye6zZD9bQVY5'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
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
  const path = '/chat/I5Pyye6zZD9bQVY5';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3'
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
  🔤/chat/I5Pyye6zZD9bQVY5🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /chat/I5Pyye6zZD9bQVY5🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "chat-delete-result",
  "requestId": "chat-delete_1776657991271",
  "success": true,
  "data": {
    "messageId": "I5Pyye6zZD9bQVY5"
  }
}
```


---

## DELETE /chat

Clear all chat messages

Flushes all chat message history. Only GMs can perform this action.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Success confirmation

### Try It Out

<ApiTester
  method="DELETE"
  path="/chat"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/chat';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3'
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
curl -X DELETE 'http://localhost:3010/chat?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/chat'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
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
  const path = '/chat';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3'
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
  🔤/chat🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /chat🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "chat-flush-result",
  "requestId": "chat-flush_1776657991287",
  "success": true,
  "data": {
    "message": "All chat messages have been deleted"
  }
}
```


---

## GET /chat/subscribe

Subscribe to real-time chat events via Server-Sent Events (SSE)

Opens a persistent SSE connection that streams chat events (create, update, delete) as they occur in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| speaker | string |  | query | Filter events by speaker name or actor ID |
| type | number |  | query | Filter events by chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll) |
| whisperOnly | boolean |  | query | Only receive whispered messages |
| userId | string |  | query | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**SSE stream** - SSE event stream

### Try It Out

<SseTester
  path="/chat/subscribe"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"speaker","type":"string","required":false,"source":"query"},{"name":"type","type":"number","required":false,"source":"query"},{"name":"whisperOnly","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const { EventSource } = require('eventsource'); // npm install eventsource

const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key-here';
const url = `${baseUrl}/chat/subscribe?clientId=your-client-id`;

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

function formatMessage(prefix, message) {
  const speaker = message.author?.name || message.speaker?.alias || '?';
  console.log(`[${prefix}] ${speaker}: ${message.content}`);
  if (message.flavor) console.log(`  Flavor: ${message.flavor}`);
  if (message.isRoll && message.rolls?.length > 0) {
    for (const roll of message.rolls) {
      const dice = roll.dice?.map(d =>
        `${d.results.map(r => `${r.result}${r.active ? '' : '(dropped)'}`).join(', ')} (d${d.faces})`
      ).join(' + ') || '';
      console.log(`  Roll: ${roll.formula} = ${roll.total}${roll.isCritical ? ' CRITICAL!' : ''}${roll.isFumble ? ' FUMBLE!' : ''}`);
      if (dice) console.log(`  Dice: ${dice}`);
    }
  }
}

eventSource.addEventListener('connected', (event) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
});

eventSource.addEventListener('chat-create', (event) => {
  const message = JSON.parse(event.data);
  formatMessage('new', message);
});

eventSource.addEventListener('chat-update', (event) => {
  const message = JSON.parse(event.data);
  formatMessage('updated', message);
});

eventSource.addEventListener('chat-delete', (event) => {
  const data = JSON.parse(event.data);
  console.log('Message deleted:', JSON.stringify(data));
});

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
};

// To disconnect later:
// eventSource.close();
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
# Connect to the SSE stream (streams events until interrupted with Ctrl+C)
curl -N 'http://localhost:3010/chat/subscribe?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Accept: text/event-stream"

# Example output:
# event: connected
# data: {"clientId":"your-client-id"}
#
# event: chat-create
# data: {"id":"abc123","content":"Hello!","author":{"id":"xyz","name":"GM"},"isRoll":false,...}
#
# event: chat-create (dice roll)
# data: {"id":"def456","content":"16","flavor":"Attack Roll","isRoll":true,"rolls":[{"formula":"1d20+5","total":16,"isCritical":false,"isFumble":false,"dice":[{"faces":20,"results":[{"result":11,"active":true}]}]}],...}
```

</TabItem>
<TabItem value="python" label="Python">

```python
import sseclient  # pip install sseclient-py
import requests
import json

base_url = 'http://localhost:3010'
url = f'{base_url}/chat/subscribe'
params = {'clientId': 'your-client-id'}
headers = {
    'x-api-key': 'your-api-key-here',
    'Accept': 'text/event-stream'
}

# Connect to the SSE stream
response = requests.get(url, params=params, headers=headers, stream=True)
client = sseclient.SSEClient(response)

for event in client.events():
    data = json.loads(event.data)

    if event.event == 'connected':
        print(f'Connected: {data["clientId"]}')
    elif event.event in ('chat-create', 'chat-update'):
        prefix = 'new' if event.event == 'chat-create' else 'updated'
        speaker = (data.get('author') or {}).get('name') or (data.get('speaker') or {}).get('alias') or '?'
        print(f'[{prefix}] {speaker}: {data.get("content", "")}')
        if data.get('flavor'):
            print(f'  Flavor: {data["flavor"]}')
        if data.get('isRoll') and data.get('rolls'):
            for roll in data['rolls']:
                dice_parts = []
                for d in roll.get('dice', []):
                    results = ', '.join(
                        f'{r["result"]}{"" if r.get("active", True) else "(dropped)"}'
                        for r in d.get('results', [])
                    )
                    dice_parts.append(f'{results} (d{d["faces"]})')
                crit = ' CRITICAL!' if roll.get('isCritical') else ''
                fumble = ' FUMBLE!' if roll.get('isFumble') else ''
                print(f'  Roll: {roll["formula"]} = {roll["total"]}{crit}{fumble}')
                if dice_parts:
                    print(f'  Dice: {" + ".join(dice_parts)}')
    elif event.event == 'chat-delete':
        print(f'Message deleted: {json.dumps(data)}')
```

</TabItem>
<TabItem value="typescript" label="TypeScript">

```typescript
// npm install eventsource
import { EventSource } from 'eventsource';

const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key-here';
const url = `${baseUrl}/chat/subscribe?clientId=your-client-id`;

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

interface ChatMessage {
  id: string;
  content: string;
  type: string;
  author: { id: string; name: string };
  speaker: any;
  timestamp: number;
  flavor: string;
  isRoll: boolean;
  rolls: {
    formula: string;
    total: number;
    isCritical: boolean;
    isFumble: boolean;
    dice: { faces: number; results: { result: number; active: boolean }[] }[];
  }[];
  whisper: string[];
  flags: Record<string, any>;
}

function formatMessage(prefix: string, message: ChatMessage) {
  const speaker = message.author?.name || message.speaker?.alias || '?';
  console.log(`[${prefix}] ${speaker}: ${message.content}`);
  if (message.flavor) console.log(`  Flavor: ${message.flavor}`);
  if (message.isRoll && message.rolls?.length > 0) {
    for (const roll of message.rolls) {
      const dice = roll.dice?.map(d =>
        `${d.results.map(r => `${r.result}${r.active ? '' : '(dropped)'}`).join(', ')} (d${d.faces})`
      ).join(' + ') || '';
      console.log(`  Roll: ${roll.formula} = ${roll.total}${roll.isCritical ? ' CRITICAL!' : ''}${roll.isFumble ? ' FUMBLE!' : ''}`);
      if (dice) console.log(`  Dice: ${dice}`);
    }
  }
}

eventSource.addEventListener('connected', (event: MessageEvent) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
});

eventSource.addEventListener('chat-create', (event: MessageEvent) => {
  const message: ChatMessage = JSON.parse(event.data);
  formatMessage('new', message);
});

eventSource.addEventListener('chat-update', (event: MessageEvent) => {
  const message: ChatMessage = JSON.parse(event.data);
  formatMessage('updated', message);
});

eventSource.addEventListener('chat-delete', (event: MessageEvent) => {
  const data = JSON.parse(event.data);
  console.log('Message deleted:', JSON.stringify(data));
});

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
};

// To disconnect: eventSource.close();
```

</TabItem>
<TabItem value="emojicode" label="Emojicode">

```emojicode
Just don't 😂
```

</TabItem>
</Tabs>

#### Response

**Status:** 200

```json
{
  "event": "connected",
  "data": {
    "clientId": "fvtt_099ad17ea199e7e3"
  }
}
```


