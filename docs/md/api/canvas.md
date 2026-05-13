---
tag: canvas
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Canvas

## GET /canvas/:documentType

Get canvas embedded documents

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | query | Scene ID to query (defaults to the active scene) |
| documentId | string |  | query | Specific document ID to retrieve |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of embedded documents

### Try It Out

<ApiTester
  method="GET"
  path="/canvas/:documentType"
  parameters={[{"name":"documentType","type":"string","required":true,"source":"params"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"documentId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/canvas/tokens';
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
curl -X GET 'http://localhost:3010/canvas/tokens?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/canvas/tokens'
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
  const path = '/canvas/tokens';
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
  🔤/canvas/tokens🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /canvas/tokens🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "get-canvas-documents-result",
  "requestId": "get-canvas-documents_1777996591565",
  "sceneId": "iI8vL6F5ett88LXH",
  "documentType": "tokens",
  "data": [
    {
      "actorId": "sMD3o6zej6ckQkpo",
      "x": 400,
      "y": 400,
      "_id": "fPVbRrzEleY88MDH",
      "name": "",
      "displayName": 0,
      "actorLink": false,
      "delta": {
        "_id": "Uta7DgvnRjkWubic",
        "system": {},
        "items": [],
        "effects": [],
        "flags": {}
      },
      "appendNumber": false,
      "prependAdjective": false,
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
      "hexagonalShape": 0,
      "elevation": 0,
      "sort": 0,
      "locked": false,
      "lockRotation": false,
      "rotation": 0,
      "alpha": 1,
      "hidden": false,
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
      "_regions": [],
      "flags": {}
    }
  ]
}
```


---

## GET /measure-distance

Measure the distance between two points or tokens

Calculates the distance between two positions on the canvas, respecting the grid type and measurement rules. Points can be specified as coordinates or by referencing tokens by UUID or name.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| originX | number |  | body, query | Origin x coordinate (optional if originUuid/originName provided) |
| originY | number |  | body, query | Origin y coordinate |
| targetX | number |  | body, query | Target x coordinate (optional if targetUuid/targetName provided) |
| targetY | number |  | body, query | Target y coordinate |
| originUuid | string |  | body, query | UUID of the origin token |
| originName | string |  | body, query | Name of the origin token |
| targetUuid | string |  | body, query | UUID of the target token |
| targetName | string |  | body, query | Name of the target token |
| sceneId | string |  | body, query | Scene ID (defaults to active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Distance measurement including units and grid spaces

### Try It Out

<ApiTester
  method="GET"
  path="/measure-distance"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"originX","type":"number","required":false,"source":"body"},{"name":"originY","type":"number","required":false,"source":"body"},{"name":"targetX","type":"number","required":false,"source":"body"},{"name":"targetY","type":"number","required":false,"source":"body"},{"name":"originUuid","type":"string","required":false,"source":"body"},{"name":"originName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/measure-distance';
const params = {
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
  originX: '0',
  originY: '0',
  targetX: '500',
  targetY: '500'
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
curl -X GET 'http://localhost:3010/measure-distance?clientId=foundry-testing-r6bXhB7k9cXa3cif&originX=0&originY=0&targetX=500&targetY=500' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/measure-distance'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif',
    'originX': '0',
    'originY': '0',
    'targetX': '500',
    'targetY': '500'
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
  const path = '/measure-distance';
  const params = {
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
    originX: '0',
    originY: '0',
    targetX: '500',
    targetY: '500'
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
  🔤/measure-distance🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤originX=0🔤 ➡️ originX
  🔤originY=0🔤 ➡️ originY
  🔤targetX=500🔤 ➡️ targetX
  🔤targetY=500🔤 ➡️ targetY
  🔤?🧲clientId🧲&🧲originX🧲&🧲originY🧲&🧲targetX🧲&🧲targetY🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /measure-distance🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "measure-distance-result",
  "requestId": "measure-distance_1775271900614",
  "data": {
    "distance": 25,
    "units": "ft",
    "origin": {
      "x": 0,
      "y": 0
    },
    "target": {
      "x": 500,
      "y": 500
    },
    "sceneId": "SkiuhdnphDtg3E9Q"
  }
}
```


---

## POST /canvas/:documentType

Create canvas embedded document(s)

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| data | object | ✓ | body | Document data object or array of objects to create |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | Scene ID to create in (defaults to the active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Created document(s)

### Try It Out

<ApiTester
  method="POST"
  path="/canvas/:documentType"
  parameters={[{"name":"documentType","type":"string","required":true,"source":"params"},{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/canvas/tokens';
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
      "data": {
        "x": 400,
        "y": 400,
        "actorId": "sMD3o6zej6ckQkpo"
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/canvas/tokens?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"data":{"x":400,"y":400,"actorId":"sMD3o6zej6ckQkpo"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/canvas/tokens'
params = {
    'clientId': 'fvtt_71dbc81bd608978a'
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
      "data": {
        "x": 400,
        "y": 400,
        "actorId": "sMD3o6zej6ckQkpo"
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
  const path = '/canvas/tokens';
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
        "data": {
          "x": 400,
          "y": 400,
          "actorId": "sMD3o6zej6ckQkpo"
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
  🔤/canvas/tokens🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"data":{"x":400,"y":400,"actorId":"sMD3o6zej6ckQkpo"}}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /canvas/tokens🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 55❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "create-canvas-document-result",
  "requestId": "create-canvas-document_1777996591475",
  "sceneId": "iI8vL6F5ett88LXH",
  "documentType": "tokens",
  "data": [
    {
      "actorId": "sMD3o6zej6ckQkpo",
      "x": 400,
      "y": 400,
      "_id": "fPVbRrzEleY88MDH",
      "name": "",
      "displayName": 0,
      "actorLink": false,
      "delta": {
        "_id": "Uta7DgvnRjkWubic",
        "system": {},
        "items": [],
        "effects": [],
        "flags": {}
      },
      "appendNumber": false,
      "prependAdjective": false,
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
      "hexagonalShape": 0,
      "elevation": 0,
      "sort": 0,
      "locked": false,
      "lockRotation": false,
      "rotation": 0,
      "alpha": 1,
      "hidden": false,
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
      "_regions": [],
      "flags": {}
    }
  ]
}
```


---

## PUT /canvas/:documentType

Update a canvas embedded document

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| documentId | string | ✓ | body, query | ID of the document to update |
| data | object | ✓ | body | Object containing the fields to update |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | Scene ID containing the document (defaults to the active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Updated document

### Try It Out

<ApiTester
  method="PUT"
  path="/canvas/:documentType"
  parameters={[{"name":"documentType","type":"string","required":true,"source":"params"},{"name":"documentId","type":"string","required":true,"source":"body"},{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/canvas/tokens';
const params = {
  clientId: 'fvtt_71dbc81bd608978a'
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
      "documentId": "fPVbRrzEleY88MDH",
      "data": {
        "x": 450,
        "y": 450
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X PUT 'http://localhost:3010/canvas/tokens?clientId=fvtt_71dbc81bd608978a' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"documentId":"fPVbRrzEleY88MDH","data":{"x":450,"y":450}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/canvas/tokens'
params = {
    'clientId': 'fvtt_71dbc81bd608978a'
}
url = f'{base_url}{path}'

response = requests.put(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here',
        'Content-Type': 'application/json'
    },
    json={
      "documentId": "fPVbRrzEleY88MDH",
      "data": {
        "x": 450,
        "y": 450
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
  const path = '/canvas/tokens';
  const params = {
    clientId: 'fvtt_71dbc81bd608978a'
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
        "documentId": "fPVbRrzEleY88MDH",
        "data": {
          "x": 450,
          "y": 450
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
  🔤/canvas/tokens🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"documentId":"fPVbRrzEleY88MDH","data":{"x":450,"y":450}}🔤 ➡️ body

  💭 Build HTTP request
  🔤PUT /canvas/tokens🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 58❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "update-canvas-document-result",
  "requestId": "update-canvas-document_1777996591571",
  "sceneId": "iI8vL6F5ett88LXH",
  "documentType": "tokens",
  "data": [
    {
      "actorId": "sMD3o6zej6ckQkpo",
      "x": 450,
      "y": 450,
      "_id": "fPVbRrzEleY88MDH",
      "name": "",
      "displayName": 0,
      "actorLink": false,
      "delta": {
        "_id": "Uta7DgvnRjkWubic",
        "system": {},
        "items": [],
        "effects": [],
        "flags": {}
      },
      "appendNumber": false,
      "prependAdjective": false,
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
      "hexagonalShape": 0,
      "elevation": 0,
      "sort": 0,
      "locked": false,
      "lockRotation": false,
      "rotation": 0,
      "alpha": 1,
      "hidden": false,
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
      "_regions": [],
      "flags": {}
    }
  ]
}
```


---

## DELETE /canvas/:documentType

Delete a canvas embedded document

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls, regions) |
| documentId | string | ✓ | query | ID of the document to delete |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | query | Scene ID containing the document (defaults to the active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Deletion result

### Try It Out

<ApiTester
  method="DELETE"
  path="/canvas/:documentType"
  parameters={[{"name":"documentType","type":"string","required":true,"source":"params"},{"name":"documentId","type":"string","required":true,"source":"query"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/canvas/tokens';
const params = {
  clientId: 'fvtt_71dbc81bd608978a',
  documentId: 'fPVbRrzEleY88MDH'
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
curl -X DELETE 'http://localhost:3010/canvas/tokens?clientId=fvtt_71dbc81bd608978a&documentId=fPVbRrzEleY88MDH' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/canvas/tokens'
params = {
    'clientId': 'fvtt_71dbc81bd608978a',
    'documentId': 'fPVbRrzEleY88MDH'
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
  const path = '/canvas/tokens';
  const params = {
    clientId: 'fvtt_71dbc81bd608978a',
    documentId: 'fPVbRrzEleY88MDH'
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
  🔤/canvas/tokens🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤documentId=fPVbRrzEleY88MDH🔤 ➡️ documentId
  🔤?🧲clientId🧲&🧲documentId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /canvas/tokens🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "delete-canvas-document-result",
  "requestId": "delete-canvas-document_1777996591587",
  "sceneId": "iI8vL6F5ett88LXH",
  "documentType": "tokens",
  "success": true
}
```


---

## POST /move-token

Move a token to specific coordinates

Moves a token on the canvas to the specified x,y position, optionally animating through waypoints. Token can be identified by UUID or name.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| x | number | ✓ | body, query | Target x coordinate |
| y | number | ✓ | body, query | Target y coordinate |
| clientId | string |  | query | Client ID for the Foundry world |
| uuid | string |  | body, query | UUID of the token to move (optional if name provided) |
| name | string |  | body, query | Name of the token to move (optional if uuid provided) |
| waypoints | array |  | body, query | Array of waypoint objects with x and y coordinates to animate through before reaching final position |
| animate | boolean |  | body, query | Whether to animate the movement (default: true) |
| sceneId | string |  | body, query | Scene ID (defaults to active scene) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the token movement including new position

### Try It Out

<ApiTester
  method="POST"
  path="/move-token"
  parameters={[{"name":"x","type":"number","required":true,"source":"body"},{"name":"y","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"uuid","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"waypoints","type":"array","required":false,"source":"body"},{"name":"animate","type":"boolean","required":false,"source":"body"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/move-token';
const params = {
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
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
      "name": "Test Token",
      "x": 500,
      "y": 500,
      "animate": false
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/move-token?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Token","x":500,"y":500,"animate":false}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/move-token'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif'
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
      "name": "Test Token",
      "x": 500,
      "y": 500,
      "animate": False
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
  const path = '/move-token';
  const params = {
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
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
        "name": "Test Token",
        "x": 500,
        "y": 500,
        "animate": false
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
  🔤/move-token🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"name":"Test Token","x":500,"y":500,"animate":false}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /move-token🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 53❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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

**Status:** 400

```json
{
  "type": "move-token-result",
  "requestId": "move-token_1775271900392",
  "error": "Token not found: Test Token"
}
```


