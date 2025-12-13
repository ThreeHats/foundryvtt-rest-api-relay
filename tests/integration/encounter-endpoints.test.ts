/**
 * @file encounter-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Combat Encounter Management Endpoint Tests
 * @endpoints POST /start-encounter, GET /encounters, POST /next-turn, POST /next-round, POST /last-turn, POST /last-round, POST /remove-from-encounter, POST /add-to-encounter, DELETE /end-encounter
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';
import { get } from 'http';
import e from 'cors';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Encounter', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/encounter-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/start-encounter (v${version})`, () => {
      test('POST /start-encounter', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/start-encounter',
            host: ['{{baseUrl}}'],
            path: ['start-encounter'],
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
              startWithSelected: true,
              rollAll: true,
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/start-encounter'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounter');
        expect(captured.response.data.encounter).toHaveProperty('id');
        expect(captured.response.data.encounter).toHaveProperty('round');
        expect(captured.response.data.encounter).toHaveProperty('turn');
        expect(captured.response.data.encounter).toHaveProperty('combatants');
        expect(captured.response.data.encounter.combatants).toBeInstanceOf(Array);
        expect(captured.response.data.encounter.combatants.length).toBeGreaterThan(0);
        expect(captured.response.data.encounter.combatants[0]).toHaveProperty('id');
        expect(captured.response.data.encounter.combatants[0]).toHaveProperty('name');
        expect(captured.response.data.encounter.combatants[0]).toHaveProperty('tokenUuid', getGlobalVariable(version, 'selectedTokenUuid'));
        expect(captured.response.data.encounter.combatants[0]).toHaveProperty('actorUuid', getGlobalVariable(version, 'selectedActorUuid'));
        expect(captured.response.data.encounter.combatants[0]).toHaveProperty('img');
        expect(captured.response.data.encounter.combatants[0]).toHaveProperty('initiative');
        expect(captured.response.data.encounter.combatants[0]).toHaveProperty('hidden');
        expect(captured.response.data.encounter.combatants[0]).toHaveProperty('defeated');

        // Set global variable for future tests
        setGlobalVariable(version, 'encounterId', captured.response.data.encounter.id);
      });
    });

    describe(`/encounters (v${version})`, () => {
      test('GET /encounters', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/encounters',
            host: ['{{baseUrl}}'],
            path: ['encounters'],
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
          '/encounters'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounters');
        expect(captured.response.data.encounters[0]).toHaveProperty('id');
        expect(captured.response.data.encounters[0]).toHaveProperty('round');
        expect(captured.response.data.encounters[0]).toHaveProperty('turn');
        expect(captured.response.data.encounters[0]).toHaveProperty('combatants');
        expect(captured.response.data.encounters[0].combatants).toBeInstanceOf(Array);
        expect(captured.response.data.encounters[0].combatants.length).toBeGreaterThan(0);
        expect(captured.response.data.encounters[0].combatants[0]).toHaveProperty('id');
        expect(captured.response.data.encounters[0].combatants[0]).toHaveProperty('name');
        expect(captured.response.data.encounters[0].combatants[0]).toHaveProperty('tokenUuid', getGlobalVariable(version, 'selectedTokenUuid'));
        expect(captured.response.data.encounters[0].combatants[0]).toHaveProperty('actorUuid', getGlobalVariable(version, 'selectedActorUuid'));
        expect(captured.response.data.encounters[0].combatants[0]).toHaveProperty('img');
        expect(captured.response.data.encounters[0].combatants[0]).toHaveProperty('initiative');
        expect(captured.response.data.encounters[0].combatants[0]).toHaveProperty('hidden');
        expect(captured.response.data.encounters[0].combatants[0]).toHaveProperty('defeated');
      });
    });

    describe(`/next-turn (v${version})`, () => {
      test('POST /next-turn', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const encounterId = getGlobalVariable(version, 'encounterId');
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/next-turn',
            host: ['{{baseUrl}}'],
            path: ['next-turn'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'encounterId',
                value: encounterId
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
          ]
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/next-turn'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounterId', encounterId);
        expect(captured.response.data).toHaveProperty('action');
        expect(captured.response.data).toHaveProperty('currentTurn');
        expect(captured.response.data).toHaveProperty('currentRound');
        expect(captured.response.data).toHaveProperty('actorTurn');
        expect(captured.response.data).toHaveProperty('tokenTurn');
        expect(captured.response.data).toHaveProperty('encounter');
        expect(captured.response.data.encounter).toHaveProperty('id', encounterId);
        expect(captured.response.data.encounter).toHaveProperty('round');
        expect(captured.response.data.encounter).toHaveProperty('turn');
      });
    });

    describe(`/next-round (v${version})`, () => {
      test('POST /next-round', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const encounterId = getGlobalVariable(version, 'encounterId');
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/next-round',
            host: ['{{baseUrl}}'],
            path: ['next-round'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'encounterId',
                value: encounterId
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
          ]
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/next-round'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounterId', encounterId);
        expect(captured.response.data).toHaveProperty('action');
        expect(captured.response.data).toHaveProperty('currentTurn');
        expect(captured.response.data).toHaveProperty('currentRound');
        expect(captured.response.data).toHaveProperty('actorTurn');
        expect(captured.response.data).toHaveProperty('tokenTurn');
        expect(captured.response.data).toHaveProperty('encounter');
        expect(captured.response.data.encounter).toHaveProperty('id', encounterId);
        expect(captured.response.data.encounter).toHaveProperty('round');
        expect(captured.response.data.encounter).toHaveProperty('turn');
      });
    });

    describe(`/last-turn (v${version})`, () => {
      test('POST /last-turn', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const encounterId = getGlobalVariable(version, 'encounterId');
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/last-turn',
            host: ['{{baseUrl}}'],
            path: ['last-turn'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'encounterId',
                value: encounterId
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
          ]
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/last-turn'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounterId', encounterId);
        expect(captured.response.data).toHaveProperty('action');
        expect(captured.response.data).toHaveProperty('currentTurn');
        expect(captured.response.data).toHaveProperty('currentRound');
        expect(captured.response.data).toHaveProperty('actorTurn');
        expect(captured.response.data).toHaveProperty('tokenTurn');
        expect(captured.response.data).toHaveProperty('encounter');
        expect(captured.response.data.encounter).toHaveProperty('id', encounterId);
        expect(captured.response.data.encounter).toHaveProperty('round');
        expect(captured.response.data.encounter).toHaveProperty('turn');
      });
    });

    describe(`/last-round (v${version})`, () => {
      test('POST /last-round', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const encounterId = getGlobalVariable(version, 'encounterId');
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/last-round',
            host: ['{{baseUrl}}'],
            path: ['last-round'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'encounterId',
                value: encounterId
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
          ]
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/last-round'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounterId', encounterId);
        expect(captured.response.data).toHaveProperty('action');
        expect(captured.response.data).toHaveProperty('currentTurn');
        expect(captured.response.data).toHaveProperty('currentRound');
        expect(captured.response.data).toHaveProperty('actorTurn');
        expect(captured.response.data).toHaveProperty('tokenTurn');
        expect(captured.response.data).toHaveProperty('encounter');
        expect(captured.response.data.encounter).toHaveProperty('id', encounterId);
        expect(captured.response.data.encounter).toHaveProperty('round');
        expect(captured.response.data.encounter).toHaveProperty('turn');
      });
    });

    describe(`/remove-from-encounter (v${version})`, () => {
      test('POST /remove-from-encounter', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const encounter = getGlobalVariable(version, 'encounterId');
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/remove-from-encounter',
            host: ['{{baseUrl}}'],
            path: ['remove-from-encounter'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'encounterId',
                value: encounter
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
              selected: true
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/remove-from-encounter'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounterId', encounter);
        expect(captured.response.data).toHaveProperty('removed');
        expect(captured.response.data.removed).toBeInstanceOf(Array);
        expect(captured.response.data.removed.length).toBeGreaterThan(0);
        expect(captured.response.data.removed).toContain(getGlobalVariable(version, 'selectedTokenUuid'));
        expect(captured.response.data).toHaveProperty('failed');
        expect(captured.response.data.failed).toBeInstanceOf(Array);
        expect(captured.response.data.failed.length).toBe(0);
      });
    });

    describe(`/add-to-encounter (v${version})`, () => {
      test('POST /add-to-encounter', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const encounter = getGlobalVariable(version, 'encounterId');
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/add-to-encounter',
            host: ['{{baseUrl}}'],
            path: ['add-to-encounter'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'encounterId',
                value: encounter
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
              selected: true,
              uuids: [],
              rollInitiative: true
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/add-to-encounter'
        );
        capturedExamples.push(captured);
        
        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounterId', encounter);
        expect(captured.response.data).toHaveProperty('added');
        expect(captured.response.data.added).toBeInstanceOf(Array);
        expect(captured.response.data.added.length).toBeGreaterThan(0);
        expect(captured.response.data.added).toContain(getGlobalVariable(version, 'selectedTokenUuid'));
        expect(captured.response.data).toHaveProperty('failed')
        expect(captured.response.data.failed).toBeInstanceOf(Array);
        expect(captured.response.data.failed.length).toBe(0);
      });
    });

    describe(`/end-encounter (v${version})`, () => {
      test('POST /end-encounter', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const encounterId = getGlobalVariable(version, 'encounterId');
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/end-encounter',
            host: ['{{baseUrl}}'],
            path: ['end-encounter'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'encounterId',
                value: encounterId
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
          ]
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/end-encounter'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('encounterId', encounterId);
        expect(captured.response.data).toHaveProperty('message');
      });
    });

  });

});
