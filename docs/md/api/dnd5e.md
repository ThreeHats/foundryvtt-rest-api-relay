---
tag: dnd5e
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


# dnd5e

## GET /dnd5e/get-actor-details

Get detailed information for a specific D&D 5e actor. Retrieves comprehensive details about an actor including stats, inventory, spells, features, and other character information based on the requested details array.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor |
| details | array | ✓ | body, query | Array of detail types to retrieve (e.g., ["resources", "items", "spells", "features"]) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Actor details object containing requested information

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/get-actor-details';
const params = {
  clientId: 'your-client-id',
  actorUuid: 'Actor.VKu2l9IdAzxaXrOo',
  details: '["resources","items","features","spells"]'
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
curl -X GET 'http://localhost:3010/dnd5e/get-actor-details?clientId=your-client-id&actorUuid=Actor.VKu2l9IdAzxaXrOo&details=%5B%22resources%22%2C%22items%22%2C%22features%22%2C%22spells%22%5D' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/get-actor-details'
params = {
    'clientId': 'your-client-id',
    'actorUuid': 'Actor.VKu2l9IdAzxaXrOo',
    'details': '["resources","items","features","spells"]'
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
  const path = '/dnd5e/get-actor-details';
  const params = {
    clientId: 'your-client-id',
    actorUuid: 'Actor.VKu2l9IdAzxaXrOo',
    details: '["resources","items","features","spells"]'
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
  🔤/dnd5e/get-actor-details🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤actorUuid=Actor.VKu2l9IdAzxaXrOo🔤 ➡️ actorUuid
  🔤details=["resources","items","features","spells"]🔤 ➡️ details
  🔤?🧲clientId🧲&🧲actorUuid🧲&🧲details🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /dnd5e/get-actor-details🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "requestId": "get-actor-details_1773999676181",
  "clientId": "your-client-id",
  "type": "get-actor-details-result",
  "data": {
    "uuid": "Actor.VKu2l9IdAzxaXrOo",
    "resources": {
      "primary": {
        "value": null,
        "max": null,
        "sr": true,
        "lr": true,
        "label": "Ki"
      },
      "secondary": {
        "value": null,
        "max": null,
        "sr": false,
        "lr": false,
        "label": ""
      },
      "tertiary": {
        "value": null,
        "max": null,
        "sr": false,
        "lr": false,
        "label": ""
      }
    },
    "spells": [
      {
        "name": "test-polymorph",
        "type": "spell",
        "system": {
          "description": {
            "value": "<p>This spell transforms a creature that you can see within range into a new form. An unwilling creature must make a Wisdom saving throw to avoid the effect. The spell has no effect on a Shapechanger or a creature with 0 Hit Points.</p><p>The transformation lasts for the Duration, or until the target drops to 0 Hit Points or dies. The new form can be any beast whose Challenge rating is equal to or less than the target's (or the target's level, if it doesn't have a challenge rating). The target's game statistics, including mental Ability Scores, are replaced by the Statistics of the chosen beast. It retains its Alignment and personality.</p><p>The target assumes the Hit Points of its new form. When it reverts to its normal form, the creature returns to the number of hit points it had before it transformed. If it reverts as a result of dropping to 0 hit points, any excess damage carries over to its normal form. As long as the excess damage doesn't reduce the creature's normal form to 0 hit points, it isn't knocked Unconscious.</p><p>The creature is limited in the Actions it can perform by the Nature of its new form, and it can't speak, cast Spells, or take any other action that requires hands or speech.</p><p>The target's gear melds into the new form. The creature can't activate, use, wield, or otherwise benefit from any of its Equipment.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "activation": {
            "type": "action",
            "condition": "",
            "value": 1
          },
          "duration": {
            "value": "1",
            "units": "hour"
          },
          "target": {
            "affects": {
              "type": "creature",
              "count": "1",
              "choice": false
            },
            "template": {
              "units": "",
              "contiguous": false
            }
          },
          "range": {
            "value": "60",
            "units": "ft"
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "level": 4,
          "school": "trs",
          "materials": {
            "value": "A caterpillar cocoon",
            "consumed": false,
            "cost": 0,
            "supply": 0
          },
          "preparation": {
            "mode": "prepared",
            "prepared": false
          },
          "properties": [
            "vocal",
            "somatic",
            "material",
            "concentration"
          ],
          "activities": {
            "dnd5eactivity000": {
              "_id": "dnd5eactivity000",
              "type": "save",
              "activation": {
                "type": "action",
                "value": null,
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
                "units": "inst",
                "concentration": false,
                "override": false
              },
              "effects": [
                {
                  "_id": "tM8Yj0dw52paxVc4",
                  "onSave": false
                }
              ],
              "range": {
                "override": false,
                "units": "self"
              },
              "target": {
                "prompt": true,
                "template": {
                  "contiguous": false,
                  "units": "ft"
                },
                "affects": {
                  "choice": false
                },
                "override": false
              },
              "damage": {
                "onSave": "half",
                "parts": []
              },
              "save": {
                "ability": [
                  "wis"
                ],
                "dc": {
                  "calculation": "spellcasting",
                  "formula": ""
                }
              },
              "uses": {
                "spent": 0,
                "recovery": []
              },
              "sort": 0
            }
          },
          "identifier": "polymorph"
        },
        "img": "icons/magic/control/energy-stream-link-large-teal.webp",
        "effects": [
          {
            "name": "Polymorph",
            "origin": "Compendium.dnd5e.spells.Item.04nMsTWkIFvkbXlY",
            "duration": {
              "startTime": null,
              "seconds": 3600,
              "combat": null,
              "rounds": null,
              "turns": null,
              "startRound": null,
              "startTurn": null
            },
            "disabled": true,
            "_id": "tM8Yj0dw52paxVc4",
            "changes": [],
            "description": "",
            "transfer": false,
            "statuses": [
              "transformed"
            ],
            "flags": {},
            "tint": "#ffffff",
            "_stats": {
              "compendiumSource": null,
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "img": "icons/magic/control/energy-stream-link-large-teal.webp",
            "type": "base",
            "system": {},
            "sort": 0
          }
        ],
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
          "createdTime": 1773999674489,
          "modifiedTime": 1773999674489,
          "lastModifiedBy": "r6bXhB7k9cXa3cif"
        },
        "_id": "vjFKm5K4EIzqgQXU"
      }
    ],
    "items": [
      {
        "name": "Hammer",
        "type": "loot",
        "img": "icons/tools/hand/hammer-cobbler-steel.webp",
        "system": {
          "description": {
            "value": "<p>A tool with a heavy metal head mounted at the end of its handle, used for jobs such as breaking things and driving in nails. </p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
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
          "identified": true,
          "type": {
            "value": "",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "properties": [],
          "identifier": "hammer"
        },
        "effects": [],
        "folder": "dlru9Hy74nSMv6fr",
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.14pNRT4sZy9rgvhb",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "YJ1P3PnFKHOdQpaP",
        "sort": 0
      },
      {
        "name": "Tinderbox",
        "type": "loot",
        "img": "icons/sundries/lights/torch-black.webp",
        "system": {
          "description": {
            "value": "<p>This small container holds flint, fire steel, and tinder (usually dry cloth soaked in light oil) used to kindle a fire. Using it to light a torch - or anything else with abundant, exposed fuel - takes an action. Lighting any other fire takes 1 minute.</p>\n<p> </p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
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
          "identified": true,
          "container": "8KWz5DJbWUpNWniP",
          "type": {
            "value": "",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "properties": [],
          "identifier": "tinderbox"
        },
        "effects": [],
        "folder": "Dx3K2y0J1wJUPP9m",
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.1FSubnBpSTDmVaYV",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "ownership": {
          "default": 0
        },
        "_id": "DDnxRCeYUhXstWU8",
        "sort": 0
      },
      {
        "name": "Waterskin",
        "type": "consumable",
        "img": "icons/sundries/survival/wetskin-leather-purple.webp",
        "system": {
          "description": {
            "value": "<p>A leather hide sewn into an enclosed skin which can contain up to 4 pints of liquid. It weighs 5 pounds when full; a pint of water is approximately 1 pound.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 1,
          "weight": {
            "value": 5,
            "units": "lb"
          },
          "price": {
            "value": 2,
            "denomination": "sp"
          },
          "attunement": "",
          "equipped": false,
          "rarity": "",
          "identified": true,
          "uses": {
            "max": "4",
            "recovery": [],
            "autoDestroy": false,
            "spent": 0
          },
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
          "container": "6OYR11aJX2dEVtOj",
          "properties": [],
          "type": {
            "value": "food",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "identifier": "waterskin"
        },
        "effects": [],
        "folder": "Dx3K2y0J1wJUPP9m",
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.1L5wkmbw0fmNAr38",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "ownership": {
          "default": 0
        },
        "_id": "5skKSSB4ShHbKoc8",
        "sort": 0
      },
      {
        "name": "Torch",
        "type": "consumable",
        "img": "icons/sundries/lights/torch-black.webp",
        "system": {
          "description": {
            "value": "<p>A torch burns for 1 hour, providing bright light in a 20-foot radius and dim light for an additional 20 feet. If you make a melee attack with a burning torch and hit, it deals 1 fire damage.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 10,
          "weight": {
            "value": 1,
            "units": "lb"
          },
          "price": {
            "value": 1,
            "denomination": "cp"
          },
          "attunement": "",
          "equipped": false,
          "rarity": "",
          "identified": true,
          "uses": {
            "max": "1",
            "recovery": [],
            "autoDestroy": false,
            "spent": 0
          },
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
          "container": "8KWz5DJbWUpNWniP",
          "properties": [],
          "type": {
            "value": "trinket",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "identifier": "torch"
        },
        "effects": [],
        "folder": "Dx3K2y0J1wJUPP9m",
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.29ZLE8PERtFVD3QU",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "ownership": {
          "default": 0
        },
        "_id": "dDuMscUuMI2bTdkj",
        "sort": 0
      },
      {
        "name": "Stick of Incense",
        "type": "loot",
        "img": "icons/consumables/grains/breadsticks-crackers-wrapped-ration-brown.webp",
        "system": {
          "description": {
            "value": "<p>When blocks of incense cannot be used or a cheaper alternative is required, people often use these to perfume the air, whether for pleasurable or religious purposes.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
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
          "identified": true,
          "type": {
            "value": "",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "properties": [],
          "identifier": "stick-of-incense"
        },
        "effects": [],
        "folder": "dlru9Hy74nSMv6fr",
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.3b0RvGi0TnTYpIxn",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "nC6VcR5JAIbR4err",
        "sort": 0
      },
      {
        "name": "Dart",
        "type": "weapon",
        "img": "icons/weapons/thrown/dart-feathered.webp",
        "system": {
          "description": {
            "value": "<p>A small thrown implement crafted with a short wooden shaft and crossed feathres with a sharp wooden or metal tip. Darts can be thrown with sufficient force to puncture the skin.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 10,
          "weight": {
            "value": 0.25,
            "units": "lb"
          },
          "price": {
            "value": 5,
            "denomination": "cp"
          },
          "attunement": "",
          "equipped": true,
          "rarity": "",
          "identified": true,
          "cover": null,
          "range": {
            "value": 20,
            "long": 60,
            "units": "ft",
            "reach": null
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "damage": {
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
            },
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
            }
          },
          "armor": {
            "value": 10
          },
          "hp": {
            "value": 0,
            "max": 0,
            "dt": null,
            "conditions": ""
          },
          "properties": [
            "fin",
            "thr"
          ],
          "proficient": null,
          "type": {
            "value": "simpleR",
            "baseItem": "dart"
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "crewed": false,
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "ammunition": {},
          "mastery": "",
          "identifier": "dart"
        },
        "effects": [],
        "folder": "MLMTCAvKsuFE3vYA",
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "WeKJI3gPUAU52WAX",
        "sort": 0
      },
      {
        "name": "Common Clothes",
        "type": "equipment",
        "img": "icons/equipment/chest/shirt-collared-brown.webp",
        "system": {
          "description": {
            "value": "<p>Clothes worn by most commoners.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 1,
          "weight": {
            "value": 3,
            "units": "lb"
          },
          "price": {
            "value": 5,
            "denomination": "sp"
          },
          "attunement": "",
          "equipped": false,
          "rarity": "",
          "identified": true,
          "cover": null,
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "armor": {
            "value": null,
            "dex": null,
            "magicalBonus": null
          },
          "hp": {
            "value": 0,
            "max": 0,
            "dt": null,
            "conditions": ""
          },
          "speed": {
            "value": null,
            "conditions": ""
          },
          "strength": null,
          "proficient": null,
          "type": {
            "value": "clothing",
            "baseItem": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "crewed": false,
          "properties": [],
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "identifier": "common-clothes"
        },
        "effects": [],
        "folder": "aJgMxnZED9XdoN2W",
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.8RXjiddJ6VGyE7vB",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "1F73YcUHbZMgePAD",
        "sort": 0
      },
      {
        "name": "Rations",
        "type": "consumable",
        "img": "icons/consumables/grains/bread-loaf-boule-rustic-brown.webp",
        "system": {
          "description": {
            "value": "<p>Rations consist of dry foods suitable for extended travel, including jerky, dried fruit, hardtack, and nuts.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 10,
          "weight": {
            "value": 2,
            "units": "lb"
          },
          "price": {
            "value": 5,
            "denomination": "sp"
          },
          "attunement": "",
          "equipped": false,
          "rarity": "",
          "identified": true,
          "uses": {
            "max": "1",
            "recovery": [],
            "autoDestroy": true,
            "spent": 0
          },
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
          "container": "XY8b594Dn7plACLL",
          "properties": [],
          "type": {
            "value": "food",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "identifier": "rations"
        },
        "effects": [],
        "folder": "Dx3K2y0J1wJUPP9m",
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.8d95YV1jHcxPygJ9",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "ownership": {
          "default": 0
        },
        "_id": "LQhGSEatJ3VK7oqW",
        "sort": 0
      },
      {
        "name": "Crowbar",
        "type": "loot",
        "img": "icons/tools/hand/pickaxe-steel-white.webp",
        "system": {
          "description": {
            "value": "<p>Using a crowbar grants advantage to Strength checks where the crowbar's leverage can be applied.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
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
          "identified": true,
          "container": "XY8b594Dn7plACLL",
          "type": {
            "value": "",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "properties": [],
          "identifier": "crowbar"
        },
        "effects": [],
        "folder": "Dx3K2y0J1wJUPP9m",
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.AkyQyonZMVcvOrXU",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "ownership": {
          "default": 0
        },
        "_id": "zF5L4xHnJBC7b2iM",
        "sort": 0
      },
      {
        "name": "Unarmed Strike",
        "type": "weapon",
        "img": "icons/skills/melee/unarmed-punch-fist.webp",
        "system": {
          "description": {
            "value": "<p>A  punch, kick, head-butt, or similar forceful blow (none of which count as weapons). On a hit, an unarmed strike deals bludgeoning damage equal to 1 + your Strength modifier. You are proficient with your unarmed strikes.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 1,
          "weight": {
            "value": 0,
            "units": "lb"
          },
          "price": {
            "value": 0,
            "denomination": "gp"
          },
          "attunement": "",
          "equipped": true,
          "rarity": "",
          "identified": true,
          "cover": null,
          "range": {
            "value": null,
            "long": null,
            "units": "ft",
            "reach": null
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "damage": {
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
            },
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
            }
          },
          "armor": {
            "value": 10
          },
          "hp": {
            "value": 0,
            "max": 0,
            "dt": null,
            "conditions": ""
          },
          "properties": [
            "fin"
          ],
          "proficient": null,
          "type": {
            "value": "simpleM",
            "baseItem": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "crewed": false,
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "ammunition": {},
          "identifier": "unarmed-strike",
          "mastery": ""
        },
        "effects": [],
        "folder": "MLMTCAvKsuFE3vYA",
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
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.GsuvwoekKZatfKwF",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "g7U3OAXVcoI4lwzf",
        "sort": 0
      },
      {
        "name": "Piton",
        "type": "consumable",
        "img": "icons/tools/fasteners/nail-steel.webp",
        "system": {
          "description": {
            "value": "<p>A metal spike that is drive into a seam in a climbing surface with a climbing hammer. It can also be used like iron spikes to block doors in a pinch.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 10,
          "weight": {
            "value": 0.25,
            "units": "lb"
          },
          "price": {
            "value": 5,
            "denomination": "cp"
          },
          "attunement": "",
          "equipped": false,
          "rarity": "",
          "identified": true,
          "uses": {
            "max": "1",
            "recovery": [],
            "autoDestroy": false,
            "spent": 0
          },
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
          "container": "xsB7Y2WI476kvOt4",
          "properties": [],
          "type": {
            "value": "trinket",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "identifier": "piton"
        },
        "effects": [],
        "folder": "Dx3K2y0J1wJUPP9m",
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.P31t6tGgt9aLAdYt",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "ownership": {
          "default": 0
        },
        "_id": "8UQanAvHUIHZXp0O",
        "sort": 0
      },
      {
        "name": "Hempen Rope (50 ft.)",
        "type": "consumable",
        "img": "icons/sundries/survival/rope-wrapped-brown.webp",
        "system": {
          "description": {
            "value": "<p>Rope, whether made of hemp or silk, has 2 hit points and can be burst with a DC 17 Strength check.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 1,
          "weight": {
            "value": 10,
            "units": "lb"
          },
          "price": {
            "value": 1,
            "denomination": "gp"
          },
          "attunement": "",
          "equipped": false,
          "rarity": "",
          "identified": true,
          "uses": {
            "max": "1",
            "recovery": [],
            "autoDestroy": false,
            "spent": 0
          },
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
          "type": {
            "value": "trinket",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "properties": [],
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "identifier": "hempen-rope-50-ft"
        },
        "effects": [],
        "folder": "UnUwTG4YIgd0kaUJ",
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
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.QXmaarJ4X8P0C1HV",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "irCoxLHje0eA6Qmu",
        "sort": 0
      },
      {
        "name": "Reliquary",
        "type": "equipment",
        "img": "icons/containers/chest/chest-reinforced-steel-red.webp",
        "system": {
          "description": {
            "value": "<p>A tiny box or other container holding a fragment of a precious relic, saint, or other historical figure that dedicated their life to walk the path of a true believer. A deity imbues the bearer of this artifact with the ability to call forth their power and in doing so spread the faith once more.</p>\n<p><strong>Spellcasting Focus</strong>. A cleric or paladin can use a holy symbol as a spellcasting focus. To use the symbol in this way, the caster must hold it in hand, wear it visibly, or bear it on a shield.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 1,
          "weight": {
            "value": 2,
            "units": "lb"
          },
          "price": {
            "value": 5,
            "denomination": "gp"
          },
          "attunement": "",
          "equipped": false,
          "rarity": "",
          "identified": true,
          "cover": null,
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "armor": {
            "value": null,
            "dex": null,
            "magicalBonus": null
          },
          "hp": {
            "value": 0,
            "max": 0,
            "dt": null,
            "conditions": ""
          },
          "speed": {
            "value": null,
            "conditions": ""
          },
          "strength": null,
          "proficient": null,
          "type": {
            "value": "trinket",
            "baseItem": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "crewed": false,
          "properties": [
            "foc"
          ],
          "activities": {},
          "attuned": false,
          "identifier": "reliquary"
        },
        "effects": [],
        "folder": "xedn1r43VWuEBcli",
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
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.gP1URGq3kVIIFHJ7",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "GfKuJYLFfw00oW9R",
        "sort": 0
      },
      {
        "name": "Vestments",
        "type": "equipment",
        "img": "icons/equipment/back/mantle-collared-black.webp",
        "system": {
          "description": {
            "value": "<p>Simple or ostentacious wear, often used by priests and other religious figures for use in rituals and ceremonies.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 1,
          "weight": {
            "value": 4,
            "units": "lb"
          },
          "price": {
            "value": 1,
            "denomination": "gp"
          },
          "attunement": "",
          "equipped": true,
          "rarity": "",
          "identified": true,
          "cover": null,
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "armor": {
            "value": null,
            "dex": null,
            "magicalBonus": null
          },
          "hp": {
            "value": 0,
            "max": 0,
            "dt": null,
            "conditions": ""
          },
          "speed": {
            "value": null,
            "conditions": ""
          },
          "strength": null,
          "proficient": null,
          "type": {
            "value": "clothing",
            "baseItem": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "crewed": false,
          "properties": [],
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "identifier": "vestments"
        },
        "effects": [],
        "folder": "aJgMxnZED9XdoN2W",
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.irtqrzaUCeshmTZp",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "NsNcJBUG5ajbk7sG",
        "sort": 0
      },
      {
        "name": "Prayer Book",
        "type": "loot",
        "img": "icons/sundries/books/book-purple-cross.webp",
        "system": {
          "description": {
            "value": "<p>A book containing prayers and incantations dedicated to a specific power for the faithful to follow.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
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
          "identified": true,
          "type": {
            "value": "",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "properties": [],
          "identifier": "prayer-book"
        },
        "effects": [],
        "folder": "dlru9Hy74nSMv6fr",
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.odV5cq2HSLSCH69k",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "XbF0CTipyqbMKHsB",
        "sort": 0
      },
      {
        "name": "Shortsword",
        "type": "weapon",
        "img": "icons/weapons/swords/sword-guard-worn-purple.webp",
        "system": {
          "description": {
            "value": "<p>A medium sized blade with a firm crossguard and a leather-wrapped handle. A versatile weapon which makes up in versatility what it lacks in reach.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 1,
          "weight": {
            "value": 2,
            "units": "lb"
          },
          "price": {
            "value": 10,
            "denomination": "gp"
          },
          "attunement": "",
          "equipped": true,
          "rarity": "",
          "identified": true,
          "cover": null,
          "range": {
            "value": null,
            "long": null,
            "units": "ft",
            "reach": null
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "damage": {
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
            },
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
            }
          },
          "armor": {
            "value": 10
          },
          "hp": {
            "value": 0,
            "max": 0,
            "dt": null,
            "conditions": ""
          },
          "properties": [
            "fin",
            "lgt"
          ],
          "proficient": null,
          "type": {
            "value": "martialM",
            "baseItem": "shortsword"
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "crewed": false,
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "ammunition": {},
          "identifier": "shortsword"
        },
        "effects": [],
        "folder": "MLMTCAvKsuFE3vYA",
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "GYyovoNdU66nxOlX",
        "sort": 0
      },
      {
        "name": "test-potion of storm giant strength",
        "type": "consumable",
        "img": "icons/consumables/potions/bottle-bulb-corked-labeled-blue.webp",
        "system": {
          "description": {
            "value": "<p><em>This potion's transparent liquid has floating in it a sliver of fingernail from a giant of the appropriate type.</em></p>\n<p>When you drink this potion, your Strength score changes to 29 for 1 hour.  The potion has no effect on you if your Strength is equal to or greater than that score.</p>\n<p> </p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "quantity": 1,
          "weight": {
            "value": 0.1,
            "units": "lb"
          },
          "price": {
            "value": 2000,
            "denomination": "gp"
          },
          "attunement": "",
          "equipped": false,
          "rarity": "legendary",
          "identified": true,
          "uses": {
            "max": "1",
            "recovery": [],
            "autoDestroy": true,
            "spent": 0
          },
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
          "type": {
            "value": "potion",
            "subtype": ""
          },
          "unidentified": {
            "description": ""
          },
          "container": null,
          "properties": [
            "mgc"
          ],
          "magicalBonus": null,
          "activities": {
            "dnd5eactivity000": {
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
                "value": "1",
                "units": "hour",
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
          },
          "attuned": false,
          "identifier": "potion-of-storm-giant-strength"
        },
        "effects": [],
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
          "createdTime": 1773999675965,
          "modifiedTime": 1773999675965,
          "lastModifiedBy": "r6bXhB7k9cXa3cif"
        },
        "_id": "URKx2sRVZ9YAyXuv"
      }
    ],
    "features": [
      {
        "name": "Priest",
        "type": "background",
        "system": {
          "description": {
            "value": "<ul><li><strong>Skill Proficiencies:</strong> Insight, Religion</li><li><strong>Languages:</strong> Two of your choice</li><li><strong>Equipment:</strong> A holy symbol, 5 sticks of incense, prayer book, vestments, a set of common clothes, and a pouch containing 15 gp.</li></ul>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "advancement": [],
          "startingEquipment": [],
          "identifier": "priest"
        },
        "img": "icons/sundries/documents/document-torn-diagram-tan.webp",
        "effects": [],
        "folder": null,
        "sort": 0,
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "compendiumSource": null,
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "_id": "q4tr1vTU8RxtU1UZ"
      },
      {
        "_id": "FtOM4QiOW5MwgcS3",
        "name": "Lucky",
        "ownership": {
          "default": 0
        },
        "type": "feat",
        "system": {
          "description": {
            "value": "<p>When you roll a 1 on the d20 for an attack roll, ability check, or saving throw, you can reroll the die and must use the new roll.</p><section class=\"secret foundry-note\" id=\"secret-S04TPyvUh05Dz0Ng\"><p><strong>Foundry Note</strong></p><p>This property can be enabled on your character sheet in the Special Traits configuration on the Attributes tab.</p></section>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "uses": {
            "max": "",
            "spent": 0,
            "recovery": []
          },
          "type": {
            "value": "race",
            "subtype": ""
          },
          "requirements": "Halfling",
          "properties": [],
          "activities": {},
          "enchant": {},
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "identifier": "lucky",
          "advancement": [],
          "crewed": false
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
        "img": "icons/sundries/gaming/dice-runed-brown.webp",
        "effects": [],
        "folder": "kbtbKofcv13crhke",
        "sort": 0,
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        }
      },
      {
        "_id": "nmmihiqphHjoE8dl",
        "name": "Brave",
        "ownership": {
          "default": 0
        },
        "type": "feat",
        "system": {
          "description": {
            "value": "<p>You have advantage on saving throws against being frightened.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "uses": {
            "max": "",
            "spent": 0,
            "recovery": []
          },
          "type": {
            "value": "race",
            "subtype": ""
          },
          "requirements": "Halfling",
          "properties": [],
          "activities": {},
          "enchant": {},
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "identifier": "brave",
          "advancement": [],
          "crewed": false
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
        "img": "icons/skills/melee/unarmed-punch-fist.webp",
        "effects": [],
        "folder": "kbtbKofcv13crhke",
        "sort": 0,
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        }
      },
      {
        "_id": "cWrETHzCRs1Ueqd3",
        "name": "Halfling Nimbleness",
        "ownership": {
          "default": 0
        },
        "type": "feat",
        "system": {
          "description": {
            "value": "<p>You can move through the space of any creature that is of a size larger than yours.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "uses": {
            "max": "",
            "spent": 0,
            "recovery": []
          },
          "type": {
            "value": "race",
            "subtype": ""
          },
          "requirements": "Halfling",
          "properties": [],
          "activities": {},
          "enchant": {},
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "identifier": "halfling-nimbleness",
          "advancement": [],
          "crewed": false
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
        "img": "icons/skills/movement/feet-winged-boots-brown.webp",
        "effects": [],
        "folder": "kbtbKofcv13crhke",
        "sort": 0,
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        }
      },
      {
        "_id": "AArhiOrSkaQUnCZS",
        "name": "Naturally Stealthy",
        "ownership": {
          "default": 0
        },
        "type": "feat",
        "system": {
          "description": {
            "value": "<p>You can attempt to hide even when you are obscured only by a creature that is at least one size larger than you.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "uses": {
            "max": "",
            "spent": 0,
            "recovery": []
          },
          "type": {
            "value": "race",
            "subtype": ""
          },
          "requirements": "Lightfoot Halfling",
          "properties": [],
          "activities": {},
          "enchant": {},
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "identifier": "naturally-stealthy",
          "advancement": [],
          "crewed": false
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
        "img": "icons/magic/perception/silhouette-stealth-shadow.webp",
        "effects": [],
        "folder": "kbtbKofcv13crhke",
        "sort": 0,
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
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
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "identifier": "monk",
          "levels": 1,
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
          },
          "startingEquipment": [
            {
              "type": "OR",
              "_id": "5yj0P4r9teJDnDtd",
              "group": "",
              "sort": 100000,
              "requiresProficiency": false
            },
            {
              "type": "linked",
              "count": null,
              "key": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
              "_id": "R5tuRtaPonfjQCVU",
              "group": "5yj0P4r9teJDnDtd",
              "sort": 200000,
              "requiresProficiency": false
            },
            {
              "type": "weapon",
              "count": null,
              "key": "simpleM",
              "_id": "Mlf6kel8ws6xgDER",
              "group": "5yj0P4r9teJDnDtd",
              "sort": 300000,
              "requiresProficiency": false
            },
            {
              "type": "OR",
              "_id": "3TbVLmLPtjVaSh5O",
              "group": "",
              "sort": 400000,
              "requiresProficiency": false
            },
            {
              "type": "linked",
              "count": null,
              "key": "Compendium.dnd5e.items.Item.XY8b594Dn7plACLL",
              "_id": "AvDYtl0uvQsDuhnb",
              "group": "3TbVLmLPtjVaSh5O",
              "sort": 500000,
              "requiresProficiency": false
            },
            {
              "type": "linked",
              "count": null,
              "key": "Compendium.dnd5e.items.Item.8KWz5DJbWUpNWniP",
              "_id": "4QKQURCmIurbTAzp",
              "group": "3TbVLmLPtjVaSh5O",
              "sort": 600000,
              "requiresProficiency": false
            },
            {
              "type": "linked",
              "count": 10,
              "key": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
              "_id": "AOYuulsULvsHbSLO",
              "group": "",
              "sort": 700000,
              "requiresProficiency": false
            }
          ],
          "wealth": "5d4",
          "primaryAbility": {
            "value": [],
            "all": true
          },
          "hd": {
            "denomination": "d8",
            "spent": 0,
            "additional": ""
          }
        },
        "effects": [],
        "folder": "HQ1Oy7HkbnxnE63o",
        "sort": 0,
        "ownership": {
          "default": 0
        },
        "flags": {},
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.classes.Item.6VoZrWxhOEKGYhnq",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        }
      },
      {
        "_id": "CwgoTDXWCD7PknIN",
        "name": "Unarmored Defense",
        "ownership": {
          "default": 0
        },
        "type": "feat",
        "system": {
          "description": {
            "value": "<p>Beginning at 1st Level, while you are wearing no armor and not wielding a Shield, your AC equals 10 + your Dexterity modifier + your Wisdom modifier.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "uses": {
            "max": "",
            "spent": 0,
            "recovery": []
          },
          "type": {
            "value": "class",
            "subtype": ""
          },
          "requirements": "Monk 1",
          "properties": [],
          "activities": {},
          "enchant": {},
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "identifier": "unarmored-defense",
          "advancement": [],
          "crewed": false
        },
        "flags": {
          "dnd5e": {
            "sourceId": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
            "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK"
          }
        },
        "img": "icons/magic/control/silhouette-hold-change-blue.webp",
        "effects": [
          {
            "_id": "R5ro4AuNjcdWD56O",
            "changes": [
              {
                "key": "system.attributes.ac.calc",
                "mode": 5,
                "value": "unarmoredMonk",
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
            "origin": "Item.cOdcNWy4hII029DT",
            "transfer": true,
            "flags": {},
            "tint": "#ffffff",
            "name": "Unarmored Defense",
            "description": "",
            "statuses": [],
            "_stats": {
              "compendiumSource": null,
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "img": "icons/magic/control/silhouette-hold-change-blue.webp",
            "type": "base",
            "system": {},
            "sort": 0
          }
        ],
        "folder": "TMmNG8ujFDBEWXRe",
        "sort": 0,
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        }
      },
      {
        "_id": "pchnXqd5C79fVlxy",
        "name": "Martial Arts",
        "ownership": {
          "default": 0
        },
        "type": "feat",
        "system": {
          "description": {
            "value": "<p>At 1st level, your practice of martial arts gives you mastery of combat styles that use and monk weapons, which are shortswords and any simple melee weapons that don't have the two-handed or heavy property. You gain the following benefits while you are unarmed or wielding only monk weapons and you aren't wearing armor or wielding a shield:</p>\n<ul>\n<li>\n<p>You can use Dexterity instead of Strength for the attack and damage rolls of your unarmed strikes and monk weapons.</p>\n</li>\n<li>\n<p>You can roll a d4 in place of the normal damage of your unarmed strike or monk weapon. This die changes as you gain monk levels, as shown in the Martial Arts column of the Monk table.</p>\n</li>\n<li>\n<p>When you use the Attack action with an unarmed strike or a monk weapon on your turn, you can make one unarmed strike as a bonus action. For example, if you take the Attack action and attack with a quarterstaff, you can also make an unarmed strike as a bonus action, assuming you haven't already taken a bonus action this turn.</p>\n</li>\n</ul>\n<p>Certain monasteries use specialized forms of the monk weapons. For example, you might use a club that is two lengths of wood connected by a short chain (called a nunchaku) or a sickle with a shorter, straighter blade (called a kama). Whatever name you use for a monk weapon, you can use the game statistics provided for the weapon.</p>",
            "chat": ""
          },
          "source": {
            "custom": "",
            "book": "SRD 5.1",
            "page": "",
            "license": "CC-BY-4.0",
            "rules": "2014",
            "revision": 1
          },
          "uses": {
            "max": "",
            "spent": 0,
            "recovery": []
          },
          "type": {
            "value": "class",
            "subtype": ""
          },
          "requirements": "Monk 1",
          "properties": [],
          "activities": {},
          "enchant": {},
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "identifier": "martial-arts",
          "advancement": [],
          "crewed": false
        },
        "flags": {
          "dnd5e": {
            "sourceId": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw",
            "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK"
          }
        },
        "img": "icons/skills/melee/unarmed-punch-fist.webp",
        "effects": [],
        "folder": "TMmNG8ujFDBEWXRe",
        "sort": 0,
        "_stats": {
          "duplicateSource": null,
          "compendiumSource": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw",
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        }
      }
    ]
  }
}
```


---

## POST /dnd5e/modify-item-charges

Modify the charges for a specific item owned by an actor. Increases or decreases the charges/uses of an item in an actor's inventory. Useful for consumable items like potions, scrolls, or charged magic items.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor who owns the item |
| amount | number | ✓ | body, query | The amount to modify charges by (positive or negative) |
| itemUuid | string |  | body, query | The UUID of the specific item (optional if itemName provided) |
| itemName | string |  | body, query | The name of the item if UUID not provided (optional if itemUuid provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the charge modification operation

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/modify-item-charges';
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "itemName": "Waterskin",
      "amount": -1
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/modify-item-charges?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","itemName":"Waterskin","amount":-1}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/modify-item-charges'
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "itemName": "Waterskin",
      "amount": -1
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
  const path = '/dnd5e/modify-item-charges';
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
        "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
        "itemName": "Waterskin",
        "amount": -1
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
  🔤/dnd5e/modify-item-charges🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","itemName":"Waterskin","amount":-1}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/modify-item-charges🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 73❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "modify-item-charges_1773999677342",
  "clientId": "your-client-id",
  "type": "modify-item-charges-result",
  "data": {
    "itemUuid": "Actor.VKu2l9IdAzxaXrOo.Item.5skKSSB4ShHbKoc8",
    "oldCharges": 4,
    "newCharges": 3
  }
}
```


---

## POST /dnd5e/use-ability

Use a general ability for an actor. Triggers the use of any ability, feature, spell, or item for an actor. This is a generic endpoint that can handle various types of abilities.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor using the ability |
| abilityUuid | string |  | body, query | The UUID of the specific ability (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the ability (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the ability use operation

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/use-ability';
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "abilityName": "Hammer"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/use-ability?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","abilityName":"Hammer"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/use-ability'
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "abilityName": "Hammer"
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
  const path = '/dnd5e/use-ability';
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
        "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
        "abilityName": "Hammer"
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
  🔤/dnd5e/use-ability🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","abilityName":"Hammer"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/use-ability🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 61❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "use-ability_1773999679339",
  "clientId": "your-client-id",
  "type": "use-ability-result",
  "data": {
    "uuid": "Actor.VKu2l9IdAzxaXrOo",
    "ability": "Hammer",
    "result": "DtvA1njIxZtLYeZN"
  }
}
```


---

## POST /dnd5e/use-feature

Use a class or racial feature for an actor. Activates class features (like Action Surge, Rage) or racial features (like Dragonborn Breath Weapon) for a character.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor using the feature |
| abilityUuid | string |  | body, query | The UUID of the specific feature (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the feature if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the feature (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the feature use operation

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/use-feature';
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "abilityName": "Priest"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/use-feature?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","abilityName":"Priest"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/use-feature'
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "abilityName": "Priest"
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
  const path = '/dnd5e/use-feature';
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
        "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
        "abilityName": "Priest"
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
  🔤/dnd5e/use-feature🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","abilityName":"Priest"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/use-feature🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 61❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "use-feature_1773999678768",
  "clientId": "your-client-id",
  "type": "use-feature-result",
  "data": {
    "uuid": "Actor.VKu2l9IdAzxaXrOo",
    "ability": "Priest",
    "result": "y7PTe3bV7FuIAfTd"
  }
}
```


---

## POST /dnd5e/use-spell

Cast a spell for an actor. Casts a spell from the actor's spell list, consuming spell slots as appropriate. Handles cantrips, leveled spells, and spell-like abilities.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor casting the spell |
| abilityUuid | string |  | body, query | The UUID of the specific spell (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the spell if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the spell (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the spell casting operation

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/use-spell';
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "abilityName": "test-polymorph"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/use-spell?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","abilityName":"test-polymorph"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/use-spell'
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "abilityName": "test-polymorph"
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
  const path = '/dnd5e/use-spell';
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
        "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
        "abilityName": "test-polymorph"
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
  🔤/dnd5e/use-spell🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","abilityName":"test-polymorph"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/use-spell🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 69❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "use-spell_1773999679047",
  "clientId": "your-client-id",
  "type": "use-spell-result",
  "data": {
    "uuid": "Actor.VKu2l9IdAzxaXrOo",
    "ability": "test-polymorph",
    "result": null
  }
}
```


---

## POST /dnd5e/use-item

Use an item for an actor. Activates an item from the actor's inventory, such as drinking a potion, using a magic item, or activating equipment with special properties.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| actorUuid | string | ✓ | body, query | UUID of the actor using the item |
| abilityUuid | string |  | body, query | The UUID of the specific item (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the item if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the item (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the item use operation

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/use-item';
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "abilityName": "Hammer"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/use-item?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","abilityName":"Hammer"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/use-item'
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "abilityName": "Hammer"
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
  const path = '/dnd5e/use-item';
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
        "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
        "abilityName": "Hammer"
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
  🔤/dnd5e/use-item🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","abilityName":"Hammer"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/use-item🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 61❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "use-item_1773999678225",
  "clientId": "your-client-id",
  "type": "use-item-result",
  "data": {
    "uuid": "Actor.VKu2l9IdAzxaXrOo",
    "ability": "Hammer",
    "result": "oHpBypRys0oASH7h"
  }
}
```


---

## POST /dnd5e/modify-experience

Modify the experience points for a specific actor. Adds or removes experience points from an actor.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string | ✓ | body, query | Client ID for the Foundry world |
| amount | number | ✓ | body, query | The amount of experience to add (can be negative) |
| actorUuid | string |  | body, query | UUID of the actor to modify |
| selected | boolean |  | body, query | Modify the selected token's actor |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the experience modification operation

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/modify-experience';
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "amount": 100
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/modify-experience?clientId=your-client-id' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","amount":100}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/modify-experience'
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
      "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
      "amount": 100
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
  const path = '/dnd5e/modify-experience';
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
        "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
        "amount": 100
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
  🔤/dnd5e/modify-experience🔤 ➡️ path

  💭 Query parameters
  🔤clientId=your-client-id🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.VKu2l9IdAzxaXrOo","amount":100}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/modify-experience🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 51❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "requestId": "modify-experience_1773999676464",
  "clientId": "your-client-id",
  "type": "modify-experience-result",
  "data": {
    "actorUuid": "Actor.VKu2l9IdAzxaXrOo",
    "oldXp": 0,
    "newXp": 100
  }
}
```


