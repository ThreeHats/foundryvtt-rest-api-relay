## **GET** /structure

## Returns the folders and compendiums in the world

### Request

#### Request URL

```
$baseUrl/structure?clientId=$clientId
```

#### Request Headers

| Key | Value | Description |
| --- | ----- | ----------- |
| x-api-key | \{\{apiKey\}\} |   |

#### Request Parameters

| Parameter Type | Key | Value | Description |
| -------------- | --- | ----- | ----------- |
| Query String Parameter | clientId | \{\{clientId\}\} | Auth token to connect to specific Foundry world |

### Response

#### Status: 200 OK

```json
{
  "requestId": "structure_1741128185565_fvzui9t",
  "clientId": "foundry-rQLkX9c1U2Tzkyh8",
  "folders": [
    {
      "id": "LB27IfNNZUR4bnom",
      "name": "Character Features",
      "type": "Compendium",
      "path": "Folder.LB27IfNNZUR4bnom",
      "sorting": 200000
    },
    {
      "id": "ni4v9Z0UmemhZMgy",
      "name": "D&D SRD Content",
      "type": "Compendium",
      "path": "Folder.ni4v9Z0UmemhZMgy",
      "sorting": 100000
    },
    {
      "id": "2hCpf26RfZ0boRGD",
      "name": "Items & Spells",
      "type": "Compendium",
      "path": "Folder.2hCpf26RfZ0boRGD",
      "sorting": 300000
    },
    {
      "id": "mq5f7JPS3ClCwnDc",
      "name": "Monsters",
      "type": "Compendium",
      "path": "Folder.mq5f7JPS3ClCwnDc",
      "sorting": 400000
    },
    {
      "id": "zmAZJmay9AxvRNqh",
      "name": "test",
      "type": "Actor",
      "path": "Folder.zmAZJmay9AxvRNqh",
      "sorting": 0
    }
  ],
  "compendiums": [
    {
      "id": "dnd5e.heroes",
      "name": "Starter Heroes",
      "path": "Compendium.dnd5e.heroes",
      "entity": "Actor",
      "packageType": "Actor",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.monsters",
      "name": "Monsters (SRD)",
      "path": "Compendium.dnd5e.monsters",
      "entity": "Actor",
      "packageType": "Actor",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.items",
      "name": "Items (SRD)",
      "path": "Compendium.dnd5e.items",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.tradegoods",
      "name": "Trade Goods (SRD)",
      "path": "Compendium.dnd5e.tradegoods",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.spells",
      "name": "Spells (SRD)",
      "path": "Compendium.dnd5e.spells",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.backgrounds",
      "name": "Backgrounds (SRD)",
      "path": "Compendium.dnd5e.backgrounds",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.classes",
      "name": "Classes (SRD)",
      "path": "Compendium.dnd5e.classes",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.subclasses",
      "name": "Subclasses (SRD)",
      "path": "Compendium.dnd5e.subclasses",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.classfeatures",
      "name": "Class & Subclass Features (SRD)",
      "path": "Compendium.dnd5e.classfeatures",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.races",
      "name": "Races (SRD)",
      "path": "Compendium.dnd5e.races",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.monsterfeatures",
      "name": "Monster Features (SRD)",
      "path": "Compendium.dnd5e.monsterfeatures",
      "entity": "Item",
      "packageType": "Item",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.rules",
      "name": "Rules (SRD)",
      "path": "Compendium.dnd5e.rules",
      "entity": "JournalEntry",
      "packageType": "JournalEntry",
      "system": "dnd5e"
    },
    {
      "id": "dnd5e.tables",
      "name": "Tables (SRD)",
      "path": "Compendium.dnd5e.tables",
      "entity": "RollTable",
      "packageType": "RollTable",
      "system": "dnd5e"
    }
  ]
}
```


