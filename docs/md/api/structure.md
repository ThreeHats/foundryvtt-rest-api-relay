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
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

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
  "requestId": "structure_1773999623981",
  "clientId": "your-client-id",
  "type": "structure-result",
  "data": {
    "folders": {
      "test-folder": {
        "id": "cytOvtbwpRX336gV",
        "uuid": "Folder.cytOvtbwpRX336gV",
        "type": "Scene"
      }
    },
    "entities": {
      "scenes": [
        {
          "_id": "NUEDEFAULTSCENE0",
          "name": "Foundry Virtual Tabletop",
          "active": false,
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
            },
            {
              "_id": "D72uPb59A6kqlM1o",
              "name": "Aboleth",
              "displayName": 0,
              "actorId": "2z0ZFTO8hYIWeEiL",
              "actorLink": false,
              "delta": {
                "_id": "8d0K6twLsEcqq9Py",
                "system": {},
                "items": [],
                "effects": [],
                "flags": {}
              },
              "width": 2,
              "height": 2,
              "texture": {
                "src": "systems/dnd5e/tokens/aberration/Aboleth.webp",
                "anchorX": 0.5,
                "anchorY": 0.5,
                "offsetX": 0,
                "offsetY": 0,
                "fit": "contain",
                "scaleX": 2,
                "scaleY": 2,
                "rotation": 0,
                "tint": "#ffffff",
                "alphaThreshold": 0.75
              },
              "shape": 4,
              "x": 1664,
              "y": 358,
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
                  "scale": 2,
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
            },
            {
              "_id": "lP1YxLyDcmAGmj4I",
              "name": "Aboleth",
              "displayName": 0,
              "actorId": "2z0ZFTO8hYIWeEiL",
              "actorLink": false,
              "delta": {
                "_id": "j4y2pzRQibl3DXGF",
                "system": {},
                "items": [],
                "effects": [],
                "flags": {}
              },
              "width": 2,
              "height": 2,
              "texture": {
                "src": "systems/dnd5e/tokens/aberration/Aboleth.webp",
                "anchorX": 0.5,
                "anchorY": 0.5,
                "offsetX": 0,
                "offsetY": 0,
                "fit": "contain",
                "scaleX": 2,
                "scaleY": 2,
                "rotation": 0,
                "tint": "#ffffff",
                "alphaThreshold": 0.75
              },
              "shape": 4,
              "x": 2077,
              "y": 557,
              "elevation": 0,
              "sort": 2,
              "locked": false,
              "lockRotation": false,
              "rotation": 306.0850730428521,
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
                  "scale": 2,
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
            "modifiedTime": 1773973282350,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
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
        },
        {
          "_id": "dgTzUeYHz3ofwqC0",
          "name": "test-scene-updated",
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
          "width": 1000,
          "height": 1000,
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
          "tokens": [
            {
              "_id": "XVxua916sjJpOpho",
              "name": "Unknown",
              "displayName": 0,
              "actorId": "c17ZnzQuaTgBEUO8",
              "actorLink": false,
              "delta": {
                "_id": "xMGIe3RylbAJZaIX",
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
              "shape": 4,
              "x": 500,
              "y": 500,
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
              "movementAction": "walk",
              "_movementHistory": [],
              "_regions": [],
              "flags": {}
            }
          ],
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1773973276903,
            "modifiedTime": 1773999592018,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "2mlLTd0S2pYR5qbW",
          "name": "test-scene-updated",
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
          "foregroundElevation": 20,
          "thumb": null,
          "width": 1000,
          "height": 1000,
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
          "tokens": [
            {
              "_id": "7Nbvl6vN27DDqxxh",
              "name": "Updated Test Actor",
              "displayName": 0,
              "actorId": "VKu2l9IdAzxaXrOo",
              "actorLink": false,
              "delta": {
                "_id": "30e64hH7BTAXDpMt",
                "type": "character",
                "system": {},
                "items": [
                  {
                    "_id": "q4tr1vTU8RxtU1UZ",
                    "name": "Priest",
                    "type": "background",
                    "img": "icons/sundries/documents/document-torn-diagram-tan.webp",
                    "system": {
                      "description": {
                        "value": "<ul><li><strong>Skill Proficiencies:</strong> Insight, Religion</li><li><strong>Languages:</strong> Two of your choice</li><li><strong>Equipment:</strong> A holy symbol, 5 sticks of incense, prayer book, vestments, a set of common clothes, and a pouch containing 15 gp.</li></ul>",
                        "chat": ""
                      },
                      "identifier": "priest",
                      "source": {
                        "book": "",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "startingEquipment": [],
                      "advancement": []
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": null,
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "O3ZjSw0GtAOPScHI",
                    "name": "Lightfoot Halfling",
                    "type": "race",
                    "img": "icons/equipment/feet/shoes-leather-simple-brown.webp",
                    "system": {
                      "description": {
                        "value": "<p>Your halfling character has a number of traits in common with all other halflings.</p><p><em><strong>Ability Score Increase.</strong></em> Your Dexterity score increases by 2.</p><p><em><strong>Age.</strong></em> A halfling reaches adulthood at the age of 20 and generally lives into the middle of his or her second century.</p><p><em><strong>Alignment.</strong></em> Most halflings are lawful good. As a rule, they are good-hearted and kind, hate to see others in pain, and have no tolerance for oppression. They are also very orderly and traditional, leaning heavily on the support of their community and the comfort of their old ways.</p><p><em><strong>Size.</strong></em> Halflings average about 3 feet tall and weigh about 40 pounds. Your size is Small.</p><p><em><strong>Speed.</strong></em> Your base walking speed is 25 feet.</p><p><em><strong>Lucky.</strong></em> When you roll a 1 on the d20 for an attack roll, ability check, or saving throw, you can reroll the die and must use the new roll.</p><p><em><strong>Brave.</strong></em> You have advantage on saving throws against being frightened.</p><p><em><strong>Halfling Nimbleness.</strong></em> You can move through the space of any creature that is of a size larger than yours.</p><p><em><strong>Languages.</strong></em> You can speak, read, and write Common and Halfling. The Halfling language isn't secret, but halflings are loath to share it with others. They write very little, so they don't have a rich body of literature. Their oral tradition, however, is very strong. Almost all halflings speak Common to converse with the people in whose lands they dwell or through which they are traveling.</p><h5>Lightfoot</h5><p>As a lightfoot halfling, you can easily hide from notice, even using other people as cover. You're inclined to be affable and get along well with others.</p><p>Lightfoots are more prone to wanderlust than other halflings, and often dwell alongside other races or take up a nomadic life.</p><p><em><strong>Ability Score Increase.</strong></em> Your Charisma score increases by 1.</p><p><em><strong>Naturally Stealthy.</strong></em> You can attempt to hide even when you are obscured only by a creature that is at least one size larger than you.</p>",
                        "chat": ""
                      },
                      "identifier": "lightfoot-halfling",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "advancement": [
                        {
                          "_id": "nInhIgkbzzJTdm8F",
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.races.LOMdcNAGWh5xpfm4",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.races.7Yoo9hG0hfPSmBoC",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.races.PqxZgcJzp1VVgP8t",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.races.GWPjKFeIthBBeCFJ",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {
                            "added": {
                              "FtOM4QiOW5MwgcS3": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
                              "nmmihiqphHjoE8dl": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC",
                              "cWrETHzCRs1Ueqd3": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
                              "AArhiOrSkaQUnCZS": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ"
                            }
                          },
                          "level": 0,
                          "title": ""
                        },
                        {
                          "_id": "Z9hvZFkWUNvowbQX",
                          "type": "AbilityScoreImprovement",
                          "configuration": {
                            "points": 0,
                            "fixed": {
                              "str": 0,
                              "dex": 2,
                              "con": 0,
                              "int": 0,
                              "wis": 0,
                              "cha": 1
                            },
                            "cap": 2,
                            "locked": []
                          },
                          "value": {
                            "type": "asi",
                            "assignments": {
                              "dex": 2,
                              "cha": 1
                            }
                          },
                          "level": 0,
                          "title": ""
                        },
                        {
                          "_id": "hv2bcANK5jEJZaAb",
                          "type": "Size",
                          "configuration": {
                            "sizes": [
                              "sm"
                            ]
                          },
                          "value": {
                            "size": "sm"
                          },
                          "level": 1,
                          "title": "",
                          "hint": "Halflings average about 3 feet tall and weigh about 40 pounds. Your size is Small."
                        },
                        {
                          "_id": "nGwMjsfNU6CXHk3A",
                          "type": "Trait",
                          "configuration": {
                            "mode": "default",
                            "allowReplacements": false,
                            "grants": [
                              "languages:standard:common",
                              "languages:standard:halfling"
                            ],
                            "choices": []
                          },
                          "level": 0,
                          "title": "",
                          "value": {
                            "chosen": [
                              "languages:standard:common",
                              "languages:standard:halfling"
                            ]
                          }
                        }
                      ],
                      "movement": {
                        "burrow": null,
                        "climb": null,
                        "fly": null,
                        "swim": null,
                        "walk": 25,
                        "units": "ft",
                        "hover": false
                      },
                      "senses": {
                        "darkvision": null,
                        "blindsight": null,
                        "tremorsense": null,
                        "truesight": null,
                        "units": "ft",
                        "special": ""
                      },
                      "type": {
                        "value": "humanoid",
                        "subtype": "halfling",
                        "custom": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.races.Item.ZgYBjYYfiUstQD6f",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "FtOM4QiOW5MwgcS3",
                    "name": "Lucky",
                    "type": "feat",
                    "img": "icons/sundries/gaming/dice-runed-brown.webp",
                    "system": {
                      "activities": [],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>When you roll a 1 on the d20 for an attack roll, ability check, or saving throw, you can reroll the die and must use the new roll.</p><section class=\"secret foundry-note\" id=\"secret-S04TPyvUh05Dz0Ng\"><p><strong>Foundry Note</strong></p><p>This property can be enabled on your character sheet in the Special Traits configuration on the Attributes tab.</p></section>",
                        "chat": ""
                      },
                      "identifier": "lucky",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "advancement": [],
                      "crewed": false,
                      "enchant": {},
                      "prerequisites": {
                        "level": null,
                        "repeatable": false
                      },
                      "properties": [],
                      "requirements": "Halfling",
                      "type": {
                        "value": "race",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "sourceId": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        }
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "nmmihiqphHjoE8dl",
                    "name": "Brave",
                    "type": "feat",
                    "img": "icons/skills/melee/unarmed-punch-fist.webp",
                    "system": {
                      "activities": [],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>You have advantage on saving throws against being frightened.</p>",
                        "chat": ""
                      },
                      "identifier": "brave",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "advancement": [],
                      "crewed": false,
                      "enchant": {},
                      "prerequisites": {
                        "level": null,
                        "repeatable": false
                      },
                      "properties": [],
                      "requirements": "Halfling",
                      "type": {
                        "value": "race",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "sourceId": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC",
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        }
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "cWrETHzCRs1Ueqd3",
                    "name": "Halfling Nimbleness",
                    "type": "feat",
                    "img": "icons/skills/movement/feet-winged-boots-brown.webp",
                    "system": {
                      "activities": [],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>You can move through the space of any creature that is of a size larger than yours.</p>",
                        "chat": ""
                      },
                      "identifier": "halfling-nimbleness",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "advancement": [],
                      "crewed": false,
                      "enchant": {},
                      "prerequisites": {
                        "level": null,
                        "repeatable": false
                      },
                      "properties": [],
                      "requirements": "Halfling",
                      "type": {
                        "value": "race",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "sourceId": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        }
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "AArhiOrSkaQUnCZS",
                    "name": "Naturally Stealthy",
                    "type": "feat",
                    "img": "icons/magic/perception/silhouette-stealth-shadow.webp",
                    "system": {
                      "activities": [],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>You can attempt to hide even when you are obscured only by a creature that is at least one size larger than you.</p>",
                        "chat": ""
                      },
                      "identifier": "naturally-stealthy",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "advancement": [],
                      "crewed": false,
                      "enchant": {},
                      "prerequisites": {
                        "level": null,
                        "repeatable": false
                      },
                      "properties": [],
                      "requirements": "Lightfoot Halfling",
                      "type": {
                        "value": "race",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "sourceId": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ",
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        }
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "8Grf7ga6JcZF0X6x",
                    "name": "Monk",
                    "type": "class",
                    "img": "icons/skills/melee/hand-grip-staff-blue.webp",
                    "system": {
                      "description": {
                        "value": "<p>As a monk, you gain the following class features.</p><h3>Hit Points</h3><p><strong>Hit Dice:</strong> 1d8 per monk level<br /><strong>Hit Points at 1st Level:</strong> 8 + your Constitution modifier<br /><strong>Hit Points at Higher Levels:</strong> 1d8 (or 5) + your Constitution modifier per monk level after 1st</p><h3>Proficiencies</h3><p><strong>Armor:</strong> None<br /><strong>Weapons:</strong> Simple weapons, shortswords<br /><strong>Tools:</strong> Choose one type of artisan's tools or one musical instrument<br /><strong>Saving Throws:</strong> Strength, Dexterity<br /><strong>Skills:</strong> Choose two from Acrobatics, Athletics, History, Insight, Religion, and Stealth</p><h1>Monk Advancement</h1><table><thead><tr><td>Level</td><td>Proficiency Bonus</td><td>Martial Arts</td><td>Ki Points</td><td>Unarmored Movement</td><td>Features</td></tr></thead><tbody><tr><td>1st</td><td>+2</td><td>1d4</td><td>—</td><td>—</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI]{Unarmored Defense}, @UUID[Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw]{Martial Arts}</td></tr><tr><td>2nd</td><td>+2</td><td>1d4</td><td>2</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.10b6z2W1txNkrGP7]{Ki}, @UUID[Compendium.dnd5e.classfeatures.Item.zCeqyQ8uIPNdYJSW]{Unarmored Movement}</td></tr><tr><td>3rd</td><td>+2</td><td>1d4</td><td>3</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.rtpQdX77dYWbDIOH]{Monastic Tradition}, @UUID[Compendium.dnd5e.classfeatures.Item.mzweVbnsJPQiVkAe]{Deflect Missiles}</td></tr><tr><td>4th</td><td>+2</td><td>1d4</td><td>4</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}, @UUID[Compendium.dnd5e.classfeatures.Item.KQz9bqxVkXjDl8gK]{Slow Fall}</td></tr><tr><td>5th</td><td>+3</td><td>1d6</td><td>5</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.XogoBnFWmCAHXppo]{Extra Attack}, @UUID[Compendium.dnd5e.classfeatures.Item.pvRc6GAu1ok6zihC]{Stunning Strike}</td></tr><tr><td>6th</td><td>+3</td><td>1d6</td><td>6</td><td>+15 ft.</td><td><p>@UUID[Compendium.dnd5e.classfeatures.Item.7flZKruSSu6dHg6D]{Ki-Empowered Strikes},</p><p>Monastic Tradition feature</p></td></tr><tr><td>7th</td><td>+3</td><td>1d6</td><td>7</td><td>+15 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.a4P4DNMmH8CqSNkC]{Evasion}, @UUID[Compendium.dnd5e.classfeatures.Item.ZmC31XKS4YNENnoc]{Stillness of Mind}</td></tr><tr><td>8th</td><td>+3</td><td>1d6</td><td>8</td><td>+15 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>9th</td><td>+4</td><td>1d6</td><td>9</td><td>+15 ft.</td><td>Unarmored Movement improvement</td></tr><tr><td>10th</td><td>+4</td><td>1d6</td><td>10</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.bqWA7t9pDELbNRkp]{Purity of Body}</td></tr><tr><td>11th</td><td>+4</td><td>1d8</td><td>11</td><td>+20 ft.</td><td>Monastic Tradition feature</td></tr><tr><td>12th</td><td>+4</td><td>1d8</td><td>12</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>13th</td><td>+5</td><td>1d8</td><td>13</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.XjuGBeB8Y0C3A5D4]{Tongue of the Sun and Moon}</td></tr><tr><td>14th</td><td>+5</td><td>1d8</td><td>14</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.7D2EkLdISwShEDlN]{Diamond Soul}</td></tr><tr><td>15th</td><td>+5</td><td>1d8</td><td>15</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.gDH8PMrKvLHaNmEI]{Timeless Body}</td></tr><tr><td>16th</td><td>+5</td><td>1d8</td><td>16</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>17th</td><td>+6</td><td>1d10</td><td>17</td><td>+25 ft.</td><td>Monastic Tradition feature</td></tr><tr><td>18th</td><td>+6</td><td>1d10</td><td>18</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.3jwFt3hSqDswBlOH]{Empty Body}</td></tr><tr><td>19th</td><td>+6</td><td>1d10</td><td>19</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>20th</td><td>+6</td><td>1d10</td><td>20</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.mQNPg89YIs7g5tG4]{Perfect Self}</td></tr></tbody></table><h1>Monastic Traditions</h1><p>Three traditions of monastic pursuit are common in the monasteries scattered across the multiverse. Most monasteries practice one tradition exclusively, but a few honor the three traditions and instruct each monk according to his or her aptitude and interest. All three traditions rely on the same basic techniques, diverging as the student grows more adept. Thus, a monk need choose a tradition only upon reaching 3rd level.</p><p>@UUID[Compendium.dnd5e.subclasses.Item.IvlpKMXX3PmW1NY2]{Way of the Open Hand}</p>",
                        "chat": ""
                      },
                      "identifier": "monk",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "startingEquipment": [
                        {
                          "_id": "5yj0P4r9teJDnDtd",
                          "group": "",
                          "sort": 100000,
                          "type": "OR",
                          "requiresProficiency": false
                        },
                        {
                          "_id": "R5tuRtaPonfjQCVU",
                          "group": "5yj0P4r9teJDnDtd",
                          "sort": 200000,
                          "type": "linked",
                          "count": null,
                          "key": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
                          "requiresProficiency": false
                        },
                        {
                          "_id": "Mlf6kel8ws6xgDER",
                          "group": "5yj0P4r9teJDnDtd",
                          "sort": 300000,
                          "type": "weapon",
                          "count": null,
                          "key": "simpleM",
                          "requiresProficiency": false
                        },
                        {
                          "_id": "3TbVLmLPtjVaSh5O",
                          "group": "",
                          "sort": 400000,
                          "type": "OR",
                          "requiresProficiency": false
                        },
                        {
                          "_id": "AvDYtl0uvQsDuhnb",
                          "group": "3TbVLmLPtjVaSh5O",
                          "sort": 500000,
                          "type": "linked",
                          "count": null,
                          "key": "Compendium.dnd5e.items.Item.XY8b594Dn7plACLL",
                          "requiresProficiency": false
                        },
                        {
                          "_id": "4QKQURCmIurbTAzp",
                          "group": "3TbVLmLPtjVaSh5O",
                          "sort": 600000,
                          "type": "linked",
                          "count": null,
                          "key": "Compendium.dnd5e.items.Item.8KWz5DJbWUpNWniP",
                          "requiresProficiency": false
                        },
                        {
                          "_id": "AOYuulsULvsHbSLO",
                          "group": "",
                          "sort": 700000,
                          "type": "linked",
                          "count": 10,
                          "key": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
                          "requiresProficiency": false
                        }
                      ],
                      "wealth": "5d4",
                      "levels": 1,
                      "primaryAbility": {
                        "value": [],
                        "all": true
                      },
                      "hd": {
                        "additional": 0,
                        "denomination": "d8",
                        "spent": 0
                      },
                      "advancement": [
                        {
                          "type": "HitPoints",
                          "configuration": {},
                          "value": {
                            "1": "max"
                          },
                          "title": "Hit Points",
                          "icon": "systems/dnd5e/icons/svg/hit-points.svg",
                          "_id": "ocxNtDFJ7YDaYaK7"
                        },
                        {
                          "_id": "mmAxx3U7FvXNAcKc",
                          "type": "Trait",
                          "configuration": {
                            "mode": "default",
                            "allowReplacements": false,
                            "grants": [
                              "weapon:sim",
                              "weapon:mar:shortsword"
                            ],
                            "choices": []
                          },
                          "level": 1,
                          "title": "",
                          "value": {
                            "chosen": [
                              "weapon:sim",
                              "weapon:mar:shortsword"
                            ]
                          }
                        },
                        {
                          "_id": "QPXy59CQGY9HB0c3",
                          "type": "Trait",
                          "configuration": {
                            "mode": "default",
                            "allowReplacements": false,
                            "grants": [],
                            "choices": [
                              {
                                "count": 1,
                                "pool": [
                                  "tool:art:*",
                                  "tool:music:*"
                                ]
                              }
                            ]
                          },
                          "level": 1,
                          "title": "",
                          "classRestriction": "primary",
                          "value": {
                            "chosen": [
                              "tool:art:brewer"
                            ]
                          }
                        },
                        {
                          "_id": "4M8MQ1E64zbcRg6B",
                          "type": "Trait",
                          "configuration": {
                            "mode": "default",
                            "allowReplacements": false,
                            "grants": [
                              "saves:str",
                              "saves:dex"
                            ],
                            "choices": []
                          },
                          "level": 1,
                          "title": "",
                          "classRestriction": "primary",
                          "value": {
                            "chosen": [
                              "saves:str",
                              "saves:dex"
                            ]
                          }
                        },
                        {
                          "_id": "7HRRCPk80Ng2Evdx",
                          "type": "Trait",
                          "configuration": {
                            "mode": "default",
                            "allowReplacements": false,
                            "grants": [],
                            "choices": [
                              {
                                "count": 2,
                                "pool": [
                                  "skills:acr",
                                  "skills:ath",
                                  "skills:his",
                                  "skills:ins",
                                  "skills:rel",
                                  "skills:ste"
                                ]
                              }
                            ]
                          },
                          "level": 1,
                          "title": "",
                          "classRestriction": "primary",
                          "value": {
                            "chosen": [
                              "skills:acr",
                              "skills:ath"
                            ]
                          }
                        },
                        {
                          "_id": "BQWHr3mt5flvkfIj",
                          "type": "Trait",
                          "configuration": {
                            "mode": "default",
                            "allowReplacements": false,
                            "grants": [
                              "di:poison",
                              "ci:diseased",
                              "ci:poisoned"
                            ],
                            "choices": []
                          },
                          "level": 10,
                          "title": "Purity of Body",
                          "value": {
                            "chosen": []
                          },
                          "hint": "Your mastery of the ki flowing through you makes you immune to disease and poison."
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.UAvV7N7T4zJhxdfI",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.l50hjTxO2r0iecKw",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.eGxoNmSMWKNzChCO",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {
                            "added": {
                              "CwgoTDXWCD7PknIN": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
                              "pchnXqd5C79fVlxy": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw",
                              "RiURabP4FDYMeuWx": "Compendium.dnd5e.classfeatures.Item.eGxoNmSMWKNzChCO"
                            }
                          },
                          "level": 1,
                          "title": "Features",
                          "_id": "n0q8XyiGA3vLPgpK"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.10b6z2W1txNkrGP7",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.5MwNlVZK7m6VolOH",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.TDglPcxIVEzvVSgK",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.yrSFIGTaQOH2PFRI",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.zCeqyQ8uIPNdYJSW",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 2,
                          "title": "Features",
                          "_id": "7TyDqpGGi3r3nsp0"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.rtpQdX77dYWbDIOH",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.mzweVbnsJPQiVkAe",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 3,
                          "title": "Features",
                          "_id": "2sLHTw6k15DSW8WB"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.KQz9bqxVkXjDl8gK",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 4,
                          "title": "Features",
                          "_id": "Zc1jOZK1b9mIKekq"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.XogoBnFWmCAHXppo",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.pvRc6GAu1ok6zihC",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 5,
                          "title": "Features",
                          "_id": "j9LeWmxlsENKaMLo"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.7flZKruSSu6dHg6D",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 6,
                          "title": "Features",
                          "_id": "psobDjMqtA2216Db"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.a4P4DNMmH8CqSNkC",
                                "optional": false
                              },
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.ZmC31XKS4YNENnoc",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 7,
                          "title": "Features",
                          "_id": "K38aFaEMxMqRB0BC"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.bqWA7t9pDELbNRkp",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 10,
                          "title": "Features",
                          "_id": "eLqmJotmwzlGNrxG"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.XjuGBeB8Y0C3A5D4",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 13,
                          "title": "Features",
                          "_id": "N0geIQiuofqYgswj"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.7D2EkLdISwShEDlN",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 14,
                          "title": "Features",
                          "_id": "N1hjizyI82UPp8UI"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.gDH8PMrKvLHaNmEI",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 15,
                          "title": "Features",
                          "_id": "TcLZS9WzC7bPETSd"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.3jwFt3hSqDswBlOH",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 18,
                          "title": "Features",
                          "_id": "tRb3a0tA5IpehPs8"
                        },
                        {
                          "type": "ItemGrant",
                          "configuration": {
                            "items": [
                              {
                                "uuid": "Compendium.dnd5e.classfeatures.mQNPg89YIs7g5tG4",
                                "optional": false
                              }
                            ],
                            "optional": false,
                            "spell": {
                              "ability": [],
                              "preparation": "",
                              "uses": {
                                "max": "",
                                "per": "",
                                "requireSlot": false
                              }
                            }
                          },
                          "value": {},
                          "level": 20,
                          "title": "Features",
                          "_id": "sEQz9c9XhWYjS9x5"
                        },
                        {
                          "type": "ScaleValue",
                          "configuration": {
                            "identifier": "die",
                            "type": "dice",
                            "distance": {
                              "units": ""
                            },
                            "scale": {
                              "1": {
                                "number": null,
                                "faces": 4,
                                "modifiers": []
                              },
                              "5": {
                                "number": null,
                                "faces": 6,
                                "modifiers": []
                              },
                              "11": {
                                "number": null,
                                "faces": 8,
                                "modifiers": []
                              },
                              "17": {
                                "number": null,
                                "faces": 10,
                                "modifiers": []
                              }
                            }
                          },
                          "value": {},
                          "title": "Martial Arts Die",
                          "_id": "MXFbf0nxMiyLdPbX"
                        },
                        {
                          "type": "ScaleValue",
                          "configuration": {
                            "identifier": "unarmored-movement",
                            "type": "distance",
                            "distance": {
                              "units": "ft"
                            },
                            "scale": {
                              "2": {
                                "value": 10
                              },
                              "6": {
                                "value": 15
                              },
                              "10": {
                                "value": 20
                              },
                              "14": {
                                "value": 25
                              },
                              "18": {
                                "value": 30
                              }
                            }
                          },
                          "value": {},
                          "title": "Unarmored Movement",
                          "_id": "1OzfWDWCquoHMeX5"
                        },
                        {
                          "type": "AbilityScoreImprovement",
                          "configuration": {
                            "points": 2,
                            "fixed": {
                              "str": 0,
                              "dex": 0,
                              "con": 0,
                              "int": 0,
                              "wis": 0,
                              "cha": 0
                            },
                            "cap": 2,
                            "locked": []
                          },
                          "value": {
                            "type": "asi"
                          },
                          "level": 4,
                          "title": "Ability Score Improvement",
                          "_id": "ofNSUhSHKhhDuPSR"
                        },
                        {
                          "type": "AbilityScoreImprovement",
                          "configuration": {
                            "points": 2,
                            "fixed": {
                              "str": 0,
                              "dex": 0,
                              "con": 0,
                              "int": 0,
                              "wis": 0,
                              "cha": 0
                            },
                            "cap": 2,
                            "locked": []
                          },
                          "value": {
                            "type": "asi"
                          },
                          "level": 8,
                          "title": "Ability Score Improvement",
                          "_id": "s3t9o57hP6iUHirr"
                        },
                        {
                          "type": "AbilityScoreImprovement",
                          "configuration": {
                            "points": 2,
                            "fixed": {
                              "str": 0,
                              "dex": 0,
                              "con": 0,
                              "int": 0,
                              "wis": 0,
                              "cha": 0
                            },
                            "cap": 2,
                            "locked": []
                          },
                          "value": {
                            "type": "asi"
                          },
                          "level": 12,
                          "title": "Ability Score Improvement",
                          "_id": "O24MWOKc1ImsKaml"
                        },
                        {
                          "type": "AbilityScoreImprovement",
                          "configuration": {
                            "points": 2,
                            "fixed": {
                              "str": 0,
                              "dex": 0,
                              "con": 0,
                              "int": 0,
                              "wis": 0,
                              "cha": 0
                            },
                            "cap": 2,
                            "locked": []
                          },
                          "value": {
                            "type": "asi"
                          },
                          "level": 16,
                          "title": "Ability Score Improvement",
                          "_id": "xdqWoLtgO3uyl3nJ"
                        },
                        {
                          "_id": "puDaUsYrlks0z5gm",
                          "type": "AbilityScoreImprovement",
                          "configuration": {
                            "points": 2,
                            "fixed": {
                              "str": 0,
                              "dex": 0,
                              "con": 0,
                              "int": 0,
                              "wis": 0,
                              "cha": 0
                            },
                            "cap": 2,
                            "locked": []
                          },
                          "value": {
                            "type": "asi"
                          },
                          "level": 19,
                          "title": ""
                        },
                        {
                          "_id": "0awj2yq115ev9u9o",
                          "type": "Subclass",
                          "configuration": {},
                          "value": {
                            "document": null,
                            "uuid": null
                          },
                          "level": 3,
                          "title": "Monastic Tradition"
                        }
                      ],
                      "spellcasting": {
                        "progression": "none",
                        "ability": "",
                        "preparation": {
                          "formula": ""
                        }
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.classes.Item.6VoZrWxhOEKGYhnq",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "CwgoTDXWCD7PknIN",
                    "name": "Unarmored Defense",
                    "type": "feat",
                    "img": "icons/magic/control/silhouette-hold-change-blue.webp",
                    "system": {
                      "activities": [],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>Beginning at 1st Level, while you are wearing no armor and not wielding a Shield, your AC equals 10 + your Dexterity modifier + your Wisdom modifier.</p>",
                        "chat": ""
                      },
                      "identifier": "unarmored-defense",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "advancement": [],
                      "crewed": false,
                      "enchant": {},
                      "prerequisites": {
                        "level": null,
                        "repeatable": false
                      },
                      "properties": [],
                      "requirements": "Monk 1",
                      "type": {
                        "value": "class",
                        "subtype": ""
                      }
                    },
                    "effects": [
                      {
                        "_id": "R5ro4AuNjcdWD56O",
                        "name": "Unarmored Defense",
                        "img": "icons/magic/control/silhouette-hold-change-blue.webp",
                        "type": "base",
                        "system": {},
                        "changes": [
                          {
                            "key": "system.attributes.ac.calc",
                            "value": "unarmoredMonk",
                            "mode": 5,
                            "priority": null
                          }
                        ],
                        "disabled": false,
                        "duration": {
                          "startTime": 0,
                          "seconds": null,
                          "combat": null,
                          "rounds": null,
                          "turns": null,
                          "startRound": null,
                          "startTurn": null
                        },
                        "description": "",
                        "origin": "Item.cOdcNWy4hII029DT",
                        "tint": "#ffffff",
                        "transfer": true,
                        "statuses": [],
                        "sort": 0,
                        "flags": {},
                        "_stats": {
                          "coreVersion": "13.348",
                          "systemId": "dnd5e",
                          "systemVersion": "5.0.4",
                          "lastModifiedBy": null,
                          "compendiumSource": null,
                          "duplicateSource": null,
                          "exportSource": null
                        }
                      }
                    ],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "sourceId": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
                        "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK"
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "pchnXqd5C79fVlxy",
                    "name": "Martial Arts",
                    "type": "feat",
                    "img": "icons/skills/melee/unarmed-punch-fist.webp",
                    "system": {
                      "activities": [],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>At 1st level, your practice of martial arts gives you mastery of combat styles that use and monk weapons, which are shortswords and any simple melee weapons that don't have the two-handed or heavy property. You gain the following benefits while you are unarmed or wielding only monk weapons and you aren't wearing armor or wielding a shield:</p>\n<ul>\n<li>\n<p>You can use Dexterity instead of Strength for the attack and damage rolls of your unarmed strikes and monk weapons.</p>\n</li>\n<li>\n<p>You can roll a d4 in place of the normal damage of your unarmed strike or monk weapon. This die changes as you gain monk levels, as shown in the Martial Arts column of the Monk table.</p>\n</li>\n<li>\n<p>When you use the Attack action with an unarmed strike or a monk weapon on your turn, you can make one unarmed strike as a bonus action. For example, if you take the Attack action and attack with a quarterstaff, you can also make an unarmed strike as a bonus action, assuming you haven't already taken a bonus action this turn.</p>\n</li>\n</ul>\n<p>Certain monasteries use specialized forms of the monk weapons. For example, you might use a club that is two lengths of wood connected by a short chain (called a nunchaku) or a sickle with a shorter, straighter blade (called a kama). Whatever name you use for a monk weapon, you can use the game statistics provided for the weapon.</p>",
                        "chat": ""
                      },
                      "identifier": "martial-arts",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "advancement": [],
                      "crewed": false,
                      "enchant": {},
                      "prerequisites": {
                        "level": null,
                        "repeatable": false
                      },
                      "properties": [],
                      "requirements": "Monk 1",
                      "type": {
                        "value": "class",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "sourceId": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw",
                        "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK"
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "YJ1P3PnFKHOdQpaP",
                    "name": "Hammer",
                    "type": "loot",
                    "img": "icons/tools/hand/hammer-cobbler-steel.webp",
                    "system": {
                      "description": {
                        "value": "<p>A tool with a heavy metal head mounted at the end of its handle, used for jobs such as breaking things and driving in nails. </p>",
                        "chat": ""
                      },
                      "identifier": "hammer",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 3,
                        "units": "lb"
                      },
                      "price": {
                        "value": 1,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "properties": [],
                      "type": {
                        "value": "",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.14pNRT4sZy9rgvhb",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "DDnxRCeYUhXstWU8",
                    "name": "Tinderbox",
                    "type": "loot",
                    "img": "icons/sundries/lights/torch-black.webp",
                    "system": {
                      "description": {
                        "value": "<p>This small container holds flint, fire steel, and tinder (usually dry cloth soaked in light oil) used to kindle a fire. Using it to light a torch - or anything else with abundant, exposed fuel - takes an action. Lighting any other fire takes 1 minute.</p>\n<p> </p>",
                        "chat": ""
                      },
                      "identifier": "tinderbox",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": "8KWz5DJbWUpNWniP",
                      "quantity": 1,
                      "weight": {
                        "value": 1,
                        "units": "lb"
                      },
                      "price": {
                        "value": 5,
                        "denomination": "sp"
                      },
                      "rarity": "",
                      "properties": [],
                      "type": {
                        "value": "",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.1FSubnBpSTDmVaYV",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "5skKSSB4ShHbKoc8",
                    "name": "Waterskin",
                    "type": "consumable",
                    "img": "icons/sundries/survival/wetskin-leather-purple.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "utility",
                          "activation": {
                            "type": "action",
                            "value": 1,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [
                              {
                                "type": "itemUses",
                                "target": "",
                                "value": "1",
                                "scaling": {
                                  "mode": "",
                                  "formula": ""
                                }
                              }
                            ],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "units": "touch",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": []
                          },
                          "sort": 0
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": 4,
                        "recovery": [],
                        "autoDestroy": false
                      },
                      "description": {
                        "value": "<p>A leather hide sewn into an enclosed skin which can contain up to 4 pints of liquid. It weighs 5 pounds when full; a pint of water is approximately 1 pound.</p>",
                        "chat": ""
                      },
                      "identifier": "waterskin",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": "6OYR11aJX2dEVtOj",
                      "quantity": 1,
                      "weight": {
                        "value": 5,
                        "units": "lb"
                      },
                      "price": {
                        "value": 2,
                        "denomination": "sp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "damage": {
                        "base": {
                          "number": null,
                          "denomination": null,
                          "types": [],
                          "custom": {
                            "enabled": false
                          },
                          "scaling": {
                            "number": 1
                          }
                        },
                        "replace": false
                      },
                      "magicalBonus": null,
                      "properties": [],
                      "type": {
                        "value": "food",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.1L5wkmbw0fmNAr38",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "dDuMscUuMI2bTdkj",
                    "name": "Torch",
                    "type": "consumable",
                    "img": "icons/sundries/lights/torch-black.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "attack",
                          "activation": {
                            "type": "action",
                            "value": 1,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [
                              {
                                "type": "itemUses",
                                "target": "",
                                "value": "1",
                                "scaling": {
                                  "mode": "",
                                  "formula": ""
                                }
                              }
                            ],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "1",
                            "units": "hour",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "units": "self",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "radius",
                              "size": "40",
                              "width": "",
                              "height": "",
                              "units": "ft"
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "attack": {
                            "ability": "str",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "value": "melee",
                              "classification": "weapon"
                            }
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": [
                              {
                                "number": null,
                                "denomination": null,
                                "bonus": "",
                                "types": [
                                  "fire"
                                ],
                                "custom": {
                                  "enabled": true,
                                  "formula": "1"
                                },
                                "scaling": {
                                  "mode": "whole",
                                  "number": null,
                                  "formula": ""
                                }
                              }
                            ]
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": []
                          },
                          "sort": 0
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": 1,
                        "recovery": [],
                        "autoDestroy": false
                      },
                      "description": {
                        "value": "<p>A torch burns for 1 hour, providing bright light in a 20-foot radius and dim light for an additional 20 feet. If you make a melee attack with a burning torch and hit, it deals 1 fire damage.</p>",
                        "chat": ""
                      },
                      "identifier": "torch",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": "8KWz5DJbWUpNWniP",
                      "quantity": 10,
                      "weight": {
                        "value": 1,
                        "units": "lb"
                      },
                      "price": {
                        "value": 1,
                        "denomination": "cp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "damage": {
                        "base": {
                          "number": null,
                          "denomination": null,
                          "bonus": "",
                          "types": [
                            "fire"
                          ],
                          "custom": {
                            "enabled": true,
                            "formula": "1"
                          },
                          "scaling": {
                            "mode": "",
                            "number": null,
                            "formula": ""
                          }
                        },
                        "replace": false
                      },
                      "magicalBonus": null,
                      "properties": [],
                      "type": {
                        "value": "trinket",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.29ZLE8PERtFVD3QU",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "nC6VcR5JAIbR4err",
                    "name": "Stick of Incense",
                    "type": "loot",
                    "img": "icons/consumables/grains/breadsticks-crackers-wrapped-ration-brown.webp",
                    "system": {
                      "description": {
                        "value": "<p>When blocks of incense cannot be used or a cheaper alternative is required, people often use these to perfume the air, whether for pleasurable or religious purposes.</p>",
                        "chat": ""
                      },
                      "identifier": "stick-of-incense",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 5,
                      "weight": {
                        "value": 0,
                        "units": "lb"
                      },
                      "price": {
                        "value": 2,
                        "denomination": "sp"
                      },
                      "rarity": "",
                      "properties": [],
                      "type": {
                        "value": "",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.3b0RvGi0TnTYpIxn",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "WeKJI3gPUAU52WAX",
                    "name": "Dart",
                    "type": "weapon",
                    "img": "icons/weapons/thrown/dart-feathered.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "attack",
                          "activation": {
                            "type": "action",
                            "value": 1,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "value": "20",
                            "units": "ft",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "value": "ranged",
                              "classification": "weapon"
                            }
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": []
                          },
                          "sort": 0
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>A small thrown implement crafted with a short wooden shaft and crossed feathres with a sharp wooden or metal tip. Darts can be thrown with sufficient force to puncture the skin.</p>",
                        "chat": ""
                      },
                      "identifier": "dart",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 10,
                      "weight": {
                        "value": 0.25,
                        "units": "lb"
                      },
                      "price": {
                        "value": 5,
                        "denomination": "cp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": true,
                      "cover": null,
                      "crewed": false,
                      "hp": {
                        "conditions": "",
                        "dt": null,
                        "max": 0,
                        "value": 0
                      },
                      "ammunition": {},
                      "armor": {
                        "value": 10
                      },
                      "damage": {
                        "base": {
                          "number": 1,
                          "denomination": 4,
                          "bonus": "",
                          "types": [
                            "piercing"
                          ],
                          "custom": {
                            "enabled": false,
                            "formula": ""
                          },
                          "scaling": {
                            "mode": "",
                            "number": null,
                            "formula": ""
                          }
                        },
                        "versatile": {
                          "number": null,
                          "denomination": null,
                          "bonus": "",
                          "types": [],
                          "custom": {
                            "enabled": false,
                            "formula": ""
                          },
                          "scaling": {
                            "mode": "",
                            "number": null,
                            "formula": ""
                          }
                        }
                      },
                      "magicalBonus": null,
                      "mastery": "",
                      "properties": [
                        "fin",
                        "thr"
                      ],
                      "proficient": null,
                      "range": {
                        "value": 20,
                        "long": 60,
                        "reach": null,
                        "units": "ft"
                      },
                      "type": {
                        "value": "simpleR",
                        "baseItem": "dart"
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "1F73YcUHbZMgePAD",
                    "name": "Common Clothes",
                    "type": "equipment",
                    "img": "icons/equipment/chest/shirt-collared-brown.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "utility",
                          "activation": {
                            "type": "",
                            "value": null,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "units": "self",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": []
                          },
                          "sort": 0
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>Clothes worn by most commoners.</p>",
                        "chat": ""
                      },
                      "identifier": "common-clothes",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 3,
                        "units": "lb"
                      },
                      "price": {
                        "value": 5,
                        "denomination": "sp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "cover": null,
                      "crewed": false,
                      "hp": {
                        "conditions": "",
                        "dt": null,
                        "max": 0,
                        "value": 0
                      },
                      "speed": {
                        "conditions": "",
                        "value": null
                      },
                      "armor": {
                        "value": 0,
                        "magicalBonus": null,
                        "dex": null
                      },
                      "proficient": null,
                      "properties": [],
                      "strength": null,
                      "type": {
                        "value": "clothing",
                        "baseItem": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.8RXjiddJ6VGyE7vB",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "LQhGSEatJ3VK7oqW",
                    "name": "Rations",
                    "type": "consumable",
                    "img": "icons/consumables/grains/bread-loaf-boule-rustic-brown.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "utility",
                          "activation": {
                            "type": "action",
                            "value": 1,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [
                              {
                                "type": "itemUses",
                                "target": "",
                                "value": "1",
                                "scaling": {
                                  "mode": "",
                                  "formula": ""
                                }
                              }
                            ],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "units": "touch",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "1",
                              "type": "creature",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": []
                          },
                          "sort": 0
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": 1,
                        "recovery": [],
                        "autoDestroy": true
                      },
                      "description": {
                        "value": "<p>Rations consist of dry foods suitable for extended travel, including jerky, dried fruit, hardtack, and nuts.</p>",
                        "chat": ""
                      },
                      "identifier": "rations",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": "XY8b594Dn7plACLL",
                      "quantity": 10,
                      "weight": {
                        "value": 2,
                        "units": "lb"
                      },
                      "price": {
                        "value": 5,
                        "denomination": "sp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "damage": {
                        "base": {
                          "number": null,
                          "denomination": null,
                          "types": [],
                          "custom": {
                            "enabled": false
                          },
                          "scaling": {
                            "number": 1
                          }
                        },
                        "replace": false
                      },
                      "magicalBonus": null,
                      "properties": [],
                      "type": {
                        "value": "food",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.8d95YV1jHcxPygJ9",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "1kqMwSwjfErSFNvl",
                    "name": "Pouch",
                    "type": "container",
                    "img": "icons/containers/bags/pouch-rounded-leather-gold-tan.webp",
                    "system": {
                      "description": {
                        "value": "<p>A cloth or leather pouch can hold up to 20 sling bullets or 50 blowgun needles, among other things. A compartmentalized pouch for holding spell components is called a component pouch. A pouch can hold up to ⅕ cubic foot or 6 pounds of gear.</p>",
                        "chat": ""
                      },
                      "identifier": "pouch",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 1,
                        "units": "lb"
                      },
                      "price": {
                        "value": 5,
                        "denomination": "sp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "currency": {
                        "pp": 0,
                        "gp": 0,
                        "ep": 0,
                        "sp": 0,
                        "cp": 0
                      },
                      "capacity": {
                        "volume": {
                          "units": "cubicFoot"
                        },
                        "weight": {
                          "value": 6,
                          "units": "lb"
                        }
                      },
                      "properties": []
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.9bWTRRDym06PzSAf",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "zF5L4xHnJBC7b2iM",
                    "name": "Crowbar",
                    "type": "loot",
                    "img": "icons/tools/hand/pickaxe-steel-white.webp",
                    "system": {
                      "description": {
                        "value": "<p>Using a crowbar grants advantage to Strength checks where the crowbar's leverage can be applied.</p>",
                        "chat": ""
                      },
                      "identifier": "crowbar",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": "XY8b594Dn7plACLL",
                      "quantity": 1,
                      "weight": {
                        "value": 5,
                        "units": "lb"
                      },
                      "price": {
                        "value": 2,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "properties": [],
                      "type": {
                        "value": "",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.AkyQyonZMVcvOrXU",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "g7U3OAXVcoI4lwzf",
                    "name": "Unarmed Strike",
                    "type": "weapon",
                    "img": "icons/skills/melee/unarmed-punch-fist.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "attack",
                          "activation": {
                            "type": "action",
                            "value": 1,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "value": "5",
                            "units": "ft",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "value": "melee",
                              "classification": "unarmed"
                            }
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": [],
                            "max": ""
                          },
                          "sort": 0,
                          "name": ""
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>A  punch, kick, head-butt, or similar forceful blow (none of which count as weapons). On a hit, an unarmed strike deals bludgeoning damage equal to 1 + your Strength modifier. You are proficient with your unarmed strikes.</p>",
                        "chat": ""
                      },
                      "identifier": "unarmed-strike",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 0,
                        "units": "lb"
                      },
                      "price": {
                        "value": 0,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": true,
                      "cover": null,
                      "crewed": false,
                      "hp": {
                        "conditions": "",
                        "dt": null,
                        "max": 0,
                        "value": 0
                      },
                      "ammunition": {},
                      "armor": {
                        "value": 10
                      },
                      "damage": {
                        "base": {
                          "number": null,
                          "denomination": null,
                          "bonus": "",
                          "types": [
                            "bludgeoning"
                          ],
                          "custom": {
                            "enabled": true,
                            "formula": "@scale.monk.die"
                          },
                          "scaling": {
                            "mode": "",
                            "number": null,
                            "formula": ""
                          }
                        },
                        "versatile": {
                          "number": null,
                          "denomination": null,
                          "bonus": "",
                          "types": [],
                          "custom": {
                            "enabled": false,
                            "formula": ""
                          },
                          "scaling": {
                            "mode": "",
                            "number": null,
                            "formula": ""
                          }
                        }
                      },
                      "magicalBonus": null,
                      "mastery": "",
                      "properties": [
                        "fin"
                      ],
                      "proficient": null,
                      "range": {
                        "value": null,
                        "long": null,
                        "reach": 5,
                        "units": "ft"
                      },
                      "type": {
                        "value": "simpleM",
                        "baseItem": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "last": {
                          "dnd5eactivity000": {
                            "attackMode": "oneHanded",
                            "damageType": {
                              "0": "bludgeoning"
                            }
                          }
                        },
                        "riders": {
                          "activity": [],
                          "effect": []
                        }
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.GsuvwoekKZatfKwF",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "WcROGD590imVj9qp",
                    "name": "Backpack",
                    "type": "container",
                    "img": "icons/containers/bags/pack-leather-white-tan.webp",
                    "system": {
                      "description": {
                        "value": "<p>A backpack can hold one cubic foot or 30 pounds of gear. You can also strap items, such as a bedroll or a coil of rope, to the outside of a backpack.</p>",
                        "chat": ""
                      },
                      "identifier": "backpack",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 5,
                        "units": "lb"
                      },
                      "price": {
                        "value": 2,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "currency": {
                        "pp": 0,
                        "gp": 0,
                        "ep": 0,
                        "sp": 0,
                        "cp": 0
                      },
                      "capacity": {
                        "volume": {
                          "units": "cubicFoot"
                        },
                        "weight": {
                          "value": 30,
                          "units": "lb"
                        }
                      },
                      "properties": []
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.H8YCd689ezlD26aT",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "8UQanAvHUIHZXp0O",
                    "name": "Piton",
                    "type": "consumable",
                    "img": "icons/tools/fasteners/nail-steel.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "utility",
                          "activation": {
                            "type": "action",
                            "value": 1,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [
                              {
                                "type": "itemUses",
                                "target": "",
                                "value": "1",
                                "scaling": {
                                  "mode": "",
                                  "formula": ""
                                }
                              }
                            ],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "value": "5",
                            "units": "ft",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": []
                          },
                          "sort": 0
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": 1,
                        "recovery": [],
                        "autoDestroy": false
                      },
                      "description": {
                        "value": "<p>A metal spike that is drive into a seam in a climbing surface with a climbing hammer. It can also be used like iron spikes to block doors in a pinch.</p>",
                        "chat": ""
                      },
                      "identifier": "piton",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": "xsB7Y2WI476kvOt4",
                      "quantity": 10,
                      "weight": {
                        "value": 0.25,
                        "units": "lb"
                      },
                      "price": {
                        "value": 5,
                        "denomination": "cp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "damage": {
                        "base": {
                          "number": null,
                          "denomination": null,
                          "types": [],
                          "custom": {
                            "enabled": false
                          },
                          "scaling": {
                            "number": 1
                          }
                        },
                        "replace": false
                      },
                      "magicalBonus": null,
                      "properties": [],
                      "type": {
                        "value": "trinket",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.P31t6tGgt9aLAdYt",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "irCoxLHje0eA6Qmu",
                    "name": "Hempen Rope (50 ft.)",
                    "type": "consumable",
                    "img": "icons/sundries/survival/rope-wrapped-brown.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "check",
                          "activation": {
                            "type": "action",
                            "value": 1,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [
                              {
                                "type": "itemUses",
                                "target": "",
                                "scaling": {},
                                "value": "1"
                              }
                            ],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "units": "self",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "check": {
                            "ability": "str",
                            "dc": {
                              "calculation": "",
                              "formula": "17"
                            },
                            "associated": []
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": [],
                            "max": ""
                          },
                          "sort": 0,
                          "name": "Burst"
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": 1,
                        "recovery": [],
                        "autoDestroy": false
                      },
                      "description": {
                        "value": "<p>Rope, whether made of hemp or silk, has 2 hit points and can be burst with a DC 17 Strength check.</p>",
                        "chat": ""
                      },
                      "identifier": "hempen-rope-50-ft",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 10,
                        "units": "lb"
                      },
                      "price": {
                        "value": 1,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "damage": {
                        "base": {
                          "number": null,
                          "denomination": null,
                          "types": [],
                          "custom": {
                            "enabled": false
                          },
                          "scaling": {
                            "number": 1
                          }
                        },
                        "replace": false
                      },
                      "magicalBonus": null,
                      "properties": [],
                      "type": {
                        "value": "trinket",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "riders": {
                          "activity": [],
                          "effect": []
                        }
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.QXmaarJ4X8P0C1HV",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "GfKuJYLFfw00oW9R",
                    "name": "Reliquary",
                    "type": "equipment",
                    "img": "icons/containers/chest/chest-reinforced-steel-red.webp",
                    "system": {
                      "activities": [],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>A tiny box or other container holding a fragment of a precious relic, saint, or other historical figure that dedicated their life to walk the path of a true believer. A deity imbues the bearer of this artifact with the ability to call forth their power and in doing so spread the faith once more.</p>\n<p><strong>Spellcasting Focus</strong>. A cleric or paladin can use a holy symbol as a spellcasting focus. To use the symbol in this way, the caster must hold it in hand, wear it visibly, or bear it on a shield.</p>",
                        "chat": ""
                      },
                      "identifier": "reliquary",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 2,
                        "units": "lb"
                      },
                      "price": {
                        "value": 5,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": false,
                      "cover": null,
                      "crewed": false,
                      "hp": {
                        "conditions": "",
                        "dt": null,
                        "max": 0,
                        "value": 0
                      },
                      "speed": {
                        "conditions": "",
                        "value": null
                      },
                      "armor": {
                        "value": 0,
                        "magicalBonus": null,
                        "dex": null
                      },
                      "proficient": null,
                      "properties": [
                        "foc"
                      ],
                      "strength": null,
                      "type": {
                        "value": "trinket",
                        "baseItem": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {
                      "dnd5e": {
                        "riders": {
                          "activity": [],
                          "effect": []
                        }
                      }
                    },
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.gP1URGq3kVIIFHJ7",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "NsNcJBUG5ajbk7sG",
                    "name": "Vestments",
                    "type": "equipment",
                    "img": "icons/equipment/back/mantle-collared-black.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "utility",
                          "activation": {
                            "type": "",
                            "value": null,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "units": "self",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": []
                          },
                          "sort": 0
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>Simple or ostentacious wear, often used by priests and other religious figures for use in rituals and ceremonies.</p>",
                        "chat": ""
                      },
                      "identifier": "vestments",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 4,
                        "units": "lb"
                      },
                      "price": {
                        "value": 1,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": true,
                      "cover": null,
                      "crewed": false,
                      "hp": {
                        "conditions": "",
                        "dt": null,
                        "max": 0,
                        "value": 0
                      },
                      "speed": {
                        "conditions": "",
                        "value": null
                      },
                      "armor": {
                        "value": 0,
                        "magicalBonus": null,
                        "dex": null
                      },
                      "proficient": null,
                      "properties": [],
                      "strength": null,
                      "type": {
                        "value": "clothing",
                        "baseItem": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.irtqrzaUCeshmTZp",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "XbF0CTipyqbMKHsB",
                    "name": "Prayer Book",
                    "type": "loot",
                    "img": "icons/sundries/books/book-purple-cross.webp",
                    "system": {
                      "description": {
                        "value": "<p>A book containing prayers and incantations dedicated to a specific power for the faithful to follow.</p>",
                        "chat": ""
                      },
                      "identifier": "prayer-book",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 5,
                        "units": "lb"
                      },
                      "price": {
                        "value": 25,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "properties": [],
                      "type": {
                        "value": "",
                        "subtype": ""
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.odV5cq2HSLSCH69k",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  },
                  {
                    "_id": "GYyovoNdU66nxOlX",
                    "name": "Shortsword",
                    "type": "weapon",
                    "img": "icons/weapons/swords/sword-guard-worn-purple.webp",
                    "system": {
                      "activities": [
                        {
                          "_id": "dnd5eactivity000",
                          "type": "attack",
                          "activation": {
                            "type": "action",
                            "value": 1,
                            "condition": "",
                            "override": false
                          },
                          "consumption": {
                            "targets": [],
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "value": "",
                            "units": "inst",
                            "special": "",
                            "override": false
                          },
                          "effects": [],
                          "range": {
                            "value": "5",
                            "units": "ft",
                            "special": "",
                            "override": false
                          },
                          "target": {
                            "template": {
                              "count": "",
                              "contiguous": false,
                              "type": "",
                              "size": "",
                              "width": "",
                              "height": "",
                              "units": ""
                            },
                            "affects": {
                              "count": "",
                              "type": "",
                              "choice": false,
                              "special": ""
                            },
                            "prompt": true,
                            "override": false
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "value": "melee",
                              "classification": "weapon"
                            }
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "uses": {
                            "spent": 0,
                            "recovery": []
                          },
                          "sort": 0
                        }
                      ],
                      "uses": {
                        "spent": 0,
                        "max": "",
                        "recovery": []
                      },
                      "description": {
                        "value": "<p>A medium sized blade with a firm crossguard and a leather-wrapped handle. A versatile weapon which makes up in versatility what it lacks in reach.</p>",
                        "chat": ""
                      },
                      "identifier": "shortsword",
                      "source": {
                        "book": "SRD 5.1",
                        "page": "",
                        "custom": "",
                        "license": "CC-BY-4.0",
                        "revision": 1,
                        "rules": "2014"
                      },
                      "identified": true,
                      "unidentified": {
                        "description": ""
                      },
                      "container": null,
                      "quantity": 1,
                      "weight": {
                        "value": 2,
                        "units": "lb"
                      },
                      "price": {
                        "value": 10,
                        "denomination": "gp"
                      },
                      "rarity": "",
                      "attunement": "",
                      "attuned": false,
                      "equipped": true,
                      "cover": null,
                      "crewed": false,
                      "hp": {
                        "conditions": "",
                        "dt": null,
                        "max": 0,
                        "value": 0
                      },
                      "ammunition": {},
                      "armor": {
                        "value": 10
                      },
                      "damage": {
                        "base": {
                          "number": 1,
                          "denomination": 6,
                          "bonus": "",
                          "types": [
                            "piercing"
                          ],
                          "custom": {
                            "enabled": false,
                            "formula": ""
                          },
                          "scaling": {
                            "mode": "",
                            "number": null,
                            "formula": ""
                          }
                        },
                        "versatile": {
                          "number": null,
                          "denomination": null,
                          "bonus": "",
                          "types": [],
                          "custom": {
                            "enabled": false,
                            "formula": ""
                          },
                          "scaling": {
                            "mode": "",
                            "number": null,
                            "formula": ""
                          }
                        }
                      },
                      "magicalBonus": null,
                      "properties": [
                        "fin",
                        "lgt"
                      ],
                      "proficient": null,
                      "range": {
                        "value": null,
                        "long": null,
                        "reach": 5,
                        "units": "ft"
                      },
                      "type": {
                        "value": "martialM",
                        "baseItem": "shortsword"
                      }
                    },
                    "effects": [],
                    "folder": null,
                    "sort": 0,
                    "ownership": {
                      "default": 0
                    },
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "lastModifiedBy": null,
                      "compendiumSource": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  }
                ],
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
              "shape": 4,
              "x": 500,
              "y": 500,
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
              "movementAction": "walk",
              "_movementHistory": [],
              "_regions": [],
              "flags": {}
            }
          ],
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1773999590697,
            "modifiedTime": 1773999592018,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
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
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

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
  "requestId": "get-folder_1773999624120",
  "clientId": "your-client-id",
  "type": "get-folder-result",
  "data": {
    "id": "cytOvtbwpRX336gV",
    "uuid": "Folder.cytOvtbwpRX336gV",
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
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

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
  "requestId": "create-folder_1773999623769",
  "clientId": "your-client-id",
  "type": "create-folder-result",
  "data": {
    "id": "cytOvtbwpRX336gV",
    "uuid": "Folder.cytOvtbwpRX336gV",
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
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

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
  folderId: 'cytOvtbwpRX336gV'
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
curl -X DELETE 'http://localhost:3010/delete-folder?clientId=your-client-id&folderId=cytOvtbwpRX336gV' \
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
    'folderId': 'cytOvtbwpRX336gV'
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
    folderId: 'cytOvtbwpRX336gV'
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
  🔤folderId=cytOvtbwpRX336gV🔤 ➡️ folderId
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
  "requestId": "delete-folder_1773999624258",
  "clientId": "your-client-id",
  "type": "delete-folder-result",
  "data": {
    "deleted": true,
    "folderId": "cytOvtbwpRX336gV",
    "entitiesDeleted": 0,
    "foldersDeleted": 1
  }
}
```


