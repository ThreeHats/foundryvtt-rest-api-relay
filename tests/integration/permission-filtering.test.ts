/**
 * @file permission-filtering.test.ts
 * @description Permission Filtering Tests
 *
 * Dynamically scans route files for endpoints that accept `userId` and tests:
 * 1. Invalid userId → 400 "User not found" on every endpoint (always runs)
 * 2. Player filtering → verifies non-GM sees scoped results (requires TEST_PLAYER_USER_ID)
 *
 * To enable player filtering tests:
 *   1. Create a non-GM player in your Foundry world (e.g., "Player1")
 *   2. Set TEST_PLAYER_USER_ID_V13=Player1 in .env.test
 */

import { describe, test, expect } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { forEachVersion } from '../helpers/multiVersion';
import { getEntityUuid } from '../helpers/testEntities';
import * as fs from 'fs';
import * as path from 'path';

// ──────────────────────────────────────────────
// Route scanning
// ──────────────────────────────────────────────

interface UserIdEndpoint {
  path: string;
  method: string;
  requiredParams: Array<{ name: string; from: string }>;
}

/**
 * Scan route files for endpoints that accept a userId parameter.
 * Also extracts required params so we can build valid requests.
 */
function extractUserIdEndpoints(): UserIdEndpoint[] {
  const endpoints: UserIdEndpoint[] = [];
  const apiDir = path.join(__dirname, '../../src/routes/api');
  const routeFiles = fs.readdirSync(apiDir).filter(f => f.endsWith('.ts'));

  // Build router mount prefix map from api.ts
  const apiTsPath = path.join(__dirname, '../../src/routes/api.ts');
  const apiTsContent = fs.readFileSync(apiTsPath, 'utf8');
  const routerMounts = new Map<string, string>();
  const mountRegex = /app\.use\(['"]([^'"]+)['"],\s*(\w+Router)\)/g;
  let m;
  while ((m = mountRegex.exec(apiTsContent)) !== null) {
    const fileName = m[2].replace('Router', '') + '.ts';
    if (m[1] !== '/') routerMounts.set(fileName, m[1]);
  }

  for (const file of routeFiles) {
    const content = fs.readFileSync(path.join(apiDir, file), 'utf8');
    const lines = content.split('\n');
    const prefix = routerMounts.get(file) || '';

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      const routeMatch = line.match(/\w+[Rr]outer\.(get|post|put|delete)\(\s*["']([^"']+)["']/);
      if (!routeMatch) continue;

      const [, method, routePath] = routeMatch;
      const fullPath = prefix + routePath;

      // Skip deprecated endpoints
      let deprecated = false;
      for (let k = Math.max(0, i - 10); k < Math.min(i + 5, lines.length); k++) {
        if (lines[k].includes('deprecated')) { deprecated = true; break; }
      }
      if (deprecated) continue;

      // Scan ahead for userId and required params
      let hasUserId = false;
      let inRequired = false;
      const requiredParams: Array<{ name: string; from: string }> = [];

      for (let j = i; j < Math.min(i + 100, lines.length); j++) {
        const scanLine = lines[j];
        if (scanLine.includes("'userId'") || scanLine.includes('"userId"')) hasUserId = true;
        if (scanLine.includes('requiredParams')) inRequired = true;
        if (scanLine.includes('optionalParams')) inRequired = false;

        if (inRequired) {
          const nameMatch = scanLine.match(/name:\s*['"]([^'"]+)['"]/);
          // Match both from: 'body' and from: ['body', 'query']
          const fromSingle = scanLine.match(/from:\s*['"]([^'"]+)['"]/);
          const fromArray = scanLine.match(/from:\s*\[['"]([^'"]+)['"]/);
          const from = fromSingle?.[1] || fromArray?.[1] || null;
          if (nameMatch && from) {
            requiredParams.push({ name: nameMatch[1], from });
          }
        }

        if (scanLine.includes('}));') && j > i + 5) break;
      }

      if (hasUserId) {
        endpoints.push({ path: fullPath, method: method.toUpperCase(), requiredParams });
      }
    }
  }

  return endpoints;
}

/** Provide a sensible dummy value for a required parameter. */
function dummyValue(name: string): any {
  switch (name) {
    case 'clientId': return '{{clientId}}';
    case 'uuid': return 'Actor.fakeId12345678';
    case 'query': return 'test';
    case 'formula': return '1d20';
    case 'entityType': return 'Actor';
    case 'content': return 'test';
    case 'actorUuid': return 'Actor.fakeId12345678';
    case 'name': return 'test';
    case 'folderId': return 'fakeFolder123';
    case 'folderType': return 'Actor';
    case 'attribute': return 'name';
    case 'amount': return 1;
    case 'sceneId': return 'fakeScene123';
    case 'script': return 'return 1';
    case 'data': return { name: 'test' };
    case 'details': return '["resources"]';
    default: return 'test-value';
  }
}

/**
 * Extra params needed for routes with custom validateParams that the scanner can't detect.
 * Keys are "METHOD /path", values are extra query/body params to add.
 */
const CUSTOM_VALIDATOR_PARAMS: Record<string, { query?: Record<string, string>; body?: Record<string, any> }> = {
  'GET /scene': { query: { all: 'true' } },
  'PUT /scene': { body: { sceneId: 'fakeScene123', data: { name: 'test' } } },
  'DELETE /scene': { query: { name: 'fakeScene' } },
  'POST /switch-scene': { body: { name: 'fakeScene' } },
  'POST /select': { body: { all: true } },
  'POST /execute-js': { body: { script: 'return 1' } },
  'GET /get-folder': { query: { name: 'test' } },
  'POST /create-folder': { body: { name: 'test', folderType: 'Actor' } },
  'DELETE /delete-folder': { body: { folderId: 'fakeFolderId' } },
};

/**
 * Build a request config for an endpoint with all required params satisfied,
 * plus the given userId.
 */
function buildRequest(endpoint: UserIdEndpoint, userId: string): ApiRequestConfig {
  // Handle path params like :documentType
  const testPath = endpoint.path.replace(/:(\w+)/g, (_m, name) => {
    if (name === 'documentType') return 'tokens';
    if (name === 'messageId') return 'fakeMsg123';
    if (name === 'path') return 'test';
    return 'dummy';
  });

  const queryParams: Array<{ key: string; value: string }> = [
    { key: 'userId', value: userId }
  ];
  const bodyData: Record<string, any> = {};

  for (const p of endpoint.requiredParams) {
    if (p.name === 'clientId') {
      queryParams.push({ key: 'clientId', value: '{{clientId}}' });
    } else if (p.from === 'query' || endpoint.method === 'GET') {
      // GET requests have no body — all params go in query
      queryParams.push({ key: p.name, value: String(dummyValue(p.name)) });
    } else {
      bodyData[p.name] = dummyValue(p.name);
    }
  }

  // Ensure clientId is always in query
  if (!queryParams.some(q => q.key === 'clientId')) {
    queryParams.push({ key: 'clientId', value: '{{clientId}}' });
  }

  // Apply custom validator overrides for routes with complex validation
  const endpointKey = `${endpoint.method} ${endpoint.path}`;
  const custom = CUSTOM_VALIDATOR_PARAMS[endpointKey];
  if (custom?.query) {
    for (const [k, v] of Object.entries(custom.query)) {
      if (!queryParams.some(q => q.key === k)) {
        queryParams.push({ key: k, value: v });
      }
    }
  }
  if (custom?.body) {
    Object.assign(bodyData, custom.body);
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

function getTestPlayerUserId(version: string): string | null {
  return process.env[`TEST_PLAYER_USER_ID_V${version}`] || null;
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

    const playerUserId = getTestPlayerUserId(version);

    if (playerUserId) {
      describe(`Player filtering (v${version})`, () => {

        // Verify the test player exists and is non-GM, and get the GM user ID
        let gmUserId: string;

        test('GET /players confirms test player is non-GM', async () => {
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

          const player = users.find(
            (u: any) => u.id === playerUserId || u.name?.toLowerCase() === playerUserId.toLowerCase()
          );
          expect(player).toBeTruthy();
          expect(player.isGM).toBe(false);

          const gm = users.find((u: any) => u.isGM);
          expect(gm).toBeTruthy();
          gmUserId = gm.id;

          console.log(`  ✓ Player: ${player.name} (${player.id}), GM: ${gm.name} (${gm.id})`);
        });

        // Player sees fewer search results than GM
        test('GET /search — player sees ≤ GM results', async () => {
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
          const playerResponse = await makeSearchRequest(playerUserId);

          expect(gmResponse.status).toBe(200);
          expect(playerResponse.status).toBe(200);

          const gmCount = gmResponse.data.results?.length ?? 0;
          const playerCount = playerResponse.data.results?.length ?? 0;
          console.log(`  ✓ Search: GM=${gmCount}, Player=${playerCount}`);
          expect(gmCount).toBeGreaterThanOrEqual(playerCount);
        });

        // Player cannot update GM-created entities
        test('PUT /update — player denied on GM-owned entity', async () => {
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
                { key: 'userId', value: playerUserId }
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
                { key: 'userId', value: playerUserId }
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
          setVariable('clientId', getClientId());

          // Step 1: Send a whisper to GM only (by user ID, not name)
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
                { key: 'userId', value: playerUserId }
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
          const playerResponse = await makeChatRequest(playerUserId);

          expect(gmResponse.status).toBe(200);
          expect(playerResponse.status).toBe(200);

          const gmTotal = gmResponse.data.data?.total ?? 0;
          const playerTotal = playerResponse.data.data?.total ?? 0;
          console.log(`  ✓ Chat messages — GM: ${gmTotal}, Player: ${playerTotal}`);
          expect(gmTotal).toBeGreaterThanOrEqual(playerTotal);
        }, 30000);
      });

    } else {
      describe(`Player filtering (v${version})`, () => {
        test.skip(`Skipped — set TEST_PLAYER_USER_ID_V${version} in .env.test`, () => {});
      });
    }

  });
});
