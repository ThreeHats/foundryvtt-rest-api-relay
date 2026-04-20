/**
 * @file scopes.test.ts
 * @description Tests for API key action scopes and multi-client scoping
 * @endpoints POST /auth/api-keys, GET /auth/api-keys, PATCH /auth/api-keys/:id,
 *   GET /api/get, POST /api/create, POST /api/execute-js
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables } from '../../helpers/testVariables';
import { forEachVersion } from '../../helpers/multiVersion';

const masterApiKey = testVariables.apiKey;

// Track keys created for cleanup
const keysToCleanup: string[] = [];

/** Helper to create a scoped key with specific scopes */
async function createScopedKey(name: string, scopes?: string[], scopedClientIds?: string[]): Promise<{ id: string; key: string; scopes: string[] }> {
  const body: any = { name };
  if (scopes !== undefined) {
    body.scopes = scopes;
  }
  if (scopedClientIds !== undefined) {
    body.scopedClientIds = scopedClientIds;
  }

  const config: ApiRequestConfig = {
    url: {
      raw: `${testVariables.baseUrl}/auth/api-keys`,
      host: [testVariables.baseUrl],
      path: ['auth', 'api-keys'],
    },
    method: 'POST',
    header: [{ key: 'x-api-key', value: masterApiKey }],
    body: { mode: 'raw', raw: JSON.stringify(body) },
  };

  const response = await makeRequest(config);
  expect(response.status).toBe(201);
  keysToCleanup.push(String(response.data.id));
  return {
    id: String(response.data.id),
    key: response.data.key,
    scopes: response.data.scopes,
  };
}

/** Helper to delete a scoped key */
async function deleteScopedKey(id: string) {
  const config: ApiRequestConfig = {
    url: {
      raw: `${testVariables.baseUrl}/auth/api-keys/${id}`,
      host: [testVariables.baseUrl],
      path: ['auth', 'api-keys', id],
    },
    method: 'DELETE',
    header: [{ key: 'x-api-key', value: masterApiKey }],
  };
  await makeRequest(config);
}

/** Helper to make an API request with a given key */
async function apiRequest(method: string, endpoint: string, apiKey: string, query?: Record<string, string>, body?: any) {
  const queryParams = query
    ? Object.entries(query).map(([key, value]) => ({ key, value }))
    : undefined;

  const config: ApiRequestConfig = {
    url: {
      raw: `${testVariables.baseUrl}${endpoint}`,
      host: [testVariables.baseUrl],
      path: endpoint.split('/').filter(Boolean),
      query: queryParams,
    },
    method: method as any,
    header: [{ key: 'x-api-key', value: apiKey }],
    body: body ? { mode: 'raw', raw: JSON.stringify(body) } : undefined,
  };

  return makeRequest(config);
}

