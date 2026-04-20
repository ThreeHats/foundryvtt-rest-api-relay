/**
 * @file scoped-keys-endpoints.test.ts
 * @description Scoped API key CRUD and auto-routing tests
 * @endpoints POST /auth/api-keys, GET /auth/api-keys, PATCH /auth/api-keys/:id,
 *   DELETE /auth/api-keys/:id
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { captureExample, appendExamples } from '../../helpers/captureExample';
import { forEachVersion } from '../../helpers/multiVersion';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

// Track scoped key state across tests
let scopedKeyId = '';
let scopedKeyToken = '';
let routingKeyId = '';
let routingKeyToken = '';

// Safety-net: track all created key IDs so we can delete any that the test flow missed
const createdKeyIds: string[] = [];

const masterApiKey = testVariables.apiKey;

describe('Scoped API Keys', () => {
  afterAll(async () => {
    const outputPath = path.join(__dirname, '../../../docs/examples/auth-examples.json');
    appendExamples(capturedExamples, outputPath);
    console.log(`\nAppended ${capturedExamples.length} scoped-key examples to ${outputPath}`);

    // Best-effort cleanup: delete any keys not already removed by the test flow
    for (const id of createdKeyIds) {
      if (id) {
        await makeRequest({
          url: { raw: `${testVariables.baseUrl}/auth/api-keys/${id}`, host: [testVariables.baseUrl], path: ['auth', 'api-keys', id] },
          method: 'DELETE',
          header: [{ key: 'x-api-key', value: masterApiKey }],
        }).catch(() => {});
      }
    }
  });

  // ==================== Group 1: Scoped Key CRUD ====================

  describe('CRUD Operations', () => {
    test('POST /auth/api-keys - create a scoped key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            name: 'Test Scoped Key',
            scopes: ['entity:read', 'structure:read'],
            monthlyLimit: '500'
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(201);
      expect(captured.response.data).toHaveProperty('id');
      expect(captured.response.data).toHaveProperty('key');
      expect(captured.response.data).toHaveProperty('name', 'Test Scoped Key');
      expect(captured.response.data).toHaveProperty('enabled', true);

      scopedKeyId = String(captured.response.data.id);
      scopedKeyToken = captured.response.data.key;
      createdKeyIds.push(scopedKeyId);
    });

    test('POST /auth/api-keys - should reject missing name', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({})
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys (missing name)');
      expect(captured.response.status).toBe(400);
      expect(captured.response.data).toHaveProperty('error');
    });

    test('GET /auth/api-keys - list scoped keys', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys'],
        },
        method: 'GET',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('keys');
      expect(Array.isArray(captured.response.data.keys)).toBe(true);
      expect(captured.response.data.keys.length).toBeGreaterThanOrEqual(1);

      // Verify key is masked in list
      const key = captured.response.data.keys.find((k: any) => String(k.id) === scopedKeyId);
      expect(key).toBeTruthy();
      expect(key.key).toMatch(/^.{8}\.\.\.$/);
      expect(key).toHaveProperty('name', 'Test Scoped Key');
    });

    test('PATCH /auth/api-keys/:id - update scoped key', async () => {
      expect(scopedKeyId).toBeTruthy();

      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys/${scopedKeyId}`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys', scopedKeyId],
        },
        method: 'PATCH',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            name: 'Updated Scoped Key',
            monthlyLimit: 1000
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys/:id');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('name', 'Updated Scoped Key');
    });

    test('PATCH /auth/api-keys/999999 - should return 404 for non-existent key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys/999999`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys', '999999'],
        },
        method: 'PATCH',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({ name: 'Does Not Exist' })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys/:id (not found)');
      expect(captured.response.status).toBe(404);
    });

    test('POST /auth/api-keys - scoped key cannot create other keys', async () => {
      expect(scopedKeyToken).toBeTruthy();

      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: scopedKeyToken }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({ name: 'Should Fail' })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys (scoped key forbidden)');
      expect(captured.response.status).toBe(403);
      expect(captured.response.data).toHaveProperty('error');
    });

    test('GET /auth/api-keys - scoped key cannot list keys', async () => {
      expect(scopedKeyToken).toBeTruthy();

      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys'],
        },
        method: 'GET',
        header: [
          { key: 'x-api-key', value: scopedKeyToken }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys (scoped key forbidden)');
      expect(captured.response.status).toBe(403);
    });

    test('DELETE /auth/api-keys/:id - delete scoped key', async () => {
      expect(scopedKeyId).toBeTruthy();

      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys/${scopedKeyId}`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys', scopedKeyId],
        },
        method: 'DELETE',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys/:id');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('success', true);
    });

    test('DELETE /auth/api-keys/999999 - should return 404 for non-existent key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/api-keys/999999`,
          host: [testVariables.baseUrl],
          path: ['auth', 'api-keys', '999999'],
        },
        method: 'DELETE',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/auth/api-keys/:id (not found)');
      expect(captured.response.status).toBe(404);
    });
  });

  // ==================== Group 2: Scoped Key Auto-Routing ====================

  describe('Auto-Routing', () => {
    forEachVersion((version, getClientId) => {
      describe(`Scoped key auto-routing (v${version})`, () => {
        test('Create scoped key with scopedClientId', async () => {
          const clientId = getClientId();
          expect(clientId).toBeTruthy();

          const requestConfig: ApiRequestConfig = {
            url: {
              raw: `${testVariables.baseUrl}/auth/api-keys`,
              host: [testVariables.baseUrl],
              path: ['auth', 'api-keys'],
            },
            method: 'POST',
            header: [
              { key: 'x-api-key', value: masterApiKey }
            ],
            body: {
              mode: 'raw',
              raw: JSON.stringify({
                name: `Auto-Route Test v${version}`,
                scopes: ['structure:read'],
                scopedClientId: clientId
              })
            }
          };

          const captured = await captureExample(requestConfig, {}, '/auth/api-keys (with scopedClientId)');

          expect(captured.response.status).toBe(201);
          expect(captured.response.data).toHaveProperty('key');
          expect(captured.response.data).toHaveProperty('id');

          routingKeyId = String(captured.response.data.id);
          routingKeyToken = captured.response.data.key;
          createdKeyIds.push(routingKeyId);
        });

        test('GET /structure - auto-routes via scoped key (no clientId param)', async () => {
          expect(routingKeyToken).toBeTruthy();

          // Use scoped key WITHOUT clientId query param — should auto-route
          const requestConfig: ApiRequestConfig = {
            url: {
              raw: `${testVariables.baseUrl}/structure`,
              host: [testVariables.baseUrl],
              path: ['structure'],
            },
            method: 'GET',
            header: [
              { key: 'x-api-key', value: routingKeyToken }
            ]
          };

          const captured = await captureExample(requestConfig, {}, '/structure (scoped key auto-route)');
          capturedExamples.push(captured);

          // The scoped key's scopedClientId should auto-resolve
          expect(captured.response.status).toBe(200);
          expect(captured.response.data).toBeTruthy();
        });

        test('Delete auto-routing scoped key', async () => {
          expect(routingKeyId).toBeTruthy();

          const requestConfig: ApiRequestConfig = {
            url: {
              raw: `${testVariables.baseUrl}/auth/api-keys/${routingKeyId}`,
              host: [testVariables.baseUrl],
              path: ['auth', 'api-keys', routingKeyId],
            },
            method: 'DELETE',
            header: [
              { key: 'x-api-key', value: masterApiKey }
            ]
          };

          const captured = await captureExample(requestConfig, {}, '/auth/api-keys/:id (cleanup routing key)');
          expect(captured.response.status).toBe(200);
        });
      });
    });
  });
});
