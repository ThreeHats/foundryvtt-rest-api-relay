---
tag: filesystem
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# FileSystem

## GET /file-system

Get file system structure

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| path | string |  | query | The path to retrieve (relative to source) |
| source | string |  | query | The source directory to use (data, systems, modules, etc.) |
| recursive | boolean |  | query | Whether to recursively list all subdirectories |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - File system structure with files and directories

### Try It Out

<ApiTester
  method="GET"
  path="/file-system"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"path","type":"string","required":false,"source":"query"},{"name":"source","type":"string","required":false,"source":"query"},{"name":"recursive","type":"boolean","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/file-system';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  source: 'data',
  recursive: 'false'
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
curl -X GET 'http://localhost:3010/file-system?clientId=fvtt_099ad17ea199e7e3&source=data&recursive=false' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/file-system'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'source': 'data',
    'recursive': 'false'
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
  const path = '/file-system';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    source: 'data',
    recursive: 'false'
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
  🔤/file-system🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤source=data🔤 ➡️ source
  🔤recursive=false🔤 ➡️ recursive
  🔤?🧲clientId🧲&🧲source🧲&🧲recursive🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /file-system🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "type": "file-system-result",
  "requestId": "file-system_1776657991001",
  "success": true,
  "path": "",
  "source": "data",
  "results": [
    {
      "name": "assets",
      "path": "assets",
      "type": "directory"
    },
    {
      "name": "maps",
      "path": "maps",
      "type": "directory"
    },
    {
      "name": "modules",
      "path": "modules",
      "type": "directory"
    },
    {
      "name": "moulinette",
      "path": "moulinette",
      "type": "directory"
    },
    {
      "name": "obsidian-files",
      "path": "obsidian-files",
      "type": "directory"
    },
    {
      "name": "rest-api-tests",
      "path": "rest-api-tests",
      "type": "directory"
    },
    {
      "name": "systems",
      "path": "systems",
      "type": "directory"
    },
    {
      "name": "tokenizer",
      "path": "tokenizer",
      "type": "directory"
    },
    {
      "name": "uploaded-chat-media",
      "path": "uploaded-chat-media",
      "type": "directory"
    },
    {
      "name": "worlds",
      "path": "worlds",
      "type": "directory"
    },
    {
      "name": "flooded-cave-test.webp",
      "path": "flooded-cave-test.webp",
      "type": "file"
    }
  ],
  "recursive": false
}
```


---

## GET /download

Download a file from Foundry's file system

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| path | string |  | query | The full path to the file to download |
| source | string |  | query | The source directory to use (data, systems, modules, etc.) |
| format | string |  | query | The format to return the file in (binary, base64) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - File contents in the requested format

### Try It Out

<ApiTester
  method="GET"
  path="/download"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"path","type":"string","required":false,"source":"query"},{"name":"source","type":"string","required":false,"source":"query"},{"name":"format","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/download';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  path: 'rest-api-tests/test-file.txt',
  source: 'data',
  format: 'base64'
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
curl -X GET 'http://localhost:3010/download?clientId=fvtt_099ad17ea199e7e3&path=rest-api-tests%2Ftest-file.txt&source=data&format=base64' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/download'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'path': 'rest-api-tests/test-file.txt',
    'source': 'data',
    'format': 'base64'
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
  const path = '/download';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    path: 'rest-api-tests/test-file.txt',
    source: 'data',
    format: 'base64'
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
  🔤/download🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤path=rest-api-tests/test-file.txt🔤 ➡️ path
  🔤source=data🔤 ➡️ source
  🔤format=base64🔤 ➡️ format
  🔤?🧲clientId🧲&🧲path🧲&🧲source🧲&🧲format🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤GET /download🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "fileData": "data:text/plain;base64,SGVsbG8gZnJvbSBSRVNUIEFQSSB0ZXN0IQ==",
  "filename": "test-file.txt",
  "mimeType": "text/plain",
  "path": "rest-api-tests/test-file.txt",
  "requestId": "download-file_1776657991006",
  "success": true,
  "type": "download-file-result"
}
```


