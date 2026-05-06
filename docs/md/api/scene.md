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
  "requestId": "get-scene_1777909249375",
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
        "modifiedTime": 1776821577409,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    },
    {
      "_id": "OoTPjYkL2GjuQ2a7",
      "_stats": {
        "compendiumSource": null,
        "duplicateSource": "Scene.Fpu0odizNnNnjlLI",
        "exportSource": null,
        "coreVersion": "13.348",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1777176992560,
        "modifiedTime": 1777176992560,
        "lastModifiedBy": "5ypAoBvOiyjDKiaZ"
      },
      "active": false,
      "background": {
        "alphaThreshold": 0,
        "anchorX": 0,
        "anchorY": 0,
        "fit": "fill",
        "offsetX": 0,
        "offsetY": 0,
        "rotation": 0,
        "scaleX": 1,
        "scaleY": 1,
        "src": "flooded-cave-test.webp",
        "tint": "#ffffff"
      },
      "backgroundColor": "#999999",
      "drawings": [
        {
          "_id": "7us1TaI1NwwGdqLs",
          "author": "5ypAoBvOiyjDKiaZ",
          "bezierFactor": 0.5,
          "elevation": 0,
          "fillAlpha": 0.5,
          "fillColor": "#ffffff",
          "fillType": 0,
          "flags": {},
          "fontFamily": "Signika",
          "fontSize": 48,
          "hidden": false,
          "interface": false,
          "locked": false,
          "rotation": 0,
          "shape": {
            "height": 342,
            "points": [],
            "radius": null,
            "type": "r",
            "width": 287
          },
          "sort": 0,
          "strokeAlpha": 1,
          "strokeColor": "#ffffff",
          "strokeWidth": 8,
          "textAlpha": 1,
          "textColor": "#ffffff",
          "texture": null,
          "x": 1942.5,
          "y": 1581.75
        },
        {
          "_id": "nrcBVMkwWw8XPcBj",
          "author": "5ypAoBvOiyjDKiaZ",
          "bezierFactor": 0.5,
          "elevation": 0,
          "fillAlpha": 0.5,
          "fillColor": "#ffffff",
          "fillType": 0,
          "flags": {},
          "fontFamily": "Signika",
          "fontSize": 48,
          "hidden": false,
          "interface": false,
          "locked": false,
          "rotation": 0,
          "shape": {
            "height": 583,
            "points": [],
            "radius": null,
            "type": "r",
            "width": 518
          },
          "sort": 1,
          "strokeAlpha": 1,
          "strokeColor": "#ffffff",
          "strokeWidth": 8,
          "textAlpha": 1,
          "textColor": "#ffffff",
          "texture": null,
          "x": 1359.75,
          "y": 2025.75
        }
      ],
      "environment": {
        "base": {
          "hue": 0,
          "intensity": 0,
          "luminosity": 0,
          "saturation": 0,
          "shadows": 0
        },
        "cycle": true,
        "dark": {
          "hue": 0.7138888888888889,
          "intensity": 0,
          "luminosity": -0.25,
          "saturation": 0,
          "shadows": 0
        },
        "darknessLevel": 0,
        "darknessLock": false,
        "globalLight": {
          "alpha": 0.5,
          "bright": false,
          "color": null,
          "coloration": 1,
          "contrast": 0,
          "darkness": {
            "max": 1,
            "min": 0
          },
          "enabled": false,
          "luminosity": 0,
          "saturation": 0,
          "shadows": 0
        }
      },
      "flags": {},
      "fog": {
        "colors": {
          "explored": null,
          "unexplored": null
        },
        "exploration": true,
        "overlay": null
      },
      "foreground": null,
      "foregroundElevation": null,
      "grid": {
        "alpha": 0.2,
        "color": "#000000",
        "distance": 5,
        "size": 74,
        "style": "solidLines",
        "thickness": 1,
        "type": 1,
        "units": "ft"
      },
      "height": 2962,
      "initial": {
        "scale": null,
        "x": null,
        "y": null
      },
      "journal": null,
      "journalEntryPage": null,
      "lights": [
        {
          "_id": "RFgXdhctPXhpSyOM",
          "config": {
            "alpha": 0.5,
            "angle": 360,
            "animation": {
              "intensity": 5,
              "reverse": false,
              "speed": 5,
              "type": null
            },
            "attenuation": 0.5,
            "bright": 24.04,
            "color": null,
            "coloration": 1,
            "contrast": 0,
            "darkness": {
              "max": 1,
              "min": 0
            },
            "dim": 48.09,
            "luminosity": 0.5,
            "negative": false,
            "priority": 0,
            "saturation": 0,
            "shadows": 0
          },
          "elevation": 0,
          "flags": {},
          "hidden": false,
          "rotation": 0,
          "vision": false,
          "walls": true,
          "x": 1628,
          "y": 1443
        }
      ],
      "name": "test",
      "navName": "",
      "navOrder": 0,
      "navigation": false,
      "notes": [
        {
          "_id": "YqzoO6yW51CQtt0k",
          "elevation": 0,
          "entryId": "u7byD1yDxgtzqeT4",
          "flags": {},
          "fontFamily": "Signika",
          "fontSize": 32,
          "global": false,
          "iconSize": 40,
          "pageId": null,
          "sort": 0,
          "text": "jhgcfkjhgcv",
          "textAnchor": 1,
          "textColor": "#ffffff",
          "texture": {
            "alphaThreshold": 0,
            "anchorX": 0.5,
            "anchorY": 0.5,
            "fit": "contain",
            "offsetX": 0,
            "offsetY": 0,
            "rotation": 0,
            "scaleX": 1,
            "scaleY": 1,
            "src": "icons/svg/book.svg",
            "tint": "#ffffff"
          },
          "x": 1665,
          "y": 2109
        }
      ],
      "ownership": {
        "5ypAoBvOiyjDKiaZ": 3,
        "default": 0,
        "lmPGreUzUZ4YxH6D": 3
      },
      "padding": 0.25,
      "playlist": null,
      "playlistSound": null,
      "regions": [
        {
          "_id": "PaZoZJk6XppvofiN",
          "behaviors": [],
          "color": "#2893cc",
          "elevation": {
            "bottom": null,
            "top": null
          },
          "flags": {},
          "locked": false,
          "name": "Region",
          "shapes": [
            {
              "height": 407,
              "hole": false,
              "rotation": 0,
              "type": "rectangle",
              "width": 582.75,
              "x": 1433.75,
              "y": 2016.5
            },
            {
              "height": 786.25,
              "hole": false,
              "rotation": 0,
              "type": "rectangle",
              "width": 471.75,
              "x": 1757.5,
              "y": 2358.75
            },
            {
              "hole": false,
              "radiusX": 296,
              "radiusY": 74,
              "rotation": 0,
              "type": "ellipse",
              "x": 1729.75,
              "y": 1877.75
            }
          ],
          "visibility": 0
        }
      ],
      "sort": 0,
      "sounds": [
        {
          "_id": "kWlgWosbxU3tbSj0",
          "darkness": {
            "max": 1,
            "min": 0
          },
          "easing": true,
          "effects": {
            "base": {
              "intensity": 5,
              "type": ""
            },
            "muffled": {
              "intensity": 5
            }
          },
          "elevation": 0,
          "flags": {},
          "hidden": false,
          "path": null,
          "radius": 48.02,
          "repeat": false,
          "volume": 0.5,
          "walls": true,
          "x": 1147,
          "y": 2405
        }
      ],
      "templates": [
        {
          "_id": "iJq4Fzo3Kvln5zVP",
          "angle": 53.13,
          "author": "5ypAoBvOiyjDKiaZ",
          "borderColor": "#000000",
          "direction": 275.5275401516562,
          "distance": 77.86205751198719,
          "elevation": 0,
          "fillColor": "#28cca2",
          "flags": {},
          "hidden": false,
          "sort": 0,
          "t": "cone",
          "texture": null,
          "width": 0,
          "x": 1332,
          "y": 2886
        }
      ],
      "thumb": "worlds/5e-tables/assets/scenes/OoTPjYkL2GjuQ2a7-thumb.webp",
      "tiles": [
        {
          "_id": "b0cYOogM1zPL9NFc",
          "alpha": 1,
          "elevation": 0,
          "flags": {},
          "height": 333,
          "hidden": false,
          "locked": false,
          "occlusion": {
            "alpha": 0,
            "mode": 0
          },
          "restrictions": {
            "light": false,
            "weather": false
          },
          "rotation": 0,
          "sort": 0,
          "texture": {
            "alphaThreshold": 0.75,
            "anchorX": 0.5,
            "anchorY": 0.5,
            "fit": "fill",
            "offsetX": 0,
            "offsetY": 0,
            "rotation": 0,
            "scaleX": 1,
            "scaleY": 1,
            "src": null,
            "tint": "#ffffff"
          },
          "video": {
            "autoplay": true,
            "loop": true,
            "volume": 0
          },
          "width": 222,
          "x": 1628,
          "y": 1369
        }
      ],
      "tokenVision": true,
      "tokens": [
        {
          "_id": "ek1v2zxBxkZEftca",
          "_movementHistory": [],
          "_regions": [],
          "actorId": "xSst5kAigAZw6wDr",
          "actorLink": true,
          "alpha": 1,
          "bar1": {
            "attribute": "attributes.hp"
          },
          "bar2": {
            "attribute": "attributes.ac.value"
          },
          "delta": {
            "_id": "ek1v2zxBxkZEftca",
            "effects": [],
            "flags": {},
            "img": null,
            "items": [],
            "name": null,
            "ownership": null,
            "system": {},
            "type": null
          },
          "detectionModes": [],
          "displayBars": 40,
          "displayName": 30,
          "disposition": 1,
          "elevation": 0,
          "flags": {},
          "height": 1,
          "hidden": false,
          "light": {
            "alpha": 1,
            "angle": 360,
            "animation": {
              "intensity": 5,
              "reverse": false,
              "speed": 5,
              "type": null
            },
            "attenuation": 0.5,
            "bright": 0,
            "color": null,
            "coloration": 1,
            "contrast": 0,
            "darkness": {
              "max": 1,
              "min": 0
            },
            "dim": 0,
            "luminosity": 0.5,
            "negative": false,
            "priority": 0,
            "saturation": 0,
            "shadows": 0
          },
          "lockRotation": false,
          "locked": false,
          "movementAction": null,
          "name": "Perrin",
          "occludable": {
            "radius": 0
          },
          "ring": {
            "colors": {
              "background": null,
              "ring": null
            },
            "effects": 1,
            "enabled": false,
            "subject": {
              "scale": 1,
              "texture": null
            }
          },
          "rotation": 0,
          "shape": 4,
          "sight": {
            "angle": 360,
            "attenuation": 0.1,
            "brightness": 0,
            "color": null,
            "contrast": 0,
            "enabled": true,
            "range": 5,
            "saturation": 0,
            "visionMode": "basic"
          },
          "sort": 0,
          "texture": {
            "alphaThreshold": 0.75,
            "anchorX": 0.5,
            "anchorY": 0.5,
            "fit": "contain",
            "offsetX": 0,
            "offsetY": 0,
            "rotation": 0,
            "scaleX": 0.8,
            "scaleY": 0.8,
            "src": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
            "tint": "#ffffff"
          },
          "turnMarker": {
            "animation": null,
            "disposition": false,
            "mode": 1,
            "src": null
          },
          "width": 1,
          "x": 1480,
          "y": 1554
        }
      ],
      "walls": [
        {
          "_id": "LSziSUTponH836zp",
          "animation": null,
          "c": [
            1286,
            1674,
            1656,
            2581
          ],
          "dir": 0,
          "door": 0,
          "ds": 0,
          "flags": {},
          "light": 20,
          "move": 20,
          "sight": 20,
          "sound": 20,
          "threshold": {
            "attenuation": false,
            "light": null,
            "sight": null,
            "sound": null
          }
        },
        {
          "_id": "QICTFX7JLm3W8PLs",
          "animation": null,
          "c": [
            1970,
            1711,
            2017,
            2840
          ],
          "dir": 0,
          "door": 0,
          "ds": 0,
          "flags": {},
          "light": 20,
          "move": 20,
          "sight": 20,
          "sound": 20,
          "threshold": {
            "attenuation": false,
            "light": null,
            "sight": null,
            "sound": null
          }
        },
        {
          "_id": "dOBYPT8ja94pUuFV",
          "animation": null,
          "c": [
            1230,
            2701,
            2424,
            1452
          ],
          "dir": 0,
          "door": 0,
          "ds": 0,
          "flags": {},
          "light": 20,
          "move": 20,
          "sight": 20,
          "sound": 20,
          "threshold": {
            "attenuation": false,
            "light": null,
            "sight": null,
            "sound": null
          }
        },
        {
          "_id": "wo7Vxm6QrHhI7ySh",
          "animation": null,
          "c": [
            1822,
            1434,
            2479,
            2951
          ],
          "dir": 0,
          "door": 0,
          "ds": 0,
          "flags": {},
          "light": 20,
          "move": 20,
          "sight": 20,
          "sound": 20,
          "threshold": {
            "attenuation": false,
            "light": null,
            "sight": null,
            "sound": null
          }
        },
        {
          "_id": "qhAnrbcpLw5UCnJN",
          "animation": null,
          "c": [
            1989,
            1591,
            1203,
            3016
          ],
          "dir": 0,
          "door": 0,
          "ds": 0,
          "flags": {},
          "light": 20,
          "move": 20,
          "sight": 20,
          "sound": 20,
          "threshold": {
            "attenuation": false,
            "light": null,
            "sight": null,
            "sound": null
          }
        },
        {
          "_id": "1wlEJl7SD94Ia2iY",
          "animation": null,
          "c": [
            1277,
            1721,
            2387,
            2035
          ],
          "dir": 0,
          "door": 0,
          "ds": 0,
          "flags": {},
          "light": 20,
          "move": 20,
          "sight": 20,
          "sound": 20,
          "threshold": {
            "attenuation": false,
            "light": null,
            "sight": null,
            "sound": null
          }
        }
      ],
      "weather": "",
      "width": 2221,
      "folder": null
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
      "name": "test-scene-updated",
      "width": 1000,
      "_id": "7iYl9ExwMdFm9POw",
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
        "createdTime": 1776824906803,
        "modifiedTime": 1776824911917,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
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
      "name": "test-scene-updated",
      "width": 1000,
      "_id": "aQADc2ek0f7ls9af",
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
      "tokens": [
        {
          "actorId": "5bW7ahcfLR0uMb9j",
          "x": 200,
          "y": 200,
          "shape": 4,
          "_id": "91dY5wnZ9pbcxcC0",
          "name": "",
          "displayName": 0,
          "actorLink": false,
          "delta": {
            "_id": "xhCuAiQoswAzsnqY",
            "system": {},
            "items": [],
            "effects": [
              {
                "name": "Dead",
                "img": "systems/dnd5e/icons/svg/statuses/dead.svg",
                "_id": "dnd5edead0000000",
                "statuses": [
                  "dead"
                ],
                "type": "base",
                "system": {},
                "changes": [],
                "disabled": false,
                "duration": {
                  "startTime": 6,
                  "combat": null
                },
                "description": "",
                "origin": null,
                "tint": "#ffffff",
                "transfer": false,
                "sort": 0,
                "flags": {},
                "_stats": {
                  "compendiumSource": null,
                  "duplicateSource": null,
                  "exportSource": null,
                  "coreVersion": "13.348",
                  "systemId": "dnd5e",
                  "systemVersion": "5.0.4",
                  "createdTime": 1776821644329,
                  "modifiedTime": 1776821644329,
                  "lastModifiedBy": "r6bXhB7k9cXa3cif"
                }
              }
            ],
            "flags": {},
            "name": null,
            "type": null,
            "img": null,
            "ownership": null
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
        },
        {
          "name": "test",
          "shape": 4,
          "_id": "a0alpONow8ppf7jZ",
          "displayName": 0,
          "actorId": null,
          "actorLink": false,
          "delta": {
            "_id": "TaWAifIKUZ3FSfEj",
            "system": {},
            "items": [],
            "effects": [],
            "flags": {},
            "name": null,
            "type": null,
            "img": null,
            "ownership": null
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
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 1,
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
        },
        {
          "name": "test",
          "shape": 4,
          "_id": "9H8zRN6jFYAKi4St",
          "displayName": 0,
          "actorId": null,
          "actorLink": false,
          "delta": {
            "_id": "RbCCqngEl6EST2GA",
            "system": {},
            "items": [],
            "effects": [],
            "flags": {},
            "name": null,
            "type": null,
            "img": null,
            "ownership": null
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
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 2,
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
      ],
      "lights": [
        {
          "_id": "CklyB5sLeOrqXnyD",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "rotation": 0,
          "walls": true,
          "vision": false,
          "config": {
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
          "hidden": false,
          "flags": {}
        },
        {
          "_id": "FtC2kpWU6Y3S6wJT",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "rotation": 0,
          "walls": true,
          "vision": false,
          "config": {
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
          "hidden": false,
          "flags": {}
        }
      ],
      "notes": [
        {
          "_id": "tRvHXdBSGlfd6egB",
          "entryId": null,
          "pageId": null,
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 0,
          "texture": {
            "src": "icons/svg/book.svg",
            "anchorX": 0.5,
            "anchorY": 0.5,
            "offsetX": 0,
            "offsetY": 0,
            "fit": "contain",
            "scaleX": 1,
            "scaleY": 1,
            "rotation": 0,
            "tint": "#ffffff",
            "alphaThreshold": 0
          },
          "iconSize": 40,
          "fontFamily": "Signika",
          "fontSize": 32,
          "textAnchor": 1,
          "textColor": "#ffffff",
          "global": false,
          "flags": {}
        },
        {
          "_id": "DEnQNeArk9YJJAV1",
          "entryId": null,
          "pageId": null,
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 1,
          "texture": {
            "src": "icons/svg/book.svg",
            "anchorX": 0.5,
            "anchorY": 0.5,
            "offsetX": 0,
            "offsetY": 0,
            "fit": "contain",
            "scaleX": 1,
            "scaleY": 1,
            "rotation": 0,
            "tint": "#ffffff",
            "alphaThreshold": 0
          },
          "iconSize": 40,
          "fontFamily": "Signika",
          "fontSize": 32,
          "textAnchor": 1,
          "textColor": "#ffffff",
          "global": false,
          "flags": {}
        }
      ],
      "sounds": [
        {
          "_id": "g9pHEMy9gBHQht4C",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "radius": 0,
          "path": null,
          "repeat": false,
          "volume": 0.5,
          "walls": true,
          "easing": true,
          "hidden": false,
          "darkness": {
            "min": 0,
            "max": 1
          },
          "effects": {
            "base": {
              "intensity": 5
            },
            "muffled": {
              "intensity": 5
            }
          },
          "flags": {}
        },
        {
          "_id": "pss5T9PLXgbQ6y7y",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "radius": 0,
          "path": null,
          "repeat": false,
          "volume": 0.5,
          "walls": true,
          "easing": true,
          "hidden": false,
          "darkness": {
            "min": 0,
            "max": 1
          },
          "effects": {
            "base": {
              "intensity": 5
            },
            "muffled": {
              "intensity": 5
            }
          },
          "flags": {}
        }
      ],
      "regions": [],
      "templates": [
        {
          "_id": "zJcVWgT5pyfWTaw8",
          "author": "r6bXhB7k9cXa3cif",
          "t": "circle",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 0,
          "distance": 0,
          "direction": 0,
          "angle": 0,
          "width": 0,
          "borderColor": "#000000",
          "fillColor": "#cc2829",
          "texture": null,
          "hidden": false,
          "flags": {}
        },
        {
          "_id": "mi0HgNMBHlwdmnO3",
          "author": "r6bXhB7k9cXa3cif",
          "t": "circle",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 1,
          "distance": 0,
          "direction": 0,
          "angle": 0,
          "width": 0,
          "borderColor": "#000000",
          "fillColor": "#cc2829",
          "texture": null,
          "hidden": false,
          "flags": {}
        }
      ],
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
        "createdTime": 1776821602279,
        "modifiedTime": 1776824906859,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
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
      "name": "test-scene-updated",
      "width": 1000,
      "_id": "pxk3rKsNgpwB6bG5",
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
        "createdTime": 1776821577383,
        "modifiedTime": 1776821602303,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
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
      "name": "test-scene-updated",
      "width": 1000,
      "_id": "vCLJOw0STWGphIWU",
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
          "enabled": true,
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
          "name": "test",
          "shape": 4,
          "_id": "TNEsC0gFgz1Nc9pg",
          "displayName": 0,
          "actorId": null,
          "actorLink": false,
          "delta": {
            "_id": "781VTntlraAbO7C3",
            "system": {},
            "items": [],
            "effects": [],
            "flags": {},
            "name": null,
            "type": null,
            "img": null,
            "ownership": null
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
          "x": 300,
          "y": 400,
          "elevation": 0,
          "sort": 0,
          "locked": false,
          "lockRotation": false,
          "rotation": 323.13010235415595,
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
        },
        {
          "name": "test",
          "shape": 4,
          "_id": "Yi9FWJpDX5TPvXK3",
          "displayName": 0,
          "actorId": null,
          "actorLink": false,
          "delta": {
            "_id": "8Wwm8EDsy6y1C9aG",
            "system": {},
            "items": [],
            "effects": [],
            "flags": {},
            "name": null,
            "type": null,
            "img": null,
            "ownership": null
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
          "x": 400,
          "y": 300,
          "elevation": 0,
          "sort": 1,
          "locked": false,
          "lockRotation": false,
          "rotation": 306.86989764584405,
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
            "effects": 0,
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
        },
        {
          "actorLink": true,
          "alpha": 1,
          "bar1": {
            "attribute": "attributes.hp"
          },
          "bar2": {
            "attribute": "attributes.ac.value"
          },
          "detectionModes": [],
          "displayBars": 40,
          "displayName": 30,
          "disposition": 1,
          "flags": {},
          "height": 1,
          "light": {
            "alpha": 1,
            "angle": 360,
            "animation": {
              "intensity": 5,
              "reverse": false,
              "speed": 5,
              "type": null
            },
            "attenuation": 0.5,
            "bright": 0,
            "color": null,
            "coloration": 1,
            "contrast": 0,
            "darkness": {
              "max": 1,
              "min": 0
            },
            "dim": 0,
            "luminosity": 0.5,
            "negative": false,
            "priority": 0,
            "saturation": 0,
            "shadows": 0
          },
          "lockRotation": false,
          "movementAction": null,
          "name": "Perrin",
          "occludable": {
            "radius": 0
          },
          "ring": {
            "colors": {
              "background": null,
              "ring": null
            },
            "effects": 1,
            "enabled": false,
            "subject": {
              "scale": 1,
              "texture": null
            }
          },
          "rotation": 0,
          "sight": {
            "angle": 360,
            "attenuation": 0.1,
            "brightness": 0,
            "color": null,
            "contrast": 0,
            "enabled": true,
            "range": 5,
            "saturation": 0,
            "visionMode": "basic"
          },
          "texture": {
            "alphaThreshold": 0.75,
            "anchorX": 0.5,
            "anchorY": 0.5,
            "fit": "contain",
            "offsetX": 0,
            "offsetY": 0,
            "rotation": 0,
            "scaleX": 0.8,
            "scaleY": 0.8,
            "src": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
            "tint": "#ffffff"
          },
          "turnMarker": {
            "animation": null,
            "disposition": false,
            "mode": 1,
            "src": null
          },
          "width": 1,
          "actorId": "xSst5kAigAZw6wDr",
          "hidden": false,
          "sort": 2,
          "shape": 4,
          "_id": "KV0Ua6gg79bZ456e",
          "delta": {
            "_id": "KV0Ua6gg79bZ456e",
            "system": {},
            "flags": {},
            "name": null,
            "type": null,
            "img": null,
            "items": [],
            "effects": [],
            "ownership": null
          },
          "x": 300,
          "y": 500,
          "elevation": 0,
          "locked": false,
          "_movementHistory": [],
          "_regions": []
        }
      ],
      "lights": [
        {
          "_id": "xhnUaDHq5geiaMpN",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "rotation": 0,
          "walls": true,
          "vision": false,
          "config": {
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
          "hidden": false,
          "flags": {}
        },
        {
          "_id": "alXpkL3wPgo5ilnE",
          "x": 400,
          "y": 400,
          "elevation": 0,
          "rotation": 0,
          "walls": true,
          "vision": false,
          "config": {
            "negative": false,
            "priority": 0,
            "alpha": 0.5,
            "angle": 360,
            "bright": 5,
            "color": "#ff0000",
            "coloration": 1,
            "dim": 5,
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
          "hidden": false,
          "flags": {}
        }
      ],
      "notes": [
        {
          "_id": "u4cgxPOUlNFG8nuo",
          "entryId": null,
          "pageId": null,
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 0,
          "texture": {
            "src": "icons/svg/book.svg",
            "anchorX": 0.5,
            "anchorY": 0.5,
            "offsetX": 0,
            "offsetY": 0,
            "fit": "contain",
            "scaleX": 1,
            "scaleY": 1,
            "rotation": 0,
            "tint": "#ffffff",
            "alphaThreshold": 0
          },
          "iconSize": 40,
          "fontFamily": "Signika",
          "fontSize": 32,
          "textAnchor": 1,
          "textColor": "#ffffff",
          "global": false,
          "flags": {}
        },
        {
          "_id": "wlwDtURPMkAE0FzQ",
          "entryId": null,
          "pageId": null,
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 1,
          "texture": {
            "src": "icons/svg/book.svg",
            "anchorX": 0.5,
            "anchorY": 0.5,
            "offsetX": 0,
            "offsetY": 0,
            "fit": "contain",
            "scaleX": 1,
            "scaleY": 1,
            "rotation": 0,
            "tint": "#ffffff",
            "alphaThreshold": 0
          },
          "iconSize": 40,
          "fontFamily": "Signika",
          "fontSize": 32,
          "textAnchor": 1,
          "textColor": "#ffffff",
          "global": false,
          "flags": {}
        }
      ],
      "sounds": [
        {
          "_id": "shSUwiJBdTy6CeeY",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "radius": 0,
          "path": null,
          "repeat": false,
          "volume": 0.5,
          "walls": true,
          "easing": true,
          "hidden": false,
          "darkness": {
            "min": 0,
            "max": 1
          },
          "effects": {
            "base": {
              "intensity": 5
            },
            "muffled": {
              "intensity": 5
            }
          },
          "flags": {}
        },
        {
          "_id": "eVHjE3FcVX13yc31",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "radius": 0,
          "path": null,
          "repeat": false,
          "volume": 0.5,
          "walls": true,
          "easing": true,
          "hidden": false,
          "darkness": {
            "min": 0,
            "max": 1
          },
          "effects": {
            "base": {
              "intensity": 5
            },
            "muffled": {
              "intensity": 5
            }
          },
          "flags": {}
        }
      ],
      "regions": [],
      "templates": [
        {
          "_id": "G1UnpWqP0sPpMJ5W",
          "author": "r6bXhB7k9cXa3cif",
          "t": "circle",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 0,
          "distance": 0,
          "direction": 0,
          "angle": 0,
          "width": 0,
          "borderColor": "#000000",
          "fillColor": "#cc2829",
          "texture": null,
          "hidden": false,
          "flags": {}
        },
        {
          "_id": "u3ECicH2osxWg3EL",
          "author": "r6bXhB7k9cXa3cif",
          "t": "circle",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "sort": 1,
          "distance": 0,
          "direction": 0,
          "angle": 0,
          "width": 0,
          "borderColor": "#000000",
          "fillColor": "#cc2829",
          "texture": null,
          "hidden": false,
          "flags": {}
        }
      ],
      "tiles": [],
      "walls": [
        {
          "light": 20,
          "sight": 20,
          "sound": 20,
          "move": 20,
          "c": [
            650,
            338,
            413,
            613
          ],
          "_id": "LDpSlXWcYlP7hudR",
          "dir": 0,
          "door": 0,
          "ds": 0,
          "threshold": {
            "light": null,
            "sight": null,
            "sound": null,
            "attenuation": false
          },
          "animation": null,
          "flags": {}
        }
      ],
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
        "createdTime": 1776824911885,
        "modifiedTime": 1777908899344,
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
      "_id": "p5lAI5vNiA5ncl9B",
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
        "createdTime": 1777909249364,
        "modifiedTime": 1777909249364,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    },
    {
      "height": 500,
      "name": "test-scene-expendable",
      "width": 500,
      "_id": "HnqqUK4yMWaerJPM",
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
        "createdTime": 1777909249371,
        "modifiedTime": 1777909249371,
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
  "requestId": "create-scene_1777909249361",
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
    "_id": "p5lAI5vNiA5ncl9B",
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
      "createdTime": 1777909249364,
      "modifiedTime": 1777909249364,
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
      "sceneId": "p5lAI5vNiA5ncl9B",
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
  -d '{"sceneId":"p5lAI5vNiA5ncl9B","data":{"name":"test-scene-updated"}}'
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
      "sceneId": "p5lAI5vNiA5ncl9B",
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
        "sceneId": "p5lAI5vNiA5ncl9B",
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
  🔤{"sceneId":"p5lAI5vNiA5ncl9B","data":{"name":"test-scene-updated"}}🔤 ➡️ body

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
  "requestId": "update-scene_1777909249383",
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
    "_id": "p5lAI5vNiA5ncl9B",
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
      "createdTime": 1777909249364,
      "modifiedTime": 1777909249384,
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
  sceneId: 'HnqqUK4yMWaerJPM'
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
curl -X DELETE 'http://localhost:3010/scene?clientId=fvtt_099ad17ea199e7e3&sceneId=HnqqUK4yMWaerJPM' \
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
    'sceneId': 'HnqqUK4yMWaerJPM'
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
    sceneId: 'HnqqUK4yMWaerJPM'
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
  🔤sceneId=HnqqUK4yMWaerJPM🔤 ➡️ sceneId
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
  "requestId": "delete-scene_1777909254408",
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
      "sceneId": "p5lAI5vNiA5ncl9B"
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
  -d '{"sceneId":"p5lAI5vNiA5ncl9B"}'
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
      "sceneId": "p5lAI5vNiA5ncl9B"
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
        "sceneId": "p5lAI5vNiA5ncl9B"
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
  🔤{"sceneId":"p5lAI5vNiA5ncl9B"}🔤 ➡️ body

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
  "requestId": "switch-scene_1777909249389",
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
    "_id": "p5lAI5vNiA5ncl9B",
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
      "createdTime": 1777909249364,
      "modifiedTime": 1777909249390,
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

