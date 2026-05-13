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
const baseUrl = 'http://localhost:3010';
const path = '/select';
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
curl -X POST 'http://localhost:3010/select?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"all":true,"overwrite":true}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/select'
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
  const baseUrl = 'http://localhost:3010';
  const path = '/select';
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
  3010 ➡️ port
  🔤/select🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"all":true,"overwrite":true}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /select🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 29❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "select_1777996626330",
  "success": true,
  "count": 1,
  "message": "1 entities selected",
  "selected": [
    "Scene.iI8vL6F5ett88LXH.Token.tYjgFFD4Zgjy7nZM"
  ]
}
```


---

## GET /selected

Get selected token(s)

Retrieves the currently selected token(s) in the Foundry VTT client.

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
const baseUrl = 'http://localhost:3010';
const path = '/selected';
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
curl -X GET 'http://localhost:3010/selected?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/selected'
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
  const path = '/selected';
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
  🔤/selected🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /selected🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "selected_1777996626333",
  "success": true,
  "selected": [
    {
      "tokenUuid": "Scene.iI8vL6F5ett88LXH.Token.tYjgFFD4Zgjy7nZM",
      "actorUuid": "Scene.iI8vL6F5ett88LXH.Token.tYjgFFD4Zgjy7nZM.Actor.sMD3o6zej6ckQkpo"
    }
  ]
}
```


---

## GET /players

Get players/users

