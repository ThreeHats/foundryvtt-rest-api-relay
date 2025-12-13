---
tag: session
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


# session

## POST /session-handshake

Create a handshake token for the client to use for secure authentication

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| x-api-key | string | ✓ | header | API key header |
| x-foundry-url | string | ✓ | header | Foundry URL header |
| x-username | string | ✓ | header | Username header |
| x-world-name | string |  | header | World name header |

### Returns

**object** - Handshake token and encryption details

---

## POST /start-session

Start a headless Foundry session using puppeteer

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| handshakeToken | string | ✓ | body | The token received from session-handshake |
| encryptedPassword | string | ✓ | body | Password encrypted with the public key |
| x-api-key | string | ✓ | header | API key header |

### Returns

**object** - Session information including sessionId and clientId

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
// Step 1: Get handshake credentials
const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key';

const handshakeResponse = await fetch(`${baseUrl}/session-handshake`, {
  method: 'POST',
  headers: {
    'x-api-key': apiKey,
    'x-foundry-url': 'http://localhost:30000',
    'x-world-name': 'my-world',
    'x-username': 'Gamemaster'
  }
});
const { token, publicKey, nonce } = await handshakeResponse.json();

// Step 2: Encrypt password using Web Crypto API (RSA-OAEP with SHA-1)
const password = 'your-password';
const payload = JSON.stringify({ password, nonce });

// Import the public key
const pemContents = publicKey
  .replace('-----BEGIN PUBLIC KEY-----', '')
  .replace('-----END PUBLIC KEY-----', '')
  .replace(/\n/g, '');
const binaryKey = Uint8Array.from(atob(pemContents), c => c.charCodeAt(0));

const cryptoKey = await crypto.subtle.importKey(
  'spki',
  binaryKey,
  { name: 'RSA-OAEP', hash: 'SHA-1' },  // Must use SHA-1 to match server
  false,
  ['encrypt']
);

// Encrypt the payload
const encrypted = await crypto.subtle.encrypt(
  { name: 'RSA-OAEP' },
  cryptoKey,
  new TextEncoder().encode(payload)
);
const encryptedPassword = btoa(String.fromCharCode(...new Uint8Array(encrypted)));

// Step 3: Start the session
const sessionResponse = await fetch(`${baseUrl}/start-session`, {
  method: 'POST',
  headers: {
    'x-api-key': apiKey,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ handshakeToken: token, encryptedPassword })
});
const { sessionId, clientId } = await sessionResponse.json();
console.log('Session started:', { sessionId, clientId });
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
# Session creation requires encryption - use the JavaScript, Python, or TypeScript examples.
# 
# The flow is:
# 1. POST /session-handshake to get token, publicKey, nonce
# 2. Encrypt JSON payload { password, nonce } with RSA-OAEP using publicKey
# 3. POST /start-session with { handshakeToken, encryptedPassword }
#
# Example handshake request:
curl -X POST 'http://localhost:3010/session-handshake' \
  -H "x-api-key: your-api-key" \
  -H "x-foundry-url: http://localhost:30000" \
  -H "x-world-name: my-world" \
  -H "x-username: Gamemaster"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import padding
import base64
import json

base_url = 'http://localhost:3010'
api_key = 'your-api-key'

# Step 1: Get handshake credentials
handshake_response = requests.post(
    f'{base_url}/session-handshake',
    headers={
        'x-api-key': api_key,
        'x-foundry-url': 'http://localhost:30000',
        'x-world-name': 'my-world',
        'x-username': 'Gamemaster'
    }
)
handshake_data = handshake_response.json()
token = handshake_data['token']
public_key_pem = handshake_data['publicKey']
nonce = handshake_data['nonce']

# Step 2: Encrypt password using RSA-OAEP with SHA-1 (must match server)
password = 'your-password'
payload = json.dumps({'password': password, 'nonce': nonce})

public_key = serialization.load_pem_public_key(public_key_pem.encode())
encrypted = public_key.encrypt(
    payload.encode(),
    padding.OAEP(
        mgf=padding.MGF1(algorithm=hashes.SHA1()),
        algorithm=hashes.SHA1(),
        label=None
    )
)
encrypted_password = base64.b64encode(encrypted).decode()

# Step 3: Start the session
session_response = requests.post(
    f'{base_url}/start-session',
    headers={'x-api-key': api_key},
    json={'handshakeToken': token, 'encryptedPassword': encrypted_password}
)
session_data = session_response.json()
print('Session started:', session_data)
```

</TabItem>
<TabItem value="typescript" label="TypeScript">

```typescript
import crypto from 'crypto';

// Step 1: Get handshake credentials
const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key';

const handshakeResponse = await fetch(`${baseUrl}/session-handshake`, {
  method: 'POST',
  headers: {
    'x-api-key': apiKey,
    'x-foundry-url': 'http://localhost:30000',
    'x-world-name': 'my-world',
    'x-username': 'Gamemaster'
  }
});
const { token, publicKey, nonce } = await handshakeResponse.json();

