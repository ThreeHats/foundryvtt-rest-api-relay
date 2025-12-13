---
tag: roll
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


# roll

## GET /rolls

Get recent rolls Retrieves a list of up to 20 recent rolls made in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| limit | number |  | query | Optional limit on the number of rolls to return (default is 20) |

### Returns

**object** - An array of recent rolls with details

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/rolls';
const params = {
  clientId: 'your-client-id',
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
curl -X GET 'http://localhost:3010/rolls?clientId=your-client-id&limit=20' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/rolls'
params = {
    'clientId': 'your-client-id',
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
    clientId: 'your-client-id',
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
  🔤clientId=your-client-id🔤 ➡️ clientId
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
  "requestId": "rolls_1765635982518",
  "clientId": "your-client-id",
  "type": "rolls-result",
  "data": [
    {
      "id": "RnuRRnuolZDORCSf",
      "messageId": "RnuRRnuolZDORCSf",
      "user": {
        "id": "r6bXhB7k9cXa3cif",
        "name": "tester"
      },
      "speaker": {
        "scene": null,
        "actor": null,
        "token": null
      },
      "flavor": "Test Roll",
      "rollTotal": 17,
      "formula": "2d20kh",
      "isCritical": false,
      "isFumble": false,
      "dice": [
        {
          "faces": 20,
          "results": [
            {
              "result": 3,
              "active": false
            },
            {
              "result": 17,
              "active": true
            }
          ]
        }
      ],
      "timestamp": 1765635982506
    }
  ]
}
```


---

## GET /lastroll

Get the last roll Retrieves the most recent roll made in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |

### Returns

**object** - The most recent roll with details

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/lastroll';
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
curl -X GET 'http://localhost:3010/lastroll?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/lastroll'
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
  const path = '/lastroll';
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
  🔤/lastroll🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
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
  "requestId": "last-roll_1765635982792",
  "clientId": "your-client-id",
  "type": "last-roll-result",
  "data": {
    "id": "RnuRRnuolZDORCSf",
    "messageId": "RnuRRnuolZDORCSf",
    "user": {
      "id": "r6bXhB7k9cXa3cif",
      "name": "tester"
    },
    "speaker": {
      "scene": null,
      "actor": null,
      "token": null
    },
    "flavor": "Test Roll",
    "rollTotal": 17,
    "formula": "2d20kh",
    "isCritical": false,
    "isFumble": false,
    "dice": [
      {
        "faces": 20,
        "results": [
          {
            "result": 3,
            "active": false
          },
          {
            "result": 17,
            "active": true
          }
        ]
      }
    ],
    "timestamp": 1765635982506
  }
}
```


---

## POST /roll

Make a roll Executes a roll with the specified formula

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| formula | string | ✓ | body | The roll formula to evaluate (e.g., "1d20 + 5") |
| flavor | string |  | body | Optional flavor text for the roll |
| createChatMessage | boolean |  | body | Whether to create a chat message for the roll |
| speaker | string |  | body | The speaker for the roll |
| whisper | array |  | body | Users to whisper the roll result to |

### Returns

**object** - Result of the roll operation

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/roll';
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
curl -X POST 'http://localhost:3010/roll?clientId=your-client-id' \
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
  🔤clientId=your-client-id🔤 ➡️ clientId
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
  "requestId": "roll_1765635981991",
  "clientId": "your-client-id",
  "type": "roll-result",
  "success": true,
  "data": {
    "id": "manual_1765635982507_cei007er5vs",
    "chatMessageCreated": true,
    "roll": {
      "formula": "2d20kh",
      "total": 17,
      "isCritical": false,
      "isFumble": false,
      "dice": [
        {
          "faces": 20,
          "results": [
            {
              "result": 3,
              "active": false
            },
            {
              "result": 17,
              "active": true
            }
          ]
        }
      ],
      "timestamp": 1765635982507
    }
  }
}
```

