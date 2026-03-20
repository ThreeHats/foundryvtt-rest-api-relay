/**
 * @file canvas-endpoints.test.ts
 * @description Canvas Endpoint Tests
 * @endpoints GET /canvas/:documentType, POST /canvas/:documentType, PUT /canvas/:documentType, DELETE /canvas/:documentType
 *
 * Tests all canvas document types: tokens, tiles, drawings, lights, sounds, notes, templates, walls
 * Each type gets full CRUD (create, read, update, delete) with captured examples for docs.
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../helpers/globalVariables';
import { getEntityUuid } from '../helpers/testEntities';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

/**
 * Canvas document type definitions with create and update data.
 * Each entry provides the minimal data needed to create a document of that type.
 */
interface CanvasDocTypeDef {
  type: string;
  createData: Record<string, any>;
  updateData: Record<string, any>;
  /** If true, the create data needs an actorId injected at test time */
  needsActorId?: boolean;
}

const CANVAS_DOC_TYPES: CanvasDocTypeDef[] = [
  {
    type: 'tokens',
    createData: { x: 400, y: 400 },
    updateData: { x: 450, y: 450 },
    needsActorId: true,
  },
  {
    type: 'tiles',
    createData: { x: 0, y: 0, width: 200, height: 200, texture: { src: '' } },
    updateData: { width: 300, height: 300 },
  },
  {
    type: 'drawings',
    createData: { x: 100, y: 100, shape: { type: 'r', width: 100, height: 100 } },
    updateData: { x: 150, y: 150 },
  },
  {
    type: 'lights',
    createData: { x: 300, y: 300, config: { dim: 20, bright: 10 } },
    updateData: { config: { dim: 30, bright: 15 } },
  },
  {
    type: 'sounds',
    createData: { x: 200, y: 200, radius: 10, path: '' },
    updateData: { radius: 20 },
  },
  {
    type: 'notes',
    createData: { x: 250, y: 250, text: 'Test Note' },
    updateData: { text: 'Updated Test Note' },
  },
  {
    type: 'templates',
    createData: { x: 350, y: 350, t: 'circle', distance: 15 },
    updateData: { distance: 20 },
  },
  {
    type: 'walls',
    createData: { c: [100, 100, 300, 100] },
    updateData: { c: [100, 100, 400, 100] },
  },
];

