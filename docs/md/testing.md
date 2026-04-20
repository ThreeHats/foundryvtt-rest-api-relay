---
id: testing
title: Testing
sidebar_position: 14
---

# Testing

This comprehensive guide covers everything you need to know about testing the Foundry REST API, including setup, running tests, understanding test infrastructure, and generating documentation examples.

## Overview

The Foundry REST API uses a testing framework that:

- Tests against **multiple Foundry VTT versions** simultaneously
- Supports both **headless browser automation** and **existing sessions**
- **Captures real API responses** for documentation generation
- Manages **entity lifecycle** with automatic cleanup
- Sources test data from **system compendiums** for properly initialized entities
- Runs **system-specific tests** (e.g., D&D 5e) only on matching Foundry instances
- Tests **permission filtering** with the `userId` parameter
- Uses a **custom test sequencer** to ensure proper execution order

## Prerequisites

Before running tests, you need:

1. **Relay Server**: The REST API relay server running (this project)
2. **Foundry Instances**: One or more Foundry VTT instances running with the REST API module installed
3. **Test Worlds**: Foundry worlds ready for testing (entities will be created/deleted during tests)

## Environment Configuration

### Setting Up the Test Environment

Create a `.env.test` file based on the example:

```bash
cp .env.test.example .env.test
```

### Configuration Options

#### Required Settings

```bash
# Relay Server Configuration
TEST_BASE_URL=http://localhost:3010
TEST_API_KEY=your-relay-api-key-here

# Foundry Instance URLs
FOUNDRY_V12_URL=http://localhost:30012
FOUNDRY_V13_URL=http://localhost:30013

# Which versions to test (comma-separated)
TEST_FOUNDRY_VERSIONS=12,13
```

#### Session Mode: Existing Session (Recommended for Development)

Using an existing session is faster and doesn't require browser automation:

```bash
USE_EXISTING_SESSION=true

# Client IDs per version (get from /clients endpoint or browser console)
TEST_CLIENT_ID_V12=foundry-abc123
TEST_CLIENT_ID_V13=foundry-xyz789
```

:::note
When using existing sessions, the session endpoint tests (`session-endpoints.test.ts`) are automatically skipped since session creation/authentication is not needed.
:::

**How to get your client ID:**
1. Open Foundry in your browser and log in to a world as a GM
2. The REST API module will connect to the relay server
3. Visit `http://localhost:3010/clients` with your API key header to see connected clients
4. Or check the browser console for the client ID

:::caution Cookie Isolation
When using existing sessions with multiple Foundry versions in the same browser, be aware that browsers share cookies across localhost ports. This can cause authentication issues. Use different browser profiles or incognito windows for each Foundry version.
:::

#### Session Mode: Headless Automation (For CI/CD)

For automated testing, the framework can create headless browser sessions:

```bash
USE_EXISTING_SESSION=false

# Foundry login credentials
FOUNDRY_USERNAME=Gamemaster
FOUNDRY_PASSWORD=your-password
```

:::caution Known Limitation
Headless mode may fail if there are already GMs connected to the Foundry world, as they may retain "primary GM" status and the relay will connect to them instead of the headless session. Ensure no other GM clients are connected when running headless tests. This will be improved in a future update.
:::

```bash

# Default world name
TEST_DEFAULT_WORLD=test-world

# Version-specific world names (optional)
FOUNDRY_V12_WORLD=test-world-v12
FOUNDRY_V13_WORLD=test-world-v13
```

#### Permission Filtering Tests (Optional)

To test `userId` permission filtering with a real non-GM player:

```bash
# Create a non-GM player in your Foundry world, then provide their ID or username
TEST_PLAYER_USER_ID_V12=Player1
TEST_PLAYER_USER_ID_V13=Player1
```

