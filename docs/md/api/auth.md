---
tag: auth
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Auth

## POST /auth/register

Register a new user account

Creates a new user with the provided email and password. Returns the new user's API key.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| email | string | ✓ | body | The user's email address |
| password | string | ✓ | body | The user's password (min 8 chars, must include uppercase, lowercase, and number) |

### Returns

**object** - User object with API key

### Try It Out

<ApiTester
  method="POST"
  path="/auth/register"
  parameters={[{"name":"email","type":"string","required":true,"source":"body"},{"name":"password","type":"string","required":true,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/register';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "email": "auth-test-1775068868827@example.com",
      "password": "TestPassword1"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/register' \
  -H "Content-Type: application/json" \
  -d '{"email":"auth-test-1775068868827@example.com","password":"TestPassword1"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/register'
url = f'{base_url}{path}'

response = requests.post(
    url,
    json={
      "email": "auth-test-1775068868827@example.com",
      "password": "TestPassword1"
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
  const path = '/auth/register';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "email": "auth-test-1775068868827@example.com",
        "password": "TestPassword1"
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
  🔤/auth/register🔤 ➡️ path

  💭 Request body
  🔤{"email":"auth-test-1775068868827@example.com","password":"TestPassword1"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/register HTTP/1.1❌r❌nHost: localhost:3010❌r❌nContent-Type: application/json❌r❌nContent-Length: 74❌r❌n❌r❌n{"email":"auth-test-1775068868827@example.com","password":"TestPassword1"}🔤 ➡️ request

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

**Status:** 201

```json
{
  "apiKey": "your-api-key-here",
  "createdAt": {
    "Time": "2026-04-01T13:41:08.884915504-05:00",
    "Valid": true
  },
  "email": "auth-test-1775068868827@example.com",
  "id": 12,
  "subscriptionStatus": "free"
}
```


---

## POST /auth/login

Log in with email and password

Authenticates a user and returns their account details including API key.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| email | string | ✓ | body | The user's email address |
| password | string | ✓ | body | The user's password |

### Returns

**object** - User object with API key

### Try It Out

<ApiTester
  method="POST"
  path="/auth/login"
  parameters={[{"name":"email","type":"string","required":true,"source":"body"},{"name":"password","type":"string","required":true,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/login';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "email": "noahblalonde@gmail.com",
      "password": "nek_jtq_awk7JZD6vwj"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/login' \
  -H "Content-Type: application/json" \
  -d '{"email":"noahblalonde@gmail.com","password":"nek_jtq_awk7JZD6vwj"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/login'
url = f'{base_url}{path}'

response = requests.post(
    url,
    json={
      "email": "noahblalonde@gmail.com",
      "password": "nek_jtq_awk7JZD6vwj"
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
  const path = '/auth/login';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "email": "noahblalonde@gmail.com",
        "password": "nek_jtq_awk7JZD6vwj"
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
  🔤/auth/login🔤 ➡️ path

  💭 Request body
  🔤{"email":"noahblalonde@gmail.com","password":"nek_jtq_awk7JZD6vwj"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/login HTTP/1.1❌r❌nHost: localhost:3010❌r❌nContent-Type: application/json❌r❌nContent-Length: 67❌r❌n❌r❌n{"email":"noahblalonde@gmail.com","password":"nek_jtq_awk7JZD6vwj"}🔤 ➡️ request

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
  "apiKey": "your-api-key-here",
  "createdAt": {
    "Time": "0001-01-01T00:00:00Z",
    "Valid": false
  },
  "email": "noahblalonde@gmail.com",
  "id": 4,
  "requestsThisMonth": 8059
}
```


---

## POST /auth/regenerate-key

Regenerate API key

Generates a new API key for the authenticated user. The old key is invalidated.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| email | string | ✓ | body | The user's email address |
| password | string | ✓ | body | The user's password |

### Returns

**object** - New API key

### Try It Out

<ApiTester
  method="POST"
  path="/auth/regenerate-key"
  parameters={[{"name":"email","type":"string","required":true,"source":"body"},{"name":"password","type":"string","required":true,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/regenerate-key';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "email": "auth-test-1775068868827@example.com",
      "password": "TestPassword1"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/regenerate-key' \
  -H "Content-Type: application/json" \
  -d '{"email":"auth-test-1775068868827@example.com","password":"TestPassword1"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/regenerate-key'
url = f'{base_url}{path}'

response = requests.post(
    url,
    json={
      "email": "auth-test-1775068868827@example.com",
      "password": "TestPassword1"
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
  const path = '/auth/regenerate-key';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "email": "auth-test-1775068868827@example.com",
        "password": "TestPassword1"
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
  🔤/auth/regenerate-key🔤 ➡️ path

  💭 Request body
  🔤{"email":"auth-test-1775068868827@example.com","password":"TestPassword1"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/regenerate-key HTTP/1.1❌r❌nHost: localhost:3010❌r❌nContent-Type: application/json❌r❌nContent-Length: 74❌r❌n❌r❌n{"email":"auth-test-1775068868827@example.com","password":"TestPassword1"}🔤 ➡️ request

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
  "apiKey": "your-api-key-here"
}
```


---

## GET /auth/user-data

Get user data

Retrieves the authenticated user's account details including usage statistics and subscription status.

### Returns

**object** - User data object

### Try It Out

<ApiTester
  method="GET"
  path="/auth/user-data"
  parameters={[]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/user-data';
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
curl -X GET 'http://localhost:3010/auth/user-data' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/user-data'
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
  const path = '/auth/user-data';
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
  🔤/auth/user-data🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /auth/user-data HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "apiKey": "your-api-key-here",
  "email": "noahblalonde@gmail.com",
  "id": 4,
  "limits": {
    "dailyLimit": 10000000,
    "monthlyLimit": 1000000,
    "unlimitedMonthly": false
  },
  "requestsThisMonth": 8059,
  "requestsToday": 272,
  "subscriptionStatus": "free"
}
```


---

## GET /auth/export-data

Export user data

Exports all stored user data for GDPR data portability compliance.

### Returns

**object** - Complete user data export

### Try It Out

<ApiTester
  method="GET"
  path="/auth/export-data"
  parameters={[]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/export-data';
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
curl -X GET 'http://localhost:3010/auth/export-data' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/export-data'
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
  const path = '/auth/export-data';
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
  🔤/auth/export-data🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /auth/export-data HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "apiAccess": {
    "apiKey": "your-api-key-here"
  },
  "exportDate": "2026-04-01T18:41:09Z",
  "scopedKeys": [
    {
      "createdAt": {
        "Time": "0001-01-01T00:00:00Z",
        "Valid": false
      },
      "dailyLimit": 1000,
      "enabled": true,
      "expiresAt": null,
      "hasFoundryCredentials": false,
      "id": 1,
      "name": "test1",
      "scopedClientId": "your-client-id",
      "scopedUserId": "tester"
    },
    {
      "createdAt": {
        "Time": "0001-01-01T00:00:00Z",
        "Valid": false
      },
      "dailyLimit": 0,
      "enabled": true,
      "expiresAt": null,
      "hasFoundryCredentials": true,
      "id": 2,
      "name": "headless-test1",
      "scopedClientId": "your-client-id",
      "scopedUserId": ""
    },
    {
      "createdAt": {
        "Time": "0001-01-01T00:00:00Z",
        "Valid": false
      },
      "dailyLimit": 0,
      "enabled": true,
      "expiresAt": null,
      "hasFoundryCredentials": true,
      "id": 3,
      "name": "test-headless2-gm",
      "scopedClientId": "foundry-testing-r6bXhB7k9cXa3cif",
      "scopedUserId": "Gamemaster"
    },
    {
      "createdAt": {
        "Time": "2026-03-25T18:08:33.030028526-05:00",
        "Valid": true
      },
      "dailyLimit": 0,
      "enabled": true,
      "expiresAt": null,
      "hasFoundryCredentials": true,
      "id": 16,
      "name": "test-session",
      "scopedClientId": "",
      "scopedUserId": ""
    }
  ],
  "subscription": {
    "status": "free",
    "stripeCustomerId": "",
    "subscriptionId": ""
  },
  "usage": {
    "requestsThisMonth": 8059,
    "requestsToday": 272
  },
  "user": {
    "createdAt": {
      "Time": "0001-01-01T00:00:00Z",
      "Valid": false
    },
    "email": "noahblalonde@gmail.com",
    "id": 4,
    "updatedAt": {
      "Time": "2026-04-01T13:41:08.336801707-05:00",
      "Valid": true
    }
  }
}
```


---

## DELETE /auth/account

Delete user account

Permanently deletes the user's account and all associated data.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| confirmEmail | string | ✓ | body | The user's email address (must match account email) |
| password | string | ✓ | body | The user's password for verification |

### Returns

**object** - Confirmation of deletion

### Try It Out

<ApiTester
  method="DELETE"
  path="/auth/account"
  parameters={[{"name":"confirmEmail","type":"string","required":true,"source":"body"},{"name":"password","type":"string","required":true,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/account';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'DELETE',
  headers: {
    'x-api-key': 'your-api-key-here',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "confirmEmail": "auth-test-1775068868827@example.com",
      "password": "ResetPassword1"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X DELETE 'http://localhost:3010/auth/account' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"confirmEmail":"auth-test-1775068868827@example.com","password":"ResetPassword1"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/account'
url = f'{base_url}{path}'

response = requests.delete(
    url,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "confirmEmail": "auth-test-1775068868827@example.com",
      "password": "ResetPassword1"
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
  const path = '/auth/account';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'delete',
    headers: {
      'x-api-key': 'your-api-key-here',
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "confirmEmail": "auth-test-1775068868827@example.com",
        "password": "ResetPassword1"
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
  🔤/auth/account🔤 ➡️ path

  💭 Request body
  🔤{"confirmEmail":"auth-test-1775068868827@example.com","password":"ResetPassword1"}🔤 ➡️ body

  💭 Build HTTP request
  🔤DELETE /auth/account HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 82❌r❌n❌r❌n{"confirmEmail":"auth-test-1775068868827@example.com","password":"ResetPassword1"}🔤 ➡️ request

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
  "message": "Account successfully deleted",
  "success": true
}
```


---

## POST /auth/change-password

Change password while logged in

Allows an authenticated user to change their password.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| currentPassword | string | ✓ | body | The user's current password |
| newPassword | string | ✓ | body | The new password (min 8 chars, must include uppercase, lowercase, and number) |

### Returns

**object** - Success message

### Try It Out

<ApiTester
  method="POST"
  path="/auth/change-password"
  parameters={[{"name":"currentPassword","type":"string","required":true,"source":"body"},{"name":"newPassword","type":"string","required":true,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/change-password';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'x-api-key': 'your-api-key-here',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "currentPassword": "nek_jtq_awk7JZD6vwj",
      "newPassword": "ChangedPassword2"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/change-password' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"currentPassword":"nek_jtq_awk7JZD6vwj","newPassword":"ChangedPassword2"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/change-password'
url = f'{base_url}{path}'

response = requests.post(
    url,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "currentPassword": "nek_jtq_awk7JZD6vwj",
      "newPassword": "ChangedPassword2"
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
  const path = '/auth/change-password';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'x-api-key': 'your-api-key-here',
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "currentPassword": "nek_jtq_awk7JZD6vwj",
        "newPassword": "ChangedPassword2"
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
  🔤/auth/change-password🔤 ➡️ path

  💭 Request body
  🔤{"currentPassword":"nek_jtq_awk7JZD6vwj","newPassword":"ChangedPassword2"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/change-password HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 74❌r❌n❌r❌n{"currentPassword":"nek_jtq_awk7JZD6vwj","newPassword":"ChangedPassword2"}🔤 ➡️ request

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
  "message": "Password changed successfully"
}
```


---

## POST /auth/forgot-password

Request a password reset

Sends a password reset email if the account exists.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| email | string | ✓ | body | The email address associated with the account |

### Returns

**object** - Generic success message

### Try It Out

<ApiTester
  method="POST"
  path="/auth/forgot-password"
  parameters={[{"name":"email","type":"string","required":true,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/forgot-password';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "email": "auth-test-1775068868827@example.com"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/forgot-password' \
  -H "Content-Type: application/json" \
  -d '{"email":"auth-test-1775068868827@example.com"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/forgot-password'
url = f'{base_url}{path}'

response = requests.post(
    url,
    json={
      "email": "auth-test-1775068868827@example.com"
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
  const path = '/auth/forgot-password';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "email": "auth-test-1775068868827@example.com"
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
  🔤/auth/forgot-password🔤 ➡️ path

  💭 Request body
  🔤{"email":"auth-test-1775068868827@example.com"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/forgot-password HTTP/1.1❌r❌nHost: localhost:3010❌r❌nContent-Type: application/json❌r❌nContent-Length: 47❌r❌n❌r❌n{"email":"auth-test-1775068868827@example.com"}🔤 ➡️ request

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
  "message": "If an account with that email exists, a password reset link has been sent."
}
```


---

## POST /auth/reset-password

Reset password with token

Resets the user's password using a valid reset token.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| token | string | ✓ | body | The password reset token from the email |
| password | string | ✓ | body | The new password (min 8 chars, must include uppercase, lowercase, and number) |

### Returns

**object** - Success message

### Try It Out

<ApiTester
  method="POST"
  path="/auth/reset-password"
  parameters={[{"name":"token","type":"string","required":true,"source":"body"},{"name":"password","type":"string","required":true,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/reset-password';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "token": "your-api-key-here",
      "password": "ResetPassword1"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/reset-password' \
  -H "Content-Type: application/json" \
  -d '{"token":"your-api-key-here","password":"ResetPassword1"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/reset-password'
url = f'{base_url}{path}'

response = requests.post(
    url,
    json={
      "token": "your-api-key-here",
      "password": "ResetPassword1"
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
  const path = '/auth/reset-password';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "token": "your-api-key-here",
        "password": "ResetPassword1"
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
  🔤/auth/reset-password🔤 ➡️ path

  💭 Request body
  🔤{"token":"your-api-key-here","password":"ResetPassword1"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/reset-password HTTP/1.1❌r❌nHost: localhost:3010❌r❌nContent-Type: application/json❌r❌nContent-Length: 104❌r❌n❌r❌n{"token":"your-api-key-here","password":"ResetPassword1"}🔤 ➡️ request

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
  "message": "Password has been reset successfully"
}
```


---

## GET /auth/validate-reset-token/:token

Validate a password reset token

Checks whether a password reset token is still valid before showing the reset form.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| token | string | ✓ | params | The password reset token to validate |

### Returns

**object** - Token validity status

### Try It Out

<ApiTester
  method="GET"
  path="/auth/validate-reset-token/:token"
  parameters={[{"name":"token","type":"string","required":true,"source":"params"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/validate-reset-token/your-api-key-here';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'GET'
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X GET 'http://localhost:3010/auth/validate-reset-token/your-api-key-here'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/validate-reset-token/your-api-key-here'
url = f'{base_url}{path}'

response = requests.get(
    url
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
  const path = '/auth/validate-reset-token/your-api-key-here';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'get',
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
  🔤/auth/validate-reset-token/your-api-key-here🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /auth/validate-reset-token/your-api-key-here HTTP/1.1❌r❌nHost: localhost:3010❌r❌n❌r❌n🔤 ➡️ request

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
  "valid": true
}
```


---

## POST /auth/api-keys

Create a new scoped API key

Creates a sub-key under the authenticated user's master key with optional restrictions.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| name | string | ✓ | body | Friendly name for the key |
| scopedClientId | string |  | body | Lock to specific Foundry client ID |
| scopedUserId | string |  | body | Lock to specific Foundry user ID |
| dailyLimit | string |  | body | Per-key daily request cap |
| expiresAt | string |  | body | Expiry timestamp (ISO 8601) |
| foundryUrl | string |  | body | Foundry instance URL for headless sessions |
| foundryUsername | string |  | body | Foundry login username |
| foundryPassword | string |  | body | Foundry login password (encrypted at rest) |

### Returns

**object** - New scoped API key details

### Try It Out

<ApiTester
  method="POST"
  path="/auth/api-keys"
  parameters={[{"name":"name","type":"string","required":true,"source":"body"},{"name":"scopedClientId","type":"string","required":false,"source":"body"},{"name":"scopedUserId","type":"string","required":false,"source":"body"},{"name":"dailyLimit","type":"string","required":false,"source":"body"},{"name":"expiresAt","type":"string","required":false,"source":"body"},{"name":"foundryUrl","type":"string","required":false,"source":"body"},{"name":"foundryUsername","type":"string","required":false,"source":"body"},{"name":"foundryPassword","type":"string","required":false,"source":"body"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/api-keys';
const url = `${baseUrl}${path}`;

const response = await fetch(url, {
  method: 'POST',
  headers: {
    'x-api-key': 'your-api-key-here',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
      "name": "Test Scoped Key",
      "dailyLimit": "500"
    })
});
const data = await response.json();
console.log(data);
```

</TabItem>
<TabItem value="curl" label="cURL">

```bash
curl -X POST 'http://localhost:3010/auth/api-keys' \
  -H "x-api-key: your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Scoped Key","dailyLimit":"500"}'
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/api-keys'
url = f'{base_url}{path}'

response = requests.post(
    url,
    headers={
        'x-api-key': 'your-api-key-here'
    },
    json={
      "name": "Test Scoped Key",
      "dailyLimit": "500"
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
  const path = '/auth/api-keys';
  const url = `${baseUrl}${path}`;

  const response = await axios({
    method: 'post',
    headers: {
      'x-api-key': 'your-api-key-here',
      'Content-Type': 'application/json'
    },
    url,
    data: {
        "name": "Test Scoped Key",
        "dailyLimit": "500"
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
  🔤/auth/api-keys🔤 ➡️ path

  💭 Request body
  🔤{"name":"Test Scoped Key","dailyLimit":"500"}🔤 ➡️ body

  💭 Build HTTP request
  🔤POST /auth/api-keys HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌nContent-Type: application/json❌r❌nContent-Length: 45❌r❌n❌r❌n{"name":"Test Scoped Key","dailyLimit":"500"}🔤 ➡️ request

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

**Status:** 201

```json
{
  "createdAt": {
    "Time": "2026-04-01T13:41:21.466281879-05:00",
    "Valid": true
  },
  "dailyLimit": 500,
  "enabled": true,
  "expiresAt": null,
  "hasFoundryCredentials": false,
  "id": 18,
  "key": "your-api-key-here",
  "name": "Test Scoped Key",
  "scopedClientId": "",
  "scopedUserId": ""
}
```


---

## GET /auth/api-keys

List all scoped API keys

Returns all scoped keys for the authenticated user.

### Returns

**array** - Array of scoped API keys

### Try It Out

<ApiTester
  method="GET"
  path="/auth/api-keys"
  parameters={[]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/api-keys';
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
curl -X GET 'http://localhost:3010/auth/api-keys' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/api-keys'
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
  const path = '/auth/api-keys';
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
  🔤/auth/api-keys🔤 ➡️ path

  💭 Build HTTP request
  🔤GET /auth/api-keys HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "keys": [
    {
      "createdAt": {
        "Time": "0001-01-01T00:00:00Z",
        "Valid": false
      },
      "dailyLimit": 1000,
      "enabled": true,
      "expiresAt": null,
      "foundryUrl": "",
      "foundryUsername": "",
      "hasFoundryCredentials": false,
      "id": 1,
      "isExpired": false,
      "key": "8ede5230...",
      "name": "test1",
      "requestsToday": 0,
      "scopedClientId": "your-client-id",
      "scopedUserId": "tester",
      "updatedAt": {
        "Time": "2026-03-27T19:00:00.055861824-05:00",
        "Valid": true
      }
    },
    {
      "createdAt": {
        "Time": "0001-01-01T00:00:00Z",
        "Valid": false
      },
      "dailyLimit": 0,
      "enabled": true,
      "expiresAt": null,
      "foundryUrl": "http://localhost:30013",
      "foundryUsername": "tester",
      "hasFoundryCredentials": true,
      "id": 2,
      "isExpired": false,
      "key": "164dd9b7...",
      "name": "headless-test1",
      "requestsToday": 0,
      "scopedClientId": "your-client-id",
      "scopedUserId": "",
      "updatedAt": {
        "Time": "2026-03-27T19:00:00.055861824-05:00",
        "Valid": true
      }
    },
    {
      "createdAt": {
        "Time": "0001-01-01T00:00:00Z",
        "Valid": false
      },
      "dailyLimit": 0,
      "enabled": true,
      "expiresAt": null,
      "foundryUrl": "http://localhost:30013",
      "foundryUsername": "tester",
      "hasFoundryCredentials": true,
      "id": 3,
      "isExpired": false,
      "key": "0a7ff547...",
      "name": "test-headless2-gm",
      "requestsToday": 0,
      "scopedClientId": "foundry-testing-r6bXhB7k9cXa3cif",
      "scopedUserId": "Gamemaster",
      "updatedAt": {
        "Time": "2026-03-27T19:00:00.055861824-05:00",
        "Valid": true
      }
    },
    {
      "createdAt": {
        "Time": "2026-03-25T18:08:33.030028526-05:00",
        "Valid": true
      },
      "dailyLimit": 0,
      "enabled": true,
      "expiresAt": null,
      "foundryUrl": "http://localhost:30013",
      "foundryUsername": "tester",
      "hasFoundryCredentials": true,
      "id": 16,
      "isExpired": false,
      "key": "62b5ae8a...",
      "name": "test-session",
      "requestsToday": 0,
      "scopedClientId": "",
      "scopedUserId": "",
      "updatedAt": {
        "Time": "2026-03-27T19:00:00.055861824-05:00",
        "Valid": true
      }
    },
    {
      "createdAt": {
        "Time": "2026-04-01T13:41:21.466281879-05:00",
        "Valid": true
      },
      "dailyLimit": 500,
      "enabled": true,
      "expiresAt": null,
      "foundryUrl": "",
      "foundryUsername": "",
      "hasFoundryCredentials": false,
      "id": 18,
      "isExpired": false,
      "key": "69de5a63...",
      "name": "Test Scoped Key",
      "requestsToday": 0,
      "scopedClientId": "",
      "scopedUserId": "",
      "updatedAt": {
        "Time": "2026-04-01T13:41:21.466281879-05:00",
        "Valid": true
      }
    }
  ]
}
```


---

## DELETE /auth/api-keys/:id

Delete a scoped API key

Permanently deletes a scoped key.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| id | string | ✓ | params | The scoped key ID |

### Returns

**object** - Success message

### Try It Out

<ApiTester
  method="DELETE"
  path="/auth/api-keys/:id"
  parameters={[{"name":"id","type":"string","required":true,"source":"params"}]}
/>

### Code Examples

<Tabs groupId="programming-language">
<TabItem value="javascript" label="JavaScript">

```javascript
const baseUrl = 'http://localhost:3010';
const path = '/auth/api-keys/18';
const url = `${baseUrl}${path}`;

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
curl -X DELETE 'http://localhost:3010/auth/api-keys/18' \
  -H "x-api-key: your-api-key-here"
```

</TabItem>
<TabItem value="python" label="Python">

```python
import requests

base_url = 'http://localhost:3010'
path = '/auth/api-keys/18'
url = f'{base_url}{path}'

response = requests.delete(
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
  const path = '/auth/api-keys/18';
  const url = `${baseUrl}${path}`;

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
  🔤/auth/api-keys/18🔤 ➡️ path

  💭 Build HTTP request
  🔤DELETE /auth/api-keys/18 HTTP/1.1❌r❌nHost: localhost:3010❌r❌nx-api-key: your-api-key-here❌r❌n❌r❌n🔤 ➡️ request

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
  "message": "API key deleted",
  "success": true
}
```


