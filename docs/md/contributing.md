---
id: contributing
title: Contributing
sidebar_position: 11
---

# Contributing

Welcome! This guide covers everything you need to know to contribute to the Foundry REST API project. Whether you're fixing bugs, adding features, or improving documentation, this guide will help you understand the codebase architecture and development patterns.

## Project Architecture

The Foundry REST API consists of two interconnected repositories:

### 1. Relay Server (`foundryvtt-rest-api-relay`)

The relay server is an Express.js application that:
- Provides HTTP REST endpoints for external clients
- Manages WebSocket connections to Foundry clients
- Handles authentication and session management
- Routes requests between HTTP clients and Foundry instances

### 2. Foundry Module (`foundryvtt-rest-api`)

The Foundry module is a TypeScript module that:
- Runs inside Foundry VTT as a GM-only module
- Connects to the relay server via WebSocket
- Handles incoming requests from the relay
- Executes Foundry API calls and returns results

```
┌─────────────────┐     HTTP      ┌─────────────────┐   WebSocket   ┌─────────────────┐
│  External App   │ ─────────────▶│  Relay Server   │◀─────────────▶│  Foundry VTT    │
│  (Your Code)    │◀───────────── │  (Express.js)   │               │  (Module)       │
└─────────────────┘               └─────────────────┘               └─────────────────┘
```

## Development Setup

### Prerequisites

- **Node.js 18+** and **pnpm** package manager
- **Foundry VTT** with a valid license
- **Git** for version control

### Setting Up the Relay Server

