/**
 * @file client-endpoints.test.ts
 * @description Client listing endpoint tests
 * @endpoints GET /clients
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import { forEachVersion } from '../../helpers/multiVersion';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Clients', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../../docs/examples/clients-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/clients (v${version})`, () => {
      test('GET /clients - list connected clients', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/clients',
            host: ['{{baseUrl}}'],
            path: ['clients'],
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
          '/clients'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('total');
        expect(captured.response.data).toHaveProperty('clients');
        expect(Array.isArray(captured.response.data.clients)).toBe(true);
        expect(captured.response.data.total).toBeGreaterThanOrEqual(1);

        // Verify client structure
        const client = captured.response.data.clients[0];
        expect(client).toHaveProperty('clientId');
        expect(client).toHaveProperty('instanceId');
        expect(client).toHaveProperty('lastSeen');
        expect(client).toHaveProperty('connectedSince');
        expect(client).toHaveProperty('worldId');
        expect(client).toHaveProperty('worldTitle');
        expect(client).toHaveProperty('foundryVersion');
        expect(client).toHaveProperty('systemId');
        expect(client).toHaveProperty('systemTitle');
        expect(client).toHaveProperty('systemVersion');
      });

      test('GET /clients - unauthorized without API key', async () => {
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/clients',
            host: ['{{baseUrl}}'],
            path: ['clients'],
          },
          method: 'GET',
          header: []
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/clients (unauthorized)'
        );

        expect(captured.response.status).toBe(401);
      });
    });
  });
});
