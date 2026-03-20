/**
 * @file scene-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings (incomplete)
 * @description Scene Endpoint Tests
 * @endpoints GET /scene, POST /scene, PUT /scene, DELETE /scene, POST /switch-scene
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../helpers/globalVariables';

import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Scene', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/scene-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/scene (v${version})`, () => {
      test('GET /scene - Record original active scene', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'active',
                value: 'true',
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

        // Store the original active scene ID so we can restore it later
        if (response.status === 200 && response.data?.data?._id) {
          setGlobalVariable(version, 'originalActiveSceneId', response.data.data._id);
          console.log(`  ✓ Recorded original active scene: ${response.data.data._id} (${response.data.data.name})`);
        } else {
          console.log(`  ⚠ No active scene found, will skip restore`);
        }
      });

      test('POST /scene - Create a scene', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
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
              data: {
                name: 'test-scene',
                width: 1000,
                height: 1000,
                grid: { size: 100 }
              }
            }, null, 2)
          }
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/scene - Create scene'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data.name).toBe('test-scene');
        expect(captured.response.data.data._id).toBeTruthy();

        // Store the scene ID for subsequent tests
        const sceneId = captured.response.data.data._id;
        setGlobalVariable(version, 'testSceneId', sceneId);

        // Note: NOT registered for entity cleanup — scene-cleanup.test.ts handles
        // restoring the original scene and deleting this one (must restore first
        // since you can't delete the active scene).
      }, 30000);

      test('POST /scene - Create expendable scene', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
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
              data: {
                name: 'test-scene-expendable',
                width: 500,
                height: 500
              }
            }, null, 2)
          }
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/scene - Create expendable scene'
        );

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data.name).toBe('test-scene-expendable');

        const sceneId = captured.response.data.data._id;
        setGlobalVariable(version, 'expendableSceneId', sceneId);
      }, 30000);

      test('GET /scene - Get all scenes', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'all',
                value: 'true',
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
          '/scene - Get all scenes'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeInstanceOf(Array);
        expect(captured.response.data.data.length).toBeGreaterThan(0);

        // Verify our test scene is in the list
        const testScene = captured.response.data.data.find((s: any) => s.name === 'test-scene');
        expect(testScene).toBeTruthy();
      });

      test('GET /scene - Get scene by name', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'name',
                value: 'test-scene',
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
          '/scene - Get scene by name'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data.name).toBe('test-scene');
      });

      test('GET /scene - Get scene by ID', async () => {
        setVariable('clientId', getClientId());

        const sceneId = getGlobalVariable(version, 'testSceneId');
        expect(sceneId).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'sceneId',
                value: sceneId,
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
          '/scene - Get scene by ID'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data._id).toBe(sceneId);
        expect(captured.response.data.data.name).toBe('test-scene');
      });

      test('PUT /scene - Update a scene', async () => {
        setVariable('clientId', getClientId());

        const sceneId = getGlobalVariable(version, 'testSceneId');
        expect(sceneId).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
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
              sceneId: sceneId,
              data: {
                name: 'test-scene-updated'
              }
            }, null, 2)
          }
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/scene - Update scene'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data.name).toBe('test-scene-updated');
      }, 30000);
    });

    describe(`/switch-scene (v${version})`, () => {
      test('POST /switch-scene - Activate test scene', async () => {
        setVariable('clientId', getClientId());

        const sceneId = getGlobalVariable(version, 'testSceneId');
        expect(sceneId).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/switch-scene',
            host: ['{{baseUrl}}'],
            path: ['switch-scene'],
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
              sceneId: sceneId
            }, null, 2)
          }
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/switch-scene - Activate scene'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data.data).toBeTruthy();

        // Wait for the scene to fully load before proceeding
        await new Promise(resolve => setTimeout(resolve, 5000));
      }, 30000);

      test('GET /scene - Verify active scene', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'active',
                value: 'true',
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
          '/scene - Get active scene'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeTruthy();

        const sceneId = getGlobalVariable(version, 'testSceneId');
        expect(captured.response.data.data._id).toBe(sceneId);
      });
    });

    describe(`/scene DELETE (v${version})`, () => {
      test('DELETE /scene - Delete expendable scene by ID', async () => {
        setVariable('clientId', getClientId());

        const sceneId = getGlobalVariable(version, 'expendableSceneId');
        expect(sceneId).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'sceneId',
                value: sceneId,
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
          '/scene - Delete scene'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.success).toBe(true);
      }, 30000);
    });

  });

});
