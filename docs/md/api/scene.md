---
tag: scene
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Scene

## GET /scene

Get scene(s)

Retrieves one or more scenes by ID, name, active status, viewed status, or all.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | query | ID of a specific scene to retrieve |
| name | string |  | query | Name of the scene to retrieve |
| active | boolean |  | query | Set to true to get the currently active scene |
| viewed | boolean |  | query | Set to true to get the currently viewed scene |
| all | boolean |  | query | Set to true to get all scenes |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Scene data

### Try It Out

<ApiTester
  method="GET"
  path="/scene"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"active","type":"boolean","required":false,"source":"query"},{"name":"viewed","type":"boolean","required":false,"source":"query"},{"name":"all","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/scene';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  all: 'true'
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
curl -X GET 'http://localhost:3010/scene?clientId=fvtt_099ad17ea199e7e3&all=true' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/scene'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'all': 'true'
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
  const path = '/scene';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    all: 'true'
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
  🔤/scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤all=true🔤 ➡️ all
  🔤?🧲clientId🧲&🧲all🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "get-scene-result",
  "requestId": "get-scene_1776657957940",
  "data": [
    {
      "name": "Scene",
      "_id": "i01IgzYFzGddbTtP",
      "active": false,
      "navigation": true,
      "navOrder": 0,
      "background": {
        "src": null,
        "anchorX": 0,
        "anchorY": 0,
        "offsetX": 0,
        "offsetY": 0,
        "fit": "fill",
        "scaleX": 1,
        "scaleY": 1,
        "rotation": 0,
        "tint": "#ffffff",
        "alphaThreshold": 0
      },
      "foreground": null,
      "foregroundElevation": null,
      "thumb": null,
      "width": 4000,
      "height": 3000,
      "padding": 0.25,
      "initial": {
        "x": null,
        "y": null,
        "scale": null
      },
      "backgroundColor": "#999999",
      "grid": {
        "type": 1,
        "size": 100,
        "style": "solidLines",
        "thickness": 1,
        "color": "#000000",
        "alpha": 0.2,
        "distance": 5,
        "units": "ft"
      },
      "tokenVision": true,
      "fog": {
        "exploration": true,
        "overlay": null,
        "colors": {
          "explored": null,
          "unexplored": null
        }
      },
      "environment": {
        "darknessLevel": 0,
        "darknessLock": false,
        "globalLight": {
          "enabled": false,
          "alpha": 0.5,
          "bright": false,
          "color": null,
          "coloration": 1,
          "luminosity": 0,
          "saturation": 0,
          "contrast": 0,
          "shadows": 0,
          "darkness": {
            "min": 0,
            "max": 1
          }
        },
        "cycle": true,
        "base": {
          "hue": 0,
          "intensity": 0,
          "luminosity": 0,
          "saturation": 0,
          "shadows": 0
        },
        "dark": {
          "hue": 0.7138888888888889,
          "intensity": 0,
          "luminosity": -0.25,
          "saturation": 0,
          "shadows": 0
        }
      },
      "drawings": [],
      "tokens": [],
      "lights": [],
      "notes": [],
      "sounds": [],
      "regions": [],
      "templates": [],
      "tiles": [],
      "walls": [],
      "playlist": null,
      "playlistSound": null,
      "journal": null,
      "journalEntryPage": null,
      "weather": "",
      "folder": null,
      "sort": 0,
      "ownership": {
        "default": 0,
        "5ypAoBvOiyjDKiaZ": 3
      },
      "flags": {},
      "_stats": {
        "compendiumSource": null,
        "duplicateSource": null,
        "exportSource": null,
        "coreVersion": "13.348",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1776529787190,
        "modifiedTime": 1776596230732,
        "lastModifiedBy": "5ypAoBvOiyjDKiaZ"
      }
    },
    {
      "name": "test",
      "_id": "fWuYVDZAnXQiKXM5",
      "active": true,
      "navigation": true,
      "navOrder": 0,
      "background": {
        "src": "flooded-cave-test.webp",
        "anchorX": 0,
        "anchorY": 0,
        "offsetX": 0,
        "offsetY": 0,
        "fit": "fill",
        "scaleX": 1,
        "scaleY": 1,
        "rotation": 0,
        "tint": "#ffffff",
        "alphaThreshold": 0
      },
      "foreground": null,
      "foregroundElevation": null,
      "thumb": "worlds/testing/assets/scenes/fWuYVDZAnXQiKXM5-thumb.webp",
      "width": 2221,
      "height": 2962,
      "padding": 0.25,
      "initial": {
        "x": null,
        "y": null,
        "scale": null
      },
      "backgroundColor": "#999999",
      "grid": {
        "type": 1,
        "size": 74,
        "style": "solidLines",
        "thickness": 1,
        "color": "#000000",
        "alpha": 0.2,
        "distance": 5,
        "units": "ft"
      },
      "tokenVision": true,
      "fog": {
        "exploration": true,
        "overlay": null,
        "colors": {
          "explored": null,
          "unexplored": null
        }
      },
      "environment": {
        "darknessLevel": 0,
        "darknessLock": false,
        "globalLight": {
          "enabled": false,
          "alpha": 0.5,
          "bright": false,
          "color": null,
          "coloration": 1,
          "luminosity": 0,
          "saturation": 0,
          "contrast": 0,
          "shadows": 0,
          "darkness": {
            "min": 0,
            "max": 1
          }
        },
        "cycle": true,
        "base": {
          "hue": 0,
          "intensity": 0,
          "luminosity": 0,
          "saturation": 0,
          "shadows": 0
        },
        "dark": {
          "hue": 0.7138888888888889,
          "intensity": 0,
          "luminosity": -0.25,
          "saturation": 0,
          "shadows": 0
        }
      },
      "drawings": [],
      "tokens": [],
      "lights": [],
      "notes": [],
      "sounds": [],
      "regions": [],
      "templates": [],
      "tiles": [],
      "walls": [],
      "playlist": null,
      "playlistSound": null,
      "journal": null,
      "journalEntryPage": null,
      "weather": "",
      "folder": null,
      "sort": 0,
      "ownership": {
        "default": 0,
        "5ypAoBvOiyjDKiaZ": 3
      },
      "flags": {},
      "_stats": {
        "compendiumSource": null,
        "duplicateSource": null,
        "exportSource": null,
        "coreVersion": "13.348",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1776596215113,
        "modifiedTime": 1776656673046,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      },
      "navName": ""
    },
    {
      "grid": {
        "size": 100,
        "type": 1,
        "style": "solidLines",
        "thickness": 1,
        "color": "#000000",
        "alpha": 0.2,
        "distance": 5,
        "units": "ft"
      },
      "height": 1000,
      "name": "test-scene",
      "width": 1000,
      "_id": "r36nfimJGHYGUGQX",
      "active": false,
      "navigation": true,
      "navOrder": 0,
      "background": {
        "src": null,
        "anchorX": 0,
        "anchorY": 0,
        "offsetX": 0,
        "offsetY": 0,
        "fit": "fill",
        "scaleX": 1,
        "scaleY": 1,
        "rotation": 0,
        "tint": "#ffffff",
        "alphaThreshold": 0
      },
      "foreground": null,
      "foregroundElevation": null,
      "thumb": null,
      "padding": 0.25,
      "initial": {
        "x": null,
        "y": null,
        "scale": null
      },
      "backgroundColor": "#999999",
      "tokenVision": true,
      "fog": {
        "exploration": true,
        "overlay": null,
        "colors": {
          "explored": null,
          "unexplored": null
        }
      },
      "environment": {
        "darknessLevel": 0,
        "darknessLock": false,
        "globalLight": {
          "enabled": false,
          "alpha": 0.5,
          "bright": false,
          "color": null,
          "coloration": 1,
          "luminosity": 0,
          "saturation": 0,
          "contrast": 0,
          "shadows": 0,
          "darkness": {
            "min": 0,
            "max": 1
          }
        },
        "cycle": true,
        "base": {
          "hue": 0,
          "intensity": 0,
          "luminosity": 0,
          "saturation": 0,
          "shadows": 0
        },
        "dark": {
          "hue": 0.7138888888888889,
          "intensity": 0,
          "luminosity": -0.25,
          "saturation": 0,
          "shadows": 0
        }
      },
      "drawings": [],
      "tokens": [],
      "lights": [],
      "notes": [],
      "sounds": [],
      "regions": [],
      "templates": [],
      "tiles": [],
      "walls": [],
      "playlist": null,
      "playlistSound": null,
      "journal": null,
      "journalEntryPage": null,
      "weather": "",
      "folder": null,
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
        "createdTime": 1776657957926,
        "modifiedTime": 1776657957926,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    },
    {
      "height": 500,
      "name": "test-scene-expendable",
      "width": 500,
      "_id": "yomijR8uu8kQrjip",
      "active": false,
      "navigation": true,
      "navOrder": 0,
      "background": {
        "src": null,
        "anchorX": 0,
        "anchorY": 0,
        "offsetX": 0,
        "offsetY": 0,
        "fit": "fill",
        "scaleX": 1,
        "scaleY": 1,
        "rotation": 0,
        "tint": "#ffffff",
        "alphaThreshold": 0
      },
      "foreground": null,
      "foregroundElevation": null,
      "thumb": null,
      "padding": 0.25,
      "initial": {
        "x": null,
        "y": null,
        "scale": null
      },
      "backgroundColor": "#999999",
      "grid": {
        "type": 1,
        "size": 100,
        "style": "solidLines",
        "thickness": 1,
        "color": "#000000",
        "alpha": 0.2,
        "distance": 5,
        "units": "ft"
      },
      "tokenVision": true,
      "fog": {
        "exploration": true,
        "overlay": null,
        "colors": {
          "explored": null,
          "unexplored": null
        }
      },
      "environment": {
        "darknessLevel": 0,
        "darknessLock": false,
        "globalLight": {
          "enabled": false,
          "alpha": 0.5,
          "bright": false,
          "color": null,
          "coloration": 1,
          "luminosity": 0,
          "saturation": 0,
          "contrast": 0,
          "shadows": 0,
          "darkness": {
            "min": 0,
            "max": 1
          }
        },
        "cycle": true,
        "base": {
          "hue": 0,
          "intensity": 0,
          "luminosity": 0,
          "saturation": 0,
          "shadows": 0
        },
        "dark": {
          "hue": 0.7138888888888889,
          "intensity": 0,
          "luminosity": -0.25,
          "saturation": 0,
          "shadows": 0
        }
      },
      "drawings": [],
      "tokens": [],
      "lights": [],
      "notes": [],
      "sounds": [],
      "regions": [],
      "templates": [],
      "tiles": [],
      "walls": [],
      "playlist": null,
      "playlistSound": null,
      "journal": null,
      "journalEntryPage": null,
      "weather": "",
      "folder": null,
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
        "createdTime": 1776657957934,
        "modifiedTime": 1776657957934,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    }
  ]
}
```


---

## GET /scene/image/raw

Get the raw background image of a scene

Returns the scene's background image file without any tokens, lights, or other canvas elements rendered on it.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | Scene ID (defaults to viewed/active scene) |
| active | boolean |  | body, query | If true, explicitly use the player-facing active scene instead of the viewed scene |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**binary** - The raw scene background image

### Try It Out

<ApiTester
  method="GET"
  path="/scene/image/raw"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"active","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /scene

Create a new scene

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| data | object | ✓ | body | Scene data object (name, width, height, grid, etc.) |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Created scene data

### Try It Out

<ApiTester
  method="POST"
  path="/scene"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/scene';
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
      "data": {
        "name": "test-scene",
        "width": 1000,
        "height": 1000,
        "grid": {
          "size": 100
        }
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/scene?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"data":{"name":"test-scene","width":1000,"height":1000,"grid":{"size":100}}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/scene'
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
      "data": {
        "name": "test-scene",
        "width": 1000,
        "height": 1000,
        "grid": {
          "size": 100
        }
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
  const path = '/scene';
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
        "data": {
          "name": "test-scene",
          "width": 1000,
          "height": 1000,
          "grid": {
            "size": 100
          }
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
  🔤/scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"data":{"name":"test-scene","width":1000,"height":1000,"grid":{"size":100}}}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 77❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "create-scene-result",
  "requestId": "create-scene_1776657957923",
  "data": {
    "grid": {
      "size": 100,
      "type": 1,
      "style": "solidLines",
      "thickness": 1,
      "color": "#000000",
      "alpha": 0.2,
      "distance": 5,
      "units": "ft"
    },
    "height": 1000,
    "name": "test-scene",
    "width": 1000,
    "_id": "r36nfimJGHYGUGQX",
    "active": false,
    "navigation": true,
    "navOrder": 0,
    "background": {
      "src": null,
      "anchorX": 0,
      "anchorY": 0,
      "offsetX": 0,
      "offsetY": 0,
      "fit": "fill",
      "scaleX": 1,
      "scaleY": 1,
      "rotation": 0,
      "tint": "#ffffff",
      "alphaThreshold": 0
    },
    "foreground": null,
    "foregroundElevation": null,
    "thumb": null,
    "padding": 0.25,
    "initial": {
      "x": null,
      "y": null,
      "scale": null
    },
    "backgroundColor": "#999999",
    "tokenVision": true,
    "fog": {
      "exploration": true,
      "overlay": null,
      "colors": {
        "explored": null,
        "unexplored": null
      }
    },
    "environment": {
      "darknessLevel": 0,
      "darknessLock": false,
      "globalLight": {
        "enabled": false,
        "alpha": 0.5,
        "bright": false,
        "color": null,
        "coloration": 1,
        "luminosity": 0,
        "saturation": 0,
        "contrast": 0,
        "shadows": 0,
        "darkness": {
          "min": 0,
          "max": 1
        }
      },
      "cycle": true,
      "base": {
        "hue": 0,
        "intensity": 0,
        "luminosity": 0,
        "saturation": 0,
        "shadows": 0
      },
      "dark": {
        "hue": 0.7138888888888889,
        "intensity": 0,
        "luminosity": -0.25,
        "saturation": 0,
        "shadows": 0
      }
    },
    "drawings": [],
    "tokens": [],
    "lights": [],
    "notes": [],
    "sounds": [],
    "regions": [],
    "templates": [],
    "tiles": [],
    "walls": [],
    "playlist": null,
    "playlistSound": null,
    "journal": null,
    "journalEntryPage": null,
    "weather": "",
    "folder": null,
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
      "createdTime": 1776657957926,
      "modifiedTime": 1776657957926,
      "lastModifiedBy": "r6bXhB7k9cXa3cif"
    }
  }
}
```


