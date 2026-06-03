/**
 * @file known-clients.test.ts
 * @description Integration tests for the get-known-clients WebSocket action,
 * including verification that publicUrl is returned in the client list.
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import WebSocket from 'ws';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';

const baseUrl = testVariables.baseUrl;
const wsRelayUrl = baseUrl.replace(/^http/, 'ws') + '/relay';

/**
 * Connects to the /relay endpoint using a connection token (Foundry module auth flow).
 * Returns a connected WebSocket that has completed auth-success.
 */
async function connectToRelay(
  clientId: string,
  token: string,
  worldId = 'test-known-clients',
  worldTitle = 'Test Known Clients'
): Promise<WebSocket> {
  const url = `${wsRelayUrl}?id=${encodeURIComponent(clientId)}&worldId=${encodeURIComponent(worldId)}&worldTitle=${encodeURIComponent(worldTitle)}`;

  return new Promise((resolve, reject) => {
    const ws = new WebSocket(url);
    const timeout = setTimeout(() => {
      ws.close();
      reject(new Error('Relay connection timed out'));
    }, 10000);

    ws.on('open', () => {
      ws.send(JSON.stringify({ type: 'auth', token }));
    });

    ws.on('message', (raw: Buffer | string) => {
      const data = JSON.parse(typeof raw === 'string' ? raw : raw.toString('utf8'));
      if (data.type === 'auth-success') {
        clearTimeout(timeout);
        resolve(ws);
      } else if (data.type === 'error') {
        clearTimeout(timeout);
        ws.close();
        reject(new Error(`Relay auth error: ${data.message}`));
      }
    });

    ws.on('error', (err) => {
      clearTimeout(timeout);
      reject(err);
    });

    ws.on('close', (code) => {
      clearTimeout(timeout);
      if (code !== 1000) {
        reject(new Error(`Relay connection closed with code ${code}`));
      }
    });
  });
}

/**
 * Sends a message on an open WebSocket and waits for a response matching the requestId.
 */
async function sendAndWait(ws: WebSocket, message: object): Promise<any> {
  const requestId = `test-${Date.now()}-${Math.random().toString(36).slice(2)}`;
  const payload = { ...message, requestId };

  return new Promise((resolve, reject) => {
    const timeout = setTimeout(() => reject(new Error('sendAndWait timed out')), 10000);

    const handler = (raw: Buffer | string) => {
      const data = JSON.parse(typeof raw === 'string' ? raw : raw.toString('utf8'));
      if (data.requestId === requestId) {
        clearTimeout(timeout);
        ws.off('message', handler);
        resolve(data);
      }
    };

    ws.on('message', handler);
    ws.send(JSON.stringify(payload));
  });
}

describe('get-known-clients', () => {
  let connectionToken = '';
  let pairedClientId = '';

  afterAll(async () => {
    if (!pairedClientId) return;
    // Clean up the KnownClient row created during testing
    try {
      const sessionToken = testVariables.sessionToken;
      if (!sessionToken) return;
      const res = await axios.get(`${baseUrl}/auth/known-clients`, {
        headers: { Authorization: `Bearer ${sessionToken}` },
      });
      const testClient = res.data?.clients?.find((c: any) => c.clientId === pairedClientId);
      if (testClient?.id) {
        await axios.delete(`${baseUrl}/auth/known-clients/${testClient.id}`, {
          headers: { Authorization: `Bearer ${sessionToken}` },
        });
      }
    } catch {
      // cleanup is best-effort
    }
  });

  test('pairing flow produces a connection token', async () => {
    const sessionToken = testVariables.sessionToken;
    if (!sessionToken) {
      console.log('  Skipping: no sessionToken available (need authenticated session)');
      return;
    }

    // Step 1: Generate pairing code
    const codeRes = await axios.post(
      `${baseUrl}/auth/connection-tokens`,
      {},
      { headers: { Authorization: `Bearer ${sessionToken}` } }
    );
    expect(codeRes.status).toBe(201);
    expect(codeRes.data).toHaveProperty('code');
    const pairingCode = codeRes.data.code;

    // Step 2: Complete pairing to get connection token + clientId
    const pairRes = await axios.post(`${baseUrl}/auth/pair`, {
      code: pairingCode,
      worldId: 'test-known-clients',
      worldTitle: 'Test Known Clients',
    });
    expect(pairRes.status).toBe(200);
    expect(pairRes.data).toHaveProperty('token');
    expect(pairRes.data).toHaveProperty('clientId');

    connectionToken = pairRes.data.token;
    pairedClientId = pairRes.data.clientId;
  }, 15000);

  test('get-known-clients requires connection-token auth', async () => {
    // Connect with a plain API key (not a connection token) — should be rejected
    const clientId = pairedClientId || 'fvtt_fallbacktest12345';
    if (!pairedClientId) {
      console.log('  Skipping: no pairedClientId (pairing test did not run)');
      return;
    }

    // Use a second client connected with API key to the /ws/api endpoint
    // and verify it cannot call get-known-clients (requires connection-token auth)
    const wsApiUrl = baseUrl.replace(/^http/, 'ws') + '/ws/api';
    const ws = new WebSocket(`${wsApiUrl}?clientId=${encodeURIComponent(clientId)}`);

    const result = await new Promise<any>((resolve, reject) => {
      const timeout = setTimeout(() => { ws.close(); reject(new Error('timeout')); }, 10000);

      ws.on('open', () => ws.send(JSON.stringify({ type: 'auth', token: testVariables.apiKey })));
      ws.on('message', (raw) => {
        const data = JSON.parse(typeof raw === 'string' ? raw : raw.toString('utf8'));
        if (data.type === 'connected') {
          // Now send get-known-clients — should fail because /ws/api uses API key, not connection token
          const requestId = 'test-no-conn-token';
          ws.send(JSON.stringify({ type: 'get-known-clients', requestId }));
        }
        if (data.type === 'known-clients-result') {
          clearTimeout(timeout);
          ws.close();
          resolve(data);
        }
      });
      ws.on('error', (e) => { clearTimeout(timeout); reject(e); });
      // Any close is a valid rejection of the API-key connection — resolve on
      // every code (not just 1000) so a policy-violation close doesn't time out.
      ws.on('close', (code) => { clearTimeout(timeout); resolve({ type: 'closed', code }); });
    });

    // The API-key connection must NOT receive a valid known-clients list.
    // Valid rejections: the connection closes, errors, or returns a result that
    // carries an error. The ONLY failure is a successful clients array leaking
    // over an x-api-key connection — assert that explicitly so the test can't
    // pass vacuously when the connection is simply closed.
    const leakedClientList =
      result.type === 'known-clients-result' && Array.isArray(result.clients) && !result.error;
    expect(leakedClientList).toBe(false);
  }, 15000);

  test('get-known-clients returns clients with publicUrl field', async () => {
    if (!connectionToken || !pairedClientId) {
      console.log('  Skipping: pairing did not complete');
      return;
    }

    const ws = await connectToRelay(pairedClientId, connectionToken);

    try {
      const resp = await sendAndWait(ws, { type: 'get-known-clients' });

      expect(resp.type).toBe('known-clients-result');
      expect(Array.isArray(resp.clients)).toBe(true);

      // Every client entry must have a publicUrl string field
      for (const client of resp.clients) {
        expect(client).toHaveProperty('publicUrl');
        expect(typeof client.publicUrl).toBe('string');
      }
    } finally {
      ws.close();
    }
  }, 15000);
});
