/**
 * @file permission-filtering.test.ts
 * @description Permission Filtering Tests
 *
 * Uses api-docs.json to discover endpoints that accept `userId` and tests:
 * 1. Invalid userId → 400 "User not found" on every endpoint (always runs)
 * 2. Player filtering → verifies non-GM sees scoped results (requires TEST_PLAYER_USER_ID)
 *
 * When new endpoints are added, re-run `cd go-relay && go run ./cmd/docgen`
 * and these tests automatically cover them.
 *
 * To enable player filtering tests:
 *   1. Create a non-GM player in your Foundry world (e.g., "Player1")
 *   2. Set TEST_PLAYER_USER_ID_V13=Player1 in .env.test
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { forEachVersion } from '../../helpers/multiVersion';
import { getEntityUuid } from '../../helpers/testEntities';
import {
  loadEndpoints, dummyParamValue, resolvePathParam,
  isAuthRoute, shouldSkipInvalidUserId, getCustomRequiredParams
} from '../../helpers/endpointMetadata';

// ──────────────────────────────────────────────
// Endpoint discovery from api-docs.json
// ──────────────────────────────────────────────

interface UserIdEndpoint {
  path: string;
  method: string;
  requiredParams: Array<{ name: string; location: string; type: string }>;
}

/**
 * Discover endpoints that accept userId from api-docs.json.
 */
function extractUserIdEndpoints(): UserIdEndpoint[] {
  const allEndpoints = loadEndpoints();
  return allEndpoints
    .filter(ep => ep.hasUserId && !ep.isSSE && !isAuthRoute(ep) && !shouldSkipInvalidUserId(ep))
    .map(ep => ({
      path: ep.path,
      method: ep.method,
      requiredParams: ep.requiredParams.map(p => ({
        name: p.name,
        location: p.location,
        type: p.type,
      })),
    }));
}

/** Provide a sensible dummy value for a required parameter. */
function dummyValue(name: string, type: string = 'string'): any {
  return dummyParamValue(name, type);
}

// Custom validator params are centralized in endpointMetadata.ts

/**
 * Build a request config for an endpoint with all required params satisfied,
 * plus the given userId.
 */
