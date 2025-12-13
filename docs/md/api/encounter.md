---
tag: encounter
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


# encounter

## GET /encounters

Get all active encounters Retrieves a list of all currently active encounters in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |

### Returns

**object** - An array of active encounters with details

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/encounters';
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
curl -X GET 'http://localhost:3010/encounters?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/encounters'
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
  const path = '/encounters';
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
  🔤/encounters🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /encounters🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "encounters_1765658033112",
  "clientId": "your-client-id",
  "type": "encounters-result",
  "encounters": [
    {
      "id": "7r5lECrTxEGQS3Kq",
      "round": 1,
      "turn": 0,
      "current": true,
      "combatants": [
        {
          "id": "v2xjIf0WCdHx9ZLz",
          "name": "test",
          "tokenUuid": "Scene.NUEDEFAULTSCENE0.Token.O4sEnBrG5I3lFNGk",
          "actorUuid": "Actor.xctPu6799LkAP6p3",
          "img": "icons/svg/mystery-man.svg",
          "initiative": 20,
          "hidden": false,
          "defeated": false
        }
      ]
    }
  ]
}
```


---

## POST /start-encounter

Start a new encounter Initiates a new encounter in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| tokens | array |  | body | Array of token UUIDs to include in the encounter |
| startWithSelected | boolean |  | body | Whether to start with selected tokens |
| startWithPlayers | boolean |  | body | Whether to start with players |
| rollNPC | boolean |  | body | Whether to roll for NPCs |
| rollAll | boolean |  | body | Whether to roll for all tokens |
| name | string |  | body | The name of the encounter (unused) |

### Returns

**object** - Details of the started encounter

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/start-encounter';
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
      "startWithSelected": true,
      "rollAll": true
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/start-encounter?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"startWithSelected":true,"rollAll":true}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/start-encounter'
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
      "startWithSelected": True,
      "rollAll": True
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
  const path = '/start-encounter';
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
        "startWithSelected": true,
        "rollAll": true
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
  🔤/start-encounter🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"startWithSelected":true,"rollAll":true}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /start-encounter🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 41❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "start-encounter_1765658031424",
  "clientId": "your-client-id",
  "type": "start-encounter-result",
  "encounterId": "7r5lECrTxEGQS3Kq",
  "encounter": {
    "id": "7r5lECrTxEGQS3Kq",
    "round": 1,
    "turn": 0,
    "combatants": [
      {
        "id": "v2xjIf0WCdHx9ZLz",
        "name": "test",
        "tokenUuid": "Scene.NUEDEFAULTSCENE0.Token.O4sEnBrG5I3lFNGk",
        "actorUuid": "Actor.xctPu6799LkAP6p3",
        "img": "icons/svg/mystery-man.svg",
        "initiative": 20,
        "hidden": false,
        "defeated": false
      }
    ]
  }
}
```


---

## POST /next-turn

Advance to the next turn in the encounter Moves the encounter to the next turn.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| encounter | string |  | query, body | The ID of the encounter to advance (optional, defaults to current encounter) |

### Returns

