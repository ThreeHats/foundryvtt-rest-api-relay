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

**Required scope:** `scene:read`

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
const baseUrl = 'http://localhost:3011';
const path = '/scene';
const params = {
  clientId: 'qsl-integration-test',
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
curl -X GET 'http://localhost:3011/scene?clientId=qsl-integration-test&all=true' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/scene'
params = {
    'clientId': 'qsl-integration-test',
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
  const baseUrl = 'http://localhost:3011';
  const path = '/scene';
  const params = {
    clientId: 'qsl-integration-test',
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
  3011 ➡️ port
  🔤/scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤all=true🔤 ➡️ all
  🔤?🧲clientId🧲&🧲all🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "get-scene_1782956912577",
  "data": [
    {
      "_id": "NUEDEFAULTSCENE0",
      "name": "Foundry Virtual Tabletop",
      "navigation": true,
      "navOrder": 0,
      "thumb": "nue/defaultscene/thumb.webp",
      "width": 3840,
      "height": 1920,
      "padding": 0,
      "shiftX": 0,
      "shiftY": 0,
      "initial": {
        "x": null,
        "y": null,
        "scale": null
      },
      "initialLevel": "defaultLevel0000",
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
        "mode": 0,
        "colors": {
          "explored": null,
          "unexplored": null
        },
        "reset": 1660769143211
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
      "transition": {
        "type": null,
        "duration": 1500,
        "activeOnly": false
      },
      "drawings": [],
      "tokens": [
        {
          "name": "test",
          "shape": 4,
          "level": "defaultLevel0000",
          "_id": "3PipH5UTtbVhES7K",
          "displayName": 0,
          "actorId": null,
          "actorLink": false,
          "delta": null,
          "x": 0,
          "y": 0,
          "elevation": 0,
          "width": 1,
          "height": 1,
          "depth": 1,
          "texture": {
            "src": "icons/svg/mystery-man.svg",
            "anchorX": 0.5,
            "anchorY": 0.5,
            "fit": "contain",
            "scaleX": 1,
            "scaleY": 1,
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
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
          "detectionModes": {},
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
      "levels": [
        {
          "_id": "defaultLevel0000",
          "name": "Foundry Virtual Tabletop",
          "elevation": {
            "bottom": 0,
            "top": 4
          },
          "background": {
            "color": "#25070d",
            "src": "nue/defaultscene/fvtt-background.webp",
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
          "foreground": {
            "src": null,
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
          "fog": {
            "src": null,
            "tint": "#ffffff"
          },
          "textures": {
            "anchorX": 0.5,
            "anchorY": 0.5,
            "offsetX": 0,
            "offsetY": 0,
            "fit": "fill",
            "scaleX": 1,
            "scaleY": 1,
            "rotation": 0
          },
          "visibility": {
            "levels": []
          },
          "sort": 0,
          "flags": {}
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
          "elevation": 0,
          "levels": [],
          "locked": false
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
          "elevation": 0,
          "levels": [],
          "locked": false
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
          "elevation": 0,
          "levels": [],
          "locked": false
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
          "elevation": 0,
          "levels": [],
          "locked": false
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
          "elevation": 0,
          "levels": [],
          "locked": false
        },
        {
          "name": "test",
          "_id": "mVCPXMBaIWRJAg6u",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "levels": [],
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
          "locked": false,
          "flags": {}
        }
      ],
      "notes": [
        {
          "_id": "Xra5DWrNscar1diu",
          "author": "cpvaGKk3hgoBCzCS",
          "entryId": null,
          "pageId": null,
          "x": 0,
          "y": 0,
          "elevation": 0,
          "levels": [],
          "sort": 0,
          "locked": false,
          "texture": {
            "src": "icons/svg/book.svg",
            "anchorX": 0.5,
            "anchorY": 0.5,
            "fit": "contain",
            "scaleX": 1,
            "scaleY": 1,
            "tint": "#ffffff",
            "alphaThreshold": 0
          },
          "iconSize": 40,
          "fontFamily": "",
          "fontSize": 32,
          "textAnchor": 1,
          "textColor": "#ffffff",
          "global": false,
          "flags": {}
        }
      ],
      "sounds": [
        {
          "name": "test",
          "_id": "8k8o4JbBVJ3n9epl",
          "x": 0,
          "y": 0,
          "elevation": 0,
          "levels": [],
          "radius": 0,
          "path": null,
          "repeat": false,
          "volume": 0.5,
          "walls": true,
          "easing": true,
          "hidden": false,
          "locked": false,
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
      "regions": [
        {
          "_id": "8ofAU5wNntboYU8l",
          "name": "Circle Template [Gamemaster]",
          "color": "#161068",
          "shapes": [
            {
              "type": "circle",
              "x": 0,
              "y": 0,
              "radius": 0,
              "gridBased": false,
              "hole": false
            }
          ],
          "elevation": {
            "bottom": 0,
            "top": null,
            "topInclusive": false
          },
          "levels": [],
          "restriction": {
            "enabled": false,
            "type": "move",
            "priority": 0
          },
          "attachment": {
            "token": null
          },
          "behaviors": [],
          "visibility": 2,
          "highlightMode": "coverage",
          "displayMeasurements": true,
          "hidden": false,
          "locked": false,
          "ownership": {
            "default": 0,
            "cpvaGKk3hgoBCzCS": 3
          },
          "flags": {
            "core": {
              "MeasuredTemplate": true
            }
          },
          "_shapeConstraints": null
        }
      ],
      "tiles": [
        {
          "texture": {
            "src": "nue/defaultscene/fvtt-logo.webp",
            "scaleX": 1,
            "scaleY": 1,
            "tint": "#ffffff",
            "anchorX": 0,
            "anchorY": 0,
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
            "alpha": 0,
            "modes": []
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
          "elevation": 0,
          "levels": []
        }
      ],
      "walls": [],
      "playlist": null,
      "playlistSound": null,
      "journal": null,
      "journalEntryPage": null,
      "weather": "",
      "folder": null,
      "flags": {},
      "navName": "",
      "_stats": {
        "coreVersion": "14.363",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1782956902919,
        "modifiedTime": 1782956902920,
        "lastModifiedBy": "cpvaGKk3hgoBCzCS",
        "compendiumSource": null,
        "duplicateSource": null,
        "exportSource": null
      },
      "active": true,
      "sort": 0,
      "ownership": {
        "default": 0,
        "cpvaGKk3hgoBCzCS": 3
      }
    },
    {
      "name": "test",
      "_id": "UTjLwXV2Zl26XPRl",
      "active": false,
      "navigation": true,
      "navOrder": 0,
      "thumb": null,
      "width": 4000,
      "height": 3000,
      "padding": 0.25,
      "shiftX": 0,
      "shiftY": 0,
      "initial": {
        "x": null,
        "y": null,
        "scale": null
      },
      "initialLevel": "defaultLevel0000",
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
        "mode": 1,
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
      "transition": {
        "type": null,
        "duration": 1500,
        "activeOnly": false
      },
      "drawings": [],
      "tokens": [],
      "levels": [
        {
          "_id": "defaultLevel0000",
          "name": "Level",
          "elevation": {
            "bottom": 0,
            "top": 20
          },
          "background": {
            "color": "#999999",
            "src": null,
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
          "foreground": {
            "src": null,
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
          "fog": {
            "src": null,
            "tint": "#ffffff"
          },
          "textures": {
            "anchorX": 0.5,
            "anchorY": 0.5,
            "offsetX": 0,
            "offsetY": 0,
            "fit": "fill",
            "scaleX": 1,
            "scaleY": 1,
            "rotation": 0
          },
          "visibility": {
            "levels": []
          },
          "sort": 0,
          "flags": {}
        }
      ],
      "lights": [],
      "notes": [],
      "sounds": [],
      "regions": [],
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
        "cpvaGKk3hgoBCzCS": 3
      },
      "flags": {},
      "_stats": {
        "coreVersion": "14.363",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1782956906426,
        "modifiedTime": 1782956906427,
        "lastModifiedBy": "cpvaGKk3hgoBCzCS",
        "compendiumSource": null,
        "duplicateSource": null,
        "exportSource": null
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
      "_id": "SM7AhDv5JgZh6IvK",
      "active": false,
      "navigation": true,
      "navOrder": 0,
      "thumb": null,
      "padding": 0.25,
      "shiftX": 0,
      "shiftY": 0,
      "initial": {
        "x": null,
        "y": null,
        "scale": null
      },
      "initialLevel": "defaultLevel0000",
      "tokenVision": true,
      "fog": {
        "mode": 1,
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
      "transition": {
        "type": null,
        "duration": 1500,
        "activeOnly": false
      },
      "drawings": [],
      "tokens": [],
      "levels": [
        {
          "_id": "defaultLevel0000",
          "name": "Level",
          "elevation": {
            "bottom": 0,
            "top": 20
          },
          "background": {
            "color": "#999999",
            "src": null,
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
          "foreground": {
            "src": null,
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
          "fog": {
            "src": null,
            "tint": "#ffffff"
          },
          "textures": {
            "anchorX": 0.5,
            "anchorY": 0.5,
            "offsetX": 0,
            "offsetY": 0,
            "fit": "fill",
            "scaleX": 1,
            "scaleY": 1,
            "rotation": 0
          },
          "visibility": {
            "levels": []
          },
          "sort": 0,
          "flags": {}
        }
      ],
      "lights": [],
      "notes": [],
      "sounds": [],
      "regions": [],
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
        "cpvaGKk3hgoBCzCS": 3
      },
      "flags": {},
      "_stats": {
        "coreVersion": "14.363",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1782956912564,
        "modifiedTime": 1782956912564,
        "lastModifiedBy": "cpvaGKk3hgoBCzCS",
        "compendiumSource": null,
        "duplicateSource": null,
        "exportSource": null
      }
    },
    {
      "height": 500,
      "name": "test-scene-expendable",
      "width": 500,
      "_id": "Hk2XrjCAU9pHZ7Gc",
      "active": false,
      "navigation": true,
      "navOrder": 0,
      "thumb": null,
      "padding": 0.25,
      "shiftX": 0,
      "shiftY": 0,
      "initial": {
        "x": null,
        "y": null,
        "scale": null
      },
      "initialLevel": "defaultLevel0000",
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
        "mode": 1,
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
      "transition": {
        "type": null,
        "duration": 1500,
        "activeOnly": false
      },
      "drawings": [],
      "tokens": [],
      "levels": [
        {
          "_id": "defaultLevel0000",
          "name": "Level",
          "elevation": {
            "bottom": 0,
            "top": 20
          },
          "background": {
            "color": "#999999",
            "src": null,
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
          "foreground": {
            "src": null,
            "tint": "#ffffff",
            "alphaThreshold": 0.75
          },
          "fog": {
            "src": null,
            "tint": "#ffffff"
          },
          "textures": {
            "anchorX": 0.5,
            "anchorY": 0.5,
            "offsetX": 0,
            "offsetY": 0,
            "fit": "fill",
            "scaleX": 1,
            "scaleY": 1,
            "rotation": 0
          },
          "visibility": {
            "levels": []
          },
          "sort": 0,
          "flags": {}
        }
      ],
      "lights": [],
      "notes": [],
      "sounds": [],
      "regions": [],
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
        "cpvaGKk3hgoBCzCS": 3
      },
      "flags": {},
      "_stats": {
        "coreVersion": "14.363",
        "systemId": "dnd5e",
        "systemVersion": "5.0.4",
        "createdTime": 1782956912572,
        "modifiedTime": 1782956912573,
        "lastModifiedBy": "cpvaGKk3hgoBCzCS",
        "compendiumSource": null,
        "duplicateSource": null,
        "exportSource": null
      }
    }
  ]
}
```


---

## GET /scene/image/raw

Get the raw background image of a scene

Returns the scene's background image file without any tokens, lights, or other canvas elements rendered on it.

**Required scope:** `scene:read`

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

**Required scope:** `scene:write`

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
const baseUrl = 'http://localhost:3011';
const path = '/scene';
const params = {
  clientId: 'qsl-integration-test'
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
curl -X POST 'http://localhost:3011/scene?clientId=qsl-integration-test' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"data":{"name":"test-scene","width":1000,"height":1000,"grid":{"size":100}}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/scene'
params = {
    'clientId': 'qsl-integration-test'
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
  const baseUrl = 'http://localhost:3011';
  const path = '/scene';
  const params = {
    clientId: 'qsl-integration-test'
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
  3011 ➡️ port
  🔤/scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"data":{"name":"test-scene","width":1000,"height":1000,"grid":{"size":100}}}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 77❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "create-scene_1782956912561",
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
    "_id": "SM7AhDv5JgZh6IvK",
    "active": false,
    "navigation": true,
    "navOrder": 0,
    "thumb": null,
    "padding": 0.25,
    "shiftX": 0,
    "shiftY": 0,
    "initial": {
      "x": null,
      "y": null,
      "scale": null
    },
    "initialLevel": "defaultLevel0000",
    "tokenVision": true,
    "fog": {
      "mode": 1,
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
    "transition": {
      "type": null,
      "duration": 1500,
      "activeOnly": false
    },
    "drawings": [],
    "tokens": [],
    "levels": [
      {
        "_id": "defaultLevel0000",
        "name": "Level",
        "elevation": {
          "bottom": 0,
          "top": 20
        },
        "background": {
          "color": "#999999",
          "src": null,
          "tint": "#ffffff",
          "alphaThreshold": 0.75
        },
        "foreground": {
          "src": null,
          "tint": "#ffffff",
          "alphaThreshold": 0.75
        },
        "fog": {
          "src": null,
          "tint": "#ffffff"
        },
        "textures": {
          "anchorX": 0.5,
          "anchorY": 0.5,
          "offsetX": 0,
          "offsetY": 0,
          "fit": "fill",
          "scaleX": 1,
          "scaleY": 1,
          "rotation": 0
        },
        "visibility": {
          "levels": []
        },
        "sort": 0,
        "flags": {}
      }
    ],
    "lights": [],
    "notes": [],
    "sounds": [],
    "regions": [],
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
      "cpvaGKk3hgoBCzCS": 3
    },
    "flags": {},
    "_stats": {
      "coreVersion": "14.363",
      "systemId": "dnd5e",
      "systemVersion": "5.0.4",
      "createdTime": 1782956912564,
      "modifiedTime": 1782956912564,
      "lastModifiedBy": "cpvaGKk3hgoBCzCS",
      "compendiumSource": null,
      "duplicateSource": null,
      "exportSource": null
    }
  }
}
```


---

## PUT /scene

Update an existing scene

**Required scope:** `scene:write`

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
const baseUrl = 'http://localhost:3011';
const path = '/scene';
const params = {
  clientId: 'qsl-integration-test'
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
      "sceneId": "SM7AhDv5JgZh6IvK",
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
curl -X PUT 'http://localhost:3011/scene?clientId=qsl-integration-test' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"sceneId":"SM7AhDv5JgZh6IvK","data":{"name":"test-scene-updated"}}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/scene'
params = {
    'clientId': 'qsl-integration-test'
}
url = f'{base_url}{path}'

response = requests.put(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "sceneId": "SM7AhDv5JgZh6IvK",
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
  const baseUrl = 'http://localhost:3011';
  const path = '/scene';
  const params = {
    clientId: 'qsl-integration-test'
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
        "sceneId": "SM7AhDv5JgZh6IvK",
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
  3011 ➡️ port
  🔤/scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"sceneId":"SM7AhDv5JgZh6IvK","data":{"name":"test-scene-updated"}}🔤 ➡️ body

  💭 Build HTTP request
  🔤PUT /scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 67❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "update-scene_1782956912583",
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
    "_id": "SM7AhDv5JgZh6IvK",
    "active": false,
    "navigation": true,
    "navOrder": 0,
    "thumb": null,
    "padding": 0.25,
    "shiftX": 0,
    "shiftY": 0,
    "initial": {
      "x": null,
      "y": null,
      "scale": null
    },
    "initialLevel": "defaultLevel0000",
    "tokenVision": true,
    "fog": {
      "mode": 1,
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
    "transition": {
      "type": null,
      "duration": 1500,
      "activeOnly": false
    },
    "drawings": [],
    "tokens": [],
    "levels": [
      {
        "_id": "defaultLevel0000",
        "name": "Level",
        "elevation": {
          "bottom": 0,
          "top": 20
        },
        "background": {
          "color": "#999999",
          "src": null,
          "tint": "#ffffff",
          "alphaThreshold": 0.75
        },
        "foreground": {
          "src": null,
          "tint": "#ffffff",
          "alphaThreshold": 0.75
        },
        "fog": {
          "src": null,
          "tint": "#ffffff"
        },
        "textures": {
          "anchorX": 0.5,
          "anchorY": 0.5,
          "offsetX": 0,
          "offsetY": 0,
          "fit": "fill",
          "scaleX": 1,
          "scaleY": 1,
          "rotation": 0
        },
        "visibility": {
          "levels": []
        },
        "sort": 0,
        "flags": {}
      }
    ],
    "lights": [],
    "notes": [],
    "sounds": [],
    "regions": [],
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
      "cpvaGKk3hgoBCzCS": 3
    },
    "flags": {},
    "_stats": {
      "coreVersion": "14.363",
      "systemId": "dnd5e",
      "systemVersion": "5.0.4",
      "createdTime": 1782956912564,
      "modifiedTime": 1782956912584,
      "lastModifiedBy": "cpvaGKk3hgoBCzCS",
      "compendiumSource": null,
      "duplicateSource": null,
      "exportSource": null
    }
  }
}
```


---

## DELETE /scene

Delete a scene

**Required scope:** `scene:write`

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
const baseUrl = 'http://localhost:3011';
const path = '/scene';
const params = {
  clientId: 'qsl-integration-test',
  sceneId: 'Hk2XrjCAU9pHZ7Gc'
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
curl -X DELETE 'http://localhost:3011/scene?clientId=qsl-integration-test&sceneId=Hk2XrjCAU9pHZ7Gc' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/scene'
params = {
    'clientId': 'qsl-integration-test',
    'sceneId': 'Hk2XrjCAU9pHZ7Gc'
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
  const baseUrl = 'http://localhost:3011';
  const path = '/scene';
  const params = {
    clientId: 'qsl-integration-test',
    sceneId: 'Hk2XrjCAU9pHZ7Gc'
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
  3011 ➡️ port
  🔤/scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤sceneId=Hk2XrjCAU9pHZ7Gc🔤 ➡️ sceneId
  🔤?🧲clientId🧲&🧲sceneId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "delete-scene_1782956917608",
  "success": true
}
```


---

## POST /switch-scene

Switch the active scene

**Required scope:** `scene:write`

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
const baseUrl = 'http://localhost:3011';
const path = '/switch-scene';
const params = {
  clientId: 'qsl-integration-test'
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
      "sceneId": "SM7AhDv5JgZh6IvK"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3011/switch-scene?clientId=qsl-integration-test' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"sceneId":"SM7AhDv5JgZh6IvK"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3011'
path = '/switch-scene'
params = {
    'clientId': 'qsl-integration-test'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "sceneId": "SM7AhDv5JgZh6IvK"
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
  const baseUrl = 'http://localhost:3011';
  const path = '/switch-scene';
  const params = {
    clientId: 'qsl-integration-test'
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
        "sceneId": "SM7AhDv5JgZh6IvK"
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
  3011 ➡️ port
  🔤/switch-scene🔤 ➡️ path

  💭 Query parameters
  🔤clientId=qsl-integration-test🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"sceneId":"SM7AhDv5JgZh6IvK"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /switch-scene🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3011❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 30❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "switch-scene_1782956912588",
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
    "_id": "SM7AhDv5JgZh6IvK",
    "active": true,
    "navigation": true,
    "navOrder": 0,
    "thumb": null,
    "padding": 0.25,
    "shiftX": 0,
    "shiftY": 0,
    "initial": {
      "x": null,
      "y": null,
      "scale": null
    },
    "initialLevel": "defaultLevel0000",
    "tokenVision": true,
    "fog": {
      "mode": 1,
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
    "transition": {
      "type": null,
      "duration": 1500,
      "activeOnly": false
    },
    "drawings": [],
    "tokens": [],
    "levels": [
      {
        "_id": "defaultLevel0000",
        "name": "Level",
        "elevation": {
          "bottom": 0,
          "top": 20
        },
        "background": {
          "color": "#999999",
          "src": null,
          "tint": "#ffffff",
          "alphaThreshold": 0.75
        },
        "foreground": {
          "src": null,
          "tint": "#ffffff",
          "alphaThreshold": 0.75
        },
        "fog": {
          "src": null,
          "tint": "#ffffff"
        },
        "textures": {
          "anchorX": 0.5,
          "anchorY": 0.5,
          "offsetX": 0,
          "offsetY": 0,
          "fit": "fill",
          "scaleX": 1,
          "scaleY": 1,
          "rotation": 0
        },
        "visibility": {
          "levels": []
        },
        "sort": 0,
        "flags": {}
      }
    ],
    "lights": [],
    "notes": [],
    "sounds": [],
    "regions": [],
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
      "cpvaGKk3hgoBCzCS": 3
    },
    "flags": {},
    "_stats": {
      "coreVersion": "14.363",
      "systemId": "dnd5e",
      "systemVersion": "5.0.4",
      "createdTime": 1782956912564,
      "modifiedTime": 1782956912596,
      "lastModifiedBy": "cpvaGKk3hgoBCzCS",
      "compendiumSource": null,
      "duplicateSource": null,
      "exportSource": null
    }
  }
}
```


---

## GET /scene/image

Get a rendered screenshot of a scene

Captures the full rendered canvas of a scene including all visible layers (tokens, lights, walls, etc.) as an image. The scene can be specified by ID or defaults to the active scene.

**Required scope:** `scene:read`

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

