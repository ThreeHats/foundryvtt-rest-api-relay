/**
 * @file remote-request.test.ts
 * @description Tests for the relay-brokered cross-world tunnel (remote-request / remote-response).
 *
 * remote-request lets a Foundry module invoke any supported action on another world via the relay
 * tunnel. Requires the source connection token to have allowedTargetClients + remoteScopes
 * configured. API key clients always receive a "no connection token" rejection.
 *
 * Full end-to-end tests (actual successful remote actions) require two Foundry module instances
 * connected with properly configured connection tokens and are skipped when unavailable.
 */

import { describe, test, expect, afterAll, afterEach } from '@jest/globals';
import { testVariables } from '../../helpers/testVariables';
import { WsTestClient } from '../../helpers/wsClient';
import { captureWsExample, saveWsExamples } from '../../helpers/captureWsExample';
import { forEachVersion } from '../../helpers/multiVersion';
import * as path from 'path';

const capturedExamples: any[] = [];
const baseUrl = testVariables.baseUrl.replace('http', 'ws');
const wsUrl = `${baseUrl}/ws/api`;

describe('Remote-Request Cross-World Tunnel', () => {
  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/remote-request-examples.json');
      saveWsExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} WS examples to ${outputPath}`);
    }
  });

  forEachVersion((version, getClientId) => {
    const clients: WsTestClient[] = [];

    function createClient(): WsTestClient {
      const client = new WsTestClient();
      clients.push(client);
      return client;
    }

    afterEach(() => {
      for (const client of clients) {
        client.close();
      }
      clients.length = 0;
    });

    describe(`remote-request validation (v${version})`, () => {
      test('remote-request from API key client is rejected (no connection token)', async () => {
        const clientId = getClientId();
        if (!clientId) {
          console.log('  Skipping: no clientId available');
          return;
        }

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const response = await client.sendAndWait({
          type: 'remote-request',
          requestId: 'test-rejection-1',
          targetClientId: 'fvtt_target',
          action: 'entity',
          payload: {},
        }, 10000);

        // API key clients cannot use remote-request — must use a module connection token
        expect(response.type).toBe('remote-response');
        expect(response.success).toBe(false);
        expect(response.error).toBeDefined();

        capturedExamples.push(
          captureWsExample(
            'remote-request',
            'Remote action tunnel — rejected (no module connection token)',
            {
              targetClientId: 'fvtt_target',
              action: 'entity',
              payload: { entityType: 'Actor', id: 'Actor.abc123' },
              autoStartIfOffline: false,
            },
            response,
            wsUrl
          )
        );
      }, 15000);

      test('remote-request missing required fields is rejected', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        // Missing targetClientId and action
        const response = await client.sendAndWait({
          type: 'remote-request',
          requestId: 'test-missing-fields-1',
          payload: {},
        }, 10000);

        expect(response.type).toBe('remote-response');
        expect(response.success).toBe(false);
        expect(response.error).toBeDefined();
      }, 15000);

      test('remote-request missing action is rejected', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const response = await client.sendAndWait({
          type: 'remote-request',
          requestId: 'test-missing-action-1',
          targetClientId: 'fvtt_target',
          // action intentionally missing
          payload: {},
        }, 10000);

        expect(response.type).toBe('remote-response');
        expect(response.success).toBe(false);
        expect(response.error).toBeDefined();
      }, 15000);
    });
  });
});
