---
tag: auth
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Auth

## POST /auth/key-request

Request a scoped API key

Initiates a key request that a user must approve in the dashboard. Supports two flows:

**Device flow (no `callbackUrl`)** — for CLI tools, scripts, and desktop apps. The response contains an `approvalUrl` to direct the user to. Poll `GET /auth/key-request/{code}/status` until `status` is `approved`, then read `apiKey` from the response.

**Web flow (`callbackUrl` provided)** — for web apps that can receive an HTTP redirect. After the user approves, the dashboard redirects to your `callbackUrl` with a `code` query parameter containing the exchange code. POST that code to `POST /auth/key-request/exchange` to retrieve the API key.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| appName | string | ✓ | body | Name of the application requesting access |
| scopes | array | ✓ | body | List of permission scopes the key requires |
| appDescription | string |  | body | Short description of what the application does |
| appUrl | string |  | body | Homepage or docs URL for the application |
| callbackUrl | string |  | body | Enables web flow: URL the dashboard redirects to after approval, with a `code` query parameter containing the exchange code |
| clientIds | array |  | body | Foundry client IDs to restrict the key to |
| suggestedMonthlyLimit | number |  | body | Suggested per-key monthly request cap (user may override) |
| suggestedExpiry | string |  | body | Suggested expiry date ISO 8601 (user may override) |

### Returns

**object** - `code`, `approvalUrl`, `expiresIn`, `expiresAt`

### Try It Out

<ApiTester
  method="POST"
  path="/auth/key-request"
  parameters={[{"name":"appName","type":"string","required":true,"source":"body"},{"name":"scopes","type":"array","required":true,"source":"body"},{"name":"appDescription","type":"string","required":false,"source":"body"},{"name":"appUrl","type":"string","required":false,"source":"body"},{"name":"callbackUrl","type":"string","required":false,"source":"body"},{"name":"clientIds","type":"array","required":false,"source":"body"},{"name":"suggestedMonthlyLimit","type":"number","required":false,"source":"body"},{"name":"suggestedExpiry","type":"string","required":false,"source":"body"}]}
/>

---

## GET /auth/key-request/:code/status

Poll key request status

Returns the current status of a pending key request. When `status` is `approved`, the response includes the newly created `apiKey`. Once the key has been retrieved, the status becomes `exchanged` and the key is no longer returned.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| code | string | ✓ | params | The code returned by POST /auth/key-request |

### Returns

**object** - `status` (`pending` | `approved` | `denied` | `expired` | `exchanged`), `apiKey`, `scopes`, `clientIds` (when approved)

### Try It Out

<ApiTester
  method="GET"
  path="/auth/key-request/:code/status"
  parameters={[{"name":"code","type":"string","required":true,"source":"params"}]}
/>

---

## POST /auth/key-request/exchange

Exchange approval code for API key (web flow)

Web flow only. After the user approves the request in the dashboard, the relay redirects to the `callbackUrl` with a one-time `code` query parameter. POST that code here to receive the API key. The code is single-use — a second call returns 410.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| code | string | ✓ | body | The exchange code delivered to your callbackUrl |

### Returns

**object** - `apiKey`, `scopes`, `clientIds`

### Try It Out

<ApiTester
  method="POST"
  path="/auth/key-request/exchange"
  parameters={[{"name":"code","type":"string","required":true,"source":"body"}]}
/>

