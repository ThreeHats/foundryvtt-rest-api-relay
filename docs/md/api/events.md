---
tag: events
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import SseTester from '@site/src/components/SseTester';

# Events

## GET /hooks/subscribe

Subscribe to real-time events via Server-Sent Events (SSE)

Opens a persistent SSE connection for the specified event type. Supported event types: hooks (all Foundry hooks), combat, actor, scene.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| hooks | string |  | query | Comma-separated hook names to filter (hooks type only) |
| encounterId | string |  | query | Filter by encounter ID (combat type only) |
| actorUuid | string |  | query | Actor UUID to subscribe to (actor type, required) |
| sceneId | string |  | query | Scene ID to filter (scene type only) |

### Returns

**SSE stream** - SSE event stream

### Try It Out

<SseTester
  path="/hooks/subscribe"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"hooks","type":"string","required":false,"source":"query"},{"name":"encounterId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /encounters/subscribe

Subscribe to real-time events via Server-Sent Events (SSE)

Opens a persistent SSE connection for the specified event type. Supported event types: hooks (all Foundry hooks), combat, actor, scene.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| hooks | string |  | query | Comma-separated hook names to filter (hooks type only) |
| encounterId | string |  | query | Filter by encounter ID (combat type only) |
| actorUuid | string |  | query | Actor UUID to subscribe to (actor type, required) |
| sceneId | string |  | query | Scene ID to filter (scene type only) |

### Returns

**SSE stream** - SSE event stream

### Try It Out

<SseTester
  path="/encounters/subscribe"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"hooks","type":"string","required":false,"source":"query"},{"name":"encounterId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /actor/subscribe

Subscribe to real-time events via Server-Sent Events (SSE)

Opens a persistent SSE connection for the specified event type. Supported event types: hooks (all Foundry hooks), combat, actor, scene.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| hooks | string |  | query | Comma-separated hook names to filter (hooks type only) |
| encounterId | string |  | query | Filter by encounter ID (combat type only) |
| actorUuid | string |  | query | Actor UUID to subscribe to (actor type, required) |
| sceneId | string |  | query | Scene ID to filter (scene type only) |

### Returns

**SSE stream** - SSE event stream

### Try It Out

<SseTester
  path="/actor/subscribe"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"hooks","type":"string","required":false,"source":"query"},{"name":"encounterId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"}]}
/>

---

## GET /scene/subscribe

Subscribe to real-time events via Server-Sent Events (SSE)

Opens a persistent SSE connection for the specified event type. Supported event types: hooks (all Foundry hooks), combat, actor, scene.

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| hooks | string |  | query | Comma-separated hook names to filter (hooks type only) |
| encounterId | string |  | query | Filter by encounter ID (combat type only) |
| actorUuid | string |  | query | Actor UUID to subscribe to (actor type, required) |
| sceneId | string |  | query | Scene ID to filter (scene type only) |

### Returns

**SSE stream** - SSE event stream

### Try It Out

<SseTester
  path="/scene/subscribe"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"hooks","type":"string","required":false,"source":"query"},{"name":"encounterId","type":"string","required":false,"source":"query"},{"name":"actorUuid","type":"string","required":false,"source":"query"},{"name":"sceneId","type":"string","required":false,"source":"query"}]}
/>

