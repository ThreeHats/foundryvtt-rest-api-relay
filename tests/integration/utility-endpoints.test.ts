/**
 * @file utility-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Utility and Canvas Interaction Endpoint Tests
 * @endpoints POST /select, GET /selected, POST /execute-js
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Utility', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/utility-examples.json');
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

});
