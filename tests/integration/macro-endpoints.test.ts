/**
 * @file macro-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Macro Endpoint Tests
 * @endpoints GET /macros, POST /macro/:uuid/execute
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { getEntityUuid } from '../helpers/testEntities';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Macro', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/macro-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/macros (v${version})`, () => {
      test('GET /macros', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/macros',
            host: ['{{baseUrl}}'],
            path: ['macros'],
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
          '/macros'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.macros).toBeTruthy();
        expect(captured.response.data.macros.length).toBeGreaterThan(0);
        expect(captured.response.data.macros[0].name).toBeTruthy();
        expect(captured.response.data.macros[0].uuid).toBeTruthy();
        expect(captured.response.data.macros[0].command).toBeTruthy();
      });
    });

    describe(`/macro (v${version})`, () => {
      test('POST /macro/:uuid/execute', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const macroUuid = getEntityUuid(version, 'Macro', 'primary');
        if (!macroUuid) {
          throw new Error(`Macro UUID not found for version ${version}`);
        }
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: `{{baseUrl}}/macro/${macroUuid}/execute`,
            host: ['{{baseUrl}}'],
            path: ['macro', macroUuid, 'execute'],
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
                args: {
                  targetName: "Goblin",
                  damage: 100000,
                  effect: "poison"
                }
              }, null, 2)
          }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/macro/:uuid/execute'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.result).toBeTruthy();
        expect(captured.response.data.result.success).toBe(true);
        expect(captured.response.data.result.damageDealt).toBe(100000);
        expect(captured.response.data.result.target).toBe("Goblin");
      });
    });

  });

});
