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
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Search results containing matching entities

### Try It Out

<ApiTester
  method="GET"
  path="/search"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"query","type":"string","required":false,"source":"query"},{"name":"filter","type":"string","required":false,"source":"query"},{"name":"excludeCompendiums","type":"boolean","required":false,"source":"query"},{"name":"limit","type":"number","required":false,"source":"query"},{"name":"minified","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/search';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
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
curl -X GET 'http://localhost:3010/search?clientId=fvtt_099ad17ea199e7e3&filter=documentType%3AActor&excludeCompendiums=true' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/search'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
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
    clientId: 'fvtt_099ad17ea199e7e3',
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
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
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
  "requestId": "search_1776657969438",
  "filter": "documentType:Actor",
  "results": [
    {
      "documentType": "Actor",
      "folder": null,
      "id": "XJ53lfbhXlqQZCbd",
      "name": "Player Character",
      "package": null,
      "packageName": null,
      "subType": "character",
      "uuid": "Actor.XJ53lfbhXlqQZCbd",
      "icon": "icons/svg/mystery-man.svg",
      "journalLink": "@UUID[Actor.XJ53lfbhXlqQZCbd]{Player Character}",
      "tagline": "Actors Directory",
      "formattedMatch": "Player Character",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "w5STPCwE3YTDztRk",
      "name": "test-perrin (halfling monk)",
      "package": null,
      "packageName": null,
      "subType": "character",
      "uuid": "Actor.w5STPCwE3YTDztRk",
      "icon": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
      "journalLink": "@UUID[Actor.w5STPCwE3YTDztRk]{test-perrin (halfling monk)}",
      "tagline": "Actors Directory",
      "formattedMatch": "test-perrin (halfling monk)",
      "resultType": "WorldEntity"
    },
    {
      "documentType": "Actor",
      "folder": null,
      "id": "q9uWyfdPwTlzbpxb",
      "name": "Updated Test Actor",
      "package": null,
      "packageName": null,
      "subType": "character",
      "uuid": "Actor.q9uWyfdPwTlzbpxb",
      "icon": "systems/dnd5e/tokens/heroes/MonkStaff.webp",
      "journalLink": "@UUID[Actor.q9uWyfdPwTlzbpxb]{Updated Test Actor}",
      "tagline": "Actors Directory",
      "formattedMatch": "Updated Test Actor",
      "resultType": "WorldEntity"
    }
  ]
}
```


