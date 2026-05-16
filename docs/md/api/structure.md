---
tag: structure
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Structure

## GET /structure

Get the structure of the Foundry world

Retrieves the folder and compendium structure for the specified Foundry world.

**Required scope:** `structure:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| includeEntityData | boolean |  | query | Whether to include full entity data or just UUIDs and names |
| path | string |  | query | Path to read structure from (null = root) |
| recursive | boolean |  | query | Whether to read down the folder tree |
| recursiveDepth | number |  | query | Depth to recurse into folders (default 5) |
| types | string |  | query | Types to return (Scene/Actor/Item/JournalEntry/RollTable/Cards/Macro/Playlist), can be comma-separated or JSON array |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The folder and compendium structure

### Try It Out

<ApiTester
  method="GET"
  path="/structure"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"includeEntityData","type":"boolean","required":false,"source":"query"},{"name":"path","type":"string","required":false,"source":"query"},{"name":"recursive","type":"boolean","required":false,"source":"query"},{"name":"recursiveDepth","type":"number","required":false,"source":"query"},{"name":"types","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/structure';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
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
curl -X GET 'http://localhost:3010/structure?clientId=fvtt_099ad17ea199e7e3&includeEntityData=true&recursive=true&types=Scene' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/structure'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
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
    clientId: 'fvtt_099ad17ea199e7e3',
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
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
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
  "type": "structure-result",
  "requestId": "structure_1778896423682",
  "data": {
    "folders": {
      "test-folder": {
        "id": "qVncHkNOxAuctXy2",
        "uuid": "Folder.qVncHkNOxAuctXy2",
        "type": "Scene"
      }
    },
    "entities": {
      "scenes": [
        {
          "_id": "i01IgzYFzGddbTtP",
          "name": "Scene",
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
            "createdTime": 1776529787190,
            "modifiedTime": 1776821577409,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "2xFy4d19bDsP08Aw",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778789514438,
            "modifiedTime": 1778789514438,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "9AxQS1AJfhieexNT",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778878339415,
            "modifiedTime": 1778878339415,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "E216rK2779fPDqiH",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778810084151,
            "modifiedTime": 1778810084151,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "FIqMGxtWklfvb80c",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778809562920,
            "modifiedTime": 1778809562920,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "HJ3cGhzRV8h0KqHp",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778787172506,
            "modifiedTime": 1778787172506,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "OoTPjYkL2GjuQ2a7",
          "name": "test",
          "active": false,
          "navigation": false,
          "navOrder": 0,
          "navName": "",
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
          "foregroundElevation": 20,
          "thumb": "worlds/5e-tables/assets/scenes/OoTPjYkL2GjuQ2a7-thumb.webp",
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
          "drawings": [
            {
              "_id": "7us1TaI1NwwGdqLs",
              "author": "5ypAoBvOiyjDKiaZ",
              "shape": {
                "type": "r",
                "width": 287,
                "height": 342,
                "radius": null,
                "points": []
              },
              "x": 1942.5,
              "y": 1581.75,
              "elevation": 0,
              "sort": 0,
              "rotation": 0,
              "bezierFactor": 0.5,
              "fillType": 0,
              "fillColor": "#ffffff",
              "fillAlpha": 0.5,
              "strokeWidth": 8,
              "strokeColor": "#ffffff",
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
            },
            {
              "_id": "nrcBVMkwWw8XPcBj",
              "author": "5ypAoBvOiyjDKiaZ",
              "shape": {
                "type": "r",
                "width": 518,
                "height": 583,
                "radius": null,
                "points": []
              },
              "x": 1359.75,
              "y": 2025.75,
              "elevation": 0,
              "sort": 1,
              "rotation": 0,
              "bezierFactor": 0.5,
              "fillType": 0,
              "fillColor": "#ffffff",
              "fillAlpha": 0.5,
              "strokeWidth": 8,
              "strokeColor": "#ffffff",
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
              "_id": "ek1v2zxBxkZEftca",
              "name": "Perrin",
              "displayName": 30,
              "actorId": "xSst5kAigAZw6wDr",
              "actorLink": true,
              "delta": null,
              "width": 1,
              "height": 1,
              "texture": {
                "src": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
                "anchorX": 0.5,
                "anchorY": 0.5,
                "offsetX": 0,
                "offsetY": 0,
                "fit": "contain",
                "scaleX": 0.8,
                "scaleY": 0.8,
                "rotation": 0,
                "tint": "#ffffff",
                "alphaThreshold": 0.75
              },
              "shape": 4,
              "x": 1480,
              "y": 1554,
              "elevation": 0,
              "sort": 0,
              "locked": false,
              "lockRotation": false,
              "rotation": 0,
              "alpha": 1,
              "hidden": false,
              "disposition": 1,
              "displayBars": 40,
              "bar1": {
                "attribute": "attributes.hp"
              },
              "bar2": {
                "attribute": "attributes.ac.value"
              },
              "light": {
                "negative": false,
                "priority": 0,
                "alpha": 1,
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
                "range": 5,
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
                  "range": 5
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
              "_id": "RFgXdhctPXhpSyOM",
              "x": 1628,
              "y": 1443,
              "elevation": 0,
              "rotation": 0,
              "walls": true,
              "vision": false,
              "config": {
                "negative": false,
                "priority": 0,
                "alpha": 0.5,
                "angle": 360,
                "bright": 24.04,
                "color": null,
                "coloration": 1,
                "dim": 48.09,
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
              "_id": "YqzoO6yW51CQtt0k",
              "entryId": "u7byD1yDxgtzqeT4",
              "pageId": null,
              "x": 1665,
              "y": 2109,
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
              "text": "jhgcfkjhgcv",
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
              "_id": "kWlgWosbxU3tbSj0",
              "x": 1147,
              "y": 2405,
              "elevation": 0,
              "radius": 48.02,
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
                  "type": "",
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
              "_id": "PaZoZJk6XppvofiN",
              "name": "Region",
              "color": "#2893cc",
              "shapes": [
                {
                  "type": "rectangle",
                  "hole": false,
                  "x": 1433.75,
                  "y": 2016.5,
                  "width": 582.75,
                  "height": 407,
                  "rotation": 0
                },
                {
                  "type": "rectangle",
                  "hole": false,
                  "x": 1757.5,
                  "y": 2358.75,
                  "width": 471.75,
                  "height": 786.25,
                  "rotation": 0
                },
                {
                  "type": "ellipse",
                  "hole": false,
                  "x": 1729.75,
                  "y": 1877.75,
                  "radiusX": 296,
                  "radiusY": 74,
                  "rotation": 0
                }
              ],
              "elevation": {
                "bottom": null,
                "top": null
              },
              "behaviors": [],
              "visibility": 0,
              "locked": false,
              "flags": {}
            }
          ],
          "templates": [
            {
              "_id": "iJq4Fzo3Kvln5zVP",
              "author": "5ypAoBvOiyjDKiaZ",
              "t": "cone",
              "x": 1332,
              "y": 2886,
              "elevation": 0,
              "sort": 0,
              "distance": 77.86205751198719,
              "direction": 275.5275401516562,
              "angle": 53.13,
              "width": 0,
              "borderColor": "#000000",
              "fillColor": "#28cca2",
              "texture": null,
              "hidden": false,
              "flags": {}
            }
          ],
          "tiles": [
            {
              "_id": "b0cYOogM1zPL9NFc",
              "texture": {
                "src": null,
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
              "width": 222,
              "height": 333,
              "x": 1628,
              "y": 1369,
              "elevation": 0,
              "sort": 0,
              "rotation": 0,
              "alpha": 1,
              "hidden": false,
              "locked": false,
              "restrictions": {
                "light": false,
                "weather": false
              },
              "occlusion": {
                "mode": 0,
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
          "walls": [
            {
              "_id": "LSziSUTponH836zp",
              "c": [
                1286,
                1674,
                1656,
                2581
              ],
              "light": 20,
              "move": 20,
              "sight": 20,
              "sound": 20,
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
            },
            {
              "_id": "QICTFX7JLm3W8PLs",
              "c": [
                1970,
                1711,
                2017,
                2840
              ],
              "light": 20,
              "move": 20,
              "sight": 20,
              "sound": 20,
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
            },
            {
              "_id": "dOBYPT8ja94pUuFV",
              "c": [
                1230,
                2701,
                2424,
                1452
              ],
              "light": 20,
              "move": 20,
              "sight": 20,
              "sound": 20,
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
            },
            {
              "_id": "wo7Vxm6QrHhI7ySh",
              "c": [
                1822,
                1434,
                2479,
                2951
              ],
              "light": 20,
              "move": 20,
              "sight": 20,
              "sound": 20,
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
            },
            {
              "_id": "qhAnrbcpLw5UCnJN",
              "c": [
                1989,
                1591,
                1203,
                3016
              ],
              "light": 20,
              "move": 20,
              "sight": 20,
              "sound": 20,
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
            },
            {
              "_id": "1wlEJl7SD94Ia2iY",
              "c": [
                1277,
                1721,
                2387,
                2035
              ],
              "light": 20,
              "move": 20,
              "sight": 20,
              "sound": 20,
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
            "5ypAoBvOiyjDKiaZ": 3,
            "default": 0,
            "lmPGreUzUZ4YxH6D": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1777176992560,
            "modifiedTime": 1777176992560,
            "lastModifiedBy": "5ypAoBvOiyjDKiaZ",
            "compendiumSource": null,
            "duplicateSource": "Scene.Fpu0odizNnNnjlLI",
            "exportSource": null
          }
        },
        {
          "_id": "X914jr7smFaAxr72",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778812415693,
            "modifiedTime": 1778812415693,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "azpFAQYHh0s3yHML",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778789859370,
            "modifiedTime": 1778789859370,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "nsoLECqbXXs1PrBZ",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778810347466,
            "modifiedTime": 1778810347466,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "uisZmFzI2Tjh1xSA",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778718765802,
            "modifiedTime": 1778718765802,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "v5E3y16p9juupLJS",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778884093752,
            "modifiedTime": 1778884093752,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "7iYl9ExwMdFm9POw",
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1776824906803,
            "modifiedTime": 1776824911917,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "aQADc2ek0f7ls9af",
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
              "_id": "91dY5wnZ9pbcxcC0",
              "name": "Unknown",
              "displayName": 0,
              "actorId": "5bW7ahcfLR0uMb9j",
              "actorLink": false,
              "delta": {
                "_id": "xhCuAiQoswAzsnqY",
                "system": {},
                "items": [],
                "effects": [
                  {
                    "_id": "dnd5edead0000000",
                    "name": "Dead",
                    "img": "systems/dnd5e/icons/svg/statuses/dead.svg",
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
                    "statuses": [
                      "dead"
                    ],
                    "sort": 0,
                    "flags": {},
                    "_stats": {
                      "coreVersion": "13.348",
                      "systemId": "dnd5e",
                      "systemVersion": "5.0.4",
                      "createdTime": 1776821644329,
                      "modifiedTime": 1776821644329,
                      "lastModifiedBy": "r6bXhB7k9cXa3cif",
                      "compendiumSource": null,
                      "duplicateSource": null,
                      "exportSource": null
                    }
                  }
                ],
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
              "x": 200,
              "y": 200,
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
            },
            {
              "_id": "a0alpONow8ppf7jZ",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "TaWAifIKUZ3FSfEj",
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
              "movementAction": "walk",
              "_movementHistory": [],
              "_regions": [],
              "flags": {}
            },
            {
              "_id": "9H8zRN6jFYAKi4St",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "RbCCqngEl6EST2GA",
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
              "movementAction": "walk",
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1776821602279,
            "modifiedTime": 1776824906859,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "cDiBNArUBdiC76Xd",
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
              "_id": "34oGFTyXvbLYWpWZ",
              "name": "Updated Test Actor",
              "displayName": 0,
              "actorId": "gtRH6ZKlyzPYdvmW",
              "actorLink": false,
              "delta": {
                "_id": "rEXYan5oaapT93Rm",
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
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.races.LOMdcNAGWh5xpfm4"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.races.7Yoo9hG0hfPSmBoC"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.races.PqxZgcJzp1VVgP8t"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.races.GWPjKFeIthBBeCFJ"
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
                          "level": 0,
                          "title": "",
                          "type": "ItemGrant",
                          "value": {
                            "added": {
                              "AArhiOrSkaQUnCZS": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ",
                              "FtOM4QiOW5MwgcS3": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
                              "cWrETHzCRs1Ueqd3": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
                              "nmmihiqphHjoE8dl": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC"
                            }
                          }
                        },
                        {
                          "_id": "Z9hvZFkWUNvowbQX",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 1,
                              "con": 0,
                              "dex": 2,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 0
                          },
                          "level": 0,
                          "title": "",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "assignments": {
                              "cha": 1,
                              "dex": 2
                            },
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "hv2bcANK5jEJZaAb",
                          "configuration": {
                            "sizes": [
                              "sm"
                            ]
                          },
                          "hint": "Halflings average about 3 feet tall and weigh about 40 pounds. Your size is Small.",
                          "level": 1,
                          "title": "",
                          "type": "Size",
                          "value": {
                            "size": "sm"
                          }
                        },
                        {
                          "_id": "nGwMjsfNU6CXHk3A",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [],
                            "grants": [
                              "languages:standard:common",
                              "languages:standard:halfling"
                            ],
                            "mode": "default"
                          },
                          "level": 0,
                          "title": "",
                          "type": "Trait",
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
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        },
                        "sourceId": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4"
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
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        },
                        "sourceId": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC"
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
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        },
                        "sourceId": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t"
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
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        },
                        "sourceId": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ"
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
                          "_id": "ocxNtDFJ7YDaYaK7",
                          "configuration": {},
                          "icon": "systems/dnd5e/icons/svg/hit-points.svg",
                          "title": "Hit Points",
                          "type": "HitPoints",
                          "value": {
                            "1": "max"
                          }
                        },
                        {
                          "_id": "mmAxx3U7FvXNAcKc",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [],
                            "grants": [
                              "weapon:sim",
                              "weapon:mar:shortsword"
                            ],
                            "mode": "default"
                          },
                          "level": 1,
                          "title": "",
                          "type": "Trait",
                          "value": {
                            "chosen": [
                              "weapon:sim",
                              "weapon:mar:shortsword"
                            ]
                          }
                        },
                        {
                          "_id": "QPXy59CQGY9HB0c3",
                          "classRestriction": "primary",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [
                              {
                                "count": 1,
                                "pool": [
                                  "tool:art:*",
                                  "tool:music:*"
                                ]
                              }
                            ],
                            "grants": [],
                            "mode": "default"
                          },
                          "level": 1,
                          "title": "",
                          "type": "Trait",
                          "value": {
                            "chosen": [
                              "tool:art:brewer"
                            ]
                          }
                        },
                        {
                          "_id": "4M8MQ1E64zbcRg6B",
                          "classRestriction": "primary",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [],
                            "grants": [
                              "saves:str",
                              "saves:dex"
                            ],
                            "mode": "default"
                          },
                          "level": 1,
                          "title": "",
                          "type": "Trait",
                          "value": {
                            "chosen": [
                              "saves:str",
                              "saves:dex"
                            ]
                          }
                        },
                        {
                          "_id": "7HRRCPk80Ng2Evdx",
                          "classRestriction": "primary",
                          "configuration": {
                            "allowReplacements": false,
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
                            ],
                            "grants": [],
                            "mode": "default"
                          },
                          "level": 1,
                          "title": "",
                          "type": "Trait",
                          "value": {
                            "chosen": [
                              "skills:acr",
                              "skills:ath"
                            ]
                          }
                        },
                        {
                          "_id": "BQWHr3mt5flvkfIj",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [],
                            "grants": [
                              "di:poison",
                              "ci:diseased",
                              "ci:poisoned"
                            ],
                            "mode": "default"
                          },
                          "hint": "Your mastery of the ki flowing through you makes you immune to disease and poison.",
                          "level": 10,
                          "title": "Purity of Body",
                          "type": "Trait",
                          "value": {
                            "chosen": []
                          }
                        },
                        {
                          "_id": "n0q8XyiGA3vLPgpK",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.UAvV7N7T4zJhxdfI"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.l50hjTxO2r0iecKw"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.eGxoNmSMWKNzChCO"
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
                          "level": 1,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {
                            "added": {
                              "CwgoTDXWCD7PknIN": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
                              "RiURabP4FDYMeuWx": "Compendium.dnd5e.classfeatures.Item.eGxoNmSMWKNzChCO",
                              "pchnXqd5C79fVlxy": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw"
                            }
                          }
                        },
                        {
                          "_id": "7TyDqpGGi3r3nsp0",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.10b6z2W1txNkrGP7"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.5MwNlVZK7m6VolOH"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.TDglPcxIVEzvVSgK"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.yrSFIGTaQOH2PFRI"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.zCeqyQ8uIPNdYJSW"
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
                          "level": 2,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "2sLHTw6k15DSW8WB",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.rtpQdX77dYWbDIOH"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.mzweVbnsJPQiVkAe"
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
                          "level": 3,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "Zc1jOZK1b9mIKekq",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.KQz9bqxVkXjDl8gK"
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
                          "level": 4,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "j9LeWmxlsENKaMLo",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.XogoBnFWmCAHXppo"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.pvRc6GAu1ok6zihC"
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
                          "level": 5,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "psobDjMqtA2216Db",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.7flZKruSSu6dHg6D"
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
                          "level": 6,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "K38aFaEMxMqRB0BC",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.a4P4DNMmH8CqSNkC"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.ZmC31XKS4YNENnoc"
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
                          "level": 7,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "eLqmJotmwzlGNrxG",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.bqWA7t9pDELbNRkp"
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
                          "level": 10,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "N0geIQiuofqYgswj",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.XjuGBeB8Y0C3A5D4"
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
                          "level": 13,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "N1hjizyI82UPp8UI",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.7D2EkLdISwShEDlN"
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
                          "level": 14,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "TcLZS9WzC7bPETSd",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.gDH8PMrKvLHaNmEI"
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
                          "level": 15,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "tRb3a0tA5IpehPs8",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.3jwFt3hSqDswBlOH"
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
                          "level": 18,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "sEQz9c9XhWYjS9x5",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.mQNPg89YIs7g5tG4"
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
                          "level": 20,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "MXFbf0nxMiyLdPbX",
                          "configuration": {
                            "distance": {
                              "units": ""
                            },
                            "identifier": "die",
                            "scale": {
                              "1": {
                                "faces": 4,
                                "modifiers": [],
                                "number": null
                              },
                              "5": {
                                "faces": 6,
                                "modifiers": [],
                                "number": null
                              },
                              "11": {
                                "faces": 8,
                                "modifiers": [],
                                "number": null
                              },
                              "17": {
                                "faces": 10,
                                "modifiers": [],
                                "number": null
                              }
                            },
                            "type": "dice"
                          },
                          "title": "Martial Arts Die",
                          "type": "ScaleValue",
                          "value": {}
                        },
                        {
                          "_id": "1OzfWDWCquoHMeX5",
                          "configuration": {
                            "distance": {
                              "units": "ft"
                            },
                            "identifier": "unarmored-movement",
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
                            },
                            "type": "distance"
                          },
                          "title": "Unarmored Movement",
                          "type": "ScaleValue",
                          "value": {}
                        },
                        {
                          "_id": "ofNSUhSHKhhDuPSR",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 4,
                          "title": "Ability Score Improvement",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "s3t9o57hP6iUHirr",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 8,
                          "title": "Ability Score Improvement",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "O24MWOKc1ImsKaml",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 12,
                          "title": "Ability Score Improvement",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "xdqWoLtgO3uyl3nJ",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 16,
                          "title": "Ability Score Improvement",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "puDaUsYrlks0z5gm",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 19,
                          "title": "",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "0awj2yq115ev9u9o",
                          "configuration": {},
                          "level": 3,
                          "title": "Monastic Tradition",
                          "type": "Subclass",
                          "value": {
                            "document": null,
                            "uuid": null
                          }
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
                        "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
                        "sourceId": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI"
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
                        "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
                        "sourceId": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw"
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {
                                  "formula": "",
                                  "mode": ""
                                },
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "touch"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "attack": {
                            "ability": "str",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "classification": "weapon",
                              "value": "melee"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {
                                  "formula": "",
                                  "mode": ""
                                },
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": [
                              {
                                "bonus": "",
                                "custom": {
                                  "enabled": true,
                                  "formula": "1"
                                },
                                "denomination": null,
                                "number": null,
                                "scaling": {
                                  "formula": "",
                                  "mode": "whole",
                                  "number": null
                                },
                                "types": [
                                  "fire"
                                ]
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "hour",
                            "value": "1"
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "self"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "40",
                              "type": "radius",
                              "units": "ft",
                              "width": ""
                            }
                          },
                          "type": "attack",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "classification": "weapon",
                              "value": "ranged"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "ft",
                            "value": "20"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "attack",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "",
                            "value": null
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "self"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {
                                  "formula": "",
                                  "mode": ""
                                },
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "touch"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "1",
                              "special": "",
                              "type": "creature"
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "classification": "unarmed",
                              "value": "melee"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "name": "",
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "ft",
                            "value": "5"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "attack",
                          "uses": {
                            "max": "",
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {
                                  "formula": "",
                                  "mode": ""
                                },
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "ft",
                            "value": "5"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "check": {
                            "ability": "str",
                            "associated": [],
                            "dc": {
                              "calculation": "",
                              "formula": "17"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {},
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "name": "Burst",
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "self"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "check",
                          "uses": {
                            "max": "",
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "",
                            "value": null
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "self"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "classification": "weapon",
                              "value": "melee"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "ft",
                            "value": "5"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "attack",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
              "x": 200,
              "y": 200,
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
            "createdTime": 1777968644264,
            "modifiedTime": 1777969049386,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "oAo1nnXqH58UK6h4",
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
              "_id": "zwTXwUu0iIZZF8WF",
              "name": "Unknown",
              "displayName": 0,
              "actorId": "S2YDx7AJKBzunIei",
              "actorLink": false,
              "delta": {
                "_id": "Us5thMXbAeLCYjPV",
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
            },
            {
              "_id": "Q8A3HTBqwpyBxvbp",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "h4CZd8cuYlKo1MAt",
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
              "movementAction": "walk",
              "_movementHistory": [],
              "_regions": [],
              "flags": {}
            },
            {
              "_id": "ATkbFAc89hrONO8B",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "9dHverB9J6UHsFcR",
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
              "movementAction": "walk",
              "_movementHistory": [],
              "_regions": [],
              "flags": {}
            },
            {
              "_id": "vW4JKDTKqJl8IW86",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "biM8UCEgO5XehwU4",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 3,
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
            },
            {
              "_id": "VcWUEgg6Dz7Qsc5k",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "F8clfJsmAzhh3ham",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 4,
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
            },
            {
              "_id": "8J5WiroyAnGOuBey",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "ipZ5C6BRYWODw9zz",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 5,
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
            },
            {
              "_id": "FFgfk07809UcslsP",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "eCXsi4Je0ERuEY65",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 6,
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
            },
            {
              "_id": "evVFcKwG7dBh6KZa",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "Gojqh1Y9yNUFLOJi",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 7,
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
            },
            {
              "_id": "v8dm2gPR5JLyVJDb",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "1o3vAedx5N1JEgCU",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 8,
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
            },
            {
              "_id": "bgCZWfkgY05XHyll",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "MmaXMA5YP4YnAwm8",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 9,
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
            },
            {
              "_id": "DRjFVQgfwWBkqlcG",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "LaVPgqD95KumJRGB",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 10,
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
            },
            {
              "_id": "IdTS49lsUW6JzIsH",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "idQljZZpsp6TgbR1",
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
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 11,
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
          "lights": [
            {
              "_id": "ksek5ZGBvTrNTFSv",
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
              "_id": "maVuHwDHxlQRa6NE",
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
              "_id": "W06Gn1gfv4IvVcB1",
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
              "_id": "zAhaHEXY69R4NGh0",
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
              "_id": "dpSaa6RhCKsO8TnK",
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
              "_id": "plMXFhz6CC4CnNLr",
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
              "_id": "NLin9vhwECgTo3ar",
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
              "_id": "21qNx4fcGnfNA1M9",
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
              "_id": "7eXjK2HpiILxPq91",
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
              "_id": "zW7B0qwoyQAu4b6q",
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
              "_id": "Dn8jycsQG4hKhjXX",
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
              "_id": "ZZ8pNbwlgVTtB28D",
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
              "_id": "XWtVQHaBAU1uPTad",
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
            },
            {
              "_id": "tGDTJBwnTEo4jjAR",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 2,
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
              "_id": "cnWzGlH8TL8oUcHU",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 3,
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
              "_id": "2SajQNZKCukOoh0U",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 4,
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
              "_id": "B9y1iY1hskADS5IG",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 5,
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
              "_id": "JOBFOFtc6iI4hUpA",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 6,
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
              "_id": "yywdJGVOdCb6oc8A",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 7,
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
              "_id": "dpLvpy6yj61pQ2Bi",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 8,
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
              "_id": "A9uq2gPCHg7LlxZJ",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 9,
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
              "_id": "fcQD4GHs8EVDC1XH",
              "entryId": null,
              "pageId": null,
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 10,
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
              "_id": "bxBAcwTer3xfSAp2",
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
              "_id": "0NToDzSgufk4Behe",
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
              "_id": "L8PuvwXnTRbsk0Hj",
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
              "_id": "yhT82OhqXpa4g5Eb",
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
              "_id": "fkfuR7SokriwwE5y",
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
              "_id": "GI80qrhkKXu1k3HJ",
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
              "_id": "p1PxzmRrqTmr0baS",
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
              "_id": "ULh6j6QBKvtIjP4J",
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
              "_id": "EAprdAOfdvx9nzsD",
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
              "_id": "JvRs9A6YWT34JbBT",
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
              "_id": "NWBXDs1MULeGaaqD",
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
              "_id": "KdDmzBOrUjGyfhqO",
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
              "_id": "zwDtmrEQsz7Ai3QB",
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
            },
            {
              "_id": "YHgZ7qQu1N8DPiWP",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 2,
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
              "_id": "JnttWgf26eS62zMU",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 3,
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
              "_id": "XzFJgMKV5mfaYnLm",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 4,
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
              "_id": "f6DA5KHyaeyjyVPy",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 5,
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
              "_id": "nnBCH4ffKl4ftglP",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 6,
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
              "_id": "KxG1dtLsOBRH0pVw",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 7,
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
              "_id": "m7hFbmbEhDQSuyZe",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 8,
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
              "_id": "nI7SRRGPm3yTKTnQ",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 9,
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
              "_id": "dvytsVdA6agkd7UR",
              "author": "r6bXhB7k9cXa3cif",
              "t": "circle",
              "x": 0,
              "y": 0,
              "elevation": 0,
              "sort": 10,
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1777969049357,
            "modifiedTime": 1778896417833,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "pxk3rKsNgpwB6bG5",
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1776821577383,
            "modifiedTime": 1776821602303,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "vCLJOw0STWGphIWU",
          "name": "test-scene-updated",
          "active": false,
          "navigation": true,
          "navOrder": 0,
          "navName": "",
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
              "_id": "TNEsC0gFgz1Nc9pg",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "781VTntlraAbO7C3",
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
              "movementAction": "walk",
              "_movementHistory": [],
              "_regions": [],
              "flags": {}
            },
            {
              "_id": "Yi9FWJpDX5TPvXK3",
              "name": "test",
              "displayName": 0,
              "actorId": null,
              "actorLink": false,
              "delta": {
                "_id": "8Wwm8EDsy6y1C9aG",
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
              "movementAction": "walk",
              "_movementHistory": [],
              "_regions": [],
              "flags": {}
            },
            {
              "_id": "KV0Ua6gg79bZ456e",
              "name": "Perrin",
              "displayName": 30,
              "actorId": "xSst5kAigAZw6wDr",
              "actorLink": true,
              "delta": null,
              "width": 1,
              "height": 1,
              "texture": {
                "src": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
                "anchorX": 0.5,
                "anchorY": 0.5,
                "offsetX": 0,
                "offsetY": 0,
                "fit": "contain",
                "scaleX": 0.8,
                "scaleY": 0.8,
                "rotation": 0,
                "tint": "#ffffff",
                "alphaThreshold": 0.75
              },
              "shape": 4,
              "x": 300,
              "y": 500,
              "elevation": 0,
              "sort": 2,
              "locked": false,
              "lockRotation": false,
              "rotation": 0,
              "alpha": 1,
              "hidden": false,
              "disposition": 1,
              "displayBars": 40,
              "bar1": {
                "attribute": "attributes.hp"
              },
              "bar2": {
                "attribute": "attributes.ac.value"
              },
              "light": {
                "negative": false,
                "priority": 0,
                "alpha": 1,
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
                "range": 5,
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
                  "range": 5
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
              "_id": "LDpSlXWcYlP7hudR",
              "c": [
                650,
                338,
                413,
                613
              ],
              "light": 20,
              "move": 20,
              "sight": 20,
              "sound": 20,
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
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1776824911885,
            "modifiedTime": 1777968644307,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "InuAccGQ1wWqwPfm",
          "name": "test",
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
            "r6bXhB7k9cXa3cif": 3
          },
          "flags": {},
          "_stats": {
            "coreVersion": "13.348",
            "systemId": "dnd5e",
            "systemVersion": "5.0.4",
            "createdTime": 1778896413295,
            "modifiedTime": 1778896413295,
            "lastModifiedBy": "r6bXhB7k9cXa3cif",
            "compendiumSource": null,
            "duplicateSource": null,
            "exportSource": null
          }
        },
        {
          "_id": "tgYnjCnq6EFiREjl",
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
              "_id": "TOInpAfNZ6GRAJaG",
              "name": "Updated Test Actor",
              "displayName": 0,
              "actorId": "S6Nh5SfYPcyyUhRg",
              "actorLink": false,
              "delta": {
                "_id": "xFsLxQX7Lmyc4TaG",
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
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.races.LOMdcNAGWh5xpfm4"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.races.7Yoo9hG0hfPSmBoC"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.races.PqxZgcJzp1VVgP8t"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.races.GWPjKFeIthBBeCFJ"
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
                          "level": 0,
                          "title": "",
                          "type": "ItemGrant",
                          "value": {
                            "added": {
                              "AArhiOrSkaQUnCZS": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ",
                              "FtOM4QiOW5MwgcS3": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
                              "cWrETHzCRs1Ueqd3": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
                              "nmmihiqphHjoE8dl": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC"
                            }
                          }
                        },
                        {
                          "_id": "Z9hvZFkWUNvowbQX",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 1,
                              "con": 0,
                              "dex": 2,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 0
                          },
                          "level": 0,
                          "title": "",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "assignments": {
                              "cha": 1,
                              "dex": 2
                            },
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "hv2bcANK5jEJZaAb",
                          "configuration": {
                            "sizes": [
                              "sm"
                            ]
                          },
                          "hint": "Halflings average about 3 feet tall and weigh about 40 pounds. Your size is Small.",
                          "level": 1,
                          "title": "",
                          "type": "Size",
                          "value": {
                            "size": "sm"
                          }
                        },
                        {
                          "_id": "nGwMjsfNU6CXHk3A",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [],
                            "grants": [
                              "languages:standard:common",
                              "languages:standard:halfling"
                            ],
                            "mode": "default"
                          },
                          "level": 0,
                          "title": "",
                          "type": "Trait",
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
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        },
                        "sourceId": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4"
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
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        },
                        "sourceId": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC"
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
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        },
                        "sourceId": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t"
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
                        "advancementOrigin": "O3ZjSw0GtAOPScHI.nInhIgkbzzJTdm8F",
                        "riders": {
                          "activity": [],
                          "effect": []
                        },
                        "sourceId": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ"
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
                          "_id": "ocxNtDFJ7YDaYaK7",
                          "configuration": {},
                          "icon": "systems/dnd5e/icons/svg/hit-points.svg",
                          "title": "Hit Points",
                          "type": "HitPoints",
                          "value": {
                            "1": "max"
                          }
                        },
                        {
                          "_id": "mmAxx3U7FvXNAcKc",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [],
                            "grants": [
                              "weapon:sim",
                              "weapon:mar:shortsword"
                            ],
                            "mode": "default"
                          },
                          "level": 1,
                          "title": "",
                          "type": "Trait",
                          "value": {
                            "chosen": [
                              "weapon:sim",
                              "weapon:mar:shortsword"
                            ]
                          }
                        },
                        {
                          "_id": "QPXy59CQGY9HB0c3",
                          "classRestriction": "primary",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [
                              {
                                "count": 1,
                                "pool": [
                                  "tool:art:*",
                                  "tool:music:*"
                                ]
                              }
                            ],
                            "grants": [],
                            "mode": "default"
                          },
                          "level": 1,
                          "title": "",
                          "type": "Trait",
                          "value": {
                            "chosen": [
                              "tool:art:brewer"
                            ]
                          }
                        },
                        {
                          "_id": "4M8MQ1E64zbcRg6B",
                          "classRestriction": "primary",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [],
                            "grants": [
                              "saves:str",
                              "saves:dex"
                            ],
                            "mode": "default"
                          },
                          "level": 1,
                          "title": "",
                          "type": "Trait",
                          "value": {
                            "chosen": [
                              "saves:str",
                              "saves:dex"
                            ]
                          }
                        },
                        {
                          "_id": "7HRRCPk80Ng2Evdx",
                          "classRestriction": "primary",
                          "configuration": {
                            "allowReplacements": false,
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
                            ],
                            "grants": [],
                            "mode": "default"
                          },
                          "level": 1,
                          "title": "",
                          "type": "Trait",
                          "value": {
                            "chosen": [
                              "skills:acr",
                              "skills:ath"
                            ]
                          }
                        },
                        {
                          "_id": "BQWHr3mt5flvkfIj",
                          "configuration": {
                            "allowReplacements": false,
                            "choices": [],
                            "grants": [
                              "di:poison",
                              "ci:diseased",
                              "ci:poisoned"
                            ],
                            "mode": "default"
                          },
                          "hint": "Your mastery of the ki flowing through you makes you immune to disease and poison.",
                          "level": 10,
                          "title": "Purity of Body",
                          "type": "Trait",
                          "value": {
                            "chosen": []
                          }
                        },
                        {
                          "_id": "n0q8XyiGA3vLPgpK",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.UAvV7N7T4zJhxdfI"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.l50hjTxO2r0iecKw"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.eGxoNmSMWKNzChCO"
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
                          "level": 1,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {
                            "added": {
                              "CwgoTDXWCD7PknIN": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
                              "RiURabP4FDYMeuWx": "Compendium.dnd5e.classfeatures.Item.eGxoNmSMWKNzChCO",
                              "pchnXqd5C79fVlxy": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw"
                            }
                          }
                        },
                        {
                          "_id": "7TyDqpGGi3r3nsp0",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.10b6z2W1txNkrGP7"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.5MwNlVZK7m6VolOH"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.TDglPcxIVEzvVSgK"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.yrSFIGTaQOH2PFRI"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.zCeqyQ8uIPNdYJSW"
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
                          "level": 2,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "2sLHTw6k15DSW8WB",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.rtpQdX77dYWbDIOH"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.mzweVbnsJPQiVkAe"
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
                          "level": 3,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "Zc1jOZK1b9mIKekq",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.KQz9bqxVkXjDl8gK"
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
                          "level": 4,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "j9LeWmxlsENKaMLo",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.XogoBnFWmCAHXppo"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.pvRc6GAu1ok6zihC"
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
                          "level": 5,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "psobDjMqtA2216Db",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.7flZKruSSu6dHg6D"
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
                          "level": 6,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "K38aFaEMxMqRB0BC",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.a4P4DNMmH8CqSNkC"
                              },
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.ZmC31XKS4YNENnoc"
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
                          "level": 7,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "eLqmJotmwzlGNrxG",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.bqWA7t9pDELbNRkp"
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
                          "level": 10,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "N0geIQiuofqYgswj",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.XjuGBeB8Y0C3A5D4"
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
                          "level": 13,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "N1hjizyI82UPp8UI",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.7D2EkLdISwShEDlN"
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
                          "level": 14,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "TcLZS9WzC7bPETSd",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.gDH8PMrKvLHaNmEI"
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
                          "level": 15,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "tRb3a0tA5IpehPs8",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.3jwFt3hSqDswBlOH"
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
                          "level": 18,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "sEQz9c9XhWYjS9x5",
                          "configuration": {
                            "items": [
                              {
                                "optional": false,
                                "uuid": "Compendium.dnd5e.classfeatures.mQNPg89YIs7g5tG4"
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
                          "level": 20,
                          "title": "Features",
                          "type": "ItemGrant",
                          "value": {}
                        },
                        {
                          "_id": "MXFbf0nxMiyLdPbX",
                          "configuration": {
                            "distance": {
                              "units": ""
                            },
                            "identifier": "die",
                            "scale": {
                              "1": {
                                "faces": 4,
                                "modifiers": [],
                                "number": null
                              },
                              "5": {
                                "faces": 6,
                                "modifiers": [],
                                "number": null
                              },
                              "11": {
                                "faces": 8,
                                "modifiers": [],
                                "number": null
                              },
                              "17": {
                                "faces": 10,
                                "modifiers": [],
                                "number": null
                              }
                            },
                            "type": "dice"
                          },
                          "title": "Martial Arts Die",
                          "type": "ScaleValue",
                          "value": {}
                        },
                        {
                          "_id": "1OzfWDWCquoHMeX5",
                          "configuration": {
                            "distance": {
                              "units": "ft"
                            },
                            "identifier": "unarmored-movement",
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
                            },
                            "type": "distance"
                          },
                          "title": "Unarmored Movement",
                          "type": "ScaleValue",
                          "value": {}
                        },
                        {
                          "_id": "ofNSUhSHKhhDuPSR",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 4,
                          "title": "Ability Score Improvement",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "s3t9o57hP6iUHirr",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 8,
                          "title": "Ability Score Improvement",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "O24MWOKc1ImsKaml",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 12,
                          "title": "Ability Score Improvement",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "xdqWoLtgO3uyl3nJ",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 16,
                          "title": "Ability Score Improvement",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "puDaUsYrlks0z5gm",
                          "configuration": {
                            "cap": 2,
                            "fixed": {
                              "cha": 0,
                              "con": 0,
                              "dex": 0,
                              "int": 0,
                              "str": 0,
                              "wis": 0
                            },
                            "locked": [],
                            "points": 2
                          },
                          "level": 19,
                          "title": "",
                          "type": "AbilityScoreImprovement",
                          "value": {
                            "type": "asi"
                          }
                        },
                        {
                          "_id": "0awj2yq115ev9u9o",
                          "configuration": {},
                          "level": 3,
                          "title": "Monastic Tradition",
                          "type": "Subclass",
                          "value": {
                            "document": null,
                            "uuid": null
                          }
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
                        "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
                        "sourceId": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI"
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
                        "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
                        "sourceId": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw"
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {
                                  "formula": "",
                                  "mode": ""
                                },
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "touch"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "attack": {
                            "ability": "str",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "classification": "weapon",
                              "value": "melee"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {
                                  "formula": "",
                                  "mode": ""
                                },
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": [
                              {
                                "bonus": "",
                                "custom": {
                                  "enabled": true,
                                  "formula": "1"
                                },
                                "denomination": null,
                                "number": null,
                                "scaling": {
                                  "formula": "",
                                  "mode": "whole",
                                  "number": null
                                },
                                "types": [
                                  "fire"
                                ]
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "hour",
                            "value": "1"
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "self"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "40",
                              "type": "radius",
                              "units": "ft",
                              "width": ""
                            }
                          },
                          "type": "attack",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "classification": "weapon",
                              "value": "ranged"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "ft",
                            "value": "20"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "attack",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "",
                            "value": null
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "self"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {
                                  "formula": "",
                                  "mode": ""
                                },
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "touch"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "1",
                              "special": "",
                              "type": "creature"
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "classification": "unarmed",
                              "value": "melee"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "name": "",
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "ft",
                            "value": "5"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "attack",
                          "uses": {
                            "max": "",
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {
                                  "formula": "",
                                  "mode": ""
                                },
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "ft",
                            "value": "5"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "check": {
                            "ability": "str",
                            "associated": [],
                            "dc": {
                              "calculation": "",
                              "formula": "17"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": [
                              {
                                "scaling": {},
                                "target": "",
                                "type": "itemUses",
                                "value": "1"
                              }
                            ]
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "name": "Burst",
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "self"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "check",
                          "uses": {
                            "max": "",
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "",
                            "value": null
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "self"
                          },
                          "roll": {
                            "formula": "",
                            "name": "",
                            "prompt": false,
                            "visible": false
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "utility",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
                          "activation": {
                            "condition": "",
                            "override": false,
                            "type": "action",
                            "value": 1
                          },
                          "attack": {
                            "ability": "",
                            "bonus": "",
                            "critical": {
                              "threshold": null
                            },
                            "flat": false,
                            "type": {
                              "classification": "weapon",
                              "value": "melee"
                            }
                          },
                          "consumption": {
                            "scaling": {
                              "allowed": false,
                              "max": ""
                            },
                            "spellSlot": true,
                            "targets": []
                          },
                          "damage": {
                            "critical": {
                              "bonus": ""
                            },
                            "includeBase": true,
                            "parts": []
                          },
                          "description": {
                            "chatFlavor": ""
                          },
                          "duration": {
                            "concentration": false,
                            "override": false,
                            "special": "",
                            "units": "inst",
                            "value": ""
                          },
                          "effects": [],
                          "range": {
                            "override": false,
                            "special": "",
                            "units": "ft",
                            "value": "5"
                          },
                          "sort": 0,
                          "target": {
                            "affects": {
                              "choice": false,
                              "count": "",
                              "special": "",
                              "type": ""
                            },
                            "override": false,
                            "prompt": true,
                            "template": {
                              "contiguous": false,
                              "count": "",
                              "height": "",
                              "size": "",
                              "type": "",
                              "units": "",
                              "width": ""
                            }
                          },
                          "type": "attack",
                          "uses": {
                            "recovery": [],
                            "spent": 0
                          }
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
              "x": 200,
              "y": 200,
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
            "createdTime": 1778896417797,
            "modifiedTime": 1778896417833,
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

## GET /get-folder

Get a specific folder by name

**Required scope:** `structure:read`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| name | string | ✓ | body, query | Name of the folder to retrieve |
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The folder information and its contents

### Try It Out

<ApiTester
  method="GET"
  path="/get-folder"
  parameters={[{"name":"name","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/get-folder';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
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
curl -X GET 'http://localhost:3010/get-folder?clientId=fvtt_099ad17ea199e7e3&name=test-folder' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/get-folder'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
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
    clientId: 'fvtt_099ad17ea199e7e3',
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
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
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
  "type": "get-folder-result",
  "requestId": "get-folder_1778896423713",
  "data": {
    "id": "qVncHkNOxAuctXy2",
    "uuid": "Folder.qVncHkNOxAuctXy2",
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

**Required scope:** `structure:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| name | string | ✓ | body, query | Name of the new folder |
| folderType | string | ✓ | body, query | Type of folder (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist) |
| clientId | string |  | query | Client ID for the Foundry world |
| parentFolderId | string |  | body, query | ID of the parent folder (optional for root level) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - The created folder information

