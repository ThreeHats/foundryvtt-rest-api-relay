---
tag: structure
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


# structure

## GET /structure

Get the structure of the Foundry world Retrieves the folder and compendium structure for the specified Foundry world.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | query | Client ID for the Foundry world |
| includeEntityData | boolean |  | query | Whether to include full entity data or just UUIDs and names |
| path | string |  | query | Path to read structure from (null = root) |
| recursive | boolean |  | query | Whether to read down the folder tree |
| recursiveDepth | number |  | query | Depth to recurse into folders (default 5) |
| types | string |  | query | Types to return (Scene/Actor/Item/JournalEntry/RollTable/Cards/Macro/Playlist), can be comma-separated or JSON array |

### Returns

**object** - The folder and compendium structure

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/structure';
const params = {
  clientId: 'your-client-id',
  includeEntityData: 'true',
  recursive: 'true',
  types: 'Scene'
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
curl -X GET 'http://localhost:3010/structure?clientId=your-client-id&includeEntityData=true&recursive=true&types=Scene' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/structure'
params = {
    'clientId': 'your-client-id',
    'includeEntityData': 'true',
    'recursive': 'true',
    'types': 'Scene'
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
  const path = '/structure';
  const params = {
    clientId: 'your-client-id',
    includeEntityData: 'true',
    recursive: 'true',
    types: 'Scene'
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
  🔤/structure🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤includeEntityData=true🔤 ➡️ includeEntityData
  🔤recursive=true🔤 ➡️ recursive
  🔤types=Scene🔤 ➡️ types
  🔤?🧲clientId🧲&🧲includeEntityData🧲&🧲recursive🧲&🧲types🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /structure🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "structure_1765635977185",
  "clientId": "your-client-id",
  "type": "structure-result",
  "data": {
    "folders": {
      "test-folder": {
        "id": "SsYlZRbKfAr9WHOP",
        "uuid": "Folder.SsYlZRbKfAr9WHOP",
        "type": "Scene"
      }
    },
    "entities": {
      "scenes": [
        {
          "_id": "NUEDEFAULTSCENE0",
          "name": "Foundry Virtual Tabletop",
          "active": true,
          "navigation": true,
          "navOrder": 0,
          "navName": "",
          "background": {
            "src": "nue/defaultscene/fvtt-background.webp",
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
          "foregroundElevation": 4,
          "thumb": "nue/defaultscene/thumb.webp",
          "width": 3840,
          "height": 1920,
          "padding": 0,
          "initial": {
            "x": null,
            "y": null,
            "scale": null
          },
          "backgroundColor": "#25070d",
          "grid": {
            "type": 0,
            "size": 100,
            "style": "solidLines",
            "thickness": 1,
            "color": "#000000",
            "alpha": 0.2,
            "distance": 1,
            "units": ""
          },
          "tokenVision": false,
          "fog": {
            "exploration": false,
            "reset": 1660769143211,
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
          "tokens": [
            {
              "_id": "O4sEnBrG5I3lFNGk",
              "name": "test",
              "displayName": 0,
              "actorId": "xctPu6799LkAP6p3",
              "actorLink": true,
              "delta": null,
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
              "shape": 4,
              "x": 1596,
              "y": 623,
              "elevation": 0,
              "sort": 0,
              "locked": false,
              "lockRotation": false,
              "rotation": 0,
              "alpha": 1,
              "hidden": false,
              "disposition": 1,
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
                "enabled": true,
                "range": 0,
                "angle": 360,
                "visionMode": "basic",
                "color": null,
                "attenuation": 0.1,
                "brightness": 0,
                "saturation": 0,
                "contrast": 0
              },
              "detectionModes": [
                {
                  "id": "lightPerception",
                  "enabled": true,
                  "range": null
                },
                {
                  "id": "basicSight",
                  "enabled": true,
                  "range": 0
                }
              ],
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
              "movementAction": "walk",
              "_movementHistory": [],
              "_regions": [],
              "flags": {}
            }
          ],
          "lights": [
            {
              "_id": "d22Cax8HDPMG4F6I",
              "x": 656,
              "y": 1473,
              "elevation": 0,
              "rotation": 0,
              "walls": true,
              "vision": false,
              "config": {
                "negative": false,
                "priority": 0,
                "alpha": 0.5,
                "angle": 360,
                "bright": 4.27,
                "color": "#ff9072",
                "coloration": 1,
                "dim": 8.54,
                "attenuation": 0.5,
                "luminosity": 0,
                "saturation": 0,
                "contrast": 0,
                "shadows": 0,
                "animation": {
                  "type": "fog",
                  "speed": 2,
                  "intensity": 5,
                  "reverse": false
                },
                "darkness": {
                  "min": 0,
                  "max": 1
                }
              },
              "hidden": false,
              "flags": {}
            },
            {
              "_id": "eGuMjw01vEYimWVX",
              "x": 1826,
              "y": 1891,
              "elevation": 0,
              "rotation": 0,
              "walls": true,
              "vision": false,
              "config": {
                "negative": false,
                "priority": 0,
                "alpha": 0.5,
                "angle": 360,
                "bright": 4.27,
                "color": "#ffffff",
                "coloration": 1,
                "dim": 8.54,
                "attenuation": 0.5,
                "luminosity": 0,
                "saturation": 0,
                "contrast": 0,
                "shadows": 0,
                "animation": {
                  "type": "fog",
                  "speed": 2,
                  "intensity": 5,
                  "reverse": false
                },
                "darkness": {
                  "min": 0,
                  "max": 1
                }
              },
              "hidden": false,
              "flags": {}
            },
            {
              "_id": "TCET4ZNPkl5oZukY",
              "x": 3057,
              "y": 1439,
              "elevation": 0,
              "rotation": 0,
              "walls": true,
              "vision": false,
              "config": {
                "negative": false,
                "priority": 0,
                "alpha": 0.5,
                "angle": 360,
                "bright": 4.27,
                "color": "#ffffff",
                "coloration": 1,
                "dim": 8.54,
                "attenuation": 0.5,
                "luminosity": 0,
                "saturation": 0,
                "contrast": 0,
                "shadows": 0,
                "animation": {
                  "type": "fog",
                  "speed": 2,
                  "intensity": 5,
                  "reverse": false
                },
                "darkness": {
                  "min": 0,
                  "max": 1
                }
              },
              "hidden": false,
              "flags": {}
            },
            {
              "_id": "cOpD0Q4AuCGiKRCb",
              "x": 2824,
              "y": 772,
              "elevation": 0,
              "rotation": 0,
              "walls": true,
              "vision": false,
              "config": {
                "negative": false,
                "priority": 0,
                "alpha": 0.5,
                "angle": 360,
                "bright": 0.26,
                "color": "#ffed79",
                "coloration": 1,
                "dim": 0.53,
                "attenuation": 0.5,
                "luminosity": 0.5,
                "saturation": 0,
                "contrast": 0,
                "shadows": 0,
                "animation": {
                  "type": "torch",
                  "speed": 5,
                  "intensity": 5,
                  "reverse": false
                },
                "darkness": {
                  "min": 0,
                  "max": 1
                }
              },
              "hidden": false,
              "flags": {}
            },
            {
              "_id": "adhkydxURYamgKKF",
              "x": 2822,
              "y": 777,
              "elevation": 0,
              "rotation": 179,
              "walls": true,
              "vision": false,
              "config": {
                "negative": false,
                "priority": 0,
                "alpha": 0.65,
                "angle": 30,
                "bright": 5,
                "color": "#ff976f",
                "coloration": 1,
                "dim": 5,
                "attenuation": 0.7,
                "luminosity": 0.3,
                "saturation": 0,
                "contrast": 0,
                "shadows": 0,
                "animation": {
                  "type": "sunburst",
                  "speed": 1,
                  "intensity": 7,
                  "reverse": false
                },
                "darkness": {
                  "min": 0,
                  "max": 1
                }
              },
              "hidden": false,
              "flags": {}
            }
          ],
          "notes": [],
          "sounds": [],
          "regions": [],
          "templates": [],
          "tiles": [
            {
              "_id": "mMxIUI1fXJmrR1zK",
              "texture": {
                "src": "nue/defaultscene/fvtt-logo.webp",
                "anchorX": 0.5,
                "anchorY": 0.5,
                "offsetX": 0,
                "offsetY": 0,
                "fit": "fill",
                "scaleX": 1,
                "scaleY": 1,
                "rotation": 0,
                "tint": "#ffffff",
                "alphaThreshold": 0.75
              },
              "width": 800,
              "height": 800,
              "x": 1520,
              "y": 480,
              "elevation": 0,
              "sort": 100,
              "rotation": 0,
              "alpha": 1,
              "hidden": false,
              "locked": false,
              "restrictions": {
                "light": false,
                "weather": false
              },
              "occlusion": {
                "mode": 1,
                "alpha": 0
              },
              "video": {
                "loop": true,
                "autoplay": true,
                "volume": 0
              },
              "flags": {}
            }
          ],
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1763765287462,
            "modifiedTime": 1763765287462,
            "lastModifiedBy": "5ypAoBvOiyjDKiaZ",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "VnnIYuJJjlZzUeRT",
          "name": "a",
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
          "foregroundElevation": 20,
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1765382040878,
            "modifiedTime": 1765382040878,
            "lastModifiedBy": "5ypAoBvOiyjDKiaZ",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        }
      ]
    }
  }
}
```


---

## GET /contents/:path

This route is deprecated - use /structure with the path query parameter instead

### Returns

**object** - Error message directing to use /structure endpoint

---

## GET /get-folder

Get a specific folder by name

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| name | string | ✓ | body, query | Name of the folder to retrieve |

### Returns

**object** - The folder information and its contents

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/get-folder';
const params = {
  clientId: 'your-client-id',
  name: 'test-folder'
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
curl -X GET 'http://localhost:3010/get-folder?clientId=your-client-id&name=test-folder' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/get-folder'
params = {
    'clientId': 'your-client-id',
    'name': 'test-folder'
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
  const path = '/get-folder';
  const params = {
    clientId: 'your-client-id',
    name: 'test-folder'
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
  🔤/get-folder🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤name=test-folder🔤 ➡️ name
  🔤?🧲clientId🧲&🧲name🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /get-folder🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "get-folder_1765635977461",
  "clientId": "your-client-id",
  "type": "get-folder-result",
  "data": {
    "id": "SsYlZRbKfAr9WHOP",
    "uuid": "Folder.SsYlZRbKfAr9WHOP",
    "name": "test-folder",
    "type": "Scene",
    "parentFolder": null,
    "contents": []
  }
}
```


---

## POST /create-folder

Create a new folder

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| name | string | ✓ | body, query | Name of the new folder |
| folderType | string | ✓ | body, query | Type of folder (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist) |
| parentFolderId | string |  | body, query | ID of the parent folder (optional for root level) |

### Returns

**object** - The created folder information

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/create-folder';
const params = {
  clientId: 'your-client-id',
  name: 'test-folder',
  folderType: 'Scene'
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
curl -X POST 'http://localhost:3010/create-folder?clientId=your-client-id&name=test-folder&folderType=Scene' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/create-folder'
params = {
    'clientId': 'your-client-id',
    'name': 'test-folder',
    'folderType': 'Scene'
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
  const path = '/create-folder';
  const params = {
    clientId: 'your-client-id',
    name: 'test-folder',
    folderType: 'Scene'
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
  🔤/create-folder🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤name=test-folder🔤 ➡️ name
  🔤folderType=Scene🔤 ➡️ folderType
  🔤?🧲clientId🧲&🧲name🧲&🧲folderType🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤POST /create-folder🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "create-folder_1765635976871",
  "clientId": "your-client-id",
  "type": "create-folder-result",
  "data": {
    "id": "SsYlZRbKfAr9WHOP",
    "uuid": "Folder.SsYlZRbKfAr9WHOP",
    "name": "test-folder",
    "type": "Scene",
    "parentFolder": null
  }
}
```


---

## DELETE /delete-folder

Delete a folder

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| folderId | string | ✓ | body, query | ID of the folder to delete |
| deleteAll | boolean |  | body, query | Whether to delete all entities in the folder |

### Returns

**object** - Confirmation of deletion

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/delete-folder';
const params = {
  clientId: 'your-client-id',
  folderId: 'SsYlZRbKfAr9WHOP'
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
curl -X DELETE 'http://localhost:3010/delete-folder?clientId=your-client-id&folderId=SsYlZRbKfAr9WHOP' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/delete-folder'
params = {
    'clientId': 'your-client-id',
    'folderId': 'SsYlZRbKfAr9WHOP'
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
  const path = '/delete-folder';
  const params = {
    clientId: 'your-client-id',
    folderId: 'SsYlZRbKfAr9WHOP'
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
  🔤/delete-folder🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤folderId=SsYlZRbKfAr9WHOP🔤 ➡️ folderId
  🔤?🧲clientId🧲&🧲folderId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /delete-folder🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "delete-folder_1765635977769",
  "clientId": "your-client-id",
  "type": "delete-folder-result",
  "data": {
    "deleted": true,
    "folderId": "SsYlZRbKfAr9WHOP",
    "entitiesDeleted": 0,
    "foldersDeleted": 1
  }
}
```

