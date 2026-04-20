/**
 * @file websocket-api.test.ts
 * @description Client-facing WebSocket API integration tests
 *
 * Tests the /ws/api endpoint including authentication, request/response,
 * event subscriptions, and error handling.
 */

import { describe, test, expect, afterAll, afterEach } from '@jest/globals';
import axios from 'axios';
import { WsTestClient } from '../../helpers/wsClient';
import { testVariables } from '../../helpers/testVariables';
import { captureWsExample, saveWsExamples } from '../../helpers/captureWsExample';
import { forEachVersion } from '../../helpers/multiVersion';
import * as path from 'path';

const capturedExamples: any[] = [];
const baseUrl = testVariables.baseUrl.replace(/^http/, 'ws');
const wsUrl = `${baseUrl}/ws/api`;

describe('WebSocket API (/ws/api)', () => {
  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/ws-core-examples.json');
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

    // ═══════════════════════════════════════════
    // Authentication Tests
    // ═══════════════════════════════════════════

    describe(`Authentication (v${version})`, () => {
      test('should reject connection without token', async () => {
        const client = createClient();
        await expect(
          client.connect(`${baseUrl}/ws/api`, '', getClientId())
        ).rejects.toThrow();
      }, 3000);

      test('should reject connection with invalid token', async () => {
        const client = createClient();
        await expect(
          client.connect(`${baseUrl}/ws/api`, 'invalid-key-12345', getClientId())
        ).rejects.toThrow();
      }, 3000);

      test('should reject connection without clientId', async () => {
        const client = createClient();
        await expect(
          client.connect(`${baseUrl}/ws/api`, testVariables.apiKey, '')
        ).rejects.toThrow();
      }, 3000);

      test('should connect successfully with valid credentials', async () => {
        const clientId = getClientId();
        if (!clientId) {
          console.log('  Skipping: no clientId available');
          return;
        }

        const client = createClient();
        const connected = await client.connect(wsUrl, testVariables.apiKey, clientId);

        expect(connected).toBeDefined();
        expect(connected.type).toBe('connected');
        expect(connected.clientId).toBe(clientId);
        expect(connected.supportedTypes).toBeDefined();
        expect(Array.isArray(connected.supportedTypes)).toBe(true);
        expect(connected.eventChannels).toContain('chat-events');
        expect(connected.eventChannels).toContain('roll-events');
      }, 3000);
    });

    // ═══════════════════════════════════════════
    // Request/Response Tests
    // ═══════════════════════════════════════════

    describe(`Request/Response (v${version})`, () => {
      test('should return error for unknown message type', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const response = await client.sendAndWait({
          type: 'nonexistent-type',
          requestId: 'test-unknown',
        });

        expect(response.type).toBe('error');
        expect(response.error).toContain('Unknown message type');
      }, 3000);

      test('should return error for missing requestId', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const errorPromise = new Promise<any>((resolve) => {
          client.on('message', (data: any) => {
            if (data.type === 'error' && data.error?.includes('requestId')) {
              resolve(data);
            }
          });
        });

        client.send({ type: 'search' });
        const error = await errorPromise;
        expect(error.type).toBe('error');
      }, 3000);

      test('should return error for invalid JSON', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const errorPromise = new Promise<any>((resolve) => {
          client.on('message', (data: any) => {
            if (data.type === 'error' && data.error?.includes('Invalid JSON')) {
              resolve(data);
            }
          });
        });

        // Access the raw WS to send invalid JSON
        (client as any).ws.send('not valid json{{{');
        const error = await errorPromise;
        expect(error.error).toBe('Invalid JSON');
      }, 3000);

      test('should handle search request', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const request = { query: 'test' };
        const response = await client.sendAndWait({ type: 'search', ...request });

        expect(response.type).toBe('search-result');
        expect(response.requestId).toBeDefined();
        expect(response.clientId).toBe(clientId);

        capturedExamples.push(
          captureWsExample('search', 'Search for entities', request, response, wsUrl)
        );
      }, 30000);

      test('should handle ping/pong', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const response = await client.sendAndWait({ type: 'ping' });
        expect(response.type).toBe('pong');
      }, 3000);
    });

    // ═══════════════════════════════════════════
    // Event Subscription Tests
    // ═══════════════════════════════════════════

    describe(`Event Subscriptions (v${version})`, () => {
      test('should subscribe to chat-events', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const response = await client.subscribe('chat-events');
        expect(response.type).toBe('subscribed');
        expect(response.channel).toBe('chat-events');
      }, 3000);

      test('should subscribe to roll-events', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const response = await client.subscribe('roll-events');
        expect(response.type).toBe('subscribed');
        expect(response.channel).toBe('roll-events');
      }, 3000);

      test('should reject subscribe to invalid channel', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const response = await client.sendAndWait({
          type: 'subscribe',
          channel: 'invalid-channel',
        });

        expect(response.type).toBe('error');
        expect(response.error).toContain('Invalid channel');
      }, 3000);

      test('should unsubscribe from chat-events', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        await client.subscribe('chat-events');
        const response = await client.unsubscribe('chat-events');
        expect(response.type).toBe('unsubscribed');
        expect(response.channel).toBe('chat-events');
      }, 3000);

      test('should subscribe with filters', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);

        const response = await client.subscribe('chat-events', {
          speaker: 'GM',
          whisperOnly: false,
        });

        expect(response.type).toBe('subscribed');
        expect(response.channel).toBe('chat-events');
      }, 3000);

      test('should receive chat-event after subscribing', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);
        await client.subscribe('chat-events');

        // Listen BEFORE sending the message
        const eventPromise = new Promise<any>((resolve, reject) => {
          const timeout = setTimeout(
            () => reject(new Error('Timed out waiting for chat-event')),
            10000
          );
          client.onEvent('chat-event', (data) => {
            clearTimeout(timeout);
            resolve(data);
          });
        });

        // Send a chat message to trigger the event
        const chatResp = await axios.post(
          `${testVariables.baseUrl}/chat?clientId=${clientId}`,
          { content: 'WS chat-event test' },
          { headers: { 'x-api-key': testVariables.apiKey }, validateStatus: () => true }
        );
        expect(chatResp.status).toBe(200);

        const event = await eventPromise;
        expect(event.type).toBe('chat-event');
        expect(event.data).toBeDefined();

        capturedExamples.push(
          captureWsExample('chat-event', 'Received chat-event via WS subscription', {}, event, wsUrl)
        );

        // Clean up
        const msgId = chatResp.data?.data?.id;
        if (msgId) {
          await axios.delete(`${testVariables.baseUrl}/chat/${msgId}?clientId=${clientId}`, {
            headers: { 'x-api-key': testVariables.apiKey },
            validateStatus: () => true,
          }).catch(() => {});
        }
      }, 15000);

      test('should receive roll-event after subscribing', async () => {
        const clientId = getClientId();
        if (!clientId) return;

        const client = createClient();
        await client.connect(wsUrl, testVariables.apiKey, clientId);
        await client.subscribe('roll-events');

        // Listen BEFORE triggering the roll
        const eventPromise = new Promise<any>((resolve, reject) => {
          const timeout = setTimeout(
            () => reject(new Error('Timed out waiting for roll-event')),
            10000
          );
          client.onEvent('roll-event', (data) => {
            clearTimeout(timeout);
            resolve(data);
          });
        });

        // Roll a die to trigger the event
        const rollResp = await axios.post(
          `${testVariables.baseUrl}/roll?clientId=${clientId}`,
          { formula: '1d6', flavor: 'WS roll-event test', createChatMessage: true },
          { headers: { 'x-api-key': testVariables.apiKey }, validateStatus: () => true }
        );
        expect(rollResp.status).toBe(200);

        const event = await eventPromise;
        expect(event.type).toBe('roll-event');
        expect(event.data).toBeDefined();

        capturedExamples.push(
          captureWsExample('roll-event', 'Received roll-event via WS subscription', {}, event, wsUrl)
        );
      }, 15000);
    });
  });
});