### Try It Out

<ApiTester
  method="POST"
  path="/create-folder"
  parameters={[{"name":"name","type":"string","required":true,"source":"body"},{"name":"folderType","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"parentFolderId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/create-folder';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
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
curl -X POST 'http://localhost:3010/create-folder?clientId=fvtt_099ad17ea199e7e3&name=test-folder&folderType=Scene' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/create-folder'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
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
    clientId: 'fvtt_099ad17ea199e7e3',
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
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
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
  "type": "create-folder-result",
  "requestId": "create-folder_1778896423673",
  "data": {
    "id": "qVncHkNOxAuctXy2",
    "uuid": "Folder.qVncHkNOxAuctXy2",
    "name": "test-folder",
    "type": "Scene",
    "parentFolder": null
  }
}
```


---

## DELETE /delete-folder

Delete a folder

**Required scope:** `structure:write`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| folderId | string | ✓ | body, query | ID of the folder to delete |
| clientId | string |  | query | Client ID for the Foundry world |
| deleteAll | boolean |  | body, query | Whether to delete all entities in the folder |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Confirmation of deletion

### Try It Out

<ApiTester
  method="DELETE"
  path="/delete-folder"
  parameters={[{"name":"folderId","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"deleteAll","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/delete-folder';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  folderId: 'qVncHkNOxAuctXy2'
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
curl -X DELETE 'http://localhost:3010/delete-folder?clientId=fvtt_099ad17ea199e7e3&folderId=qVncHkNOxAuctXy2' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/delete-folder'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'folderId': 'qVncHkNOxAuctXy2'
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
    clientId: 'fvtt_099ad17ea199e7e3',
    folderId: 'qVncHkNOxAuctXy2'
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
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤folderId=qVncHkNOxAuctXy2🔤 ➡️ folderId
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
  "type": "delete-folder-result",
  "requestId": "delete-folder_1778896423715",
  "data": {
    "deleted": true,
    "folderId": "qVncHkNOxAuctXy2",
    "entitiesDeleted": 0,
    "foldersDeleted": 1
  }
}
```


---

## GET /contents/:path

This route is deprecated

Use /structure with the path query parameter instead.

**Required scope:** `structure:read`

### Returns

**object** - Error message directing to use /structure endpoint

### Try It Out

<ApiTester
  method="GET"
  path="/contents/:path"
  parameters={[]}
/>

