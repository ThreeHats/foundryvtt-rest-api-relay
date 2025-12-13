/**
 * @file search-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Entity Search Endpoint Tests
 * @endpoints GET /search
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { getEntityUuid } from '../helpers/testEntities';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Search', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/search-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/search (v${version})`, () => {
      test('GET /search', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/search',
            host: ['{{baseUrl}}'],
            path: ['search'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'query',
                value: 'test-item',
              },
              {
                key: 'filter',
                value: 'documentType:Item,subType:base',
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
          '/search'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('query');
        expect(captured.response.data).toHaveProperty('results');
        expect(captured.response.data.results.length).toBeGreaterThan(0);
        expect(captured.response.data.results[0]).toHaveProperty('uuid', getEntityUuid(version, 'Item', 'primary'));
        expect(captured.response.data.results[0]).toHaveProperty('documentType', 'Item');
      });
    });

  });

});
