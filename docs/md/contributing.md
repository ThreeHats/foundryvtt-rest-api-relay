---
id: contributing
title: Contributing
sidebar_position: 15
---

# Contributing

Welcome! This guide covers everything you need to know to contribute to the Foundry REST API project. Whether you're fixing bugs, adding features, or improving documentation, this guide will help you understand the codebase architecture and development patterns.

## Project Architecture

The Foundry REST API consists of two interconnected repositories:

### 1. Relay Server (`foundryvtt-rest-api-relay`)

The relay server is a Go application that:
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
│  (Your Code)    │◀───────────── │  (Go/Chi)       │               │  (Module)       │
└─────────────────┘               └─────────────────┘               └─────────────────┘
```

## Development Setup

### Prerequisites

- **Go 1.22+** for the relay server
- **Node.js 18+** and **pnpm** for frontend, docs, and tests
- **Foundry VTT** with a valid license
- **Git** for version control

### Setting Up the Relay Server

1. **Fork the repository** on GitHub: [ThreeHats/foundryvtt-rest-api-relay](https://github.com/ThreeHats/foundryvtt-rest-api-relay)
2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/foundryvtt-rest-api-relay.git
cd foundryvtt-rest-api-relay

# Start development server (SQLite)
pnpm run local:sqlite
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
├── go-relay/
│   ├── cmd/
│   │   ├── server/main.go       # Application entry point
│   │   └── docgen/main.go       # API documentation generator
│   └── internal/
│       ├── config/              # Environment variable loading
│       ├── database/            # DB initialization, migrations (SQLite + Postgres)
│       ├── server/              # Chi router setup, middleware wiring
│       ├── handler/             # HTTP route handlers
│       │   ├── helpers/         # Response helpers, parameter extraction
│       │   ├── entity.go        # Entity CRUD, search, rolls, encounters, etc.
│       │   ├── dnd5e.go         # D&D 5e system-specific endpoints
│       │   ├── auth.go          # Auth routes + scoped API key management
│       │   ├── stripe.go        # Stripe billing
│       │   └── ...
│       ├── middleware/          # Auth, rate limiting, request forwarding
│       ├── ws/                  # WebSocket client management, message routing
│       ├── model/               # Database models (User, ApiKey, etc.)
│       ├── service/             # Email, encryption, API key validation
│       ├── worker/              # Headless browser session management
│       └── cron/                # Scheduled jobs (daily/monthly resets)
├── frontend/                    # Astro/Svelte dashboard frontend
├── docs/                        # Docusaurus documentation site
├── tests/                       # Jest integration tests
└── scripts/                     # Build and utility scripts
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

### Step 1: Add the Route Handler (Relay Server)

Create or update a handler in `go-relay/internal/handler/`:

```go
// go-relay/internal/handler/entity.go (or a new file)

// My feature endpoint description
//
// Detailed description of what this endpoint does.
// Go doc comments with @tag and @param annotations are used to auto-generate API documentation.
// @tag MyFeature
// @param {string} clientId [query] Client ID for the Foundry world
// @param {string} targetUuid [body] UUID of the target entity
// @param {number} amount [body] Amount to modify (default: 1)
// @returns Result of the operation
var myFeatureAction = helpers.EndpointConfig{
    Type: "my-action",
    RequiredParams: []helpers.ParamDef{
        {Name: "targetUuid", From: "body"},
    },
    OptionalParams: []helpers.ParamDef{
        {Name: "amount", From: "body", Type: "number"},
    },
}
```

### Step 2: Register the Route (Relay Server)

Add the route in `go-relay/internal/handler/routes.go`:

```go
// In RegisterAPIRoutes function
r.Post("/my-feature/do-something", h(mgr, pending, myFeatureAction))
```

:::caution Important
New routes **must** be registered in `routes.go` or they won't be accessible.
:::

### Step 3: Add the Handler (Module)

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

### Step 4: Register the Router in the Module

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

### Step 5: Write Tests

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

**Relay:** `go-relay/internal/handler/dnd5e.go` - HTTP endpoints for D&D 5e
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

**Relay** (Go - zerolog):
```go
import "github.com/rs/zerolog/log"

log.Info().Str("key", "value").Msg("Processing request")
log.Warn().Msg("Potential issue detected")
log.Error().Err(err).Msg("Error occurred")
```

Use the appropriate logger for each project. Never use `console.log` (module) or `fmt.Println` (relay) directly.

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
- [ ] Doc comments on all new endpoints (with correct `@route` paths)
- [ ] Route registered in `routes.go` (relay) and router added to `all.ts` (module)
- [ ] No `fmt.Println` statements in Go (use zerolog)
- [ ] No `console.log` statements in TypeScript (use ModuleLogger)
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
