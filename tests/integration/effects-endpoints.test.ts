/**
 * @file effects-endpoints.test.ts
 * @description Active Effects Endpoint Tests (system-agnostic)
 * @endpoints GET /effects, POST /effects, DELETE /effects
 *
 * These tests verify the ActiveEffect CRUD operations on actors.
 * They use the primary test actor created by entity-endpoints.test.ts.
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { getEntityUuid } from '../helpers/testEntities';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('Effects', () => {
  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../docs/examples/effects-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }
  });

  forEachVersion((version, getClientId) => {

    // ═══════════════════════════════════════════
    // GET /effects
    // ═══════════════════════════════════════════

    describe(`GET /effects (v${version})`, () => {
      test('GET /effects - list active effects on an actor', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/effects',
            host: ['{{baseUrl}}'],
            path: ['effects'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'uuid', value: actorUuid! }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const captured = await captureExample(requestConfig, testVariables, '/effects - list');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('uuid');
        expect(captured.response.data.data).toHaveProperty('effects');
        expect(captured.response.data.data.effects).toBeInstanceOf(Array);
        console.log(`  ✓ Actor has ${captured.response.data.data.effects.length} active effects`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /effects + DELETE /effects (add then remove)
    // ═══════════════════════════════════════════

    describe(`POST /effects + DELETE /effects (v${version})`, () => {
      test('POST /effects - add a custom effect, then DELETE to remove it', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        // Step 1: Add a custom effect
        const addConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/effects',
            host: ['{{baseUrl}}'],
            path: ['effects'],
            query: [
              { key: 'clientId', value: '{{clientId}}' }
            ]
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
              effectData: {
                name: 'Test Effect',
                icon: 'icons/svg/aura.svg',
                changes: []
              }
            })
          }
        };

        const addCaptured = await captureExample(addConfig, testVariables, '/effects - add custom');
        capturedExamples.push(addCaptured);

        expect(addCaptured.response.status).toBe(200);
        expect(addCaptured.response.data).toHaveProperty('data');
        expect(addCaptured.response.data.data).toHaveProperty('effect');
        expect(addCaptured.response.data.data.effect).toHaveProperty('id');
        expect(addCaptured.response.data.data.effect).toHaveProperty('name', 'Test Effect');

        const effectId = addCaptured.response.data.data.effect.id;
        console.log(`  ✓ Added custom effect: id=${effectId}, name=${addCaptured.response.data.data.effect.name}`);

        // Step 2: Verify the effect appears in GET /effects
        const getConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/effects',
            host: ['{{baseUrl}}'],
            path: ['effects'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'uuid', value: actorUuid! }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const getResolved = replaceVariables(getConfig, testVariables);
        const getResponse = await makeRequest(getResolved);

        expect(getResponse.status).toBe(200);
        const effects = getResponse.data.data.effects;
        const found = effects.find((e: any) => e.id === effectId);
        expect(found).toBeTruthy();
        console.log(`  ✓ Effect confirmed in GET /effects (${effects.length} total)`);

        // Step 3: Remove the effect
        const deleteConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/effects',
            host: ['{{baseUrl}}'],
            path: ['effects'],
            query: [
              { key: 'clientId', value: '{{clientId}}' }
            ]
          },
          method: 'DELETE',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              uuid: actorUuid,
              effectId: effectId
            })
          }
        };

        const deleteCaptured = await captureExample(deleteConfig, testVariables, '/effects - remove');
        capturedExamples.push(deleteCaptured);

        expect(deleteCaptured.response.status).toBe(200);
        expect(deleteCaptured.response.data).toHaveProperty('data');
        expect(deleteCaptured.response.data.data).toHaveProperty('removedEffectId', effectId);
        console.log(`  ✓ Removed effect: id=${effectId}`);

        // Step 4: Verify the effect is gone
        const verifyResolved = replaceVariables(getConfig, testVariables);
        const verifyResponse = await makeRequest(verifyResolved);

        expect(verifyResponse.status).toBe(200);
        const remainingEffects = verifyResponse.data.data.effects;
        const stillFound = remainingEffects.find((e: any) => e.id === effectId);
        expect(stillFound).toBeFalsy();
        console.log(`  ✓ Effect no longer present (${remainingEffects.length} remaining)`);
      }, 30000);
    });

  });
});
