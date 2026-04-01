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
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls) |
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
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
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
curl -X GET 'http://localhost:3010/canvas/tokens?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/canvas/tokens'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif'
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
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
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
  "requestId": "get-canvas-documents_1775068875744",
  "sceneId": "kxsUY8QzM6ClEZqt",
  "documentType": "tokens",
  "data": [
    {
      "actorId": "pxZTVHItjx6GgPgC",
      "x": 400,
      "y": 400,
      "shape": 4,
      "_id": "dkx4UQ1aawVpxXbM",
      "name": "",
      "displayName": 0,
      "actorLink": false,
      "delta": {
        "_id": "1PQGmRrCINutTzTZ",
        "system": {},
        "items": [],
        "effects": [],
        "flags": {}
      },
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
      "turnMarker": {
        "mode": 1,
        "animation": null,
        "src": null,
        "disposition": false
      },
      "movementAction": null,
      "_movementHistory": [],
      "_regions": [],
      "flags": {}
    }
  ]
}
```


---

## POST /canvas/:documentType

Create canvas embedded document(s)

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls) |
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
      "data": {
        "x": 400,
        "y": 400,
        "actorId": "pxZTVHItjx6GgPgC"
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/canvas/tokens?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"data":{"x":400,"y":400,"actorId":"pxZTVHItjx6GgPgC"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/canvas/tokens'
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
      "data": {
        "x": 400,
        "y": 400,
        "actorId": "pxZTVHItjx6GgPgC"
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
        "data": {
          "x": 400,
          "y": 400,
          "actorId": "pxZTVHItjx6GgPgC"
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"data":{"x":400,"y":400,"actorId":"pxZTVHItjx6GgPgC"}}🔤 ➡️ body

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
  "requestId": "create-canvas-document_1775068875614",
  "sceneId": "kxsUY8QzM6ClEZqt",
  "documentType": "tokens",
  "data": [
    {
      "actorId": "pxZTVHItjx6GgPgC",
      "x": 400,
      "y": 400,
      "shape": 4,
      "_id": "dkx4UQ1aawVpxXbM",
      "name": "",
      "displayName": 0,
      "actorLink": false,
      "delta": {
        "_id": "1PQGmRrCINutTzTZ",
        "system": {},
        "items": [],
        "effects": [],
        "flags": {}
      },
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
      "turnMarker": {
        "mode": 1,
        "animation": null,
        "src": null,
        "disposition": false
      },
      "movementAction": null,
      "_movementHistory": [],
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
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls) |
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
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
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
      "documentId": "dkx4UQ1aawVpxXbM",
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
curl -X PUT 'http://localhost:3010/canvas/tokens?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"documentId":"dkx4UQ1aawVpxXbM","data":{"x":450,"y":450}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/canvas/tokens'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif'
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
      "documentId": "dkx4UQ1aawVpxXbM",
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
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif'
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
        "documentId": "dkx4UQ1aawVpxXbM",
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"documentId":"dkx4UQ1aawVpxXbM","data":{"x":450,"y":450}}🔤 ➡️ body

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
  "requestId": "update-canvas-document_1775068875754",
  "sceneId": "kxsUY8QzM6ClEZqt",
  "documentType": "tokens",
  "data": [
    {
      "actorId": "pxZTVHItjx6GgPgC",
      "x": 450,
      "y": 450,
      "shape": 4,
      "_id": "dkx4UQ1aawVpxXbM",
      "name": "",
      "displayName": 0,
      "actorLink": false,
      "delta": {
        "_id": "1PQGmRrCINutTzTZ",
        "system": {},
        "items": [],
        "effects": [],
        "flags": {}
      },
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
      "turnMarker": {
        "mode": 1,
        "animation": null,
        "src": null,
        "disposition": false
      },
      "movementAction": null,
      "_movementHistory": [],
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
| documentType | string | ✓ | params | Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls) |
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
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
  documentId: 'dkx4UQ1aawVpxXbM'
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
curl -X DELETE 'http://localhost:3010/canvas/tokens?clientId=foundry-testing-r6bXhB7k9cXa3cif&documentId=dkx4UQ1aawVpxXbM' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/canvas/tokens'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif',
    'documentId': 'dkx4UQ1aawVpxXbM'
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
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
    documentId: 'dkx4UQ1aawVpxXbM'
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤documentId=dkx4UQ1aawVpxXbM🔤 ➡️ documentId
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
  "requestId": "delete-canvas-document_1775068875847",
  "sceneId": "kxsUY8QzM6ClEZqt",
  "documentType": "tokens",
  "success": true
}
```


