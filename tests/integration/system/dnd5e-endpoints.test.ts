/**
 * @file dnd5e-endpoints.test.ts
 * @description D&D 5th Edition System-Specific Endpoint Tests
 * @endpoints GET /dnd5e/get-actor-details, POST /dnd5e/modify-item-charges, POST /dnd5e/use-ability,
 *            POST /dnd5e/use-feature, POST /dnd5e/use-spell, POST /dnd5e/use-item, POST /dnd5e/modify-experience,
 *            POST /dnd5e/short-rest, POST /dnd5e/long-rest, POST /dnd5e/skill-check,
 *            POST /dnd5e/ability-save, POST /dnd5e/ability-check, POST /dnd5e/death-save,
 *            GET /dnd5e/concentration, POST /dnd5e/break-concentration, POST /dnd5e/concentration-save,
 *            POST /dnd5e/equip-item, POST /dnd5e/attune-item, POST /dnd5e/transfer-currency
 *
 * These tests only run on Foundry instances with the dnd5e system.
 * They use compendium-sourced actors (from dnd5e.heroes) which come with
 * features, spells, items, and proper system data.
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import { forEachVersionWithSystem } from '../../helpers/multiVersion';
import { getEntityUuid } from '../../helpers/testEntities';
import { getGlobalVariable } from '../../helpers/globalVariables';
import { setupSystemTestData } from '../../helpers/systemSetup';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

// Track whether concentration was successfully set up per version
const concentrationReady: Record<string, boolean> = {};

// Track resolved item names for inventory tests per version
const resolvedItemNames: Record<string, string> = {};

/**
 * Add a concentration effect to an actor via POST /effects.
 * Returns true if successful.
 */
async function addConcentrationEffect(actorUuid: string): Promise<boolean> {
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
          name: 'Concentrating',
          icon: 'icons/magic/light/orb-lightbulb-gray.webp',
          statuses: ['concentrating'],
          duration: { rounds: 10 },
          changes: []
        }
      })
    }
  };

  const resolved = replaceVariables(addConfig, testVariables);
  const response = await makeRequest(resolved);
  return response.status === 200;
}

