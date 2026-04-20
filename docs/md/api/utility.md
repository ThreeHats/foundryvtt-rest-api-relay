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
curl -X POST 'http://localhost:3010/select?clientId=fvtt_099ad17ea199e7e3' \
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
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
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
  "requestId": "select_1776657990079",
  "success": true,
  "count": 1,
  "message": "1 entities selected",
  "selected": [
    "Scene.r36nfimJGHYGUGQX.Token.Q9y6lc2dPYs2WRQT"
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
curl -X GET 'http://localhost:3010/selected?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/selected'
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
  const path = '/selected';
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
  🔤/selected🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
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
  "requestId": "selected_1776657990085",
  "success": true,
  "selected": [
    {
      "tokenUuid": "Scene.r36nfimJGHYGUGQX.Token.Q9y6lc2dPYs2WRQT",
      "actorUuid": "Scene.r36nfimJGHYGUGQX.Token.Q9y6lc2dPYs2WRQT.Actor.q9uWyfdPwTlzbpxb"
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
curl -X GET 'http://localhost:3010/players?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/players'
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
  const path = '/players';
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
  🔤/players🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
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
  "requestId": "players_1776657990090",
  "users": [
    {
      "id": "5ypAoBvOiyjDKiaZ",
      "name": "Gamemaster",
      "role": 4,
      "isGM": true,
      "active": false,
      "color": "#28cca2",
      "avatar": "icons/svg/mystery-man.svg"
    },
    {
      "id": "r6bXhB7k9cXa3cif",
      "name": "tester",
      "role": 4,
      "isGM": true,
      "active": true,
      "color": "#cc2829",
      "avatar": "icons/svg/mystery-man.svg"
    },
    {
      "id": "JLfKTeTgCDpAdDfw",
      "name": "some-cool-guy",
      "role": 1,
      "isGM": false,
      "active": false,
      "color": "#6328cc",
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
curl -X GET 'http://localhost:3010/world-info?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/world-info'
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
  const path = '/world-info';
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
  🔤/world-info🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
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
  "requestId": "world-info_1776657990130",
  "data": {
    "world": {
      "id": "testing",
      "title": "testing"
    },
    "system": {
      "id": "dnd5e",
      "title": "Dungeons & Dragons Fifth Edition",
      "version": "5.0.4"
    },
    "foundryVersion": "13.348",
    "modules": [
      {
        "id": "5e-world-currency",
        "title": "World Currency for D&D 5e",
        "version": "4.0.6",
        "active": false
      },
      {
        "id": "ATL",
        "title": "Active Token Effects",
        "version": "v0.8.1",
        "active": false
      },
      {
        "id": "ActiveAuras",
        "title": "Active Auras",
        "version": "0.12.3",
        "active": false
      },
      {
        "id": "LockView",
        "title": "Lock View",
        "version": "2.0.1",
        "active": false
      },
      {
        "id": "Rideable",
        "title": "Rideable",
        "version": "4.0.7",
        "active": false
      },
      {
        "id": "_CodeMirror",
        "title": "CodeMirror",
        "version": "5.58.3-fvtt5",
        "active": false
      },
      {
        "id": "alternative-rotation",
        "title": "Alternative Rotation",
        "version": "2.3.1",
        "active": false
      },
      {
        "id": "always-hp",
        "title": "Always HP",
        "version": "13.02",
        "active": false
      },
      {
        "id": "arbron-hp-bar",
        "title": "Arbron’s Improved HP Bar",
        "version": "2.2.7",
        "active": false
      },
      {
        "id": "arms-reach",
        "title": "Arms Reach",
        "version": "13.0.1",
        "active": false
      },
      {
        "id": "autoanimations",
        "title": "Automated Animations",
        "version": "6.2.3",
        "active": false
      },
      {
        "id": "backgroundless-pins",
        "title": "Backgroundless Pins",
        "version": "3.0.0",
        "active": false
      },
      {
        "id": "baileywiki-cabal-dungeon",
        "title": "Baileywiki Cabal Dungeon",
        "version": "0.3.3",
        "active": false
      },
      {
        "id": "baileywiki-city-district-01",
        "title": "Baileywiki Modular City - District 01",
        "version": "0.3.5",
        "active": false
      },
      {
        "id": "baileywiki-city-district-02-docks",
        "title": "Baileywiki Modular City - District 02: Docks",
        "version": "0.3.1",
        "active": false
      },
      {
        "id": "baileywiki-city-district-03-temples",
        "title": "Baileywiki Modular City - District 03: Temples",
        "version": "0.3.0",
        "active": false
      },
      {
        "id": "baileywiki-city-district-04-ruins-and-slums",
        "title": "Baileywiki Modular City - District 04: Ruins and Slums",
        "version": "0.3.1",
        "active": false
      },
      {
        "id": "baileywiki-city-district-05-affluent",
        "title": "Baileywiki Modular City - District 05: Affluent",
        "version": "0.3.0",
        "active": false
      },
      {
        "id": "baileywiki-maps-premium",
        "title": "Baileywiki Maps Premium",
        "version": "0.6.1",
        "active": false
      },
      {
        "id": "baileywiki-maps-premium-towns",
        "title": "Baileywiki Maps Towns",
        "version": "0.6.12",
        "active": false
      },
      {
        "id": "baileywiki-nuts-and-bolts",
        "title": "Baileywiki Nuts and Bolts",
        "version": "0.4.1",
        "active": false
      },
      {
        "id": "barbrawl",
        "title": "Bar Brawl",
        "version": "1.8.12",
        "active": false
      },
      {
        "id": "betterroofs",
        "title": "Better Roofs",
        "version": "3.0.0",
        "active": false
      },
      {
        "id": "boss-loot-assets-premium",
        "title": "BLFX Assets & Animation Editor Premium",
        "version": "2.7.2",
        "active": false
      },
      {
        "id": "caydens-dice",
        "title": "Cayden's Dice",
        "version": "1.1.3",
        "active": false
      },
      {
        "id": "character-actions-list-5e",
        "title": "Character Actions List dnd5e",
        "version": "7.1.1",
        "active": false
      },
      {
        "id": "chat-media",
        "title": "Chat Media",
        "version": "13.0.1",
        "active": false
      },
      {
        "id": "cleaner-sheet-title-bar",
        "title": "Cleaner sheet title bar",
        "version": "1.5.5",
        "active": false
      },
      {
        "id": "color-picker",
        "title": "Color Picker",
        "version": "1.7",
        "active": false
      },
      {
        "id": "colorsettings",
        "title": "lib - Color Settings",
        "version": "3.0.4",
        "active": false
      },
      {
        "id": "combat-tracker-dock",
        "title": "Carousel Combat Tracker",
        "version": "4.0.2",
        "active": false
      },
      {
        "id": "combatready",
        "title": "Combat Ready!",
        "version": "5.0.2",
        "active": false
      },
      {
        "id": "compact-scene-navigation",
        "title": "Compact Scene Navigation",
        "version": "1.0.2",
        "active": false
      },
      {
        "id": "cosmic-dice",
        "title": "Cosmic Dice",
        "version": "1.2.2",
        "active": false
      },
      {
        "id": "custom-css",
        "title": "Custom CSS",
        "version": "2.4.3",
        "active": false
      },
      {
        "id": "dae",
        "title": "Dynamic effects using Active Effects (DAE)",
        "version": "13.0.11",
        "active": false
      },
      {
        "id": "dark-mode-5e",
        "title": "Zeta's Dark Mode for DnD 5th Edition",
        "version": "2.3.5",
        "active": false
      },
      {
        "id": "dd-import",
        "title": "Universal Battlemap Importer",
        "version": "5.0.0",
        "active": false
      },
      {
        "id": "dfreds-convenient-effects",
        "title": "DFreds Convenient Effects",
        "version": "8.1.1",
        "active": false
      },
      {
        "id": "dice-calculator",
        "title": "Dice Tray",
        "version": "3.4.6",
        "active": false
      },
      {
        "id": "dice-so-nice",
        "title": "Dice So Nice!",
        "version": "5.2.0",
        "active": false
      },
      {
        "id": "dnd-randomizer",
        "title": "Stochastic, Fantastic! - Random Encounter Generator",
        "version": "1.0.67",
        "active": false
      },
      {
        "id": "dnd5e-animations",
        "title": "D&D5e Animations",
        "version": "3.1.1",
        "active": false
      },
      {
        "id": "dnd5e-character-monitor",
        "title": "D&D5e Character Monitor",
        "version": "2.1.0",
        "active": false
      },
      {
        "id": "drag-ruler",
        "title": "Drag Ruler",
        "version": "1.13.7",
        "active": false
      },
      {
        "id": "dungeon-draw",
        "title": "Dungeon Draw",
        "version": "4.0.0",
        "active": false
      },
      {
        "id": "elevation-drag-ruler",
        "title": "DnD5e Drag Ruler Integration",
        "version": "1.11.0",
        "active": false
      },
      {
        "id": "enhancedcombathud",
        "title": "Argon - Combat HUD (CORE)",
        "version": "4.0.7",
        "active": false
      },
      {
        "id": "enhancedcombathud-dnd5e",
        "title": "Argon - Combat HUD (DND5E)",
        "version": "5.0.1",
        "active": false
      },
      {
        "id": "filepicker-plus",
        "title": "Filepicker Plus",
        "version": "5.0.5",
        "active": false
      },
      {
        "id": "find-the-culprit",
        "title": "Find the Culprit",
        "version": "3.2.2",
        "active": false
      },
      {
        "id": "forien-ammo-swapper",
        "title": "Forien's Ammo Swapper",
        "version": "0.4.1",
        "active": false
      },
      {
        "id": "forien-copy-environment",
        "title": "Forien's Copy Environment",
        "version": "v3.0.0",
        "active": false
      },
      {
        "id": "foundry-rest-api",
        "title": "Foundry REST API",
        "version": "3.0.0",
        "active": true
      },
      {
        "id": "foundry-server-to-server",
        "title": "Server-to-Server Transfer",
        "version": "0.1.0",
        "active": true
      },
      {
        "id": "foundry-tables-to-items",
        "title": "Foundry Tables to Items",
        "version": "1.0.0",
        "active": true
      },
      {
        "id": "foundry-taskbar",
        "title": "Taskbar",
        "version": "5.0.5",
        "active": false
      },
      {
        "id": "foundry-vtt-content-parser",
        "title": "Foundry VTT Content Parser",
        "version": "0.2.50",
        "active": false
      },
      {
        "id": "foundry_community_tables",
        "title": "Foundry Community Tables",
        "version": "13.0",
        "active": false
      },
      {
        "id": "foundryvtt-simple-calendar",
        "title": "Simple Calendar",
        "version": "2.4.18",
        "active": false
      },
      {
        "id": "fvtt-paper-doll-ui",
        "title": "Paper Doll",
        "version": "2.0.1",
        "active": false
      },
      {
        "id": "gatherer",
        "title": "Gatherer",
        "version": "4.0.2",
        "active": false
      },
      {
        "id": "give-item",
        "title": "Give item to another player",
        "version": "2.1.1",
        "active": false
      },
      {
        "id": "gm-notes",
        "title": "GM Notes",
        "version": "10.5",
        "active": false
      },
      {
        "id": "gm-vision",
        "title": "GM Vision",
        "version": "2.0.2",
        "active": false
      },
      {
        "id": "group-initiative",
        "title": "Group Initiative",
        "version": "2.1.7",
        "active": false
      },
      {
        "id": "hexplorer",
        "title": "Hexplorer",
        "version": "2.0.4",
        "active": false
      },
      {
        "id": "image-context",
        "title": "Image Context",
        "version": "3.1.3",
        "active": false
      },
      {
        "id": "item-piles",
        "title": "Item Piles",
        "version": "3.2.20",
        "active": false
      },
      {
        "id": "itemacro",
        "title": "Item Macro",
        "version": "3.0.1",
        "active": false
      },
      {
        "id": "jb2a_patreon",
        "title": "JB2A - Patreon Complete Collection",
        "version": "0.8.1",
        "active": false
      },
      {
        "id": "journal-icon-numbers",
        "title": "Automatic Journal Icon Numbers",
        "version": "v1.10.0",
        "active": false
      },
      {
        "id": "legend-lore",
        "title": "Legend Lore",
        "version": "v1.0.2",
        "active": false
      },
      {
        "id": "levels",
        "title": "Levels",
        "version": "6.0.16",
        "active": false
      },
      {
        "id": "levelsautocover",
        "title": "AutoCover - Automatic Cover Calculator",
        "version": "3.0.2",
        "active": false
      },
      {
        "id": "lib-dfreds-ui-extender",
        "title": "Lib: DFreds UI Extender",
        "version": "2.1.2",
        "active": false
      },
      {
        "id": "lib-find-the-path-12",
        "title": "Library - Find the Path",
        "version": "2.0.5",
        "active": false
      },
      {
        "id": "lib-wrapper",
        "title": "libWrapper",
        "version": "1.13.3.0",
        "active": false
      },
      {
        "id": "limits",
        "title": "Limits",
        "version": "2.0.5",
        "active": false
      },
      {
        "id": "lootsheet-simple",
        "title": "Loot Sheet NPC 5e (Simple Version)",
        "version": "13.504.1",
        "active": false
      },
      {
        "id": "lordudice",
        "title": "Lordu's Custom Dice for Dice So Nice",
        "version": "0.41",
        "active": false
      },
      {
        "id": "macro-wheel",
        "title": "Macro Wheel",
        "version": "3.0.0",
        "active": false
      },
      {
        "id": "magicitems",
        "title": "Magic Items",
        "version": "5.0.0-beta.2",
        "active": false
      },
      {
        "id": "mastercrafted",
        "title": "Mastercrafted - Crafting Manager",
        "version": "4.0.5",
        "active": false
      },
      {
        "id": "midi-qol",
        "title": "Midi QOL",
        "version": "13.0.15",
        "active": false
      },
      {
        "id": "mkah-compendium-importer",
        "title": "Mana's Compendium Importer",
        "version": "3.0.1",
        "active": false
      },
      {
        "id": "mmm",
        "title": "Maxwell's Manual of Malicious Maladies",
        "version": "5.0.1",
        "active": false
      },
      {
        "id": "mob-attack-tool",
        "title": "Mob Attack Tool",
        "version": "0.13.0",
        "active": false
      },
      {
        "id": "module-credits",
        "title": "Module Management+",
        "version": "2.2.3",
        "active": false
      },
      {
        "id": "monks-active-tiles",
        "title": "Monk's Active Tile Triggers",
        "version": "13.05",
        "active": false
      },
      {
        "id": "monks-bloodsplats",
        "title": "Monk's Bloodsplats",
        "version": "13.01",
        "active": false
      },
      {
        "id": "monks-combat-marker",
        "title": "Monk's Combat Marker",
        "version": "12.01",
        "active": false
      },
      {
        "id": "monks-enhanced-journal",
        "title": "Monk's Enhanced Journal",
        "version": "12.02",
        "active": false
      },
      {
        "id": "monks-player-settings",
        "title": "Monk's Player Settings",
        "version": "12.01",
        "active": false
      },
      {
        "id": "monks-tokenbar",
        "title": "Monk's TokenBar",
        "version": "12.04",
        "active": false
      },
      {
        "id": "monks-wall-enhancement",
        "title": "Monk's Wall Enhancement",
        "version": "13.02",
        "active": false
      },
      {
        "id": "mookAI-12",
        "title": "mookAI - An AI for your mooks",
        "version": "1.0.5",
        "active": false
      },
      {
        "id": "moulinette",
        "title": "Moulinette Media Search",
        "version": "1.8.3",
        "active": false
      },
      {
        "id": "moulinette-soundboards",
        "title": "Soundboard & Soundpad (by Moulinette)",
        "version": "1.2.3",
        "active": false
      },
      {
        "id": "multi-token-edit",
        "title": "Baileywiki Mass Edit",
        "version": "2.1.1",
        "active": false
      },
      {
        "id": "multiface-tiles",
        "title": "Multiface Tiles",
        "version": "13.002",
        "active": false
      },
      {
        "id": "multiple-document-selection",
        "title": "Multiple Document Selection",
        "version": "13.02",
        "active": false
      },
      {
        "id": "neutro-dice-theme-bitd",
        "title": "Blades in the Dark theme for Dice So Nice",
        "version": "1.1.0",
        "active": false
      },
      {
        "id": "party-monitor-dock",
        "title": "Party HUD",
        "version": "3.0.1",
        "active": false
      },
      {
        "id": "patrol",
        "title": "Patrol",
        "version": "3.0.2",
        "active": false
      },
      {
        "id": "pdf-pager",
        "title": "PDF Pager",
        "version": "13.0.1",
        "active": false
      },
      {
        "id": "permission_viewer",
        "title": "Ownership Viewer",
        "version": "13.1",
        "active": true
      },
      {
        "id": "plutonium",
        "title": "Plutonium",
        "version": "3.11.0-noble-prerelease-25-36",
        "active": false
      },
      {
        "id": "poi-teleport",
        "title": "Point of Interest Teleporter, by Beneos",
        "version": "13.0.0",
        "active": false
      },
      {
        "id": "potion-crafting-and-gathering",
        "title": "Potion Crafting & Gathering",
        "version": "2.0.0",
        "active": false
      },
      {
        "id": "psfx",
        "title": "PSFX - Peri's Sound Effects",
        "version": "0.4.0",
        "active": false
      },
      {
        "id": "puzzle-locks",
        "title": "Puzzle Locks",
        "version": "3.0.3",
        "active": false
      },
      {
        "id": "quick-insert",
        "title": "Quick Insert - Search Widget",
        "version": "3.6.0",
        "active": false
      },
      {
        "id": "region-attacher",
        "title": "Region Attacher",
        "version": "1.9.1",
        "active": false
      },
      {
        "id": "remote-highlight-ui",
        "title": "Remote Highlight UI",
        "version": "2.3.0",
        "active": false
      },
      {
        "id": "rest-recovery",
        "title": "Rest Recovery for 5E",
        "version": "5.0.0",
        "active": false
      },
      {
        "id": "ripper-premium-dice",
        "title": "Wave Dice - TheRipper93 Premium Dice (Supporter)",
        "version": "2.0.0",
        "active": false
      },
      {
        "id": "ripper-premium-dice-ea",
        "title": "Animated Dice Pack - TheRipper93 Premium Dice (Early Access)",
        "version": "2.0.0",
        "active": false
      },
      {
        "id": "routinglib",
        "title": "routinglib",
        "version": "1.1.0",
        "active": false
      },
      {
        "id": "scene-packer",
        "title": "Library: Scene Packer",
        "version": "2.8.9",
        "active": false
      },
      {
        "id": "sequencer",
        "title": "Sequencer",
        "version": "3.6.9",
        "active": false
      },
      {
        "id": "simbuls-athenaeum",
        "title": "Simbul's Athenaeum",
        "version": "1.1.0",
        "active": false
      },
      {
        "id": "simbuls-creature-aide",
        "title": "Simbul's Creature Aide",
        "version": "1.2.1",
        "active": false
      },
      {
        "id": "simple-quest",
        "title": "Simple Quest",
        "version": "3.0.11",
        "active": false
      },
      {
        "id": "smarttarget",
        "title": "Smart Target",
        "version": "3.0.1",
        "active": false
      },
      {
        "id": "socketlib",
        "title": "socketlib",
        "version": "v1.1.3",
        "active": false
      },
      {
        "id": "soncraft-audio",
        "title": "Soncraft Audio Module",
        "version": "0.2.1",
        "active": false
      },
      {
        "id": "spotlight-omnisearch",
        "title": "Spotlight Omnisearch",
        "version": "3.0.1",
        "active": false
      },
      {
        "id": "stealthy",
        "title": "Stealthy",
        "version": "13.0.1",
        "active": false
      },
      {
        "id": "table-safety",
        "title": "Table Safety",
        "version": "1.3.2",
        "active": false
      },
      {
        "id": "tagger",
        "title": "Tagger",
        "version": "1.5.3",
        "active": false
      },
      {
        "id": "theripper-premium-hub",
        "title": "TheRipper93's Module Hub",
        "version": "5.0.6",
        "active": false
      },
      {
        "id": "tidbits",
        "title": "Tidbits",
        "version": "3.0.5",
        "active": false
      },
      {
        "id": "tile-scroll",
        "title": "Tile Scroll",
        "version": "4.0.0",
        "active": false
      },
      {
        "id": "times-up",
        "title": "Times Up",
        "version": "13.1.0",
        "active": false
      },
      {
        "id": "token-frames",
        "title": "Token Frames",
        "version": "1.7.0",
        "active": false
      },
      {
        "id": "token-hud-wildcard",
        "title": "Token HUD Wildcard",
        "version": "3.3.0",
        "active": false
      },
      {
        "id": "torch",
        "title": "Torch",
        "version": "3.1.0",
        "active": false
      },
      {
        "id": "trinium-chat-buttons",
        "title": "Trinium Chat Buttons",
        "version": "v2.0.2",
        "active": false
      },
      {
        "id": "trs-foundryvtt-dice-set",
        "title": "Foundry VTT Anniversary Dice Set by The Rollsmith",
        "version": "1.0.1",
        "active": false
      },
      {
        "id": "vision-5e",
        "title": "Vision 5e",
        "version": "3.0.6",
        "active": false
      },
      {
        "id": "visual-active-effects",
        "title": "Visual Active Effects",
        "version": "13.0.3",
        "active": false
      },
      {
        "id": "vtta-tokenizer",
        "title": "Tokenizer",
        "version": "4.5.6",
        "active": false
      },
      {
        "id": "wall-height",
        "title": "Wall Height",
        "version": "7.0.8",
        "active": false
      },
      {
        "id": "world-setting-sync",
        "title": "World Setting Sync",
        "version": "2.0.0",
        "active": false
      },
      {
        "id": "zoom-pan-options",
        "title": "Zoom/Pan Options",
        "version": "2.0.1",
        "active": false
      }
    ],
    "users": [
      {
        "id": "5ypAoBvOiyjDKiaZ",
        "name": "Gamemaster",
        "role": 4,
        "isGM": true,
        "active": false,
        "color": "#28cca2",
        "avatar": "icons/svg/mystery-man.svg"
      },
      {
        "id": "r6bXhB7k9cXa3cif",
        "name": "tester",
        "role": 4,
        "isGM": true,
        "active": true,
        "color": "#cc2829",
        "avatar": "icons/svg/mystery-man.svg"
      },
      {
        "id": "JLfKTeTgCDpAdDfw",
        "name": "some-cool-guy",
        "role": 1,
        "isGM": false,
        "active": false,
        "color": "#6328cc",
        "avatar": "icons/svg/mystery-man.svg"
      }
    ],
    "activeScene": {
      "id": "r36nfimJGHYGUGQX",
      "name": "test-scene-updated"
    }
  }
}
```


