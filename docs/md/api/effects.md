---
tag: effects
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Effects

## GET /effects

Get all active effects on an actor or token

Returns the collection of ActiveEffect documents currently applied to the specified actor or token.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to query |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of active effects

### Try It Out

<ApiTester
  method="GET"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/effects';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  uuid: 'Actor.q9uWyfdPwTlzbpxb'
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
curl -X GET 'http://localhost:3010/effects?clientId=fvtt_099ad17ea199e7e3&uuid=Actor.q9uWyfdPwTlzbpxb' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/effects'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'uuid': 'Actor.q9uWyfdPwTlzbpxb'
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
  const path = '/effects';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    uuid: 'Actor.q9uWyfdPwTlzbpxb'
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
  🔤/effects🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤uuid=Actor.q9uWyfdPwTlzbpxb🔤 ➡️ uuid
  🔤?🧲clientId🧲&🧲uuid🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /effects🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "get-effects-result",
  "requestId": "get-effects_1776657998905",
  "data": {
    "uuid": "Actor.q9uWyfdPwTlzbpxb",
    "effects": []
  }
}
```


---

## GET /effects/list

List all available status effects

Returns all status effects defined by the game system's configuration. Useful for discovering valid statusId values for the add/remove effect endpoints.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of available status effects with id, name, and icon

### Try It Out

<ApiTester
  method="GET"
  path="/effects/list"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/effects/list';
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
curl -X GET 'http://localhost:3010/effects/list?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/effects/list'
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
  const path = '/effects/list';
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
  🔤/effects/list🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /effects/list🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "get-status-effects-result",
  "requestId": "get-status-effects_1776657998898",
  "data": {
    "effects": [
      {
        "id": "dead",
        "name": "Dead",
        "icon": "systems/dnd5e/icons/svg/statuses/dead.svg"
      },
      {
        "id": "coverHalf",
        "name": "Half Cover",
        "icon": "systems/dnd5e/icons/svg/statuses/cover-half.svg"
      },
      {
        "id": "coverThreeQuarters",
        "name": "Three-Quarters Cover",
        "icon": "systems/dnd5e/icons/svg/statuses/cover-three-quarters.svg"
      },
      {
        "id": "coverTotal",
        "name": "Total Cover",
        "icon": "systems/dnd5e/icons/svg/statuses/cover-total.svg"
      },
      {
        "id": "bleeding",
        "name": "Bleeding",
        "icon": "systems/dnd5e/icons/svg/statuses/bleeding.svg"
      },
      {
        "id": "blinded",
        "name": "Blinded",
        "icon": "systems/dnd5e/icons/svg/statuses/blinded.svg"
      },
      {
        "id": "burning",
        "name": "Burning",
        "icon": "systems/dnd5e/icons/svg/statuses/burning.svg"
      },
      {
        "id": "burrowing",
        "name": "Burrowing",
        "icon": "systems/dnd5e/icons/svg/statuses/burrowing.svg"
      },
      {
        "id": "charmed",
        "name": "Charmed",
        "icon": "systems/dnd5e/icons/svg/statuses/charmed.svg"
      },
      {
        "id": "concentrating",
        "name": "Concentrating",
        "icon": "systems/dnd5e/icons/svg/statuses/concentrating.svg"
      },
      {
        "id": "cursed",
        "name": "Cursed",
        "icon": "systems/dnd5e/icons/svg/statuses/cursed.svg"
      },
      {
        "id": "deafened",
        "name": "Deafened",
        "icon": "systems/dnd5e/icons/svg/statuses/deafened.svg"
      },
      {
        "id": "dehydration",
        "name": "Dehydration",
        "icon": "systems/dnd5e/icons/svg/statuses/dehydration.svg"
      },
      {
        "id": "diseased",
        "name": "Diseased",
        "icon": "systems/dnd5e/icons/svg/statuses/diseased.svg"
      },
      {
        "id": "dodging",
        "name": "Dodging",
        "icon": "systems/dnd5e/icons/svg/statuses/dodging.svg"
      },
      {
        "id": "encumbered",
        "name": "Encumbered",
        "icon": "systems/dnd5e/icons/svg/statuses/encumbered.svg"
      },
      {
        "id": "ethereal",
        "name": "Ethereal",
        "icon": "systems/dnd5e/icons/svg/statuses/ethereal.svg"
      },
      {
        "id": "exceedingCarryingCapacity",
        "name": "Exceeding Carrying Capacity",
        "icon": "systems/dnd5e/icons/svg/statuses/exceeding-carrying-capacity.svg"
      },
      {
        "id": "exhaustion",
        "name": "Exhaustion",
        "icon": "systems/dnd5e/icons/svg/statuses/exhaustion.svg"
      },
      {
        "id": "falling",
        "name": "Falling",
        "icon": "systems/dnd5e/icons/svg/statuses/falling.svg"
      },
      {
        "id": "flying",
        "name": "Flying",
        "icon": "systems/dnd5e/icons/svg/statuses/flying.svg"
      },
      {
        "id": "frightened",
        "name": "Frightened",
        "icon": "systems/dnd5e/icons/svg/statuses/frightened.svg"
      },
      {
        "id": "grappled",
        "name": "Grappled",
        "icon": "systems/dnd5e/icons/svg/statuses/grappled.svg"
      },
      {
        "id": "heavilyEncumbered",
        "name": "Heavily Encumbered",
        "icon": "systems/dnd5e/icons/svg/statuses/heavily-encumbered.svg"
      },
      {
        "id": "hiding",
        "name": "Hiding",
        "icon": "systems/dnd5e/icons/svg/statuses/hiding.svg"
      },
      {
        "id": "hovering",
        "name": "Hovering",
        "icon": "systems/dnd5e/icons/svg/statuses/hovering.svg"
      },
      {
        "id": "incapacitated",
        "name": "Incapacitated",
        "icon": "systems/dnd5e/icons/svg/statuses/incapacitated.svg"
      },
      {
        "id": "invisible",
        "name": "Invisible",
        "icon": "systems/dnd5e/icons/svg/statuses/invisible.svg"
      },
      {
        "id": "malnutrition",
        "name": "Malnutrition",
        "icon": "systems/dnd5e/icons/svg/statuses/malnutrition.svg"
      },
      {
        "id": "marked",
        "name": "Marked",
        "icon": "systems/dnd5e/icons/svg/statuses/marked.svg"
      },
      {
        "id": "paralyzed",
        "name": "Paralyzed",
        "icon": "systems/dnd5e/icons/svg/statuses/paralyzed.svg"
      },
      {
        "id": "petrified",
        "name": "Petrified",
        "icon": "systems/dnd5e/icons/svg/statuses/petrified.svg"
      },
      {
        "id": "poisoned",
        "name": "Poisoned",
        "icon": "systems/dnd5e/icons/svg/statuses/poisoned.svg"
      },
      {
        "id": "prone",
        "name": "Prone",
        "icon": "systems/dnd5e/icons/svg/statuses/prone.svg"
      },
      {
        "id": "restrained",
        "name": "Restrained",
        "icon": "systems/dnd5e/icons/svg/statuses/restrained.svg"
      },
      {
        "id": "silenced",
        "name": "Silenced",
        "icon": "systems/dnd5e/icons/svg/statuses/silenced.svg"
      },
      {
        "id": "sleeping",
        "name": "Sleeping",
        "icon": "systems/dnd5e/icons/svg/statuses/sleeping.svg"
      },
      {
        "id": "stable",
        "name": "Stable",
        "icon": "systems/dnd5e/icons/svg/statuses/stable.svg"
      },
      {
        "id": "stunned",
        "name": "Stunned",
        "icon": "systems/dnd5e/icons/svg/statuses/stunned.svg"
      },
      {
        "id": "suffocation",
        "name": "Suffocation",
        "icon": "systems/dnd5e/icons/svg/statuses/suffocation.svg"
      },
      {
        "id": "surprised",
        "name": "Surprised",
        "icon": "systems/dnd5e/icons/svg/statuses/surprised.svg"
      },
      {
        "id": "transformed",
        "name": "Transformed",
        "icon": "systems/dnd5e/icons/svg/statuses/transformed.svg"
      },
      {
        "id": "unconscious",
        "name": "Unconscious",
        "icon": "systems/dnd5e/icons/svg/statuses/unconscious.svg"
      }
    ]
  }
}
```