Without this, only invalid-userId error handling tests run. With it, additional tests verify:
- Search result filtering (player sees fewer results than GM)
- Write permission denial (player can't update/delete GM-owned entities)
- Chat whisper visibility (player can't see GM-only whispers)

## Running Tests

### Full Test Suite

```bash
pnpm test
```

This runs all integration tests in the correct order using a custom test sequencer.

### What Happens During a Test Run

1. **Session Setup** (`session-endpoints.test.ts`)
   - Creates or validates Foundry sessions
   - Stores client IDs, session data, and system IDs for other tests

2. **Entity Creation** (`entity-endpoints.test.ts`)
   - Fetches entity data from system compendiums (e.g., `dnd5e.heroes` for actors)
   - Falls back to minimal `type: 'base'` entities if compendium lookup fails
   - Creates test actors, items, journals, macros
   - Registers entities for automatic cleanup

3. **Auth Validation** (`auth-requirements.test.ts`)
   - Dynamically scans all route files for endpoints
   - Verifies API key and clientId requirements on every endpoint

4. **Scene + Canvas** (`scene-endpoints.test.ts`, `canvas-endpoints.test.ts`)
   - Creates and activates a test scene
   - Tests full CRUD on all 8 canvas document types (tokens, tiles, drawings, lights, sounds, notes, templates, walls)
   - Creates a persistent token for utility/encounter tests

5. **Core Functionality Tests**
   - Structure, search, roll (including SSE streaming), sheet, macro endpoints
   - Utility endpoints (select, players, world info, token movement, measure distance)
   - Encounter management (uses persistent token on the test scene; rolls initiative when compendium actors are available)
   - File system operations (upload, download, browse)
   - Chat (including SSE streaming, whisper, flush)
   - Playlist control and sound playback

6. **Permission Filtering** (`permission-filtering.test.ts`)
   - Dynamically scans route files for all endpoints accepting `userId`
   - Tests invalid userId rejection on every endpoint
   - Tests player permission scoping (if `TEST_PLAYER_USER_ID` is configured)

6b. **Feature Tests**
   - Scoped API key management and auto-routing
   - Hook and combat event subscriptions (WebSocket)
   - Scene screenshots and raw background images
   - User CRUD operations
   - WebSocket API relay

7. **System-Specific Tests** (`dnd5e-endpoints.test.ts`)
   - Only runs on Foundry instances with the matching system
   - Sets up system-specific test data (gives actor spells and consumables from compendiums)
   - Tests actor details, modify experience, use abilities/spells/items
   - Skill checks, ability saves, death saves, short/long rest
   - Concentration tracking and saves
   - Inventory management (equip, attune, transfer currency)

8. **Cleanup** (`cleanup-entities.test.ts`, `scene-cleanup.test.ts`, `end-sessions.test.ts`, `auth-cleanup.test.ts`)
   - Deletes all created test entities
   - Restores the original active scene and deletes the test scene
   - Closes headless sessions (if applicable)

### Running Specific Tests

:::caution Test Dependencies
Most test files cannot be run individually. Tests depend on shared state (client IDs, entity UUIDs) that is established by earlier tests in the sequence. For example, `entity-endpoints.test.ts` requires `session-endpoints.test.ts` to have run first to set up the client ID.

**Always run the full test suite with `pnpm test`** to ensure proper test ordering and state management.

This may be improved in the future.
:::

```bash
# Run the full test suite
pnpm test
```

### Skipping Cleanup (For Debugging)

To keep test entities for inspection:

```bash
SKIP_CLEANUP=true pnpm test
```

### Validate Setup

Before running full tests, validate your configuration:

```bash
pnpm validate-setup
```

This checks:
- `.env.test` file exists and has required variables
- Relay server connectivity (at `TEST_BASE_URL`)
- Foundry instance availability (for each `FOUNDRY_VX_URL`)
- API key authentication

## Test Architecture

### Directory Structure

```
tests/
├── integration/                        # API integration tests
│   ├── session-endpoints.test.ts       # Session management (runs first)
│   ├── client-endpoints.test.ts        # Connected client listing
│   ├── entity-endpoints.test.ts        # Entity CRUD (creates test data)
│   ├── auth-requirements.test.ts       # Auth validation (scans all routes)
│   ├── auth-endpoints.test.ts          # Auth endpoint CRUD tests
│   ├── scene-endpoints.test.ts         # Scene CRUD + activate test scene
│   ├── canvas-endpoints.test.ts        # All canvas document types
│   ├── structure-endpoints.test.ts     # World structure queries
│   ├── search-endpoints.test.ts        # Search functionality
│   ├── roll-endpoints.test.ts          # Dice rolling + SSE streaming
│   ├── sheet-endpoints.test.ts         # Actor sheet screenshots
│   ├── macro-endpoints.test.ts         # Macro execution
│   ├── utility-endpoints.test.ts       # Utility endpoints (select, players, world-info, etc.)
│   ├── encounter-endpoints.test.ts     # Combat/encounters
│   ├── fileSystem-endpoints.test.ts    # File operations
│   ├── chat-endpoints.test.ts          # Chat messages + SSE streaming
│   ├── permission-filtering.test.ts    # userId permission scoping
│   ├── scoped-keys-endpoints.test.ts   # Scoped API key CRUD + auto-routing
│   ├── playlist-endpoints.test.ts      # Playlist control + sound playback
│   ├── hooks-subscribe.test.ts         # Hook event subscriptions (SSE/WS)
│   ├── combat-subscribe.test.ts        # Combat event subscriptions (WS)
│   ├── scene-image.test.ts             # Scene screenshots + raw backgrounds
│   ├── user-endpoints.test.ts          # User CRUD (self-contained)
│   ├── websocket-api.test.ts           # Raw WebSocket API relay
│   ├── effects-endpoints.test.ts       # ActiveEffect CRUD + status effects list
│   ├── dnd5e-endpoints.test.ts         # D&D 5e system tests (includes concentration + inventory)
│   ├── cleanup-entities.test.ts        # Deletes test entities
│   ├── scene-cleanup.test.ts           # Restores original scene
│   ├── end-sessions.test.ts            # Session cleanup (runs last)
│   └── auth-cleanup.test.ts            # Auth test user cleanup
├── helpers/                            # Test utilities
│   ├── apiRequest.ts                   # HTTP request helper
│   ├── captureExample.ts               # Documentation capture (HTTP)
│   ├── captureWsExample.ts             # Documentation capture (WebSocket)
│   ├── compendiumData.ts               # Compendium entity fetching
│   ├── endpointMetadata.ts             # Dynamic endpoint discovery
│   ├── globalVariables.ts              # Cross-file state
│   ├── multiVersion.ts                 # Multi-version test runner
│   ├── systemSetup.ts                  # Per-system test data setup
│   ├── testEntities.ts                 # Entity lifecycle management
│   ├── testSequencer.ts                # Test execution order
│   ├── testVariables.ts                # Environment variables
│   └── wsClient.ts                     # WebSocket test client
├── setup.ts                            # Jest setup (loads .env.test)
└── globalTeardown.ts                   # Cleanup after all tests
```

### Test Execution Order

Tests run in a specific order defined in `tests/helpers/testSequencer.ts`. The order follows these phases:

1. **Session Setup** - Creates/validates Foundry sessions
2. **Entity Creation** - Creates test actors, items, journals, etc. (compendium-sourced when possible)
3. **Auth Validation** - Verifies authentication requirements
4. **Scene + Canvas** - Creates a test scene, activates it, tests canvas CRUD for all document types
5. **Core Functionality** - Tests all remaining API endpoints (structure, search, roll, sheet, macro, utility, encounter, fileSystem, chat)
6. **Permission Filtering** - Tests userId validation and player permission scoping
7. **System-Specific** - D&D 5e tests (only on dnd5e systems)
8. **Cleanup** - Deletes test entities, restores original scene, closes sessions

:::caution Important
New test files **must** be added to `TEST_ORDER` in `testSequencer.ts` at an appropriate position. Tests not in this array won't run as part of the suite. See the actual file for the current test order.
:::

## Test Helpers

### Multi-Version Testing (`multiVersion.ts`)

Run tests across all configured Foundry versions:

```typescript
import { forEachVersion } from '../helpers/multiVersion';

forEachVersion((version, getClientId) => {
  describe(`/my-endpoint (v${version})`, () => {
    test('GET /my-endpoint', async () => {
      const clientId = getClientId();
      // ... test code
    });
  });
});
```

For system-specific tests:

```typescript
import { forEachVersionWithSystem } from '../helpers/multiVersion';

forEachVersionWithSystem('dnd5e', (version, getClientId) => {
  describe(`D&D 5e tests (v${version})`, () => {
    // Only runs on clients with dnd5e system
  });
});
```

### Compendium Data (`compendiumData.ts`)

Fetch properly initialized entity data from the active system's compendiums:

```typescript
import { fetchCompendiumEntityData } from '../helpers/compendiumData';

// Fetch any actor from the system's compendiums
const actorData = await fetchCompendiumEntityData(version, 'Actor');

// Fetch from a specific pack
const heroData = await fetchCompendiumEntityData(version, 'Actor', {
  packFilter: 'heroes',
});

// Fetch a specific entry by name
const spellData = await fetchCompendiumEntityData(version, 'Item', {
  packFilter: 'spells',
  entryFilter: 'fireball',
});
```

Compendium data is:
- **Cached** per version+entityType+options (one `/structure` call per combination per run)
- **Cleaned** for creation (`_id`, `_stats`, `folder`, `ownership`, `flags` removed)
- **Prefixed** with `test-` in the name for identification
- **System-agnostic** — works with any system (dnd5e, pf2e, etc.) by filtering packs by `systemId`

### System Setup (`systemSetup.ts`)

Set up system-specific test data on actors (spells, consumables, etc.):

```typescript
import { setupSystemTestData } from '../helpers/systemSetup';

const result = await setupSystemTestData(version, actorUuid);
// result: { spellName: 'Acid Splash' | null, consumableName: 'Potion of Healing' | null }
```

To add support for a new system, add a `SystemTestConfig` entry to the `SYSTEM_CONFIGS` map in `systemSetup.ts`:

```typescript
const pf2eConfig: SystemTestConfig = {
  id: 'pf2e',
  compendiumPacks: { actors: 'heroes', spells: 'spells-srd', consumables: 'equipment-srd' },
  giveSpell: async (version, actorUuid) => { /* pf2e-specific logic */ },
  giveConsumable: async (version, actorUuid) => { /* pf2e-specific logic */ },
};
SYSTEM_CONFIGS['pf2e'] = pf2eConfig;
```

### Entity Management (`testEntities.ts`)

Create test entities with automatic cleanup. Actors and Items are automatically sourced from system compendiums when available, falling back to minimal `type: 'base'` entities:

```typescript
import { createTestEntities, getEntityUuid } from '../helpers/testEntities';

// Create multiple entities (actors/items auto-sourced from compendiums)
await createTestEntities(version, [
  { key: 'primary', entityType: 'Actor', captureForDocs: true },
  { key: 'secondary', entityType: 'Actor' },
  { key: 'expendable', entityType: 'Actor' },  // For delete tests
], { capturedExamples });

// Later, retrieve the UUID
const actorUuid = getEntityUuid(version, 'Actor', 'primary');
```

Entity types supported:
- `Actor`, `Item`, `JournalEntry`, `Scene`, `Macro`, `RollTable`, `Playlist`

### API Request Helper (`apiRequest.ts`)

Make API requests with a configuration:

```typescript
import { ApiRequestConfig, makeRequest, replaceVariables } from '../helpers/apiRequest';
import { testVariables } from '../helpers/testVariables';

const requestConfig: ApiRequestConfig = {
  url: {
    raw: '{{baseUrl}}/get',
    host: ['{{baseUrl}}'],
    path: ['get'],
    query: [
      { key: 'clientId', value: '{{clientId}}' },
      { key: 'uuid', value: 'Actor.abc123' }
    ]
  },
  method: 'GET',
  header: [
    { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
  ]
};

const resolvedConfig = replaceVariables(requestConfig, testVariables);
const response = await makeRequest(resolvedConfig);
```

### Global Variables (`globalVariables.ts`)

Share state across test files using file-based persistence:

```typescript
import { setGlobalVariable, getGlobalVariable } from '../helpers/globalVariables';

// Store a value
setGlobalVariable('13', 'clientId', 'foundry-abc123');

// Retrieve it in another test file
const clientId = getGlobalVariable('13', 'clientId');
```

### Capturing Examples (`captureExample.ts`)

Capture API requests and responses for documentation:

```typescript
import { captureExample, saveExamples } from '../helpers/captureExample';

const captured = await captureExample(
  requestConfig,
  testVariables,
  '/my-endpoint - Description'
);

// Assertions on the captured response
expect(captured.response.status).toBe(200);

// Save examples for documentation generation
capturedExamples.push(captured);

afterAll(() => {
  saveExamples(capturedExamples, 'docs/examples/my-examples.json');
});
```

Captured examples include auto-generated code snippets in:
- JavaScript (fetch)
- TypeScript (axios)
- Python (requests)
- cURL
- Emojicode (for fun!)

## Writing New Tests

### Basic Test Structure

```typescript
/**
 * @file my-endpoints.test.ts
 * @description My Endpoint Tests
 * @endpoints GET /my-endpoint, POST /my-endpoint
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import * as path from 'path';

const capturedExamples: any[] = [];

describe('My Endpoint', () => {
  afterAll(() => {
    const outputPath = path.join(__dirname, '../../docs/examples/my-examples.json');
    saveExamples(capturedExamples, outputPath);
  });

  forEachVersion((version, getClientId) => {
    describe(`/my-endpoint (v${version})`, () => {
      test('GET /my-endpoint', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          // ... request configuration
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/my-endpoint'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
      });
    });
  });
});
```

### Adding to Test Order

Add your test file to `tests/helpers/testSequencer.ts` in the `TEST_ORDER` array at an appropriate position:

- Add **after** any tests that create data your tests depend on
- If your tests need canvas documents (tokens, lights, walls), add them **after** `scene-endpoints.test.ts` and **before** `scene-cleanup.test.ts`
- System-specific tests go in **Phase 7** (after permission filtering, before cleanup)
- Add **before** cleanup tests (`cleanup-entities.test.ts`, `scene-cleanup.test.ts`, `end-sessions.test.ts`)

:::caution Important
Tests not in `TEST_ORDER` won't run as part of `pnpm test`.
:::

### Adding a System-Specific Test

1. Create `tests/integration/mystem-endpoints.test.ts` using `forEachVersionWithSystem('mystem', ...)`
2. Add a system config to `tests/helpers/systemSetup.ts` in the `SYSTEM_CONFIGS` map
3. Add the test file to `TEST_ORDER` in Phase 7
4. The test will auto-skip on Foundry instances that don't have your system

## Documentation Generation

### Capturing Examples During Tests

Tests automatically capture API examples when using `captureExample()`. These are saved to `docs/examples/`.

### Generating API Documentation

```bash
# Generate API docs from route files
pnpm docs:generate

# Update docs with captured examples
pnpm docs:examples

# Full documentation update
pnpm docs:update

# Full documentation update and build (does all of the above + builds the site)
pnpm docs:full
```

The doc generator automatically expands parameterized routes (like `/canvas/:documentType`) into concrete endpoints (e.g., `/canvas/tokens`, `/canvas/walls`, etc.) using the `ROUTE_PARAM_EXPANSIONS` table in `generateApiDocs.js`.

:::caution Documentation Changes
Running tests captures real API responses, which updates `docs/examples/*.json` files. This regenerates code examples with current response data, potentially modifying documentation for endpoints you didn't touch.

**Before committing**, discard changes to documentation files for endpoints you're not working on:
```bash
# Review what changed
git diff docs/

# Discard changes to files you didn't intentionally modify
git checkout docs/examples/entity-examples.json  # example
```

Only commit documentation changes that are relevant to your PR.
:::

### Building Documentation Site

```bash
# Install docs dependencies
pnpm docs:install

# Build docs site
pnpm docs:build
```

## Troubleshooting

### Common Issues

**"No clientId found for version X"**
- Ensure session tests ran successfully
- Check that `USE_EXISTING_SESSION` is properly set
- Verify `TEST_CLIENT_ID_VX` environment variable

**"Invalid client ID"**
- The Foundry client may have disconnected
- Restart Foundry and re-run tests
- Check relay server logs for connection issues

**"Session creation failed"**
- Verify Foundry is running and world is loaded
- Check username/password credentials
- Ensure no popup dialogs are blocking

**"Request timed out" (408)**
- Increase timeout in test: `}, 30000);`
- Check network connectivity between relay and Foundry
- Verify the REST API module is active in Foundry
- For system-specific tests: ensure the Foundry module is rebuilt and the page is refreshed

**"User not found" on permission tests**
- The `TEST_PLAYER_USER_ID` must match a real Foundry user ID or username
- Check the `/players` endpoint to see available users

**Compendium data not loading**
- Session tests must run first (they store `systemId`)
- Check that the Foundry world has system compendium packs
- Look for `⚠ No compendiumPacks` warnings in test output

**Cookie isolation issues (multi-version testing)**
- Use different browser profiles for each Foundry version
- Or use headless mode (`USE_EXISTING_SESSION=false`)

### Debugging

:::tip Debugging with Examples
Check the `docs/examples/` folder to see the full requests, responses, and auto-generated code examples captured during tests. This is very helpful for debugging API issues.
:::

### Browser Console Logging

When debugging headless session issues, you can capture the browser's console output to log files by setting the `CAPTURE_BROWSER_CONSOLE` environment variable:

```bash
CAPTURE_BROWSER_CONSOLE=debug pnpm test
```

Valid levels (each includes levels above it):
- `error` — page errors and failed requests only
- `warn` — warnings + errors
- `debug` — all console output

Logs are written to `data/browser-logs/` with filenames like `{world}_{username}_{foundryVersion}_{timestamp}.log`. Previous logs for the same world/username are automatically cleaned up when a new session starts.

You can also pass `captureBrowserConsole` in the request body of `POST /start-session` to enable logging for individual sessions.

### Debugging Test Output

After running tests, check the generated outputs for debugging:

**Test Reports** (generated in `test-results/`):
- `test-report.html` - Visual HTML report of all test results
- `junit.xml` - JUnit format report for CI/CD integration

**Captured Examples** (generated in `docs/examples/`):
- JSON files containing full request/response data for each endpoint
- Auto-generated code snippets (JavaScript, Python, cURL, etc.)
- Useful for verifying API behavior and debugging failures

:::note
The `tests/.global-vars.json` file stores state during test execution but is automatically deleted after tests complete.
:::

## Best Practices

1. **Use semantic entity keys**: `primary`, `secondary`, `expendable` instead of `test1`, `test2`
2. **Always register for cleanup**: Use `createTestEntities()` to automatically track entities
3. **Capture examples for docs**: Use `captureExample()` for endpoints that should be documented
4. **Test all versions**: Use `forEachVersion()` to ensure compatibility
5. **Extended timeouts**: Use longer timeouts for complex operations: `}, 30000);`
6. **Descriptive assertions**: Test specific response properties, not just status codes
7. **Data-driven tests**: Use tables/loops for endpoints with multiple variants (see canvas tests)
8. **Graceful degradation**: System-specific tests should skip gracefully when data isn't available

## Generating Test Files

For new endpoints, you can auto-generate test file boilerplate:

```bash
pnpm test:generate
```

This script parses route files and generates integration test stubs for endpoints that don't have tests yet. It extracts:

- Endpoint paths and methods
- Required and optional parameters from `createApiRoute` configuration
- JSDoc descriptions for test documentation (from the `@route` tag)

### Manual Steps Required

Generated test files are **starting points** that need customization:

1. **Replace placeholder parameter values** with actual test data
2. **Use test helpers** for dynamic values:
   - `getEntityUuid(version, 'Actor', 'primary')` - Get UUIDs from test entities
   - `getGlobalVariable(version, 'key')` - Get values stored by other tests
   - `setGlobalVariable(version, 'key', value)` - Store values for other tests
3. **Add meaningful assertions** that verify the response data
4. **Add the test file to `testSequencer.ts`** in the appropriate phase

### Creating Additional Test Entities

If your tests need specific entities that don't exist, add them in `entity-endpoints.test.ts`:

```typescript
await createTestEntities(version, [
  { key: 'my-special-actor', entityType: 'Actor', data: { /* custom data */ } },
], { capturedExamples });
```

These entities will be automatically cleaned up after tests complete.

## Useful Resources

- **Foundry VTT API Documentation**: [https://foundryvtt.com/api/](https://foundryvtt.com/api/)
- **Project Discord**: [https://discord.gg/U634xNGRAC](https://discord.gg/U634xNGRAC)
- **Contributing Guide**: [Contributing Documentation](./contributing.md)
