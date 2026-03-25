---
tag: scene
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# scene

## GET /scene

Get scene(s) Retrieves one or more scenes by ID, name, active status, viewed status, or all.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
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
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
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
curl -X GET 'http://localhost:3010/scene?clientId=foundry-testing-r6bXhB7k9cXa3cif&all=true' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/scene'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif',
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
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
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
  "requestId": "get-scene_1774367584375",
  "data": [
    {
      "_id": "NUEDEFAULTSCENE0",
      "name": "Foundry Virtual Tabletop",
      "navigation": true,
      "navOrder": 0,
      "navName": "",
      "background": {
        "src": "nue/defaultscene/fvtt-background.webp",
        "scaleX": 1,
        "scaleY": 1,
        "offsetX": 0,
        "offsetY": 0,
        "rotation": 0,
        "tint": "#ffffff",
        "anchorX": 0,
        "anchorY": 0,
        "fit": "fill",
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
        "color": "#000000",
        "alpha": 0.2,
        "distance": 1,
        "units": "",
        "style": "solidLines",
        "thickness": 1
      },
      "tokenVision": false,
      "drawings": [
        {
          "shape": {
            "height": 100,
            "type": "r",
            "width": 100,
            "radius": null,
            "points": []
          },
          "x": 100,
          "y": 100,
          "_id": "YXXBrFbROdLYoYvm",
          "author": "r6bXhB7k9cXa3cif",
          "elevation": 0,
          "sort": 0,
          "rotation": 0,
          "bezierFactor": 0,
          "fillType": 0,
          "fillColor": "#cc2829",
          "fillAlpha": 0.5,
          "strokeWidth": 8,
          "strokeColor": "#cc2829",
          "strokeAlpha": 1,
          "texture": null,
          "fontFamily": "Signika",
          "fontSize": 48,
          "textColor": "#ffffff",
          "textAlpha": 1,
          "hidden": false,
          "locked": false,
          "interface": false,
          "flags": {}
        }
      ],
      "tokens": [
        {
          "name": "actor",
          "displayName": 0,
          "actorLink": true,
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
          "actorId": "sj7XQgYVZv8SpEQG",
          "hidden": false,
          "sort": 3,
          "shape": 4,
          "_id": "ZROaJlBucaNukYyI",
          "delta": {
            "_id": "ZROaJlBucaNukYyI",
            "system": {},
            "flags": {},
            "name": null,
            "type": null,
            "img": null,
            "items": [],
            "effects": [],
            "ownership": null
          },
          "x": 1390,
          "y": 418,
          "elevation": 0,
          "locked": false,
          "_movementHistory": [],
          "_regions": []
        }
      ],
      "lights": [
        {
          "_id": "d22Cax8HDPMG4F6I",
          "x": 656,
          "y": 1473,
          "rotation": 0,
          "walls": true,
          "vision": false,
          "config": {
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
            },
            "negative": false,
            "priority": 0
          },
          "hidden": false,
          "flags": {},
          "elevation": 0
        },
        {
          "x": 1826,
          "y": 1891,
          "rotation": 0,
          "walls": true,
          "vision": false,
          "config": {
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
            },
            "negative": false,
            "priority": 0
          },
          "hidden": false,
          "flags": {},
          "_id": "eGuMjw01vEYimWVX",
          "elevation": 0
        },
        {
          "x": 3057,
          "y": 1439,
          "rotation": 0,
          "walls": true,
          "vision": false,
          "config": {
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
            },
            "negative": false,
            "priority": 0
          },
          "hidden": false,
          "flags": {},
          "_id": "TCET4ZNPkl5oZukY",
          "elevation": 0
        },
        {
          "_id": "cOpD0Q4AuCGiKRCb",
          "x": 2824,
          "y": 772,
          "rotation": 0,
          "walls": true,
          "vision": false,
          "config": {
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
            },
            "negative": false,
            "priority": 0
          },
          "hidden": false,
          "flags": {},
          "elevation": 0
        },
        {
          "_id": "adhkydxURYamgKKF",
          "x": 2822,
          "y": 777,
          "rotation": 179,
          "walls": true,
          "vision": false,
          "config": {
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
            },
            "negative": false,
            "priority": 0
          },
          "hidden": false,
          "flags": {},
          "elevation": 0
        }
      ],
      "notes": [],
      "sounds": [],
      "templates": [],
      "tiles": [
        {
          "texture": {
            "src": "nue/defaultscene/fvtt-logo.webp",
            "scaleX": 1,
            "scaleY": 1,
            "offsetX": 0,
            "offsetY": 0,
            "rotation": 0,
            "tint": "#ffffff",
            "anchorX": 0.5,
            "anchorY": 0.5,
            "fit": "fill",
            "alphaThreshold": 0.75
          },
          "x": 1520,
          "y": 480,
          "width": 800,
          "height": 800,
          "_id": "mMxIUI1fXJmrR1zK",
          "rotation": 0,
          "alpha": 1,
          "hidden": false,
          "locked": false,
          "occlusion": {
            "mode": 1,
            "alpha": 0
          },
          "video": {
            "loop": true,
            "autoplay": true,
            "volume": 0
          },
          "flags": {},
          "sort": 100,
          "restrictions": {
            "light": false,
            "weather": false
          },
          "elevation": 0
        }
      ],
      "walls": [],
      "playlist": null,
      "playlistSound": null,
      "journal": null,
      "journalEntryPage": null,
      "weather": "",
      "flags": {},
      "active": false,
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
        "globalLight": {
          "enabled": false,
          "darkness": {
            "max": 1,
            "min": 0
          },
          "alpha": 0.5,
          "bright": false,
          "color": null,
          "coloration": 1,
          "luminosity": 0,
          "saturation": 0,
          "contrast": 0,
          "shadows": 0
        },
        "darknessLevel": 0,
        "darknessLock": false,
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
      "regions": [],
      "folder": null,
      "sort": 0,
      "ownership": {
        "default": 0,
        "5ypAoBvOiyjDKiaZ": 3
      },
      "_stats": {
        "compendiumSource": null,
        "duplicateSource": null,
        "exportSource": null,
        "coreVersion": "13.348",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1763765287462,
        "modifiedTime": 1774161338704,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    },
    {
      "name": "test-scene-expendable",
      "width": 500,
      "height": 500,
      "_id": "zA0IVJWezqEadUaA",
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
        "createdTime": 1774161335815,
        "modifiedTime": 1774161335815,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    },
    {
      "name": "test-scene-updated",
      "width": 1000,
      "height": 1000,
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
      "_id": "lQBzUavCuBjGxRW8",
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
        "createdTime": 1774161335040,
        "modifiedTime": 1774365414522,
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
      "name": "test-scene",
      "width": 1000,
      "_id": "u2dOm1Uzbx9CT9jn",
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
        "createdTime": 1774367584350,
        "modifiedTime": 1774367584350,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    },
    {
      "height": 500,
      "name": "test-scene-expendable",
      "width": 500,
      "_id": "jUAhtzPlUbAb5Jfh",
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
        "createdTime": 1774367584363,
        "modifiedTime": 1774367584363,
        "lastModifiedBy": "r6bXhB7k9cXa3cif"
      }
    }
  ]
}
```


---

## POST /scene

Create a new scene

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| data | object | ✓ | body | Scene data object (name, width, height, grid, etc.) |
| clientId | string |  | body, query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Created scene data

### Try It Out

<ApiTester
  method="POST"
  path="/scene"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/scene';
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
curl -X POST 'http://localhost:3010/scene?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
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
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif'
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
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
  "requestId": "create-scene_1774367584346",
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
    "_id": "u2dOm1Uzbx9CT9jn",
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
      "createdTime": 1774367584350,
      "modifiedTime": 1774367584350,
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
|------|------|----------|--------|--------------|
| data | object | ✓ | body | Object containing the scene fields to update |
| clientId | string |  | body, query | Client ID for the Foundry world |
| sceneId | string |  | query, body | ID of the scene to update |
| name | string |  | query, body | Name of the scene to update (alternative to sceneId) |
| active | boolean |  | query, body | Set to true to target the active scene |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Updated scene data

### Try It Out

<ApiTester
  method="PUT"
  path="/scene"
  parameters={[{"name":"data","type":"object","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"sceneId","type":"string","required":false,"source":"query"},{"name":"name","type":"string","required":false,"source":"query"},{"name":"active","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/scene';
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
      "sceneId": "u2dOm1Uzbx9CT9jn",
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
curl -X PUT 'http://localhost:3010/scene?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"sceneId":"u2dOm1Uzbx9CT9jn","data":{"name":"test-scene-updated"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/scene'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif'
}
url = f'{base_url}{path}'

response = requests.put(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "sceneId": "u2dOm1Uzbx9CT9jn",
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
        "sceneId": "u2dOm1Uzbx9CT9jn",
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"sceneId":"u2dOm1Uzbx9CT9jn","data":{"name":"test-scene-updated"}}🔤 ➡️ body

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
  "requestId": "update-scene_1774367584390",
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
    "_id": "u2dOm1Uzbx9CT9jn",
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
      "createdTime": 1774367584350,
      "modifiedTime": 1774367584392,
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
|------|------|----------|--------|--------------|
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
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
  sceneId: 'jUAhtzPlUbAb5Jfh'
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
curl -X DELETE 'http://localhost:3010/scene?clientId=foundry-testing-r6bXhB7k9cXa3cif&sceneId=jUAhtzPlUbAb5Jfh' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/scene'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif',
    'sceneId': 'jUAhtzPlUbAb5Jfh'
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
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
    sceneId: 'jUAhtzPlUbAb5Jfh'
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤sceneId=jUAhtzPlUbAb5Jfh🔤 ➡️ sceneId
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
  "requestId": "delete-scene_1774367589416",
  "success": true
}
```


---

## POST /switch-scene

Switch the active scene

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string |  | body, query | Client ID for the Foundry world |
| sceneId | string |  | body, query | ID of the scene to activate |
| name | string |  | body, query | Name of the scene to activate (alternative to sceneId) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the scene switch

### Try It Out

<ApiTester
  method="POST"
  path="/switch-scene"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"sceneId","type":"string","required":false,"source":"body"},{"name":"name","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/switch-scene';
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
      "sceneId": "u2dOm1Uzbx9CT9jn"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/switch-scene?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"sceneId":"u2dOm1Uzbx9CT9jn"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/switch-scene'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "sceneId": "u2dOm1Uzbx9CT9jn"
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
        "sceneId": "u2dOm1Uzbx9CT9jn"
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"sceneId":"u2dOm1Uzbx9CT9jn"}🔤 ➡️ body

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
  "requestId": "switch-scene_1774367584400",
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
    "_id": "u2dOm1Uzbx9CT9jn",
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
      "createdTime": 1774367584350,
      "modifiedTime": 1774367584402,
      "lastModifiedBy": "r6bXhB7k9cXa3cif"
    }
  }
}
```