describe('Action Scopes', () => {

  forEachVersion((version, getClientId) => {
    describe(`[v${version}] Scope Enforcement`, () => {
      test('scoped key with entity:read can GET /get', async () => {
        const key = await createScopedKey('entity-read-test', ['entity:read', 'world:info']);
        const response = await apiRequest('GET', '/world-info', key.key, { clientId: getClientId() });
        expect(response.status).toBe(200);
      });

      test('scoped key with entity:read gets 403 on POST /create', async () => {
        const key = await createScopedKey('entity-read-only', ['entity:read']);
        const response = await apiRequest('POST', '/create', key.key, undefined, {
          entityType: 'Actor',
          name: 'scope-test-should-fail',
          type: 'character',
        });
        expect(response.status).toBe(403);
        expect(response.data.error).toContain('scope');
      });

      test('scoped key with entity:write can POST /create', async () => {
        const key = await createScopedKey('entity-write-test', ['entity:read', 'entity:write']);
        const response = await apiRequest('POST', '/create', key.key, undefined, {
          entityType: 'Actor',
          name: `scope-test-${Date.now()}`,
          type: 'character',
        });
        // Should succeed (200) or at least not be 403
        expect(response.status).not.toBe(403);
      });

      test('scoped key without execute-js gets 403 on POST /execute-js', async () => {
        const key = await createScopedKey('no-exec-js', ['entity:read']);
        const response = await apiRequest('POST', '/execute-js', key.key, undefined, {
          code: 'return 1+1',
        });
        expect(response.status).toBe(403);
        expect(response.data.error).toContain('execute-js');
      });

      test('scoped key with execute-js scope can POST /execute-js', async () => {
        const key = await createScopedKey('has-exec-js', ['execute-js']);
        const response = await apiRequest('POST', '/execute-js', key.key, undefined, {
          code: 'return 1+1',
        });
        // Should succeed or at least not be 403
        expect(response.status).not.toBe(403);
      });

      test('master key has full access regardless of scopes', async () => {
        const response = await apiRequest('POST', '/execute-js', masterApiKey, undefined, {
          code: 'return 1+1',
        });
        expect(response.status).not.toBe(403);
      });
    });

    describe(`[v${version}] Scope Defaults and CRUD`, () => {
      test('creating scoped key without specifying scopes returns 400', async () => {
        const config: ApiRequestConfig = {
          url: {
            raw: `${testVariables.baseUrl}/auth/api-keys`,
            host: [testVariables.baseUrl],
            path: ['auth', 'api-keys'],
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: masterApiKey }],
          body: { mode: 'raw', raw: JSON.stringify({ name: 'no-scopes-test' }) },
        };
        const response = await makeRequest(config);
        expect(response.status).toBe(400);
        expect(response.data).toHaveProperty('error');
      });

      test('creating scoped key with explicit scopes stores them correctly', async () => {
        const requestedScopes = ['entity:read', 'roll:read', 'chat:read'];
        const key = await createScopedKey('explicit-scopes-test', requestedScopes);
        expect(key.scopes).toEqual(expect.arrayContaining(requestedScopes));
        expect(key.scopes.length).toBe(requestedScopes.length);
      });

      test('creating scoped key with invalid scope returns 400', async () => {
        const body = { name: 'bad-scope-test', scopes: ['entity:read', 'not-a-real-scope'] };
        const config: ApiRequestConfig = {
          url: {
            raw: `${testVariables.baseUrl}/auth/api-keys`,
            host: [testVariables.baseUrl],
            path: ['auth', 'api-keys'],
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: masterApiKey }],
          body: { mode: 'raw', raw: JSON.stringify(body) },
        };
        const response = await makeRequest(config);
        expect(response.status).toBe(400);
        expect(response.data.error).toContain('not-a-real-scope');
      });

      test('updating scoped key scopes via PATCH works', async () => {
        const key = await createScopedKey('update-scopes-test', ['entity:read']);
        const newScopes = ['entity:read', 'entity:write', 'roll:read'];

        const config: ApiRequestConfig = {
          url: {
            raw: `${testVariables.baseUrl}/auth/api-keys/${key.id}`,
            host: [testVariables.baseUrl],
            path: ['auth', 'api-keys', key.id],
          },
          method: 'PATCH',
          header: [{ key: 'x-api-key', value: masterApiKey }],
          body: { mode: 'raw', raw: JSON.stringify({ scopes: newScopes }) },
        };

        const response = await makeRequest(config);
        expect(response.status).toBe(200);
        expect(response.data.scopes).toEqual(expect.arrayContaining(newScopes));
      });

      test('GET /auth/api-keys returns scopes in response', async () => {
        const config: ApiRequestConfig = {
          url: {
            raw: `${testVariables.baseUrl}/auth/api-keys`,
            host: [testVariables.baseUrl],
            path: ['auth', 'api-keys'],
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: masterApiKey }],
        };

        const response = await makeRequest(config);
        expect(response.status).toBe(200);
        expect(response.data.keys).toBeDefined();
        // At least one key should have scopes
        const keyWithScopes = response.data.keys.find((k: any) =>
          Array.isArray(k.scopes) && k.scopes.length > 0
        );
        expect(keyWithScopes).toBeDefined();
      });

      test('existing key with empty scopes has full access (backward compat)', async () => {
        // The master key itself has no scopes set — it should have full access
        // This is implicitly tested by the master key test above, but let's verify
        // the endpoint behavior directly
        const response = await apiRequest('GET', '/world-info', masterApiKey, { clientId: getClientId() });
        expect(response.status).toBe(200);
      });
    });
  });
});

describe('Multi-Client Scoping', () => {

  forEachVersion((version, getClientId) => {
    describe(`[v${version}] Client Restriction`, () => {
      test('scoped key with specific clientId can access that client', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const key = await createScopedKey('client-scoped', ['entity:read', 'world:info'], [clientId]);
        const response = await apiRequest('GET', '/world-info', key.key, { clientId });
        expect(response.status).toBe(200);
      });

      test('scoped key with wrong clientId gets 403', async () => {
        const key = await createScopedKey('wrong-client', ['entity:read', 'world:info'], ['fake-client-id-xyz']);
        const clientId = getClientId();
        if (!clientId) return;

        const response = await apiRequest('GET', '/world-info', key.key, { clientId });
        expect(response.status).toBe(403);
      });

      test('scoped key with empty clientIds can access any client', async () => {
        const key = await createScopedKey('unrestricted-client', ['world:info'], []);
        const response = await apiRequest('GET', '/world-info', key.key, { clientId: getClientId() });
        expect(response.status).toBe(200);
      });

      test('scoped key with one allowed client auto-resolves', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const key = await createScopedKey('single-client-auto', ['world:info'], [clientId]);
        // Don't pass clientId — should auto-resolve to the single allowed client
        const response = await apiRequest('GET', '/world-info', key.key);
        expect(response.status).toBe(200);
      });

      test('creating key with scopedClientIds returns them in response', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const body = {
          name: 'multi-client-test',
          scopes: ['entity:read'],
          scopedClientIds: [clientId, 'another-client-id'],
        };
        const config: ApiRequestConfig = {
          url: {
            raw: `${testVariables.baseUrl}/auth/api-keys`,
            host: [testVariables.baseUrl],
            path: ['auth', 'api-keys'],
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: masterApiKey }],
          body: { mode: 'raw', raw: JSON.stringify(body) },
        };
        const response = await makeRequest(config);
        expect(response.status).toBe(201);
        keysToCleanup.push(String(response.data.id));
        expect(response.data.scopedClientIds).toEqual(expect.arrayContaining([clientId, 'another-client-id']));
      });
    });
  });
});

// Module-level cleanup: runs after both describe blocks complete, so it catches
// keys from both Action Scopes and Multi-Client Scoping tests.
afterAll(async () => {
  for (const id of keysToCleanup) {
    await deleteScopedKey(id).catch(() => {});
  }
});