---

## PUT /scene

Update an existing scene

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| data | object | ✓ | body | Object containing the scene fields to update |
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | ID of the scene to update |
| name | string |  | body, query | Name of the scene to update (alternative to sceneId) |
| active | boolean |  | body, query | Set to true to target the active scene |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Updated scene data

### Try It Out

<ApiTester
  method="PUT"
  path="/scene"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"active","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/scene';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3'
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
      "sceneId": "r36nfimJGHYGUGQX",
      "data": {
        "name": "test-scene-updated"
      }
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X PUT 'http://localhost:3010/scene?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"sceneId":"r36nfimJGHYGUGQX","data":{"name":"test-scene-updated"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/scene'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3'
}
url = f'{base_url}{path}'

response = requests.put(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "sceneId": "r36nfimJGHYGUGQX",
      "data": {
        "name": "test-scene-updated"
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
  const path = '/scene';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3'
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
        "sceneId": "r36nfimJGHYGUGQX",
        "data": {
          "name": "test-scene-updated"
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
  🔤/scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"sceneId":"r36nfimJGHYGUGQX","data":{"name":"test-scene-updated"}}🔤 ➡️ body

  💭 Build HTTP request
  🔤PUT /scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 67❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "update-scene-result",
  "requestId": "update-scene_1776657957945",
  "data": {
    "grid": {
      "size": 100,
      "type": 1,
      "style": "solidLines",
      "thickness": 1,
      "color": "#000000",
      "alpha": 0.2,
      "distance": 5,
      "units": "ft"
    },
    "height": 1000,
    "name": "test-scene-updated",
    "width": 1000,
    "_id": "r36nfimJGHYGUGQX",
    "active": false,
    "navigation": true,
    "navOrder": 0,
    "background": {
      "src": null,
      "anchorX": 0,
      "anchorY": 0,
      "offsetX": 0,
      "offsetY": 0,
      "fit": "fill",
      "scaleX": 1,
      "scaleY": 1,
      "rotation": 0,
      "tint": "#ffffff",
      "alphaThreshold": 0
    },
    "foreground": null,
    "foregroundElevation": null,
    "thumb": null,
    "padding": 0.25,
    "initial": {
      "x": null,
      "y": null,
      "scale": null
    },
    "backgroundColor": "#999999",
    "tokenVision": true,
    "fog": {
      "exploration": true,
      "overlay": null,
      "colors": {
        "explored": null,
        "unexplored": null
      }
    },
    "environment": {
      "darknessLevel": 0,
      "darknessLock": false,
      "globalLight": {
        "enabled": false,
        "alpha": 0.5,
        "bright": false,
        "color": null,
        "coloration": 1,
        "luminosity": 0,
        "saturation": 0,
        "contrast": 0,
        "shadows": 0,
        "darkness": {
          "min": 0,
          "max": 1
        }
      },
      "cycle": true,
      "base": {
        "hue": 0,
        "intensity": 0,
        "luminosity": 0,
        "saturation": 0,
        "shadows": 0
      },
      "dark": {
        "hue": 0.7138888888888889,
        "intensity": 0,
        "luminosity": -0.25,
        "saturation": 0,
        "shadows": 0
      }
    },
    "drawings": [],
    "tokens": [],
    "lights": [],
    "notes": [],
    "sounds": [],
    "regions": [],
    "templates": [],
    "tiles": [],
    "walls": [],
    "playlist": null,
    "playlistSound": null,
    "journal": null,
    "journalEntryPage": null,
    "weather": "",
    "folder": null,
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
      "createdTime": 1776657957926,
      "modifiedTime": 1776657957946,
      "lastModifiedBy": "r6bXhB7k9cXa3cif"
    }
  }
}
```


---

## DELETE /scene

Delete a scene

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | query | ID of the scene to delete |
| name | string |  | query | Name of the scene to delete (alternative to sceneId) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Deletion result

### Try It Out

<ApiTester
  method="DELETE"
  path="/scene"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/scene';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  sceneId: 'yomijR8uu8kQrjip'
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
curl -X DELETE 'http://localhost:3010/scene?clientId=fvtt_099ad17ea199e7e3&sceneId=yomijR8uu8kQrjip' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/scene'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'sceneId': 'yomijR8uu8kQrjip'
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
  const path = '/scene';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    sceneId: 'yomijR8uu8kQrjip'
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
  🔤/scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤sceneId=yomijR8uu8kQrjip🔤 ➡️ sceneId
  🔤?🧲clientId🧲&🧲sceneId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "delete-scene-result",
  "requestId": "delete-scene_1776657962962",
  "success": true
}
```


