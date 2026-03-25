---
tag: dnd5e
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# dnd5e

## GET /dnd5e/get-actor-details

Get detailed information for a specific D&D 5e actor. Retrieves comprehensive details about an actor including stats, inventory, spells, features, and other character information based on the requested details array.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| actorUuid | string | ✓ | body, query | UUID of the actor |
| details | array | ✓ | body, query | Array of detail types to retrieve (e.g., ["resources", "items", "spells", "features"]) |
| clientId | string |  | body, query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Actor details object containing requested information

### Try It Out

<ApiTester
  method="GET"
  path="/dnd5e/get-actor-details"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"details","type":"array","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/get-actor-details';
const params = {
  clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
  actorUuid: 'Actor.ioZexonJDGVuU8zl',
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
curl -X GET 'http://localhost:3010/dnd5e/get-actor-details?clientId=foundry-testing-r6bXhB7k9cXa3cif&actorUuid=Actor.ioZexonJDGVuU8zl&details=%5B%22resources%22%2C%22items%22%2C%22features%22%2C%22spells%22%5D' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/get-actor-details'
params = {
    'clientId': 'foundry-testing-r6bXhB7k9cXa3cif',
    'actorUuid': 'Actor.ioZexonJDGVuU8zl',
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
    clientId: 'foundry-testing-r6bXhB7k9cXa3cif',
    actorUuid: 'Actor.ioZexonJDGVuU8zl',
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤actorUuid=Actor.ioZexonJDGVuU8zl🔤 ➡️ actorUuid
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
  "type": "get-actor-details-result",
  "requestId": "get-actor-details_1774367605072",
  "data": {
    "uuid": "Actor.ioZexonJDGVuU8zl",
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
        "effects": [
          {
            "_id": "tM8Yj0dw52paxVc4",
            "_stats": {
              "compendiumSource": null,
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "changes": [],
            "description": "",
            "disabled": true,
            "duration": {
              "combat": null,
              "rounds": null,
              "seconds": 3600,
              "startRound": null,
              "startTime": null,
              "startTurn": null,
              "turns": null
            },
            "flags": {},
            "img": "icons/magic/control/energy-stream-link-large-teal.webp",
            "name": "Polymorph",
            "origin": "Compendium.dnd5e.spells.Item.04nMsTWkIFvkbXlY",
            "sort": 0,
            "statuses": [
              "transformed"
            ],
            "system": {},
            "tint": "#ffffff",
            "transfer": false,
            "type": "base"
          }
        ],
        "img": "icons/magic/control/energy-stream-link-large-teal.webp",
        "name": "test-polymorph",
        "system": {
          "activation": {
            "condition": "",
            "type": "action",
            "value": 1
          },
          "activities": {
            "dnd5eactivity000": {
              "_id": "dnd5eactivity000",
              "activation": {
                "override": false,
                "type": "action",
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
              "damage": {
                "onSave": "half",
                "parts": []
              },
              "description": {
                "chatFlavor": ""
              },
              "duration": {
                "concentration": false,
                "override": false,
                "units": "inst"
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
              "save": {
                "ability": [
                  "wis"
                ],
                "dc": {
                  "calculation": "spellcasting",
                  "formula": ""
                }
              },
              "sort": 0,
              "target": {
                "affects": {
                  "choice": false
                },
                "override": false,
                "prompt": true,
                "template": {
                  "contiguous": false,
                  "units": "ft"
                }
              },
              "type": "save",
              "uses": {
                "recovery": [],
                "spent": 0
              }
            }
          },
          "description": {
            "chat": "",
            "value": "<p>This spell transforms a creature that you can see within range into a new form. An unwilling creature must make a Wisdom saving throw to avoid the effect. The spell has no effect on a Shapechanger or a creature with 0 Hit Points.</p><p>The transformation lasts for the Duration, or until the target drops to 0 Hit Points or dies. The new form can be any beast whose Challenge rating is equal to or less than the target's (or the target's level, if it doesn't have a challenge rating). The target's game statistics, including mental Ability Scores, are replaced by the Statistics of the chosen beast. It retains its Alignment and personality.</p><p>The target assumes the Hit Points of its new form. When it reverts to its normal form, the creature returns to the number of hit points it had before it transformed. If it reverts as a result of dropping to 0 hit points, any excess damage carries over to its normal form. As long as the excess damage doesn't reduce the creature's normal form to 0 hit points, it isn't knocked Unconscious.</p><p>The creature is limited in the Actions it can perform by the Nature of its new form, and it can't speak, cast Spells, or take any other action that requires hands or speech.</p><p>The target's gear melds into the new form. The creature can't activate, use, wield, or otherwise benefit from any of its Equipment.</p>"
          },
          "duration": {
            "units": "hour",
            "value": "1"
          },
          "identifier": "polymorph",
          "level": 4,
          "materials": {
            "consumed": false,
            "cost": 0,
            "supply": 0,
            "value": "A caterpillar cocoon"
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
          "range": {
            "units": "ft",
            "value": "60"
          },
          "school": "trs",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "target": {
            "affects": {
              "choice": false,
              "count": "1",
              "type": "creature"
            },
            "template": {
              "contiguous": false,
              "units": ""
            }
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          }
        },
        "type": "spell",
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
          "createdTime": 1774367604775,
          "modifiedTime": 1774367604775,
          "lastModifiedBy": "r6bXhB7k9cXa3cif"
        },
        "_id": "cuN0zv7OrhLW1LKO"
      }
    ],
    "items": [
      {
        "_id": "YJ1P3PnFKHOdQpaP",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.14pNRT4sZy9rgvhb",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "dlru9Hy74nSMv6fr",
        "img": "icons/tools/hand/hammer-cobbler-steel.webp",
        "name": "Hammer",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "container": null,
          "description": {
            "chat": "",
            "value": "<p>A tool with a heavy metal head mounted at the end of its handle, used for jobs such as breaking things and driving in nails. </p>"
          },
          "identified": true,
          "identifier": "hammer",
          "price": {
            "denomination": "gp",
            "value": 1
          },
          "properties": [],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": ""
          },
          "unidentified": {
            "description": ""
          },
          "weight": {
            "units": "lb",
            "value": 3
          }
        },
        "type": "loot"
      },
      {
        "_id": "DDnxRCeYUhXstWU8",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.1FSubnBpSTDmVaYV",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "Dx3K2y0J1wJUPP9m",
        "img": "icons/sundries/lights/torch-black.webp",
        "name": "Tinderbox",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "container": "8KWz5DJbWUpNWniP",
          "description": {
            "chat": "",
            "value": "<p>This small container holds flint, fire steel, and tinder (usually dry cloth soaked in light oil) used to kindle a fire. Using it to light a torch - or anything else with abundant, exposed fuel - takes an action. Lighting any other fire takes 1 minute.</p>\n<p> </p>"
          },
          "identified": true,
          "identifier": "tinderbox",
          "price": {
            "denomination": "sp",
            "value": 5
          },
          "properties": [],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": ""
          },
          "unidentified": {
            "description": ""
          },
          "weight": {
            "units": "lb",
            "value": 1
          }
        },
        "type": "loot"
      },
      {
        "_id": "5skKSSB4ShHbKoc8",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.1L5wkmbw0fmNAr38",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "Dx3K2y0J1wJUPP9m",
        "img": "icons/sundries/survival/wetskin-leather-purple.webp",
        "name": "Waterskin",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "attunement": "",
          "container": "6OYR11aJX2dEVtOj",
          "damage": {
            "base": {
              "custom": {
                "enabled": false
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "number": 1
              },
              "types": []
            },
            "replace": false
          },
          "description": {
            "chat": "",
            "value": "<p>A leather hide sewn into an enclosed skin which can contain up to 4 pints of liquid. It weighs 5 pounds when full; a pint of water is approximately 1 pound.</p>"
          },
          "equipped": false,
          "identified": true,
          "identifier": "waterskin",
          "magicalBonus": null,
          "price": {
            "denomination": "sp",
            "value": 2
          },
          "properties": [],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "food"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "autoDestroy": false,
            "max": "4",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 5
          }
        },
        "type": "consumable"
      },
      {
        "_id": "dDuMscUuMI2bTdkj",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.29ZLE8PERtFVD3QU",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "Dx3K2y0J1wJUPP9m",
        "img": "icons/sundries/lights/torch-black.webp",
        "name": "Torch",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "attunement": "",
          "container": "8KWz5DJbWUpNWniP",
          "damage": {
            "base": {
              "bonus": "",
              "custom": {
                "enabled": true,
                "formula": "1"
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "formula": "",
                "mode": "",
                "number": null
              },
              "types": [
                "fire"
              ]
            },
            "replace": false
          },
          "description": {
            "chat": "",
            "value": "<p>A torch burns for 1 hour, providing bright light in a 20-foot radius and dim light for an additional 20 feet. If you make a melee attack with a burning torch and hit, it deals 1 fire damage.</p>"
          },
          "equipped": false,
          "identified": true,
          "identifier": "torch",
          "magicalBonus": null,
          "price": {
            "denomination": "cp",
            "value": 1
          },
          "properties": [],
          "quantity": 10,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "trinket"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "autoDestroy": false,
            "max": "1",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 1
          }
        },
        "type": "consumable"
      },
      {
        "_id": "nC6VcR5JAIbR4err",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.3b0RvGi0TnTYpIxn",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "dlru9Hy74nSMv6fr",
        "img": "icons/consumables/grains/breadsticks-crackers-wrapped-ration-brown.webp",
        "name": "Stick of Incense",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "container": null,
          "description": {
            "chat": "",
            "value": "<p>When blocks of incense cannot be used or a cheaper alternative is required, people often use these to perfume the air, whether for pleasurable or religious purposes.</p>"
          },
          "identified": true,
          "identifier": "stick-of-incense",
          "price": {
            "denomination": "sp",
            "value": 2
          },
          "properties": [],
          "quantity": 5,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": ""
          },
          "unidentified": {
            "description": ""
          },
          "weight": {
            "units": "lb",
            "value": 0
          }
        },
        "type": "loot"
      },
      {
        "_id": "WeKJI3gPUAU52WAX",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "MLMTCAvKsuFE3vYA",
        "img": "icons/weapons/thrown/dart-feathered.webp",
        "name": "Dart",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "ammunition": {},
          "armor": {
            "value": 10
          },
          "attuned": false,
          "attunement": "",
          "container": null,
          "cover": null,
          "crewed": false,
          "damage": {
            "base": {
              "bonus": "",
              "custom": {
                "enabled": false,
                "formula": ""
              },
              "denomination": 4,
              "number": 1,
              "scaling": {
                "formula": "",
                "mode": "",
                "number": null
              },
              "types": [
                "piercing"
              ]
            },
            "versatile": {
              "bonus": "",
              "custom": {
                "enabled": false,
                "formula": ""
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "formula": "",
                "mode": "",
                "number": null
              },
              "types": []
            }
          },
          "description": {
            "chat": "",
            "value": "<p>A small thrown implement crafted with a short wooden shaft and crossed feathres with a sharp wooden or metal tip. Darts can be thrown with sufficient force to puncture the skin.</p>"
          },
          "equipped": true,
          "hp": {
            "conditions": "",
            "dt": null,
            "max": 0,
            "value": 0
          },
          "identified": true,
          "identifier": "dart",
          "magicalBonus": null,
          "mastery": "",
          "price": {
            "denomination": "cp",
            "value": 5
          },
          "proficient": null,
          "properties": [
            "fin",
            "thr"
          ],
          "quantity": 10,
          "range": {
            "long": 60,
            "reach": null,
            "units": "ft",
            "value": 20
          },
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "baseItem": "dart",
            "value": "simpleR"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 0.25
          }
        },
        "type": "weapon"
      },
      {
        "_id": "1F73YcUHbZMgePAD",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.8RXjiddJ6VGyE7vB",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "aJgMxnZED9XdoN2W",
        "img": "icons/equipment/chest/shirt-collared-brown.webp",
        "name": "Common Clothes",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "armor": {
            "dex": null,
            "magicalBonus": null,
            "value": null
          },
          "attuned": false,
          "attunement": "",
          "container": null,
          "cover": null,
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>Clothes worn by most commoners.</p>"
          },
          "equipped": false,
          "hp": {
            "conditions": "",
            "dt": null,
            "max": 0,
            "value": 0
          },
          "identified": true,
          "identifier": "common-clothes",
          "price": {
            "denomination": "sp",
            "value": 5
          },
          "proficient": null,
          "properties": [],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "speed": {
            "conditions": "",
            "value": null
          },
          "strength": null,
          "type": {
            "baseItem": "",
            "value": "clothing"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 3
          }
        },
        "type": "equipment"
      },
      {
        "_id": "LQhGSEatJ3VK7oqW",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.8d95YV1jHcxPygJ9",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "Dx3K2y0J1wJUPP9m",
        "img": "icons/consumables/grains/bread-loaf-boule-rustic-brown.webp",
        "name": "Rations",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "attunement": "",
          "container": "XY8b594Dn7plACLL",
          "damage": {
            "base": {
              "custom": {
                "enabled": false
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "number": 1
              },
              "types": []
            },
            "replace": false
          },
          "description": {
            "chat": "",
            "value": "<p>Rations consist of dry foods suitable for extended travel, including jerky, dried fruit, hardtack, and nuts.</p>"
          },
          "equipped": false,
          "identified": true,
          "identifier": "rations",
          "magicalBonus": null,
          "price": {
            "denomination": "sp",
            "value": 5
          },
          "properties": [],
          "quantity": 10,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "food"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "autoDestroy": true,
            "max": "1",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 2
          }
        },
        "type": "consumable"
      },
      {
        "_id": "zF5L4xHnJBC7b2iM",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.AkyQyonZMVcvOrXU",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "Dx3K2y0J1wJUPP9m",
        "img": "icons/tools/hand/pickaxe-steel-white.webp",
        "name": "Crowbar",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "container": "XY8b594Dn7plACLL",
          "description": {
            "chat": "",
            "value": "<p>Using a crowbar grants advantage to Strength checks where the crowbar's leverage can be applied.</p>"
          },
          "identified": true,
          "identifier": "crowbar",
          "price": {
            "denomination": "gp",
            "value": 2
          },
          "properties": [],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": ""
          },
          "unidentified": {
            "description": ""
          },
          "weight": {
            "units": "lb",
            "value": 5
          }
        },
        "type": "loot"
      },
      {
        "_id": "g7U3OAXVcoI4lwzf",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.GsuvwoekKZatfKwF",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
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
        "folder": "MLMTCAvKsuFE3vYA",
        "img": "icons/skills/melee/unarmed-punch-fist.webp",
        "name": "Unarmed Strike",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "ammunition": {},
          "armor": {
            "value": 10
          },
          "attuned": false,
          "attunement": "",
          "container": null,
          "cover": null,
          "crewed": false,
          "damage": {
            "base": {
              "bonus": "",
              "custom": {
                "enabled": true,
                "formula": "@scale.monk.die"
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "formula": "",
                "mode": "",
                "number": null
              },
              "types": [
                "bludgeoning"
              ]
            },
            "versatile": {
              "bonus": "",
              "custom": {
                "enabled": false,
                "formula": ""
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "formula": "",
                "mode": "",
                "number": null
              },
              "types": []
            }
          },
          "description": {
            "chat": "",
            "value": "<p>A  punch, kick, head-butt, or similar forceful blow (none of which count as weapons). On a hit, an unarmed strike deals bludgeoning damage equal to 1 + your Strength modifier. You are proficient with your unarmed strikes.</p>"
          },
          "equipped": true,
          "hp": {
            "conditions": "",
            "dt": null,
            "max": 0,
            "value": 0
          },
          "identified": true,
          "identifier": "unarmed-strike",
          "magicalBonus": null,
          "mastery": "",
          "price": {
            "denomination": "gp",
            "value": 0
          },
          "proficient": null,
          "properties": [
            "fin"
          ],
          "quantity": 1,
          "range": {
            "long": null,
            "reach": null,
            "units": "ft",
            "value": null
          },
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "baseItem": "",
            "value": "simpleM"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 0
          }
        },
        "type": "weapon"
      },
      {
        "_id": "8UQanAvHUIHZXp0O",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.P31t6tGgt9aLAdYt",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "Dx3K2y0J1wJUPP9m",
        "img": "icons/tools/fasteners/nail-steel.webp",
        "name": "Piton",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "attunement": "",
          "container": "xsB7Y2WI476kvOt4",
          "damage": {
            "base": {
              "custom": {
                "enabled": false
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "number": 1
              },
              "types": []
            },
            "replace": false
          },
          "description": {
            "chat": "",
            "value": "<p>A metal spike that is drive into a seam in a climbing surface with a climbing hammer. It can also be used like iron spikes to block doors in a pinch.</p>"
          },
          "equipped": false,
          "identified": true,
          "identifier": "piton",
          "magicalBonus": null,
          "price": {
            "denomination": "cp",
            "value": 5
          },
          "properties": [],
          "quantity": 10,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "trinket"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "autoDestroy": false,
            "max": "1",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 0.25
          }
        },
        "type": "consumable"
      },
      {
        "_id": "irCoxLHje0eA6Qmu",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.QXmaarJ4X8P0C1HV",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {
          "dnd5e": {
            "riders": {
              "activity": [],
              "effect": []
            }
          }
        },
        "folder": "UnUwTG4YIgd0kaUJ",
        "img": "icons/sundries/survival/rope-wrapped-brown.webp",
        "name": "Hempen Rope (50 ft.)",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "attuned": false,
          "attunement": "",
          "container": null,
          "damage": {
            "base": {
              "custom": {
                "enabled": false
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "number": 1
              },
              "types": []
            },
            "replace": false
          },
          "description": {
            "chat": "",
            "value": "<p>Rope, whether made of hemp or silk, has 2 hit points and can be burst with a DC 17 Strength check.</p>"
          },
          "equipped": false,
          "identified": true,
          "identifier": "hempen-rope-50-ft",
          "magicalBonus": null,
          "price": {
            "denomination": "gp",
            "value": 1
          },
          "properties": [],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "trinket"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "autoDestroy": false,
            "max": "1",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 10
          }
        },
        "type": "consumable"
      },
      {
        "_id": "GfKuJYLFfw00oW9R",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.gP1URGq3kVIIFHJ7",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {
          "dnd5e": {
            "riders": {
              "activity": [],
              "effect": []
            }
          }
        },
        "folder": "xedn1r43VWuEBcli",
        "img": "icons/containers/chest/chest-reinforced-steel-red.webp",
        "name": "Reliquary",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {},
          "armor": {
            "dex": null,
            "magicalBonus": null,
            "value": null
          },
          "attuned": false,
          "attunement": "",
          "container": null,
          "cover": null,
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>A tiny box or other container holding a fragment of a precious relic, saint, or other historical figure that dedicated their life to walk the path of a true believer. A deity imbues the bearer of this artifact with the ability to call forth their power and in doing so spread the faith once more.</p>\n<p><strong>Spellcasting Focus</strong>. A cleric or paladin can use a holy symbol as a spellcasting focus. To use the symbol in this way, the caster must hold it in hand, wear it visibly, or bear it on a shield.</p>"
          },
          "equipped": false,
          "hp": {
            "conditions": "",
            "dt": null,
            "max": 0,
            "value": 0
          },
          "identified": true,
          "identifier": "reliquary",
          "price": {
            "denomination": "gp",
            "value": 5
          },
          "proficient": null,
          "properties": [
            "foc"
          ],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "speed": {
            "conditions": "",
            "value": null
          },
          "strength": null,
          "type": {
            "baseItem": "",
            "value": "trinket"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 2
          }
        },
        "type": "equipment"
      },
      {
        "_id": "NsNcJBUG5ajbk7sG",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.irtqrzaUCeshmTZp",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "aJgMxnZED9XdoN2W",
        "img": "icons/equipment/back/mantle-collared-black.webp",
        "name": "Vestments",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "armor": {
            "dex": null,
            "magicalBonus": null,
            "value": null
          },
          "attuned": false,
          "attunement": "",
          "container": null,
          "cover": null,
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>Simple or ostentacious wear, often used by priests and other religious figures for use in rituals and ceremonies.</p>"
          },
          "equipped": true,
          "hp": {
            "conditions": "",
            "dt": null,
            "max": 0,
            "value": 0
          },
          "identified": true,
          "identifier": "vestments",
          "price": {
            "denomination": "gp",
            "value": 1
          },
          "proficient": null,
          "properties": [],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "speed": {
            "conditions": "",
            "value": null
          },
          "strength": null,
          "type": {
            "baseItem": "",
            "value": "clothing"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 4
          }
        },
        "type": "equipment"
      },
      {
        "_id": "XbF0CTipyqbMKHsB",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.odV5cq2HSLSCH69k",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "dlru9Hy74nSMv6fr",
        "img": "icons/sundries/books/book-purple-cross.webp",
        "name": "Prayer Book",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "container": null,
          "description": {
            "chat": "",
            "value": "<p>A book containing prayers and incantations dedicated to a specific power for the faithful to follow.</p>"
          },
          "identified": true,
          "identifier": "prayer-book",
          "price": {
            "denomination": "gp",
            "value": 25
          },
          "properties": [],
          "quantity": 1,
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": ""
          },
          "unidentified": {
            "description": ""
          },
          "weight": {
            "units": "lb",
            "value": 5
          }
        },
        "type": "loot"
      },
      {
        "_id": "GYyovoNdU66nxOlX",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "MLMTCAvKsuFE3vYA",
        "img": "icons/weapons/swords/sword-guard-worn-purple.webp",
        "name": "Shortsword",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
          },
          "ammunition": {},
          "armor": {
            "value": 10
          },
          "attuned": false,
          "attunement": "",
          "container": null,
          "cover": null,
          "crewed": false,
          "damage": {
            "base": {
              "bonus": "",
              "custom": {
                "enabled": false,
                "formula": ""
              },
              "denomination": 6,
              "number": 1,
              "scaling": {
                "formula": "",
                "mode": "",
                "number": null
              },
              "types": [
                "piercing"
              ]
            },
            "versatile": {
              "bonus": "",
              "custom": {
                "enabled": false,
                "formula": ""
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "formula": "",
                "mode": "",
                "number": null
              },
              "types": []
            }
          },
          "description": {
            "chat": "",
            "value": "<p>A medium sized blade with a firm crossguard and a leather-wrapped handle. A versatile weapon which makes up in versatility what it lacks in reach.</p>"
          },
          "equipped": true,
          "hp": {
            "conditions": "",
            "dt": null,
            "max": 0,
            "value": 0
          },
          "identified": true,
          "identifier": "shortsword",
          "magicalBonus": null,
          "price": {
            "denomination": "gp",
            "value": 10
          },
          "proficient": null,
          "properties": [
            "fin",
            "lgt"
          ],
          "quantity": 1,
          "range": {
            "long": null,
            "reach": null,
            "units": "ft",
            "value": null
          },
          "rarity": "",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "baseItem": "shortsword",
            "value": "martialM"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 2
          }
        },
        "type": "weapon"
      },
      {
        "effects": [],
        "img": "icons/consumables/potions/bottle-bulb-corked-labeled-blue.webp",
        "name": "test-potion of storm giant strength",
        "system": {
          "activities": {
            "dnd5eactivity000": {
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
                "units": "hour",
                "value": "1"
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
          },
          "attuned": false,
          "attunement": "",
          "container": null,
          "damage": {
            "base": {
              "custom": {
                "enabled": false
              },
              "denomination": null,
              "number": null,
              "scaling": {
                "number": 1
              },
              "types": []
            },
            "replace": false
          },
          "description": {
            "chat": "",
            "value": "<p><em>This potion's transparent liquid has floating in it a sliver of fingernail from a giant of the appropriate type.</em></p>\n<p>When you drink this potion, your Strength score changes to 29 for 1 hour.  The potion has no effect on you if your Strength is equal to or greater than that score.</p>\n<p> </p>"
          },
          "equipped": false,
          "identified": true,
          "identifier": "potion-of-storm-giant-strength",
          "magicalBonus": null,
          "price": {
            "denomination": "gp",
            "value": 2000
          },
          "properties": [
            "mgc"
          ],
          "quantity": 1,
          "rarity": "legendary",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "potion"
          },
          "unidentified": {
            "description": ""
          },
          "uses": {
            "autoDestroy": true,
            "max": "1",
            "recovery": [],
            "spent": 0
          },
          "weight": {
            "units": "lb",
            "value": 0.1
          }
        },
        "type": "consumable",
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
          "createdTime": 1774367604969,
          "modifiedTime": 1774367604969,
          "lastModifiedBy": "r6bXhB7k9cXa3cif"
        },
        "_id": "xzMsr0JNmiadK4Cr"
      }
    ],
    "features": [
      {
        "_id": "q4tr1vTU8RxtU1UZ",
        "_stats": {
          "compendiumSource": null,
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": null,
        "img": "icons/sundries/documents/document-torn-diagram-tan.webp",
        "name": "Priest",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "advancement": [],
          "description": {
            "chat": "",
            "value": "<ul><li><strong>Skill Proficiencies:</strong> Insight, Religion</li><li><strong>Languages:</strong> Two of your choice</li><li><strong>Equipment:</strong> A holy symbol, 5 sticks of incense, prayer book, vestments, a set of common clothes, and a pouch containing 15 gp.</li></ul>"
          },
          "identifier": "priest",
          "source": {
            "book": "",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "startingEquipment": []
        },
        "type": "background"
      },
      {
        "_id": "FtOM4QiOW5MwgcS3",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
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
        "folder": "kbtbKofcv13crhke",
        "img": "icons/sundries/gaming/dice-runed-brown.webp",
        "name": "Lucky",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {},
          "advancement": [],
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>When you roll a 1 on the d20 for an attack roll, ability check, or saving throw, you can reroll the die and must use the new roll.</p><section class=\"secret foundry-note\" id=\"secret-S04TPyvUh05Dz0Ng\"><p><strong>Foundry Note</strong></p><p>This property can be enabled on your character sheet in the Special Traits configuration on the Attributes tab.</p></section>"
          },
          "enchant": {},
          "identifier": "lucky",
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "properties": [],
          "requirements": "Halfling",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "race"
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          }
        },
        "type": "feat"
      },
      {
        "_id": "nmmihiqphHjoE8dl",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
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
        "folder": "kbtbKofcv13crhke",
        "img": "icons/skills/melee/unarmed-punch-fist.webp",
        "name": "Brave",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {},
          "advancement": [],
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>You have advantage on saving throws against being frightened.</p>"
          },
          "enchant": {},
          "identifier": "brave",
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "properties": [],
          "requirements": "Halfling",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "race"
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          }
        },
        "type": "feat"
      },
      {
        "_id": "cWrETHzCRs1Ueqd3",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
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
        "folder": "kbtbKofcv13crhke",
        "img": "icons/skills/movement/feet-winged-boots-brown.webp",
        "name": "Halfling Nimbleness",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {},
          "advancement": [],
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>You can move through the space of any creature that is of a size larger than yours.</p>"
          },
          "enchant": {},
          "identifier": "halfling-nimbleness",
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "properties": [],
          "requirements": "Halfling",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "race"
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          }
        },
        "type": "feat"
      },
      {
        "_id": "AArhiOrSkaQUnCZS",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
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
        "folder": "kbtbKofcv13crhke",
        "img": "icons/magic/perception/silhouette-stealth-shadow.webp",
        "name": "Naturally Stealthy",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {},
          "advancement": [],
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>You can attempt to hide even when you are obscured only by a creature that is at least one size larger than you.</p>"
          },
          "enchant": {},
          "identifier": "naturally-stealthy",
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "properties": [],
          "requirements": "Lightfoot Halfling",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "race"
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          }
        },
        "type": "feat"
      },
      {
        "_id": "8Grf7ga6JcZF0X6x",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.classes.Item.6VoZrWxhOEKGYhnq",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {},
        "folder": "HQ1Oy7HkbnxnE63o",
        "img": "icons/skills/melee/hand-grip-staff-blue.webp",
        "name": "Monk",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
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
          "description": {
            "chat": "",
            "value": "<p>As a monk, you gain the following class features.</p><h3>Hit Points</h3><p><strong>Hit Dice:</strong> 1d8 per monk level<br /><strong>Hit Points at 1st Level:</strong> 8 + your Constitution modifier<br /><strong>Hit Points at Higher Levels:</strong> 1d8 (or 5) + your Constitution modifier per monk level after 1st</p><h3>Proficiencies</h3><p><strong>Armor:</strong> None<br /><strong>Weapons:</strong> Simple weapons, shortswords<br /><strong>Tools:</strong> Choose one type of artisan's tools or one musical instrument<br /><strong>Saving Throws:</strong> Strength, Dexterity<br /><strong>Skills:</strong> Choose two from Acrobatics, Athletics, History, Insight, Religion, and Stealth</p><h1>Monk Advancement</h1><table><thead><tr><td>Level</td><td>Proficiency Bonus</td><td>Martial Arts</td><td>Ki Points</td><td>Unarmored Movement</td><td>Features</td></tr></thead><tbody><tr><td>1st</td><td>+2</td><td>1d4</td><td>—</td><td>—</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI]{Unarmored Defense}, @UUID[Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw]{Martial Arts}</td></tr><tr><td>2nd</td><td>+2</td><td>1d4</td><td>2</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.10b6z2W1txNkrGP7]{Ki}, @UUID[Compendium.dnd5e.classfeatures.Item.zCeqyQ8uIPNdYJSW]{Unarmored Movement}</td></tr><tr><td>3rd</td><td>+2</td><td>1d4</td><td>3</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.rtpQdX77dYWbDIOH]{Monastic Tradition}, @UUID[Compendium.dnd5e.classfeatures.Item.mzweVbnsJPQiVkAe]{Deflect Missiles}</td></tr><tr><td>4th</td><td>+2</td><td>1d4</td><td>4</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}, @UUID[Compendium.dnd5e.classfeatures.Item.KQz9bqxVkXjDl8gK]{Slow Fall}</td></tr><tr><td>5th</td><td>+3</td><td>1d6</td><td>5</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.XogoBnFWmCAHXppo]{Extra Attack}, @UUID[Compendium.dnd5e.classfeatures.Item.pvRc6GAu1ok6zihC]{Stunning Strike}</td></tr><tr><td>6th</td><td>+3</td><td>1d6</td><td>6</td><td>+15 ft.</td><td><p>@UUID[Compendium.dnd5e.classfeatures.Item.7flZKruSSu6dHg6D]{Ki-Empowered Strikes},</p><p>Monastic Tradition feature</p></td></tr><tr><td>7th</td><td>+3</td><td>1d6</td><td>7</td><td>+15 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.a4P4DNMmH8CqSNkC]{Evasion}, @UUID[Compendium.dnd5e.classfeatures.Item.ZmC31XKS4YNENnoc]{Stillness of Mind}</td></tr><tr><td>8th</td><td>+3</td><td>1d6</td><td>8</td><td>+15 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>9th</td><td>+4</td><td>1d6</td><td>9</td><td>+15 ft.</td><td>Unarmored Movement improvement</td></tr><tr><td>10th</td><td>+4</td><td>1d6</td><td>10</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.bqWA7t9pDELbNRkp]{Purity of Body}</td></tr><tr><td>11th</td><td>+4</td><td>1d8</td><td>11</td><td>+20 ft.</td><td>Monastic Tradition feature</td></tr><tr><td>12th</td><td>+4</td><td>1d8</td><td>12</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>13th</td><td>+5</td><td>1d8</td><td>13</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.XjuGBeB8Y0C3A5D4]{Tongue of the Sun and Moon}</td></tr><tr><td>14th</td><td>+5</td><td>1d8</td><td>14</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.7D2EkLdISwShEDlN]{Diamond Soul}</td></tr><tr><td>15th</td><td>+5</td><td>1d8</td><td>15</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.gDH8PMrKvLHaNmEI]{Timeless Body}</td></tr><tr><td>16th</td><td>+5</td><td>1d8</td><td>16</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>17th</td><td>+6</td><td>1d10</td><td>17</td><td>+25 ft.</td><td>Monastic Tradition feature</td></tr><tr><td>18th</td><td>+6</td><td>1d10</td><td>18</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.3jwFt3hSqDswBlOH]{Empty Body}</td></tr><tr><td>19th</td><td>+6</td><td>1d10</td><td>19</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>20th</td><td>+6</td><td>1d10</td><td>20</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.mQNPg89YIs7g5tG4]{Perfect Self}</td></tr></tbody></table><h1>Monastic Traditions</h1><p>Three traditions of monastic pursuit are common in the monasteries scattered across the multiverse. Most monasteries practice one tradition exclusively, but a few honor the three traditions and instruct each monk according to his or her aptitude and interest. All three traditions rely on the same basic techniques, diverging as the student grows more adept. Thus, a monk need choose a tradition only upon reaching 3rd level.</p><p>@UUID[Compendium.dnd5e.subclasses.Item.IvlpKMXX3PmW1NY2]{Way of the Open Hand}</p>"
          },
          "hd": {
            "additional": "",
            "denomination": "d8",
            "spent": 0
          },
          "identifier": "monk",
          "levels": 1,
          "primaryAbility": {
            "all": true,
            "value": []
          },
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "spellcasting": {
            "ability": "",
            "preparation": {
              "formula": ""
            },
            "progression": "none"
          },
          "startingEquipment": [
            {
              "_id": "5yj0P4r9teJDnDtd",
              "group": "",
              "requiresProficiency": false,
              "sort": 100000,
              "type": "OR"
            },
            {
              "_id": "R5tuRtaPonfjQCVU",
              "count": null,
              "group": "5yj0P4r9teJDnDtd",
              "key": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
              "requiresProficiency": false,
              "sort": 200000,
              "type": "linked"
            },
            {
              "_id": "Mlf6kel8ws6xgDER",
              "count": null,
              "group": "5yj0P4r9teJDnDtd",
              "key": "simpleM",
              "requiresProficiency": false,
              "sort": 300000,
              "type": "weapon"
            },
            {
              "_id": "3TbVLmLPtjVaSh5O",
              "group": "",
              "requiresProficiency": false,
              "sort": 400000,
              "type": "OR"
            },
            {
              "_id": "AvDYtl0uvQsDuhnb",
              "count": null,
              "group": "3TbVLmLPtjVaSh5O",
              "key": "Compendium.dnd5e.items.Item.XY8b594Dn7plACLL",
              "requiresProficiency": false,
              "sort": 500000,
              "type": "linked"
            },
            {
              "_id": "4QKQURCmIurbTAzp",
              "count": null,
              "group": "3TbVLmLPtjVaSh5O",
              "key": "Compendium.dnd5e.items.Item.8KWz5DJbWUpNWniP",
              "requiresProficiency": false,
              "sort": 600000,
              "type": "linked"
            },
            {
              "_id": "AOYuulsULvsHbSLO",
              "count": 10,
              "group": "",
              "key": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
              "requiresProficiency": false,
              "sort": 700000,
              "type": "linked"
            }
          ],
          "wealth": "5d4"
        },
        "type": "class"
      },
      {
        "_id": "CwgoTDXWCD7PknIN",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [
          {
            "_id": "R5ro4AuNjcdWD56O",
            "_stats": {
              "compendiumSource": null,
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "changes": [
              {
                "key": "system.attributes.ac.calc",
                "mode": 5,
                "priority": null,
                "value": "unarmoredMonk"
              }
            ],
            "description": "",
            "disabled": false,
            "duration": {
              "combat": null,
              "rounds": null,
              "seconds": null,
              "startRound": null,
              "startTime": 0,
              "startTurn": null,
              "turns": null
            },
            "flags": {},
            "img": "icons/magic/control/silhouette-hold-change-blue.webp",
            "name": "Unarmored Defense",
            "origin": "Item.cOdcNWy4hII029DT",
            "sort": 0,
            "statuses": [],
            "system": {},
            "tint": "#ffffff",
            "transfer": true,
            "type": "base"
          }
        ],
        "flags": {
          "dnd5e": {
            "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
            "sourceId": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI"
          }
        },
        "folder": "TMmNG8ujFDBEWXRe",
        "img": "icons/magic/control/silhouette-hold-change-blue.webp",
        "name": "Unarmored Defense",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {},
          "advancement": [],
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>Beginning at 1st Level, while you are wearing no armor and not wielding a Shield, your AC equals 10 + your Dexterity modifier + your Wisdom modifier.</p>"
          },
          "enchant": {},
          "identifier": "unarmored-defense",
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "properties": [],
          "requirements": "Monk 1",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "class"
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          }
        },
        "type": "feat"
      },
      {
        "_id": "pchnXqd5C79fVlxy",
        "_stats": {
          "compendiumSource": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw",
          "duplicateSource": null,
          "exportSource": null,
          "coreVersion": "13.348",
          "systemId": "dnd5e",
          "systemVersion": "5.0.4",
          "lastModifiedBy": null
        },
        "effects": [],
        "flags": {
          "dnd5e": {
            "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
            "sourceId": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw"
          }
        },
        "folder": "TMmNG8ujFDBEWXRe",
        "img": "icons/skills/melee/unarmed-punch-fist.webp",
        "name": "Martial Arts",
        "ownership": {
          "default": 0
        },
        "sort": 0,
        "system": {
          "activities": {},
          "advancement": [],
          "crewed": false,
          "description": {
            "chat": "",
            "value": "<p>At 1st level, your practice of martial arts gives you mastery of combat styles that use and monk weapons, which are shortswords and any simple melee weapons that don't have the two-handed or heavy property. You gain the following benefits while you are unarmed or wielding only monk weapons and you aren't wearing armor or wielding a shield:</p>\n<ul>\n<li>\n<p>You can use Dexterity instead of Strength for the attack and damage rolls of your unarmed strikes and monk weapons.</p>\n</li>\n<li>\n<p>You can roll a d4 in place of the normal damage of your unarmed strike or monk weapon. This die changes as you gain monk levels, as shown in the Martial Arts column of the Monk table.</p>\n</li>\n<li>\n<p>When you use the Attack action with an unarmed strike or a monk weapon on your turn, you can make one unarmed strike as a bonus action. For example, if you take the Attack action and attack with a quarterstaff, you can also make an unarmed strike as a bonus action, assuming you haven't already taken a bonus action this turn.</p>\n</li>\n</ul>\n<p>Certain monasteries use specialized forms of the monk weapons. For example, you might use a club that is two lengths of wood connected by a short chain (called a nunchaku) or a sickle with a shorter, straighter blade (called a kama). Whatever name you use for a monk weapon, you can use the game statistics provided for the weapon.</p>"
          },
          "enchant": {},
          "identifier": "martial-arts",
          "prerequisites": {
            "level": null,
            "repeatable": false
          },
          "properties": [],
          "requirements": "Monk 1",
          "source": {
            "book": "SRD 5.1",
            "custom": "",
            "license": "CC-BY-4.0",
            "page": "",
            "revision": 1,
            "rules": "2014"
          },
          "type": {
            "subtype": "",
            "value": "class"
          },
          "uses": {
            "max": "",
            "recovery": [],
            "spent": 0
          }
        },
        "type": "feat"
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
| actorUuid | string | ✓ | body, query | UUID of the actor who owns the item |
| amount | number | ✓ | body, query | The amount to modify charges by (positive or negative) |
| clientId | string |  | body, query | Client ID for the Foundry world |
| itemUuid | string |  | body, query | The UUID of the specific item (optional if itemName provided) |
| itemName | string |  | body, query | The name of the item if UUID not provided (optional if itemUuid provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the charge modification operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/modify-item-charges"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"amount","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"itemUuid","type":"string","required":false,"source":"body"},{"name":"itemName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/modify-item-charges';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
curl -X POST 'http://localhost:3010/dnd5e/modify-item-charges?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","itemName":"Waterskin","amount":-1}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/modify-item-charges'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","itemName":"Waterskin","amount":-1}🔤 ➡️ body

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
  "type": "modify-item-charges-result",
  "requestId": "modify-item-charges_1774367605322",
  "data": {
    "itemUuid": "Actor.ioZexonJDGVuU8zl.Item.5skKSSB4ShHbKoc8",
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
| actorUuid | string | ✓ | body, query | UUID of the actor using the ability |
| clientId | string |  | body, query | Client ID for the Foundry world |
| abilityUuid | string |  | body, query | The UUID of the specific ability (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the ability if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the ability (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the ability use operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/use-ability"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/use-ability';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "abilityName": "Hammer"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/use-ability?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","abilityName":"Hammer"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/use-ability'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","abilityName":"Hammer"}🔤 ➡️ body

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
  "type": "use-ability-result",
  "requestId": "use-ability_1774367605665",
  "data": {
    "uuid": "Actor.ioZexonJDGVuU8zl",
    "ability": "Hammer",
    "result": "gj1tFIlRBZHn2lZE"
  }
}
```


---

## POST /dnd5e/use-feature

Use a class or racial feature for an actor. Activates class features (like Action Surge, Rage) or racial features (like Dragonborn Breath Weapon) for a character.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| actorUuid | string | ✓ | body, query | UUID of the actor using the feature |
| clientId | string |  | body, query | Client ID for the Foundry world |
| abilityUuid | string |  | body, query | The UUID of the specific feature (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the feature if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the feature (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the feature use operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/use-feature"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/use-feature';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "abilityName": "Priest"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/use-feature?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","abilityName":"Priest"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/use-feature'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","abilityName":"Priest"}🔤 ➡️ body

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
  "type": "use-feature-result",
  "requestId": "use-feature_1774367605619",
  "data": {
    "uuid": "Actor.ioZexonJDGVuU8zl",
    "ability": "Priest",
    "result": "oeo1w38mmP8wBl16"
  }
}
```


---

## POST /dnd5e/use-spell

Cast a spell for an actor. Casts a spell from the actor's spell list, consuming spell slots as appropriate. Handles cantrips, leveled spells, and spell-like abilities.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| actorUuid | string | ✓ | body, query | UUID of the actor casting the spell |
| clientId | string |  | body, query | Client ID for the Foundry world |
| abilityUuid | string |  | body, query | The UUID of the specific spell (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the spell if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the spell (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the spell casting operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/use-spell"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/use-spell';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "abilityName": "test-polymorph"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/use-spell?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","abilityName":"test-polymorph"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/use-spell'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","abilityName":"test-polymorph"}🔤 ➡️ body

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
  "type": "use-spell-result",
  "requestId": "use-spell_1774367605649",
  "data": {
    "uuid": "Actor.ioZexonJDGVuU8zl",
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
| actorUuid | string | ✓ | body, query | UUID of the actor using the item |
| clientId | string |  | body, query | Client ID for the Foundry world |
| abilityUuid | string |  | body, query | The UUID of the specific item (optional if abilityName provided) |
| abilityName | string |  | body, query | The name of the item if UUID not provided (optional if abilityUuid provided) |
| targetUuid | string |  | body, query | The UUID of the target for the item (optional) |
| targetName | string |  | body, query | The name of the target if UUID not provided (optional) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the item use operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/use-item"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"abilityUuid","type":"string","required":false,"source":"body"},{"name":"abilityName","type":"string","required":false,"source":"body"},{"name":"targetUuid","type":"string","required":false,"source":"body"},{"name":"targetName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/use-item';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "abilityName": "Hammer"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/use-item?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","abilityName":"Hammer"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/use-item'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","abilityName":"Hammer"}🔤 ➡️ body

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
  "type": "use-item-result",
  "requestId": "use-item_1774367605569",
  "data": {
    "uuid": "Actor.ioZexonJDGVuU8zl",
    "ability": "Hammer",
    "result": "MAccCmSp62Gr90hH"
  }
}
```


---

## POST /dnd5e/short-rest

Perform a short rest for an actor. Triggers the D&D 5e short rest workflow including hit dice recovery, class feature resets, and HP recovery.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string |  | body, query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor to rest (optional if selected is true) |
| selected | boolean |  | body, query | Use the currently selected token's actor |
| autoHD | boolean |  | body, query | Automatically spend hit dice during short rest |
| autoHDThreshold | number |  | body, query | HP threshold below which to auto-spend hit dice (0-1 as fraction of max HP) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the short rest operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/short-rest"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"body"},{"name":"autoHD","type":"boolean","required":false,"source":"body"},{"name":"autoHDThreshold","type":"number","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/short-rest';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/short-rest?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/short-rest'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl"
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
  const path = '/dnd5e/short-rest';
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl"
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
  🔤/dnd5e/short-rest🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/short-rest🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 38❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "short-rest-result",
  "requestId": "short-rest_1774367606152",
  "data": {
    "actorUuid": "Actor.ioZexonJDGVuU8zl",
    "result": {
      "type": "short",
      "deltas": {
        "hitPoints": 0,
        "hitDice": 0
      },
      "updateData": {
        "system": {
          "spells": {
            "pact": {
              "value": 0
            }
          }
        },
        "_id": "ioZexonJDGVuU8zl",
        "type": "character"
      },
      "updateItems": [],
      "newDay": false,
      "rolls": [],
      "clone": {
        "effects": [
          {
            "_id": "dnd5ebloodied000",
            "name": "Bloodied",
            "img": "systems/dnd5e/icons/svg/statuses/bloodied.svg",
            "statuses": [
              "bloodied"
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
              "createdTime": 1774367605935,
              "modifiedTime": 1774367605935,
              "lastModifiedBy": "r6bXhB7k9cXa3cif"
            }
          }
        ],
        "img": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
        "items": [
          {
            "_id": "q4tr1vTU8RxtU1UZ",
            "_stats": {
              "compendiumSource": null,
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": null,
            "img": "icons/sundries/documents/document-torn-diagram-tan.webp",
            "name": "Priest",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "advancement": [],
              "description": {
                "chat": "",
                "value": "<ul><li><strong>Skill Proficiencies:</strong> Insight, Religion</li><li><strong>Languages:</strong> Two of your choice</li><li><strong>Equipment:</strong> A holy symbol, 5 sticks of incense, prayer book, vestments, a set of common clothes, and a pouch containing 15 gp.</li></ul>"
              },
              "identifier": "priest",
              "source": {
                "book": "",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "startingEquipment": []
            },
            "type": "background"
          },
          {
            "_id": "O3ZjSw0GtAOPScHI",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.ZgYBjYYfiUstQD6f",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "uQNrdSpcXeoB5EnX",
            "img": "icons/equipment/feet/shoes-leather-simple-brown.webp",
            "name": "Lightfoot Halfling",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
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
              "description": {
                "chat": "",
                "value": "<p>Your halfling character has a number of traits in common with all other halflings.</p><p><em><strong>Ability Score Increase.</strong></em> Your Dexterity score increases by 2.</p><p><em><strong>Age.</strong></em> A halfling reaches adulthood at the age of 20 and generally lives into the middle of his or her second century.</p><p><em><strong>Alignment.</strong></em> Most halflings are lawful good. As a rule, they are good-hearted and kind, hate to see others in pain, and have no tolerance for oppression. They are also very orderly and traditional, leaning heavily on the support of their community and the comfort of their old ways.</p><p><em><strong>Size.</strong></em> Halflings average about 3 feet tall and weigh about 40 pounds. Your size is Small.</p><p><em><strong>Speed.</strong></em> Your base walking speed is 25 feet.</p><p><em><strong>Lucky.</strong></em> When you roll a 1 on the d20 for an attack roll, ability check, or saving throw, you can reroll the die and must use the new roll.</p><p><em><strong>Brave.</strong></em> You have advantage on saving throws against being frightened.</p><p><em><strong>Halfling Nimbleness.</strong></em> You can move through the space of any creature that is of a size larger than yours.</p><p><em><strong>Languages.</strong></em> You can speak, read, and write Common and Halfling. The Halfling language isn't secret, but halflings are loath to share it with others. They write very little, so they don't have a rich body of literature. Their oral tradition, however, is very strong. Almost all halflings speak Common to converse with the people in whose lands they dwell or through which they are traveling.</p><h5>Lightfoot</h5><p>As a lightfoot halfling, you can easily hide from notice, even using other people as cover. You're inclined to be affable and get along well with others.</p><p>Lightfoots are more prone to wanderlust than other halflings, and often dwell alongside other races or take up a nomadic life.</p><p><em><strong>Ability Score Increase.</strong></em> Your Charisma score increases by 1.</p><p><em><strong>Naturally Stealthy.</strong></em> You can attempt to hide even when you are obscured only by a creature that is at least one size larger than you.</p>"
              },
              "identifier": "lightfoot-halfling",
              "movement": {
                "burrow": null,
                "climb": null,
                "fly": null,
                "hover": false,
                "swim": null,
                "units": "ft",
                "walk": 25
              },
              "senses": {
                "blindsight": null,
                "darkvision": null,
                "special": "",
                "tremorsense": null,
                "truesight": null,
                "units": "ft"
              },
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "custom": "",
                "subtype": "halfling",
                "value": "humanoid"
              }
            },
            "type": "race"
          },
          {
            "_id": "FtOM4QiOW5MwgcS3",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "kbtbKofcv13crhke",
            "img": "icons/sundries/gaming/dice-runed-brown.webp",
            "name": "Lucky",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>When you roll a 1 on the d20 for an attack roll, ability check, or saving throw, you can reroll the die and must use the new roll.</p><section class=\"secret foundry-note\" id=\"secret-S04TPyvUh05Dz0Ng\"><p><strong>Foundry Note</strong></p><p>This property can be enabled on your character sheet in the Special Traits configuration on the Attributes tab.</p></section>"
              },
              "enchant": {},
              "identifier": "lucky",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Halfling",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "race"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "nmmihiqphHjoE8dl",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "kbtbKofcv13crhke",
            "img": "icons/skills/melee/unarmed-punch-fist.webp",
            "name": "Brave",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>You have advantage on saving throws against being frightened.</p>"
              },
              "enchant": {},
              "identifier": "brave",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Halfling",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "race"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "cWrETHzCRs1Ueqd3",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "kbtbKofcv13crhke",
            "img": "icons/skills/movement/feet-winged-boots-brown.webp",
            "name": "Halfling Nimbleness",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>You can move through the space of any creature that is of a size larger than yours.</p>"
              },
              "enchant": {},
              "identifier": "halfling-nimbleness",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Halfling",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "race"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "AArhiOrSkaQUnCZS",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "kbtbKofcv13crhke",
            "img": "icons/magic/perception/silhouette-stealth-shadow.webp",
            "name": "Naturally Stealthy",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>You can attempt to hide even when you are obscured only by a creature that is at least one size larger than you.</p>"
              },
              "enchant": {},
              "identifier": "naturally-stealthy",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Lightfoot Halfling",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "race"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "8Grf7ga6JcZF0X6x",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.classes.Item.6VoZrWxhOEKGYhnq",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "HQ1Oy7HkbnxnE63o",
            "img": "icons/skills/melee/hand-grip-staff-blue.webp",
            "name": "Monk",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
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
              "description": {
                "chat": "",
                "value": "<p>As a monk, you gain the following class features.</p><h3>Hit Points</h3><p><strong>Hit Dice:</strong> 1d8 per monk level<br /><strong>Hit Points at 1st Level:</strong> 8 + your Constitution modifier<br /><strong>Hit Points at Higher Levels:</strong> 1d8 (or 5) + your Constitution modifier per monk level after 1st</p><h3>Proficiencies</h3><p><strong>Armor:</strong> None<br /><strong>Weapons:</strong> Simple weapons, shortswords<br /><strong>Tools:</strong> Choose one type of artisan's tools or one musical instrument<br /><strong>Saving Throws:</strong> Strength, Dexterity<br /><strong>Skills:</strong> Choose two from Acrobatics, Athletics, History, Insight, Religion, and Stealth</p><h1>Monk Advancement</h1><table><thead><tr><td>Level</td><td>Proficiency Bonus</td><td>Martial Arts</td><td>Ki Points</td><td>Unarmored Movement</td><td>Features</td></tr></thead><tbody><tr><td>1st</td><td>+2</td><td>1d4</td><td>—</td><td>—</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI]{Unarmored Defense}, @UUID[Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw]{Martial Arts}</td></tr><tr><td>2nd</td><td>+2</td><td>1d4</td><td>2</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.10b6z2W1txNkrGP7]{Ki}, @UUID[Compendium.dnd5e.classfeatures.Item.zCeqyQ8uIPNdYJSW]{Unarmored Movement}</td></tr><tr><td>3rd</td><td>+2</td><td>1d4</td><td>3</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.rtpQdX77dYWbDIOH]{Monastic Tradition}, @UUID[Compendium.dnd5e.classfeatures.Item.mzweVbnsJPQiVkAe]{Deflect Missiles}</td></tr><tr><td>4th</td><td>+2</td><td>1d4</td><td>4</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}, @UUID[Compendium.dnd5e.classfeatures.Item.KQz9bqxVkXjDl8gK]{Slow Fall}</td></tr><tr><td>5th</td><td>+3</td><td>1d6</td><td>5</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.XogoBnFWmCAHXppo]{Extra Attack}, @UUID[Compendium.dnd5e.classfeatures.Item.pvRc6GAu1ok6zihC]{Stunning Strike}</td></tr><tr><td>6th</td><td>+3</td><td>1d6</td><td>6</td><td>+15 ft.</td><td><p>@UUID[Compendium.dnd5e.classfeatures.Item.7flZKruSSu6dHg6D]{Ki-Empowered Strikes},</p><p>Monastic Tradition feature</p></td></tr><tr><td>7th</td><td>+3</td><td>1d6</td><td>7</td><td>+15 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.a4P4DNMmH8CqSNkC]{Evasion}, @UUID[Compendium.dnd5e.classfeatures.Item.ZmC31XKS4YNENnoc]{Stillness of Mind}</td></tr><tr><td>8th</td><td>+3</td><td>1d6</td><td>8</td><td>+15 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>9th</td><td>+4</td><td>1d6</td><td>9</td><td>+15 ft.</td><td>Unarmored Movement improvement</td></tr><tr><td>10th</td><td>+4</td><td>1d6</td><td>10</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.bqWA7t9pDELbNRkp]{Purity of Body}</td></tr><tr><td>11th</td><td>+4</td><td>1d8</td><td>11</td><td>+20 ft.</td><td>Monastic Tradition feature</td></tr><tr><td>12th</td><td>+4</td><td>1d8</td><td>12</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>13th</td><td>+5</td><td>1d8</td><td>13</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.XjuGBeB8Y0C3A5D4]{Tongue of the Sun and Moon}</td></tr><tr><td>14th</td><td>+5</td><td>1d8</td><td>14</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.7D2EkLdISwShEDlN]{Diamond Soul}</td></tr><tr><td>15th</td><td>+5</td><td>1d8</td><td>15</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.gDH8PMrKvLHaNmEI]{Timeless Body}</td></tr><tr><td>16th</td><td>+5</td><td>1d8</td><td>16</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>17th</td><td>+6</td><td>1d10</td><td>17</td><td>+25 ft.</td><td>Monastic Tradition feature</td></tr><tr><td>18th</td><td>+6</td><td>1d10</td><td>18</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.3jwFt3hSqDswBlOH]{Empty Body}</td></tr><tr><td>19th</td><td>+6</td><td>1d10</td><td>19</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>20th</td><td>+6</td><td>1d10</td><td>20</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.mQNPg89YIs7g5tG4]{Perfect Self}</td></tr></tbody></table><h1>Monastic Traditions</h1><p>Three traditions of monastic pursuit are common in the monasteries scattered across the multiverse. Most monasteries practice one tradition exclusively, but a few honor the three traditions and instruct each monk according to his or her aptitude and interest. All three traditions rely on the same basic techniques, diverging as the student grows more adept. Thus, a monk need choose a tradition only upon reaching 3rd level.</p><p>@UUID[Compendium.dnd5e.subclasses.Item.IvlpKMXX3PmW1NY2]{Way of the Open Hand}</p>"
              },
              "hd": {
                "additional": "",
                "denomination": "d8",
                "spent": 0
              },
              "identifier": "monk",
              "levels": 1,
              "primaryAbility": {
                "all": true,
                "value": []
              },
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "spellcasting": {
                "ability": "",
                "preparation": {
                  "formula": ""
                },
                "progression": "none"
              },
              "startingEquipment": [
                {
                  "_id": "5yj0P4r9teJDnDtd",
                  "group": "",
                  "requiresProficiency": false,
                  "sort": 100000,
                  "type": "OR"
                },
                {
                  "_id": "R5tuRtaPonfjQCVU",
                  "count": null,
                  "group": "5yj0P4r9teJDnDtd",
                  "key": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
                  "requiresProficiency": false,
                  "sort": 200000,
                  "type": "linked"
                },
                {
                  "_id": "Mlf6kel8ws6xgDER",
                  "count": null,
                  "group": "5yj0P4r9teJDnDtd",
                  "key": "simpleM",
                  "requiresProficiency": false,
                  "sort": 300000,
                  "type": "weapon"
                },
                {
                  "_id": "3TbVLmLPtjVaSh5O",
                  "group": "",
                  "requiresProficiency": false,
                  "sort": 400000,
                  "type": "OR"
                },
                {
                  "_id": "AvDYtl0uvQsDuhnb",
                  "count": null,
                  "group": "3TbVLmLPtjVaSh5O",
                  "key": "Compendium.dnd5e.items.Item.XY8b594Dn7plACLL",
                  "requiresProficiency": false,
                  "sort": 500000,
                  "type": "linked"
                },
                {
                  "_id": "4QKQURCmIurbTAzp",
                  "count": null,
                  "group": "3TbVLmLPtjVaSh5O",
                  "key": "Compendium.dnd5e.items.Item.8KWz5DJbWUpNWniP",
                  "requiresProficiency": false,
                  "sort": 600000,
                  "type": "linked"
                },
                {
                  "_id": "AOYuulsULvsHbSLO",
                  "count": 10,
                  "group": "",
                  "key": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
                  "requiresProficiency": false,
                  "sort": 700000,
                  "type": "linked"
                }
              ],
              "wealth": "5d4"
            },
            "type": "class"
          },
          {
            "_id": "CwgoTDXWCD7PknIN",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [
              {
                "_id": "R5ro4AuNjcdWD56O",
                "_stats": {
                  "compendiumSource": null,
                  "duplicateSource": null,
                  "exportSource": null,
                  "coreVersion": "13.348",
                  "systemId": "dnd5e",
                  "systemVersion": "5.0.4",
                  "lastModifiedBy": null
                },
                "changes": [
                  {
                    "key": "system.attributes.ac.calc",
                    "mode": 5,
                    "priority": null,
                    "value": "unarmoredMonk"
                  }
                ],
                "description": "",
                "disabled": false,
                "duration": {
                  "combat": null,
                  "rounds": null,
                  "seconds": null,
                  "startRound": null,
                  "startTime": 0,
                  "startTurn": null,
                  "turns": null
                },
                "flags": {},
                "img": "icons/magic/control/silhouette-hold-change-blue.webp",
                "name": "Unarmored Defense",
                "origin": "Item.cOdcNWy4hII029DT",
                "sort": 0,
                "statuses": [],
                "system": {},
                "tint": "#ffffff",
                "transfer": true,
                "type": "base"
              }
            ],
            "flags": {
              "dnd5e": {
                "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
                "sourceId": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI"
              }
            },
            "folder": "TMmNG8ujFDBEWXRe",
            "img": "icons/magic/control/silhouette-hold-change-blue.webp",
            "name": "Unarmored Defense",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>Beginning at 1st Level, while you are wearing no armor and not wielding a Shield, your AC equals 10 + your Dexterity modifier + your Wisdom modifier.</p>"
              },
              "enchant": {},
              "identifier": "unarmored-defense",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Monk 1",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "class"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "pchnXqd5C79fVlxy",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {
              "dnd5e": {
                "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
                "sourceId": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw"
              }
            },
            "folder": "TMmNG8ujFDBEWXRe",
            "img": "icons/skills/melee/unarmed-punch-fist.webp",
            "name": "Martial Arts",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>At 1st level, your practice of martial arts gives you mastery of combat styles that use and monk weapons, which are shortswords and any simple melee weapons that don't have the two-handed or heavy property. You gain the following benefits while you are unarmed or wielding only monk weapons and you aren't wearing armor or wielding a shield:</p>\n<ul>\n<li>\n<p>You can use Dexterity instead of Strength for the attack and damage rolls of your unarmed strikes and monk weapons.</p>\n</li>\n<li>\n<p>You can roll a d4 in place of the normal damage of your unarmed strike or monk weapon. This die changes as you gain monk levels, as shown in the Martial Arts column of the Monk table.</p>\n</li>\n<li>\n<p>When you use the Attack action with an unarmed strike or a monk weapon on your turn, you can make one unarmed strike as a bonus action. For example, if you take the Attack action and attack with a quarterstaff, you can also make an unarmed strike as a bonus action, assuming you haven't already taken a bonus action this turn.</p>\n</li>\n</ul>\n<p>Certain monasteries use specialized forms of the monk weapons. For example, you might use a club that is two lengths of wood connected by a short chain (called a nunchaku) or a sickle with a shorter, straighter blade (called a kama). Whatever name you use for a monk weapon, you can use the game statistics provided for the weapon.</p>"
              },
              "enchant": {},
              "identifier": "martial-arts",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Monk 1",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "class"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "YJ1P3PnFKHOdQpaP",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.14pNRT4sZy9rgvhb",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "dlru9Hy74nSMv6fr",
            "img": "icons/tools/hand/hammer-cobbler-steel.webp",
            "name": "Hammer",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": null,
              "description": {
                "chat": "",
                "value": "<p>A tool with a heavy metal head mounted at the end of its handle, used for jobs such as breaking things and driving in nails. </p>"
              },
              "identified": true,
              "identifier": "hammer",
              "price": {
                "denomination": "gp",
                "value": 1
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 3
              }
            },
            "type": "loot"
          },
          {
            "_id": "DDnxRCeYUhXstWU8",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.1FSubnBpSTDmVaYV",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/sundries/lights/torch-black.webp",
            "name": "Tinderbox",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": "8KWz5DJbWUpNWniP",
              "description": {
                "chat": "",
                "value": "<p>This small container holds flint, fire steel, and tinder (usually dry cloth soaked in light oil) used to kindle a fire. Using it to light a torch - or anything else with abundant, exposed fuel - takes an action. Lighting any other fire takes 1 minute.</p>\n<p> </p>"
              },
              "identified": true,
              "identifier": "tinderbox",
              "price": {
                "denomination": "sp",
                "value": 5
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 1
              }
            },
            "type": "loot"
          },
          {
            "_id": "5skKSSB4ShHbKoc8",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.1L5wkmbw0fmNAr38",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": "r6bXhB7k9cXa3cif",
              "modifiedTime": 1774367605451
            },
            "effects": [],
            "flags": {
              "dnd5e": {
                "migratedUses": 4,
                "riders": {
                  "activity": [],
                  "effect": []
                }
              }
            },
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/sundries/survival/wetskin-leather-purple.webp",
            "name": "Waterskin",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": "6OYR11aJX2dEVtOj",
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>A leather hide sewn into an enclosed skin which can contain up to 4 pints of liquid. It weighs 5 pounds when full; a pint of water is approximately 1 pound.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "waterskin",
              "magicalBonus": null,
              "price": {
                "denomination": "sp",
                "value": 2
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "food"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": false,
                "max": "4",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 5
              }
            },
            "type": "consumable"
          },
          {
            "_id": "dDuMscUuMI2bTdkj",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.29ZLE8PERtFVD3QU",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/sundries/lights/torch-black.webp",
            "name": "Torch",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": "8KWz5DJbWUpNWniP",
              "damage": {
                "base": {
                  "bonus": "",
                  "custom": {
                    "enabled": true,
                    "formula": "1"
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": [
                    "fire"
                  ]
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>A torch burns for 1 hour, providing bright light in a 20-foot radius and dim light for an additional 20 feet. If you make a melee attack with a burning torch and hit, it deals 1 fire damage.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "torch",
              "magicalBonus": null,
              "price": {
                "denomination": "cp",
                "value": 1
              },
              "properties": [],
              "quantity": 10,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "trinket"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": false,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 1
              }
            },
            "type": "consumable"
          },
          {
            "_id": "nC6VcR5JAIbR4err",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.3b0RvGi0TnTYpIxn",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "dlru9Hy74nSMv6fr",
            "img": "icons/consumables/grains/breadsticks-crackers-wrapped-ration-brown.webp",
            "name": "Stick of Incense",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": null,
              "description": {
                "chat": "",
                "value": "<p>When blocks of incense cannot be used or a cheaper alternative is required, people often use these to perfume the air, whether for pleasurable or religious purposes.</p>"
              },
              "identified": true,
              "identifier": "stick-of-incense",
              "price": {
                "denomination": "sp",
                "value": 2
              },
              "properties": [],
              "quantity": 5,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 0
              }
            },
            "type": "loot"
          },
          {
            "_id": "WeKJI3gPUAU52WAX",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "MLMTCAvKsuFE3vYA",
            "img": "icons/weapons/thrown/dart-feathered.webp",
            "name": "Dart",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "ammunition": {},
              "armor": {
                "value": 10
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "damage": {
                "base": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": 4,
                  "number": 1,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": [
                    "piercing"
                  ]
                },
                "versatile": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": []
                }
              },
              "description": {
                "chat": "",
                "value": "<p>A small thrown implement crafted with a short wooden shaft and crossed feathres with a sharp wooden or metal tip. Darts can be thrown with sufficient force to puncture the skin.</p>"
              },
              "equipped": true,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "dart",
              "magicalBonus": null,
              "mastery": "",
              "price": {
                "denomination": "cp",
                "value": 5
              },
              "proficient": null,
              "properties": [
                "fin",
                "thr"
              ],
              "quantity": 10,
              "range": {
                "long": 60,
                "reach": null,
                "units": "ft",
                "value": 20
              },
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "baseItem": "dart",
                "value": "simpleR"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 0.25
              }
            },
            "type": "weapon"
          },
          {
            "_id": "1F73YcUHbZMgePAD",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.8RXjiddJ6VGyE7vB",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "aJgMxnZED9XdoN2W",
            "img": "icons/equipment/chest/shirt-collared-brown.webp",
            "name": "Common Clothes",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "armor": {
                "dex": null,
                "magicalBonus": null,
                "value": null
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>Clothes worn by most commoners.</p>"
              },
              "equipped": false,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "common-clothes",
              "price": {
                "denomination": "sp",
                "value": 5
              },
              "proficient": null,
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "speed": {
                "conditions": "",
                "value": null
              },
              "strength": null,
              "type": {
                "baseItem": "",
                "value": "clothing"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 3
              }
            },
            "type": "equipment"
          },
          {
            "_id": "LQhGSEatJ3VK7oqW",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.8d95YV1jHcxPygJ9",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/consumables/grains/bread-loaf-boule-rustic-brown.webp",
            "name": "Rations",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": "XY8b594Dn7plACLL",
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>Rations consist of dry foods suitable for extended travel, including jerky, dried fruit, hardtack, and nuts.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "rations",
              "magicalBonus": null,
              "price": {
                "denomination": "sp",
                "value": 5
              },
              "properties": [],
              "quantity": 10,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "food"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": true,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 2
              }
            },
            "type": "consumable"
          },
          {
            "_id": "1kqMwSwjfErSFNvl",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.9bWTRRDym06PzSAf",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "PgtyWdfkDTDbc0ov",
            "img": "icons/containers/bags/pouch-rounded-leather-gold-tan.webp",
            "name": "Pouch",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "attuned": false,
              "attunement": "",
              "capacity": {
                "volume": {
                  "units": "cubicFoot"
                },
                "weight": {
                  "units": "lb",
                  "value": 6
                }
              },
              "container": null,
              "currency": {
                "cp": 0,
                "ep": 0,
                "gp": 0,
                "pp": 0,
                "sp": 0
              },
              "description": {
                "chat": "",
                "value": "<p>A cloth or leather pouch can hold up to 20 sling bullets or 50 blowgun needles, among other things. A compartmentalized pouch for holding spell components is called a component pouch. A pouch can hold up to ⅕ cubic foot or 6 pounds of gear.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "pouch",
              "price": {
                "denomination": "sp",
                "value": 5
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 1
              }
            },
            "type": "container"
          },
          {
            "_id": "zF5L4xHnJBC7b2iM",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.AkyQyonZMVcvOrXU",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/tools/hand/pickaxe-steel-white.webp",
            "name": "Crowbar",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": "XY8b594Dn7plACLL",
              "description": {
                "chat": "",
                "value": "<p>Using a crowbar grants advantage to Strength checks where the crowbar's leverage can be applied.</p>"
              },
              "identified": true,
              "identifier": "crowbar",
              "price": {
                "denomination": "gp",
                "value": 2
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 5
              }
            },
            "type": "loot"
          },
          {
            "_id": "g7U3OAXVcoI4lwzf",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.GsuvwoekKZatfKwF",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "MLMTCAvKsuFE3vYA",
            "img": "icons/skills/melee/unarmed-punch-fist.webp",
            "name": "Unarmed Strike",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "ammunition": {},
              "armor": {
                "value": 10
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "damage": {
                "base": {
                  "bonus": "",
                  "custom": {
                    "enabled": true,
                    "formula": "@scale.monk.die"
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": [
                    "bludgeoning"
                  ]
                },
                "versatile": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": []
                }
              },
              "description": {
                "chat": "",
                "value": "<p>A  punch, kick, head-butt, or similar forceful blow (none of which count as weapons). On a hit, an unarmed strike deals bludgeoning damage equal to 1 + your Strength modifier. You are proficient with your unarmed strikes.</p>"
              },
              "equipped": true,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "unarmed-strike",
              "magicalBonus": null,
              "mastery": "",
              "price": {
                "denomination": "gp",
                "value": 0
              },
              "proficient": null,
              "properties": [
                "fin"
              ],
              "quantity": 1,
              "range": {
                "long": null,
                "reach": null,
                "units": "ft",
                "value": null
              },
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "baseItem": "",
                "value": "simpleM"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 0
              }
            },
            "type": "weapon"
          },
          {
            "_id": "WcROGD590imVj9qp",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.H8YCd689ezlD26aT",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "PgtyWdfkDTDbc0ov",
            "img": "icons/containers/bags/pack-leather-white-tan.webp",
            "name": "Backpack",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "attuned": false,
              "attunement": "",
              "capacity": {
                "volume": {
                  "units": "cubicFoot"
                },
                "weight": {
                  "units": "lb",
                  "value": 30
                }
              },
              "container": null,
              "currency": {
                "cp": 0,
                "ep": 0,
                "gp": 0,
                "pp": 0,
                "sp": 0
              },
              "description": {
                "chat": "",
                "value": "<p>A backpack can hold one cubic foot or 30 pounds of gear. You can also strap items, such as a bedroll or a coil of rope, to the outside of a backpack.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "backpack",
              "price": {
                "denomination": "gp",
                "value": 2
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 5
              }
            },
            "type": "container"
          },
          {
            "_id": "8UQanAvHUIHZXp0O",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.P31t6tGgt9aLAdYt",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/tools/fasteners/nail-steel.webp",
            "name": "Piton",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": "xsB7Y2WI476kvOt4",
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>A metal spike that is drive into a seam in a climbing surface with a climbing hammer. It can also be used like iron spikes to block doors in a pinch.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "piton",
              "magicalBonus": null,
              "price": {
                "denomination": "cp",
                "value": 5
              },
              "properties": [],
              "quantity": 10,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "trinket"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": false,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 0.25
              }
            },
            "type": "consumable"
          },
          {
            "_id": "irCoxLHje0eA6Qmu",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.QXmaarJ4X8P0C1HV",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {
              "dnd5e": {
                "riders": {
                  "activity": [],
                  "effect": []
                }
              }
            },
            "folder": "UnUwTG4YIgd0kaUJ",
            "img": "icons/sundries/survival/rope-wrapped-brown.webp",
            "name": "Hempen Rope (50 ft.)",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>Rope, whether made of hemp or silk, has 2 hit points and can be burst with a DC 17 Strength check.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "hempen-rope-50-ft",
              "magicalBonus": null,
              "price": {
                "denomination": "gp",
                "value": 1
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "trinket"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": false,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 10
              }
            },
            "type": "consumable"
          },
          {
            "_id": "GfKuJYLFfw00oW9R",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.gP1URGq3kVIIFHJ7",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {
              "dnd5e": {
                "riders": {
                  "activity": [],
                  "effect": []
                }
              }
            },
            "folder": "xedn1r43VWuEBcli",
            "img": "icons/containers/chest/chest-reinforced-steel-red.webp",
            "name": "Reliquary",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "armor": {
                "dex": null,
                "magicalBonus": null,
                "value": null
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>A tiny box or other container holding a fragment of a precious relic, saint, or other historical figure that dedicated their life to walk the path of a true believer. A deity imbues the bearer of this artifact with the ability to call forth their power and in doing so spread the faith once more.</p>\n<p><strong>Spellcasting Focus</strong>. A cleric or paladin can use a holy symbol as a spellcasting focus. To use the symbol in this way, the caster must hold it in hand, wear it visibly, or bear it on a shield.</p>"
              },
              "equipped": false,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "reliquary",
              "price": {
                "denomination": "gp",
                "value": 5
              },
              "proficient": null,
              "properties": [
                "foc"
              ],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "speed": {
                "conditions": "",
                "value": null
              },
              "strength": null,
              "type": {
                "baseItem": "",
                "value": "trinket"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 2
              }
            },
            "type": "equipment"
          },
          {
            "_id": "NsNcJBUG5ajbk7sG",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.irtqrzaUCeshmTZp",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "aJgMxnZED9XdoN2W",
            "img": "icons/equipment/back/mantle-collared-black.webp",
            "name": "Vestments",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "armor": {
                "dex": null,
                "magicalBonus": null,
                "value": null
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>Simple or ostentacious wear, often used by priests and other religious figures for use in rituals and ceremonies.</p>"
              },
              "equipped": true,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "vestments",
              "price": {
                "denomination": "gp",
                "value": 1
              },
              "proficient": null,
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "speed": {
                "conditions": "",
                "value": null
              },
              "strength": null,
              "type": {
                "baseItem": "",
                "value": "clothing"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 4
              }
            },
            "type": "equipment"
          },
          {
            "_id": "XbF0CTipyqbMKHsB",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.odV5cq2HSLSCH69k",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "dlru9Hy74nSMv6fr",
            "img": "icons/sundries/books/book-purple-cross.webp",
            "name": "Prayer Book",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": null,
              "description": {
                "chat": "",
                "value": "<p>A book containing prayers and incantations dedicated to a specific power for the faithful to follow.</p>"
              },
              "identified": true,
              "identifier": "prayer-book",
              "price": {
                "denomination": "gp",
                "value": 25
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 5
              }
            },
            "type": "loot"
          },
          {
            "_id": "GYyovoNdU66nxOlX",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "MLMTCAvKsuFE3vYA",
            "img": "icons/weapons/swords/sword-guard-worn-purple.webp",
            "name": "Shortsword",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "ammunition": {},
              "armor": {
                "value": 10
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "damage": {
                "base": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": 6,
                  "number": 1,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": [
                    "piercing"
                  ]
                },
                "versatile": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": []
                }
              },
              "description": {
                "chat": "",
                "value": "<p>A medium sized blade with a firm crossguard and a leather-wrapped handle. A versatile weapon which makes up in versatility what it lacks in reach.</p>"
              },
              "equipped": true,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "shortsword",
              "magicalBonus": null,
              "price": {
                "denomination": "gp",
                "value": 10
              },
              "proficient": null,
              "properties": [
                "fin",
                "lgt"
              ],
              "quantity": 1,
              "range": {
                "long": null,
                "reach": null,
                "units": "ft",
                "value": null
              },
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "baseItem": "shortsword",
                "value": "martialM"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 2
              }
            },
            "type": "weapon"
          },
          {
            "effects": [
              {
                "_id": "tM8Yj0dw52paxVc4",
                "_stats": {
                  "compendiumSource": null,
                  "duplicateSource": null,
                  "exportSource": null,
                  "coreVersion": "13.348",
                  "systemId": "dnd5e",
                  "systemVersion": "5.0.4",
                  "lastModifiedBy": null
                },
                "changes": [],
                "description": "",
                "disabled": true,
                "duration": {
                  "combat": null,
                  "rounds": null,
                  "seconds": 3600,
                  "startRound": null,
                  "startTime": null,
                  "startTurn": null,
                  "turns": null
                },
                "flags": {},
                "img": "icons/magic/control/energy-stream-link-large-teal.webp",
                "name": "Polymorph",
                "origin": "Compendium.dnd5e.spells.Item.04nMsTWkIFvkbXlY",
                "sort": 0,
                "statuses": [
                  "transformed"
                ],
                "system": {},
                "tint": "#ffffff",
                "transfer": false,
                "type": "base"
              }
            ],
            "img": "icons/magic/control/energy-stream-link-large-teal.webp",
            "name": "test-polymorph",
            "system": {
              "activation": {
                "condition": "",
                "type": "action",
                "value": 1
              },
              "activities": {
                "dnd5eactivity000": {
                  "_id": "dnd5eactivity000",
                  "activation": {
                    "override": false,
                    "type": "action",
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
                  "damage": {
                    "onSave": "half",
                    "parts": []
                  },
                  "description": {
                    "chatFlavor": ""
                  },
                  "duration": {
                    "concentration": false,
                    "override": false,
                    "units": "inst"
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
                  "save": {
                    "ability": [
                      "wis"
                    ],
                    "dc": {
                      "calculation": "spellcasting",
                      "formula": ""
                    }
                  },
                  "sort": 0,
                  "target": {
                    "affects": {
                      "choice": false
                    },
                    "override": false,
                    "prompt": true,
                    "template": {
                      "contiguous": false,
                      "units": "ft"
                    }
                  },
                  "type": "save",
                  "uses": {
                    "recovery": [],
                    "spent": 0
                  }
                }
              },
              "description": {
                "chat": "",
                "value": "<p>This spell transforms a creature that you can see within range into a new form. An unwilling creature must make a Wisdom saving throw to avoid the effect. The spell has no effect on a Shapechanger or a creature with 0 Hit Points.</p><p>The transformation lasts for the Duration, or until the target drops to 0 Hit Points or dies. The new form can be any beast whose Challenge rating is equal to or less than the target's (or the target's level, if it doesn't have a challenge rating). The target's game statistics, including mental Ability Scores, are replaced by the Statistics of the chosen beast. It retains its Alignment and personality.</p><p>The target assumes the Hit Points of its new form. When it reverts to its normal form, the creature returns to the number of hit points it had before it transformed. If it reverts as a result of dropping to 0 hit points, any excess damage carries over to its normal form. As long as the excess damage doesn't reduce the creature's normal form to 0 hit points, it isn't knocked Unconscious.</p><p>The creature is limited in the Actions it can perform by the Nature of its new form, and it can't speak, cast Spells, or take any other action that requires hands or speech.</p><p>The target's gear melds into the new form. The creature can't activate, use, wield, or otherwise benefit from any of its Equipment.</p>"
              },
              "duration": {
                "units": "hour",
                "value": "1"
              },
              "identifier": "polymorph",
              "level": 4,
              "materials": {
                "consumed": false,
                "cost": 0,
                "supply": 0,
                "value": "A caterpillar cocoon"
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
              "range": {
                "units": "ft",
                "value": "60"
              },
              "school": "trs",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "target": {
                "affects": {
                  "choice": false,
                  "count": "1",
                  "type": "creature"
                },
                "template": {
                  "contiguous": false,
                  "units": ""
                }
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "spell",
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
              "createdTime": 1774367604775,
              "modifiedTime": 1774367604775,
              "lastModifiedBy": "r6bXhB7k9cXa3cif"
            },
            "_id": "cuN0zv7OrhLW1LKO"
          },
          {
            "effects": [],
            "img": "icons/consumables/potions/bottle-bulb-corked-labeled-blue.webp",
            "name": "test-potion of storm giant strength",
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
                    "units": "hour",
                    "value": "1"
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
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p><em>This potion's transparent liquid has floating in it a sliver of fingernail from a giant of the appropriate type.</em></p>\n<p>When you drink this potion, your Strength score changes to 29 for 1 hour.  The potion has no effect on you if your Strength is equal to or greater than that score.</p>\n<p> </p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "potion-of-storm-giant-strength",
              "magicalBonus": null,
              "price": {
                "denomination": "gp",
                "value": 2000
              },
              "properties": [
                "mgc"
              ],
              "quantity": 1,
              "rarity": "legendary",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "potion"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": true,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 0.1
              }
            },
            "type": "consumable",
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
              "createdTime": 1774367604969,
              "modifiedTime": 1774367604969,
              "lastModifiedBy": "r6bXhB7k9cXa3cif"
            },
            "_id": "xzMsr0JNmiadK4Cr"
          }
        ],
        "name": "Updated Test Actor",
        "prototypeToken": {
          "actorLink": true,
          "alpha": 1,
          "appendNumber": false,
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
          "prependAdjective": false,
          "randomImg": false,
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
          "width": 1
        },
        "system": {
          "abilities": {
            "cha": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 0,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 11
            },
            "con": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 0,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 13
            },
            "dex": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 1,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 16
            },
            "int": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 0,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 12
            },
            "str": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 1,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 8
            },
            "wis": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 0,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 15
            }
          },
          "attributes": {
            "ac": {
              "calc": "unarmoredMonk",
              "flat": null,
              "formula": ""
            },
            "attunement": {
              "max": 3
            },
            "concentration": {
              "ability": "",
              "bonuses": {
                "save": ""
              },
              "limit": 1,
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              }
            },
            "death": {
              "bonuses": {
                "save": ""
              },
              "failure": 0,
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "success": 1
            },
            "exhaustion": 0,
            "hp": {
              "bonuses": {
                "level": "",
                "overall": ""
              },
              "max": 9,
              "temp": null,
              "tempmax": null,
              "value": 0
            },
            "init": {
              "ability": "",
              "bonus": "0",
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              }
            },
            "inspiration": false,
            "loyalty": {},
            "movement": {
              "burrow": 0,
              "climb": 0,
              "fly": 0,
              "hover": false,
              "swim": 0,
              "units": "ft",
              "walk": 25
            },
            "senses": {
              "blindsight": 0,
              "darkvision": 0,
              "special": "",
              "tremorsense": 0,
              "truesight": 0,
              "units": "ft"
            },
            "spellcasting": "int"
          },
          "bastion": {
            "description": "",
            "name": ""
          },
          "bonuses": {
            "abilities": {
              "check": "",
              "save": "",
              "skill": ""
            },
            "msak": {
              "attack": "",
              "damage": ""
            },
            "mwak": {
              "attack": "",
              "damage": ""
            },
            "rsak": {
              "attack": "",
              "damage": ""
            },
            "rwak": {
              "attack": "",
              "damage": ""
            },
            "spell": {
              "dc": ""
            }
          },
          "currency": {
            "cp": 0,
            "ep": 0,
            "gp": 15,
            "pp": 0,
            "sp": 0
          },
          "details": {
            "age": "",
            "alignment": "True Neutral",
            "appearance": "",
            "background": "q4tr1vTU8RxtU1UZ",
            "biography": {
              "public": "",
              "value": "<p>A devoted student of his local monastery, Perrin has been tasked with bringing balance to the surrounding lands.</p><p><em>Token artwork by <a href=\"https://www.forgotten-adventures.net/\" target=\"_blank\" rel=\"noopener\">Forgotten Adventures</a>.</em></p>"
            },
            "bond": "",
            "eyes": "",
            "faith": "",
            "flaw": "",
            "gender": "",
            "hair": "",
            "height": "",
            "ideal": "",
            "originalClass": "8Grf7ga6JcZF0X6x",
            "race": "O3ZjSw0GtAOPScHI",
            "skin": "",
            "trait": "",
            "weight": "",
            "xp": {
              "value": 0
            }
          },
          "favorites": [
            {
              "id": ".Item.FtOM4QiOW5MwgcS3",
              "sort": 500000,
              "type": "item"
            },
            {
              "id": ".Item.WeKJI3gPUAU52WAX",
              "sort": 600000,
              "type": "item"
            },
            {
              "id": ".Item.g7U3OAXVcoI4lwzf",
              "sort": 550000,
              "type": "item"
            },
            {
              "id": ".Item.GYyovoNdU66nxOlX",
              "sort": 700000,
              "type": "item"
            }
          ],
          "resources": {
            "primary": {
              "label": "Ki",
              "lr": true,
              "max": null,
              "sr": true,
              "value": null
            },
            "secondary": {
              "label": "",
              "lr": false,
              "max": null,
              "sr": false,
              "value": null
            },
            "tertiary": {
              "label": "",
              "lr": false,
              "max": null,
              "sr": false,
              "value": null
            }
          },
          "skills": {
            "acr": {
              "ability": "dex",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "ani": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "arc": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "ath": {
              "ability": "str",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "dec": {
              "ability": "cha",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "his": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "ins": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "inv": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "itm": {
              "ability": "cha",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "med": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "nat": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "per": {
              "ability": "cha",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "prc": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "prf": {
              "ability": "cha",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "rel": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "slt": {
              "ability": "dex",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "ste": {
              "ability": "dex",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "sur": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            }
          },
          "spells": {
            "pact": {
              "override": null,
              "value": 0
            },
            "spell0": {
              "override": null,
              "value": 0
            },
            "spell1": {
              "override": null,
              "value": 0
            },
            "spell2": {
              "override": null,
              "value": 0
            },
            "spell3": {
              "override": null,
              "value": 0
            },
            "spell4": {
              "override": null,
              "value": 0
            },
            "spell5": {
              "override": null,
              "value": 0
            },
            "spell6": {
              "override": null,
              "value": 0
            },
            "spell7": {
              "override": null,
              "value": 0
            },
            "spell8": {
              "override": null,
              "value": 0
            },
            "spell9": {
              "override": null,
              "value": 0
            }
          },
          "tools": {
            "art": {
              "ability": "int",
              "bonuses": {
                "check": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "brewer": {
              "ability": "int",
              "bonuses": {
                "check": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            }
          },
          "traits": {
            "armorProf": {
              "custom": "",
              "value": []
            },
            "ci": {
              "custom": "Advantage against being frightenned",
              "value": []
            },
            "di": {
              "bypasses": [],
              "custom": "",
              "value": []
            },
            "dm": {
              "amount": {},
              "bypasses": []
            },
            "dr": {
              "bypasses": [],
              "custom": "",
              "value": []
            },
            "dv": {
              "bypasses": [],
              "custom": "",
              "value": []
            },
            "languages": {
              "communication": {},
              "custom": "",
              "value": [
                "celestial",
                "draconic",
                "common",
                "halfling"
              ]
            },
            "size": "sm",
            "weaponProf": {
              "custom": "",
              "mastery": {
                "bonus": [],
                "value": []
              },
              "value": [
                "sim",
                "shortsword"
              ]
            }
          }
        },
        "type": "character",
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
          "createdTime": 1774367580132,
          "modifiedTime": 1774367606052,
          "lastModifiedBy": "r6bXhB7k9cXa3cif"
        },
        "_id": null
      },
      "dhd": 0,
      "dhp": 0,
      "longRest": false
    }
  }
}
```


---

## POST /dnd5e/long-rest

Perform a long rest for an actor. Triggers the D&D 5e long rest workflow including full HP recovery, spell slot restoration, hit dice recovery, and feature resets.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| clientId | string |  | body, query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor to rest (optional if selected is true) |
| selected | boolean |  | body, query | Use the currently selected token's actor |
| newDay | boolean |  | body, query | Whether the long rest marks a new day (default: true) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the long rest operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/long-rest"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"body"},{"name":"newDay","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/long-rest';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "newDay": true
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/long-rest?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","newDay":true}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/long-rest'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "newDay": True
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
  const path = '/dnd5e/long-rest';
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
        "newDay": true
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
  🔤/dnd5e/long-rest🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","newDay":true}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/long-rest🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 52❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "long-rest-result",
  "requestId": "long-rest_1774367606211",
  "data": {
    "actorUuid": "Actor.ioZexonJDGVuU8zl",
    "result": {
      "type": "long",
      "deltas": {
        "hitPoints": 9,
        "hitDice": 0
      },
      "updateData": {
        "system": {
          "attributes": {
            "hp": {
              "value": 9
            },
            "death": {
              "success": 0,
              "failure": 0
            }
          },
          "spells": {
            "pact": {
              "value": 0
            },
            "spell1": {
              "value": 0
            },
            "spell2": {
              "value": 0
            },
            "spell3": {
              "value": 0
            },
            "spell4": {
              "value": 0
            },
            "spell5": {
              "value": 0
            },
            "spell6": {
              "value": 0
            },
            "spell7": {
              "value": 0
            },
            "spell8": {
              "value": 0
            },
            "spell9": {
              "value": 0
            }
          }
        },
        "_id": "ioZexonJDGVuU8zl",
        "type": "character"
      },
      "updateItems": [],
      "newDay": true,
      "rolls": [],
      "clone": {
        "effects": [
          {
            "_id": "dnd5ebloodied000",
            "name": "Bloodied",
            "img": "systems/dnd5e/icons/svg/statuses/bloodied.svg",
            "statuses": [
              "bloodied"
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
              "createdTime": 1774367605935,
              "modifiedTime": 1774367605935,
              "lastModifiedBy": "r6bXhB7k9cXa3cif"
            }
          }
        ],
        "img": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
        "items": [
          {
            "_id": "q4tr1vTU8RxtU1UZ",
            "_stats": {
              "compendiumSource": null,
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": null,
            "img": "icons/sundries/documents/document-torn-diagram-tan.webp",
            "name": "Priest",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "advancement": [],
              "description": {
                "chat": "",
                "value": "<ul><li><strong>Skill Proficiencies:</strong> Insight, Religion</li><li><strong>Languages:</strong> Two of your choice</li><li><strong>Equipment:</strong> A holy symbol, 5 sticks of incense, prayer book, vestments, a set of common clothes, and a pouch containing 15 gp.</li></ul>"
              },
              "identifier": "priest",
              "source": {
                "book": "",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "startingEquipment": []
            },
            "type": "background"
          },
          {
            "_id": "O3ZjSw0GtAOPScHI",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.ZgYBjYYfiUstQD6f",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "uQNrdSpcXeoB5EnX",
            "img": "icons/equipment/feet/shoes-leather-simple-brown.webp",
            "name": "Lightfoot Halfling",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
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
              "description": {
                "chat": "",
                "value": "<p>Your halfling character has a number of traits in common with all other halflings.</p><p><em><strong>Ability Score Increase.</strong></em> Your Dexterity score increases by 2.</p><p><em><strong>Age.</strong></em> A halfling reaches adulthood at the age of 20 and generally lives into the middle of his or her second century.</p><p><em><strong>Alignment.</strong></em> Most halflings are lawful good. As a rule, they are good-hearted and kind, hate to see others in pain, and have no tolerance for oppression. They are also very orderly and traditional, leaning heavily on the support of their community and the comfort of their old ways.</p><p><em><strong>Size.</strong></em> Halflings average about 3 feet tall and weigh about 40 pounds. Your size is Small.</p><p><em><strong>Speed.</strong></em> Your base walking speed is 25 feet.</p><p><em><strong>Lucky.</strong></em> When you roll a 1 on the d20 for an attack roll, ability check, or saving throw, you can reroll the die and must use the new roll.</p><p><em><strong>Brave.</strong></em> You have advantage on saving throws against being frightened.</p><p><em><strong>Halfling Nimbleness.</strong></em> You can move through the space of any creature that is of a size larger than yours.</p><p><em><strong>Languages.</strong></em> You can speak, read, and write Common and Halfling. The Halfling language isn't secret, but halflings are loath to share it with others. They write very little, so they don't have a rich body of literature. Their oral tradition, however, is very strong. Almost all halflings speak Common to converse with the people in whose lands they dwell or through which they are traveling.</p><h5>Lightfoot</h5><p>As a lightfoot halfling, you can easily hide from notice, even using other people as cover. You're inclined to be affable and get along well with others.</p><p>Lightfoots are more prone to wanderlust than other halflings, and often dwell alongside other races or take up a nomadic life.</p><p><em><strong>Ability Score Increase.</strong></em> Your Charisma score increases by 1.</p><p><em><strong>Naturally Stealthy.</strong></em> You can attempt to hide even when you are obscured only by a creature that is at least one size larger than you.</p>"
              },
              "identifier": "lightfoot-halfling",
              "movement": {
                "burrow": null,
                "climb": null,
                "fly": null,
                "hover": false,
                "swim": null,
                "units": "ft",
                "walk": 25
              },
              "senses": {
                "blindsight": null,
                "darkvision": null,
                "special": "",
                "tremorsense": null,
                "truesight": null,
                "units": "ft"
              },
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "custom": "",
                "subtype": "halfling",
                "value": "humanoid"
              }
            },
            "type": "race"
          },
          {
            "_id": "FtOM4QiOW5MwgcS3",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.LOMdcNAGWh5xpfm4",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "kbtbKofcv13crhke",
            "img": "icons/sundries/gaming/dice-runed-brown.webp",
            "name": "Lucky",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>When you roll a 1 on the d20 for an attack roll, ability check, or saving throw, you can reroll the die and must use the new roll.</p><section class=\"secret foundry-note\" id=\"secret-S04TPyvUh05Dz0Ng\"><p><strong>Foundry Note</strong></p><p>This property can be enabled on your character sheet in the Special Traits configuration on the Attributes tab.</p></section>"
              },
              "enchant": {},
              "identifier": "lucky",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Halfling",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "race"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "nmmihiqphHjoE8dl",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.7Yoo9hG0hfPSmBoC",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "kbtbKofcv13crhke",
            "img": "icons/skills/melee/unarmed-punch-fist.webp",
            "name": "Brave",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>You have advantage on saving throws against being frightened.</p>"
              },
              "enchant": {},
              "identifier": "brave",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Halfling",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "race"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "cWrETHzCRs1Ueqd3",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.PqxZgcJzp1VVgP8t",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "kbtbKofcv13crhke",
            "img": "icons/skills/movement/feet-winged-boots-brown.webp",
            "name": "Halfling Nimbleness",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>You can move through the space of any creature that is of a size larger than yours.</p>"
              },
              "enchant": {},
              "identifier": "halfling-nimbleness",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Halfling",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "race"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "AArhiOrSkaQUnCZS",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.races.Item.GWPjKFeIthBBeCFJ",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "kbtbKofcv13crhke",
            "img": "icons/magic/perception/silhouette-stealth-shadow.webp",
            "name": "Naturally Stealthy",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>You can attempt to hide even when you are obscured only by a creature that is at least one size larger than you.</p>"
              },
              "enchant": {},
              "identifier": "naturally-stealthy",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Lightfoot Halfling",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "race"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "8Grf7ga6JcZF0X6x",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.classes.Item.6VoZrWxhOEKGYhnq",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "HQ1Oy7HkbnxnE63o",
            "img": "icons/skills/melee/hand-grip-staff-blue.webp",
            "name": "Monk",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
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
              "description": {
                "chat": "",
                "value": "<p>As a monk, you gain the following class features.</p><h3>Hit Points</h3><p><strong>Hit Dice:</strong> 1d8 per monk level<br /><strong>Hit Points at 1st Level:</strong> 8 + your Constitution modifier<br /><strong>Hit Points at Higher Levels:</strong> 1d8 (or 5) + your Constitution modifier per monk level after 1st</p><h3>Proficiencies</h3><p><strong>Armor:</strong> None<br /><strong>Weapons:</strong> Simple weapons, shortswords<br /><strong>Tools:</strong> Choose one type of artisan's tools or one musical instrument<br /><strong>Saving Throws:</strong> Strength, Dexterity<br /><strong>Skills:</strong> Choose two from Acrobatics, Athletics, History, Insight, Religion, and Stealth</p><h1>Monk Advancement</h1><table><thead><tr><td>Level</td><td>Proficiency Bonus</td><td>Martial Arts</td><td>Ki Points</td><td>Unarmored Movement</td><td>Features</td></tr></thead><tbody><tr><td>1st</td><td>+2</td><td>1d4</td><td>—</td><td>—</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI]{Unarmored Defense}, @UUID[Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw]{Martial Arts}</td></tr><tr><td>2nd</td><td>+2</td><td>1d4</td><td>2</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.10b6z2W1txNkrGP7]{Ki}, @UUID[Compendium.dnd5e.classfeatures.Item.zCeqyQ8uIPNdYJSW]{Unarmored Movement}</td></tr><tr><td>3rd</td><td>+2</td><td>1d4</td><td>3</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.rtpQdX77dYWbDIOH]{Monastic Tradition}, @UUID[Compendium.dnd5e.classfeatures.Item.mzweVbnsJPQiVkAe]{Deflect Missiles}</td></tr><tr><td>4th</td><td>+2</td><td>1d4</td><td>4</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}, @UUID[Compendium.dnd5e.classfeatures.Item.KQz9bqxVkXjDl8gK]{Slow Fall}</td></tr><tr><td>5th</td><td>+3</td><td>1d6</td><td>5</td><td>+10 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.XogoBnFWmCAHXppo]{Extra Attack}, @UUID[Compendium.dnd5e.classfeatures.Item.pvRc6GAu1ok6zihC]{Stunning Strike}</td></tr><tr><td>6th</td><td>+3</td><td>1d6</td><td>6</td><td>+15 ft.</td><td><p>@UUID[Compendium.dnd5e.classfeatures.Item.7flZKruSSu6dHg6D]{Ki-Empowered Strikes},</p><p>Monastic Tradition feature</p></td></tr><tr><td>7th</td><td>+3</td><td>1d6</td><td>7</td><td>+15 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.a4P4DNMmH8CqSNkC]{Evasion}, @UUID[Compendium.dnd5e.classfeatures.Item.ZmC31XKS4YNENnoc]{Stillness of Mind}</td></tr><tr><td>8th</td><td>+3</td><td>1d6</td><td>8</td><td>+15 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>9th</td><td>+4</td><td>1d6</td><td>9</td><td>+15 ft.</td><td>Unarmored Movement improvement</td></tr><tr><td>10th</td><td>+4</td><td>1d6</td><td>10</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.bqWA7t9pDELbNRkp]{Purity of Body}</td></tr><tr><td>11th</td><td>+4</td><td>1d8</td><td>11</td><td>+20 ft.</td><td>Monastic Tradition feature</td></tr><tr><td>12th</td><td>+4</td><td>1d8</td><td>12</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>13th</td><td>+5</td><td>1d8</td><td>13</td><td>+20 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.XjuGBeB8Y0C3A5D4]{Tongue of the Sun and Moon}</td></tr><tr><td>14th</td><td>+5</td><td>1d8</td><td>14</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.7D2EkLdISwShEDlN]{Diamond Soul}</td></tr><tr><td>15th</td><td>+5</td><td>1d8</td><td>15</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.gDH8PMrKvLHaNmEI]{Timeless Body}</td></tr><tr><td>16th</td><td>+5</td><td>1d8</td><td>16</td><td>+25 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>17th</td><td>+6</td><td>1d10</td><td>17</td><td>+25 ft.</td><td>Monastic Tradition feature</td></tr><tr><td>18th</td><td>+6</td><td>1d10</td><td>18</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.3jwFt3hSqDswBlOH]{Empty Body}</td></tr><tr><td>19th</td><td>+6</td><td>1d10</td><td>19</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.s0Cc2zcX0JzIgam5]{Ability Score Improvement}</td></tr><tr><td>20th</td><td>+6</td><td>1d10</td><td>20</td><td>+30 ft.</td><td>@UUID[Compendium.dnd5e.classfeatures.Item.mQNPg89YIs7g5tG4]{Perfect Self}</td></tr></tbody></table><h1>Monastic Traditions</h1><p>Three traditions of monastic pursuit are common in the monasteries scattered across the multiverse. Most monasteries practice one tradition exclusively, but a few honor the three traditions and instruct each monk according to his or her aptitude and interest. All three traditions rely on the same basic techniques, diverging as the student grows more adept. Thus, a monk need choose a tradition only upon reaching 3rd level.</p><p>@UUID[Compendium.dnd5e.subclasses.Item.IvlpKMXX3PmW1NY2]{Way of the Open Hand}</p>"
              },
              "hd": {
                "additional": "",
                "denomination": "d8",
                "spent": 0
              },
              "identifier": "monk",
              "levels": 1,
              "primaryAbility": {
                "all": true,
                "value": []
              },
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "spellcasting": {
                "ability": "",
                "preparation": {
                  "formula": ""
                },
                "progression": "none"
              },
              "startingEquipment": [
                {
                  "_id": "5yj0P4r9teJDnDtd",
                  "group": "",
                  "requiresProficiency": false,
                  "sort": 100000,
                  "type": "OR"
                },
                {
                  "_id": "R5tuRtaPonfjQCVU",
                  "count": null,
                  "group": "5yj0P4r9teJDnDtd",
                  "key": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
                  "requiresProficiency": false,
                  "sort": 200000,
                  "type": "linked"
                },
                {
                  "_id": "Mlf6kel8ws6xgDER",
                  "count": null,
                  "group": "5yj0P4r9teJDnDtd",
                  "key": "simpleM",
                  "requiresProficiency": false,
                  "sort": 300000,
                  "type": "weapon"
                },
                {
                  "_id": "3TbVLmLPtjVaSh5O",
                  "group": "",
                  "requiresProficiency": false,
                  "sort": 400000,
                  "type": "OR"
                },
                {
                  "_id": "AvDYtl0uvQsDuhnb",
                  "count": null,
                  "group": "3TbVLmLPtjVaSh5O",
                  "key": "Compendium.dnd5e.items.Item.XY8b594Dn7plACLL",
                  "requiresProficiency": false,
                  "sort": 500000,
                  "type": "linked"
                },
                {
                  "_id": "4QKQURCmIurbTAzp",
                  "count": null,
                  "group": "3TbVLmLPtjVaSh5O",
                  "key": "Compendium.dnd5e.items.Item.8KWz5DJbWUpNWniP",
                  "requiresProficiency": false,
                  "sort": 600000,
                  "type": "linked"
                },
                {
                  "_id": "AOYuulsULvsHbSLO",
                  "count": 10,
                  "group": "",
                  "key": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
                  "requiresProficiency": false,
                  "sort": 700000,
                  "type": "linked"
                }
              ],
              "wealth": "5d4"
            },
            "type": "class"
          },
          {
            "_id": "CwgoTDXWCD7PknIN",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [
              {
                "_id": "R5ro4AuNjcdWD56O",
                "_stats": {
                  "compendiumSource": null,
                  "duplicateSource": null,
                  "exportSource": null,
                  "coreVersion": "13.348",
                  "systemId": "dnd5e",
                  "systemVersion": "5.0.4",
                  "lastModifiedBy": null
                },
                "changes": [
                  {
                    "key": "system.attributes.ac.calc",
                    "mode": 5,
                    "priority": null,
                    "value": "unarmoredMonk"
                  }
                ],
                "description": "",
                "disabled": false,
                "duration": {
                  "combat": null,
                  "rounds": null,
                  "seconds": null,
                  "startRound": null,
                  "startTime": 0,
                  "startTurn": null,
                  "turns": null
                },
                "flags": {},
                "img": "icons/magic/control/silhouette-hold-change-blue.webp",
                "name": "Unarmored Defense",
                "origin": "Item.cOdcNWy4hII029DT",
                "sort": 0,
                "statuses": [],
                "system": {},
                "tint": "#ffffff",
                "transfer": true,
                "type": "base"
              }
            ],
            "flags": {
              "dnd5e": {
                "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
                "sourceId": "Compendium.dnd5e.classfeatures.Item.UAvV7N7T4zJhxdfI"
              }
            },
            "folder": "TMmNG8ujFDBEWXRe",
            "img": "icons/magic/control/silhouette-hold-change-blue.webp",
            "name": "Unarmored Defense",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>Beginning at 1st Level, while you are wearing no armor and not wielding a Shield, your AC equals 10 + your Dexterity modifier + your Wisdom modifier.</p>"
              },
              "enchant": {},
              "identifier": "unarmored-defense",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Monk 1",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "class"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "pchnXqd5C79fVlxy",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {
              "dnd5e": {
                "advancementOrigin": "8Grf7ga6JcZF0X6x.n0q8XyiGA3vLPgpK",
                "sourceId": "Compendium.dnd5e.classfeatures.Item.l50hjTxO2r0iecKw"
              }
            },
            "folder": "TMmNG8ujFDBEWXRe",
            "img": "icons/skills/melee/unarmed-punch-fist.webp",
            "name": "Martial Arts",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "advancement": [],
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>At 1st level, your practice of martial arts gives you mastery of combat styles that use and monk weapons, which are shortswords and any simple melee weapons that don't have the two-handed or heavy property. You gain the following benefits while you are unarmed or wielding only monk weapons and you aren't wearing armor or wielding a shield:</p>\n<ul>\n<li>\n<p>You can use Dexterity instead of Strength for the attack and damage rolls of your unarmed strikes and monk weapons.</p>\n</li>\n<li>\n<p>You can roll a d4 in place of the normal damage of your unarmed strike or monk weapon. This die changes as you gain monk levels, as shown in the Martial Arts column of the Monk table.</p>\n</li>\n<li>\n<p>When you use the Attack action with an unarmed strike or a monk weapon on your turn, you can make one unarmed strike as a bonus action. For example, if you take the Attack action and attack with a quarterstaff, you can also make an unarmed strike as a bonus action, assuming you haven't already taken a bonus action this turn.</p>\n</li>\n</ul>\n<p>Certain monasteries use specialized forms of the monk weapons. For example, you might use a club that is two lengths of wood connected by a short chain (called a nunchaku) or a sickle with a shorter, straighter blade (called a kama). Whatever name you use for a monk weapon, you can use the game statistics provided for the weapon.</p>"
              },
              "enchant": {},
              "identifier": "martial-arts",
              "prerequisites": {
                "level": null,
                "repeatable": false
              },
              "properties": [],
              "requirements": "Monk 1",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "class"
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "feat"
          },
          {
            "_id": "YJ1P3PnFKHOdQpaP",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.14pNRT4sZy9rgvhb",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "dlru9Hy74nSMv6fr",
            "img": "icons/tools/hand/hammer-cobbler-steel.webp",
            "name": "Hammer",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": null,
              "description": {
                "chat": "",
                "value": "<p>A tool with a heavy metal head mounted at the end of its handle, used for jobs such as breaking things and driving in nails. </p>"
              },
              "identified": true,
              "identifier": "hammer",
              "price": {
                "denomination": "gp",
                "value": 1
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 3
              }
            },
            "type": "loot"
          },
          {
            "_id": "DDnxRCeYUhXstWU8",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.1FSubnBpSTDmVaYV",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/sundries/lights/torch-black.webp",
            "name": "Tinderbox",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": "8KWz5DJbWUpNWniP",
              "description": {
                "chat": "",
                "value": "<p>This small container holds flint, fire steel, and tinder (usually dry cloth soaked in light oil) used to kindle a fire. Using it to light a torch - or anything else with abundant, exposed fuel - takes an action. Lighting any other fire takes 1 minute.</p>\n<p> </p>"
              },
              "identified": true,
              "identifier": "tinderbox",
              "price": {
                "denomination": "sp",
                "value": 5
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 1
              }
            },
            "type": "loot"
          },
          {
            "_id": "5skKSSB4ShHbKoc8",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.1L5wkmbw0fmNAr38",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": "r6bXhB7k9cXa3cif",
              "modifiedTime": 1774367605451
            },
            "effects": [],
            "flags": {
              "dnd5e": {
                "migratedUses": 4,
                "riders": {
                  "activity": [],
                  "effect": []
                }
              }
            },
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/sundries/survival/wetskin-leather-purple.webp",
            "name": "Waterskin",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": "6OYR11aJX2dEVtOj",
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>A leather hide sewn into an enclosed skin which can contain up to 4 pints of liquid. It weighs 5 pounds when full; a pint of water is approximately 1 pound.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "waterskin",
              "magicalBonus": null,
              "price": {
                "denomination": "sp",
                "value": 2
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "food"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": false,
                "max": "4",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 5
              }
            },
            "type": "consumable"
          },
          {
            "_id": "dDuMscUuMI2bTdkj",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.29ZLE8PERtFVD3QU",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/sundries/lights/torch-black.webp",
            "name": "Torch",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": "8KWz5DJbWUpNWniP",
              "damage": {
                "base": {
                  "bonus": "",
                  "custom": {
                    "enabled": true,
                    "formula": "1"
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": [
                    "fire"
                  ]
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>A torch burns for 1 hour, providing bright light in a 20-foot radius and dim light for an additional 20 feet. If you make a melee attack with a burning torch and hit, it deals 1 fire damage.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "torch",
              "magicalBonus": null,
              "price": {
                "denomination": "cp",
                "value": 1
              },
              "properties": [],
              "quantity": 10,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "trinket"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": false,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 1
              }
            },
            "type": "consumable"
          },
          {
            "_id": "nC6VcR5JAIbR4err",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.3b0RvGi0TnTYpIxn",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "dlru9Hy74nSMv6fr",
            "img": "icons/consumables/grains/breadsticks-crackers-wrapped-ration-brown.webp",
            "name": "Stick of Incense",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": null,
              "description": {
                "chat": "",
                "value": "<p>When blocks of incense cannot be used or a cheaper alternative is required, people often use these to perfume the air, whether for pleasurable or religious purposes.</p>"
              },
              "identified": true,
              "identifier": "stick-of-incense",
              "price": {
                "denomination": "sp",
                "value": 2
              },
              "properties": [],
              "quantity": 5,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 0
              }
            },
            "type": "loot"
          },
          {
            "_id": "WeKJI3gPUAU52WAX",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.3rCO8MTIdPGSW6IJ",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "MLMTCAvKsuFE3vYA",
            "img": "icons/weapons/thrown/dart-feathered.webp",
            "name": "Dart",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "ammunition": {},
              "armor": {
                "value": 10
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "damage": {
                "base": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": 4,
                  "number": 1,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": [
                    "piercing"
                  ]
                },
                "versatile": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": []
                }
              },
              "description": {
                "chat": "",
                "value": "<p>A small thrown implement crafted with a short wooden shaft and crossed feathres with a sharp wooden or metal tip. Darts can be thrown with sufficient force to puncture the skin.</p>"
              },
              "equipped": true,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "dart",
              "magicalBonus": null,
              "mastery": "",
              "price": {
                "denomination": "cp",
                "value": 5
              },
              "proficient": null,
              "properties": [
                "fin",
                "thr"
              ],
              "quantity": 10,
              "range": {
                "long": 60,
                "reach": null,
                "units": "ft",
                "value": 20
              },
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "baseItem": "dart",
                "value": "simpleR"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 0.25
              }
            },
            "type": "weapon"
          },
          {
            "_id": "1F73YcUHbZMgePAD",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.8RXjiddJ6VGyE7vB",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "aJgMxnZED9XdoN2W",
            "img": "icons/equipment/chest/shirt-collared-brown.webp",
            "name": "Common Clothes",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "armor": {
                "dex": null,
                "magicalBonus": null,
                "value": null
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>Clothes worn by most commoners.</p>"
              },
              "equipped": false,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "common-clothes",
              "price": {
                "denomination": "sp",
                "value": 5
              },
              "proficient": null,
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "speed": {
                "conditions": "",
                "value": null
              },
              "strength": null,
              "type": {
                "baseItem": "",
                "value": "clothing"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 3
              }
            },
            "type": "equipment"
          },
          {
            "_id": "LQhGSEatJ3VK7oqW",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.8d95YV1jHcxPygJ9",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/consumables/grains/bread-loaf-boule-rustic-brown.webp",
            "name": "Rations",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": "XY8b594Dn7plACLL",
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>Rations consist of dry foods suitable for extended travel, including jerky, dried fruit, hardtack, and nuts.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "rations",
              "magicalBonus": null,
              "price": {
                "denomination": "sp",
                "value": 5
              },
              "properties": [],
              "quantity": 10,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "food"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": true,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 2
              }
            },
            "type": "consumable"
          },
          {
            "_id": "1kqMwSwjfErSFNvl",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.9bWTRRDym06PzSAf",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "PgtyWdfkDTDbc0ov",
            "img": "icons/containers/bags/pouch-rounded-leather-gold-tan.webp",
            "name": "Pouch",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "attuned": false,
              "attunement": "",
              "capacity": {
                "volume": {
                  "units": "cubicFoot"
                },
                "weight": {
                  "units": "lb",
                  "value": 6
                }
              },
              "container": null,
              "currency": {
                "cp": 0,
                "ep": 0,
                "gp": 0,
                "pp": 0,
                "sp": 0
              },
              "description": {
                "chat": "",
                "value": "<p>A cloth or leather pouch can hold up to 20 sling bullets or 50 blowgun needles, among other things. A compartmentalized pouch for holding spell components is called a component pouch. A pouch can hold up to ⅕ cubic foot or 6 pounds of gear.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "pouch",
              "price": {
                "denomination": "sp",
                "value": 5
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 1
              }
            },
            "type": "container"
          },
          {
            "_id": "zF5L4xHnJBC7b2iM",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.AkyQyonZMVcvOrXU",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/tools/hand/pickaxe-steel-white.webp",
            "name": "Crowbar",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": "XY8b594Dn7plACLL",
              "description": {
                "chat": "",
                "value": "<p>Using a crowbar grants advantage to Strength checks where the crowbar's leverage can be applied.</p>"
              },
              "identified": true,
              "identifier": "crowbar",
              "price": {
                "denomination": "gp",
                "value": 2
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 5
              }
            },
            "type": "loot"
          },
          {
            "_id": "g7U3OAXVcoI4lwzf",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.GsuvwoekKZatfKwF",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
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
            "folder": "MLMTCAvKsuFE3vYA",
            "img": "icons/skills/melee/unarmed-punch-fist.webp",
            "name": "Unarmed Strike",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "ammunition": {},
              "armor": {
                "value": 10
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "damage": {
                "base": {
                  "bonus": "",
                  "custom": {
                    "enabled": true,
                    "formula": "@scale.monk.die"
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": [
                    "bludgeoning"
                  ]
                },
                "versatile": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": []
                }
              },
              "description": {
                "chat": "",
                "value": "<p>A  punch, kick, head-butt, or similar forceful blow (none of which count as weapons). On a hit, an unarmed strike deals bludgeoning damage equal to 1 + your Strength modifier. You are proficient with your unarmed strikes.</p>"
              },
              "equipped": true,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "unarmed-strike",
              "magicalBonus": null,
              "mastery": "",
              "price": {
                "denomination": "gp",
                "value": 0
              },
              "proficient": null,
              "properties": [
                "fin"
              ],
              "quantity": 1,
              "range": {
                "long": null,
                "reach": null,
                "units": "ft",
                "value": null
              },
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "baseItem": "",
                "value": "simpleM"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 0
              }
            },
            "type": "weapon"
          },
          {
            "_id": "WcROGD590imVj9qp",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.H8YCd689ezlD26aT",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "PgtyWdfkDTDbc0ov",
            "img": "icons/containers/bags/pack-leather-white-tan.webp",
            "name": "Backpack",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "attuned": false,
              "attunement": "",
              "capacity": {
                "volume": {
                  "units": "cubicFoot"
                },
                "weight": {
                  "units": "lb",
                  "value": 30
                }
              },
              "container": null,
              "currency": {
                "cp": 0,
                "ep": 0,
                "gp": 0,
                "pp": 0,
                "sp": 0
              },
              "description": {
                "chat": "",
                "value": "<p>A backpack can hold one cubic foot or 30 pounds of gear. You can also strap items, such as a bedroll or a coil of rope, to the outside of a backpack.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "backpack",
              "price": {
                "denomination": "gp",
                "value": 2
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 5
              }
            },
            "type": "container"
          },
          {
            "_id": "8UQanAvHUIHZXp0O",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.P31t6tGgt9aLAdYt",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "Dx3K2y0J1wJUPP9m",
            "img": "icons/tools/fasteners/nail-steel.webp",
            "name": "Piton",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": "xsB7Y2WI476kvOt4",
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>A metal spike that is drive into a seam in a climbing surface with a climbing hammer. It can also be used like iron spikes to block doors in a pinch.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "piton",
              "magicalBonus": null,
              "price": {
                "denomination": "cp",
                "value": 5
              },
              "properties": [],
              "quantity": 10,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "trinket"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": false,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 0.25
              }
            },
            "type": "consumable"
          },
          {
            "_id": "irCoxLHje0eA6Qmu",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.QXmaarJ4X8P0C1HV",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {
              "dnd5e": {
                "riders": {
                  "activity": [],
                  "effect": []
                }
              }
            },
            "folder": "UnUwTG4YIgd0kaUJ",
            "img": "icons/sundries/survival/rope-wrapped-brown.webp",
            "name": "Hempen Rope (50 ft.)",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p>Rope, whether made of hemp or silk, has 2 hit points and can be burst with a DC 17 Strength check.</p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "hempen-rope-50-ft",
              "magicalBonus": null,
              "price": {
                "denomination": "gp",
                "value": 1
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "trinket"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": false,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 10
              }
            },
            "type": "consumable"
          },
          {
            "_id": "GfKuJYLFfw00oW9R",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.gP1URGq3kVIIFHJ7",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {
              "dnd5e": {
                "riders": {
                  "activity": [],
                  "effect": []
                }
              }
            },
            "folder": "xedn1r43VWuEBcli",
            "img": "icons/containers/chest/chest-reinforced-steel-red.webp",
            "name": "Reliquary",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {},
              "armor": {
                "dex": null,
                "magicalBonus": null,
                "value": null
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>A tiny box or other container holding a fragment of a precious relic, saint, or other historical figure that dedicated their life to walk the path of a true believer. A deity imbues the bearer of this artifact with the ability to call forth their power and in doing so spread the faith once more.</p>\n<p><strong>Spellcasting Focus</strong>. A cleric or paladin can use a holy symbol as a spellcasting focus. To use the symbol in this way, the caster must hold it in hand, wear it visibly, or bear it on a shield.</p>"
              },
              "equipped": false,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "reliquary",
              "price": {
                "denomination": "gp",
                "value": 5
              },
              "proficient": null,
              "properties": [
                "foc"
              ],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "speed": {
                "conditions": "",
                "value": null
              },
              "strength": null,
              "type": {
                "baseItem": "",
                "value": "trinket"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 2
              }
            },
            "type": "equipment"
          },
          {
            "_id": "NsNcJBUG5ajbk7sG",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.irtqrzaUCeshmTZp",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "aJgMxnZED9XdoN2W",
            "img": "icons/equipment/back/mantle-collared-black.webp",
            "name": "Vestments",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "armor": {
                "dex": null,
                "magicalBonus": null,
                "value": null
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "description": {
                "chat": "",
                "value": "<p>Simple or ostentacious wear, often used by priests and other religious figures for use in rituals and ceremonies.</p>"
              },
              "equipped": true,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "vestments",
              "price": {
                "denomination": "gp",
                "value": 1
              },
              "proficient": null,
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "speed": {
                "conditions": "",
                "value": null
              },
              "strength": null,
              "type": {
                "baseItem": "",
                "value": "clothing"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 4
              }
            },
            "type": "equipment"
          },
          {
            "_id": "XbF0CTipyqbMKHsB",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.odV5cq2HSLSCH69k",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "dlru9Hy74nSMv6fr",
            "img": "icons/sundries/books/book-purple-cross.webp",
            "name": "Prayer Book",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "container": null,
              "description": {
                "chat": "",
                "value": "<p>A book containing prayers and incantations dedicated to a specific power for the faithful to follow.</p>"
              },
              "identified": true,
              "identifier": "prayer-book",
              "price": {
                "denomination": "gp",
                "value": 25
              },
              "properties": [],
              "quantity": 1,
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": ""
              },
              "unidentified": {
                "description": ""
              },
              "weight": {
                "units": "lb",
                "value": 5
              }
            },
            "type": "loot"
          },
          {
            "_id": "GYyovoNdU66nxOlX",
            "_stats": {
              "compendiumSource": "Compendium.dnd5e.items.Item.osLzOwQdPtrK3rQH",
              "duplicateSource": null,
              "exportSource": null,
              "coreVersion": "13.348",
              "systemId": "dnd5e",
              "systemVersion": "5.0.4",
              "lastModifiedBy": null
            },
            "effects": [],
            "flags": {},
            "folder": "MLMTCAvKsuFE3vYA",
            "img": "icons/weapons/swords/sword-guard-worn-purple.webp",
            "name": "Shortsword",
            "ownership": {
              "default": 0
            },
            "sort": 0,
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
              },
              "ammunition": {},
              "armor": {
                "value": 10
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "cover": null,
              "crewed": false,
              "damage": {
                "base": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": 6,
                  "number": 1,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": [
                    "piercing"
                  ]
                },
                "versatile": {
                  "bonus": "",
                  "custom": {
                    "enabled": false,
                    "formula": ""
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "formula": "",
                    "mode": "",
                    "number": null
                  },
                  "types": []
                }
              },
              "description": {
                "chat": "",
                "value": "<p>A medium sized blade with a firm crossguard and a leather-wrapped handle. A versatile weapon which makes up in versatility what it lacks in reach.</p>"
              },
              "equipped": true,
              "hp": {
                "conditions": "",
                "dt": null,
                "max": 0,
                "value": 0
              },
              "identified": true,
              "identifier": "shortsword",
              "magicalBonus": null,
              "price": {
                "denomination": "gp",
                "value": 10
              },
              "proficient": null,
              "properties": [
                "fin",
                "lgt"
              ],
              "quantity": 1,
              "range": {
                "long": null,
                "reach": null,
                "units": "ft",
                "value": null
              },
              "rarity": "",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "baseItem": "shortsword",
                "value": "martialM"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 2
              }
            },
            "type": "weapon"
          },
          {
            "effects": [
              {
                "_id": "tM8Yj0dw52paxVc4",
                "_stats": {
                  "compendiumSource": null,
                  "duplicateSource": null,
                  "exportSource": null,
                  "coreVersion": "13.348",
                  "systemId": "dnd5e",
                  "systemVersion": "5.0.4",
                  "lastModifiedBy": null
                },
                "changes": [],
                "description": "",
                "disabled": true,
                "duration": {
                  "combat": null,
                  "rounds": null,
                  "seconds": 3600,
                  "startRound": null,
                  "startTime": null,
                  "startTurn": null,
                  "turns": null
                },
                "flags": {},
                "img": "icons/magic/control/energy-stream-link-large-teal.webp",
                "name": "Polymorph",
                "origin": "Compendium.dnd5e.spells.Item.04nMsTWkIFvkbXlY",
                "sort": 0,
                "statuses": [
                  "transformed"
                ],
                "system": {},
                "tint": "#ffffff",
                "transfer": false,
                "type": "base"
              }
            ],
            "img": "icons/magic/control/energy-stream-link-large-teal.webp",
            "name": "test-polymorph",
            "system": {
              "activation": {
                "condition": "",
                "type": "action",
                "value": 1
              },
              "activities": {
                "dnd5eactivity000": {
                  "_id": "dnd5eactivity000",
                  "activation": {
                    "override": false,
                    "type": "action",
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
                  "damage": {
                    "onSave": "half",
                    "parts": []
                  },
                  "description": {
                    "chatFlavor": ""
                  },
                  "duration": {
                    "concentration": false,
                    "override": false,
                    "units": "inst"
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
                  "save": {
                    "ability": [
                      "wis"
                    ],
                    "dc": {
                      "calculation": "spellcasting",
                      "formula": ""
                    }
                  },
                  "sort": 0,
                  "target": {
                    "affects": {
                      "choice": false
                    },
                    "override": false,
                    "prompt": true,
                    "template": {
                      "contiguous": false,
                      "units": "ft"
                    }
                  },
                  "type": "save",
                  "uses": {
                    "recovery": [],
                    "spent": 0
                  }
                }
              },
              "description": {
                "chat": "",
                "value": "<p>This spell transforms a creature that you can see within range into a new form. An unwilling creature must make a Wisdom saving throw to avoid the effect. The spell has no effect on a Shapechanger or a creature with 0 Hit Points.</p><p>The transformation lasts for the Duration, or until the target drops to 0 Hit Points or dies. The new form can be any beast whose Challenge rating is equal to or less than the target's (or the target's level, if it doesn't have a challenge rating). The target's game statistics, including mental Ability Scores, are replaced by the Statistics of the chosen beast. It retains its Alignment and personality.</p><p>The target assumes the Hit Points of its new form. When it reverts to its normal form, the creature returns to the number of hit points it had before it transformed. If it reverts as a result of dropping to 0 hit points, any excess damage carries over to its normal form. As long as the excess damage doesn't reduce the creature's normal form to 0 hit points, it isn't knocked Unconscious.</p><p>The creature is limited in the Actions it can perform by the Nature of its new form, and it can't speak, cast Spells, or take any other action that requires hands or speech.</p><p>The target's gear melds into the new form. The creature can't activate, use, wield, or otherwise benefit from any of its Equipment.</p>"
              },
              "duration": {
                "units": "hour",
                "value": "1"
              },
              "identifier": "polymorph",
              "level": 4,
              "materials": {
                "consumed": false,
                "cost": 0,
                "supply": 0,
                "value": "A caterpillar cocoon"
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
              "range": {
                "units": "ft",
                "value": "60"
              },
              "school": "trs",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "target": {
                "affects": {
                  "choice": false,
                  "count": "1",
                  "type": "creature"
                },
                "template": {
                  "contiguous": false,
                  "units": ""
                }
              },
              "uses": {
                "max": "",
                "recovery": [],
                "spent": 0
              }
            },
            "type": "spell",
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
              "createdTime": 1774367604775,
              "modifiedTime": 1774367604775,
              "lastModifiedBy": "r6bXhB7k9cXa3cif"
            },
            "_id": "cuN0zv7OrhLW1LKO"
          },
          {
            "effects": [],
            "img": "icons/consumables/potions/bottle-bulb-corked-labeled-blue.webp",
            "name": "test-potion of storm giant strength",
            "system": {
              "activities": {
                "dnd5eactivity000": {
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
                    "units": "hour",
                    "value": "1"
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
              },
              "attuned": false,
              "attunement": "",
              "container": null,
              "damage": {
                "base": {
                  "custom": {
                    "enabled": false
                  },
                  "denomination": null,
                  "number": null,
                  "scaling": {
                    "number": 1
                  },
                  "types": []
                },
                "replace": false
              },
              "description": {
                "chat": "",
                "value": "<p><em>This potion's transparent liquid has floating in it a sliver of fingernail from a giant of the appropriate type.</em></p>\n<p>When you drink this potion, your Strength score changes to 29 for 1 hour.  The potion has no effect on you if your Strength is equal to or greater than that score.</p>\n<p> </p>"
              },
              "equipped": false,
              "identified": true,
              "identifier": "potion-of-storm-giant-strength",
              "magicalBonus": null,
              "price": {
                "denomination": "gp",
                "value": 2000
              },
              "properties": [
                "mgc"
              ],
              "quantity": 1,
              "rarity": "legendary",
              "source": {
                "book": "SRD 5.1",
                "custom": "",
                "license": "CC-BY-4.0",
                "page": "",
                "revision": 1,
                "rules": "2014"
              },
              "type": {
                "subtype": "",
                "value": "potion"
              },
              "unidentified": {
                "description": ""
              },
              "uses": {
                "autoDestroy": true,
                "max": "1",
                "recovery": [],
                "spent": 0
              },
              "weight": {
                "units": "lb",
                "value": 0.1
              }
            },
            "type": "consumable",
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
              "createdTime": 1774367604969,
              "modifiedTime": 1774367604969,
              "lastModifiedBy": "r6bXhB7k9cXa3cif"
            },
            "_id": "xzMsr0JNmiadK4Cr"
          }
        ],
        "name": "Updated Test Actor",
        "prototypeToken": {
          "actorLink": true,
          "alpha": 1,
          "appendNumber": false,
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
          "prependAdjective": false,
          "randomImg": false,
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
          "width": 1
        },
        "system": {
          "abilities": {
            "cha": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 0,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 11
            },
            "con": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 0,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 13
            },
            "dex": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 1,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 16
            },
            "int": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 0,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 12
            },
            "str": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 1,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 8
            },
            "wis": {
              "bonuses": {
                "check": "",
                "save": ""
              },
              "check": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "max": 20,
              "proficient": 0,
              "save": {
                "roll": {
                  "max": null,
                  "min": null,
                  "mode": 0
                }
              },
              "value": 15
            }
          },
          "attributes": {
            "ac": {
              "calc": "unarmoredMonk",
              "flat": null,
              "formula": ""
            },
            "attunement": {
              "max": 3
            },
            "concentration": {
              "ability": "",
              "bonuses": {
                "save": ""
              },
              "limit": 1,
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              }
            },
            "death": {
              "bonuses": {
                "save": ""
              },
              "failure": 0,
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "success": 1
            },
            "exhaustion": 0,
            "hp": {
              "bonuses": {
                "level": "",
                "overall": ""
              },
              "max": 9,
              "temp": null,
              "tempmax": null,
              "value": 0
            },
            "init": {
              "ability": "",
              "bonus": "0",
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              }
            },
            "inspiration": false,
            "loyalty": {},
            "movement": {
              "burrow": 0,
              "climb": 0,
              "fly": 0,
              "hover": false,
              "swim": 0,
              "units": "ft",
              "walk": 25
            },
            "senses": {
              "blindsight": 0,
              "darkvision": 0,
              "special": "",
              "tremorsense": 0,
              "truesight": 0,
              "units": "ft"
            },
            "spellcasting": "int"
          },
          "bastion": {
            "description": "",
            "name": ""
          },
          "bonuses": {
            "abilities": {
              "check": "",
              "save": "",
              "skill": ""
            },
            "msak": {
              "attack": "",
              "damage": ""
            },
            "mwak": {
              "attack": "",
              "damage": ""
            },
            "rsak": {
              "attack": "",
              "damage": ""
            },
            "rwak": {
              "attack": "",
              "damage": ""
            },
            "spell": {
              "dc": ""
            }
          },
          "currency": {
            "cp": 0,
            "ep": 0,
            "gp": 15,
            "pp": 0,
            "sp": 0
          },
          "details": {
            "age": "",
            "alignment": "True Neutral",
            "appearance": "",
            "background": "q4tr1vTU8RxtU1UZ",
            "biography": {
              "public": "",
              "value": "<p>A devoted student of his local monastery, Perrin has been tasked with bringing balance to the surrounding lands.</p><p><em>Token artwork by <a href=\"https://www.forgotten-adventures.net/\" target=\"_blank\" rel=\"noopener\">Forgotten Adventures</a>.</em></p>"
            },
            "bond": "",
            "eyes": "",
            "faith": "",
            "flaw": "",
            "gender": "",
            "hair": "",
            "height": "",
            "ideal": "",
            "originalClass": "8Grf7ga6JcZF0X6x",
            "race": "O3ZjSw0GtAOPScHI",
            "skin": "",
            "trait": "",
            "weight": "",
            "xp": {
              "value": 0
            }
          },
          "favorites": [
            {
              "id": ".Item.FtOM4QiOW5MwgcS3",
              "sort": 500000,
              "type": "item"
            },
            {
              "id": ".Item.WeKJI3gPUAU52WAX",
              "sort": 600000,
              "type": "item"
            },
            {
              "id": ".Item.g7U3OAXVcoI4lwzf",
              "sort": 550000,
              "type": "item"
            },
            {
              "id": ".Item.GYyovoNdU66nxOlX",
              "sort": 700000,
              "type": "item"
            }
          ],
          "resources": {
            "primary": {
              "label": "Ki",
              "lr": true,
              "max": null,
              "sr": true,
              "value": null
            },
            "secondary": {
              "label": "",
              "lr": false,
              "max": null,
              "sr": false,
              "value": null
            },
            "tertiary": {
              "label": "",
              "lr": false,
              "max": null,
              "sr": false,
              "value": null
            }
          },
          "skills": {
            "acr": {
              "ability": "dex",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "ani": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "arc": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "ath": {
              "ability": "str",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "dec": {
              "ability": "cha",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "his": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "ins": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "inv": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "itm": {
              "ability": "cha",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "med": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "nat": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "per": {
              "ability": "cha",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "prc": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "prf": {
              "ability": "cha",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "rel": {
              "ability": "int",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "slt": {
              "ability": "dex",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "ste": {
              "ability": "dex",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            },
            "sur": {
              "ability": "wis",
              "bonuses": {
                "check": "",
                "passive": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 0
            }
          },
          "spells": {
            "pact": {
              "override": null,
              "value": 0
            },
            "spell0": {
              "override": null,
              "value": 0
            },
            "spell1": {
              "override": null,
              "value": 0
            },
            "spell2": {
              "override": null,
              "value": 0
            },
            "spell3": {
              "override": null,
              "value": 0
            },
            "spell4": {
              "override": null,
              "value": 0
            },
            "spell5": {
              "override": null,
              "value": 0
            },
            "spell6": {
              "override": null,
              "value": 0
            },
            "spell7": {
              "override": null,
              "value": 0
            },
            "spell8": {
              "override": null,
              "value": 0
            },
            "spell9": {
              "override": null,
              "value": 0
            }
          },
          "tools": {
            "art": {
              "ability": "int",
              "bonuses": {
                "check": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            },
            "brewer": {
              "ability": "int",
              "bonuses": {
                "check": ""
              },
              "roll": {
                "max": null,
                "min": null,
                "mode": 0
              },
              "value": 1
            }
          },
          "traits": {
            "armorProf": {
              "custom": "",
              "value": []
            },
            "ci": {
              "custom": "Advantage against being frightenned",
              "value": []
            },
            "di": {
              "bypasses": [],
              "custom": "",
              "value": []
            },
            "dm": {
              "amount": {},
              "bypasses": []
            },
            "dr": {
              "bypasses": [],
              "custom": "",
              "value": []
            },
            "dv": {
              "bypasses": [],
              "custom": "",
              "value": []
            },
            "languages": {
              "communication": {},
              "custom": "",
              "value": [
                "celestial",
                "draconic",
                "common",
                "halfling"
              ]
            },
            "size": "sm",
            "weaponProf": {
              "custom": "",
              "mastery": {
                "bonus": [],
                "value": []
              },
              "value": [
                "sim",
                "shortsword"
              ]
            }
          }
        },
        "type": "character",
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
          "createdTime": 1774367580132,
          "modifiedTime": 1774367606052,
          "lastModifiedBy": "r6bXhB7k9cXa3cif"
        },
        "_id": null
      },
      "dhp": 9,
      "dhd": 0,
      "longRest": true
    }
  }
}
```


---

## POST /dnd5e/skill-check

Roll a skill check for an actor. Rolls a D&D 5e skill check with all applicable modifiers including proficiency, expertise, Jack of All Trades, and conditional bonuses.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| actorUuid | string | ✓ | body, query | UUID of the actor rolling the check |
| skill | string | ✓ | body, query | Skill abbreviation (e.g., "acr", "ath", "ste", "prc") |
| clientId | string |  | body, query | Client ID for the Foundry world |
| advantage | boolean |  | body, query | Roll with advantage |
| disadvantage | boolean |  | body, query | Roll with disadvantage |
| bonus | string |  | body, query | Extra bonus formula to add (e.g., "1d4", "+2") |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the skill check roll

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/skill-check"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"skill","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"disadvantage","type":"boolean","required":false,"source":"body"},{"name":"bonus","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/skill-check';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "skill": "prc"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/skill-check?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","skill":"prc"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/skill-check'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "skill": "prc"
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
  const path = '/dnd5e/skill-check';
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
        "skill": "prc"
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
  🔤/dnd5e/skill-check🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","skill":"prc"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/skill-check🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 52❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "skill-check-result",
  "requestId": "skill-check_1774367605684",
  "data": {
    "actorUuid": "Actor.ioZexonJDGVuU8zl",
    "skill": "prc",
    "total": 17,
    "formula": "1d20 + 2",
    "result": "15 + 2"
  }
}
```


---

## POST /dnd5e/ability-save

Roll an ability saving throw for an actor. Rolls a D&D 5e ability saving throw with all applicable modifiers.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| actorUuid | string | ✓ | body, query | UUID of the actor rolling the save |
| ability | string | ✓ | body, query | Ability abbreviation (e.g., "str", "dex", "con", "int", "wis", "cha") |
| clientId | string |  | body, query | Client ID for the Foundry world |
| advantage | boolean |  | body, query | Roll with advantage |
| disadvantage | boolean |  | body, query | Roll with disadvantage |
| bonus | string |  | body, query | Extra bonus formula to add (e.g., "1d4", "+2") |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the saving throw roll

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/ability-save"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"ability","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"disadvantage","type":"boolean","required":false,"source":"body"},{"name":"bonus","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/ability-save';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "ability": "dex"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/ability-save?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","ability":"dex"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/ability-save'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "ability": "dex"
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
  const path = '/dnd5e/ability-save';
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
        "ability": "dex"
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
  🔤/dnd5e/ability-save🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","ability":"dex"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/ability-save🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 54❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "ability-save-result",
  "requestId": "ability-save_1774367605731",
  "data": {
    "actorUuid": "Actor.ioZexonJDGVuU8zl",
    "ability": "dex",
    "total": 8,
    "formula": "1d20 + 3 + 2 + 0",
    "result": "3 + 3 + 2 + 0"
  }
}
```


---

## POST /dnd5e/ability-check

Roll an ability check for an actor. Rolls a D&D 5e ability check (raw ability test, not a skill check) with all applicable modifiers.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| actorUuid | string | ✓ | body, query | UUID of the actor rolling the check |
| ability | string | ✓ | body, query | Ability abbreviation (e.g., "str", "dex", "con", "int", "wis", "cha") |
| clientId | string |  | body, query | Client ID for the Foundry world |
| advantage | boolean |  | body, query | Roll with advantage |
| disadvantage | boolean |  | body, query | Roll with disadvantage |
| bonus | string |  | body, query | Extra bonus formula to add (e.g., "1d4", "+2") |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the ability check roll

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/ability-check"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"ability","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"disadvantage","type":"boolean","required":false,"source":"body"},{"name":"bonus","type":"string","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/ability-check';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "ability": "str"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/ability-check?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","ability":"str"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/ability-check'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "ability": "str"
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
  const path = '/dnd5e/ability-check';
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
        "ability": "str"
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
  🔤/dnd5e/ability-check🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","ability":"str"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/ability-check🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 54❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "ability-check-result",
  "requestId": "ability-check_1774367605758",
  "data": {
    "actorUuid": "Actor.ioZexonJDGVuU8zl",
    "ability": "str",
    "total": 3,
    "formula": "1d20 - 1",
    "result": "4 - 1"
  }
}
```


---

## POST /dnd5e/death-save

Roll a death saving throw for an actor. Rolls a D&D 5e death saving throw, handling DC 10 CON save, three successes/failures tracking, nat 20 healing, and nat 1 double failure.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| actorUuid | string | ✓ | body, query | UUID of the actor rolling the death save |
| clientId | string |  | body, query | Client ID for the Foundry world |
| advantage | boolean |  | body, query | Roll with advantage |
| createChatMessage | boolean |  | body, query | Whether to post the roll to chat (default: true) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the death saving throw

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/death-save"
  parameters={[{"name":"actorUuid","type":"string","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"advantage","type":"boolean","required":false,"source":"body"},{"name":"createChatMessage","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/death-save';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/death-save?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/death-save'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl"
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
  const path = '/dnd5e/death-save';
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl"
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
  🔤/dnd5e/death-save🔤 ➡️ path

  💭 Query parameters
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /dnd5e/death-save🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 38❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "death-save-result",
  "requestId": "death-save_1774367605925",
  "data": {
    "actorUuid": "Actor.ioZexonJDGVuU8zl",
    "total": 10,
    "formula": "1d20",
    "result": "10",
    "deathSaves": {
      "success": 1,
      "failure": 0
    }
  }
}
```


---

## POST /dnd5e/modify-experience

Modify the experience points for a specific actor. Adds or removes experience points from an actor.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| amount | number | ✓ | body, query | The amount of experience to add (can be negative) |
| clientId | string |  | body, query | Client ID for the Foundry world |
| actorUuid | string |  | body, query | UUID of the actor to modify |
| selected | boolean |  | body, query | Modify the selected token's actor |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the experience modification operation

### Try It Out

<ApiTester
  method="POST"
  path="/dnd5e/modify-experience"
  parameters={[{"name":"amount","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"body"},{"name":"actorUuid","type":"string","required":false,"source":"body"},{"name":"selected","type":"boolean","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/dnd5e/modify-experience';
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
      "amount": 100
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/dnd5e/modify-experience?clientId=foundry-testing-r6bXhB7k9cXa3cif' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"actorUuid":"Actor.ioZexonJDGVuU8zl","amount":100}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/dnd5e/modify-experience'
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
      "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
        "actorUuid": "Actor.ioZexonJDGVuU8zl",
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
  🔤clientId=foundry-testing-r6bXhB7k9cXa3cif🔤 ➡️ clientId
  🔤?🧲clientId🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"actorUuid":"Actor.ioZexonJDGVuU8zl","amount":100}🔤 ➡️ body

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
  "type": "modify-experience-result",
  "requestId": "modify-experience_1774367605094",
  "data": {
    "actorUuid": "Actor.ioZexonJDGVuU8zl",
    "oldXp": 0,
    "newXp": 100
  }
}
```