Retrieves a list of all users configured in the Foundry VTT world. Useful for discovering valid userId values for permission-scoped API calls.

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
const baseUrl = 'http://localhost:3010';
const path = '/players';
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
curl -X GET 'http://localhost:3010/players?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/players'
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
  const path = '/players';
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
  🔤/players🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /players🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "players_1777996626347",
  "users": [
    {
      "id": "vXdyKYLgpmko3kHx",
      "name": "Gamemaster",
      "role": 4,
      "isGM": true,
      "active": false,
      "color": "#becc28",
      "avatar": "icons/svg/mystery-man.svg"
    },
    {
      "id": "fCfNJPT9Atc26yyv",
      "name": "tester",
      "role": 4,
      "isGM": true,
      "active": true,
      "color": "#cc28b5",
      "avatar": "icons/svg/mystery-man.svg"
    },
    {
      "id": "mCIKww86Dh7ogc9v",
      "name": "Player1",
      "role": 1,
      "isGM": false,
      "active": false,
      "color": "#cc8128",
      "avatar": "icons/svg/mystery-man.svg"
    },
    {
      "id": "0NqudQfTbKOdSgxB",
      "name": "test",
      "role": 1,
      "isGM": false,
      "active": false,
      "color": "#283dcc",
      "avatar": "icons/svg/mystery-man.svg"
    }
  ]
}
```


---

## GET /world-info

Get comprehensive world information

Returns a single object with world name, game system, Foundry version, all modules (with active status), all users (with online status), and the active scene. Useful for API clients to discover the world state.

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
const baseUrl = 'http://localhost:3010';
const path = '/world-info';
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
curl -X GET 'http://localhost:3010/world-info?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/world-info'
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
  const path = '/world-info';
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
  🔤/world-info🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /world-info🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "world-info_1777996626372",
  "data": {
    "world": {
      "id": "rest-api",
      "title": "rest-api"
    },
    "system": {
      "id": "dnd5e",
      "title": "Dungeons & Dragons Fifth Edition",
      "version": "4.3.8"
    },
    "foundryVersion": "12.331",
    "modules": [
      {
        "id": "_CodeMirror",
        "title": "CodeMirror",
        "version": "5.58.3-fvtt5",
        "active": false
      },
      {
        "id": "about-face",
        "title": "About Face",
        "version": "3.26.5",
        "active": false
      },
      {
        "id": "aedifs-token-sounds",
        "title": "Sound of Token",
        "version": "2.1.0",
        "active": false
      },
      {
        "id": "barbrawl",
        "title": "Bar Brawl",
        "version": "1.8.9",
        "active": false
      },
      {
        "id": "barbrawl-profiles",
        "title": "Bar Brawl Profiles",
        "version": "1.0.0",
        "active": false
      },
      {
        "id": "boss-splash",
        "title": "Boss Splash Screen",
        "version": "1.1.11",
        "active": false
      },
      {
        "id": "colorsettings",
        "title": "lib - Color Settings",
        "version": "3.0.3",
        "active": false
      },
      {
        "id": "combat-enhancements",
        "title": "Combat Enhancements",
        "version": "1.2.4",
        "active": false
      },
      {
        "id": "combat-tracker-dock",
        "title": "Carousel Combat Tracker",
        "version": "3.1.8",
        "active": false
      },
      {
        "id": "combatbooster",
        "title": "Combat Booster: Turn Marker, Recent Actions and more",
        "version": "4.0.4",
        "active": false
      },
      {
        "id": "condition-lab-triggler",
        "title": "Condition Lab & Triggler",
        "version": "2.0.1",
        "active": false
      },
      {
        "id": "dae",
        "title": "Dynamic effects using Active Effects (DAE)",
        "version": "12.0.29",
        "active": false
      },
      {
        "id": "dd-import",
        "title": "Universal Battlemap Importer",
        "version": "4.0.0",
        "active": true
      },
      {
        "id": "df-architect",
        "title": "DF Architect",
        "version": "5.0.0",
        "active": false
      },
      {
        "id": "dice-calculator",
        "title": "Dice Tray",
        "version": "2.27",
        "active": false
      },
      {
        "id": "dice-so-nice",
        "title": "Dice So Nice!",
        "version": "5.1.8",
        "active": false
      },
      {
        "id": "find-the-culprit",
        "title": "Find the Culprit",
        "version": "3.1.2",
        "active": false
      },
      {
        "id": "forien-copy-environment",
        "title": "Forien's Copy Environment",
        "version": "v2.2.4",
        "active": false
      },
      {
        "id": "foundry-module-bundle-m93mw9jstj6rag9m88g",
        "title": "Foundry Module Bundle",
        "version": "1.0.0",
        "active": false
      },
      {
        "id": "foundry-rest-api",
        "title": "Foundry REST API",
        "version": "3.0.2",
        "active": true
      },
      {
        "id": "foundry-server-to-server",
        "title": "Server-to-Server Transfer",
        "version": "0.1.0",
        "active": false
      },
      {
        "id": "fxmaster",
        "title": "FXMaster",
        "version": "4.1.0",
        "active": false
      },
      {
        "id": "illandril-token-hud-scale",
        "title": "Illandril's Token HUD scaler",
        "version": "3.0.0",
        "active": false
      },
      {
        "id": "image-hover",
        "title": "Image Hover",
        "version": "3.0.5",
        "active": false
      },
      {
        "id": "just-popcorn-initiative",
        "title": "Just - Popcorn Initiative",
        "version": "1.1.0",
        "active": false
      },
      {
        "id": "lib-wrapper",
        "title": "libWrapper",
        "version": "1.13.2.0",
        "active": false
      },
      {
        "id": "midi-qol",
        "title": "Midi QOL",
        "version": "12.4.64",
        "active": false
      },
      {
        "id": "monks-little-details",
        "title": "Monk's Little Details",
        "version": "12.04",
        "active": false
      },
      {
        "id": "pings",
        "title": "Pings",
        "version": "1.4.2",
        "active": false
      },
      {
        "id": "playerListStatus",
        "title": "Playerlist Status - Lib for View status on Playerlist",
        "version": "3.1.3",
        "active": false
      },
      {
        "id": "playerStatus",
        "title": "Player Status",
        "version": "3.1.4",
        "active": false
      },
      {
        "id": "quick-insert",
        "title": "Quick Insert - Search Widget",
        "version": "3.6.0",
        "active": true
      },
      {
        "id": "rounded-measurement-templates",
        "title": "Rounded Measurement Templates",
        "version": "1.0.0",
        "active": false
      },
      {
        "id": "settings-extender",
        "title": "Settings Extender",
        "version": "1.2.3",
        "active": false
      },
      {
        "id": "sidebar-resizer",
        "title": "Sidebar and Windows Resizer",
        "version": "0.8.0",
        "active": false
      },
      {
        "id": "smarttarget",
        "title": "Smart Target",
        "version": "2.0.0",
        "active": false
      },
      {
        "id": "socketlib",
        "title": "socketlib",
        "version": "1.1.2",
        "active": false
      },
      {
        "id": "soncraft-audio",
        "title": "Soncraft Audio Module",
        "version": "0.2.1",
        "active": false
      },
      {
        "id": "sorted-status-effects",
        "title": "Sorted Status Effects",
        "version": "#{VERSION}#",
        "active": false
      },
      {
        "id": "splatter",
        "title": "Splatter",
        "version": "4.0.5",
        "active": false
      },
      {
        "id": "spotlight-omnisearch",
        "title": "Spotlight Omnisearch",
        "version": "2.0.7",
        "active": false
      },
      {
        "id": "statuscounter",
        "title": "Status Icon Counters",
        "version": "3.0.2",
        "active": false
      },
      {
        "id": "token-aura-ring",
        "title": "Token Aura Ring",
        "version": "2.8.1",
        "active": false
      },
      {
        "id": "token-ease",
        "title": "Token Ease",
        "version": "1.3.0",
        "active": false
      },
      {
        "id": "token-mold",
        "title": "Token Mold",
        "version": "2.22.3",
        "active": false
      },
      {
        "id": "token-tooltip-alt",
        "title": "Token Tooltip Alt",
        "version": "3.3.4",
        "active": false
      },
      {
        "id": "token-variants",
        "title": "Token Variant Art",
        "version": "5.1.6",
        "active": false
      },
      {
        "id": "visual-active-effects",
        "title": "Visual Active Effects",
        "version": "12.1.5",
        "active": false
      },
      {
        "id": "what-the-flag",
        "title": "WTF: What the Flag?",
        "version": "1.0.5",
        "active": false
      }
    ],
    "users": [
      {
        "id": "vXdyKYLgpmko3kHx",
        "name": "Gamemaster",
        "role": 4,
        "isGM": true,
        "active": false,
        "color": "#becc28",
        "avatar": "icons/svg/mystery-man.svg"
      },
      {
        "id": "fCfNJPT9Atc26yyv",
        "name": "tester",
        "role": 4,
        "isGM": true,
        "active": true,
        "color": "#cc28b5",
        "avatar": "icons/svg/mystery-man.svg"
      },
      {
        "id": "mCIKww86Dh7ogc9v",
        "name": "Player1",
        "role": 1,
        "isGM": false,
        "active": false,
        "color": "#cc8128",
        "avatar": "icons/svg/mystery-man.svg"
      },
      {
        "id": "0NqudQfTbKOdSgxB",
        "name": "test",
        "role": 1,
        "isGM": false,
        "active": false,
        "color": "#283dcc",
        "avatar": "icons/svg/mystery-man.svg"
      }
    ],
    "activeScene": {
      "id": "iI8vL6F5ett88LXH",
      "name": "test-scene-updated"
    }
  }
}
```


---

## POST /execute-js

Execute JavaScript

Executes a JavaScript script in the Foundry VTT client.

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
const baseUrl = 'http://localhost:3010';
const path = '/execute-js';
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
      "script": "const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/execute-js?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"script":"const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/execute-js'
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
  const baseUrl = 'http://localhost:3010';
  const path = '/execute-js';
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
  3010 ➡️ port
  🔤/execute-js🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"script":"const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /execute-js🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 104❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "execute-js_1777996626348",
  "success": true,
  "result": "ws://localhost:3010/relay"
}
```


