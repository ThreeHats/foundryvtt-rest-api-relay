/**
 * @file roll-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Dice Rolling Endpoint Tests
 * @endpoints POST /roll, GET /rolls, GET /lastroll
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Roll', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/roll-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/roll (v${version})`, () => {
      test('POST /roll', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/roll',
            host: ['{{baseUrl}}'],
            path: ['roll'],
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
              formula: '2d20kh',
              flavor: 'Test Roll',
              createChatMessage: true
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/roll'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('roll');
        expect(captured.response.data.data.roll).toHaveProperty('formula');
        expect(captured.response.data.data.roll).toHaveProperty('total');
        expect(captured.response.data.data.roll).toHaveProperty('isCritical');
        expect(captured.response.data.data.roll).toHaveProperty('isFumble');
        expect(captured.response.data.data.roll).toHaveProperty('dice');

        // Save variables for use in subsequent tests
        setGlobalVariable(version, 'lastRollTotal', captured.response.data.data.roll.total);
      });
    });

    describe(`/rolls (v${version})`, () => {
      test('GET /rolls', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/rolls',
            host: ['{{baseUrl}}'],
            path: ['rolls'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'limit',
                value: '20',
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
          '/rolls'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(Array.isArray(captured.response.data.data)).toBe(true);
        expect(captured.response.data.data.length).toBeLessThanOrEqual(20);
        expect(captured.response.data.data.length).toBeGreaterThan(0);
        expect(captured.response.data.data[0]).toHaveProperty('formula');
        expect(captured.response.data.data[0]).toHaveProperty('rollTotal', getGlobalVariable(version, 'lastRollTotal'));
        expect(captured.response.data.data[0]).toHaveProperty('isCritical');
        expect(captured.response.data.data[0]).toHaveProperty('isFumble');
        expect(captured.response.data.data[0]).toHaveProperty('dice');
      });
    });

    describe(`/lastroll (v${version})`, () => {
      test('GET /lastroll', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/lastroll',
            host: ['{{baseUrl}}'],
            path: ['lastroll'],
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

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/lastroll'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('formula');
        expect(captured.response.data.data).toHaveProperty('rollTotal', getGlobalVariable(version, 'lastRollTotal'));
        expect(captured.response.data.data).toHaveProperty('isCritical');
        expect(captured.response.data.data).toHaveProperty('isFumble');
        expect(captured.response.data.data).toHaveProperty('dice');
      });
    });

  });

});