1. **Fork the repository** on GitHub: [ThreeHats/foundryvtt-rest-api-relay](https://github.com/ThreeHats/foundryvtt-rest-api-relay)
2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/foundryvtt-rest-api-relay.git
cd foundryvtt-rest-api-relay
pnpm install

# Start development server (in-memory database)
pnpm local
```

### Setting Up the Module

1. **Fork the repository** on GitHub: [ThreeHats/foundryvtt-rest-api](https://github.com/ThreeHats/foundryvtt-rest-api)
2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/foundryvtt-rest-api.git
cd foundryvtt-rest-api

pnpm install
```

3. Create a `.env` file to specify your Foundry modules directory:

```bash
# .env
FOUNDRY_VTT_DATA_MODULES_PATH="/path/to/your/FoundryVTT/Data/modules"
```

4. Build the module (it will be placed in your specified modules directory):

```bash
pnpm build
```

:::note
The build process automatically copies the module to your Foundry modules directory based on the `FOUNDRY_VTT_DATA_MODULES_PATH` environment variable. If not set, it defaults to the Windows AppData location.
:::

## Repository Structure

### Relay Server Structure

```
foundryvtt-rest-api-relay/
├── src/
│   ├── index.ts              # Application entry point
│   ├── routes/
│   │   ├── api/              # API route handlers
│   │   │   ├── entity.ts     # Entity CRUD endpoints
│   │   │   ├── roll.ts       # Dice rolling endpoints
│   │   │   ├── session.ts    # Session management
│   │   │   └── ...
│   │   ├── route-helpers.ts  # createApiRoute helper
│   │   └── shared.ts         # Pending requests tracking
│   ├── core/
│   │   └── ClientManager.ts  # WebSocket client management
│   ├── middleware/
│   │   ├── auth.ts           # API key authentication
│   │   └── requestForwarder.ts
│   └── utils/
│       └── logger.ts
├── tests/
│   ├── integration/          # API integration tests
│   └── helpers/              # Test utilities
├── docs/                     # Documentation site (Docusaurus)
└── scripts/                  # Build and utility scripts
```

### Module Structure

```
foundryvtt-rest-api/
├── src/
│   ├── ts/
│   │   ├── module.ts         # Module entry point
│   │   ├── settings.ts       # Module settings
│   │   ├── constants.ts      # Module ID, settings keys
│   │   ├── types.ts          # TypeScript types
│   │   ├── network/
│   │   │   ├── webSocketManager.ts    # WebSocket connection
│   │   │   ├── webSocketEndpoints.ts  # Router registration
│   │   │   └── routers/
│   │   │       ├── baseRouter.ts      # Router base class
│   │   │       ├── all.ts             # Router exports
│   │   │       ├── entity.ts          # Entity handlers
│   │   │       └── ...
│   │   ├── systems/          # Game system integrations
│   │   │   ├── IRestApiSystem.ts
│   │   │   ├── dnd5e.ts
│   │   │   └── a5e.ts
│   │   └── utils/
│   │       ├── logger.ts
│   │       ├── serialization.ts
│   │       └── search.ts
│   ├── module.json           # Foundry module manifest
│   └── styles/
└── tests/                    # Unit tests
```

## Adding a New Endpoint

Adding a new API endpoint requires changes to both repositories. Here's the complete process:

### Step 1: Add the Route in the Relay Server

Create or update a route file in `src/routes/api/`:

```typescript
// src/routes/api/myFeature.ts
import { Router } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';

export const myFeatureRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

/**
 * My feature endpoint
 *
 * Detailed description of what this endpoint does.
 * JSDoc comments are REQUIRED for all endpoints - they are used to auto-generate API documentation.
 * 
 * @route POST /my-feature/do-something
 * @returns {object} Result object
 * 
 * NOTE: The @route tag MUST include the full path including any router prefix.
 * For example, if this router is mounted at '/dnd5e', use: @route POST /dnd5e/do-something
 */
myFeatureRouter.post("/do-something", ...commonMiddleware, express.json(), createApiRoute({
    type: 'my-action',  // Message type sent to Foundry
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Always required
        { name: 'targetUuid', from: 'body', type: 'string' }
    ],
    optionalParams: [
        { name: 'amount', from: 'body', type: 'number' },
        { name: 'options', from: 'body', type: 'object' }
    ],
    validateParams: (params) => {
        // Custom validation logic
        if (params.amount && params.amount < 0) {
            return { error: "'amount' must be positive" };
        }
        return null;
    },
    timeout: 15000  // Optional: custom timeout in ms
}));
```

### JSDoc and Parameter Documentation

**All endpoints MUST have JSDoc comments.** These comments are parsed by `scripts/generateApiDocs.js` to auto-generate API documentation.

JSDoc comments should include:
- A brief description of what the endpoint does
- `@route` tag with the HTTP method and **full path** (including router prefix)
- `@returns` tag describing the response

#### createApiRoute Endpoints

For endpoints using `createApiRoute`, parameter documentation is done via **inline comments** next to each parameter definition. The documentation generator automatically extracts these:

```typescript
createApiRoute({
    type: 'my-action',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'targetUuid', from: 'body', type: 'string' } // UUID of the target entity
    ],
    optionalParams: [
        { name: 'amount', from: 'body', type: 'number' } // Amount to modify (default: 1)
    ]
})
```

The `from`, `type`, and required/optional status are auto-generated from the configuration - you just need the description comment.

#### Traditional Endpoints (Non-createApiRoute)

For traditional Express handlers, use full JSDoc `@param` tags with this format:

```typescript
/**
 * Get actor sheet HTML
 * 
 * @route GET /sheet
 * @param {string} uuid - [query,?] The UUID of the entity
 * @param {boolean} selected - [query,?] Whether to use the selected entity
 * @param {string} clientId - [query] The ID of the Foundry client (required)
 * @returns {object} The sheet HTML or data
 */
```

The bracket notation `[source,?]` indicates:
- `source`: Where the param comes from (`query`, `body`, `params`)
- `?`: Optional (omit for required params)

For system-specific endpoints (like D&D 5e), the `@route` tag must include the system prefix:
```typescript
/**
 * @route GET /dnd5e/get-actor-details  // Include the /dnd5e prefix!
 */
```

### Traditional Endpoints

While `createApiRoute` is preferred, traditional Express route handlers are still valid if (and only if) `createApiRoute` cannot handle the functionality. They must:
- Include JSDoc comments with the same requirements
- Handle parameter validation manually
- Follow the same response patterns
- Include proper error handling

### Understanding `createApiRoute`

The `createApiRoute` helper standardizes API route handling:

```typescript
interface ApiRouteConfig {
  // Message type identifier - must match the handler in the module
  type: PendingRequestType;
  
  // Required parameters - request fails if missing
  requiredParams?: ParamDef[];
  
  // Optional parameters
  optionalParams?: ParamDef[];
  
  // Custom timeout (default: 10000ms)
  timeout?: number;
  
  // Custom validation function (can be sync or async)
  validateParams?: (params, req) => { error?: string; howToUse?: string } | null | Promise<...>;
  
  // Transform payload before sending to Foundry (can be sync or async)
  buildPayload?: (params, req) => Record<string, any> | Promise<Record<string, any>>;
  
  // Add extra data to pending request tracking
  buildPendingRequest?: (params) => Partial<PendingRequest>;
}

interface ParamDef {
  name: string;
  from: 'body' | 'query' | 'params' | ('body' | 'query' | 'params')[];
  type?: 'string' | 'number' | 'boolean' | 'array' | 'object';
}
```

### Step 2: Register the Router (Relay)

Add your router to `src/routes/api.ts`:

```typescript
// Add import at the top with other router imports
import { myFeatureRouter } from './api/myFeature';

// Add to router registration section (around line 376)
app.use('/my-feature', myFeatureRouter);
// Or use '/' prefix if the routes already include the path:
// app.use('/', myFeatureRouter);
```

:::caution Important
New routers **must** be registered in `api.ts` or they won't be accessible.
:::

### Step 3: Add the Handler in the Module

Create or update a router file in `src/ts/network/routers/`:

```typescript
// src/ts/network/routers/myFeature.ts
import { Router } from "./baseRouter";
import { ModuleLogger } from "../../utils/logger";
import { deepSerializeEntity } from "../../utils/serialization";

export const router = new Router("myFeatureRouter");

router.addRoute({
  actionType: "my-action",  // Must match 'type' in relay createApiRoute
  handler: async (data, context) => {
    const socketManager = context?.socketManager;
    ModuleLogger.info(`Received my-action request:`, data);

    try {
      // Your Foundry API logic here
      const entity = await fromUuid(data.targetUuid);
      
      if (!entity) {
        socketManager?.send({
          type: "my-action-result",
          requestId: data.requestId,
          error: "Entity not found",
        });
        return;
      }

      // Do something with the entity
      const result = await entity.someMethod(data.amount);

      // Send success response
      socketManager?.send({
        type: "my-action-result",
        requestId: data.requestId,
        data: deepSerializeEntity(result),
      });
    } catch (error) {
      ModuleLogger.error(`Error in my-action:`, error);
      socketManager?.send({
        type: "my-action-result",
        requestId: data.requestId,
        error: (error as Error).message,
      });
    }
  },
});
```

### Understanding `addRoute`

The module's router system is simpler:

```typescript
interface RouteI {
  // Message type to listen for (from relay)
  actionType: string;
  
  // Handler function
  handler: (data: any, context: HandlerContext | undefined) => void;
}

interface HandlerContext {
  socketManager: WebSocketManager;
}
```

**Important patterns:**
- Always include `requestId` in responses for request correlation
- Response type should be `{actionType}-result`
- Use `deepSerializeEntity()` for entity data to handle circular references
- Log errors with `ModuleLogger.error()`

### Step 4: Register the Router (Module)

Add your router to `src/ts/network/routers/all.ts`:

```typescript
import { router as MyFeatureRouter } from "./myFeature";

export const routers: Router[] = [
    // ... existing routers
    MyFeatureRouter,
];
```

:::caution Important
New module routers **must** be added to the `routers` array in `all.ts` or they won't receive WebSocket messages.
:::

### Step 5: Add the Request Type (Relay)

Update `src/routes/shared.ts` to add your new request type to the `PENDING_REQUEST_TYPES` array:

```typescript
export const PENDING_REQUEST_TYPES = [
    'search', 'entity', 'structure', 'contents', 'create', 'update', 'delete',
    // ... existing types ...
    'my-action'  // Add your new type
] as const;
```

:::caution Important
New request types **must** be added to `PENDING_REQUEST_TYPES` or the relay won't be able to track pending requests and match responses.
:::

The response handling is automatic through `createApiRoute`. When Foundry sends back a `{type}-result` message, the relay's WebSocket handler matches it to the pending request and sends the HTTP response.

### Step 6: Write Tests

Add integration tests for your new endpoint:

```typescript
// tests/integration/myFeature-endpoints.test.ts
import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { getEntityUuid } from '../helpers/testEntities';
import * as path from 'path';

const capturedExamples: any[] = [];

describe('My Feature', () => {
  afterAll(() => {
    const outputPath = path.join(__dirname, '../../docs/examples/myFeature-examples.json');
    saveExamples(capturedExamples, outputPath);
  });

  forEachVersion((version, getClientId) => {
    describe(`/my-feature (v${version})`, () => {
      test('POST /my-feature/do-something', async () => {
        setVariable('clientId', getClientId());
        
        const targetUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(targetUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/my-feature/do-something',
            host: ['{{baseUrl}}'],
            path: ['my-feature', 'do-something'],
            query: [
              { key: 'clientId', value: '{{clientId}}' }
            ]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              targetUuid: targetUuid,
              amount: 10
            })
          }
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/my-feature/do-something'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
      });
    });
  });
});
```

Add your test to the test sequencer in `tests/helpers/testSequencer.ts` at an appropriate position (before cleanup tests).

:::tip Test Generation
You can auto-generate test file boilerplate with `pnpm test:generate`. The generated tests need manual input for correct parameter values and assertions. See the [Testing Documentation](./testing.md) for complete details on running and writing tests.
:::

:::caution Important
New test files **must** be added to `TEST_ORDER` in `testSequencer.ts` or they won't run as part of the test suite.
:::

## Game System Integration

The module supports system-specific functionality in two ways:

### 1. System Configuration (`src/ts/systems/`)

The systems architecture provides configuration values (like attribute paths) that differ between game systems. Currently, this is minimal:

```typescript
// src/ts/systems/IRestApiSystem.ts
export interface IRestApiSystem {
    ACTOR_CURRENCY_ATTRIBUTE: string;
}
```

:::note Work in Progress
The system configuration architecture exists but is largely unused. Most system-specific functionality is handled in the router files directly. This area needs further development.
:::

### 2. System-Specific Routers (Recommended Approach)

The primary way to add system-specific functionality is through dedicated router files:

**Relay:** `src/routes/api/dnd5e.ts` - HTTP endpoints for D&D 5e
**Module:** `src/ts/network/routers/dnd5e.ts` - WebSocket handlers for D&D 5e

System-specific routers typically:
1. Check `game.system.id` before registering routes
2. Use `Hooks.once('init', ...)` to defer registration until the system is loaded
3. Access system-specific data structures directly

```typescript
// Example from dnd5e.ts in the module
Hooks.once('init', () => {
    const isDnd5e = game.system.id === "dnd5e";
    if (isDnd5e) {
        router.addRoute({
            actionType: "get-actor-details",
            handler: async (data, context) => {
                // D&D 5e specific logic
            }
        });
    }
});
```

To add support for a new game system, create new router files following the D&D 5e pattern rather than modifying the systems configuration.

## Utility Functions

### Serialization (`utils/serialization.ts` - Module)

```typescript
import { deepSerializeEntity } from "../../utils/serialization";

// Safely serialize Foundry entities with circular reference handling
const serialized = deepSerializeEntity(actor);
```

### Search Utilities (`utils/search.ts` - Module)

```typescript
import { parseFilterString, matchesAllFilters } from "../../utils/search";

// Parse a filter string like "documentType:Actor,folder:zmAZJmay9AxvRNqh"
const filters = parseFilterString("documentType:Actor,folder:zmAZJmay9AxvRNqh");

// Check if a search result matches all filters
const matches = matchesAllFilters(searchResult, filters);
```

### Logging

Both projects have their own logging implementations:

**Module** (`src/ts/utils/logger.ts`):
```typescript
import { ModuleLogger } from "../../utils/logger";

ModuleLogger.info("Processing request:", data);
ModuleLogger.warn("Potential issue detected");
ModuleLogger.error("Error occurred:", error);
```

**Relay** (`src/utils/logger.ts`):
```typescript
import { log } from '../utils/logger';

log.info("Processing request:", data);
log.warn("Potential issue detected");
log.error("Error occurred:", error);
```

Use the appropriate logger for each project. Never use `console.log` directly.

## Pull Request Guidelines

### Before Submitting

1. **Fork the repository** and create a feature branch
2. **Write tests** for new functionality (see [Testing Documentation](./testing.md))
3. **Run the full test suite** - see [Testing Documentation](./testing.md) for setup and execution
4. **Update documentation** for API changes (auto-generated via `pnpm docs:full`)
5. **Follow code style** (consistent with existing code)

:::tip Documentation Changes
Running tests regenerates documentation examples, which will modify files for endpoints you didn't change. **Discard changes to documentation files for endpoints you're not working on** before committing.
:::

### PR Checklist

- [ ] Tests pass locally
- [ ] New endpoints have corresponding tests
- [ ] Test file added to `TEST_ORDER` in `testSequencer.ts`
- [ ] JSDoc comments on all new endpoints (with correct `@route` paths)
- [ ] New routers registered in `api.ts` (relay) and `all.ts` (module)
- [ ] New request types added to `PENDING_REQUEST_TYPES` in `shared.ts`
- [ ] No `console.log` statements (use appropriate logger)
- [ ] TypeScript types are complete
- [ ] Both relay and module changes coordinated
- [ ] Unrelated documentation changes discarded

### Commit Message Format

This part is a guideline, not necessarily strict rules.

```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Examples:
```
feat(entity): add bulk update endpoint
fix(session): handle disconnect during handshake
docs(api): update authentication examples
```

## Code Style Guidelines

### TypeScript

- Use explicit types (avoid `any` where possible)
- Use `async/await` over raw promises
- Document public functions with JSDoc
- Handle errors appropriately

### Relay Routes

- Use `createApiRoute` for standardized handling
- Include detailed JSDoc comments for documentation generation
- Validate all user input
- Return consistent error formats

### Module Handlers

- Always include `requestId` in responses
- Use `ModuleLogger` for logging
- Serialize entities with `deepSerializeEntity`
- Handle both success and error cases

## Questions and Support

- **Issues**: Use GitHub Issues for bugs and feature requests
- **Discord**: [Join our Discord](https://discord.gg/U634xNGRAC) for community support and discussions

## Useful Resources

- **Foundry VTT API Documentation**: [https://foundryvtt.com/api/](https://foundryvtt.com/api/)
- **Foundry VTT Wiki**: [https://foundryvtt.wiki/](https://foundryvtt.wiki/)

Thank you for contributing to the Foundry REST API project!
