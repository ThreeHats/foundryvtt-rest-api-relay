/**
 * @file combat-subscribe.test.ts
 * @description Combat SSE / WebSocket Subscription Tests
 * @endpoints WS subscribe to "combat-events" channel
 *
 * Tests subscribing to combat event notifications via the WebSocket API.
 */

import { describe, test, expect, afterAll, afterEach } from '@jest/globals';
import { makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { WsTestClient } from '../../helpers/wsClient';
import { captureWsExample, saveWsExamples } from '../../helpers/captureWsExample';
import { forEachVersion } from '../../helpers/multiVersion';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

const baseUrl = testVariables.baseUrl.replace(/^http/, 'ws');
const wsUrl = `${baseUrl}/ws/api`;

function getApiKey(): string {
  return testVariables.apiKey;
}

describe('Combat Subscribe', () => {
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
      const outputPath = path.join(__dirname, '../../../docs/examples/combat-subscribe-examples.json');
      saveWsExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }
  });

  forEachVersion((version, getClientId) => {

    describe(`Combat event subscription (v${version})`, () => {
      test('WS subscribe to combat-events channel', async () => {
        const clientId = getClientId();
        if (!clientId) {
          console.log('  Skipping: no clientId available');
          return;
        }

        const client = createClient();
        const connected = await client.connect(wsUrl, getApiKey(), clientId);
        expect(connected).toBeDefined();
        expect(connected.type).toBe('connected');

        // Subscribe to combat-events channel
        const request = { channel: 'combat-events' };
        const response = await client.sendAndWait({ type: 'subscribe', ...request });

        capturedExamples.push(
          captureWsExample('subscribe-combat-events', 'Subscribe to combat events', request, response, wsUrl)
        );

        // If supported, expect subscribed confirmation
        if (response.type === 'subscribed') {
          expect(response.channel).toBe('combat-events');
          console.log('  Successfully subscribed to combat-events');
        } else {
          // If not supported, it should return an error
          expect(response.type).toBe('error');
          console.log(`  Combat events channel not supported: ${response.error}`);
        }
      }, 15000);

      test('WS unsubscribe from combat-events channel', async () => {
        const clientId = getClientId();
        if (!clientId) {
          console.log('  Skipping: no clientId available');
          return;
        }

        const client = createClient();
        const connected = await client.connect(wsUrl, getApiKey(), clientId);
        expect(connected).toBeDefined();

        // Subscribe first
        const subResponse = await client.sendAndWait({ type: 'subscribe', channel: 'combat-events' });

        if (subResponse.type !== 'subscribed') {
          console.log('  Skipping unsubscribe test: combat-events channel not available');
          return;
        }

        // Now unsubscribe
        const unsubResponse = await client.sendAndWait({ type: 'unsubscribe', channel: 'combat-events' });

        capturedExamples.push(
          captureWsExample('unsubscribe-combat-events', 'Unsubscribe from combat events', { channel: 'combat-events' }, unsubResponse, wsUrl)
        );

        expect(unsubResponse.type).toBe('unsubscribed');
        expect(unsubResponse.channel).toBe('combat-events');
      }, 15000);

      test('WS subscribe and receive a combat event', async () => {
        const clientId = getClientId();
        if (!clientId) {
          console.log('  Skipping: no clientId available');
          return;
        }

        const client = createClient();
        await client.connect(wsUrl, getApiKey(), clientId);

        // Subscribe to combat-events — must succeed
        const subResponse = await client.sendAndWait({ type: 'subscribe', channel: 'combat-events' });
        expect(subResponse.type).toBe('subscribed');
        expect(subResponse.channel).toBe('combat-events');

        // Set up listener BEFORE triggering the action
        const combatEventPromise = new Promise<any>((resolve, reject) => {
          const timeout = setTimeout(
            () => reject(new Error('Timed out waiting for combat-event after POST /start-encounter')),
            10000
          );
          client.on('message', (data: any) => {
            if (data.type === 'combat-event') {
              clearTimeout(timeout);
              resolve(data);
            }
          });
        });

        // Trigger a combat event by starting an encounter
        setVariable('clientId', clientId);
        const startResponse = await makeRequest(replaceVariables({
          url: {
            raw: '{{baseUrl}}/start-encounter',
            host: ['{{baseUrl}}'],
            path: ['start-encounter'],
            query: [{ key: 'clientId', value: '{{clientId}}' }],
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}' }],
          body: { mode: 'raw', raw: JSON.stringify({ startWithSelected: false }) },
        }, testVariables));

        expect(startResponse.status).toBe(200);

        // Wait for the combat-event — required
        const combatEvent = await combatEventPromise;
        expect(combatEvent).toBeDefined();
        expect(combatEvent.type).toBe('combat-event');

        console.log(`  Received combat event: ${combatEvent.eventType || combatEvent.hookName || 'unknown'}`);
        capturedExamples.push(
          captureWsExample('combat-event', 'Combat event received', {}, combatEvent, wsUrl)
        );

        // Clean up: end the encounter
        await makeRequest(replaceVariables({
          url: {
            raw: '{{baseUrl}}/end-encounter',
            host: ['{{baseUrl}}'],
            path: ['end-encounter'],
            query: [{ key: 'clientId', value: '{{clientId}}' }],
          },
          method: 'POST',
          header: [{ key: 'x-api-key', value: '{{apiKey}}' }],
          body: { mode: 'raw', raw: '{}' },
        }, testVariables)).catch(() => {});
      }, 20000);
    });

  });
});
