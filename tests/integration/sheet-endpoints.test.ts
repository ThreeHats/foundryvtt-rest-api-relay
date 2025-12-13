/**
 * @file sheet-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Actor/Item Sheet Endpoint Tests
 * @endpoints GET /sheet
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { getEntityUuid } from '../helpers/testEntities';
import * as path from 'path';
import { get } from 'http';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Sheet', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/sheet-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/sheet (v${version})`, () => {
      test('GET /sheet', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/sheet',
            host: ['{{baseUrl}}'],
            path: ['sheet'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
                description: 'The ID of the Foundry client to connect to'
              },
              {
                key: 'uuid',
                value: getEntityUuid(version, 'Actor', 'primary') || '',
                description: 'The UUID of the entity to get the sheet for'
              },
              {
                key: 'format',
                value: 'json',
                description: 'The format to return the sheet in (html, json)'
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
          '/sheet'
        );
        capturedExamples.push(captured);

        // Assertions
        // Deprecated in v13 and greater
        expect([200, 400]).toContain(captured.response.status);
        if (captured.response.status === 200) {
          expect(captured.response.data).toHaveProperty('uuid');
          expect(captured.response.data).toHaveProperty('html');
          expect(captured.response.data).toHaveProperty('css');
        }
      });
    });

  });

});
