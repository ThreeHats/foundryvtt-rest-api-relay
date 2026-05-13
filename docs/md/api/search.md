---
tag: search
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Search

## GET /search

Search entities

This endpoint allows searching for entities in the Foundry world based on a query string. Search world entities and compendiums using the native built-in search engine. No third-party modules required. Results are ranked by relevance: exact match, prefix match, substring match, word-prefix match, and subsequence match.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| query | string |  | query | Search query string (omit to browse all entities matching filter) |
| filter | string |  | query | Filter string — simple: filter="Actor"; compound: filter="documentType:Item,subType:weapon". Supported keys: documentType, subType, folder, package, resultType |
| excludeCompendiums | boolean |  | query | Exclude compendium entries from results (default: false — compendiums are included by default) |
| limit | number |  | query | Maximum number of results to return (default: 200, max: 500) |
| minified | boolean |  | query | Return minimal fields only — uuid, id, name, img, documentType (default: false) |
| ownedByUserId | string |  | query | Filter results to only documents the specified Foundry user (ID or username) has Owner permission on |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Search results containing matching entities

### Try It Out

<ApiTester
  method="GET"
  path="/search"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"query","type":"string","required":false,"source":"query"},{"name":"filter","type":"string","required":false,"source":"query"},{"name":"excludeCompendiums","type":"boolean","required":false,"source":"query"},{"name":"limit","type":"number","required":false,"source":"query"},{"name":"minified","type":"boolean","required":false,"source":"query"},{"name":"ownedByUserId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/search';
const params = {
  clientId: 'fvtt_71dbc81bd608978a',
  filter: 'documentType:Actor',
  excludeCompendiums: 'true'
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
curl -X GET 'http://localhost:3010/search?clientId=fvtt_71dbc81bd608978a&filter=documentType%3AActor&excludeCompendiums=true' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/search'
params = {
    'clientId': 'fvtt_71dbc81bd608978a',
    'filter': 'documentType:Actor',
    'excludeCompendiums': 'true'
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
  const path = '/search';
  const params = {
    clientId: 'fvtt_71dbc81bd608978a',
    filter: 'documentType:Actor',
    excludeCompendiums: 'true'
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
  🔤/search🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_71dbc81bd608978a🔤 ➡️ clientId
  🔤filter=documentType:Actor🔤 ➡️ filter
  🔤excludeCompendiums=true🔤 ➡️ excludeCompendiums
  🔤?🧲clientId🧲&🧲filter🧲&🧲excludeCompendiums🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /search🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "search-result",
  "requestId": "search_1777996592905",
  "filter": "documentType:Actor",
  "results": [
    {
      "documentType": "Actor",
      "folder": null,
      "id": "2midVQ1laJFMrN4D",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.2midVQ1laJFMrN4D",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.2midVQ1laJFMrN4D]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "IGDNKF9amHdCQYGX",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.IGDNKF9amHdCQYGX",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.IGDNKF9amHdCQYGX]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "Q8f1DRjKZJmtZNvD",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.Q8f1DRjKZJmtZNvD",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.Q8f1DRjKZJmtZNvD]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "V162ZiBjzf489Ekz",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.V162ZiBjzf489Ekz",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.V162ZiBjzf489Ekz]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "Xp4mNa6B0fBeKAAF",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.Xp4mNa6B0fBeKAAF",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.Xp4mNa6B0fBeKAAF]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "ZTJRxZIdnhFQalD0",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.ZTJRxZIdnhFQalD0",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.ZTJRxZIdnhFQalD0]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "dKQ5TvZ1ySmZ5zC1",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.dKQ5TvZ1ySmZ5zC1",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.dKQ5TvZ1ySmZ5zC1]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "haaLIgOLEYGT0nqK",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.haaLIgOLEYGT0nqK",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.haaLIgOLEYGT0nqK]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "iOn5ZWzR0Hv1gN06",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.iOn5ZWzR0Hv1gN06",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.iOn5ZWzR0Hv1gN06]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "liCyZwJ5BRnryiXU",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.liCyZwJ5BRnryiXU",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.liCyZwJ5BRnryiXU]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "mNMNalBcWKhoQ8QE",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.mNMNalBcWKhoQ8QE",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.mNMNalBcWKhoQ8QE]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "wZR5R8cIs1dUpKC7",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.wZR5R8cIs1dUpKC7",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.wZR5R8cIs1dUpKC7]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "ywz3tLzEpjMCdXhK",
      "name": "hook-test-actor",
      "package": null,
      "packageName": null,
      "subType": "base",
      "uuid": "Actor.ywz3tLzEpjMCdXhK",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.ywz3tLzEpjMCdXhK]{hook-test-actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "hook-test-actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "4f57kTea0R0ZTOta",
      "name": "test-perrin (halfling monk)",
      "package": null,
      "packageName": null,
      "subType": "character",
      "uuid": "Actor.4f57kTea0R0ZTOta",
      "icon": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
      "journalLink": "@UUID[Actor.4f57kTea0R0ZTOta]{test-perrin (halfling monk)}",
      "tagline": "Actors Directory",
      "formattedMatch": "test-perrin (halfling monk)",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "Nr1kY13FmjUZCHjw",
      "name": "test-perrin (halfling monk)",
      "package": null,
      "packageName": null,
      "subType": "character",
      "uuid": "Actor.Nr1kY13FmjUZCHjw",
      "icon": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
      "journalLink": "@UUID[Actor.Nr1kY13FmjUZCHjw]{test-perrin (halfling monk)}",
      "tagline": "Actors Directory",
      "formattedMatch": "test-perrin (halfling monk)",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "yqVY0l6AY0HDtq7Z",
      "name": "Updated Test Actor",
      "package": null,
      "packageName": null,
      "subType": "character",
      "uuid": "Actor.yqVY0l6AY0HDtq7Z",
      "icon": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
      "journalLink": "@UUID[Actor.yqVY0l6AY0HDtq7Z]{Updated Test Actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "Updated Test Actor",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "sMD3o6zej6ckQkpo",
      "name": "Updated Test Actor",
      "package": null,
      "packageName": null,
      "subType": "character",
      "uuid": "Actor.sMD3o6zej6ckQkpo",
      "icon": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
      "journalLink": "@UUID[Actor.sMD3o6zej6ckQkpo]{Updated Test Actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "Updated Test Actor",
      "resultType": "WorldEntity"
    }
  ]
}
```