**object** - Details of the next turn

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/next-turn';
const params = {
  clientId: 'your-client-id',
  encounterId: '7r5lECrTxEGQS3Kq'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'POST',
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
curl -X POST 'http://localhost:3010/next-turn?clientId=your-client-id&encounterId=7r5lECrTxEGQS3Kq' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/next-turn'
params = {
    'clientId': 'your-client-id',
    'encounterId': '7r5lECrTxEGQS3Kq'
}
url = f'{base_url}{path}'

response = requests.post(
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
  const path = '/next-turn';
  const params = {
    clientId: 'your-client-id',
    encounterId: '7r5lECrTxEGQS3Kq'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'post',
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
  🔤/next-turn🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤encounterId=7r5lECrTxEGQS3Kq🔤 ➡️ encounterId
  🔤?🧲clientId🧲&🧲encounterId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤POST /next-turn🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "next-turn_1765658033388",
  "clientId": "your-client-id",
  "type": "next-turn-result",
  "encounterId": "7r5lECrTxEGQS3Kq",
  "action": "nextTurn",
  "currentTurn": 0,
  "currentRound": 2,
  "actorTurn": "Actor.xctPu6799LkAP6p3",
  "tokenTurn": "Scene.NUEDEFAULTSCENE0.Token.O4sEnBrG5I3lFNGk",
  "encounter": {
    "id": "7r5lECrTxEGQS3Kq",
    "round": 2,
    "turn": 0
  }
}
```


---

## POST /next-round

Advance to the next round in the encounter Moves the encounter to the next round.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| encounter | string |  | query, body | The ID of the encounter to advance (optional, defaults to current encounter) |

### Returns

**object** - Details of the next round

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/next-round';
const params = {
  clientId: 'your-client-id',
  encounterId: '7r5lECrTxEGQS3Kq'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'POST',
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
curl -X POST 'http://localhost:3010/next-round?clientId=your-client-id&encounterId=7r5lECrTxEGQS3Kq' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/next-round'
params = {
    'clientId': 'your-client-id',
    'encounterId': '7r5lECrTxEGQS3Kq'
}
url = f'{base_url}{path}'

response = requests.post(
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
  const path = '/next-round';
  const params = {
    clientId: 'your-client-id',
    encounterId: '7r5lECrTxEGQS3Kq'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'post',
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
  🔤/next-round🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤encounterId=7r5lECrTxEGQS3Kq🔤 ➡️ encounterId
  🔤?🧲clientId🧲&🧲encounterId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤POST /next-round🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "next-round_1765658034163",
  "clientId": "your-client-id",
  "type": "next-round-result",
  "encounterId": "7r5lECrTxEGQS3Kq",
  "action": "nextRound",
  "currentTurn": 0,
  "currentRound": 3,
  "actorTurn": "Actor.xctPu6799LkAP6p3",
  "tokenTurn": "Scene.NUEDEFAULTSCENE0.Token.O4sEnBrG5I3lFNGk",
  "encounter": {
    "id": "7r5lECrTxEGQS3Kq",
    "round": 3,
    "turn": 0
  }
}
```


---

## POST /last-turn

Advance to the last turn in the encounter Moves the encounter to the last turn.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| encounter | string |  | query, body | The ID of the encounter to advance (optional, defaults to current encounter) |

### Returns

**object** - Details of the last turn

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/last-turn';
const params = {
  clientId: 'your-client-id',
  encounterId: '7r5lECrTxEGQS3Kq'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'POST',
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
curl -X POST 'http://localhost:3010/last-turn?clientId=your-client-id&encounterId=7r5lECrTxEGQS3Kq' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/last-turn'
params = {
    'clientId': 'your-client-id',
    'encounterId': '7r5lECrTxEGQS3Kq'
}
url = f'{base_url}{path}'

response = requests.post(
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
  const path = '/last-turn';
  const params = {
    clientId: 'your-client-id',
    encounterId: '7r5lECrTxEGQS3Kq'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'post',
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
  🔤/last-turn🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤encounterId=7r5lECrTxEGQS3Kq🔤 ➡️ encounterId
  🔤?🧲clientId🧲&🧲encounterId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤POST /last-turn🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "last-turn_1765658034949",
  "clientId": "your-client-id",
  "type": "last-turn-result",
  "encounterId": "7r5lECrTxEGQS3Kq",
  "action": "previousTurn",
  "currentTurn": 0,
  "currentRound": 2,
  "actorTurn": "Actor.xctPu6799LkAP6p3",
  "tokenTurn": "Scene.NUEDEFAULTSCENE0.Token.O4sEnBrG5I3lFNGk",
  "encounter": {
    "id": "7r5lECrTxEGQS3Kq",
    "round": 2,
    "turn": 0
  }
}
```


---

## POST /last-round

Advance to the last round in the encounter Moves the encounter to the last round.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| encounter | string |  | query, body | The ID of the encounter to advance (optional, defaults to current encounter) |

### Returns

**object** - Details of the last round

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/last-round';
const params = {
  clientId: 'your-client-id',
  encounterId: '7r5lECrTxEGQS3Kq'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'POST',
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
curl -X POST 'http://localhost:3010/last-round?clientId=your-client-id&encounterId=7r5lECrTxEGQS3Kq' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/last-round'
params = {
    'clientId': 'your-client-id',
    'encounterId': '7r5lECrTxEGQS3Kq'
}
url = f'{base_url}{path}'

response = requests.post(
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
  const path = '/last-round';
  const params = {
    clientId: 'your-client-id',
    encounterId: '7r5lECrTxEGQS3Kq'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'post',
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
  🔤/last-round🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤encounterId=7r5lECrTxEGQS3Kq🔤 ➡️ encounterId
  🔤?🧲clientId🧲&🧲encounterId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤POST /last-round🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "last-round_1765658035717",
  "clientId": "your-client-id",
  "type": "last-round-result",
  "encounterId": "7r5lECrTxEGQS3Kq",
  "action": "previousRound",
  "currentTurn": 0,
  "currentRound": 1,
  "actorTurn": "Actor.xctPu6799LkAP6p3",
  "tokenTurn": "Scene.NUEDEFAULTSCENE0.Token.O4sEnBrG5I3lFNGk",
  "encounter": {
    "id": "7r5lECrTxEGQS3Kq",
    "round": 1,
    "turn": 0
  }
}
```


---

## POST /end-encounter

End an encounter Ends the current encounter in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| encounter | string |  | query, body | The ID of the encounter to end (optional, defaults to current encounter) |

### Returns

**object** - Details of the ended encounter

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/end-encounter';
const params = {
  clientId: 'your-client-id',
  encounterId: '7r5lECrTxEGQS3Kq'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'POST',
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
curl -X POST 'http://localhost:3010/end-encounter?clientId=your-client-id&encounterId=7r5lECrTxEGQS3Kq' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/end-encounter'
params = {
    'clientId': 'your-client-id',
    'encounterId': '7r5lECrTxEGQS3Kq'
}
url = f'{base_url}{path}'

response = requests.post(
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
  const path = '/end-encounter';
  const params = {
    clientId: 'your-client-id',
    encounterId: '7r5lECrTxEGQS3Kq'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'post',
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
  🔤/end-encounter🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤encounterId=7r5lECrTxEGQS3Kq🔤 ➡️ encounterId
  🔤?🧲clientId🧲&🧲encounterId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤POST /end-encounter🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "end-encounter_1765658037498",
  "clientId": "your-client-id",
  "type": "end-encounter-result",
  "encounterId": "7r5lECrTxEGQS3Kq",
  "message": "Encounter successfully ended"
}
```


---

## POST /add-to-encounter

Add tokens to an encounter Adds selected tokens or specified UUIDs to the current encounter.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| encounter | string |  | query, body | The ID of the encounter to add tokens to (optional, defaults to current encounter) |
| selected | boolean |  | body | Whether to add selected tokens (optional, defaults to false) |
| uuids | array |  | body | The UUIDs of the tokens to add (optional, defaults to empty array) |
| rollInitiative | boolean |  | body | Whether to roll initiative for the added tokens (optional, defaults to false) |

### Returns

**object** - Details of the updated encounter

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/add-to-encounter';
const params = {
  clientId: 'your-client-id',
  encounterId: '7r5lECrTxEGQS3Kq'
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
      "selected": true,
      "uuids": [],
      "rollInitiative": true
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/add-to-encounter?clientId=your-client-id&encounterId=7r5lECrTxEGQS3Kq' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"selected":true,"uuids":[],"rollInitiative":true}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/add-to-encounter'
params = {
    'clientId': 'your-client-id',
    'encounterId': '7r5lECrTxEGQS3Kq'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "selected": True,
      "uuids": [],
      "rollInitiative": True
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
  const path = '/add-to-encounter';
  const params = {
    clientId: 'your-client-id',
    encounterId: '7r5lECrTxEGQS3Kq'
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
        "selected": true,
        "uuids": [],
        "rollInitiative": true
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
  🔤/add-to-encounter🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤encounterId=7r5lECrTxEGQS3Kq🔤 ➡️ encounterId
  🔤?🧲clientId🧲&🧲encounterId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"selected":true,"uuids":[],"rollInitiative":true}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /add-to-encounter🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 50❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "add-to-encounter_1765658036978",
  "clientId": "your-client-id",
  "type": "add-to-encounter-result",
  "encounterId": "7r5lECrTxEGQS3Kq",
  "added": [
    "Scene.NUEDEFAULTSCENE0.Token.O4sEnBrG5I3lFNGk"
  ],
  "failed": []
}
```


---

## POST /remove-from-encounter

Remove tokens from an encounter Removes selected tokens or specified UUIDs from the current encounter.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | The ID of the Foundry client to connect to |
| encounter | string |  | query, body | The ID of the encounter to remove tokens from (optional, defaults to current encounter) |
| selected | boolean |  | body | Whether to remove selected tokens (optional, defaults to false) |
| uuids | array |  | body | The UUIDs of the tokens to remove (optional, defaults to empty array) |

### Returns

**object** - Details of the updated encounter

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/remove-from-encounter';
const params = {
  clientId: 'your-client-id',
  encounterId: '7r5lECrTxEGQS3Kq'
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
      "selected": true
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/remove-from-encounter?clientId=your-client-id&encounterId=7r5lECrTxEGQS3Kq' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"selected":true}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/remove-from-encounter'
params = {
    'clientId': 'your-client-id',
    'encounterId': '7r5lECrTxEGQS3Kq'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "selected": True
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
  const path = '/remove-from-encounter';
  const params = {
    clientId: 'your-client-id',
    encounterId: '7r5lECrTxEGQS3Kq'
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
        "selected": true
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
  🔤/remove-from-encounter🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤encounterId=7r5lECrTxEGQS3Kq🔤 ➡️ encounterId
  🔤?🧲clientId🧲&🧲encounterId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"selected":true}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /remove-from-encounter🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 17❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "remove-from-encounter_1765658036463",
  "clientId": "your-client-id",
  "type": "remove-from-encounter-result",
  "encounterId": "7r5lECrTxEGQS3Kq",
  "removed": [
    "Scene.NUEDEFAULTSCENE0.Token.O4sEnBrG5I3lFNGk"
  ],
  "failed": []
}
```


