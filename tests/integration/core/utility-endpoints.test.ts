/**
 * @file utility-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Utility and Canvas Interaction Endpoint Tests
 * @endpoints POST /select, GET /selected, POST /execute-js, GET /players,
 *            GET /world-info, POST /move-token, GET /measure-distance
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import { forEachVersion } from '../../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../../helpers/globalVariables';
import { getEntityUuid } from '../../helpers/testEntities';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Utility', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../../docs/examples/utility-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/select (v${version})`, () => {
      test('POST /select', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/select',
            host: ['{{baseUrl}}'],
            path: ['select'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
                description: 'Client ID for the Foundry world'
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
              all: true,
              overwrite: true
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/select'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data).toHaveProperty('selected');
        expect(captured.response.data.selected).toBeInstanceOf(Array);
        expect(captured.response.data.selected.length).toBeGreaterThan(0);
        expect(captured.response.data).toHaveProperty('count');
      });
    });

    describe(`/selected (v${version})`, () => {
      test('GET /selected', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/selected',
            host: ['{{baseUrl}}'],
            path: ['selected'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
                description: 'Client ID for the Foundry world'
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

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/selected'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data).toHaveProperty('selected');
        expect(captured.response.data.selected).toBeInstanceOf(Array);
        expect(captured.response.data.selected.length).toBeGreaterThan(0);
        expect(captured.response.data.selected[0]).toHaveProperty('tokenUuid');
        expect(captured.response.data.selected[0]).toHaveProperty('actorUuid');

        // Set global variables for future tests
        setGlobalVariable(version, 'selectedTokenUuid', captured.response.data.selected[0].tokenUuid);
        setGlobalVariable(version, 'selectedActorUuid', captured.response.data.selected[0].actorUuid);
      });
    });

    describe(`/players (v${version})`, () => {
      test('GET /players', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/players',
            host: ['{{baseUrl}}'],
            path: ['players'],
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
          '/players'
        );
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('users');
        expect(captured.response.data.users).toBeInstanceOf(Array);
        expect(captured.response.data.users.length).toBeGreaterThan(0);

        // Verify user structure
        const user = captured.response.data.users[0];
        expect(user).toHaveProperty('id');
        expect(user).toHaveProperty('name');
        expect(user).toHaveProperty('role');
        expect(user).toHaveProperty('isGM');
        expect(user).toHaveProperty('active');
      });
    });

    describe(`/execute-js (v${version})`, () => {
      test('POST /execute-js', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/execute-js',
            host: ['{{baseUrl}}'],
            path: ['execute-js'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
                description: 'Client ID for the Foundry world'
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
              script: 'const wsRelayUrl=game.settings.get(\"foundry-rest-api\", \"wsRelayUrl\");return wsRelayUrl;'
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/execute-js'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data).toHaveProperty('result');
        expect(typeof captured.response.data.result).toBe('string');
      });
    });

  });

  forEachVersion((version, getClientId) => {
    describe(`GET /world-info (v${version})`, () => {
      test('GET /world-info - get comprehensive world information', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/world-info',
            host: ['{{baseUrl}}'],
            path: ['world-info'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };

        const captured = await captureExample(requestConfig, testVariables, '/world-info');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        const data = captured.response.data.data;
        expect(data).toHaveProperty('world');
        expect(data.world).toHaveProperty('id');
        expect(data.world).toHaveProperty('title');
        expect(data).toHaveProperty('system');
        expect(data.system).toHaveProperty('id');
        expect(data).toHaveProperty('foundryVersion');
        expect(data).toHaveProperty('modules');
        expect(data.modules).toBeInstanceOf(Array);
        expect(data).toHaveProperty('users');
        expect(data.users).toBeInstanceOf(Array);
        expect(data).toHaveProperty('activeScene');
        console.log(`  World: ${data.world.title}, System: ${data.system.id}, Modules: ${data.modules.length}, Users: ${data.users.length}`);
      }, 15000);
    });
  });

  forEachVersion((version, getClientId) => {
    describe(`POST /move-token (v${version})`, () => {
      test('POST /move-token - move a token by actor UUID', async () => {
        setVariable('clientId', getClientId());

        // Use the primary actor UUID — the move-token endpoint accepts an Actor UUID
        // and finds the corresponding token on the current scene
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/move-token',
            host: ['{{baseUrl}}'],
            path: ['move-token'],
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
              uuid: actorUuid,
              x: 200,
              y: 200,
              animate: false
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/move-token');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toHaveProperty('x', 200);
        expect(captured.response.data.data).toHaveProperty('y', 200);
        console.log(`  ✓ Moved token for actor ${actorUuid} to (200, 200)`);
      }, 15000);
    });
  });

  forEachVersion((version, getClientId) => {
    describe(`GET /measure-distance (v${version})`, () => {
      test('GET /measure-distance - measure between coordinates', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/measure-distance',
            host: ['{{baseUrl}}'],
            path: ['measure-distance'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'originX', value: '0' },
              { key: 'originY', value: '0' },
              { key: 'targetX', value: '500' },
              { key: 'targetY', value: '500' }
            ]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };

        const captured = await captureExample(requestConfig, testVariables, '/measure-distance');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        const data = captured.response.data.data;
        expect(data).toHaveProperty('distance');
        expect(data).toHaveProperty('units');
        expect(typeof data.distance).toBe('number');
        console.log(`  Distance: ${data.distance} ${data.units}`);
      }, 15000);
    });
  });

});