---

## POST /upload

Upload a file to Foundry's file system (handles both base64 and binary data)

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| path | string |  | query, body | The directory path to upload to |
| filename | string |  | query, body | The filename to save as |
| source | string |  | query, body | The source directory to use (data, systems, modules, etc.) |
| mimeType | string |  | query, body | The MIME type of the file |
| overwrite | boolean |  | query, body | Whether to overwrite an existing file |
| fileData | string |  | body | Base64 encoded file data (if sending as JSON) 250MB limit |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Result of the file upload operation

### Try It Out

<ApiTester
  method="POST"
  path="/upload"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"path","type":"string","required":false,"source":"query"},{"name":"filename","type":"string","required":false,"source":"query"},{"name":"source","type":"string","required":false,"source":"query"},{"name":"mimeType","type":"string","required":false,"source":"query"},{"name":"overwrite","type":"boolean","required":false,"source":"query"},{"name":"fileData","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/upload';
const params = {
  clientId: 'fvtt_099ad17ea199e7e3',
  path: 'rest-api-tests',
  source: 'data',
  filename: 'test-file.txt',
  mimeType: 'text/plain'
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
      "fileData": "data:text/plain;base64,SGVsbG8gZnJvbSBSRVNUIEFQSSB0ZXN0IQ==",
      "mimeType": "text/plain",
      "overwrite": true
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/upload?clientId=fvtt_099ad17ea199e7e3&path=rest-api-tests&source=data&filename=test-file.txt&mimeType=text%2Fplain' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"fileData":"data:text/plain;base64,SGVsbG8gZnJvbSBSRVNUIEFQSSB0ZXN0IQ==","mimeType":"text/plain","overwrite":true}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/upload'
params = {
    'clientId': 'fvtt_099ad17ea199e7e3',
    'path': 'rest-api-tests',
    'source': 'data',
    'filename': 'test-file.txt',
    'mimeType': 'text/plain'
}
url = f'{base_url}{path}'

response = requests.post(
    url,
    params=params,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "fileData": "data:text/plain;base64,SGVsbG8gZnJvbSBSRVNUIEFQSSB0ZXN0IQ==",
      "mimeType": "text/plain",
      "overwrite": True
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
  const path = '/upload';
  const params = {
    clientId: 'fvtt_099ad17ea199e7e3',
    path: 'rest-api-tests',
    source: 'data',
    filename: 'test-file.txt',
    mimeType: 'text/plain'
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
        "fileData": "data:text/plain;base64,SGVsbG8gZnJvbSBSRVNUIEFQSSB0ZXN0IQ==",
        "mimeType": "text/plain",
        "overwrite": true
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
  🔤/upload🔤 ➡️ path

  💭 Query parameters
  🔤clientId=fvtt_099ad17ea199e7e3🔤 ➡️ clientId
  🔤path=rest-api-tests🔤 ➡️ path
  🔤source=data🔤 ➡️ source
  🔤filename=test-file.txt🔤 ➡️ filename
  🔤mimeType=text/plain🔤 ➡️ mimeType
  🔤?🧲clientId🧲&🧲path🧲&🧲source🧲&🧲filename🧲&🧲mimeType🧲🔤 ➡️ queryString

  💭 Request body
  🔤{"fileData":"data:text/plain;base64,SGVsbG8gZnJvbSBSRVNUIEFQSSB0ZXN0IQ==","mimeType":"text/plain","overwrite":true}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /upload🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 115❌r❌n❌r❌n🧲body🧲🔤 ➡️ request

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
  "type": "upload-file-result",
  "requestId": "upload-file_1776657990987",
  "success": true,
  "path": "rest-api-tests/test-file.txt"
}
```