---

## POST /switch-scene

Switch the active scene

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| sceneId | string |  | body, query | ID of the scene to activate |
| name | string |  | body, query | Name of the scene to activate (alternative to sceneId) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the scene switch

### Try It Out

<ApiTester
  method="POST"
  path="/switch-scene"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/switch-scene';
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
      "sceneId": "r36nfimJGHYGUGQX"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/switch-scene?clientId=fvtt_099ad17ea199e7e3' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"sceneId":"r36nfimJGHYGUGQX"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/switch-scene'
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
      "sceneId": "r36nfimJGHYGUGQX"
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
  const path = '/switch-scene';
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
        "sceneId": "r36nfimJGHYGUGQX"
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
  🔤/switch-scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"sceneId":"r36nfimJGHYGUGQX"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /switch-scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 30❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "switch-scene-result",
  "requestId": "switch-scene_1776657957950",
  "success": true,
  "data": {
    "grid": {
      "size": 100,
      "type": 1,
      "style": "solidLines",
      "thickness": 1,
      "color": "#000000",
      "alpha": 0.2,
      "distance": 5,
      "units": "ft"
    },
    "height": 1000,
    "name": "test-scene-updated",
    "width": 1000,
    "_id": "r36nfimJGHYGUGQX",
    "active": true,
    "navigation": true,
    "navOrder": 0,
    "background": {
      "src": null,
      "anchorX": 0,
      "anchorY": 0,
      "offsetX": 0,
      "offsetY": 0,
      "fit": "fill",
      "scaleX": 1,
      "scaleY": 1,
      "rotation": 0,
      "tint": "#ffffff",
      "alphaThreshold": 0
    },
    "foreground": null,
    "foregroundElevation": null,
    "thumb": null,
    "padding": 0.25,
    "initial": {
      "x": null,
      "y": null,
      "scale": null
    },
    "backgroundColor": "#999999",
    "tokenVision": true,
    "fog": {
      "exploration": true,
      "overlay": null,
      "colors": {
        "explored": null,
        "unexplored": null
      }
    },
    "environment": {
      "darknessLevel": 0,
      "darknessLock": false,
      "globalLight": {
        "enabled": false,
        "alpha": 0.5,
        "bright": false,
        "color": null,
        "coloration": 1,
        "luminosity": 0,
        "saturation": 0,
        "contrast": 0,
        "shadows": 0,
        "darkness": {
          "min": 0,
          "max": 1
        }
      },
      "cycle": true,
      "base": {
        "hue": 0,
        "intensity": 0,
        "luminosity": 0,
        "saturation": 0,
        "shadows": 0
      },
      "dark": {
        "hue": 0.7138888888888889,
        "intensity": 0,
        "luminosity": -0.25,
        "saturation": 0,
        "shadows": 0
      }
    },
    "drawings": [],
    "tokens": [],
    "lights": [],
    "notes": [],
    "sounds": [],
    "regions": [],
    "templates": [],
    "tiles": [],
    "walls": [],
    "playlist": null,
    "playlistSound": null,
    "journal": null,
    "journalEntryPage": null,
    "weather": "",
    "folder": null,
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
      "createdTime": 1776657957926,
      "modifiedTime": 1776657957952,
      "lastModifiedBy": "r6bXhB7k9cXa3cif"
    }
  }
}
```


---

## GET /scene/image

Get a rendered screenshot of a scene

Captures the full rendered canvas of a scene including all visible layers (tokens, lights, walls, etc.) as an image. The scene can be specified by ID or defaults to the active scene.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| sceneId | string |  | query | Scene ID (defaults to viewed/active scene) |
| active | boolean |  | query | If true, explicitly use the player-facing active scene instead of the viewed scene |
| clientId | string |  | query | Client ID for the Foundry world |
| format | string |  | query | Image format: png or jpeg (default: png) |
| quality | number |  | query | Image quality 0-1 for JPEG (default: 0.9) |
| viewport | boolean |  | query | If true, capture exactly what the browser currently shows instead of the full scene |
| width | number |  | query | Output image width in pixels (default: scene width) |
| height | number |  | query | Output image height in pixels (default: scene height) |
| showGrid | boolean |  | query | Include grid lines in capture (default: false) |
| hideOverlays | boolean |  | query | Hide fog of war, weather, vision, and UI overlays (default: false) |
| userId | string |  | query | Foundry user ID or username |

### Returns

**binary** - The scene screenshot as an image

### Try It Out

<ApiTester
  method="GET"
  path="/scene/image"
  parameters={[{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"active","type":"boolean","required":false,"source":"query"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"format","type":"string","required":false,"source":"query"},{"name":"quality","type":"number","required":false,"source":"query"},{"name":"viewport","type":"boolean","required":false,"source":"query"},{"name":"width","type":"number","required":false,"source":"query"},{"name":"height","type":"number","required":false,"source":"query"},{"name":"showGrid","type":"boolean","required":false,"source":"query"},{"name":"hideOverlays","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

