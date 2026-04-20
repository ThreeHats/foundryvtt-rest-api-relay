---
tag: roll
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';
import SseTester from '@site/src/components/SseTester';

# Roll

## GET /rolls

Get recent rolls

Retrieves a list of up to 20 recent rolls made in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| limit | number |  | query | Optional limit on the number of rolls to return (default is 20) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - An array of recent rolls with details

### Try It Out

<ApiTester
  method="GET"
  path="/rolls"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"limit","type":"number","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/rolls';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  limit: '20'
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
curl -X GET 'http://localhost:3010/rolls?clientId=fvtt_099ad17ea199e7e3&limit=20' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/rolls'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'limit': '20'
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
  const path = '/rolls';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    limit: '20'
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
  🔤/rolls🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤limit=20🔤 ➡️ limit
  🔤?🧲clientId🧲&🧲limit🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /rolls🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "rolls-result",
  "requestId": "rolls_1776657970172",
  "data": [
    {
      "id": "manual_1776657969666_24z58qfq2z1",
      "messageId": "manual_1776657969666_24z58qfq2z1",
      "user": null,
      "speaker": {},
      "flavor": "Test Roll",
      "rollTotal": 13,
      "formula": "2d20kh",
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
      ],
      "timestamp": 1776657969666
    }
  ]
}
```


---

## GET /lastroll

Get the last roll

Retrieves the most recent roll made in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The most recent roll with details

### Try It Out

<ApiTester
  method="GET"
  path="/lastroll"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/lastroll';
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
curl -X GET 'http://localhost:3010/lastroll?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/lastroll'
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
  const path = '/lastroll';
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
  🔤/lastroll🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /lastroll🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "last-roll-result",
  "requestId": "last-roll_1776657970175",
  "data": {
    "id": "manual_1776657969666_24z58qfq2z1",
    "messageId": "manual_1776657969666_24z58qfq2z1",
    "user": null,
    "speaker": {},
    "flavor": "Test Roll",
    "rollTotal": 13,
    "formula": "2d20kh",
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
    ],
    "timestamp": 1776657969666
  }
}
```


---

## POST /roll

Make a roll

Executes a roll with the specified formula.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| formula | string | ✓ | body | The roll formula to evaluate (e.g., "1d20 + 5") |
| clientId | string |  | query | Client ID for the Foundry world |
| flavor | string |  | body | Optional flavor text for the roll |
| createChatMessage | boolean |  | body | Whether to create a chat message for the roll |
| speaker | string |  | body | The speaker for the roll |
| whisper | array |  | body | Users to whisper the roll result to |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the roll operation

### Try It Out

<ApiTester
  method="POST"
  path="/roll"
  parameters={[{"name":"formula","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"flavor","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"},{"name":"speaker","type":"string","required":false,"source":"body"},{"name":"whisper","type":"array","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/roll';
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
      "formula": "2d20kh",
      "flavor": "Test Roll",
      "createChatMessage": true
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/roll?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"formula":"2d20kh","flavor":"Test Roll","createChatMessage":true}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/roll'
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
      "formula": "2d20kh",
      "flavor": "Test Roll",
      "createChatMessage": True
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
  const path = '/roll';
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
        "formula": "2d20kh",
        "flavor": "Test Roll",
        "createChatMessage": true
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
  🔤/roll🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"formula":"2d20kh","flavor":"Test Roll","createChatMessage":true}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /roll🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 66❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "roll-result",
  "requestId": "roll_1776657969658",
  "success": true,
  "data": {
    "id": "manual_1776657969666_24z58qfq2z1",
    "chatMessageCreated": true,
    "roll": {
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
      ],
      "timestamp": 1776657969666
    }
  }
}
```


---

## GET /rolls/subscribe

Subscribe to real-time roll events via Server-Sent Events (SSE)

Opens a persistent SSE connection that streams roll events as they occur.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**SSE stream** - SSE event stream

### Try It Out

<SseTester
  path="/rolls/subscribe"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const { EventSource } = require('eventsource'); // npm install eventsource

const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key-here';
const url = `${baseUrl}/rolls/subscribe?clientId=your-client-id`;

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

eventSource.addEventListener('connected', (event) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
});

eventSource.addEventListener('roll', (event) => {
  const roll = JSON.parse(event.data);
  const dice = roll.dice?.map(d =>
    `${d.results.map(r => `${r.result}${r.active ? '' : '(dropped)'}`).join(', ')} (d${d.faces})`
  ).join(' + ') || '';
  console.log(`[${roll.user?.name}] ${roll.formula} = ${roll.rollTotal}${roll.isCritical ? ' CRITICAL!' : ''}${roll.isFumble ? ' FUMBLE!' : ''}`);
  if (roll.flavor) console.log(`  Flavor: ${roll.flavor}`);
  if (dice) console.log(`  Dice: ${dice}`);
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
# Connect to the roll SSE stream (streams events until interrupted with Ctrl+C)
curl -N 'http://localhost:3010/rolls/subscribe?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Accept: text/event-stream"

# Example output:
# event: connected
# data: {"clientId":"your-client-id"}
#
# event: roll
# data: {"id":"abc123","user":{"id":"xyz","name":"GM"},"formula":"1d20+5","rollTotal":18,"isCritical":false,"isFumble":false,"dice":[{"faces":20,"results":[{"result":13,"active":true}]}],"flavor":"Attack Roll","timestamp":1234567890}
```

</TabItem>
<TabItem value="python" label="Python">

```python
import sseclient  # pip install sseclient-py
import requests
import json

base_url = 'http://localhost:3010'
url = f'{base_url}/rolls/subscribe'
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
    elif event.event == 'roll':
        user = (data.get('user') or {}).get('name', '?')
        crit = ' CRITICAL!' if data.get('isCritical') else ''
        fumble = ' FUMBLE!' if data.get('isFumble') else ''
        print(f'[{user}] {data["formula"]} = {data["rollTotal"]}{crit}{fumble}')
        if data.get('flavor'):
            print(f'  Flavor: {data["flavor"]}')
        for d in data.get('dice', []):
            results = ', '.join(
                f'{r["result"]}{"" if r.get("active", True) else "(dropped)"}'
                for r in d.get('results', [])
            )
            print(f'  Dice: {results} (d{d["faces"]})')
```

</TabItem>
<TabItem value="typescript" label="TypeScript">

```typescript
// npm install eventsource
import { EventSource } from 'eventsource';

const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key-here';
const url = `${baseUrl}/rolls/subscribe?clientId=your-client-id`;

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

interface RollEvent {
  id: string;
  messageId: string;
  user: { id: string; name: string };
  speaker: any;
  flavor: string;
  rollTotal: number;
  formula: string;
  isCritical: boolean;
  isFumble: boolean;
  dice: { faces: number; results: { result: number; active: boolean }[] }[];
  timestamp: number;
}

eventSource.addEventListener('connected', (event: MessageEvent) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
});

eventSource.addEventListener('roll', (event: MessageEvent) => {
  const roll: RollEvent = JSON.parse(event.data);
  const dice = roll.dice?.map(d =>
    `${d.results.map(r => `${r.result}${r.active ? '' : '(dropped)'}`).join(', ')} (d${d.faces})`
  ).join(' + ') || '';
  console.log(`[${roll.user?.name}] ${roll.formula} = ${roll.rollTotal}${roll.isCritical ? ' CRITICAL!' : ''}${roll.isFumble ? ' FUMBLE!' : ''}`);
  if (roll.flavor) console.log(`  Flavor: ${roll.flavor}`);
  if (dice) console.log(`  Dice: ${dice}`);
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


