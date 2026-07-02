---
tag: session
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Session

## POST /session-handshake

Create a handshake token for secure authentication

**Required scope:** `session:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| x-api-key | string | ✓ | header | API key header |
| x-foundry-url | string | ✓ | header | Foundry URL header |
| x-username | string | ✓ | header | Username header |
| x-world-name | string |  | header | World name header |

### Returns

**object** - Handshake token and encryption details

### Try It Out

<ApiTester
  method="POST"
  path="/session-handshake"
  parameters={[{"name":"x-api-key","type":"string","required":true,"source":"header"},{"name":"x-foundry-url","type":"string","required":true,"source":"header"},{"name":"x-username","type":"string","required":true,"source":"header"},{"name":"x-world-name","type":"string","required":false,"source":"header"}]}
/>

---

## POST /start-session

Start a headless Foundry session using puppeteer

**Required scope:** `session:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| handshakeToken | string | ✓ | body | The token received from session-handshake |
| encryptedPassword | string | ✓ | body | Password encrypted with the public key |
| x-api-key | string | ✓ | header | API key header |
| captureBrowserConsole | string |  | body | Log level for browser console capture ("error", "warn", or "debug") |

### Returns

**object** - Session information including sessionId and clientId

### Try It Out

<ApiTester
  method="POST"
  path="/start-session"
  parameters={[{"name":"handshakeToken","type":"string","required":true,"source":"body"},{"name":"encryptedPassword","type":"string","required":true,"source":"body"},{"name":"x-api-key","type":"string","required":true,"source":"header"},{"name":"captureBrowserConsole","type":"string","required":false,"source":"body"}]}
/>

---

## DELETE /end-session

Stop a headless Foundry session

**Required scope:** `session:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| sessionId | string | ✓ | query | The ID of the session to end |
| x-api-key | string | ✓ | header | API key header |

### Returns

**object** - Status of the operation

### Try It Out

<ApiTester
  method="DELETE"
  path="/end-session"
  parameters={[{"name":"sessionId","type":"string","required":true,"source":"query"},{"name":"x-api-key","type":"string","required":true,"source":"header"}]}
/>

---

## GET /session

Get all active headless Foundry sessions

**Required scope:** `session:manage`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| x-api-key | string | ✓ | header | API key header |

### Returns

**object** - List of active sessions for the current API key

### Try It Out

<ApiTester
  method="GET"
  path="/session"
  parameters={[{"name":"x-api-key","type":"string","required":true,"source":"header"}]}
/>

