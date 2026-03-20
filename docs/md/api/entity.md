---
tag: entity
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# entity

## GET /get

Get entity details This endpoint retrieves the details of a specific entity.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to retrieve (optional if selected=true) |
| selected | boolean |  | query | Whether to get the selected entity |
| actor | boolean |  | query | Return the actor of specified entity |

### Returns

**object** - Entity details object containing requested information

### Try It Out

<ApiTester
  method="GET"
  path="/get"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"actor","type":"boolean","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/get';
const params = {
  clientId: 'your-client-id',
  uuid: 'Actor.E3Ve2pzkDJYYj9eZ'
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
curl -X GET 'http://localhost:3010/get?clientId=your-client-id&uuid=Actor.E3Ve2pzkDJYYj9eZ' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/get'
params = {
    'clientId': 'your-client-id',
    'uuid': 'Actor.E3Ve2pzkDJYYj9eZ'
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
  const path = '/get';
  const params = {
    clientId: 'your-client-id',
    uuid: 'Actor.E3Ve2pzkDJYYj9eZ'
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
  🔤/get🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤uuid=Actor.E3Ve2pzkDJYYj9eZ🔤 ➡️ uuid
  🔤?🧲clientId🧲&🧲uuid🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /get🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "entity_1773794684057",
  "clientId": "your-client-id",
  "type": "entity-result",
  "uuid": "Actor.E3Ve2pzkDJYYj9eZ",
  "data": {
    "name": "test-actor",
    "type": "base",
    "folder": null,
    "_id": "E3Ve2pzkDJYYj9eZ",
    "img": "icons/svg/mystery-man.svg",
    "system": {},
    "prototypeToken": {
      "name": "test-actor",
      "displayName": 0,
      "actorLink": false,
      "width": 1,
      "height": 1,
      "texture": {
        "src": "icons/svg/mystery-man.svg",
        "anchorX": 0.5,
        "anchorY": 0.5,
        "offsetX": 0,
        "offsetY": 0,
        "fit": "contain",
        "scaleX": 1,
        "scaleY": 1,
        "rotation": 0,
        "tint": "#ffffff",
        "alphaThreshold": 0.75
      },
      "lockRotation": false,
      "rotation": 0,
      "alpha": 1,
      "disposition": -1,
      "displayBars": 0,
      "bar1": {
        "attribute": "attributes.hp"
      },
      "bar2": {
        "attribute": null
      },
      "light": {
        "negative": false,
        "priority": 0,
        "alpha": 0.5,
        "angle": 360,
        "bright": 0,
        "color": null,
        "coloration": 1,
        "dim": 0,
        "attenuation": 0.5,
        "luminosity": 0.5,
        "saturation": 0,
        "contrast": 0,
        "shadows": 0,
        "animation": {
          "type": null,
          "speed": 5,
          "intensity": 5,
          "reverse": false
        },
        "darkness": {
          "min": 0,
          "max": 1
        }
      },
      "sight": {
        "enabled": false,
        "range": 0,
        "angle": 360,
        "visionMode": "basic",
        "color": null,
        "attenuation": 0.1,
        "brightness": 0,
        "saturation": 0,
        "contrast": 0
      },
      "detectionModes": [],
      "occludable": {
        "radius": 0
      },
      "ring": {
        "enabled": false,
        "colors": {
          "ring": null,
          "background": null
        },
        "effects": 1,
        "subject": {
          "scale": 1,
          "texture": null
        }
      },
      "turnMarker": {
        "mode": 1,
        "animation": null,
        "src": null,
        "disposition": false
      },
      "movementAction": null,
      "flags": {},
      "randomImg": false,
      "appendNumber": false,
      "prependAdjective": false
    },
    "items": [],
    "effects": [],
    "sort": 0,
    "ownership": {
      "default": 0,
      "r6bXhB7k9cXa3cif": 3
    },
    "flags": {},
    "_stats": {
      "compendiumSource": null,
      "duplicateSource": null,
      "exportSource": null,
      "coreVersion": "13.348",
      "systemId": "dnd5e",
      "systemVersion": "5.0.4",
      "createdTime": 1773794681164,
      "modifiedTime": 1773794681164,
      "lastModifiedBy": "r6bXhB7k9cXa3cif"
    }
  }
}
```


---

## POST /create

Create a new entity This endpoint creates a new entity in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| entityType | string | ✓ | body | Document type of entity to create (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist, ext.) |
| data | object | ✓ | body | Data for the new entity |
| folder | string |  | body | Optional folder UUID to place the new entity in |

### Returns

**object** - Result of the entity creation operation

### Try It Out

<ApiTester
  method="POST"
  path="/create"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"query"},{"name":"entityType","type":"string","required":true,"source":"body"},{"name":"data","type":"object","required":true,"source":"body"},{"name":"folder","type":"string","required":false,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/create';
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
      "entityType": "Actor",
      "data": {
        "name": "test-actor",
        "type": "base"
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/create?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"entityType":"Actor","data":{"name":"test-actor","type":"base"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/create'
params = {
    'clientId': 'your-client-id'
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
      "entityType": "Actor",
      "data": {
        "name": "test-actor",
        "type": "base"
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
  const path = '/create';
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
        "entityType": "Actor",
        "data": {
          "name": "test-actor",
          "type": "base"
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
  🔤/create🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"entityType":"Actor","data":{"name":"test-actor","type":"base"}}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /create🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 65❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "create_1773794680897",
  "clientId": "your-client-id",
  "type": "create-result",
  "uuid": "Actor.E3Ve2pzkDJYYj9eZ",
  "entity": {
    "name": "test-actor",
    "type": "base",
    "folder": null,
    "_id": "E3Ve2pzkDJYYj9eZ",
    "img": "icons/svg/mystery-man.svg",
    "system": {},
    "prototypeToken": {
      "name": "test-actor",
      "displayName": 0,
      "actorLink": false,
      "width": 1,
      "height": 1,
      "texture": {
        "src": "icons/svg/mystery-man.svg",
        "anchorX": 0.5,
        "anchorY": 0.5,
        "offsetX": 0,
        "offsetY": 0,
        "fit": "contain",
        "scaleX": 1,
        "scaleY": 1,
        "rotation": 0,
        "tint": "#ffffff",
        "alphaThreshold": 0.75
      },
      "lockRotation": false,
      "rotation": 0,
      "alpha": 1,
      "disposition": -1,
      "displayBars": 0,
      "bar1": {
        "attribute": "attributes.hp"
      },
      "bar2": {
        "attribute": null
      },
      "light": {
        "negative": false,
        "priority": 0,
        "alpha": 0.5,
        "angle": 360,
        "bright": 0,
        "color": null,
        "coloration": 1,
        "dim": 0,
        "attenuation": 0.5,
        "luminosity": 0.5,
        "saturation": 0,
        "contrast": 0,
        "shadows": 0,
        "animation": {
          "type": null,
          "speed": 5,
          "intensity": 5,
          "reverse": false
        },
        "darkness": {
          "min": 0,
          "max": 1
        }
      },
      "sight": {
        "enabled": false,
        "range": 0,
        "angle": 360,
        "visionMode": "basic",
        "color": null,
        "attenuation": 0.1,
        "brightness": 0,
        "saturation": 0,
        "contrast": 0
      },
      "detectionModes": [],
      "occludable": {
        "radius": 0
      },
      "ring": {
        "enabled": false,
        "colors": {
          "ring": null,
          "background": null
        },
        "effects": 1,
        "subject": {
          "scale": 1,
          "texture": null
        }
      },
      "turnMarker": {
        "mode": 1,
        "animation": null,
        "src": null,
        "disposition": false
      },
      "movementAction": null,
      "flags": {},
      "randomImg": false,
      "appendNumber": false,
      "prependAdjective": false
    },
    "items": [],
    "effects": [],
    "sort": 0,
    "ownership": {
      "default": 0,
      "r6bXhB7k9cXa3cif": 3
    },
    "flags": {},
    "_stats": {
      "compendiumSource": null,
      "duplicateSource": null,
      "exportSource": null,
      "coreVersion": "13.348",
      "systemId": "dnd5e",
      "systemVersion": "5.0.4",
      "createdTime": 1773794681164,
      "modifiedTime": 1773794681164,
      "lastModifiedBy": "r6bXhB7k9cXa3cif"
    }
  }
}
```


---

## PUT /update

Update an existing entity This endpoint updates an existing entity in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| data | object | ✓ | body | Data to update the entity with |
| uuid | string |  | query | UUID of the entity to update (optional if selected=true) |
| selected | boolean |  | query | Whether to update the selected entity |
| actor | boolean |  | query | Update the actor of selected entity when selected=true |

### Returns

**object** - Result of the entity update operation

### Try It Out

<ApiTester
  method="PUT"
  path="/update"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"query"},{"name":"data","type":"object","required":true,"source":"body"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"},{"name":"actor","type":"boolean","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/update';
const params = {
  clientId: 'your-client-id',
  uuid: 'Actor.E3Ve2pzkDJYYj9eZ'
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
      "data": {
        "name": "Updated Test Actor"
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X PUT 'http://localhost:3010/update?clientId=your-client-id&uuid=Actor.E3Ve2pzkDJYYj9eZ' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"data":{"name":"Updated Test Actor"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/update'
params = {
    'clientId': 'your-client-id',
    'uuid': 'Actor.E3Ve2pzkDJYYj9eZ'
}
url = f'{base_url}{path}'

response = requests.put(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "data": {
        "name": "Updated Test Actor"
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
  const path = '/update';
  const params = {
    clientId: 'your-client-id',
    uuid: 'Actor.E3Ve2pzkDJYYj9eZ'
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
        "data": {
          "name": "Updated Test Actor"
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
  🔤/update🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤uuid=Actor.E3Ve2pzkDJYYj9eZ🔤 ➡️ uuid
  🔤?🧲clientId🧲&🧲uuid🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"data":{"name":"Updated Test Actor"}}🔤 ➡️ body

  💭 Build HTTP request
  🔤PUT /update🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 38❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "update_1773794684324",
  "clientId": "your-client-id",
  "type": "update-result",
  "uuid": "Actor.E3Ve2pzkDJYYj9eZ",
  "entity": [
    {
      "name": "Updated Test Actor",
      "type": "base",
      "folder": null,
      "_id": "E3Ve2pzkDJYYj9eZ",
      "img": "icons/svg/mystery-man.svg",
      "system": {},
      "prototypeToken": {
        "name": "test-actor",
        "displayName": 0,
        "actorLink": false,
        "width": 1,
        "height": 1,
        "texture": {
          "src": "icons/svg/mystery-man.svg",
          "anchorX": 0.5,
          "anchorY": 0.5,
          "offsetX": 0,
          "offsetY": 0,
          "fit": "contain",
          "scaleX": 1,
          "scaleY": 1,
          "rotation": 0,
          "tint": "#ffffff",
          "alphaThreshold": 0.75
        },
        "lockRotation": false,
        "rotation": 0,
        "alpha": 1,
        "disposition": -1,
        "displayBars": 0,
        "bar1": {
          "attribute": "attributes.hp"
        },
        "bar2": {
          "attribute": null
        },
        "light": {
          "negative": false,
          "priority": 0,
          "alpha": 0.5,
          "angle": 360,
          "bright": 0,
          "color": null,
          "coloration": 1,
          "dim": 0,
          "attenuation": 0.5,
          "luminosity": 0.5,
          "saturation": 0,
          "contrast": 0,
          "shadows": 0,
          "animation": {
            "type": null,
            "speed": 5,
            "intensity": 5,
            "reverse": false
          },
          "darkness": {
            "min": 0,
            "max": 1
          }
        },
        "sight": {
          "enabled": false,
          "range": 0,
          "angle": 360,
          "visionMode": "basic",
          "color": null,
          "attenuation": 0.1,
          "brightness": 0,
          "saturation": 0,
          "contrast": 0
        },
        "detectionModes": [],
        "occludable": {
          "radius": 0
        },
        "ring": {
          "enabled": false,
          "colors": {
            "ring": null,
            "background": null
          },
          "effects": 1,
          "subject": {
            "scale": 1,
            "texture": null
          }
        },
        "turnMarker": {
          "mode": 1,
          "animation": null,
          "src": null,
          "disposition": false
        },
        "movementAction": null,
        "flags": {},
        "randomImg": false,
        "appendNumber": false,
        "prependAdjective": false
      },
      "items": [],
      "effects": [],
      "sort": 0,
      "ownership": {
        "default": 0,
        "r6bXhB7k9cXa3cif": 3
      },
      "flags": {},
      "_stats": {
        "compendiumSource": null,
        "duplicateSource": null,
        "exportSource": null,
        "coreVersion": "13.348",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1773794681164,
        "modifiedTime": 1773794684551,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    }
  ]
}
```


---

## DELETE /delete

Delete an entity This endpoint deletes an entity from the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to delete (optional if selected=true) |
| selected | boolean |  | query | Whether to delete the selected entity |

### Returns

**object** - Result of the entity deletion operation

### Try It Out

<ApiTester
  method="DELETE"
  path="/delete"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/delete';
const params = {
  clientId: 'your-client-id',
  uuid: 'Actor.qItWfYTYwt88LqyU'
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
curl -X DELETE 'http://localhost:3010/delete?clientId=your-client-id&uuid=Actor.qItWfYTYwt88LqyU' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/delete'
params = {
    'clientId': 'your-client-id',
    'uuid': 'Actor.qItWfYTYwt88LqyU'
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
  const path = '/delete';
  const params = {
    clientId: 'your-client-id',
    uuid: 'Actor.qItWfYTYwt88LqyU'
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
  🔤/delete🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤uuid=Actor.qItWfYTYwt88LqyU🔤 ➡️ uuid
  🔤?🧲clientId🧲&🧲uuid🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /delete🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "delete_1773794684830",
  "clientId": "your-client-id",
  "type": "delete-result",
  "uuid": "Actor.qItWfYTYwt88LqyU",
  "success": true
}
```


---

## POST /give

Give an item to an entity This endpoint gives an item to a specified entity. Optionally, removes the item from the giver.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| fromUuid | string |  | body | UUID of the entity giving the item |
| toUuid | string |  | body | UUID of the entity receiving the item |
| selected | boolean |  | body | Whether to give to the selected token's actor |
| itemUuid | string |  | body | UUID of the item to give (optional if itemName provided) |
| itemName | string |  | body | Name of the item to give (search with Quick Insert if UUID not provided) |
| quantity | number |  | body | Quantity of the item to give (negative values decrease quantity to 0) |

### Returns

**object** - Result of the item giving operation

### Try It Out

<ApiTester
  method="POST"
  path="/give"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"fromUuid","type":"string","required":false,"source":"body"},{"name":"toUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"body"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"},{"name":"quantity","type":"number","required":false,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/give';
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
      "toUuid": "Actor.E3Ve2pzkDJYYj9eZ",
      "itemUuid": "Item.wIAD6PBEmIKEaEMX",
      "quantity": 1
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/give?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"toUuid":"Actor.E3Ve2pzkDJYYj9eZ","itemUuid":"Item.wIAD6PBEmIKEaEMX","quantity":1}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/give'
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
      "toUuid": "Actor.E3Ve2pzkDJYYj9eZ",
      "itemUuid": "Item.wIAD6PBEmIKEaEMX",
      "quantity": 1
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
  const path = '/give';
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
        "toUuid": "Actor.E3Ve2pzkDJYYj9eZ",
        "itemUuid": "Item.wIAD6PBEmIKEaEMX",
        "quantity": 1
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
  🔤/give🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"toUuid":"Actor.E3Ve2pzkDJYYj9eZ","itemUuid":"Item.wIAD6PBEmIKEaEMX","quantity":1}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /give🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 83❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "give_1773794685347",
  "clientId": "your-client-id",
  "type": "give-result",
  "toUuid": "Actor.E3Ve2pzkDJYYj9eZ",
  "quantity": 1,
  "itemUuid": "Item.wIAD6PBEmIKEaEMX",
  "newItemId": "XEoN8awYLpLzloh8",
  "success": true
}
```


---

## POST /remove

Remove an item from an entity This endpoint removes an item from a specified entity.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string |  | body | UUID of the actor to remove the item from (optional if selected=true) |
| selected | boolean |  | body | Whether to remove from the selected token's actor |
| itemUuid | string |  | body | UUID of the item to remove |
| itemName | string |  | body | Name of the item to remove (search with Quick Insert if UUID not provided) |
| quantity | number |  | body | Quantity of the item to remove |

### Returns

**object** - Result of the item removal operation

### Try It Out

<ApiTester
  method="POST"
  path="/remove"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"body"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"body"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"},{"name":"quantity","type":"number","required":false,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/remove';
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
      "actorUuid": "Actor.E3Ve2pzkDJYYj9eZ",
      "itemUuid": "Actor.E3Ve2pzkDJYYj9eZ.Item.XEoN8awYLpLzloh8",
      "quantity": 1
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/remove?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.E3Ve2pzkDJYYj9eZ","itemUuid":"Actor.E3Ve2pzkDJYYj9eZ.Item.XEoN8awYLpLzloh8","quantity":1}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/remove'
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
      "actorUuid": "Actor.E3Ve2pzkDJYYj9eZ",
      "itemUuid": "Actor.E3Ve2pzkDJYYj9eZ.Item.XEoN8awYLpLzloh8",
      "quantity": 1
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
  const path = '/remove';
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
        "actorUuid": "Actor.E3Ve2pzkDJYYj9eZ",
        "itemUuid": "Actor.E3Ve2pzkDJYYj9eZ.Item.XEoN8awYLpLzloh8",
        "quantity": 1
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
  🔤/remove🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.E3Ve2pzkDJYYj9eZ","itemUuid":"Actor.E3Ve2pzkDJYYj9eZ.Item.XEoN8awYLpLzloh8","quantity":1}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /remove🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 109❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "remove_1773794685876",
  "clientId": "your-client-id",
  "type": "remove-result",
  "actorUuid": "Actor.E3Ve2pzkDJYYj9eZ",
  "itemUuid": "Actor.E3Ve2pzkDJYYj9eZ.Item.XEoN8awYLpLzloh8",
  "quantity": 0,
  "success": true
}
```


---

## POST /decrease

Decrease an attribute This endpoint decreases an attribute of a specified entity.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| attribute | string | ✓ | body | The attribute data path to decrease (e.g., "system.attributes.hp.value") |
| amount | number | ✓ | body | The amount to decrease the attribute by |
| uuid | string |  | query | UUID of the entity to decrease the attribute for (optional if selected=true) |
| selected | boolean |  | query | Whether to decrease the attribute for the selected entity |

### Returns

**object** - Result of the attribute decrease operation

### Try It Out

<ApiTester
  method="POST"
  path="/decrease"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"query"},{"name":"attribute","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/decrease';
const params = {
  clientId: 'your-client-id',
  uuid: 'Actor.E3Ve2pzkDJYYj9eZ'
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
      "attribute": "prototypeToken.height",
      "amount": 5
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/decrease?clientId=your-client-id&uuid=Actor.E3Ve2pzkDJYYj9eZ' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"attribute":"prototypeToken.height","amount":5}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/decrease'
params = {
    'clientId': 'your-client-id',
    'uuid': 'Actor.E3Ve2pzkDJYYj9eZ'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "attribute": "prototypeToken.height",
      "amount": 5
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
  const path = '/decrease';
  const params = {
    clientId: 'your-client-id',
    uuid: 'Actor.E3Ve2pzkDJYYj9eZ'
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
        "attribute": "prototypeToken.height",
        "amount": 5
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
  🔤/decrease🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤uuid=Actor.E3Ve2pzkDJYYj9eZ🔤 ➡️ uuid
  🔤?🧲clientId🧲&🧲uuid🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"attribute":"prototypeToken.height","amount":5}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /decrease🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 48❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "decrease_1773794686901",
  "clientId": "your-client-id",
  "type": "decrease-result",
  "results": [
    {
      "uuid": "Actor.E3Ve2pzkDJYYj9eZ",
      "attribute": "prototypeToken.height",
      "oldValue": 6,
      "newValue": 1
    }
  ],
  "success": true
}
```


---

## POST /increase

Increase an attribute This endpoint increases an attribute of a specified entity.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| attribute | string | ✓ | body | The attribute data path to increase (e.g., "system.attributes.hp.value") |
| amount | number | ✓ | body | The amount to increase the attribute by |
| uuid | string |  | query | UUID of the entity to increase the attribute for (optional if selected=true) |
| selected | boolean |  | query | Whether to increase the attribute for the selected entity |

### Returns

**object** - Result of the attribute increase operation

### Try It Out

<ApiTester
  method="POST"
  path="/increase"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"query"},{"name":"attribute","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/increase';
const params = {
  clientId: 'your-client-id',
  uuid: 'Actor.E3Ve2pzkDJYYj9eZ'
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
      "attribute": "prototypeToken.height",
      "amount": 5
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/increase?clientId=your-client-id&uuid=Actor.E3Ve2pzkDJYYj9eZ' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"attribute":"prototypeToken.height","amount":5}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/increase'
params = {
    'clientId': 'your-client-id',
    'uuid': 'Actor.E3Ve2pzkDJYYj9eZ'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "attribute": "prototypeToken.height",
      "amount": 5
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
  const path = '/increase';
  const params = {
    clientId: 'your-client-id',
    uuid: 'Actor.E3Ve2pzkDJYYj9eZ'
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
        "attribute": "prototypeToken.height",
        "amount": 5
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
  🔤/increase🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤uuid=Actor.E3Ve2pzkDJYYj9eZ🔤 ➡️ uuid
  🔤?🧲clientId🧲&🧲uuid🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"attribute":"prototypeToken.height","amount":5}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /increase🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 48❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "increase_1773794686408",
  "clientId": "your-client-id",
  "type": "increase-result",
  "results": [
    {
      "uuid": "Actor.E3Ve2pzkDJYYj9eZ",
      "attribute": "prototypeToken.height",
      "oldValue": 1,
      "newValue": 6
    }
  ],
  "success": true
}
```


---

## POST /kill

Kill an entity Marks an entity as killed in the combat tracker, gives it the "dead" status, and sets its health to 0 in the Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| uuid | string |  | query | UUID of the entity to kill (optional if selected=true) |
| selected | boolean |  | query | Whether to kill the selected entity |

### Returns

**object** - Result of the entity kill operation

### Try It Out

<ApiTester
  method="POST"
  path="/kill"
  parameters={[{"name":"clientId","type":"string","required":true,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"query"},{"name":"selected","type":"boolean","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/kill';
const params = {
  clientId: 'your-client-id',
  uuid: 'Actor.KqBznQo2iKKd84sK'
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
curl -X POST 'http://localhost:3010/kill?clientId=your-client-id&uuid=Actor.KqBznQo2iKKd84sK' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/kill'
params = {
    'clientId': 'your-client-id',
    'uuid': 'Actor.KqBznQo2iKKd84sK'
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
  const path = '/kill';
  const params = {
    clientId: 'your-client-id',
    uuid: 'Actor.KqBznQo2iKKd84sK'
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
  🔤/kill🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤uuid=Actor.KqBznQo2iKKd84sK🔤 ➡️ uuid
  🔤?🧲clientId🧲&🧲uuid🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤POST /kill🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "kill_1773794687406",
  "clientId": "your-client-id",
  "type": "kill-result",
  "results": [
    {
      "uuid": "Actor.KqBznQo2iKKd84sK",
      "success": true,
      "message": "Actor marked as defeated, HP set to 0, and dead effect applied to 0 tokens"
    }
  ]
}
```


