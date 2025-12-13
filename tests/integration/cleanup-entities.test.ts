/**
 * @file cleanup-entities.test.ts
 * @generated false
 * @description Cleanup Test Entities
 * @endpoints DELETE /delete (batch cleanup of created entities)
 */

import { describe, test, expect } from '@jest/globals';
import { makeRequest, replaceVariables, ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { forEachVersion } from '../helpers/multiVersion';
import { getCleanupList, clearCleanupList } from '../helpers/testEntities';

// Check if cleanup should be skipped
const skipCleanup = process.env.SKIP_CLEANUP === 'true';

describe('Cleanup Test Entities', () => {
  if (skipCleanup) {
    test.skip('Cleanup skipped (SKIP_CLEANUP=true)', () => {});
    return;
  }

  forEachVersion((version, getClientId) => {
    describe(`Cleanup (v${version})`, () => {
      test('DELETE all created test entities', async () => {
        const clientId = getClientId();
        if (!clientId) {
          console.log(`⚠️ No clientId for v${version}, skipping cleanup`);
          return;
        }
        
        setVariable('clientId', clientId);
        
        const uuids = getCleanupList(version);
        
        if (uuids.length === 0) {
          console.log(`ℹ️ No entities to clean up for v${version}`);
          return;
        }
        
        console.log(`\n🧹 Cleaning up ${uuids.length} entities for v${version}...`);
        
        let deleted = 0;
        let alreadyGone = 0;
        let failed = 0;
        
        // Delete in reverse order (in case of any dependencies)
        for (const uuid of [...uuids].reverse()) {
          const requestConfig: ApiRequestConfig = {
            url: {
              raw: '{{baseUrl}}/delete',
              host: ['{{baseUrl}}'],
              path: ['delete'],
              query: [
                { key: 'clientId', value: '{{clientId}}' },
                { key: 'uuid', value: uuid }
              ]
            },
            method: 'DELETE',
            header: [
              { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
            ]
          };
          
          try {
            const resolvedConfig = replaceVariables(requestConfig, testVariables);
            const response = await makeRequest(resolvedConfig);
            
            if (response.status === 200) {
              deleted++;
              console.log(`  ✓ Deleted ${uuid}`);
            } else if (response.status === 404) {
              alreadyGone++;
              console.log(`  ○ Already gone: ${uuid}`);
            } else {
              failed++;
              console.warn(`  ✗ Failed to delete ${uuid}: ${response.status}`);
            }
          } catch (error) {
            failed++;
            console.warn(`  ✗ Error deleting ${uuid}: ${error}`);
          }
        }
        
        // Clear the cleanup list
        clearCleanupList(version);
        
        console.log(`\n✅ Cleanup complete for v${version}:`);
        console.log(`   Deleted: ${deleted}`);
        console.log(`   Already gone: ${alreadyGone}`);
        if (failed > 0) {
          console.log(`   Failed: ${failed}`);
        }
        
        // Test passes as long as we attempted cleanup
        expect(true).toBe(true);
      }, 60000); // Extended timeout for cleanup
    });
  });
});
