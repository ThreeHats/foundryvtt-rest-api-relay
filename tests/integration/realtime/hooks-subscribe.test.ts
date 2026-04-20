/**
 * @file hooks-subscribe.test.ts
 * @description Hooks SSE / WebSocket Subscription Tests
 * @endpoints GET /hooks/subscribe, WS subscribe to "hooks" channel
 *
 * Tests subscribing to Foundry hook events via the WebSocket API.
 */

import { describe, test, expect, afterAll, afterEach } from '@jest/globals';
import { makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { forEachVersion } from '../../helpers/multiVersion';
import { WsTestClient } from '../../helpers/wsClient';
import { captureWsExample, saveWsExamples } from '../../helpers/captureWsExample';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

const baseUrl = testVariables.baseUrl.replace(/^http/, 'ws');
const wsUrl = `${baseUrl}/ws/api`;

function getApiKey(): string {
  return testVariables.apiKey;
}

describe('Hooks Subscribe', () => {
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

  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/hooks-examples.json');
      saveWsExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }
  });

  forEachVersion((version, getClientId) => {
    let createdActorUuid: string | null = null;

    describe(`Hooks subscription (v${version})`, () => {
      afterEach(async () => {
        if (createdActorUuid) {
          const uuid = createdActorUuid;
          createdActorUuid = null;
          const clientId = getClientId();
          if (clientId) {
            setVariable('clientId', clientId);
            await makeRequest(replaceVariables({
              url: {
                raw: `{{baseUrl}}/delete?clientId={{clientId}}&uuid=${encodeURIComponent(uuid)}`,
                host: ['{{baseUrl}}'],
                path: ['delete'],
              },
              method: 'DELETE',
              header: [{ key: 'x-api-key', value: '{{apiKey}}' }],
            }, testVariables)).catch(() => {});
          }
        }
      });

      test('WS subscribe to hooks channel', async () => {
        const clientId = getClientId();
        if (!clientId) {
          console.log('  Skipping: no clientId available');
          return;
        }

        const client = createClient();
        const connected = await client.connect(wsUrl, getApiKey(), clientId);
        expect(connected).toBeDefined();
        expect(connected.type).toBe('connected');

        // Subscribe to hooks channel
        const request = { type: 'subscribe', channel: 'hooks' };
        const response = await client.sendAndWait(request);

        // The server may support or reject 'hooks' as a channel
        // Capture the result either way for documentation
        capturedExamples.push(
          captureWsExample('subscribe-hooks', 'Subscribe to hooks channel', { channel: 'hooks' }, response, wsUrl)
        );

        // If supported, expect subscribed confirmation
        if (response.type === 'subscribed') {
          expect(response.channel).toBe('hooks');
        } else {
          // If not supported, it should return an error with a message
          expect(response.type).toBe('error');
          console.log(`  Hooks channel not supported: ${response.error}`);
        }
      }, 15000);

      test('WS subscribe and receive a hook event', async () => {
        const clientId = getClientId();
        if (!clientId) {
          console.log('  Skipping: no clientId available');
          return;
        }

        const client = createClient();
        const connected = await client.connect(wsUrl, getApiKey(), clientId);
        expect(connected).toBeDefined();

        // Subscribe to hooks — must succeed
        const subResponse = await client.sendAndWait({ type: 'subscribe', channel: 'hooks' });
        expect(subResponse.type).toBe('subscribed');
        expect(subResponse.channel).toBe('hooks');

        // Brief pause: the relay sends event-subscription-update to the Foundry module
        // asynchronously. Give Foundry time to register its hook listeners before we
        // fire the trigger.
        await new Promise(r => setTimeout(r, 300));

        // Set up a listener for hook events BEFORE triggering the action
        const hookEventPromise = new Promise<any>((resolve, reject) => {
          const timeout = setTimeout(
            () => reject(new Error('Timed out waiting for hook-event after POST /create')),
            10000
          );
          client.on('message', (data: any) => {
            if (data.type === 'hook-event') {
              clearTimeout(timeout);
              resolve(data);
            }
          });
        });

        // Trigger a hook by creating an actor — fires createActor, which IS in FORWARDED_HOOKS.
        // (createChatMessage is intentionally excluded from FORWARDED_HOOKS.)
        setVariable('clientId', clientId);
        const createResponse = await makeRequest(replaceVariables({
          url: {
            raw: '{{baseUrl}}/create',
            host: ['{{baseUrl}}'],
            path: ['create'],
            query: [{ key: 'clientId', value: '{{clientId}}' }],
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}' }],
          body: { mode: 'raw', raw: JSON.stringify({ entityType: 'Actor', data: { name: 'hook-test-actor', type: 'base' } }) },
        }, testVariables));

        expect(createResponse.status).toBe(200);
        createdActorUuid = createResponse.data?.uuid ?? null;

        // Wait for the hook-event — required, not optional.
        // The fanout sends { type: "hook-event", hook: "<name>", data: {...} }
        const hookEvent = await hookEventPromise;
        expect(hookEvent).toBeDefined();
        expect(hookEvent.type).toBe('hook-event');
        expect(hookEvent.hook).toBeDefined();

        console.log(`  Received hook event: ${hookEvent.hook}`);
        capturedExamples.push(
          captureWsExample('hook-event', 'Hook event received', {}, hookEvent, wsUrl)
        );
      }, 15000);
    });

  });
});