describe('Canvas', () => {
  afterAll(() => {
    const outputPath = path.join(__dirname, '../../docs/examples/canvas-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {

    // ═══════════════════════════════════════════
    // Full CRUD for each canvas document type
    // ═══════════════════════════════════════════

    for (const docDef of CANVAS_DOC_TYPES) {
      describe(`/canvas/${docDef.type} (v${version})`, () => {
        const globalIdKey = `canvas_${docDef.type}_id`;

        test(`POST /canvas/${docDef.type} - Create`, async () => {
          setVariable('clientId', getClientId());

          let createData = { ...docDef.createData };

          // Inject actorId for tokens
          if (docDef.needsActorId) {
            const actorUuid = getEntityUuid(version, 'Actor', 'primary');
            expect(actorUuid).toBeTruthy();
            createData.actorId = actorUuid!.split('.').pop();
          }

          const requestConfig: ApiRequestConfig = {
            url: {
              raw: `{{baseUrl}}/canvas/${docDef.type}`,
              host: ['{{baseUrl}}'],
              path: ['canvas', docDef.type],
              query: [{ key: 'clientId', value: '{{clientId}}' }]
            },
            method: 'POST',
            header: [
              { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
              { key: 'Content-Type', value: 'application/json', type: 'text' }
            ],
            body: {
              mode: 'raw',
              raw: JSON.stringify({ data: createData }, null, 2)
            }
          };

          const captured = await captureExample(
            requestConfig, testVariables,
            `/canvas/${docDef.type} - Create`
          );
          capturedExamples.push(captured);

          expect(captured.response.status).toBe(200);
          expect(captured.response.data.data).toBeInstanceOf(Array);
          expect(captured.response.data.data.length).toBe(1);
          expect(captured.response.data.data[0]._id).toBeTruthy();
          expect(captured.response.data.documentType).toBe(docDef.type);

          // Store for subsequent tests
          setGlobalVariable(version, globalIdKey, captured.response.data.data[0]._id);
        }, 30000);

        test(`GET /canvas/${docDef.type} - Get all`, async () => {
          setVariable('clientId', getClientId());

          const requestConfig: ApiRequestConfig = {
            url: {
              raw: `{{baseUrl}}/canvas/${docDef.type}`,
              host: ['{{baseUrl}}'],
              path: ['canvas', docDef.type],
              query: [{ key: 'clientId', value: '{{clientId}}' }]
            },
            method: 'GET',
            header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
          };

          const captured = await captureExample(
            requestConfig, testVariables,
            `/canvas/${docDef.type} - Get all`
          );
          capturedExamples.push(captured);

          expect(captured.response.status).toBe(200);
          expect(captured.response.data.data).toBeInstanceOf(Array);
          expect(captured.response.data.documentType).toBe(docDef.type);
          expect(captured.response.data.data.length).toBeGreaterThan(0);
        });

        test(`GET /canvas/${docDef.type} - Get by ID`, async () => {
          setVariable('clientId', getClientId());

          const docId = getGlobalVariable(version, globalIdKey);
          expect(docId).toBeTruthy();

          const requestConfig: ApiRequestConfig = {
            url: {
              raw: `{{baseUrl}}/canvas/${docDef.type}`,
              host: ['{{baseUrl}}'],
              path: ['canvas', docDef.type],
              query: [
                { key: 'clientId', value: '{{clientId}}' },
                { key: 'documentId', value: docId }
              ]
            },
            method: 'GET',
            header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
          };

          const captured = await captureExample(
            requestConfig, testVariables,
            `/canvas/${docDef.type} - Get by ID`
          );
          capturedExamples.push(captured);

          expect(captured.response.status).toBe(200);
          expect(captured.response.data.data).toBeInstanceOf(Array);
          expect(captured.response.data.data.length).toBe(1);
          expect(captured.response.data.data[0]._id).toBe(docId);
        });

        test(`PUT /canvas/${docDef.type} - Update`, async () => {
          setVariable('clientId', getClientId());

          const docId = getGlobalVariable(version, globalIdKey);
          expect(docId).toBeTruthy();

          const requestConfig: ApiRequestConfig = {
            url: {
              raw: `{{baseUrl}}/canvas/${docDef.type}`,
              host: ['{{baseUrl}}'],
              path: ['canvas', docDef.type],
              query: [{ key: 'clientId', value: '{{clientId}}' }]
            },
            method: 'PUT',
            header: [
              { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
              { key: 'Content-Type', value: 'application/json', type: 'text' }
            ],
            body: {
              mode: 'raw',
              raw: JSON.stringify({
                documentId: docId,
                data: docDef.updateData
              }, null, 2)
            }
          };

          const captured = await captureExample(
            requestConfig, testVariables,
            `/canvas/${docDef.type} - Update`
          );
          capturedExamples.push(captured);

          expect(captured.response.status).toBe(200);
          expect(captured.response.data.data).toBeInstanceOf(Array);
          expect(captured.response.data.data.length).toBe(1);
          expect(captured.response.data.documentType).toBe(docDef.type);
        }, 30000);

        test(`DELETE /canvas/${docDef.type} - Delete`, async () => {
          setVariable('clientId', getClientId());

          const docId = getGlobalVariable(version, globalIdKey);
          expect(docId).toBeTruthy();

          const requestConfig: ApiRequestConfig = {
            url: {
              raw: `{{baseUrl}}/canvas/${docDef.type}`,
              host: ['{{baseUrl}}'],
              path: ['canvas', docDef.type],
              query: [
                { key: 'clientId', value: '{{clientId}}' },
                { key: 'documentId', value: docId }
              ]
            },
            method: 'DELETE',
            header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
          };

          const captured = await captureExample(
            requestConfig, testVariables,
            `/canvas/${docDef.type} - Delete`
          );
          capturedExamples.push(captured);

          expect(captured.response.status).toBe(200);
          expect(captured.response.data.success).toBe(true);
          expect(captured.response.data.documentType).toBe(docDef.type);
        });
      });
    }

    // ═══════════════════════════════════════════
    // Persistent token for utility/encounter tests
    // ═══════════════════════════════════════════

    describe(`/canvas/tokens - persistent token (v${version})`, () => {
      test('POST /canvas/tokens - Create persistent token for utility/encounter tests', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();
        const actorId = actorUuid!.split('.').pop();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/canvas/tokens',
            host: ['{{baseUrl}}'],
            path: ['canvas', 'tokens'],
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
              data: { actorId, x: 500, y: 500 }
            }, null, 2)
          }
        };

        const captured = await captureExample(
          requestConfig, testVariables,
          '/canvas/tokens - Create persistent token'
        );

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toBeInstanceOf(Array);
        expect(captured.response.data.data.length).toBe(1);

        const tokenId = captured.response.data.data[0]._id;
        setGlobalVariable(version, 'persistentTokenId', tokenId);
        console.log(`  ✓ Created persistent token ${tokenId} for utility/encounter tests`);
      }, 30000);
    });

    // ═══════════════════════════════════════════
    // Validation
    // ═══════════════════════════════════════════

    describe(`/canvas - validation (v${version})`, () => {
      test('GET /canvas/invalid - Invalid document type returns error', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/canvas/invalid',
            host: ['{{baseUrl}}'],
            path: ['canvas', 'invalid'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };

        const captured = await captureExample(
          requestConfig, testVariables,
          '/canvas/invalid - Invalid document type'
        );

        expect(captured.response.status).toBe(400);
        expect(captured.response.data.error).toContain('Invalid documentType');
      });
    });

  });
});