function buildRequest(endpoint: UserIdEndpoint, userId: string): ApiRequestConfig {
  // Handle path params like {documentType} (Chi style)
  const testPath = endpoint.path.replace(/\{(\w+)\}/g, (_m, name) => resolvePathParam(name));

  const queryParams: Array<{ key: string; value: string }> = [
    { key: 'userId', value: userId }
  ];
  const bodyData: Record<string, any> = {};

  for (const p of endpoint.requiredParams) {
    if (p.name === 'clientId') {
      queryParams.push({ key: 'clientId', value: '{{clientId}}' });
    } else if (p.location === 'query' || endpoint.method === 'GET') {
      // GET requests have no body — all params go in query
      const val = dummyValue(p.name, p.type);
      queryParams.push({ key: p.name, value: typeof val === 'object' ? JSON.stringify(val) : String(val) });
    } else {
      bodyData[p.name] = dummyValue(p.name, p.type);
    }
  }

  // Ensure clientId is always in query
  if (!queryParams.some(q => q.key === 'clientId')) {
    queryParams.push({ key: 'clientId', value: '{{clientId}}' });
  }

  // Apply custom required params for endpoints with custom validation.
  // These go in query for GET/DELETE (DELETE handlers typically read from query),
  // and in body for POST/PUT.
  const customParams = getCustomRequiredParams(endpoint.method, endpoint.path);
  for (const [k, v] of Object.entries(customParams)) {
    if (endpoint.method === 'GET' || endpoint.method === 'DELETE') {
      if (!queryParams.some(q => q.key === k)) {
        queryParams.push({ key: k, value: typeof v === 'object' ? JSON.stringify(v) : String(v) });
      }
    } else {
      if (!(k in bodyData)) {
        bodyData[k] = v;
      }
    }
  }

  const config: ApiRequestConfig = {
    url: {
      raw: `{{baseUrl}}${testPath}`,
      host: ['{{baseUrl}}'],
      path: [testPath.replace(/^\//, '')],
      query: queryParams
    },
    method: endpoint.method as any,
    header: [
      { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
      { key: 'Content-Type', value: 'application/json', type: 'text' }
    ]
  };

  if (['POST', 'PUT', 'DELETE'].includes(endpoint.method)) {
    config.body = { mode: 'raw', raw: JSON.stringify(bodyData) };
  }

  return config;
}

// ──────────────────────────────────────────────
// Tests
// ──────────────────────────────────────────────

const userIdEndpoints = extractUserIdEndpoints();

// Endpoints to skip for the invalid-userId sweep:
// (none currently — all userId-accepting endpoints should validate)
const SKIP_INVALID_USERID: string[] = [];

describe('Permission Filtering', () => {

  forEachVersion((version, getClientId) => {

    // ═══════════════════════════════════════════
    // Invalid userId — all endpoints that accept userId
    // ═══════════════════════════════════════════

    describe(`Invalid userId (v${version})`, () => {
      const testableEndpoints = userIdEndpoints.filter(
        e => !SKIP_INVALID_USERID.includes(e.path)
      );

      testableEndpoints.forEach(endpoint => {
        test(`${endpoint.method} ${endpoint.path} - rejects invalid userId`, async () => {
          setVariable('clientId', getClientId());

          const config = buildRequest(endpoint, 'nonexistent-user-xyz-99999');
          const resolved = replaceVariables(config, testVariables);
          const response = await makeRequest(resolved);

          expect(response.status).toBe(400);
          expect(response.data).toHaveProperty('error');
          expect(response.data.error).toContain('User not found');
        });
      });
    });

    // ═══════════════════════════════════════════
    // Player filtering — conditional on env config
    // ═══════════════════════════════════════════

    describe(`Player filtering (v${version})`, () => {
      let gmUserId: string = '';
      let foundPlayerUserId: string = '';
      const testPlayerName = `permission-test-player-v${version}`;

      beforeAll(async () => {
        const clientId = getClientId();
        if (!clientId) return;
        setVariable('clientId', clientId);

        // Create ephemeral non-GM player for permission tests
        const createConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/user',
            host: ['{{baseUrl}}'],
            path: ['user'],
            query: [{ key: 'clientId', value: '{{clientId}}' }],
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}' }],
          body: { mode: 'raw', raw: JSON.stringify({ name: testPlayerName, role: 1 }) },
        };
        const createResp = await makeRequest(replaceVariables(createConfig, testVariables));
        if (createResp.status === 200) {
          foundPlayerUserId = createResp.data?.data?.id ?? '';
        }

        // Discover the GM user ID from /players
        const playersConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/players',
            host: ['{{baseUrl}}'],
            path: ['players'],
            query: [{ key: 'clientId', value: '{{clientId}}' }],
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}' }],
        };
        const playersResp = await makeRequest(replaceVariables(playersConfig, testVariables));
        if (playersResp.status === 200) {
          const gm = (playersResp.data?.users ?? []).find((u: any) => u.isGM);
          if (gm) gmUserId = gm.id;
        }
      });

      afterAll(async () => {
        if (!foundPlayerUserId) return;
        const clientId = getClientId();
        if (!clientId) return;
        setVariable('clientId', clientId);
        const deleteConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/user',
            host: ['{{baseUrl}}'],
            path: ['user'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'id', value: foundPlayerUserId },
            ],
          },
          method: 'DELETE',
          header: [{ key: 'x-api-key', value: '{{apiKey}}' }],
        };
        await makeRequest(replaceVariables(deleteConfig, testVariables)).catch(() => {});
      });

      // Verify the created player exists and is non-GM
      test('GET /players confirms test player is non-GM', async () => {
        if (!foundPlayerUserId) {
          console.log('  Skipping: test player creation failed in beforeAll');
          return;
        }
        setVariable('clientId', getClientId());

        const config: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/players',
            host: ['{{baseUrl}}'],
            path: ['players'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };
        const resolved = replaceVariables(config, testVariables);
        const response = await makeRequest(resolved);

        expect(response.status).toBe(200);
        const users = response.data.users;
        expect(users).toBeInstanceOf(Array);

        const player = users.find((u: any) => u.id === foundPlayerUserId);
        expect(player).toBeTruthy();
        expect(player.isGM).toBe(false);

        expect(gmUserId).toBeTruthy();
        console.log(`  ✓ Player: ${player.name} (${player.id}), GM ID: ${gmUserId}`);
      });

      // Player sees fewer search results than GM
      test('GET /search — player sees ≤ GM results', async () => {
        if (!gmUserId || !foundPlayerUserId) {
          console.log('  Skipping: setup not complete');
          return;
        }
        setVariable('clientId', getClientId());

        const makeSearchRequest = (userId: string) => {
          const config: ApiRequestConfig = {
            url: {
              raw: '{{baseUrl}}/search',
              host: ['{{baseUrl}}'],
              path: ['search'],
              query: [
                { key: 'clientId', value: '{{clientId}}' },
                { key: 'query', value: 'test' },
                { key: 'userId', value: userId }
              ]
            },
            method: 'GET',
            header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
          };
          return makeRequest(replaceVariables(config, testVariables));
        };

        // Sequential to avoid overwhelming the WebSocket connection
        const gmResponse = await makeSearchRequest(gmUserId);
        const playerResponse = await makeSearchRequest(foundPlayerUserId);

        expect(gmResponse.status).toBe(200);
        expect(playerResponse.status).toBe(200);

        const gmCount = gmResponse.data.results?.length ?? 0;
        const playerCount = playerResponse.data.results?.length ?? 0;
        console.log(`  ✓ Search: GM=${gmCount}, Player=${playerCount}`);
        expect(gmCount).toBeGreaterThanOrEqual(playerCount);
      });

      // Player cannot update GM-created entities
      test('PUT /update — player denied on GM-owned entity', async () => {
        if (!foundPlayerUserId) {
          console.log('  Skipping: setup not complete');
          return;
        }
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const config: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/update',
            host: ['{{baseUrl}}'],
            path: ['update'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'uuid', value: actorUuid! },
              { key: 'userId', value: foundPlayerUserId }
            ]
          },
          method: 'PUT',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          body: { mode: 'raw', raw: JSON.stringify({ data: { name: 'Should Not Update' } }) }
        };
        const response = await makeRequest(replaceVariables(config, testVariables));

        expect(response.status).toBe(400);
        expect(response.data.error).toContain('permission');
        console.log(`  ✓ Update denied: ${response.data.error}`);
      });

      // Player cannot delete GM-created entities
      test('DELETE /delete — player denied on GM-owned entity', async () => {
        if (!foundPlayerUserId) {
          console.log('  Skipping: setup not complete');
          return;
        }
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'secondary');
        expect(actorUuid).toBeTruthy();

        const config: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/delete',
            host: ['{{baseUrl}}'],
            path: ['delete'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'uuid', value: actorUuid! },
              { key: 'userId', value: foundPlayerUserId }
            ]
          },
          method: 'DELETE',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };
        const response = await makeRequest(replaceVariables(config, testVariables));

        expect(response.status).toBe(400);
        expect(response.data.error).toContain('permission');
        console.log(`  ✓ Delete denied: ${response.data.error}`);
      });

      // ── Chat whisper visibility ──

      test('Whispered message not visible to non-recipient player', async () => {
        if (!gmUserId || !foundPlayerUserId) {
          console.log('  Skipping: setup not complete');
          return;
        }
        setVariable('clientId', getClientId());

        // Step 1: Send a whisper to GM only (by user ID)
        const sendConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              content: 'Secret GM whisper - player should not see this',
              whisper: [gmUserId]
            })
          }
        };
        const sendResponse = await makeRequest(replaceVariables(sendConfig, testVariables));
        expect(sendResponse.status).toBe(200);
        const whisperId = sendResponse.data.data?.id;
        expect(whisperId).toBeTruthy();

        // Step 2: Player should not see the whisper
        const playerChatConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'limit', value: '50' },
              { key: 'userId', value: foundPlayerUserId }
            ]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };
        const playerChat = await makeRequest(replaceVariables(playerChatConfig, testVariables));
        expect(playerChat.status).toBe(200);

        const playerMessages = playerChat.data.data?.messages || [];
        const playerSeesWhisper = playerMessages.some((m: any) => m.id === whisperId);
        expect(playerSeesWhisper).toBe(false);
        console.log(`  ✓ Whisper hidden from player (${playerMessages.length} messages visible)`);

        // Step 3: GM should see the whisper
        const gmChatConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'limit', value: '50' },
              { key: 'userId', value: gmUserId }
            ]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };
        const gmChat = await makeRequest(replaceVariables(gmChatConfig, testVariables));
        expect(gmChat.status).toBe(200);

        const gmMessages = gmChat.data.data?.messages || [];
        expect(gmMessages.some((m: any) => m.id === whisperId)).toBe(true);
        console.log(`  ✓ Whisper visible to GM (${gmMessages.length} messages visible)`);

        // Cleanup
        try {
          const del: ApiRequestConfig = {
            url: {
              raw: `{{baseUrl}}/chat/${whisperId}`,
              host: ['{{baseUrl}}'],
              path: ['chat', whisperId],
              query: [{ key: 'clientId', value: '{{clientId}}' }]
            },
            method: 'DELETE',
            header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
          };
          await makeRequest(replaceVariables(del, testVariables));
        } catch { /* best effort */ }
      }, 15000);

      // ── Player sees fewer chat messages than GM ──

      test('GET /chat — player sees ≤ GM messages', async () => {
        if (!gmUserId || !foundPlayerUserId) {
          console.log('  Skipping: setup not complete');
          return;
        }
        setVariable('clientId', getClientId());

        const makeChatRequest = (userId: string) => {
          const config: ApiRequestConfig = {
            url: {
              raw: '{{baseUrl}}/chat',
              host: ['{{baseUrl}}'],
              path: ['chat'],
              query: [
                { key: 'clientId', value: '{{clientId}}' },
                { key: 'userId', value: userId }
              ]
            },
            method: 'GET',
            header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
          };
          return makeRequest(replaceVariables(config, testVariables));
        };

        // Sequential to avoid overwhelming the WebSocket connection
        const gmResponse = await makeChatRequest(gmUserId);
        const playerResponse = await makeChatRequest(foundPlayerUserId);

        expect(gmResponse.status).toBe(200);
        expect(playerResponse.status).toBe(200);

        const gmTotal = gmResponse.data.data?.total ?? 0;
        const playerTotal = playerResponse.data.data?.total ?? 0;
        console.log(`  ✓ Chat messages — GM: ${gmTotal}, Player: ${playerTotal}`);
        expect(gmTotal).toBeGreaterThanOrEqual(playerTotal);
      }, 30000);
    });

  });
});
