/**
 * @file user-endpoints.test.ts
 * @description User Management Endpoint Tests
 * @endpoints GET /users, GET /user, POST /user, PUT /user, DELETE /user
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import { forEachVersion } from '../../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../../helpers/globalVariables';

import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('User', () => {
  afterAll(() => {
    const outputPath = path.join(__dirname, '../../../docs/examples/user-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/user (v${version})`, () => {
      test('GET /users - List all users', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/users',
            host: ['{{baseUrl}}'],
            path: ['users'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
          },
          method: 'GET',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ]
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/users - List all users'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeInstanceOf(Array);
        expect(captured.response.data.data.length).toBeGreaterThan(0);

        // Verify user objects have expected fields
        const firstUser = captured.response.data.data[0];
        expect(firstUser).toHaveProperty('id');
        expect(firstUser).toHaveProperty('name');
        expect(firstUser).toHaveProperty('role');
        expect(firstUser).toHaveProperty('isGM');
        expect(firstUser).toHaveProperty('active');
      });

      test('POST /user - Create a test user', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/user',
            host: ['{{baseUrl}}'],
            path: ['user'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
          },
          method: 'POST',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              name: 'test-api-user',
              role: 1,
              password: 'testpassword123'
            }, null, 2)
          }
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/user - Create a new user'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data.name).toBe('test-api-user');
        expect(captured.response.data.data.role).toBe(1);
        expect(captured.response.data.data.id).toBeTruthy();

        // Store the created user ID for subsequent tests
        const userId = captured.response.data.data.id;
        setGlobalVariable(version, 'testUserId', userId);
        console.log(`  ✓ Created test user: ${userId} (test-api-user)`);
      }, 30000);

      test('GET /user - Get user by ID', async () => {
        setVariable('clientId', getClientId());
        const userId = getGlobalVariable(version, 'testUserId');
        expect(userId).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/user',
            host: ['{{baseUrl}}'],
            path: ['user'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'id',
                value: userId,
              }
            ]
          },
          method: 'GET',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ]
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/user - Get user by ID'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data.id).toBe(userId);
        expect(captured.response.data.data.name).toBe('test-api-user');
        expect(captured.response.data.data.role).toBe(1);
      });

      test('GET /user - Get user by name', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/user',
            host: ['{{baseUrl}}'],
            path: ['user'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'name',
                value: 'test-api-user',
              }
            ]
          },
          method: 'GET',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ]
        };

        const resolvedConfig = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolvedConfig);

        expect(response.status).toBe(200);
        expect(response.data.data).toBeTruthy();
        expect(response.data.data.name).toBe('test-api-user');
      });

      test('PUT /user - Update user role', async () => {
        setVariable('clientId', getClientId());
        const userId = getGlobalVariable(version, 'testUserId');
        expect(userId).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/user',
            host: ['{{baseUrl}}'],
            path: ['user'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
          },
          method: 'PUT',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              id: userId,
              data: {
                role: 2
              }
            }, null, 2)
          }
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/user - Update user role'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data.role).toBe(2);
      }, 30000);

      test('DELETE /user - Delete test user', async () => {
        setVariable('clientId', getClientId());
        const userId = getGlobalVariable(version, 'testUserId');
        expect(userId).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/user',
            host: ['{{baseUrl}}'],
            path: ['user'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'id',
                value: userId,
              }
            ]
          },
          method: 'DELETE',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ]
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/user - Delete user'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.success).toBe(true);
        console.log(`  ✓ Deleted test user: ${userId}`);
      }, 30000);

      test('GET /users - Verify user was deleted', async () => {
        setVariable('clientId', getClientId());
        const userId = getGlobalVariable(version, 'testUserId');

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/users',
            host: ['{{baseUrl}}'],
            path: ['users'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
          },
          method: 'GET',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ]
        };

        const resolvedConfig = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolvedConfig);

        expect(response.status).toBe(200);
        expect(response.data.data).toBeInstanceOf(Array);

        // Verify the test user is no longer in the list
        const testUser = response.data.data.find((u: any) => u.id === userId);
        expect(testUser).toBeUndefined();
        console.log(`  ✓ Verified test user ${userId} no longer exists`);
      });
    });
  });
});