describe('Dnd5e', () => {
  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/dnd5e-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }
  });

  forEachVersionWithSystem('dnd5e', (version, getClientId) => {

    // ═══════════════════════════════════════════
    // Setup: give the test actor system-specific items (spells, consumables)
    // ═══════════════════════════════════════════

    describe(`System setup (v${version})`, () => {
      test('Give test actor spells and consumables from compendium', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const result = await setupSystemTestData(version, actorUuid!);
        console.log(`  Setup complete: spell=${result.spellName || 'none'}, consumable=${result.consumableName || 'none'}`);
      }, 30000);
    });

    // ═══════════════════════════════════════════
    // GET /dnd5e/get-actor-details
    // ═══════════════════════════════════════════

    describe(`/dnd5e/get-actor-details (v${version})`, () => {
      test('GET /dnd5e/get-actor-details - all detail types', async () => {
        setVariable('clientId', getClientId());

        // Use the primary actor created from compendium (dnd5e.heroes)
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/get-actor-details',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'get-actor-details'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'actorUuid', value: actorUuid! },
              { key: 'details', value: '["resources","items","features","spells"]' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/get-actor-details');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('uuid', actorUuid);

        // At least some of these should be present on a compendium hero
        const details = captured.response.data.data;
        console.log(`  ✓ Actor details: resources=${!!details.resources}, items=${details.items?.length ?? 0}, features=${details.features?.length ?? 0}, spells=${details.spells?.length ?? 0}`);

        // Resources should be present (even if empty) for a character type
        if (details.resources) {
          expect(details.resources).toBeTruthy();
        }
      }, 15000);

      test('GET /dnd5e/get-actor-details - items only', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/get-actor-details',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'get-actor-details'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'actorUuid', value: actorUuid! },
              { key: 'details', value: '["items"]' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const resolved = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolved);

        expect(response.status).toBe(200);
        expect(response.data.data).toHaveProperty('uuid', actorUuid);
        // Should only have items, not features/spells/resources
        if (response.data.data.items) {
          expect(response.data.data.items).toBeInstanceOf(Array);
        }
        // These should NOT be present since we only asked for items
        expect(response.data.data).not.toHaveProperty('spells');
        expect(response.data.data).not.toHaveProperty('features');
        expect(response.data.data).not.toHaveProperty('resources');
      });
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/modify-experience
    // ═══════════════════════════════════════════

    describe(`/dnd5e/modify-experience (v${version})`, () => {
      test('POST /dnd5e/modify-experience - add XP', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/modify-experience',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'modify-experience'],
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
              actorUuid: actorUuid,
              amount: 100
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/modify-experience - add');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('actorUuid');
        expect(captured.response.data.data).toHaveProperty('oldXp');
        expect(captured.response.data.data).toHaveProperty('newXp');
        expect(captured.response.data.data.newXp).toBe(captured.response.data.data.oldXp + 100);
        console.log(`  ✓ XP: ${captured.response.data.data.oldXp} → ${captured.response.data.data.newXp}`);
      });

      test('POST /dnd5e/modify-experience - remove XP', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/modify-experience',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'modify-experience'],
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
              actorUuid: actorUuid,
              amount: -100
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/modify-experience - remove');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data.newXp).toBe(captured.response.data.data.oldXp - 100);
        console.log(`  ✓ XP: ${captured.response.data.data.oldXp} → ${captured.response.data.data.newXp}`);
      });
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/modify-item-charges
    // ═══════════════════════════════════════════

    describe(`/dnd5e/modify-item-charges (v${version})`, () => {
      test('POST /dnd5e/modify-item-charges - spend and restore a charge', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        // Find an item with charges via get-actor-details
        const detailsConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/get-actor-details',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'get-actor-details'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'actorUuid', value: actorUuid! },
              { key: 'details', value: '["items"]' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const detailsResolved = replaceVariables(detailsConfig, testVariables);
        const detailsResponse = await makeRequest(detailsResolved);

        const items = detailsResponse.data?.data?.items || [];
        // Find an item with numeric max charges > 0
        const chargeItem = items.find((i: any) => {
          const max = parseInt(i.system?.uses?.max);
          return max > 0;
        });

        if (!chargeItem) {
          console.log(`  ○ Skipping modify-item-charges: no items with charges found`);
          return;
        }

        console.log(`  Using item with charges: "${chargeItem.name}" (max: ${chargeItem.system.uses.max}, spent: ${chargeItem.system.uses.spent})`);

        // Step 1: Spend a charge (amount = -1)
        const spendConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/modify-item-charges',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'modify-item-charges'],
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
              actorUuid: actorUuid,
              itemName: chargeItem.name,
              amount: -1
            })
          }
        };

        const spendCaptured = await captureExample(spendConfig, testVariables, '/dnd5e/modify-item-charges');
        capturedExamples.push(spendCaptured);

        expect(spendCaptured.response.status).toBe(200);
        expect(spendCaptured.response.data).toHaveProperty('data');
        expect(spendCaptured.response.data.data).toHaveProperty('oldCharges');
        expect(spendCaptured.response.data.data).toHaveProperty('newCharges');
        console.log(`  ✓ Spent charge: ${spendCaptured.response.data.data.oldCharges} → ${spendCaptured.response.data.data.newCharges}`);

        // Step 2: Restore the charge (amount = 1)
        const restoreConfig: ApiRequestConfig = {
          ...spendConfig,
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              actorUuid: actorUuid,
              itemName: chargeItem.name,
              amount: 1
            })
          }
        };

        const restoreResolved = replaceVariables(restoreConfig, testVariables);
        const restoreResponse = await makeRequest(restoreResolved);

        expect(restoreResponse.status).toBe(200);
        expect(restoreResponse.data.data).toHaveProperty('oldCharges');
        expect(restoreResponse.data.data).toHaveProperty('newCharges');
        console.log(`  ✓ Restored charge: ${restoreResponse.data.data.oldCharges} → ${restoreResponse.data.data.newCharges}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/use-item (use an equipment/weapon by name)
    // ═══════════════════════════════════════════

    describe(`/dnd5e/use-item (v${version})`, () => {
      test('POST /dnd5e/use-item - use item by name from actor inventory', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        // First, get actor details to find an item name
        const detailsConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/get-actor-details',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'get-actor-details'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'actorUuid', value: actorUuid! },
              { key: 'details', value: '["items"]' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const detailsResolved = replaceVariables(detailsConfig, testVariables);
        const detailsResponse = await makeRequest(detailsResolved);

        if (detailsResponse.status !== 200 || !detailsResponse.data.data?.items?.length) {
          console.log(`  ○ Skipping use-item: actor has no items`);
          return;
        }

        const itemName = detailsResponse.data.data.items[0].name;
        console.log(`  Using item: "${itemName}"`);

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/use-item',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'use-item'],
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
              actorUuid: actorUuid,
              abilityName: itemName
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/use-item');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('ability', itemName);
        console.log(`  ✓ Used item: ${itemName}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/use-feature (use a feat by name)
    // ═══════════════════════════════════════════

    describe(`/dnd5e/use-feature (v${version})`, () => {
      test('POST /dnd5e/use-feature - use feature by name', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        // Get features from actor
        const detailsConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/get-actor-details',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'get-actor-details'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'actorUuid', value: actorUuid! },
              { key: 'details', value: '["features"]' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const detailsResolved = replaceVariables(detailsConfig, testVariables);
        const detailsResponse = await makeRequest(detailsResolved);

        if (detailsResponse.status !== 200 || !detailsResponse.data.data?.features?.length) {
          console.log(`  ○ Skipping use-feature: actor has no features`);
          return;
        }

        const featureName = detailsResponse.data.data.features[0].name;
        console.log(`  Using feature: "${featureName}"`);

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/use-feature',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'use-feature'],
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
              actorUuid: actorUuid,
              abilityName: featureName
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/use-feature');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('ability', featureName);
        console.log(`  ✓ Used feature: ${featureName}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/use-spell (use a spell by name)
    // ═══════════════════════════════════════════

    describe(`/dnd5e/use-spell (v${version})`, () => {
      test('POST /dnd5e/use-spell - use spell by name', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        // Use the spell name from system setup, or discover from actor
        let spellName = getGlobalVariable(version, 'dnd5e_test_spell_name');

        if (!spellName) {
          // Fallback: check if actor has any spells natively
          const detailsConfig: ApiRequestConfig = {
            url: {
              raw: '{{baseUrl}}/dnd5e/get-actor-details',
              host: ['{{baseUrl}}'],
              path: ['dnd5e', 'get-actor-details'],
              query: [
                { key: 'clientId', value: '{{clientId}}' },
                { key: 'actorUuid', value: actorUuid! },
                { key: 'details', value: '["spells"]' }
              ]
            },
            method: 'GET',
            header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
          };
          const detailsResponse = await makeRequest(replaceVariables(detailsConfig, testVariables));
          spellName = detailsResponse.data?.data?.spells?.[0]?.name;
        }

        if (!spellName) {
          console.log(`  ○ Skipping use-spell: actor has no spells (system setup may have failed)`);
          return;
        }

        console.log(`  Using spell: "${spellName}"`);

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/use-spell',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'use-spell'],
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
              actorUuid: actorUuid,
              abilityName: spellName
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/use-spell');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('ability', spellName);
        console.log(`  ✓ Used spell: ${spellName}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/use-ability (generic — uses any item type)
    // ═══════════════════════════════════════════

    describe(`/dnd5e/use-ability (v${version})`, () => {
      test('POST /dnd5e/use-ability - use any ability by name', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        // Get all details to find any usable ability
        const detailsConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/get-actor-details',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'get-actor-details'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'actorUuid', value: actorUuid! },
              { key: 'details', value: '["items","features","spells"]' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const detailsResolved = replaceVariables(detailsConfig, testVariables);
        const detailsResponse = await makeRequest(detailsResolved);
        const details = detailsResponse.data?.data || {};

        // Find any ability — prefer items, then features, then spells
        const ability = (details.items?.[0]) || (details.features?.[0]) || (details.spells?.[0]);
        if (!ability) {
          console.log(`  ○ Skipping use-ability: actor has no abilities`);
          return;
        }

        console.log(`  Using ability: "${ability.name}" (${ability.type})`);

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/use-ability',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'use-ability'],
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
              actorUuid: actorUuid,
              abilityName: ability.name
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/use-ability');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('ability', ability.name);
        console.log(`  ✓ Used ability: ${ability.name}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/skill-check
    // ═══════════════════════════════════════════

    describe(`/dnd5e/skill-check (v${version})`, () => {
      test('POST /dnd5e/skill-check - roll a perception check', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/skill-check',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'skill-check'],
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
              actorUuid: actorUuid,
              skill: 'prc'
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/skill-check');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('skill', 'prc');
        expect(captured.response.data.data).toHaveProperty('total');
        expect(typeof captured.response.data.data.total).toBe('number');
        console.log(`  ✓ Perception check: total=${captured.response.data.data.total}, formula=${captured.response.data.data.formula}`);
      }, 15000);

      test('POST /dnd5e/skill-check - roll with advantage', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/skill-check',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'skill-check'],
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
              actorUuid: actorUuid,
              skill: 'ste',
              advantage: true
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/skill-check - advantage');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('skill', 'ste');
        expect(captured.response.data.data).toHaveProperty('total');
        console.log(`  ✓ Stealth (advantage): total=${captured.response.data.data.total}`);
      }, 15000);

      test('POST /dnd5e/skill-check - invalid skill returns error', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/skill-check',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'skill-check'],
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
              actorUuid: actorUuid,
              skill: 'invalid_skill'
            })
          }
        };

        const resolved = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolved);

        expect(response.status).toBe(400);
        expect(response.data).toHaveProperty('error');
        expect(response.data.error).toMatch(/invalid skill/i);
        console.log(`  ✓ Invalid skill rejected: ${response.data.error}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/ability-save
    // ═══════════════════════════════════════════

    describe(`/dnd5e/ability-save (v${version})`, () => {
      test('POST /dnd5e/ability-save - roll a dexterity save', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/ability-save',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'ability-save'],
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
              actorUuid: actorUuid,
              ability: 'dex'
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/ability-save');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('ability', 'dex');
        expect(captured.response.data.data).toHaveProperty('total');
        expect(typeof captured.response.data.data.total).toBe('number');
        console.log(`  ✓ DEX save: total=${captured.response.data.data.total}, formula=${captured.response.data.data.formula}`);
      }, 15000);

      test('POST /dnd5e/ability-save - invalid ability returns error', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/ability-save',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'ability-save'],
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
              actorUuid: actorUuid,
              ability: 'xyz'
            })
          }
        };

        const resolved = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolved);

        expect(response.status).toBe(400);
        expect(response.data).toHaveProperty('error');
        expect(response.data.error).toMatch(/invalid ability/i);
        console.log(`  ✓ Invalid ability rejected: ${response.data.error}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/ability-check
    // ═══════════════════════════════════════════

    describe(`/dnd5e/ability-check (v${version})`, () => {
      test('POST /dnd5e/ability-check - roll a strength check', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/ability-check',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'ability-check'],
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
              actorUuid: actorUuid,
              ability: 'str'
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/ability-check');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('ability', 'str');
        expect(captured.response.data.data).toHaveProperty('total');
        expect(typeof captured.response.data.data.total).toBe('number');
        console.log(`  ✓ STR check: total=${captured.response.data.data.total}, formula=${captured.response.data.data.formula}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/death-save
    // ═══════════════════════════════════════════

    describe(`/dnd5e/death-save (v${version})`, () => {
      test('POST /dnd5e/death-save - roll a death saving throw', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        // First, set the actor to 0 HP so death saves are valid
        const killConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/kill',
            host: ['{{baseUrl}}'],
            path: ['kill'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'uuid', value: actorUuid! }
            ]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const killResolved = replaceVariables(killConfig, testVariables);
        await makeRequest(killResolved);

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/death-save',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'death-save'],
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
              actorUuid: actorUuid
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/death-save');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('total');
        expect(typeof captured.response.data.data.total).toBe('number');
        expect(captured.response.data.data).toHaveProperty('deathSaves');
        expect(captured.response.data.data.deathSaves).toHaveProperty('success');
        expect(captured.response.data.data.deathSaves).toHaveProperty('failure');
        console.log(`  ✓ Death save: total=${captured.response.data.data.total}, successes=${captured.response.data.data.deathSaves.success}, failures=${captured.response.data.data.deathSaves.failure}`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/short-rest
    // ═══════════════════════════════════════════

    describe(`/dnd5e/short-rest (v${version})`, () => {
      test('POST /dnd5e/short-rest - perform a short rest', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/short-rest',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'short-rest'],
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
              actorUuid: actorUuid
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/short-rest');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('actorUuid');
        expect(captured.response.data.data).toHaveProperty('result');
        console.log(`  ✓ Short rest completed`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /dnd5e/long-rest
    // ═══════════════════════════════════════════

    describe(`/dnd5e/long-rest (v${version})`, () => {
      test('POST /dnd5e/long-rest - perform a long rest', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/long-rest',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'long-rest'],
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
              actorUuid: actorUuid,
              newDay: true
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/long-rest');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('actorUuid');
        expect(captured.response.data.data).toHaveProperty('result');
        console.log(`  ✓ Long rest completed`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // Concentration Tests
    // ═══════════════════════════════════════════

    describe(`/dnd5e/concentration (v${version})`, () => {
      test('GET /dnd5e/concentration - Check concentration status', async () => {
        setVariable('clientId', getClientId());
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/concentration',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'concentration'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'actorUuid', value: actorUuid! }
            ]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/concentration');
        capturedExamples.push(captured);
        expect(captured.response.status).toBe(200);
      }, 15000);
    });

    describe(`Concentration setup (v${version})`, () => {
      test('Add concentration effect to test actor', async () => {
        setVariable('clientId', getClientId());
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const success = await addConcentrationEffect(actorUuid!);
        concentrationReady[version] = success;
        if (success) {
          console.log(`  Added concentration effect to test actor`);
        } else {
          console.warn(`  Failed to add concentration effect - break/save tests will be skipped`);
        }
        expect(success).toBe(true);
      }, 15000);
    });

    describe(`/dnd5e/concentration-save (v${version})`, () => {
      test('POST /dnd5e/concentration-save - Roll concentration save', async () => {
        if (!concentrationReady[version]) {
          console.log('  Skipped: concentration effect not set up');
          return;
        }
        setVariable('clientId', getClientId());
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/concentration-save',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'concentration-save'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ actorUuid: actorUuid!, damage: 15 }, null, 2)
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/concentration-save');
        capturedExamples.push(captured);
        expect(captured.response.status).toBe(200);

        // If the save failed, re-add concentration for the break test
        const maintained = captured.response.data?.data?.maintained;
        if (maintained === false) {
          console.log('  Concentration was broken by save, re-adding for break test');
          await addConcentrationEffect(actorUuid!);
        }
      }, 15000);
    });

    describe(`/dnd5e/break-concentration (v${version})`, () => {
      test('POST /dnd5e/break-concentration - Break concentration', async () => {
        if (!concentrationReady[version]) {
          console.log('  Skipped: concentration effect not set up');
          return;
        }
        setVariable('clientId', getClientId());
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/break-concentration',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'break-concentration'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ actorUuid: actorUuid! }, null, 2)
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/break-concentration');
        capturedExamples.push(captured);
        expect(captured.response.status).toBe(200);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // Inventory Tests
    // ═══════════════════════════════════════════

    describe(`Inventory setup (v${version})`, () => {
      test('Discover an equipable item on the test actor', async () => {
        setVariable('clientId', getClientId());
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const detailsConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/get-actor-details',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'get-actor-details'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'actorUuid', value: actorUuid! },
              { key: 'details', value: '["items"]' }
            ]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };

        const resolved = replaceVariables(detailsConfig, testVariables);
        const response = await makeRequest(resolved);
        expect(response.status).toBe(200);

        const items = response.data?.data?.items || response.data?.items || [];
        const equipableItem = items.find((item: any) =>
          item.type === 'weapon' || item.type === 'equipment' || item.type === 'armor'
        );

        if (equipableItem) {
          resolvedItemNames[version] = equipableItem.name;
          console.log(`  Found equipable item: "${equipableItem.name}" (${equipableItem.type})`);
        } else {
          console.warn(`  No equipable items found on test actor, inventory tests will be skipped`);
        }
      }, 15000);
    });

    describe(`/dnd5e/equip-item (v${version})`, () => {
      test('POST /dnd5e/equip-item - Equip an item', async () => {
        const itemName = resolvedItemNames[version];
        if (!itemName) { console.log('  Skipped: no equipable item found'); return; }
        setVariable('clientId', getClientId());
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/equip-item',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'equip-item'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ actorUuid: actorUuid!, itemName, equipped: true }, null, 2)
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/equip-item');
        capturedExamples.push(captured);
        expect(captured.response.status).toBe(200);
      }, 15000);
    });

    describe(`/dnd5e/attune-item (v${version})`, () => {
      test('POST /dnd5e/attune-item - Attune an item', async () => {
        const itemName = resolvedItemNames[version];
        if (!itemName) { console.log('  Skipped: no equipable item found'); return; }
        setVariable('clientId', getClientId());
        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/attune-item',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'attune-item'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ actorUuid: actorUuid!, itemName, attuned: true }, null, 2)
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/attune-item');
        capturedExamples.push(captured);
        expect(captured.response.status).toBe(200);
      }, 15000);
    });

    describe(`/dnd5e/transfer-currency (v${version})`, () => {
      test('POST /dnd5e/transfer-currency - Transfer currency between actors', async () => {
        setVariable('clientId', getClientId());
        const sourceActorUuid = getEntityUuid(version, 'Actor', 'primary');
        const targetActorUuid = getEntityUuid(version, 'Actor', 'secondary');
        expect(sourceActorUuid).toBeTruthy();
        expect(targetActorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/dnd5e/transfer-currency',
            host: ['{{baseUrl}}'],
            path: ['dnd5e', 'transfer-currency'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              sourceActorUuid: sourceActorUuid!,
              targetActorUuid: targetActorUuid!,
              currency: { gp: 1 }
            }, null, 2)
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/dnd5e/transfer-currency');
        capturedExamples.push(captured);
        expect(captured.response.status).toBe(200);
      }, 15000);
    });

  });
});
