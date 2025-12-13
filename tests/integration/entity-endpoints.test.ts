/**
 * @file entity-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Entity CRUD Endpoint Tests
 * @endpoints POST /create, GET /get, PUT /update, DELETE /delete, POST /give, POST /remove, POST /increase, POST /decrease, POST /kill
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { createTestEntities, getEntityUuid } from '../helpers/testEntities';
import { getGlobalVariable, setGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Entity', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/entity-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/create (v${version})`, () => {
      test('POST /create', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        console.log(`\nCreating test entities for v${version}...`);
        
        // Create actors for various tests
        await createTestEntities(version, [
          { key: 'primary', entityType: 'Actor', captureForDocs: true },
          { key: 'secondary', entityType: 'Actor' },
          { key: 'expendable', entityType: 'Actor' },  // For delete test
        ], { capturedExamples });
        
        // Create items for tests
        await createTestEntities(version, [
          { key: 'primary', entityType: 'Item' },
        ], { capturedExamples });
        
        // Create a journal entry
        await createTestEntities(version, [
          { key: 'primary', entityType: 'JournalEntry' },
        ], { capturedExamples });

        // Create macro
        await createTestEntities(version, [
          { key: 'primary', entityType: 'Macro' },
        ], { capturedExamples });
        
        // Verify primary actor was created
        const primaryActorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(primaryActorUuid).toBeTruthy();
      }, 30000);
    });

    describe(`/get (v${version})`, () => {
      test('GET /get', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Use the primary actor we created
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/get',
            host: ['{{baseUrl}}'],
            path: ['get'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'uuid',
                value: actorUuid!,
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
          '/get'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.data).toBeTruthy();
        expect(captured.response.data.data.name).toBe('test-actor');
        expect(captured.response.data.data.type).toBe('base');
      });
    });

    describe(`/update (v${version})`, () => {
      test('PUT /update', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Use the primary actor we created
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/update',
            host: ['{{baseUrl}}'],
            path: ['update'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'uuid',
                value: actorUuid!,
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
              data: {
                name: 'Updated Test Actor'
              }
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/update'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.entity).toBeTruthy();
        expect(captured.response.data.entity).toBeInstanceOf(Array);
        expect(captured.response.data.entity[0].name).toBe('Updated Test Actor');
      });
    });

    describe(`/delete (v${version})`, () => {
      test('DELETE /delete', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Use the expendable actor (created specifically for deletion)
        const actorUuid = getEntityUuid(version, 'Actor', 'expendable');
        expect(actorUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/delete',
            host: ['{{baseUrl}}'],
            path: ['delete'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'uuid',
                
                value: actorUuid!,
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

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/delete'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data.uuid).toBe(actorUuid);
      });
    });

    describe(`/give (v${version})`, () => {
      test('POST /give', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Use receiver actor and primary item
        const receiverUuid = getEntityUuid(version, 'Actor', 'primary');
        const itemUuid = getEntityUuid(version, 'Item', 'primary');
        
        expect(receiverUuid).toBeTruthy();
        expect(itemUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/give',
            host: ['{{baseUrl}}'],
            path: ['give'],
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
              toUuid: receiverUuid,
              itemUuid: itemUuid,
              quantity: 1
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/give'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data.itemUuid).toBe(itemUuid);
        expect(captured.response.data.toUuid).toBe(receiverUuid);
        expect(captured.response.data.quantity).toBe(1);
        expect(captured.response.data.newItemId).toBeTruthy();

        // Set global variables for other tests
        setGlobalVariable(version, 'expendableActorSubItem', `${receiverUuid}.Item.${captured.response.data.newItemId}`);
      });
    });

    describe(`/remove (v${version})`, () => {
      test('POST /remove', async () => {
        setVariable('clientId', getClientId());
        
        // Use primary actor and expendable item
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        const itemUuid = getGlobalVariable(version, 'expendableActorSubItem')!;
        
        expect(actorUuid).toBeTruthy();
        expect(itemUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/remove',
            host: ['{{baseUrl}}'],
            path: ['remove'],
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
              actorUuid: actorUuid,
              itemUuid: itemUuid,
              quantity: 1
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/remove'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data.itemUuid).toBe(itemUuid);
      });
    });

    describe(`/increase (v${version})`, () => {
      test('POST /increase', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/increase',
            host: ['{{baseUrl}}'],
            path: ['increase'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'uuid',
                value: actorUuid!,
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
              attribute: 'prototypeToken.height',
              amount: 5
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/increase'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data).toHaveProperty('results');
        expect(captured.response.data.results[0]).toHaveProperty('uuid', actorUuid);
        expect(captured.response.data.results[0]).toHaveProperty('attribute', 'prototypeToken.height');
        expect(captured.response.data.results[0]).toHaveProperty('oldValue', 1);
        expect(captured.response.data.results[0]).toHaveProperty('newValue', 6);
      });
    });

    describe(`/decrease (v${version})`, () => {
      test('POST /decrease', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/decrease',
            host: ['{{baseUrl}}'],
            path: ['decrease'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}'
              },
              {
                key: 'uuid',
                value: actorUuid!,
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
              attribute: 'prototypeToken.height',
              amount: 5
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/decrease'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data.success).toBe(true);
        expect(captured.response.data).toHaveProperty('results');
        expect(captured.response.data.results[0]).toHaveProperty('uuid', actorUuid);
        expect(captured.response.data.results[0]).toHaveProperty('attribute', 'prototypeToken.height');
        expect(captured.response.data.results[0]).toHaveProperty('oldValue', 6);
        expect(captured.response.data.results[0]).toHaveProperty('newValue', 1);
      });
    });

    describe(`/kill (v${version})`, () => {
      test('POST /kill', async () => {
        setVariable('clientId', getClientId());
        
        // Use secondary actor for kill test (don't kill primary)
        const actorUuid = getEntityUuid(version, 'Actor', 'secondary');
        expect(actorUuid).toBeTruthy();
        
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/kill',
            host: ['{{baseUrl}}'],
            path: ['kill'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'uuid',
                value: actorUuid!,
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
          '/kill'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();
        expect(captured.response.data).toHaveProperty('results');
        expect(captured.response.data.results[0]).toHaveProperty('uuid', actorUuid);
        expect(captured.response.data.results[0]).toHaveProperty('success', true);
      });
    });

  });

});