// Step 2: Encrypt password using Node.js crypto (RSA-OAEP)
const password = 'your-password';
const payload = JSON.stringify({ password, nonce });
const encryptedPassword = crypto.publicEncrypt(
  { key: publicKey, padding: crypto.constants.RSA_PKCS1_OAEP_PADDING },
  Buffer.from(payload, 'utf8')
).toString('base64');

// Step 3: Start the session
const sessionResponse = await fetch(`${baseUrl}/start-session`, {
  method: 'POST',
  headers: {
    'x-api-key': apiKey,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ handshakeToken: token, encryptedPassword })
});
const { sessionId, clientId } = await sessionResponse.json();
console.log('Session started:', { sessionId, clientId });
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
  🔤/start-session🔤 ➡️ path

  💭 Request body
  🔤{"handshakeToken":"your-api-key-hereyour-api-key-here","encryptedPassword":"vS9rgCQ0plAzjXCvHDObn11XHyv+ScQdOiUjfUgiS9l91QgjDf3snZN1BJX8QZxvy6OJlUr8CLEg2mgzxfM2T3AJ4uWG0hzLzXB308Tx8pVYj1F0TCpIlHhrLZHlOQaHTKxmITMy6L+Mom0ZRXe/bAna+hkC/QeUi2U8LMupTGXa5ecA5PQBXCW3i3KBIJqee0GV2usSU6y2LZv6KpV+IKqEACUjiRrKXNqx/LV/jOOrFLUUKfatZP8IPlkXB6h71SsJvNG5/BzLzj0/MGLy+hrIYmI8bpMSLz7uJxXX2JAZ+gUenSmijBwRxu6kAvBMxL4vIteKFxfo/PpBRpDMEA=="}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /start-session HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 452❌r❌n❌r❌n{"handshakeToken":"your-api-key-hereyour-api-key-here","encryptedPassword":"vS9rgCQ0plAzjXCvHDObn11XHyv+ScQdOiUjfUgiS9l91QgjDf3snZN1BJX8QZxvy6OJlUr8CLEg2mgzxfM2T3AJ4uWG0hzLzXB308Tx8pVYj1F0TCpIlHhrLZHlOQaHTKxmITMy6L+Mom0ZRXe/bAna+hkC/QeUi2U8LMupTGXa5ecA5PQBXCW3i3KBIJqee0GV2usSU6y2LZv6KpV+IKqEACUjiRrKXNqx/LV/jOOrFLUUKfatZP8IPlkXB6h71SsJvNG5/BzLzj0/MGLy+hrIYmI8bpMSLz7uJxXX2JAZ+gUenSmijBwRxu6kAvBMxL4vIteKFxfo/PpBRpDMEA=="}🔤 ➡️ request

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
  "success": true,
  "message": "Foundry session started successfully",
  "sessionId": "46f06659-72d7-450f-82f5-537a4cd4bb8b",
  "clientId": "your-client-id"
}
```


---

## DELETE /end-session

Stop a headless Foundry session

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| sessionId | string | ✓ | query | The ID of the session to end |
| x-api-key | string | ✓ | header | API key header |

### Returns

**object** - Status of the operation

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/end-session';
const params = {
  sessionId: '46f06659-72d7-450f-82f5-537a4cd4bb8b'
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
curl -X DELETE 'http://localhost:3010/end-session?sessionId=46f06659-72d7-450f-82f5-537a4cd4bb8b' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/end-session'
params = {
    'sessionId': '46f06659-72d7-450f-82f5-537a4cd4bb8b'
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
  const path = '/end-session';
  const params = {
    sessionId: '46f06659-72d7-450f-82f5-537a4cd4bb8b'
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
  🔤/end-session🔤 ➡️ path

  💭 Query parameters
  🔤sessionId=46f06659-72d7-450f-82f5-537a4cd4bb8b🔤 ➡️ sessionId
  🔤?🧲sessionId🧲🔤 ➡️ queryString

  💭 Build HTTP request
  🔤DELETE /end-session🧲queryString🧲 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "success": true,
  "message": "Foundry session terminated (partial cleanup)"
}
```


---

## GET /session

Get all active headless Foundry sessions

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|--------------|
| x-api-key | string | ✓ | header | API key header |

### Returns

**object** - List of active sessions for the current API key

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/session';
const url = `${baseUrl}${path}`;

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
curl -X GET 'http://localhost:3010/session' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/session'
url = f'{base_url}{path}'

response = requests.get(
    url,
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
  const path = '/session';
  const url = `${baseUrl}${path}`;

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
  🔤/session🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /session HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "activeSessions": [
    {
      "id": "46f06659-72d7-450f-82f5-537a4cd4bb8b",
      "clientId": "your-client-id",
      "lastActivity": 1765635936184,
      "idleMinutes": 0,
      "instanceId": "local",
      "worldId": "testing",
      "worldTitle": "testing",
      "foundryVersion": "13.348",
      "systemId": "dnd5e",
      "systemTitle": "Dungeons & Dragons Fifth Edition",
      "systemVersion": "5.0.4",
      "customName": ""
    }
  ],
  "pendingSessions": []
}
```

