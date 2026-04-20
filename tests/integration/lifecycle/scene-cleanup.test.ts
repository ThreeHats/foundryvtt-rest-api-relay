/**
 * @file scene-cleanup.test.ts
 * @description Restores the original active scene and deletes the test scene.
 * Must run after canvas-endpoints.test.ts (which depends on the test scene being active).
 */

import { describe, test, expect } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { forEachVersion } from '../../helpers/multiVersion';
import { getGlobalVariable } from '../../helpers/globalVariables';

describe('Scene Cleanup', () => {
  forEachVersion((version, getClientId) => {
    describe(`Restore original scene (v${version})`, () => {
      test('POST /switch-scene - Restore original active scene', async () => {
        setVariable('clientId', getClientId());

        const originalSceneId = getGlobalVariable(version, 'originalActiveSceneId');
        if (!originalSceneId) {
          console.log(`  ⚠ No original active scene to restore, skipping`);
          return;
        }

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/switch-scene',
            host: ['{{baseUrl}}'],
            path: ['switch-scene'],
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
              sceneId: originalSceneId
            }, null, 2)
          }
        };

        const resolvedConfig = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolvedConfig);

        expect(response.status).toBe(200);
        expect(response.data.success).toBe(true);
        console.log(`  ✓ Restored original active scene: ${originalSceneId}`);

        // Wait for the scene switch to fully complete before deleting
        await new Promise(resolve => setTimeout(resolve, 5000));
      }, 30000);

      test('DELETE /scene - Delete test scene', async () => {
        setVariable('clientId', getClientId());

        const sceneId = getGlobalVariable(version, 'testSceneId');
        if (!sceneId) {
          console.log(`  ⚠ No test scene to delete, skipping`);
          return;
        }

        // Verify we're not trying to delete the active scene
        const originalSceneId = getGlobalVariable(version, 'originalActiveSceneId');
        if (sceneId === originalSceneId) {
          console.log(`  ⚠ Test scene is the original scene, skipping delete`);
          return;
        }

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'sceneId',
                value: sceneId,
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

        const resolvedConfig = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolvedConfig);

        if (response.status !== 200) {
          console.log(`  ⚠ Failed to delete test scene (${response.status}):`, response.data);
        }
        expect(response.status).toBe(200);
        expect(response.data.success).toBe(true);
        console.log(`  ✓ Deleted test scene: ${sceneId}`);
      }, 30000);
    });
  });
});