---

## POST /effects

Add an active effect to an actor or token

Adds a status condition (by statusId) or a custom ActiveEffect (via effectData) to the specified actor or token.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to add the effect to |
| clientId | string |  | query | Client ID for the Foundry world |
| statusId | string |  | body, query | Standard status condition ID (e.g., "poisoned", "blinded", "prone") |
| effectData | object |  | body, query | Custom ActiveEffect data object (name, icon, duration, changes) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the add operation

### Try It Out

<ApiTester
  method="POST"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"statusId","type":"string","required":false,"source":"body"},{"name":"effectData","type":"object","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/effects';
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
      "uuid": "Actor.q9uWyfdPwTlzbpxb",
      "effectData": {
        "name": "Test Effect",
        "icon": "icons/svg/aura.svg",
        "changes": []
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/effects?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"uuid":"Actor.q9uWyfdPwTlzbpxb","effectData":{"name":"Test Effect","icon":"icons/svg/aura.svg","changes":[]}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/effects'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here',
        'Content-Type': 'application/json'
    },
    json={
      "uuid": "Actor.q9uWyfdPwTlzbpxb",
      "effectData": {
        "name": "Test Effect",
        "icon": "icons/svg/aura.svg",
        "changes": []
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
  const path = '/effects';
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
        "uuid": "Actor.q9uWyfdPwTlzbpxb",
        "effectData": {
          "name": "Test Effect",
          "icon": "icons/svg/aura.svg",
          "changes": []
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
  🔤/effects🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"uuid":"Actor.q9uWyfdPwTlzbpxb","effectData":{"name":"Test Effect","icon":"icons/svg/aura.svg","changes":[]}}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /effects🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 110❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "add-effect-result",
  "requestId": "add-effect_1776657998907",
  "data": {
    "uuid": "Actor.q9uWyfdPwTlzbpxb",
    "effect": {
      "id": "QTBr8UFiVsPT7pdI",
      "uuid": "Actor.q9uWyfdPwTlzbpxb.ActiveEffect.QTBr8UFiVsPT7pdI",
      "name": "Test Effect",
      "icon": "icons/svg/aura.svg",
      "statuses": []
    }
  }
}
```


---

## DELETE /effects

Remove an active effect from an actor or token

Removes an effect by its document ID (effectId) or by status condition identifier (statusId).

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| uuid | string | ✓ | body, query | UUID of the actor or token to remove the effect from |
| clientId | string |  | query | Client ID for the Foundry world |
| effectId | string |  | body, query | The ActiveEffect document ID to remove |
| statusId | string |  | body, query | Standard status condition ID to remove (e.g., "poisoned") |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the remove operation

### Try It Out

<ApiTester
  method="DELETE"
  path="/effects"
  parameters={[{"name":"uuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"effectId","type":"string","required":false,"source":"body"},{"name":"statusId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/effects';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3'
};
const queryString = new URLSearchParams(params).toString();
const url = `${baseUrl}${path}?${queryString}`;

const response = await fetch(url, {
  method: 'DELETE',
  headers: {
    'x-api-key': 'your-api-key-here',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "uuid": "Actor.q9uWyfdPwTlzbpxb",
      "effectId": "QTBr8UFiVsPT7pdI"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X DELETE 'http://localhost:3010/effects?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"uuid":"Actor.q9uWyfdPwTlzbpxb","effectId":"QTBr8UFiVsPT7pdI"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/effects'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
}
url = f'{base_url}{path}'

response = requests.delete(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here',
        'Content-Type': 'application/json'
    },
    json={
      "uuid": "Actor.q9uWyfdPwTlzbpxb",
      "effectId": "QTBr8UFiVsPT7pdI"
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
  const path = '/effects';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3'
  };
  const queryString = new URLSearchParams(params).toString();
  const url = `${baseUrl}${path}?${queryString}`;

  const response = await axios({
    method: 'delete',
    headers: {
      'x-api-key': 'your-api-key-here',
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "uuid": "Actor.q9uWyfdPwTlzbpxb",
        "effectId": "QTBr8UFiVsPT7pdI"
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
  🔤/effects🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"uuid":"Actor.q9uWyfdPwTlzbpxb","effectId":"QTBr8UFiVsPT7pdI"}🔤 ➡️ body

  💭 Build HTTP request
  🔤DELETE /effects🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 63❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "remove-effect-result",
  "requestId": "remove-effect_1776657998989",
  "data": {
    "uuid": "Actor.q9uWyfdPwTlzbpxb",
    "removedEffectId": "QTBr8UFiVsPT7pdI"
  }
}
```


